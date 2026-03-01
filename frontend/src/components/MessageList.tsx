import { Check, CheckCheck } from "lucide-react";
import { Virtuoso } from "react-virtuoso";
import useAuthStore from "@/stores/authStore";
import useMessagesStore from "@/stores/messagesStore";
import { formatChatTime } from "@/utils/dateUtils";

interface Props {
	roomId: number;
}

export default function MessageList({ roomId }: Props) {
	const { messages, hasMore } = useMessagesStore((s) => s.getMessage(roomId));
	const fetchMessages = useMessagesStore((s) => s.fetchMessages);
	const currentUserId = useAuthStore((s) => s.user?.id);

	return (
		<Virtuoso
			className="flex-1"
			data={messages}
			alignToBottom
			followOutput="auto"
			startReached={() => {
				if (hasMore) fetchMessages(roomId);
			}}
			computeItemKey={(index, msg) =>
				msg.id ? String(msg.id) : (msg.nonce ?? String(index))
			}
			itemContent={(_, msg) => {
				const isMine = msg.senderId === currentUserId;
				const isSending = Boolean(msg.nonce);
				const isEdited = Boolean(msg.editedAt);

				const initials = msg.senderName
					.split(" ")
					.map((n) => n[0])
					.join("")
					.toUpperCase()
					.slice(0, 2);

				const formattedTime = formatChatTime(msg.sentAt);

				return (
					<div
						className={`px-6 py-2 chat ${isMine ? "chat-end" : "chat-start"}`}
					>
						{/* Avatar */}
						<div className="chat-image avatar avatar-placeholder">
							<div
								className={`w-10 rounded-full ${
									isMine ? "bg-primary text-primary-content" : "bg-base-300"
								}`}
							>
								<span className="text-sm font-medium">{initials}</span>
							</div>
						</div>

						{/* Header */}
						<div className="chat-header pb-1">
							{!isMine && msg.senderName}
							<time className="text-xs opacity-50 ml-2">{formattedTime}</time>
						</div>

						{/* Bubble */}
						<div
							className={`chat-bubble ${isMine ? "chat-bubble-primary" : ""}`}
						>
							{msg.content}
							{isEdited && (
								<span className="text-xs opacity-60 ml-2">(edited)</span>
							)}
						</div>

						{/* Footer */}
						{isMine && (
							<div className="chat-footer opacity-50 text-xs flex items-center gap-1 justify-end">
								{isSending ? (
									"Sending..."
								) : msg.read ? (
									<>
										<Check className="w-4 h-4" />
										<CheckCheck className="w-4 h-4" />
										Seen
									</>
								) : (
									<>
										<Check className="w-4 h-4" />
										Delivered
									</>
								)}
							</div>
						)}
					</div>
				);
			}}
		/>
	);
}
