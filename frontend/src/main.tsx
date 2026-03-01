import { RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import TanStackQueryProvider, {
	queryClient,
} from "./integrations/tanstack-query/root-provider";
import { router } from "./router";
import useAuthStore from "./stores/authStore";

import "./styles.css";

const App = () => {
	const auth = useAuthStore();

	return (
		<TanStackQueryProvider>
			<RouterProvider router={router} context={{ queryClient, auth }} />
		</TanStackQueryProvider>
	);
};

const rootElement = document.getElementById("root");
if (rootElement && !rootElement.innerHTML) {
	const root = ReactDOM.createRoot(rootElement);
	root.render(
		<StrictMode>
			<App />
		</StrictMode>,
	);
}
