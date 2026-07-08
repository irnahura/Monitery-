package monitor

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"peekaping/backend/internal/analytics"
	"peekaping/backend/internal/dto"
	"peekaping/backend/internal/models"
	"peekaping/backend/internal/notification"
	"peekaping/backend/internal/repository"
)

type Service struct {
	store          *repository.Store
	notifier       notification.Notifier
	analytics      analytics.Service
	defaultTimeout time.Duration
}

func NewService(store *repository.Store, notifier notification.Notifier, analyticsService analytics.Service, defaultTimeout time.Duration) Service {
	return Service{store: store, notifier: notifier, analytics: analyticsService, defaultTimeout: defaultTimeout}
}

func (s Service) List(ctx context.Context, userID uint) ([]dto.MonitorSummary, error) {
	monitors, err := s.store.ListMonitorsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.MonitorSummary, 0, len(monitors))
	for _, monitor := range monitors {
		out = append(out, repository.MonitorToSummary(monitor))
	}
	return out, nil
}

func (s Service) Create(ctx context.Context, userID uint, input dto.MonitorRequest) (dto.MonitorSummary, error) {
	if err := validatePayload(input); err != nil {
		return dto.MonitorSummary{}, err
	}
	monitor := toModel(input, userID)
	if err := s.store.CreateMonitor(ctx, &monitor); err != nil {
		return dto.MonitorSummary{}, err
	}
	return repository.MonitorToSummary(monitor), nil
}

func (s Service) Update(ctx context.Context, userID, monitorID uint, input dto.MonitorRequest) (dto.MonitorSummary, error) {
	if err := validatePayload(input); err != nil {
		return dto.MonitorSummary{}, err
	}
	monitor, err := s.store.FindMonitorByIDAndUser(ctx, monitorID, userID)
	if err != nil {
		return dto.MonitorSummary{}, err
	}
	updated := toModel(input, userID)
	fields := map[string]any{
		"name":             updated.Name,
		"url":              updated.URL,
		"interval_seconds": updated.IntervalSeconds,
		"timeout_seconds":  updated.TimeoutSeconds,
		"method":           updated.Method,
		"headers":          updated.Headers,
		"follow_redirects": updated.FollowRedirects,
		"validate_ssl":     updated.ValidateSSL,
		"retry_count":      updated.RetryCount,
	}
	if err := s.store.UpdateMonitor(ctx, &monitor, fields); err != nil {
		return dto.MonitorSummary{}, err
	}
	monitor, err = s.store.FindMonitorByID(ctx, monitorID)
	if err != nil {
		return dto.MonitorSummary{}, err
	}
	return repository.MonitorToSummary(monitor), nil
}

func (s Service) Delete(ctx context.Context, userID, monitorID uint) error {
	deleted, err := s.store.DeleteMonitorWithChecks(ctx, monitorID, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("monitor not found")
	}
	return nil
}

func (s Service) History(ctx context.Context, userID, monitorID uint, limit int) (dto.MonitorHistoryResponse, error) {
	if _, err := s.store.FindMonitorByIDAndUser(ctx, monitorID, userID); err != nil {
		return dto.MonitorHistoryResponse{}, err
	}
	checks, err := s.store.ListHealthChecks(ctx, monitorID, limit)
	if err != nil {
		return dto.MonitorHistoryResponse{}, err
	}
	out := make([]dto.HealthCheckResponse, 0, len(checks))
	for _, check := range checks {
		out = append(out, repository.HealthCheckToResponse(check))
	}
	summary, err := s.analytics.Summary(ctx, monitorID)
	if err != nil {
		return dto.MonitorHistoryResponse{}, err
	}
	return dto.MonitorHistoryResponse{History: out, Summary: summary}, nil
}

func (s Service) Latest(ctx context.Context, userID, monitorID uint) (dto.MonitorLatestResponse, error) {
	monitor, err := s.store.FindMonitorByIDAndUser(ctx, monitorID, userID)
	if err != nil {
		return dto.MonitorLatestResponse{}, err
	}
	summary, err := s.analytics.Summary(ctx, monitorID)
	if err != nil {
		return dto.MonitorLatestResponse{}, err
	}
	var latest *dto.HealthCheckResponse
	check, err := s.store.LatestHealthCheck(ctx, monitorID)
	if err == nil {
		response := repository.HealthCheckToResponse(check)
		latest = &response
	}
	return dto.MonitorLatestResponse{Monitor: repository.MonitorToSummary(monitor), Latest: latest, Summary: summary}, nil
}

func (s Service) RunDueChecks(ctx context.Context) {
	monitors, err := s.store.ListMonitorsForScheduler(ctx)
	if err != nil {
		log.Printf("could not load monitors for scheduler: %v", err)
		return
	}
	now := time.Now()
	for _, monitor := range monitors {
		if !due(now, monitor) {
			continue
		}
		s.check(ctx, monitor)
	}
}

