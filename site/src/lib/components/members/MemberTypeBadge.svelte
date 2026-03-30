<script lang="ts">
  let { value }: { value: string } = $props();

  const isUser = value === 'user';
</script>

<span class="badge">
  {isUser ? 'User' : 'Group'}
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    padding: 0.125rem 0.625rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
    background: color-mix(in srgb, var(--primary) 15%, transparent);
    color: var(--primary);
  }
</style>
