import React, { useState } from "react";
import type { LessonItemArrange } from "../../types";

function classNames(...xs: Array<string | false | null | undefined>) {
  return xs.filter(Boolean).join(" ");
}

export function Arrange({ item, onNext }: { item: LessonItemArrange; onNext: () => Promise<void> | void }) {
  const shuffle = <T,>(arr: T[]): T[] => {
    const a = [...arr];
    for (let i = a.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [a[i], a[j]] = [a[j], a[i]];
    }
    return a;
  };

  const [bank, setBank] = useState<string[]>(() => shuffle(item.words));
  const [built, setBuilt] = useState<string[]>([]);
  const [locked, setLocked] = useState(false);
  const [isCorrect, setIsCorrect] = useState<boolean | null>(null);

  const add = (w: string, i: number) => { if (locked) return; setBuilt((b) => [...b, w]); setBank((b) => b.filter((_, idx) => idx !== i)); };
  const remove = (i: number) => { if (locked) return; setBank((b) => [...b, built[i]]); setBuilt((b) => b.filter((_, idx) => idx !== i)); };
  const submit = () => { if (locked) return; const correct = built.join(" ") === item.words.join(" "); setIsCorrect(correct); setLocked(true); };

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-blue-50 to-cyan-50 rounded-xl p-4 border border-blue-100">
        <div className="text-sm font-medium text-blue-700 mb-2">📝 Arrange the words to form the sentence:</div>
      </div>
      
      <div>
        <div className="text-sm font-medium text-gray-700 mb-3 flex items-center gap-2">
          <span>🔤 Available Words</span>
          {bank.length === 0 && <span className="text-xs text-green-600 font-normal">(all words used)</span>}
        </div>
        <div className="flex flex-wrap gap-2 p-4 bg-gray-50 rounded-xl border-2 border-dashed border-gray-200 min-h-[60px]">
          {bank.map((w, i) => (
            <button 
              key={`${w}-${i}`} 
              onClick={() => add(w, i)} 
              disabled={locked} 
              className="px-4 py-2 rounded-lg border-2 border-blue-300 bg-white text-sm font-medium text-blue-700 hover:bg-blue-50 hover:border-blue-400 hover:shadow-md transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {w}
            </button>
          ))}
        </div>
      </div>
      
      <div>
        <div className="text-sm font-medium text-gray-700 mb-3 flex items-center gap-2">
          <span>✨ Your Sentence</span>
          {built.length === 0 && <span className="text-xs text-gray-500 font-normal">(drag words here or click to build)</span>}
        </div>
        <div className="flex flex-wrap gap-2 p-4 bg-indigo-50 rounded-xl border-2 border-indigo-200 min-h-[60px] items-center">
          {built.map((w, i) => (
            <button 
              key={`${w}-built-${i}`} 
              onClick={() => remove(i)} 
              disabled={locked} 
              className="px-4 py-2 rounded-lg border-2 border-indigo-400 bg-white text-sm font-medium text-indigo-700 hover:bg-indigo-100 hover:border-indigo-500 hover:shadow-md transition-all disabled:cursor-not-allowed"
              title="Click to remove"
            >
              {w}
            </button>
          ))}
        </div>
      </div>
      
      <div className="flex items-center justify-between p-4 bg-gray-50 rounded-xl border border-gray-200">
        <div className="flex items-center gap-3">
          {locked && (
            <span className={classNames(
              "px-4 py-2 rounded-lg font-semibold text-sm shadow-sm",
              isCorrect ? "bg-green-100 text-green-700 border border-green-300" : "bg-red-100 text-red-700 border border-red-300"
            )}>
              {isCorrect ? "✓ Correct!" : "✗ Wrong"}
            </span>
          )}
        </div>
        {!locked ? (
          <button 
            onClick={submit} 
            disabled={built.length === 0} 
            className="px-6 py-2 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium shadow-md hover:shadow-lg transition-all transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
          >
            Submit Answer
          </button>
        ) : (
          <button 
            onClick={onNext} 
            className="px-6 py-2 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium shadow-md hover:shadow-lg transition-all transform hover:scale-105"
          >
            Next →
          </button>
        )}
      </div>
    </div>
  );
}


