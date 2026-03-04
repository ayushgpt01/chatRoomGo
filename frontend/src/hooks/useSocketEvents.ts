import { useEffect } from "react";
import useAuthStore from "@/stores/authStore";
import useMessagesStore from "@/stores/messagesStore";
import useSocketStore from "@/stores/socketStore";
import {
	type ErrorEvent,
	IncomingEventTypes,
	type MessageCreatedEvent,
	type MessageDeletedEvent,
	type MessageUpdatedEvent,
} from "@/types/events";

export default function useSocketEvents(roomId: number) {
	const userId = useAuthStore((s) => s.user?.id);
	const subscribe = useSocketStore((s) => s.subscribe);
	const connect = useSocketStore((s) => s.connect);
	const disconnect = useSocketStore((s) => s.disconnect);

	useEffect(() => {
		if (!userId || !roomId) return;

		connect(
			`ws://${import.meta.env.VITE_WS_URL}/ws?room=${roomId}&user=${userId}`,
		);

		return () => disconnect();
	}, [connect, disconnect, roomId, userId]);

	useEffect(() => {
		const store = useMessagesStore.getState();

		const unsubCreated = subscribe(
			IncomingEventTypes.EventMessageCreated,
			(event: MessageCreatedEvent) => {
				const msg = event.payload;
				if (msg.roomId !== roomId) return;

				const existing = store
					.getMessage(roomId)
					.messages.some((m) => m.id === msg.id);

				if (existing) return;

				store.upsertMessage(roomId, msg);
			},
		);

		const unsubUpdated = subscribe(
			IncomingEventTypes.EventMessageUpdated,
			(event: MessageUpdatedEvent) => {
				const msg = event.payload;
				if (msg.roomId !== roomId) return;

				const current = store
					.getMessage(roomId)
					.messages.find((m) => m.id === msg.id);

				if (!current) return;

				if (current.content === msg.content) return;

				store.upsertMessage(roomId, msg);
			},
		);

		const unsubDeleted = subscribe(
			IncomingEventTypes.EventMessageDeleted,
			(event: MessageDeletedEvent) => {
				if (event.payload.roomId !== roomId) return;

				const exists = store
					.getMessage(roomId)
					.messages.some((m) => m.id === event.payload.messageId);

				if (!exists) return;

				store.removeMessage(roomId, event.payload.messageId);
			},
		);

		const unsubError = subscribe(
			IncomingEventTypes.EventError,
			(event: ErrorEvent) => {
				console.error("Socket error:", event.payload);
			},
		);

		return () => {
			unsubCreated();
			unsubUpdated();
			unsubDeleted();
			unsubError();
		};
	}, [roomId, subscribe]);
}
