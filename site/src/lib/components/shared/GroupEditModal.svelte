<script lang="ts">
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { autofocus } from "$lib/utils/autofocus";

    let {
        groupName,
        description = "",
        onSave,
        onClose,
    }: {
        groupName: string;
        description: string;
        onSave: (data: { description: string }) => Promise<void>;
        onClose: () => void;
    } = $props();

    let editDescription = $state(description);
    let saving = $state(false);
    let dialogEl: HTMLDialogElement;

    async function handleSubmit(event: Event) {
        event.preventDefault();

        saving = true;

        try {
            await onSave({
                description: editDescription.trim(),
            });
        } catch (err) {
            handleInlineError(err, "Unable to Update Group");
        } finally {
            saving = false;
        }
    }

    function handleClose() {
        if (!saving) {
            dialogEl?.close();
            onClose();
        }
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Escape" && !saving) {
            handleClose();
        }
    }

    // Open dialog on mount
    import { onMount } from "svelte";
    onMount(() => {
        dialogEl?.showModal();
    });
</script>

<svelte:window on:keydown={handleKeydown} />

<dialog bind:this={dialogEl} onclose={handleClose}>
    <form onsubmit={handleSubmit}>
        <header>
            <h3>Edit Group</h3>
        </header>

        <section>
            <!-- Name Field (read-only) -->
            <div data-field>
                <label for="group-name">Name</label>
                <input
                    type="text"
                    id="group-name"
                    value={groupName}
                    disabled
                />
            </div>

            <!-- Description Field -->
            <div data-field>
                <label for="group-description">Description</label>
                <textarea
                    id="group-description"
                    bind:value={editDescription}
                    disabled={saving}
                    rows="3"
                    maxlength="500"
                    placeholder="Enter group description"
                    use:autofocus
                ></textarea>
            </div>
        </section>

        <footer>
            <button
                type="button"
                onclick={handleClose}
                disabled={saving}
                data-variant="secondary"
            >
                Cancel
            </button>
            <button
                type="submit"
                disabled={saving}
                aria-busy={saving}
            >
                Update
            </button>
        </footer>
    </form>
</dialog>
