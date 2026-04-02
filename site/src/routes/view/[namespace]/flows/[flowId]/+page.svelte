<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import FlowInputForm from "$lib/components/flow-input/FlowInputForm.svelte";
    import Table from "$lib/components/shared/Table.svelte";
    import Header from "$lib/components/shared/Header.svelte";
    import StatusBadge from "$lib/components/shared/StatusBadge.svelte";
    import Tabs from "$lib/components/shared/Tabs.svelte";
    import Pagination from "$lib/components/shared/Pagination.svelte";
    import FlowHero from "$lib/components/flows/FlowHero.svelte";
    import FlowActionsSummary from "$lib/components/flow-input/FlowActionsSummary.svelte";
    import FlowSchedulesList from "$lib/components/flows/FlowSchedulesList.svelte";
    import ScheduledExecutionsList from "$lib/components/flows/ScheduledExecutionsList.svelte";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import type { PageData } from "./$types";
    import type { TableColumn, ScheduledExecution } from "$lib/types";
    import { DEFAULT_PAGE_SIZE } from "$lib/constants";
    import { permissionChecker } from "$lib/utils/permissions";
    import { formatDateTime, getStartTime } from "$lib/utils";
    import { apiClient } from "$lib/apiClient";
    import { IconPencil, IconEye } from "@tabler/icons-svelte";

    let { data }: { data: PageData } = $props();

    let activeTab = $state<"run" | "schedules" | "history">("run");
    let historyLoading = $state(false);
    let flowExecutions = $state<any[]>([]);
    let historyCurrentPage = $state(1);
    let historyItemsPerPage = $state(DEFAULT_PAGE_SIZE);
    let historyTotalCount = $state(0);
    let historyPageCount = $state(0);
    let canUpdateFlow = $state(false);
    let canViewFlowConfig = $state(false);
    let scheduledExecutions = $state<ScheduledExecution[]>(
        data.flowMeta?.scheduled_executions || [],
    );
    let userSchedules = $state<any[]>(data.userSchedules || []);

    let namespace = $derived(page.params.namespace);
    let flowId = $derived(page.params.flowId);
    let rerunFromExecId = $derived(data.rerunFromExecId);
    let showRerunBanner = $state(!!rerunFromExecId);

    // Check update and view_config permissions on mount
    permissionChecker(data.user!, "flow", data.namespaceId, ["update", "view_config"], "_").then(
        (permissions) => {
            canUpdateFlow = permissions.canUpdate;
            canViewFlowConfig = permissions.canViewConfig;
        },
    );

    // Reload schedules after updates
    const reloadSchedules = async () => {
        try {
            const res = await apiClient.flows.schedules.list(
                namespace!,
                flowId!,
            );
            userSchedules = res.schedules || [];
        } catch (error) {
            handleInlineError(error, "Failed to reload schedules");
        }
    };

    const refreshScheduledExecutions = async () => {
        try {
            const flowMeta = await apiClient.flows.getMeta(namespace!, flowId!);
            scheduledExecutions = flowMeta.scheduled_executions || [];
        } catch (error) {
            console.error("Failed to refresh scheduled executions:", error);
        }
    };

    const loadFlowHistory = async () => {
        historyLoading = true;

        try {
            const response = await fetch(
                `/api/v1/${namespace}/flows/${flowId}/executions?page=${historyCurrentPage}&count_per_page=${historyItemsPerPage}`,
                {
                    credentials: "include",
                },
            );
            const result = await response.json();

            if (!response.ok) {
                const apiError = new Error(
                    result.error || "Failed to fetch execution history",
                );
                handleInlineError(apiError, "Unable to Load Flow History");
                flowExecutions = [];
                return;
            }

            flowExecutions = result.executions || [];
            historyTotalCount = result.total_count || 0;
            historyPageCount = result.page_count || 1;
        } catch (error) {
            handleInlineError(error, "Unable to Load Flow History");
            flowExecutions = [];
        } finally {
            historyLoading = false;
        }
    };

    const goToHistoryPage = (pageNum: number) => {
        if (pageNum < 1 || pageNum > historyPageCount) return;
        historyCurrentPage = pageNum;
        loadFlowHistory();
    };

    const handleHistoryPageChange = (event: CustomEvent<{ page: number }>) => {
        goToHistoryPage(event.detail.page);
    };

    const formatDuration = (startedAt: string, completedAt?: string) => {
        if (!startedAt) return "Unknown";

        const start = new Date(startedAt);
        const end = completedAt ? new Date(completedAt) : new Date();
        const durationMs = end.getTime() - start.getTime();

        if (durationMs < 1000) return "<1s";

        const seconds = Math.floor(durationMs / 1000);
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);

        if (hours > 0) {
            return `${hours}h ${minutes % 60}m`;
        } else if (minutes > 0) {
            return `${minutes}m ${seconds % 60}s`;
        } else {
            return `${seconds}s`;
        }
    };

    // Watch for tab changes and load history when needed
    $effect(() => {
        if (activeTab === "history") {
            loadFlowHistory();
        }
    });

    // Table configuration
    const tableColumns: TableColumn<any>[] = [
        {
            key: "id",
            header: "Exec ID",
            render: (value) => `
        <a
          href="/view/${namespace}/results/${flowId}/${value}"
          class="exec-link"
        >
          ${value.substring(0, 8)}
        </a>
      `,
        },
        {
            key: "started_at",
            header: "Started At",
            width: "w-40",
            render: (_value, row) =>
                `<span class="text-lighter text-sm">${formatDateTime(getStartTime(row))}</span>`,
        },
        {
            key: "duration",
            header: "Duration",
            render: (_value, row) =>
                `<span class="text-lighter text-sm">${formatDuration(getStartTime(row), row.completed_at)}</span>`,
        },
        {
            key: "status",
            header: "Status",
            component: StatusBadge,
        },
        {
            key: "triggered_by",
            header: "Triggered By",
            width: "w-32",
            render: (value) =>
                `<span class="text-sm">${value || "System"}</span>`,
        },
        {
            key: "trigger_type",
            header: "Trigger Type",
            render: (value, row) =>
                `<span class="badge ${
                    row.trigger_type === "manual"
                        ? ""
                        : "success"
                }">${row.trigger_type}</span>`,
        },
    ];

    // Tab configuration
    const tabs = [
        { id: "run", label: "Run" },
        { id: "schedule", label: "Schedule" },
        { id: "history", label: "History" },
    ];
