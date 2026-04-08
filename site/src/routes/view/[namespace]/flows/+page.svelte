<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import { apiClient } from "$lib/apiClient";
    import Header from "$lib/components/shared/Header.svelte";
    import Table from "$lib/components/shared/Table.svelte";
    import Pagination from "$lib/components/shared/Pagination.svelte";
    import SearchInput from "$lib/components/shared/SearchInput.svelte";
    import PageHeader from "$lib/components/shared/PageHeader.svelte";
    import { handleInlineError, showSuccess } from "$lib/utils/errorHandling";
    import type { TableColumn, TableAction, FlowListItem, FlowsPaginateResponse, FlowGroupResp } from "$lib/types";
    import { FLOWS_PER_PAGE } from "$lib/constants";
    import {
        permissionChecker,
        type ResourcePermissions,
    } from "$lib/utils/permissions";
    import DeleteModal from "$lib/components/shared/DeleteModal.svelte";
    import GroupEditModal from "$lib/components/shared/GroupEditModal.svelte";
    import { IconPlus, IconBolt, IconFolder, IconChevronLeft } from "@tabler/icons-svelte";
    import NameLinkCell from "$lib/components/shared/cells/NameLinkCell.svelte";
    import MutedTextCell from "$lib/components/shared/cells/MutedTextCell.svelte";

    interface FlowTableRow {
        _kind: 'group' | 'flow';
        name: string;
        description: string;
        slug: string;
        id: string;
        prefix: string;
        step_count: number;
        flow_count: number;
    }

    let { data } = $props();
    let searchValue = $state("");
    let flows = $state<FlowListItem[]>([]);
    let groups = $state<FlowGroupResp[]>([]);
    let pageCount = $state(0);
    let totalCount = $state(0);
    let currentPage = $state(data.currentPage);
    let loading = $state(true);
    let activeGroup = $state<string | null>(data.group || null);
    let permissions = $state<ResourcePermissions>({
        canCreate: false,
        canRead: false,
        canUpdate: false,
        canDelete: false,
    });
    let showDeleteModal = $state(false);
    let flowToDelete = $state<FlowTableRow | null>(null);
    let showEditGroupModal = $state(false);
    let groupToEdit = $state<FlowTableRow | null>(null);

    // Handle the async data from load function
    $effect(() => {
        let cancelled = false;
        activeGroup = data.group || null;

        if (data.group) {
            // Inside a group — load group flows
            data.groupFlowsPromise
                ?.then((result) => {
                    if (!cancelled) {
                        flows = result.flows || [];
                        groups = [];
                        pageCount = 0;
                        totalCount = flows.length;
                        loading = false;
                    }
                })
                .catch((err: Error) => {
                    if (!cancelled) {
                        flows = [];
                        handleInlineError(err, "Unable to Load Group Flows");
                        loading = false;
                    }
                });
        } else {
            // Root view — load groups and paginated flows
            data.groupsPromise
                ?.then((result) => {
                    if (!cancelled) {
                        groups = result.groups || [];
                    }
                })
                .catch(() => {
                    if (!cancelled) groups = [];
                });

            data.flowsPromise
                .then((result: FlowsPaginateResponse) => {
                    if (!cancelled) {
                        flows = result.flows;
                        pageCount = result.page_count;
                        totalCount = result.total_count;
                        loading = false;
                    }
                })
                .catch((err: Error) => {
                    if (!cancelled) {
                        flows = [];
                        pageCount = 0;
                        totalCount = 0;
                        handleInlineError(err, "Unable to Load Flows");
                        loading = false;
                    }
                });
        }

        return () => {
            cancelled = true;
        };
    });

    const goToEditFlow = (flowSlug: string) => {
        goto(`/view/${page.params.namespace}/flows/${flowSlug}/edit`);
    };

    const handleDeleteFlow = (flow: FlowTableRow) => {
        flowToDelete = flow;
        showDeleteModal = true;
    };

    const confirmDeleteFlow = async () => {
        if (!flowToDelete) return;

        try {
            if (flowToDelete._kind === 'group') {
                await apiClient.flows.groups.delete(
                    page.params.namespace!,
                    flowToDelete.id,
                );
                showSuccess(
                    "Group Deleted",
                    `Group "${flowToDelete.name}" and all its flows have been deleted successfully`,
                );
                const result = await apiClient.flows.groups.list(page.params.namespace!);
                groups = result.groups || [];
            } else {
                await apiClient.flows.delete(
                    page.params.namespace!,
                    flowToDelete.slug,
                );
                showSuccess(
                    "Flow Deleted",
                    `Flow "${flowToDelete.name}" has been deleted successfully`,
                );
                if (activeGroup) {
                    await loadGroupFlows(activeGroup);
                } else {
                    await loadFlows(searchValue, currentPage);
                }
            }
        } catch (err) {
            handleInlineError(err, flowToDelete._kind === 'group' ? "Unable to Delete Group" : "Unable to Delete Flow");
        } finally {
            showDeleteModal = false;
            flowToDelete = null;
        }
    };

    const cancelDelete = () => {
        showDeleteModal = false;
        flowToDelete = null;
    };

    const handleEditGroup = (row: FlowTableRow) => {
        groupToEdit = row;
        showEditGroupModal = true;
    };

    const saveGroupEdit = async (data: { description: string }) => {
        if (!groupToEdit) return;

        await apiClient.flows.groups.update(
            page.params.namespace!,
            groupToEdit.id,
            { name: groupToEdit.prefix, description: data.description },
        );
        showSuccess(
            "Group Updated",
            `Group "${groupToEdit.prefix}" has been updated successfully`,
        );
        // Reload groups
        const result = await apiClient.flows.groups.list(page.params.namespace!);
        groups = result.groups || [];
        showEditGroupModal = false;
        groupToEdit = null;
    };

    const cancelEditGroup = () => {
        showEditGroupModal = false;
        groupToEdit = null;
    };

    // Check permissions on mount
    const checkPermissions = async () => {
        permissions = await permissionChecker(
            data.user!,
            "flow",
            data.namespaceId,
            ["create", "update", "delete"],
            "_",
        );
    };

    const handleAdd = () => {
        goto(`/view/${page.params.namespace}/flows/create`);
    };

    const handleDuplicateFlow = (row: FlowTableRow) => {
        goto(`/view/${page.params.namespace}/flows/create?duplicate_from=${row.slug}`);
    };

    checkPermissions();

    const navigateToGroup = (prefix: string) => {
        goto(`/view/${page.params.namespace}/flows?group=${encodeURIComponent(prefix)}`);
    };

    const navigateToRoot = () => {
        goto(`/view/${page.params.namespace}/flows`);
    };

    const loadGroupFlows = async (group: string) => {
        loading = true;
        try {
            const result = await apiClient.flows.groups.get(page.params.namespace!, group);
            flows = result.flows || [];
            totalCount = flows.length;
            pageCount = 0;
        } catch (err) {
            handleInlineError(err, "Unable to Load Group Flows");
        } finally {
            loading = false;
        }
    };

    const loadFlows = async (filter: string = "", pageNumber: number = 1) => {
        loading = true;

        try {
            const result = await apiClient.flows.list(page.params.namespace!, {
                filter,
                page: pageNumber,
                count_per_page: FLOWS_PER_PAGE,
            });

            flows = result.flows;
            pageCount = result.page_count;
            totalCount = result.total_count;
            currentPage = pageNumber;
        } catch (err) {
            handleInlineError(err, "Unable to Load Flows List");
        } finally {
            loading = false;
        }
    };

    const handleSearch = (query: string) => {
        searchValue = query;
        loadFlows(query, 1);
    };

    const goToPage = (pageNum: number) => {
        loadFlows(searchValue.trim(), pageNum);
    };

    const handlePageChange = (event: CustomEvent<{ page: number }>) => {
        goToPage(event.detail.page);
    };

    const isSearching = $derived(searchValue.trim().length > 0);

    const tableData = $derived.by(() => {
        const rows: FlowTableRow[] = [];

        if (activeGroup || isSearching) {
            for (const f of flows) {
                rows.push({
                    _kind: 'flow',
                    name: f.name,
                    description: f.description,
                    slug: f.slug,
                    id: f.id,
                    prefix: f.prefix,
                    step_count: f.step_count,
                    flow_count: 0,
                });
            }
        } else {
            // Build a lookup from the groups API response for metadata (description, flow_count, id)
            const groupMeta = new Map(groups.map(g => [g.prefix, g]));
            const seenPrefixes = new Set<string>();

            for (const f of flows) {
                if (f.prefix) {
                    if (!seenPrefixes.has(f.prefix)) {
                        seenPrefixes.add(f.prefix);
                        const meta = groupMeta.get(f.prefix);
                        rows.push({
                            _kind: 'group',
                            name: f.prefix,
                            description: meta?.description || '',
                            prefix: f.prefix,
                            flow_count: meta?.flow_count ?? 0,
                            slug: '',
                            id: meta?.id || '',
                            step_count: 0,
                        });
                    }
                } else {
                    rows.push({
                        _kind: 'flow',
                        name: f.name,
                        description: f.description,
                        slug: f.slug,
                        id: f.id,
                        prefix: f.prefix,
                        step_count: f.step_count,
                        flow_count: 0,
                    });
                }
            }
        }

        return rows;
    });

    const columns: TableColumn<FlowTableRow>[] = [
        {
            key: "name",
            header: "Name",
            sortable: true,
            component: NameLinkCell,
            componentProps: {
                getIcon: (row: FlowTableRow) => row._kind === 'group' ? IconFolder : IconBolt,
                href: (row: FlowTableRow) => row._kind === 'group'
                    ? `/view/${page.params.namespace}/flows?group=${encodeURIComponent(row.prefix)}`
                    : `/view/${page.params.namespace}/flows/${row.slug}`,
                subtitle: (row: FlowTableRow) => row._kind === 'group'
                    ? `${row.flow_count} flow${row.flow_count !== 1 ? 's' : ''}`
                    : row.prefix ? row.prefix : undefined,
            }
        },
        {
            key: "description",
            header: "Description",
            component: MutedTextCell,
            componentProps: { truncate: 40 }
        },
    ];

    const actions = $derived.by(() => {
        const actionsList: TableAction<FlowTableRow>[] = [];

        const isFlow = (row: FlowTableRow) => row._kind === 'flow';

        if (permissions.canUpdate) {
            actionsList.push({
                label: "Edit",
                onClick: (row: FlowTableRow) => {
                    if (row._kind === 'group') {
                        handleEditGroup(row);
                    } else {
                        goToEditFlow(row.slug);
                    }
                },
                className: "text-link",
            });
        }

        if (permissions.canCreate) {
            actionsList.push({
                label: "Duplicate",
                onClick: (row: FlowTableRow) => handleDuplicateFlow(row),
                visible: (row: FlowTableRow) => row._kind === 'flow',
            });
        }

        if (permissions.canDelete) {
            actionsList.push({
                label: "Delete",
                onClick: (row: FlowTableRow) => handleDeleteFlow(row),
                className: "text-danger-600",
            });
        }

        return actionsList;
    });

    // Breadcrumbs
    const breadcrumbs = $derived.by(() => {
        const crumbs = [
            { label: page.params.namespace! },
            { label: "Flows", url: `/view/${page.params.namespace}/flows` },
        ];
        if (data.group) {
            crumbs.push({ label: data.group });
        }
        return crumbs;
    });
