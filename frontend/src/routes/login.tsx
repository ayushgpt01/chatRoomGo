import { zodResolver } from "@hookform/resolvers/zod";
import {
	createFileRoute,
	Link,
	redirect,
	useNavigate,
} from "@tanstack/react-router";
import { type SubmitHandler, useForm } from "react-hook-form";
import useAuthStore from "@/stores/authStore";
import useToastStore from "@/stores/toastStore";
import { type LoginCredentials, LoginSchema } from "@/types/auth";

export const Route = createFileRoute("/login")({
	beforeLoad: ({ context }) => {
		console.log(context.auth.isAuthenticated, "context.auth.isAuthenticated");
		if (context.auth.isAuthenticated) {
			throw redirect({
				to: "/",
			});
		}
	},
	component: RouteComponent,
});

function RouteComponent() {
	const login = useAuthStore((state) => state.login);
	const isAuthenticating = useAuthStore((state) => state.isAuthenticating);
	const navigate = useNavigate();
	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm({
		resolver: zodResolver(LoginSchema),
	});
	const showToast = useToastStore((s) => s.show);

	const onSubmit: SubmitHandler<LoginCredentials> = async (data) => {
		try {
			await login(data);
			showToast("Logged in successfully", "success");
			navigate({ to: "/" });
		} catch (err) {
			showToast("Invalid credentials", "error");
			console.error(err);
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-base-200">
			<div className="card w-full max-w-md bg-base-100 shadow-xl">
				<form onSubmit={handleSubmit(onSubmit)} className="card-body space-y-2">
					<h2 className="text-2xl font-bold text-center">Login</h2>

					<label className="form-control w-full mt-4">
						<div className="label">
							<span className="label-text">Username</span>
						</div>
						<input
							type="text"
							placeholder="Username"
							className={`input input-bordered w-full mt-1 ${
								errors.username ? "input-error" : ""
							}`}
							{...register("username")}
						/>
					</label>
					{errors.username && (
						<p className="text-error text-sm mt-1">
							{errors.username.message || "Please enter a valid username"}
						</p>
					)}

					<label className="form-control w-full mt-2">
						<div className="label">
							<span className="label-text">Password</span>
						</div>
						<input
							type="password"
							placeholder="Password"
							className={`input input-bordered w-full mt-1 ${
								errors.password ? "input-error" : ""
							}`}
							{...register("password")}
						/>
					</label>
					{errors.password && (
						<p className="text-error text-sm mt-1">
							{errors.password.message || "Please enter a valid password"}
						</p>
					)}

					<button
						type="submit"
						className="btn btn-primary mt-4"
						disabled={isAuthenticating}
					>
						{isAuthenticating ? (
							<>
								<span className="loading loading-ring loading-sm"></span>
								Logging in...
							</>
						) : (
							"Login"
						)}
					</button>

					<p className="text-center text-sm mt-4">
						Don&apos;t have an account?{" "}
						<Link to="/signup" className="link link-primary">
							Sign up
						</Link>
					</p>
				</form>
			</div>
		</div>
	);
}
