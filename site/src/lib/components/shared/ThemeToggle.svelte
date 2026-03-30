<script lang="ts">
    import { theme, type Theme } from '$lib/stores/theme';
    import { IconSun, IconMoon, IconDeviceDesktop } from '@tabler/icons-svelte';

    let { collapsed = false }: { collapsed?: boolean } = $props();

    const icons = { light: IconSun, dark: IconMoon, system: IconDeviceDesktop };
    const order: Theme[] = ['light', 'dark', 'system'];

    function cycle() {
        const next = order[(order.indexOf($theme) + 1) % 3];
        theme.set(next);
    }
</script>

<button
    type="button"
    onclick={cycle}
    class="ghost hstack gap-3"
    class:w-100={!collapsed}
    title="{$theme} theme"
>
    <svelte:component this={icons[$theme]} size={20} />
    {#if !collapsed}
        <span class="theme-label">{$theme}</span>
    {/if}
</button>

<style>
    .theme-label {
        text-transform: capitalize;
    }
</style>
