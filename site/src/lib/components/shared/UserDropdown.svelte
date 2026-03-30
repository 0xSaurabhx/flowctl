<script lang="ts">
  import { currentUser } from '$lib/stores/auth';
  import { apiClient } from '$lib/apiClient';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import { clearPermissionCache } from '$lib/utils/permissions';
  import { IconChevronDown } from '@tabler/icons-svelte';

  let { isCollapsed = false }: { isCollapsed?: boolean } = $props();

  let userSettingsOpen = $state(false);

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
      // Force redirect even if logout fails
      window.location.href = '/login';
    }
  };

  // Handle outside clicks
  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.user-dropdown-container')) {
      userSettingsOpen = false;
    }
  }
</script>

<svelte:window on:click={handleOutsideClick} />

<!-- User Menu -->
<div class="user-dropdown-container">
  <button
    type="button"
    onclick={() => userSettingsOpen = !userSettingsOpen}
    class="user-btn"
    class:icon-only={isCollapsed}
    aria-label="User menu toggle"
    title={isCollapsed ? $currentUser?.name || 'User menu' : ''}
  >
    <div class="avatar">
      <span>{$currentUser ? getUserInitials($currentUser.name) : 'U'}</span>
    </div>
    {#if !isCollapsed}
      <div class="user-info">
        <div class="user-name">{$currentUser?.name || 'Loading...'}</div>
        <div class="user-role text-lighter">{$currentUser?.role || ''}</div>
      </div>
      <IconChevronDown
        class="chevron {userSettingsOpen ? 'open' : ''}"
        size={16}
      />
    {/if}
  </button>

  <!-- Dropdown Menu -->
  {#if userSettingsOpen}
    <div
      class="user-menu"
      class:narrow={isCollapsed}
      role="menu"
      aria-label="User menu"
    >
      <button
        type="button"
        onclick={logout}
        role="menuitem"
      >
        Logout
      </button>
    </div>
  {/if}
</div>

<style>
  .user-dropdown-container {
    position: relative;
  }

  .user-btn {
    width: 100%;
    display: flex;
    align-items: center;
    font-size: var(--text-7);
    font-weight: var(--font-medium);
    color: var(--foreground);
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius-medium);
    padding: var(--space-2) var(--space-3);
    cursor: pointer;
    transition: background 0.15s;
  }

  .user-btn:hover {
    background: var(--faint);
  }

  .user-btn.icon-only {
    justify-content: center;
    padding: var(--space-2);
  }

  .avatar {
    width: 2rem;
    height: 2rem;
    background: var(--primary);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: white;
    font-weight: 600;
    font-size: var(--text-7);
  }

  .user-info {
    margin-left: var(--space-3);
    flex: 1;
    text-align: left;
  }

  .user-name {
    font-size: var(--text-7);
    font-weight: var(--font-medium);
    color: var(--foreground);
  }

  .user-role {
    font-size: var(--text-8);
    text-transform: capitalize;
  }

  :global(.chevron) {
    color: var(--muted-foreground);
    transition: transform 0.2s;
    flex-shrink: 0;
  }

  :global(.chevron.open) {
    transform: rotate(180deg);
  }

  .user-menu {
    position: absolute;
    bottom: 100%;
    left: 0;
    width: 100%;
    margin-bottom: var(--space-1);
    background: var(--card);
    border-radius: var(--radius-medium);
    box-shadow: var(--shadow-large);
    border: 1px solid var(--border);
    padding: var(--space-1) 0;
  }

  .user-menu.narrow {
    width: 8rem;
  }

  .user-menu button {
    width: 100%;
    text-align: left;
    padding: var(--space-2) var(--space-3);
    font-size: var(--text-7);
    color: var(--foreground);
    background: none;
    border: none;
    border-radius: 0;
    cursor: pointer;
  }

  .user-menu button:hover {
    background: var(--faint);
  }
</style>
