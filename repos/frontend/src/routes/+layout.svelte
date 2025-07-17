<script lang="ts">
  import "../app.css";
  let { children } = $props();
  // Detects if the device is a mobile device or not:
  // const isMobile = /Mobi|Android|iPhone|iPad/i.test(navigator.userAgent);
  // console.log(isMobile);
</script>

<!-- keeps toasts alive across page navigation -->
{@render children?.()}
