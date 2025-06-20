export function load({ url }) {
  const userId = url.searchParams.get("user");
  return { userId };
}
