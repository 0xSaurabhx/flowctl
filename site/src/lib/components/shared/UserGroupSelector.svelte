<script lang="ts">
  import { apiClient } from '$lib/apiClient';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import type { User, Group } from '$lib/types';
  import { IconUsers, IconUser } from '@tabler/icons-svelte';

  let {
    type = $bindable('user'),
    selectedSubject = $bindable(null),
    placeholder = 'Search...',
    disabled = false
  }: {
    type: 'user' | 'group';
    selectedSubject: User | Group | null;
    placeholder?: string;
    disabled?: boolean;
  } = $props();

  let searchQuery = $state('');
  let searchResults = $state<(User | Group)[]>([]);
  let loading = $state(false);
  let inputEl: HTMLInputElement;
  let popoverEl: HTMLDivElement;

  async function loadSubjects() {
    loading = true;
    try {
      if (type === 'user') {
        const response = await apiClient.users.list({
          filter: searchQuery,
          count_per_page: 20
        });
        searchResults = response.users || [];
      } else {
        const response = await apiClient.groups.list({
          filter: searchQuery,
          count_per_page: 20
        });
        searchResults = response.groups || [];
      }
    } catch (error) {
      handleInlineError(error, type === 'user' ? 'Unable to Load Users' : 'Unable to Load Groups');
      searchResults = [];
      closeDropdown();
    } finally {
      loading = false;
    }
  }

  function positionPopover() {
    if (!inputEl || !popoverEl) return;
    const r = inputEl.getBoundingClientRect();
    popoverEl.style.top = `${r.bottom + 4}px`;
    popoverEl.style.left = `${r.left}px`;
    popoverEl.style.width = `${r.width}px`;
  }

  function openDropdown() {
    try { popoverEl?.showPopover(); } catch {}
    positionPopover();
  }

  function closeDropdown() {
    try { popoverEl?.hidePopover(); } catch {}
  }

  async function handleInput() {
    await loadSubjects();
    openDropdown();
  }

  async function handleFocus() {
    if (searchResults.length === 0) {
      await loadSubjects();
    }
    openDropdown();
  }

  function selectSubject(subject: User | Group) {
    selectedSubject = subject;
    searchQuery = '';
    closeDropdown();
    searchResults = [];
  }

  function clearSelection() {
    selectedSubject = null;
    searchQuery = '';
    searchResults = [];
    closeDropdown();
  }

  // Reset when type changes
  $effect(() => {
    clearSelection();
  });

  // Close dropdown when clicking outside
  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.user-group-selector') && !target.closest('[popover]')) {
      closeDropdown();
    }
  }
</script>

<svelte:window onclick={handleOutsideClick} />

<div class="user-group-selector">
  <div class="selector-field">
    <input
      type="text"
      bind:this={inputEl}
      bind:value={searchQuery}
      oninput={handleInput}
      onfocus={handleFocus}
      {placeholder}
      aria-busy={loading}
      autocomplete="off"
      {disabled}
    />

    {#if !loading}
      <svg class="search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
      </svg>
    {/if}
  </div>

  <!-- Popover dropdown — renders in top layer, escapes dialog overflow -->
  <div
    bind:this={popoverEl}
    popover="manual"
    class="selector-popover"
    role="listbox"
  >
    {#if searchResults.length > 0}
      {#each searchResults as subject}
        <button
          type="button"
          class="selector-item"
          onclick={() => selectSubject(subject)}
          role="option"
        >
          <div class="hstack">
            <div class="item-icon">
              {#if type === 'user'}
                <IconUser size={16} />
              {:else}
                <IconUsers size={16} />
              {/if}
            </div>
            <div>
              <div class="item-name">{'name' in subject ? subject.name : subject.username}</div>
              <div class="text-lighter item-desc">{subject.id}</div>
            </div>
          </div>
        </button>
      {/each}
    {:else if !loading}
      <div class="selector-empty text-lighter">
        {type === 'user' ? 'No users found' : 'No groups found'}
      </div>
    {/if}
  </div>

  <!-- Selected subject display -->
  {#if selectedSubject}
    <div class="selected-subject mt-2">
      <div class="hstack justify-between">
        <div class="hstack">
          <div class="item-icon">
            {#if type === 'user'}
              <IconUser size={16} />
            {:else}
              <IconUsers size={16} />
            {/if}
          </div>
          <div>
            <div class="item-name">{'name' in selectedSubject ? selectedSubject.name : selectedSubject.username}</div>
            <div class="text-lighter item-desc">{selectedSubject.id}</div>
          </div>
        </div>
        <button type="button" onclick={clearSelection} class="clear-btn" {disabled} aria-label="Clear selection">
          <svg class="icon-sm" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  .selector-field {
    position: relative;
  }

  .selector-field input {
    width: 100%;
    padding-right: 2.5rem;
  }

  .search-icon {
    width: 1.25rem;
    height: 1.25rem;
    position: absolute;
    right: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--muted-foreground);
    pointer-events: none;
  }

  .selector-popover {
    position: fixed;
    margin: 0;
    padding: 0;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius-medium);
    box-shadow: var(--shadow-large);
    max-height: 14rem;
    overflow-y: auto;
  }

  .selector-item {
    all: unset;
    box-sizing: border-box;
    display: flex;
    width: 100%;
    padding: var(--space-2) var(--space-4);
    cursor: pointer;
    border-bottom: 1px solid var(--border);
    font-size: var(--text-7);
  }

  .selector-item:last-child {
    border-bottom: none;
  }

  .selector-item:hover {
    background: var(--faint);
  }

  .item-icon {
    width: 2rem;
    height: 2rem;
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 0.75rem;
    background: var(--faint);
    color: var(--primary);
  }

  .item-name {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--foreground);
  }

  .item-desc {
    font-size: 0.75rem;
  }

  .selector-empty {
    padding: var(--space-3) var(--space-4);
    font-size: var(--text-7);
    text-align: center;
  }

  .selected-subject {
    padding: 0.5rem;
    background: var(--faint);
    border-radius: 0.5rem;
    border: 1px solid var(--border);
  }

  .clear-btn {
    all: unset;
    cursor: pointer;
    display: flex;
    color: var(--muted-foreground);
  }

  .clear-btn:hover {
    color: var(--foreground);
  }

  .icon-sm {
    width: 1rem;
    height: 1rem;
  }
</style>
