import {
	differenceInHours,
	format,
	formatDistanceToNow,
	isSameDay as isSameDayFn,
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

export function formatDateSeparator(dateString: string) {
	const date = new Date(dateString);
	const now = new Date();

	// Today → "Today"
	if (isToday(date)) {
		return "Today";
	}

	// Yesterday → "Yesterday"
	if (isYesterday(date)) {
		return "Yesterday";
	}

	// This week → Day name
	if (differenceInHours(now, date) < 168) {
		// 7 days
		return format(date, "EEEE");
	}

	// Older → "12 January 2026"
	return format(date, "dd MMMM yyyy");
}

export function formatLastSeen(dateString: string) {
	const date = new Date(dateString);
	const now = new Date();

	// Less than a minute ago → "just now"
	if (differenceInHours(now, date) === 0) {
		return "last seen recently";
	}

	// Today → "last seen at 14:32"
	if (isToday(date)) {
		return `last seen at ${format(date, "HH:mm")}`;
	}

	// Yesterday → "last seen yesterday at 14:32"
	if (isYesterday(date)) {
		return `last seen yesterday at ${format(date, "HH:mm")}`;
	}

	// This week → "last seen on Monday"
	if (differenceInHours(now, date) < 168) {
		return `last seen on ${format(date, "EEEE")}`;
	}

	// Older → "last seen on 12 January 2026"
	return `last seen on ${format(date, "dd MMMM yyyy")}`;
}

export function isSameDay(date1: string, date2: string): boolean {
	return isSameDayFn(new Date(date1), new Date(date2));
}
