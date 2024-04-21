import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent }) => {
	const { user } = await parent();

	// TODO: pass some param here to open the login modal in this case
	if (!user) redirect(302, '/');

	return { user };
};
