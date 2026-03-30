<script lang="ts">
  import { onMount } from 'svelte';
  import { apiClient } from '$lib/apiClient';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import type { Group } from '$lib/types';
  import { IconChevronDown, IconUsersGroup, IconX } from '@tabler/icons-svelte';

  let {
    selectedGroups = $bindable([]),
    placeholder = 'Search groups...',
    disabled = false,
    multiple = true
  }: {
    selectedGroups: Group[];
    placeholder?: string;
    disabled?: boolean;
    multiple?: boolean;
  } = $props();

  let searchQuery = $state('');
  let searchResults = $state<Group[]>([]);
  let allGroups = $state<Group[]>([]);
  let showDropdown = $state(false);
  let loading = $state(false);
  let initialized = $state(false);

  async function loadAllGroups() {
    if (initialized) return;

    loading = true;
    try {
      const response = await apiClient.groups.list({
        count_per_page: 100 // Get more groups for selection
      });
      allGroups = response.groups || [];
      searchResults = allGroups;
      initialized = true;
    } catch (error) {
      handleInlineError(error, 'Unable to Load Groups');
      allGroups = [];
      searchResults = [];
    } finally {
      loading = false;
    }
  }

  function filterGroups() {
    if (!searchQuery.trim()) {
      searchResults = allGroups;
    } else {
      searchResults = allGroups.filter(group =>
        group.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (group.description && group.description.toLowerCase().includes(searchQuery.toLowerCase()))
      );
    }
  }

  async function handleInput() {
    filterGroups();
  }

  async function handleFocus() {
    filterGroups();
    showDropdown = searchResults.length > 0;
  }

  function selectGroup(group: Group) {
    if (multiple) {
      if (!selectedGroups.some(g => g.id === group.id)) {
        selectedGroups = [...selectedGroups, group];
      }
    } else {
      selectedGroups = [group];
    }
    searchQuery = '';
    showDropdown = false;
  }

  function removeGroup(groupId: string) {
    selectedGroups = selectedGroups.filter(g => g.id !== groupId);
  }


  // Load groups when component mounts
  onMount(() => {
    loadAllGroups();
  });

  // Close dropdown when clicking outside
  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.group-selector')) {
      showDropdown = false;
    }
  }
</script>

<svelte:window on:click={handleOutsideClick} />

<div class="group-selector">
  <div class="selector-field">
    <input
      type="text"
      bind:value={searchQuery}
      oninput={handleInput}
      onfocus={handleFocus}
      {placeholder}
      aria-busy={loading}
      autocomplete="off"
      {disabled}
    />

    {#if loading}
      <span class="field-icon" aria-busy="true"></span>
    {:else}
      <IconChevronDown class="field-icon" size={16} />
    {/if}

    <!-- Dropdown -->
    {#if showDropdown}
      <div class="dropdown-menu" role="listbox">
        {#if searchResults.length > 0}
          {#each searchResults as group}
            <!-- Skip already selected groups -->
            {#if !selectedGroups.some(g => g.id === group.id)}
              <div
                class="dropdown-item"
                onclick={() => selectGroup(group)}
                role="option"
                tabindex="0"
                onkeydown={(e) => e.key === 'Enter' && selectGroup(group)}
              >
                <div class="hstack">
                  <div class="item-icon">
                    <IconUsersGroup size={16} />
                  </div>
                  <div>
                    <div class="item-name">{group.name}</div>
                    <div class="text-lighter item-desc">{group.description || 'No description'}</div>
                  </div>
                </div>
              </div>
            {/if}
          {/each}
        {:else if initialized && !loading}
          <div class="dropdown-empty text-lighter">
            {searchQuery ? 'No groups found' : 'No groups available'}
          </div>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Selected groups display -->
  {#if selectedGroups.length > 0}
    <div class="selected-items mt-2">
      {#each selectedGroups as group (group.id)}
        <span class="badge">
          {group.name}
          {#if !disabled}
            <button
              type="button"
              onclick={() => removeGroup(group.id)}
              aria-label="Remove {group.name}"
            >
              <IconX size={12} />
            </button>
          {/if}
        </span>
      {/each}
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

  .selector-field :global(.field-icon) {
    position: absolute;
    right: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--muted-foreground);
    pointer-events: none;
  }

  .dropdown-menu {
    position: absolute;
    z-index: 20;
    width: 100%;
    margin-top: 0.5rem;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    max-height: 12rem;
    overflow-y: auto;
  }

  .dropdown-item {
    padding: 0.5rem 1rem;
    cursor: pointer;
    border-bottom: 1px solid var(--border);
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
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    text-align: center;
  }

  .selected-items {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .selected-items .badge {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
  }

  .selected-items .badge button {
    all: unset;
    cursor: pointer;
    display: flex;
    color: var(--primary);
  }

  .selected-items .badge button:hover {
    color: var(--foreground);
  }
</style>