</script>

<svelte:head>
    <title>Run Flow - {data.flowMeta?.meta?.name || "Loading..."}</title>
</svelte:head>

<Header
    breadcrumbs={[
        { label: namespace!, url: `/view/${namespace}/flows` },
        { label: "Flows", url: `/view/${namespace}/flows` },
        ...(data.flowMeta?.meta?.prefix
            ? [{ label: data.flowMeta.meta.prefix, url: `/view/${namespace}/flows?group=${encodeURIComponent(data.flowMeta.meta.prefix)}` }]
            : []),
        { label: data.flowMeta?.meta?.name || "Loading..." },
    ]}
    actions={[
        ...(canUpdateFlow
            ? [
                  {
                      icon: IconPencil,
                      label: "Edit",
                      onClick: () =>
                          goto(`/view/${namespace}/flows/${flowId}/edit`),
                      variant: "primary" as const,
                  },
              ]
            : canViewFlowConfig
              ? [
                    {
                        icon: IconEye,
                        label: "View Config",
                        onClick: () =>
                            goto(`/view/${namespace}/flows/${flowId}/edit`),
                        variant: "secondary" as const,
                    },
                ]
              : []),
    ]}
/>

<FlowHero
    name={data.flowMeta?.meta?.name || "Loading..."}
    description={data.flowMeta?.meta?.description || ""}
/>

<div class="tab-bar">
    <div class="tab-bar-inner">
        <Tabs {tabs} bind:activeTab />
    </div>
</div>

