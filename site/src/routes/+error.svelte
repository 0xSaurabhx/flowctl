<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth';
  import { getDefaultNamespace } from '$lib/utils/navigation';
  import {
    IconLock,
    IconShieldX,
    IconFileX,
    IconAlertTriangle,
    IconLogin,
    IconHome
  } from '@tabler/icons-svelte';

  const error = page.error;
  const status = page.status;
  const errorCode = error?.code;

  const getErrorDetails = () => {
    if (status === 401) {
      return {
        IconComponent: IconLock,
        iconClass: 'icon-warning',
        bgClass: 'icon-bg-warning',
        title: 'Authentication Required',
        message: error?.message || 'Please log in to access this resource',
        showLoginButton: true
      };
    }

    if (status === 403) {
      return {
        IconComponent: IconShieldX,
        iconClass: 'icon-danger',
        bgClass: 'icon-bg-danger',
        title: 'Access Denied',
        message: error?.message || 'You do not have permission to access this resource',
        showLoginButton: false
      };
    }

    if (status === 404) {
      return {
        IconComponent: IconFileX,
        iconClass: 'icon-muted',
        bgClass: 'icon-bg-muted',
        title: 'Page Not Found',
        message: error?.message || 'The page you are looking for does not exist',
        showLoginButton: false
      };
    }

    return {
      IconComponent: IconAlertTriangle,
      iconClass: 'icon-danger',
      bgClass: 'icon-bg-danger',
      title: 'Something went wrong',
      message: error?.message || 'An unexpected error occurred',
      showLoginButton: false
    };
  };

  const errorDetails = getErrorDetails();

  const handleGoHome = async () => {
    if ($isAuthenticated) {
      const namespace = await getDefaultNamespace();
      goto(`/view/${namespace}/flows`);
    } else {
      goto('/');
    }
  };

  const handleLogin = () => {
    goto('/login');
  };
</script>

<svelte:head>
  <title>Error - Flowctl</title>
</svelte:head>

<div class="error-page">
  <div class="error-content">
    <!-- Error Icon -->
    <div class="mb-8">
      <div class="error-icon-circle {errorDetails.bgClass}">
        <errorDetails.IconComponent class={errorDetails.iconClass} size={48} />
      </div>
    </div>

    <!-- Error Content -->
    <div class="mb-8">
      <h1 class="mb-4">
        {errorDetails.title}
      </h1>
      <p class="text-light mb-2" style="font-size: var(--text-5);">
        {errorDetails.message}
      </p>

      <!-- Show error code if available -->
      {#if errorCode}
        <p class="text-light mt-2 font-mono text-sm">
          Error Code: {errorCode}
        </p>
      {/if}
    </div>

    <!-- Action Buttons -->
    <div class="hstack gap-3 justify-center flex-wrap">
      {#if errorDetails.showLoginButton}
        <button onclick={handleLogin}>
          <span class="hstack gap-2">
            <IconLogin size={18} />
            Log In
          </span>
        </button>
      {/if}

      <button
        onclick={handleGoHome}
        data-variant="secondary"
      >
        <span class="hstack gap-2">
          <IconHome size={18} />
          {$isAuthenticated ? 'Dashboard' : 'Home'}
        </span>
      </button>
    </div>
  </div>
</div>

<style>
  .error-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--muted);
    padding-inline: var(--space-4);
  }

  .error-content {
    max-width: 32rem;
    width: 100%;
    text-align: center;
  }

  .error-icon-circle {
    margin-inline: auto;
    width: 6rem;
    height: 6rem;
    border-radius: 50%;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .icon-bg-warning { background: color-mix(in srgb, var(--warning) 15%, transparent); }
  .icon-bg-danger  { background: color-mix(in srgb, var(--danger) 15%, transparent); }
  .icon-bg-muted   { background: var(--faint); }

  :global(.icon-warning) { color: var(--warning); }
  :global(.icon-danger)  { color: var(--danger); }
  :global(.icon-muted)   { color: var(--muted-foreground); }
</style>
