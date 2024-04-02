<script lang="ts">
	import SignDialog from './SignDialog.svelte';
	import AccountDropdown from './AccountDropdown.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { getRecommendations, me } from "../requests/user"
	import { type Resource, type User } from '../schemas/user';

	const user = createQuery<User>({
		queryKey: ['me'],
		queryFn: me,
	})

	const recommendations = createQuery<Resource[]>({
		queryKey: ['recommendations'],
		queryFn: getRecommendations
	})
</script>

<svelte:head>
	<title>Banoffee</title>
	<meta name="description" content="Home page" />
</svelte:head>

<section>
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
	  <div class="mx-auto max-w-4xl flex flex-col gap-16">
			<!-- TODO: create a header component -->
			<header class="grid grid-cols-4 w-full justify-between py-8 items-center">
				<span class="font-primary font-bold text-5xl">Banoffee</span>
				<div class="col-span-2 w-full flex items-center justify-center gap-10">
					<div class="relative flex flex-col items-center">
						<a href="/discover" class="text-2xl">For you</a>
						<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
					</div>
					<a href="/bookmarks" class="text-2xl text-[#7A7974]">Discover</a>
				</div>
				
				{@render signIn()}
			</header>

			<section class="flex flex-col gap-8">
				<h1 class="font-primary font-bold text-6xl max-w-[50%] tracking-wider">You might like</h1>		
				<div class="flex flex-col gap-4">
					{#if !$recommendations.isPending && $recommendations.data}
						{#each $recommendations.data as recommendation, i}
							<!-- TODO: move to a card component -->
							<a href={"/resource-" + i} class="flex gap-4 h-40 shadow px-8 py-6 rounded-lg bg-white">
								<span class="bg-primary-400 h-full aspect-video rounded-lg" />
								<div class="flex flex-col justify-between h-full">
									<div class="flex flex-col">
										<p class="font-primary text-2xl tracking-wider font-semibold">{recommendation.name}</p>
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

	  </div>
	</div>
</section>

{#snippet signIn()}
	<div class="flex w-full justify-end">
		{#if !$user.isPending}
			{#if $user.data}
				<AccountDropdown user={$user.data} />
			{:else}
				<SignDialog />
			{/if}
		{/if}
	</div>
{/snippet}
