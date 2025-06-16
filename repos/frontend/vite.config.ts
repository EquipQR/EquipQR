import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import fs from 'node:fs';
import path from 'node:path';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(process.cwd(), '../../'), '');

  return {
    plugins: [react()],
    server: {
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
