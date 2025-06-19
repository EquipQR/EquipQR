import type { Issue } from "./types/issue";

export let mockIssues: Issue[] = [
  {
    id: "1",
    created_at: "2024-12-01T14:33:00Z",
    status: "resolved",
    description: "Fuel line replaced due to corrosion.",
  },
  {
    id: "2",
    created_at: "2025-01-15T09:20:00Z",
    status: "open",
    description: "Unusual vibration detected during operation.",
  },
];
