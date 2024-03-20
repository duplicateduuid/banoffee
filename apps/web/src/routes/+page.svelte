<script>
	import { MoveRight } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Home page" />
</svelte:head>

<section>
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
	  <div class="mx-auto max-w-4xl">
			<div class="flex justify-between py-8 items-center">
				<span class="font-primary font-bold text-5xl">Banoffee</span>
				<div class="flex items-center gap-10">
					<div class="relative flex flex-col items-center">
						<a href="/discover" class="text-2xl">Discover</a>
						<span class="bg-[#FFB800] h-1 absolute mt-1 top-full w-[60%] rounded-full" />
					</div>
					<a href="/bookmarks" class="text-2xl text-[#7A7974]">Bookmarks</a>
				</div>
				<a href="/signin" class="flex items-center gap-6">
					<span class="text-primary font-bold text-2xl">Sign in</span>
					<div class="relative w-12 h-12">
						<div class="bg-[#18171C] w-12 h-12 flex items-center justify-center rounded-full z-20 absolute">
							<MoveRight color="#F7F6F1" />
						</div>
						<div class="bg-[#F7F6F1] w-12 h-12 rounded-full absolute top-0 -left-0.5 z-10" />
						<div class="bg-[#18171C] w-12 h-12 rounded-full absolute top-0 -left-1" />
					</div>
				</a>
			</div>
	  </div>
	</div>
</section>
