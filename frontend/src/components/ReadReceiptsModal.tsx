import { Check, CheckCheck, X } from "lucide-react";
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

	const readByUsers = roomMembers.filter(
		(member) =>
			message.readBy?.includes(member.id) && member.id !== message.senderId,
	);

	const unreadByUsers = roomMembers.filter(
		(member) =>
			!message.readBy?.includes(member.id) && member.id !== message.senderId,
	);

	const totalRecipients = roomMembers.length - 1;
	const readCount = readByUsers.length;
	const unreadCount = unreadByUsers.length;
	const progress =
		totalRecipients === 0 ? 0 : Math.round((readCount / totalRecipients) * 100);

	const avatarInitials = (name: string) =>
		name
			.split(" ")
			.map((n) => n[0])
			.join("")
			.slice(0, 2)
			.toUpperCase();

	return (
		<div className="modal modal-open">
			<div className="modal-box max-w-lg p-0 overflow-hidden">
				<div className="flex items-center justify-between px-6 py-4 border-b bg-base-200">
					<h3 className="font-semibold text-lg">Read Receipts</h3>
					<button
						onClick={onClose}
						className="btn btn-sm btn-circle btn-ghost"
						aria-label="Close modal"
						type="button"
					>
						<X className="w-4 h-4" />
					</button>
				</div>

				<div className="px-6 py-4 border-b bg-base-100">
					<p className="text-sm opacity-70 mb-2">
						Sent at {formatChatTime(message.sentAt)}
					</p>

					<div className="bg-base-200 rounded-xl p-3 text-sm">
						{message.content}
					</div>

					<div className="mt-3">
						<div className="flex justify-between text-xs opacity-70 mb-1">
							<span>{readCount} read</span>
							<span>{unreadCount} unread</span>
						</div>

						<progress
							className="progress progress-primary w-full"
							value={progress}
							max="100"
						/>
					</div>
				</div>

				<div className="max-h-80 overflow-y-auto">
					{readByUsers.length > 0 && (
						<div className="px-6 py-4 border-b">
							<h4 className="text-sm font-semibold text-success mb-3 flex items-center gap-2">
								<CheckCheck className="w-4 h-4" />
								Read ({readByUsers.length})
							</h4>

							<div className="space-y-3">
								{readByUsers.map((user) => (
									<div key={user.id} className="flex items-center gap-3">
										<div className="avatar placeholder">
											<div className="w-9 rounded-full bg-primary text-primary-content">
												<span className="text-xs">
													{avatarInitials(user.name)}
												</span>
											</div>
										</div>

										<div className="flex-1">
											<p className="text-sm font-medium">{user.name}</p>
											<p className="text-xs opacity-60">@{user.username}</p>
										</div>

										<CheckCheck className="w-4 h-4 text-success" />
									</div>
								))}
							</div>
						</div>
					)}

					{unreadByUsers.length > 0 && (
						<div className="px-6 py-4">
							<h4 className="text-sm font-semibold text-base-content/60 mb-3 flex items-center gap-2">
								<Check className="w-4 h-4" />
								Not read yet ({unreadByUsers.length})
							</h4>

							<div className="space-y-3">
								{unreadByUsers.map((user) => (
									<div
										key={user.id}
										className="flex items-center gap-3 opacity-70"
									>
										<div className="avatar placeholder">
											<div className="w-9 rounded-full bg-base-300">
												<span className="text-xs">
													{avatarInitials(user.name)}
												</span>
											</div>
										</div>

										<div className="flex-1">
											<p className="text-sm font-medium">{user.name}</p>
											<p className="text-xs opacity-60">@{user.username}</p>
										</div>

										<Check className="w-4 h-4 text-base-content/40" />
									</div>
								))}
							</div>
						</div>
					)}

					{readByUsers.length === 0 && unreadByUsers.length === 0 && (
						<div className="py-10 text-center opacity-60">
							No other members in this room
						</div>
					)}
				</div>

				<div className="modal-action px-6 py-4 border-t">
					<button
						className="btn btn-primary w-full"
						onClick={onClose}
						type="button"
					>
						Close
					</button>
				</div>
			</div>
		</div>
	);
}
