<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** A net-new reference in this family.
> Sibling topics in this family are reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: enterprise-poc-technical-validation
description: >-
  Run an enterprise technical proof-of-concept / proof-of-value that CONVERTS to production — the
  pre/early-sales technical-validation MOTION a MongoDB SE/TAM owns: POC vs POV vs pilot vs trial vs
  competitive bake-off (when each fits); success & exit criteria tied to a quantified outcome; the mutual
  action plan (MAP) as the POC spine; scoping (use-case, time-box, environment, test data, access); roles
  (SE/TAM, champion, economic buyer, evaluators); failure modes (no criteria, science projects, scope
  creep, happy ears, proving the wrong thing, technical-win-but-no-commercial-win); POC-to-production
  conversion; POCs under FSI/regulated constraints. TRIGGER: "POC vs POV vs pilot", "POC success/exit
  criteria", "build the MAP for an evaluation", "technical win but no deal / pilot purgatory", "run a POC
  at a bank". SKIP (owns MAP only as the pre-sales POC spine; champion/EB only as POC roles): post-sale
  MAP / value-realization / value hypothesis / TTV → value-realization-outcome-cs; MEDDPICC scoring,
  NRR/GRR → tam-commercial-metrics; EBR/QBR & health-score → tam-expertise; buyer-side TPRM/SIG-CAIQ/DORA
  → enterprise-vendor-management-and-tprm; POC-plan prose → writing-expert / executive-comms.
origin: local
version: "1.1.0"
updated: "2026-06-18"
category: custom
tags: [tam, sales-engineering, presales, poc, pov, technical-validation, enterprise, fsi]
keywords:
  - proof of concept POC proof of value POV pilot trial competitive bake-off
  - POC success criteria exit criteria definition of done time-box
  - mutual action plan MAP mutual success plan close plan reverse timeline
  - POC scoping use-case selection sandbox test data synthetic data
  - sales engineer SE TAM champion economic buyer technical evaluator roles
  - POC failure modes science project scope creep happy ears pilot purgatory
  - technical win commercial win gap pilot-to-production conversion
  - FSI bank regulated POC security review TPRM SR 11-7 DORA NYDFS
  - innovation budget tourist money POC qualification
  - value engineering proof of value cost of inaction business case
whenToUse:
  - "decide whether a POC, POV, pilot, trial, or competitive bake-off is the right validation instrument"
  - "scope a POC: pick the use case, time-box it, choose the environment and test data, control scope creep"
  - "write success criteria and exit criteria that tie the technical test to a quantified business outcome"
  - "build the mutual action plan (MAP) that is the POC's spine and assign role ownership (SE/TAM/champion/EB)"
  - "diagnose why a technically successful POC did not convert to a production purchase"
  - "run a technical evaluation at a regulated / FSI buyer (security-review sequencing, data handling, bank timelines)"
  - "qualify a POC before committing SE time (avoid free-consulting / tourist-money science projects)"
whenNotToUse:
  - "the MAP as a post-sale CS artifact, value-realization theory, value hypothesis, or time-to-value (use value-realization-outcome-cs)"
  - "MEDDPICC qualification scoring, NRR/GRR formulas, or commercial-metric benchmarks (use tam-commercial-metrics)"
  - "EBR/QBR structure, health-score framework selection, or account-deliverable review (use tam-expertise / tam-account-reports)"
  - "buyer-side TPRM lifecycle, the SIG/CAIQ questionnaire mechanics, or DORA/OCC vendor governance (use enterprise-vendor-management-and-tprm)"
  - "the sentence-level prose craft of a POC plan, readout deck, or customer email (use writing-expert / executive-comms / technical-writing-craft)"
related_skills:
  - tam-operations
  - value-realization-outcome-cs
  - tam-commercial-metrics
  - enterprise-vendor-management-and-tprm
---

# Enterprise POC / POV Technical-Validation Motion

The **pre/early-sales technical-validation motion** a MongoDB SE/TAM runs to prove a solution works
*in the buyer's world* and turn that proof into a production purchase. This skill is the **execution
playbook**: how to choose, scope, run, and close a POC so a technical win becomes a commercial win,
including the harder realities at a regulated / financial-services buyer.

**This file holds** §1 playbook, §2 taxonomy, §3 criteria, §4 MAP, §5 roles, §6 scoping, §7
failure-mode catalog, and §10 checklists. **The deeper material is one Read away:** §8 (POC-to-production
conversion + the conversion-statistics provenance table) and §9 (running POCs under FSI/regulated
constraints) live in `references/conversion-and-fsi.md` — load it for those topics. The full citation
list also lives in that sub-reference.

**Scope boundary (read first).** This skill owns the **POC-EXECUTION** motion. It deliberately does
**not** re-derive:

- **The MAP as a post-sale customer-success artifact, value-realization theory, the value hypothesis,
  value scorecards, or time-to-value** → `value-realization-outcome-cs`. Here, "value" appears only as
  *success criteria tied to a dollar outcome that drives the buying decision* (the pre-sales use of
  value), and the MAP appears only as the *pre-sales POC spine* (not the post-sale CS plan).
