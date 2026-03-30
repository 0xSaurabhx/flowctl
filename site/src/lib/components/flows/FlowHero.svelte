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
  <div class="hero-inner">
    <div class="hstack gap-4 items-start">
      <div class="hero-icon-wrapper">
        {@html icon}
      </div>
      <div class="hero-content">
        <h1 class="hero-title">{name}</h1>
        {#if description}
          <p class="text-light hero-description">{description}</p>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .hero-section {
    background: var(--card);
    border-bottom: 1px solid var(--border);
    padding: 1.5rem;
  }
  .hero-inner {
    max-width: 56rem;
    margin: 0 auto;
  }
  .hero-icon-wrapper {
    width: 3rem;
    height: 3rem;
    background: var(--faint);
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  .hero-icon-wrapper :global(svg) {
    width: 1.5rem;
    height: 1.5rem;
    color: var(--muted-foreground);
  }
  .hero-content {
    flex: 1;
    min-width: 0;
  }
  .hero-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--foreground);
    margin-bottom: 0.25rem;
  }
  .hero-description {
    font-size: 0.875rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
