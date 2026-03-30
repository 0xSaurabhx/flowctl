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
<div class="card p-4 mb-4">
  <div class="info-layout">
    <div>
      <div class="hstack gap-2">
        <h1 class="flow-title">{flowName}</h1>
        {#if triggerType}
          <span class="badge {triggerType === 'manual' ? '' : 'success'}">
            {triggerType}
          </span>
        {/if}
      </div>
      <p class="text-light mt-1">Started at {startTime}</p>
      {#if triggeredBy}
        <p class="text-lighter label mt-2">Triggered By</p>
        <p>{extractName(triggeredBy)}</p>
      {/if}
    </div>
    <div class="info-right">
      <p class="text-lighter label">Execution ID</p>
      <p class="mono">{executionId}</p>
      {#if scheduledAt}
        <p class="text-lighter label mt-2">Scheduled At</p>
        <p>{scheduledAt}</p>
      {/if}
    </div>
  </div>
</div>

<style>
  .card {
    background: var(--card);
    border-radius: 0.5rem;
    border: 1px solid var(--border);
  }
  .info-layout {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }
  .flow-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--foreground);
  }
  .info-right {
    text-align: right;
  }
  .label {
    font-size: 0.875rem;
  }
  .mono {
    font-family: monospace;
    font-size: 0.875rem;
  }
  .badge {
    display: inline-flex;
    align-items: center;
    padding: 0.125rem 0.625rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
    background: var(--faint);
    color: var(--foreground);
  }
  .badge.success {
    background: var(--success);
    color: white;
  }
  .p-4 { padding: 1.5rem; }
  .mb-4 { margin-bottom: 1.5rem; }
  .mt-1 { margin-top: 0.25rem; }
  .mt-2 { margin-top: 0.75rem; }
</style>