- **MEDDPICC qualification scoring, NRR/GRR formulas, commercial-metric benchmarks** →
  `tam-commercial-metrics`. This skill uses MEDDPICC's *Decision Criteria / Champion / Economic Buyer*
  vocabulary as evaluation **roles**, not as a scoring rubric.
- **Buyer-side TPRM lifecycle, the SIG/CAIQ questionnaire mechanics, GRC/VRM tooling, DORA/OCC vendor
  governance** → `enterprise-vendor-management-and-tprm`. §9 here is only the *seller's POC-execution
  overlay* on those constraints (sequencing, data, timelines), not the buyer-side governance discipline.
- **EBR/QBR structure, health-score framework selection** → `tam-expertise`.

> **Evidence labelling.** This domain is overwhelmingly **practitioner lore** (sales-engineering books,
> presales communities, vendor blogs), not split-test/RCT-validated science. Claims are tagged:
> **[consensus]** = 3+ independent practitioner sources agree; **[common]** = 2 sources or one strong
> authority; **[contested]** = sources genuinely disagree (contradiction preserved); **[stat: verify]** =
> a specific number whose provenance is weak or vendor-sourced (directional only). Where something *is*
> closer to data (Gartner/RAND/CEB), it is called out. Volatile market claims are dated **as of 2026**.
> Citations are footnote keys (e.g. `[^care]`); the full citation list is in `references/conversion-and-fsi.md`.

## Contents

