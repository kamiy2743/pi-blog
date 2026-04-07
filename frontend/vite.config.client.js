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
      input: {
        app: './src/app.ts',
        error: './src/error/error.html'
      },
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name][extname]'
      }
    }
  }
})
