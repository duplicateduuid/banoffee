import axios, { type AxiosResponse } from "axios";

export const api = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL,
  withCredentials: true,
});

export class RequestError extends Error {
  constructor(response?: AxiosResponse | undefined) {
    super();

    this.message = response?.data.message
      ? response.data.message
      : "unexpected error";

    Object.setPrototypeOf(this, RequestError.prototype)
  }
}
