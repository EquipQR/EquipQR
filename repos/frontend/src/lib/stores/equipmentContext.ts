import { writable } from 'svelte/store';

export const scannedEquipmentId = writable<string | null>(null);
