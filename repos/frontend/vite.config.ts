import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, loadEnv } from "vite";
import fs from "node:fs";
import path from "node:path";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(process.cwd(), "../../"), "");

  return {
    plugins: [tailwindcss(), sveltekit()],
    server: {
      host: env.VITE_HOST || "localhost",
      port: Number(env.VITE_PORT) || 3000,
      https: {
        key: fs.readFileSync(path.resolve(__dirname, env.VITE_SSL_KEY_PATH)),
        cert: fs.readFileSync(path.resolve(__dirname, env.VITE_SSL_CERT_PATH)),
      },
      proxy: {
        "/api": {
          target: env.VITE_PROXY_TARGET,
          changeOrigin: true,
          secure: false,
          rewrite: (path) => path.replace(/^\/api/, ""),
        },
      },
    },
    optimizeDeps: {
      exclude: ["lucide-svelte"],
    },
  };
});
