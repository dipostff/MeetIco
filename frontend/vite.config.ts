import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          http:   ['axios'],
        }
      }
    },
    minify: 'terser',
    chunkSizeWarningLimit: 500,
  },
  resolve: {
    alias: { '@': '/src' }
  }
})
