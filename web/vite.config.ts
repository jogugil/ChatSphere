import { defineConfig } from "vite";

// vite.config.ts
export default defineConfig({
    root: './',
    publicDir: 'public', 
    server: {
      host: '0.0.0.0',
      port: 5173,
      proxy: {
        '/api': 'http://localhost:8081',  // Hay que Cambia el puerto según modifique el servidor GoChat 
      },
    },
  });
  