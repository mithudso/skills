<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** This is the **applied VLM-critique technique** spoke: how to drive a
> vision-capable LLM (GPT-4o/4V, Claude vision, Gemini) to *critique a design or UI screenshot* — prompt patterns, structured findings,
> region grounding, before/after comparison, judge↔human calibration, the failure modes, and the reliability practices that make the
> output trustworthy. It is **not** about how the model sees: vision encoders, projectors, fusion, tiling, training stages, and the
> benchmark internals (MMMU/POPE/MMHal at the architecture level) live in `ai-llm-model-layer` (`references/multimodal-llm-architecture.md`)
> — cross-reference that, do not duplicate it. The *general-text* LLM-as-judge calibration theory (kappa, binary-vs-Likert, position/
> verbosity/self-preference bias) lives in the eval family (`ai-agents-orchestration` → `references/eval-driven-development.md`); this
> spoke specializes it to the **visual** case. Loop stop-conditions & self-consistency mechanics → `references/iterative-self-refinement-loops.md`.

---

---
name: vision-model-design-critique
title: "Vision-Model / Multimodal-LLM Design Critique Technique"
description: >-
  Using a vision-capable LLM (GPT-4o/4V, Claude vision, Gemini) to critique a design
  or UI screenshot. TRIGGER: prompt a vision model to review/grade a UI/mockup/graphic
  against a rubric (Nielsen heuristics, hierarchy, contrast, a11y); critique findings
  as structured JSON (severity/heuristic/element/fix); ground findings to a region
  (Set-of-Mark, boxes); before/after or A/B multi-image comparison; trust/calibrate a
  VLM "visual judge" vs human designers; diagnose a critique that hallucinated an
  element, misread UI text, botched spatial relations, or just agreed with you
  (sycophancy); make a vision-critique pipeline reliable (anchors, describe-then-judge,
  self-consistency, panels, abstention). SKIP: how the model encodes pixels —
  encoders/tiling/training + MMMU/POPE/MMHal internals → multimodal-llm-architecture.md;
  generic text LLM-as-judge calibration → eval-driven-development.md; the human practice
  of UI critique / WCAG → frontend-ui; generic CoT/few-shot/JSON-mode → prompt-engineering.
origin: local
category: developer
version: "1.0.0"
updated: "2026-06-16"
tags:
  - vlm-as-judge
  - mllm-as-judge
  - design-critique
  - ui-ux-evaluation
  - visual-grounding
  - set-of-mark
  - object-hallucination
  - structured-output
keywords:
  - vision model design critique
  - VLM-as-judge
  - MLLM-as-judge
  - GPT-4V UI evaluation
  - heuristic evaluation with LLM
  - structured JSON findings
  - Set-of-Mark prompting
  - before/after design comparison
  - object hallucination POPE
  - describe-then-judge
whenToUse:
  - "prompt GPT-4o / Claude / Gemini to critique a UI screenshot or mockup against a rubric"
  - "get design-critique findings back as schema-guaranteed JSON (severity, heuristic, element, fix)"
  - "ground each finding to a region so a human can verify it (Set-of-Mark, bounding boxes)"
  - "run a before/after or A/B design comparison with multiple images in one prompt"
  - "decide how far to trust a VLM visual judge vs a human designer"
  - "debug a vision-model critique that hallucinated, misread text, or just agreed with me"
  - "make a vision-critique pipeline reliable (anchors, self-consistency, panels, abstention)"
related_skills:
  - multimodal-llm-architecture
  - eval-driven-development
  - prompt-engineering
  - iterative-self-refinement-loops
  - accessibility-ux-reviewer
metadata:
  changelog:
    - "2026-06-16 /dr v1.0.0 — initial research build (≈50 sources; 9 concepts; verify-claims gate ON)"
---

# Vision-Model / Multimodal-LLM Design Critique Technique

**`verified-as-of: 2026-06-16`** for all model-specific behavior, vendor API parameters, and benchmark scores below (a fast-moving area — re-verify version-pinned claims on refresh).

## Overview

A *vision-capable* LLM can be pointed at a design — a UI screenshot, a Figma frame, a landing page, a brand asset — and asked to **critique** it: find usability problems, rate it against a rubric, compare two variants, suggest fixes. This is the multimodal analogue of `document-critique`, and it is the implementation backbone for any "upload a screenshot, get a design review" skill.

