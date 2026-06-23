import { apiClient } from '$lib/apiClient.js';
import type { LayoutLoad } from './$types';

export const ssr = false; // Disable SSR for the entire app

export const load: LayoutLoad = () => {
  // Return promises without awaiting - layout renders immediately while data loads
  const userPromise = apiClient.users.getProfile().catch((error) => {
    console.error('[Auth] Failed to fetch user profile:', error);
    return null;
  });

  const infoPromise = apiClient.info.get().catch((error) => {
    console.error('[AppInfo] Failed to fetch app info:', error);
    return null;
  });

  return { userPromise, infoPromise };
};