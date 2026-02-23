import { nanoid } from "nanoid";
import { create } from "zustand";

export type ToastType = "success" | "error" | "info";

export interface Toast {
	id: string;
	message: string;
	type: ToastType;
}

interface ToastState {
	toasts: Toast[];
	show: (message: string, type?: ToastType) => void;
	remove: (id: string) => void;
}

const useToastStore = create<ToastState>((set) => ({
	toasts: [],
	show: (message, type = "info") =>
		set((state) => ({
			toasts: [...state.toasts, { id: nanoid(), message, type }],
		})),

	remove: (id) =>
		set((state) => ({
			toasts: state.toasts.filter((t) => t.id !== id),
		})),
}));

export default useToastStore;
