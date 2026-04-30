<script lang="ts">
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { autofocus } from "$lib/utils/autofocus";
    import OatSelect from "$lib/components/shared/OatSelect.svelte";
    import type { CredentialReq, CredentialResp } from "$lib/types";

    interface Props {
        isEditMode?: boolean;
        credentialData?: CredentialResp | null;
        onSave: (credentialData: CredentialReq) => void;
        onClose: () => void;
    }

    let {
        isEditMode = false,
        credentialData = null,
        onSave,
        onClose,
    }: Props = $props();

    let formData = $state({
        name: "",
        key_type: "" as "private_key" | "password" | "",
        key_data: "",
    });

    let loading = $state(false);

    let dialogEl: HTMLDialogElement;

    $effect(() => {
        if (isEditMode && credentialData) {
            formData = {
                name: credentialData.name || "",
                key_type: credentialData.key_type as "private_key" | "password",
                key_data: "",
            };
        } else if (!isEditMode) {
            formData = {
                name: "",
                key_type: "",
                key_data: "",
            };
        }
    });

    $effect(() => {
        if (dialogEl) {
            dialogEl.showModal();
        }
    });

    async function handleSubmit() {
        try {
            loading = true;

            const credentialFormData: CredentialReq = {
                name: formData.name,
                key_type: formData.key_type as "private_key" | "password",
                key_data: formData.key_data,
            };

            await onSave(credentialFormData);
        } catch (err) {
            handleInlineError(
                err,
                isEditMode
                    ? "Unable to Update Credential"
                    : "Unable to Create Credential",
            );
        } finally {
            loading = false;
        }
    }
</script>

<dialog bind:this={dialogEl} onclose={onClose}>
    <form
        onsubmit={(e) => {
            e.preventDefault();
            handleSubmit();
        }}
    >
        <header>
            <h3>{isEditMode ? "Edit Credential" : "Add New Credential"}</h3>
        </header>

        <section>
            <div class="row">
                <div class="col-6" data-field>
                    <label>Credential Name *</label>
                    <input
                        type="text"
                        bind:value={formData.name}
                        use:autofocus
                        placeholder="my-ssh-key"
                        required
                        disabled={loading}
                    />
                </div>

                <div class="col-6" data-field>
                    <label>Type *</label>
                    <OatSelect
                        bind:value={formData.key_type}
                        options={[
                            { value: 'private_key', label: 'SSH Key' },
                            { value: 'password', label: 'Password' }
                        ]}
                        placeholder="Select type..."
                        required
                        disabled={loading}
                    />
                </div>
            </div>

            {#if formData.key_type === "private_key"}
                <div data-field>
                    <label>Private Key *</label>
                    <textarea
                        bind:value={formData.key_data}
                        placeholder="-----BEGIN OPENSSH PRIVATE KEY-----"
                        required
                        disabled={loading}
                        rows="5"
                        class="font-mono text-xs"
                    ></textarea>
                </div>
            {/if}

            {#if formData.key_type === "password"}
                <div data-field>
                    <label>Password *</label>
                    <input
                        type="password"
                        bind:value={formData.key_data}
                        placeholder="Enter password"
                        required
                        disabled={loading}
                    />
                </div>
            {/if}
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
            >
                <span class="hstack gap-2 justify-center" aria-busy={loading} data-spinner="small">
                    {isEditMode ? "Update" : "Create"}
                </span>
            </button>
        </footer>
    </form>
</dialog>

<style>
    dialog {
        max-width: 42rem;
        width: 100%;
    }
</style>
