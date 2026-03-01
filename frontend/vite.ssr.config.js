import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  build: {
    outDir: './server',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        entryFileNames: 'blog-front.js'
      }
    }
  }
})
