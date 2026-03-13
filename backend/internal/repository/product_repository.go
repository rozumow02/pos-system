package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"pos-system/backend/internal/apperrors"
	"pos-system/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) List(ctx context.Context, filter domain.ProductListFilter) ([]domain.Product, int, error) {
	offset := (filter.Page - 1) * filter.Limit
	query := strings.TrimSpace(filter.Query)

	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM products
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%' OR COALESCE(sku, '') ILIKE '%' || $1 || '%' OR COALESCE(barcode, '') ILIKE '%' || $1 || '%')
		  AND ($2::BOOLEAN IS NULL OR is_active = $2)
	`, query, filter.Active).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, name, sku, barcode, price::double precision, stock, is_active, created_at, updated_at
		FROM products
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%' OR COALESCE(sku, '') ILIKE '%' || $1 || '%' OR COALESCE(barcode, '') ILIKE '%' || $1 || '%')
		  AND ($2::BOOLEAN IS NULL OR is_active = $2)
		ORDER BY name ASC
		LIMIT $3 OFFSET $4
	`, query, filter.Active, filter.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0, filter.Limit)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, rows.Err()
}

func (r *ProductRepository) Search(ctx context.Context, query string, limit int) ([]domain.Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, sku, barcode, price::double precision, stock, is_active, created_at, updated_at
		FROM products
		WHERE is_active = TRUE
		  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR COALESCE(sku, '') ILIKE '%' || $1 || '%' OR COALESCE(barcode, '') ILIKE '%' || $1 || '%')
		ORDER BY name ASC
		LIMIT $2
	`, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, sku, barcode, price::double precision, stock, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id)

	product, err := scanProduct(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, fmt.Errorf("%w: product %d not found", apperrors.ErrNotFound, id)
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("get product by id: %w", err)
	}

	return product, nil
}

func (r *ProductRepository) CreateTx(ctx context.Context, tx pgx.Tx, input domain.CreateProductInput) (domain.Product, error) {
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO products (name, sku, barcode, price, stock, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, sku, barcode, price::double precision, stock, is_active, created_at, updated_at
	`, input.Name, input.SKU, input.Barcode, input.Price, input.Stock, isActive)

	product, err := scanProduct(row)
	if err != nil {
		return domain.Product{}, mapDatabaseError("create product", err)
	}

	return product, nil
}

func (r *ProductRepository) UpdateTx(ctx context.Context, tx pgx.Tx, id int64, input domain.UpdateProductInput) (domain.Product, error) {
	row := tx.QueryRow(ctx, `
		UPDATE products
		SET
			name = COALESCE($2, name),
			sku = COALESCE($3, sku),
			barcode = COALESCE($4, barcode),
			price = COALESCE($5, price),
			stock = COALESCE($6, stock),
			is_active = COALESCE($7, is_active),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, sku, barcode, price::double precision, stock, is_active, created_at, updated_at
	`, id, input.Name, input.SKU, input.Barcode, input.Price, input.Stock, input.IsActive)

	product, err := scanProduct(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, fmt.Errorf("%w: product %d not found", apperrors.ErrNotFound, id)
	}
	if err != nil {
		return domain.Product{}, mapDatabaseError("update product", err)
	}

	return product, nil
}

func (r *ProductRepository) GetByIDsForUpdate(ctx context.Context, tx pgx.Tx, ids []int64) ([]domain.Product, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, name, sku, barcode, price::double precision, stock, is_active, created_at, updated_at
		FROM products
		WHERE id = ANY($1)
		FOR UPDATE
	`, ids)
	if err != nil {
		return nil, fmt.Errorf("select products for update: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanProduct(s scanner) (domain.Product, error) {
	var product domain.Product
	var sku sql.NullString
	var barcode sql.NullString

	err := s.Scan(
		&product.ID,
		&product.Name,
		&sku,
		&barcode,
		&product.Price,
		&product.Stock,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return domain.Product{}, err
	}

	if sku.Valid {
		product.SKU = &sku.String
	}
	if barcode.Valid {
		product.Barcode = &barcode.String
	}

	return product, nil
}

func mapDatabaseError(action string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return fmt.Errorf("%w: duplicate product field", apperrors.ErrConflict)
	}
	return fmt.Errorf("%s: %w", action, err)
}
