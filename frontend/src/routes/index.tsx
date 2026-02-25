import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import useAuthStore from "@/stores/authStore";
import useRoomStore from "@/stores/roomStore";
import useToastStore from "@/stores/toastStore";

export const Route = createFileRoute("/")({ component: App });

function App() {
	const navigate = useNavigate();
	const [roomId, setRoomId] = useState("");
	const [error, setError] = useState<string | null>(null);

	const isLoggedIn = useAuthStore((s) => s.isAuthenticated);
	const logout = useAuthStore((s) => s.logout);
	const setAuth = useAuthStore((s) => s.setAuth);

	const join = useRoomStore((s) => s.join);

	const showToast = useToastStore((s) => s.show);

	const joinRoom = async () => {
		const id = Number(roomId);

		if (Number.isNaN(id)) {
			setError("Invalid room id");
			return;
		}

		try {
			const res = await join(id);

			if (res.login) {
				setAuth(res.login);
			}

			showToast("Logged in successfully", "success");
			navigate({ to: "/rooms/$roomId", params: { roomId } });
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Could not join the room",
				"error",
			);
			console.error(err);
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-base-200">
			<div className="card w-full max-w-md bg-base-100 shadow-xl">
				<div className="card-body">
					<h1 className="text-3xl font-bold text-center">Go Chat</h1>

					<div className="form-control mt-4">
						<input
							type="text"
							placeholder="Enter Room ID"
							className={`input input-bordered w-full ${
								error ? "input-error" : ""
							}`}
							value={roomId}
							onChange={(e) => setRoomId(e.target.value)}
						/>
						{error && (
							<p className="text-error text-sm mt-1">
								{error || "Please enter a valid room id"}
							</p>
						)}
					</div>

					<button
						type="button"
						onClick={joinRoom}
						className="btn btn-primary mt-4"
					>
						Join Room
					</button>

					<div className="divider">or</div>
					{!isLoggedIn ? (
						<div className="flex flex-col gap-2">
							<button
								type="button"
								onClick={() => navigate({ to: "/login" })}
								className="btn btn-outline w-full"
							>
								Login
							</button>
							<button
								type="button"
								onClick={() => navigate({ to: "/signup" })}
								className="btn btn-secondary w-full"
							>
								Sign Up
							</button>
						</div>
					) : (
						<div className="flex flex-col gap-2">
							<button
								type="button"
								onClick={logout}
								className="btn btn-secondary w-full"
							>
								Logout
							</button>
						</div>
					)}
				</div>
			</div>
		</div>
	);
}
