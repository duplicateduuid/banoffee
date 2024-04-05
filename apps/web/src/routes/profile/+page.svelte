<script lang="ts">
	import { HeartHandshakeIcon, UserRound } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Card from '$lib/components/Card.svelte';

	let state = $state<'completed' | 'on-going' | 'bookmarked'>('completed');

	$effect(() => {
		if (!$page.data.user) goto('/');
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
					state = 'completed';
				}}
			>
				<p class="text-md font-primary font-semibold">Completed</p>
				{#if state === 'completed'}
					<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
				{/if}
			</button>
			<button
				class="relative flex flex-col items-center outline-none border-none cursor-pointer"
				on:click={() => {
					state = 'on-going';
				}}
			>
				<p class="text-md font-primary font-semibold">On-going</p>
				{#if state === 'on-going'}
					<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
				{/if}
			</button>
			<button
				class="relative flex flex-col items-center outline-none border-none cursor-pointer"
				on:click={() => {
					state = 'bookmarked';
				}}
			>
				<p class="text-md font-primary font-semibold">Bookmarks</p>
				{#if state === 'bookmarked'}
					<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
				{/if}
			</button>
		</div>
		<div class="w-full flex flex-wrap gap-4 mt-4 items-center justify-center">
			{#each { length: 8 } as _}
				<Card name="Test" redirect="/" url="" className="flex-col w-64 h-64" />
			{/each}
		</div>
	</section>
{/if}
