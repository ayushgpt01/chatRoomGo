import { Link } from "@tanstack/react-router";
import { Virtuoso } from "react-virtuoso";
import useRoomStore from "@/stores/roomStore";

const RoomsFooter = () => {
	const isFetching = useRoomStore((s) => s.isFetching);

	if (!isFetching) return null;

	return <div className="p-4 text-center text-sm opacity-60">Loading...</div>;
};

export default function RoomsSidebar() {
	const roomsList = useRoomStore((s) => s.roomsList);
	const getRooms = useRoomStore((s) => s.getRooms);

	return (
		<div className="w-64 bg-base-100 border-r flex flex-col">
			{/* Header */}
			<div className="p-4 h-16 border-b font-semibold text-lg">Rooms</div>

			{/* Room List */}
			<Virtuoso
				className="flex-1"
				data={roomsList}
				endReached={() => getRooms()}
				computeItemKey={(_, room) => room.id}
				itemContent={(_, room) => {
					return (
						<Link
							to="/rooms/$roomId"
							params={{ roomId: String(room.id) }}
							className="w-full block px-4 py-3 transition-colors hover:bg-base-200"
							activeProps={{
								className:
									"w-full block px-4 py-3 bg-primary text-primary-content hover:text-white",
							}}
						>
							<div className="font-medium truncate">{room.name}</div>

							<div className="text-xs opacity-60 truncate">
								Room ID: {room.id}
							</div>
						</Link>
					);
				}}
				components={{ Footer: RoomsFooter }}
			/>
		</div>
	);
}
