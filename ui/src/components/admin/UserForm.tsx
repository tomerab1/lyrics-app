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

export function UserForm() {
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const [createdId, setCreatedId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setCreatedId(null);
    try {
      const res = await fetch(`${API_BASE}/users`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name }),
      });
      const data = await json<{ data: { id: string } }>(res);
      setCreatedId(data.data.id);
      setName("");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="Create User">
      <form onSubmit={onSubmit} className="space-y-3">
        <input className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500" placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} required />
        <button type="submit" className="px-4 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white shadow focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50" disabled={loading}>
          {loading ? "Creating..." : "Create User"}
        </button>
      </form>
      {createdId && (
        <p className="text-sm text-green-700 mt-2">
          Created user id: <span className="font-mono">{createdId}</span>
        </p>
      )}
      {error && <p className="text-sm text-red-700 mt-2">{error}</p>}
    </Card>
  );
}


