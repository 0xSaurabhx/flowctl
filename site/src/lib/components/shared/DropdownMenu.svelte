<script lang="ts">
  import { IconDotsVertical } from '@tabler/icons-svelte';

  interface MenuItem {
    label: string;
    onClick: () => void;
    variant?: 'default' | 'danger';
  }

  let {
    items,
    position = 'right'
  }: {
    items: MenuItem[];
    position?: 'left' | 'right';
  } = $props();

  let isOpen = $state(false);
  let buttonElement: HTMLButtonElement | undefined = $state();
  let dropdownPosition = $state({ top: 0, left: 0, right: 0 });

  function updateDropdownPosition() {
    if (!buttonElement) return;

    const rect = buttonElement.getBoundingClientRect();
    dropdownPosition = {
      top: rect.bottom + 4,
      left: position === 'left' ? rect.left : 0,
      right: position === 'right' ? window.innerWidth - rect.right : 0
    };
  }

  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.dropdown-menu-container')) {
      isOpen = false;
    }
  }

  function handleItemClick(item: MenuItem) {
    item.onClick();
    isOpen = false;
  }

  function toggleDropdown() {
    isOpen = !isOpen;
    if (isOpen) {
      updateDropdownPosition();
    }
  }
</script>

<svelte:window on:click={handleOutsideClick} on:scroll={updateDropdownPosition} on:resize={updateDropdownPosition} />

<div class="dropdown-menu-container">
  <button
    type="button"
    bind:this={buttonElement}
    onclick={toggleDropdown}
    class="dropdown-trigger"
    aria-label="Actions menu"
  >
    <IconDotsVertical class="icon" />
  </button>

  {#if isOpen}
    <div
      class="dropdown-panel"
      style="top: {dropdownPosition.top}px; {position === 'right' ? `right: ${dropdownPosition.right}px` : `left: ${dropdownPosition.left}px`}"
      role="menu"
    >
      {#each items as item}
        <button
          type="button"
          onclick={() => handleItemClick(item)}
          class="dropdown-item {item.variant === 'danger' ? 'danger' : ''}"
          role="menuitem"
        >
          {item.label}
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .dropdown-menu-container {
    position: relative;
  }

  .dropdown-trigger {
    all: unset;
    cursor: pointer;
    padding: var(--space-1);
    border-radius: var(--radius-small);
    display: inline-flex;
    color: var(--muted-foreground);
  }

  .dropdown-trigger:hover {
    background: var(--muted);
  }

  .dropdown-panel {
    position: fixed;
    z-index: 50;
    width: 9rem;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius-medium);
    box-shadow: var(--shadow-large);
    padding: var(--space-1) 0;
    display: flex;
    flex-direction: column;
  }

  .dropdown-item {
    all: unset;
    display: block;
    width: 100%;
    text-align: left;
    padding: var(--space-2) var(--space-4);
    font-size: 0.875rem;
    color: var(--foreground);
    cursor: pointer;
    box-sizing: border-box;
  }

  .dropdown-item:hover {
    background: var(--muted);
  }

  .dropdown-item.danger {
    color: var(--danger);
  }

  .dropdown-item.danger:hover {
    background: color-mix(in srgb, var(--danger) 10%, transparent);
  }
</style>
