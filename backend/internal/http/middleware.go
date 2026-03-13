package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func requestLoggingMiddleware(logger zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		event := logger.Info()
		if err != nil {
			event = logger.Error().Err(err)
		}

		event.
			Str("request_id", c.GetRespHeader(fiber.HeaderXRequestID)).
			Str("method", c.Method()).
			Str("path", c.OriginalURL()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Msg("request completed")

		return err
	}
}
