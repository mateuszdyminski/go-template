package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")
var xAPIVersion = http.CanonicalHeaderKey("X-API-Version")

type LoggingMiddleware struct {
	logger *zap.SugaredLogger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger.Sugar(),
	}
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		interceptor := &interceptor{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(interceptor, r)
		var (
			status = strconv.Itoa(interceptor.statusCode)
			took   = time.Since(begin)
		)
		m.logger.Debugw(
			"req",
			zap.String("requestId", GetReqID(r.Context())),
			zap.String("proto", r.Proto),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("remote", getRealIP(r)),
			zap.String("status", status),
			zap.String("user-agent", r.UserAgent()),
			zap.String("took", took.String()),
		)
	})
}

func getRealIP(r *http.Request) string {
	var ip string

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else {
		ip = r.RemoteAddr
	}

	return ip
}

func VersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(xAPIVersion, APIVersion)

		next.ServeHTTP(w, r)
	})
}
