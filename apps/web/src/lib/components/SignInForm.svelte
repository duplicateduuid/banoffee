<script lang="ts">
	import SuperDebug, { superForm } from 'sveltekit-superforms';
  import { page } from "$app/stores";
	import { enhance } from '$app/forms';

  const { form } = superForm($page.data.form);
</script>

<form
  use:enhance
  method="POST"
  class="flex flex-col items-center w-full gap-8 pt-8 pb-6"
>
  <h2 class="font-bold font-primary text-3xl w-full px-8">
    Welcome back!
  </h2>

  <div class="flex flex-col gap-2.5 w-full px-8">
    <SuperDebug data={$form} />
    <button
      class="w-full w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center hover:bg-stone-50 transition"
    >
      <div class="flex items-center gap-4">
        <img src="/icons/google.svg" alt="google icon" class="h-6 w-6" />
        Sign in with Google
      </div>
    </button>


    <button
      class="w-full w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center hover:bg-stone-50 transition"
    >
      <div class="flex items-center gap-4">
        <img src="/icons/apple.svg" alt="google icon" class="h-6 w-6" />
        Sign in with Apple
      </div>
    </button>
  </div>

  <div class="relative w-full flex items-center justify-center px-8">
    <span class="absolute bg-white px-3 text-sm text-stone-600"> or </span>
    <hr class="h-px bg-stone-200 w-full">
  </div>

  <div class="flex flex-col gap-4 w-full px-8">
    <fieldset class="flex flex-col gap-0.5 text-stone-800">
      <label class="text-stone-800 tracking-tight" for="login"> Username or Email </label>
      <input
        id="login"
        name="login"
        bind:value={$form.login}
        class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
          outline-0 ring-0 hover:ring-2 focus:ring-2 ring-primary-100 transition
          border border-solid focus:border-primary-300"
      />
    </fieldset>

    <fieldset class="flex flex-col gap-0.5 text-stone-800">
      <label class="text-stone-800 tracking-tight" for="password"> Password </label>
      <input
        id="password"
        name="password"
        type="password"
        bind:value={$form.password}
        class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
          outline-0 ring-0 hover:ring-2 focus:ring-2 ring-primary-100 transition
          border border-solid focus:border-primary-300"
      />    
    </fieldset>
  </div>

  <hr class="w-full">

  <div class="flex justify-between gap-4 w-full px-8">
    <span class="flex items-center gap-1.5 text-stone-500 text-sm">
      Don't have an account?
      <button class="underline text-stone-800 hover:text-stone-600 font-semibold">Sign up</button>
    </span>

    <button
      type="submit"
      class="py-3 items-center justify-center rounded-lg bg-primary-400 hover:bg-primary-300
        text-white font-semibold px-4 font-medium leading-none text-magnum-900 shadow"
    >
      Sign in
    </button>
  </div>
</form>