func (s Service) check(ctx context.Context, monitor models.Monitor) {
	previousUp := monitor.LatestIsUp
	result := s.ping(ctx, monitor)
	if err := s.store.CreateHealthCheck(ctx, &result); err != nil {
		log.Printf("could not save health check: monitor=%d error=%v", monitor.ID, err)
		return
	}

	update := map[string]any{
		"latest_status_code":   result.StatusCode,
		"latest_response_time": result.ResponseTimeMS,
		"latest_is_up":         result.IsUp,
		"last_checked_at":      result.CheckedAt,
	}
	if err := s.store.UpdateMonitor(ctx, &monitor, update); err != nil {
		log.Printf("could not update monitor status: monitor=%d error=%v", monitor.ID, err)
	}

	if previousUp != nil && *previousUp != result.IsUp {
		if result.IsUp {
			s.notifier.NotifyUp(monitor)
		} else {
			s.notifier.NotifyDown(monitor)
		}
	}
}

func (s Service) ping(ctx context.Context, monitor models.Monitor) models.HealthCheck {
	timeout := s.defaultTimeout
	if monitor.TimeoutSeconds > 0 {
		timeout = time.Duration(monitor.TimeoutSeconds) * time.Second
	}

	client := &http.Client{Timeout: timeout}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !monitor.ValidateSSL}
	client.Transport = transport
	if !monitor.FollowRedirects {
		client.CheckRedirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	var statusCode int
	var responseTime int64
	var err error
	for attempt := 0; attempt <= monitor.RetryCount; attempt++ {
		statusCode, responseTime, err = s.send(ctx, client, monitor)
		if err == nil && statusCode >= 200 && statusCode < 400 {
			break
		}
	}

	return models.HealthCheck{
		MonitorID:      monitor.ID,
		StatusCode:     statusCode,
		ResponseTimeMS: responseTime,
		IsUp:           err == nil && statusCode >= 200 && statusCode < 400,
		CheckedAt:      time.Now().UTC(),
	}
}

func (s Service) send(ctx context.Context, client *http.Client, monitor models.Monitor) (int, int64, error) {
	req, err := http.NewRequestWithContext(ctx, monitor.Method, monitor.URL, nil)
	if err != nil {
		return 0, 0, err
	}

	headers := map[string]string{}
	if len(monitor.Headers) > 0 {
		if err := json.Unmarshal(monitor.Headers, &headers); err != nil {
			log.Printf("could not parse monitor headers: monitor=%d error=%v", monitor.ID, err)
		}
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	start := time.Now()
	res, err := client.Do(req)
	responseTime := time.Since(start).Milliseconds()
	if err != nil {
		return 0, responseTime, err
	}
	defer res.Body.Close()
	return res.StatusCode, responseTime, nil
}

func validatePayload(input dto.MonitorRequest) error {
	if input.Name == "" {
		return errors.New("name is required")
	}
	parsed, err := http.NewRequest(http.MethodGet, input.URL, nil)
	if err != nil || parsed.URL.Scheme == "" || parsed.URL.Host == "" {
		return errors.New("url must be a valid absolute URL")
	}
	if parsed.URL.Scheme != "http" && parsed.URL.Scheme != "https" {
		return errors.New("url must use http or https")
	}
	if input.IntervalSeconds != 0 && input.IntervalSeconds < 60 {
		return errors.New("interval_seconds must be at least 60")
	}
	if input.TimeoutSeconds != 0 && (input.TimeoutSeconds < 1 || input.TimeoutSeconds > 60) {
		return errors.New("timeout_seconds must be between 1 and 60")
	}
	if input.RetryCount < 0 || input.RetryCount > 3 {
		return errors.New("retry_count must be between 0 and 3")
	}
	return nil
}

func toModel(input dto.MonitorRequest, userID uint) models.Monitor {
	method := strings.ToUpper(input.Method)
	if method == "" {
		method = http.MethodGet
	}
	interval := input.IntervalSeconds
	if interval == 0 {
		interval = 60
	}
	timeout := input.TimeoutSeconds
	if timeout == 0 {
		timeout = 10
	}
	followRedirects := true
	if input.FollowRedirects != nil {
		followRedirects = *input.FollowRedirects
	}
	validateSSL := true
	if input.ValidateSSL != nil {
		validateSSL = *input.ValidateSSL
	}
	return models.Monitor{
		UserID:          userID,
		Name:            input.Name,
		URL:             input.URL,
		IntervalSeconds: interval,
		TimeoutSeconds:  timeout,
		Method:          method,
		Headers:         headersToJSON(input.Headers),
		FollowRedirects: followRedirects,
		ValidateSSL:     validateSSL,
		RetryCount:      input.RetryCount,
	}
}

func headersToJSON(headers map[string]string) []byte {
	if len(headers) == 0 {
		return []byte("{}")
	}
	bytes, err := json.Marshal(headers)
	if err != nil {
		return []byte("{}")
	}
	return bytes
}

func due(now time.Time, monitor models.Monitor) bool {
	if monitor.LastCheckedAt == nil {
		return true
	}
	interval := monitor.IntervalSeconds
	if interval < 60 {
		interval = 60
	}
	return now.Sub(*monitor.LastCheckedAt) >= time.Duration(interval)*time.Second
}
