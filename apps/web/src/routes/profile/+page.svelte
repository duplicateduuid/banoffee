<script lang="ts">
	import { HeartHandshakeIcon, UserRound } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Card from '$lib/components/Card.svelte';
	import { getMyResources } from '../../requests/user';
	import { createInfiniteQuery } from '@tanstack/svelte-query';
	import type { Resource } from '../../schemas/resource';
	import classNames from 'classnames';
	import Spinner from '$lib/components/Spinner.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';

	let status = $state<'completed' | 'ongoing' | 'bookmarked'>('completed');

	const resourcesQuery = createInfiniteQuery<Resource[]>({
		queryKey: ['get-my-resources'],
		queryFn: ({ pageParam }) => getMyResources(10, pageParam as number, status),
		initialPageParam: 0,
		getNextPageParam: (lastPage, allPages) => {
			if (!lastPage) return undefined;

			return lastPage.length === 10 ? allPages.length * 10 : undefined;
		},
		staleTime: 0,
		refetchOnMount: false
	});

	let observerElem: HTMLDivElement | null = $state(null);

	$effect(() => {
		if (!$page.data.user) goto('/');

		const observer = new IntersectionObserver(
			(entries) => {
				const [target] = entries;

				if (target.isIntersecting && $resourcesQuery.hasNextPage) {
					$resourcesQuery.fetchNextPage();
				}
			},
			{ threshold: 0 }
		);

		if (observerElem) {
			observer.observe(observerElem);

			return () => {
				observer.unobserve(observerElem!);
			};
		}
	});
</script>

{#if $page.data.user}
	<section class="flex flex-col gap-4 items-center justify-center">
		<div class="user-info flex flex-col items-center gap-2.5 font-semibold text-xl w-full">
			<div class="avatar p-4 rounded-full bg-primary-400">
				<UserRound size={48} />
			</div>
			<p class="text-2xl font-primary font-bold">{$page.data.user.username}</p>
		</div>
		<div class="flex flex-row gap-2 items-center">
			<HeartHandshakeIcon size={16} />
			<p>
				Member since {new Intl.DateTimeFormat('en-US', {
					day: '2-digit',
					month: '2-digit'
				}).format(new Date($page.data.user.created_at))}
			</p>
		</div>
		<div class="flex flex-row gap-12 items-center mt-4">
			<button
				class="relative flex flex-col items-center outline-none border-none cursor-pointer"
				on:click={() => {
					status = 'completed';
					$resourcesQuery.refetch();
				}}
			>
				<p class="text-md font-primary font-semibold">Completed</p>
				{#if status === 'completed'}
					<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
				{/if}
			</button>
			<button
				class="relative flex flex-col items-center outline-none border-none cursor-pointer"
				on:click={() => {
					status = 'ongoing';
					$resourcesQuery.refetch();
				}}
			>
				<p class="text-md font-primary font-semibold">On-going</p>
				{#if status === 'ongoing'}
					<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
				{/if}
			</button>
			<button
				class="relative flex flex-col items-center outline-none border-none cursor-pointer"
				on:click={() => {
					status = 'bookmarked';
					$resourcesQuery.refetch();
				}}
			>
				<p class="text-md font-primary font-semibold">Bookmarks</p>
				{#if status === 'bookmarked'}
					<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
				{/if}
			</button>
		</div>

		{#if $resourcesQuery.isPending && !$resourcesQuery.data}
			<div class="flex flex-col gap-4">
				{#each { length: 2 } as _}
					<SkeletonCard className="w-full" />
				{/each}
			</div>
		{/if}

		{#if !$resourcesQuery.isPending && $resourcesQuery.data}
			<div class="w-full flex flex-wrap gap-4 mt-4 items-center justify-center">
				{#each $resourcesQuery.data.pages as resources}
					{#each resources as resource}
						<Card
							redirect="/resource"
							url={resource.url}
							name={resource.name}
							description={resource.description}
							author={resource.author}
							className="flex-col w-64 h-64"
						/>
					{/each}
				{/each}
			</div>
			<div class={classNames(!$resourcesQuery.hasNextPage && 'hidden')} bind:this={observerElem}>
				{#if $resourcesQuery.isFetchingNextPage && $resourcesQuery.hasNextPage}
					<Spinner />
				{/if}
			</div>
		{/if}
	</section>
{/if}
