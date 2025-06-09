// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss',"@pinia/nuxt",'pinia-plugin-persistedstate/nuxt','@nuxt/icon'],
  plugins: [
    { src: '~/plugins/matomo.client.js', mode: 'client' },
  ],
  runtimeConfig: {
    public: {
      BACKEND_URL: process.env.NUXT_PUBLIC_BACKEND_URL || 'http://localhost:6050',
      BASE_URL: process.env.NUXT_PUBLIC_BASE_URL || 'http://localhost:3000',
      matomo_host: process.env.MATOMO_HOST || "https://matomo.srv630927.hstgr.cloud/",
      matomo_site_id: process.env.MATOMO_SITE_ID || 8,
    },
  },
})
