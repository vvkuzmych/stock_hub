package http

import (
	"net/http"
)

// Handler holds HTTP handler dependencies
type Handler struct {
	// Add dependencies here if needed (e.g., services)
}

// NewHandler creates a new HTTP handler
func NewHandler() *Handler {
	return &Handler{}
}

// ServeIndex serves the index.html file
func (h *Handler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// ServeAssets serves static assets
func (h *Handler) ServeAssets() http.Handler {
	return http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets")))
}

// HealthCheck returns server health status
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
