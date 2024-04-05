<script lang="ts">
	import classNames from 'classnames';
	import { createDialog, melt } from '@melt-ui/svelte';
	import { api } from '../../api';
	import type { Resource } from '../../schemas/resource';
	import { fade } from 'svelte/transition';
	import { X } from 'lucide-svelte';
	import { createQuery } from '@tanstack/svelte-query';

	type Props = {
		redirect: string;
		url: string;
		name: string;
		description?: string | null;
		author?: string | null;
		className?: string;
	};
	const { url, name, description, author, className }: Props = $props();

	const {
		elements: { trigger, overlay, content, close, portalled },
		states: { open }
	} = createDialog({
		forceVisible: true
	});

	const resource = createQuery<{ resource: Resource | null; user_holds?: boolean }>({
		queryKey: ['resource', url],
		queryFn: async () => {
			const { data } = await api.get(`/user/resource?url=${url}`);

			return data;
		},
		retry: false,
		enabled: $open
	});

	// TODO: handle errors & request loading
	const handleSaveResource = async (input: { resource: Resource | null; user_holds?: boolean }) => {
		try {
			if (input.resource?.id) {
				await api.post(`/user/resource/${input.resource.id}`, { status: 'bookmarked' });
				return;
			}

			const {
				data: { resource: newResource }
			} = await api.post<{ resource: Resource }>('/resource', {
				url,
				name
			});

			await api.post(`/user/resource/${newResource.id}`, { status: 'bookmarked' });
		} catch (err) {
			throw new Error(`Unexpected error calling the API: ${JSON.stringify(err)}`);
		}
	};
</script>

<div class={classNames('flex gap-4 h-40 shadow px-8 py-6 rounded-lg bg-white', className)}>
	<span class="bg-primary-400 h-full aspect-video rounded-lg" />
	<div class="flex flex-col justify-between h-full">
		<div class="flex flex-col">
			<p class="font-primary text-2xl tracking-wider font-semibold line-clamp-1">
				{name}
			</p>
			{#if description}
				<p class="text-sm line-clamp-2">{description}</p>
			{/if}
		</div>
		<div class={classNames('flex justify-end', author && 'justify-between')}>
			{#if author}
				<p>{author}</p>
			{/if}

			<button
				use:melt={$trigger}
				on:click={(e) => {
					e.stopPropagation();
					$resource.refetch();
				}}>Save</button
			>
		</div>
	</div>
</div>

<div class="" use:melt={$portalled}>
	{#if $open}
		<div
			use:melt={$overlay}
			class="fixed inset-0 z-50 bg-black/50"
			transition:fade={{ duration: 150 }}
		/>
		<div
			class="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-50 min-w-[30em] rounded-xl flex flex-col justify-between gap-4 p-8 font-secondary bg-white shadow-xl"
			use:melt={$content}
		>
			<button
				use:melt={$close}
				class="p-2 border bg-white hover:bg-stone-50 shadow rounded-lg absolute -top-4 -right-4 transition"
			>
				<X size={16} />
			</button>

			<div
				class="max-w-full flex flex-col gap-[2px] text-ellipsis whitespace-nowrap overflow-hidden"
			>
				<p class="text-sm font-normal">Title:</p>
				<p class="text-base text-ellipsis whitespace-nowrap overflow-hidden font-medium">
					{name || 'Unknown title'}
				</p>
			</div>

			<button
				class="w-full h-12 rounded-md border-none bg-[#4e473b] text-[#F7F6F1] text-base inline-block disabled:opacity-50 disabled:cursor-not-allowed"
				on:click={() => {
					if ($resource.data) {
						handleSaveResource($resource.data).then(() => $resource.refetch());
					}
				}}
				disabled={($resource.isPending && !$resource.data) ||
					(!$resource.isPending && $resource.data && $resource.data.user_holds)}
			>
				{#if $resource.isPending && !$resource.data}
					Loading...
				{:else if !$resource.isPending && $resource.data && $resource.data.user_holds}
					Bookmarked
				{:else}
					Bookmark
				{/if}
			</button>
		</div>
	{/if}
</div>
