<script lang="ts">
	import classNames from 'classnames';
	import { createDialog, melt } from '@melt-ui/svelte';
	import { api } from '../../api';
	import type { Resource } from '../../schemas/resource';
	import { fade } from 'svelte/transition';
	import { BookmarkIcon, CircleCheck, Clock, Star, X } from 'lucide-svelte';
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';

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

	const queryClient = useQueryClient();

	// LMAO
	const convertReviewRatingToNumber = (reviewRating?: string | null): number => {
		if (reviewRating === 'one') return 1;
		if (reviewRating === 'two') return 2;
		if (reviewRating === 'three') return 3;
		if (reviewRating === 'four') return 4;
		if (reviewRating === 'five') return 5;

		return 0;
	};

	const convertNumberToReviewRating = (rating: number): string | null => {
		if (rating === 1) return 'one';
		if (rating === 2) return 'two';
		if (rating === 3) return 'three';
		if (rating === 4) return 'four';
		if (rating === 5) return 'five';

		return null;
	};

	let resourceStatus = $state('bookmarked');
	let resourceRating = $state(0);

	const resource = createQuery<{ resource: Resource | null; user_holds?: boolean }>({
		queryKey: ['resource', url],
		queryFn: async () => {
			const { data } = await api.get(`/user/resource?url=${url}`);

			if (data?.resource?.status) {
				resourceStatus = data.resource.status;
			}

			if (data?.resource?.review_rating) {
				resourceRating = convertReviewRatingToNumber(data.resource.review_rating);
			}

			return data;
		},
		retry: false,
		enabled: $open,
		staleTime: 0
	});

	const saveMutation = createMutation({
		mutationKey: ['save', url],
		mutationFn: async (input: { resource: Resource | null; user_holds?: boolean }) => {
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
		},
		onSuccess: () => $resource.refetch()
	});

	const handleSaveResource = (input: { resource: Resource | null; user_holds?: boolean }) =>
		$saveMutation.mutate(input);

	const updateMutation = createMutation({
		mutationKey: ['update', url],
		mutationFn: async (input: { resource: Resource; status: string; rating: number }) => {
			try {
				await api.post(`/user/resource/${input.resource.id}`, {
					status: input.status,
					review_rating: convertNumberToReviewRating(input.rating)
				});
			} catch (err) {
				throw new Error(`Unexpected error calling the API: ${JSON.stringify(err)}`);
			}
		},
		onSuccess: () => $resource.refetch()
	});

	const handleUpdateResource = (input: { resource: Resource; status: string; rating: number }) => {
		$updateMutation.mutate(input);

		queryClient.refetchQueries({ queryKey: ['get-my-resources'] });

		// TODO: for some weird reason refetch does not trigger rerender (lol)
		window.location.reload();
	};
</script>

