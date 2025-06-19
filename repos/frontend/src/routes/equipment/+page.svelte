<script lang="ts">
  import { onMount, tick } from "svelte";
  import { getEquipmentById } from "$lib/api/equipment";
  import { ArrowLeft, FileText } from "lucide-svelte";
  import { formatKey } from "$lib/utils";
  import type { Equipment } from "$lib/types/equipment";

  import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
  } from "$lib/components/ui/card";
  import { Separator } from "$lib/components/ui/separator/index.js";
  import { Button } from "$lib/components/ui/button";
  import { Badge } from "$lib/components/ui/badge";
  import { Skeleton } from "$lib/components/ui/skeleton";

  let equipmentId: string | null = null;
  let equipment: Equipment | null = null;
  let error: string | null = null;
  let imageLoaded = false;
  let imageRef: HTMLImageElement | null = null;

  onMount(async () => {
    const urlParams = new URLSearchParams(window.location.search);
    equipmentId = urlParams.get("scanned") ?? null;

    if (equipmentId) {
      try {
        equipment = await getEquipmentById(equipmentId);
        error = null;
      } catch (e) {
        error = "Failed to load equipment.";
      }
    }

    await tick();
    if (imageRef?.complete) {
      imageLoaded = true;
    }
  });
</script>

<div class="relative w-full">
  <div class="fixed top-0 left-0 w-full h-48 z-10">
    {#if !imageLoaded}
      <Skeleton class="w-full h-full rounded-none animate-pulse" />
    {/if}

    <img
      bind:this={imageRef}
      src="https://csdieselgenerators.com/wp-content/uploads/2017/03/Blog02-767x377.jpg"
      alt="Equipment"
      class={`w-full h-full object-cover transition-opacity duration-300 ${imageLoaded ? "" : "hidden"}`}
      on:load={() => (imageLoaded = true)}
    />

    <div class="absolute top-4 right-4 flex space-x-3">
      <Button
        variant="ghost"
        size="icon"
        class="bg-black/60 text-white hover:bg-black/80 transition transform active:scale-90"
      >
        <FileText class="w-5 h-5" />
      </Button>
      <Button
        variant="ghost"
        size="icon"
        class="bg-black/60 text-white hover:bg-black/80 transition transform active:scale-90"
        onclick={() => (window.location.href = "/")}
      >
        <ArrowLeft class="w-5 h-5" />
      </Button>
    </div>
  </div>

  <div class="pt-56 max-w-3xl mx-auto px-4">
    <Card class="border border-border bg-background shadow-md">
      {#if error}
        <CardContent class="p-6 text-red-500">{error}</CardContent>
      {:else if equipment}
        <CardHeader class="px-6 pb-0">
          <CardTitle class="text-2xl font-bold">
            {equipment.type}
          </CardTitle>
        </CardHeader>
        <Separator />
        <CardContent class="px-6 pb-6 space-y-6">
          <div>
            <span class="text-sm font-medium text-muted-foreground">Status</span
            >
            <div class="mt-1">
              <Badge
                class={equipment.status === "in service"
                  ? "bg-green-600 text-white"
                  : equipment.status === "not in service"
                    ? "bg-red-600 text-white"
                    : "bg-yellow-500 text-black"}
              >
                {equipment.status}
              </Badge>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-muted-foreground">Type</span>
              <p class="font-medium">{equipment.type}</p>
            </div>
            <div>
              <span class="text-muted-foreground">Location</span>
              <p class="font-medium">{equipment.location}</p>
            </div>
            {#if equipment.moreFields}
              {#each Object.entries(equipment.moreFields) as [key, value]}
                <div>
                  <span class="text-muted-foreground">{formatKey(key)}</span>
                  <p class="font-medium">{value}</p>
                </div>
              {/each}
            {/if}
          </div>
        </CardContent>
      {:else if equipmentId}
        <CardContent class="space-y-4 px-6 py-8">
          <Skeleton class="h-5 w-2/5" />
          <Skeleton class="h-4 w-3/4" />
          <Skeleton class="h-4 w-1/3" />
          <Skeleton class="h-4 w-2/4" />
        </CardContent>
      {:else}
        <CardContent class="px-6 py-8 text-muted-foreground text-sm">
          No scanned equipment ID found.
        </CardContent>
      {/if}
    </Card>
  </div>
</div>
