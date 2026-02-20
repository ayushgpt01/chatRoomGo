import { Link } from "@tanstack/react-router";

export default function Header() {
	return (
		<div className="navbar bg-base-100 shadow-sm">
			<div className="navbar-start">
				<div className="dropdown">
					<button
						tabIndex={0}
						type="button"
						className="btn btn-ghost btn-circle"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							className="h-5 w-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<title>Breadcrumb Menu</title>{" "}
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth="2"
								d="M4 6h16M4 12h16M4 18h7"
							/>{" "}
						</svg>
					</button>
					<ul
						tabIndex={-1}
						className="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow"
					>
						<li>
							<Link to="/">Homepage</Link>
						</li>
					</ul>
				</div>
			</div>
			<div className="navbar-center">
				<Link to="/" className="btn btn-ghost text-xl">
					Chat Room
				</Link>
			</div>
		</div>
	);
}
