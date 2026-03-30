<script lang="ts">
  import type { ApprovalResp } from '$lib/types';

  let { row, onClick }: { row: ApprovalResp; onClick?: (approval: ApprovalResp) => void } = $props();
</script>

<div class="hstack gap-2">
  <div class="cell-icon">
    <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
    </svg>
  </div>
  <div>
    <a href="#" class="cell-link" onclick={(e) => { e.preventDefault(); onClick?.(row); }}>{row.action_id}</a>
    <div class="text-lighter cell-id">{row.id}</div>
  </div>
</div>

<style>
  .cell-icon {
    width: 2.5rem;
    height: 2.5rem;
    background: var(--faint);
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .icon {
    width: 1.25rem;
    height: 1.25rem;
    color: var(--primary);
  }
  .cell-link {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--foreground);
  }
  .cell-link:hover {
    color: var(--primary);
    text-decoration: underline;
  }
  .cell-id {
    font-size: 0.875rem;
  }
</style>
