<script lang="ts">
  import { ArrowLeft, Send } from "lucide-svelte";
  import { Button } from "$lib/components/ui/button";
  import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
  } from "$lib/components/ui/card";
  import { Input } from "$lib/components/ui/input";
  import { Textarea } from "$lib/components/ui/textarea";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import { Separator } from "$lib/components/ui/separator";
  import { goto } from "$app/navigation";
  import type { Issue } from "$lib/types/issue";
  import type { CreateIssueRequest } from "$lib/types/issue";
  import { createIssue } from "$lib/api/issue";

  export let data: {
    equipmentId: string | null;
  };

  let submitting = false;
  let error: string | null = null;
  let success = false;

  const formData: Omit<Issue, "id" | "date_submitted"> = {
    title: "",
    description: "",
    progress: "new",
    equipment_id: data.equipmentId ?? "",
  };

  
  async function submitIssue(): Promise<void> {
    console.log("[submitIssue] ⚙️ Submit function called");

    // Trim and validate basic inputs
    const title = formData.title.trim();
    const description = formData.description.trim();
    const equipmentId = formData.equipment_id?.trim();

    if (!title || !description) {
      error = "Title and description are required.";
      console.warn("[submitIssue] ❌ Missing title/description");
      return;
    }

    if (!equipmentId) {
      error = "Missing equipment ID.";
      console.warn("[submitIssue] ❌ Missing equipmentId");
      return;
    }

    submitting = true;
    error = null;

    try {
      const payload: CreateIssueRequest = {
        title,
        description,
        equipmentId,
      };

      console.log("[submitIssue] 📦 Sending payload:", payload);

      const res = await createIssue(payload);
      console.log("[submitIssue] ✅ Server response:", res);

      success = true;
    } catch (e) {
      console.error("[submitIssue] ❌ Error:", e);
      error = (e as Error).message ?? "Failed to submit issue.";
    } finally {
      submitting = false;
    }
  }
</script>

<div class="relative w-full">
  <div
    class="relative w-full h-[200px] sm:h-[240px] md:h-[280px] bg-muted flex items-center justify-center"
  >
    <div class="absolute top-4 left-4 z-10">
      <Button
        variant="ghost"
        size="icon"
        class="bg-black/60 text-white hover:bg-black/80"
        onclick={async () => {
          goto(`/equipment?scanned=${data.equipmentId}`);
        }}
      >
        <ArrowLeft class="w-5 h-5" />
      </Button>
    </div>
    <h1 class="text-xl md:text-2xl font-semibold text-white z-10">
      Submit New Issue
    </h1>
  </div>

  <div class="mt-6 max-w-2xl mx-auto px-4">
    <Card class="border border-border bg-background shadow-md">
      <CardHeader class="px-6 pb-0">
        <CardTitle class="text-2xl font-bold">Report an Issue</CardTitle>
      </CardHeader>
      <Separator />
      <CardContent class="px-6 py-6 space-y-6">
        {#if data.equipmentId}
          <form on:submit|preventDefault={submitIssue} class="space-y-5">
            <div class="space-y-1">
              <label
                for="title"
                class="block text-sm font-medium text-muted-foreground"
                >Title</label
              >
              <Input
                id="title"
                placeholder="Short summary"
                bind:value={formData.title}
              />
            </div>
            <div class="space-y-1">
              <label
                for="description"
                class="block text-sm font-medium text-muted-foreground"
                >Description</label
              >
              <Textarea
                id="description"
                rows={5}
                placeholder="Detailed issue description"
                bind:value={formData.description}
              />
            </div>
            <div class="space-y-1">
              <label
                for="progress"
                class="block text-sm font-medium text-muted-foreground"
                >Progress</label
              >
              <Input
                id="progress"
                bind:value={formData.progress}
                placeholder="e.g., new, in progress, resolved"
              />
            </div>
            <Button
              type="submit"
              disabled={submitting}
              class="w-full justify-center text-sm font-medium gap-2 py-3"
            >
              <Send class="w-4 h-4" />
              {submitting ? "Submitting..." : "Submit Issue"}
            </Button>
            {#if error}
              <p class="text-red-600 text-sm">{error}</p>
            {/if}
            {#if success}
              <p class="text-green-600 text-sm">
                Issue submitted successfully.
              </p>
            {/if}
          </form>
        {:else}
          <Skeleton class="h-5 w-3/5" />
          <p class="text-sm text-muted-foreground">No equipment selected.</p>
        {/if}
      </CardContent>
    </Card>
  </div>
</div>
