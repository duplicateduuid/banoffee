import type { PageLoad } from "./$types";

export const load = (async () => {
  const { user } = { user: "hehe"};

  if (user) {
    return {
      user: user
    };
  }
}) satisfies PageLoad
