<script lang="ts">
  let count = $state(0);
</script>

<div>
  <button on:click={() => count++}>
    clicks: {count}
  </button>
</div>

<style>
</style>
