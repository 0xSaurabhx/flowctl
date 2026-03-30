<script lang="ts" generics="T">
    import type { TableColumn, TableAction } from "$lib/types";
    import type { ComponentType } from "svelte";
    import LoadingSpinner from "./LoadingSpinner.svelte";

    type SortDirection = "asc" | "desc" | null;

    type Props = {
        columns: TableColumn<T>[];
        data: T[];
        onRowClick?: (row: T) => void;
        actions?: TableAction<T>[];
        loading?: boolean;
        emptyMessage?: string;
        emptyIcon?: string;
        EmptyIconComponent?: ComponentType;
        emptyIconSize?: number;
        title?: string;
        subtitle?: string;
    };

    let {
        columns,
        data,
        onRowClick,
        actions = [],
        loading = false,
        emptyMessage = "No data available",
        emptyIcon,
        EmptyIconComponent,
        emptyIconSize = 64,
        title,
        subtitle,
    }: Props = $props();

    let sortKey = $state<string | null>(null);
    let sortDirection = $state<SortDirection>(null);
    let openMenuIndex = $state<number | null>(null);
    let menuPosition = $state<{ top: number; left: number }>({ top: 0, left: 0 });

    const toggleMenu = (index: number, event: Event) => {
        event.stopPropagation();
        if (openMenuIndex === index) {
            openMenuIndex = null;
            return;
        }
        const button = event.currentTarget as HTMLElement;
        const rect = button.getBoundingClientRect();
        menuPosition = {
            top: rect.bottom + 4,
            left: rect.right - 144, // 144 = w-36 (9rem)
        };
        openMenuIndex = index;
    };

    const closeMenu = () => {
        openMenuIndex = null;
    };

    const handleClickOutside = (event: MouseEvent) => {
        if (openMenuIndex !== null) {
            const target = event.target as HTMLElement;
            if (!target.closest('.actions-menu')) {
                closeMenu();
            }
        }
    };

    const getValue = (row: T, column: TableColumn<T>) => {
        const keys = column.key.split(".");
        let value = row as any;

        for (const key of keys) {
            if (value && typeof value === "object") {
                value = value[key];
            } else {
                return undefined;
            }
        }

        return value;
    };

    const renderValue = (row: T, column: TableColumn<T>) => {
        const value = getValue(row, column);

        if (column.render) {
            return column.render(value, row);
        }

        return value ?? "";
    };

    const handleRowClick = (row: T) => {
        if (onRowClick) {
            onRowClick(row);
        }
    };

    const handleActionClick = (
        action: TableAction<T>,
        row: T,
        event: Event,
    ) => {
        event.stopPropagation();
        action.onClick(row, event);
    };

    const handleSort = (column: TableColumn<T>) => {
        if (!column.sortable) return;

        if (sortKey === column.key) {
            // Cycle through: asc -> desc -> null
            if (sortDirection === "asc") {
                sortDirection = "desc";
            } else if (sortDirection === "desc") {
                sortDirection = null;
                sortKey = null;
            } else {
                sortDirection = "asc";
            }
        } else {
            sortKey = column.key;
            sortDirection = "asc";
        }
    };

    const sortedData = $derived.by(() => {
        if (!sortKey || !sortDirection) return data;

        const column = columns.find((c) => c.key === sortKey);
        if (!column) return data;

        return [...data].sort((a, b) => {
            const aValue = getValue(a, column);
            const bValue = getValue(b, column);

            // Handle null/undefined values
            if (aValue == null && bValue == null) return 0;
            if (aValue == null) return sortDirection === "asc" ? 1 : -1;
            if (bValue == null) return sortDirection === "asc" ? -1 : 1;

            // Convert to strings for comparison if not already numbers
            const aCompare =
                typeof aValue === "number"
                    ? aValue
                    : String(aValue).toLowerCase();
            const bCompare =
                typeof bValue === "number"
                    ? bValue
                    : String(bValue).toLowerCase();

            if (aCompare < bCompare) return sortDirection === "asc" ? -1 : 1;
            if (aCompare > bCompare) return sortDirection === "asc" ? 1 : -1;
            return 0;
        });
    });
</script>

