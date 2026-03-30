<script lang="ts">
  import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
  import { autofocus } from '$lib/utils/autofocus';
  import { getTimezones } from '$lib/utils/timezone';
  import { isValidCronExpression } from '$lib/utils';
  import FlowInputFields from '$lib/components/shared/FlowInputFields.svelte';
  import type { FlowInput, UserSchedule, ScheduleCreateReq, ScheduleUpdateReq } from '$lib/types';

  let {
    isEditMode = false,
    schedule = null,
    flowInputs,
    namespace,
    flowId,
    onSave,
    onClose
  }: {
    isEditMode?: boolean;
    schedule?: UserSchedule | null;
    flowInputs: FlowInput[];
    namespace: string;
    flowId: string;
    onSave: (data: ScheduleCreateReq | ScheduleUpdateReq) => Promise<void>;
    onClose: () => void;
  } = $props();

  // Initialize form data based on edit mode
  function getInitialFormData() {
    if (isEditMode && schedule) {
      return {
        cron: schedule.cron,
        timezone: schedule.timezone,
        is_active: schedule.is_active,
        inputs: { ...schedule.inputs }
      };
    }
    // Initialize inputs with defaults for create mode
    const initialInputs: Record<string, any> = {};
    flowInputs.forEach(input => {
      if (input.default) {
        initialInputs[input.name] = input.default;
      }
    });
    return {
      cron: '',
      timezone: 'UTC',
      is_active: true,
      inputs: initialInputs
    };
  }

  let formData = $state(getInitialFormData());
  let loading = $state(false);
  let cronError = $state('');
  let hasScheduleWarning = $derived(
    flowInputs.some(input => input.type === 'file' || !input.default || input.default.trim() === '')
  );
  const timezones = getTimezones();

  let dialogEl: HTMLDialogElement;

  function validateCron() {
    if (!formData.cron.trim()) {
      cronError = 'Cron expression is required';
      return false;
    }
    if (!isValidCronExpression(formData.cron)) {
      cronError = 'Invalid cron expression';
      return false;
    }
    cronError = '';
    return true;
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    if (!validateCron()) return;

    loading = true;
    try {
      await onSave(formData);
      showSuccess(
        isEditMode ? 'Schedule Updated' : 'Schedule Created',
        'Operation completed successfully'
      );
      onClose();
    } catch (error) {
      handleInlineError(error, 'Unable to Save Schedule');
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (dialogEl) {
      dialogEl.showModal();
    }
  });
</script>

<dialog bind:this={dialogEl} onclose={onClose}>
  <form onsubmit={handleSubmit}>
    <header>
      <h3>{isEditMode ? 'Edit Schedule' : 'Create Schedule'}</h3>
    </header>

    <section>
      <div data-field>
        <label for="cron-expr">Cron Expression *</label>
        <input
          id="cron-expr"
          type="text"
          bind:value={formData.cron}
          onblur={validateCron}
          class:error={!!cronError}
          placeholder="0 2 * * *"
          required
          use:autofocus
        />
        {#if cronError}
          <p class="field-error text-sm">{cronError}</p>
        {/if}
        <p class="text-lighter text-xs field-hint">
          Examples: <code>0 2 * * *</code> (daily 2AM),
          <code>0 */6 * * *</code> (every 6 hours)
        </p>
      </div>

      <div data-field>
        <label for="tz-input">Timezone *</label>
        <input
          id="tz-input"
          type="text"
          list="tz-list"
          bind:value={formData.timezone}
          required
        />
        <datalist id="tz-list">
          {#each timezones as tz}
            <option value={tz.tzCode}>{tz.label}</option>
          {/each}
        </datalist>
      </div>

      {#if isEditMode}
        <div class="hstack justify-between toggle-row">
          <span>Active</span>
          <input
            type="checkbox"
            role="switch"
            checked={formData.is_active}
            onchange={() => formData.is_active = !formData.is_active}
          />
        </div>
      {/if}

      {#if flowInputs.length > 0}
        <div class="inputs-section">
          <h3 class="text-sm font-medium">Flow Inputs</h3>
          <FlowInputFields inputs={flowInputs} bind:values={formData.inputs} />
        </div>
      {/if}
    </section>

    <footer>
      <button
        type="button"
        data-variant="secondary"
        onclick={onClose}
        disabled={loading}
      >
        Cancel
      </button>
      <button
        type="submit"
        disabled={loading}
        aria-busy={loading}
      >
        {loading ? 'Saving...' : (isEditMode ? 'Update' : 'Create')}
      </button>
    </footer>
  </form>
</dialog>

<style>
  dialog {
    max-width: 42rem;
    width: 100%;
  }
  .error {
    border-color: var(--danger) !important;
  }
  .field-error {
    color: var(--danger);
    margin-top: 0.25rem;
  }
  .field-hint {
    margin-top: 0.25rem;
  }
  .toggle-row {
    padding: 0.5rem 0;
    border-top: 1px solid var(--border);
  }
  .inputs-section {
    padding-top: 1rem;
    border-top: 1px solid var(--border);
  }
  .inputs-section h3 {
    margin-bottom: 0.75rem;
  }
</style>
