import type { Equipment } from "$lib/types/equipment";
import type { Issue } from "$lib/types/issue";

export async function getEquipmentById(id: string): Promise<Equipment | null> {
  const res = await fetch(`/api/equipment/${id}`);
  if (!res.ok) {
    console.error(`Failed to fetch equipment with id ${id}`);
    return null;
  }
  return res.json();
}

export async function getEquipmentIssuesById(
  id: string
): Promise<Issue[] | null> {
  const res = await fetch(`/api/equipment/${id}/issues`);
  if (!res.ok) {
    console.error(`Failed to fetch issues for equipment with id ${id}`);
    return null;
  }
  return res.json();
}
