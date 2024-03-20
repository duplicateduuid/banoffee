<svelte:head>
  <title>Sign in</title>
</svelte:head>

<div class="flex h-full w-full">
  <section class="w-[600px] h-full hidden lg:block">
    <video src="videofoda.mp4" class="h-full object-cover" autoplay />
  </section>
  <section class="flex flex-col items-center lg:items-start lg:ml-[100px] justify-center w-full w-full">
    <div class="flex flex-col gap-10 w-full max-w-[400px]">
      <h1 class="font-primary font-semibold text-primary-500 text-3xl">
        Sign in to
        <span class="font-bold">Banoffee</span>
      </h1>
      <div class="flex flex-col gap-8">
        <button class="border border-stone-200 hover:border-stone-300 text-stone-900 rounded-full py-3.5 text-center w-full font-semibold transition">
          Sign in with Google
        </button>

        <div class="relative w-full flex items-center justify-center">
          <span class="absolute bg-white px-3 text-sm text-stone-600">or sign in with email</span>
          <hr class="h-px bg-stone-200 w-full">
        </div>

        <div class="flex flex-col gap-6 mt-2">
          <div class="flex flex-col gap-1">
            <label for="login" class="font-bold text-stone-800">Username or Email</label>
            <input
              id="login"
              class="rounded-xl border border-stone-200 focus:border-primary-300 py-3 px-5 text-sm transition outline-none ring-primary-100/90 hover:ring-4 focus:ring-4"
            />
          </div>

          <div class="flex flex-col gap-1">
            <div class="flex justify-between justify-center">
              <label for="password" class="font-bold text-stone-800">Password</label>
              <a class="underline" href="/reset">Forgot?</a>
            </div>
            <input
              id="password"
              class="rounded-xl border border-stone-200 focus:border-primary-300 py-3 px-5 text-sm transition outline-none ring-primary-100/90 hover:ring-4 focus:ring-4"
              type="password"
            />
          </div>
        </div>

        <div class="flex flex-col gap-5 mt-2 items-center">
          <button class="bg-primary-400 hover:bg-primary-300 text-white rounded-full py-3 text-center w-full font-primary font-bold text-xl transition">
            Login
          </button>

          <p class="text-sm text-stone-700">
            Don't have an account?
            <a class="underline" href="signup">Sign up</a>
          </p>
        </div>
      </div>
    </div>
  </section>
</div>
