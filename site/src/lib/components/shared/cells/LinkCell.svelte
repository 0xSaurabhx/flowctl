<script lang="ts">
  let {
    row,
    value,
    href,
    mono = false,
    truncate = 0
  }: {
    row: any;
    value: any;
    href: (row: any) => string;
    mono?: boolean;
    truncate?: number;
  } = $props();

  const display = $derived(truncate > 0 && typeof value === 'string' ? value.substring(0, truncate) : value);
</script>

<a href={href(row)} class={mono ? 'cell-link-mono' : 'cell-link'}>{display}</a>
