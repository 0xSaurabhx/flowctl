<script lang="ts">
  let {
    name,
    description,
    icon = `<svg class="hero-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/>
    </svg>`
  }: {
    name: string,
    description?: string,
    icon?: string
  } = $props();
</script>

<!-- Hero Section -->
<div class="hero-section">
  <div class="hero-inner hstack gap-4 items-center">
    <div class="icon-box shrink-0" style="width:2.5rem;height:2.5rem">
      {@html icon}
    </div>
    <div class="flex-1 min-w-0">
      <h1 class="hero-title">{name}</h1>
      {#if description}
        <p class="text-light hero-description">{description}</p>
      {/if}
    </div>
  </div>
</div>

<style>
  .hero-section {
    background: var(--card);
    padding: var(--space-5) var(--space-6) var(--space-4);
  }
  .hero-inner {
    max-width: 72rem;
    margin: 0 auto;
  }
  .hero-title {
    font-size: var(--text-4);
    font-weight: var(--font-semibold);
    margin: 0;
    line-height: 1.3;
  }
  .hero-description {
    font-size: var(--text-7);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin: var(--space-1) 0 0;
  }
</style>
