import React from "react";
import { Card } from "../Card";

export function StartLessonForm({ userId, setUserId, busy, startLesson, error }: { userId: string; setUserId: (v: string) => void; busy: boolean; startLesson: () => Promise<void>; error: string | null }) {
  return (
    <Card title="🎓 Start a New Lesson">
      <div className="space-y-4">
        <div className="bg-gradient-to-r from-indigo-50 to-purple-50 rounded-lg p-3 border border-indigo-100">
          <p className="text-sm text-gray-700">
            Enter your User ID to begin. Each lesson contains 6 exercises mixing fill-in-the-blank and word arrangement tasks.
          </p>
        </div>
        
        <div className="flex flex-col md:flex-row gap-3 items-start md:items-end">
          <div className="flex-1 w-full">
            <label className="block text-sm font-medium text-gray-700 mb-2">User ID</label>
            <input 
              className="w-full border-2 border-gray-300 rounded-lg px-4 py-3 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all shadow-sm" 
              placeholder="Enter your user ID (e.g., U123)" 
              value={userId} 
              onChange={(e) => setUserId(e.target.value)} 
            />
          </div>
          <button 
            className="px-6 py-3 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium shadow-md hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none" 
            onClick={startLesson} 
            disabled={!userId || busy}
          >
            {busy ? "⏳ Starting..." : "🚀 Start Lesson"}
          </button>
        </div>
        
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-3">
            <p className="text-sm text-red-700 font-medium">{error}</p>
          </div>
        )}
      </div>
    </Card>
  );
}


