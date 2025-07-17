import type { PageLoad } from './$types';
import { getUserCurrent } from '$lib/api/auth.js';

export const load: PageLoad = async ({ fetch, url }) => {
  try {
    // Get current user to verify admin access
    const user = await getUserCurrent(fetch);
    
    // Extract business ID from URL params or user context
    // For now, we'll use a placeholder - this should be adjusted based on your auth system
    const businessId = url.searchParams.get('businessId') || 'default-business-id';
    
    return {
      user,
      businessId
    };
  } catch (error) {
    // If user is not authenticated or doesn't have admin access, redirect
    throw new Error('Unauthorized access to admin panel');
  }
};