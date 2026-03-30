<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { apiClient } from "$lib/apiClient.js";
    import Header from "$lib/components/shared/Header.svelte";
    import FlowMetadata from "$lib/components/flow-create/FlowMetadata.svelte";
    import FlowInputs from "$lib/components/flow-create/FlowInputs.svelte";
    import FlowActions from "$lib/components/flow-create/FlowActions.svelte";
    import FlowNotifications from "$lib/components/flow-create/FlowNotifications.svelte";
    import ValidationModal from "$lib/components/flow-create/ValidationModal.svelte";
    import Tabs from "$lib/components/shared/Tabs.svelte";
    import SecretsTab from "$lib/components/secrets/SecretsTab.svelte";
    import type { PageData } from "./$types";
    import type {
        FlowCreateReq,
        FlowInputReq,
        FlowActionReq,
        Schedule,
    } from "$lib/types.js";
    import { goto } from "$app/navigation";
    import { handleInlineError, showSuccess } from "$lib/utils/errorHandling";
    import { IconInfoCircle } from "@tabler/icons-svelte";

    let { data }: { data: PageData } = $props();
    const namespace = $page.params.namespace as string;

    // Flow state
    let flow = $state({
        metadata: {
            id: "",
            name: "",
            description: "",
            prefix: "",
            schedules: [] as Schedule[],
            namespace: namespace,
            allow_overlap: false,
            user_schedulable: false,
        },
        inputs: [] as any[],
        actions: [] as any[],
        notifications: [] as any[],
    });

    // Modal states
    let showValidation = $state(false);
    let validationResult = $state({
        success: false,
        errors: [] as string[],
    });

    // Loading states
    let saving = $state(false);
    const availableExecutors = data.availableExecutors;
    const availableMessengers = data.availableMessengers || [];

    // Executor configs for actions
    let executorConfigs = $state({} as Record<string, any>);

    // Messenger configs for notifications (pre-loaded in page loader)
    const messengerConfigs = data.messengerConfigs || {};

    // Tab state
    let activeTab = $state("metadata");
    let formElement: HTMLFormElement;

    const tabs = [
        { id: "metadata", label: "General" },
        { id: "inputs", label: "Inputs" },
        { id: "actions", label: "Actions" },
        { id: "notifications", label: "Notifications" },
        { id: "secrets", label: "Secrets" },
    ];

    async function loadExecutorConfigs(actions: any[]) {
        const executorTypes = [...new Set(actions.map((a) => a.executor).filter(Boolean))];
        for (const executor of executorTypes) {
            try {
                const config = await apiClient.executors.getConfig(executor);
                if (config.$defs && config.$ref) {
                    const refPath = config.$ref.replace("#/$defs/", "");
                    executorConfigs[executor] = config.$defs[refPath] || config;
                } else {
                    executorConfigs[executor] = config;
                }
            } catch (error) {
                handleInlineError(error, `Error loading config for executor ${executor}`);
            }
        }
    }

    onMount(async () => {
        if (data.prefillFlow) {
            flow.metadata = data.prefillFlow.metadata;
            flow.inputs = data.prefillFlow.inputs;
            flow.actions = data.prefillFlow.actions;
            flow.notifications = data.prefillFlow.notifications;
            if (flow.actions.length > 0) {
                await loadExecutorConfigs(flow.actions);
            }
        } else {
            addAction();
        }
    });

    function addInput() {
        flow.inputs.push({
            name: "",
            type: "string",
            label: "",
            description: "",
            required: false,
            default: "",
            validation: "",
            options: [],
            optionsText: "",
        });
    }

    function addAction() {
        const tempId = Date.now() + Math.random();
        flow.actions.push({
            tempId: tempId,
            id: "",
            name: "",
            executor: "",
            with: {},
            selectedNodes: [],
            variables: [],
            approval: false,
            artifacts: [],
            condition: "",
            collapsed: false,
        });
    }

    function addNotification() {
        flow.notifications.push({
            channel: "",
            events: [],
            config: {},
        });
    }

    async function saveFlow() {
        saving = true;

        try {
            // Transform the flow data to match the API schema
            const flowData: FlowCreateReq = {
                metadata: {
                    name: flow.metadata.name,
                    description: flow.metadata.description || undefined,
                    prefix: flow.metadata.prefix || undefined,
                    schedules:
                        flow.metadata.schedules?.filter((s) => s.cron.trim()) ||
                        undefined,
                    allow_overlap: flow.metadata.allow_overlap || false,
                    user_schedulable: flow.metadata.user_schedulable || false,
                },
                inputs: flow.inputs
                    .filter((i) => i.name)
                    .map(
                        (input): FlowInputReq => ({
                            name: input.name,
                            type: input.type,
                            label: input.label || undefined,
                            description: input.description || undefined,
                            validation: input.validation || undefined,
                            required: input.required || false,
                            default: input.default || undefined,
                            options:
                                input.type === "select" && !input.useRemoteOptions && input.optionsText
                                    ? input.optionsText
                                          .split("\n")
                                          .filter((o: string) => o.trim())
                                    : undefined,
                            remote_options:
                                input.type === "select" && input.useRemoteOptions && input.remote_options?.url
                                    ? {
                                          url: input.remote_options.url,
                                          method: input.remote_options.method || undefined,
                                          headers: Object.keys(input.remote_options.headers ?? {}).length > 0
                                              ? input.remote_options.headers
                                              : undefined,
                                      }
                                    : undefined,
                            max_file_size: input.max_file_size || undefined,
                        }),
                    ),
                actions: flow.actions
                    .filter((a) => a.name)
                    .map(
                        (action): FlowActionReq => ({
                            name: action.name,
                            executor: action.executor as "script" | "docker",
                            with: action.with || {},
                            approval: action.approval || false,
                            variables: action.variables
                                ?.filter((v: any) => v.name && v.name.trim())
                                .map((v: any) => ({ [v.name]: v.value })),
                            artifacts:
                                action.artifacts && action.artifacts.length > 0
                                    ? action.artifacts.filter((a: string) =>
                                          a.trim(),
                                      )
                                    : undefined,
                            condition: action.condition || undefined,
                            on: action.selectedNodes?.length
                                ? action.selectedNodes
                                : undefined,
                        }),
                    ),
                notify: flow.notifications
                    .filter((n) => n.channel)
                    .map((notification) => ({
                        channel: notification.channel,
                        events: notification.events || [],
                        config: notification.config || {},
                    })),
            };

            const result = await apiClient.flows.create(namespace, flowData);

            showSuccess(
                "Flow Created",
                `Flow "${flow.metadata.name}" has been created successfully.`,
            );

            // Redirect to flow detail page
            await goto(`/view/${namespace}/flows/${result.id}`);
        } catch (error: any) {
            handleInlineError(error, "Unable to Create Flow");
        } finally {
            saving = false;
        }
    }

    // Remove header actions - buttons moved to bottom
