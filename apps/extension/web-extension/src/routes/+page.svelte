<script lang="ts">
  import { api } from "$lib";
  import type { Resource } from "../types";
  import { getBrowserEnv } from "../utils/getBrowserEnv";
  import icon from "../../static/icon_128px.png";

  let browserEnv = getBrowserEnv();
  let loading: boolean = $state(true);
  let url: string | null = $state(null);
  let title: string | null = $state(null);
  let resourceId: string | null = $state(null);
  let userHolds: boolean = $state(false);

  const handleSaveResource = async () => {
    try {
      if (resourceId) {
        // TODO: handle case of the user already has the resource saved
        if (userHolds) return;

        await api.post(`/user/resource/${resourceId}`, { status: "bookmarked" });
        return;
      }

      const {
        data: { resource: newResource },
      } = await api.post<{ resource: Resource }>("/resource", {
        url,
        name: title,
      });

      await api.post(`/user/resource/${newResource.id}`, { status: "bookmarked" });
    } catch (err) {
      throw new Error(
        `Unexpected error calling the API: ${JSON.stringify(err)}`
      );
    }
  };

  // TODO: test it on firefox
  $effect(() => {
    if (browserEnv) {
      browserEnv.storage.local.get("session", (data) => {
        if (data.session) {
          browserEnv!.tabs.query({ active: true }, async (tabs) => {
            const currentTab = tabs?.[0];

            if (currentTab.url) {
              api
                .get<{
                  resource: Resource | null;
                  user_holds?: boolean;
                }>(`/user/resource?url=${currentTab.url}`)
                .then((res) => {
                  const {
                    data: { resource, user_holds },
                  } = res;

                  if (resource) {
                    url = resource.url;
                    title = resource.name;
                    resourceId = resource.id;
                    userHolds = user_holds || false;

                    loading = false;
                    return;
                  }

                  url = currentTab.url!;
                  title = currentTab.title || null;

                  loading = false;
                })
                .catch((err) => {
                  throw new Error(
                    `Unexpected error calling the API: ${JSON.stringify(err)}`
                  );
                });
            }
          });
        } else {
          // TODO: replace mocked URL
          browserEnv!.tabs.create({
            url: `${import.meta.env.VITE_WEB_PAGE_URL}/extension-login`,
          });
        }
      });
    }
  });
</script>

<div class="w-full h-full">
  {#if loading}
    <div class="min-w-[30em] min-h-[20em]">
      <div id="load">
        <div>G</div>
        <div>N</div>
        <div>I</div>
        <div>D</div>
        <div>A</div>
        <div>O</div>
        <div>L</div>
      </div>
    </div>
  {/if}

  {#if !loading}
    <div class="min-w-[30em] min-h-[20em] flex flex-col justify-between p-8">
      <div class="w-full h-full flex items-center gap-2">
        <img src={icon} alt="Banoffee Icon" class="rounded-full w-8 h-8" />

        <p class="text-lg font-semibold">Banoffee</p>
      </div>

      <div
        class="max-w-full flex flex-col gap-[2px] text-ellipsis whitespace-nowrap overflow-hidden"
      >
        <p class="text-sm font-normal">Title:</p>
        <p
          class="text-base text-ellipsis whitespace-nowrap overflow-hidden font-medium"
        >
          {title || "Unknown title"}
        </p>
      </div>

      <button
        class="w-full h-12 rounded-md border-none bg-[#4e473b] text-[#F7F6F1] text-base inline-block"
        on:click={handleSaveResource}
      >
        Bookmark
      </button>
    </div>
  {/if}
</div>

<!-- TODO: maybe rewrite it using tailwind (I'm lazy, sorry) -->
<style>
  #load {
    position: absolute;
    width: 600px;
    height: 36px;
    left: 50%;
    top: 45%;
    margin-left: -300px;
    overflow: visible;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    cursor: default;
  }

  #load div {
    position: absolute;
    width: 20px;
    height: 36px;
    opacity: 0;
    font-family: Helvetica, Arial, sans-serif;
    animation: move 2.5s linear infinite;
    -o-animation: move 2.5s linear infinite;
    -moz-animation: move 2.5s linear infinite;
    -webkit-animation: move 2.5s linear infinite;
    transform: rotate(180deg);
    -o-transform: rotate(180deg);
    -moz-transform: rotate(180deg);
    -webkit-transform: rotate(180deg);
    color: #4e473b;
  }

  #load div:nth-child(2) {
    animation-delay: 0.2s;
    -o-animation-delay: 0.2s;
    -moz-animation-delay: 0.2s;
    -webkit-animation-delay: 0.2s;
  }
  #load div:nth-child(3) {
    animation-delay: 0.4s;
    -o-animation-delay: 0.4s;
    -webkit-animation-delay: 0.4s;
    -webkit-animation-delay: 0.4s;
  }
  #load div:nth-child(4) {
    animation-delay: 0.6s;
    -o-animation-delay: 0.6s;
    -moz-animation-delay: 0.6s;
    -webkit-animation-delay: 0.6s;
  }
  #load div:nth-child(5) {
    animation-delay: 0.8s;
    -o-animation-delay: 0.8s;
    -moz-animation-delay: 0.8s;
    -webkit-animation-delay: 0.8s;
  }
  #load div:nth-child(6) {
    animation-delay: 1s;
    -o-animation-delay: 1s;
    -moz-animation-delay: 1s;
    -webkit-animation-delay: 1s;
  }
  #load div:nth-child(7) {
    animation-delay: 1.2s;
    -o-animation-delay: 1.2s;
    -moz-animation-delay: 1.2s;
    -webkit-animation-delay: 1.2s;
  }

  @keyframes move {
    0% {
      left: 0;
      opacity: 0;
    }
    35% {
      left: 41%;
      -moz-transform: rotate(0deg);
      -webkit-transform: rotate(0deg);
      -o-transform: rotate(0deg);
      transform: rotate(0deg);
      opacity: 1;
    }
    65% {
      left: 59%;
      -moz-transform: rotate(0deg);
      -webkit-transform: rotate(0deg);
      -o-transform: rotate(0deg);
      transform: rotate(0deg);
      opacity: 1;
    }
    100% {
      left: 100%;
      -moz-transform: rotate(-180deg);
      -webkit-transform: rotate(-180deg);
      -o-transform: rotate(-180deg);
      transform: rotate(-180deg);
      opacity: 0;
    }
  }

  @-moz-keyframes move {
    0% {
      left: 0;
      opacity: 0;
    }
    35% {
      left: 41%;
      -moz-transform: rotate(0deg);
      transform: rotate(0deg);
      opacity: 1;
    }
    65% {
      left: 59%;
      -moz-transform: rotate(0deg);
      transform: rotate(0deg);
      opacity: 1;
    }
    100% {
      left: 100%;
      -moz-transform: rotate(-180deg);
      transform: rotate(-180deg);
      opacity: 0;
    }
  }

  @-webkit-keyframes move {
    0% {
      left: 0;
      opacity: 0;
    }
    35% {
      left: 41%;
      -webkit-transform: rotate(0deg);
      transform: rotate(0deg);
      opacity: 1;
    }
    65% {
      left: 59%;
      -webkit-transform: rotate(0deg);
      transform: rotate(0deg);
      opacity: 1;
    }
    100% {
      left: 100%;
      -webkit-transform: rotate(-180deg);
      transform: rotate(-180deg);
      opacity: 0;
    }
  }

  @-o-keyframes move {
    0% {
      left: 0;
      opacity: 0;
    }
  }
</style>
