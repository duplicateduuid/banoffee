<script lang="ts">
	import SearchBar from '$lib/components/SearchBar.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { getPopularThisWeek } from '../../requests/user';
	import type { Resource } from '../../schemas/resource';
	import Card from '$lib/components/Card.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';

	const popularResources = createQuery<Resource[]>({
		queryKey: ['popular-this-week'],
		queryFn: getPopularThisWeek,
		staleTime: 0
	});
</script>

<svelte:head>
	<title>Banoffee</title>
	<meta name="description" content="Discover page" />
</svelte:head>

<section class="flex flex-col gap-4">
	<h1 class="font-primary font-bold text-2xl sm:text-3xl tracking-wider">Discover the best of the web</h1>

	<SearchBar className="mb-4" />

	{#if $popularResources.isPending || ($popularResources.data && $popularResources.data.length > 0)}
		<p class="font-primary font-bold text-xl sm:text-2xl tracking-wider">Popular this week</p>
	{/if}

	{#if $popularResources.isPending && !$popularResources.data}
		<div class="flex flex-col gap-4">
			{#each { length: 2 } as _}
				<SkeletonCard />
			{/each}
		</div>
	{/if}

	{#if !$popularResources.isPending && $popularResources.data && $popularResources.data.length > 0}
		<div class="flex gap-2 flex-col sm:flex-row">
			<Card
				className="w-full sm:w-[50%] flex flex-col h-[392px] shadow px-8 py-6 rounded-lg bg-white"
				redirect="/resource"
				url={$popularResources.data[0].url}
				name={$popularResources.data[0].name}
				description={$popularResources.data[0].description}
				author={$popularResources.data[0].author}
			/>
			<div class="w-full sm:w-[50%] h-[392px] grid grid-rows-1 grid-cols-2 gap-2">
				{#each $popularResources.data.splice(1) as resource, i}
					<Card
						className="flex flex-col h-48 shadow px-8 py-6 rounded-lg bg-white"
						redirect="/resource"
						url={resource.url}
						name={resource.name}
						description={resource.description}
						author={resource.author}
					/>
				{/each}
			</div>
		</div>
	{/if}

	<!-- TODO: Categories -->
	<!-- <p class="font-primary font-bold text-2xl tracking-wider">Discover more topics</p>

	<div class="w-full flex flex-wrap items-center justify-center gap-2 px-4 sm:px-8 lg:px-32">
		{#each { length: 15 } as _, i}
			<div
				class="inline-block p-4 bg-transparent text-textPrimary text-sm font-semibold leading-none border border-gray-300 border-dividerTertiary rounded-md transition duration-150 ease-out font-primary hover:cursor-pointer hover:border-black hover:font-bold"
			>
				Category
			</div>
		{/each}
	</div> -->
</section>
