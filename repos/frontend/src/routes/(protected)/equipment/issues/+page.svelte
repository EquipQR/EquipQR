<script lang="ts">
  import { ArrowLeft } from "lucide-svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
  } from "$lib/components/ui/card";
  import { Separator } from "$lib/components/ui/separator";
  import { Button } from "$lib/components/ui/button";
  import { Badge, badgeVariants } from "$lib/components/ui/badge/index.js";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import type { Issue } from "$lib/types/issue";

  export let data: {
    equipmentId: string;
  };

  let issues: list[Issue] = [];
  let loading = true;

  onMount(() => {
    issues = [
      {
        id: "issue-1",
        title: "Battery not holding charge",
        description: "The battery drains within 30 minutes of usage.",
        progress: "new",
        equipment_id: data.equipmentId,
        assignee_id: "tech-123",
        date_submitted: new Date().toISOString(),
      },
      {
        id: "issue-2",
        title: "Loose wiring on control panel",
        description: "Exposed wires behind the panel may be a hazard.",
        progress: "in progress",
        equipment_id: data.equipmentId,
        assignee_id: "tech-456",
        date_submitted: new Date(Date.now() - 86400000).toISOString(),
      },
      {
        id: "issue-3",
        title: "Missing safety label",
        description: "The caution label has peeled off.",
        progress: "resolved",
        equipment_id: data.equipmentId,
        assignee_id: "tech-789",
        date_submitted: new Date(Date.now() - 172800000).toISOString(),
      },
    ];
    loading = false;
  });

  function statusColor(progress: string): string {
    switch (progress.toLowerCase()) {
      case "new":
        return badgeVariants({ variant: "default" });
      case "in progress":
        return badgeVariants({ variant: "secondary" });
      case "resolved":
        return badgeVariants({ variant: "outline" });
      default:
        return badgeVariants({ variant: "outline" });
    }
  }

  function updateStatus(id: string, newStatus: string): void {
    issues = issues.map((issue) =>
      issue.id === id ? { ...issue, progress: newStatus } : issue
    );
  }
</script>

<div class="relative w-full">
  <div class="relative w-full h-48 sm:h-60 md:h-72 bg-gradient-to-br from-muted to-background flex items-center justify-center shadow-inner">
    <div class="absolute top-4 left-4 z-10">
      <Button
        variant="ghost"
        size="icon"
        class="bg-black/50 text-white hover:bg-black/70"
        onclick={() => goto(`/equipment?scanned=${data.equipmentId}`)}
      >
        <ArrowLeft class="w-5 h-5" />
      </Button>
    </div>
    <h1 class="text-2xl md:text-3xl font-bold text-white drop-shadow z-10">
      Reported Issues
    </h1>
  </div>

  <div class="mt-8 max-w-3xl mx-auto px-4 space-y-6">
    {#if loading}
      <p class="text-center text-muted-foreground">Loading issues...</p>
    {:else if issues.length > 0}
      {#each issues as issue}
        <Card class="border border-border bg-background rounded-xl shadow-sm hover:shadow-md transition">
          <CardHeader class="pb-2">
            <div class="flex flex-col space-y-1">
              <CardTitle class="text-lg font-semibold">{issue.title}</CardTitle>
              <span class="text-xs text-muted-foreground">
                {new Date(issue.date_submitted).toLocaleString()}
              </span>
            </div>
          </CardHeader>
          <CardContent class="pt-0 space-y-4 text-sm">
            <p class="text-foreground leading-relaxed">{issue.description}</p>
            <Separator />
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 text-xs text-muted-foreground">
              <div class="flex items-center gap-2">
                <span class="font-medium text-foreground">Status:</span>
                <span class={statusColor(issue.progress)}>{issue.progress}</span>

                <DropdownMenu.Root>
                  <DropdownMenu.Trigger asChild>
                    <Button variant="ghost" size="xs" class="px-2 text-xs">
                      Update
                    </Button>
                  </DropdownMenu.Trigger>
                  <DropdownMenu.Content class="z-50">
                    <DropdownMenu.Item onSelect={() => updateStatus(issue.id, "new")}>
                      Mark as New
                    </DropdownMenu.Item>
                    <DropdownMenu.Item onSelect={() => updateStatus(issue.id, "in progress")}>
                      Mark as In Progress
                    </DropdownMenu.Item>
                    <DropdownMenu.Item onSelect={() => updateStatus(issue.id, "resolved")}>
                      Mark as Resolved
                    </DropdownMenu.Item>
                  </DropdownMenu.Content>
                </DropdownMenu.Root>
              </div>
              <div>
                <span class="font-medium text-foreground">Assignee:</span>
                <span class="ml-1 text-foreground">{issue.assignee_id}</span>
              </div>
            </div>
            <div class="pt-2 text-right">
              <Button
                variant="secondary"
                size="sm"
                class="text-xs font-medium"
                onclick={() => goto(`/issue/${issue.id}`)}
              >
                Work on this issue
              </Button>
            </div>
          </CardContent>
        </Card>
      {/each}
    {:else}
      <p class="text-center text-muted-foreground">No issues submitted yet.</p>
    {/if}
  </div>
</div>
