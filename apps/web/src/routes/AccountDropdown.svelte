<script lang="ts">
  import { UserRound } from "lucide-svelte";
  import { type User } from "../schemas/user";
	import { createDropdownMenu, melt } from "@melt-ui/svelte";
	import { writable } from "svelte/store";
	import { fade, fly } from "svelte/transition";

  type Props = {
    user: User
  }
  const { user }: Props = $props();

  const {
    elements: { menu, item, trigger },
    states: { open }
  } = createDropdownMenu({ forceVisible: true, loop: true, positioning: { placement: "bottom-end" } })
</script>

<button type="button" class="trigger" use:melt={$trigger}>
  <UserRound />
</button>

{#if $open}
  <div class="menu" use:melt={$menu} transition:fly={{duration: 100, y: -20}}>
    <div class="user-info">
      <div class="avatar">
        <UserRound size={32} />
      </div>
      <p>{user.username}</p>
    </div>

    <div class="items">
      <button
        type="button"
        class="item"
        use:melt={$item}
      >
        Bookmarks
      </button>
      <button
        type="button"
        class="item"
        use:melt={$item}
      >
        Settings
      </button>
      <hr>
      <button
        type="button"
        class="item"
        use:melt={$item}
      >
        Sign out
      </button>
    </div>
  </div>
{/if}
<style>
  .trigger {
    @apply bg-primary-400 rounded-full p-2.5 text-stone-800;
    @apply hover:bg-primary-300 transition;
  }

  .menu {
    @apply mt-4 min-w-[300px] p-6;
    @apply bg-white rounded-2xl outline-none;
    @apply flex flex-col gap-4;
    @apply border border-stone-100 shadow-lg;
  }

  .user-info {
    @apply flex flex-col items-center gap-2.5 font-semibold;
    @apply text-xl w-full;
  }

  .avatar {
    @apply p-4 rounded-full bg-primary-400;
  }

  .items {
    @apply flex flex-col items-start w-full;
  }  

  .item {
    @apply text-lg text-stone-700 tracking-wide py-1.5 outline-none;
    @apply text-stone-700 hover:text-stone-400 focus:text-stone-400;
    @apply transition duration-75;
  }

  hr {
    @apply my-2 w-full;
  }
</style>
