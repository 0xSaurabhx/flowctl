<script lang="ts">
  let {
    value = $bindable(''),
    placeholder = 'Search...',
    loading = false,
    onSearch,
    debounceMs = 300
  }: {
    value: string,
    placeholder?: string,
    loading?: boolean,
    onSearch: (query: string) => void,
    debounceMs?: number
  } = $props();

  let debounceTimer: number;

  const handleInput = (event: Event) => {
    const target = event.target as HTMLInputElement;
    value = target.value;

    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      onSearch(value.trim());
    }, debounceMs);
  };
</script>

<div class="search-wrapper">
  <div class="search-field">
    <input
      type="search"
      {placeholder}
      {value}
      oninput={handleInput}
    />
    {#if loading}
      <div class="search-spinner" aria-busy="true"></div>
    {/if}
  </div>
</div>

<style>
  .search-wrapper {
    max-width: 28rem;
  }

  .search-field {
    position: relative;
  }

  .search-spinner {
    position: absolute;
    right: var(--space-3);
    top: 50%;
    transform: translateY(-50%);
    width: 1rem;
    height: 1rem;
  }
</style>
