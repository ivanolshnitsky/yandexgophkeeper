package logger

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

var writerPool = sync.Pool{
	New: func() any {
		return &responseWriter{}
	},
}

func Initialize(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level %q: %w", level, err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	return cfg.Build()
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func Middleware(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := writerPool.Get().(*responseWriter)
		rw.ResponseWriter = w

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Info("request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", rw.status),
			zap.Int("size", rw.size),
			zap.Duration("duration", duration),
		)

		rw.ResponseWriter = nil
		writerPool.Put(rw)
	})
}
