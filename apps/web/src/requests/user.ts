import { z } from "zod";
import { RequestError, api } from "../api";
import { resourceSchema, userSchema } from "../schemas/user";

export const me = async () => {
  try {
    const { data } = await api.get("/me");
    return userSchema.passthrough().parse(data.user);
  } catch (error) {
    if (error instanceof Error) {
      throw new RequestError(error);
    }

    throw new RequestError(new Error("unexpected error"));
  }
}

export const getRecommendations = async () => {
  try {
    const { data } = await api.get("/recommendations");

    return z.array(resourceSchema).parse(data.recommendations);
  } catch (error) {
    if (error instanceof Error) {
      throw new RequestError(error);
    }

    throw new RequestError(new Error("unexpected error"));
  }
}