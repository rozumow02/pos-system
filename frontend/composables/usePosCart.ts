import type { Product } from "~/types"

export interface CartItem {
  product_id: number
  name: string
  price: number
  stock: number
  quantity: number
}

export function usePosCart() {
  const cart = useState<CartItem[]>("pos-cart", () => [])

  function addProduct(product: Product) {
    const existing = cart.value.find((item) => item.product_id === product.id)
    if (existing) {
      if (existing.quantity < product.stock) {
        existing.quantity += 1
      }
      return
    }

    cart.value.push({
      product_id: product.id,
      name: product.name,
      price: product.price,
      stock: product.stock,
      quantity: 1
    })
  }

  function updateQuantity(productId: number, quantity: number) {
    const item = cart.value.find((entry) => entry.product_id === productId)
    if (!item) return
    if (quantity <= 0) {
      removeProduct(productId)
      return
    }
    item.quantity = Math.min(quantity, item.stock)
  }

  function removeProduct(productId: number) {
    cart.value = cart.value.filter((item) => item.product_id !== productId)
  }

  function clearCart() {
    cart.value = []
  }

  const total = computed(() => cart.value.reduce((sum, item) => sum + item.price * item.quantity, 0))
  const totalItems = computed(() => cart.value.reduce((sum, item) => sum + item.quantity, 0))

  return {
    cart,
    total,
    totalItems,
    addProduct,
    updateQuantity,
    removeProduct,
    clearCart
  }
}
