package middleware

import (
	"github.com/andybalholm/brotli"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

func createCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {

	return &customResponseWriter{w, 0, 200}
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

func (lrw *compressWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *compressWriter) Write(body []byte) (int, error) {
	n, err := lrw.Writer.Write(body)
	lrw.size += n
	return n, err
}

func Config(h http.Handler) http.Handler {
	logFn := func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		customWriter := createCustomResponseWriter(writer)
		uri := request.RequestURI
		method := request.Method
		var responseStatus, responseSize int
		contentEncoding := request.Header.Get("Content-Encoding")
		sendCompress := strings.Contains(contentEncoding, "br")
		if sendCompress {
			cr := brotli.NewReader(request.Body)
			request.Body = io.NopCloser(cr)
		}
		if strings.Contains(request.Header.Get("Accept-Encoding"), "br") {
			brWriter := brotli.NewWriterLevel(customWriter, brotli.BestCompression)
			defer brWriter.Close()
			customCompressWriter := compressWriter{*customWriter, brWriter}
			customCompressWriter.Header().Set("Content-Encoding", "br")
			h.ServeHTTP(&customCompressWriter, request)
			responseStatus = customCompressWriter.status
			responseSize = customCompressWriter.size
		} else {
			h.ServeHTTP(customWriter, request)
			responseStatus = customWriter.status
			responseSize = customWriter.size
		}
		duration := time.Since(start)
		sugar := zap.S()
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
		sugar.Infoln(
			"responseStatus", responseStatus,
			"responseSize", responseSize,
		)

	}
	return http.HandlerFunc(logFn)
}
