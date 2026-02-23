import { z } from "zod";

export const UserSchema = z.object({
	id: z.number(),
	username: z.string().max(255),
	name: z.string().min(2),
});

export type User = z.infer<typeof UserSchema>;

export const LoginSchema = z.object({
	username: z
		.string()
		.max(255, "Username cannot be longer than 255 characters")
		.nonempty("Username is required"),
	password: z.string().min(6, "Password must be at least 6 characters"),
});

export type LoginCredentials = z.infer<typeof LoginSchema>;

export const SignupSchema = z
	.object({
		name: z.string().min(2, "Name must be at least 2 characters"),
		username: z
			.string()
			.max(255, "Username cannot be longer than 255 characters")
			.nonempty("Username is required"),
		password: z.string().min(6, "Password must be at least 6 characters"),
		confirmPassword: z.string(),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "Passwords don't match",
		path: ["confirmPassword"],
	});

export type SignupCredentials = z.infer<typeof SignupSchema>;
