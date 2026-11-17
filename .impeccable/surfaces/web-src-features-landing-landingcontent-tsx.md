---
version: 1
slug: "web-src-features-landing-landingcontent-tsx"
primary_target: "web/src/features/landing/LandingContent.tsx"
related_targets: []
---

# Surface brief: landing ("/" for guests)

Scope: `web/src/features/landing/LandingContent.tsx`, prerendered to static HTML at build. Mode: Persuade. Established world: The Study Notebook (DESIGN.md), no new identity.

Audience and job: CPNS/ASN applicants deciding whether to start practising. Action: "Mulai berlatih" to /login. Proof on hand: one sample question with steps and error type (own data, labelled contoh), a verbatim TWK source quote with its human review status, the exam spec (110 soal, 100 menit, ambang 65/80/166). Constraints: voice "kamu"; no testimonials, user counts or pass guarantees; disclaimer without naming official bodies in landing markup (landing.test.tsx); no per-pattern threshold is shown because the app has none.

Unresolved: font loading is a separate step; footer in App.tsx repeats the disclaimer.

## Direction contract

THESIS: The first screen is the diagnosis itself, a worked result sheet, not a pitch about diagnosis. Refuses the category default of headline + button + row of equal feature cards.
OWN-WORLD: Warm paper ground, ink text, one highlighter yellow reserved for the weakest row and the primary action. Flat card sheet with 1px lines, mono numerals, bars as ink on a line-coloured track, status as icon + text. Recognisable with content removed: a ruled-by-rows result sheet with one yellow band.
STORY: Visitor understands that each mistake gets a named type and a pattern score; believes it because the key provenance and exam spec are shown; acts by starting practice.
FIRST VIEWPORT: Desktop 1280: left column H1 (52px, display, no eyebrow above it) with highlighted "jenis soal", lead, primary CTA; right column (about 520px) the sheet "Titik lemahmu" with a Contoh chip, three pattern rows (bar + mono percent, weakest row highlighted, error-type line under it), and "Latih pola ini" at its foot. Mobile 390: H1, lead, CTA, then the sheet at full width. Signature motion: bars fill once on load, exponential ease-out, staggered, disabled under reduced motion.
FORM: Peta titik lemah, position 7 of my 7 grounded structures (dealt lead), seed:af13b5045b27776c.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
