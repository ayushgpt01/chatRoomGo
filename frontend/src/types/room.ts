import { z } from "zod";
import { UserSchema } from "./auth";

export const RoomSchema = z.object({
	id: z.coerce.number(),
	name: z.string(),
	participantCount: z.number().default(1),
	updatedAt: z.string(),
	members: z.array(UserSchema).optional(),
});

export type Room = z.infer<typeof RoomSchema>;

export const GetRoomsSchema = z.object({
	limit: z.coerce.number().min(1).max(100).default(50),
	cursor: z.string().nullable(),
});

export type GetRooms = z.infer<typeof GetRoomsSchema>;
