package service

import (
	"context"
	"fmt"

	"pos-system/backend/internal/apperrors"
	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService struct {
	db       *pgxpool.Pool
	products *repository.ProductRepository
	orders   *repository.OrderRepository
}

func NewOrderService(db *pgxpool.Pool, products *repository.ProductRepository, orders *repository.OrderRepository) *OrderService {
	return &OrderService{
		db:       db,
		products: products,
		orders:   orders,
	}
}

func (s *OrderService) Create(ctx context.Context, input domain.CreateOrderInput) (domain.OrderReceipt, error) {
	if len(input.Items) == 0 {
		return domain.OrderReceipt{}, fmt.Errorf("%w: order must include at least one item", apperrors.ErrValidation)
	}

	mergedItems := make(map[int64]int)
	var productIDs []int64
	for _, item := range input.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return domain.OrderReceipt{}, fmt.Errorf("%w: invalid order item", apperrors.ErrValidation)
		}
		if _, exists := mergedItems[item.ProductID]; !exists {
			productIDs = append(productIDs, item.ProductID)
		}
		mergedItems[item.ProductID] += item.Quantity
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.OrderReceipt{}, err
	}
	defer tx.Rollback(ctx)

	products, err := s.products.GetByIDsForUpdate(ctx, tx, productIDs)
	if err != nil {
		return domain.OrderReceipt{}, err
	}
	if len(products) != len(productIDs) {
		return domain.OrderReceipt{}, fmt.Errorf("%w: one or more products do not exist", apperrors.ErrNotFound)
	}

	productsByID := make(map[int64]domain.Product, len(products))
	for _, product := range products {
		productsByID[product.ID] = product
	}

	var total float64
	orderItems := make([]domain.OrderItem, 0, len(productIDs))
	for _, productID := range productIDs {
		product := productsByID[productID]
		quantity := mergedItems[productID]

		if !product.IsActive {
			return domain.OrderReceipt{}, fmt.Errorf("%w: product %d is inactive", apperrors.ErrValidation, productID)
		}
		if product.Stock < quantity {
			return domain.OrderReceipt{}, fmt.Errorf("%w: product %s has only %d units left", apperrors.ErrInsufficientStock, product.Name, product.Stock)
		}

		lineTotal := product.Price * float64(quantity)
		total += lineTotal
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			Quantity:  quantity,
			UnitPrice: product.Price,
			LineTotal: lineTotal,
		})
	}

	order, err := s.orders.CreateOrderTx(ctx, tx, total)
	if err != nil {
		return domain.OrderReceipt{}, err
	}

	if err := s.orders.InsertOrderItemsTx(ctx, tx, order.ID, orderItems); err != nil {
		return domain.OrderReceipt{}, err
	}

	for _, item := range orderItems {
		if err := s.orders.UpdateProductStockTx(ctx, tx, item.ProductID, -item.Quantity); err != nil {
			return domain.OrderReceipt{}, err
		}
		if err := s.orders.InsertInventoryMovementTx(ctx, tx, domain.InventoryMovementInput{
			ProductID: item.ProductID,
			ChangeQty: -item.Quantity,
			Reason:    "sale",
		}); err != nil {
			return domain.OrderReceipt{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.OrderReceipt{}, err
	}

	order.TotalAmount = total
	return domain.OrderReceipt{Order: order, Items: orderItems}, nil
}

func (s *OrderService) Today(ctx context.Context) (domain.TodayOrdersResponse, error) {
	return s.orders.ListToday(ctx)
}
