<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { page } from '$app/stores';
  import type { PageData } from './$types';
  
  // UI Components
  import { Button } from '$lib/components/ui/button/index.js';
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
  import { Badge } from '$lib/components/ui/badge/index.js';
  import { Separator } from '$lib/components/ui/separator/index.js';
  import { toast } from 'svelte-sonner';
  
  // Admin Components
  import RegistrationTable from './components/RegistrationTable.svelte';
  import SearchFilters from './components/SearchFilters.svelte';
  import BulkActions from './components/BulkActions.svelte';
  import InviteModal from './components/InviteModal.svelte';
  import ConfirmationDialog from './components/ConfirmationDialog.svelte';
  import UndoSnackbar from './components/UndoSnackbar.svelte';
  
  // Stores
  import { 
    registrations, 
    loading, 
    error, 
    adminStats, 
    selectedIds,
    registrationActions 
  } from './stores/registrations.js';
  
  export let data: PageData & { businessId?: string };
  
  // Component state
  let showInviteModal = false;
  let showConfirmDialog = false;
  let confirmAction: (() => void) | null = null;
  let confirmMessage = '';
  let undoAction: (() => void) | null = null;
  let undoMessage = '';
  let showUndo = false;
  
  // Reactive statements
  $: hasSelectedItems = $selectedIds.length > 0;
  $: stats = $adminStats;
  
  onMount(() => {
    const businessId = (data as any).businessId || 'default-business-id';
    
    // Load initial data if not already loaded
    if ($registrations.length === 0) {
      registrationActions.loadRegistrations(businessId);
    }
    
    // Set up periodic refresh
    const refreshInterval = setInterval(() => {
      registrationActions.loadRegistrations(businessId);
    }, 30000); // Refresh every 30 seconds
    
    return () => {
      clearInterval(refreshInterval);
    };
  });
  
  // Error handling
  $: if ($error) {
    toast.error($error);
  }
  
  // Action handlers
  function handleApprove(registrationId: string) {
    confirmAction = () => {
      registrationActions.approveRegistration(registrationId);
      toast.success('Registration approved successfully');
    };
    confirmMessage = 'Are you sure you want to approve this registration?';
    showConfirmDialog = true;
  }
  
  function handleDeny(registrationId: string) {
    confirmAction = () => {
      registrationActions.denyRegistration(registrationId);
      
      // Show undo option
      undoAction = () => {
        // This would need to be implemented in the backend
        toast.info('Undo functionality coming soon');
      };
      undoMessage = 'Registration denied';
      showUndo = true;
      
      // Auto-hide undo after 5 seconds
      setTimeout(() => {
        showUndo = false;
        undoAction = null;
      }, 5000);
      
      toast.success('Registration denied');
    };
    confirmMessage = 'Are you sure you want to deny this registration? This action can be undone within 5 seconds.';
    showConfirmDialog = true;
  }
  
  function handleBulkApprove() {
    if ($selectedIds.length === 0) return;
    
    confirmAction = () => {
      registrationActions.bulkAction({
        type: 'approve',
        registration_ids: $selectedIds
      });
      toast.success(`${$selectedIds.length} registrations approved`);
    };
    confirmMessage = `Are you sure you want to approve ${$selectedIds.length} selected registrations?`;
    showConfirmDialog = true;
  }
  
  function handleBulkDeny() {
    if ($selectedIds.length === 0) return;
    
    confirmAction = () => {
      registrationActions.bulkAction({
        type: 'deny',
        registration_ids: $selectedIds
      });
      toast.success(`${$selectedIds.length} registrations denied`);
    };
    confirmMessage = `Are you sure you want to deny ${$selectedIds.length} selected registrations?`;
    showConfirmDialog = true;
  }
  
  function handleGenerateInvite() {
    showInviteModal = true;
  }
  
  function handleRefresh() {
    const businessId = (data as any).businessId || 'default-business-id';
    registrationActions.loadRegistrations(businessId);
    toast.info('Data refreshed');
  }
  
  function confirmActionHandler() {
    if (confirmAction) {
      confirmAction();
      confirmAction = null;
    }
    showConfirmDialog = false;
  }
  
  function cancelConfirmation() {
    confirmAction = null;
    showConfirmDialog = false;
  }
