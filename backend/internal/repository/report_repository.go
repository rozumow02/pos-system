package repository

import (
	"context"
	"fmt"

	"pos-system/backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	pool *pgxpool.Pool
}

func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{pool: pool}
}

func (r *ReportRepository) GetDashboard(ctx context.Context, threshold int) (domain.DashboardMetrics, error) {
	metrics := domain.DashboardMetrics{}

	err := r.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(total_amount)::double precision, 0),
			COUNT(*)
		FROM orders
		WHERE created_at >= date_trunc('day', NOW())
	`).Scan(&metrics.RevenueToday, &metrics.OrdersToday)
	if err != nil {
		return domain.DashboardMetrics{}, fmt.Errorf("dashboard totals: %w", err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(quantity), 0)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE o.created_at >= date_trunc('day', NOW())
	`).Scan(&metrics.ItemsSoldToday)
	if err != nil {
		return domain.DashboardMetrics{}, fmt.Errorf("dashboard item totals: %w", err)
	}

	topProducts, err := r.GetTopProducts(ctx, "", "", 5)
	if err != nil {
		return domain.DashboardMetrics{}, err
	}
	metrics.TopProducts = topProducts

	lowStock, err := r.GetLowStockProducts(ctx, threshold)
	if err != nil {
		return domain.DashboardMetrics{}, err
	}
	metrics.LowStock = lowStock
	metrics.LowStockCount = len(lowStock)

	return metrics, nil
}

func (r *ReportRepository) GetTopProducts(ctx context.Context, dateFrom, dateTo string, limit int) ([]domain.TopProduct, error) {
	args := []any{limit}
	where := `
		WHERE o.created_at >= date_trunc('day', NOW())
		  AND o.created_at < date_trunc('day', NOW()) + INTERVAL '1 day'
	`

	if dateFrom != "" && dateTo != "" {
		where = `
			WHERE o.created_at >= $1::date
			  AND o.created_at < ($2::date + INTERVAL '1 day')
		`
		args = []any{dateFrom, dateTo, limit}
	}

	query := `
		SELECT
			p.id,
			p.name,
			COALESCE(SUM(oi.quantity), 0) AS quantity_sold,
			COALESCE(SUM(oi.line_total)::double precision, 0) AS revenue
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN products p ON p.id = oi.product_id
	` + where + `
		GROUP BY p.id, p.name
		ORDER BY quantity_sold DESC, revenue DESC
		LIMIT $` + fmt.Sprintf("%d", len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("top products query: %w", err)
	}
	defer rows.Close()

	var result []domain.TopProduct
	for rows.Next() {
		var row domain.TopProduct
		if err := rows.Scan(&row.ProductID, &row.Name, &row.QuantitySold, &row.Revenue); err != nil {
			return nil, fmt.Errorf("scan top products: %w", err)
		}
		result = append(result, row)
	}

	return result, rows.Err()
}

func (r *ReportRepository) GetLowStockProducts(ctx context.Context, threshold int) ([]domain.Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, sku, barcode, price::double precision, stock, is_active, created_at, updated_at
		FROM products
		WHERE stock <= $1
		ORDER BY stock ASC, name ASC
	`, threshold)
	if err != nil {
		return nil, fmt.Errorf("low stock products: %w", err)
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
