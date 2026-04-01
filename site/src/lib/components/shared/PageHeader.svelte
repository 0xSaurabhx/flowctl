<script lang="ts">
  import type { ComponentType } from 'svelte';

  let {
    title,
    subtitle,
    actions = []
  }: {
    title: string,
    subtitle?: string,
    actions?: Array<{
      label: string,
      onClick: () => void,
      variant?: 'primary' | 'secondary',
      icon?: string,
      IconComponent?: ComponentType,
      iconSize?: number
    }>
  } = $props();
</script>

<header class="page-header hstack justify-between mb-6">
  <div class="title-block">
    <h1>{title}</h1>
    {#if subtitle}
      <p class="text-light">{subtitle}</p>
    {/if}
  </div>

  {#if actions.length > 0}
    <div class="hstack gap-3">
      {#each actions as action}
        <button
          onclick={action.onClick}
          data-variant={action.variant === 'secondary' ? 'secondary' : undefined}
          aria-label={action.label}
          class="hstack gap-2"
        >
          {#if action.IconComponent}
            <action.IconComponent size={action.iconSize || 16} aria-hidden="true" />
          {:else if action.icon}
            {@html action.icon}
          {/if}
          {action.label}
        </button>
      {/each}
    </div>
  {/if}
</header>

<style>
  .title-block h1 {
    font-size: var(--text-3);
    margin: 0;
  }
  .title-block p {
    margin: var(--space-1) 0 0;
    font-size: var(--text-7);
  }
</style>
