<script lang="ts">
  let {
    row,
    value,
    format,
    mono = false,
    lighter = false,
    truncate = 0
  }: {
    row: any;
    value: any;
    format?: (value: any, row: any) => string;
    mono?: boolean;
    lighter?: boolean;
    truncate?: number;
  } = $props();

  const formatted = $derived(format ? format(value, row) : (value ?? ''));
  const display = $derived(truncate > 0 && typeof formatted === 'string' ? formatted.substring(0, truncate) : formatted);
</script>

{#if truncate > 0}
  <span class="{lighter ? 'text-lighter' : 'cell-muted'} {mono ? 'font-mono' : ''}" style="max-width:20rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;display:block">{display}</span>
{:else}
  <span class="{lighter ? 'text-lighter' : 'cell-muted'} {mono ? 'font-mono' : ''}">{display}</span>
{/if}
