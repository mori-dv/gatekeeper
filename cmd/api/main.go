package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	// chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gatekeeper/internal/config"
	"gatekeeper/internal/handler"
	"gatekeeper/internal/limiter"
	"gatekeeper/internal/logger"
	"gatekeeper/internal/metrics"
	"gatekeeper/internal/middleware"
	"gatekeeper/internal/storage"
)

func main() {
	cfg := config.Load()

	logg := logger.New()

	metrics.Init()

	redisClient := storage.NewRedis(cfg.RedisURL)

	ctx := context.Background()

	if err := storage.PingRedis(ctx, redisClient); err != nil {
		log.Fatal(err)
	}

	limit := limiter.New(
		redisClient,
		cfg.GetRateLimit(),
		time.Minute,
	)

	r := chi.NewRouter()

	//r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logg))
	r.Use(middleware.RateLimit(limit))

	r.Use(metrics.MetricsMiddleware)

	r.Get("/", handler.Root)
	r.Get("/health", handler.Healthcheck)
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logg.Info("server started", "port", cfg.AppPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	logg.Info("shutting down server")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal(err)
	}

	logg.Info("server stopped")
}
