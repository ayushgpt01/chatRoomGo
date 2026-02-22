import { z } from "zod";

export const UserSchema = z.object({
	id: z.string(),
	email: z.email(),
	name: z.string().min(2),
	role: z.enum(["admin", "user"]),
});

export type User = z.infer<typeof UserSchema>;

export const LoginSchema = z.object({
	email: z.email("Invalid email address"),
	password: z.string().min(6, "Password must be at least 6 characters"),
});

export type LoginCredentials = z.infer<typeof LoginSchema>;

export const SignupSchema = z.object({
	name: z.string().min(2, "Name must be at least 2 characters"),
	email: z.email("Invalid email address"),
	password: z.string().min(6, "Password must be at least 6 characters"),
});

export type SignupCredentials = z.infer<typeof SignupSchema>;
