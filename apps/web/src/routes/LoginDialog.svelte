<script lang="ts">
  import { createDialog, melt } from '@melt-ui/svelte';
  import { X } from "lucide-svelte";
  /** Internal helpers */
  import { fade } from 'svelte/transition';

  const {
    elements: {
      trigger,
      overlay,
      content,
      title,
      description,
      close,
      portalled,
    },
    states: { open },
  } = createDialog({
    forceVisible: true,
  });
</script>

<button
  use:melt={$trigger}
  class="inline-flex items-center justify-center rounded-xl bg-white px-4 py-3
  font-medium leading-none text-magnum-700 shadow hover:opacity-75"
>
  Sign in
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
        <X size=20 />
      </button>
      <form method="POST" action="/sign-in" class="flex flex-col items-center w-full gap-8 pt-8 pb-6">
        <h2 use:melt={$title} class="font-bold font-primary text-3xl w-full px-8">
          Welcome back!
        </h2>

        <div class="flex flex-col gap-2.5 w-full px-8">
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
              class="px-3 py-2 w-full items-center justify-center rounded-lg text-black
                outline-0 ring-0 hover:ring-2 focus:ring-2 ring-primary-100 transition
                border border-solid focus:border-primary-300"
            />
          </fieldset>

          <fieldset class="flex flex-col gap-0.5 text-stone-800">
            <label class="text-stone-800 tracking-tight" for="password"> Password </label>
            <input
              id="password"
              type="password"
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
            use:melt={$close}
            class="py-3 items-center justify-center rounded-lg bg-primary-400 hover:bg-primary-300
              text-white font-semibold px-4 font-medium leading-none text-magnum-900 shadow"
          >
            Sign in
          </button>
        </div>
      </form>
    </div>
  {/if}
</div>
