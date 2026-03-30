<script lang="ts">
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { autofocus } from "$lib/utils/autofocus";
    import type { CredentialResp, NodeReq, NodeResp } from "$lib/types";

    interface Props {
        isEditMode?: boolean;
        nodeData?: NodeResp | null;
        credentials: CredentialResp[];
        onSave: (nodeData: NodeReq) => void;
        onClose: () => void;
    }

    let {
        isEditMode = false,
        nodeData = null,
        credentials,
        onSave,
        onClose,
    }: Props = $props();

    let formData = $state({
        name: "",
        hostname: "",
        port: 22,
        username: "",
        connection_type: "ssh",
        auth: {
            credential_id: "",
            method: "",
        },
        tags: [] as string[],
        tagsString: "",
    });

    let loading = $state(false);
    let dialogEl: HTMLDialogElement;

    $effect(() => {
        if (isEditMode && nodeData) {
            formData.name = nodeData.name || "";
            formData.hostname = nodeData.hostname || "";
            formData.port = nodeData.port || 22;
            formData.username = nodeData.username || "";
            formData.connection_type = nodeData.connection_type || "ssh";
            formData.auth.credential_id = nodeData.auth?.credential_id || "";
            formData.auth.method = nodeData.auth?.method || "";
            formData.tags = nodeData.tags || [];
            formData.tagsString = (nodeData.tags || []).join(", ");
        } else if (!isEditMode) {
            formData.name = "";
            formData.hostname = "";
            formData.port = 22;
            formData.username = "";
            formData.connection_type = "ssh";
            formData.auth.credential_id = "";
            formData.auth.method = "";
            formData.tags = [];
            formData.tagsString = "";
        }
    });

    $effect(() => {
        if (dialogEl) {
            dialogEl.showModal();
        }
    });

    function onCredentialChange() {
        if (formData.auth.credential_id) {
            const selectedCredential = credentials.find(
                (c) => c.id === formData.auth.credential_id,
            );
            if (selectedCredential) {
                formData.auth.method = selectedCredential.key_type;
            }
        } else {
            formData.auth.method = "";
        }
    }

    function getAuthMethodDisplay(method: string) {
        switch (method) {
            case "private_key":
                return "SSH Key";
            case "password":
                return "Password";
            default:
                return "";
        }
    }

    async function handleSubmit() {
        try {
            loading = true;

            const tags = formData.tagsString
                .split(",")
                .map((tag) => tag.trim())
                .filter((tag) => tag.length > 0);

            const nodeFormData: NodeReq = {
                name: formData.name,
                hostname: formData.hostname,
                port: formData.port,
                username: formData.username,
                connection_type: formData.connection_type,
                tags: tags,
                auth: {
                    credential_id: formData.auth.credential_id,
                    method: formData.auth.method,
                },
            };

            await onSave(nodeFormData);
        } catch (err) {
            handleInlineError(
                err,
                isEditMode ? "Unable to Update Node" : "Unable to Create Node",
            );
        } finally {
            loading = false;
        }
    }
</script>

<dialog bind:this={dialogEl} onclose={onClose}>
    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
        <header>
            <h3>{isEditMode ? "Edit Node" : "Add Node"}</h3>
        </header>

        <section>
            <div data-field>
                <label>Name</label>
                <input
                    type="text"
                    bind:value={formData.name}
                    required
                    disabled={loading}
                    use:autofocus
                />
            </div>

            <div data-field>
                <label>Hostname</label>
                <input
                    type="text"
                    bind:value={formData.hostname}
                    required
                    disabled={loading}
                />
            </div>

            <div data-field>
                <label>Port</label>
                <input
                    type="number"
                    bind:value={formData.port}
                    min="1"
                    max="65535"
                    required
                    disabled={loading}
                />
            </div>

            <div data-field>
                <label>Username</label>
                <input
                    type="text"
                    bind:value={formData.username}
                    required
                    disabled={loading}
                />
            </div>

            <div data-field>
                <label>Connection Type</label>
                <select
                    bind:value={formData.connection_type}
                    required
                    disabled={loading}
                >
                    <option value="">Select connection type</option>
                    <option value="ssh">SSH</option>
                    <option value="qssh">QSSH</option>
                </select>
            </div>

            <div data-field>
                <label>Credential</label>
                <select
                    bind:value={formData.auth.credential_id}
                    onchange={onCredentialChange}
                    required
                    disabled={loading}
                >
                    <option value="">Select credential</option>
                    {#each credentials as credential}
                        <option value={credential.id}>
                            {credential.name} ({credential.key_type})
                        </option>
                    {/each}
                </select>
            </div>

            <div data-field>
                <label>Tags (comma-separated)</label>
                <input
                    type="text"
                    bind:value={formData.tagsString}
                    placeholder="production, web, east"
                    disabled={loading}
                />
            </div>
        </section>

        <footer>
            <button
                type="button"
                data-variant="secondary"
                onclick={onClose}
                disabled={loading}
            >
                Cancel
            </button>
            <button
                type="submit"
                disabled={loading}
                aria-busy={loading}
            >
                {isEditMode ? "Update" : "Create"}
            </button>
        </footer>
    </form>
</dialog>

<style>
    dialog {
        max-width: 32rem;
        width: 100%;
    }
</style>
