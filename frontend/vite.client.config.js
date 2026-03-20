import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    svelte(),
    tailwindcss()
  ],
  base: '/dist/client/',
  build: {
    outDir: './dist/client/',
    emptyOutDir: true,
    rollupOptions: {
      input: './src/app.ts',
      output: {
        entryFileNames: 'assets/app.js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name][extname]'
      }
    }
  }
})
