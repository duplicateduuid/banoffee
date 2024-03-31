import { superValidate } from "sveltekit-superforms";
import type { PageLoad } from "./$types";
import { zod } from "sveltekit-superforms/adapters";
import { signInRequestSchema } from "../requests/auth";
import { me } from "../requests/user";

export const load: PageLoad = async ({ parent }) => {
  const { queryClient } = await parent();

  const signInForm = await superValidate(zod(signInRequestSchema));
  const signUpForm = await superValidate(zod(signInRequestSchema));

  await queryClient.prefetchQuery({
    queryKey: ['me'],
    queryFn: me
  })

  return { signInForm, signUpForm };
};
