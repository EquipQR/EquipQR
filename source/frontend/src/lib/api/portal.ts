import { goto } from "$app/navigation";

export async function registerBusiness(
  username: string,
  email: string,
  password: string,
  companyName: string,
  businessEmail: string,
  phone: string,
  industry: string,
  companySize: string,
  country: string
): Promise<void> {
  const res = await fetch("/api/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      username,
      email,
      password,
      businessName: companyName,
      businessEmail,
      phone,
      businessType: industry,
      companySize,
      country,
      countryCode: country.slice(0, 2).toUpperCase(),
    }),
    credentials: "include"
  });

  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.error || "Registration failed");
  }

  const data = await res.json();
  console.log("Success:", data);
  await goto("/portal/login");
}
