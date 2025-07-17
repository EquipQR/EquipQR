<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Checkbox } from '$lib/components/ui/checkbox/index.js';
  import { Badge } from '$lib/components/ui/badge/index.js';
  import { Avatar, AvatarFallback } from '$lib/components/ui/avatar/index.js';
  import { Skeleton } from '$lib/components/ui/skeleton/index.js';
  
  import { 
    paginatedRegistrations, 
    loading, 
    actionLoading, 
    selectedIds, 
    sort,
    pagination,
    registrationActions 
  } from '../stores/registrations.js';
  
  import type { PendingRegistration } from '../types/admin.js';
  
  const dispatch = createEventDispatcher<{
    approve: string;
    deny: string;
    generateInvite: string;
  }>();
  
  function handleSort(column: 'created_at' | 'email' | 'username' | 'days_pending') {
    registrationActions.updateSort(column);
  }
  
  function handleSelectAll() {
    if ($selectedIds.length === $paginatedRegistrations.length) {
      registrationActions.clearSelection();
    } else {
      registrationActions.selectAll();
    }
  }
  
  function handleSelectRow(registrationId: string) {
    registrationActions.toggleSelection(registrationId);
  }
  
  function getUrgencyClass(urgency: string): string {
    switch (urgency) {
      case 'critical': return 'urgency-critical';
      case 'urgent': return 'urgency-urgent';
      default: return 'urgency-normal';
    }
  }
  
  function getUrgencyBadge(urgency: string) {
    switch (urgency) {
      case 'critical': return { variant: 'destructive' as const, text: 'Critical' };
      case 'urgent': return { variant: 'secondary' as const, text: 'Urgent' };
      default: return { variant: 'outline' as const, text: 'Normal' };
    }
  }
  
  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;
    
    return date.toLocaleDateString();
  }
  
  function getUserInitials(username: string): string {
    return username.slice(0, 2).toUpperCase();
  }
  
  $: allSelected = $selectedIds.length === $paginatedRegistrations.length && $paginatedRegistrations.length > 0;
  $: someSelected = $selectedIds.length > 0 && $selectedIds.length < $paginatedRegistrations.length;
</script>

