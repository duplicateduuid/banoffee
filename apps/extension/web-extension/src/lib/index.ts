import axios from "axios";
import { getBrowserEnv } from "../utils/getBrowserEnv";

// place files you want to import through the `$lib` alias in this folder.
export const api = axios.create({
    baseURL: import.meta.env.VITE_BACKEND_URL,
    withCredentials: true,
});

api.interceptors.request.use(config => {
    const browserEnv = getBrowserEnv();
    
    if (browserEnv) {
        browserEnv.storage.local.get("session", (data) => {
            if (data.session) {
                config.headers['Cookie'] = `sessionId=${data.session.id}; Max-Age=${data.session.expiration}; Domain=localhost; Path=/`;
            }
        });
    }

    return config;
}, error => {
    return Promise.reject(error);
});