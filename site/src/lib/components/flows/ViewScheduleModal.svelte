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
		<article class="card p-4 mb-4" style="background: var(--faint)">
			<div class="row">
				<div class="col-6">
					<div class="text-xs font-medium text-lighter" style="text-transform: uppercase; margin-bottom: 0.25rem">Cron Expression</div>
					<code>{schedule.cron}</code>
				</div>
				<div class="col-6">
					<div class="text-xs font-medium text-lighter" style="text-transform: uppercase; margin-bottom: 0.25rem">Timezone</div>
					<div>{schedule.timezone}</div>
				</div>
			</div>
		</article>

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
</style>
