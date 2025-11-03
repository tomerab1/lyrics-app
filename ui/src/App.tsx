import React, { useState } from "react";
import { AdminPanel } from "./components/AdminPanel";
import { LessonPanel } from "./components/LessonPanel";

function classNames(...xs: Array<string | false | null | undefined>) {
  return xs.filter(Boolean).join(" ");
}

export default function App() {
  const [activeTab, setActiveTab] = useState<"admin" | "ui">("admin");
  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-pink-50 text-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-5xl">
        <header className="mb-6 text-center">
          <h1 className="text-2xl md:text-3xl font-semibold text-gray-900">Lyrics Practice Mini App</h1>
        </header>
        <nav className="mb-6 flex justify-center">
          <div className="inline-flex rounded-full overflow-hidden bg-white/70 backdrop-blur ring-1 ring-black/5 shadow-sm">
            <button className={classNames("px-4 py-2 text-sm md:text-base transition", activeTab === "admin" ? "bg-white text-indigo-700 shadow" : "bg-transparent text-gray-600 hover:bg-white/60")} onClick={() => setActiveTab("admin")}>
              Admin
            </button>
            <button className={classNames("px-4 py-2 text-sm md:text-base transition", activeTab === "ui" ? "bg-white text-indigo-700 shadow" : "bg-transparent text-gray-600 hover:bg-white/60")} onClick={() => setActiveTab("ui")}>
              UI
            </button>
          </div>
        </nav>
        {activeTab === "admin" ? <AdminPanel /> : <LessonPanel />}
      </div>
    </div>
  );
}
