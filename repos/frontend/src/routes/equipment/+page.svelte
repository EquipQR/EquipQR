<script lang="ts">
  import { onMount, tick } from "svelte";
  import { ArrowLeft, FileText, PlusCircle } from "lucide-svelte";
  import { getEquipmentById, getEquipmentIssuesById } from "$lib/api/equipment";
  import { formatKey } from "$lib/utils";
  import type { Equipment } from "$lib/types/equipment";
  import type { Issue } from "$lib/types/issue";

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
  import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
  } from "$lib/components/ui/accordion";

  let equipmentId: string | null = null;
  let equipment: Equipment | null = null;
  let error: string | null = null;
  let imageLoaded = false;
  let imageRef: HTMLImageElement | null = null;

  let issues: Issue[] = [];

  export let data: { equipmentId: string | null };
  onMount(async () => {
    if (data.equipmentId) {
      try {
        equipment = await getEquipmentById(data.equipmentId);
        issues = (await getEquipmentIssuesById(data.equipmentId)) ?? [];
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

  function handleReportClick(): void {
    if (data.equipmentId) {
      window.location.href = `/issue?scanned=${encodeURIComponent(data.equipmentId)}`;
    }
  }
</script>

<div class="relative w-full">
  <div
    class="relative w-full h-[200px] sm:h-[280px] md:h-[320px] overflow-hidden"
  >
    {#if !imageLoaded}
      <Skeleton
        class="absolute inset-0 w-full h-full rounded-none animate-pulse"
      />
    {/if}

    <img
      bind:this={imageRef}
      src="https://csdieselgenerators.com/wp-content/uploads/2017/03/Blog02-767x377.jpg"
      alt="Equipment"
      class={`absolute inset-0 w-full h-full object-cover transition-opacity duration-300 ${imageLoaded ? "" : "hidden"}`}
      on:load={() => (imageLoaded = true)}
    />

    <div class="absolute top-4 right-4 z-10 flex space-x-3">
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

  <div class="mt-6 max-w-3xl mx-auto px-4">
    <Card class="border border-border bg-background shadow-md">
      {#if error}
        <CardContent class="p-6 text-red-500">{error}</CardContent>
      {:else if equipment}
        <CardHeader class="px-6 pb-0">
          <div class="flex items-start justify-between">
            <CardTitle class="text-2xl font-bold">{equipment.type}</CardTitle>
            <Badge
              class={`mt-1 ${
                equipment.status === "in service"
                  ? "bg-green-600 text-white"
                  : equipment.status === "not in service"
                    ? "bg-red-600 text-white"
                    : "bg-yellow-500 text-black"
              }`}
            >
              {equipment.status}
            </Badge>
          </div>
        </CardHeader>
        <Separator />
        <CardContent class="px-6 pb-6 space-y-6">
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
          <div class="mt-6">
            <Button
              onclick={handleReportClick}
              variant="outline"
              class="w-full justify-center text-sm font-medium gap-2 py-3 border-muted hover:bg-muted/50 transition"
            >
              <PlusCircle class="w-4 h-4" />
              Report Issue
            </Button>
          </div>
          {#if issues.length}
            <Separator />
            <div class="space-y-2">
              <h2 class="text-lg font-semibold">Issue History</h2>
              <Accordion type="single" class="w-full">
                {#each issues as issue}
                  <AccordionItem value={issue.id}>
                    <AccordionTrigger>
                      {new Date(issue.date_submitted).toLocaleDateString()} – {issue.progress}
                    </AccordionTrigger>
                    <AccordionContent>
                      <p class="text-sm text-muted-foreground">
                        {issue.description}
                      </p>
                    </AccordionContent>
                  </AccordionItem>
                {/each}
              </Accordion>
            </div>
          {/if}
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
