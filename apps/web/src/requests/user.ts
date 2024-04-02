import type { Cookies } from "@sveltejs/kit";
import { RequestError, api } from "../api";
import { userSchema } from "../schemas/user";

export const me = async (cookies?: Cookies) => {
  const sessionId = cookies?.get("sessionId");
  
  const config = sessionId
    ? {
      headers: { Cookie: `sessionId=${sessionId}`}
    }
    : undefined;  

  try {
    const { data } = await api.get("/me", config);
    return userSchema.passthrough().parse(data.user);
  } catch (error) {
    if (error instanceof Error) {
      throw new RequestError(error);
    }

    throw new RequestError(new Error("unexpected error"));
  }
}
