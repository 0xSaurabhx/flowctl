<script lang="ts">
  let { value }: { value: string } = $props();

  const isUser = value === 'user';
</script>

<span class="badge outline">
  {isUser ? 'User' : 'Group'}
</span>
