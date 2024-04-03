<script lang="ts">
	import { page } from '$app/stores';
	import { getResourceByURL, searchResourcesByName } from '../../requests/resource';
	import type { Resource } from '../../schemas/resource';
	import { isValidUrl } from '../../utils';
	import { createQuery, createInfiniteQuery } from '@tanstack/svelte-query';
	import Card from '$lib/components/Card.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';

	const search = $page.url.searchParams.get('search');
	const isURL = isValidUrl(search || '');

	const resource = createQuery<Resource | null>({
		queryKey: ['get-resource-by-url'],
		queryFn: () => getResourceByURL(search || ''),
		enabled: isURL,
		staleTime: 0,
		select: (data) => {
			if (!isURL) return null;

			return data;
		}
	});

	const searchQuery = createInfiniteQuery<Resource[]>({
		queryKey: ['search'],
		queryFn: ({ pageParam }) => searchResourcesByName(search || '', pageParam as number),
		initialPageParam: 0,
		getNextPageParam: (lastPage) => {
			return undefined;
		},
		enabled: !isURL,
		staleTime: 0,
		select: (data) => {
			if (isURL) {
				return {
					pages: [],
					pageParams: []
				};
			}

			return data;
		}
	});
</script>

<section>
	{#if isURL ? $resource.isPending && !$resource.data : $searchQuery.isPending && !$searchQuery.data}
		<div class="flex flex-col gap-4">
			{#each { length: 2 } as _}
				<SkeletonCard />
			{/each}
		</div>
	{/if}

	{#if !$resource.isPending && $resource.data}
		<Card
			redirect="/resource"
			name={$resource.data.name}
			description={$resource.data.description}
			author={$resource.data.author}
		/>
	{/if}

	{#if !$searchQuery.isPending && $searchQuery.data}
		<div class="flex flex-col gap-4">
			{#each $searchQuery.data.pages as resources}
				{#each resources as resource}
					<Card
						redirect="/resource"
						name={resource.name}
						description={resource.description}
						author={resource.author}
					/>
				{/each}
			{/each}
		</div>
	{/if}
</section>
