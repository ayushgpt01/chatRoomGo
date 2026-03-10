import { X } from "lucide-react";
import type { User } from "@/types/auth";
import type { Message } from "@/types/message";
import { formatChatTime } from "@/utils/dateUtils";

interface Props {
	message: Message;
	roomMembers: User[];
	isOpen: boolean;
	onClose: () => void;
}

export default function ReadReceiptsModal({
	message,
	roomMembers,
	isOpen,
	onClose,
}: Props) {
	if (!isOpen) return null;

	// Get users who have read the message
	const readByUsers = roomMembers.filter(
		(member) =>
			message.readBy?.includes(member.id) && member.id !== message.senderId,
	);

	// Get users who haven't read the message (excluding sender)
	const unreadByUsers = roomMembers.filter(
		(member) =>
			!message.readBy?.includes(member.id) && member.id !== message.senderId,
	);

	const totalRecipients = roomMembers.length - 1; // Exclude sender
	const readCount = readByUsers.length;
	const unreadCount = unreadByUsers.length;

	return (
		<div className="modal modal-open">
			<div className="modal-box max-w-md">
				<div className="flex items-center justify-between mb-4">
					<h3 className="font-bold text-lg">Read Receipts</h3>
					<button
						onClick={onClose}
						className="btn btn-sm btn-circle btn-ghost"
						aria-label="Close modal"
						type="button"
					>
						<X className="w-4 h-4" />
					</button>
				</div>

				<div className="mb-4">
					<p className="text-sm opacity-70 mb-2">
						Sent at {formatChatTime(message.sentAt)}
					</p>
					<div className="flex gap-4 text-sm">
						<span className="text-green-600">
							✓ Read by {readCount}/{totalRecipients}
						</span>
						{unreadCount > 0 && (
							<span className="text-gray-500">{unreadCount} unread</span>
						)}
					</div>
				</div>

				{/* Read by section */}
				{readByUsers.length > 0 && (
					<div className="mb-6">
						<h4 className="font-medium text-sm mb-2 text-green-600">
							Read by ({readByUsers.length})
						</h4>
						<div className="space-y-2 max-h-40 overflow-y-auto">
							{readByUsers.map((user) => (
								<div key={user.id} className="flex items-center gap-3">
									<div className="avatar placeholder">
										<div className="w-8 rounded-full bg-primary text-primary-content">
											<span className="text-xs">
												{user.name.charAt(0).toUpperCase()}
											</span>
										</div>
									</div>
									<div className="flex-1">
										<p className="text-sm font-medium">{user.name}</p>
										<p className="text-xs opacity-60">@{user.username}</p>
									</div>
									<div className="text-green-500">✓✓</div>
								</div>
							))}
						</div>
					</div>
				)}

				{/* Unread by section */}
				{unreadByUsers.length > 0 && (
					<div>
						<h4 className="font-medium text-sm mb-2 text-gray-500">
							Not read yet ({unreadByUsers.length})
						</h4>
						<div className="space-y-2 max-h-40 overflow-y-auto">
							{unreadByUsers.map((user) => (
								<div
									key={user.id}
									className="flex items-center gap-3 opacity-60"
								>
									<div className="avatar placeholder">
										<div className="w-8 rounded-full bg-base-300">
											<span className="text-xs">
												{user.name.charAt(0).toUpperCase()}
											</span>
										</div>
									</div>
									<div className="flex-1">
										<p className="text-sm font-medium">{user.name}</p>
										<p className="text-xs opacity-60">@{user.username}</p>
									</div>
									<div className="text-gray-400">✓</div>
								</div>
							))}
						</div>
					</div>
				)}

				{readByUsers.length === 0 && unreadByUsers.length === 0 && (
					<div className="text-center py-8 opacity-60">
						<p>No other members in this room</p>
					</div>
				)}

				<div className="modal-action">
					<button className="btn btn-primary" onClick={onClose} type="button">
						Close
					</button>
				</div>
			</div>
		</div>
	);
}
