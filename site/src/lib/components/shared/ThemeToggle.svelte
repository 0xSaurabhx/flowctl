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
    data-variant="secondary"
    class="theme-btn"
    class:full-width={!collapsed}
    title="{$theme} theme"
>
    <svelte:component this={icons[$theme]} size={20} />
    {#if !collapsed}
        <span class="theme-label">{$theme}</span>
    {/if}
</button>

<style>
    .theme-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0.5rem;
        background: none;
        border: none;
        color: var(--muted-foreground);
        cursor: pointer;
        border-radius: 0.5rem;
    }

    .theme-btn:hover {
        background: var(--faint);
    }

    .theme-btn.full-width {
        width: 100%;
        padding: 0.5rem 1rem;
    }

    .theme-label {
        margin-left: 0.75rem;
        text-transform: capitalize;
    }
</style>
