# Lyrics Practice Mini App


## Overview
Two tabs:

1) **Admin**
- Create users (required `name`).
- Create songs (`title`, `artist`, `lyrics`).
- Show short lists of users and songs (minimal fields).

2) **UI**
- User enters `userId` and starts a lesson.
- Server picks a song for the user.
- Each lesson has **6 assignments** in a **stepper** flow.
- Assignment types:
  - **Fill Blanks**: one word in a line is replaced by `___`. User chooses from **4 options**.
  - **Arrange Words**: the same line’s words are provided **in correct order**; the UI shuffles and compares.
- Only **Fill Blanks** mistakes are saved and used to prioritize future lessons.
- End of lesson: show stats (correct, wrong, accuracy, and words to re‑practice).

Keep styling minimal. No animations. No retries.

---

## Endpoints

### Admin
- `POST /users` `{ name }` → `{ id }`
- `GET /users` → `[ { id, name } ]`
- `POST /songs` `{ title, artist, lyrics }` → `{ id, lineCount }`
- `GET /songs` → `[ { id, title } ]`

### Lessons
- `POST /lessons` `{ userId }` → lesson of 6 items  
  **Response**
  ```json
  {
    "lessonId": "L1",
    "items": [
      {
        "type": "fillblanks",
        "lineIndex": 3,
        "renderedLine": "When you ___ your best but you don't succeed",
        "words": ["try", "fail", "run", "fall"]  // exactly 4 options: 1 correct + 3 distractors
      },
      {
        "type": "arrange",
        "lineIndex": 4,
        "words": ["and", "I", "will", "try", "to", "fix", "you"]  // CORRECT ORDER; UI shuffles and compares
      }
    ]
  }
  ```
  Notes:
  - For **fillblanks**, `words` is a 4‑item options bank (1 correct + 3 distractors).
  - For **arrange**, `words` is the **correct order**; the **UI** shuffles and validates exact match.

- `POST /answers` `{ lessonId, itemIndex, type, correct, userInput? }` → `{ ok: true }`
  - **Persist** answers only when `type === "fillblanks"`.
  - **Ignore persistence** for `arrange`.
  - Duplicate submissions for the same item should return **409**.

- `GET /lessons/:lessonId/summary` →
  ```json
  {
    "total": 6,
    "correct": 4,
    "wrong": 2,
    "accuracy": 66.7,
    "scheduledForRepractice": ["try", "succeed"]
  }
  ```
`scheduledForRepractice` is derived only from fillblanks mistakes.

---

## Lesson Generation (Server)
- Choose a song for the user (prefer ones with open **fillblanks** mistakes; otherwise any).
- Build **6** items from lines in that song.
- Include both assignment types.
- Fix the 6 items at creation time (no mid‑lesson regeneration).

---

## UI Requirements
- **Admin tab**
  - Create User form.
  - Create Song form.
  - Small lists: users (id, name) and songs (id, title).

- **UI tab**
  - Input for `userId`; Start Lesson button.
  - Stepper: “Step X of 6”.
  - **Fill Blanks**: show `renderedLine` and 4 buttons from `words`.
  - **Arrange Words**: show shuffled `words` as a bank, user builds the sentence once; compare to **original order**.
  - After step 6: show summary and “Start New Lesson” button.

---

## Persistence
- Persist: users, songs, lessons, and **fillblanks** answers/history.
- Do **not** persist arrange outcomes.

---

## Database Options
- Recommended: **MongoDB**.
- Lightweight: **local JSON file** store.
- Other valid choices: **SQLite**, **Postgres**, etc.

---

## Monorepo Structure

```
api/        # Go backend (Chi, MongoDB)
ui/         # React + Vite frontend
README.md   # This file
```

## Getting Started

### Prerequisites
- Go 1.22+
- Node.js 18+
- Docker (for MongoDB via compose)

### Backend (API)

1) Start MongoDB:
```bash
cd api
docker compose up -d
```

2) Create `api/.env`:
```env
SERVER_ADDR=:5555
MONGO_ADDR=mongodb://localhost:27017
```

3) Run API:
```bash
cd api
go run ./cmd
```

The API will be available at `http://localhost:5555/api`.

### Frontend (UI)

```bash
cd ui
npm install
npm run dev
```

Open the printed localhost URL (typically `http://localhost:5173`).

## API Details

Base URL: `http://localhost:5555/api`

- POST `/users` body `{ name }` → `{ data: { id } }`
- GET `/users` → `{ data: [ { id, name } ] }`
- POST `/songs` body `{ title, artist, lyrics }` → `{ data: { id, lineCount } }`
- GET `/songs` → `{ data: [ { id, title } ] }`
- POST `/lessons` body `{ userId }` → `{ data: { lessonId, items } }`
- POST `/answers` body `{ lessonId, itemIndex, type, userInput }` → `{ data: { ok, correct } }`
  - Only persisted for `type === "fillblanks"`.
  - Duplicate answer per item returns 409.
- GET `/lessons/{lessonId}/summary` → bare JSON `{ total, correct, wrong, accuracy, scheduledForRepractice }`

Notes:
- For fillblanks, the server sends 4 options (1 correct + 3 distractors). Correctness is validated server‑side on submission.
- For arrange, the server sends the correct order; UI shuffles and validates locally; results are not persisted.

## Implementation Notes

- Backend uses Chi router and MongoDB official driver.
- Lessons are stored with fixed items at creation. Fillblanks items embed the hidden word for server validation; UI never receives it.
- Answers are stored once per item; duplicates are rejected.
- Summary aggregates total items, correctness across fillblanks, and lists mistaken words for re‑practice scheduling.

## Frontend Notes

- Components are split into smaller pieces under `ui/src/components` and app constants in `ui/src/constants.ts`.
- The UI is intentionally minimal: no animations and single attempt per item.
