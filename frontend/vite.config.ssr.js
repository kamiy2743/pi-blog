import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  publicDir: false,
  build: {
    outDir: './dist/ssr/',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        entryFileNames: 'ssr.js'
      }
    }
  }
})
