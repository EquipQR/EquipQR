import { writable, derived, get } from 'svelte/store';
import type { 
  PendingRegistration, 
  FilterState, 
  PaginationState, 
  SortState, 
  BulkAction, 
  InviteOptions,
  AdminStats 
} from '../types/admin.js';

// Core data stores
export const registrations = writable<PendingRegistration[]>([]);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Filter and pagination stores
export const filters = writable<FilterState>({
  search: '',
  email: '',
  dateRange: { start: null, end: null },
  urgencyLevel: 'all',
  status: 'all'
});

export const pagination = writable<PaginationState>({
  page: 1,
  pageSize: 20,
  total: 0,
  totalPages: 0
});

export const sort = writable<SortState>({
  column: 'created_at',
  direction: 'desc'
});

// Selection store for bulk actions
export const selectedIds = writable<string[]>([]);

// Loading states for individual actions
export const actionLoading = writable<Record<string, boolean>>({});

// Derived store for filtered and sorted registrations
export const filteredRegistrations = derived(
  [registrations, filters, sort],
  ([$registrations, $filters, $sort]) => {
    let filtered = [...$registrations];

    // Apply search filter
    if ($filters.search) {
      const searchLower = $filters.search.toLowerCase();
      filtered = filtered.filter(reg => 
        reg.user.username.toLowerCase().includes(searchLower) ||
        reg.user.email.toLowerCase().includes(searchLower)
      );
    }

    // Apply email filter
    if ($filters.email) {
      const emailLower = $filters.email.toLowerCase();
      filtered = filtered.filter(reg => 
        reg.user.email.toLowerCase().includes(emailLower)
      );
    }

    // Apply urgency filter
    if ($filters.urgencyLevel !== 'all') {
      filtered = filtered.filter(reg => reg.urgency_level === $filters.urgencyLevel);
    }

    // Apply date range filter
    if ($filters.dateRange.start) {
      const startDate = new Date($filters.dateRange.start);
      filtered = filtered.filter(reg => new Date(reg.created_at) >= startDate);
    }
    if ($filters.dateRange.end) {
      const endDate = new Date($filters.dateRange.end);
      filtered = filtered.filter(reg => new Date(reg.created_at) <= endDate);
    }

    // Apply sorting
    filtered.sort((a, b) => {
      let aValue: any, bValue: any;
      
      switch ($sort.column) {
        case 'created_at':
          aValue = new Date(a.created_at);
          bValue = new Date(b.created_at);
          break;
        case 'email':
          aValue = a.user.email;
          bValue = b.user.email;
          break;
        case 'username':
          aValue = a.user.username;
          bValue = b.user.username;
          break;
        case 'days_pending':
          aValue = a.days_pending;
          bValue = b.days_pending;
          break;
        default:
          return 0;
      }

      if (aValue < bValue) return $sort.direction === 'asc' ? -1 : 1;
      if (aValue > bValue) return $sort.direction === 'asc' ? 1 : -1;
      return 0;
    });

    return filtered;
  }
);

// Derived store for paginated results
export const paginatedRegistrations = derived(
  [filteredRegistrations, pagination],
  ([$filteredRegistrations, $pagination]) => {
    const start = ($pagination.page - 1) * $pagination.pageSize;
    const end = start + $pagination.pageSize;
    return $filteredRegistrations.slice(start, end);
  }
);

// Derived store for admin stats
export const adminStats = derived(
  registrations,
  ($registrations) => {
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const weekAgo = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);

    const stats: AdminStats = {
      totalPending: $registrations.length,
      urgentCount: $registrations.filter(r => r.urgency_level === 'urgent' || r.urgency_level === 'critical').length,
      todayCount: $registrations.filter(r => new Date(r.created_at) >= today).length,
      weeklyCount: $registrations.filter(r => new Date(r.created_at) >= weekAgo).length
    };

    return stats;
  }
);

