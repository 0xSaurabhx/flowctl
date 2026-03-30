<script lang="ts">
  import type { NamespaceMemberResp } from '$lib/types';
  import { IconUsers, IconUser } from '@tabler/icons-svelte';

  let { row, onClick }: { row: NamespaceMemberResp; onClick?: (member: NamespaceMemberResp) => void } = $props();
</script>

<div class="hstack gap-2">
  <div class="cell-icon">
    {#if row.subject_type === 'user'}
      <IconUser size={20} />
    {:else}
      <IconUsers size={20} />
    {/if}
  </div>
  <div>
    {#if onClick}
      <a href="#" class="cell-link" onclick={(e) => { e.preventDefault(); onClick(row); }}>{row.subject_name}</a>
    {:else}
      <div class="cell-name">{row.subject_name}</div>
    {/if}
    <div class="text-lighter cell-id">{row.subject_id}</div>
  </div>
</div>

<style>
  .cell-icon {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--faint);
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
  .cell-name {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--foreground);
  }
  .cell-id {
    font-size: 0.875rem;
  }
</style>
