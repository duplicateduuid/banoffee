<script lang="ts">
	import { getCookies } from '../../utils';

	$effect(() => {
		const cookies = getCookies(document.cookie);

		const sessionId = cookies['sessionId'];

		if (sessionId)
			document.dispatchEvent(
				new CustomEvent('loginAttempt', { detail: { sessionId, expiration: 3600 * 24 } })
			);
	});
</script>

<!-- TODO: beautify this page -->
<section class="w-full h-full flex items-center justify-center">
	<p>You can close this tab now</p>
</section>
