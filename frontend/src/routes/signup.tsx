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
import { type SignupCredentials, SignupSchema } from "@/types/auth";

export const Route = createFileRoute("/signup")({
	beforeLoad: ({ context }) => {
		if (context.auth.isAuthenticated) {
			throw redirect({
				to: "/",
			});
		}
	},
	component: RouteComponent,
});

function RouteComponent() {
	const signup = useAuthStore((state) => state.signup);
	const isCreating = useAuthStore((state) => state.isCreating);
	const navigate = useNavigate();
	const {
		register,
		handleSubmit,
		formState: { errors },
		resetField,
	} = useForm({
		resolver: zodResolver(SignupSchema),
	});
	const showToast = useToastStore((s) => s.show);

	const onSubmit: SubmitHandler<SignupCredentials> = async (data) => {
		try {
			await signup(data);
			showToast("Signed up successfully", "success");
			navigate({ to: "/" });
		} catch (err) {
			showToast("Error creating user", "error");
			console.error(err);
			resetField("password");
			resetField("confirmPassword");
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-base-200">
			<div className="card w-full max-w-md bg-base-100 shadow-xl">
				<form onSubmit={handleSubmit(onSubmit)} className="card-body space-y-2">
					<h2 className="text-2xl font-bold text-center">Create Account</h2>

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
							<span className="label-text">Display Name</span>
						</div>
						<input
							type="text"
							placeholder="Name"
							className={`input input-bordered w-full mt-1 ${
								errors.name ? "input-error" : ""
							}`}
							{...register("name")}
						/>
					</label>
					{errors.name && (
						<p className="text-error text-sm mt-1">
							{errors.name.message || "Please enter a valid name"}
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

					<label className="form-control w-full mt-2">
						<div className="label">
							<span className="label-text">Confirm Password</span>
						</div>
						<input
							type="password"
							placeholder="Confirm Password"
							className={`input input-bordered w-full mt-1 ${
								errors.confirmPassword ? "input-error" : ""
							}`}
							{...register("confirmPassword")}
						/>
					</label>
					{errors.confirmPassword && (
						<p className="text-error text-sm mt-1">
							{errors.confirmPassword.message}
						</p>
					)}

					<button
						type="submit"
						className="btn btn-primary mt-4"
						disabled={isCreating}
					>
						{isCreating ? (
							<>
								<span className="loading loading-ring loading-sm"></span>
								Creating your account...
							</>
						) : (
							"Sign Up"
						)}
					</button>

					<p className="text-center text-sm mt-4">
						Already have an account?{" "}
						<Link to="/login" className="link link-primary">
							Login
						</Link>
					</p>
				</form>
			</div>
		</div>
	);
}