1. [The fast playbook (TL;DR)](#1-the-fast-playbook-tldr)
2. [Modality taxonomy: POC vs POV vs pilot vs trial vs bake-off](#2-modality-taxonomy-poc-vs-pov-vs-pilot-vs-trial-vs-bake-off)
3. [Success criteria & exit criteria (tied to a quantified outcome)](#3-success-criteria--exit-criteria-tied-to-a-quantified-outcome)
4. [The mutual action plan (MAP) as the POC spine](#4-the-mutual-action-plan-map-as-the-poc-spine)
5. [Roles & ownership: who owns what](#5-roles--ownership-who-owns-what)
6. [Scoping discipline](#6-scoping-discipline)
7. [Failure modes → root cause → fix](#7-failure-modes--root-cause--fix)
8. POC-to-production conversion & the technical-win → commercial-win gap → **`references/conversion-and-fsi.md` §8**
9. Running POCs under FSI / regulated constraints → **`references/conversion-and-fsi.md` §9**
10. [Checklists](#10-checklists)

---

## 1. The fast playbook (TL;DR)

A defensible enterprise technical validation, start to finish:

1. **Qualify before you commit an SE.** Run a POC only when there is a *real technical question a demo
   cannot answer*, a funded initiative, a named economic buyer, and a timeline. If the buyer cannot name
   the question, "they need a better demo, not a POC." [consensus][^homerun][^30mpc][^itbrief]
2. **Pick the right instrument** (§2). Most evaluations should be lighter than a full POC; reserve the
   heavy POC for genuine technical uncertainty on a deal large enough to justify the SE cost. [consensus]
3. **Write 3–5 falsifiable success criteria + explicit exit/kill criteria, in the buyer's metrics,
   signed before any build** (§3). This is the single most-cited predictor of conversion. [common][^guideflow][^rework]
4. **Pre-agree the commercial next step:** "If we hit these by [date], what happens?" → a *named* step
   (contract review, procurement kickoff), not "we'll evaluate." This makes the POC a **conditional
   close**. [consensus][^30mpc][^nedl]
5. **Build the MAP as the spine** (§4): reverse-timeline from the buyer's go-live date; the POC is one
   swim lane alongside security, legal, business case, procurement.
6. **Scope ruthlessly** (§6): one land use case, hard time-box (2–6 weeks enterprise), production-like
   environment, representative (default-synthetic) data, a written out-of-scope list and parking lot.
7. **Run the business-case track in parallel:** cost of inaction + future value + payback, in the
   buyer's numbers, refined with the economic buyer, so the champion can carry it (the value MODEL is
   `value-realization-outcome-cs`; here it is success criteria expressed in dollars).
8. **Close at the readout, with the economic buyer in the room** (§8): run the criteria pass/fail and
   transition to the business decision *the same day*.
9. **At a regulated/FSI buyer** (§9): start the security review **in parallel on day 1**, scope to an
   isolated sandbox on synthetic data, pre-stage your attestation pack, and plan for **months**, not weeks.

The recurring truth across every source: **POCs rarely fail on the technology. They fail on process:
qualification, criteria, ownership, and the technical-win → commercial-win gap.** [consensus]

---

## 2. Modality taxonomy: POC vs POV vs pilot vs trial vs bake-off

**The crux:** these terms name distinct points on two axes — *what doubt is being resolved* (technical
feasibility vs business value vs operational readiness vs self-directed discovery) and *how much
commitment exists / where in the buying cycle*. [consensus]

**Hold two facts at once** [contested, but the *tension itself* is consensus]: the POC≠POV distinction
is **real and actionable when you are designing an evaluation**, AND the labels are **used
interchangeably in the field** (Vivun and Mailchimp explicitly say so, while Steerlab/Guideflow treat
them as sharply distinct).[^vivun][^steerlab][^guideflow] **Resolution: distinguish the *concept*
(technical-doubt vs value-doubt) from the *vocabulary* (unreliable in the wild).** When a customer says
"POC," confirm which doubt they actually need resolved.

| Modality | Question it answers | What it proves | Who runs it | Commitment | Typical duration | Buying-cycle stage | Wrong-tool failure mode |
|---|---|---|---|---|---|---|---|
| **Trial / free trial** | "Can *I* use it?" | Self-directed usability; the "aha" moment (only when time-to-value is short) | Buyer, self-serve (low-touch) | Lowest | 7–30 days | Early / consideration; SMB & simple products | Complex product → blank-screen, wrong-fit funnel; enterprise buyers disengage |
| **Proof of Technology (POT)** | "Can the tech support this?" | One isolated technical question (integration / throughput); no business users | Vendor + buyer technical staff | Low | Days–weeks | Earliest technical gate; often *inside* a POC | Mistaken for a full POC; over- or under-scoped |
| **POC** | "Can this work *here*, with our data/integrations/constraints?" | **Technical feasibility** in the buyer's environment | Vendor SE-led + buyer IT/eng/security | Medium | 2–6 weeks | Mid-late, *after* a demo earns it, *before* purchase | Used when the doubt is *value* → "unpaid consulting"; no criteria → POC purgatory |
| **POV (proof of value)** | "Is it *worth it*?" | **Measurable business value / ROI** vs the buyer's objectives | Vendor SE + AE + champion + **economic buyer** | Medium-high | 2–12 weeks (longer at enterprise) | Later; build the commercial case | Vague criteria = "worse than no POV"; skipped → spend unjustified |
| **Pilot** | "Can we *scale* this?" | **Operational readiness & adoption** in production | Buyer end-users + vendor CS | High (often post-decision) | 30–90 days | At/just past the decision boundary [contested] | Run during selection → vendor lock-in before terms negotiated; success discounted at scale |
| **Competitive bake-off** | "Which one *wins*?" | Comparative fit, head-to-head in the buyer's simulated-production env | Buyer-led, cross-functional team + all finalists' SEs | High (finalists) | Multi-week finalist exercise | Latest; *after* an RFP narrows to ~2 | Spec-sheet criteria → incumbent wins by default; challenger must re-frame the rubric |

**POC vs POV: the actionable test** [common][^guideflow]: *where does the buyer's doubt sit?* If it is
"will this technically work in our environment," that is a **POC** (audience: technical evaluators who
can *block*). If it is "is it worth the money / will it deliver ROI," that is a **POV** (audience:
economic buyer / VP who *funds* it). Output of a POC = a technical validation report; output of a POV =
**a business case in the customer's own numbers**. In most enterprise deals you need **both** audiences
convinced; the common sequence is **POC first to clear the technical gate, then (or in parallel) POV to
build the business case.** [common]

**Competitive bake-off: when it helps vs hurts** [consensus, 6+ sources][^pulse][^itsjustrevenue][^infoweek]:

- A bake-off has *one* job for the challenger: **not to out-feature the rival but to change *what gets
  evaluated*.** Work with the champion to shape the rubric around buyer-relevant questions the incumbent
  answers poorly, *before the incumbent knows they are being scored*.
- The **incumbent wins a vanilla feature-grid bake-off by default** ("the signed contract with
  auto-renew is what beats you"). Leading with a feature comparison validates the incumbent's playfield.
- Lock **scoring weights and the tie-breaker rule before scores come in**; reweighting after demos is
  "rationalizing a gut call, not deciding." Classic seller mistakes: badmouthing the competitor,
  over-talking your product, reflexive discounting.

**Adjacents:** *"Evaluation"* is the generic parent term (John Care lists "pilot, proof of value, trial,
or evaluation" as alternative names for the same post-selection high-touch step).[^care] *"Sandbox"* is
the **environment**, not a modality; do not treat it as a peer of POC/POV. *Guided/interactive demo*
increasingly runs *before* (or replaces) a trial/POC. [consensus]

> **As of 2026:** multiple sources document erosion of the pure-technical POC's influence: buyers
> "assume baseline technical capability," so feasibility is less often the core risk; *operational fit,
> procurement survival, and adoption* are. Forrester ("Proof Is the Product," Apr 2026) frames trials
> and POCs as a real GTM motion that "fails when treated as technical experiments or polished demos
> rather than decision-ready, outcome-driven engagements." Gravity is pulling every modality toward the
> **POV (outcome) end**, sharpened by generative-AI buys where buyers refuse to commit without
> validation on *their own data*. [common: directional trend, not settled][^forrester]

---

## 3. Success criteria & exit criteria (tied to a quantified outcome)

**What makes a good success criterion** [consensus][^rework][^nedl][^guideflow]: specific, **measurable,
falsifiable (binary pass/fail)**, controllable, and **agreed in writing before work starts**. A good
criterion has a number, a threshold, and an observation window:

- *Bad:* "We want to see if it's easy to use." / "The AI should improve our claims process."
- *Good:* "80% of pilot users complete their first core workflow without support within 5 business
  days." / "The model identifies ≥15% more payment leakage than the manual process at a <5%
  false-positive rate, over a 10-day window."

**How many:** **3–5 measurable criteria** is the consensus band (some run 2–3 for velocity deals). The
discipline matters more than the exact count; keep it tight enough to reach a clean go/no-go. [common]

**Define *failure*, not just success.** A strong enablement standard requires the rep to "document what
counts as success, what counts as failure, and what conversion decision follows each outcome." [common][^pointer]

**Who signs off:** the buyer's evaluation owner — **someone with the authority to say no** — must accept
the criteria as valid, in writing, *before the SE starts building*. [common][^rework][^care] Several
sources call written-criteria-up-front "the single biggest predictor of POC-to-close conversion"
(repeated across Guideflow, b2bsalestraining, rework; strong practitioner consensus, **not** an RCT).

**Tension to hold** [contested]: the *sales-side* instinct is to set "2–3 controllable exit criteria you
know you can win." The *buyer-integrity* view warns that criteria gamed to be trivially winnable produce
"success theater" — a fuzzy outcome that surfaces real problems only during implementation, where they
are "always more expensive." **Both are true: criteria must be winnable *and* genuinely diagnostic.**

### Exit / kill criteria

"Exit criteria" is used two ways; keep them distinct:

1. **Success/exit bar (definition of done):** the POC is *complete* when these are demonstrated.
2. **Stop / kill conditions:** when to *abandon* regardless of technical progress.

**Time-box exit (non-negotiable).** A POC needs a hard end date set up front; "an open-ended timeline is
not a POC; it's a free trial you're managing." Typical: **2–6 weeks** enterprise (scales with
complexity; ~1–2 weeks velocity). [consensus] [stat: verify] One source cites "POCs past 45 days see 60%
higher abandonment" — single-sourced, directional.

**When to kill / walk away** [common][^presalespulse], the cleanest enumerated list:

- Requirements are **fundamentally outside product capability** and not on the near-term roadmap.
- **Scope keeps expanding** with no end ("free consulting").
- **Key stakeholders disengage** and the champion cannot re-engage them (the deal is dead even if the
  POC succeeds).
- **Evaluation criteria shift repeatedly** (a moving target no product satisfies).
- **Timeline keeps slipping** without justification (delay = low priority = won't close).
- **Pre-kill signal:** if the customer is *hesitant to define success metrics at all*, do not start; it
  is a red flag for "tourist money" (see §8 in `references/conversion-and-fsi.md`).

### Tying criteria to a quantified business outcome (value-engineering, pre-sales use)

This is the spine that converts a technical proof into a commercial decision. The MEDDPICC **Decision
Criteria** split is the canonical model [consensus][^meddicc][^forcemgmt]:

- **Technical Decision Criteria:** what must work to make the solution viable at scale; owned by
  architects/engineering/security. *These are what a POC literally tests.*
- **Business / Economic Decision Criteria:** why the company should invest now, in revenue/risk/strategic
  terms; owned by the economic buyer. "If your decision criteria live only in technical language, you do
  not yet have business criteria."

The translation discipline (pre-sales "value engineering" / proof of value) — quantify four things or the
business case stalls in procurement [consensus][^spotlight][^growthhub][^forcemgmt]:

1. **Cost of the current state** (the most-skipped and most-persuasive number; it makes the status quo
   expensive; a "Quantified Risk-of-Inaction").
2. **Expected future-state value.**
3. **Payback period.**
4. **The risk / assumptions** around both.

Use the **buyer's own numbers from discovery**; separate **hard savings vs soft savings vs revenue
impact** (economic buyers weight them differently); pair **above-the-line** metrics (revenue,
time-to-market; get executive attention) with **below-the-line** (cost/FTE savings; get mid-management
buy-in). Build a **preliminary value hypothesis early** to create urgency, refine the full assessment
with economic-buyer input, and arm the champion to circulate it as their own.

> **Cross-reference, not duplication:** the *theory* of value realization, the value hypothesis as a CS
> construct, and value-scorecard cadence live in `value-realization-outcome-cs`. This section is only the
> *pre-sales* application: success criteria expressed as a dollar outcome that drives the purchase.

---

## 4. The mutual action plan (MAP) as the POC spine

**What it is** [consensus][^prolifiq][^rafiki][^itsjustrevenue]: a **jointly owned, living document**
mapping every step from the current evaluation state to a signed contract and successful go-live, with a
**named owner and target date on every task, for both buyer and seller**. (Also: mutual success plan,
close plan, joint execution plan.) The load-bearing word is *mutual*: co-authored, not a seller checklist
emailed for acknowledgment. **The POC is one swim lane inside the MAP,** alongside security review,
business-case approval, legal redlines, and procurement.

**Why it works** [common→consensus]:

- It makes the path to close **explicit, so stalls become visible immediately**: a missed date is a
  signal; refusing to build a plan at all is a bigger one.
- It is **the best late-stage qualification signal**: "the plan is the qualification. A buyer who
  engages with it, fights with it, and updates it is doing the internal work to get the deal done."
- [stat: verify] A widely-circulated "MAPs close deals ~38% faster" claim is **vendor-sourced with no
  independent corroboration**; treat as marketing, not data.

### Sample MAP skeleton (cite-able)

| Section | Contents |
|---|---|
| **Shared goal / exec summary** | The buyer's primary business objective + metric with baseline→target; the **desired go-live date and the business reason behind it** (board commitment, fiscal year, contract expiry); quantified cost of inaction |
| **Stakeholders & roles** | Named owners on **both** sides — economic buyer/exec sponsor, champion, technical evaluators, procurement, security, legal, SE, AE |
| **Workstreams / swim lanes** | Each with dated milestones: **technical validation/POC**, security & compliance review, business-case/ROI approval, legal redlines, procurement, implementation/onboarding |
| **Decision gates / go-no-go** | Each milestone names the artifact + sign-off needed; one gates the next ("security review complete *with sign-off* before legal redlines begin"); sequential, not aspirational |
| **The POC sub-plan** | 3–5 measurable success criteria, fixed start/end dates, named owners both sides, scheduled readout + explicit **"what happens next if we succeed"** |
| **Resources** | Links to the exact artifacts the champion needs (ROI model, references, security docs); the MAP *reduces* buyer effort |
| **Cadence** | A recurring review rhythm (a weekly 15-min sync in the final stretch) |

**How to build it** [consensus]:

- **Co-build live** on a shared screen; start with *their* milestones, then layer in yours. Do not send a
  pre-filled template for feedback.
- **Earn the right / timing:** introduce it *after* discovery or post-demo/technical-validation, when the
  buyer has signalled intent; not in discovery, not in month three (a late MAP is "a retrospective
  exercise, not a planning tool").
- **Reverse-timeline / backward planning is the core construction technique** [consensus][^rafiki]: anchor
  on the buyer's desired live date, then back into contract signature → legal review → security review →
  business-case approval → technical validation. "Most sellers work forward from today; the best work
  backward from the outcome." It doubles as an urgency-creating close technique. Buyer-centric framing:
  *"I want to make sure we're not the reason this slips. Can we spend 30 minutes mapping everything that
  has to happen, on your side and ours, to hit your go-live date?"*
- **Litmus test:** "If your champion cannot open the document, edit a date, and forward it to their CFO
  without translation, you don't have a mutual action plan; you have a seller's wishlist."

**Criticism (preserve)** [contested]: most vendor sources treat the MAP as near-universally good; skeptics
(rework, itsjustrevenue) treat it as *conditionally* useful and sometimes counterproductive. Failure
modes: seller-centric "theater" (every owner is a seller name), introduced too late, no real
accountability, spreadsheet-attached-to-email invisibility, and — importantly — **MAPs work best only
*after* real multi-threading** (you cannot name all the milestones/owners until the stakeholders are in
the conversation). In a procurement-driven enterprise deal that runs its own structured process, **mirror
their process rather than introducing a parallel MAP.**

> **Cross-reference:** the MAP as a **post-sale customer-success / value-realization artifact** (joint
> impact plan, value scorecard cadence) is `value-realization-outcome-cs`. Here it is the **pre/early-sales
> POC execution spine**.

---

## 5. Roles & ownership: who owns what

**The cleanest, best-corroborated rule** [consensus, 6+ sources][^guideflow][^prospeo][^nobel]: **the AE
owns the deal, the forecast, and the commercial relationship; the SE owns the technical narrative, the
demo, and the POC.** The SE's deliverable is the **technical win**; the AE keeps the POC tethered to the
buying decision. Discovery is *shared* (AE leads commercial qualification, SE deepens technical). A clean
**handoff both ways** matters: AE briefs the SE before technical meetings; SE feeds insight back after.

**Role glossary across the lifecycle** [consensus]:

- **SE / Solutions Engineer (a.k.a. Solutions Consultant):** *pre-sale*; owns discovery, demos, POC, RFP
  technical sections, security questionnaires; hands off at signature.
- **SA / Solutions Architect:** either a presales title (hyperscalers) *or* a delivery-side role designing
  production architecture; increasingly a hybrid spanning late-presales → implementation to reduce handoff
  loss; may join late-stage deals for architecture/security depth.
- **TAM / Technical Account Manager:** *post-sale* — ongoing technical relationship, adoption, escalations,
  QBRs, renewals. **In a new-logo POC the TAM is rarely the owner**; the TAM is typically pulled in by the
  SE for *expansion* evaluations on an existing account, and may run those solo. (For MongoDB Premium
  Services role taxonomy — TAM/NTSE/IR/DCE — see `tam-reference`.)

**Customer-side roles** (MEDDPICC vocabulary, used here as roles, not a scoring rubric) [consensus][^meddicc][^forcemgmt]:

- **Champion:** has **power, influence, and credibility** AND **actively sells on your behalf when you are
  not in the room** AND has a **vested interest in your success**. "No champion, no deal." You must
  *prepare* the champion for internal pushback and *test* them (will they introduce you to the EB?).
- **Champion vs Coach:** a **Coach** gives helpful insight/intros but **lacks the power/influence of a
  champion, or is not yet selling for you**. Discriminator: **evidence, not emotion** — a true champion
  helps build the business case and shape decision criteria. (*Power* = authority to make something
  happen; *influence* = ability to move people; distinct, and you need both.)
- **Economic Buyer:** the person with overall authority — "the power to say yes when others say no, and no
  when others say yes." Controls discretionary funds, has **veto**, signs/approves the contract, cares
  about cost, time-to-value, and team confidence. Usually will not self-identify; find them, often via the
  champion.
- **Technical Buyer / technical evaluator:** evaluates feasibility, fit, security, integration; **can veto
  on technical grounds but cannot approve** — "they block, they don't approve." Their question is "why
  you?" (vs the EB's "why now?"). [common]

**Mobilizers (research-grounded cousin of the Champion)** [consensus: CEB/Gartner data][^challenger]: CEB
(now Gartner) profiled thousands of B2B stakeholders into seven types in three groups. **Mobilizers**
(Go-Getter, Teacher, Skeptic) drive organizational change and can build consensus; **Talkers** (Friend,
Guide, Climber) engage readily but do not drive action; the **Blocker** favors status quo. High performers
engage Mobilizers; average performers engage Talkers — **engagement metrics produce false positives**
(Talkers talk but do not act). The seller's job is **Commercial Coaching**: arm the Mobilizer with insight
to wrangle the buying group. (Deeper buyer-committee psychology: `applied-psychology`.)

**Multi-threading beyond a single champion** [consensus: Gartner data + practitioner]:

- The enterprise B2B buying group is **6–11 stakeholders across ~5 functions** (Gartner); newer
  Forrester/6sense figures run higher. **95% of buying groups revisit earlier decisions.** The journey is
  **non-linear** (six "buying jobs," incl. *validation* — where the POC lives — and *consensus creation*).
- **Single-threading is the most common reversible failure** [stat: verify]: ~1-in-3 champions change
  roles, leave, or go silent during an average cycle (Chorus data). Go **with** the champion, not around
  them; asking for an intro to their VP "is the entire point of having a champion"; drafting the intro
  email for them to forward lifts intro conversion materially. Then **synthesize**: circulate one shared
  view — and **a short MAP is the synthesis artifact** that protects a multi-threaded deal from
  single-champion collapse.

### POC RACI / DACI

**RACI**: **R**esponsible (does the work), **A**ccountable (owns the outcome), **C**onsulted,
**I**nformed. The classic rule is **exactly one Accountable per row**. In a vendor↔buyer POC,
accountability legitimately **splits across the table** (a vendor-side owner for the seller's tasks, a
buyer-side owner for the buyer's), so the convention below is **one Accountable per side per row** —
when a row has both an `A (vendor)` and an `A (buyer)` it is naming the two distinct sides, not breaking
the rule. For decision moments, practitioners often prefer **DACI** (Driver, Approver [one],
Contributor, Informed; bounded timeline). A PreSales Collective article recommends a **DACI matrix to
delineate POC responsibilities between buyer teams** (BU, procurement, IT, PMO). [common]

Sample POC RACI (one Accountable per side; tailor to the deal):

| POC activity | SE | AE | Champion | Economic Buyer | Tech Evaluators | Proc/Sec/Legal |
|---|---|---|---|---|---|---|
| Qualify POC is warranted | C | **A/R** | C | I | C | — |
| Define & sign off success criteria | **R** | C | **A** (buyer-side) | I | C | C (security) |
| Build/configure POC environment | **A/R** | I | I | — | C | — |
| Run the evaluation / test cases | **R** (vendor side) | I | C | I | **A/R** (buyer side) | C |
| Secure exec access & "what if we win" | C | **A/R** (vendor) | **R** (internal) | **A** (buyer) | I | — |
| Security / compliance review | C | I | C | I | C | **A/R** |
| Business case / ROI | C | **R** | **R** (internal sell) | **A** | I | — |
| Readout + go/no-go | **R** (present) | **A** (vendor: commercial) | **R** (champions) | **A** (buyer: final yes/no) | C | C |
| Contract / procurement | I | **A/R** (vendor) | C | **A** (buyer) | I | **R** |

**Criticism (preserve)** [contested: McKinsey + HBR + practitioners]: McKinsey finds "RACI often makes
things worse" and advises "give more people a voice, but fewer people a vote, and don't use RACI"; HBR
(2026) finds decision-rights frameworks "often fail" in practice. The steel-man rebuttal: the *concept* of
clear single-owner accountability is useful; the failure is mechanical/cultural over-application.
**Implication for a POC:** the lightweight DACI (one Approver) or simply the **MAP's
named-owner-per-milestone structure** is often a better fit than a heavy RACI grid; the MAP *is* the de
facto responsibility matrix for the evaluation.

---

## 6. Scoping discipline

**Governing idea** [consensus]: a POC is the most expensive instrument in the technical-sales motion, so
**bound it ruthlessly to prove the fewest things that will actually trigger the purchase.**

**Use-case selection** [consensus][^presalespulse][^care][^cohan]:

- Validate **3–5 critical capabilities**, not 15. "The biggest POC mistake is trying to evaluate
  everything."
- Anchor on the **compelling event**: the timeframe-forcing business pressure (compliance deadline, new
  product launch, end-of-year budget). [common] Care reports practitioner data that a *fully qualified*
  customer taken to POC converts at "about five out of six," heavily caveated on qualification. [stat: verify]
- Pick the **land use case / beachhead**, not the full vision. Chris White's **Minimum Viable Demo**
  logic: for every step ask **"so what?"**; if you cannot answer, remove it. Cohan separates **Vision
  Generation** ("what's possible," before discovery) from **Technical Proof** (after discovery, only as
  long as needed); a POC is a Technical-Proof instrument, and a *narrower beachhead* dramatically reduces
  the burden of proof. "The more contained you make it, the easier it is to close." [common]
- Anti-pattern use cases for a *first* POC: high-stakes decisions, extensive integration, immature data,
  unclear success criteria. [common]

**Time-boxing** [consensus]: a hard end date set before kickoff. Bands: enterprise POC **2–6 weeks**
(beyond 6 weeks the scope is too broad); POV longer; "if you cannot prove the hypothesis in ~4 weeks, the
scope is wrong or the hypothesis is not testable." Open-ended POCs die via loss of urgency → champion
disengagement → "POC purgatory." A long POC also *costs*: it ties up an SE for weeks and burns goodwill;
treat each POC as a bet on the deal being worth the team's time. **Technique:** set the customer's
worst-case expectation generously while targeting internally shorter ("buffer and beat").

**Environment: the realism ↔ speed ↔ risk trilemma** [contested, preserve]:

- Spectrum: **vendor-hosted sandbox → customer-hosted sandbox/staging → production**. The SE should own
  setup; the most-repeated operational failure is **the environment not being ready at kickoff**, burning
  week one.
- *"Sandbox is misleading" camp:* "no sandbox operates at the scale and complexity of production"; a
  vendor who *insists* on a sandbox may be hiding poor production performance; cleansed-sandbox POCs
  produce surprises at implementation.
- *"Production is too risky" camp:* production / production-derived data in eval environments is a
  compliance and breach hazard (non-production environments are "up to 80% of an enterprise's attack
  surface").
- *Defensible middle:* a **production-like / representative** environment — schema, volume, and
  characteristics matching production, *without raw production data* — and a **staged approach**
  (abbreviated sandbox to shortlist, then production-like/integrated run for the finalist). Make it
  **believable** (no "demo"/"test" labels, no fake names, realistic exceptions). [consensus]
- *Note:* frictionless **self-service trials** are a growing alternative to the provisioned POC for
  hands-on technical buyers, but they shift the realism question onto the buyer.

**Test data** [consensus]:

- **Representative beats clean:** preserve cardinality, distribution, null patterns, join complexity,
  seasonality, and edge cases. "Toy rows make demos look clean but useless." Match production on
  **structure, volume, and characteristics**.
- The **data-prep burden is where POCs are quietly doomed**: make a **data audit (volume, quality, access
  constraints) a formal first-phase deliverable**, not an assumption.
- **Default to synthetic**; reserve **masked production** for shapes only production has; **real customer
  data only masked/anonymized/pseudonymized**, never raw exports (deep regulated rules → §9 in
  `references/conversion-and-fsi.md`).

**Getting access (general)** [consensus]: **provisioning lead time is the silent scope-killer.** Start
provisioning **the day after scoping, not the day the POC starts** (enterprise setup can take 3–5 business
days). Send an access checklist *before* kickoff (accounts, SSO, API keys, VPN, permissions). "No access"
— to users or systems — is itself a disqualifier ("you might be stuck in a lab, not a pipeline").

**Scope-creep control** [consensus]: scope creep is "the #1 killer of POC timelines"; each addition seems
reasonable in isolation. Controls: a **signed charter/scoping doc before day 1**; an **explicit
out-of-scope list** (3–4 items) on the POC deck so the customer self-edits; an **"In-POC or Out-of-Scope"
decision tree** (on charter? blocks a success metric? → else **park it**); a **visible parking lot /
Phase-2 backlog** so you say "not now" instead of "no"; **change control** (post-signature additions
require a written amendment that *extends the timeline*; no free expansions); reframe new asks as
**trade-offs** ("that's a Phase 2"), never silent yeses or flat refusals. [contested, preserve]: skeptics
argue creep is *structural* to enterprises that "never intend to buy" (extract free consulting), so demand
fixed scope + a buying-intent test; the contract camp argues creep "is almost always a contract problem;
you did not define done." Synthesis: creep is **predictable but controllable** — discipline lives in
*both* the contract and the qualification.

> **Cross-reference:** the modality choice that frames a scope (POC vs POV vs pilot) is §2; the criteria
> that define "done" are §3.

---

## 7. Failure modes → root cause → fix

POCs rarely fail on technology; they fail on process. The named-pattern catalog [consensus across the
presales corpus]:

| Failure mode (named) | Root cause | Fix (pointer) |
|---|---|---|
| **No / vague success criteria** | Vendor skips the step, or writes criteria too fuzzy to pass-fail → "a curiosity problem, not a buying problem" | §3: signed, falsifiable, measured criteria before any build; define failure too |
| **"Science project" / POC purgatory** | No decision date; no economic buyer; champion lacks internal buy-in; "just exploring" | §3 time-box + §8 readout-as-gate; "if the customer won't commit to a closing meeting, don't start the POC" |
| **Scope creep ("death by a thousand asks")** | No documented scope boundary; each add feels small | §6: charter + out-of-scope list + parking lot + timeline-linked change control |
| **"Happy ears"** | Optimism replaces qualification; enthusiasm mistaken for buying intent; inflates forecasts | Disqualifying questions ("what could stop you?"); verify economic buyer + decision process; reward early disqualification; "be evidence-driven, not optimistic" |
| **Proving the WRONG thing** | POC built around an *interesting* technical question never tied to the metric the exec is accountable for | §3: connect finding → business outcome → the stakeholder under pressure; "the technical win is the foundation of the commercial case, not the case" |
| **Technical-win-but-no-commercial-win** | A successful POC tests whether the *product works*, never whether the *org can buy it*; procurement/legal/finance never looped in | §8: run the business-case track in parallel; EB at the readout; tie results to dollars |
| **Single-threading / champion-only** | "The person who loves your product usually cannot buy it"; single point of failure | §5: multi-thread the buying group; validate/test the champion; "if your champion left tomorrow, would the deal survive?" |
| **Free-consulting / tire-kicker** | No funded initiative, no criteria, no timeline; SE labor extracted | §1/§9: qualify before committing SE time; convert to a *paid* POV; "running >2 free POCs/quarter for sub-$50K deals = broken qualification" [stat: verify] |
| **Over-engineering / boiling the ocean** | "Analysis paralysis dressed as thoroughness"; building production-grade in the sandbox (and giving away the rollout, so procurement rejects the services quote) | §6: narrow to top-3 risks / one workflow-one number-one quarter; build the smallest test that proves the principle |

> §8 (the technical-win → commercial-win gap, conversion bridges, the tourist-money problem, and the
> conversion-statistics provenance table) is in `references/conversion-and-fsi.md`.

---

## 10. Checklists

### Before you agree to a POC (qualification)
- [ ] A **real technical question a demo cannot answer** (if the buyer cannot name it → better demo, not a POC)
- [ ] Funded initiative, named **economic buyer**, a timeline / compelling event
- [ ] The deal is large enough to justify SE cost (some set a rough ACV floor; **[stat: verify]**)
- [ ] Decided the right instrument (POC vs POV vs pilot vs trial vs bake-off; §2)

### Charter (signed day 0)
- [ ] **3–5** falsifiable, measured success criteria in the buyer's metrics, and a **failure** definition
- [ ] **Explicit out-of-scope list** (3–4 items) on the POC deck
- [ ] **Fixed end date** (2–6 weeks enterprise) + a pre-committed go/no-go rubric
- [ ] **Pre-agreed commercial next step** ("if we hit these by [date], what happens?")
- [ ] Named evaluation owner with authority to say no; recurring sync booked
- [ ] The POC is one swim lane in a **co-built MAP** reverse-timelined from go-live (§4)

### Environment & data
- [ ] Environment chosen on the realism↔speed↔risk spectrum *consciously*, not by default
- [ ] **Data audit** as a formal first-phase deliverable (volume, quality, access)
- [ ] Test data **representative** (distribution, nulls, joins, edge cases, realistic scale); **default-synthetic**; real data only masked
- [ ] SE owns setup; environment + data ready *before* kickoff

### Access (start day 1 of scoping)
- [ ] Access checklist issued pre-kickoff (accounts, SSO, API keys, VPN, permissions)
- [ ] Provisioning started; lead time budgeted; security-review lead time costed into the timeline
- [ ] SE has pre-tested every use case before the customer touches it

### Creep control (continuous)
- [ ] "In-POC or Out-of-Scope" decision tree in use; visible **parking lot / Phase-2 backlog**
- [ ] Change control: post-signature additions require a written amendment that *extends the timeline*
- [ ] New asks reframed as trade-offs ("that's a Phase 2"), never silent yeses or flat refusals

### Close (readout)
- [ ] Readout booked at kickoff; **economic buyer in the room**
- [ ] Run criteria pass/fail live; transition to the business decision the **same day**
- [ ] Business case (cost of inaction + value + payback, in the buyer's numbers) ready for the champion to carry

### FSI overlay (detail in `references/conversion-and-fsi.md` §9)
- [ ] Security review started **in parallel on day 1**; attestation pack pre-staged
- [ ] Isolated sandbox + synthetic/masked data; no production NPI/PCI/PII
- [ ] AI/ML: explainability + SR 11-7-style validation designed in
- [ ] Timeline planned in **months**; workstreams parallelized

---

**Citations & deep sections:** §8 (conversion + statistics provenance) and §9 (FSI/regulated) plus the
full footnote citation list are in `references/conversion-and-fsi.md`. Cross-references:
`value-realization-outcome-cs` (value model / post-sale MAP), `tam-commercial-metrics` (MEDDPICC scoring,
NRR/GRR), `enterprise-vendor-management-and-tprm` (buyer-side TPRM mechanics), `applied-psychology`
(buyer-committee psychology), `fsi-banking-regulatory-context` (FSI regulatory map), `tam-reference`
(MongoDB Premium Services roles).
