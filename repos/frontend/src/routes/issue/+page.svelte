<script lang="ts">
  import { ArrowLeft, Send } from "lucide-svelte";
  import { v4 as uuidv4 } from "uuid";
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
  import type { Issue } from "$lib/types/issue";
  import { goto } from "$app/navigation";

  export let data: {
    equipmentId: string | null;
  };

  let submitting = false;
  let error: string | null = null;
  let success = false;

  let formData: Omit<Issue, "id" | "date_submitted" | "assignee_id"> & {
    assignee_id: string | null;
  } = {
    title: "",
    description: "",
    progress: "new",
    equipment_id: data.equipmentId ?? "",
    assignee_id: null,
  };

  async function submitIssue(): Promise<void> {
    if (!formData.title.trim() || !formData.description.trim()) {
      error = "Title and description are required.";
      return;
    }

    submitting = true;
    error = null;
    try {
      const newIssue: Issue = {
        id: uuidv4(),
        title: formData.title,
        description: formData.description,
        progress: formData.progress,
        equipment_id: formData.equipment_id,
        assignee_id: formData.assignee_id ?? "unassigned",
        date_submitted: new Date().toISOString(),
      };

      // await postNewIssue(newIssue); // implement this in your API
      console.log("Submitted Issue:", newIssue);
      success = true;
    } catch (e) {
      error = "Failed to submit issue.";
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
        onclick={() =>
          (goto(`/equipment?scanned=${data.equipmentId}`))}
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
