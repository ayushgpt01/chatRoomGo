import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { LoginResponse } from "@/services/authService";
import { roomService } from "@/services/roomService";
import type { Room } from "@/types/room";
import { getErrorMessage } from "@/utils/errorHandler";

export interface RoomState {
	room: Room | null;
	isJoining: boolean;
	isLeaving: boolean;
	isCreating: boolean;
	error: string | null;
	create: (roomName: string) => Promise<number>;
	join: (roomId: number) => Promise<{ login: LoginResponse | undefined }>;
	leave: () => Promise<void>;
}

const useRoomStore = create<RoomState>()(
	persist(
		(set, get) => ({
			room: null,
			isJoining: false,
			isLeaving: false,
			isCreating: false,
			error: null,

			join: async (payload) => {
				set({ isJoining: true, error: null });
				try {
					const { room, login } = await roomService.join(payload);
					set({ isJoining: false, room });
					return { login };
				} catch (err) {
					set({ error: getErrorMessage(err), isJoining: false });
					throw err;
				}
			},
			leave: async () => {
				const room = get().room;
				if (!room) return;
				set({ isLeaving: true, error: null });
				try {
					await roomService.leave(room.id);

					set({ isLeaving: false, room: null });
				} catch (err) {
					set({
						error: getErrorMessage(err),
						isLeaving: false,
					});
					throw err;
				}
			},
			create: async (payload) => {
				set({ isCreating: true, error: null });
				try {
					const res = await roomService.create(payload);
					set({ isCreating: false, room: res.room });

					return res.room.id;
				} catch (err) {
					set({
						error: getErrorMessage(err),
						isLeaving: false,
					});
					throw err;
				}
			},
		}),
		{
			name: "room-storage",
			partialize: (state) => ({ room: state.room }),
		},
	),
);

export default useRoomStore;
