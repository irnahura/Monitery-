package monitor

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var allowedMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodOptions: true,
}

type Payload struct {
	Name            string            `json:"name" binding:"required"`
	URL             string            `json:"url" binding:"required"`
	IntervalSeconds int               `json:"interval_seconds"`
	TimeoutSeconds  int               `json:"timeout_seconds"`
	Method          string            `json:"method"`
	Headers         map[string]string `json:"headers"`
	FollowRedirects *bool             `json:"follow_redirects"`
	ValidateSSL     *bool             `json:"validate_ssl"`
	RetryCount      int               `json:"retry_count"`
}

func (p Payload) Validate() error {
	parsed, err := url.ParseRequestURI(p.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("url must be a valid absolute URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("url must use http or https")
	}
	if p.IntervalSeconds != 0 && p.IntervalSeconds < 60 {
		return errors.New("interval_seconds must be at least 60")
	}
	if p.TimeoutSeconds != 0 && (p.TimeoutSeconds < 1 || p.TimeoutSeconds > 60) {
		return errors.New("timeout_seconds must be between 1 and 60")
	}
	if p.RetryCount < 0 || p.RetryCount > 3 {
		return errors.New("retry_count must be between 0 and 3")
	}
	method := strings.ToUpper(defaultString(p.Method, http.MethodGet))
	if !allowedMethods[method] {
		return errors.New("unsupported method")
	}
	return nil
}

func (p Payload) HeadersJSON() []byte {
	bytes, _ := json.Marshal(p.Headers)
	if string(bytes) == "null" {
		return []byte("{}")
	}
	return bytes
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
