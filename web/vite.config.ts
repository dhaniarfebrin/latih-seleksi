import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

// Dev server stays on localhost (never expose it to the LAN); the nginx
// container serves the production build.
export default defineConfig({
  plugins: [react()],
  server: {
    host: 'localhost',
    proxy: { '/v1': 'http://localhost:8080' },
  },
  test: { environment: 'node' },
})
