package handler

import (
	"strconv"

	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	var active *bool
	if raw := c.Query("active"); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err == nil {
			active = &value
		}
	}

	result, err := h.service.List(c.Context(), domain.ProductListFilter{
		Query:  c.Query("q"),
		Page:   page,
		Limit:  limit,
		Active: active,
	})
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(result)
}

func (h *ProductHandler) Search(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	products, err := h.service.Search(c.Context(), c.Query("q"), limit)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(fiber.Map{"items": products})
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var input domain.CreateProductInput
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, err)
	}

	product, err := h.service.Create(c.Context(), input)
	if err != nil {
		return writeError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return writeError(c, err)
	}

	var input domain.UpdateProductInput
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, err)
	}

	product, err := h.service.Update(c.Context(), id, input)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(product)
}
