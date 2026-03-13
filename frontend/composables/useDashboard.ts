import type { DashboardResponse } from "~/types"

export function useDashboard() {
  const data = useState<DashboardResponse | null>("dashboard-data", () => null)
  const loading = useState<boolean>("dashboard-loading", () => false)

  async function refresh() {
    loading.value = true
    try {
      data.value = await useApi<DashboardResponse>("/reports/dashboard")
    } finally {
      loading.value = false
    }
  }

  return {
    data,
    loading,
    refresh
  }
}
