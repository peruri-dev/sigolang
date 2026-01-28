import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: parseInt(process.env.VITE_PORT || "5173"), // Use PORT if available, otherwise default to 5173
    headers: {
      // Security Headers for development
      "X-Content-Type-Options": "nosniff",
      "X-XSS-Protection": "1; mode=block",
      "X-Frame-Options": "SAMEORIGIN",
      "Strict-Transport-Security":
        "max-age=31536000; includeSubDomains; preload",
      "Content-Security-Policy": "frame-ancestors 'self';",
      "Referrer-Policy": "strict-origin-when-cross-origin",
      "Permissions-Policy": "fullscreen=(), camera=(), microphone=()",
      "Cross-Origin-Resource-Policy": "same-site",
      "Cross-Origin-Opener-Policy": "same-origin",
      "Cross-Origin-Embedder-Policy": "require-corp",
      "X-DNS-Prefetch-Control": "off",
    },
  },
  build: {
    outDir: 'frontend/dist',
    // generates .vite/manifest.json in outDir
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      input: "/src/main.tsx",
    },
  },
})
