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

<input
  type="search"
  {placeholder}
  {value}
  oninput={handleInput}
  style="min-width: 16rem"
/>
