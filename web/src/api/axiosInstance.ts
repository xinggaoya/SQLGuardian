import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

const axiosInstance: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_BASE_AXIOS_URL as string,
    timeout: 5000, // 请求超时时间（单位：毫秒）
});
// 请求拦截器
axiosInstance.interceptors.request.use(
    (config: AxiosRequestConfig) => {
        // 在发送请求之前做一些处理，例如添加请求头
        // config.headers['Authorization'] = 'Bearer ' + getToken();
        return config;
    },
    (error: any) => {
        return Promise.reject(error);
    }
);

// 响应拦截器
axiosInstance.interceptors.response.use(
    (response: AxiosResponse) => {
        // 在响应之前做一些处理
        return response.data;
    },
    (error: any) => {
        return Promise.reject(error);
    }
);

export default axiosInstance;
