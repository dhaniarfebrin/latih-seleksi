---
version: 1
slug: "web-src-features-exam-components-running-tsx"
primary_target: "web/src/features/exam/components/Running.tsx"
related_targets: ["web/src/features/drill/components/Question.tsx","web/src/features/exam/components/Result.tsx"]
---

# Surface brief: Drill dan Ujian (answering screens)

Scope: `web/src/features/drill/components/{DrillStart,Question,Summary}.tsx` and `DrillPage.tsx` (primary), `web/src/features/exam/StartPage.tsx`, `components/Running.tsx`, `components/Result.tsx`. Mode: Operate. Established world: The Study Notebook (DESIGN.md), no new identity.

Audience and job: a signed-in SKD candidate answering questions under time pressure. They need the time left, where they are, and what is still open, without leaving the question. Drill adds per-question time against an optional target and an explanation after each answer. Data: existing API fields only; no new endpoints.
Constraints: voice "kamu"; the server owns the exam clock (the timer stays `remainingSeconds` with the server offset, `role="timer"`); every current control, answer flow, flag, confirm-finish step, report form and route survives; weighted TKP is scored by weight, never shown as right or wrong; no eyebrows, no colored side borders, no motion beyond state changes; the floating mobile nav stays hidden during an exam.

## Direction contract

THESIS: Time and position are one band above the question: big time left, the subtests as tabs, a ruler of question numbers, and the question alone below it. Refuses the category default of a side panel of number buttons plus a separate timer widget.
OWN-WORLD: Warm paper, a flat card for the question, ink-filled ruler cells for answered, dashed for doubtful, one highlighter-yellow cell for the current question. Recognisable with content removed: a ruler of numbered cells under a mono clock.
STORY: The user glances up to read the clock and the ruler, answers, and moves on; the ruler shows what is still empty or doubtful, and the same band drives Drill with 8 cells and per-question time.
FIRST VIEWPORT: Desktop 1280 Ujian: band (clock, subtest tabs with answered counts, Selesai, progress track, ruler of the active subtest), question card about 720px centered with stem, five options, Sebelumnya / Tandai ragu / Berikutnya. Mobile 390: clock and Selesai on one line, subtest tabs scroll sideways, ruler folded into a "Nomor soal (x/y terjawab)" disclosure, question card full width. Drill: the same band with an 8-cell read-only ruler and optional target, explanation panel under the card after answering. Hasil: score, per-subtest ruled rows with bars, then the review.
FORM: Pita waktu, assigned roll of my 7 grounded structures (index 7 leads; dealt 7, 1, 6), seed reference c6ea606730331a781dd4e51d488c32de.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md drift report, and no shipping raster without provenance (none planned).
