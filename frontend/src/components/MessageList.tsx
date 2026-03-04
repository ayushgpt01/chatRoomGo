import { Check, CheckCheck, MoreVertical } from "lucide-react";
import { useId, useState } from "react";
import { Virtuoso } from "react-virtuoso";
import useAuthStore from "@/stores/authStore";
import useMessagesStore from "@/stores/messagesStore";
import type { Message } from "@/types/message";
import { formatChatTime } from "@/utils/dateUtils";

interface Props {
	roomId: number;
}

function MessageItem({
	message,
	roomId,
}: {
	message: Message;
	roomId: number;
}) {
	const editMessage = useMessagesStore((s) => s.editMessage);
	const deleteMessage = useMessagesStore((s) => s.deleteMessage);
	const currentUserId = useAuthStore((s) => s.user?.id);

	const [editingId, setEditingId] = useState<number | null>(null);
	const [editValue, setEditValue] = useState("");
	const popoverId = useId();
	const anchorId = useId();

	const isMine = message.senderId === currentUserId;
	const isSending = Boolean(message.nonce);
	const isEdited = Boolean(message.editedAt);

	const initials = message.senderName
		.split(" ")
		.map((n) => n[0])
		.join("")
		.toUpperCase()
		.slice(0, 2);

	const formattedTime = formatChatTime(message.sentAt);

	return (
		<div className={`px-6 py-2 chat ${isMine ? "chat-end" : "chat-start"}`}>
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
				{!isMine && message.senderName}
				<time className="text-xs opacity-50 ml-2">{formattedTime}</time>
			</div>

			{/* Bubble */}
			<div
				className={`chat-bubble relative ${
					isMine ? "chat-bubble-primary" : ""
				}`}
			>
				{editingId === message.id ? (
					<div className="flex items-center gap-2 w-full">
						<input
							className="input input-sm input-bordered bg-base-100 text-base-content flex-1"
							value={editValue}
							onChange={(e) => setEditValue(e.target.value)}
							onKeyDown={async (e) => {
								if (e.key === "Enter") {
									await editMessage(roomId, message.id, editValue);
									setEditingId(null);
								}
								if (e.key === "Escape") {
									setEditingId(null);
								}
							}}
						/>

						<div className="flex gap-1">
							<button
								className="btn btn-xs btn-ghost"
								onClick={() => {
									setEditingId(null);
									setEditValue("");
								}}
								type="button"
							>
								Cancel
							</button>

							<button
								className="btn btn-xs btn-primary"
								onClick={async () => {
									await editMessage(roomId, message.id, editValue);
									setEditingId(null);
									setEditValue("");
								}}
								type="button"
							>
								Save
							</button>
						</div>
					</div>
				) : (
					<div className="flex justify-between items-start gap-2">
						<div className="flex-1 wrap-break-word">
							{message.content}
							{isEdited && (
								<span className="text-xs opacity-60 ml-2">(edited)</span>
							)}
						</div>

						{/* Dropdown Menu */}
						{isMine && !isSending && (
							<div>
								<button
									popoverTarget={popoverId}
									type="button"
									style={{ anchorName: anchorId }}
									className="btn btn-ghost btn-xs btn-circle"
								>
									<MoreVertical className="w-4 h-4" />
								</button>
								<ul
									className="dropdown menu p-2 shadow-xl text-base-content bg-base-100 border border-base-300 rounded-box w-36"
									popover="auto"
									id={popoverId}
									style={{ positionAnchor: anchorId }}
								>
									<li>
										<button
											onClick={() => {
												setEditingId(message.id);
												setEditValue(message.content);
											}}
											type="button"
										>
											Edit
										</button>
									</li>
									<li>
										<button
											onClick={() => {
												navigator.clipboard.writeText(message.content);
											}}
											type="button"
										>
											Copy
										</button>
									</li>
									<li>
										<button
											className="text-error"
											onClick={async () => {
												await deleteMessage(roomId, message.id);
											}}
											type="button"
										>
											Delete
										</button>
									</li>
								</ul>
							</div>
						)}
					</div>
				)}
			</div>
			{/* Footer */}
			{isMine && (
				<div className="chat-footer opacity-50 text-xs flex items-center gap-1 justify-end">
					{isSending ? (
						"Sending..."
					) : message.read ? (
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
}

export default function MessageList({ roomId }: Props) {
	const { messages, hasMore } = useMessagesStore((s) => s.getMessage(roomId));
	const fetchMessages = useMessagesStore((s) => s.fetchMessages);

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
			itemContent={(_, msg) => <MessageItem message={msg} roomId={roomId} />}
		/>
	);
}
