<script lang="ts">
  interface Props {
    data: any;
    title?: string;
    expanded?: boolean;
  }

  let { data, title = 'JSON Data', expanded = false }: Props = $props();

  const formatJson = (obj: any): string => {
    if (obj === null || obj === undefined) return '';
    try {
      return JSON.stringify(obj, null, 2);
    } catch {
      return String(obj);
    }
  };

  let jsonString = $derived(formatJson(data));
  let hasData = $derived(data !== null && data !== undefined && jsonString.trim() !== '');
</script>

{#if hasData}
  <details open={expanded}>
    <summary><strong>{title}</strong> <span class="badge">JSON</span></summary>
    <pre><code>{jsonString}</code></pre>
  </details>
{/if}
