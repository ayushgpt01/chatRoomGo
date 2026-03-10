import { create } from "zustand";
import type { TypingEvent } from "@/types/message";

interface TypingState {
	// Map of roomId to Set of typing user names
	typingUsers: Record<number, Set<string>>;

	// Actions
	addTypingUser: (roomId: number, userName: string) => void;
	removeTypingUser: (roomId: number, userName: string) => void;
	getTypingUsers: (roomId: number) => string[];
	clearRoomTyping: (roomId: number) => void;
	handleTypingEvent: (event: TypingEvent) => void;
}

export const useTypingStore = create<TypingState>((set, get) => ({
	typingUsers: {},

	addTypingUser: (roomId: number, userName: string) => {
		set((state) => {
			const currentTyping = state.typingUsers[roomId] || new Set();
			const newTyping = new Set(currentTyping);
			newTyping.add(userName);

			return {
				typingUsers: {
					...state.typingUsers,
					[roomId]: newTyping,
				},
			};
		});
	},

	removeTypingUser: (roomId: number, userName: string) => {
		set((state) => {
			const currentTyping = state.typingUsers[roomId] || new Set();
			const newTyping = new Set(currentTyping);
			newTyping.delete(userName);

			// Clean up empty sets to save memory
			if (newTyping.size === 0) {
				const { [roomId]: _, ...rest } = state.typingUsers;
				return { typingUsers: rest };
			}

			return {
				typingUsers: {
					...state.typingUsers,
					[roomId]: newTyping,
				},
			};
		});
	},

	getTypingUsers: (roomId: number) => {
		return Array.from(get().typingUsers[roomId] || []);
	},

	clearRoomTyping: (roomId: number) => {
		set((state) => {
			const { [roomId]: _, ...rest } = state.typingUsers;
			return { typingUsers: rest };
		});
	},

	handleTypingEvent: (event: TypingEvent) => {
		const { type, payload } = event;
		const { roomId, userName } = payload;

		if (type === "user_started_typing") {
			get().addTypingUser(roomId, userName);
		} else if (type === "user_stopped_typing") {
			get().removeTypingUser(roomId, userName);
		}
	},
}));

// Auto-clear typing after 3 seconds of inactivity
export const setupTypingTimeouts = () => {
	const timeouts: Record<string, NodeJS.Timeout> = {};

	return {
		setTypingTimeout: (
			roomId: number,
			userName: string,
			callback: () => void,
		) => {
			const key = `${roomId}-${userName}`;

			// Clear existing timeout for this user
			if (timeouts[key]) {
				clearTimeout(timeouts[key]);
			}

			// Set new timeout
			timeouts[key] = setTimeout(() => {
				callback();
				delete timeouts[key];
			}, 3000);
		},

		clearTypingTimeout: (roomId: number, userName: string) => {
			const key = `${roomId}-${userName}`;
			if (timeouts[key]) {
				clearTimeout(timeouts[key]);
				delete timeouts[key];
			}
		},
	};
};
