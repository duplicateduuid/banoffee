<script lang="ts">
	import { UserRound } from 'lucide-svelte';
	import { type User } from '../../schemas/user';
	import { createDropdownMenu, melt } from '@melt-ui/svelte';
	import { fly } from 'svelte/transition';

	type Props = {
		user: User;
	};
	const { user }: Props = $props();

	const {
		elements: { menu, item, trigger },
		states: { open }
	} = createDropdownMenu({
		forceVisible: true,
		loop: true,
		positioning: { placement: 'bottom-end' }
	});
</script>

<button
	type="button"
	class="trigger bg-primary-400 rounded-full p-2.5 text-stone-800 hover:bg-primary-300 transition"
	use:melt={$trigger}
>
	<UserRound />
</button>

{#if $open}
	<div
		class="menu mt-4 min-w-[300px] p-6 bg-white rounded-2xl outline-none flex flex-col gap-4 border border-stone-100 shadow-lg"
		use:melt={$menu}
		transition:fly={{ duration: 100, y: -20 }}
	>
		<div class="user-info flex flex-col items-center gap-2.5 font-semibold text-xl w-full">
			<div class="avatar p-4 rounded-full bg-primary-400">
				<UserRound size={32} />
			</div>
			<p>{user.username}</p>
		</div>

		<div class="items flex flex-col items-start w-full">
			<button
				type="button"
				class="item text-lg text-stone-700 tracking-wide py-1.5 outline-none hover:text-stone-400 focus:text-stone-400 transition duration-75"
				use:melt={$item}
			>
				Bookmarks
			</button>
			<button
				type="button"
				class="item text-lg text-stone-700 tracking-wide py-1.5 outline-none hover:text-stone-400 focus:text-stone-400 transition duration-75"
				use:melt={$item}
			>
				Settings
			</button>
			<hr class="my-2 w-full" />
			<button
				type="button"
				class="item text-lg text-stone-700 tracking-wide py-1.5 outline-none hover:text-stone-400 focus:text-stone-400 transition duration-75"
				use:melt={$item}
			>
				Sign out
			</button>
		</div>
	</div>
{/if}
