import { createFileRoute, Link } from "@tanstack/react-router";

export const Route = createFileRoute("/login")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<div className="min-h-screen flex items-center justify-center bg-base-200">
			<div className="card w-full max-w-md bg-base-100 shadow-xl">
				<div className="card-body">
					<h2 className="text-2xl font-bold text-center">Login</h2>

					<input
						type="email"
						placeholder="Email"
						className="input input-bordered w-full mt-4"
					/>
					<input
						type="password"
						placeholder="Password"
						className="input input-bordered w-full mt-2"
					/>

					<button type="button" className="btn btn-primary mt-4">
						Login
					</button>

					<p className="text-center text-sm mt-4">
						Don’t have an account?{" "}
						<Link to="/signup" className="link link-primary">
							Sign up
						</Link>
					</p>
				</div>
			</div>
		</div>
	);
}
