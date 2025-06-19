<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { Card, CardContent } from "$lib/components/ui/card";
	import { Avatar, AvatarFallback, AvatarImage } from "$lib/components/ui/avatar";
	import { Sheet, SheetContent, SheetTrigger } from "$lib/components/ui/sheet";
	import { Menu } from "lucide-svelte";

	const userName = "Justin";
	const role = "Admin";
	const accountType = "Master Account";

	const navItems: { label: string; href: string }[] = [
		{ label: "Assets", href: "/assets" },
		{ label: "Issues", href: "/issues" },
		{ label: "Scan QR", href: "/" }, // TODO: Redirect to /scan
		{ label: "Users", href: "/users" },
		{ label: "Equipment", href: "/equipment" }
	];
</script>

<div class="min-h-screen bg-black text-white font-comic p-4 space-y-6">
	<!-- Top Bar -->
	<div class="flex items-center justify-between">
		<h1 class="text-xl">ExampleTrucks</h1>
		<Sheet>
			<SheetTrigger>
				<Button variant="ghost" size="icon" class="text-white">
					<Menu class="w-5 h-5" />
				</Button>
			</SheetTrigger>
			<SheetContent side="right" class="bg-black text-white w-64 p-6 space-y-4">
				<h2 class="text-lg font-semibold">Navigation</h2>
				{#each navItems as item}
					<a href={item.href} class="block w-full">
						<Button variant="outline" class="w-full border-white text-white">{item.label}</Button>
					</a>
				{/each}
			</SheetContent>
		</Sheet>
	</div>

	<!-- Profile -->
	<div class="flex items-center gap-4">
		<Avatar>
			<AvatarImage src="https://github.com/shadcn.png" alt="@justin" />
			<AvatarFallback>JP</AvatarFallback>
		</Avatar>
		<div class="leading-tight">
			<p class="text-base">Welcome, {userName}!</p>
			<p class="text-sm text-white/70">{role} • {accountType}</p>
		</div>
	</div>

	<!-- Navigation Grid -->
	<div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
		{#each navItems as item}
			<Card class="bg-transparent border-white text-white hover:bg-white/5 transition">
				<CardContent class="flex items-center justify-center h-20">
					<a href={item.href} class="text-center w-full text-md font-normal">{item.label}</a>
				</CardContent>
			</Card>
		{/each}
	</div>
</div>
