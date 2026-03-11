import { useEffect } from "react";
import useAuthStore from "@/stores/authStore";
import useMessagesStore from "@/stores/messagesStore";
import useSocketStore from "@/stores/socketStore";
import { useTypingStore } from "@/stores/typingStore";
import {
	type ErrorEvent,
	IncomingEventTypes,
	type MessageCreatedEvent,
	type MessageDeletedEvent,
	type MessageUpdatedEvent,
	type UserStartedTypingEvent,
	type UserStoppedTypingEvent,
} from "@/types/events";

export default function useSocketEvents(roomId: number) {
	const userId = useAuthStore((s) => s.user?.id);
	const subscribe = useSocketStore((s) => s.subscribe);
	const connect = useSocketStore((s) => s.connect);
	const disconnect = useSocketStore((s) => s.disconnect);
	const handleTypingEvent = useTypingStore((s) => s.handleTypingEvent);

	useEffect(() => {
		if (!userId || !roomId) return;

		const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
		const wsUrl = `${protocol}//${
			import.meta.env.VITE_ENV === "development"
				? import.meta.env.VITE_WS_URL
				: window.location.host
		}/ws?room=${roomId}&user=${userId}`;

		connect(wsUrl);

		return () => {
			disconnect();
		};
	}, [connect, disconnect, roomId, userId]);

	useEffect(() => {
		const unsubCreated = subscribe(
			IncomingEventTypes.EventMessageCreated,
			(event: MessageCreatedEvent) => {
				const msg = event.payload.message;
				if (msg.roomId !== roomId) return;

				// Use fresh state from the store instead of stale state
				const currentStore = useMessagesStore.getState();
				const existing = currentStore
					.getMessage(roomId)
					.messages.some((m) => m.id === msg.id);

				if (existing) return;

				currentStore.upsertMessage(roomId, msg);
			},
		);

		const unsubUpdated = subscribe(
			IncomingEventTypes.EventMessageUpdated,
			(event: MessageUpdatedEvent) => {
				const msg = event.payload.message;
				if (msg.roomId !== roomId) return;

				const currentStore = useMessagesStore.getState();
				const current = currentStore
					.getMessage(roomId)
					.messages.find((m) => m.id === msg.id);

				if (!current) return;

				if (current.content === msg.content) return;

				currentStore.upsertMessage(roomId, msg);
			},
		);

		const unsubDeleted = subscribe(
			IncomingEventTypes.EventMessageDeleted,
			(event: MessageDeletedEvent) => {
				if (event.payload.roomId !== roomId) return;

				const currentStore = useMessagesStore.getState();
				const exists = currentStore
					.getMessage(roomId)
					.messages.some((m) => m.id === event.payload.messageId);

				if (!exists) return;

				currentStore.removeMessage(roomId, event.payload.messageId);
			},
		);

		const unsubError = subscribe(
			IncomingEventTypes.EventError,
			(event: ErrorEvent) => {
				console.error("Socket error:", event.payload);
			},
		);

		const unsubStartedTyping = subscribe(
			IncomingEventTypes.EventUserStartedTyping,
			(event: UserStartedTypingEvent) => {
				if (event.payload.roomId !== roomId) return;
				if (event.payload.userId === userId) return;
				handleTypingEvent({
					type: "user_started_typing",
					payload: event.payload,
				});
			},
		);

		const unsubStoppedTyping = subscribe(
			IncomingEventTypes.EventUserStoppedTyping,
			(event: UserStoppedTypingEvent) => {
				if (event.payload.roomId !== roomId) return;
				if (event.payload.userId === userId) return;
				handleTypingEvent({
					type: "user_stopped_typing",
					payload: event.payload,
				});
			},
		);

		return () => {
			unsubCreated();
			unsubUpdated();
			unsubDeleted();
			unsubError();
			unsubStartedTyping();
			unsubStoppedTyping();
		};
	}, [roomId, subscribe, handleTypingEvent, userId]);
}