// Actions
export const registrationActions = {
  async loadRegistrations(businessId: string) {
    loading.set(true);
    error.set(null);
    
    try {
      const response = await fetch(`/api/pending/${businessId}`, {
        credentials: 'include'
      });
      
      if (!response.ok) {
        throw new Error('Failed to load registrations');
      }
      
      const data = await response.json();
      
      // Calculate urgency and days pending
      const processedData = data.map((reg: any) => ({
        ...reg,
        days_pending: Math.floor((Date.now() - new Date(reg.created_at).getTime()) / (1000 * 60 * 60 * 24)),
        urgency_level: calculateUrgencyLevel(reg.created_at)
      }));
      
      registrations.set(processedData);
      
      // Update pagination
      const currentPagination = get(pagination);
      const totalPages = Math.ceil(processedData.length / currentPagination.pageSize);
      pagination.update(p => ({
        ...p,
        total: processedData.length,
        totalPages
      }));
      
    } catch (err) {
      error.set(err instanceof Error ? err.message : 'Unknown error occurred');
    } finally {
      loading.set(false);
    }
  },

  async approveRegistration(registrationId: string) {
    actionLoading.update(state => ({ ...state, [registrationId]: true }));
    
    try {
      const response = await fetch('/api/pending/approve', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ request_id: registrationId })
      });
      
      if (!response.ok) {
        throw new Error('Failed to approve registration');
      }
      
      // Remove from local state
      registrations.update(regs => regs.filter(r => r.id !== registrationId));
      selectedIds.update(ids => ids.filter(id => id !== registrationId));
      
    } catch (err) {
      error.set(err instanceof Error ? err.message : 'Failed to approve registration');
    } finally {
      actionLoading.update(state => ({ ...state, [registrationId]: false }));
    }
  },

  async denyRegistration(registrationId: string) {
    actionLoading.update(state => ({ ...state, [registrationId]: true }));
    
    try {
      const response = await fetch('/api/pending/deny', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ request_id: registrationId })
      });
      
      if (!response.ok) {
        throw new Error('Failed to deny registration');
      }
      
      // Remove from local state
      registrations.update(regs => regs.filter(r => r.id !== registrationId));
      selectedIds.update(ids => ids.filter(id => id !== registrationId));
      
    } catch (err) {
      error.set(err instanceof Error ? err.message : 'Failed to deny registration');
    } finally {
      actionLoading.update(state => ({ ...state, [registrationId]: false }));
    }
  },

  async bulkAction(action: BulkAction) {
    loading.set(true);
    
    try {
      const response = await fetch('/api/pending/bulk-action', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(action)
      });
      
      if (!response.ok) {
        throw new Error(`Failed to perform bulk ${action.type}`);
      }
      
      // Remove processed registrations from local state
      registrations.update(regs => 
        regs.filter(r => !action.registration_ids.includes(r.id))
      );
      selectedIds.set([]);
      
    } catch (err) {
      error.set(err instanceof Error ? err.message : `Failed to perform bulk ${action.type}`);
    } finally {
      loading.set(false);
    }
  },

  async generateInvite(options: InviteOptions): Promise<string> {
    try {
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
      
    } catch (err) {
      error.set(err instanceof Error ? err.message : 'Failed to generate invite');
      throw err;
    }
  },

  // Utility functions
  toggleSelection(registrationId: string) {
    selectedIds.update(ids => {
      if (ids.includes(registrationId)) {
        return ids.filter(id => id !== registrationId);
      } else {
        return [...ids, registrationId];
      }
    });
  },

  selectAll() {
    const currentRegistrations = get(paginatedRegistrations);
    selectedIds.set(currentRegistrations.map(r => r.id));
  },

  clearSelection() {
    selectedIds.set([]);
  },

  updateFilters(newFilters: Partial<FilterState>) {
    filters.update(current => ({ ...current, ...newFilters }));
    // Reset to first page when filters change
    pagination.update(p => ({ ...p, page: 1 }));
  },

  updateSort(column: SortState['column']) {
    sort.update(current => ({
      column,
      direction: current.column === column && current.direction === 'asc' ? 'desc' : 'asc'
    }));
  },

  updatePagination(updates: Partial<PaginationState>) {
    pagination.update(current => ({ ...current, ...updates }));
  }
};

// Helper function to calculate urgency level
function calculateUrgencyLevel(createdAt: string): 'normal' | 'urgent' | 'critical' {
  const daysPending = Math.floor((Date.now() - new Date(createdAt).getTime()) / (1000 * 60 * 60 * 24));
  
  if (daysPending >= 7) return 'critical';
  if (daysPending >= 3) return 'urgent';
  return 'normal';
}