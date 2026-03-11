import { Check, CheckCheck, MoreVertical } from "lucide-react";
import { useId, useState } from "react";
import useAuthStore from "@/stores/authStore";
import useMessagesStore from "@/stores/messagesStore";
import type { User } from "@/types/auth";
import type { Message } from "@/types/message";
import { formatChatTime } from "@/utils/dateUtils";

const MAX_TIME = 5 * 60 * 1000;

export default function MessageItem({
	message,
	prevMessage,
	roomId,
	roomMembers,
	onShowReadReceipts,
}: {
	message: Message;
	prevMessage?: Message;
	roomId: number;
	roomMembers: User[];
	onShowReadReceipts: (message: Message) => void;
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
	const isConsecutive =
		prevMessage &&
		prevMessage.senderId === message.senderId &&
		new Date(message.sentAt).getTime() -
			new Date(prevMessage.sentAt).getTime() <
			MAX_TIME;

	// Calculate read status
	const getReadStatus = () => {
		if (isSending) return "sending";
		if (!message.delivered) return "sent";

		const otherMembers = roomMembers.filter((m) => m.id !== message.senderId);
		if (otherMembers.length === 0) return "delivered";

		const readCount =
			message.readBy?.filter((id) => otherMembers.some((m) => m.id === id))
				.length || 0;

		if (readCount === otherMembers.length) return "read";
		if (readCount > 0) return "delivered";
		return "delivered";
	};

	const readStatus = getReadStatus();

	const initials = message.senderName
		.split(" ")
		.map((n) => n[0])
		.join("")
		.toUpperCase()
		.slice(0, 2);

	const formattedTime = formatChatTime(message.sentAt);
	const showAvatar = !isConsecutive;

	return (
		<article
			className={`chat group px-6 py-1 ${isMine ? "chat-end" : "chat-start"}`}
			aria-label={`${isMine ? "Your" : `${message.senderName}'s`} message${isEdited ? " (edited)" : ""}`}
		>
			{showAvatar && (
				<div className="chat-image avatar avatar-placeholder">
					<div
						className={`w-10 rounded-full transition-transform duration-200 group-hover:scale-105 ${
							isMine
								? "bg-primary text-primary-content"
								: "bg-secondary text-primary-content"
						}`}
					>
						<span className="text-sm font-medium">{initials}</span>
					</div>
				</div>
			)}

			{!isConsecutive && !isMine && (
				<div className="chat-header pb-1 text-sm opacity-80">
					{message.senderName}
					<time className="text-xs opacity-50 ml-2" dateTime={message.sentAt}>
						{formattedTime}
					</time>
				</div>
			)}

			<div
				className={`chat-bubble relative transition-all duration-200 group-hover:shadow-lg ${
					isMine ? "chat-bubble-primary" : "chat-bubble-secondary"
				} ${isConsecutive ? "rounded-2xl" : "rounded-3xl"} before:hidden after:hidden`}
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
								<span className="text-xs opacity-60 ml-2 italic">(edited)</span>
							)}
						</div>

						{isMine && !isSending && (
							<div className="opacity-0 group-hover:opacity-100 transition-opacity duration-200">
								<button
									popoverTarget={popoverId}
									type="button"
									style={{ anchorName: anchorId }}
									className="btn btn-ghost btn-xs btn-circle hover:bg-base-300 transition-colors duration-200"
									aria-label="Message options"
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

			{isMine && (
				<div className="chat-footer opacity-60 text-xs flex items-center gap-1 justify-end mt-1">
					{isSending ? (
						<span className="text-xs opacity-70 animate-pulse">Sending...</span>
					) : (
						<div className="flex items-center gap-1">
							<button
								type="button"
								onClick={() => onShowReadReceipts(message)}
								className={`transition-colors duration-200 hover:opacity-80 ${
									readStatus === "read" ? "text-blue-500" : "opacity-70"
								}`}
								aria-label={
									readStatus === "sent"
										? "Sent"
										: readStatus === "delivered"
											? "Delivered"
											: "Read by everyone"
								}
							>
								{readStatus === "sent" && <Check className="w-4 h-4" />}
								{readStatus === "delivered" && (
									<CheckCheck className="w-4 h-4" />
								)}
								{readStatus === "read" && (
									<CheckCheck className="w-4 h-4 text-blue-500" />
								)}
							</button>
							{!isConsecutive && (
								<time className="text-xs opacity-50" dateTime={message.sentAt}>
									{formattedTime}
								</time>
							)}
						</div>
					)}
				</div>
			)}
		</article>
	);
}
