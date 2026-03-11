import { create } from "zustand";
import { useShallow } from "zustand/react/shallow";
import type { User } from "@/types/auth";
import type { TypingEvent } from "@/types/message";

interface TypingState {
	typingUsers: Record<number, number[]>;
	userMap: Record<number, { name: string; username: string }>;

	// Actions
	addTypingUser: (
		roomId: number,
		userId: number,
		userInfo: { name: string; username: string },
	) => void;
	removeTypingUser: (roomId: number, userId: number) => void;
	clearRoomTyping: (roomId: number) => void;
	handleTypingEvent: (event: TypingEvent) => void;
	updateUserMap: (users: User[]) => void;
}

export const useTypingStore = create<TypingState>((set, get) => ({
	typingUsers: {},
	userMap: {},

	addTypingUser: (roomId, userId, userInfo) => {
		set((state) => {
			const current = state.typingUsers[roomId] ?? [];

			if (current.includes(userId)) return state;

			return {
				typingUsers: {
					...state.typingUsers,
					[roomId]: [...current, userId],
				},
				userMap: {
					...state.userMap,
					[userId]: userInfo,
				},
			};
		});
	},

	removeTypingUser: (roomId, userId) => {
		set((state) => {
			const current = state.typingUsers[roomId];
			if (!current) return state;

			const next = current.filter((id) => id !== userId);

			if (next.length === 0) {
				const { [roomId]: _, ...rest } = state.typingUsers;
				return { typingUsers: rest };
			}

			return {
				typingUsers: {
					...state.typingUsers,
					[roomId]: next,
				},
			};
		});
	},

	clearRoomTyping: (roomId) => {
		set((state) => {
			const { [roomId]: _, ...rest } = state.typingUsers;
			return { typingUsers: rest };
		});
	},

	handleTypingEvent: (event) => {
		const { type, payload } = event;
		const { roomId, userId, userName } = payload;

		if (type === "user_started_typing") {
			get().addTypingUser(roomId, userId, {
				name: userName,
				username: userName,
			});
		}

		if (type === "user_stopped_typing") {
			get().removeTypingUser(roomId, userId);
		}
	},

	updateUserMap: (users) => {
		set((state) => {
			const next = { ...state.userMap };

			users.forEach((u) => {
				next[u.id] = {
					name: u.name,
					username: u.username,
				};
			});

			return { userMap: next };
		});
	},
}));

export function useTypingUsers(roomId: number) {
	return useTypingStore(
		useShallow((state) => {
			const ids = state.typingUsers[roomId] ?? [];
			return ids.map((id) => state.userMap[id]).filter(Boolean);
		}),
	);
}
