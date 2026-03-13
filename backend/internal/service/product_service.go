package service

import (
	"context"
	"fmt"
	"strings"

	"pos-system/backend/internal/apperrors"
	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	db        *pgxpool.Pool
	products  *repository.ProductRepository
	inventory *repository.OrderRepository
}

func NewProductService(db *pgxpool.Pool, products *repository.ProductRepository, inventory *repository.OrderRepository) *ProductService {
	return &ProductService{
		db:        db,
		products:  products,
		inventory: inventory,
	}
}

func (s *ProductService) List(ctx context.Context, filter domain.ProductListFilter) (domain.ProductListResult, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 200 {
		filter.Limit = 50
	}

	items, total, err := s.products.List(ctx, filter)
	if err != nil {
		return domain.ProductListResult{}, err
	}

	return domain.ProductListResult{
		Items: items,
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}, nil
}

func (s *ProductService) Search(ctx context.Context, query string, limit int) ([]domain.Product, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.products.Search(ctx, strings.TrimSpace(query), limit)
}

func (s *ProductService) Create(ctx context.Context, input domain.CreateProductInput) (domain.Product, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return domain.Product{}, fmt.Errorf("%w: product name is required", apperrors.ErrValidation)
	}
	if input.Price < 0 {
		return domain.Product{}, fmt.Errorf("%w: product price must be non-negative", apperrors.ErrValidation)
	}
	if input.Stock < 0 {
		return domain.Product{}, fmt.Errorf("%w: product stock must be non-negative", apperrors.ErrValidation)
	}
	normalizeCreateProductInput(&input)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.Product{}, err
	}
	defer tx.Rollback(ctx)

	product, err := s.products.CreateTx(ctx, tx, input)
	if err != nil {
		return domain.Product{}, err
	}

	if input.Stock > 0 {
		if err := s.inventory.InsertInventoryMovementTx(ctx, tx, domain.InventoryMovementInput{
			ProductID: product.ID,
			ChangeQty: input.Stock,
			Reason:    "initial_stock",
		}); err != nil {
			return domain.Product{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (s *ProductService) Update(ctx context.Context, id int64, input domain.UpdateProductInput) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, fmt.Errorf("%w: invalid product id", apperrors.ErrValidation)
	}

	if err := validateUpdateProductInput(input); err != nil {
		return domain.Product{}, err
	}

	existing, err := s.products.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	normalizeUpdateProductInput(&input)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.Product{}, err
	}
	defer tx.Rollback(ctx)

	product, err := s.products.UpdateTx(ctx, tx, id, input)
	if err != nil {
		return domain.Product{}, err
	}

	if input.Stock != nil {
		delta := *input.Stock - existing.Stock
		if delta != 0 {
			if err := s.inventory.InsertInventoryMovementTx(ctx, tx, domain.InventoryMovementInput{
				ProductID: id,
				ChangeQty: delta,
				Reason:    "manual_adjustment",
			}); err != nil {
				return domain.Product{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func normalizeCreateProductInput(input *domain.CreateProductInput) {
	if input.SKU != nil {
		trimmed := strings.TrimSpace(*input.SKU)
		if trimmed == "" {
			input.SKU = nil
		} else {
			input.SKU = &trimmed
		}
	}
	if input.Barcode != nil {
		trimmed := strings.TrimSpace(*input.Barcode)
		if trimmed == "" {
			input.Barcode = nil
		} else {
			input.Barcode = &trimmed
		}
	}
}

func normalizeUpdateProductInput(input *domain.UpdateProductInput) {
	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		input.Name = &trimmed
	}
	if input.SKU != nil {
		trimmed := strings.TrimSpace(*input.SKU)
		input.SKU = &trimmed
	}
	if input.Barcode != nil {
		trimmed := strings.TrimSpace(*input.Barcode)
		input.Barcode = &trimmed
	}
}

func validateUpdateProductInput(input domain.UpdateProductInput) error {
	if input.Name == nil && input.SKU == nil && input.Barcode == nil && input.Price == nil && input.Stock == nil && input.IsActive == nil {
		return fmt.Errorf("%w: no fields to update", apperrors.ErrValidation)
	}
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return fmt.Errorf("%w: product name cannot be empty", apperrors.ErrValidation)
	}
	if input.Price != nil && *input.Price < 0 {
		return fmt.Errorf("%w: product price must be non-negative", apperrors.ErrValidation)
	}
	if input.Stock != nil && *input.Stock < 0 {
		return fmt.Errorf("%w: product stock must be non-negative", apperrors.ErrValidation)
	}
	return nil
}
