import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, import.meta.dirname);

  return {
    build: { outDir: "internal/vite/build" },
    clearScreen: false,
    server: {
      proxy: { "/api": { target: env.VITE_API_URL, changeOrigin: true } },
    },
  };
});
