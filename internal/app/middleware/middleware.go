package middleware

import (
	"compress/gzip"
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func createCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {

	return &customResponseWriter{ResponseWriter: w, Size: 0, Status: 200}
}

// Config updates response, check authorization, performs compression
// and add custom logs with response status, size, duration
func Config(h http.Handler) http.Handler {
	logFn := func(writer http.ResponseWriter, request *http.Request) {
		var err error
		var responseStatus, responseSize int

		start := time.Now()
		customWriter := createCustomResponseWriter(writer)
		uri := request.RequestURI
		method := request.Method

		var token string
		var userID int
		token, userID, err = authorize(request)
		if err != nil {
			customWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(request.Context(), UserIDKey, userID)
		http.SetCookie(customWriter, &http.Cookie{Name: "Authorization", Value: token})

		contentEncoding := request.Header.Get("Content-Encoding")
		sendCompress := strings.Contains(contentEncoding, "gzip")
		if sendCompress {
			var cr *gzip.Reader
			cr, err = gzip.NewReader(request.Body)
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
			customCompressWriter := compressWriter{
				ResponseWriter: customWriter.ResponseWriter,
				Size:           customWriter.Size,
				Status:         customWriter.Status,
				Writer:         gzipWriter,
			}
			customCompressWriter.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(&customCompressWriter, request.WithContext(ctx))
			responseStatus = customCompressWriter.Status
			responseSize = customCompressWriter.Size
		} else {
			h.ServeHTTP(customWriter, request.WithContext(ctx))
			responseStatus = customWriter.Status
			responseSize = customWriter.Size
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
