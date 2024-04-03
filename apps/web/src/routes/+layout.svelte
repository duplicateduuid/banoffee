<script lang="ts">
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import '../app.css';
	import type { LayoutData } from './$types';
	import AccountDropdown from '../lib/components/AccountDropdown.svelte';
	import SignDialog from '../lib/components/SignDialog.svelte';
	import { page } from '$app/stores';
	import classnames from 'classnames';

	export let data: LayoutData;
	const { user } = data;
</script>

<main class="w-full min-h-screen font-secondary bg-[#F7F6F1] text-[#18171C]">
	<QueryClientProvider client={data.queryClient}>
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="mx-auto max-w-4xl flex flex-col gap-16 lg:max-w-6xl">
				<!-- TODO: create a header component -->
				<header class="grid grid-cols-4 w-full justify-between py-8 items-center">
					<span class="font-primary font-bold text-5xl">Banoffee</span>
					<div class="col-span-2 w-full flex items-center justify-center gap-10">
						<div class="relative flex flex-col items-center">
							<a
								href="/"
								class={classnames('text-2xl', { 'text-[#7A7974]': $page.url.pathname !== '/' })}
								>For you</a
							>
							{#if $page.url.pathname === '/'}
								<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
							{/if}
						</div>
						<div class="relative flex flex-col items-center">
							<a
								href="/discover"
								class={classnames('text-2xl', {
									'text-[#7A7974]': $page.url.pathname !== '/discover'
								})}>Discover</a
							>
							{#if $page.url.pathname === '/discover'}
								<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
							{/if}
						</div>
					</div>

					{@render signIn()}
				</header>

				<slot />

				<footer class="py-8">
					<hr class="w-full border-gray-500" />
					<div class="mx-auto max-w-7xl overflow-hidden lg:px-8">
						<div class="mt-10 flex justify-center space-x-10">
							<a href="#" class="text-gray-400 hover:text-gray-500">
								<span class="sr-only">Discord</span>
								<svg
									class="w-6 h-6"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="currentColor"
									viewBox="0 0 21 16"
								>
									<path
										d="M16.942 1.556a16.3 16.3 0 0 0-4.126-1.3 12.04 12.04 0 0 0-.529 1.1 15.175 15.175 0 0 0-4.573 0 11.585 11.585 0 0 0-.535-1.1 16.274 16.274 0 0 0-4.129 1.3A17.392 17.392 0 0 0 .182 13.218a15.785 15.785 0 0 0 4.963 2.521c.41-.564.773-1.16 1.084-1.785a10.63 10.63 0 0 1-1.706-.83c.143-.106.283-.217.418-.33a11.664 11.664 0 0 0 10.118 0c.137.113.277.224.418.33-.544.328-1.116.606-1.71.832a12.52 12.52 0 0 0 1.084 1.785 16.46 16.46 0 0 0 5.064-2.595 17.286 17.286 0 0 0-2.973-11.59ZM6.678 10.813a1.941 1.941 0 0 1-1.8-2.045 1.93 1.93 0 0 1 1.8-2.047 1.919 1.919 0 0 1 1.8 2.047 1.93 1.93 0 0 1-1.8 2.045Zm6.644 0a1.94 1.94 0 0 1-1.8-2.045 1.93 1.93 0 0 1 1.8-2.047 1.918 1.918 0 0 1 1.8 2.047 1.93 1.93 0 0 1-1.8 2.045Z"
									/>
								</svg>
							</a>
							<a href="#" class="text-gray-400 hover:text-gray-500">
								<span class="sr-only">GitHub</span>
								<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
									<path
										fill-rule="evenodd"
										d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
										clip-rule="evenodd"
									/>
								</svg>
							</a>
							<a href="#" class="text-gray-400 hover:text-gray-500">
								<span class="sr-only">X</span>
								<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
									<path
										d="M13.6823 10.6218L20.2391 3H18.6854L12.9921 9.61788L8.44486 3H3.2002L10.0765 13.0074L3.2002 21H4.75404L10.7663 14.0113L15.5685 21H20.8131L13.6819 10.6218H13.6823ZM11.5541 13.0956L10.8574 12.0991L5.31391 4.16971H7.70053L12.1742 10.5689L12.8709 11.5655L18.6861 19.8835H16.2995L11.5541 13.096V13.0956Z"
									/>
								</svg>
							</a>
						</div>
						<p class="mt-10 text-center text-xs leading-5 text-gray-500">&copy; Banoffee.</p>
					</div>
				</footer>
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
