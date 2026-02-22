import { createRouter as createTanStackRouter } from "@tanstack/react-router";
import { queryClient } from "./integrations/tanstack-query/root-provider";
import { routeTree } from "./routeTree.gen";
import type { AuthState } from "./stores/authStore";

export const router = createTanStackRouter({
	routeTree,
	context: {
		queryClient: queryClient,
		auth: {} as AuthState,
	},
	scrollRestoration: true,
	defaultPreload: "intent",
	defaultPreloadStaleTime: 0,
});

declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}
