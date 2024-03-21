import { z } from "zod";
import { RequestError, api } from "../api";
import { userSchema } from "../schemas/user";
import { isAxiosError } from "axios";

// TODO: properly add validation rules
export const signInRequestSchema = z.object({
  login: z.string(),
  password: z.string(),
});

export type SignInRequestType = z.infer<typeof signInRequestSchema>;

export const signInRequest = async (payload: SignInRequestType) => {
  try {
    const { data } = await api.post("/login", payload);
    return userSchema.passthrough().parse(data);
  } catch (e) {
    if (isAxiosError(e)) {
      throw new RequestError(e.response);
    }

    throw new RequestError();
  }
}
