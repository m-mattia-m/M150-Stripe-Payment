import { AxiosRequestConfig } from 'axios';
import { Configuration } from '../api';

export const generateAxiosConfig = (): Configuration => {
    if (!import.meta.env.VITE_API_URL) {
        throw new Error('Axios Config could not be generated!');
    }
    const basePath = import.meta.env.VITE_API_URL;

    return new Configuration({
        basePath,
    });
};
