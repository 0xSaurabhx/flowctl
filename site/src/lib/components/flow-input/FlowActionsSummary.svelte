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
  <div class="card mt-4">
    <div class="card-header">
      <h3 class="section-title">
        {title} ({actions.length})
      </h3>
    </div>
    <div class="action-list">
      {#each actions as action, index}
        <div class="action-item hstack justify-between">
          <div class="hstack gap-2">
            <div class="action-number">
              {index + 1}
            </div>
            <div>{action.name}</div>
          </div>
          <span class="badge">
            {action.executor}
          </span>
        </div>
      {/each}
    </div>
  </div>
{/if}

<style>
  .card {
    background: var(--card);
    border-radius: 0.5rem;
    border: 1px solid var(--border);
  }
  .card-header {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--border);
  }
  .section-title {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--foreground);
  }
  .action-list {
    display: flex;
    flex-direction: column;
  }
  .action-list > :not(:last-child) {
    border-bottom: 1px solid var(--border);
  }
  .action-item {
    padding: 0.75rem 1rem;
    font-size: 0.875rem;
    color: var(--foreground);
  }
  .action-item:hover {
    background: var(--faint);
  }
  .action-number {
    width: 1.5rem;
    height: 1.5rem;
    background: var(--faint);
    color: var(--primary);
    border-radius: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 500;
  }
  .badge {
    display: inline-flex;
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
    font-weight: 500;
    border-radius: 0.25rem;
    background: var(--faint);
    color: var(--foreground);
  }
  .mt-4 { margin-top: 1.5rem; }
</style>
