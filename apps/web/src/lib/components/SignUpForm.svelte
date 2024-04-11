<script lang="ts">
	import { ChevronLeft } from 'lucide-svelte';
	import { superForm } from 'sveltekit-superforms';
	import { page } from '$app/stores';

	type Props = {
		onSignIn: () => void;
	};

	let { onSignIn }: Props = $props();

	const { form, errors, constraints, message, enhance } = superForm($page.data.signUpForm);

	type Step = 'ask-sign-alternative' | 'sign-with-email';

	let step: Step = $state('ask-sign-alternative');
</script>

<div class="flex flex-col items-center w-full gap-8 pt-8 pb-6 px-8">
	<div class="flex gap-4 w-full items-center">
		{#if step === 'sign-with-email'}
			<button
				type="button"
				onclick={() => (step = 'ask-sign-alternative')}
				class="shadow rounded-lg p-1"
			>
				<ChevronLeft />
			</button>
		{/if}

		<h2 class="font-bold font-primary text-3xl w-full">Nice to meet you!</h2>
	</div>

	{#if step === 'ask-sign-alternative'}
		{@render alternatives()}
	{:else if step === 'sign-with-email'}
		{@render emailForm()}
	{/if}
</div>

{#snippet alternatives()}
	<div class="flex flex-col gap-2.5 w-full">
		<button
			class="w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center bg-stone-800 hover:bg-stone-700 text-white transition"
		>
			<div class="flex items-center gap-4">
				<img src="/icons/google.svg" alt="google icon" class="h-6 w-6" />
				Sign in with Google
			</div>
		</button>

		<button
			class="w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center bg-stone-800 hover:bg-stone-700 text-white transition"
		>
			<div class="flex items-center gap-4">
				<img src="/icons/apple-white.svg" alt="google icon" class="h-6 w-6" />
				Sign in with Apple
			</div>
		</button>
	</div>

	<div class="relative w-full flex items-center justify-center">
		<span class="absolute bg-white px-3 text-sm text-stone-600"> or </span>
		<hr class="h-px bg-stone-200 w-full" />
	</div>

	<button
		onclick={() => (step = 'sign-with-email')}
		class="w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center transition"
	>
		Continue with email
	</button>

	{#if $message}
		<div class="w-full">
			<p class="bg-red-200 rounded-lg px-4 py-2 text-red-600 font-semibold text-center">
				{$message}
			</p>
		</div>
	{/if}

	{@render signIn()}
{/snippet}

{#snippet signIn()}
	<span class="flex items-center gap-1.5 text-stone-500 text-sm">
		Already have an account?
		<button
			type="submit"
			class="underline text-stone-800 hover:text-stone-600 font-semibold"
			onclick={onSignIn}
		>
			Sign in
		</button>
	</span>
{/snippet}

{#snippet emailForm()}
	<form
		use:enhance={{
			onResult: () => {
				// TODO: find a better way to reload the page.
				// had tried invalidateAll(), but the user is loaded in the server side;
				// had also tried goto("/"), but im already at "/", so nothing happens;
				window.location.reload();
			}
		}}
		method="POST"
		action="?/signUp"
		class="w-full flex flex-col gap-8"
	>
		<div class="flex flex-col gap-4 w-full">
			<fieldset class="flex flex-col gap-0.5 text-stone-800">
				<label class="text-stone-800 tracking-tight" for="login"> Username </label>
				<input
					id="username"
					name="username"
					aria-invalid={$errors.username ? 'true' : undefined}
					bind:value={$form.username}
					{...$constraints.username}
					class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
            outline-none transition border-solid border {$errors.username
						? 'border-rose-500'
						: 'focus:border-primary-300 hover:ring-2 focus:ring-2 ring-primary-100 ring-0'}"
				/>

				{#if $errors.username}
					<p class="text-rose-500 text-xs font-semibold">
						{$errors.username}
					</p>
				{/if}
			</fieldset>

			<fieldset class="flex flex-col gap-0.5 text-stone-800">
				<label class="text-stone-800 tracking-tight" for="login"> Email </label>
				<input
					id="email"
					name="email"
					aria-invalid={$errors.email ? 'true' : undefined}
					bind:value={$form.email}
					{...$constraints.email}
					class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
            outline-none transition border-solid border {$errors.email
						? 'border-rose-500'
						: 'focus:border-primary-300 hover:ring-2 focus:ring-2 ring-primary-100 ring-0'}"
				/>

				{#if $errors.email}
					<p class="text-rose-500 text-xs font-semibold">
						{$errors.email}
					</p>
				{/if}
			</fieldset>

			<fieldset class="flex flex-col gap-0.5 text-stone-800">
				<label class="text-stone-800 tracking-tight" for="password"> Password </label>
				<input
					id="password"
					name="password"
					type="password"
					aria-invalid={$errors.password ? 'true' : undefined}
					bind:value={$form.password}
					{...$constraints.password}
					class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
            outline-none transition border-solid border {$errors.password
						? 'border-rose-500'
						: 'focus:border-primary-300 hover:ring-2 focus:ring-2 ring-primary-100 ring-0'}"
				/>
				{#if $errors.password}
					<p class="text-rose-500 text-xs font-semibold">
						{$errors.password}
					</p>
				{/if}
			</fieldset>
		</div>

		<hr class="w-[calc(100%+64px)] mx-[-32px]" />

		<div class="flex justify-between gap-4 w-full">
			{@render signIn()}

			<button
				type="submit"
				class="py-3 items-center justify-center rounded-lg bg-primary-400 hover:bg-primary-300
          text-white font-semibold px-4 leading-none text-magnum-900 shadow"
			>
				Sign in
			</button>
		</div>
	</form>
{/snippet}
