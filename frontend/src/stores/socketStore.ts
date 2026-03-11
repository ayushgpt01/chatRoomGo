import { create } from "zustand";
import type { IncomingSocketEvent } from "@/types/events";

type Listener<T extends IncomingSocketEvent["type"]> = (
	event: Extract<IncomingSocketEvent, { type: T }>,
) => void;

interface SocketState {
	socket: WebSocket | null;
	status: "connecting" | "open" | "closed";
	error: string | null;
	currentUrl: string | null;
	connect: (url: string) => void;
	disconnect: () => void;
	send: (data: unknown) => void;
	subscribe: <T extends IncomingSocketEvent["type"]>(
		type: T,
		listener: Listener<T>,
	) => () => void;
}

// biome-ignore lint/suspicious/noExplicitAny: the 'subscribe' method provides the type safety at the boundary.
const listeners = new Map<IncomingSocketEvent["type"], Set<Listener<any>>>();

let reconnectTimeout: ReturnType<typeof setTimeout>;
let retryCount = 0;

const CLOSE_CODE = 1000;
const MAX_RECONNECT_DELAY = 30000;
const messageQueue: string[] = [];
const MAX_QUEUE_SIZE = 100;

const useSocketStore = create<SocketState>((set, get) => ({
	socket: null,
	status: "closed" as const,
	error: null,
	currentUrl: null,

	connect: (url) => {
		const { socket: current, status, currentUrl } = get();
		if (
			currentUrl === url &&
			(current?.readyState === WebSocket.OPEN ||
				current?.readyState === WebSocket.CONNECTING ||
				status === "connecting")
		) {
			return;
		}

		const socket = new WebSocket(url);
		set({ socket, status: "connecting", error: null, currentUrl: url });

		socket.onopen = () => {
			console.log("WS Connected", url);
			set({ status: "open" });
			retryCount = 0;
			clearTimeout(reconnectTimeout);
			while (messageQueue.length > 0) {
				const msg = messageQueue.shift();
				if (msg) socket.send(msg);
			}
		};

		socket.onmessage = (event) => {
			console.log("WS Event recieved", event);
			try {
				const parsed: IncomingSocketEvent = JSON.parse(event.data);
				const eventListeners = listeners.get(parsed.type);
				eventListeners?.forEach((callback) => {
					callback(parsed);
				});
			} catch (e) {
				console.error("Invalid WS message", e);
			}
		};

		socket.onclose = (event) => {
			if (get().socket === socket) {
				set({ socket: null, status: "closed" });
			}

			if (event.code !== CLOSE_CODE) {
				const delay = Math.min(1000 * 2 ** retryCount, MAX_RECONNECT_DELAY);
				reconnectTimeout = setTimeout(() => {
					retryCount++;
					get().connect(url);
				}, delay);
			}
		};

		socket.onerror = () => set({ error: "WebSocket connection failed" });
	},

	disconnect: () => {
		const { socket } = get();
		clearTimeout(reconnectTimeout);
		messageQueue.length = 0;
		console.log("WS Disconnected", socket?.url);
		socket?.close(CLOSE_CODE);

		if (get().socket === socket) {
			set({ socket: null, status: "closed", currentUrl: null });
		}
	},

	send: (data) => {
		const { socket, status } = get();
		const payload = JSON.stringify(data);

		if (socket && status === "open") {
			socket.send(payload);
		} else {
			if (messageQueue.length >= MAX_QUEUE_SIZE) {
				messageQueue.shift();
			}
			messageQueue.push(payload);
		}
	},

	subscribe: <T extends IncomingSocketEvent["type"]>(
		type: T,
		listener: Listener<T>,
	) => {
		if (!listeners.has(type)) {
			listeners.set(type, new Set());
		}

		const setOfListeners = listeners.get(type);
		if (!setOfListeners) {
			return () => {};
		}

		setOfListeners.add(listener);

		return () => {
			setOfListeners.delete(listener);
			if (setOfListeners.size === 0) {
				listeners.delete(type);
			}
		};
	},
}));

export default useSocketStore;
