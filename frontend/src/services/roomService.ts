import axiosClient from "@/integrations/axios/axiosClient";
import { UserSchema } from "@/types/auth";
import { type Room, RoomSchema } from "@/types/room";
import type { LoginResponse } from "./authService";

type JoinRoomResponse = {
	room: Room;
	login?: LoginResponse;
};

type CreateRoomResponse = {
	room: Room;
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

	create: async (roomName: string) => {
		const response = await axiosClient.post<CreateRoomResponse>(
			"/room/create",
			{ name: roomName },
		);

		const data = response.data;
		RoomSchema.parse(data.room);
		return data;
	},
};
