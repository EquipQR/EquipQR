import type { Business } from '$lib/types/business';

export type Equipment = {
  id: string;
  businessId: string;
  status: 'in service' | 'not in service';
  type: string;
  location: string;
  moreFields: Record<string, unknown>;
  business: Business;
};

export type EquipmentCreatePayload = {
  businessId: string;
  status: 'in service' | 'not in service';
  type: string;
  location?: string;
  moreFields?: Record<string, unknown>;
};
