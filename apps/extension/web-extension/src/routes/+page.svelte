<script lang="ts">
  let isLoggedIn: boolean = $state(false);

  // TODO: test it on firefox
  $effect(() => {
    if (typeof chrome !== "undefined") {
      chrome.storage.local.get("sessionId", (data) => {
        if (data.sessionId) {
          isLoggedIn = true;
        }
      });
    } else if (typeof browser !== "undefined") {
      browser.storage.local.get("sessionId", (data) => {
        if (data.sessionId) {
          isLoggedIn = true;
        }
      });
    }
  });
</script>

<div>
  {#if isLoggedIn}
    <p class="text-red-500">isLoggedIn</p>
  {:else}
    <p class="text-red-500">isNotLoggedIn</p>
  {/if}
</div>

<style>
</style>
