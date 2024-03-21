<script lang="ts">
	import { api } from "../../api";
	import { getCookies } from "../../utils";

    // TODO: remove mocked values
    const handleLogin = async () => {
        await api.post('/login', {
            email: "email@test.com",
            password: "test1234"
        });

        const cookies = getCookies(document.cookie);

        const sessionId = cookies['sessionId'];

        document.dispatchEvent(new CustomEvent('loginAttempt', { detail: { sessionId, expiration: 3600*24 } }));
    };
</script>

<section>
	<button on:click={handleLogin}>
        Login
    </button>
</section>
