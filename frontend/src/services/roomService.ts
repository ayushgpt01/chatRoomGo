import { z } from "zod";
import axiosClient from "@/integrations/axios/axiosClient";
import { UserSchema } from "@/types/auth";
import {
	type GetRooms,
	GetRoomsSchema,
	type Room,
	RoomSchema,
} from "@/types/room";
import type { LoginResponse } from "./authService";

type JoinRoomResponse = {
	room: Room;
	login?: LoginResponse;
};

type CreateRoomResponse = {
	room: Room;
};

type GetRoomsResponse = {
	rooms: Room[];
	nextCursor: string | null;
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

	getRooms: async (payload: GetRooms) => {
		const params = GetRoomsSchema.parse(payload);

		const response = await axiosClient.get<GetRoomsResponse>("/room/getAll", {
			params,
		});

		const data = response.data;
		z.array(RoomSchema).parse(data.rooms);
		return data;
	},
};
