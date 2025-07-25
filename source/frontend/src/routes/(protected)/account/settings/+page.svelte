<script lang="ts">
  import { onMount } from "svelte";
  import { startRegistration } from "@simplewebauthn/browser";
  import { currentUser } from "$lib/api/user";
  import { getUserCurrent } from "$lib/api/auth";
  import { goto } from "$app/navigation";

  let user: Awaited<ReturnType<typeof getUserCurrent>> | null = null;
  let newImage: File | null = null;
  let imagePreview: string | null = null;
  let webauthnRegistered = false;
  let loading = false;
  let error: string | null = null;

  onMount(async () => {
    try {
      user = await getUserCurrent(fetch);
      currentUser.set(user);
      webauthnRegistered = user?.has_webauthn ?? false;
    } catch (err) {
      error = "Failed to load user";
    }
  });

  async function handleAvatarUpload() {
    if (!newImage) return;
    const form = new FormData();
    form.append("avatar", newImage);

    try {
      await fetch("/api/user/avatar", {
        method: "POST",
        body: form,
      });
      location.reload();
    } catch {
      error = "Upload failed";
    }
  }

async function setupWebAuthn(): Promise<void> {
  loading = true;
  error = null;

  try {
    if (!user?.id) throw new Error("User ID is missing");

    console.log("🟡 Starting WebAuthn registration for user:", user.id);

    const res = await fetch("/api/auth/webauthn/register/begin", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user_id: user.id }),
    });

    if (!res.ok) {
      const msg = await res.text();
      console.error("❌ Failed to begin registration:", msg);
      throw new Error("Failed to begin WebAuthn registration");
    }

    const { publicKey } = await res.json();
    console.log("🛠️ Received registration options:", publicKey);

    if (!publicKey || !publicKey.challenge) {
      throw new Error("Invalid WebAuthn registration options: missing challenge");
    }

    const credential = await startRegistration(publicKey);
    console.log("✅ Got registration credential:", credential);

    const finishRes = await fetch("/api/auth/webauthn/register/finish", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(credential),
    });

    if (!finishRes.ok) {
      const msg = await finishRes.text();
      console.error("❌ Registration finalization failed:", msg);
      throw new Error("WebAuthn registration failed");
    }

    console.log("🎉 WebAuthn registration complete");
    webauthnRegistered = true;
  } catch (err) {
    console.error("❌ setupWebAuthn failed:", err);
    error = (err as Error).message;
  } finally {
    loading = false;
  }
}
  function onImageChange(event: Event) {
    const fileInput = event.target as HTMLInputElement;
    if (fileInput?.files?.[0]) {
      newImage = fileInput.files[0];
      imagePreview = URL.createObjectURL(newImage);
    }
  }
</script>

<h1 class="text-2xl font-bold mb-6">Account Settings</h1>

{#if error}
  <p class="text-red-500">{error}</p>
{/if}

<section class="mb-6">
  <h2 class="text-xl font-semibold mb-2">Profile Picture</h2>
  <img
    src={imagePreview || user?.avatar_url || "/default-avatar.png"}
    alt="Profile"
    class="w-24 h-24 rounded-full mb-2 object-cover border"
  />
  <input type="file" accept="image/*" on:change={onImageChange} />
  {#if newImage}
    <button class="btn btn-primary mt-2" on:click={handleAvatarUpload}
      >Upload</button
    >
  {/if}
</section>

<section class="mb-6">
  <h2 class="text-xl font-semibold mb-2">WebAuthn Security Key</h2>
  {#if webauthnRegistered}
    <p class="text-green-600">Registered ✅</p>
  {:else}
    <button
      class="btn btn-secondary"
      on:click={setupWebAuthn}
      disabled={loading}
    >
      {loading ? "Setting up..." : "Register Security Key"}
    </button>
  {/if}
</section>
