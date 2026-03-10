import { z } from "zod";

export const MessageSchema = z.object({
	id: z.coerce.number(),
	roomId: z.coerce.number(),
	content: z.string(),
	senderName: z.string(),
	senderId: z.number(),
	sentAt: z.string(),
	editedAt: z.string().nullable(),
	nonce: z.string().nullable().optional(),
	delivered: z.boolean().default(true),
	readBy: z.array(z.number()).nullable().optional(),
});

export type Message = z.infer<typeof MessageSchema>;

export const GetMessageSchema = z.object({
	limit: z.coerce.number().min(1).max(100).default(50),
	roomId: z.coerce.number(),
	cursor: z.string().nullable(),
});

export type GetMessages = z.infer<typeof GetMessageSchema>;

// Typing event types
export const TypingEventSchema = z.object({
	type: z.enum(["user_started_typing", "user_stopped_typing"]),
	payload: z.object({
		roomId: z.number(),
		userId: z.number(),
		userName: z.string(),
	}),
});

export type TypingEvent = z.infer<typeof TypingEventSchema>;

export const TypingPayloadSchema = z.object({
	roomId: z.number(),
});

export type TypingPayload = z.infer<typeof TypingPayloadSchema>;
