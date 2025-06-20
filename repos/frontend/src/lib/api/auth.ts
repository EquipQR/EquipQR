import { goto } from "$app/navigation";

export async function loginUser(
  loginEmail: string,
  loginPassword: string
): Promise<void> {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: loginEmail, password: loginPassword }),
  });

  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.message || "Login failed");
  }

  await goto("/");
}

export async function registerUser(
  registerUsername: string,
  registerEmail: string,
  registerPassword: string
): Promise<void> {
  const res = await fetch("/api/auth/register", {
    method: "POST",
    headers: { "Content-Type": "applications/json" },
    body: JSON.stringify({
      username: registerUsername,
      email: registerEmail,
      password: registerPassword,
    }),
  });

  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.message || "Registration failed");
  }

  await goto("/");
}

export async function getUserCurrent(fetchFn: typeof fetch): Promise<{
  id: string;
  username: string;
  email: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}> {
  const res = await fetchFn("/api/user", {
    method: "GET",
    headers: { Accept: "applications/json" },
  });

  if (res.status === 401 || res.status === 403) {
    throw new Error("Unauthorized");
  }

  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.error || "Failed to fetch user");
  }

  const data = await res.json();
  return data.user;
}

export async function logout(): Promise<void> {
  try {
    await fetch("/api/auth/logout", {
      method: "POST",
      credentials: "include",
    });
  } catch (err) {
    console.error("Logout failed:", err);
  } finally {
    await goto("/auth/login");
  }
}