</script>

<svelte:head>
    <title>{activeGroup ? `${activeGroup} - ` : ''}Flows - {page.params.namespace} - Flowctl</title>
</svelte:head>

<Header breadcrumbs={breadcrumbs}>
    {#snippet children()}
        {#if !activeGroup}
            <SearchInput
                bind:value={searchValue}
                placeholder="Search flows..."
                {loading}
                onSearch={handleSearch}
            />
        {/if}
    {/snippet}
</Header>

<!-- Page Content -->
<div class="page-content">
    <PageHeader
        title={activeGroup ? activeGroup : "Flows"}
        subtitle={activeGroup ? `Flows in the ${activeGroup} group` : "Manage and run your workflows"}
        actions={permissions.canCreate && !activeGroup
            ? [
                  {
                      label: "Add",
                      onClick: handleAdd,
                      variant: "primary",
                      IconComponent: IconPlus,
                      iconSize: 16,
                  },
              ]
            : []}
    />

    <!-- Back to root button (inside group) -->
    {#if activeGroup}
        <button
            onclick={navigateToRoot}
            class="back-btn mb-4"
            data-variant="secondary"
        >
            <IconChevronLeft size={16} />
            Back to all flows
        </button>
    {/if}

    <!-- Flows Table -->
    <Table
        {columns}
        data={tableData}
        actions={actions}
        {loading}
        emptyMessage={activeGroup
            ? `No flows in the "${activeGroup}" group`
            : searchValue
                ? "Try adjusting your search"
                : "No flows are available in this namespace"}
        emptyActionLabel={!activeGroup && !searchValue && permissions.canCreate ? "Create your first flow" : undefined}
        onEmptyAction={!activeGroup && !searchValue && permissions.canCreate ? handleAdd : undefined}
        EmptyIconComponent={IconBolt}
        emptyIconSize={48}
    />

    <!-- Pagination and Count (root view only) -->
    {#if !activeGroup && flows.length > 0}
        <div class="mt-6 hstack justify-between items-center">
            <div class="text-light text-sm">
                Showing {flows.length} of {totalCount} flows
            </div>
            <Pagination
                {currentPage}
                totalPages={pageCount}
                {loading}
                on:page-change={handlePageChange}
            />
        </div>
    {/if}

    {#if activeGroup && flows.length > 0}
        <div class="mt-6 text-light text-sm">
            {flows.length} flow{flows.length !== 1 ? 's' : ''} in this group
        </div>
    {/if}
</div>

<!-- Delete Modal -->
{#if showDeleteModal && flowToDelete}
    <DeleteModal
        title={flowToDelete._kind === 'group' ? "Delete Group" : "Delete Flow"}
        itemName={flowToDelete.name}
        description={flowToDelete._kind === 'group' ? "This will permanently delete all flows in this group." : null}
        onConfirm={confirmDeleteFlow}
        onClose={cancelDelete}
    />
{/if}

<!-- Edit Group Modal -->
{#if showEditGroupModal && groupToEdit}
    <GroupEditModal
        groupName={groupToEdit.prefix}
        description={groupToEdit.description}
        onSave={saveGroupEdit}
        onClose={cancelEditGroup}
    />
{/if}

