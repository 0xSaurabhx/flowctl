<script lang="ts">
	type Size = 'sm' | 'md' | 'lg';

	type Props = {
		size?: Size;
		label?: string;
		class?: string;
	};

	let { size = 'md', label, class: className = '' }: Props = $props();
</script>

<div class="hstack gap-3 justify-center {className}">
	<div aria-busy="true" data-spinner={size === 'sm' ? 'small' : size === 'lg' ? 'large' : undefined}></div>
	{#if label}
		<span class="text-sm text-light">{label}</span>
	{/if}
</div>
