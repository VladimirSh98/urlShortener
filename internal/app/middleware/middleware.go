package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func createCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{w, 0, 220}
}

func (lrw *customResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *customResponseWriter) Write(body []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(body)
	lrw.size += n
	return n, err
}

func Config(h http.Handler) http.Handler {
	logFn := func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		customWriter := createCustomResponseWriter(writer)
		uri := request.RequestURI
		method := request.Method
		h.ServeHTTP(customWriter, request)
		duration := time.Since(start)
		sugar := zap.S()
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
		sugar.Infoln(
			"responseStatus", customWriter.status,
			"responseSize", customWriter.size,
		)

	}
	return http.HandlerFunc(logFn)
}
