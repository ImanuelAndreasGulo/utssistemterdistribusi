package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-counter-go/internal/api"
	"distributed-counter-go/internal/cluster"
	"distributed-counter-go/internal/config"
	"distributed-counter-go/internal/counter"
	"distributed-counter-go/internal/storage"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	store := storage.NewJSONStore(cfg.DataPath, logger)
	ctr, err := counter.NewCounter(cfg.NodeID, store, logger)
	if err != nil {
		logger.Fatal("failed to initialize counter", zap.Error(err))
	}

	mgr := cluster.NewManager(cfg, ctr, logger)

	router := api.NewRouter(cfg, ctr, mgr, logger)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Heartbeat(ctx)
	go mgr.SyncLoop(ctx)

	go func() {
		logger.Info("server started", zap.String("node", cfg.NodeID), zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}
}
