import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authService } from "@/services/authService";
import type { LoginCredentials, SignupCredentials, User } from "@/types/auth";
import { getErrorMessage } from "@/utils/errorHandler";

export interface AuthState {
	user: User | null;
	isAuthenticated: boolean;
	isAuthenticating: boolean;
	isCreating: boolean;
	error: string | null;
	login: (credentials: LoginCredentials) => Promise<void>;
	logout: () => void;
	signup: (credentials: SignupCredentials) => Promise<void>;
	checkAuth: () => Promise<void>;
}

const useAuthStore = create<AuthState>()(
	persist(
		(set, get) => ({
			user: null,
			isAuthenticated: false,
			isAuthenticating: false,
			isCreating: false,
			error: null,

			login: async (credentials) => {
				set({ isAuthenticating: true, error: null });
				try {
					const { user, token, refreshToken } =
						await authService.login(credentials);

					localStorage.setItem("token", token);
					localStorage.setItem("refresh_token", refreshToken);
					set({ user, isAuthenticated: true, isAuthenticating: false });
				} catch (err) {
					set({ error: getErrorMessage(err), isAuthenticating: false });
					throw err;
				}
			},

			signup: async (credentials) => {
				set({ isCreating: true, error: null });
				try {
					const { user, token, refreshToken } =
						await authService.signup(credentials);

					localStorage.setItem("token", token);
					localStorage.setItem("refresh_token", refreshToken);

					set({ user, isAuthenticated: true, isCreating: false });
				} catch (e) {
					set({ error: getErrorMessage(e), isCreating: false });
					throw e;
				}
			},

			logout: async () => {
				const refreshToken = localStorage.getItem("refresh_token");
				if (refreshToken) {
					await authService.logout(refreshToken);
				}

				localStorage.removeItem("token");
				localStorage.removeItem("refresh_token");
				set({ user: null, isAuthenticated: false, error: null });
			},

			checkAuth: async () => {
				if (!localStorage.getItem("token")) return;
				try {
					const user = await authService.getCurrentUser();
					set({ user, isAuthenticated: true });
				} catch {
					get().logout();
				}
			},
		}),
		{
			name: "auth-storage",
			partialize: (state) => ({
				user: state.user,
				isAuthenticated: state.isAuthenticated,
			}),
		},
	),
);

export default useAuthStore;
