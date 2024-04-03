import { z } from 'zod';
import { RequestError, api } from '../api';
import { resourceSchema } from '../schemas/resource';

export const getResourceByURL = async (url: string) => {
	try {
		const { data } = await api.get(`/resource?url=${url}`);

		return resourceSchema.parse(data.resource);
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
};

export const searchResourcesByName = async (name: string, offset: number) => {
	try {
		console.log(offset);

		const { data } = await api.get(`/resource/search?name=${name}&limit=10&offset=${offset}`);

		return z.array(resourceSchema).parse(data.resources);
	} catch (error) {
		if (error instanceof Error) {
			throw new RequestError(error);
		}

		throw new RequestError(new Error('unexpected error'));
	}
};
