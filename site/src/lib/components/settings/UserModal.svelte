<script lang="ts">
    import GroupSelector from "$lib/components/shared/GroupSelector.svelte";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { autofocus } from "$lib/utils/autofocus";
    import type { User, Group, UserWithGroups } from "$lib/types";

    let {
        isEditMode = false,
        userData = null,
        availableGroups = [],
        onSave,
        onClose,
    }: {
        isEditMode: boolean;
        userData: UserWithGroups | null;
        availableGroups: Group[];
        onSave: (data: any) => Promise<void>;
        onClose: () => void;
    } = $props();

    let name = $state(userData?.name || "");
    let username = $state(userData?.username || "");
    let selectedGroups = $state<Group[]>(userData?.groups || []);
    let saving = $state(false);
    let dialogEl: HTMLDialogElement;

    async function handleSubmit(event: Event) {
        event.preventDefault();

        saving = true;

        try {
            await onSave({
                name: name.trim(),
                username: username.trim(),
                groups: selectedGroups.map((g) => g.id),
            });
        } catch (err) {
            handleInlineError(
                err,
                isEditMode ? "Unable to Update User" : "Unable to Create User",
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
            <h3>{isEditMode ? "Edit User" : "Add New User"}</h3>
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
                    use:autofocus
                />
            </div>

            <div data-field>
                <label for="username">Email</label>
                <input
                    type="text"
                    id="username"
                    bind:value={username}
                    required
                    disabled={saving}
                />
            </div>

            <div data-field>
                <label>Groups</label>
                <GroupSelector
                    bind:selectedGroups
                    disabled={saving}
                    placeholder="Search and select groups..."
                />
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
                disabled={saving}
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
