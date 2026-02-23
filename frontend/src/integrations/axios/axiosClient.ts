import axios from "axios";
import { router } from "@/router";
import useAuthStore from "@/stores/authStore";

// 1. Determine the Base URL
// This logic checks if we are in dev; if so, it uses the dev URL,
// otherwise it defaults to the production string from your env.
const baseURL = `${import.meta.env.VITE_API_URL}/api`;

const axiosClient = axios.create({
	baseURL: baseURL,
	headers: {
		"Content-Type": "application/json",
		Accept: "application/json",
	},
	// Setting a timeout is good practice to prevent hanging requests
	timeout: 10000,
});

// 2. Request Interceptor
// Perfect for injecting Authorization headers (like JWT tokens) dynamically
axiosClient.interceptors.request.use(
	(config) => {
		const token = localStorage.getItem("token");
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	},
);

let isRefreshing = false;
let failedQueue: Array<{
	resolve: (value: unknown) => void;
	reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: unknown, token: string | null = null) => {
	failedQueue.forEach((prom) => {
		if (error) prom.reject(error);
		else prom.resolve(token);
	});
	failedQueue = [];
};

axiosClient.interceptors.response.use(
	(response) => response,
	async (error) => {
		const originalRequest = error.config;

		// 401 Unauthorized: Trigger Token Refresh
		if (error.response?.status === 401 && !originalRequest._retry) {
			if (isRefreshing) {
				return new Promise((resolve, reject) => {
					failedQueue.push({ resolve, reject });
				})
					.then((token) => {
						originalRequest.headers.Authorization = `Bearer ${token}`;
						return axiosClient(originalRequest);
					})
					.catch((err) => Promise.reject(err));
			}

			const refreshToken = localStorage.getItem("refresh_token");

			if (!refreshToken) {
				useAuthStore.getState().logout();
				return Promise.reject(error);
			}

			originalRequest._retry = true;
			isRefreshing = true;

			try {
				// Call your refresh endpoint
				const { data } = await axios.post(`${baseURL}/auth/refresh`, {
					refreshToken,
				});

				const newToken = data.token;
				localStorage.setItem("token", newToken);

				// Update the original request and the default header
				axiosClient.defaults.headers.common.Authorization = `Bearer ${newToken}`;
				originalRequest.headers.Authorization = `Bearer ${newToken}`;

				processQueue(null, newToken);
				return axiosClient(originalRequest);
			} catch (refreshError) {
				processQueue(refreshError, null);
				// If refresh fails, clear store and redirect
				useAuthStore.getState().logout();
				return Promise.reject(refreshError);
			} finally {
				isRefreshing = false;
			}
		}

		// 403 Forbidden: Immediate Logout
		if (error.response?.status === 403) {
			router.navigate({ to: "/unauthorized" });
		}

		return Promise.reject(error);
	},
);

export default axiosClient;
