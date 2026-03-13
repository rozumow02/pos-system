package repository

import (
	"context"
	"fmt"

	"pos-system/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) CreateOrderTx(ctx context.Context, tx pgx.Tx, total float64) (domain.Order, error) {
	var order domain.Order
	err := tx.QueryRow(ctx, `
		INSERT INTO orders (total_amount)
		VALUES ($1)
		RETURNING id, total_amount::double precision, created_at
	`, total).Scan(&order.ID, &order.TotalAmount, &order.CreatedAt)
	if err != nil {
		return domain.Order{}, fmt.Errorf("create order: %w", err)
	}
	return order, nil
}

func (r *OrderRepository) InsertOrderItemsTx(ctx context.Context, tx pgx.Tx, orderID int64, items []domain.OrderItem) error {
	for _, item := range items {
		if _, err := tx.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price, line_total)
			VALUES ($1, $2, $3, $4, $5)
		`, orderID, item.ProductID, item.Quantity, item.UnitPrice, item.LineTotal); err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	return nil
}

func (r *OrderRepository) UpdateProductStockTx(ctx context.Context, tx pgx.Tx, productID int64, delta int) error {
	tag, err := tx.Exec(ctx, `
		UPDATE products
		SET stock = stock + $2, updated_at = NOW()
		WHERE id = $1
	`, productID, delta)
	if err != nil {
		return fmt.Errorf("update product stock: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update product stock: no rows affected")
	}
	return nil
}

func (r *OrderRepository) InsertInventoryMovementTx(ctx context.Context, tx pgx.Tx, movement domain.InventoryMovementInput) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO inventory_movements (product_id, change_qty, reason)
		VALUES ($1, $2, $3)
	`, movement.ProductID, movement.ChangeQty, movement.Reason)
	if err != nil {
		return fmt.Errorf("insert inventory movement: %w", err)
	}
	return nil
}

func (r *OrderRepository) ListToday(ctx context.Context) (domain.TodayOrdersResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, total_amount::double precision, created_at
		FROM orders
		WHERE created_at >= date_trunc('day', NOW())
		ORDER BY created_at DESC
	`)
	if err != nil {
		return domain.TodayOrdersResponse{}, fmt.Errorf("list today's orders: %w", err)
	}
	defer rows.Close()

	response := domain.TodayOrdersResponse{
		Orders: make([]domain.Order, 0),
	}

	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.TotalAmount, &order.CreatedAt); err != nil {
			return domain.TodayOrdersResponse{}, fmt.Errorf("scan today's orders: %w", err)
		}
		response.Orders = append(response.Orders, order)
		response.TotalRevenue += order.TotalAmount
	}
	if err := rows.Err(); err != nil {
		return domain.TodayOrdersResponse{}, err
	}

	response.TotalOrders = len(response.Orders)
	return response, nil
}
