package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		windowStart := now.Add(-rl.window)

		// Clean old requests
		if requests, exists := rl.requests[ip]; exists {
			validRequests := make([]time.Time, 0)
			for _, reqTime := range requests {
				if reqTime.After(windowStart) {
					validRequests = append(validRequests, reqTime)
				}
			}
			rl.requests[ip] = validRequests
		}

		// Check rate limit
		if len(rl.requests[ip]) >= rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": rl.window.Seconds(),
			})
			c.Abort()
			return
		}

		// Add current request
		rl.requests[ip] = append(rl.requests[ip], now)
		c.Next()
	}
}

// Global rate limiter instance
var GlobalRateLimiter = NewRateLimiter(100, time.Minute) // 100 requests per minute

func RateLimitMiddleware() gin.HandlerFunc {
	return GlobalRateLimiter.RateLimitMiddleware()
}
