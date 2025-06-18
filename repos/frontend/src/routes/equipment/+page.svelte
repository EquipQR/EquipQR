<script lang="ts">
  import { page } from '$app/stores';
  import { onDestroy } from 'svelte';

  let equipmentId: string | null = null;

  const unsubscribe = page.subscribe(($page) => {
    const scanned = $page.url.searchParams.get('scanned');
    equipmentId = scanned ?? null;
  });

  onDestroy(() => {
    unsubscribe();
  });
</script>

<div>
  <h1 class="text-2xl font-bold mb-4">Equipment Page</h1>

  {#if equipmentId}
    <p>Loaded equipment ID: {equipmentId}</p>
  {:else}
    <p>No scanned equipment ID found.</p>
  {/if}
</div>
