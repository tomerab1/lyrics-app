export type User = { id: string; name: string };
export type Song = { id: string; title: string };

export type LessonItemFill = {
  type: "fillblanks";
  lineIndex: number;
  renderedLine: string;
  words: string[]; // exactly 4 options
  correct_word?: string; // provided by server, use for local validation
};

export type LessonItemArrange = {
  type: "arrange";
  lineIndex: number;
  words: string[]; // correct order; UI shuffles and compares
};

export type LessonItem = LessonItemFill | LessonItemArrange;

export type Lesson = {
  lessonId: string;
  items: LessonItem[];
};

export type Summary = {
  total: number;
  correct: number;
  wrong: number;
  accuracy: number; // percentage
  scheduledForRepractice: string[];
};


