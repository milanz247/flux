import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'

// Flux frontend build. Dev: the Go server injects script tags pointing at
// this dev server (VITE_DEV_SERVER). Prod: `flux build` emits a manifest
// into public/build that the Go server reads.
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  // The Go server serves public/ itself; Vite only owns public/build.
  publicDir: false,
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./resources/js', import.meta.url)),
    },
  },
  build: {
    manifest: true,
    outDir: 'public/build',
    emptyOutDir: true,
    rollupOptions: {
      input: 'resources/js/app.ts',
    },
  },
  server: {
    port: 5173,
    strictPort: true,
    origin: 'http://localhost:5173',
    cors: true,
  },
})
