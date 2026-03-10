import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { ArrowLeft, Search, Users } from "lucide-react";
import { useState } from "react";
import MessageList from "@/components/MessageList";
import RoomsSidebar from "@/components/RoomsSidebar";
import useSocketEvents from "@/hooks/useSocketEvents";
import useMessagesStore from "@/stores/messagesStore";
import useRoomStore from "@/stores/roomStore";
import useToastStore from "@/stores/toastStore";
import { useTypingStore } from "@/stores/typingStore";

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
	const [showSearch, setShowSearch] = useState(false);
	const getTypingUsers = useTypingStore((s) => s.getTypingUsers);

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
				{/* Enhanced Header */}
				<div className="navbar h-16 bg-base-100 border-b px-4">
					<div className="flex items-center gap-3 flex-1">
						{/* Back Button */}
						<button
							type="button"
							className="btn btn-ghost btn-circle btn-sm lg:hidden"
							onClick={() => navigate({ to: "/" })}
							aria-label="Go back"
						>
							<ArrowLeft className="w-4 h-4" />
						</button>

						{/* Room Info */}
						<div className="flex items-center gap-3 flex-1">
							<div className="hidden sm:block">
								<div className="font-semibold text-lg">
									{room?.name || `Room ${roomId}`}
								</div>
								<div className="text-xs opacity-60 flex items-center gap-1">
									<Users className="w-3 h-3" />
									{room?.participantCount || 1} participants
								</div>
							</div>
						</div>

						{/* Header Actions */}
						<div className="flex items-center gap-2">
							<button
								type="button"
								className="btn btn-ghost btn-circle btn-sm"
								onClick={() => setShowSearch(!showSearch)}
								aria-label="Search messages"
							>
								<Search className="w-4 h-4" />
							</button>

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
					</div>
				</div>

				{/* Search Bar */}
				{showSearch && (
					<div className="bg-base-200 border-b px-4 py-2">
						<input
							type="text"
							className="input input-sm input-bordered w-full"
							placeholder="Search messages..."
							// TODO: Implement search functionality
						/>
					</div>
				)}

				<MessageList
					roomId={Number(roomId)}
					typingUsers={getTypingUsers(Number(roomId))}
				/>

				{/* Enhanced Input */}
				<div className="p-4 bg-base-100 border-t">
					<div className="flex gap-2 items-end">
						<input
							type="text"
							className="input input-bordered flex-1"
							placeholder="Type a message..."
							value={message}
							onChange={(e) => setMessage(e.target.value)}
							onKeyDown={(e) => {
								if (e.key === "Enter" && !e.shiftKey) {
									e.preventDefault();
									sendMessage();
								}
							}}
							aria-label="Message input"
						/>

						{/* TODO: Add emoji picker, file upload buttons */}
						<button
							type="button"
							onClick={sendMessage}
							className="btn btn-primary btn-circle"
							aria-label="Send message"
							disabled={!message.trim()}
						>
							<svg
								className="w-5 h-5"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
								role="img"
								aria-label="Send icon"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
								/>
							</svg>
						</button>
					</div>
				</div>
			</div>
		</div>
	);
}
