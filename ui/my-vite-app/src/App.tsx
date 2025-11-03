import React, { useEffect, useMemo, useState } from "react";

// --- Types matching your API ---

type User = { id: string; name: string };
type Song = { id: string; title: string };

type LessonItemFill = {
  type: "fillblanks";
  lineIndex: number;
  renderedLine: string;
  words: string[]; // exactly 4 options
  correct_word: string
};

type LessonItemArrange = {
  type: "arrange";
  lineIndex: number;
  words: string[]; // correct order; UI shuffles and compares
};

type LessonItem = LessonItemFill | LessonItemArrange;

type Lesson = {
  lessonId: string;
  items: LessonItem[];
};

type Summary = {
  total: number;
  correct: number;
  wrong: number;
  accuracy: number; // percentage, e.g., 66.7
  scheduledForRepractice: string[];
};

// --- Helpers ---

function classNames(...xs: Array<string | false | null | undefined>) {
  return xs.filter(Boolean).join(" ");
}

function shuffle<T>(arr: T[]): T[] {
  const a = [...arr];
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [a[i], a[j]] = [a[j], a[i]];
  }
  return a;
}

async function json<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

// --- UI Components ---

export default function App() {
  const [activeTab, setActiveTab] = useState<"admin" | "ui">("admin");

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-5xl">
        <header className="mb-6 text-center">
          <h1 className="text-2xl md:text-3xl font-semibold">Lyrics Practice Mini App</h1>
        </header>

        <nav className="mb-6 flex justify-center">
          <div className="inline-flex rounded-2xl overflow-hidden shadow">
            <button
              className={classNames(
                "px-4 py-2 text-sm md:text-base",
                activeTab === "admin" ? "bg-white" : "bg-gray-200"
              )}
              onClick={() => setActiveTab("admin")}
            >
              Admin
            </button>
            <button
              className={classNames(
                "px-4 py-2 text-sm md:text-base",
                activeTab === "ui" ? "bg-white" : "bg-gray-200"
              )}
              onClick={() => setActiveTab("ui")}
            >
              UI
            </button>
          </div>
        </nav>

        {activeTab === "admin" ? <AdminPanel /> : <LessonPanel />}
      </div>
    </div>
  );
}

// --- Admin Tab ---

function AdminPanel() {
  return (
    <div className="grid gap-6 md:grid-cols-2">
      <UserFormCard />
      <SongFormCard />
      <UsersListCard />
      <SongsListCard />
    </div>
  );
}

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="bg-white rounded-2xl shadow p-4 md:p-6">
      <h2 className="text-lg font-semibold mb-3 text-center">{title}</h2>
      {children}
    </section>
  );
}


function UserFormCard() {
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
      const res = await fetch("http://localhost:5555/api/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name }),
      });
      const data = await json<{ id: string }>(res);
      setCreatedId(data.id);
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
        <input
          className="w-full border rounded-lg px-3 py-2"
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <button
          type="submit"
          className="px-4 py-2 rounded-xl bg-gray-900 text-white disabled:opacity-50"
          disabled={loading}
        >
          {loading ? "Creating..." : "Create User"}
        </button>
      </form>
      {createdId && <p className="text-sm text-green-700 mt-2">Created user id: <span className="font-mono">{createdId}</span></p>}
      {error && <p className="text-sm text-red-700 mt-2">{error}</p>}
    </Card>
  );
}

function SongFormCard() {
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
      const res = await fetch("http://localhost:5555/api/songs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, artist, lyrics }),
      });
      const data = await json<{ id: string; lineCount: number }>(res);
      setCreated(data);
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
        <input
          className="w-full border rounded-lg px-3 py-2"
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />
        <input
          className="w-full border rounded-lg px-3 py-2"
          placeholder="Artist"
          value={artist}
          onChange={(e) => setArtist(e.target.value)}
          required
        />
        <textarea
          className="w-full border rounded-lg px-3 py-2 h-40 font-mono"
          placeholder={"Paste lyrics here. One line per lyric line."}
          value={lyrics}
          onChange={(e) => setLyrics(e.target.value)}
          required
        />
        <button
          type="submit"
          className="px-4 py-2 rounded-xl bg-gray-900 text-white disabled:opacity-50"
          disabled={loading}
        >
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

