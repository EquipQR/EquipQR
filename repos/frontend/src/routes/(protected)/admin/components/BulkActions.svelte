<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Badge } from '$lib/components/ui/badge/index.js';
  import { Separator } from '$lib/components/ui/separator/index.js';
  
  export let selectedCount: number;
  
  const dispatch = createEventDispatcher<{
    bulkApprove: void;
    bulkDeny: void;
    clearSelection: void;
    exportCsv: void;
    generateInvites: void;
  }>();
  
  function handleBulkApprove() {
    dispatch('bulkApprove');
  }
  
  function handleBulkDeny() {
    dispatch('bulkDeny');
  }
  
  function handleClearSelection() {
    dispatch('clearSelection');
  }
  
  function handleExportCsv() {
    dispatch('exportCsv');
  }
</script>

<div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
  <!-- Selection Info -->
  <div class="flex items-center gap-3">
    <div class="flex items-center gap-2">
      <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="font-medium">
        <Badge variant="secondary" class="mr-2">{selectedCount}</Badge>
        {selectedCount === 1 ? 'registration' : 'registrations'} selected
      </span>
    </div>
    
    <Separator orientation="vertical" class="h-6" />
    
    <Button 
      variant="ghost" 
      size="sm" 
      onclick={handleClearSelection}
      class="text-muted-foreground hover:text-foreground"
    >
      Clear selection
    </Button>
  </div>
  
  <!-- Bulk Actions -->
  <div class="flex items-center gap-2">
    <!-- Approve All -->
    <Button 
      variant="default"
      size="sm"
      onclick={handleBulkApprove}
      class="bg-green-600 hover:bg-green-700"
    >
      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
      </svg>
      Approve All ({selectedCount})
    </Button>
    
    <!-- Deny All -->
    <Button 
      variant="destructive"
      size="sm"
      onclick={handleBulkDeny}
    >
      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
      Deny All ({selectedCount})
    </Button>
    
    <Separator orientation="vertical" class="h-6" />
    
    <!-- Export CSV -->
    <Button 
      variant="outline"
      size="sm"
      onclick={handleExportCsv}
    >
      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      Export CSV
    </Button>
    
    <!-- Generate Invites for Selected -->
    <Button 
      variant="outline"
      size="sm"
      onclick={() => dispatch('generateInvites')}
    >
      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
      </svg>
      Generate Invites
    </Button>
  </div>
</div>

<!-- Bulk Action Tips -->
<div class="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
  <div class="flex items-start gap-2">
    <svg class="w-4 h-4 text-blue-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    <div class="text-sm text-blue-800">
      <p class="font-medium mb-1">Bulk Action Tips:</p>
      <ul class="space-y-1 text-blue-700">
        <li>• <strong>Approve All:</strong> Instantly approve all selected registrations and add users to your business</li>
        <li>• <strong>Deny All:</strong> Reject selected registrations (can be undone within 5 seconds)</li>
        <li>• <strong>Generate Invites:</strong> Create invitation links for selected users instead of direct approval</li>
        <li>• <strong>Export CSV:</strong> Download selected registration data for external processing</li>
      </ul>
    </div>
  </div>
</div>