import type { Product, ProductListResponse } from '~/types'

export function useProducts() {
  const items = useState<Product[]>('products-cache', () => [])
  const loading = useState<boolean>('products-loading', () => false)
  const loaded = useState<boolean>('products-loaded', () => false)

  async function fetchProducts(query = '') {
    loading.value = true
    try {
      const response = await useApi<ProductListResponse>(`/products?q=${encodeURIComponent(query)}&limit=200`)
      items.value = response.items
      loaded.value = true
      return response.items
    }
    finally {
      loading.value = false
    }
  }

  async function ensureProducts() {
    if (!loaded.value) {
      await fetchProducts()
    }
    return items.value
  }

  async function searchProducts(query: string) {
    if (!query.trim()) {
      return items.value
    }

    const normalized = query.toLowerCase()
    return items.value.filter(product =>
      product.name.toLowerCase().includes(normalized)
      || (product.sku || '').toLowerCase().includes(normalized)
      || (product.barcode || '').toLowerCase().includes(normalized),
    )
  }

  return {
    items,
    loading,
    loaded,
    ensureProducts,
    fetchProducts,
    searchProducts,
  }
}
