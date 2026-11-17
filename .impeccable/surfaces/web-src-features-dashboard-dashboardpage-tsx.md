---
version: 1
slug: "web-src-features-dashboard-dashboardpage-tsx"
primary_target: "web/src/features/dashboard/DashboardPage.tsx"
related_targets: ["web/src/features/home/HomePage.tsx","web/src/features/history/HistoryPage.tsx"]
---

# Surface brief: Dasbor, Beranda, Riwayat (signed-in work screens)

Scope: `web/src/features/dashboard/DashboardPage.tsx` (primary), `web/src/features/home/HomePage.tsx`, `web/src/features/history/HistoryPage.tsx`. Mode: Operate. Established world: The Study Notebook (DESIGN.md), no new identity.

Audience and job: a signed-in SKD candidate deciding what to practise next. Dasbor answers four questions at a glance: which pattern is weakest, what to practise next, is it improving, why mistakes happen. Beranda leads with the next practice (or a running simulation, which outranks it). Riwayat shows simulation results over time above the list. Data: only existing API fields (dashboard sessions/patterns/queue/sections, exam results); derived summaries are allowed, new endpoints are not.
Constraints: voice "kamu"; the dashboard rests on the last 2000 answers and causes are only estimated from answer time, so the note stays; a trend from fewer than 3 points is not drawn (empty state teaches instead); no per-pattern threshold exists, so none is shown; all current fields, links and routes survive; no eyebrows, no motion beyond the bar fill.

## Direction contract

THESIS: The dashboard reads as a grade book: trend first, one highlighted next step, then the ruled pattern list weakest-first. Refuses the category default of a KPI-tile row over equal cards.
OWN-WORLD: Warm paper, flat cards with 1px lines, ink bars on a line-coloured track, mono numerals, one highlighter yellow reserved for the next-practice card and the weakest row. Recognisable with content removed: ruled rows plus one yellow band.
STORY: The user sees whether recent sessions are rising, sees which single pattern to practise and why, and starts it in one click; the causes sit under the weak rows for the reader who wants them.
FIRST VIEWPORT: Desktop 1280: H1 "Titik lemah" and the data note; left a card with 10 session columns (newest ink, older sand, value above, date below); right the yellow "Latih berikutnya" card (pattern, reason, ink button); below, the "Per pola soal" ruled list, weakest row highlighted with its causes. Mobile 390: the same order in one column, columns narrower, causes wrap. Beranda: H1, then the yellow next-practice card (or running-simulation card), Drill and Simulasi choices, last result. Riwayat: H1, score columns for finished simulations, then the existing list.
FORM: Buku nilai, dealt option of my 7 grounded structures (the card at page position 1), seed:e9078ecf08eda865.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md drift report, and no shipping raster without provenance (none planned).
