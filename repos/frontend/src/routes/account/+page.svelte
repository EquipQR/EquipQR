<script lang="ts">
  import { onMount } from "svelte";

  import * as Tabs from '$lib/components/ui/tabs/index.js';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$lib/components/ui/card';
  onMount(() => {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  document.documentElement.classList.add(prefersDark ? 'dark' : 'light');
});

</script>

<Tabs.Root value="account" class="min-h-screen w-full bg-background text-foreground flex flex-col items-center justify-start pt-12">
  <Tabs.List class="flex space-x-2 bg-muted p-1 rounded-md border border-border mb-8">
    <Tabs.Trigger
      value="account"
      class="px-6 py-2 text-sm font-medium uppercase tracking-wide border border-transparent data-[state=active]:bg-background data-[state=active]:text-foreground rounded-none focus:outline-none"
    >
      Account
    </Tabs.Trigger>
    <Tabs.Trigger
      value="password"
      class="px-6 py-2 text-sm font-medium uppercase tracking-wide border border-transparent data-[state=active]:bg-background data-[state=active]:text-foreground rounded-none focus:outline-none"
    >
      Password
    </Tabs.Trigger>
  </Tabs.List>

  <Tabs.Content value="account" class="w-full max-w-md px-4">
    <Card class="border border-border rounded-md bg-muted/40 shadow-sm">
      <CardHeader>
        <CardTitle class="text-lg font-semibold">Account</CardTitle>
        <CardDescription class="text-sm text-muted-foreground">
          Make changes to your account here. Click save when you're done.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-1">
          <Label for="name" class="text-sm font-medium">Name</Label>
          <Input id="name" value="Pedro Duarte" class="rounded-md" />
        </div>

        <div class="space-y-1">
          <Label for="username" class="text-sm font-medium">Username</Label>
          <Input id="username" value="@peduarte" class="rounded-md" />
        </div>

        <Button class="w-full mt-4 bg-white text-black border border-border hover:bg-neutral-100 rounded-md font-medium">
          Save changes
        </Button>
      </CardContent>
    </Card>
  </Tabs.Content>

  <Tabs.Content value="password" class="w-full max-w-md px-4">
    <Card class="border border-border rounded-md bg-muted/40 shadow-sm">
      <CardHeader>
        <CardTitle class="text-lg font-semibold">Password</CardTitle>
        <CardDescription class="text-sm text-muted-foreground">
          Change your password here.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <!-- Add password change fields here -->
      </CardContent>
    </Card>
  </Tabs.Content>
</Tabs.Root>
