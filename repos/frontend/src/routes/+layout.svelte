<script lang="ts">
  import "../app.css";
  let { children } = $props();
</script>

<!-- keeps toasts alive across page navigation -->
{@render children?.()}
