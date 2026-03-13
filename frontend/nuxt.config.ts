export default defineNuxtConfig({
  compatibilityDate: "2025-12-01",
  devtools: { enabled: false },
  modules: ["@nuxt/eslint"],
  css: ["~/assets/css/main.css"],
  eslint: {
    config: {
      stylistic: true
    },
    checker: false
  },
  runtimeConfig: {
    apiServerBase: process.env.NUXT_API_SERVER_BASE || "http://backend:8080/api",
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || "/api"
    }
  },
  typescript: {
    strict: true,
    typeCheck: true
  }
})
