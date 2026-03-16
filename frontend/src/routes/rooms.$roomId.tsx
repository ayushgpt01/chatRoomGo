import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { ArrowLeft, Users } from "lucide-react";
import { useEffect, useState } from "react";
import MessageList from "@/components/MessageList";
import RoomsSidebar from "@/components/RoomsSidebar";
import useSocketEvents from "@/hooks/useSocketEvents";
import useMessagesStore from "@/stores/messagesStore";
import useRoomStore from "@/stores/roomStore";
import useSocketStore from "@/stores/socketStore";
import useToastStore from "@/stores/toastStore";
import { useTypingStore } from "@/stores/typingStore";
import { OutgoingEventTypes } from "@/types/events";

export const Route = createFileRoute("/rooms/$roomId")({
	component: RoomComponent,
	beforeLoad: ({ context }) => {
		if (!context.auth.isAuthenticated) {
			throw redirect({ to: "/login" });
		}
	},
	loader: async ({ params }) => {
		const roomId = Number(params.roomId);

		const roomStore = useRoomStore.getState();
		const messageStore = useMessagesStore.getState();

		if (roomStore.roomsList.length === 0 && roomStore.hasMore) {
			await roomStore.getRooms();
		}

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
	const updateUserMap = useTypingStore((s) => s.updateUserMap);
	const userMap = useTypingStore((s) => s.userMap);
	const socketSend = useSocketStore((s) => s.send);

	const [isTyping, setIsTyping] = useState(false);
	const [stopTypingTimeoutId, setStopTypingTimeoutId] = useState<number | null>(
		null,
	);

	useSocketEvents(Number(roomId));

	const roomsList = useRoomStore((s) => s.roomsList);
	const leave = useRoomStore((s) => s.leave);
	const isLeaving = useRoomStore((s) => s.isLeaving);
	const showToast = useToastStore((s) => s.show);
	const room = roomsList.find((r) => r.id === Number(roomId));

	useEffect(() => {
		const currentRoom = roomsList.find((r) => r.id === Number(roomId));
		if (
			currentRoom?.members &&
			currentRoom.members.length !== Object.values(userMap).length
		) {
			updateUserMap(currentRoom.members);
		}
	}, [roomsList, roomId, updateUserMap, userMap]);

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
		if (isTyping) {
			socketSend({
				type: OutgoingEventTypes.EventStopTyping,
				data: { roomId: Number(roomId) },
			});
			setIsTyping(false);
		}
	};

	const scheduleStopTyping = () => {
		if (stopTypingTimeoutId !== null) {
			window.clearTimeout(stopTypingTimeoutId);
		}

		const id = window.setTimeout(() => {
			socketSend({
				type: OutgoingEventTypes.EventStopTyping,
				data: { roomId: Number(roomId) },
			});
			setIsTyping(false);
			setStopTypingTimeoutId(null);
		}, 1200);
		setStopTypingTimeoutId(id);
	};

	const handleTyping = () => {
		if (!isTyping) {
			socketSend({
				type: OutgoingEventTypes.EventStartTyping,
				data: { roomId: Number(roomId) },
			});
			setIsTyping(true);
		}
		scheduleStopTyping();
	};

	return (
		<div className="h-screen flex bg-base-200">
			<RoomsSidebar />

			<div className="flex-1 flex flex-col">
				<div className="navbar h-16 bg-base-100 border-b px-4">
					<div className="flex items-center gap-3 flex-1">
						<button
							type="button"
							className="btn btn-ghost btn-circle btn-sm lg:hidden"
							onClick={() => navigate({ to: "/" })}
							aria-label="Go back"
						>
							<ArrowLeft className="w-4 h-4" />
						</button>

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

						<div className="flex items-center gap-2">
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

				<MessageList roomId={Number(roomId)} />

				<div className="p-4 bg-base-100 border-t">
					<div className="flex gap-2 items-end">
						<input
							type="text"
							className="input input-bordered flex-1"
							placeholder="Type a message..."
							value={message}
							onChange={(e) => {
								setMessage(e.target.value);
								handleTyping();
							}}
							onKeyDown={(e) => {
								if (e.key === "Enter" && !e.shiftKey) {
									e.preventDefault();
									sendMessage();
								}
							}}
							aria-label="Message input"
						/>

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
