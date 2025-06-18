<script lang="ts">
  import QRScanner from "$lib/components/QRScanner.svelte";

  let loading = false;
  let error: string | undefined;

</script>

<QRScanner {loading} {error} />
