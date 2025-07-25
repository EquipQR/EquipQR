<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '$lib/components/ui/button/index.js';
  import { fly } from 'svelte/transition';
  
  export let show = false;
  export let message = 'Action completed';
  export let duration = 5000;
  
  const dispatch = createEventDispatcher<{
    undo: void;
    dismiss: void;
  }>();
  
  let timeoutId: number;
  
  $: if (show) {
    timeoutId = setTimeout(() => {
      show = false;
      dispatch('dismiss');
    }, duration);
  } else {
    clearTimeout(timeoutId);
  }
  
  function handleUndo() {
    dispatch('undo');
    show = false;
    clearTimeout(timeoutId);
  }
  
  function handleDismiss() {
    show = false;
    dispatch('dismiss');
    clearTimeout(timeoutId);
  }
</script>

{#if show}
  <div 
    class="fixed bottom-4 left-1/2 transform -translate-x-1/2 z-50"
    transition:fly={{ y: 100, duration: 200 }}
  >
    <div class="bg-gray-900 text-white px-4 py-3 rounded-lg shadow-lg flex items-center space-x-3 min-w-80">
      <div class="flex-1">
        <p class="text-sm font-medium">{message}</p>
      </div>
      
      <Button 
        variant="ghost" 
        size="sm"
        class="text-blue-400 hover:text-blue-300 hover:bg-gray-800"
        onclick={handleUndo}
      >
        Undo
      </Button>
      
      <button 
        class="text-gray-400 hover:text-white"
        onclick={handleDismiss}
        aria-label="Dismiss"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  </div>
{/if}