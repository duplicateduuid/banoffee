<script>
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
	import '../app.css';

	const queryClient = new QueryClient();
</script>

<main class="w-full min-h-screen font-secondary bg-[#F7F6F1] text-[#18171C]">
	<QueryClientProvider client={queryClient}>
		<slot />
	</QueryClientProvider>
</main>
