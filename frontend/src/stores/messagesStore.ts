import { create } from "zustand";
import messageService from "@/services/messageService";
import type { Message } from "@/types/message";
import { getErrorMessage } from "@/utils/errorHandler";
import useAuthStore from "./authStore";

interface MessageStore {
	messages: Message[];
	cursor: string | null;
	hasMore: boolean;
}

export interface MessagesState {
	messagesPerRoom: Record<number, MessageStore>;
	loading: boolean;
	error: string | null;

	getMessage: (roomId: number | string) => MessageStore;
	fetchMessages: (roomId: number) => Promise<void>;

	sendMessage: (roomId: number, content: string) => Promise<void>;
	editMessage: (
		roomId: number,
		messageId: number,
		content: string,
	) => Promise<void>;
	deleteMessage: (roomId: number, messageId: number) => Promise<void>;

	// internal helpers
	replaceMessageByNonce: (
		roomId: number,
		nonce: string,
		message: Message,
	) => void;
	upsertMessage: (roomId: number, message: Message) => void;
	removeMessage: (roomId: number, messageId: number) => void;

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

	reset: () => set({ messagesPerRoom: {}, loading: false, error: null }),

	getMessage: (roomId) =>
		get().messagesPerRoom[Number(roomId)] || emptyMessageStore,

	fetchMessages: async (roomId) => {
		const { loading } = get();
		const current = get().messagesPerRoom[roomId] || emptyMessageStore;
		if (loading || !current.hasMore) return;

		try {
			set({ loading: true, error: null });

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

	sendMessage: async (roomId, content) => {
		const user = useAuthStore.getState().user;
		if (!user) return;

		const nonce = crypto.randomUUID();

		const optimistic: Message = {
			id: 0,
			roomId,
			content,
			senderId: user.id,
			senderName: user.username,
			sentAt: new Date().toISOString(),
			read: false,
			editedAt: null,
			nonce,
		};

		// optimistic insert
		set((state) => ({
			messagesPerRoom: {
				...state.messagesPerRoom,
				[roomId]: {
					...state.messagesPerRoom[roomId],
					messages: [
						...(state.messagesPerRoom[roomId]?.messages ?? []),
						optimistic,
					],
				},
			},
		}));

		try {
			const created = await messageService.sendMessage({
				roomId,
				content,
			});

			// Replace optimistic via nonce
			get().replaceMessageByNonce(roomId, nonce, created);
		} catch (err) {
			// remove failed optimistic
			set((state) => ({
				messagesPerRoom: {
					...state.messagesPerRoom,
					[roomId]: {
						...state.messagesPerRoom[roomId],
						messages: state.messagesPerRoom[roomId].messages.filter(
							(m) => m.nonce !== nonce,
						),
					},
				},
			}));

			set({ error: getErrorMessage(err) });
		}
	},

	editMessage: async (roomId, messageId, content) => {
		const current =
			get()
				.getMessage(roomId)
				.messages.find((m) => m.id === messageId) || null;

		if (!current) return;

		const previous = { ...current };

		// optimistic update
		get().upsertMessage(roomId, {
			...current,
			content,
			editedAt: new Date().toISOString(),
		});

		try {
			const updated = await messageService.editMessage({
				roomId,
				messageId,
				content,
			});

			get().upsertMessage(roomId, updated);
		} catch (err) {
			// revert
			get().upsertMessage(roomId, previous);
			set({ error: getErrorMessage(err) });
		}
	},

	deleteMessage: async (roomId, messageId) => {
		const current =
			get()
				.getMessage(roomId)
				.messages.find((m) => m.id === messageId) || null;

		if (!current) return;

		// optimistic remove
		get().removeMessage(roomId, messageId);

		try {
			await messageService.deleteMessage({
				roomId,
				messageId,
			});
		} catch (err) {
			// revert
			get().upsertMessage(roomId, current);
			set({ error: getErrorMessage(err) });
		}
	},

	replaceMessageByNonce: (roomId, nonce, message) => {
		set((state) => ({
			messagesPerRoom: {
				...state.messagesPerRoom,
				[roomId]: {
					...state.messagesPerRoom[roomId],
					messages: state.messagesPerRoom[roomId].messages.map((m) =>
						m.nonce === nonce ? message : m,
					),
				},
			},
		}));
	},

	upsertMessage: (roomId, message) => {
		set((state) => {
			const store = state.messagesPerRoom[roomId];
			if (!store) return state;

			const exists = store.messages.some((m) => m.id === message.id);

			return {
				messagesPerRoom: {
					...state.messagesPerRoom,
					[roomId]: {
						...store,
						messages: exists
							? store.messages.map((m) => (m.id === message.id ? message : m))
							: [...store.messages, message],
					},
				},
			};
		});
	},

	removeMessage: (roomId, messageId) => {
		set((state) => ({
			messagesPerRoom: {
				...state.messagesPerRoom,
				[roomId]: {
					...state.messagesPerRoom[roomId],
					messages: state.messagesPerRoom[roomId].messages.filter(
						(m) => m.id !== messageId,
					),
				},
			},
		}));
	},
}));

export default useMessagesStore;