<!-- Tab Content -->
<div class="tab-content">
    {#if activeTab === "run"}
        <div class="content-narrow">
            {#if showRerunBanner}
                <div class="mb-6">
                    <div role="alert" class="hstack gap-2 items-start justify-between">
                        <div class="hstack gap-2 flex-1 items-start">
                            <svg style="width: 1.25rem; height: 1.25rem; color: var(--primary); margin-top: 0.125rem; flex-shrink: 0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            <div class="flex-1">
                                <h3 class="text-sm font-medium">Rerunning execution</h3>
                                <p class="text-sm text-light mt-2">
                                    Inputs have been prepopulated from execution
                                    <a href="/view/{namespace}/results/{flowId}/{rerunFromExecId}" class="font-mono">
                                        {rerunFromExecId.substring(0, 8)}
                                    </a>
                                </p>
                            </div>
                        </div>
                        <button onclick={() => (showRerunBanner = false)} class="ghost icon small" aria-label="Dismiss">
                            <svg style="width: 1.25rem; height: 1.25rem" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>
            {/if}
            <FlowInputForm
                inputs={data.flowInputs || []}
                namespace={namespace!}
                flowId={flowId!}
                executionInput={data.executionInput}
                optionsRequestId={data.optionsRequestId}
                onScheduled={refreshScheduledExecutions}
            />
            <FlowActionsSummary actions={data.flowMeta?.actions || []} />
        </div>
    {/if}

    <!-- Schedules Tab -->
    {#if activeTab === "schedule"}
        <div class="content-wide">
            <div class="vstack gap-4">
                <ScheduledExecutionsList
                    schedules={scheduledExecutions}
                    cronSchedules={userSchedules}
                    namespace={namespace!}
                    flowId={flowId!}
                />

                <FlowSchedulesList
                    namespace={namespace!}
                    flowId={flowId!}
                    flowInputs={data.flowInputs || []}
                    userSchedulable={data.flowMeta?.meta?.user_schedulable ||
                        false}
                    user={data.user}
                    schedules={userSchedules}
                    onUpdate={reloadSchedules}
                    {canUpdateFlow}
                />

            </div>
        </div>
    {/if}

    <!-- History Tab -->
    {#if activeTab === "history"}
        <div class="content-max">
            <Table
                columns={tableColumns}
                data={flowExecutions}
                loading={historyLoading}
                title="Execution History for {data.flowMeta?.meta?.name ||
                    'Flow'}"
                subtitle="Past executions of this flow"
                emptyMessage="No execution history"
            />

            {#if historyPageCount > 1}
                <div class="mt-6 hstack justify-between items-center">
                    <div class="text-light text-sm">
                        Showing {(historyCurrentPage - 1) *
                            historyItemsPerPage +
                            1} to {Math.min(
                            historyCurrentPage * historyItemsPerPage,
                            historyTotalCount,
                        )} of {historyTotalCount} results
                    </div>
                    <Pagination
                        currentPage={historyCurrentPage}
                        totalPages={historyPageCount}
                        loading={historyLoading}
                        on:page-change={handleHistoryPageChange}
                    />
                </div>
            {/if}
        </div>
    {/if}
</div>

<style>
    .tab-bar {
        background: var(--card);
        border-bottom: 1px solid var(--border);
        padding: 0 1.5rem;
    }

    .tab-bar-inner {
        max-width: 56rem;
        margin: 0 auto;
    }

    .tab-content {
        padding: 2rem 1.5rem;
        background: var(--muted);
    }

    .content-narrow {
        max-width: 42rem;
        margin: 0 auto;
    }

    .content-wide {
        max-width: 64rem;
        margin: 0 auto;
    }

    .content-max {
        max-width: 72rem;
        margin: 0 auto;
    }

    :global(.exec-link) {
        font-size: 0.875rem;
        font-family: monospace;
        color: var(--primary);
        text-decoration: none;
        display: block;
    }

    :global(.exec-link:hover) {
        text-decoration: underline;
    }
</style>
