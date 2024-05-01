<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import Spinner from "$lib/components/Spinner.svelte";
	import { oauthGoogleExchange } from "../../../requests/user";

    const code = $page.url.searchParams.get('code');

    $effect(() => {
        if (code) {
            // TODO: error handling & remove window reload
            oauthGoogleExchange(code).then(() => goto('/', { invalidateAll: true }).then(() => window.location.reload())).catch(() => {});
        }
    });
</script>

<section class="w-full h-full flex items-center justify-center">
	<Spinner />
</section>