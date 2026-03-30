<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { apiClient } from "$lib/apiClient";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { currentUser } from "$lib/stores/auth";
    import { selectedNamespace } from "$lib/stores/namespace";
    import type { Namespace } from "$lib/types";
    import { DEFAULT_PAGE_SIZE } from "$lib/constants";
    import { setContext } from "svelte";
    import {
        permissionChecker,
        clearPermissionCache,
        type ResourcePermissions,
    } from "$lib/utils/permissions";
    import {
        IconChevronDown,
        IconGridDots,
        IconServer,
        IconKey,
        IconLock,
        IconUsers,
        IconCircleCheck,
        IconClock,
        IconSettings,
        IconChevronsLeft,
        IconChevronsRight,
    } from "@tabler/icons-svelte";
    import UserDropdown from "./UserDropdown.svelte";
    import ThemeToggle from "./ThemeToggle.svelte";
    import Logo from "./Logo.svelte";
    import { APP_VERSION, APP_COMMIT } from "$lib/constants";

    let { namespace }: { namespace: string } = $props();

    let namespaceDropdownOpen = $state(false);
    let namespaces = $state<Namespace[]>([]);
    let searchQuery = $state("");
    let searchResults = $state<Namespace[]>([]);
    let currentNamespace = $state(page.params.namespace || namespace);
    let currentNamespaceId = $state<string>("");
    let searchLoading = $state(false);
    let isCollapsed = $state(false);
    let permissions = $state<{
        flows: ResourcePermissions;
        nodes: ResourcePermissions;
        credentials: ResourcePermissions;
        secrets: ResourcePermissions;
        members: ResourcePermissions;
        approvals: ResourcePermissions;
        history: ResourcePermissions;
    }>({
        flows: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        nodes: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        credentials: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        secrets: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        members: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        approvals: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        history: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
    });

    const isActiveLink = (section: string): boolean => {
        const currentPath = page.url.pathname;

        if (section === "flows") {
            return currentPath.includes("/flows");
        } else if (section === "nodes") {
            return currentPath.includes("/nodes");
        } else if (section === "credentials") {
            return currentPath.includes("/credentials");
        } else if (section === "secrets") {
            return currentPath.includes("/secrets");
        } else if (section === "members") {
            return currentPath.includes("/members");
        } else if (section === "approvals") {
            return currentPath.includes("/approvals");
        } else if (section === "history") {
            return currentPath.includes("/history");
        } else if (section === "settings") {
            return currentPath.includes("/settings");
        }

        return false;
    };

    const fetchNamespaces = async (filter = "") => {
        try {
            searchLoading = true;
            const data = await apiClient.namespaces.list({
                count_per_page: DEFAULT_PAGE_SIZE,
                filter: filter,
            });
            const results = data.namespaces || [];

            if (filter) {
                searchResults = results;
            } else {
                namespaces = results;
                searchResults = results;

                // Set current namespace ID
                const currentNs = namespaces.find(
                    (ns) => ns.name === namespace,
                );
                if (currentNs) {
                    currentNamespaceId = currentNs.id;
                } else if (namespaces.length > 0) {
                    // If namespace not found, use first available namespace
                    currentNamespaceId = namespaces[0].id;
                }
            }
        } catch (error) {
            handleInlineError(error, "Unable to Load Namespaces");
            if (filter) {
                searchResults = [];
            } else {
                namespaces = [];
                searchResults = [];
            }
        } finally {
            searchLoading = false;
        }
    };

    const handleSearchInput = async () => {
        if (searchQuery.trim()) {
            await fetchNamespaces(searchQuery);
            namespaceDropdownOpen = true;
        } else {
            searchResults = namespaces;
            namespaceDropdownOpen = false;
        }
    };

    const handleSearchFocus = () => {
        if (!searchQuery.trim()) {
            searchResults = namespaces;
        }
        namespaceDropdownOpen = true;
    };

    const checkAllPermissions = async () => {
        if (!$currentUser || !currentNamespaceId) return;

        const resourceMappings: Record<string, { resource: string; prefix?: string }> = {
            flows: { resource: "flow", prefix: "_" },
            nodes: { resource: "node" },
            credentials: { resource: "credential" },
            secrets: { resource: "namespace_secret" },
            members: { resource: "member" },
            approvals: { resource: "approval" },
            history: { resource: "execution" },
        };

        try {
            const permissionPromises = Object.entries(resourceMappings).map(
                async ([frontendKey, { resource: backendResource, prefix }]) => {
                    const perms = await permissionChecker(
                        $currentUser,
                        backendResource,
                        currentNamespaceId,
                        ["view"],
                        prefix,
                    );
                    return { frontendKey, perms };
                },
            );

            const results = await Promise.all(permissionPromises);

            results.forEach(({ frontendKey, perms }) => {
                permissions[frontendKey as keyof typeof permissions] = perms;
            });
        } catch (error) {
            handleInlineError(error, "Unable to Check Sidebar Permissions");
        }
    };

    const selectNamespace = (ns: Namespace) => {
        namespaceDropdownOpen = false;
        searchQuery = "";

        // Don't navigate if already on the same namespace
        if (ns.name === namespace) {
            return;
        }

        // Clear cached permissions so the new namespace gets fresh checks
        clearPermissionCache();

        // Save selected namespace to store
        selectedNamespace.set(ns.name);

        // Force a full page reload by using window.location
        window.location.href = `/view/${ns.name}/flows`;
    };

    const toggleCollapse = () => {
        isCollapsed = !isCollapsed;
        if (isCollapsed) {
            namespaceDropdownOpen = false;
            searchQuery = "";
        }
    };

    // Handle escape key and outside clicks
    function handleOutsideClick(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".namespace-dropdown")) {
            namespaceDropdownOpen = false;
            searchQuery = "";
            searchResults = namespaces;
        }
    }

    // Set initial context
    setContext("namespace", namespace);

    onMount(() => {
        fetchNamespaces();
        checkAllPermissions();
    });

    // Update currentNamespace when namespace prop changes
    $effect(() => {
        currentNamespace = page.params.namespace || namespace;
        // Also save to store whenever namespace changes
        if (currentNamespace) {
            selectedNamespace.set(currentNamespace);
        }
    });

    // Re-check permissions when currentUser or namespace changes
    $effect(() => {
        if ($currentUser && currentNamespaceId) {
            checkAllPermissions();
        }
    });
