<script lang="ts">
	import type { ComponentType } from 'svelte';

	interface Props {
		title: string;
		value: number | string;
		icon?: string;
		IconComponent?: ComponentType;
		iconSize?: number;
		color?: 'blue' | 'green' | 'purple' | 'red' | 'yellow' | 'gray';
	}

	let { title, value, icon, IconComponent, iconSize = 24, color = 'blue' }: Props = $props();

	const colorMap: Record<string, string> = {
		blue: '',
		green: 'success',
		purple: '',
		red: 'danger',
		yellow: 'warning',
		gray: 'muted'
	};
</script>

<article class="card p-4">
	<div class="hstack justify-between">
		<div>
			<p class="text-sm text-light">{title}</p>
			<p class="stat-value">{value}</p>
		</div>
		<div class="stat-icon {colorMap[color] || ''}">
			{#if IconComponent}
				<IconComponent size={iconSize} />
			{:else if icon}
				{@html icon}
			{/if}
		</div>
	</div>
</article>

<style>
	.stat-value {
		font-size: var(--text-3);
		font-weight: 700;
	}
	.stat-icon {
		color: var(--primary);
	}
	.stat-icon.success { color: var(--success); }
	.stat-icon.danger { color: var(--danger); }
	.stat-icon.warning { color: var(--warning); }
	.stat-icon.muted { color: var(--muted-foreground); }
</style>
