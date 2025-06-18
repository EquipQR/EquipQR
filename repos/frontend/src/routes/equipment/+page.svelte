<script lang="ts">
  import { onMount } from "svelte";
  import { getEquipmentById } from "$lib/api/equipment";
  import { ArrowLeft, FileText } from "lucide-svelte";
  import type { Equipment } from "$lib/types/equipment";

  let equipmentId: string | null = null;
  let equipment: Equipment | null = null;
  let error: string | null = null;

  onMount(async () => {
    const urlParams = new URLSearchParams(window.location.search);
    equipmentId = urlParams.get("scanned") ?? null;

    if (equipmentId) {
      equipment = await getEquipmentById(equipmentId);
      error = null;
    }
  });
</script>

<div class="max-w-4xl mx-auto p-4 space-y-6">
  <div class="fixed top-0 left-0 w-screen h-48 z-10">
    <img
      src="https://csdieselgenerators.com/wp-content/uploads/2017/03/Blog02-767x377.jpg"
      alt="Equipment"
      class="w-full h-full object-cover"
    />
    <div class="absolute top-4 right-4 flex space-x-3">
      <button
        class="bg-black/70 text-white p-3 rounded-full backdrop-blur hover:bg-black/90 transition transform active:scale-95"
      >
        <FileText size={24} />
      </button>
      <button
        class="bg-black/70 text-white p-3 rounded-full backdrop-blur hover:bg-black/90 transition transform active:scale-95"
        on:click={() => (window.location.href = "/")}
      >
        <ArrowLeft size={24} />
      </button>
    </div>
  </div>

  <div class="pt-52">
    <div class="bg-white p-6 shadow-lg rounded-lg">
      <h1 class="text-3xl font-bold text-gray-800 mb-4">Equipment Details</h1>
      {#if error}
        <p class="text-red-500">{error}</p>
      {:else if equipment}
        <p><strong class="text-gray-700">ID:</strong> {equipment.id}</p>
        <p><strong class="text-gray-700">Status:</strong> {equipment.status}</p>
        <p><strong class="text-gray-700">Type:</strong> {equipment.type}</p>
        <p>
          <strong class="text-gray-700">Location:</strong>
          {equipment.location}
        </p>

        {#if equipment.moreFields}
          <ul class="mt-4 space-y-2">
            {#each Object.entries(equipment.moreFields) as [key, value]}
              <li class="text-gray-600">
                <strong class="font-semibold">{key}:</strong>
                {value}
              </li>
            {/each}
          </ul>
        {/if}
      {:else if equipmentId}
        <p class="text-gray-500">Loading equipment...</p>
      {:else}
        <p class="text-gray-500">No scanned equipment ID found.</p>
      {/if}
    </div>
  </div>
</div>
