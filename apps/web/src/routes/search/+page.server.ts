import { message, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { signInRequest, signInRequestSchema, signUpRequest, signUpRequestSchema } from "../../requests/auth";
import { fail, type Actions } from "@sveltejs/kit";
import { RequestError } from "../../api";

export const actions: Actions = {
  signIn: async ({ request, cookies }) => {
    const form = await superValidate(request, zod(signInRequestSchema));

    if (!form.valid) {
      return fail(400, { form });
    }

    try {
      const { sessionId } = await signInRequest(form.data);
     // TODO: setting cookies manually? this doesn't seems right.
      cookies.set('sessionId', sessionId, { path: "/" })

      return { form }
    } catch (error) {
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
  },
  signUp: async ({ request, cookies }) => {
    const form = await superValidate(request, zod(signUpRequestSchema));

    if (!form.valid) {
      return fail(400, { form });
    }

    try {
      const { sessionId } = await signUpRequest(form.data);
      // TODO: same here
      cookies.set("sessionId", sessionId, { path: "/" });

      return { form }
    } catch (error) {
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
