import { z } from "zod";
import { RequestError, api } from "../api";
import { userSchema } from "../schemas/user";

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
    if (e instanceof Error) {
      throw new RequestError(e);
    }

    throw new RequestError(new Error("unexpected error"));
  }
}
