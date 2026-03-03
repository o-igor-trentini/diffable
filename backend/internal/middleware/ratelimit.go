package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ipEntry struct {
	tokens    float64
	lastCheck time.Time
}

// RateLimiter implements a token bucket rate limiter per IP.
type RateLimiter struct {
	mu       sync.Mutex
	ips      map[string]*ipEntry
	rpm      int
	interval time.Duration
}

// NewRateLimiter creates a rate limiter with the given requests per minute.
func NewRateLimiter(rpm int) *RateLimiter {
	rl := &RateLimiter{
		ips:      make(map[string]*ipEntry),
		rpm:      rpm,
		interval: time.Minute,
	}

	go rl.cleanup()

	return rl
}

// Handler returns an HTTP middleware that enforces rate limiting.
func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !rl.allow(ip) {
			w.Header().Set("Retry-After", "60")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "rate_limited",
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.ips[ip]

	if !exists {
		rl.ips[ip] = &ipEntry{
			tokens:    float64(rl.rpm) - 1,
			lastCheck: now,
		}
		return true
	}

	elapsed := now.Sub(entry.lastCheck)
	entry.lastCheck = now

	// Replenish tokens based on elapsed time
	rate := float64(rl.rpm) / rl.interval.Seconds()
	entry.tokens += elapsed.Seconds() * rate

	if entry.tokens > float64(rl.rpm) {
		entry.tokens = float64(rl.rpm)
	}

	if entry.tokens < 1 {
		return false
	}

	entry.tokens--
	return true
}

// cleanup removes stale IP entries every 5 minutes.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for ip, entry := range rl.ips {
			if entry.lastCheck.Before(cutoff) {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RetryAfterHeader returns the Retry-After value in seconds.
func RetryAfterHeader(rpm int) string {
	return strconv.Itoa(60)
}
