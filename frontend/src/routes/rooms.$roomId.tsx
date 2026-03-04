import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import MessageList from "@/components/MessageList";
import RoomsSidebar from "@/components/RoomsSidebar";
import useSocketEvents from "@/hooks/useSocketEvents";
import useMessagesStore from "@/stores/messagesStore";
import useRoomStore from "@/stores/roomStore";
import useToastStore from "@/stores/toastStore";

export const Route = createFileRoute("/rooms/$roomId")({
	component: RoomComponent,
	beforeLoad: ({ context }) => {
		// Access auth via context, not the store directly
		if (!context.auth.isAuthenticated) {
			throw redirect({ to: "/login" });
		}
	},
	loader: async ({ params }) => {
		const roomId = Number(params.roomId);

		const roomStore = useRoomStore.getState();
		const messageStore = useMessagesStore.getState();

		// Fetch rooms list once
		if (roomStore.roomsList.length === 0 && roomStore.hasMore) {
			await roomStore.getRooms();
		}

		// Fetch initial messages for this room
		const roomMessages = messageStore.getMessage(roomId);

		if (roomMessages.messages.length === 0 && roomMessages.hasMore) {
			await messageStore.fetchMessages(roomId);
		}

		return null;
	},

	staleTime: 30,
	gcTime: 1000 * 60 * 5,
});

function RoomComponent() {
	const navigate = useNavigate();
	const { roomId } = Route.useParams();
	const [message, setMessage] = useState("");

	useSocketEvents(Number(roomId));

	const room = useRoomStore((s) => s.room);
	const leave = useRoomStore((s) => s.leave);
	const isLeaving = useRoomStore((s) => s.isLeaving);
	const showToast = useToastStore((s) => s.show);

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

	const send = useMessagesStore((s) => s.sendMessage);

	const sendMessage = async () => {
		if (!message.trim()) return;

		await send(Number(roomId), message.trim());
		setMessage("");
	};

	return (
		<div className="h-screen flex bg-base-200">
			{/* Sidebar */}
			<RoomsSidebar />

			{/* Chat Area */}
			<div className="flex-1 flex flex-col">
				{/* Header */}
				<div className="navbar h-16 bg-base-100 border-b px-4">
					<div className="flex-1 font-semibold">
						Room: {room?.name || `Room ${roomId}`}
					</div>
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

				<MessageList roomId={Number(roomId)} />

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
