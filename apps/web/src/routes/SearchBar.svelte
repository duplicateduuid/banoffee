<script lang="ts">
	import classNames from 'classnames';
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import * as z from 'zod';
	import { goto } from '$app/navigation';

	type Props = {
		className?: string;
	};
	const { className }: Props = $props();

	const schema = z.object({ search: z.string().min(1) });

	const { form } = createForm<z.infer<typeof schema>>({
		extend: validator({ schema }),
		onSubmit: (values) => goto(`/search?url=${values.search}`)
	});
</script>

<form class={classNames('w-full', className)} use:form>
	<div class="relative w-full">
		<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
			<svg
				class="w-4 h-4 text-gray-500 dark:text-gray-400"
				aria-hidden="true"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 20 20"
			>
				<path
					stroke="currentColor"
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
				/>
			</svg>
		</div>
		<input
			type="search"
			name="search"
			id="default-search"
			class="block w-full p-4 ps-10 text-sm font-primary text-gray-900 border border-gray-300 rounded-lg bg-gray-50 outline-none"
			placeholder="Insert the URL here"
			required
		/>
		<button
			type="submit"
			class="text-gray-900 font-primary absolute end-2.5 bottom-2.5 bg-primary-400 hover:bg-primary-300 font-medium rounded-lg text-sm px-4 py-2 outline-none"
			>Search</button
		>
	</div>
</form>
