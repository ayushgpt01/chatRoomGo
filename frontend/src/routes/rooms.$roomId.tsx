import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import useRoomStore from "@/stores/roomStore";
import useToastStore from "@/stores/toastStore";

export const Route = createFileRoute("/rooms/$roomId")({
	component: RouteComponent,
	beforeLoad: ({ context }) => {
		// Access auth via context, not the store directly
		if (!context.auth.isAuthenticated) {
			throw redirect({ to: "/login" });
		}
	},
});

let counter = 3;

function RouteComponent() {
	const navigate = useNavigate();
	const { roomId } = Route.useParams();
	const leave = useRoomStore((s) => s.leave);
	const isLeaving = useRoomStore((s) => s.isLeaving);

	const showToast = useToastStore((s) => s.show);
	const [message, setMessage] = useState("");
	const [messages, setMessages] = useState([
		{ id: 1, message: "Welcome to the room!" },
		{ id: 2, message: "This UI is clean." },
	]);

	const sendMessage = () => {
		if (!message.trim()) return;
		setMessages((prev) => [...prev, { id: counter++, message }]);
		setMessage("");
	};

	const handleLeave = async () => {
		try {
			await leave();
			navigate({ to: "/" });
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Could not leave the room",
				"error",
			);
			console.error(err);
		}
	};

	return (
		<div className="h-screen flex bg-base-200">
			{/* Sidebar */}
			<div className="w-64 bg-base-100 border-r hidden md:block">
				<div className="p-4 font-bold text-lg">Room: {roomId}</div>
				<ul className="menu p-2">
					<li>
						<a href="#s">General</a>
					</li>
					<li>
						<a href="#s">Random</a>
					</li>
				</ul>
			</div>

			{/* Chat Area */}
			<div className="flex-1 flex flex-col">
				{/* Header */}
				<div className="navbar bg-base-100 border-b px-4">
					<div className="flex-1 font-semibold">Room: {roomId}</div>
					<button
						type="button"
						className="btn btn-sm btn-outline"
						onClick={handleLeave}
						disabled={isLeaving}
					>
						{isLeaving ? (
							<>
								<span className="loading loading-ring loading-sm"></span>
								Leaving...
							</>
						) : (
							"Leave"
						)}
					</button>
				</div>

				{/* Messages */}
				<div className="flex-1 overflow-y-auto p-4 space-y-3">
					{messages.map((msg) => (
						<div key={msg.id} className="chat chat-start">
							<div className="chat-bubble">{msg.message}</div>
						</div>
					))}
				</div>

				{/* Input */}
				<div className="p-4 bg-base-100 border-t flex gap-2">
					<input
						type="text"
						className="input input-bordered flex-1"
						placeholder="Type a message..."
						value={message}
						onChange={(e) => setMessage(e.target.value)}
						onKeyDown={(e) => e.key === "Enter" && sendMessage()}
					/>
					<button
						type="button"
						onClick={sendMessage}
						className="btn btn-primary"
					>
						Send
					</button>
				</div>
			</div>
		</div>
	);
}
