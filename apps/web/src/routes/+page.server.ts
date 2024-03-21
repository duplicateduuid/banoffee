import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { signInRequest, signInRequestSchema } from "../requests/auth";
import { fail, redirect, type Actions } from "@sveltejs/kit";

export const actions: Actions = {
  "sign-in": async ({ request, cookies }) => {
    const form = await superValidate(request, zod(signInRequestSchema));

    if (!form.valid) {
      return fail(400, { form });
    }

    const user = await signInRequest(form.data);

    cookies.set('sessionId', user.id, {
      path: "/",
    });

    throw redirect(303, "/");
  }
}

export const load = (async () => {
  const form = await superValidate(zod(signInRequestSchema));

  return { form };
});
