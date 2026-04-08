<script lang="ts">
  let { value }: { value: string } = $props();
</script>

<span class="badge outline">
  {value === 'user' ? 'User' : 'Group'}
</span>
