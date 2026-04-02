<script lang="ts">
    import { onMount } from "svelte";
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
    } from "@tabler/icons-svelte";
    import UserDropdown from "./UserDropdown.svelte";
    import ThemeToggle from "./ThemeToggle.svelte";
    import Logo from "./Logo.svelte";
    import { APP_VERSION, APP_COMMIT } from "$lib/constants";

    let { namespace }: { namespace: string } = $props();

    let namespaces = $state<Namespace[]>([]);
    let searchQuery = $state("");
    let searchResults = $state<Namespace[]>([]);
    let currentNamespace = $state(page.params.namespace || namespace);
    let currentNamespaceId = $state<string>("");
    let searchLoading = $state(false);
    let nsSearchInput: HTMLInputElement;
    let permissions = $state<{
        flows: ResourcePermissions;
        nodes: ResourcePermissions;
        credentials: ResourcePermissions;
        secrets: ResourcePermissions;
        members: ResourcePermissions;
        approvals: ResourcePermissions;
        history: ResourcePermissions;
    }>({
        flows: { canCreate: false, canUpdate: false, canDelete: false, canRead: false },
        nodes: { canCreate: false, canUpdate: false, canDelete: false, canRead: false },
        credentials: { canCreate: false, canUpdate: false, canDelete: false, canRead: false },
        secrets: { canCreate: false, canUpdate: false, canDelete: false, canRead: false },
        members: { canCreate: false, canUpdate: false, canDelete: false, canRead: false },
        approvals: { canCreate: false, canUpdate: false, canDelete: false, canRead: false },
        history: { canCreate: false, canUpdate: false, canDelete: false, canRead: false },
    });

    const isActiveLink = (section: string): boolean => {
        const currentPath = page.url.pathname;
        return currentPath.includes(`/${section}`);
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

                const currentNs = namespaces.find((ns) => ns.name === namespace);
                if (currentNs) {
                    currentNamespaceId = currentNs.id;
                } else if (namespaces.length > 0) {
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
        } else {
            searchResults = namespaces;
        }
    };

    const handleNsToggle = (e: ToggleEvent) => {
        if (e.newState === "open") {
            searchQuery = "";
            searchResults = namespaces;
            setTimeout(() => nsSearchInput?.focus(), 0);
        }
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
        searchQuery = "";

        // Close the popover
        document.getElementById("ns-select")?.hidePopover();

        if (ns.name === namespace) return;

        clearPermissionCache();
        selectedNamespace.set(ns.name);
        window.location.href = `/view/${ns.name}/flows`;
    };

    setContext("namespace", namespace);

    onMount(() => {
        fetchNamespaces();
        checkAllPermissions();
    });

    $effect(() => {
        currentNamespace = page.params.namespace || namespace;
        if (currentNamespace) {
            selectedNamespace.set(currentNamespace);
        }
    });

    $effect(() => {
        if ($currentUser && currentNamespaceId) {
            checkAllPermissions();
        }
    });
</script>

<aside data-sidebar aria-label="Main navigation">
    <header>
        <a href="/" class="unstyled vstack items-center gap-1">
            <Logo height="2rem" />
            <span class="text-lighter" style="font-size: var(--text-8)">{APP_VERSION}-{APP_COMMIT}</span>
        </a>

        <!-- Namespace Dropdown -->
        <ot-dropdown class="mt-4" style="display: block">
            <button popovertarget="ns-select" class="outline small w-100 hstack justify-between">
                <span>{currentNamespace || "Select namespace"}</span>
                <IconChevronDown size={14} />
            </button>
            <div popover id="ns-select" ontoggle={handleNsToggle}>
                <div style="padding: var(--space-2); border-bottom: 1px solid var(--border)">
                    <input
                        type="search"
                        bind:this={nsSearchInput}
                        bind:value={searchQuery}
                        oninput={handleSearchInput}
                        placeholder="Search namespaces..."
                        aria-busy={searchLoading}
                        autocomplete="off"
                    />
                </div>
                {#each searchResults as ns (ns.id)}
                    <button
                        role="menuitem"
                        aria-selected={ns.name === namespace}
                        onclick={() => selectNamespace(ns)}
                    >
                        {ns.name}
                    </button>
                {/each}
                {#if searchResults.length === 0 && !searchLoading}
                    <div class="text-lighter align-center" style="padding: var(--space-3); font-size: var(--text-7)">
                        {searchQuery ? "No namespaces found" : "No namespaces available"}
                    </div>
                {/if}
                {#if searchLoading}
                    <div class="text-lighter align-center" style="padding: var(--space-3); font-size: var(--text-7)">
                        Searching...
                    </div>
                {/if}
            </div>
        </ot-dropdown>
    </header>

    <nav>
        <ul>
            {#if permissions.flows.canRead}
                <li><a href="/view/{namespace}/flows" aria-current={isActiveLink("flows") ? "page" : undefined}>
                    <IconGridDots size={20} aria-hidden="true" /> Flows
                </a></li>
            {/if}
            {#if permissions.nodes.canRead}
                <li><a href="/view/{namespace}/nodes" aria-current={isActiveLink("nodes") ? "page" : undefined}>
                    <IconServer size={20} aria-hidden="true" /> Nodes
                </a></li>
            {/if}
            {#if permissions.credentials.canRead}
                <li><a href="/view/{namespace}/credentials" aria-current={isActiveLink("credentials") ? "page" : undefined}>
                    <IconKey size={20} aria-hidden="true" /> Credentials
                </a></li>
            {/if}
            {#if permissions.secrets.canRead}
                <li><a href="/view/{namespace}/secrets" aria-current={isActiveLink("secrets") ? "page" : undefined}>
                    <IconLock size={20} aria-hidden="true" /> Secrets
                </a></li>
            {/if}
            {#if permissions.members.canRead}
                <li><a href="/view/{namespace}/members" aria-current={isActiveLink("members") ? "page" : undefined}>
                    <IconUsers size={20} aria-hidden="true" /> Members
                </a></li>
            {/if}
            {#if permissions.approvals.canRead}
                <li><a href="/view/{namespace}/approvals" aria-current={isActiveLink("approvals") ? "page" : undefined}>
                    <IconCircleCheck size={20} aria-hidden="true" /> Approvals
                </a></li>
            {/if}
            {#if permissions.history.canRead}
                <li><a href="/view/{namespace}/history" aria-current={isActiveLink("history") ? "page" : undefined}>
                    <IconClock size={20} aria-hidden="true" /> History
                </a></li>
            {/if}
        </ul>
    </nav>

    <footer>
        <div class="hstack justify-between items-center">
            {#if $currentUser && $currentUser.role === "superuser"}
                <a href="/settings" class="unstyled" aria-current={isActiveLink("settings") ? "page" : undefined} title="Settings">
                    <IconSettings size={20} aria-hidden="true" />
                </a>
            {/if}
            <ThemeToggle collapsed={true} />
        </div>
        {#if $currentUser}
            <UserDropdown />
        {/if}
    </footer>
</aside>
