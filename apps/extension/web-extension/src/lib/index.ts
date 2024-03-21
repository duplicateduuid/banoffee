import axios from "axios";

// place files you want to import through the `$lib` alias in this folder.
export const api = axios.create({
    baseURL: import.meta.env.VITE_BACKEND_URL,
    withCredentials: true,
});