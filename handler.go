package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var config *Config

// loggingMiddleware wraps an HTTP handler with request logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next(wrapped, r)

		// Log the request
		duration := time.Since(start)
		LogRequest(r.Method, r.URL.Path, r.RemoteAddr, wrapped.statusCode, duration)
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// handlerM3U handles GET /m3u requests
func handlerM3U(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		LogWarn("Invalid method for /m3u: %s from %s", r.Method, r.RemoteAddr)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	LogInfo("Starting M3U playlist generation request from %s", r.RemoteAddr)

	content, err := generateM3UContent(&config.Xtream)
	if err != nil {
		LogError("M3U generation failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to generate M3U playlist",
			"message": err.Error(),
		})
		return
	}

	LogInfo("M3U playlist generated successfully (%d bytes)", len(content))
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Content-Disposition", "attachment; filename=\"playlist.m3u\"")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

// handlerHealth handles GET /health requests
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	LogDebug("Health check request from %s", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"version":   version,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// handlerConfig handles GET /config requests (debug endpoint)
func handlerConfig(w http.ResponseWriter, r *http.Request) {
	LogDebug("Config request from %s", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	// Create a copy and redact sensitive information
	redactedCfg := *config
	redactedCfg.Xtream.Password = "***REDACTED***"
	redactedCfg.Xtream.Username = "***REDACTED***"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(redactedCfg)
}
