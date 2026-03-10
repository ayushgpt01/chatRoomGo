import {
	Check,
	CheckCheck,
	ChevronDown as ChevronDownIcon,
	MoreVertical,
} from "lucide-react";
import { useCallback, useId, useRef, useState } from "react";
import { Virtuoso, type VirtuosoHandle } from "react-virtuoso";
import ReadReceiptsModal from "@/components/ReadReceiptsModal";
import useAuthStore from "@/stores/authStore";
import useMessagesStore from "@/stores/messagesStore";
import useRoomStore from "@/stores/roomStore";
import type { User } from "@/types/auth";
import type { Message } from "@/types/message";
import { formatChatTime } from "@/utils/dateUtils";

interface Props {
	roomId: number;
	typingUsers?: string[];
	showUnreadIndicator?: boolean;
}

function MessageItem({
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
			300000; // 5 minutes

	// Calculate read status
	const getReadStatus = () => {
		if (isSending) return "sending";
		if (!message.delivered) return "sent";

		// Get room members excluding sender
		const otherMembers = roomMembers.filter((m) => m.id !== message.senderId);
		if (otherMembers.length === 0) return "delivered"; // No one else to read

		// Check if all other members have read it
		const readCount =
			message.readBy?.filter((id) => otherMembers.some((m) => m.id === id))
				.length || 0;

		if (readCount === otherMembers.length) return "read";
		if (readCount > 0) return "delivered"; // Some but not all have read
		return "delivered"; // No one has read yet
	};

	const readStatus = getReadStatus();

	const initials = message.senderName
		.split(" ")
		.map((n) => n[0])
		.join("")
		.toUpperCase()
		.slice(0, 2);

	const formattedTime = formatChatTime(message.sentAt);
	const showAvatar = !isConsecutive || !isMine;

	return (
		<article
			className={`group px-6 py-1 ${isMine ? "chat-end" : "chat-start"} ${isConsecutive ? "mt-1" : "mt-4"} transition-all duration-200`}
			aria-label={`${isMine ? "Your" : `${message.senderName}'s`} message${isEdited ? " (edited)" : ""}`}
		>
			{showAvatar && (
				<div className="chat-image avatar avatar-placeholder">
					<div
						className={`w-10 rounded-full transition-transform duration-200 group-hover:scale-105 ${
							isMine ? "bg-primary text-primary-content" : "bg-base-300"
						}`}
					>
						<span className="text-sm font-medium">{initials}</span>
					</div>
				</div>
			)}

			{/* Header */}
			{!isConsecutive && (
				<div className="chat-header pb-1">
					{!isMine && (
						<>
							<span className="font-medium">{message.senderName}</span>
							<time
								className="text-xs opacity-50 ml-2"
								dateTime={message.sentAt}
							>
								{formattedTime}
							</time>
						</>
					)}
				</div>
			)}

			{/* Bubble */}
			<div
				className={`chat-bubble relative transition-all duration-200 group-hover:shadow-lg ${
					isMine ? "chat-bubble-primary" : "bg-base-200 text-base-content"
				} ${isConsecutive ? "rounded-2xl" : "rounded-3xl"}`}
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

						{/* Dropdown Menu */}
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
			{/* Footer with read receipts */}
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

export default function MessageList({
	roomId,
	typingUsers = [],
	showUnreadIndicator = false,
}: Props) {
	const { messages, hasMore } = useMessagesStore((s) => s.getMessage(roomId));
	const fetchMessages = useMessagesStore((s) => s.fetchMessages);
	const { roomsList } = useRoomStore((s) => s);
	const virtuosoRef = useRef<VirtuosoHandle>(null);
	const [showScrollButton, setShowScrollButton] = useState(false);
	const [isAtBottom, setIsAtBottom] = useState(true);
	const [readReceiptsModal, setReadReceiptsModal] = useState<{
		isOpen: boolean;
		message: Message | null;
	}>({
		isOpen: false,
		message: null,
	});

	// Get current room members
	const currentRoom = roomsList.find((r) => r.id === roomId);
	const roomMembers = currentRoom?.members || [];

	// Modal handlers
	const handleShowReadReceipts = (message: Message) => {
		setReadReceiptsModal({
			isOpen: true,
			message,
		});
	};

	const handleCloseReadReceipts = () => {
		setReadReceiptsModal({
			isOpen: false,
			message: null,
		});
	};

	// Scroll handlers
	const scrollToBottom = useCallback(() => {
		virtuosoRef.current?.scrollToIndex({ index: "LAST", behavior: "smooth" });
	}, []);

	const handleScroll = useCallback((isScrolling: boolean) => {
		if (!isScrolling) {
			setIsAtBottom(true);
			setShowScrollButton(false);
		}
	}, []);

	const handleScrollToBottomChange = useCallback(
		(atBottom: boolean) => {
			setIsAtBottom(atBottom);
			setShowScrollButton(!atBottom && messages.length > 0);
		},
		[messages.length],
	);

	// Group messages and add date separators
	const processedData = messages.reduce(
		(acc: (Message | "separator")[], msg, index) => {
			const prevMsg = messages[index - 1];

			// Add date separator if needed
			if (
				index === 0 ||
				new Date(msg.sentAt).toDateString() !==
					new Date(prevMsg.sentAt).toDateString()
			) {
				acc.push("separator");
			}

			acc.push(msg);
			return acc;
		},
		[],
	);

	// Empty state
	if (messages.length === 0) {
		return (
			<div className="flex-1 flex items-center justify-center">
				<div className="text-center max-w-md">
					<div
						className="w-20 h-20 bg-base-300 rounded-full flex items-center justify-center mx-auto mb-4"
						role="img"
						aria-label="Chat icon"
					>
						<svg
							className="w-10 h-10 opacity-50"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<title>Chat icon</title>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
							/>
						</svg>
					</div>
					<h3 className="text-lg font-medium mb-2">No messages yet</h3>
					<p className="text-sm opacity-60">
						Start the conversation with a friendly message!
					</p>
				</div>
			</div>
		);
	}

	return (
		<div className="flex-1 relative">
			<Virtuoso
				ref={virtuosoRef}
				className="flex-1"
				data={processedData}
				alignToBottom
				followOutput="auto"
				atBottomStateChange={handleScrollToBottomChange}
				isScrolling={handleScroll}
				startReached={() => {
					if (hasMore) fetchMessages(roomId);
				}}
				computeItemKey={(index, item) => {
					if (item === "separator") return `sep-${index}`;
					const msg = item as Message;
					return msg.id ? String(msg.id) : (msg.nonce ?? String(index));
				}}
				itemContent={(index, item) => {
					if (item === "separator") {
						const msgIndex = index > 0 ? index - 1 : 0;
						const msg = messages[msgIndex] || messages[0];
						if (!msg) return null;

						return (
							<div className="flex items-center justify-center my-4">
								<div className="bg-base-300 px-3 py-1 rounded-full">
									<span className="text-xs font-medium opacity-70">
										{new Date(msg.sentAt).toLocaleDateString("en-US", {
											weekday: "long",
											month: "short",
											day: "numeric",
										})}
									</span>
								</div>
							</div>
						);
					}

					const msg = item as Message;
					const prevMsg =
						index > 0 && processedData[index - 1] !== "separator"
							? (processedData[index - 1] as Message)
							: undefined;

					return (
						<MessageItem
							message={msg}
							prevMessage={prevMsg}
							roomId={roomId}
							roomMembers={roomMembers}
							onShowReadReceipts={handleShowReadReceipts}
						/>
					);
				}}
			/>

			{/* Scroll to bottom button */}
			{showScrollButton && (
				<button
					type="button"
					onClick={scrollToBottom}
					className="absolute bottom-20 right-4 btn btn-circle btn-primary shadow-lg transition-all duration-200 hover:scale-110"
					aria-label="Scroll to bottom"
				>
					<ChevronDownIcon className="w-5 h-5" />
				</button>
			)}

			{/* Unread indicator */}
			{showUnreadIndicator && !isAtBottom && (
				<div className="absolute top-0 left-0 right-0 bg-linear-to-b from-primary/20 to-transparent h-8 pointer-events-none" />
			)}

			{/* Typing Indicator */}
			{typingUsers.length > 0 && (
				<div className="absolute bottom-0 left-0 right-0 bg-linear-to-t from-base-100 to-transparent p-4">
					<div className="flex items-center gap-2">
						<div className="flex gap-1">
							<div
								className="w-2 h-2 bg-primary rounded-full animate-bounce"
								style={{ animationDelay: "0ms" }}
							></div>
							<div
								className="w-2 h-2 bg-primary rounded-full animate-bounce"
								style={{ animationDelay: "150ms" }}
							></div>
							<div
								className="w-2 h-2 bg-primary rounded-full animate-bounce"
								style={{ animationDelay: "300ms" }}
							></div>
						</div>
						<span className="text-sm opacity-70">
							{typingUsers.length === 1
								? `${typingUsers[0]} is typing`
								: `${typingUsers.join(", ")} are typing`}
						</span>
					</div>
				</div>
			)}

			{/* Read Receipts Modal */}
			{readReceiptsModal.message && (
				<ReadReceiptsModal
					message={readReceiptsModal.message}
					roomMembers={roomMembers}
					isOpen={readReceiptsModal.isOpen}
					onClose={handleCloseReadReceipts}
				/>
			)}
		</div>
	);
}
