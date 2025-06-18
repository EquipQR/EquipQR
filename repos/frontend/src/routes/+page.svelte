<script lang="ts">
  import QRScanner from "$lib/components/QRScanner.svelte";

  let loading = false;
  let error: string | undefined;
  let equipmentId: string | null = null;

  const handleScan = async (result: string) => {
    loading = true;
    error = undefined;

    try {
      // fake fetch for demo
      await new Promise((r) => setTimeout(r, 1000));
      equipmentId = result;
      // Optional: redirect after scan
      // goto(`/equipment/${result}`);
    } catch (err) {
      error = "Failed to fetch equipment data.";
      console.error(err);
    } finally {
      loading = false;
    }
  };
</script>

{#if !equipmentId}
  <QRScanner {loading} {error} on:scan={(e) => handleScan(e.detail)} />
{:else}
  <div class="p-6 text-white">
    <h1 class="text-2xl font-bold mb-4">Equipment Scanned</h1>
    <p>
      Scanned ID: <span class="font-mono text-green-400">{equipmentId}</span>
    </p>
    <button
      class="mt-6 px-4 py-2 bg-blue-600 rounded hover:bg-blue-700 transition"
      on:click={() => (equipmentId = null)}
    >
      Scan Again
    </button>
  </div>
{/if}