<div class="card">
    {#if title || subtitle}
        <div class="table-header">
            {#if title}
                <h3>{title}</h3>
            {/if}
            {#if subtitle}
                <p class="text-light text-sm">{subtitle}</p>
            {/if}
        </div>
    {/if}

    {#if loading}
        <div class="empty-state" role="status" aria-live="polite">
            <LoadingSpinner label="Loading..." />
        </div>
    {:else if data.length === 0}
        <div class="empty-state">
            {#if EmptyIconComponent}
                <div class="text-light mb-4">
                    <EmptyIconComponent size={emptyIconSize} />
                </div>
            {:else if emptyIcon}
                {@html emptyIcon}
            {:else}
                <svg
                    class="empty-icon"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
                    />
                </svg>
            {/if}
            <h3>{emptyMessage}</h3>
        </div>
    {:else}
        <div class="table">
            <table>
                <thead>
                <tr>
                    {#each columns as column}
                        <th
                            scope="col"
                            class={column.sortable ? 'sortable' : ''}
                            onclick={() => handleSort(column)}
                            aria-sort={sortKey === column.key
                                ? sortDirection === "asc"
                                    ? "ascending"
                                    : "descending"
                                : undefined}
                        >
                            <div class="hstack gap-1">
                                <span>{column.header}</span>
                                {#if column.sortable}
                                    <div class="sort-indicators" aria-hidden="true">
                                        <svg
                                            class="sort-arrow {sortKey === column.key && sortDirection === 'asc' ? 'active' : ''}"
                                            fill="currentColor"
                                            viewBox="0 0 20 20"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                        <svg
                                            class="sort-arrow sort-arrow-down {sortKey === column.key && sortDirection === 'desc' ? 'active' : ''}"
                                            fill="currentColor"
                                            viewBox="0 0 20 20"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                    </div>
                                {/if}
                            </div>
                        </th>
                    {/each}
                    {#if actions.length > 0}
                        <th scope="col">
                            <span class="sr-only">Actions</span>
                        </th>
                    {/if}
                </tr>
            </thead>
            <tbody>
                {#each sortedData as row}
                    <tr
                        class={onRowClick ? 'clickable' : ''}
                        onclick={() => handleRowClick(row)}
                    >
                        {#each columns as column}
                            <td>
                                {#if column.component}
                                    {@const Component = column.component}
                                    <Component
                                        {row}
                                        value={getValue(row, column)}
                                        {...(column.componentProps || {})}
                                    />
                                {:else}
                                    {@html renderValue(row, column)}
                                {/if}
                            </td>
                        {/each}
                        {#if actions.length > 0}
                            {@const rowIndex = sortedData.indexOf(row)}
                            {@const visibleActions = actions.filter(a => !a.visible || a.visible(row))}
                            <td class="actions-cell">
                                {#if visibleActions.length > 0}
                                    <div class="actions-menu">
                                        <button
                                            onclick={(e) => toggleMenu(rowIndex, e)}
                                            class="actions-trigger"
                                            aria-label="Actions"
                                            aria-haspopup="true"
                                            aria-expanded={openMenuIndex === rowIndex}
                                        >
                                            <svg width="20" height="20" fill="currentColor" viewBox="0 0 20 20">
                                                <circle cx="10" cy="4" r="1.5" />
                                                <circle cx="10" cy="10" r="1.5" />
                                                <circle cx="10" cy="16" r="1.5" />
                                            </svg>
                                        </button>
                                    </div>
                                {/if}
                            </td>
                        {/if}
                    </tr>
                {/each}
            </tbody>
        </table>
        </div>
    {/if}
</div>

{#if openMenuIndex !== null}
    {@const row = sortedData[openMenuIndex]}
    {@const visibleActions = actions.filter(a => !a.visible || a.visible(row))}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="menu-backdrop"
        onclick={closeMenu}
        onkeydown={(e) => { if (e.key === 'Escape') closeMenu(); }}
    ></div>
    <div
        class="actions-menu context-menu"
        style="top: {menuPosition.top}px; left: {menuPosition.left}px;"
    >
        {#each visibleActions as action}
            {#if action.href}
                <a
                    href={action.href(row)}
                    class="context-menu-item {action.className || ''}"
                    onclick={() => closeMenu()}
                >
                    {action.label}
                </a>
            {:else}
                <button
                    onclick={(e) => { handleActionClick(action, row, e); closeMenu(); }}
                    class="context-menu-item {action.className || ''}"
                >
                    {action.label}
                </button>
            {/if}
        {/each}
    </div>
{/if}

<style>
    .table-header {
        padding: var(--space-4) var(--space-6);
        border-bottom: 1px solid var(--border);
        background: var(--muted);
    }

    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 16rem;
        text-align: center;
    }

    .empty-icon {
        width: 4rem;
        height: 4rem;
        color: var(--muted-foreground);
        margin-bottom: var(--space-4);
    }

    th.sortable {
        cursor: pointer;
        user-select: none;
    }

    .sort-indicators {
        display: flex;
        flex-direction: column;
    }

    .sort-arrow {
        width: 0.75rem;
        height: 0.75rem;
        color: var(--muted-foreground);
    }

    .sort-arrow-down {
        margin-top: -0.25rem;
    }

    .sort-arrow.active {
        color: var(--primary);
    }

    tr.clickable {
        cursor: pointer;
    }

    .actions-cell {
        width: 4rem;
        text-align: right;
    }

    .actions-trigger {
        all: unset;
        cursor: pointer;
        padding: var(--space-1);
        border-radius: var(--radius-small);
        color: var(--muted-foreground);
        display: inline-flex;
    }

    .actions-trigger:hover {
        background: var(--muted);
        color: var(--foreground);
    }

    .menu-backdrop {
        position: fixed;
        inset: 0;
        z-index: 40;
    }

    .context-menu {
        position: fixed;
        z-index: 50;
        width: 9rem;
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius-medium);
        box-shadow: var(--shadow-large);
        padding: var(--space-1) 0;
    }

    .context-menu-item {
        all: unset;
        display: block;
        width: 100%;
        text-align: left;
        padding: var(--space-2) var(--space-4);
        font-size: 0.875rem;
        color: var(--foreground);
        cursor: pointer;
        box-sizing: border-box;
    }

    .context-menu-item:hover {
        background: var(--muted);
    }

    .sr-only {
        position: absolute;
        width: 1px;
        height: 1px;
        padding: 0;
        margin: -1px;
        overflow: hidden;
        clip: rect(0, 0, 0, 0);
        white-space: nowrap;
        border-width: 0;
    }
</style>
