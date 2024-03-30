<script lang="ts">
  type Props = { onSignUp: () => void, onSubmitted: () => void };
  let { onSignUp, onSubmitted } = $props<Props>();
  
	import { createForm } from 'felte';
  import { validator } from "@felte/validator-zod";
	import { type SignInRequestType, signInRequestSchema, signInRequest } from '../../requests/auth';

  const { form, errors } = createForm<SignInRequestType>({
    extend: [validator({ schema: signInRequestSchema })],
    onSubmit: async (fields) => {
      const { user } = await signInRequest(fields);
      localStorage.setItem("user", JSON.stringify(user))

      onSubmitted();
    }
  })
</script>

<form
  use:form
  class="flex flex-col items-center w-full gap-8 pt-8 pb-6"
>
  <h2 class="font-bold font-primary text-3xl w-full px-8">
    Welcome back!
  </h2>

  <div class="flex flex-col gap-2.5 w-full px-8">
    <button
      type="button"
      class="w-full w-full rounded-lg py-2.5 font-semibold shadow flex items-center justify-center hover:bg-stone-50 transition"
    >
      <div class="flex items-center gap-4">
        <img src="/icons/google.svg" alt="google icon" class="h-6 w-6" />
        Sign in with Google
      </div>
    </button>


    <button
      type="button"
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
        class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
          outline-0 transition border-solid border {
            $errors.login
                ? 'border-rose-500'
                : 'focus:border-primary-300 hover:ring-2 focus:ring-2 ring-primary-100 ring-0'
          }"
      />

      {#if $errors.login}
        <p class="text-rose-500 text-xs font-semibold">
          {$errors.login[0]}
        </p>
      {/if}
    </fieldset>

    <fieldset class="flex flex-col gap-0.5 text-stone-800">
      <label class="text-stone-800 tracking-tight" for="password"> Password </label>
      <input
        id="password"
        name="password"
        type="password"

        class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
          outline-0 transition border-solid border {
            $errors.password
                ? 'border-rose-500'
                : 'focus:border-primary-300 hover:ring-2 focus:ring-2 ring-primary-100 ring-0'
          }"
      />
      {#if $errors.password}
        <p class="text-rose-500 text-xs font-semibold">
          {$errors.password[0]}
        </p>
      {/if}
    </fieldset>
  </div>  

  <!--
  TODO: render a server error message here
  {#if $message}
    <div class="w-full px-8">
      <p class="bg-red-200 rounded-lg px-4 py-2 text-red-600 font-semibold text-center">
        {$message}
      </p>
    </div>
  {/if}
  -->

  <hr class="w-full">

  <div class="flex justify-between gap-4 w-full px-8">
    <span class="flex items-center gap-1.5 text-stone-500 text-sm">
      Don't have an account?
      <button
        type="button"
        onclick={onSignUp}
        class="underline text-stone-800 hover:text-stone-600 font-semibold"
      >
        Sign up
      </button>
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
