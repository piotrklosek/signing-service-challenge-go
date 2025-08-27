package api

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// create type to avoid collision
type jtwKey string
type loggerKey string
type reqID string

const (
	JWTDataKey jtwKey    = "JWTData"
	LoggerKey  loggerKey = "Logger"
	RequestID  reqID     = "reqID"
)

type JWTData struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type apiHandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func errorCheckerMiddleware(f func(http.ResponseWriter, *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := f(w, r)
		// handle error for handlers
		if err != nil {

		}
	}
}

func WithJWTData(ctx context.Context, v *JWTData) context.Context {
	return context.WithValue(ctx, JWTDataKey, v)
}

func GetJWTData(ctx context.Context) *JWTData {
	return ctx.Value(JWTDataKey).(*JWTData)
}

func WithLogger(ctx context.Context, v *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, v)
}

func GetLoger(ctx context.Context) *zap.Logger {
	return ctx.Value(LoggerKey).(*zap.Logger)
}

func middleware(logger *zap.Logger, fn func(http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := genRequestID()
		ip := clientIP(r)

		// propagate context with request ID
		ctx := context.WithValue(r.Context(), RequestID, reqID)
		r = r.WithContext(ctx)

		err := fn(w, r)
		if err != nil {
			logger.Error("HTTP_REQUEST_ERROR",
				zap.String("request_id", reqID),
				zap.String("method", r.Method),
				zap.String("route", r.URL.Path),
				zap.String("ip", ip),
				zap.Error(err),
				zap.Duration("latency", time.Since(start)),
			)
			return
		}

		// // TODO: custom response encoder — np. JSON
		// if result != nil {
		// 	_ = json.NewEncoder(w).Encode(result)
		// }

		// log successful request
		logger.Info("HTTP_REQUEST",
			zap.String("request_id", reqID),
			zap.String("method", r.Method),
			zap.String("route", r.URL.Path),
			zap.String("ip", ip),
			zap.Duration("latency", time.Since(start)),
		)
	})
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func genRequestID() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