<div class="space-y-4">
  <!-- Desktop Table View -->
  <div class="hidden md:block">
    <div class="rounded-md border">
      <table class="w-full">
        <thead>
          <tr class="border-b bg-muted/50">
            <th class="h-12 px-4 text-left align-middle font-medium">
              <Checkbox 
                checked={allSelected}
                indeterminate={someSelected}
                onCheckedChange={handleSelectAll}
              />
            </th>
            
            <th class="h-12 px-4 text-left align-middle font-medium">
              <button 
                class="flex items-center gap-2 hover:text-foreground"
                onclick={() => handleSort('username')}
              >
                User
                {#if $sort.column === 'username'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                          d={$sort.direction === 'asc' ? 'M5 15l7-7 7 7' : 'M19 9l-7 7-7-7'} />
                  </svg>
                {/if}
              </button>
            </th>
            
            <th class="h-12 px-4 text-left align-middle font-medium">
              <button 
                class="flex items-center gap-2 hover:text-foreground"
                onclick={() => handleSort('email')}
              >
                Email
                {#if $sort.column === 'email'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                          d={$sort.direction === 'asc' ? 'M5 15l7-7 7 7' : 'M19 9l-7 7-7-7'} />
                  </svg>
                {/if}
              </button>
            </th>
            
            <th class="h-12 px-4 text-left align-middle font-medium">
              <button 
                class="flex items-center gap-2 hover:text-foreground"
                onclick={() => handleSort('created_at')}
              >
                Request Date
                {#if $sort.column === 'created_at'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                          d={$sort.direction === 'asc' ? 'M5 15l7-7 7 7' : 'M19 9l-7 7-7-7'} />
                  </svg>
                {/if}
              </button>
            </th>
            
            <th class="h-12 px-4 text-left align-middle font-medium">
              <button 
                class="flex items-center gap-2 hover:text-foreground"
                onclick={() => handleSort('days_pending')}
              >
                Days Pending
                {#if $sort.column === 'days_pending'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                          d={$sort.direction === 'asc' ? 'M5 15l7-7 7 7' : 'M19 9l-7 7-7-7'} />
                  </svg>
                {/if}
              </button>
            </th>
            
            <th class="h-12 px-4 text-left align-middle font-medium">Status</th>
            <th class="h-12 px-4 text-right align-middle font-medium">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if $loading}
            {#each Array(5) as _}
              <tr class="border-b">
                <td class="p-4"><Skeleton class="h-4 w-4" /></td>
                <td class="p-4"><Skeleton class="h-8 w-32" /></td>
                <td class="p-4"><Skeleton class="h-4 w-48" /></td>
                <td class="p-4"><Skeleton class="h-4 w-24" /></td>
                <td class="p-4"><Skeleton class="h-4 w-16" /></td>
                <td class="p-4"><Skeleton class="h-6 w-20" /></td>
                <td class="p-4"><Skeleton class="h-8 w-32" /></td>
              </tr>
            {/each}
          {:else if $paginatedRegistrations.length === 0}
            <tr>
              <td colspan="7" class="h-24 text-center">
                <div class="flex flex-col items-center justify-center space-y-2">
                  <svg class="w-8 h-8 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-2.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 009.586 13H7" />
                  </svg>
                  <p class="text-muted-foreground">No registration requests found</p>
                </div>
              </td>
            </tr>
          {:else}
            {#each $paginatedRegistrations as registration (registration.id)}
              <tr class="border-b hover:bg-muted/50 {getUrgencyClass(registration.urgency_level)}">
                <td class="p-4">
                  <Checkbox 
                    checked={$selectedIds.includes(registration.id)}
                    onCheckedChange={() => handleSelectRow(registration.id)}
                  />
                </td>
                
                <td class="p-4">
                  <div class="flex items-center space-x-3">
                    <Avatar class="h-8 w-8">
                      <AvatarFallback>{getUserInitials(registration.user.username)}</AvatarFallback>
                    </Avatar>
                    <div>
                      <div class="font-medium">{registration.user.username}</div>
                      <div class="text-sm text-muted-foreground">ID: {registration.user_id.slice(0, 8)}...</div>
                    </div>
                  </div>
                </td>
                
                <td class="p-4">
                  <div class="font-medium">{registration.user.email}</div>
                </td>
                
                <td class="p-4">
                  <div class="text-sm">
                    {formatDate(registration.created_at)}
                  </div>
                </td>
                
                <td class="p-4">
                  <div class="flex items-center space-x-2">
                    <span class="text-sm font-medium">{registration.days_pending}</span>
                    <Badge variant={getUrgencyBadge(registration.urgency_level).variant}>
                      {getUrgencyBadge(registration.urgency_level).text}
                    </Badge>
                  </div>
                </td>
                
                <td class="p-4">
                  <Badge variant="secondary">Pending</Badge>
                </td>
                
                <td class="p-4 text-right">
                  <div class="flex items-center justify-end space-x-2">
                    <Button 
                      size="sm" 
                      variant="outline"
                      disabled={$actionLoading[registration.id]}
                      onclick={() => dispatch('approve', registration.id)}
                    >
                      {#if $actionLoading[registration.id]}
                        <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                      {:else}
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                      {/if}
                      Approve
                    </Button>
                    
                    <Button 
                      size="sm" 
                      variant="destructive"
                      disabled={$actionLoading[registration.id]}
                      onclick={() => dispatch('deny', registration.id)}
                    >
                      {#if $actionLoading[registration.id]}
                        <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                      {:else}
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      {/if}
                      Deny
                    </Button>
                    
                    <Button 
                      size="sm" 
                      variant="outline"
                      onclick={() => dispatch('generateInvite', registration.user.email)}
                    >
                      <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                      </svg>
                      Invite
                    </Button>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
  
  <!-- Mobile Card View -->
  <div class="md:hidden space-y-4">
    {#if $loading}
      {#each Array(3) as _}
        <div class="border rounded-lg p-4 space-y-3">
          <Skeleton class="h-4 w-3/4" />
          <Skeleton class="h-4 w-1/2" />
          <Skeleton class="h-8 w-full" />
        </div>
      {/each}
    {:else if $paginatedRegistrations.length === 0}
      <div class="text-center py-8">
        <svg class="w-12 h-12 text-muted-foreground mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-2.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 009.586 13H7" />
        </svg>
        <p class="text-muted-foreground">No registration requests found</p>
      </div>
    {:else}
      {#each $paginatedRegistrations as registration (registration.id)}
        <div class="border rounded-lg p-4 space-y-3 {getUrgencyClass(registration.urgency_level)}">
          <div class="flex items-center justify-between">
            <Checkbox 
              checked={$selectedIds.includes(registration.id)}
              onCheckedChange={() => handleSelectRow(registration.id)}
            />
            <Badge variant={getUrgencyBadge(registration.urgency_level).variant}>
              {getUrgencyBadge(registration.urgency_level).text}
            </Badge>
          </div>
          
          <div class="flex items-center space-x-3">
            <Avatar class="h-10 w-10">
              <AvatarFallback>{getUserInitials(registration.user.username)}</AvatarFallback>
            </Avatar>
            <div class="flex-1">
              <div class="font-medium">{registration.user.username}</div>
              <div class="text-sm text-muted-foreground">{registration.user.email}</div>
            </div>
          </div>
          
          <div class="flex justify-between text-sm text-muted-foreground">
            <span>Requested: {formatDate(registration.created_at)}</span>
            <span>{registration.days_pending} days pending</span>
          </div>
          
          <div class="flex space-x-2">
            <Button 
              size="sm" 
              variant="outline" 
              class="flex-1"
              disabled={$actionLoading[registration.id]}
              onclick={() => dispatch('approve', registration.id)}
            >
              Approve
            </Button>
            <Button 
              size="sm" 
              variant="destructive" 
              class="flex-1"
              disabled={$actionLoading[registration.id]}
              onclick={() => dispatch('deny', registration.id)}
            >
              Deny
            </Button>
            <Button 
              size="sm" 
              variant="outline"
              onclick={() => dispatch('generateInvite', registration.user.email)}
            >
              Invite
            </Button>
          </div>
        </div>
      {/each}
    {/if}
  </div>
  
  <!-- Pagination -->
  {#if $pagination.totalPages > 1}
    <div class="flex items-center justify-between">
      <div class="text-sm text-muted-foreground">
        Showing {($pagination.page - 1) * $pagination.pageSize + 1} to {Math.min($pagination.page * $pagination.pageSize, $pagination.total)} of {$pagination.total} results
      </div>
      
      <div class="flex items-center space-x-2">
        <Button 
          variant="outline" 
          size="sm"
          disabled={$pagination.page === 1}
          onclick={() => registrationActions.updatePagination({ page: $pagination.page - 1 })}
        >
          Previous
        </Button>
        
        <div class="flex items-center space-x-1">
          {#each Array($pagination.totalPages) as _, i}
            {#if i + 1 === $pagination.page}
              <Button variant="default" size="sm" disabled>
                {i + 1}
              </Button>
            {:else if Math.abs(i + 1 - $pagination.page) <= 2 || i === 0 || i === $pagination.totalPages - 1}
              <Button 
                variant="outline" 
                size="sm"
                onclick={() => registrationActions.updatePagination({ page: i + 1 })}
              >
                {i + 1}
              </Button>
            {:else if Math.abs(i + 1 - $pagination.page) === 3}
              <span class="px-2">...</span>
            {/if}
          {/each}
        </div>
        
        <Button 
          variant="outline" 
          size="sm"
          disabled={$pagination.page === $pagination.totalPages}
          onclick={() => registrationActions.updatePagination({ page: $pagination.page + 1 })}
        >
          Next
        </Button>
      </div>
    </div>
  {/if}
</div>