</script>

<svelte:window on:click={handleOutsideClick} />

<!-- Sidebar Navigation -->
<aside
    class="sidebar"
    class:collapsed={isCollapsed}
    data-sidebar
    aria-label="Main navigation"
>
    <!-- Logo -->
    <a href="/" class="logo-link">
        <div class="logo-area">
            {#if isCollapsed}
                <Logo height="1.5rem" iconOnly={true} />
            {:else}
                <Logo height="2rem" />
                <div class="text-lighter version-text">{APP_VERSION}-{APP_COMMIT}</div>
            {/if}
        </div>
    </a>

    <!-- Namespace Dropdown -->
    {#if !isCollapsed}
        <div class="namespace-section namespace-dropdown">
            <div class="namespace-field">
                <label for="namespace-search" class="namespace-label">Namespace</label>
                <div class="input-wrapper">
                    <input
                        type="text"
                        id="namespace-search"
                        bind:value={searchQuery}
                        oninput={handleSearchInput}
                        onfocus={handleSearchFocus}
                        placeholder={currentNamespace || "Search namespaces..."}
                        aria-busy={searchLoading}
                        autocomplete="off"
                    />

                    {#if searchLoading}
                        <span class="input-icon" aria-busy="true"></span>
                    {:else}
                        <IconChevronDown
                            class="input-icon chevron {namespaceDropdownOpen ? 'open' : ''}"
                            size={16}
                        />
                    {/if}
                </div>

                <!-- Dropdown Menu -->
                {#if namespaceDropdownOpen}
                    <div
                        class="dropdown-menu"
                        role="listbox"
                        aria-label="Namespace selection"
                    >
                        {#each searchResults as ns (ns.id)}
                            <button
                                type="button"
                                role="option"
                                aria-selected={ns.name === namespace}
                                onclick={() => selectNamespace(ns)}
                                class="dropdown-item"
                                class:active={ns.name === namespace}
                            >
                                {ns.name}
                            </button>
                        {/each}
                        {#if searchResults.length === 0 && !searchLoading}
                            <div class="dropdown-empty text-lighter">
                                {searchQuery
                                    ? "No namespaces found"
                                    : "No namespaces available"}
                            </div>
                        {/if}
                        {#if searchLoading}
                            <div class="dropdown-empty text-lighter">
                                Searching...
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
        </div>
    {/if}

    <!-- Navigation -->
    <ul class="nav-list" role="list">
        {#if permissions.flows.canRead}
            <li>
                <a
                    href="/view/{namespace}/flows"
                    class="nav-link"
                    class:icon-only={isCollapsed}
                    aria-current={isActiveLink("flows") ? "page" : undefined}
                    title={isCollapsed ? "Flows" : ""}
                >
                    <IconGridDots size={20} aria-hidden="true" />
                    {#if !isCollapsed}
                        Flows
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.nodes.canRead}
            <li>
                <a
                    href="/view/{namespace}/nodes"
                    class="nav-link"
                    class:icon-only={isCollapsed}
                    aria-current={isActiveLink("nodes") ? "page" : undefined}
                    title={isCollapsed ? "Nodes" : ""}
                >
                    <IconServer size={20} aria-hidden="true" />
                    {#if !isCollapsed}
                        Nodes
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.credentials.canRead}
            <li>
                <a
                    href="/view/{namespace}/credentials"
                    class="nav-link"
                    class:icon-only={isCollapsed}
                    aria-current={isActiveLink("credentials") ? "page" : undefined}
                    title={isCollapsed ? "Credentials" : ""}
                >
                    <IconKey size={20} aria-hidden="true" />
                    {#if !isCollapsed}
                        Credentials
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.secrets.canRead}
            <li>
                <a
                    href="/view/{namespace}/secrets"
                    class="nav-link"
                    class:icon-only={isCollapsed}
                    aria-current={isActiveLink("secrets") ? "page" : undefined}
                    title={isCollapsed ? "Secrets" : ""}
                >
                    <IconLock size={20} aria-hidden="true" />
                    {#if !isCollapsed}
                        Secrets
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.members.canRead}
            <li>
                <a
                    href="/view/{namespace}/members"
                    class="nav-link"
                    class:icon-only={isCollapsed}
                    aria-current={isActiveLink("members") ? "page" : undefined}
                    title={isCollapsed ? "Members" : ""}
                >
                    <IconUsers size={20} aria-hidden="true" />
                    {#if !isCollapsed}
                        Members
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.approvals.canRead}
            <li>
                <a
                    href="/view/{namespace}/approvals"
                    class="nav-link"
                    class:icon-only={isCollapsed}
                    aria-current={isActiveLink("approvals") ? "page" : undefined}
                    title={isCollapsed ? "Approvals" : ""}
                >
                    <IconCircleCheck size={20} aria-hidden="true" />
                    {#if !isCollapsed}
                        Approvals
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.history.canRead}
            <li>
                <a
                    href="/view/{namespace}/history"
                    class="nav-link"
                    class:icon-only={isCollapsed}
                    aria-current={isActiveLink("history") ? "page" : undefined}
                    title={isCollapsed ? "History" : ""}
                >
                    <IconClock size={20} aria-hidden="true" />
                    {#if !isCollapsed}
                        History
                    {/if}
                </a>
            </li>
        {/if}
    </ul>

    <div class="sidebar-footer">
        <div class="hstack {isCollapsed ? 'footer-collapsed' : ''}">
            {#if $currentUser && $currentUser.role === "superuser"}
                <a
                    href="/settings"
                    class="nav-link icon-only"
                    aria-current={isActiveLink("settings") ? "page" : undefined}
                    title="Settings"
                >
                    <IconSettings size={20} aria-hidden="true" />
                </a>
            {/if}

            <ThemeToggle collapsed={true} />

            <button
                type="button"
                onclick={toggleCollapse}
                data-variant="secondary"
                class="collapse-btn"
                aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
                title={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
            >
                {#if isCollapsed}
                    <IconChevronsRight size={20} aria-hidden="true" />
                {:else}
                    <IconChevronsLeft size={20} aria-hidden="true" />
                {/if}
            </button>
        </div>
    </div>

    <!-- User Profile Section -->
    {#if $currentUser}
        <div class="user-section">
            <UserDropdown {isCollapsed} />
        </div>
    {/if}
</aside>

<style>
    .sidebar {
        position: relative;
        background: var(--card);
        border-right: 1px solid var(--border);
        display: flex;
        flex-direction: column;
        transition: width 0.3s ease-in-out;
        width: 15rem;
    }

    .sidebar.collapsed {
        width: 5rem;
    }

    .logo-link {
        text-decoration: none;
    }

    .logo-area {
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .version-text {
        font-size: 0.75rem;
        margin-top: 0.25rem;
    }

    /* Namespace section */
    .namespace-section {
        padding: 0 1rem;
        margin-bottom: 1rem;
    }

    .namespace-field {
        position: relative;
    }

    .namespace-label {
        display: block;
        font-size: 0.75rem;
        font-weight: 500;
        color: var(--muted-foreground);
        margin-bottom: 0.25rem;
        text-transform: uppercase;
    }

    .input-wrapper {
        position: relative;
    }

    .input-wrapper input {
        width: 100%;
        padding-right: 2rem;
    }

    .input-wrapper :global(.input-icon) {
        position: absolute;
        right: 0.75rem;
        top: 50%;
        transform: translateY(-50%);
        color: var(--muted-foreground);
        pointer-events: none;
    }

    .input-wrapper :global(.input-icon.chevron) {
        transition: transform 0.2s;
    }

    .input-wrapper :global(.input-icon.chevron.open) {
        transform: translateY(-50%) rotate(180deg);
    }

    /* Dropdown menu */
    .dropdown-menu {
        position: absolute;
        z-index: 50;
        width: 100%;
        margin-top: 0.25rem;
        background: var(--card);
        border-radius: 0.5rem;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        border: 1px solid var(--border);
        max-height: 12rem;
        overflow-y: auto;
    }

    .dropdown-item {
        width: 100%;
        text-align: left;
        padding: 0.5rem 0.75rem;
        font-size: 0.875rem;
        color: var(--foreground);
        background: none;
        border: none;
        border-radius: 0;
        cursor: pointer;
    }

    .dropdown-item:hover {
        background: var(--faint);
    }

    .dropdown-item.active {
        background: var(--faint);
        color: var(--primary);
    }

    .dropdown-empty {
        padding: 0.5rem 0.75rem;
        font-size: 0.875rem;
        text-align: center;
    }

    /* Navigation */
    .nav-list {
        flex: 1;
        padding: 0 1rem;
        list-style: none;
        margin: 0;
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .nav-link {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        font-size: 0.875rem;
        font-weight: 500;
        padding: 0.75rem 1rem;
        border-radius: 0.5rem;
        color: var(--foreground);
        text-decoration: none;
        transition: background 0.15s;
    }

    .nav-link:hover {
        background: var(--faint);
    }

    .nav-link[aria-current="page"] {
        background: var(--faint);
        color: var(--primary);
    }

    .nav-link.icon-only {
        justify-content: center;
        padding: 0.75rem;
    }

    /* Footer */
    .sidebar-footer {
        padding: 0.5rem 1rem;
        border-top: 1px solid var(--border);
    }

    .sidebar-footer .hstack {
        justify-content: space-between;
        align-items: center;
    }

    .footer-collapsed {
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
    }

    .collapse-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        padding: 0.5rem;
        background: none;
        border: none;
        color: var(--muted-foreground);
        cursor: pointer;
        border-radius: 0.5rem;
    }

    .collapse-btn:hover {
        background: var(--faint);
    }

    .user-section {
        padding: 1rem;
        border-top: 1px solid var(--border);
    }
</style>
