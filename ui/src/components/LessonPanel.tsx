import { useState } from "react";
import { API_BASE } from "../constants";
import type { Lesson, Summary } from "../types";
import { StartLessonForm } from "./lesson/StartLessonForm";
import { LessonStepper } from "./lesson/LessonStepper";
import { SummaryCard } from "./lesson/SummaryCard";

async function json<T>(res: Response): Promise<T> {
	if (!res.ok) {
		const text = await res.text().catch(() => "");
		throw new Error(text || `HTTP ${res.status}`);
	}
	return res.json();
}

export function LessonPanel() {
	const [userId, setUserId] = useState("");
	const [busy, setBusy] = useState(false);
	const [lesson, setLesson] = useState<Lesson | null>(null);
	const [index, setIndex] = useState(0);
	const [summary, setSummary] = useState<Summary | null>(null);
	const [error, setError] = useState<string | null>(null);

	const startLesson = async () => {
		setBusy(true);
		setError(null);
		setSummary(null);
		try {
			const res = await fetch(`${API_BASE}/lessons`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ userId }),
			});
			const data = await json<{ data: Lesson }>(res);
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
		const res = await fetch(`${API_BASE}/lessons/${lesson.lessonId}/summary`);
		const data = await json<{ data: Summary }>(res);
		setSummary(data.data);
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
			<StartLessonForm
				userId={userId}
				setUserId={setUserId}
				busy={busy}
				startLesson={startLesson}
				error={error}
			/>

			{lesson && !summary && (
				<LessonStepper lesson={lesson} index={index} onNext={onNext} />
			)}
			{lesson && summary && <SummaryCard summary={summary} onRestart={reset} />}
		</div>
	);
}
