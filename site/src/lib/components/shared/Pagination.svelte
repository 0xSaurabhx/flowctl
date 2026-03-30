<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  type Props = {
    currentPage: number;
    totalPages: number;
    loading?: boolean;
    disabled?: boolean;
  };

  let {
    currentPage,
    totalPages,
    loading = false,
    disabled = false
  }: Props = $props();

  const dispatch = createEventDispatcher<{
    'page-change': { page: number };
  }>();

  const handlePageChange = (page: number) => {
    if (page !== currentPage && page >= 1 && page <= totalPages && !disabled && !loading) {
      dispatch('page-change', { page });
    }
  };

  const isPreviousDisabled = $derived(currentPage === 1 || disabled || loading);
  const isNextDisabled = $derived(currentPage === totalPages || disabled || loading);

  const getVisiblePages = () => {
    const pages: number[] = [];
    const maxVisible = 7;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      const start = Math.max(1, currentPage - 3);
      const end = Math.min(totalPages, start + maxVisible - 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }
    }

    return pages;
  };

  let visiblePages = $derived(getVisiblePages());
</script>

{#if totalPages > 1}
  <div class="pagination-wrapper">
    <nav class="hstack gap-2">
      <button
        onclick={() => handlePageChange(currentPage - 1)}
        disabled={isPreviousDisabled}
        data-variant="secondary"
      >
        Previous
      </button>

      {#each visiblePages as page}
        <button
          onclick={() => handlePageChange(page)}
          disabled={disabled || loading}
          data-variant={page === currentPage ? undefined : 'secondary'}
          aria-current={page === currentPage ? 'page' : undefined}
        >
          {page}
        </button>
      {/each}

      <button
        onclick={() => handlePageChange(currentPage + 1)}
        disabled={isNextDisabled}
        data-variant="secondary"
      >
        Next
      </button>
    </nav>
  </div>
{/if}

<style>
  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: var(--space-8);
  }
</style>
