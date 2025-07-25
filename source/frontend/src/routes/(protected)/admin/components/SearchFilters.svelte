<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Input } from '$lib/components/ui/input/index.js';
  import { Label } from '$lib/components/ui/label/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select/index.js';
  
  import { filters, registrationActions } from '../stores/registrations.js';
  
  const dispatch = createEventDispatcher();
  
  let searchDebounceTimer: number;
  let emailDebounceTimer: number;
  
  // Local state for form inputs
  let searchValue = $filters.search;
  let emailValue = $filters.email;
  let startDate = $filters.dateRange.start || '';
  let endDate = $filters.dateRange.end || '';
  let urgencyLevel = $filters.urgencyLevel;
  let urgencyLevelArray = [urgencyLevel];
  
  // Debounced search handlers
  function handleSearchInput() {
    clearTimeout(searchDebounceTimer);
    searchDebounceTimer = setTimeout(() => {
      registrationActions.updateFilters({ search: searchValue });
    }, 300);
  }
  
  function handleEmailInput() {
    clearTimeout(emailDebounceTimer);
    emailDebounceTimer = setTimeout(() => {
      registrationActions.updateFilters({ email: emailValue });
    }, 300);
  }
  
  function handleDateRangeChange() {
    registrationActions.updateFilters({
      dateRange: {
        start: startDate || null,
        end: endDate || null
      }
    });
  }
  
  function handleUrgencyChange(value: string[]) {
    if (value && value.length > 0) {
      urgencyLevel = value[0] as typeof urgencyLevel;
      urgencyLevelArray = [urgencyLevel];
      registrationActions.updateFilters({ urgencyLevel });
    }
  }
  
  function clearFilters() {
    searchValue = '';
    emailValue = '';
    startDate = '';
    endDate = '';
    urgencyLevel = 'all';
    urgencyLevelArray = ['all'];
    
    registrationActions.updateFilters({
      search: '',
      email: '',
      dateRange: { start: null, end: null },
      urgencyLevel: 'all',
      status: 'all'
    });
  }
  
  function setDatePreset(preset: 'today' | 'week' | 'month') {
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    
    let start: Date;
    
    switch (preset) {
      case 'today':
        start = today;
        break;
      case 'week':
        start = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);
        break;
      case 'month':
        start = new Date(today.getFullYear(), today.getMonth() - 1, today.getDate());
        break;
    }
    
    startDate = start.toISOString().split('T')[0];
    endDate = now.toISOString().split('T')[0];
    handleDateRangeChange();
  }
</script>

<div class="space-y-4">
  <!-- Search Row -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
    <!-- Name/Username Search -->
    <div class="space-y-2">
      <Label for="search">Search Name</Label>
      <Input
        id="search"
        type="text"
        placeholder="Search by username..."
        bind:value={searchValue}
        oninput={handleSearchInput}
      />
    </div>
    
    <!-- Email Search -->
    <div class="space-y-2">
      <Label for="email">Search Email</Label>
      <Input
        id="email"
        type="email"
        placeholder="Search by email..."
        bind:value={emailValue}
        oninput={handleEmailInput}
      />
    </div>
    
    <!-- Urgency Filter -->
    <div class="space-y-2">
      <Label for="urgency">Urgency Level</Label>
      <Select bind:value={urgencyLevelArray} onValueChange={handleUrgencyChange} type="multiple">
        <SelectTrigger>
          {urgencyLevel === 'all' ? 'All Levels' :
           urgencyLevel === 'normal' ? 'Normal' :
           urgencyLevel === 'urgent' ? 'Urgent' :
           urgencyLevel === 'critical' ? 'Critical' : 'All Levels'}
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Levels</SelectItem>
          <SelectItem value="normal">Normal</SelectItem>
          <SelectItem value="urgent">Urgent</SelectItem>
          <SelectItem value="critical">Critical</SelectItem>
        </SelectContent>
      </Select>
    </div>
    
    <!-- Clear Filters -->
    <div class="space-y-2">
      <Label>&nbsp;</Label>
      <Button variant="outline" onclick={clearFilters} class="w-full">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        Clear Filters
      </Button>
    </div>
  </div>
  
  <!-- Date Range Row -->
  <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
    <!-- Start Date -->
    <div class="space-y-2">
      <Label for="startDate">From Date</Label>
      <Input
        id="startDate"
        type="date"
        bind:value={startDate}
        onchange={handleDateRangeChange}
      />
    </div>
    
    <!-- End Date -->
    <div class="space-y-2">
      <Label for="endDate">To Date</Label>
      <Input
        id="endDate"
        type="date"
        bind:value={endDate}
        onchange={handleDateRangeChange}
      />
    </div>
    
    <!-- Date Presets -->
    <div class="space-y-2">
      <Label>&nbsp;</Label>
      <Button variant="outline" size="sm" onclick={() => setDatePreset('today')} class="w-full">
        Today
      </Button>
    </div>
    
    <div class="space-y-2">
      <Label>&nbsp;</Label>
      <Button variant="outline" size="sm" onclick={() => setDatePreset('week')} class="w-full">
        This Week
      </Button>
    </div>
    
    <div class="space-y-2">
      <Label>&nbsp;</Label>
      <Button variant="outline" size="sm" onclick={() => setDatePreset('month')} class="w-full">
        This Month
      </Button>
    </div>
  </div>
  
  <!-- Active Filters Display -->
  {#if $filters.search || $filters.email || $filters.dateRange.start || $filters.urgencyLevel !== 'all'}
    <div class="flex flex-wrap gap-2 pt-2 border-t">
      <span class="text-sm text-muted-foreground">Active filters:</span>
      
      {#if $filters.search}
        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-blue-100 text-blue-800">
          Name: {$filters.search}
          <button
            class="ml-1 hover:text-blue-600"
            onclick={() => registrationActions.updateFilters({ search: '' })}
          >
            ×
          </button>
        </span>
      {/if}
      
      {#if $filters.email}
        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-green-100 text-green-800">
          Email: {$filters.email}
          <button
            class="ml-1 hover:text-green-600"
            onclick={() => registrationActions.updateFilters({ email: '' })}
          >
            ×
          </button>
        </span>
      {/if}
      
      {#if $filters.urgencyLevel !== 'all'}
        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-orange-100 text-orange-800">
          Urgency: {$filters.urgencyLevel}
          <button
            class="ml-1 hover:text-orange-600"
            onclick={() => registrationActions.updateFilters({ urgencyLevel: 'all' })}
          >
            ×
          </button>
        </span>
      {/if}
      
      {#if $filters.dateRange.start}
        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-purple-100 text-purple-800">
          Date: {$filters.dateRange.start} - {$filters.dateRange.end || 'now'}
          <button
            class="ml-1 hover:text-purple-600"
            onclick={() => registrationActions.updateFilters({ dateRange: { start: null, end: null } })}
          >
            ×
          </button>
        </span>
      {/if}
    </div>
  {/if}
</div>