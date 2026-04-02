<script lang="ts">
  import type { NamespaceMemberResp } from '$lib/types';
  import { IconUsers, IconUser } from '@tabler/icons-svelte';

  let { row, onClick }: { row: NamespaceMemberResp; onClick?: (member: NamespaceMemberResp) => void } = $props();
</script>

<div class="hstack gap-2">
  <div class="icon-box cell-icon">
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
      <div class="cell-name font-medium">{row.subject_name}</div>
    {/if}
    <div class="text-lighter text-sm">{row.subject_id}</div>
  </div>
</div>

<style>
  .cell-icon {
    width: 2.5rem;
    height: 2.5rem;
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
  .cell-name {
    font-size: var(--text-7);
  }
</style>