The technique works, but only inside a narrow envelope. The dominant empirical finding across 2024-2026 is that vision LLMs are **good at *ranking/comparing* and *generating plausible critique text*, and poor at *absolute scoring*, *fine-grained spatial perception*, and *not making things up*.**[^mllmjudge][^vlmrank] A VLM critique is a fast, cheap *first-pass triage and idea generator*, not a calibrated measurement and not a replacement for a human designer or a real usability test. The job of this spoke is to extract the useful signal while suppressing the well-documented failure modes.

**The boundary you must hold:** this spoke is the *applied prompting + evaluation technique*. The model-internal "how does it see pixels" story — vision encoders, the projector, fusion, high-res tiling, training stages, and the benchmark suite at the architecture level — is owned by `multimodal-llm-architecture.md`. When a critique fails because the model literally cannot resolve small text, the *fix mechanism* (resolution/tiling/encoder) is explained there; the *prompt-level workaround* (escalate `detail`, external OCR pre-pass, abstention) is here.

## Core concepts

### 1. The central calibration result: rank ≫ score
The foundational benchmark is **MLLM-as-a-Judge** (Chen et al., ICML 2024 Oral), which evaluates GPT-4V and peers across three tasks — Scoring Evaluation, Pair Comparison, Batch Ranking. The headline: MLLMs show "remarkable human-like discernment in **Pair Comparison**" but "significant divergence from human preferences in **Scoring Evaluation and Batch Ranking**," alongside persistent bias, hallucination, and inconsistency.[^mllmjudge] A 2026 quantification reports SOTA VLM judges hitting only **~32-34% exact agreement** with human ratings on a 5-point scale, 24-30% of predictions off by ≥2 points, and Spearman ρ only ~0.30-0.46 at the pointwise level.[^vlmrank] **Practical consequence: prefer A-vs-B comparisons and relative rankings; treat any absolute 1-10 "design score" as low-confidence.**

The split is consistent and task-dependent, not uniformly bad: **comparative + coarse = human-aligned; pointwise + fine-grained/aesthetic = poorly calibrated.** Pairwise/Elo evaluators like GPTEval3D (CVPR 2024) and VisionPrefer (NeurIPS 2024) report strong human alignment precisely because they are comparative.[^gpteval3d][^visionprefer]

### 2. Binary/pairwise beats Likert — and *why* (it's a human fact, not just a model quirk)
The text-LLM-judge rule "binary/pairwise beats Likert" carries into the visual case; MJ-Bench shows multimodal/text-to-image judges give more accurate, stable feedback on a Likert scale than on raw numerical scores.[^mjbench] The deeper cause: on the **Visual Aesthetic Benchmark**, with 8 expert annotators, *direct comparative ranking* yielded ~42 percentage-points higher inter-annotator agreement than *score-derived* rankings — pointwise scoring is a worse elicitation method even for humans.[^vab] So: decompose into yes/no checks and A-vs-B comparisons rather than asking for a number on a scale.

> **Contested nuance (preserve it):** a counter-current argues *direct* pairwise can *amplify* verbosity/style/position bias relative to pointwise, motivating hybrids that put pointwise reasoning *inside* a pairwise decision (PREPAIR) or calibrate pairwise (PairS).[^prepair][^pairs] The dominant practice is still "compare, don't score," but pairwise is not bias-free.

### 3. The design/UI case specifically: useful generator, weak detector
When studied as **heuristic evaluation / critique generation**, vision LLMs align *poorly* with human experts on *issue identification* but are *useful idea generators*:
- **UICrit** (Duan et al., UIST 2024): a dataset of 3,059 critiques over 983 mobile UIs; using it lifted LLM UI feedback ~55% over zero-shot — i.e., the raw model is weak and needs grounding/examples.[^uicrit]
- **GPT-4o usability evaluation** (INTERACT 2025): GPT-4o caught only **~21% of issues** human experts found, while generating many false positives via hallucination.[^interact] A Nielsen-cited replication found two AI tools overlapped on only ~20% of problems.[^nielsen]
- **LLM-as-Design-Critic** on UICrit data: human critiques predicted design quality well (Spearman ≈ 0.556) while LLM-only critiques were far weaker (≈0.194), with low human↔LLM semantic alignment.[^designcritic] *Single-source, low-tier journal — treat as directional.*

**Read this as: the VLM surfaces candidate issues fast and writes them up well, but it misses most real problems and invents some. Use it to *widen* a human's review, never to *replace* it or to *gate* a release on a score.**

## Tools / frameworks (techniques, not products)

