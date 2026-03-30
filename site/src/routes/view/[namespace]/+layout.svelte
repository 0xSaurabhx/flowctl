<script lang="ts">
  import Sidebar from '$lib/components/shared/Sidebar.svelte';
  import Header from '$lib/components/shared/Header.svelte';

  let { children, data } = $props();
</script>

<svelte:head>
  <meta name="viewport" content="width=1024, initial-scale=1.0, user-scalable=yes" />
</svelte:head>

<div class="app-shell">
  <Sidebar namespace={data.namespace} />

  <div class="app-main">
    <main class="app-content">
      {@render children()}
    </main>
  </div>
</div>

<style>
  .app-shell {
    display: flex;
    height: 100vh;
    background: var(--muted);
  }

  .app-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .app-content {
    flex: 1;
    overflow-y: auto;
  }
</style>
