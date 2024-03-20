<svelte:head>
  <title>Sign in</title>
</svelte:head>

<div class="flex h-full w-full">
  <section class="w-[600px] h-full hidden lg:block">
    <video src="videofoda.mp4" class="h-full object-cover" autoplay />
  </section>
  <section class="flex flex-col items-center lg:items-start lg:ml-[100px] justify-center w-full w-full">
    <div class="flex flex-col gap-10 w-full max-w-[400px]">
      <h1 class="font-primary font-bold text-primary-500 text-3xl">Sign in to Banoffee</h1>
      <div class="flex flex-col gap-8">
        <button class="bg-primary-400 hover:bg-primary-300 text-white rounded-full py-3 text-center w-full font-primary font-semibold text-xl transition">
          Sign up with Google
        </button>

        <div class="relative w-full flex items-center justify-center">
          <span class="absolute bg-white px-3 text-sm text-stone-600">or</span>
          <hr class="h-px bg-stone-200 w-full">
        </div>

        <button class="border border-stone-200 hover:border-stone-300 text-stone-900 rounded-full py-3.5 text-center w-full font-ligth transition">
          Continue with email
        </button>

        <p class="text-sm text-stone-700 flex w-full justify-center gap-1">
          Already have an account?
          <a class="underline" href="/signin">Sign In</a>
        </p>
      </div>
    </div>
  </section>
</div>
