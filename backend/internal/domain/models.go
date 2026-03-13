package domain

import "time"

type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	SKU       *string   `json:"sku,omitempty"`
	Barcode   *string   `json:"barcode,omitempty"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductListFilter struct {
	Query  string
	Page   int
	Limit  int
	Active *bool
}

type ProductListResult struct {
	Items []Product `json:"items"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
	Total int       `json:"total"`
}

type CreateProductInput struct {
	Name     string  `json:"name"`
	SKU      *string `json:"sku"`
	Barcode  *string `json:"barcode"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	IsActive *bool   `json:"is_active"`
}

type UpdateProductInput struct {
	Name     *string  `json:"name"`
	SKU      *string  `json:"sku"`
	Barcode  *string  `json:"barcode"`
	Price    *float64 `json:"price"`
	Stock    *int     `json:"stock"`
	IsActive *bool    `json:"is_active"`
}

type Order struct {
	ID          int64     `json:"id"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	LineTotal float64 `json:"line_total"`
}

type CreateOrderInput struct {
	Items []CreateOrderItemInput `json:"items"`
}

type CreateOrderItemInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type OrderReceipt struct {
	Order Order       `json:"order"`
	Items []OrderItem `json:"items"`
}

type TodayOrdersResponse struct {
	Orders       []Order `json:"orders"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalOrders  int     `json:"total_orders"`
}

type InventoryMovementInput struct {
	ProductID int64
	ChangeQty int
	Reason    string
}

type TopProduct struct {
	ProductID    int64   `json:"product_id"`
	Name         string  `json:"name"`
	QuantitySold int     `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
}

type DashboardMetrics struct {
	RevenueToday   float64      `json:"revenue_today"`
	ItemsSoldToday int          `json:"items_sold_today"`
	OrdersToday    int          `json:"orders_today"`
	TopProducts    []TopProduct `json:"top_products"`
	LowStock       []Product    `json:"low_stock"`
	LowStockCount  int          `json:"low_stock_count"`
}
