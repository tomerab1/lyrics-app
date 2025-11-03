import { useEffect, useState } from "react";
import { Card } from "../Card";
import { API_BASE } from "../../constants";
import type { Song } from "../../types";

async function json<T>(res: Response): Promise<T> {
	if (!res.ok) {
		const text = await res.text().catch(() => "");
		throw new Error(text || `HTTP ${res.status}`);
	}
	return res.json();
}

export function SongsList() {
	const [songs, setSongs] = useState<Song[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const refresh = async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await fetch(`${API_BASE}/songs`);
			const data = await json<{ data: Song[] }>(res);
			setSongs(data.data);
		} catch (err: any) {
			setError(err.message);
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		refresh();
	}, []);

	return (
		<Card title="Songs (id, title)">
			<div className="flex items-center gap-2 mb-2">
				<button
					onClick={refresh}
					className="px-3 py-1 rounded-lg bg-white border border-gray-300 hover:bg-gray-50"
				>
					Refresh
				</button>
				{loading && <span className="text-xs text-gray-500">Loading...</span>}
			</div>
			{error && <p className="text-sm text-red-700 mb-2">{error}</p>}
			<ul className="divide-y">
				{songs.map((s) => (
					<li key={s.id} className="py-2 flex items-center justify-between">
						<span className="font-mono text-xs md:text-sm">{s.id}</span>
						<span> </span>
						<span className="text-sm md:text-base">{s.title}</span>
					</li>
				))}
				{songs.length === 0 && !loading && (
					<li className="py-2 text-sm text-gray-500">No songs yet.</li>
				)}
			</ul>
		</Card>
	);
}
