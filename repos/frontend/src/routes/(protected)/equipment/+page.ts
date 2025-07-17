export function load({ url }) {
  const equipmentId = url.searchParams.get("scanned");
  return { equipmentId };
}
