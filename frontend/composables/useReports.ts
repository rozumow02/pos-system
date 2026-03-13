import type { Product, TopProduct } from "~/types"

export function useReports() {
  const topProducts = ref<TopProduct[]>([])
  const lowStock = ref<Product[]>([])
  const loading = ref(false)

  async function refreshTopProducts(dateFrom = "", dateTo = "", limit = 10) {
    loading.value = true
    try {
      const params = new URLSearchParams()
      if (dateFrom) params.set("date_from", dateFrom)
      if (dateTo) params.set("date_to", dateTo)
      params.set("limit", String(limit))
      const response = await useApi<{ items: TopProduct[] }>(`/reports/top-products?${params.toString()}`)
      topProducts.value = response.items
    } finally {
      loading.value = false
    }
  }

  async function refreshLowStock(threshold = 0) {
    const suffix = threshold > 0 ? `?threshold=${threshold}` : ""
    const response = await useApi<{ items: Product[] }>(`/reports/low-stock${suffix}`)
    lowStock.value = response.items
  }

  return {
    topProducts,
    lowStock,
    loading,
    refreshTopProducts,
    refreshLowStock
  }
}
