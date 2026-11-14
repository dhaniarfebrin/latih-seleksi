# AI and adaptive learning for a CPNS SKD / BUMN test-prep app (as of 2026-11-14)

Access date for all URLs: 2026-11-14. Publication dates are given where the page showed one. Many findings below come from search-result summaries plus a few full fetches; those are marked "(search summary)" and should be checked against the primary page before being quoted. Items without a source are listed under Gaps, not as facts.

## 1. Automatic item generation (AIG): templates vs LLM, error rates, validation

### Takeaway
Template/parametric AIG is deterministic and cheap per item but needs heavy up-front expert modelling, so it fits only narrow, repetitive skills (numeric series, ratios, percentages). LLM generation covers verbal and factual items but published checks show a meaningful share of flawed items (roughly 5% to 36% depending on study and rubric), so a verification layer is mandatory.

### Cited Findings
- Template-based AIG (Gierl): build a cognitive model, then an item model ("mould"), then software generates items; one item model can yield hundreds or thousands of items — [NCME Module 34, Gierl & Lai](https://ncme.org/wp-content/uploads/2025/10/Module-34-Automated-Item-Generation-Gierl-Lai.pdf) (search summary)
- Main drawback of template AIG is heavy up-front subject-matter-expert effort; cost-effectiveness depends on most items belonging to the same content area with a limited number of skills modelled by one cognitive model (Kosh et al. 2019, as cited; original not read) — [NCME Module 34 and related AIG review](https://ncme.org/wp-content/uploads/2025/10/Module-34-Automated-Item-Generation-Gierl-Lai.pdf) (search summary)
- Example scale in one Gierl-coauthored study: 35 cognitive models, 30 items generated per model — [UAlberta ERA record](https://era.library.ualberta.ca/items/19affa9f-c4b5-4278-9110-82904f5a0bd3/download/1c6157b6-05a1-47ee-bf5c-846cd2c3352b) (search summary)
- LLM-generated MCQs for cancer-trial education (GPT-4, 1,335 items): 78 items (5%) were flagged by difficulty/agreement thresholds for manual error analysis; among flagged items, 10 had more than one correct answer and 7 had an incorrect keyed answer — [arXiv 2412.01955](https://arxiv.org/pdf/2412.01955) (search summary quoting the paper). Note: this is a flagged subset, not a full audit of all 1,335.
- NBME-affiliated study (ACL, AIME-Con, Oct 2025): abstract says experts judged AI-generated clinical MCQ drafts as helpful/non-helpful starting points and classifiers (feature-based, fine-tuned transformers, GPT-4 few-shot) were tested; the abstract has no numbers — [ACL Anthology 2025.aimecon-main.19](https://aclanthology.org/2025.aimecon-main.19/). A figure of "64% of 512 items unanimously deemed correct, 36% flagged" appeared only in a search-result summary of that paper; I could not read the PDF text, so treat as **unverified**.
- Surface quality is no longer the issue; implausible distractors and vague wording remain common failure modes of LLM MCQs — [arXiv 2405.20529](https://arxiv.org/pdf/2405.20529) (search summary)
- Distractor research 2025: D-GEN (open-source 8B/70B distractor generator, ACL Findings 2025), student-choice-prediction distractors (ACL 2025), PlausibleQA (10,000 questions, 100,000 candidate answers, SIGIR 2025); a Japanese nursing-exam study (CSEDU 2025) had 15 experts write distractors with LLM candidates offered to 8 of them and says real-exam validation is still needed — [D-GEN arXiv 2504.13439](https://arxiv.org/html/2504.13439v2), [Lee et al. arXiv 2501.13125](https://arxiv.org/pdf/2501.13125), [SciTePress 2025](https://scitepress.org/PublishedPapers/2025/134603) (search summary)
- Self-consistency (majority vote over multiple sampled solutions) works for items with directly comparable answers (a number, an option letter), but agreement is only a weak/conditional correctness signal and fails when the model shares the same misconception every run — [Ding 2026 audit, alphaXiv 2607.08065](https://www.alphaxiv.org/abs/2607.08065) (search summary)
- LLMs can change their answer to the same MCQ when option order is shuffled, so a solver check should be repeated with shuffled options — [LREC-COLING 2024, 2024.lrec-main.1229](https://aclanthology.org/2024.lrec-main.1229.pdf) (search summary)
- Duolingo English Test production pattern: generative AI creates items, humans review for bias/fairness, items are calibrated with IRT and item embeddings — [NeurIPS 2024 talk](https://neurips.cc/virtual/2024/107597); a Frontiers in AI pilot used 713 listening tasks with hundreds of responses each, supporting human-in-the-loop generation — [Frontiers 2024](https://www.frontiersin.org/journals/artificial-intelligence/articles/10.3389/frai.2024.1474019/epub) (search summary)

### Inferences
- For a solo developer, the split that fits the evidence: parametric code generators for numeric TIU (answer computed by code, so the key is correct by construction), LLM + source grounding + automated checks for TWK and verbal TIU, and a human skim of everything (or at least a sampled/flagged subset) before publication. This is my synthesis, not a published recommendation.
- Because no study tested "solver agreement + shuffled options" on exam items directly, treat it as a filter that removes ambiguous/mis-keyed items, not as proof of correctness.
- A planned human-review rate of at least a sample plus all items where solvers disagree is consistent with the 5%-36% flaw rates above; the exact rate for Indonesian CPNS-style items is unknown.

### Gaps
- No published error rates found for LLM-generated items in Indonesian or for civil-service/aptitude exams specifically.
- No source found on LLM generation of TKP (situational judgment) items or their scoring keys (TKP uses graded 1-5 option scores); no source found on LLM generation of figural reasoning items (needs image generation or programmatic SVG); both are open.
- The 64%/36% NBME figure above is unverified.

## 2. LLM math/reasoning errors in explanations and mitigations

### Takeaway
LLMs make arithmetic and logic slips even when the decomposition is right; offloading computation to code (PAL) and a calculator tool are the documented mitigations, and Khan Academy uses a calculator plus retrieval of human-written solutions in Khanmigo. Compute and fix the verified answer first, then ask the LLM only to explain it.

### Cited Findings
- PAL (program-aided LMs): the LLM writes code and a Python interpreter does the arithmetic; with Codex it beat PaLM-540B chain-of-thought on GSM8K by about 15 absolute points (72.0% vs 65.6% in a secondary summary); motivation is that LLMs make arithmetic/logic mistakes in the solution step even when decomposition is correct — [arXiv 2211.10435](https://arxiv.org/abs/2211.10435), [ICML 2023 proceedings](https://proceedings.mlr.press/v202/gao23f.html) (search summary; the 72.0/65.6 numbers come from a secondary source). These are 2022-era models; current models are stronger but the mechanism still holds.
- Khan Academy (blog first published 2025-02-26, updated 2025-03-12): built a calculator for Khanmigo instead of relying on predictive arithmetic; moved model from GPT-4 Turbo to GPT-4o; made the model consult human-written exercises, steps, hints and solutions first (manual review found better accuracy); generated text descriptions of graphics; added a tutoring-quality dashboard. The post gives no error rates — [Khan Academy blog](https://blog.khanacademy.org/khanmigo-math-computation-and-tutoring-updates/)
- Khan Academy's released dataset found AI models often mis-grade student work (telling correct students they are wrong or vice versa), partly from calculation errors — [Khan Academy blog](https://blog.khanacademy.org/?p=16736) (search summary)
- A 2025 arXiv study classifies step-level computation/reasoning errors of GPT-4o, o1, DeepSeek-V3, DeepSeek-R1 on arithmetic, algebra, number theory — [arXiv 2508.09932](https://arxiv.org/pdf/2508.09932) (search summary; error rates not retrieved)
- A 2024 PLOS ONE study reportedly found ChatGPT-generated tutoring help still had errors about 13% of the time in statistics — reported only in a search-result summary; the paper was not opened and no URL was returned, so treat as **unverified**
- A 2026 arXiv paper reports that AI tutors have "correct answer trap" blind spots when analysing student reasoning — [arXiv 2605.23925](https://arxiv.org/pdf/2605.23925) (search summary)

### Inferences
- Pipeline for numeric/logic items: the generator emits structured data (problem, answer computed by code); the LLM only writes the explanation from the given worked solution; a second pass checks that the explanation's final number matches the code's answer. This is an inference from PAL + Khan's approach.
- For deterministic skills (series, ratios, speed/distance, percentages), do not use the LLM as solver at all.

### Gaps
- No current (2026-model) measured error rate for LLM explanations of aptitude-style numeric items was found.
- No source found comparing "code-verified-answer then explain" against plain LLM explanations on exam-prep items.

## 3. Factual question generation grounded in a source corpus (RAG)

### Takeaway
Retrieval grounding reduces hallucination in generated exam material and can lift cheaper models toward parity, but it does not remove reasoning errors or unfaithful generation, and faithfulness and correctness must be checked separately. Keep each TWK item tied to an explicit quoted passage (pasal/ayat) so it can be audited.

### Cited Findings
- A 2025 Cambridge study had GPT-3.5 and GPT-4 Turbo create and answer undergraduate exams from one textbook; adding RAG cut hallucinations and raised accuracy more than a more advanced base model alone — [Cambridge Core](https://www.cambridge.org/core/product/498B8CB29AA37695AC5EACBE1B2089B4) (search summary)
- RAG addresses hallucination from missing knowledge but not from reasoning errors or drift from retrieved content; a preprint finds high-quality retrieval often co-occurs with unfaithful generation — [Barking Iguana guide (lower authority)](https://barkingiguana.com/writing/pop-quiz-faithfulness-vs-correctness/), [arXiv 2505.21072](https://arxiv.org/pdf/2505.21072) (search summary)
- An answer can be faithful to the wrong passage; score faithfulness and correctness separately — [Barking Iguana](https://barkingiguana.com/writing/pop-quiz-faithfulness-vs-correctness/) (blog, lower authority)
- A proposed RAG-MCQ pipeline has reviewers approve questions before they enter the bank and requires each question to be answerable from explicit evidence (slides, lower authority) — [ICETAI slides](https://icetai.neu.edu.vn/slides/rag-mcq-generation.pdf)
- Base LLM knowledge of Indonesian content is weak: on IndoMMLU (about 14,900 questions, 63 tasks, EMNLP 2023) GPT-3.5 scored about 53.2% zero-shot and 39.3% on local languages/cultures; these are 2023 models, so current numbers will be higher but were not retrieved — [arXiv 2310.04928](https://arxiv.org/pdf/2310.04928)

### Inferences
- Do not let the model "remember" UUD 1945 or Pancasila; retrieve the exact pasal text from an authoritative copy (not model memory), require the item to cite the pasal, and verify automatically that the keyed answer is a substring/paraphrase supported by the cited passage. The corpus must also be kept current for amended articles; the item bank should store the source reference for later audit.
- TWK also covers history, national integration, and current-affairs style facts; those depend on corpus coverage, which I have not checked.

### Gaps
- No Indonesian-language RAG question-generation evaluation found; no current IndoMMLU-style scores for 2026 models retrieved.
- Official CPNS TWK syllabus/corpus sources are outside this note's scope (other researchers).

## 4. Adaptive practice: knowledge tracing, IRT/CAT, spaced repetition, mastery learning

### Takeaway
With small data, simple models win: logistic-regression-style or extended BKT beat deep knowledge tracing until datasets are large; IRT needs hundreds of responses per calibration set and cannot be properly calibrated for a new app, so use LLM/feature-based difficulty priors and simple Elo-like updating first. FSRS with default parameters is a drop-in, evidence-backed scheduler for the fact-memory part (TWK), but it was validated on Anki flashcard data, not MCQ exam practice.

### Cited Findings
- Gervet et al. (JEDM 2020, nine datasets): logistic regression with good features leads on moderate-sized datasets, DKT leads on large ones, BKT/Markov methods lag — [JEDM](https://jedm.educationaldatamining.org/index.php/JEDM/article/download/451/123) (search summary)
- Khajah et al. 2016: BKT with published extensions performs indistinguishably from DKT — [How deep is knowledge tracing?](https://cs.paperswithcode.com/paper/how-deep-is-knowledge-tracing) (search summary)
- Shih et al. 2023: a simpler decoder with fewer parameters can predict student performance better than DKT's on several small and large datasets — [arXiv 2309.12334](https://arxiv.org/pdf/2309.12334v1) (search summary)
- IRT sample-size rules of thumb: several hundred examinees for 1PL/2PL, several thousand for 3PL (Drasgow 1989; Thissen & Wainer 1982, as cited); 500-1,000 for 3PL under ML estimation; one simulation found 500 enough for a 27-item test and 1,000 for a 70-item test; another cited study: 1PL acceptable at N=150, 2PL needs 250-750 depending on test length — [PMC6023095](https://pmc.ncbi.nlm.nih.gov/articles/PMC6023095), [AIME-Con 2026 paper](https://aclanthology.org/2026.aimecon-main.5.pdf), [Jurnal UNY](https://jurnal.uny.ac.id/index.php/reid/article/download/73055/22264/216675) (search summary; rules of thumb vary by population and content)
- AutoIRT (Duolingo-affiliated, arXiv 2409.08823, 2024): calibrates new items using item features with AutoML; "jump-start" with only 20 responses per item gave item-difficulty correlation 0.936 vs 0.855 for standard 2PL (Y/N vocab, 3,290 items) and 0.954 vs 0.847 (ViC, 585 items); cold-start (zero responses) was weaker and in one setting BERT-IRT correlated better; traditional IRT needs hundreds of responses per item — [arXiv 2409.08823](https://arxiv.org/html/2409.08823v1) (fetched)
- Duolingo English Test delivers items with Thompson sampling / contextual-bandit adaptivity and explanatory IRT — [NeurIPS 2024 talk](https://neurips.cc/virtual/2024/107597) (search summary)
- FSRS: DSR memory model (Difficulty, Stability, Retrievability; stability = days for recall probability to fall from 100% to 90%); FSRS-6 has 21 parameters whose defaults were trained on several hundred million reviews from about 10k users; if a user has too few reviews the defaults are used; page says defaults already beat SM-2 — [awesome-fsrs wiki, ABC of FSRS](https://github.com/open-spaced-repetition/awesome-fsrs/wiki/ABC-of-FSRS) (fetched). No minimum-review threshold is stated on that page.
- FSRS vs SM-2 benchmark: 9,999 collections, 349,923,850 filtered reviews; FSRS-6 (recency) has lower log loss than SM-2 for 99.6% of users; the author cautions SM-2 was not designed to predict recall probability so the comparison is not fully fair; page does not state the run date — [Expertium benchmark](https://expertium.github.io/Benchmark.html) (fetched). Anki's own FAQ is more modest, saying tests suggest FSRS is roughly on par with SM-17 — [Anki FAQ](https://docs.ankiweb.net/faqs/what-spaced-repetition-algorithm.md) (search summary)
- Half-Life Regression (Duolingo) is a forgetting-curve model; a 2026 paper revisits it for explainability; I found no study isolating its contribution to learning gains — [PMLR v339 Ilidio 2026](https://proceedings.mlr.press/v339/ilidio26a.html) (search summary)
- FSRS is implemented as open-source libraries; confirmed ports were not enumerated in this research (Gap).

### Inferences
- Realistic stack for a solo dev: (a) FSRS-6 default parameters for TWK fact cards/items, scheduling per item or per fact; (b) per-subtopic mastery from a simple model (e.g., Elo/logistic update or BKT with a few parameters) for TIU/TKP drill selection; (c) log all responses from day one so IRT (1PL/Rasch first) can be fit once a few hundred answers per item exist; (d) use LLM/feature-based difficulty estimates as priors per AutoIRT's idea. No source shows a true CAT for a small bank works well; with a small bank and few users, adaptive item selection by estimated mastery and difficulty bands is the feasible form.
- FSRS evidence is for flashcard recall; applying it to MCQ practice (where answer position, guessing and item variants matter) is an extrapolation and should be treated as unvalidated.
- Mastery learning (gate on accuracy threshold per subtopic) was not researched with sources here.

### Gaps
- No source found on FSRS or BKT performance specifically for MCQ test-prep or for datasets of a few hundred students; Gervet et al.'s "moderate" datasets are not precisely quantified in the retrieved summary.
- No source retrieved for effectiveness of mastery-learning gating in test prep, or SM-2 implementation details.
- Minimum reviews before FSRS parameter optimization: not stated in fetched sources.

## 5. Error diagnosis: why an answer was wrong, and AI tutoring on mistakes

### Takeaway
Research on LLM misconception diagnosis exists (Eedi data, algebra error classification) but accuracy drops as the label set grows and LLM tutors do not reliably use diagnostic labels; a modest, rule-assisted taxonomy (careless/time pressure/knowledge gap using response time and history) is more reliable than asking an LLM to infer misconceptions for every wrong answer.

### Cited Findings
- Misconception prediction from an incorrect answer is an established task; LLM diagnostic accuracy drops sharply as the number of candidate misconceptions grows — [arXiv 2605.23925](https://arxiv.org/pdf/2605.23925) (search summary)
- Eedi-based work: an LLM generates a misconception hypothesis, which is matched to labelled misconceptions in embedding space and reranked, rather than asking the model to pick the label directly — [arXiv 2602.02414](https://arxiv.org/html/2602.02414v1) (search summary)
- LLMs choose the wrong answer a student with a given misconception would pick, and name misconceptions, well above random but not fully satisfactorily — [arXiv 2310.02439](https://arxiv.org/pdf/2310.02439v1.pdf) (search summary)
- A hybrid ensemble of fine-tuned Gemma 2/DeepSeek/Gemma 3 with weighted voting reached MAP@3 of 0.68 on an Eedi-style task (preprint, not peer-reviewed) — [Preprints.org 202510.0463](https://www.preprints.org/manuscript/202510.0463) (search summary)
- A 2026 AACL paper reports LLM tutors show little evidence of reliably using given diagnostic labels (simulated students) — [ACL Anthology 2026.aacl-long.112](https://aclanthology.org/2026.aacl-long.112.pdf) (search summary)
- Fine-tuned LLMs classified algebra step errors better than earlier tree-embedding methods — [arXiv 2305.06163](https://arxiv.org/pdf/2305.06163) (search summary)
- Khanmigo RCT (NBER WP 35620, Oreopoulos & Low, Aug 2026; 18 Tennessee middle schools): students messaged the tutor on about a third of practice days and in only 17% of exercise sessions where they made a mistake; most messages were bare answers or suggested-prompt clicks; effect on math about 0.06-0.08 SD per school year, similar to Khan Academy practice without the AI — [NBER](https://www.nber.org/papers/w35620) (fetched)

### Inferences
- Avoid opt-in "ask the tutor" as the main error-feedback mechanism; the Khanmigo trial suggests students rarely use it after mistakes. Show a short pre-computed explanation automatically after every wrong answer instead (inference from the 17% usage figure).
- Simplest reliable diagnosis: tag errors with cheap signals (answer time vs. item median, same-subtopic accuracy, whether the chosen distractor is tied to a known pre-authored misconception). Distractors authored with a specific misconception label make diagnosis a lookup instead of an LLM inference. This is a design suggestion, not a published method.

### Gaps
- No study found on classifying "time pressure vs careless vs knowledge gap" for MCQ exams.
- No evidence found for diagnosing TKP (situational judgment) mistakes.

## 6. Global products: what AI/adaptive features exist and the evidence

### Takeaway
Products mostly use AI for explanation chat and calibration/scheduling, while independent effectiveness evidence is thin: the best independent trial (Khanmigo, NBER 2026) found gains no better than plain practice; Riiid, Duolingo, and UWorld evidence is mostly company-reported.

### Cited Findings
- Khanmigo: see Sections 2 and 5; the Khan Academy Newark case study reports Yearly Proficient Learners gained +6 points on the NJSLA vs +2 state average (platform-level, vendor-run, not Khanmigo-specific) — [Khan Academy Newark case study](https://www.khanacademy.org/schools/case-studies/newark-public-schools) (search summary). Common Sense Media says Khanmigo's math accuracy remains limited — [Common Sense Media](https://www.commonsensemedia.org/ai-ratings/khanmigo) (search summary)
- Riiid Santa TOEIC: company claims a 124-point score gain after about 20 hours, no methodology in the retrieved source — [AI Business](https://aibusiness.com/riiid-disrupting-education-with-ai-one-test-at-a-time) (trade article, company claim). A Riiid-authored study reports A/B tests on more than 20K students improving engagement factors by up to 25.13% (engagement, not scores) — [arXiv 2009.08976](https://arxiv.org/pdf/2009.08976.pdf). Riiid's SAINT transformer knowledge-tracing model was reported at NeurIPS 2020 on the EdNet dataset, per a third-party summary only (unverified; the primary paper was not opened) — [Riiid analysis, third-party](https://deepaksood619.github.io/education/competitors/riiid-analysis)
- Duolingo: Birdbrain tailors lesson difficulty; efficacy claims (five units of French/Spanish matching four university semesters on reading/listening) come from Duolingo's own studies, with critics noting the speaking gap — [Duolingo efficacy blog](https://blog.duolingo.com/results-duolingo-efficacy-studies) (search summary)
- Duolingo English Test: generative-AI item creation with human review, IRT calibration (AutoIRT), adaptive administration — see Section 1 and 4 sources.
- UWorld: AI study assistant "UAsk" has been added across product lines (NCLEX-RN, FNP, CPA, CFA, medical); CPA UAsk was beta-tested with more than 5,000 candidates; UWorld's FAQ says it may occasionally make mistakes (via third-party review); UAsk sits on top of expert-written content — [UWorld CFA UAsk](https://finance.uworld.com/cfa/uask/), [AI Journal press release](https://aijourn.com/?p=541057), [iatroX review](https://www.iatrox.com/blog/uworld-usmle-step-3-ai-tutor-review-reasoning-support-source-grounding-and-exam-fidelity) (third-party, lower authority). Sources disagree on whether UWorld has adaptive difficulty; unresolved — [BenchPrep review](https://ai.benchprep.com/resources/high-stakes-exam-preparation-platform-reviews-2026) vs [iatroX](https://www.iatrox.com/blog/uworld-alternative-ai-adaptive-usmle)
- Magoosh: little authoritative 2026 detail found; one listing claims an AI tutor for SAT/ACT — [Gaurav Tiwari list, unverified](https://gauravtiwari.org/best-test-prep-websites/)
- Note: UWorld's core value is expert-written items with detailed explanations; AI is layered on top of them, per the sources above.

### Inferences
- The pattern across UWorld, Khan, and DET is human-authored or human-reviewed content with AI on top, not unreviewed LLM generation. A solo developer matching this should keep an authored/verified bank and use the LLM for generation-with-review, explanation text, and on-demand help.
- Effectiveness claims by product vendors should not be used as marketing benchmarks for this app.

### Gaps
- No independent effectiveness study found for Riiid Santa, UWorld AI features, Magoosh, SAT/GMAT/LSAT adaptive prep; official product pages for Magoosh/UWorld were not fetched.
- No info on Indonesian CPNS/BUMN prep apps' AI features (outside this note's scope).

## 7. Cost: pre-generating and verifying 5,000 items vs on-the-fly

### Takeaway
Pre-generating and verifying 5,000 items costs on the order of $4 to $340 depending on model tier (my calculation from official prices), i.e., trivial next to human review time; on-the-fly generation per user multiplies cost with every request and removes the chance to verify, so pre-generation is clearly cheaper and safer.

### Cited Findings
Official per-million-token prices (USD) as read on 2026-11-14:
- Anthropic: Haiku 5.5 $0.10 in / $0.50 out (prompts up to 100k tokens); Sonnet 5.5 $2 / $10; Opus 5.5 $4 / $20; Batch API gives a 50% discount on input and output (Haiku 5.5 batch $0.05 / $0.25; Sonnet 5.5 batch $1 / $5; Opus 5.5 batch $2 / $10); cache reads are 0.05x-0.1x of input price depending on model (Sonnet 5.5 cache hit $0.10) — [Claude pricing](https://platform.claude.com/docs/en/about-claude/pricing)
- Google Gemini API: 3.1 Flash-Lite $0.25 in / $1.50 out (batch $0.125 / $0.75); 3.5 Flash-Lite $0.30 / $2.50; 3.6/3.8 Flash $0.75 / $3.75 until 2026-12-31 then $1.50 / $7.50 from 2027-01-01; 3.1 Pro Preview $2 / $12; free tier exists for the Flash/Flash-Lite models (not Pro 3.1, not batch) — [Gemini API pricing](https://aistudio.google.com/docs/pricing). Free-tier limits and data-use terms were not checked.
- DeepSeek: deepseek-flash $0.15 in (cache miss, off-peak) / $0.60 out off-peak, double at peak; deepseek-v4-pro $0.66 / $1.98 off-peak (cache-miss input); off-peak is half of peak; peak is 01:00-04:00 and 06:00-10:00 UTC Mon-Fri — [DeepSeek pricing](https://api-docs.deepseek.com/quick_start/pricing)

My calculation (assumption-based, not sourced): per item, generation with RAG context about 1,500 in / 500 out; 3 independent solver passes about 1,200 in / 1,800 out total; explanation plus critic pass about 800 in / 400 out. Total about 3,500 in / 2,700 out per item, so for 5,000 items about 17.5M input and 13.5M output tokens. Hidden "thinking" tokens are not included and could raise output several times.

| Model (all stages) | Standard | Batch |
|---|---|---|
| Claude Haiku 5.5 | about $8.5 | about $4.3 |
| Gemini 3.1 Flash-Lite | about $24.6 | about $12.3 |
| DeepSeek flash (off-peak) | about $10.7 | n/a |
| DeepSeek v4-pro (off-peak) | about $38.3 | n/a |
| Claude Sonnet 5.5 | about $170 | about $85 |
| Claude Opus 5.5 | about $340 | about $170 |

On-the-fly comparison (same assumption, one item with verification): Haiku 5.5 about $0.0017 per item; Sonnet 5.5 about $0.034 per item. Example: 100 daily-active users x 30 items x 30 days = 90,000 items per month, about $153 (Haiku) or $3,060 (Sonnet), recurring, versus a one-time $4 to $170 for a pre-generated 5,000-item bank. On-demand explanations of mistakes only (about 800 in / 400 out): about $0.0003 per call on Haiku 5.5, about $0.0056 on Sonnet 5.5 (my calculation).

### Inferences
- A mixed pipeline (cheap model for bulk generation and multi-pass solving, stronger model only for items where solvers disagree) should land near the low end; I did not test this.
- The real cost is the human review time, not API fees.
- Model names/prices change often (Gemini 3.8 Flash price doubles 2027-01-01 per the page); recheck before budgeting.

### Gaps
- No measured token counts from a real run; the table depends on my token assumptions.
- No published per-item verification pipeline cost for LLM-generated exam items found.
- Free-tier limits for Gemini and any Indonesian-language quality comparison across these models were not checked.
