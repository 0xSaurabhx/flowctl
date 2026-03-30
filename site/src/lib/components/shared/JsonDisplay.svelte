<script lang="ts">
  interface Props {
    data: any;
    title?: string;
    expanded?: boolean;
  }

  let { data, title = 'JSON Data', expanded = false }: Props = $props();

  let isExpanded = $state(expanded);

  const toggleExpanded = () => {
    isExpanded = !isExpanded;
  };

  const formatJson = (obj: any): string => {
    if (obj === null || obj === undefined) {
      return '';
    }

    try {
      return JSON.stringify(obj, null, 2);
    } catch (error) {
      return String(obj);
    }
  };

  let jsonString = $derived(formatJson(data));
  let hasData = $derived(data !== null && data !== undefined && jsonString.trim() !== '');
</script>

{#if hasData}
  <article class="card">
    <button
      class="json-toggle"
      onclick={toggleExpanded}
      type="button"
    >
      <strong>{title}</strong>
      <div class="hstack gap-2">
        <span class="badge">JSON</span>
        <svg
          class="chevron {isExpanded ? 'expanded' : ''}"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
        </svg>
      </div>
    </button>

    {#if isExpanded}
      <div class="json-content">
        <pre><code>{jsonString}</code></pre>
      </div>
    {/if}
  </article>
{/if}

<style>
  .json-toggle {
    all: unset;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: var(--space-3) var(--space-6);
    box-sizing: border-box;
  }

  .json-toggle:hover {
    background: var(--muted);
  }

  .chevron {
    width: 1rem;
    height: 1rem;
    color: var(--muted-foreground);
    transition: transform 0.2s;
  }

  .chevron.expanded {
    transform: rotate(180deg);
  }

  .json-content {
    padding: var(--space-4);
  }

  pre {
    padding: var(--space-4);
    border-radius: var(--radius-medium);
    font-size: var(--text-7);
    overflow-x: auto;
    line-height: 1.6;
    white-space: pre-wrap;
  }
</style>
