package http

import (
	"time"

	"pos-system/backend/internal/config"
	"pos-system/backend/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
)

func NewApp(
	cfg config.Config,
	logger zerolog.Logger,
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
	reportHandler *handler.ReportHandler,
	healthHandler *handler.HealthHandler,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "POS System API",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error().Err(err).Str("path", c.Path()).Msg("unhandled fiber error")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.FrontendOrigin,
		AllowMethods: "GET,POST,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))
	app.Use(requestLoggingMiddleware(logger))

	app.Get("/healthz", healthHandler.Health)

	api := app.Group("/api")
	api.Get("/products", productHandler.List)
	api.Post("/products", productHandler.Create)
	api.Patch("/products/:id", productHandler.Update)
	api.Get("/products/search", productHandler.Search)

	api.Post("/orders", orderHandler.Create)
	api.Get("/orders/today", orderHandler.Today)

	api.Get("/reports/dashboard", reportHandler.Dashboard)
	api.Get("/reports/top-products", reportHandler.TopProducts)
	api.Get("/reports/low-stock", reportHandler.LowStock)

	return app
}
