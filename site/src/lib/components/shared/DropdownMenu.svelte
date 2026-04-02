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

  const popoverId = `dropdown-${Math.random().toString(36).slice(2, 9)}`;

  function handleItemClick(item: MenuItem) {
    item.onClick();
    document.getElementById(popoverId)?.hidePopover();
  }
</script>

<ot-dropdown>
  <button popovertarget={popoverId} class="ghost icon small" aria-label="Actions menu">
    <IconDotsVertical size={16} />
  </button>
  <div popover id={popoverId} role="menu">
    {#each items as item}
      <button
        role="menuitem"
        onclick={() => handleItemClick(item)}
        style={item.variant === 'danger' ? 'color: var(--danger)' : ''}
      >
        {item.label}
      </button>
    {/each}
  </div>
</ot-dropdown>
