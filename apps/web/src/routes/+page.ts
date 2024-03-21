// since there's no dynamic data here, we can prerender
// it so that it gets served as a static asset in production
export const prerender = true;

export const load = async ({ parent }) => {
  const { user } = await parent();
  if (user) {
    return {
      user: user
    };
  }
}
