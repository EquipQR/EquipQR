// src/stores/user.ts
import { writable } from "svelte/store";

export const currentUser = writable<{
  id: string;
  email: string;
  username: string;
} | null>(null);
