import { UserForm } from "./admin/UserForm";
import { SongForm } from "./admin/SongForm";
import { UsersList } from "./admin/UsersList";
import { SongsList } from "./admin/SongsList";

export function AdminPanel() {
	return (
		<div className="grid gap-6 md:grid-cols-2">
			<UserForm />
			<SongForm />
			<UsersList />
			<SongsList />
		</div>
	);
}
