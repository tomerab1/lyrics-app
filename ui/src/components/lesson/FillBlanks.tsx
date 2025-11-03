import React, { useState } from "react";
import { API_BASE } from "../../constants";
import type { LessonItemFill } from "../../types";

function classNames(...xs: Array<string | false | null | undefined>) {
	return xs.filter(Boolean).join(" ");
}

export function FillBlanks({
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

	const clickOption = async (w: string) => {
		if (submitting || chosen) return;
		setChosen(w);
		setSubmitting(true);
		setError(null);
		try {
			const correctLocal = item.correct_word
				? w.toLowerCase() === item.correct_word.toLowerCase()
				: false;
			setResult(correctLocal ? "correct" : "wrong");
			const res = await fetch(`${API_BASE}/answers`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					lessonId,
					itemIndex,
					type: "fillblanks",
					userInput: w,
				}),
			});
			if (res.status === 409) {
				setError("Already answered. Moving on...");
				setTimeout(() => onNext(), 600);
				return;
			}
			const payload = await res.json().catch(() => ({ ok: true } as any));
			if (typeof payload.correct === "boolean") {
				setResult(payload.correct ? "correct" : "wrong");
			}
		} catch (err: any) {
			setError(err.message);
		} finally {
			setSubmitting(false);
		}
	};

	return (
		<div className="space-y-6">
			<div className="bg-linear-to-r from-indigo-50 to-purple-50 rounded-xl p-4 border border-indigo-100">
				<p className="text-lg md:text-xl text-gray-900 font-medium leading-relaxed">
					{item.renderedLine.split("___").map((part, i, arr) => (
						<React.Fragment key={i}>
							{part}
							{i < arr.length - 1 && (
								<span
									className={classNames(
										"inline-flex mx-2 items-center justify-center rounded-md font-semibold min-w-[60px] transition-all",
										"h-8",
										!chosen &&
											"bg-white border-2 border-dashed border-indigo-300 text-indigo-600",
										chosen &&
											result === "correct" &&
											"bg-green-100 border-2 border-green-400 text-green-700",
										chosen &&
											result === "wrong" &&
											"bg-red-100 border-2 border-red-400 text-red-700"
									)}
								>
									{chosen ? chosen : "?"}
								</span>
							)}
						</React.Fragment>
					))}
				</p>
			</div>

			<div>
				<div className="text-sm font-medium text-gray-700 mb-3">
					Choose the correct word:
				</div>
				<div className="grid grid-cols-2 md:grid-cols-4 gap-3">
					{item.words.map((w) => (
						<button
							key={w}
							disabled={!!chosen}
							onClick={() => clickOption(w)}
							className={classNames(
								"px-4 py-3 rounded-xl border-2 font-medium text-sm transition-all shadow-sm",
								"disabled:cursor-not-allowed disabled:opacity-50",
								!chosen &&
									"bg-white border-gray-200 hover:border-indigo-400 hover:bg-indigo-50 hover:shadow-md hover:scale-105",
								chosen === w &&
									result === "correct" &&
									"border-green-500 bg-green-50 text-green-700 shadow-md scale-105",
								chosen === w &&
									result === "wrong" &&
									"border-red-500 bg-red-50 text-red-700 shadow-md scale-105",
								chosen && chosen !== w && "opacity-50"
							)}
						>
							{w}
						</button>
					))}
				</div>
			</div>

			{error && (
				<div className="bg-red-50 border border-red-200 rounded-lg p-3">
					<p className="text-sm text-red-700">{error}</p>
				</div>
			)}

			{chosen && (
				<div className="flex items-center justify-between p-4 bg-gray-50 rounded-xl border border-gray-200">
					<div className="flex items-center gap-3">
						{result && (
							<span
								className={classNames(
									"px-4 py-2 rounded-lg font-semibold text-sm shadow-sm",
									result === "correct"
										? "bg-green-100 text-green-700 border border-green-300"
										: "bg-red-100 text-red-700 border border-red-300"
								)}
							>
								{result === "correct" ? "✓ Correct!" : "✗ Wrong"}
							</span>
						)}
					</div>
					<button
						onClick={onNext}
						className="px-6 py-2 rounded-xl bg-linear-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium shadow-md hover:shadow-lg transition-all transform hover:scale-105"
					>
						Next →
					</button>
				</div>
			)}
		</div>
	);
}
