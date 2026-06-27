package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"reflect-backend/internal/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.RWMutex
	clients  map[string]*client
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(rps int, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rate:    rate.Limit(rps),
		burst:   burst,
	}

	go rl.cleanup()

	return rl
}

func NewRateLimiterFromConfig(cfg *config.Config, endpoint string) *RateLimiter {
	var rps int

	switch endpoint {
	case "register":
		rps = cfg.RateLimitRegister
	case "login":
		rps = cfg.RateLimitLogin
	case "checkout":
		rps = cfg.RateLimitCheckout
	case "upload":
		rps = cfg.RateLimitUpload
	case "admin":
		rps = cfg.RateLimitAdmin
	default:
		rps = 10
	}

	burst := rps * cfg.RateLimitBurst

	if burst < 1 {
		burst = 1
	}

	return NewRateLimiter(rps, burst)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)

	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()

		for ip, client := range rl.clients {
			if time.Since(client.lastSeen) > 10*time.Minute {
				delete(rl.clients, ip)
			}
		}

		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[key]

	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.clients[key] = &client{limiter: limiter, lastSeen: time.Now()}

		return limiter
	}

	c.lastSeen = time.Now()

	return c.limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		headers := rl.getRateLimitHeaders(limiter)
		for k, v := range headers {
			c.Header(k, v)
		}

		if !limiter.Allow() {
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Rate limit exceeded, please slow down",
			})
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetCurrentUserID(c)

		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authentication required for rate limiting",
			})
			return
		}

		limiter := rl.getLimiter(userID)

		headers := rl.getRateLimitHeaders(limiter)
		for k, v := range headers {
			c.Header(k, v)
		}

		if !limiter.Allow() {
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Rate limit exceeded, please slow down",
			})
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) getRateLimitHeaders(limiter *rate.Limiter) map[string]string {
	tokens := limiter.Tokens()

	return map[string]string{
		"X-RateLimit-Limit":     strconv.Itoa(int(rl.rate)),
		"X-RateLimit-Remaining": strconv.Itoa(int(tokens)),
		"X-RateLimit-Reset":     strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10),
	}
}