</script>

<svelte:head>
  <title>Admin Dashboard - Registration Management</title>
</svelte:head>

<div class="container mx-auto p-6 space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Registration Management</h1>
      <p class="text-muted-foreground">
        Manage user registration requests and generate invites
      </p>
    </div>
    
    <div class="flex gap-2">
      <Button variant="outline" onclick={handleRefresh} disabled={$loading}>
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Refresh
      </Button>
      
      <Button onclick={handleGenerateInvite}>
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        Generate Invite
      </Button>
    </div>
  </div>
  
  <!-- Stats Cards -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">Total Pending</CardTitle>
        <svg class="h-4 w-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">{stats.totalPending}</div>
        <p class="text-xs text-muted-foreground">
          Registration requests
        </p>
      </CardContent>
    </Card>
    
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">Urgent</CardTitle>
        <Badge variant="destructive" class="h-4 w-4 p-0 text-xs">!</Badge>
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold text-orange-600">{stats.urgentCount}</div>
        <p class="text-xs text-muted-foreground">
          Require attention
        </p>
      </CardContent>
    </Card>
    
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">Today</CardTitle>
        <svg class="h-4 w-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">{stats.todayCount}</div>
        <p class="text-xs text-muted-foreground">
          New requests today
        </p>
      </CardContent>
    </Card>
    
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">This Week</CardTitle>
        <svg class="h-4 w-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">{stats.weeklyCount}</div>
        <p class="text-xs text-muted-foreground">
          Weekly total
        </p>
      </CardContent>
    </Card>
  </div>
  
  <!-- Search and Filters -->
  <Card>
    <CardHeader>
      <CardTitle>Search & Filters</CardTitle>
      <CardDescription>
        Find specific registration requests
      </CardDescription>
    </CardHeader>
    <CardContent>
      <SearchFilters />
    </CardContent>
  </Card>
  
  <!-- Bulk Actions -->
  {#if hasSelectedItems}
    <Card class="border-blue-200 bg-blue-50">
      <CardContent class="pt-6">
        <BulkActions 
          selectedCount={$selectedIds.length}
          on:bulkApprove={handleBulkApprove}
          on:bulkDeny={handleBulkDeny}
          on:clearSelection={() => registrationActions.clearSelection()}
        />
      </CardContent>
    </Card>
  {/if}
  
  <!-- Registration Table -->
  <Card>
    <CardHeader>
      <CardTitle>Registration Requests</CardTitle>
      <CardDescription>
        {#if $loading}
          Loading registration requests...
        {:else}
          {stats.totalPending} pending registration{stats.totalPending !== 1 ? 's' : ''}
        {/if}
      </CardDescription>
    </CardHeader>
    <CardContent>
      <RegistrationTable 
        on:approve={(e) => handleApprove(e.detail)}
        on:deny={(e) => handleDeny(e.detail)}
        on:generateInvite={handleGenerateInvite}
      />
    </CardContent>
  </Card>
</div>

<!-- Modals and Dialogs -->
<InviteModal 
  bind:open={showInviteModal}
  businessId={(data as any).businessId || 'default-business-id'}
/>

<ConfirmationDialog
  bind:open={showConfirmDialog}
  title="Confirm Action"
  message={confirmMessage}
  on:confirm={confirmActionHandler}
  on:cancel={cancelConfirmation}
/>

<UndoSnackbar
  bind:show={showUndo}
  message={undoMessage}
  on:undo={() => undoAction?.()}
/>

<style>
  :global(.urgency-normal) {
    border-left: 4px solid #10b981;
  }
  
  :global(.urgency-urgent) {
    border-left: 4px solid #f59e0b;
  }
  
  :global(.urgency-critical) {
    border-left: 4px solid #ef4444;
  }
  
  :global(.status-pending) {
    background: #fef3c7;
    color: #92400e;
  }
  
  :global(.status-approved) {
    background: #d1fae5;
    color: #065f46;
  }
  
  :global(.status-denied) {
    background: #fee2e2;
    color: #991b1b;
  }
</style>