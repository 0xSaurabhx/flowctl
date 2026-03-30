<script lang="ts">
  import type { FlowAction } from '$lib/types';

  let {
    actions,
    title = 'Flow Actions'
  }: {
    actions: FlowAction[],
    title?: string
  } = $props();
</script>

{#if actions && actions.length > 0}
  <article class="card mt-4">
    <header class="card-list-header">
      <h3 class="text-sm font-medium">
        {title} ({actions.length})
      </h3>
    </header>
    <div class="vstack gap-1">
      {#each actions as action, index}
        <div class="action-item hstack justify-between">
          <div class="hstack gap-2">
            <div class="action-number">
              {index + 1}
            </div>
            <div class="text-sm">{action.name}</div>
          </div>
          <span class="badge">
            {action.executor}
          </span>
        </div>
      {/each}
    </div>
  </article>
{/if}

<style>
  .card-list-header {
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--border);
  }
  .action-item {
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--border);
  }
  .action-item:last-child {
    border-bottom: none;
  }
  .action-item:hover {
    background: var(--faint);
  }
  .action-number {
    width: 1.5rem;
    height: 1.5rem;
    background: var(--faint);
    color: var(--primary);
    border-radius: var(--radius-small);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--text-8);
    font-weight: var(--font-medium);
  }
</style>
