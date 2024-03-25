
<script lang="ts">  
  let { onSignIn } = $props();
  
	import { superForm } from 'sveltekit-superforms';
  import { page } from "$app/stores";
	import { enhance } from '$app/forms';

  const { form, errors, message, constraints } = superForm($page.data.signUpForm);
</script>

<div
  class="flex flex-col items-center w-full gap-8 pt-8 pb-6 px-8"
>
  <h2 class="font-bold font-primary text-3xl w-full">
    Nice to meet you!
  </h2>

  <div class="flex flex-col gap-2.5 w-full">
    <button
      class="w-full w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center bg-stone-800 hover:bg-stone-700 text-white transition"
    >
      <div class="flex items-center gap-4">
        <img src="/icons/google.svg" alt="google icon" class="h-6 w-6" />
        Sign in with Google
      </div>
    </button>


    <button
      class="w-full w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center bg-stone-800 hover:bg-stone-700 text-white transition"
    >
      <div class="flex items-center gap-4">
        <img src="/icons/apple-white.svg" alt="google icon" class="h-6 w-6" />
        Sign in with Apple
      </div>
    </button>
  </div>

  <div class="relative w-full flex items-center justify-center">
    <span class="absolute bg-white px-3 text-sm text-stone-600"> or </span>
    <hr class="h-px bg-stone-200 w-full">
  </div>

  <button
    class="w-full w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center transition"
  >
    <div class="flex items-center gap-4">
      <img src="/icons/apple-white.svg" alt="google icon" class="h-6 w-6" />
      Continue with email
    </div>
  </button>

  {#if $message}
    <div class="w-full">
      <p class="bg-red-200 rounded-lg px-4 py-2 text-red-600 font-semibold text-center">
        {$message}
      </p>
    </div>
  {/if}

  <span class="flex items-center gap-1.5 text-stone-500 text-sm">
    Already have an account?
    <button
      onclick={onSignIn}
      class="underline text-stone-800 hover:text-stone-600 font-semibold"
    >
      Sign in
    </button>
  </span>
</div>
