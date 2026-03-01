import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  base: '/build/',
  build: {
    outDir: './build',
    emptyOutDir: true,
    rollupOptions: {
      input: '/home/kamiy2743/workspace/blog/frontend/src/app.ts',
      output: {
        entryFileNames: 'assets/app.js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name][extname]'
      }
    }
  }
})
