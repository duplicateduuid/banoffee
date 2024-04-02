import { QueryClient } from "@tanstack/svelte-query";
import { browser } from "$app/environment";
import type { LayoutLoad, LayoutLoadEvent } from "./$types";

export const load: LayoutLoad = async ({ data }: LayoutLoadEvent) => {
  const { user } = data;
  
  const queryClient = new QueryClient({
  	defaultOptions: {
  		queries: { 
        enabled: browser,
        staleTime: Infinity
      },
  	}
  });

  return { queryClient, user }
}
