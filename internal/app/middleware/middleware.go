package middleware

import (
	"compress/gzip"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

func createCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {

	return &customResponseWriter{ResponseWriter: w, size: 0, status: 200}
}

func (lrw *customResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.once.Do(func() { lrw.ResponseWriter.WriteHeader(code) })
}

func (lrw *customResponseWriter) Write(body []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(body)
	lrw.size += n
	return n, err
}

func (lrw *compressWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.once.Do(func() { lrw.ResponseWriter.WriteHeader(code) })
}

func (lrw *compressWriter) Write(body []byte) (int, error) {
	n, err := lrw.Writer.Write(body)
	lrw.size += n
	return n, err
}

func Config(h http.Handler) http.Handler {
	logFn := func(writer http.ResponseWriter, request *http.Request) {
		var responseStatus, responseSize int

		start := time.Now()
		customWriter := createCustomResponseWriter(writer)
		uri := request.RequestURI
		method := request.Method

		contentEncoding := request.Header.Get("Content-Encoding")
		sendCompress := strings.Contains(contentEncoding, "gzip")
		if sendCompress {
			cr, err := gzip.NewReader(request.Body)
			if err != nil {
				customWriter.WriteHeader(http.StatusInternalServerError)
				return
			}
			request.Body = cr
			defer cr.Close()
		}

		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			gzipWriter := gzip.NewWriter(customWriter)
			defer gzipWriter.Close()
			customCompressWriter := compressWriter{customResponseWriter: *customWriter, Writer: gzipWriter}
			customCompressWriter.Header().Set("Content-Encoding", "gzip")
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
