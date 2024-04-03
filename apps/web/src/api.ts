import axios, { AxiosError, isAxiosError, type AxiosResponse } from 'axios';
import { LucideChrome } from 'lucide-svelte';

export const api = axios.create({
	baseURL: import.meta.env.VITE_BACKEND_URL,
	withCredentials: true
});

type Session = {
	id: string;
	expiresAt: string;
};

export class RequestError extends Error {
	detail?: string;
	status: number;

	constructor(error: Error) {
		super();

		this.message = error.message;

		if (isAxiosError(error) && error.response) {
			this.detail = error.response.data;
			this.status = error.response.status;

			Object.setPrototypeOf(this, RequestError.prototype);
			return;
		}

		this.status = 500;

		Object.setPrototypeOf(this, RequestError.prototype);
	}
}
