import { superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import { QueryClient } from '@tanstack/svelte-query';
import { browser } from '$app/environment';
import type { LayoutLoad, LayoutLoadEvent } from './$types';
import { signInRequestSchema, signUpRequestSchema } from '../requests/auth';

export const load: LayoutLoad = async ({ data }: LayoutLoadEvent) => {
	const { user } = data;
	const signInForm = await superValidate(zod(signInRequestSchema));
	const signUpForm = await superValidate(zod(signUpRequestSchema));

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser,
				staleTime: Infinity
			}
		}
	});

	return { queryClient, user, signInForm, signUpForm };
};
