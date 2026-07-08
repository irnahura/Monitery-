package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"peekaping/backend/internal/auth"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func RequireAuth(validator interface {
	ValidateAPIKey(context.Context, string) (uint, bool)
}, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := bearerUser(c, jwtSecret)
		if !ok {
			userID, ok = validator.ValidateAPIKey(c.Request.Context(), c.GetHeader("X-API-Key"))
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) uint {
	value, _ := c.Get(UserIDKey)
	userID, _ := value.(uint)
	return userID
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	type bucket struct {
		count int
		reset time.Time
	}
	var mu sync.Mutex
	buckets := map[string]bucket{}

	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		now := time.Now()
		mu.Lock()
		current := buckets[key]
		if now.After(current.reset) {
			current = bucket{reset: now.Add(window)}
		}
		current.count++
		buckets[key] = current
		mu.Unlock()

		if current.count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

func bearerUser(c *gin.Context, secret string) (uint, bool) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return 0, false
	}
	userID, err := auth.ParseJWT(strings.TrimPrefix(header, "Bearer "), secret)
	return userID, err == nil
}
