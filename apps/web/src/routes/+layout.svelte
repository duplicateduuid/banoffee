<script lang="ts">
	import { QueryClientProvider } from "@tanstack/svelte-query";
	import '../app.css';
	import type { LayoutData } from "./$types";

	export let data: LayoutData;
</script>

<main class="w-full min-h-screen font-secondary bg-[#F7F6F1] text-[#18171C]">
	<QueryClientProvider client={data.queryClient}>
		<slot />
	</QueryClientProvider>
</main>
