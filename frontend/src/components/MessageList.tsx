import { useRef, useState } from "react";
import { Virtuoso, type VirtuosoHandle } from "react-virtuoso";
import ReadReceiptsModal from "@/components/ReadReceiptsModal";
import useMessagesStore from "@/stores/messagesStore";
import useRoomStore from "@/stores/roomStore";
import { useTypingUsers } from "@/stores/typingStore";
import type { Message } from "@/types/message";
import MessageItem from "./MessageItem";

interface Props {
	roomId: number;
	showUnreadIndicator?: boolean;
}

export default function MessageList({
	roomId,
	showUnreadIndicator = false,
}: Props) {
	const { messages, hasMore } = useMessagesStore((s) => s.getMessage(roomId));
	const fetchMessages = useMessagesStore((s) => s.fetchMessages);
	const { roomsList } = useRoomStore((s) => s);
	const virtuosoRef = useRef<VirtuosoHandle>(null);
	const [isAtBottom, setIsAtBottom] = useState(true);
	const typingUsers = useTypingUsers(roomId);
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

	const handleScroll = (isScrolling: boolean) => {
		if (!isScrolling) {
			setIsAtBottom(true);
		}
	};

	const handleScrollToBottomChange = (atBottom: boolean) => {
		setIsAtBottom(atBottom);
	};

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
		<div className="h-full relative">
			<Virtuoso
				ref={virtuosoRef}
				style={{ height: "100%" }}
				data={processedData}
				followOutput="smooth"
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
						const nextItem = processedData[index + 1] as Message;
						const msg = nextItem;
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
								? `${typingUsers[0].name} is typing`
								: `${typingUsers.map((u) => u.name).join(", ")} are typing`}
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
