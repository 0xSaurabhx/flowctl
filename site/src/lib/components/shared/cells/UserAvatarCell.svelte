<script lang="ts">
  import { IconUser } from '@tabler/icons-svelte';

  let { value }: { row: any; value: any } = $props();
</script>

<div class="hstack gap-3">
  <div class="icon-box" style="width:2rem;height:2rem;border-radius:50%">
    <IconUser size={16} />
  </div>
  <span class="text-sm">{value || 'System'}</span>
</div>
