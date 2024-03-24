<script lang="ts">    
  import SignInForm from '$lib/components/SignInForm.svelte';
import { createDialog, melt } from '@melt-ui/svelte';
  import { X } from "lucide-svelte";
  /** Internal helpers */
  import { fade } from 'svelte/transition';

  const {
    elements: {
      trigger,
      overlay,
      content,
      close,
      portalled,
    },
    states: { open },
  } = createDialog({
    forceVisible: true,
  });
</script>

<button
  use:melt={$trigger}
  class="inline-flex items-center justify-center rounded-xl bg-white px-4 py-3
  font-medium leading-none text-magnum-700 shadow hover:opacity-75"
>
  Sign in
</button>

<div class="" use:melt={$portalled}>
  {#if $open || true}
    <div
      use:melt={$overlay}
      class="fixed inset-0 z-50 bg-black/50"
      transition:fade={{ duration: 150 }}
    />
    <div
      class="fixed left-[50%] top-[50%] z-50 max-h-[85vh] w-[90vw]
            max-w-[500px] translate-x-[-50%] translate-y-[-50%] rounded-xl bg-white
            shadow-lg relarive"
      use:melt={$content}
    >
      <button
        use:melt={$close}
        class="p-3 border bg-white hover:bg-stone-50 shadow rounded-lg absolute -top-4 -right-4 transition"
      >
        <X size=20 />
      </button>
      <SignInForm onSubmit={$close} />    
    </div>
  {/if}
</div>
