export default defineNuxtConfig({
  modules: ['@nuxt/eslint'],
  devtools: { enabled: false },
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    apiServerBase: import.meta.env.NUXT_API_SERVER_BASE || 'http://backend:8080/api',
    public: {
      apiBase: import.meta.env.NUXT_PUBLIC_API_BASE || '/api',
    },
  },
  compatibilityDate: '2025-12-01',
  eslint: {
    config: {
      import: false,
      stylistic: false,
      standalone: false,
    },
    checker: false,
  },
})
