<script lang="ts">
	import { autofocus } from '$lib/utils/autofocus';
	import JsonDisplay from '$lib/components/shared/JsonDisplay.svelte';
	import type { UserSchedule } from '$lib/types';

	interface Props {
		schedule: UserSchedule;
		onClose: () => void;
	}

	let { schedule, onClose }: Props = $props();

	let dialogEl: HTMLDialogElement;

	$effect(() => {
		if (dialogEl) {
			dialogEl.showModal();
		}
	});
</script>

<dialog bind:this={dialogEl} onclose={onClose}>
	<header>
		<h3>Schedule Details</h3>
	</header>

	<section>
		<!-- Schedule Info Box -->
		<div class="info-box">
			<div class="info-grid">
				<div>
					<div class="info-label text-xs font-medium text-lighter">Cron Expression</div>
					<code>{schedule.cron}</code>
				</div>
				<div>
					<div class="info-label text-xs font-medium text-lighter">Timezone</div>
					<div>{schedule.timezone}</div>
				</div>
			</div>
		</div>

		<!-- Schedule Inputs -->
		<div>
			<JsonDisplay data={schedule.inputs} title="Inputs" expanded={true} />
		</div>
	</section>

	<footer>
		<button
			type="button"
			data-variant="secondary"
			onclick={onClose}
			use:autofocus
		>
			Close
		</button>
	</footer>
</dialog>

<style>
	dialog {
		max-width: 32rem;
		width: 100%;
	}
	.info-box {
		background: var(--faint);
		border: 1px solid var(--border);
		border-radius: 0.5rem;
		padding: 1rem;
		margin-bottom: 1rem;
	}
	.info-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.75rem;
	}
	.info-label {
		text-transform: uppercase;
		margin-bottom: 0.25rem;
	}
</style>
