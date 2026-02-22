import { createFileRoute, Link } from "@tanstack/react-router";
import { ShieldAlert } from "lucide-react";

export const Route = createFileRoute("/unauthorized")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<div className="min-h-screen flex items-center justify-center bg-base-200">
			<div className="card w-full max-w-md bg-base-100 shadow-xl text-center">
				<div className="card-body">
					<div className="flex justify-center">
						<ShieldAlert className="w-12 h-12 text-error" />
					</div>

					<h1 className="text-3xl font-bold mt-4">Access Denied</h1>

					<p className="text-base-content/70 mt-2">
						You do not have permission to access this resource.
					</p>

					<div className="flex flex-col gap-2 mt-6">
						<Link to="/" className="btn btn-primary w-full">
							Go Home
						</Link>
						<Link to="/login" className="btn btn-outline w-full">
							Login
						</Link>
					</div>
				</div>
			</div>
		</div>
	);
}
