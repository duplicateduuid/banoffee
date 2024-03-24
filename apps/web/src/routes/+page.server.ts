import { message, superValidate } from "sveltekit-superforms";
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
    } catch (error) {
      console.log(error)      
      const isRequestError = error instanceof RequestError;

      if (!isRequestError) {
        return fail(500, { error });
      }

      if (error.status === 403) {
        return message(
          form,
          "We couldn't find an account that matches the login and password you entered."
        )
      }

      return message(
        form,
        "Something went wrong signing you in."
      );
    }
  }
} 
