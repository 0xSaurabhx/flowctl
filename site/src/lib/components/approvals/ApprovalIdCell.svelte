<script lang="ts">
  import type { ApprovalResp } from '$lib/types';

  let { row, onClick }: { row: ApprovalResp; onClick?: (approval: ApprovalResp) => void } = $props();
</script>

<div class="hstack gap-2">
  <div class="icon-box cell-icon">
    <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
    </svg>
  </div>
  <div>
    <a href="#" class="cell-link" onclick={(e) => { e.preventDefault(); onClick?.(row); }}>{row.action_id}</a>
    <div class="text-lighter text-sm">{row.id}</div>
  </div>
</div>

<style>
  .cell-icon {
    width: 2.5rem;
    height: 2.5rem;
  }
  .icon {
    width: 1.25rem;
    height: 1.25rem;
  }
  .cell-link {
    font-size: var(--text-7);
    font-weight: var(--font-medium);
    color: var(--primary);
    text-decoration: none;
  }
  .cell-link:hover {
    color: var(--primary);
    text-decoration: underline;
  }
</style>
