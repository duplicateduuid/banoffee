import { superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import { z } from 'zod';

// Define outside the load function so the adapter can be cached
const schema = z.object({
  login: z.string().email(),
  password: z.string().min(8).max(255)
})

export const load = (async () => {
  const form = await superValidate(zod(schema));

  // Always return { form } in load functions
  return { form };
});
