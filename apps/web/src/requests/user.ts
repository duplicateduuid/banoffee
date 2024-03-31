import { RequestError, api } from "../api";
import { userSchema } from "../schemas/user";

export const me = async () => {
  try {
    const { data } = await api.get("/me");
    const user = userSchema.passthrough().parse(data.user);

    return { user };
  } catch (error) {
    if (error instanceof Error) {
      throw new RequestError(error);
    }

    throw new RequestError(new Error("unexpected error"));
  }
}
