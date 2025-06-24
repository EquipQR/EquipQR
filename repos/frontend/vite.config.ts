import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, loadEnv, Plugin } from "vite";
import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";

function postBuildHashPlugin(): Plugin {
  return {
    name: "post-build-hash",
    closeBundle: async () => {
      const srcDir = path.resolve(process.cwd(), "src");
      const hash = crypto.createHash("sha256");

      const walk = (dir: string): void => {
        for (const entry of fs.readdirSync(dir)) {
          const fullPath = path.join(dir, entry);
          const stat = fs.statSync(fullPath);
          if (stat.isDirectory()) {
            walk(fullPath);
          } else {
            const contentBuffer = fs.readFileSync(fullPath);
            const contentArray = new Uint8Array(contentBuffer);
            hash.update(contentArray);
          }
        }
      };

      walk(srcDir);
      const finalHash = hash.digest("hex");
      const outputPath = path.resolve(
        process.cwd(),
        "../backend/web/.frontend_build_hash"
      );

      fs.writeFileSync(outputPath, finalHash + "\n");
      console.log(`✅ Frontend build hash written to ${outputPath}`);
    },
  };
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(process.cwd(), "../../"), "");

  return {
    plugins: [tailwindcss(), sveltekit(), postBuildHashPlugin()],
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
        },
      },
    },
    optimizeDeps: {
      exclude: ["lucide-svelte"],
    },
  };
});
