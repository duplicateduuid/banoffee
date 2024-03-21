import axios, { type AxiosResponse } from "axios";

export const api = axios.create({
  // TODO: create a config and get URL from .env
  baseURL: "http://localhost:6969"
})

export class RequestError extends Error {
  constructor(response?: AxiosResponse | undefined) {
    super();

    this.message = response?.data.message
      ? response.data.message
      : "unexpected error";

    Object.setPrototypeOf(this, RequestError.prototype)
  }
}
