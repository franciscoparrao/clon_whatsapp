package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LoggerMiddleware handles HTTP request logging
type LoggerMiddleware struct {
	logger   *log.Logger
	logLevel LogLevel
}

// responseWriter is a wrapper around http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

// NewLoggerMiddleware creates a new logger middleware instance
func NewLoggerMiddleware(logger *log.Logger, logLevel LogLevel) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger:   logger,
		logLevel: logLevel,
	}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Log is the middleware function that logs HTTP requests
func (m *LoggerMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Get request ID or generate one
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Add request ID to response headers
		wrapped.Header().Set("X-Request-ID", requestID)

		// Log incoming request
		if m.logLevel <= DEBUG {
			m.logger.Printf("[%s] [REQUEST] %s %s %s from %s",
				requestID,
				r.Method,
				r.URL.Path,
				r.Proto,
				r.RemoteAddr,
			)
		}

		// Process the request
		next.ServeHTTP(wrapped, r)

		// Calculate request duration
		duration := time.Since(start)

		// Extract user information if available
		userID, _ := r.Context().Value("userID").(string)
		if userID == "" {
			userID = "anonymous"
		}

		// Log the response
		level := m.getLogLevel(wrapped.statusCode)
		if level >= m.logLevel {
			m.logger.Printf("[%s] [RESPONSE] %s %s %d %s %d bytes user:%s duration:%v",
				requestID,
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				http.StatusText(wrapped.statusCode),
				wrapped.written,
				userID,
				duration,
			)
		}

		// Log errors in detail
		if wrapped.statusCode >= 400 && m.logLevel <= ERROR {
			m.logger.Printf("[%s] [ERROR] %s %s failed with status %d for user %s",
				requestID,
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				userID,
			)
		}
	})
}

// LogRequest logs a specific request with custom level and message
func (m *LoggerMiddleware) LogRequest(r *http.Request, level LogLevel, format string, args ...interface{}) {
	if level < m.logLevel {
		return
	}

	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	levelStr := m.getLevelString(level)
	message := fmt.Sprintf(format, args...)

	m.logger.Printf("[%s] [%s] %s", requestID, levelStr, message)
}

// getLogLevel determines the log level based on status code
func (m *LoggerMiddleware) getLogLevel(statusCode int) LogLevel {
	switch {
	case statusCode >= 500:
		return ERROR
	case statusCode >= 400:
		return WARN
	case statusCode >= 300:
		return INFO
	default:
		return DEBUG
	}
}

// getLevelString returns string representation of log level
func (m *LoggerMiddleware) getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().Unix(), time.Now().Nanosecond())
}

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Recovery middleware to recover from panics
func Recovery(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("[PANIC] %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}