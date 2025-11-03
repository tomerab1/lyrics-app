import React from "react";
import { Card } from "../Card";
import type { Summary } from "../../types";

export function SummaryCard({ summary, onRestart }: { summary: Summary; onRestart: () => void }) {
  const scheduled = summary.scheduledForRepractice || [];
  const accuracyColor = summary.accuracy >= 80 ? "text-green-600" : summary.accuracy >= 60 ? "text-yellow-600" : "text-red-600";
  
  return (
    <Card title="🎯 Lesson Complete!">
      <div className="grid gap-6 md:grid-cols-2">
        {/* Performance Metrics */}
        <div className="space-y-4">
          <h3 className="text-md font-semibold text-gray-800 border-b pb-2">Performance</h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-blue-50 rounded-lg p-3 border border-blue-100">
              <div className="text-xs text-blue-600 font-medium mb-1">Total Items</div>
              <div className="text-2xl font-bold text-blue-700">{summary.total}</div>
            </div>
            <div className="bg-green-50 rounded-lg p-3 border border-green-100">
              <div className="text-xs text-green-600 font-medium mb-1">Correct</div>
              <div className="text-2xl font-bold text-green-700">{summary.correct}</div>
            </div>
            <div className="bg-red-50 rounded-lg p-3 border border-red-100">
              <div className="text-xs text-red-600 font-medium mb-1">Wrong</div>
              <div className="text-2xl font-bold text-red-700">{summary.wrong}</div>
            </div>
            <div className={`bg-gradient-to-br ${summary.accuracy >= 80 ? 'from-green-50 to-emerald-50 border-green-100' : summary.accuracy >= 60 ? 'from-yellow-50 to-amber-50 border-yellow-100' : 'from-red-50 to-rose-50 border-red-100'} rounded-lg p-3 border`}>
              <div className="text-xs font-medium mb-1" style={{color: summary.accuracy >= 80 ? '#16a34a' : summary.accuracy >= 60 ? '#d97706' : '#dc2626'}}>Accuracy</div>
              <div className={`text-2xl font-bold ${accuracyColor}`}>{summary.accuracy.toFixed(1)}%</div>
            </div>
          </div>
        </div>

        {/* Words to Practice */}
        <div className="space-y-4">
          <h3 className="text-md font-semibold text-gray-800 border-b pb-2">Words to Practice</h3>
          {scheduled.length > 0 ? (
            <div className="space-y-3">
              <p className="text-sm text-gray-600">
                You missed these words in the fill-in-the-blank exercises. These will be prioritized in future lessons to help you improve.
              </p>
              <div className="flex flex-wrap gap-2 p-3 bg-amber-50 rounded-lg border border-amber-200">
                {scheduled.map((w, i) => (
                  <span 
                    key={`${w}-${i}`} 
                    className="px-3 py-1.5 rounded-full text-sm font-medium bg-white border-2 border-amber-300 text-amber-800 shadow-sm hover:bg-amber-100 transition-colors"
                  >
                    {w}
                  </span>
                ))}
              </div>
            </div>
          ) : (
            <div className="p-6 bg-gradient-to-br from-green-50 to-emerald-50 rounded-lg border border-green-200 text-center">
              <div className="text-3xl mb-2">🎉</div>
              <div className="text-sm font-medium text-green-800 mb-1">Perfect Score!</div>
              <div className="text-xs text-green-600">You got all fill-in-the-blank questions correct. Great job!</div>
            </div>
          )}
        </div>
      </div>
      
      <div className="mt-6 pt-4 border-t border-gray-200 flex justify-center">
        <button 
          onClick={onRestart} 
          className="px-6 py-3 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium shadow-lg hover:shadow-xl transition-all transform hover:scale-105"
        >
          Start New Lesson
        </button>
      </div>
    </Card>
  );
}


