<script lang="ts">
  import type { ScheduledExecution, UserSchedule } from '$lib/types';
  import { getNextCronRun } from '$lib/utils/cronParser';

  interface UpcomingRun {
    type: 'cron' | 'scheduled';
    label: string;
    scheduledAt: Date;
    execId?: string;
  }

  let {
    schedules = [],
    cronSchedules = [],
    namespace,
    flowId,
    title = 'Upcoming Scheduled Runs'
  }: {
    schedules: ScheduledExecution[];
    cronSchedules?: UserSchedule[];
    namespace: string;
    flowId: string;
    title?: string;
  } = $props();

  // Compute combined list of upcoming runs
  let upcomingRuns = $derived.by(() => {
    const runs: UpcomingRun[] = [];

    // Add cron-based runs
    for (const cron of cronSchedules) {
      const nextRun = getNextCronRun(cron.cron, cron.timezone);
      if (nextRun) {
        runs.push({
          type: 'cron',
          label: cron.cron,
          scheduledAt: nextRun
        });
      }
    }

    // Add manually scheduled runs
    for (const schedule of schedules) {
      runs.push({
        type: 'scheduled',
        label: 'Scheduled',
        scheduledAt: new Date(schedule.scheduled_at),
        execId: schedule.exec_id
      });
    }

    // Sort by scheduled time ascending
    return runs.sort((a, b) => a.scheduledAt.getTime() - b.scheduledAt.getTime());
  });

  function formatScheduledTime(date: Date): string {
    return date.toLocaleString();
  }
</script>

{#if upcomingRuns.length > 0}
  <div class="card">
    <div class="card-header">
      <h3 class="section-title">{title}</h3>
      <p class="text-lighter section-subtitle">{upcomingRuns.length} {upcomingRuns.length === 1 ? 'run' : 'runs'} scheduled</p>
    </div>
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Type</th>
            <th>Scheduled Time</th>
            <th>Exec ID</th>
          </tr>
        </thead>
        <tbody>
          {#each upcomingRuns as run}
            <tr>
              <td>
                {#if run.type === 'cron'}
                  <code>{run.label}</code>
                {:else}
                  {run.label}
                {/if}
              </td>
              <td>
                {formatScheduledTime(run.scheduledAt)}
              </td>
              <td>
                {#if run.execId}
                  <a href="/view/{namespace}/results/{flowId}/{run.execId}" class="mono-link">
                    {run.execId.substring(0, 8)}
                  </a>
                {:else}
                  <span class="text-lighter">-</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<style>
  .card {
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    background: var(--card);
  }
  .card-header {
    padding: 1rem;
    border-bottom: 1px solid var(--border);
  }
  .section-title {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--foreground);
  }
  .section-subtitle {
    font-size: 0.75rem;
    margin-top: 0.125rem;
  }
  .table-wrap {
    overflow-x: auto;
  }
  .mono-link {
    font-family: monospace;
    font-size: 0.875rem;
  }
</style>
