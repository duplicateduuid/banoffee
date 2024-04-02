<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { getRecommendations } from '../requests/user';
	import { type Resource } from '../schemas/user';

	const recommendations = createQuery<Resource[]>({
		queryKey: ['recommendations'],
		queryFn: getRecommendations
	});
</script>

<svelte:head>
	<title>Banoffee</title>
	<meta name="description" content="Home page" />
</svelte:head>

<section class="flex flex-col gap-8">
	<h1 class="font-primary font-bold text-6xl max-w-[50%] tracking-wider">You might like</h1>
	<div class="flex flex-col gap-4">
		<!-- TODO: add loading state/UI -->
		{#if !$recommendations.isPending && $recommendations.data}
			{#each $recommendations.data as recommendation, i}
				<!-- TODO: move to a card component -->
				<a href={'/resource-' + i} class="flex gap-4 h-40 shadow px-8 py-6 rounded-lg bg-white">
					<span class="bg-primary-400 h-full aspect-video rounded-lg" />
					<div class="flex flex-col justify-between h-full">
						<div class="flex flex-col">
							<p class="font-primary text-2xl tracking-wider font-semibold">
								{recommendation.name}
							</p>
							{#if recommendation.description}
								<p class="text-sm line-clamp-2">{recommendation.description}</p>
							{/if}
						</div>
						{#if recommendation.author}
							<p>{recommendation.author}</p>
						{/if}
					</div>
				</a>
			{/each}
		{/if}
	</div>
</section>
