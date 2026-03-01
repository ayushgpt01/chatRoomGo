import { create } from "zustand";
import messageService from "@/services/messageService";
import type { Message } from "@/types/message";
import { getErrorMessage } from "@/utils/errorHandler";

interface MessageStore {
	messages: Message[];
	cursor: string | null;
	hasMore: boolean;
}

export const sampleMessages: Message[] = [
	{
		id: 1,
		content: "Hey everyone 👋",
		senderName: "Alice Johnson",
		senderId: 2,
		sentAt: new Date(Date.now() - 1000 * 60 * 12).toISOString(),
		editedAt: null,
		read: true,
	},
	{
		id: 2,
		content: "Welcome to the room!",
		senderName: "Bob Smith",
		senderId: 3,
		sentAt: new Date(Date.now() - 1000 * 60 * 10).toISOString(),
		editedAt: null,
		read: true,
	},
	{
		id: 3,
		content: "This one was edited after sending.",
		senderName: "Alice Johnson",
		senderId: 2,
		sentAt: new Date(Date.now() - 1000 * 60 * 8).toISOString(),
		editedAt: new Date(Date.now() - 1000 * 60 * 7).toISOString(),
		read: true,
	},
	{
		id: 4,
		content: "Looks clean in night mode.",
		senderName: "Charlie Brown",
		senderId: 4,
		sentAt: new Date(Date.now() - 1000 * 60 * 6).toISOString(),
		editedAt: null,
		read: false,
	},
	{
		id: 5,
		content: "Message sent but not read yet.",
		senderName: "You",
		senderId: 1, // assume currentUserId = 1
		sentAt: new Date(Date.now() - 1000 * 60 * 4).toISOString(),
		editedAt: null,
		read: false,
	},
	{
		id: 6,
		content: "This one is read by everyone.",
		senderName: "You",
		senderId: 1,
		sentAt: new Date(Date.now() - 1000 * 60 * 3).toISOString(),
		editedAt: null,
		read: true,
	},
	{
		id: 0, // temporary ID before server assigns real one
		content: "Sending this right now...",
		senderName: "You",
		senderId: 1,
		sentAt: new Date().toISOString(),
		editedAt: null,
		nonce: "temp-12345",
		read: false,
	},
];

export interface MessagesState {
	messagesPerRoom: Record<number, MessageStore>;

	loading: boolean;
	error: string | null;

	getMessage: (roomId: number | string) => MessageStore;
	fetchMessages: (roomId: number) => Promise<void>;
	reset: () => void;
}

const emptyMessageStore: MessageStore = {
	messages: [],
	cursor: null,
	hasMore: true,
};

const useMessagesStore = create<MessagesState>()((set, get) => ({
	messagesPerRoom: {},
	loading: false,
	error: null,

	reset: () => {
		set({ messagesPerRoom: {}, loading: false, error: null });
	},

	getMessage: (roomId) =>
		get().messagesPerRoom[Number(roomId)] || emptyMessageStore,

	fetchMessages: async (roomId) => {
		const { loading } = get();
		const current = get().messagesPerRoom[roomId] || emptyMessageStore;

		if (loading || !current.hasMore) return;

		try {
			set({ error: null, loading: true });

			const { messages, nextCursor } = await messageService.getHistory({
				limit: 50,
				roomId,
				cursor: current.cursor,
			});

			set((state) => ({
				messagesPerRoom: {
					...state.messagesPerRoom,
					[roomId]: {
						messages: [
							...messages,
							...(state.messagesPerRoom[roomId]?.messages ?? []),
						],
						cursor: nextCursor,
						hasMore: nextCursor !== null,
					},
				},
				loading: false,
			}));
		} catch (err) {
			set({ error: getErrorMessage(err), loading: false });
		}
	},
}));

export default useMessagesStore;
