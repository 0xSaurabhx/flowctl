<script lang="ts">
  let { value }: { value: string } = $props();
</script>

<div class="hstack">
  <div class="icon-box">
    <svg class="icon-sm" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
    </svg>
  </div>
  <div class="exec-value min-w-0 flex-1">
    <code class="text-lighter">{value}</code>
  </div>
</div>

<style>
  .icon-sm {
    width: 1rem;
    height: 1rem;
    color: var(--muted-foreground);
  }

  .exec-value code {
    font-size: var(--text-8);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    display: block;
  }
</style>
