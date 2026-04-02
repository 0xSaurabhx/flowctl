<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { apiClient } from "$lib/apiClient.js";
    import Header from "$lib/components/shared/Header.svelte";
    import PageHeader from "$lib/components/shared/PageHeader.svelte";
    import FlowMetadata from "$lib/components/flow-create/FlowMetadata.svelte";
    import FlowInputs from "$lib/components/flow-create/FlowInputs.svelte";
    import FlowActions from "$lib/components/flow-create/FlowActions.svelte";
    import FlowNotifications from "$lib/components/flow-create/FlowNotifications.svelte";
    import ValidationModal from "$lib/components/flow-create/ValidationModal.svelte";
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

    let formElement: HTMLFormElement;

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

<Header
    breadcrumbs={[
        { label: namespace, url: `/view/${namespace}/flows` },
        { label: "Flows", url: `/view/${namespace}/flows` },
        { label: "Create" },
    ]}
/>

<div class="page-content">
        <PageHeader title="Create Flow" subtitle="Define a new workflow" />

                {#if data.prefillFlow}
                    <div role="alert" class="mb-4 hstack gap-2" style="align-items: flex-start">
                        <IconInfoCircle class="shrink-0" style="margin-top: var(--space-1);" size={16} />
                        <span>Secrets are not copied. You will need to re-add any secrets under the <strong>Secrets</strong> tab after creating this flow.</span>
                    </div>
                {/if}

                <form bind:this={formElement}>
                    <ot-tabs>
                        <div role="tablist">
                            <button role="tab" aria-selected="true">General</button>
                            <button role="tab">Inputs</button>
                            <button role="tab">Actions</button>
                            <button role="tab">Notifications</button>
                            <button role="tab">Secrets</button>
                        </div>
                        <div role="tabpanel">
                            <FlowMetadata
                                bind:metadata={flow.metadata}
                                {namespace}
                                inputs={flow.inputs}
                            />
                        </div>
                        <div role="tabpanel">
                            <FlowInputs bind:inputs={flow.inputs} {addInput} />
                        </div>
                        <div role="tabpanel">
                            <FlowActions
                                {namespace}
                                bind:actions={flow.actions}
                                {addAction}
                                {availableExecutors}
                                bind:executorConfigs
                            />
                        </div>
                        <div role="tabpanel">
                            <FlowNotifications
                                bind:notifications={flow.notifications}
                                {addNotification}
                                {availableMessengers}
                                {messengerConfigs}
                            />
                        </div>
                        <div role="tabpanel">
                            <SecretsTab {namespace} disabled={true} />
                        </div>
                    </ot-tabs>

                    <!-- Action Buttons -->
                    <div class="hstack gap-2 justify-end mt-6">
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
                </form>
</div>

{#if showValidation}
    <ValidationModal bind:show={showValidation} {validationResult} />
{/if}

<style>
    form {
        max-width: 64rem;
    }

    :global([role="tabpanel"]) {
        border: 1px solid var(--border);
        border-radius: var(--radius-medium);
        padding: var(--space-6);
        background: var(--card);
    }
</style>
