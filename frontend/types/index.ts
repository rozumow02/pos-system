export interface Product {
  id: number
  name: string
  sku?: string | null
  barcode?: string | null
  price: number
  stock: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface ProductListResponse {
  items: Product[]
  page: number
  limit: number
  total: number
}

export interface TopProduct {
  product_id: number
  name: string
  quantity_sold: number
  revenue: number
}

export interface DashboardResponse {
  revenue_today: number
  items_sold_today: number
  orders_today: number
  top_products: TopProduct[]
  low_stock: Product[]
  low_stock_count: number
}

export interface Order {
  id: number
  total_amount: number
  created_at: string
}

export interface OrderReceiptItem {
  product_id: number
  name: string
  quantity: number
  unit_price: number
  line_total: number
}

export interface OrderReceipt {
  order: Order
  items: OrderReceiptItem[]
}
