package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	maxReqs  int
	window   time.Duration
}

func NewRateLimiter(maxReqs int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		maxReqs:  maxReqs,
		window:   window,
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	timestamps := rl.requests[ip]
	var valid []time.Time
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.maxReqs {
		rl.requests[ip] = valid
		return false
	}

	valid = append(valid, now)
	rl.requests[ip] = valid
	return true
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	for ip, timestamps := range rl.requests {
		var valid []time.Time
		for _, t := range timestamps {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = valid
		}
	}
}

var globalLimiter = NewRateLimiter(100, time.Minute)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c.Request)

		if !globalLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Terlalu banyak permintaan. Silakan coba lagi dalam beberapa saat.",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP != "" {
		if idx := strings.IndexByte(clientIP, ','); idx != -1 {
			clientIP = clientIP[:idx]
		}
		clientIP = strings.TrimSpace(clientIP)
		return clientIP
	}

	clientIP = r.Header.Get("X-Real-IP")
	if clientIP != "" {
		return clientIP
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
