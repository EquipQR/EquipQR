<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '$lib/components/ui/dialog/index.js';
  
  export let open = false;
  export let title = 'Confirm Action';
  export let message = 'Are you sure you want to proceed?';
  export let confirmText = 'Confirm';
  export let cancelText = 'Cancel';
  export let variant: 'default' | 'destructive' = 'default';
  
  const dispatch = createEventDispatcher<{
    confirm: void;
    cancel: void;
  }>();
  
  function handleConfirm() {
    dispatch('confirm');
    open = false;
  }
  
  function handleCancel() {
    dispatch('cancel');
    open = false;
  }
</script>

<Dialog bind:open>
  <DialogContent class="sm:max-w-md">
    <DialogHeader>
      <DialogTitle>{title}</DialogTitle>
      <DialogDescription>
        {message}
      </DialogDescription>
    </DialogHeader>
    
    <DialogFooter>
      <Button variant="outline" onclick={handleCancel}>
        {cancelText}
      </Button>
      <Button variant={variant} onclick={handleConfirm}>
        {confirmText}
      </Button>
    </DialogFooter>
  </DialogContent>
</Dialog>