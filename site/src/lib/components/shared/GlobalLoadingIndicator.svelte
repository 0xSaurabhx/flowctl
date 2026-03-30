<script lang="ts">
	import { beforeNavigate, afterNavigate } from '$app/navigation';
	import { page } from '$app/stores';
	import { isLoading } from '$lib/stores/auth';
	import Logo from './Logo.svelte';

	let navigating = $state(false);

	// Use $derived for reactive store values (Svelte 5 pattern)
	let initialLoading = $derived($isLoading);
	let currentPath = $derived($page.url.pathname);
	let isRootPage = $derived(currentPath === '/');
	let isLoginPage = $derived(currentPath === '/login' || currentPath.startsWith('/login/'));

	beforeNavigate(() => {
		navigating = true;
	});

	afterNavigate(() => {
		navigating = false;
	});
</script>

{#if !isRootPage && !isLoginPage && initialLoading}
	<!-- Splash screen for initial page load -->
	<div
		role="status"
		aria-live="polite"
		aria-label="Loading application"
		class="loading-overlay"
	>
		<div class="vstack items-center gap-6">
			<Logo height="4rem" />
			<div aria-busy="true"></div>
		</div>
	</div>
{:else if !isRootPage && !isLoginPage && navigating}
	<!-- Loading indicator for navigation between pages -->
	<div
		role="status"
		aria-live="polite"
		aria-label="Loading page"
		class="loading-overlay translucent"
	>
		<div class="vstack items-center gap-3">
			<div aria-busy="true"></div>
		</div>
	</div>
{/if}

<style>
	.loading-overlay {
		position: fixed;
		inset: 0;
		z-index: 50;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--card);
	}

	.loading-overlay.translucent {
		background: color-mix(in srgb, var(--card) 60%, transparent);
		backdrop-filter: blur(4px);
	}
</style>
