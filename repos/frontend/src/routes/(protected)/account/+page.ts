export function load({ url }: { url: URL }): { userId: string | null } {
  const userId = url.searchParams.get("user");
  return { userId };
}
