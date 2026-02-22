import { AxiosError } from "axios";
import { ZodError } from "zod";

export const getErrorMessage = (
	error: unknown,
	fallbackMessage = "Server connection failed",
): string => {
	if (error instanceof ZodError) {
		return error.issues[0].message;
	}
	if (error instanceof AxiosError) {
		return error.response?.data?.message || fallbackMessage;
	}
	return "An unknown error occurred";
};
