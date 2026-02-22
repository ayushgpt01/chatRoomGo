import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";

export const Route = createFileRoute("/rooms/$roomId")({
	component: RouteComponent,
});

let counter = 3;

function RouteComponent() {
	const { roomId } = Route.useParams();
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
					<button type="button" className="btn btn-sm btn-outline">
						Leave
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
