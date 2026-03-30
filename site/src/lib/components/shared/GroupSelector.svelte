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
  let loading = $state(false);
  let initialized = $state(false);
  let inputEl: HTMLInputElement;
  let popoverEl: HTMLDivElement;

  async function loadAllGroups() {
    if (initialized) return;
    loading = true;
    try {
      const response = await apiClient.groups.list({ count_per_page: 100 });
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

  function positionPopover() {
    if (!inputEl || !popoverEl) return;
    const r = inputEl.getBoundingClientRect();
    popoverEl.style.top = `${r.bottom + 4}px`;
    popoverEl.style.left = `${r.left}px`;
    popoverEl.style.width = `${r.width}px`;
  }

  function openDropdown() {
    filterGroups();
    try { popoverEl?.showPopover(); } catch {}
    positionPopover();
  }

  function closeDropdown() {
    try { popoverEl?.hidePopover(); } catch {}
  }

  function handleInput() {
    filterGroups();
    openDropdown();
  }

  function handleFocus() {
    openDropdown();
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
    closeDropdown();
  }

  function removeGroup(groupId: string) {
    selectedGroups = selectedGroups.filter(g => g.id !== groupId);
  }

  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.group-selector') && !target.closest('[popover]')) {
      closeDropdown();
    }
  }

  onMount(() => {
    loadAllGroups();
  });
</script>

<svelte:window on:click={handleOutsideClick} />

<div class="group-selector">
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
    {#if loading}
      <span class="field-icon" aria-busy="true"></span>
    {:else}
      <IconChevronDown class="field-icon" size={16} />
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
      {#each searchResults as group}
        {#if !selectedGroups.some(g => g.id === group.id)}
          <button
            type="button"
            class="selector-item"
            onclick={() => selectGroup(group)}
            role="option"
          >
            <div class="hstack">
              <div class="icon-box">
                <IconUsersGroup size={16} />
              </div>
              <div>
                <div class="item-name">{group.name}</div>
                <div class="text-lighter item-desc">{group.description || 'No description'}</div>
              </div>
            </div>
          </button>
        {/if}
      {/each}
    {:else if initialized && !loading}
      <div class="selector-empty text-lighter">
        {searchQuery ? 'No groups found' : 'No groups available'}
      </div>
    {/if}
  </div>

  {#if selectedGroups.length > 0}
    <div class="selected-items mt-2">
      {#each selectedGroups as group (group.id)}
        <span class="badge">
          {group.name}
          {#if !disabled}
            <button type="button" onclick={() => removeGroup(group.id)} aria-label="Remove {group.name}">
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
  .selector-item:last-child { border-bottom: none; }
  .selector-item:hover { background: var(--faint); }

  .item-name {
    font-size: var(--text-7);
    font-weight: var(--font-medium);
    color: var(--foreground);
  }
  .item-desc { font-size: var(--text-8); }

  .selector-empty {
    padding: var(--space-3) var(--space-4);
    font-size: var(--text-7);
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
