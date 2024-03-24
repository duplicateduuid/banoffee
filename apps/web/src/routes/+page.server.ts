import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { signInRequest, signInRequestSchema } from "../requests/auth";
import { fail, redirect, type Actions } from "@sveltejs/kit";
import { RequestError } from "../api";

export const actions: Actions = {
  default: async ({ request }) => {
    const form = await superValidate(request, zod(signInRequestSchema));

    if (!form.valid) {
      return fail(400, { form });
    }

    try {
      const user = await signInRequest(form.data);

      throw redirect(303, "/");      
    } catch (e) {
      console.log(e)
      if (e instanceof RequestError) {
        return fail(e.status, { message: e.message });
      }

      return fail(500, { error: e });
    }
  }
} 