| Technique | What it does | When to reach for it |
|---|---|---|
| **Describe-then-judge** (analyze-then-judge) | Force the model to *describe visible elements first*, then map them to criteria, then emit a verdict | Default for every critique — the single highest-leverage move |
| **Rubric-as-prompt** | Paste the actual criteria definitions + per-level anchors into the prompt; tell it what to ignore | Any heuristic / scored evaluation |
| **Structured Outputs (JSON schema)** | Schema-guaranteed findings list | When findings feed code/a dashboard — emit *after* free reasoning |
| **Set-of-Mark (SoM)** | Overlay numbered marks on regions; model references region IDs | Grounding findings to verifiable locations |
| **OCR + icon-detector pre-pass** (OmniParser-style) | External detector supplies boxes + captions; VLM reasons over them | Production UI grounding; dense/small text |
| **Self-consistency** | Sample N, take majority/median | Borderline verdicts; cheap reliability (diminishing returns — see anti-patterns) |
| **LLM-as-jury (panel)** | Several *diverse-family* judges vote | High-stakes critique; cancels single-model bias |
| **Order-swap** | Present A/B in both orders, discard order-inconsistent verdicts | Any multi-image comparison |
| **Abstention option** | Allow "not visible / cannot determine" per finding | Suppresses fabricated findings everywhere |

## Methodology — a reliable vision-critique prompt

A robust critique request layers these, in order:

1. **Place the image and its instructions tightly together.** Anthropic's vision docs: Claude "works best when images come before text" — put the image earlier than the question about it.[^anthropicvision] *Exception:* for **coordinate/click-grounded** tasks, Anthropic's computer-use guidance recommends **text-instruction-first** (tell it what to look for before it looks).[^anthropiccu] OpenAI practitioner guidance also favors instructions-adjacent. **Rule: keep them grouped; image-then-text for Claude critique, text-first for grounding.**
2. **Set resolution deliberately.** OpenAI's `detail` knob (`low`/`high`/`auto`) is the cost/perception lever: `low` downsamples to 512×512 (~85 tokens, fine for layout triage but **misses small text and subtle defects**); `high` tiles at 512px and is the regime screenshots/forms/tables/OCR-style critique actually need.[^openaivision] Start `low` for first-pass classification, escalate to `high` when small text or fine defects matter.
3. **Give it the rubric, grounded in *supplied* definitions.** Don't say "review heuristically" — each model has a different internalized Nielsen and you get generic, inconsistent output. Paste the actual heuristic definitions and instruct "use only these definitions, not your prior knowledge" to create a controlled vocabulary so findings are comparable.[^groundrubric] Define each scale point concretely (what a 1 vs 3 vs 5 looks like).
4. **Describe before you judge.** Require a visible-elements description *first*, then per-criterion analysis, then the verdict — the score is then conditioned on self-generated evidence, which "significantly reduces blind approvals."[^analyzejudge] A GPT-4V study reported ~50% improvement from a "describe first, then answer" CoT prefix.[^desc2dec]
5. **Allow abstention.** Tell the model it may answer "not visible / cannot determine" and to omit low-confidence findings rather than invent violations. Anthropic's reduce-hallucinations guidance explicitly lists "allow Claude to say 'I don't know'" and "ask for evidence before answering."[^anthropichalluc] OpenAI argues accuracy-only scoring rewards guessing and that abstention deserves credit.[^openaihalluc] Uncertainty-based abstention avoids ~50% of hallucinations on unanswerable items.[^abstain]
6. **Emit JSON last (reason first, then format).** All three vendors offer schema-guaranteed output: OpenAI Structured Outputs (`response_format: {type:"json_schema", strict:true}`, plus a `refusal` field), Gemini `responseMimeType:"application/json"` + `responseSchema`, Anthropic via tool-use.[^openaistruct][^geministruct] But forcing structure *during* reasoning degrades it (the "format tax"; see Anti-patterns), so let the model critique freely, then serialize. A proven findings schema: a list of `{severity, heuristic, element, description, recommendation, region_id?, confidence}`.[^uxcode]
7. **Ground each finding to a region.** See next section — so a human can verify "the contrast issue is *here*," not hunt for it.
8. **Low temperature + (for stakes) self-consistency or a small diverse panel.** Temperature 0 for determinism; sample-N majority or a 3-judge panel for borderline calls — with the caveats in anti-patterns.

