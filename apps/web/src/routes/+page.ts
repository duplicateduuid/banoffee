import { superValidate } from "sveltekit-superforms";
import type { PageLoad } from "./$types";
import { zod } from "sveltekit-superforms/adapters";
import { signInRequestSchema } from "../requests/auth";

export const load: PageLoad = (async () => {
  const signInForm = await superValidate(zod(signInRequestSchema));
  const signUpForm = await superValidate(zod(signInRequestSchema));

  return { signInForm, signUpForm };
});
