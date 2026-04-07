<script lang="ts">
  import { goto } from '$app/navigation';
  import { isAuthenticated, isLoading } from '$lib/stores/auth';
  import { getDefaultNamespace } from '$lib/utils/navigation';
  import Logo from '$lib/components/shared/Logo.svelte';

  // Wait for auth loading to complete
  $effect(() => {
    if ($isLoading) {
      return;
    }

    if (!$isAuthenticated) {
      goto('/login');
      return;
    }

    getDefaultNamespace()
      .then((namespace) => goto(`/view/${namespace}/flows`))
      .catch(() => goto('/login'));
  });
</script>

<svelte:head>
  <title>Flowctl</title>
</svelte:head>

<div class="loading-page">
  <div class="vstack items-center gap-6">
    <Logo height="4rem" />
    <div class="hstack justify-center gap-4" aria-busy="true">Loading...</div>
  </div>
</div>

<style>
  .loading-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--card);
  }
</style>
