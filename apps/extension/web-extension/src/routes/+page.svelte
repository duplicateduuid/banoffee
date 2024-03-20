<script lang="ts">
  let isLoggedIn: boolean = $state(false);

  // TODO: make it work for firefox
  $effect(() => {
    if (typeof chrome !== "undefined") {
      chrome.storage.local.get("sessionId", (data) => {
        if (data.sessionId) {
          isLoggedIn = true;
        }
      });
    }
  });
</script>

<div>
  {#if isLoggedIn}
    <p>isLoggedIn</p>
  {:else}
    <p>isNotLoggedIn</p>
  {/if}
</div>

<style>
</style>
