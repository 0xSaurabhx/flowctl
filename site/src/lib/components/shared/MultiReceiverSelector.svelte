<script lang="ts">
    import { apiClient } from "$lib/apiClient";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import type { User, Group } from "$lib/types";
    import { IconUsers, IconUser, IconX, IconMail } from "@tabler/icons-svelte";

    let {
        selectedReceivers = $bindable([]),
        disabled = false,
    }: {
        selectedReceivers: string[];
        disabled?: boolean;
    } = $props();

    let searchQuery = $state("");
    let searchType = $state<"user" | "group">("user");
    let searchResults = $state<(User | Group)[]>([]);
    let loading = $state(false);
    let inputEl: HTMLInputElement;
    let popoverEl: HTMLDivElement;

    interface SelectedItem {
        type: "user" | "group";
        id: string;
        name: string;
        value: string; // The formatted string: email or group:name
    }

    function isValidEmail(email: string): boolean {
        if (email.length < 3 || email.length > 254) return false;
        return /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(
            email,
        );
    }

    // Show "Add email" option when user types a valid email that's not already selected
    let showAddEmailOption = $derived(
        searchType === "user" &&
            isValidEmail(searchQuery) &&
            !selectedReceivers.includes(searchQuery),
    );

    // Derived state: parse receivers into items for display
    let selectedItems = $derived.by(() => {
        if (!selectedReceivers || selectedReceivers.length === 0) {
            return [];
        }
        return selectedReceivers.map((r) => {
            // Check if it's a group (has "group:" prefix) or a user (no prefix)
            if (r.startsWith("group:")) {
                const groupName = r.substring(6); // Remove "group:" prefix
                return {
                    type: "group" as const,
                    id: groupName,
                    name: groupName,
                    value: r,
                };
            } else {
                // User email (no prefix needed)
                return {
                    type: "user" as const,
                    id: r,
                    name: r,
                    value: r,
                };
            }
        });
    });

    async function loadSubjects() {
        loading = true;
        try {
            if (searchType === "user") {
                const response = await apiClient.users.list({
                    filter: searchQuery,
                    count_per_page: 20,
                });
                searchResults = response.users || [];
            } else {
                const response = await apiClient.groups.list({
                    filter: searchQuery,
                    count_per_page: 20,
                });
                searchResults = response.groups || [];
            }
        } catch (error) {
            handleInlineError(
                error,
                searchType === "user"
                    ? "Unable to Load Users"
                    : "Unable to Load Groups",
            );
            searchResults = [];
            closeDropdown();
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
        await loadSubjects();
        openDropdown();
    }

    async function handleFocus() {
        if (searchResults.length === 0) {
            await loadSubjects();
        }
        openDropdown();
    }

    async function handleTypeChange() {
        searchQuery = "";
        searchResults = [];
        await loadSubjects();
        openDropdown();
    }

    function selectSubject(subject: User | Group) {
        const isUser = "username" in subject;
        const name = isUser
            ? (subject as User).username
            : (subject as Group).name;
        // Users don't need a prefix, groups use "group:" prefix
        const value = isUser ? name : `group:${name}`;

        // Check if already selected
        if (selectedReceivers.includes(value)) {
            return;
        }

        // Update the bindable array
        selectedReceivers = [...selectedReceivers, value];

        searchQuery = "";
        closeDropdown();
        searchResults = [];
    }

    function addCustomEmail() {
        if (
            !isValidEmail(searchQuery) ||
            selectedReceivers.includes(searchQuery)
        )
            return;
        selectedReceivers = [...selectedReceivers, searchQuery];
        searchQuery = "";
        closeDropdown();
    }

    function removeReceiver(index: number) {
        selectedReceivers = selectedReceivers.filter((_, i) => i !== index);
    }

    // Close dropdown when clicking outside
    function handleOutsideClick(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".multi-receiver-selector") && !target.closest('[popover]')) {
            closeDropdown();
        }
    }
</script>

<svelte:window onclick={handleOutsideClick} />

