process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, loadEnv } from "vite";
import type { Plugin } from "vite";
import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";

export function postBuildHashPlugin(): Plugin {
  return {
    name: "post-build-hash",
    closeBundle: async (): Promise<void> => {
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

      const outputDir = path.resolve(process.cwd(), "../backend/web");
      const outputPath = path.join(outputDir, ".frontend_build_hash");

      fs.mkdirSync(outputDir, { recursive: true });
      fs.writeFileSync(outputPath, `${finalHash}\n`);

      console.log(`✅ Frontend build hash written to ${outputPath}`);
    },
  };
}

export function compileServiceWorkerPlugin(): Plugin {
  return {
    name: "compile-sw",
    buildStart: () => {
      const swInput = path.resolve(process.cwd(), "src/sw.ts");
      const swOutput = path.resolve(process.cwd(), "static/sw.js");

      const tsc = `npx tsc ${swInput} --target esnext --module esnext --lib webworker --outFile ${swOutput}`;

      try {
        console.log("🔧 Compiling sw.ts to sw.js...");
        require("node:child_process").execSync(tsc, { stdio: "inherit" });
        console.log("✅ sw.ts compiled to static/sw.js");
      } catch (err) {
        console.error("❌ Failed to compile sw.ts:", err);
      }
    },
  };
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(process.cwd(), "../../"), "");

  return {
    plugins: [
      tailwindcss(),
      sveltekit(),
      compileServiceWorkerPlugin(),
      postBuildHashPlugin(),
    ],
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
