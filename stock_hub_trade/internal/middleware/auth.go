package middleware

import (
	"context"
	"net/http"
	"strings"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// UsernameKey is the context key for username
	UsernameKey ContextKey = "username"
)

// Auth validates authentication token
func (m *Manager) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// TODO: Validate token (JWT, session, etc.)
		// For now, just check if token is not empty
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// TODO: Extract user info from token
		// For now, mock user data
		ctx := context.WithValue(r.Context(), UserIDKey, int64(1))
		ctx = context.WithValue(ctx, UsernameKey, "test_user")

		// Call next handler with updated context
		next(w, r.WithContext(ctx))
	}
}

// OptionalAuth validates token if present, but allows request without it
func (m *Manager) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No auth header - continue without user context
			next(w, r)
			return
		}

		// Has auth header - validate it
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			if token != "" {
				// TODO: Validate and extract user info
				ctx := context.WithValue(r.Context(), UserIDKey, int64(1))
				ctx = context.WithValue(ctx, UsernameKey, "test_user")
				r = r.WithContext(ctx)
			}
		}

		next(w, r)
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// GetUsername extracts username from context
func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}
