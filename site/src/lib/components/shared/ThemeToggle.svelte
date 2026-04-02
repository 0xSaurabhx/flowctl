<script lang="ts">
    import { theme, type Theme } from '$lib/stores/theme';
    import { IconSun, IconMoon, IconDeviceDesktop } from '@tabler/icons-svelte';

    let { collapsed = false }: { collapsed?: boolean } = $props();

    const icons = { light: IconSun, dark: IconMoon, system: IconDeviceDesktop };
    const order: Theme[] = ['light', 'dark', 'system'];
    const nextLabels: Record<Theme, string> = { system: 'Switch to light mode', light: 'Switch to dark mode', dark: 'Switch to system theme' };

    function cycle() {
        const next = order[(order.indexOf($theme) + 1) % 3];
        theme.set(next);
    }
</script>

<button
    type="button"
    onclick={cycle}
    class={collapsed ? "ghost icon small" : "ghost hstack gap-3 w-100"}
    aria-label={nextLabels[$theme]}
    title={nextLabels[$theme]}
>
    <span class="theme-icon" aria-hidden="true">
        <svelte:component this={icons[$theme]} size={collapsed ? 18 : 20} />
    </span>
    {#if !collapsed}
        <span class="theme-label">{$theme}</span>
    {/if}
</button>

<style>
    .theme-label {
        text-transform: capitalize;
    }

    .theme-icon {
        display: inline-flex;
        transition: transform 0.25s ease;
    }

    button:active .theme-icon {
        transform: rotate(60deg);
    }
</style>
