// TODO: use the generated type instead of manually defining one.
import { z } from "zod";

export const userSchema = z.object({
  id: z.string(),
  email: z.string().email(),
  username: z.string(),
  avatar_url: z.string().nullable(),
  header_url: z.string().nullable(),
  bio: z.string().nullable(),
  created_at: z.string().datetime()
});

export const resourceSchema = z.object({
  id: z.string(),
  url: z.string(),
  name: z.string(),
  image_url: z.string().nullish(),
  author: z.string().nullish(),
  description: z.string().nullish(),
  status: z.string().nullish(),
  review_rating: z.string().nullish(),
  review_comment: z.string().nullish(),
  created_at: z.date().or(z.string()).nullish()
}).passthrough()

export type User = z.infer<typeof userSchema>;

export type Resource = z.infer<typeof resourceSchema>
