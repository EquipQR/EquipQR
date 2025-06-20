export function load({ url }) {
  return {
    equipmentId: url.searchParams.get("scanned"),
  };
}
