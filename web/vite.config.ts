import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const devProxyServer = "http://localhost:8080"

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: "0.0.0.0",
    port: 8888,
    proxy: {
      "^/api": {
        target: devProxyServer,
        xfwd: true,
      },
      "^/api.v1.*": {
        target: devProxyServer,
        xfwd: true,
      },
      "^/file": {
        target: devProxyServer,
        xfwd: true,
      },
    },
  },
})
