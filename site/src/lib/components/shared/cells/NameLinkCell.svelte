<script lang="ts">
  import type { ComponentType } from 'svelte';

  let {
    row,
    value,
    icon,
    getIcon,
    iconVariant = '',
    nameKey,
    subtitleKey,
    subtitle,
    onClick,
    href
  }: {
    row: any;
    value: any;
    icon?: ComponentType;
    getIcon?: (row: any) => ComponentType;
    iconVariant?: string | ((row: any) => string);
    nameKey?: string;
    subtitleKey?: string;
    subtitle?: string | ((row: any) => string | undefined);
    onClick?: (row: any) => void;
    href?: string | ((row: any) => string);
  } = $props();

  const name = $derived(nameKey ? row[nameKey] : value);
  const resolvedSubtitle = $derived(
    typeof subtitle === 'function' ? subtitle(row)
    : subtitle ? subtitle
    : subtitleKey ? row[subtitleKey]
    : undefined
  );
  const resolvedHref = $derived(typeof href === 'function' ? href(row) : href);
  const resolvedVariant = $derived(typeof iconVariant === 'function' ? iconVariant(row) : iconVariant);
  const resolvedIcon = $derived(getIcon ? getIcon(row) : icon);
</script>

<div class="hstack gap-2">
  {#if resolvedIcon}
    {@const Icon = resolvedIcon}
    <div class="icon-box {resolvedVariant}" style="width:2.5rem;height:2.5rem">
      <Icon size={20} />
    </div>
  {/if}
  <div>
    {#if onClick}
      <a href="#" class="cell-link" onclick={(e) => { e.preventDefault(); onClick(row); }}>{name}</a>
    {:else if resolvedHref}
      <a href={resolvedHref} class="cell-link">{name}</a>
    {:else}
      <div class="font-medium text-sm">{name}</div>
    {/if}
    {#if resolvedSubtitle}
      <div class="cell-muted">{resolvedSubtitle}</div>
    {/if}
  </div>
</div>
