import React, { useState } from "react";
import { Card } from "../Card";
import { API_BASE } from "../../constants";

async function json<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

export function SongForm() {
  const [title, setTitle] = useState("");
  const [artist, setArtist] = useState("");
  const [lyrics, setLyrics] = useState("");
  const [loading, setLoading] = useState(false);
  const [created, setCreated] = useState<{ id: string; lineCount: number } | null>(null);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setCreated(null);
    try {
      const res = await fetch(`${API_BASE}/songs`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, artist, lyrics }),
      });
      const data = await json<{ data: { id: string; lineCount: number } }>(res);
      setCreated(data.data);
      setTitle("");
      setArtist("");
      setLyrics("");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="Create Song">
      <form onSubmit={onSubmit} className="space-y-3">
        <input className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500" placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} required />
        <input className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500" placeholder="Artist" value={artist} onChange={(e) => setArtist(e.target.value)} required />
        <textarea className="w-full border border-gray-300 rounded-lg px-3 py-2 h-40 font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500" placeholder={"Paste lyrics here. One line per lyric line."} value={lyrics} onChange={(e) => setLyrics(e.target.value)} required />
        <button type="submit" className="px-4 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white shadow focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50" disabled={loading}>
          {loading ? "Creating..." : "Create Song"}
        </button>
      </form>
      {created && (
        <p className="text-sm text-green-700 mt-2">
          Created song <span className="font-mono">{created.id}</span> with {created.lineCount} lines.
        </p>
      )}
      {error && <p className="text-sm text-red-700 mt-2">{error}</p>}
    </Card>
  );
}


