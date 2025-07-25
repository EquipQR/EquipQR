<script lang="ts">
  import { onMount, tick } from "svelte";
  import { ArrowLeft, PlusCircle, Plus } from "lucide-svelte";
  import { getEquipmentById, getEquipmentIssuesById } from "$lib/api/equipment";
  import { createIssue } from "$lib/api/issue";
  import { formatKey } from "$lib/utils";
  import type { Equipment } from "$lib/types/equipment";
  import type { Issue } from "$lib/types/issue";
  import type { CreateIssueRequest } from "$lib/types/issue";

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
  import * as Drawer from "$lib/components/ui/drawer/index.js";
  import { Input } from "$lib/components/ui/input";
  import { Textarea } from "$lib/components/ui/textarea";
  import { goto } from "$app/navigation";

  let equipmentId: string | null = null;
  let equipment: Equipment | null = null;
  let error: string | null = null;
  let imageLoaded = false;
  let imageRef: HTMLImageElement | null = null;
  let issues: Issue[] = [];

  export let data: { equipmentId: string | null };

  let drawerOpen = false;
  let title = "";
  let description = "";
  let submitting = false;
  let formError: string | null = null;

  let mediaFiles: File[] = [];

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

  async function handleReportClick(): Promise<void> {
    drawerOpen = true;
  }

  async function handleSubmitIssue(): Promise<void> {
    if (!title || !description || !data.equipmentId) {
      formError = "Please fill out all fields.";
      return;
    }

    submitting = true;
    formError = null;

    const payload: CreateIssueRequest = {
      title,
      description,
      equipmentId: data.equipmentId,
    };

    try {
      await createIssue(payload);
      issues = (await getEquipmentIssuesById(data.equipmentId)) ?? [];
      title = "";
      description = "";
      drawerOpen = false;
    } catch (e) {
      formError = "Failed to submit issue.";
    } finally {
      submitting = false;
    }
  }

  const MAX_FILES = 5;
  const MAX_SIZE_MB = 100;

  function handleFileChange(event: Event): void {
    const target = event.target as HTMLInputElement;
    if (!target.files) return;

    const selected = Array.from(target.files);
    const totalCount = mediaFiles.length + selected.length;

    if (totalCount > MAX_FILES) {
      formError = `You can only upload up to ${MAX_FILES} files.`;
      return;
    }

    const oversize = selected.find((f) => f.size > MAX_SIZE_MB * 1024 * 1024);
    if (oversize) {
      formError = `Each file must be under ${MAX_SIZE_MB}MB. '${oversize.name}' is too large.`;
      return;
    }

    mediaFiles = [...mediaFiles, ...selected];
    formError = null;
  }

  function removeFile(index: number): void {
    mediaFiles = mediaFiles.toSpliced(index, 1);
  }

  let previewUrl: string | null = null;
  let previewType: "image" | "video" | "unknown" | null = null;

  function openPreview(file: File): void {
    const url = URL.createObjectURL(file);
    previewUrl = url;
    if (file.type.startsWith("image/")) {
      previewType = "image";
    } else if (file.type.startsWith("video/")) {
      previewType = "video";
    } else {
      previewType = "unknown";
    }
  }

  function closePreview(): void {
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    previewUrl = null;
    previewType = null;
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

    <div class="absolute top-4 left-4 z-10">
      <Button
        variant="ghost"
        size="icon"
        class="bg-black/50 text-white hover:bg-black/70"
        onclick={async () => {
          await goto("/");
        }}
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

<Drawer.Root open={drawerOpen} onOpenChange={(v) => (drawerOpen = v)}>
  <Drawer.Trigger />
  <Drawer.Content>
    <Drawer.Header>
      <Drawer.Title>Report an Issue</Drawer.Title>
      <Drawer.Description>
        Fill out the form below to report a problem with this equipment.
      </Drawer.Description>
    </Drawer.Header>

    <div class="px-4 py-2 space-y-4">
      <Input bind:value={title} placeholder="Issue title" />
      <Textarea
        bind:value={description}
        rows={4}
        placeholder="Describe the issue..."
      />

      <!-- File Upload -->
      <div class="space-y-2">
        <label class="block text-sm font-medium text-gray-200">
          Attach Photos or Videos
        </label>
        <div class="flex flex-wrap gap-4">
          {#each mediaFiles as file, i}
            <div
              class="relative w-24 h-24 rounded overflow-hidden border border-neutral-700 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-white/80"
              role="button"
              tabindex="0"
              on:click={() => openPreview(file)}
              on:keydown={(e) => {
                if (e.key === "Enter" || e.key === " ") openPreview(file);
              }}
              aria-label="Preview file"
            >
              {#if file.type.startsWith("image/")}
                <img
                  src={URL.createObjectURL(file)}
                  alt="preview"
                  class="object-cover w-full h-full pointer-events-none"
                />
              {:else if file.type.startsWith("video/")}
                <video
                  src={URL.createObjectURL(file)}
                  class="object-cover w-full h-full pointer-events-none"
                  muted
                ></video>
              {:else}
                <div
                  class="flex items-center justify-center w-full h-full text-sm text-gray-400 pointer-events-none"
                >
                  Unsupported
                </div>
              {/if}

              <button
                type="button"
                class="absolute top-0 right-0 bg-black bg-opacity-70 text-white text-xs px-1 py-0.5 z-10"
                on:click|stopPropagation={() => removeFile(i)}
                aria-label="Remove file"
              >
                ✕
              </button>
            </div>
          {/each}

          {#if mediaFiles.length < MAX_FILES}
            <label
              for="media"
              class="flex items-center justify-center w-24 h-24 border-2 border-dashed border-neutral-700 rounded hover:border-neutral-500 cursor-pointer"
            >
              <Plus class="w-6 h-6 text-neutral-500" />
              <input
                id="media"
                type="file"
                multiple
                accept="image/*,video/*"
                capture
                on:change={handleFileChange}
                class="hidden"
              />
            </label>
          {/if}
        </div>
      </div>

      {#if formError}
        <p class="text-sm text-red-500">{formError}</p>
      {/if}
    </div>

    <Drawer.Footer class="px-4 pb-4">
      <Button onclick={handleSubmitIssue} disabled={submitting} class="w-full">
        {submitting ? "Submitting..." : "Submit"}
      </Button>
      <Drawer.Close class="mt-2 w-full">Cancel</Drawer.Close>
    </Drawer.Footer>
  </Drawer.Content>
</Drawer.Root>

{#if previewUrl}
  <div
    class="fixed inset-0 z-[100] bg-black/80 backdrop-blur-sm flex items-center justify-center p-4"
    role="dialog"
    aria-modal="true"
    tabindex="0"
    on:click={closePreview}
    on:keydown={(e) => {
      if (e.key === "Escape" || e.key === "Enter" || e.key === " ") {
        closePreview();
      }
    }}
  >
    {#if previewType === "image"}
      <img
        src={previewUrl}
        class="max-h-[90vh] max-w-[90vw] rounded shadow-lg pointer-events-none"
        alt="Full size preview"
      />
    {:else if previewType === "video"}
      <video
        src={previewUrl}
        controls
        autoplay
        class="max-h-[90vh] max-w-[90vw] rounded shadow-lg"
        aria-label="Preview video, captions not available"
      >
        <track kind="captions" src="" label="No captions available" />
      </video>
    {:else}
      <p class="text-white text-lg">Unsupported preview</p>
    {/if}
  </div>
{/if}
