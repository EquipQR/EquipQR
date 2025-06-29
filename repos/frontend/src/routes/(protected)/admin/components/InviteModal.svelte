<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '$lib/components/ui/button/index.js';
  import { Input } from '$lib/components/ui/input/index.js';
  import { Label } from '$lib/components/ui/label/index.js';
  import { Checkbox } from '$lib/components/ui/checkbox/index.js';
  import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '$lib/components/ui/dialog/index.js';
  import { Badge } from '$lib/components/ui/badge/index.js';
  import { toast } from 'svelte-sonner';
  
  import { registrationActions } from '../stores/registrations.js';
  import type { InviteOptions } from '../types/admin.js';
  
  export let open = false;
  export const businessId: string = '';
  export let prefilledEmail = '';
  
  const dispatch = createEventDispatcher<{
    close: void;
    inviteGenerated: { email: string; link: string };
  }>();
  
  // Form state
  let email = prefilledEmail;
  let expirationDays = 7;
  let isAdmin = false;
  let sendEmail = true;
  let loading = false;
  let generatedLink = '';
  let showResult = false;
  
  // Validation
  $: emailValid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  $: daysValid = expirationDays >= 1 && expirationDays <= 30;
  $: formValid = emailValid && daysValid;
  
  // Reset form when modal opens/closes
  $: if (open) {
    email = prefilledEmail;
    expirationDays = 7;
    isAdmin = false;
    sendEmail = true;
    loading = false;
    generatedLink = '';
    showResult = false;
  }
  
  async function handleGenerateInvite() {
    if (!formValid) return;
    
    loading = true;
    
    try {
      const options: InviteOptions = {
        email,
        expirationDays,
        isAdmin,
        sendEmail
      };
      
      const link = await registrationActions.generateInvite(options);
      generatedLink = link;
      showResult = true;
      
      dispatch('inviteGenerated', { email, link });
      
      if (sendEmail) {
        toast.success(`Invite sent to ${email}`);
      } else {
        toast.success('Invite link generated successfully');
      }
      
    } catch (error) {
      toast.error('Failed to generate invite');
    } finally {
      loading = false;
    }
  }
  
  function handleClose() {
    open = false;
    dispatch('close');
  }
  
  function copyToClipboard() {
    navigator.clipboard.writeText(generatedLink);
    toast.success('Link copied to clipboard');
  }
  
  function handleNewInvite() {
    showResult = false;
    email = '';
    generatedLink = '';
  }
</script>

<Dialog bind:open>
  <DialogContent class="sm:max-w-md">
    <DialogHeader>
      <DialogTitle>
        {showResult ? 'Invite Generated' : 'Generate User Invite'}
      </DialogTitle>
      <DialogDescription>
        {showResult 
          ? 'Your invite link has been generated successfully.'
          : 'Create an invitation link for a new user to join your business.'
        }
      </DialogDescription>
    </DialogHeader>
    
    {#if showResult}
      <!-- Success Result -->
      <div class="space-y-4">
        <div class="p-4 bg-green-50 border border-green-200 rounded-lg">
          <div class="flex items-center space-x-2 mb-2">
            <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span class="font-medium text-green-800">Invite Created Successfully</span>
          </div>
          <p class="text-sm text-green-700">
            {sendEmail 
              ? `An invitation email has been sent to ${email}`
              : 'Share the link below with the user'
            }
          </p>
        </div>
        
        <div class="space-y-2">
          <Label>Invitation Link</Label>
          <div class="flex space-x-2">
            <Input 
              value={generatedLink} 
              readonly 
              class="font-mono text-sm"
            />
            <Button variant="outline" size="sm" onclick={copyToClipboard}>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </Button>
          </div>
        </div>
        
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-muted-foreground">Email:</span>
            <div class="font-medium">{email}</div>
          </div>
          <div>
            <span class="text-muted-foreground">Expires:</span>
            <div class="font-medium">{expirationDays} days</div>
          </div>
          <div>
            <span class="text-muted-foreground">Admin Access:</span>
            <Badge variant={isAdmin ? 'default' : 'secondary'}>
              {isAdmin ? 'Yes' : 'No'}
            </Badge>
          </div>
          <div>
            <span class="text-muted-foreground">Email Sent:</span>
            <Badge variant={sendEmail ? 'default' : 'secondary'}>
              {sendEmail ? 'Yes' : 'No'}
            </Badge>
          </div>
        </div>
      </div>
      
      <DialogFooter class="flex-col sm:flex-row gap-2">
        <Button variant="outline" onclick={handleNewInvite} class="w-full sm:w-auto">
          Generate Another
        </Button>
        <Button onclick={handleClose} class="w-full sm:w-auto">
          Done
        </Button>
      </DialogFooter>
    {:else}
      <!-- Invite Form -->
      <div class="space-y-4">
        <div class="space-y-2">
          <Label for="email">Email Address *</Label>
          <Input
            id="email"
            type="email"
            placeholder="user@example.com"
            bind:value={email}
            class={emailValid || !email ? '' : 'border-red-500'}
          />
          {#if email && !emailValid}
            <p class="text-sm text-red-600">Please enter a valid email address</p>
          {/if}
        </div>
        
        <div class="space-y-2">
          <Label for="expiration">Expiration (days) *</Label>
          <Input
            id="expiration"
            type="number"
            min="1"
            max="30"
            bind:value={expirationDays}
            class={daysValid ? '' : 'border-red-500'}
          />
          {#if !daysValid}
            <p class="text-sm text-red-600">Must be between 1 and 30 days</p>
          {/if}
          <p class="text-sm text-muted-foreground">
            Link will expire in {expirationDays} day{expirationDays !== 1 ? 's' : ''}
          </p>
        </div>
        
        <div class="space-y-3">
          <div class="flex items-center space-x-2">
            <Checkbox 
              id="admin"
              bind:checked={isAdmin}
            />
            <Label for="admin" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              Grant admin permissions
            </Label>
          </div>
          {#if isAdmin}
            <p class="text-sm text-orange-600 ml-6">
              ⚠️ Admin users can approve/deny registrations and manage other users
            </p>
          {/if}
          
          <div class="flex items-center space-x-2">
            <Checkbox 
              id="sendEmail"
              bind:checked={sendEmail}
            />
            <Label for="sendEmail" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              Send invitation email automatically
            </Label>
          </div>
          {#if !sendEmail}
            <p class="text-sm text-muted-foreground ml-6">
              You'll need to share the generated link manually
            </p>
          {/if}
        </div>
        
        <!-- Preview -->
        <div class="p-3 bg-muted rounded-lg">
          <h4 class="text-sm font-medium mb-2">Invitation Preview:</h4>
          <div class="text-sm space-y-1">
            <div><span class="text-muted-foreground">To:</span> {email || 'user@example.com'}</div>
            <div><span class="text-muted-foreground">Access Level:</span> {isAdmin ? 'Admin' : 'Standard User'}</div>
            <div><span class="text-muted-foreground">Valid for:</span> {expirationDays} days</div>
            <div><span class="text-muted-foreground">Delivery:</span> {sendEmail ? 'Email' : 'Manual link sharing'}</div>
          </div>
        </div>
      </div>
      
      <DialogFooter>
        <Button variant="outline" onclick={handleClose}>
          Cancel
        </Button>
        <Button 
          onclick={handleGenerateInvite}
          disabled={!formValid || loading}
        >
          {#if loading}
            <svg class="w-4 h-4 mr-2 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Generating...
          {:else}
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Generate Invite
          {/if}
        </Button>
      </DialogFooter>
    {/if}
  </DialogContent>
</Dialog>