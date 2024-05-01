<script lang="ts">
	import type { Resource } from '../../schemas/resource';
	import { PackagePlusIcon, UserRound, X } from 'lucide-svelte';
	import { type User } from '../../schemas/user';
	import { createDialog, createDropdownMenu, melt } from '@melt-ui/svelte';
	import { fade, fly } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import * as z from 'zod';
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import classNames from 'classnames';
	import { getResourceByURL } from '../../requests/resource';
	import { api } from '../../api';
	import { deleteCookie, removeQueryParams } from '../../utils';

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

	const {
		elements: {
			trigger: dialogTrigger,
			overlay: dialogOverlay,
			content: dialogContent,
			close: dialogClose,
			portalled: dialogPortalled
		},
		states: { open: dialogOpen }
	} = createDialog({
		forceVisible: true
	});

	const schema = z.object({
		url: z.string().url()
	});

	const {
		form,
		data: formData,
		errors: formErrors
	} = createForm<z.infer<typeof schema>>({
		extend: validator({ schema }),
		onSubmit: (values) =>
			getResourceByURL(values.url)
				.then(async (resource) => {
					await api.post(`/user/resource/${resource.id}`, { status: 'bookmarked' });
					dialogOpen.set(false);
				})
				.catch(async () => {
					const title = await fetch(`/api/getTitleFromURL?url=${removeQueryParams(values.url)}`)
						.then((res) => res.text())
						.catch(() => null);

					const {
						data: { resource: newResource }
					} = await api.post<{ resource: Resource }>('/resource', {
						url: values.url,
						// TODO: Define a better fallback to missing title
						name: title || values.url
					});

					await api.post(`/user/resource/${newResource.id}`, { status: 'bookmarked' });
					dialogOpen.set(false);
				})
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
				on:click={() => goto('/profile')}
			>
				Profile
			</button>
			<button
				type="button"
				class="item text-lg text-stone-700 tracking-wide py-1.5 outline-none hover:text-stone-400 focus:text-stone-400 transition duration-75"
				use:melt={$item}
				on:click={() => goto('/discover')}
			>
				Discover
			</button>
			<button
				type="button"
				class="item text-lg text-stone-700 tracking-wide py-1.5 outline-none hover:text-stone-400 focus:text-stone-400 transition duration-75"
				use:melt={$item}
			>
				Settings
			</button>
			<button
				type="button"
				class="item text-lg text-stone-700 tracking-wide py-1.5 outline-none hover:text-stone-400 focus:text-stone-400 transition duration-75 flex gap-1 items-center"
				use:melt={$item}
				use:melt={$dialogTrigger}
			>
				Save a URL
				<PackagePlusIcon />
			</button>
			<hr class="my-2 w-full" />
			<button
				type="button"
				class="item text-lg text-stone-700 tracking-wide py-1.5 outline-none hover:text-stone-400 focus:text-stone-400 transition duration-75"
				use:melt={$item}
				on:click={() => {
					deleteCookie('sessionId');
					window.location.reload();
				}}
			>
				Sign out
			</button>
		</div>
	</div>
{/if}

<div class="" use:melt={$dialogPortalled}>
	{#if $dialogOpen}
		<div
			use:melt={$dialogOverlay}
			class="fixed inset-0 z-50 bg-black/50"
			transition:fade={{ duration: 150 }}
		/>
		<form
			class="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-50 min-w-[90%] sm:min-w-[30em] rounded-xl flex flex-col justify-between gap-4 p-8 font-secondary bg-white shadow-xl"
			use:melt={$dialogContent}
			use:form
		>
			<button
				use:melt={$dialogClose}
				class="p-2 border bg-white hover:bg-stone-50 shadow rounded-lg absolute -top-4 -right-4 transition"
			>
				<X size={16} />
			</button>

			<div class="w-full flex flex-col gap-0.5 text-ellipsis whitespace-nowrap overflow-hidden">
				<p class="text-sm font-normal">URL:</p>
				<input
					placeholder="Paste the URL here"
					id="url"
					name="url"
					class={classNames(
						'px-3 py-2 w-full items-center justify-center rounded-lg text-black outline-none transition border-solid border focus:border-primary-300 hover:ring-2 focus:ring-2 ring-primary-100 ring-0',
						{ 'border-rose-500': !!$formErrors.url }
					)}
				/>
				{#if $formErrors.url}
					<p class="text-rose-500 text-xs font-semibold">
						{$formErrors.url}
					</p>
				{/if}
			</div>

			<button
				class="w-full h-12 rounded-md border-none bg-[#4e473b] text-[#F7F6F1] text-base inline-block disabled:opacity-50 disabled:cursor-not-allowed"
				disabled={!$formData.url || !!$formErrors.url}
			>
				Bookmark
			</button>
		</form>
	{/if}
</div>
