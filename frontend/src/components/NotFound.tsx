import { Link } from "@tanstack/react-router";

export default function NotFound() {
	return (
		<div className="min-h-screen flex items-center justify-center bg-base-200">
			<div className="card w-full max-w-md bg-base-100 shadow-xl text-center">
				<div className="card-body">
					<h1 className="text-6xl font-bold text-primary">404</h1>
					<p className="text-lg mt-2">Page not found</p>

					<div className="mt-6">
						<Link to="/" className="btn btn-primary">
							Go Home
						</Link>
					</div>
				</div>
			</div>
		</div>
	);
}
