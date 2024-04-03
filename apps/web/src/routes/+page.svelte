<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { getRecommendations } from '../requests/user';
	import { type Resource } from '../schemas/resource';
	import Card from '../lib/components/Card.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';

	const recommendations = createQuery<Resource[]>({
		queryKey: ['recommendations'],
		queryFn: getRecommendations,
		retry: false
	});
</script>

<svelte:head>
	<title>Banoffee</title>
	<meta name="description" content="Home page" />
</svelte:head>

<section class="flex flex-col gap-8">
	<h1 class="font-primary font-bold text-6xl max-w-[50%] tracking-wider">You might like</h1>
	<div class="flex flex-col gap-4">
		{#if $recommendations.isPending && !$recommendations.data}
			{#each { length: 4 } as _}
				<SkeletonCard />
			{/each}
		{/if}

		{#if !$recommendations.isPending && $recommendations.data}
			{#each $recommendations.data as recommendation, i}
				<Card
					redirect="/resource"
					name={recommendation.name}
					description={recommendation.description}
					author={recommendation.author}
				/>
			{/each}
		{/if}
	</div>
</section>
