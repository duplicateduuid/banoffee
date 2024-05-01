import { z } from 'zod';
import type { Cookies } from '@sveltejs/kit';
import { RequestError, api } from '../api';
import { userSchema } from '../schemas/user';
import { resourceSchema } from '../schemas/resource';

export const me = async (cookies?: Cookies) => {
	const sessionId = cookies?.get('sessionId');

	const config = sessionId
		? {
				headers: { Cookie: `sessionId=${sessionId}` }
			}
		: undefined;

	try {
		const { data } = await api.get('/me', config);
		return userSchema.passthrough().parse(data.user);
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
};

export const getRecommendations = async () => {
	try {
		const { data } = await api.get('/recommendations');

		return z.array(resourceSchema).parse(data.recommendations);
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
};

export const getPopularThisWeek = async () => {
	try {
		const { data } = await api.get('/popular');

		return z.array(resourceSchema).parse(data.resources);
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
};

export const getMyResources = async (limit: number, offset: number, status: string) => {
	try {
		const { data } = await api.get(
			`/user/resources?limit=${limit}&offset=${offset}&status=${status}`
		);

		return z.array(resourceSchema).parse(data.resources);
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
};

export const oauthGoogle = async (): Promise<string> => {
	try {
		const { data } = await api.get('/oauth/google');

		return data.url;
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
}

export const oauthGoogleExchange = async (code: string) => {
	try {
		const { data } = await api.post(`/oauth/google/exchange?code=${code}`);

		return data;
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
}
