export const formatKey = (key: string): string =>
  key.replace(/_/g, " ").replace(/\b\w/g, (char) => char.toUpperCase());
