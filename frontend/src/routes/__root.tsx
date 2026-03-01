import { TanStackDevtools } from "@tanstack/react-devtools";
import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import NotFound from "@/components/NotFound";
import ToastContainer from "@/components/ToastContainer";
import TanStackQueryDevtools from "@/integrations/tanstack-query/devtools";
import type { AuthState } from "@/stores/authStore";
import useAuthStore from "@/stores/authStore";

interface MyRouterContext {
	queryClient: QueryClient;
	auth: AuthState;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
	loader: async () => {
		const authStore = useAuthStore.getState();

		if (!authStore.initialized) {
			await authStore.checkAuth();
		}

		return null;
	},
	component: () => (
		<>
			<Outlet />
			<ToastContainer />
			<TanStackDevtools
				config={{
					position: "bottom-right",
				}}
				plugins={[
					{
						name: "Tanstack Router",
						render: <TanStackRouterDevtoolsPanel />,
					},
					TanStackQueryDevtools,
				]}
			/>
		</>
	),
	notFoundComponent: NotFound,
});
