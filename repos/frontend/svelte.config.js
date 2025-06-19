import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";
import path from "node:path";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      fallback: "index.html",
      strict: false,
      pages: "../backend/web",      // ✅ HTML & assets
      assets: "../backend/web",     // ✅ JS chunks, CSS, etc.
    }),
    alias: {
      $lib: path.resolve("./src/lib"),
    },
    prerender: {
      entries: [],
    },
  },
};

export default config;
