package service

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/handler"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
)

// Run service
func Run(ctx context.Context) error {
	sugar := zap.S()
	sugar.Infoln("Prefill data success")
	repo := dbRepo.NewShortenRepository(database.DBConnection.Conn)
	service := shorten.NewShortenService(repo)
	customHandler := handler.NewHandler(service)
	err := Prefill(service)
	if err != nil {
		sugar.Warnln("Prefill data error %v", err)
	}
	router := chi.NewMux()
	router.Use(middleware.Config)
	router.Mount("/debug", chiMiddleware.Profiler())
	router.Get("/ping", customHandler.Ping)
	router.Post("/", customHandler.ManagerCreateShortURL)
	router.Post("/api/shorten", customHandler.ManagerCreateShortURLByJSON)
	router.Post("/api/shorten/batch", customHandler.ManagerCreateShortURLBatch)
	router.Get("/{id}", customHandler.ManagerReturnFullURL)
	router.Get("/api/user/urls", customHandler.ManagerGetURLsByUser)
	router.Delete("/api/user/urls", customHandler.ManagerDeleteURLsByID)
	subnetGroup := router.Group(nil)
	subnetGroup.Use(middleware.CheckTrustedSubnet)
	subnetGroup.Get("/api/internal/stats", customHandler.GetStats)

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		sugar.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	m := cmux.New(lis)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	server := &http.Server{
		Handler: router,
	}

	go func() {
		sugar.Infof("Starting gRPC server on %s", config.FlagRunAddr)
		grpcServer.Serve(grpcL)
	}()

	go func() {
		if config.EnableHTTPS {
			err = server.ListenAndServeTLS(config.CertFile, config.KeyFile)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			sugar.Errorf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	sugar.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		sugar.Errorf("Shutdown failed: %v", err)
		return err
	}

	sugar.Info("Server stopped gracefully")
	return nil
}