function UsersListCard() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refresh = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch("http://localhost:5555/api/users");
      const data = await json<{data: User[]}>(res);
      setUsers(data.data);
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
    <Card title="Users (id, name)">
      <div className="flex items-center gap-2 mb-2">
        <button onClick={refresh} className="px-3 py-1 rounded-lg bg-gray-200">Refresh</button>
        {loading && <span className="text-xs text-gray-500">Loading...</span>}
      </div>
      {error && <p className="text-sm text-red-700 mb-2">{error}</p>}
      <ul className="divide-y">
        {users.map((u) => (
          <li key={u.id} className="py-2 flex items-center justify-between">
            <span className="font-mono text-xs md:text-sm">{u.id}</span>
            <span>{' '}</span>
            <span className="text-sm md:text-base">{u.name}</span>
          </li>
        ))}
        {users.length === 0 && !loading && <li className="py-2 text-sm text-gray-500">No users yet.</li>}
      </ul>
    </Card>
  );
}

function SongsListCard() {
  const [songs, setSongs] = useState<Song[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refresh = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch("http://localhost:5555/api/songs");
      const data = await json<{data: Song[]}>(res);
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
        <button onClick={refresh} className="px-3 py-1 rounded-lg bg-gray-200">Refresh</button>
        {loading && <span className="text-xs text-gray-500">Loading...</span>}
      </div>
      {error && <p className="text-sm text-red-700 mb-2">{error}</p>}
      <ul className="divide-y">
        {songs.map((s) => (
          <li key={s.id} className="py-2 flex items-center justify-between">
            <span className="font-mono text-xs md:text-sm">{s.id}</span>
            <span>{' '}</span>
            <span className="text-sm md:text-base">{s.title}</span>
          </li>
        ))}
        {songs.length === 0 && !loading && <li className="py-2 text-sm text-gray-500">No songs yet.</li>}
      </ul>
    </Card>
  );
}

// --- UI Tab ---

function LessonPanel() {
  const [userId, setUserId] = useState("");
  const [busy, setBusy] = useState(false);
  const [lesson, setLesson] = useState<Lesson | null>(null);
  const [index, setIndex] = useState(0); // 0..5
  const [summary, setSummary] = useState<Summary | null>(null);
  const [error, setError] = useState<string | null>(null);

  const startLesson = async () => {
    setBusy(true);
    setError(null);
    setSummary(null);
    try {
      const res = await fetch("http://localhost:5555/api/lessons", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userId }),
      });
      const data = await json<{data: Lesson}>(res);
      setLesson(data.data);
      setIndex(0);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  };

  const onFinishedLesson = async () => {
    if (!lesson) return;
    try {
      const res = await fetch(`http://localhost:5555/api/lessons/${lesson.lessonId}/summary`);
      const data = await json<Summary>(res);
      setSummary(data);
    } catch (err: any) {
      setError(err.message);
    }
  };

  const onNext = async () => {
    if (!lesson) return;
    if (index < lesson.items.length - 1) {
      setIndex(index + 1);
    } else {
      await onFinishedLesson();
    }
  };

  const reset = () => {
    setLesson(null);
    setSummary(null);
    setIndex(0);
  };

  return (
    <div className="grid gap-6">
      <Card title="Start Lesson">
        <div className="flex flex-col md:flex-row gap-3 items-start md:items-end">
          <div className="flex-1 w-full">
            <label className="block text-sm mb-1">User ID</label>
            <input
              className="w-full border rounded-lg px-3 py-2 font-mono"
              placeholder="U123"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
            />
          </div>
          <button
            className="px-4 py-2 rounded-xl bg-gray-900 text-white disabled:opacity-50"
            onClick={startLesson}
            disabled={!userId || busy}
          >
            {busy ? "Starting..." : "Start Lesson"}
          </button>
        </div>
        {error && <p className="text-sm text-red-700 mt-2">{error}</p>}
      </Card>

      {lesson && !summary && (
        <LessonStepper
          lesson={lesson}
          index={index}
          onNext={onNext}
          fetchSummary={onFinishedLesson}
        />
      )}

      {lesson && summary && (
        <SummaryCard summary={summary} onRestart={reset} />
      )}
    </div>
  );
}

