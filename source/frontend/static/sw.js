/**
 * Service Worker for EquipQR
 *
 * This file handles:
 * - Pre-caching core static assets
 * - Cleanup of outdated caches on activation
 * - Runtime caching of all GET requests
 * - Offline fallback for navigations
 *
 * Compiles to `static/sw.js` and should be registered in the client.
 */

const SW_VERSION = "dev"; // will be replaced at build time
console.log(`🚀 EquipQR Service Worker version: ${SW_VERSION}`);

const CACHE_NAME = "equipqr-cache-v1";

const PRECACHE_URLS = [
  "/", // for SPA routing
  "/index.html",
  "/favicon.ico",
  "/manifest.webmanifest",
];

// Precache critical static assets
self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(PRECACHE_URLS);
    })
  );
  self.skipWaiting();
});

// Clean up old caches
self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(
          keys
            .filter((key) => key !== CACHE_NAME)
            .map((key) => caches.delete(key))
        )
      )
  );
  self.clients.claim();
});

// Runtime cache-first strategy
self.addEventListener("fetch", (event) => {
  const url = new URL(event.request.url);

  if (
    url.protocol.startsWith("chrome-extension:") ||
    url.protocol === "about:" ||
    event.request.method !== "GET"
  ) {
    return; // 🛑 Don't handle unsupported requests
  }

  event.respondWith(
    (async () => {
      const cache = await caches.open(CACHE_NAME);
      const cached = await cache.match(event.request);

      // Use network-first for immutable assets (like JS chunks)
      if (event.request.url.includes("/_app/immutable/")) {
        try {
          const networkResponse = await fetch(event.request);
          cache.put(event.request, networkResponse.clone());
          return networkResponse;
        } catch {
          if (cached) return cached;
          throw new Error("Dynamic module failed to load and is not cached.");
        }
      }

      // Normal cache-first for everything else
      if (cached) return cached;

      try {
        const networkResponse = await fetch(event.request);
        if (event.request.url.startsWith(self.location.origin)) {
          cache.put(event.request, networkResponse.clone());
        }
        return networkResponse;
      } catch {
        if (event.request.mode === "navigate") {
          const fallback = await cache.match("/");
          if (fallback) return fallback;
        }

        return new Response("Offline and not cached", {
          status: 503,
          statusText: "Service Unavailable",
        });
      }
    })()
  );
});
