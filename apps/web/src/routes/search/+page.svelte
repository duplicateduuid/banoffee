<script lang="ts">
	import { page } from '$app/stores';
	import { getResourceByURL, searchResourcesByName } from '../../requests/resource';
	import type { Resource } from '../../schemas/resource';
	import { isValidUrl } from '../../utils';
	import { createQuery, createInfiniteQuery } from '@tanstack/svelte-query';

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
                }
            }

            return data;
        }
	});
</script>

<!-- TODO: add loading state/UI & no resources founded state/UI -->
<section>
	{#if !$resource.isPending && $resource.data}
		<p>Resource Found: {$resource.data.name}</p>
	{/if}

	{#if !$searchQuery.isPending && $searchQuery.data}
		<div>
			{#each $searchQuery.data.pages as resources}
				{#each resources as resource}
					<p>Resource Found: {resource.name}</p>
				{/each}
			{/each}
		</div>
	{/if}
</section>
