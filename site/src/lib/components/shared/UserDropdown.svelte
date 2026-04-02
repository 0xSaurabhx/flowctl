<script lang="ts">
  import { currentUser } from '$lib/stores/auth';
  import { apiClient } from '$lib/apiClient';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import { clearPermissionCache } from '$lib/utils/permissions';

  const getUserInitials = (username: string): string => {
    return username.charAt(0).toUpperCase();
  };

  const logout = async () => {
    clearPermissionCache();
    try {
      await apiClient.auth.logout();
      window.location.href = '/login';
    } catch (error) {
      handleInlineError(error, 'Unable to Log Out');
      window.location.href = '/login';
    }
  };
</script>

<ot-dropdown style="display: block" class="mt-2">
  <button popovertarget="user-menu" class="ghost w-100 hstack gap-2" style="justify-content: flex-start">
    <figure data-variant="avatar" class="small" aria-label={$currentUser?.name || 'User'}>
      <abbr title={$currentUser?.name || 'User'}>{$currentUser ? getUserInitials($currentUser.name) : 'U'}</abbr>
    </figure>
    <div style="text-align: left; flex: 1; min-width: 0;">
      <div class="text-sm font-medium">{$currentUser?.name || 'Loading...'}</div>
      <div class="text-lighter text-xs" style="text-transform: capitalize">{$currentUser?.role || ''}</div>
    </div>
  </button>
  <div popover id="user-menu">
    <button role="menuitem" onclick={logout}>Logout</button>
  </div>
</ot-dropdown>

<style>
  figure[data-variant="avatar"] {
    background: var(--primary);
    color: white;
  }
</style>
