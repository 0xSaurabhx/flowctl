<script lang="ts">
	import type { FlowInput } from '$lib/types';

	let {
		inputs = [],
		values = $bindable({}),
		errors = {},
		useFormData = false
	}: {
		inputs?: FlowInput[];
		values?: Record<string, any>;
		errors?: Record<string, string>;
		useFormData?: boolean;
	} = $props();
</script>

{#if inputs && inputs.length > 0}
	<div class="vstack gap-4">
	{#each inputs as input (input.name)}
		<div data-field>
			<label for={input.name}>
				{input.label || input.name}
				{#if input.required}
					<span class="required-mark">*</span>
				{/if}
			</label>

			{#if errors[input.name]}
				<p class="error-text" data-hint>{errors[input.name]}</p>
			{/if}

			{#if input.type === 'string' || input.type === 'number'}
				{#if useFormData}
					<input
						type={input.type === 'string' ? 'text' : 'number'}
						id={input.name}
						name={input.name}
						value={values[input.name] ?? input.default ?? ''}
						placeholder={input.description || ''}
						required={input.required}
					/>
				{:else}
					<input
						type={input.type === 'string' ? 'text' : 'number'}
						bind:value={values[input.name]}
						placeholder={input.description || ''}
						required={input.required}
					/>
				{/if}
			{:else if input.type === 'checkbox'}
				<div>
					{#if useFormData}
						<input
							type="checkbox"
							id={input.name}
							name={input.name}
							value="true"
							checked={values[input.name] ?? input.default === 'true'}
						/>
					{:else}
						<input
							type="checkbox"
							bind:checked={values[input.name]}
						/>
					{/if}
				</div>
			{:else if input.type === 'select' && input.options}
				{#if useFormData}
					<select
						id={input.name}
						name={input.name}
						required={input.required}
						value={values[input.name] ?? input.default ?? ''}
					>
						<option value="">Select an option</option>
						{#each input.options as option}
							<option value={option} selected={option === (values[input.name] ?? input.default)}
								>{option}</option
							>
						{/each}
					</select>
				{:else}
					<select
						bind:value={values[input.name]}
						required={input.required}
					>
						<option value="">Select an option</option>
						{#each input.options as option}
							<option value={option}>{option}</option>
						{/each}
					</select>
				{/if}
			{:else if input.type === 'file'}
				<div class="vstack">
					<input
						type="file"
						id={input.name}
						name={input.name}
						required={input.required}
					/>
					{#if input.max_file_size}
						<small class="text-lighter" data-hint>
							Max size: {Math.round(input.max_file_size / (1024 * 1024))}MB
						</small>
					{/if}
				</div>
			{:else if input.type === 'datetime'}
				<div>
					{#if useFormData}
						<input
							type="datetime-local"
							id={input.name}
							name={input.name}
							value={values[input.name] ?? input.default ?? ''}
							required={input.required}
						/>
					{:else}
						<input
							type="datetime-local"
							bind:value={values[input.name]}
							required={input.required}
						/>
					{/if}
				</div>
			{:else if input.type === 'password'}
				<div>
					{#if useFormData}
						<input
							type="password"
							id={input.name}
							name={input.name}
							value={values[input.name] ?? input.default ?? ''}
							required={input.required}
						/>
					{:else}
						<input
							type="password"
							bind:value={values[input.name]}
							required={input.required}
						/>
					{/if}
				</div>
			{:else}
				<!-- Fallback for other input types -->
				{#if useFormData}
					<input
						type="text"
						id={input.name}
						name={input.name}
						value={values[input.name] ?? input.default ?? ''}
						placeholder={input.description || ''}
						required={input.required}
					/>
				{:else}
					<input
						type="text"
						bind:value={values[input.name]}
						placeholder={input.description || ''}
						required={input.required}
					/>
				{/if}
			{/if}

			{#if input.description}
				<small class="text-lighter" data-hint>{input.description}</small>
			{/if}
		</div>
	{/each}
	</div>
{/if}

<style>
	.required-mark {
		color: var(--danger);
	}

	.error-text {
		font-size: 0.875rem;
		color: var(--danger);
		margin-bottom: 0.5rem;
	}
</style>
