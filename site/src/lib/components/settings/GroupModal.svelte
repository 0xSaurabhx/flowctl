<script lang="ts">
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { autofocus } from "$lib/utils/autofocus";
    import type { GroupWithUsers, User } from "$lib/types";

    let {
        isEditMode = false,
        groupData = null,
        onSave,
        onClose,
    }: {
        isEditMode: boolean;
        groupData: GroupWithUsers | null;
        onSave: (data: any) => Promise<void>;
        onClose: () => void;
    } = $props();

    let name = $state(groupData?.name || "");
    let description = $state(groupData?.description || "");
    let users: User[] = $derived(groupData?.users || []);
    let saving = $state(false);
    let dialogEl: HTMLDialogElement;

    async function handleSubmit(event: Event) {
        event.preventDefault();

        saving = true;

        try {
            await onSave({
                name: name.trim(),
                description: description.trim(),
            });
        } catch (err) {
            handleInlineError(
                err,
                isEditMode
                    ? "Unable to Update Group"
                    : "Unable to Create Group",
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
            <h3>{isEditMode ? "Edit Group" : "Add New Group"}</h3>
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
                <label for="description">Description</label>
                <input
                    type="text"
                    id="description"
                    bind:value={description}
                    disabled={saving}
                    placeholder="Optional description"
                />
            </div>

            {#if isEditMode}
                <div>
                    <label>Members ({users.length})</label>
                    {#if users.length > 0}
                        <div class="members-list">
                            {#each users as user}
                                <div class="member-row hstack justify-between">
                                    <span>{user.name}</span>
                                    <span class="text-lighter">{user.username}</span>
                                </div>
                            {/each}
                        </div>
                    {:else}
                        <p class="text-lighter italic">No members in this group</p>
                    {/if}
                </div>
            {/if}
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
    .members-list {
        border: 1px solid var(--border);
        border-radius: 0.5rem;
        max-height: 10rem;
        overflow-y: auto;
    }
    .member-row {
        padding: 0.5rem 0.75rem;
        font-size: 0.875rem;
    }
    .members-list > :not(:last-child) {
        border-bottom: 1px solid var(--border);
    }
    .italic {
        font-style: italic;
        font-size: 0.875rem;
    }
</style>
