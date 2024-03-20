<script>
	import { MoveRight } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Home page" />
</svelte:head>

<section>
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
	  <div class="mx-auto max-w-4xl flex flex-col gap-16">
			<!-- TODO: create a header component -->
			<header class="flex justify-between py-8 items-center">
				<span class="font-primary font-bold text-5xl">Banoffee</span>
				<div class="flex items-center gap-10">
					<div class="relative flex flex-col items-center">
						<a href="/discover" class="text-2xl">Discover</a>
						<span class="bg-primary-400 h-1 absolute mt-1 top-full w-[60%] rounded-full" />
					</div>
					<a href="/bookmarks" class="text-2xl text-[#7A7974]">Bookmarks</a>
				</div>
				<a href="/signin" class="flex items-center gap-6">
					<span class="text-primary font-bold text-2xl">Sign in</span>
					<MoveRight />
				</a>
			</header>

			<section class="flex flex-col gap-8">
				<h1 class="font-primary font-bold text-6xl max-w-[50%] tracking-wider">Popular this week</h1>		
				<div class="flex flex-col gap-4">

					{#each {length: 4} as _, i}
						<!-- TODO: move to a card component -->
						<a href={"/resource-" + i} class="flex gap-4 h-40 shadow px-8 py-6 rounded-lg bg-white">
							<span class="bg-primary-400 h-full aspect-video rounded-lg" />
							<div class="flex flex-col justify-between h-full">
								<div class="flex flex-col">
									<p class="font-primary text-2xl tracking-wider font-semibold">Designing Data-Intensive Applications</p>
									<p class="text-sm line-clamp-2">In this practical and comprehensive guide, author Martin Kleppmann helps you navigate this diverse landscape by examining the pros and cons of various technologies for processing and storing data. </p>
								</div>
								<p>Martin Kleppmann</p>
							</div>
						</a>
					{/each}
				</div>
			</section>

	  </div>
	</div>
</section>
