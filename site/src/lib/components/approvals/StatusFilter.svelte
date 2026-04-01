<script lang="ts">
  let {
    value = $bindable(''),
    onChange
  }: {
    value: string;
    onChange: (status: string) => void;
  } = $props();

  const options = [
    { value: '', label: 'All Statuses' },
    { value: 'pending', label: 'Pending' },
    { value: 'approved', label: 'Approved' },
    { value: 'rejected', label: 'Rejected' },
  ];

  const selectedLabel = $derived(options.find(o => o.value === value)?.label ?? 'All Statuses');

  const popoverId = 'status-filter-menu';

  function select(opt: string) {
    value = opt;
    onChange(opt);
    // Close the popover
    document.getElementById(popoverId)?.hidePopover();
  }
</script>

<ot-dropdown>
  <button popovertarget={popoverId} data-variant="secondary" class="hstack gap-2 text-sm">
    {selectedLabel}
  </button>
  <div id={popoverId} popover="auto" role="menu">
    {#each options as opt}
      <button role="menuitem" aria-selected={value === opt.value} onclick={() => select(opt.value)}>
        {opt.label}
      </button>
    {/each}
  </div>
</ot-dropdown>
