import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import fs from 'node:fs';
import path from 'node:path';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(process.cwd(), '../../'), '');

  return {
    plugins: [react()],
    server: {
      host: env.VITE_HOST || 'localhost',  // Default to 'localhost' if not defined
      port: Number(env.VITE_PORT) || 3000, // Default to 3000 if not defined
      https: {
        key: fs.readFileSync(path.resolve(__dirname, env.VITE_SSL_KEY_PATH)),
        cert: fs.readFileSync(path.resolve(__dirname, env.VITE_SSL_CERT_PATH)),
      },
    },
    optimizeDeps: {
      exclude: ['lucide-react'],
    },
  };
});
