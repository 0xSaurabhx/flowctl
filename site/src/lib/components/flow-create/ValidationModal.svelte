<script lang="ts">
  import { autofocus } from '$lib/utils/autofocus';

  let {
    show = $bindable(),
    validationResult
  }: {
    show: boolean;
    validationResult: {
      success: boolean;
      errors: string[];
    };
  } = $props();

  let dialogEl: HTMLDialogElement;

  function close() {
    show = false;
    dialogEl?.close();
  }

  $effect(() => {
    if (show && dialogEl) {
      dialogEl.showModal();
    }
  });
</script>

{#if show}
  <dialog bind:this={dialogEl} onclose={() => show = false}>
    <header>
      <div class="hstack gap-4">
        <div class="result-icon" class:success={validationResult.success} class:danger={!validationResult.success}>
          {#if validationResult.success}
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          {:else}
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          {/if}
        </div>
        <div>
          <h3>{validationResult.success ? 'Validation Passed' : 'Validation Failed'}</h3>
          <p class="text-light">
            {validationResult.success ? 'Your flow definition is valid.' : 'Please fix the following issues:'}
          </p>
        </div>
      </div>
    </header>

    {#if !validationResult.success && validationResult.errors.length > 0}
      <section>
        <div class="vstack gap-2">
          {#each validationResult.errors as error}
            <div class="hstack gap-2 items-start text-sm error-item">
              <svg class="error-icon shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{error}</span>
            </div>
          {/each}
        </div>
      </section>
    {/if}

    <footer>
      <button
        data-variant="secondary"
        onclick={close}
        use:autofocus
      >
        Close
      </button>
    </footer>
  </dialog>
{/if}

<style>
  dialog {
    max-width: 42rem;
    width: 100%;
  }
  .result-icon {
    width: 3rem;
    height: 3rem;
    border-radius: 9999px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .result-icon.success {
    background: var(--success);
    color: white;
  }
  .result-icon.danger {
    background: var(--danger);
    color: white;
  }
  .icon {
    width: 1.5rem;
    height: 1.5rem;
  }
  .error-icon {
    width: 1rem;
    height: 1rem;
    color: var(--danger);
    margin-top: 0.125rem;
  }
</style>
