import { z } from "zod";

export const MessageSchema = z.object({
	id: z.coerce.number(),
	content: z.string(),
	senderName: z.string(),
	senderId: z.number(),
	sentAt: z.string(),
	editedAt: z.string().nullable(),
	nonce: z.string().optional(),
	read: z.boolean().default(false),
});

export type Message = z.infer<typeof MessageSchema>;

export const GetMessageSchema = z.object({
	limit: z.coerce.number().min(1).max(100).default(50),
	roomId: z.coerce.number(),
	cursor: z.string().nullable(),
});

export type GetMessages = z.infer<typeof GetMessageSchema>;
