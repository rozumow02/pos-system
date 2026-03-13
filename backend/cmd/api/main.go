package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pos-system/backend/internal/config"
	"pos-system/backend/internal/handler"
	apphttp "pos-system/backend/internal/http"
	"pos-system/backend/internal/platform"
	"pos-system/backend/internal/repository"
	"pos-system/backend/internal/service"
)

func main() {
	cfg := config.Load()
	logger := platform.NewLogger(cfg.LogLevel)

	ctx := context.Background()
	dbpool, err := platform.ConnectDatabase(ctx, cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer dbpool.Close()

	if err := platform.RunMigrations(ctx, dbpool, cfg.MigrationsPath, logger); err != nil {
		logger.Fatal().Err(err).Msg("failed to run migrations")
	}

	productRepo := repository.NewProductRepository(dbpool)
	orderRepo := repository.NewOrderRepository(dbpool)
	reportRepo := repository.NewReportRepository(dbpool)

	productService := service.NewProductService(dbpool, productRepo, orderRepo)
	orderService := service.NewOrderService(dbpool, productRepo, orderRepo)
	reportService := service.NewReportService(reportRepo, cfg.LowStockThreshold)

	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	reportHandler := handler.NewReportHandler(reportService)
	healthHandler := handler.NewHealthHandler()

	app := apphttp.NewApp(cfg, logger, productHandler, orderHandler, reportHandler, healthHandler)

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- app.Listen(":" + cfg.BackendPort)
	}()

	logger.Info().Str("port", cfg.BackendPort).Msg("server started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
	case err := <-listenErr:
		if err != nil {
			logger.Fatal().Err(err).Msg("fiber server stopped")
		}
		return
	}

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	}

	if err := <-listenErr; err != nil {
		logger.Info().Err(err).Msg("server shutdown complete")
	}
}