function LessonStepper({
  lesson,
  index,
  onNext,
  fetchSummary,
}: {
  lesson: Lesson;
  index: number;
  onNext: () => Promise<void> | void;
  fetchSummary: () => Promise<void> | void;
}) {
  const item = lesson.items[index];

  return (
    <Card title={`Step ${index + 1} of ${lesson.items.length}`}>
      {item.type === "fillblanks" ? (
        <FillBlanks item={item} lessonId={lesson.lessonId} itemIndex={index} onNext={onNext} />
      ) : (
        <Arrange item={item} lessonId={lesson.lessonId} itemIndex={index} onNext={onNext} />
      )}
    </Card>
  );
}

function FillBlanks({
  item,
  lessonId,
  itemIndex,
  onNext,
}: {
  item: LessonItemFill;
  lessonId: string;
  itemIndex: number;
  onNext: () => Promise<void> | void;
}) {
  const [submitting, setSubmitting] = useState(false);
  const [chosen, setChosen] = useState<string | null>(null);
  const [result, setResult] = useState<"correct" | "wrong" | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Determine the correct word by comparing renderedLine (with ___) to original words list
  // The server guarantees that exactly one word is hidden; it's among item.words.
  // We find which word would complete the blank by checking which candidate makes a valid line word split.
  const correctWord = useMemo(() => {
    // Heuristic: from options, pick the one that actually appears in the original four-option bank
    // Server already included the correct one. We can't reconstruct the original line here,
    // so assume the first one that was originally correct is unknown. Better: ask server to mark it.
    // BUT the spec states: words[] contains exactly 1 correct + 3 distractors.
    // We must not know which is correct; correctness comes from server evaluation via /answers.
    return null as unknown as string | null;
  }, [item.words]);

  const clickOption = async (w: string) => {
    if (submitting || chosen) return;
    setChosen(w);
    setSubmitting(true);
    setError(null);
    try {
      const res = await fetch("localhost:5555/answers", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          lessonId,
          itemIndex,
          type: "fillblanks",
          // The server should compute correctness and persist only if type === "fillblanks".
          // We still send a flag for convenience if your server accepts it, but it's optional.
          userInput: w,
        }),
      });

      if (res.status === 409) {
        // Duplicate submission for same item
        setError("Already answered. Moving on...");
        setTimeout(() => onNext(), 600);
        return;
      }

      // Some servers return { ok: true, correct: boolean } — handle both minimal and richer responses.
      let correct = false;
      try {
        const payload: any = await res.json();
        correct = !!payload.correct;
      } catch {
        // If no JSON returned, assume server checked and we cannot know; display neutral state.
      }

      setResult(correct ? "correct" : "wrong");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-4">
      <p className="text-lg">{item.renderedLine}</p>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        {item.words.map((w) => (
          <button
            key={w}
            disabled={!!chosen}
            onClick={() => clickOption(w)}
            className={classNames(
              "px-3 py-2 rounded-xl border",
              chosen === w && result === "correct" && "border-green-600 bg-green-50",
              chosen === w && result === "wrong" && "border-red-600 bg-red-50",
              !chosen && "hover:bg-gray-50"
            )}
          >
            {w}
          </button>
        ))}
      </div>
      {error && <p className="text-sm text-red-700">{error}</p>}
      {chosen && (
        <div className="flex items-center gap-3">
          <button
            onClick={onNext}
            className="px-4 py-2 rounded-xl bg-gray-900 text-white"
          >
            Next
          </button>
          {result && (
            <span className={classNames(
              "text-sm",
              result === "correct" ? "text-green-700" : "text-red-700"
            )}>
              {result === "correct" ? "Correct" : "Wrong"}
            </span>
          )}
        </div>
      )}
    </div>
  );
}

