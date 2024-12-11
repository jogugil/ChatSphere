// vite.config.ts
export default defineConfig({
    server: {
      proxy: {
        '/api': 'http://localhost:4000',  // Cambia el puerto según tu servidor Go o Node.js
      },
    },
  });
  