<script lang="ts">
  import { onMount } from "svelte";
  import {
    Tabs,
    TabsList,
    TabsTrigger,
    TabsContent,
  } from "$lib/components/ui/tabs";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Label } from "$lib/components/ui/label";
  import { Separator } from "$lib/components/ui/separator";
  import { loginUser, registerUser } from "$lib/api/auth";

  let tab: "login" | "register" = "login";

  let loginEmail = "";
  let loginPassword = "";

  let registerUsername = "";
  let registerEmail = "";
  let registerPassword = "";
  let agreedToTOS = false;

  let loginStep: "email" | "loading" | "password" = "email";

  let loginError = "";
  let registerError = "";

  let altLoginMethods: Array<{ id: string; label: string }> = [
    { id: "magic_link", label: "Sign in with Magic Link" },
    { id: "hardware_token", label: "Sign in with Hardware Token" },
    { id: "sso_internal", label: "Internal SSO Login" },
    { id: "ldap", label: "LDAP Credentials" },
  ];

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  const specialCharRegex = /[!@#$%^&*(),.?":{}|<>]/;

  function validateEmail(email: string): boolean {
    return emailRegex.test(email);
  }

  function validatePassword(password: string): boolean {
    return password.length >= 8 && specialCharRegex.test(password);
  }

  async function continueToPassword(): Promise<void> {
    loginError = "";
    if (!loginEmail.trim()) {
      loginError = "Email is required.";
      return;
    }
    if (!validateEmail(loginEmail)) {
      loginError = "Invalid email address.";
      return;
    }
    loginStep = "loading";
    await new Promise((r) => setTimeout(r, 200));
    loginStep = "password";
  }

  async function attemptLogin(): Promise<void> {
    loginError = "";
    if (!loginPassword) {
      loginError = "Password is required.";
      return;
    }
    if (!validatePassword(loginPassword)) {
      loginError =
        "Password must be at least 8 characters and include a special character.";
      return;
    }

    try {
      await loginUser(loginEmail, loginPassword);
    } catch (err) {
      loginError = (err as Error).message;
    }
  }

  async function attemptRegister(): Promise<void> {
    registerError = "";
    if (!registerUsername || registerUsername.length < 3) {
      registerError = "Username must be at least 3 characters long.";
      return;
    }
    if (!validateEmail(registerEmail)) {
      registerError = "Please enter a valid email address.";
      return;
    }
    if (!validatePassword(registerPassword)) {
      registerError =
        "Password must be at least 8 characters and include a special character.";
      return;
    }
    if (!agreedToTOS) {
      registerError = "You must agree to the Terms of Service.";
      return;
    }

    try {
      await registerUser(registerUsername, registerEmail, registerPassword);
    } catch (err) {
      registerError = (err as Error).message;
    }
  }

  async function handleLoginKeydown(event: KeyboardEvent): Promise<void> {
    if (event.key === "Enter") {
      if (loginStep === "email") {
        await continueToPassword();
      } else if (loginStep === "password") {
        await attemptLogin();
      }
    }
  }

  async function handleRegisterKeydown(event: KeyboardEvent): Promise<void> {
    if (event.key === "Enter") {
      await attemptRegister();
    }
  }
</script>

<div
  class="dark min-h-screen flex items-center justify-center bg-black px-4 py-12 text-white"
>
  <div
    class="w-full max-w-md bg-black rounded-xl shadow-2xl p-8 space-y-6 border border-neutral-800"
  >
    <div class="text-center space-y-1">
      <h3 class="text-3xl font-bold">EquipQR</h3>
    </div>

    <Tabs bind:value={tab}>
      <TabsList
        class="w-full grid grid-cols-2 gap-2 bg-neutral-900 rounded-md p-1"
      >
        <TabsTrigger
          value="login"
          class="data-[state=active]:bg-black data-[state=active]:text-white rounded-md py-1"
        >
          Login
        </TabsTrigger>
        <TabsTrigger
          value="register"
          class="data-[state=active]:bg-black data-[state=active]:text-white rounded-md py-1"
        >
          Register
        </TabsTrigger>
      </TabsList>

      <TabsContent value="login" class="space-y-4 pt-4">
        {#if loginStep === "email"}
          <div class="space-y-2">
            <Label for="email">Email</Label>
            <Input
              id="email"
              type="email"
              bind:value={loginEmail}
              oninput={() => (loginError = "")}
              onkeydown={handleLoginKeydown}
              class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
            />
            {#if loginError}<p class="text-sm error-text">{loginError}</p>{/if}
          </div>
          <Button
            class="w-full mt-2 transition duration-150 active:scale-95"
            onclick={continueToPassword}
          >
            Continue
          </Button>
        {:else if loginStep === "loading"}
          <div class="flex justify-center items-center h-20">
            <svg
              class="h-6 w-6 animate-spin text-neutral-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 4v4m0 8v4m4-4h4m-16 0H4m1.64-5.64l2.83 2.83m8.48 0l2.83-2.83M6.34 6.34l2.83 2.83m8.48 0l2.83-2.83"
              />
            </svg>
          </div>
        {:else if loginStep === "password"}
          <div class="space-y-2">
            <Label for="password">Password</Label>
            <Input
              id="password"
              type="password"
              bind:value={loginPassword}
              oninput={() => (loginError = "")}
              onkeydown={handleLoginKeydown}
              class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
            />
            {#if loginError}<p class="text-sm error-text">{loginError}</p>{/if}
          </div>
          <Button
            class="w-full mt-2 transition duration-150 active:scale-95"
            onclick={attemptLogin}
          >
            Sign In
          </Button>
        {/if}

        <Separator class="my-6 border-neutral-800" />
        <p class="text-center text-sm text-neutral-500">
          Choose another way to login
        </p>
        <div class="flex flex-col gap-2">
          {#each altLoginMethods as method}
            <Button
              variant="outline"
              class="w-full border-neutral-800 text-white hover:bg-neutral-900 transition duration-150 active:scale-95"
            >
              {method.label}
            </Button>
          {/each}
        </div>
      </TabsContent>

      <TabsContent value="register" class="space-y-4 pt-4">
        <div class="space-y-2">
          <Label for="register-username">Username</Label>
          <Input
            id="register-username"
            type="text"
            bind:value={registerUsername}
            oninput={() => (registerError = "")}
            onkeydown={handleRegisterKeydown}
            class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
          />
        </div>
        <div class="space-y-2">
          <Label for="register-email">Email</Label>
          <Input
            id="register-email"
            type="email"
            bind:value={registerEmail}
            oninput={() => (registerError = "")}
            onkeydown={handleRegisterKeydown}
            class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
          />
        </div>
        <div class="space-y-2">
          <Label for="register-password">Password</Label>
          <Input
            id="register-password"
            type="password"
            bind:value={registerPassword}
            oninput={() => (registerError = "")}
            onkeydown={handleRegisterKeydown}
            class="bg-neutral-900 border-neutral-800 text-white placeholder-neutral-500"
          />
        </div>
        <div class="flex items-center space-x-2 pt-2">
          <Checkbox id="tos" bind:checked={agreedToTOS} />
          <label for="tos" class="text-sm text-neutral-400">
            I agree to the <a
              href="https://legal.equipqr.io/tos"
              target="_blank"
              class="underline hover:text-white">Terms of Service</a
            >
          </label>
        </div>
        {#if registerError}
          <p class="text-sm error-text">{registerError}</p>
        {/if}
        <Button
          class="w-full mt-2 transition duration-150 active:scale-95"
          onclick={attemptRegister}
          disabled={!agreedToTOS}
        >
          Create Account
        </Button>
      </TabsContent>
    </Tabs>
  </div>
</div>
