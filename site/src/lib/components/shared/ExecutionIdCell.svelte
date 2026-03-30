<script lang="ts">
  let { value }: { value: string } = $props();
</script>

<div class="hstack">
  <div class="exec-icon">
    <svg class="icon-sm" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
    </svg>
  </div>
  <div class="exec-value">
    <code class="text-lighter">{value}</code>
  </div>
</div>

<style>
  .exec-icon {
    width: 2rem;
    height: 2rem;
    border-radius: 0.5rem;
    background: var(--faint);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 0.75rem;
    flex-shrink: 0;
  }

  .icon-sm {
    width: 1rem;
    height: 1rem;
    color: var(--muted-foreground);
  }

  .exec-value {
    min-width: 0;
    flex: 1;
  }

  .exec-value code {
    font-size: 0.75rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    display: block;
  }
</style>