<div class={classNames('flex gap-2 h-40 shadow px-8 py-6 rounded-lg bg-white', className)}>
	<span class="bg-primary-400 h-full aspect-video rounded-lg" />
	<div class="w-full flex flex-col justify-between h-full gap-4">
		<div class="w-full flex flex-col">
			<p class="font-primary text-2xl tracking-wider font-semibold line-clamp-1">
				{name}
			</p>
			{#if description}
				<p class="text-sm line-clamp-2">{description}</p>
			{/if}
		</div>
		<div class={classNames('w-full flex justify-end gap-2', author && 'justify-between')}>
			{#if author}
				<p class="line-clamp-1">{author}</p>
			{/if}

			<button
				use:melt={$trigger}
				on:click={(e) => {
					e.stopPropagation();
					$resource.refetch();
				}}
			>
				{#if $resource.data?.user_holds}
					Update
				{:else}
					Save
				{/if}
			</button>
		</div>
	</div>
</div>

<!-- TODO: beautify loading state -->
<div class="" use:melt={$portalled}>
	{#if $open}
		<div
			use:melt={$overlay}
			class="fixed inset-0 z-50 bg-black/50"
			transition:fade={{ duration: 150 }}
		/>
		<div
			class="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-50 min-w-[90%] sm:min-w-[30em] max-w-[90%] lg:max-w-xl rounded-xl flex flex-col justify-between gap-6 p-8 font-secondary bg-white shadow-xl"
			use:melt={$content}
		>
			<button
				use:melt={$close}
				class="p-2 border bg-white hover:bg-stone-50 shadow rounded-lg absolute -top-4 -right-4 transition"
			>
				<X size={16} />
			</button>

			<div class="max-w-full flex flex-col gap-[2px]">
				<p class="text-sm font-normal">Title:</p>
				<p class="text-base tracking-wider line-clamp-1 font-medium">
					{name || 'Unknown title'}
				</p>
			</div>

			{#if !$resource.isPending && !$resource.data}
				<p>Loading...</p>
			{:else if $resource.data && $resource.data.user_holds}
				<div class="flex flex-col sm:flex-row sm:items-center gap-4">
					<button
						class="flex items-center gap-2 cursor-pointer"
						on:click={() => {
							resourceStatus = 'bookmarked';
						}}
					>
						<BookmarkIcon
							size={20}
							class={classNames('stroke-gray-600', {
								'stroke-blue-600': resourceStatus === 'bookmarked'
							})}
						/>
						<span
							class={classNames('text-sm font-medium text-gray-600', {
								'text-blue-600': resourceStatus === 'bookmarked'
							})}>Bookmarked</span
						>
					</button>

					<button
						class="flex items-center gap-2 cursor-pointer"
						on:click={() => {
							resourceStatus = 'ongoing';
						}}
					>
						<Clock
							size={20}
							class={classNames('stroke-gray-600', {
								'stroke-yellow-600': resourceStatus === 'ongoing'
							})}
						/>
						<span
							class={classNames('text-sm font-medium text-gray-600', {
								'text-yellow-600': resourceStatus === 'ongoing'
							})}>On-going</span
						>
					</button>

					<button
						class="flex items-center gap-2 cursor-pointer"
						on:click={() => {
							resourceStatus = 'completed';
						}}
					>
						<CircleCheck
							size={20}
							class={classNames('stroke-gray-600', {
								'stroke-green-600': resourceStatus === 'completed'
							})}
						/>
						<span
							class={classNames('text-sm font-medium text-gray-600', {
								'text-green-600': resourceStatus === 'completed'
							})}>Completed</span
						>
					</button>
				</div>
				<div class="flex items-center gap-2">
					<p class="text-sm font-medium text-gray-600">Rating:</p>
					{#each { length: 5 } as _, i}
						<button
							class="cursor-pointer"
							on:click={() => {
								resourceRating = i + 1;
							}}
						>
							<Star
								size={20}
								class={classNames({
									'fill-yellow-400': resourceRating !== 0 && i + 1 <= resourceRating
								})}
							/>
						</button>
					{/each}
				</div>
				<button
					class="w-full h-12 rounded-md border-none bg-[#4e473b] text-[#F7F6F1] text-base inline-block disabled:opacity-50 disabled:cursor-not-allowed"
					on:click={() => {
						if ($resource.data?.resource) {
							handleUpdateResource({
								resource: $resource.data.resource,
								status: resourceStatus,
								rating: resourceRating
							});
						}
					}}
					disabled={$updateMutation.isPending}
				>
					{#if $updateMutation.isPending}
						Updating...
					{:else}
						Update
					{/if}
				</button>
			{:else}
				<button
					class="w-full h-12 rounded-md border-none bg-[#4e473b] text-[#F7F6F1] text-base inline-block disabled:opacity-50 disabled:cursor-not-allowed"
					on:click={() => {
						if ($resource.data) {
							handleSaveResource($resource.data);
						}
					}}
					disabled={($resource.isPending && !$resource.data) ||
						(!$resource.isPending && $resource.data && $resource.data.user_holds) ||
						$saveMutation.isPending}
				>
					{#if $saveMutation.isPending}
						Saving...
					{:else if $resource.isPending && !$resource.data}
						Loading...
					{:else}
						Bookmark
					{/if}
				</button>
			{/if}
		</div>
	{/if}
</div>
