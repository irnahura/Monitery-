package api

import (
	"context"
	"net/http"
	"time"

	"peekaping/backend/internal/config"
	"peekaping/backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type APIKeyValidator interface {
	ValidateAPIKey(context.Context, string) (uint, bool)
}

func NewRouter(validator APIKeyValidator, handler Handler, cfg config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authLimit := middleware.RateLimit(10, time.Minute)
	apiLimit := middleware.RateLimit(120, time.Minute)

	authGroup := router.Group("/auth", authLimit)
	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)

	protected := router.Group("/", apiLimit, middleware.RequireAuth(validator, cfg.JWTSecret))
	protected.GET("/auth/profile", handler.Profile)
	protected.GET("/monitors", handler.ListMonitors)
	protected.POST("/monitors", handler.CreateMonitor)
	protected.PUT("/monitors/:id", handler.UpdateMonitor)
	protected.DELETE("/monitors/:id", handler.DeleteMonitor)
	protected.GET("/monitors/:id/history", handler.MonitorHistory)
	protected.GET("/monitors/:id/latest", handler.MonitorLatest)
	protected.GET("/apikeys", handler.ListAPIKeys)
	protected.POST("/apikeys", handler.CreateAPIKey)
	protected.DELETE("/apikeys/:id", handler.DeleteAPIKey)

	return router
}
