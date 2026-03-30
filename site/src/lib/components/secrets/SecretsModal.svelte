<script lang="ts">
	import { handleInlineError } from '$lib/utils/errorHandling';
	import { autofocus } from '$lib/utils/autofocus';
	import type { FlowSecretReq, FlowSecretUpdateReq, FlowSecretResp } from '$lib/types';

	interface Props {
		isEditMode?: boolean;
		secretData?: FlowSecretResp | null;
		onSave: (secretData: FlowSecretReq | FlowSecretUpdateReq) => void;
		onClose: () => void;
	}

	let {
		isEditMode = false,
		secretData = null,
		onSave,
		onClose
	}: Props = $props();

	let formData = $state({
		key: '',
		value: '',
		description: ''
	});

	let loading = $state(false);
	let dialogEl: HTMLDialogElement;

	$effect(() => {
		if (isEditMode && secretData) {
			formData = {
				key: secretData.key || '',
				value: '',
				description: secretData.description || ''
			};
		} else if (!isEditMode) {
			formData = {
				key: '',
				value: '',
				description: ''
			};
		}
	});

	$effect(() => {
		if (dialogEl) {
			dialogEl.showModal();
		}
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();
		try {
			loading = true;

			const secretFormData: FlowSecretReq | FlowSecretUpdateReq = isEditMode
				? {
					value: formData.value,
					description: formData.description || undefined
				}
				: {
					key: formData.key,
					value: formData.value,
					description: formData.description || undefined
				};

			await onSave(secretFormData);
		} catch (error) {
			handleInlineError(error);
		} finally {
			loading = false;
		}
	}
</script>

<dialog bind:this={dialogEl} onclose={onClose}>
	<form onsubmit={handleSubmit}>
		<header>
			<h3>{isEditMode ? 'Edit Secret' : 'Add New Secret'}</h3>
		</header>

		<section>
			<div data-field>
				<label for="key">
					Key <span class="required">*</span>
				</label>
				<input
					type="text"
					id="key"
					bind:value={formData.key}
					required
					disabled={loading || isEditMode}
					placeholder="SECRET_KEY"
					use:autofocus
				/>
				{#if isEditMode}
					<p class="text-lighter text-xs hint">Key cannot be changed. Delete and recreate the secret to use a different key.</p>
				{/if}
			</div>

			<div data-field>
				<label for="value">
					Value <span class="required">*</span>
				</label>
				<textarea
					id="value"
					bind:value={formData.value}
					required
					disabled={loading}
					class="font-mono text-xs"
					rows="5"
					placeholder={isEditMode ? 'Enter new value to update' : 'Enter value'}
				></textarea>
				{#if isEditMode}
					<p class="text-lighter text-xs hint">Enter a new value to update. Previous value is not shown for security.</p>
				{/if}
			</div>

			<div data-field>
				<label for="description">Description</label>
				<textarea
					id="description"
					bind:value={formData.description}
					disabled={loading}
					rows="3"
					placeholder="Optional description"
				></textarea>
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
				disabled={loading || !formData.key || !formData.value}
				aria-busy={loading}
			>
				{loading ? 'Saving...' : isEditMode ? 'Update' : 'Create'}
			</button>
		</footer>
	</form>
</dialog>

<style>
	dialog {
		max-width: 28rem;
		width: 100%;
	}
	.required {
		color: var(--danger);
	}
	.hint {
		margin-top: 0.25rem;
	}
</style>
