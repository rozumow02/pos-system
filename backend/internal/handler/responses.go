package handler

import (
	"errors"
	"net/http"

	"pos-system/backend/internal/apperrors"

	"github.com/gofiber/fiber/v2"
)

func writeError(c *fiber.Ctx, err error) error {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		status = http.StatusBadRequest
	case errors.Is(err, apperrors.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, apperrors.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, apperrors.ErrInsufficientStock):
		status = http.StatusUnprocessableEntity
	}

	return c.Status(status).JSON(fiber.Map{
		"error": err.Error(),
	})
}
