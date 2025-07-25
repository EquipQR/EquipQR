export interface PendingRegistration {
  id: string;
  user_id: string;
  business_id: string;
  created_at: string;
  user: {
    username: string;
    email: string;
    is_active: boolean;
  };
  urgency_level: 'normal' | 'urgent' | 'critical';
  days_pending: number;
}

export interface InviteLink {
  id: string;
  email: string;
  expires_at: string;
  is_admin: boolean;
  status: 'pending' | 'used' | 'expired' | 'revoked';
  created_at: string;
  used_at?: string;
}

export interface BulkAction {
  type: 'approve' | 'deny' | 'generate_invite';
  registration_ids: string[];
  admin_permissions?: boolean;
}

export interface FilterState {
  search: string;
  email: string;
  dateRange: {
    start: string | null;
    end: string | null;
  };
  urgencyLevel: 'all' | 'normal' | 'urgent' | 'critical';
  status: 'all' | 'pending' | 'approved' | 'denied';
}

export interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface SortState {
  column: 'created_at' | 'email' | 'username' | 'days_pending';
  direction: 'asc' | 'desc';
}

export interface AdminStats {
  totalPending: number;
  urgentCount: number;
  todayCount: number;
  weeklyCount: number;
}

export interface InviteOptions {
  email: string;
  expirationDays: number;
  isAdmin: boolean;
  sendEmail: boolean;
}

export interface RegistrationUpdate {
  type: 'new' | 'approved' | 'denied' | 'updated';
  registration: PendingRegistration;
  timestamp: string;
}

export interface AdminPermission {
  level: 'view' | 'standard' | 'super' | 'owner';
  canApprove: boolean;
  canDeny: boolean;
  canGenerateInvites: boolean;
  canManageAdmins: boolean;
  canBulkActions: boolean;
}