<script lang="ts">
	import SignInForm from '$lib/components/SignInForm.svelte';
	import SignUpForm from '$lib/components/SignUpForm.svelte';
	import { createDialog, melt } from '@melt-ui/svelte';
	import { MoveRight, X } from 'lucide-svelte';
	import { fade } from 'svelte/transition';
	import { page } from '$app/stores';

	const {
		elements: { trigger, overlay, content, close, portalled },
		states: { open }
	} = createDialog({
		forceVisible: true
	});

	let formType: 'sign-in' | 'sign-up' = $state('sign-in');
</script>

<button use:melt={$trigger} class="flex items-center gap-2 font-bold text-2xl font-primary px-2">
	Sign in
	<MoveRight />
</button>

<div class="" use:melt={$portalled}>
	{#if $open}
		<div
			use:melt={$overlay}
			class="fixed inset-0 z-50 bg-black/50"
			transition:fade={{ duration: 150 }}
		/>
		<div
			class="fixed left-[50%] top-[50%] z-50 max-h-[85vh] w-[90vw]
            max-w-[500px] translate-x-[-50%] translate-y-[-50%] rounded-xl bg-white
            shadow-lg relarive"
			use:melt={$content}
		>
			<button
				use:melt={$close}
				class="p-3 border bg-white hover:bg-stone-50 shadow rounded-lg absolute -top-4 -right-4 transition"
			>
				<X size="20" />
			</button>
			{#if formType === 'sign-in'}
				<SignInForm onSignUp={() => (formType = 'sign-up')} />
			{:else if formType === 'sign-up'}
				<SignUpForm onSignIn={() => (formType = 'sign-in')} />
			{/if}
		</div>
	{/if}
</div>
