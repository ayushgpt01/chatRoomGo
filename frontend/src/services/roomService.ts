import axiosClient from "@/integrations/axios/axiosClient";
import { UserSchema } from "@/types/auth";
import { type Room, RoomSchema } from "@/types/room";
import type { LoginResponse } from "./authService";

type JoinRoomResponse = {
	room: Room;
	login?: LoginResponse;
};

export const roomService = {
	join: async (roomId: number) => {
		const response = await axiosClient.post<JoinRoomResponse>("/room/join", {
			roomId,
		});

		const data = response.data;

		// Returns an anonymous user if no user id present
		if (data.login) {
			UserSchema.parse(data.login.user);
		}

		RoomSchema.parse(data.room);

		return data;
	},

	leave: async (roomId: number) => {
		await axiosClient.post("/room/leave", { roomId });
	},
};