<div class="multi-receiver-selector">
    <!-- Search Input -->
    <div class="mb-2">
        <div class="hstack gap-2 mb-2">
            <select
                bind:value={searchType}
                onchange={handleTypeChange}
                {disabled}
            >
                <option value="user">User</option>
                <option value="group">Group</option>
            </select>

            <div class="search-field">
                <input
                    type="text"
                    bind:this={inputEl}
                    bind:value={searchQuery}
                    oninput={handleInput}
                    onfocus={handleFocus}
                    placeholder={searchType === "user"
                        ? "Search users or enter email..."
                        : "Search groups..."}
                    aria-busy={loading}
                    autocomplete="off"
                    {disabled}
                />

                {#if !loading}
                    <svg class="search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                {/if}
            </div>
        </div>
    </div>

    <!-- Popover dropdown — renders in top layer, escapes dialog overflow -->
    <div
        bind:this={popoverEl}
        popover="manual"
        class="selector-popover"
        role="listbox"
    >
        {#if showAddEmailOption}
            <button
                type="button"
                class="selector-item email-add"
                onclick={addCustomEmail}
                role="option"
            >
                <div class="hstack">
                    <div class="item-icon">
                        <IconMail size={16} />
                    </div>
                    <div>
                        <div class="item-name">Add "{searchQuery}"</div>
                        <div class="text-lighter item-desc">Add external email</div>
                    </div>
                </div>
            </button>
        {/if}
        {#if searchResults.length > 0}
            {#each searchResults as subject}
                <button
                    type="button"
                    class="selector-item"
                    onclick={() => selectSubject(subject)}
                    role="option"
                >
                    <div class="hstack">
                        <div class="item-icon">
                            {#if searchType === "user"}
                                <IconUser size={16} />
                            {:else}
                                <IconUsers size={16} />
                            {/if}
                        </div>
                        <div>
                            <div class="item-name">
                                {"name" in subject
                                    ? subject.name
                                    : subject.username}
                            </div>
                            <div class="text-lighter item-desc">{subject.id}</div>
                        </div>
                    </div>
                </button>
            {/each}
        {:else if !loading && !showAddEmailOption}
            <div class="selector-empty text-lighter">
                {searchType === "user"
                    ? "No users found. Enter a valid email to add."
                    : "No groups found"}
            </div>
        {/if}
    </div>

    <!-- Selected Receivers -->
    {#if selectedItems.length > 0}
        <div class="selected-items">
            {#each selectedItems as item, index (item.value)}
                <div class="receiver-chip">
                    <span class="chip-icon">
                        {#if item.type === "user"}
                            <IconUser size={12} />
                        {:else}
                            <IconUsers size={12} />
                        {/if}
                    </span>
                    <span>{item.name}</span>
                    <button
                        type="button"
                        onclick={() => removeReceiver(index)}
                        {disabled}
                        aria-label="Remove {item.name}"
                    >
                        <IconX size={12} />
                    </button>
                </div>
            {/each}
        </div>
    {:else}
        <div class="text-lighter empty-text">
            No receivers selected
        </div>
    {/if}
</div>

<style>
    .search-field {
        position: relative;
        flex: 1;
    }

    .search-field input {
        width: 100%;
        padding-right: 2.5rem;
    }

    .search-icon {
        width: 1.25rem;
        height: 1.25rem;
        position: absolute;
        right: 0.75rem;
        top: 50%;
        transform: translateY(-50%);
        color: var(--muted-foreground);
        pointer-events: none;
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

    .selector-item.email-add {
        background: var(--faint);
    }

    .selector-item.email-add:hover {
        background: color-mix(in srgb, var(--primary) 10%, transparent);
    }

    .item-icon {
        width: 2rem;
        height: 2rem;
        border-radius: 0.5rem;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 0.75rem;
        background: var(--faint);
        color: var(--primary);
    }

    .item-name {
        font-size: 0.875rem;
        font-weight: 500;
        color: var(--foreground);
    }

    .item-desc {
        font-size: 0.75rem;
    }

    .selector-empty {
        padding: var(--space-3) var(--space-4);
        font-size: var(--text-7);
        text-align: center;
    }

    .selected-items {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
    }

    .receiver-chip {
        display: inline-flex;
        align-items: center;
        gap: 0.25rem;
        padding: 0.25rem 0.75rem;
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: 0.375rem;
        font-size: 0.875rem;
    }

    .chip-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--muted-foreground);
    }

    .receiver-chip button {
        all: unset;
        cursor: pointer;
        display: flex;
        margin-left: 0.25rem;
        color: var(--muted-foreground);
    }

    .receiver-chip button:hover {
        color: var(--foreground);
    }

    .empty-text {
        text-align: center;
        font-size: 0.875rem;
    }
</style>