### Grounding findings to regions
Raw coordinate output from a general VLM is **unreliable for precise UI localization** — GPT-4o/GPT-4V "consistently struggle to determine exactly where," because transformers do semantic prediction, not coordinate arithmetic, and API boxes are often "off" or clipped by the tiling/downscale pipeline.[^gptbox][^cnnvlm] Two robust patterns:
- **Set-of-Mark (SoM)** (Yang et al., 2023): overlay numbered marks (from SAM/SEEM segmentation) on regions so the model says "mark 9" instead of emitting pixels; SoM-prompted GPT-4V beat fully-finetuned referring-expression models zero-shot.[^som] **Constraint:** marks must be legible and uncluttered — tiny marks on a 4K screenshot collapsed accuracy (5% vs 39.6%) until glyphs were scaled to ~2% of image height; degrades past ~3-10 entities.[^somfail] Cap mark density, scale glyphs, place labels outside the element.
- **OCR + icon-detector pre-pass** (OmniParser-style, Microsoft): a dedicated detector + OCR supplies boxes and *functional captions*, then SoM with unique IDs; the captions let the model match on text even when visuals fail.[^omniparser] This is the production pattern for UI grounding.
- **Points > boxes:** even grounding-trained VLMs (Molmo's pointing, Qwen-VL, Gemini boxes) localize *points near* targets more reliably than they draw accurate *boxes*.[^regionfocus][^molmo] If you only need "where roughly," ask for a point.

> Gemini natively documents a `box_2d` field (normalized `[ymin,xmin,ymax,xmax]` 0-1000) in its image-understanding docs[^geminivision] — usable for *coarse* regions, still below a specialized detector for precision.

### Multi-image & before/after comparison
All major vendors accept multiple images per prompt with documented limits — Anthropic up to ~100 images (200k-context models; 20 on claude.ai; max 8000×8000px), Azure-hosted OpenAI ~10/call, Gemini up to thousands.[^anthropicvision][^azurevision] The universal practice:
- **Label and order images explicitly** — "Image 1:", "Image 2:" (or "Image A = before, Image B = after") and reference those labels — to prevent **cross-image attention leakage** (attributes from one image bleeding into another).[^multiimg]
- **Randomize or dual-order** because multi-image VLMs have measurable **image-order/position bias**: swapping order changes predictions even in GPT-4o, with ~30% inconsistency; present A/B in both orders and discard/flag order-inconsistent verdicts.[^posbias] (This is the *practice*; the calibration math lives in the eval spoke.)

## Practical patterns

- **Triage, then human.** Run the VLM as a fast first pass to *widen* a human review (surface candidate issues, draft the write-up), never as the gate. It finds ~1 in 5 real issues and invents some.[^interact][^nielsen]
- **Compare, don't score.** Frame as "which of these two is better on contrast/hierarchy, and why," or yes/no heuristic checks — not "rate this 1-10."[^mllmjudge][^vab]
- **Ground every finding.** SoM region ID or an external-detector box on each finding so it's verifiable; this also makes false positives obvious to the human reviewer.[^som][^omniparser]
- **Describe-then-judge, JSON-last, abstention-on** as the default prompt skeleton.[^analyzejudge][^anthropichalluc]
- **External OCR for any text-critical check** (contrast ratios, label wording, truncation) — feed recognized text alongside the pixels rather than trusting the model to read small UI text.[^omniparser]
- **Diverse panel for stakes.** A panel of 3 *different-family* judges (PoLL) beats a single strong judge, reduces intra-model bias, and is ~7× cheaper than one GPT-4-class call.[^poll] Re-running *one* model is a repeated trial, not an ensemble.
- **Cross-references (hold the boundary):** the resolution/encoder *reason* a text-read fails → `multimodal-llm-architecture.md` §High-resolution tiling; loop stop-conditions for an iterate-until-clean critique → `iterative-self-refinement-loops.md`. The rubric/heuristics *themselves* — Nielsen's 10, the 0-4 severity scale, the Laws of UX, WCAG — and *human-run* heuristic evaluation belong to `frontend-ui` (`usability-heuristics-laws-of-ux.md`, `accessibility-ux-reviewer.md`). **This spoke is the technique of making a vision MODEL apply that instrument; that spoke is the instrument itself.**

## Anti-patterns

- **Trusting an absolute design score.** Pointwise VLM scores are poorly calibrated (ρ≈0.3-0.46, ~32% exact agreement).[^mllmjudge][^vlmrank] Never gate a release on "GPT-4o gave it 7/10."
- **Stating your opinion in the prompt.** Sycophancy is *worse* in the visual modality — a "sycophantic modality gap" — and leading/deceptive queries significantly *increase* hallucination across LVLMs; one (medical-domain) benchmark, EchoBench, measured >60% sycophancy for many models, with even strong models at 45-59%.[^syco][^sycolvlm][^echobench] Ask "what's wrong with this design," **not** "this design is great, right?"
- **Asking for precise spatial judgments.** Fine-grained spatial reasoning (left/right, alignment, overlap, orientation) is near-random: VSR sits >25pts below human; BLINK best models ~45-51% vs ~96% human; SpatialEval shows models "fall behind random guessing" and even *under*-perform their own text backbones.[^vsr][^blink][^spatialeval] Don't ask "is this 8px or 12px off"; ask a human or measure it.
- **Trusting element counts.** Counting degrades monotonically past the subitizing range (~4-10); accuracy can collapse from 93% to <10% as count rises (model-dependent).[^count] Don't ask "how many CTAs are above the fold."
- **Trusting the model to read small/dense/rotated text.** OCRBench shows LMMs do basic recognition but fail complex layouts; rotating a doc cut one model 90.9%→35.2%; GPT-4 Vision couldn't match its own scores pixel-only without an external OCR engine.[^ocrbench][^gpt4doc] Use an OCR pre-pass.
- **Believing emitted bounding boxes are precise.** They're coarse/often invented for UI; use SoM or an external detector.[^gptbox][^cnnvlm]
- **Forcing JSON during reasoning.** The "format tax": requiring structured output *while* generating hurts reasoning-heavy tasks and produces *valid JSON with wrong values*, so a schema-validity dashboard masks quality regressions. Reason first, serialize second; penalty is largest for weaker models.[^formatfree][^formattax]
- **Object/co-occurrence hallucination unchecked.** Models confirm objects that *co-occur* in training (a chair when only a table is present) and skew to "yes" answers (POPE, AMBER); spurious/salient cues amplify hallucination ~26× (SpurLens).[^pope][^amber][^spurlens] Force grounded description + abstention; don't accept a finding without a region.
- **Stacking many identical/correlated judges and calling it robust.** Self-consistency gains plateau early and can decline on capable models; correlated judge errors cap "effective votes" at ~2-2.5 (9 judges ≈ 2 votes), and best-of-N single-model often beats multi-agent deliberation at far better cost.[^selfcons][^ninejudges][^deliberation] Use a *small diverse* panel, not a big homogeneous one.
- **Re-explaining the encoder/training here.** That's `multimodal-llm-architecture.md`'s job — cross-reference, don't duplicate.

## Troubleshooting

- **"It described a button/section that isn't there."** Object hallucination + co-occurrence/affirmative bias.[^pope][^amber] Add describe-then-judge, require a region per finding, allow "not visible," lower temperature, and verify against an external detector's element list.
- **"It missed obvious problems a designer caught instantly."** Expected — VLMs find ~1 in 5 real issues.[^interact] Use it to widen, not replace, human review; add a UICrit-style few-shot exemplar set.[^uicrit]
- **"Its spacing/alignment/pixel-offset judgments are wrong."** Fine-grained spatial reasoning is near-random.[^vsr][^blink] Don't ask the model; measure programmatically or ask a human.
- **"It can't read the label / got the contrast wrong."** Resolution/OCR limit. Escalate `detail:high`/tiling, crop to the region, or run an external OCR pre-pass and feed the text.[^openaivision][^ocrbench] (Encoder/resolution *mechanism* → `multimodal-llm-architecture.md`.)
- **"Scores swing run-to-run / flip when I reorder the images."** Self-inconsistency + position/order bias.[^posbias][^selfcons] Temperature 0, self-consistency majority, and order-swap with order-inconsistent verdicts discarded.
- **"It just agreed with my framing."** Sycophancy.[^syco][^sycolvlm] Strip opinions from the prompt; ask neutrally or adversarially ("list what a harsh reviewer would flag").
- **"Valid JSON, wrong content."** Format tax.[^formattax] Move reasoning before serialization; validate *content* against grounded evidence, not just schema.

## References

[^mllmjudge]: Chen et al., "MLLM-as-a-Judge: Assessing Multimodal LLM-as-a-Judge with Vision-Language Benchmark," ICML 2024 (Oral), PMLR 235:6562-6595; arXiv:2402.04788. https://proceedings.mlr.press/v235/chen24h.html · project: https://mllm-judge.github.io/
[^vlmrank]: "VLM Judges Can Rank but Cannot Score: Task-Dependent Uncertainty in Multimodal Evaluation," arXiv:2604.25235. https://arxiv.org/html/2604.25235
[^mjbench]: Chen, Wu, Chen et al., "MJ-Bench: Is Your Multimodal Reward Model Really a Good Judge for Text-to-Image Generation?," NeurIPS 2025 D&B; arXiv:2407.04842. https://arxiv.org/abs/2407.04842
[^vab]: Feng, Li, Liu et al., "Visual Aesthetic Benchmark: Can Frontier Models Judge Beauty?" arXiv:2605.12684. https://ar5iv.labs.arxiv.org/html/2605.12684
[^gpteval3d]: Wu, Yang, Li et al., "GPT-4V(ision) is a Human-Aligned Evaluator for Text-to-3D Generation," CVPR 2024. https://openaccess.thecvf.com/content/CVPR2024/html/Wu_GPT-4Vision_is_a_Human-Aligned_Evaluator_for_Text-to-3D_Generation_CVPR_2024_paper.html
[^visionprefer]: Wu et al., "VisionPrefer: Multimodal LLMs Make Text-to-Image Generative Models Align Better," NeurIPS 2024; arXiv:2404.15100. https://arxiv.org/html/2404.15100v1
[^prepair]: Jeong et al., "The Comparative Trap: Pairwise Comparisons Amplify Biased Preferences of LLM Evaluators (PREPAIR)," BlackboxNLP 2025. https://aclanthology.org/2025.blackboxnlp-1.5.pdf
[^pairs]: Liu et al., "Aligning with Human Judgement: The Role of Pairwise Preference in LLM Evaluators (PairS)," arXiv:2403.16950. https://arxiv.org/pdf/2403.16950
[^uicrit]: Duan, Cheng, Li, Hartmann, Li, "UICrit: Enhancing Automated Design Evaluation with a UI Critique Dataset," UIST 2024; doi:10.1145/3654777.3676381. https://dl.acm.org/doi/fullHtml/10.1145/3654777.3676381
[^interact]: Campos et al., "Can GPT-4o Evaluate Usability Like Human Experts? A Comparative Study on Issue Identification in Heuristic Evaluation," INTERACT 2025; arXiv:2506.16345. https://arxiv.org/html/2506.16345
[^nielsen]: Nielsen, "UX Roundup" (AI heuristic-evaluation overlap ~20%), 2025-10-27. https://jakobnielsenphd.substack.com/p/ux-roundup-20251027
[^designcritic]: "LLM-as-Design-Critic: Aligning AI-Generated UI Feedback with Human Graphic Design Judgment," Intl. Journal of Graphic Design, 2025 (single-source, low-tier; directional). https://journal.stekom.ac.id/index.php/ijgd/article/view/3661
[^anthropicvision]: Anthropic, "Vision — Claude API Docs" (official). https://docs.anthropic.com/en/docs/build-with-claude/vision
[^anthropiccu]: Anthropic, "Best practices for computer and browser use with Claude" (official; text-before-image for click tasks). https://claude.com/blog/best-practices-for-computer-and-browser-use-with-claude
[^openaivision]: OpenAI, "Images and vision | OpenAI API" (official; `detail` low/high/auto, tiling token math). https://developers.openai.com/api/docs/guides/images-vision
[^groundrubric]: Sojudi, "AI Panel for Heuristic Evaluations" (grounding rubric in supplied definitions), Medium 2026. https://medium.com/@ladansojudi2/ · corroborated by "Beyond the Illusion of Consensus: Knowledge-Grounded Evaluation in LLM-as-a-Judge," arXiv:2603.11027. https://arxiv.org/html/2603.11027
[^analyzejudge]: Kang, "How to Build Reliable Multimodal AI Evaluators Using VLM Judges" (Analyze-then-Judge), Medium 2026. https://medium.com/@jiyang.kang/how-to-build-reliable-multimodal-ai-evaluators-using-vlm-judges-ca5663e3272a
[^desc2dec]: "Description-then-Decision: CoT for VLMs," arXiv:2311.09193. https://arxiv.org/pdf/2311.09193
[^anthropichalluc]: Anthropic, "Reduce hallucinations — Claude API Docs" (official; allow "I don't know," ask for evidence). https://platform.claude.com/docs/en/test-and-evaluate/strengthen-guardrails/reduce-hallucinations
[^openaihalluc]: OpenAI, "Why language models hallucinate" (official; abstention vs guessing), 2025. https://openai.com/index/why-language-models-hallucinate/
[^abstain]: "Uncertainty-Based Abstention in LLMs Improves Safety and Reduces Hallucinations," arXiv:2404.10960. https://arxiv.org/html/2404.10960v1
[^openaistruct]: OpenAI, "Structured model outputs | OpenAI API" (official; `json_schema` strict, `refusal` field). https://developers.openai.com/api/docs/guides/structured-outputs · "Introducing Structured Outputs in the API," 2024. https://openai.com/index/introducing-structured-outputs-in-the-api/
[^geministruct]: Google, "Structured outputs (generateContent)" (official; `responseMimeType`/`responseSchema`). https://ai.google.dev/gemini-api/docs/structured-output
[^uxcode]: "Catching UX Flaws in Code: Leveraging LLMs to Identify Usability Flaws at the Development Stage" (GPT-4o Nielsen→JSON pipeline: SeverityRating/IssueFound/Recommendation; Cohen's κ≈0.50, severity Krippendorff's α≈0), arXiv:2512.04262. https://arxiv.org/html/2512.04262
[^gptbox]: "Why GPT Vision Struggles with Bounding Boxes (and How We Fixed It)," Silversky 2026. https://silverskytechnology.com/why-gpt-vision-struggles-with-bounding-boxes-and-how-we-fixed-it/ · OpenAI Dev Community, "GPT API can not do image coordinates right," 2025. https://community.openai.com/t/gpt-api-can-not-do-image-coordinates-right/1363453
[^cnnvlm]: "CNN + VLM > VLM" (VLM boxes coarse/invented; detector boxes pixel-tight), Interfaze 2026. https://interfaze.ai/blog/cnn-plus-vlm-more-than-vlm
[^som]: Yang, Zhang, Li, Zou, Li, Gao, "Set-of-Mark Prompting Unleashes Extraordinary Visual Grounding in GPT-4V," arXiv:2310.11441. https://arxiv.org/abs/2310.11441 · project: https://som-gpt4v.github.io/
[^somfail]: SoMatic commit (scale annotation font by image resolution; 5% vs 39.6% legibility failure), GitHub. https://github.com/Smyan1909/SoMatic · "Graph-of-Mark," arXiv:2603.06663 (SoM 3-10 entity sweet spot, degradation). https://ar5iv.labs.arxiv.org/html/2603.06663
[^omniparser]: Lu et al. (Microsoft), "OmniParser for Pure Vision Based GUI Agent," arXiv:2408.00203. https://arxiv.org/pdf/2408.00203 · repo: https://github.com/microsoft/OmniParser/
[^regionfocus]: Luo et al., "Visual Test-time Scaling for GUI Agent Grounding / RegionFocus" (points reliable, boxes not), ICCV 2025. https://openaccess.thecvf.com/content/ICCV2025/papers/Luo_Visual_Test-time_Scaling_for_GUI_Agent_Grounding_ICCV_2025_paper.pdf
[^molmo]: Ai2, "Molmo / PixMo-Points" (pointing as pixel-grounded explanation). https://allenai.org/blog/molmo · Qwen-VL: Bai et al., arXiv:2308.12966. https://arxiv.org/pdf/2308.12966v2
[^geminivision]: Google, "Image understanding (generateContent)" (official; `box_2d` normalized 0-1000). https://ai.google.dev/gemini-api/docs/image-understanding
[^azurevision]: Microsoft Learn, "How to use vision-enabled chat models" (Azure OpenAI; 10 images/call, 20MB, detail levels). https://learn.microsoft.com/en-us/azure/ai-foundry/openai/how-to/gpt-with-vision
[^multiimg]: "Multimodal Prompting Patterns" (cross-image confusion fix: label "Image 1/2"). https://wiki.charleschen.ai/ai/processed/wiki/llm-core/prompts/raw/web/multimodal-prompting-patterns · corroborated by SurePrompts "Multimodal AI Prompting Guide," 2026. https://sureprompts.com/blog/ai-multimodal-prompting-complete-guide-2026
[^posbias]: "Identifying and Mitigating Position Bias of Multi-image Vision-Language Models," arXiv:2503.13792 / CVPR 2025 (order swap changes predictions; ~30% inconsistency; U-shaped vs recency). https://arxiv.org/html/2503.13792
[^poll]: Verga et al., "Replacing Judges with Juries (PoLL): panel of small diverse judges beats GPT-4, ~7× cheaper, less intra-model bias," arXiv:2404.18796. https://arxiv.org/html/2404.18796v2
[^syco]: Rahman et al., "Pointing to a Llama and Call it a Camel: On the Sycophancy of Multimodal LLMs" (sycophantic modality gap), EMNLP 2025. https://aclanthology.org/2025.emnlp-main.1020/
[^sycolvlm]: "Sycophancy in Vision-Language Models: A Systematic Analysis and Inference-Time Mitigation" (leading queries increase hallucination), arXiv:2408.11261. https://arxiv.org/pdf/2408.11261
[^echobench]: "EchoBench: Benchmarking Sycophancy in Medical Large Vision-Language Models" (medical-domain benchmark; >60% sycophancy for many models; Claude 3.7 45.98%, GPT-4.1 59.15%), arXiv:2509.20146. https://arxiv.org/html/2509.20146
[^vsr]: Liu, Emerson, Collier, "Visual Spatial Reasoning (VSR)," TACL 2023; arXiv:2205.00363. https://arxiv.org/abs/2205.00363 · Kamath et al., "What's 'up' with vision-language models?," EMNLP 2023. https://aclanthology.org/2023.emnlp-main.568.pdf
[^blink]: Fu, Hu et al., "BLINK: Multimodal Large Language Models Can See but Not Perceive," ECCV 2024; arXiv:2404.12390. https://arxiv.org/abs/2404.12390
[^spatialeval]: Wang, Ming, Shi et al., "Is A Picture Worth A Thousand Words? Delving Into Spatial Reasoning for VLMs (SpatialEval)," NeurIPS 2024. https://openreview.net/forum?id=cvaSru8LeO
[^count]: Paiss et al., "Teaching CLIP to Count to Ten (CountBench)," ICCV 2023. https://openaccess.thecvf.com/content/ICCV2023/papers/Paiss_Teaching_CLIP_to_Count_to_Ten_ICCV_2023_paper.pdf · "Your Vision-Language Model Can't Even Count to 20," arXiv:2510.04401. https://arxiv.org/html/2510.04401
[^ocrbench]: Liu et al., "OCRBench / OCRBench v2: On the Hidden Mystery of OCR in Large Multimodal Models," arXiv:2501.00321. https://arxiv.org/html/2501.00321v2
[^gpt4doc]: "Notes on Applicability of GPT-4 to Document Understanding" (pixel-only GPT-4V can't match scores without external OCR), arXiv:2405.18433. https://arxiv.org/html/2405.18433
[^pope]: Li, Du, Zhou, Wang, Zhao, Wen, "Evaluating Object Hallucination in Large Vision-Language Models (POPE)," EMNLP 2023; arXiv:2305.10355. https://arxiv.org/abs/2305.10355
[^amber]: Wang et al., "AMBER: An LLM-free Multi-dimensional Benchmark for MLLM Hallucination Evaluation," arXiv:2311.07397. https://arxiv.org/abs/2311.07397
[^spurlens]: "SpurLens: Automatic Detection of Spurious Cues in Multimodal LLMs" (spurious cues amplify hallucination ~26×), arXiv:2503.08884. https://arxiv.org/html/2503.08884
[^formatfree]: Tam et al., "Let Me Speak Freely? A Study on the Impact of Format Restrictions on LLM Performance," EMNLP 2024 Industry; arXiv:2408.02442. https://aclanthology.org/2024.emnlp-industry.91.pdf
[^formattax]: "The Constraint Tax: Validity-Correctness Tradeoffs in Structured Outputs," arXiv:2605.26128. https://arxiv.org/html/2605.26128 · "Capacity, Not Format," arXiv:2606.09410 (delayed structure recovers accuracy). https://arxiv.org/html/2606.09410
[^selfcons]: "Self-Consistency Is Losing Its Edge: Diminishing Returns and Rising Costs," arXiv:2511.00751. https://arxiv.org/html/2511.00751v2
[^ninejudges]: "Nine Judges, Two Effective Votes: Correlated Errors Undermine LLM Evaluation Panels," arXiv:2605.29800. https://arxiv.org/html/2605.29800v1
[^deliberation]: "DeliberationBench: When Do More Voices Hurt?" (best-single beats deliberation ~15× on cost-quality), arXiv:2601.08835. https://arxiv.org/html/2601.08835v1
