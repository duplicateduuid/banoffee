<script lang="ts">
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import '../app.css';
	import type { LayoutData } from './$types';
	import AccountDropdown from './AccountDropdown.svelte';
	import SignDialog from './SignDialog.svelte';

	export let data: LayoutData;
	const { user } = data;
</script>

<main class="w-full min-h-screen font-secondary bg-[#F7F6F1] text-[#18171C]">
	<QueryClientProvider client={data.queryClient}>
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="mx-auto max-w-4xl flex flex-col gap-16">
				<!-- TODO: create a header component -->
				<header class="grid grid-cols-4 w-full justify-between py-8 items-center">
					<span class="font-primary font-bold text-5xl">Banoffee</span>
					<div class="col-span-2 w-full flex items-center justify-center gap-10">
						<div class="relative flex flex-col items-center">
							<a href="/" class="text-2xl">For you</a>
							<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
						</div>
						<a href="/discover" class="text-2xl text-[#7A7974]">Discover</a>
					</div>

					{@render signIn()}
				</header>

				<slot />
			</div>
		</div>
	</QueryClientProvider>
</main>

{#snippet signIn()}
	<div class="flex w-full justify-end">
		{#if user}
			<AccountDropdown {user} />
		{:else}
			<SignDialog />
		{/if}
	</div>
{/snippet}
