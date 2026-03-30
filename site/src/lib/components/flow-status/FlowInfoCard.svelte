<script lang="ts">
  let {
    flowName,
    startTime,
    executionId,
    status,
    scheduledAt,
    triggerType,
    triggeredBy
  }: {
    flowName: string,
    startTime: string,
    executionId: string,
    status?: string,
    scheduledAt?: string,
    triggerType?: string,
    triggeredBy?: string
  } = $props();

  // Extract just the name from "Name <username>" format
  function extractName(triggeredBy: string): string {
    const match = triggeredBy.match(/^(.+?)\s*</);
    return match ? match[1].trim() : triggeredBy;
  }
</script>

<!-- Flow Info Card -->
<article class="card p-4 mb-4">
  <div class="hstack justify-between items-start">
    <div>
      <div class="hstack gap-2">
        <h1 class="flow-title">{flowName}</h1>
        {#if triggerType}
          <span class="badge {triggerType === 'manual' ? '' : 'success'}">
            {triggerType}
          </span>
        {/if}
      </div>
      <p class="text-light mt-2">Started at {startTime}</p>
      {#if triggeredBy}
        <p class="text-lighter text-sm mt-2">Triggered By</p>
        <p>{extractName(triggeredBy)}</p>
      {/if}
    </div>
    <div class="align-right">
      <p class="text-lighter text-sm">Execution ID</p>
      <p class="font-mono text-sm">{executionId}</p>
      {#if scheduledAt}
        <p class="text-lighter text-sm mt-2">Scheduled At</p>
        <p>{scheduledAt}</p>
      {/if}
    </div>
  </div>
</article>

<style>
  .flow-title {
    font-size: var(--text-3);
    font-weight: var(--font-semibold);
  }
</style>
