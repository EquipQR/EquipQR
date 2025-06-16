import { Equipment } from '../types/equipment';

const API_BASE_URL = 'https://192.168.0.102:8080';

export const fetchEquipment = async (id: string): Promise<Equipment> => {
  try {
    const response = await fetch(`${API_BASE_URL}/equipment/${id}`);
    
    if (!response.ok) {
      throw new Error(`Equipment not found: ${response.status}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch equipment: ${error.message}`);
    }
    throw new Error('Failed to fetch equipment: Unknown error');
  }
};