import type { Message } from "./message";

// ---- BASE TYPES ---------

export enum OutgoingEventTypes {
	EventSendMessage = "message.send",
	EventEditMessage = "message.edit",
	EventDeleteMessage = "message.delete",
}

export enum IncomingEventTypes {
	EventMessageCreated = "message_created",
	EventMessageUpdated = "message_updated",
	EventMessageDeleted = "message_deleted",
	EventUserJoinedRoom = "user_joined_room",
	EventUserLeftRoom = "user_left_room",
	EventError = "error",
}

export type ServerEvent<T extends IncomingEventTypes, P> = {
	type: T;
	payload: P;
};

export type ClientEvent<T extends OutgoingEventTypes, D> = {
	type: T;
	data: D;
};

// ---- Server TYPES ---------

export type MessageCreatedEvent = ServerEvent<
	IncomingEventTypes.EventMessageCreated,
	Message
>;

export type MessageUpdatedEvent = ServerEvent<
	IncomingEventTypes.EventMessageUpdated,
	Message
>;

export type MessageDeletedEvent = ServerEvent<
	IncomingEventTypes.EventMessageDeleted,
	{
		messageId: number;
		roomId: number;
	}
>;

export type UserJoinedRoomEvent = ServerEvent<
	IncomingEventTypes.EventUserJoinedRoom,
	{
		roomId: number;
		userId: number;
	}
>;

export type UserLeftRoomEvent = ServerEvent<
	IncomingEventTypes.EventUserLeftRoom,
	{
		roomId: number;
		userId: number;
	}
>;

export enum ErrorCodes {
	InvalidPayload = "invalid_payload",
	NotRoomMember = "not_room_member",
	Forbidden = "forbidden",
	UnsupportedEvent = "unsupported_event",
	InternalError = "internal_error",
}

export type ErrorEvent = ServerEvent<
	IncomingEventTypes.EventError,
	{
		message: string;
		code: ErrorCodes;
	}
>;

export type IncomingSocketEvent =
	| MessageCreatedEvent
	| MessageUpdatedEvent
	| MessageDeletedEvent
	| UserJoinedRoomEvent
	| UserLeftRoomEvent
	| ErrorEvent;

// ---- Client TYPES ---------

export type SendMessageEvent = ClientEvent<
	OutgoingEventTypes.EventSendMessage,
	{
		content: string;
		nonce: string;
	}
>;

export type EditMessageEvent = ClientEvent<
	OutgoingEventTypes.EventEditMessage,
	{
		messageId: number;
		content: string;
	}
>;

export type DeleteMessageEvent = ClientEvent<
	OutgoingEventTypes.EventDeleteMessage,
	{
		messageId: number;
	}
>;

export type OutgoingSocketEvent =
	| SendMessageEvent
	| EditMessageEvent
	| DeleteMessageEvent;
