<script lang="ts">
    const handleLogin = () => {
        document.dispatchEvent(new CustomEvent('loginAttempt', { detail: { sessionId: '123' } }));
    };
</script>

<section>
	<button on:click={handleLogin}>
        Login
    </button>
</section>
