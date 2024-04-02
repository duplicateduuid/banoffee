import { superValidate } from "sveltekit-superforms";
import type { PageLoad } from "./$types";
import { zod } from "sveltekit-superforms/adapters";
import { signInRequestSchema, signUpRequestSchema } from "../requests/auth";

export const load: PageLoad = async ({ parent }) => {
  const { user }= await parent();
  const signInForm = await superValidate(zod(signInRequestSchema));
  const signUpForm = await superValidate(zod(signUpRequestSchema));
  
  return { user, signInForm, signUpForm };
};
