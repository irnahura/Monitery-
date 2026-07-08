import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/auth": "http://backend:8080",
      "/monitors": "http://backend:8080",
      "/apikeys": "http://backend:8080"
    }
  }
});
