import { z } from "zod";

export const RoomSchema = z.object({
	id: z.coerce.number(),
	name: z.string(),
});

export type Room = z.infer<typeof RoomSchema>;

export const GetRoomsSchema = z.object({
	limit: z.coerce.number().min(1).max(100).default(50),
	cursor: z.string().nullable(),
});

export type GetRooms = z.infer<typeof GetRoomsSchema>;
