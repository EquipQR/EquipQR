import type { Equipment } from "$lib/types/equipment";

export async function getEquipmentById(id: string): Promise<Equipment | null> {
  const res = await fetch(`/api/equipment/${id}`);
  console.log(res)
  if (!res.ok) {
    console.error(`Failed to fetch equipment with id ${id}`);
    return null;
  }
  return res.json();
}
