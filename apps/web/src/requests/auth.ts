import { z } from "zod";
import { RequestError, api } from "../api";
import { userSchema } from "../schemas/user";

export const signInRequestSchema = z.object({
  login: z.string()
    .min(5, "Username or email must have at least 5 characters")
    .max(20, "Username or email can't have more than 20 characters"),
  password: z.string()
    .min(8, "Password must have at least 8 characters")
    .max(255, "Password can't have more than 255 characters"),
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
