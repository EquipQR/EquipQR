<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Label } from "$lib/components/ui/label";
  import { getNames } from "country-list";
  import { registerBusiness } from "$lib/api/portal";

  let username = "";
  let password = "";
  let companyName = "";
  let businessEmail = "";
  let phoneNumber = "";
  let industry = "";
  let companySize = "";
  let country = "";

  const industries = [
    "Aviation",
    "Transport",
    "Mechanics",
    "Factory",
    "Manufacturing",
    "Construction",
    "Mining",
    "Oil & Gas",
    "Marine",
    "Automotive",
    "Healthcare",
    "Hospitality",
    "Other",
  ];
  const companySizes = [
    "1-10 employees",
    "11-50 employees",
    "51-200 employees",
    "201-500 employees",
    "500+ employees",
  ];
  const countries: string[] = [...getNames(), "Other"];
</script>

<div class="dark min-h-screen flex items-center justify-center bg-black px-4 py-12 text-white">
  <div class="w-full max-w-3xl bg-black rounded-xl shadow-2xl p-8 space-y-8 border border-neutral-800">
    <div class="text-center space-y-1">
      <img src="/app_logo_512_412.png" alt="EquipQR Logo" class="mx-auto h-28 w-auto" />
      <h1 class="text-3xl font-bold tracking-tight">Get Started with EquipQR</h1>
      <p class="text-neutral-500 text-sm">
        Set up your business account in minutes and unlock early access features.
      </p>
    </div>

    <form
      on:submit|preventDefault={() =>
        registerBusiness(
          username,
          businessEmail,
          password,
          companyName,
          businessEmail,
          phoneNumber,
          industry,
          companySize,
          country
        )
      }
      class="space-y-6"
    >
      <div class="grid md:grid-cols-2 gap-6">
        <div>
          <Label for="username">Username</Label>
          <Input
            id="username"
            type="text"
            bind:value={username}
            required
            placeholder="johndoe"
            class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
          />
        </div>
        <div>
          <Label for="password">Password</Label>
          <Input
            id="password"
            type="password"
            bind:value={password}
            required
            placeholder="••••••••"
            class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
          />
        </div>
      </div>

      <div class="grid md:grid-cols-2 gap-6">
        <div>
          <Label for="companyName">Company Name</Label>
          <Input
            id="companyName"
            type="text"
            bind:value={companyName}
            placeholder="Acme Corp"
            class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
          />
        </div>
        <div>
          <Label for="businessEmail">Business Email</Label>
          <Input
            id="businessEmail"
            type="email"
            bind:value={businessEmail}
            placeholder="admin@company.com"
            class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
          />
        </div>
      </div>

      <div>
        <Label for="phoneNumber">Phone Number</Label>
        <Input
          id="phoneNumber"
          type="tel"
          bind:value={phoneNumber}
          placeholder="+1 555 123 4567"
          class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
        />
      </div>

      <div class="grid md:grid-cols-3 gap-6">
        <div>
          <Label for="industry">Industry</Label>
          <select
            id="industry"
            bind:value={industry}
            class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm rounded-md px-4 py-2.5"
          >
            <option value="" disabled selected>Select industry</option>
            {#each industries as option}
              <option value={option}>{option}</option>
            {/each}
          </select>
        </div>
        <div>
          <Label for="companySize">Company Size</Label>
          <select
            id="companySize"
            bind:value={companySize}
            class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm rounded-md px-4 py-2.5"
          >
            <option value="" disabled selected>Select size</option>
            {#each companySizes as option}
              <option value={option}>{option}</option>
            {/each}
          </select>
        </div>
        <div>
          <Label for="country">Country</Label>
          <select
            id="country"
            bind:value={country}
            class="w-full bg-neutral-900 border border-neutral-800 text-white text-sm rounded-md px-4 py-2.5"
          >
            <option value="" disabled selected>Select country</option>
            {#each countries as option}
              <option value={option}>{option}</option>
            {/each}
          </select>
        </div>
      </div>

      <div class="pt-6 border-t border-neutral-800">
        <h3 class="text-lg font-semibold mb-2">What EquipQR Offers</h3>
        <ul class="list-disc pl-5 text-sm text-neutral-400 space-y-1">
          <li>Unlimited QR codes for your equipment</li>
          <li>Error reporting with photo/video support</li>
          <li>Live dashboard & alerts</li>
          <li>Streamlined maintenance team communication</li>
          <li>Custom workflows per industry</li>
        </ul>
      </div>

      <div class="pt-4">
        <Button type="submit" class="w-full transition duration-150 active:scale-95">
          Join Early Access
        </Button>
      </div>

      <p class="text-xs text-neutral-500 text-center mt-4">
        By joining, you agree to our
        <a href="/terms" class="underline">Terms</a> and
        <a href="/privacy" class="underline">Privacy Policy</a>.
      </p>
    </form>
  </div>
</div>
