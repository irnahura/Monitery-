package api

import (
	"context"
	"net/http"
	"strconv"

	"peekaping/backend/internal/dto"
	"peekaping/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(context.Context, dto.AuthRequest) (dto.AuthResponse, error)
	Login(context.Context, dto.LoginRequest) (dto.AuthResponse, error)
	Profile(context.Context, uint) (dto.UserSummary, error)
}

type MonitorService interface {
	List(context.Context, uint) ([]dto.MonitorSummary, error)
	Create(context.Context, uint, dto.MonitorRequest) (dto.MonitorSummary, error)
	Update(context.Context, uint, uint, dto.MonitorRequest) (dto.MonitorSummary, error)
	Delete(context.Context, uint, uint) error
	History(context.Context, uint, uint, int) (dto.MonitorHistoryResponse, error)
	Latest(context.Context, uint, uint) (dto.MonitorLatestResponse, error)
}

type APIKeyService interface {
	List(context.Context, uint) ([]dto.APIKeySummary, error)
	Create(context.Context, uint, dto.CreateAPIKeyRequest) (dto.APIKeyCreateResponse, error)
	Delete(context.Context, uint, uint) error
}

type Handler struct {
	Auth     AuthService
	Monitors MonitorService
	APIKeys  APIKeyService
}

func (h Handler) Register(c *gin.Context) {
	var input dto.AuthRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "invalid registration payload")
		return
	}
	result, err := h.Auth.Register(c.Request.Context(), input)
	if err != nil {
		respondError(c, http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h Handler) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "invalid login payload")
		return
	}
	result, err := h.Auth.Login(c.Request.Context(), input)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handler) Profile(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	result, err := h.Auth.Profile(c.Request.Context(), userID)
	if err != nil {
		respondError(c, http.StatusNotFound, "user not found")
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handler) ListMonitors(c *gin.Context) {
	monitors, err := h.Monitors.List(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not list monitors")
		return
	}
	c.JSON(http.StatusOK, monitors)
}

func (h Handler) CreateMonitor(c *gin.Context) {
	var input dto.MonitorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "invalid monitor payload")
		return
	}
	result, err := h.Monitors.Create(c.Request.Context(), middleware.CurrentUserID(c), input)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h Handler) UpdateMonitor(c *gin.Context) {
	var input dto.MonitorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "invalid monitor payload")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid monitor id")
		return
	}
	result, err := h.Monitors.Update(c.Request.Context(), middleware.CurrentUserID(c), uint(id), input)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handler) DeleteMonitor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid monitor id")
		return
	}
	if err := h.Monitors.Delete(c.Request.Context(), middleware.CurrentUserID(c), uint(id)); err != nil {
		respondError(c, http.StatusInternalServerError, "could not delete monitor")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) MonitorHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid monitor id")
		return
	}
	limit := queryLimit(c, 100)
	result, err := h.Monitors.History(c.Request.Context(), middleware.CurrentUserID(c), uint(id), limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not load history")
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handler) MonitorLatest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid monitor id")
		return
	}
	result, err := h.Monitors.Latest(c.Request.Context(), middleware.CurrentUserID(c), uint(id))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not load latest check")
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handler) ListAPIKeys(c *gin.Context) {
	keys, err := h.APIKeys.List(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not list API keys")
		return
	}
	c.JSON(http.StatusOK, keys)
}

func (h Handler) CreateAPIKey(c *gin.Context) {
	var input dto.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "invalid API key payload")
		return
	}
	result, err := h.APIKeys.Create(c.Request.Context(), middleware.CurrentUserID(c), input)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not create API key")
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h Handler) DeleteAPIKey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid API key id")
		return
	}
	if err := h.APIKeys.Delete(c.Request.Context(), middleware.CurrentUserID(c), uint(id)); err != nil {
		respondError(c, http.StatusInternalServerError, "could not delete API key")
		return
	}
	c.Status(http.StatusNoContent)
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, dto.ErrorResponse{Error: message})
}

func queryLimit(c *gin.Context, fallback int) int {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(fallback)))
	if err != nil || limit <= 0 || limit > 500 {
		return fallback
	}
	return limit
}
