export function load({ url }) {
  const businessId = url.searchParams.get("business");
  return { businessId };
}
