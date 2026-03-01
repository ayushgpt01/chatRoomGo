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

const messageService = {
	getHistory: async (payload: GetMessages) => {
		const { cursor, limit, roomId } = GetMessageSchema.parse(payload);

		const response = await axiosClient.get<GetMessagesResponse>(
			`/rooms/${roomId}/messages`,
			{ params: { cursor, limit } },
		);

		const data = response.data;

		z.array(MessageSchema).parse(data.messages);

		return data;
	},
};

export default messageService;
