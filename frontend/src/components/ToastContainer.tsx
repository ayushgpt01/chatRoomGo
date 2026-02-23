import { useEffect } from "react";
import useToastStore from "@/stores/toastStore";

export default function ToastContainer() {
	const { toasts, remove } = useToastStore();

	useEffect(() => {
		if (toasts.length === 0) return;

		const timers = toasts.map((toast) =>
			setTimeout(() => remove(toast.id), 3000),
		);

		return () => timers.forEach(clearTimeout);
	}, [toasts, remove]);

	return (
		<div className="toast toast-top toast-end z-50">
			{toasts.map((toast) => (
				<div
					key={toast.id}
					className={`alert shadow-lg ${
						toast.type === "success"
							? "alert-success"
							: toast.type === "error"
								? "alert-error"
								: "alert-info"
					}`}
				>
					<span>{toast.message}</span>
				</div>
			))}
		</div>
	);
}
