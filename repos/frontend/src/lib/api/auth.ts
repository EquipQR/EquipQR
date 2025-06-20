export async function loginUser(loginEmail: string, loginPassword: string) {
  const res = await fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email: loginEmail,
      password: loginPassword,
    }),
  });
  const data = await res.json();
  console.log(data);
}

export async function registerUser(
  registerUsername: string,
  registerEmail: string,
  registerPassword: string
) {
  const res = await fetch("/api/user", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      username: registerUsername,
      email: registerEmail,
      password: registerPassword,
    }),
  });
  const data = await res.json();
  console.log(data);
}
