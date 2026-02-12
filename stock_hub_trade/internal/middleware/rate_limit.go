package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	mu              sync.Mutex
	clients         map[string]*bucket
	rate            int           // requests per window
	window          time.Duration // time window
	cleanupInterval time.Duration // cleanup interval
}

type bucket struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:         make(map[string]*bucket),
		rate:            rate,
		window:          window,
		cleanupInterval: window * 10,
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// RateLimit middleware
func (m *Manager) RateLimit(limiter *RateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if !limiter.Allow(clientIP) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next(w, r)
		}
	}
}

// Allow checks if request is allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get or create bucket
	b, exists := rl.clients[clientID]
	if !exists {
		b = &bucket{
			tokens:   rl.rate - 1,
			lastSeen: now,
		}
		rl.clients[clientID] = b
		return true
	}

	// Refill tokens based on time passed
	elapsed := now.Sub(b.lastSeen)
	if elapsed >= rl.window {
		b.tokens = rl.rate
		b.lastSeen = now
	}

	// Check if tokens available
	if b.tokens > 0 {
		b.tokens--
		b.lastSeen = now
		return true
	}

	return false
}

// cleanupLoop removes old clients
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanupOldClients()
	}
}

// cleanupOldClients removes inactive clients
func (rl *RateLimiter) cleanupOldClients() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for clientID, b := range rl.clients {
		if now.Sub(b.lastSeen) > rl.cleanupInterval {
			delete(rl.clients, clientID)
		}
	}
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// Check X-Real-IP header
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Use RemoteAddr
	return r.RemoteAddr
}
