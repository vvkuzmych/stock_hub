package middleware

import (
	"net/http"
)

// Manager manages all middleware
type Manager struct {
	// Add configuration here if needed
}

// NewManager creates a new middleware manager
func NewManager() *Manager {
	return &Manager{}
}

// Chain chains multiple middleware together
func (m *Manager) Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	// Apply middleware in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
