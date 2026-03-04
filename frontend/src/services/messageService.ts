import { z } from "zod";
import axiosClient from "@/integrations/axios/axiosClient";
import {
	GetMessageSchema,
	type GetMessages,
	type Message,
	MessageSchema,
} from "@/types/message";

interface GetMessagesResponse {
	messages: Message[];
	nextCursor: string | null;
}

const SendMessageSchema = z.object({
	roomId: z.coerce.number(),
	content: z.string().min(1),
});

const EditMessageSchema = z.object({
	roomId: z.coerce.number(),
	messageId: z.coerce.number(),
	content: z.string().min(1),
});

const DeleteMessageSchema = z.object({
	roomId: z.coerce.number(),
	messageId: z.coerce.number(),
});

const messageService = {
	getHistory: async (payload: GetMessages) => {
		const { cursor, limit, roomId } = GetMessageSchema.parse(payload);

		const response = await axiosClient.get<GetMessagesResponse>(
			`/room/${roomId}/messages`,
			{ params: { cursor, limit } },
		);

		const data = response.data;

		z.array(MessageSchema).parse(data.messages);
		data.messages.reverse();

		return data;
	},

	sendMessage: async (payload: z.infer<typeof SendMessageSchema>) => {
		const { roomId, content } = SendMessageSchema.parse(payload);

		const response = await axiosClient.post<Message>(
			`/room/${roomId}/messages`,
			{
				content,
			},
		);

		return MessageSchema.parse(response.data);
	},

	editMessage: async (payload: z.infer<typeof EditMessageSchema>) => {
		const { roomId, messageId, content } = EditMessageSchema.parse(payload);

		const response = await axiosClient.patch<Message>(
			`/room/${roomId}/messages/${messageId}`,
			{ content },
		);

		return MessageSchema.parse(response.data);
	},

	deleteMessage: async (payload: z.infer<typeof DeleteMessageSchema>) => {
		const { roomId, messageId } = DeleteMessageSchema.parse(payload);

		await axiosClient.delete(`/room/${roomId}/messages/${messageId}`);
	},
};

export default messageService;
