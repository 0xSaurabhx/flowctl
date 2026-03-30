<script lang="ts">
    import { apiClient } from "$lib/apiClient";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import type { FlowGroupResp } from "$lib/types";
    import { IconFolder, IconX, IconPlus } from "@tabler/icons-svelte";

    let {
        namespace,
        value = $bindable(""),
        allowCreate = true,
    }: {
        namespace: string;
        value: string;
        allowCreate?: boolean;
    } = $props();

    let searchQuery = $state("");
    let groups = $state<FlowGroupResp[]>([]);
    let loading = $state(false);
    let inputEl: HTMLInputElement;
    let popoverEl: HTMLDivElement;

    let filteredGroups = $derived(
        groups.filter((g) =>
            g.prefix.toLowerCase().includes(searchQuery.toLowerCase()),
        ),
    );

    let showCreateOption = $derived(
        allowCreate &&
            searchQuery.trim() !== "" &&
            !groups.some(
                (g) => g.prefix.toLowerCase() === searchQuery.toLowerCase(),
            ),
    );

    async function loadGroups() {
        loading = true;
        try {
            const result = await apiClient.flows.groups.list(namespace);
            groups = result.groups || [];
        } catch (error) {
            handleInlineError(error, "Unable to load flow groups");
            groups = [];
        } finally {
            loading = false;
        }
    }

    function positionPopover() {
        if (!inputEl || !popoverEl) return;
        const r = inputEl.getBoundingClientRect();
        popoverEl.style.top = `${r.bottom + 4}px`;
        popoverEl.style.left = `${r.left}px`;
        popoverEl.style.width = `${r.width}px`;
    }

    function openDropdown() {
        try { popoverEl?.showPopover(); } catch {}
        positionPopover();
    }

    function closeDropdown() {
        try { popoverEl?.hidePopover(); } catch {}
    }

    async function handleInput() {
        await loadGroups();
        openDropdown();
    }

    async function handleFocus() {
        if (groups.length === 0) {
            await loadGroups();
        }
        openDropdown();
    }

    function selectGroup(name: string) {
        value = name;
        searchQuery = "";
        closeDropdown();
    }

    function clear() {
        value = "";
        searchQuery = "";
    }

    function handleOutsideClick(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".flow-group-selector") && !target.closest('[popover]')) {
            closeDropdown();
        }
    }
</script>

<svelte:window onclick={handleOutsideClick} />

<div class="flow-group-selector">
    <label for="flow-group">
        Group
        <span class="text-lighter label-hint">(optional)</span>
    </label>

    {#if value}
        <div class="selected-group hstack gap-2">
            <IconFolder size={16} class="text-lighter" />
            <span class="selected-name">{value}</span>
            <button
                type="button"
                data-variant="secondary"
                class="clear-btn"
                onclick={clear}
            >
                <IconX size={16} />
            </button>
        </div>
    {:else}
        <div class="dropdown-wrapper">
            <input
                type="text"
                id="flow-group"
                bind:this={inputEl}
                bind:value={searchQuery}
                oninput={handleInput}
                onfocus={handleFocus}
                placeholder={allowCreate ? "Search or create a group..." : "Search groups..."}
                autocomplete="off"
            />

            {#if loading}
                <div class="spinner-overlay" aria-busy="true"></div>
            {/if}
        </div>

        <!-- Popover dropdown — renders in top layer, escapes dialog overflow -->
        <div
            bind:this={popoverEl}
            popover="manual"
            class="selector-popover"
            role="listbox"
        >
            {#if showCreateOption}
                <button
                    type="button"
                    class="selector-item create-item"
                    onclick={() => selectGroup(searchQuery.trim())}
                    role="option"
                >
                    <IconPlus size={16} class="create-icon" />
                    <span>Create "{searchQuery.trim()}"</span>
                </button>
            {/if}
            {#if filteredGroups.length > 0}
                {#each filteredGroups as group}
                    <button
                        type="button"
                        class="selector-item"
                        onclick={() => selectGroup(group.prefix)}
                        role="option"
                    >
                        <IconFolder size={16} class="text-lighter" />
                        <div>
                            <div class="group-name">{group.prefix}</div>
                            {#if group.flow_count > 0}
                                <div class="text-lighter group-count">
                                    {group.flow_count} flow{group.flow_count !== 1 ? "s" : ""}
                                </div>
                            {/if}
                        </div>
                    </button>
                {/each}
            {:else if !loading && !showCreateOption}
                <div class="selector-empty text-lighter">
                    {allowCreate ? "No groups found. Type to create one." : "No groups found."}
                </div>
            {/if}
        </div>
    {/if}

    <p class="text-lighter field-hint">
        Assign this flow to a group for organization. Groups are created
        automatically.
    </p>
</div>

<style>
    .flow-group-selector {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }
    .label-hint {
        font-weight: normal;
        font-size: 0.875rem;
    }
    .selected-group {
        padding: 0.5rem 0.75rem;
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: 0.375rem;
    }
    .selected-name {
        flex: 1;
        font-size: 0.875rem;
        color: var(--foreground);
    }
    .clear-btn {
        padding: 0.125rem;
        border: none;
        background: none;
        color: var(--muted-foreground);
        cursor: pointer;
    }
    .clear-btn:hover {
        color: var(--foreground);
    }
    .dropdown-wrapper {
        position: relative;
    }
    .spinner-overlay {
        position: absolute;
        right: 0.75rem;
        top: 50%;
        transform: translateY(-50%);
    }
    .selector-popover {
        position: fixed;
        margin: 0;
        padding: 0;
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius-medium);
        box-shadow: var(--shadow-large);
        max-height: 14rem;
        overflow-y: auto;
    }
    .selector-item {
        all: unset;
        box-sizing: border-box;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        width: 100%;
        padding: var(--space-2) var(--space-4);
        cursor: pointer;
        border-bottom: 1px solid var(--border);
        font-size: var(--text-7);
    }
    .selector-item:last-child {
        border-bottom: none;
    }
    .selector-item:hover {
        background: var(--faint);
    }
    .create-item {
        border-bottom: 1px solid var(--border);
    }
    .create-icon {
        color: var(--primary);
    }
    .group-name {
        font-size: 0.875rem;
        font-weight: 500;
        color: var(--foreground);
    }
    .group-count {
        font-size: 0.75rem;
    }
    .selector-empty {
        padding: var(--space-3) var(--space-4);
        font-size: var(--text-7);
        text-align: center;
    }
    .field-hint {
        font-size: 0.75rem;
        margin-top: 0.25rem;
    }
</style>
