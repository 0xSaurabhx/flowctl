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

<article class="card mb-4" style="padding: 0;">
  <div class="row card-inner">
    <!-- Left: title, badge, started -->
    <div class="col-6">
      <div class="hstack gap-3 items-center">
        <h1 class="flow-title">{flowName}</h1>
        {#if triggerType}
          <span class="badge {triggerType === 'manual' ? '' : 'success'}">
            {triggerType}
          </span>
        {/if}
      </div>
      <p class="meta-text mt-2">Started at {startTime}</p>
    </div>
    <!-- Right: exec id, triggered by -->
    <div class="col-6 align-right">
      <div class="meta-item">
        <span class="meta-label">Execution ID</span>
        <span class="font-mono">{executionId}</span>
      </div>
      {#if triggeredBy}
        <div class="meta-item mt-2">
          <span class="meta-label">Triggered by</span>
          <span>{extractName(triggeredBy)}</span>
        </div>
      {/if}
      {#if scheduledAt}
        <div class="meta-item mt-2">
          <span class="meta-label">Scheduled at</span>
          <span>{scheduledAt}</span>
        </div>
      {/if}
    </div>
  </div>
</article>

<style>
  .card-inner {
    padding: var(--space-4) var(--space-6);
    align-items: start;
  }
  .flow-title {
    font-size: var(--text-4);
    font-weight: var(--font-semibold);
    margin: 0;
  }
  .meta-text {
    font-size: var(--text-7);
    color: var(--muted-foreground);
    margin-bottom: 0;
  }
  .align-right {
    text-align: right;
  }
  .align-right .meta-item {
    justify-content: flex-end;
  }
  .meta-item {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    font-size: var(--text-7);
  }
  .meta-label {
    color: var(--muted-foreground);
  }
</style>
