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
</script>

<div class="card">
	<div class="hstack justify-between">
		<div>
			<p class="text-sm text-light">{title}</p>
			<p class="stat-value">{value}</p>
		</div>
		<div class="stat-icon {color}">
			{#if IconComponent}
				<IconComponent size={iconSize} />
			{:else if icon}
				{@html icon}
			{/if}
		</div>
	</div>
</div>

<style>
	.card {
		padding: var(--space-6);
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
	}

	.stat-icon {
		width: 3rem;
		height: 3rem;
		border-radius: var(--radius-medium);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.stat-icon.blue {
		background: color-mix(in srgb, var(--primary) 15%, transparent);
		color: var(--primary);
	}

	.stat-icon.green {
		background: color-mix(in srgb, var(--success) 15%, transparent);
		color: var(--success);
	}

	.stat-icon.purple {
		background: color-mix(in srgb, var(--primary) 15%, transparent);
		color: var(--primary);
	}

	.stat-icon.red {
		background: color-mix(in srgb, var(--danger) 15%, transparent);
		color: var(--danger);
	}

	.stat-icon.yellow {
		background: color-mix(in srgb, var(--warning) 15%, transparent);
		color: var(--warning);
	}

	.stat-icon.gray {
		background: var(--muted);
		color: var(--muted-foreground);
	}
</style>
