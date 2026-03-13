package handler

import (
	"strconv"

	"pos-system/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) Dashboard(c *fiber.Ctx) error {
	response, err := h.service.Dashboard(c.Context())
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(response)
}

func (h *ReportHandler) TopProducts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	response, err := h.service.TopProducts(c.Context(), c.Query("date_from"), c.Query("date_to"), limit)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(fiber.Map{"items": response})
}

func (h *ReportHandler) LowStock(c *fiber.Ctx) error {
	threshold, _ := strconv.Atoi(c.Query("threshold", "0"))
	response, err := h.service.LowStock(c.Context(), threshold)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(fiber.Map{"items": response})
}
