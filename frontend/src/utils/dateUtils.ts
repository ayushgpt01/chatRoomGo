import {
	differenceInHours,
	format,
	formatDistanceToNow,
	isToday,
	isYesterday,
} from "date-fns";

export function formatChatTime(dateString: string) {
	const date = new Date(dateString);
	const now = new Date();

	const hoursDiff = differenceInHours(now, date);

	// Within last 6 hours → "2 hours ago"
	if (hoursDiff < 6) {
		return formatDistanceToNow(date, { addSuffix: true });
	}

	// Today → 14:32
	if (isToday(date)) {
		return format(date, "HH:mm");
	}

	// Yesterday → Yesterday at 14:32
	if (isYesterday(date)) {
		return `Yesterday at ${format(date, "HH:mm")}`;
	}

	// Older → 12 Jan 2026, 14:32
	return format(date, "dd MMM yyyy, HH:mm");
}
