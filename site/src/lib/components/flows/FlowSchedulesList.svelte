<script lang="ts">
  import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
  import { apiClient } from '$lib/apiClient';
  import ScheduleModal from './ScheduleModal.svelte';
  import ViewScheduleModal from './ViewScheduleModal.svelte';
  import DeleteModal from '$lib/components/shared/DeleteModal.svelte';
  import DropdownMenu from '$lib/components/shared/DropdownMenu.svelte';
  import { IconClock, IconPlus } from '@tabler/icons-svelte';
  import type { UserSchedule, FlowInput, ScheduleCreateReq, ScheduleUpdateReq } from '$lib/types';

  let {
    namespace,
    flowId,
    flowInputs,
    userSchedulable = false,
    user,
    schedules = [],
    onUpdate,
    canUpdateFlow = false
  }: {
    namespace: string;
    flowId: string;
    flowInputs: FlowInput[];
    userSchedulable: boolean;
    user: any;
    schedules?: UserSchedule[];
    onUpdate?: () => Promise<void>;
    canUpdateFlow?: boolean;
  } = $props();

  let showModal = $state(false);
  let showDelete = $state(false);
  let showViewModal = $state(false);
  let editSchedule = $state<UserSchedule | null>(null);
  let deleteSchedule = $state<UserSchedule | null>(null);
  let viewSchedule = $state<UserSchedule | null>(null);
  let canCreateSchedule = $derived(userSchedulable);

  function canEdit(schedule: UserSchedule): boolean {
    return canUpdateFlow || (schedule.is_user_created && schedule.created_by === user.id);
  }

  async function handleSave(data: ScheduleCreateReq | ScheduleUpdateReq) {
    if (editSchedule) {
      await apiClient.flows.schedules.update(namespace, flowId, editSchedule.uuid, data);
    } else {
      await apiClient.flows.schedules.create(namespace, flowId, data as ScheduleCreateReq);
    }
    if (onUpdate) await onUpdate();
  }

  async function handleDelete() {
    if (!deleteSchedule) return;
    await apiClient.flows.schedules.delete(namespace, flowId, deleteSchedule.uuid);
    showSuccess('Schedule Deleted', 'Schedule deleted successfully');
    showDelete = false;
    deleteSchedule = null;
    if (onUpdate) await onUpdate();
  }

  async function toggleActive(schedule: UserSchedule) {
    try {
      await apiClient.flows.schedules.update(namespace, flowId, schedule.uuid, {
        cron: schedule.cron,
        timezone: schedule.timezone,
        inputs: schedule.inputs,
        is_active: !schedule.is_active
      });
      showSuccess('Updated', 'Schedule status updated');
      if (onUpdate) await onUpdate();
    } catch (error) {
      handleInlineError(error, 'Failed to update schedule');
    }
  }

  function getMenuItems(schedule: UserSchedule) {
    return [
      {
        label: 'View',
        onClick: () => { viewSchedule = schedule; showViewModal = true; }
      },
      {
        label: 'Edit',
        onClick: () => { editSchedule = schedule; showModal = true; }
      },
      {
        label: schedule.is_active ? 'Deactivate' : 'Activate',
        onClick: () => toggleActive(schedule)
      },
      {
        label: 'Delete',
        onClick: () => { deleteSchedule = schedule; showDelete = true; },
        variant: 'danger' as const
      }
    ];
  }
</script>

<div class="card">
  <div class="card-header hstack justify-between">
    <div>
      <h3 class="section-title">Schedules</h3>
      <p class="text-lighter section-subtitle">{schedules.length} {schedules.length === 1 ? 'schedule' : 'schedules'}</p>
    </div>
    {#if canCreateSchedule}
      <button
        type="button"
        onclick={() => { editSchedule = null; showModal = true; }}
      >
        <IconPlus size={16} />
        Add
      </button>
    {/if}
  </div>

  {#if schedules.length === 0}
    <div class="vstack empty-state">
      <IconClock size={48} class="text-lighter" />
      <p class="text-light">No schedules configured</p>
      {#if canCreateSchedule}
        <button
          type="button"
          data-variant="secondary"
          onclick={() => { editSchedule = null; showModal = true; }}
        >
          Create your first schedule
        </button>
      {/if}
    </div>
  {:else}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Cron</th>
            <th>Timezone</th>
            <th>Type</th>
            <th>Status</th>
            <th class="actions-col">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each schedules as schedule}
            <tr>
              <td>
                <code>{schedule.cron}</code>
              </td>
              <td>{schedule.timezone}</td>
              <td>
                <span class="badge {schedule.is_user_created ? 'success' : 'info'}">
                  {schedule.is_user_created ? 'User' : 'System'}
                </span>
              </td>
              <td>
                {#if schedule.is_user_created}
                  <span class="badge {schedule.is_active ? 'success' : ''}">
                    {schedule.is_active ? 'Active' : 'Inactive'}
                  </span>
                {:else}
                  <span class="text-lighter">-</span>
                {/if}
              </td>
              <td class="actions-col">
                {#if canCreateSchedule && canEdit(schedule)}
                  <div class="actions-wrapper">
                    <DropdownMenu items={getMenuItems(schedule)} />
                  </div>
                {:else}
                  <span class="text-lighter">-</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

{#if showModal}
  <ScheduleModal
    isEditMode={!!editSchedule}
    schedule={editSchedule}
    {flowInputs}
    {namespace}
    {flowId}
    onSave={handleSave}
    onClose={() => { showModal = false; editSchedule = null; }}
  />
{/if}

{#if showDelete && deleteSchedule}
  <DeleteModal
    title="Delete Schedule"
    itemName={`${deleteSchedule.cron} (${deleteSchedule.timezone})`}
    onConfirm={handleDelete}
    onClose={() => { showDelete = false; deleteSchedule = null; }}
  />
{/if}

{#if showViewModal && viewSchedule}
  <ViewScheduleModal
    schedule={viewSchedule}
    onClose={() => { showViewModal = false; viewSchedule = null; }}
  />
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
  .empty-state {
    align-items: center;
    padding: 3rem 1rem;
    gap: 0.75rem;
  }
  .table-wrap {
    overflow-x: auto;
  }
  .actions-col {
    text-align: right;
    width: 5rem;
  }
  .actions-wrapper {
    display: inline-flex;
    justify-content: flex-end;
  }
  .badge {
    display: inline-flex;
    padding: 0.125rem 0.5rem;
    font-size: 0.75rem;
    font-weight: 500;
    border-radius: 0.25rem;
    background: var(--faint);
    color: var(--foreground);
  }
  .badge.success {
    background: var(--success);
    color: white;
  }
  .badge.info {
    background: var(--primary);
    color: white;
  }
</style>
