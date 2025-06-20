import { goto } from "$app/navigation";

export async function loginUser(loginEmail: string, loginPassword: string): Promise<void> {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: loginEmail, password: loginPassword }),
  });

  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.message || "Login failed");
  }

  goto("/");
}

export async function registerUser(
  registerUsername: string,
  registerEmail: string,
  registerPassword: string
): Promise<void> {
  const res = await fetch("/api/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
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

  goto("/");
}
