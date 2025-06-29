// Admin API functions for registration management

export async function getPendingRegistrations(businessId: string) {
  const response = await fetch(`/api/pending/${businessId}`, {
    credentials: 'include'
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch pending registrations');
  }
  
  return response.json();
}

export async function approveRegistration(requestId: string) {
  const response = await fetch('/api/pending/approve', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ request_id: requestId })
  });
  
  if (!response.ok) {
    throw new Error('Failed to approve registration');
  }
  
  return response.status === 204;
}

export async function denyRegistration(requestId: string) {
  const response = await fetch('/api/pending/deny', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ request_id: requestId })
  });
  
  if (!response.ok) {
    throw new Error('Failed to deny registration');
  }
  
  return response.status === 204;
}

export async function bulkAction(action: {
  type: 'approve' | 'deny' | 'generate_invite';
  registration_ids: string[];
  admin_permissions?: boolean;
}) {
  const response = await fetch('/api/pending/bulk-action', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(action)
  });
  
  if (!response.ok) {
    throw new Error(`Failed to perform bulk ${action.type}`);
  }
  
  return response.json();
}

export async function generateInvite(options: {
  email: string;
  expirationDays: number;
  isAdmin: boolean;
  sendEmail: boolean;
}) {
  const response = await fetch('/api/invites/generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(options)
  });
  
  if (!response.ok) {
    throw new Error('Failed to generate invite');
  }
  
  const data = await response.json();
  return data.invite_link;
}