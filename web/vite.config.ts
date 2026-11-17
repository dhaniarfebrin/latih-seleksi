import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

// Dev server stays on localhost (never expose it to the LAN); the nginx
// container serves the production build.
export default defineConfig({
  plugins: [react()],
  resolve: { alias: { '@': '/src' } },
  server: {
    host: 'localhost',
    // DEV_API_URL points the proxy at another API, e.g. when 8080 belongs to a different project.
    proxy: { '/v1': process.env.DEV_API_URL ?? 'http://localhost:8080' },
  },
  test: { environment: 'node' },
})
