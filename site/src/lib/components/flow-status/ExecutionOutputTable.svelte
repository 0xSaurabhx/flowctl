<script lang="ts">
  let {
    results
  }: {
    results: Record<string, any>
  } = $props();
</script>

{#if Object.keys(results).length > 0}
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>Variable</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        {#each Object.entries(results) as [key, value]}
          <tr>
            <td class="mono">{key}</td>
            <td class="mono">{value}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}

<style>
  .table-wrap {
    overflow-x: auto;
    border-radius: 0.5rem;
    border: 1px solid var(--border);
  }
  .mono {
    font-family: monospace;
  }
</style>
