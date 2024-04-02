<script lang="ts">
	import SearchBar from '../SearchBar.svelte';
</script>

<svelte:head>
	<title>Banoffee</title>
	<meta name="description" content="Discover page" />
</svelte:head>

<section class="flex flex-col gap-4">
	<h1 class="font-primary font-bold text-3xl tracking-wider">Discover the best of the web</h1>

	<SearchBar className="mb-4" />

	<p class="font-primary font-bold text-2xl tracking-wider">Popular this week</p>

	<div class="flex gap-2">
		<a
			href={'/resource-'}
			class="w-[50%] flex flex-col gap-4 h-[392px] shadow px-8 py-6 rounded-lg bg-white"
		>
			<span class="bg-primary-400 h-full aspect-video rounded-lg" />
			<div class="flex flex-col justify-between h-full">
				<div class="flex flex-col">
					<p class="font-primary text-2xl tracking-wider font-semibold">Title</p>
					<!-- {#if recommendation.description}
						<p class="text-sm line-clamp-2">{recommendation.description}</p>
					{/if} -->
				</div>
				<!-- {#if recommendation.author}
					<p>{recommendation.author}</p>
				{/if} -->
			</div>
		</a>
		<div class="w-[50%] h-[392px] grid grid-rows-1 grid-cols-2 gap-2">
			{#each { length: 4 } as _, i}
				<a
					href={'/resource-' + i}
					class="flex flex-col gap-4 h-48 shadow px-8 py-6 rounded-lg bg-white"
				>
					<span class="bg-primary-400 h-full aspect-video rounded-lg" />
					<div class="flex flex-col justify-between h-full">
						<div class="flex flex-col">
							<p class="font-primary text-2xl tracking-wider font-semibold">Title</p>
							<!-- {#if recommendation.description}
							<p class="text-sm line-clamp-2">{recommendation.description}</p>
						{/if} -->
						</div>
						<!-- {#if recommendation.author}
						<p>{recommendation.author}</p>
					{/if} -->
					</div>
				</a>
			{/each}
		</div>
	</div>

	<p class="font-primary font-bold text-2xl tracking-wider">Discover more topics</p>

	<div class="w-full flex flex-wrap items-center justify-center gap-2 px-4 sm:px-8 lg:px-32">
		{#each { length: 15 } as _, i}
			<div
				class="inline-block p-4 bg-transparent text-textPrimary text-sm font-semibold leading-none border border-gray-300 border-dividerTertiary rounded-md transition duration-150 ease-out font-primary hover:cursor-pointer hover:border-black hover:font-bold"
			>
				Category
			</div>
		{/each}
	</div>
</section>
