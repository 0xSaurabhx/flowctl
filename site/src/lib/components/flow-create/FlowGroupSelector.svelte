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
    let showDropdown = $state(false);
    let loading = $state(false);

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

    async function handleFocus() {
        if (groups.length === 0) {
            await loadGroups();
        }
        showDropdown = true;
    }

    function selectGroup(name: string) {
        value = name;
        searchQuery = "";
        showDropdown = false;
    }

    function clear() {
        value = "";
        searchQuery = "";
    }

    function handleOutsideClick(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".flow-group-selector")) {
            showDropdown = false;
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
                bind:value={searchQuery}
                oninput={loadGroups}
                onfocus={handleFocus}
                placeholder={allowCreate ? "Search or create a group..." : "Search groups..."}
                autocomplete="off"
            />

            {#if loading}
                <div class="spinner-overlay" aria-busy="true"></div>
            {/if}

            {#if showDropdown}
                <div class="dropdown-menu">
                    {#if showCreateOption}
                        <button
                            type="button"
                            class="dropdown-item create-item"
                            onclick={() => selectGroup(searchQuery.trim())}
                        >
                            <IconPlus size={16} class="create-icon" />
                            <span>Create "{searchQuery.trim()}"</span>
                        </button>
                    {/if}
                    {#if filteredGroups.length > 0}
                        {#each filteredGroups as group}
                            <button
                                type="button"
                                class="dropdown-item"
                                onclick={() => selectGroup(group.prefix)}
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
                        <div class="dropdown-empty text-lighter">
                            {allowCreate ? "No groups found. Type to create one." : "No groups found."}
                        </div>
                    {/if}
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
    .dropdown-menu {
        position: absolute;
        z-index: 10;
        width: 100%;
        margin-top: 0.25rem;
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: 0.5rem;
        box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1);
        max-height: 12rem;
        overflow-y: auto;
    }
    .dropdown-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        width: 100%;
        padding: 0.5rem 1rem;
        border: none;
        background: none;
        cursor: pointer;
        text-align: left;
        border-bottom: 1px solid var(--border);
    }
    .dropdown-item:last-child {
        border-bottom: none;
    }
    .dropdown-item:hover {
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
    .dropdown-empty {
        padding: 0.75rem 1rem;
        font-size: 0.875rem;
        text-align: center;
    }
    .field-hint {
        font-size: 0.75rem;
        margin-top: 0.25rem;
    }
</style>
