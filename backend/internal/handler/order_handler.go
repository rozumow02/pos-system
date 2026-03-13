package handler

import (
	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var input domain.CreateOrderInput
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, err)
	}

	receipt, err := h.service.Create(c.Context(), input)
	if err != nil {
		return writeError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(receipt)
}

func (h *OrderHandler) Today(c *fiber.Ctx) error {
	response, err := h.service.Today(c.Context())
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(response)
}
