<script lang="ts">
  import { onMount } from 'svelte';
  import { apiClient } from '$lib/apiClient';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import type { NodeResp } from '$lib/types';
  import { IconChevronDown, IconServer, IconTag, IconX } from '@tabler/icons-svelte';

  let {
    namespace,
    selectedNodes = $bindable([]),
    placeholder = 'Search nodes or use tag:name...',
    disabled = false,
    multiple = true
  }: {
    namespace: string;
    selectedNodes: string[];
    placeholder?: string;
    disabled?: boolean;
    multiple?: boolean;
  } = $props();

  let searchQuery = $state('');
  let searchResults = $state<NodeResp[]>([]);
  let showDropdown = $state(false);
  let loading = $state(false);
  let searchTimeout: ReturnType<typeof setTimeout> | null = null;
  let isTagSearch = $state(false);
  let currentTagQuery = $state('');

  async function searchNodes(filter: string = '') {
    loading = true;
    isTagSearch = false;
    currentTagQuery = '';

    // Check if this is a tag search
    if (filter.startsWith('tag:')) {
      const tagName = filter.slice(4).trim();
      isTagSearch = true;
      currentTagQuery = tagName;

      if (tagName) {
        try {
          const response = await apiClient.nodes.list(namespace, {
            count_per_page: 100,
            tags: [tagName]
          });
          searchResults = response.nodes || [];
        } catch (error) {
          handleInlineError(error, 'Unable to Load Nodes');
          searchResults = [];
        } finally {
          loading = false;
        }
        return;
      } else {
        searchResults = [];
        loading = false;
        return;
      }
    }

    try {
      const response = await apiClient.nodes.list(namespace, {
        count_per_page: 100,
        filter: filter
      });
      searchResults = response.nodes || [];
    } catch (error) {
      handleInlineError(error, 'Unable to Load Nodes');
      searchResults = [];
    } finally {
      loading = false;
    }
  }

  function isTagSelection(item: string): boolean {
    return item.startsWith('tag:');
  }

  async function handleInput() {
    // Debounce search
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }

    searchTimeout = setTimeout(() => {
      searchNodes(searchQuery);
    }, 300);

    showDropdown = true;
  }

  async function handleFocus() {
    if (disabled) return;
    await searchNodes(searchQuery);
    showDropdown = true;
  }

  function selectNode(node: NodeResp) {
    if (multiple) {
      if (!selectedNodes.includes(node.name)) {
        selectedNodes = [...selectedNodes, node.name];
      }
    } else {
      selectedNodes = [node.name];
    }
    searchQuery = '';
    showDropdown = false;
  }

  function selectTag(tagName: string) {
    const tagValue = `tag:${tagName}`;
    if (multiple) {
      if (!selectedNodes.includes(tagValue)) {
        selectedNodes = [...selectedNodes, tagValue];
      }
    } else {
      selectedNodes = [tagValue];
    }
    searchQuery = '';
    showDropdown = false;
  }

  function removeNode(nodeName: string) {
    selectedNodes = selectedNodes.filter(n => n !== nodeName);
  }


  // Load nodes when component mounts (skip if disabled -- no permission to list nodes)
  onMount(() => {
    if (!disabled) {
      searchNodes();
    }
  });

  // Close dropdown when clicking outside
  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.node-selector')) {
      showDropdown = false;
    }
  }
</script>

<svelte:window on:click={handleOutsideClick} />

<div class="node-selector">
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
        <!-- Tag selection option when searching by tag -->
        {#if isTagSearch && currentTagQuery && !selectedNodes.includes(`tag:${currentTagQuery}`)}
          <div
            class="dropdown-item tag-select"
            onclick={() => selectTag(currentTagQuery)}
            role="option"
            tabindex="0"
            onkeydown={(e) => e.key === 'Enter' && selectTag(currentTagQuery)}
          >
            <div class="hstack">
              <div class="item-icon success">
                <IconTag size={16} />
              </div>
              <div>
                <div class="item-name">Select tag: {currentTagQuery}</div>
                <div class="text-lighter item-desc">{searchResults.length} node{searchResults.length !== 1 ? 's' : ''} with this tag</div>
              </div>
            </div>
          </div>
        {/if}

        {#if searchResults.length > 0}
          {#if isTagSearch}
            <div class="dropdown-divider text-lighter">
              Or select individual nodes:
            </div>
          {/if}
          {#each searchResults as node}
            <!-- Skip already selected nodes -->
            {#if !selectedNodes.includes(node.name)}
              <div
                class="dropdown-item"
                onclick={() => selectNode(node)}
                role="option"
                tabindex="0"
                onkeydown={(e) => e.key === 'Enter' && selectNode(node)}
              >
                <div class="hstack">
                  <div class="item-icon">
                    <IconServer size={16} />
                  </div>
                  <div>
                    <div class="item-name">{node.name}</div>
                    <div class="text-lighter item-desc">{node.hostname}:{node.port}</div>
                  </div>
                </div>
              </div>
            {/if}
          {/each}
        {:else if !loading}
          <div class="dropdown-empty text-lighter">
            {searchQuery ? 'No nodes found' : 'No nodes available'}
          </div>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Selected nodes display -->
  {#if selectedNodes.length > 0}
    <div class="selected-items mt-2">
      {#each selectedNodes as nodeName (nodeName)}
        {#if isTagSelection(nodeName)}
          <span class="badge success">
            <IconTag size={12} />
            {nodeName.slice(4)}
            {#if !disabled}
              <button
                type="button"
                onclick={() => removeNode(nodeName)}
                aria-label="Remove {nodeName}"
              >
                <IconX size={12} />
              </button>
            {/if}
          </span>
        {:else}
          <span class="badge">
            {nodeName}
            {#if !disabled}
              <button
                type="button"
                onclick={() => removeNode(nodeName)}
                aria-label="Remove {nodeName}"
              >
                <IconX size={12} />
              </button>
            {/if}
          </span>
        {/if}
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

  .dropdown-item.tag-select:hover {
    background: color-mix(in srgb, var(--success) 10%, transparent);
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

  .item-icon.success {
    background: color-mix(in srgb, var(--success) 15%, transparent);
    color: var(--success);
  }

  .item-name {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--foreground);
  }

  .item-desc {
    font-size: 0.75rem;
  }

  .dropdown-divider {
    padding: 0.25rem 1rem;
    font-size: 0.75rem;
    background: var(--faint);
    border-bottom: 1px solid var(--border);
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

  .selected-items .badge.success button {
    color: var(--success);
  }

  .selected-items .badge.success button:hover {
    color: var(--foreground);
  }
</style>
