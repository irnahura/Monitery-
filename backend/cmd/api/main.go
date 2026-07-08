package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"peekaping/backend/internal/analytics"
	apihandler "peekaping/backend/internal/api"
	"peekaping/backend/internal/apikey"
	"peekaping/backend/internal/auth"
	"peekaping/backend/internal/config"
	"peekaping/backend/internal/database"
	"peekaping/backend/internal/monitor"
	"peekaping/backend/internal/notification"
	"peekaping/backend/internal/repository"
	"peekaping/backend/internal/scheduler"
)

func main() {
	cfg := config.Load()
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	store := repository.New(db)
	mailer := notification.NewEmailNotifier(cfg)
	analyticsService := analytics.NewService(store)
	authService := auth.NewService(store, cfg.JWTSecret)
	apiKeyService := apikey.NewService(store)
	monitorService := monitor.NewService(store, mailer, analyticsService, cfg.DefaultTimeout)
	handler := apihandler.Handler{Auth: authService, Monitors: monitorService, APIKeys: apiKeyService}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduler.New(monitorService, time.Minute).Run(ctx)

	router := apihandler.NewRouter(store, handler, cfg)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("backend listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}
}
