import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";

export const Route = createFileRoute("/")({ component: App });

function App() {
	const navigate = useNavigate();
	const [roomId, setRoomId] = useState("");
	const isLoggedIn = false;

	const joinRoom = () => {
		if (!roomId.trim()) return;
		navigate({ to: "/rooms/$roomId", params: { roomId } });
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
							className="input input-bordered w-full"
							value={roomId}
							onChange={(e) => setRoomId(e.target.value)}
						/>
					</div>

					<button
						type="button"
						onClick={joinRoom}
						className="btn btn-primary mt-4"
					>
						Join Room
					</button>

					{!isLoggedIn && <>
					<div className="divider">or</div>

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
					</>}

				</div>
			</div>
		</div>
	);
}
