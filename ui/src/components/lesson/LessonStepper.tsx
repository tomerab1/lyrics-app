import { Card } from "../Card";
import type { Lesson } from "../../types";
import { FillBlanks } from "./FillBlanks";
import { Arrange } from "./Arrange";

export function LessonStepper({
	lesson,
	index,
	onNext,
}: {
	lesson: Lesson;
	index: number;
	onNext: () => Promise<void> | void;
}) {
	const item = lesson.items[index] as any;
	const progress = ((index + 1) / lesson.items.length) * 100;

	return (
		<Card
			title={
				<div className="space-y-3">
					<div className="text-lg md:text-xl font-semibold text-gray-900">
						Step {index + 1} of {lesson.items.length}
					</div>
					<div className="w-full bg-gray-200 rounded-full h-2.5">
						<div
							className="bg-gradient-to-r from-indigo-600 to-purple-600 h-2.5 rounded-full transition-all duration-300"
							style={{ width: `${progress}%` }}
						/>
					</div>
				</div>
			}
		>
			{item.type === "fillblanks" ? (
				<FillBlanks
					key={`F-${lesson.lessonId}-${index}`}
					item={item}
					lessonId={lesson.lessonId}
					itemIndex={index}
					onNext={onNext}
				/>
			) : (
				<Arrange
					key={`A-${lesson.lessonId}-${index}`}
					item={item}
					onNext={onNext}
				/>
			)}
		</Card>
	);
}
