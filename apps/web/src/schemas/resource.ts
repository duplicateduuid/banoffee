import { z } from 'zod';

export const resourceSchema = z
	.object({
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
	})
	.passthrough();

export type Resource = z.infer<typeof resourceSchema>;
