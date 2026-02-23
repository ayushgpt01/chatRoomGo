import axiosClient from "@/integrations/axios/axiosClient";
import {
	type LoginCredentials,
	type SignupCredentials,
	type User,
	UserSchema,
} from "@/types/auth";

// Define a specific interface for the API response
interface LoginResponse {
	user: User;
	token: string;
	refreshToken: string;
}

export const authService = {
	login: async (credentials: LoginCredentials): Promise<LoginResponse> => {
		// The Axios interceptor handles the .data unwrapping if configured,
		// but here we ensure we get the right type.
		const response = await axiosClient.post<LoginResponse>(
			"/auth/login",
			credentials,
		);
		const data = response.data;

		// Runtime validation: Ensure the user object matches our schema
		UserSchema.parse(data.user);

		return data;
	},

	logout: async (refreshToken: string): Promise<void> => {
		await axiosClient.post<LoginResponse>("/auth/logout", { refreshToken });
	},

	getCurrentUser: async (): Promise<User> => {
		const response = await axiosClient.get<User>("/auth/me");
		return UserSchema.parse(response.data);
	},

	signup: async (credentials: SignupCredentials): Promise<LoginResponse> => {
		const { confirmPassword: _, ...payload } = credentials;
		const response = await axiosClient.post<LoginResponse>(
			"/auth/signup",
			payload,
		);
		const data = response.data;

		// Runtime validation: Ensure the user object matches our schema
		UserSchema.parse(data.user);

		return data;
	},
};
