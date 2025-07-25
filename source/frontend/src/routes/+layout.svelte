<script lang="ts">
  import "../app.css";
  let { children } = $props();
  // Detects if the device is a mobile device or not:
  // const isMobile = /Mobi|Android|iPhone|iPad/i.test(navigator.userAgent);
  // console.log(isMobile);

  if ("serviceWorker" in navigator) {
    navigator.serviceWorker
      .register("/sw.js")
      .then((r) => {
        console.log("✅ Service Worker registered:", r.scope);
      })
      .catch((err) => {
        console.error("❌ Service Worker registration failed:", err);
      });
  }
</script>

<!-- keeps toasts alive across page navigation -->
{@render children?.()}
