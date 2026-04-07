<script lang="ts">
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { autofocus } from "$lib/utils/autofocus";
    import type { NamespaceResp } from "$lib/types";

    let {
        isEditMode = false,
        namespaceData = null,
        onSave,
        onClose,
    }: {
        isEditMode: boolean;
        namespaceData: NamespaceResp | null;
        onSave: (data: any) => Promise<void>;
        onClose: () => void;
    } = $props();

    let name = $state(namespaceData?.name || "");
    let saving = $state(false);
    let dialogEl: HTMLDialogElement;

    async function handleSubmit(event: Event) {
        event.preventDefault();

        saving = true;

        try {
            await onSave({
                name: name.trim(),
            });
        } catch (err) {
            handleInlineError(
                err,
                isEditMode
                    ? "Unable to Update Namespace"
                    : "Unable to Create Namespace",
            );
        } finally {
            saving = false;
        }
    }

    function handleClose() {
        if (!saving) {
            onClose();
        }
    }

    $effect(() => {
        if (dialogEl) {
            dialogEl.showModal();
        }
    });
</script>

<dialog bind:this={dialogEl} onclose={handleClose}>
    <form onsubmit={handleSubmit}>
        <header>
            <h3>{isEditMode ? "Edit Namespace" : "Add New Namespace"}</h3>
        </header>

        <section>
            <div data-field>
                <label for="name">Name</label>
                <input
                    type="text"
                    id="name"
                    bind:value={name}
                    required
                    disabled={saving}
                    maxlength="150"
                    placeholder="Enter namespace name"
                    use:autofocus
                />
                <p class="text-lighter text-xs mt-1">
                    Namespace names should be unique and descriptive.
                </p>
            </div>
        </section>

        <footer>
            <button
                type="button"
                data-variant="secondary"
                onclick={handleClose}
                disabled={saving}
            >
                Cancel
            </button>
            <button
                type="submit"
                disabled={saving || !name.trim()}
                aria-busy={saving}
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
