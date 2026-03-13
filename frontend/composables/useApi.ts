export function useApi<T>(path: string, options: Parameters<typeof $fetch<T>>[1] = {}) {
  const config = useRuntimeConfig()
  const base = import.meta.server ? config.apiServerBase : config.public.apiBase
  return $fetch<T>(`${base}${path}`, {
    ...options,
  })
}
