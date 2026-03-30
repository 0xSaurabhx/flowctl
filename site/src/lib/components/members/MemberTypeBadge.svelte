<script lang="ts">
  let { value }: { value: string } = $props();

  const isUser = value === 'user';
</script>

<span class="badge">
  {isUser ? 'User' : 'Group'}
</span>

<style>
  .badge {
    background: color-mix(in srgb, var(--primary) 15%, transparent);
    color: var(--primary);
  }
</style>