function Arrange({
  item,
  lessonId,
  itemIndex,
  onNext,
}: {
  item: LessonItemArrange;
  lessonId: string;
  itemIndex: number;
  onNext: () => Promise<void> | void;
}) {
  const [bank, setBank] = useState<string[]>(() => shuffle(item.words));
  const [built, setBuilt] = useState<string[]>([]);
  const [locked, setLocked] = useState(false);
  const [isCorrect, setIsCorrect] = useState<boolean | null>(null);

  const add = (w: string, i: number) => {
    if (locked) return;
    setBuilt((b) => [...b, w]);
    setBank((b) => b.filter((_, idx) => idx !== i));
  };

  const remove = (i: number) => {
    if (locked) return;
    setBank((b) => [...b, built[i]]);
    setBuilt((b) => b.filter((_, idx) => idx !== i));
  };

  const submit = () => {
    if (locked) return;
    const correct = built.join(" ") === item.words.join(" ");
    setIsCorrect(correct);
    setLocked(true);
    // Per spec: do NOT persist arrange outcomes. We only allow moving on.
  };

  return (
    <div className="space-y-4">
      <div>
        <div className="text-sm text-gray-600 mb-1">Bank (shuffled by UI)</div>
        <div className="flex flex-wrap gap-2">
          {bank.map((w, i) => (
            <button
              key={`${w}-${i}`}
              onClick={() => add(w, i)}
              disabled={locked}
              className="px-3 py-1 rounded-full border text-sm bg-white hover:bg-gray-50"
            >
              {w}
            </button>
          ))}
          {bank.length === 0 && <span className="text-sm text-gray-400">(empty)</span>}
        </div>
      </div>

      <div>
        <div className="text-sm text-gray-600 mb-1">Your sentence</div>
        <div className="flex flex-wrap gap-2 min-h-[2.25rem] items-center">
          {built.map((w, i) => (
            <button
              key={`${w}-built-${i}`}
              onClick={() => remove(i)}
              disabled={locked}
              className="px-3 py-1 rounded-full border text-sm bg-white"
              title="Click to remove back to bank"
            >
              {w}
            </button>
          ))}
          {built.length === 0 && <span className="text-sm text-gray-400">(start picking words)</span>}
        </div>
      </div>

      <div className="flex items-center gap-3">
        {!locked ? (
          <button
            onClick={submit}
            disabled={built.length === 0}
            className="px-4 py-2 rounded-xl bg-gray-900 text-white disabled:opacity-50"
          >
            Submit
          </button>
        ) : (
          <button onClick={onNext} className="px-4 py-2 rounded-xl bg-gray-900 text-white">Next</button>
        )}
        {locked && (
          <span className={classNames(
            "text-sm",
            isCorrect ? "text-green-700" : "text-red-700"
          )}>
            {isCorrect ? "Correct" : "Wrong"}
          </span>
        )}
      </div>
    </div>
  );
}

function SummaryCard({ summary, onRestart }: { summary: Summary; onRestart: () => void }) {
  return (
    <Card title="Lesson Summary">
      <div className="grid gap-2 md:grid-cols-2">
        <div className="space-y-1">
          <div><span className="text-gray-600">Total:</span> {summary.total}</div>
          <div><span className="text-gray-600">Correct:</span> {summary.correct}</div>
          <div><span className="text-gray-600">Wrong:</span> {summary.wrong}</div>
          <div><span className="text-gray-600">Accuracy:</span> {summary.accuracy}%</div>
        </div>
        <div>
          <div className="text-gray-600 mb-1">Scheduled for re‑practice (from fillblanks mistakes):</div>
          {summary.scheduledForRepractice.length ? (
            <div className="flex flex-wrap gap-2">
              {summary.scheduledForRepractice.map((w, i) => (
                <span key={`${w}-${i}`} className="px-2 py-1 rounded-full border text-sm bg-white">{w}</span>
              ))}
            </div>
          ) : (
            <div className="text-sm text-gray-500">None 🎉</div>
          )}
        </div>
      </div>
      <div className="mt-4">
        <button onClick={onRestart} className="px-4 py-2 rounded-xl bg-gray-900 text-white">Start New Lesson</button>
      </div>
    </Card>
  );
}
