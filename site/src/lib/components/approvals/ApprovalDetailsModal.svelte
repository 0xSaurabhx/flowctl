<script lang="ts">
    import { apiClient } from "$lib/apiClient";
    import type { ApprovalDetailsResp } from "$lib/types";
    import JsonDisplay from "$lib/components/shared/JsonDisplay.svelte";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { autofocus } from "$lib/utils/autofocus";
    import { formatDateTime } from "$lib/utils";

    let {
        open = $bindable(),
        approvalId,
        namespace,
        onApprove,
        onReject,
    }: {
        open: boolean;
        approvalId: string;
        namespace: string;
        onApprove: (approvalId: string) => Promise<void>;
        onReject: (approvalId: string) => Promise<void>;
    } = $props();

    let approval: ApprovalDetailsResp | null = $state(null);
    let loading = $state(false);
    let error = $state<string | null>(null);
    let actionLoading = $state(false);
    let dialogEl: HTMLDialogElement;

    $effect(() => {
        if (open && approvalId) {
            fetchApprovalDetails();
        }
    });

    $effect(() => {
        if (open && dialogEl) {
            dialogEl.showModal();
        }
    });

    async function fetchApprovalDetails() {
        loading = true;
        error = null;
        try {
            approval = await apiClient.approvals.get(namespace, approvalId);
        } catch (err) {
            error = "Failed to load approval details";
            handleInlineError(err, "Unable to Load Approval Details");
        } finally {
            loading = false;
        }
    }

    function closeModal() {
        open = false;
        approval = null;
        error = null;
        dialogEl?.close();
    }

    async function handleApprove() {
        if (!approval) return;
        actionLoading = true;
        try {
            await onApprove(approval.id);
            await fetchApprovalDetails();
        } catch (err) {
            handleInlineError(err, "Unable to Approve Request");
        } finally {
            actionLoading = false;
        }
    }

    async function handleReject() {
        if (!approval) return;
        actionLoading = true;
        try {
            await onReject(approval.id);
            await fetchApprovalDetails();
        } catch (err) {
            handleInlineError(err, "Unable to Reject Request");
        } finally {
            actionLoading = false;
        }
    }
</script>

{#if open}
    <dialog bind:this={dialogEl} onclose={() => { open = false; approval = null; error = null; }}>
        <header>
            <div class="hstack gap-2">
                <div class="header-icon">
                    <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                </div>
                <h3>Approval Details</h3>
            </div>
        </header>

        <section>
            {#if loading}
                <div class="centered-msg" aria-busy="true">Loading...</div>
            {:else if error}
                <div class="centered-msg">
                    <div class="error-text">{error}</div>
                    <button data-variant="secondary" onclick={fetchApprovalDetails}>
                        Try again
                    </button>
                </div>
            {:else if approval}
                <div class="vstack gap-4">
                    <!-- Approval Overview -->
                    <div class="info-box">
                        <div class="info-grid">
                            <div>
                                <span class="info-label">Action ID</span>
                                <span>{approval.action_id}</span>
                            </div>
                            <div>
                                <span class="info-label">Status</span>
                                <span class="capitalize">{approval.status}</span>
                            </div>
                            <div>
                                <span class="info-label">Execution ID</span>
                                <a href="/view/{namespace}/results/{approval.flow_id}/{approval.exec_id}">
                                    <span class="mono-link">{approval.exec_id}</span>
                                </a>
                            </div>
                            <div>
                                <span class="info-label">Requested By</span>
                                <span>{approval.requested_by}</span>
                            </div>
                            <div>
                                <span class="info-label">Reviewer</span>
                                <span>{approval.approved_by || "N/A"}</span>
                            </div>
                            <div>
                                <span class="info-label">Flow</span>
                                <span>{approval.flow_name}</span>
                            </div>
                            <div>
                                <span class="info-label">Created</span>
                                <span>{formatDateTime(approval.created_at, "Never")}</span>
                            </div>
                        </div>
                    </div>

                    <!-- Execution Inputs -->
                    {#if approval.inputs}
                        <div>
                            <h4>Execution Inputs</h4>
                            <JsonDisplay data={approval.inputs} />
                        </div>
                    {/if}
                </div>
            {/if}
        </section>

        {#if approval && approval.status === "pending"}
            <footer>
                <button
                    data-variant="secondary"
                    onclick={handleReject}
                    disabled={actionLoading}
                    aria-busy={actionLoading}
                >
                    Reject
                </button>
                <button
                    onclick={handleApprove}
                    disabled={actionLoading}
                    aria-busy={actionLoading}
                >
                    Approve
                </button>
            </footer>
        {:else}
            <footer>
                <button
                    data-variant="secondary"
                    onclick={closeModal}
                    use:autofocus
                >
                    Close
                </button>
            </footer>
        {/if}
    </dialog>
{/if}

<style>
    dialog {
        max-width: 56rem;
        width: 100%;
    }
    .header-icon {
        width: 2rem;
        height: 2rem;
        background: var(--faint);
        border-radius: 0.5rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .header-icon .icon {
        width: 1.25rem;
        height: 1.25rem;
        color: var(--primary);
    }
    .centered-msg {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
        padding: 2rem;
    }
    .error-text {
        color: var(--danger);
    }
    .info-box {
        background: var(--faint);
        border-radius: 0.5rem;
        padding: 1rem;
    }
    .info-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }
    @media (max-width: 768px) {
        .info-grid {
            grid-template-columns: 1fr;
        }
    }
    .info-label {
        display: block;
        font-size: 0.875rem;
        font-weight: 500;
        color: var(--foreground);
        margin-bottom: 0.25rem;
    }
    .info-grid span:not(.info-label) {
        font-size: 0.875rem;
        color: var(--foreground);
    }
    .capitalize {
        text-transform: capitalize;
    }
    .mono-link {
        font-family: monospace;
        font-size: 0.875rem;
    }
    .mono-link:hover {
        text-decoration: underline;
    }
    h4 {
        font-size: 1rem;
        font-weight: 600;
        color: var(--foreground);
        margin-bottom: 0.75rem;
    }
</style>
