// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: import.meta.env.DEV },
  runtimeConfig: {
    public: {
      baseApi: '',
    },
  },
  vite: {
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: "@use '~/styles/_variables.scss' as *;",
        },
      },
    },
  },
  css: ['~/styles/main.scss'],
  modules: [
    '@nuxt/eslint',
    '@pinia/nuxt',
    '@nuxtjs/color-mode',
    '@nuxt/icon'
  ]
})