</script>

<svelte:head>
    <title>Create Flow - {namespace} | Flowctl</title>
</svelte:head>

<div class="create-layout">
    <!-- Main Content -->
    <div class="create-main">
        <Header
            breadcrumbs={[
                { label: namespace, url: `/view/${namespace}/flows` },
                { label: "Flows", url: `/view/${namespace}/flows` },
                { label: "Create" },
            ]}
        />

        <!-- Page Content -->
        <div class="create-content">
            <div class="create-container">
                <!-- Page Title -->
                <div class="mb-6">
                    <h1>Create Flow</h1>
                    <p class="text-lighter mt-2">
                        Define a new workflow
                    </p>
                </div>

                {#if data.prefillFlow}
                    <div class="info-banner mb-4 hstack gap-2">
                        <IconInfoCircle class="shrink-0" style="margin-top: var(--space-1);" size={16} />
                        <span>Secrets are not copied. You will need to re-add any secrets under the <strong>Secrets</strong> tab after creating this flow.</span>
                    </div>
                {/if}

                <!-- Main Card -->
                <div class="card">
                    <!-- Tab Navigation -->
                    <div class="card-tabs">
                        <Tabs bind:activeTab {tabs} />
                    </div>

                    <!-- Tab Content -->
                    <form bind:this={formElement} class="p-4">
                        {#if activeTab === "metadata"}
                            <FlowMetadata
                                bind:metadata={flow.metadata}
                                {namespace}
                                inputs={flow.inputs}
                            />
                        {:else if activeTab === "inputs"}
                            <FlowInputs bind:inputs={flow.inputs} {addInput} />
                        {:else if activeTab === "actions"}
                            <FlowActions
                                {namespace}
                                bind:actions={flow.actions}
                                {addAction}
                                {availableExecutors}
                                bind:executorConfigs
                            />
                        {:else if activeTab === "notifications"}
                            <FlowNotifications
                                bind:notifications={flow.notifications}
                                {addNotification}
                                {availableMessengers}
                                {messengerConfigs}
                            />
                        {:else if activeTab === "secrets"}
                            <SecretsTab {namespace} disabled={true} />
                        {/if}
                    </form>

                    <!-- Action Buttons -->
                    <div class="card-actions hstack gap-2">
                        <button
                            type="button"
                            onclick={() => goto(`/view/${namespace}/flows`)}
                            data-variant="secondary"
                        >
                            Cancel
                        </button>
                        <button
                            type="button"
                            onclick={() => {
                                if (formElement?.reportValidity()) {
                                    saveFlow();
                                }
                            }}
                            disabled={saving}
                            aria-busy={saving}
                        >
                            {saving ? "Creating..." : "Create"}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

{#if showValidation}
    <ValidationModal bind:show={showValidation} {validationResult} />
{/if}

<style>
    .create-layout {
        display: flex;
        height: 100vh;
        background: var(--muted);
    }

    .create-main {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    .create-content {
        flex: 1;
        overflow-y: auto;
        background: var(--muted);
    }

    .create-container {
        max-width: 72rem;
        margin: 0 auto;
        padding: 2rem 1.5rem;
    }

    .card-tabs {
        border-bottom: 1px solid var(--border);
    }

    .card-actions {
        padding: 1rem 1.5rem;
        background: var(--muted);
        border-top: 1px solid var(--border);
        justify-content: flex-end;
    }

    .info-banner {
        font-size: 0.875rem;
        border: 1px solid var(--border);
        background: var(--faint);
        border-radius: 0.5rem;
        padding: 0.75rem 1rem;
        color: var(--foreground);
        align-items: flex-start;
    }
</style>
