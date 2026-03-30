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
  let showDropdown = $state(false);
  let loading = $state(false);

  async function loadSubjects() {
    loading = true;
    showDropdown = true;
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
      showDropdown = false;
    } finally {
      loading = false;
    }
  }

  async function handleFocus() {
    if (searchResults.length === 0) {
      await loadSubjects();
    } else {
      showDropdown = true;
    }
  }

  function selectSubject(subject: User | Group) {
    selectedSubject = subject;
    searchQuery = '';
    showDropdown = false;
    searchResults = [];
  }

  function clearSelection() {
    selectedSubject = null;
    searchQuery = '';
    searchResults = [];
    showDropdown = false;
  }

  // Reset when type changes
  $effect(() => {
    clearSelection();
  });

  // Close dropdown when clicking outside
  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.user-group-selector')) {
      showDropdown = false;
    }
  }
</script>

<svelte:window onclick={handleOutsideClick} />

<div class="user-group-selector">
  <div class="selector-field">
    <input
      type="text"
      bind:value={searchQuery}
      oninput={loadSubjects}
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

    <!-- Dropdown -->
    {#if showDropdown}
      <div class="dropdown-menu">
        {#if searchResults.length > 0}
          {#each searchResults as subject}
            <button
              type="button"
              class="dropdown-item"
              onclick={() => selectSubject(subject)}
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
          <div class="dropdown-empty text-lighter">
            {type === 'user' ? 'No users found' : 'No groups found'}
          </div>
        {/if}
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

  .dropdown-menu {
    position: absolute;
    z-index: 10;
    width: 100%;
    margin-top: 0.25rem;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    max-height: 12rem;
    overflow-y: auto;
  }

  .dropdown-item {
    width: 100%;
    text-align: left;
    padding: 0.5rem 1rem;
    cursor: pointer;
    border-bottom: 1px solid var(--border);
    background: none;
    border-left: none;
    border-right: none;
    border-top: none;
    border-radius: 0;
    color: var(--foreground);
  }

  .dropdown-item:last-child {
    border-bottom: none;
  }

  .dropdown-item:hover {
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

  .dropdown-empty {
    padding: 0.75rem 1rem;
    font-size: 0.875rem;
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
