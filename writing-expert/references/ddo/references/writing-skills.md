# Writing-Skill Catalog & Routing Matrix

Authoritative map of the writing-craft family for `/ddo` Step 2.5. The family is
**5 hub skills**, each owning a set of document-type reference files. ddo routes a
document to its owning hub plus the specific reference, then loads that reference's
**current** checks. It does **not** copy those checks here — style rules (ban
lists, tone tables, hub conventions) drift, and a copy would go stale.

> **Taxonomy note (2026-06).** Former standalone skills — `runbook-craft`,
> `postmortem-writing`, `knowledge-base-authoring`, `tutorial-writing`,
> `rfc-and-design-docs`, `support-ticket-writing`, `kill-the-AI-ism`,
> `document-critique`, and the rest — are now **reference files under these hubs**.
> Nothing was deleted; the path changed. Reference a folded topic as
> `<hub-id> (references/<name>.md)`, never as a bare standalone skill id.

## The 5 hubs

| Hub (skill id) | Owns | Route here when the document is… |
|---|---|---|
| `writing-expert` | General prose craft, voice, editing; also hosts the `document-critique` engine and the `kill-the-AI-ism` voice layer | general / unknown prose, email, real-time posts — and the voice layer for **every** document |
| `technical-writing-craft` | Software / product / engineering docs | runbook, post-mortem, RFC, spec, PRD, KB, how-to, tutorial, reference, explanation, API docs, changelog, commit/PR, error/microcopy, meeting minutes, user stories, incident comms, code/agent plans |
| `executive-comms` | Leadership / business / persuasion | exec summary, one-pager, board memo, OKRs, proposal/grant, pitch deck, whitepaper, case study, speech, founder letter, negotiation prep |
| `content-and-marketing-writing` | External / marketing / PR | sales copy, press release, crisis PR, newsletter, op-ed, launch blog, audio script, chatbot, NPS/CSAT response, support/escalation reply |
| `career-and-formal-writing` | Career / academic / legal / policy | resume, cover letter, JD, performance review, academic/citation, legal-adjacent, policy/governance, survey questions |

**Always-on voice layer.** For *every* document type, also load
`writing-expert/references/kill-the-AI-ism.md` — it carries the live Tier-1/2/3
AI-ism ban list and structural robot-tells that Pass 13 enforces. Load it; do not
transcribe it.

**The critique engine.** `document-critique` is
`writing-expert/references/document-critique.md` (Passes 0–14). ddo runs it as its
analysis engine (SKILL.md Step 3).

## How to load (SKILL.md Step 2.5 mechanics)

1. Prefer **activating the hub** via the Skill tool — it loads the hub's core
   conventions and tone-calibration table.
2. Then **`Read` the document-type reference** `<hub>/references/<name>.md` for the
   authoritative, current checks.
3. If a hub is not in `available-skills`, `Read` the reference file directly — the
   file carries the checks even when the hub itself is not activatable.
4. Cap at ~4 sources: owning-hub reference + the `kill-the-AI-ism` voice layer +
   (for customer- or exec-facing docs) `writing-expert` tone + (for a hybrid doc)
   one second-type reference.

The **Stable checks** column below is only a structural skeleton — a routing hint.
The loaded reference is authoritative for the volatile checks.

## Full routing matrix

### `technical-writing-craft` — software / product / engineering docs

| Document type | Reference | Stable structural checks (skeleton) |
|---|---|---|
| Runbook / playbook | `runbook-craft.md` | purpose · prerequisites · numbered imperative steps · per-step verification · rollback · troubleshooting |
| Incident post-mortem | `postmortem-writing.md` | summary · timeline w/ timestamps · blameless root cause · contributing factors · remediation w/ owner+date |
| Incident / status-page comms | `incident-comms.md` | impact · current status · ETA / next update time |
| RFC / design doc | `rfc-and-design-docs.md` | context · problem · proposal · alternatives w/ rejection rationale · risks · decision |
| Spec | `spec-writing.md` | scope · requirements (RFC 2119) · interfaces · out-of-scope |
| PRD | `prd-writing.md` | problem · goals / non-goals · requirements · success metrics |
| KB article | `knowledge-base-authoring.md` | symptom · cause · resolution · scope / visibility |
| How-to / task guide | `howto-writing.md` | single stated goal · ordered steps · expected result |
| Tutorial / training | `tutorial-writing.md` | learning outcome · prerequisites · runnable steps · checkpoints |
| Reference doc | `reference-doc-writing.md` | complete params · uniform entry structure · no narrative |
| Explanation / concept | `explanation-doc-writing.md` | one concept · why-before-how · no step list |
| API docs | `api-docs-craft.md` | endpoint · params · auth · example request/response · error table |
| Changelog / release notes | `changelog-and-release-notes.md` (`changelogs-for-humans.md`) | grouped by change type · user-facing impact · version + date |
| Commit message | `commit-message-craft.md` | imperative subject ≤ 50 chars · why in body |
| PR description | `pr-description-craft.md` | what · why · how to test · risk |
| Error message | `error-message-craft.md` | what happened · why · what to do next (no blame) |
| UI microcopy | `microcopy-and-ui-writing.md` | action-first · consistent terms · no dead ends |
| Meeting minutes / decision log | `meeting-minutes-and-decision-log.md` | decisions · action items w/ owner+date · open questions |
| User story / acceptance criteria | `user-story-and-acceptance-criteria.md` | role-goal-benefit · testable criteria |
| Agent plan / code plan | `agent-plan-writing.md` (`code-plan-writing.md`) | ordered steps · verification · rollback |

### `executive-comms` — leadership / business / persuasion

| Document type | Reference | Stable structural checks (skeleton) |
|---|---|---|
| Exec summary / one-pager | `one-pager-writing.md` | BLUF · the ask · owner · decision needed |
| Weekly update / status / board memo | hub core (`one-pager-writing.md`) | delta since last · risks (max 3) · next-period plan |
| OKRs | `okr-writing.md` | measurable KR · ambitious-but-scoped · outcome not output |
| Proposal / business case / grant | `proposal-and-grant-writing.md` | quantified problem · solution · alternatives · cost · risks · recommendation |
| Pitch / investor deck | `pitch-deck-writing.md` | problem · solution · market · traction · ask |
| Whitepaper | `whitepaper-writing.md` | thesis · evidence · implications · CTA |
| Case study | `case-study-writing.md` | challenge · solution · quantified result |
| Speech / talk / all-hands | `speech-writing.md` (`public-speaking-and-presentations.md`) | one core message · spoken cadence · clear close |
| Founder / shareholder letter | `founder-letter-writing.md` | candor · narrative arc · metrics |
| Negotiation / persuasion prep | `negotiation-and-persuasion.md` | interests · BATNA · framing |

### `content-and-marketing-writing` — external / marketing / PR

| Document type | Reference | Stable structural checks (skeleton) |
|---|---|---|
| Sales / marketing / landing copy | `sales-and-marketing-copy.md` | single value prop · benefit-first · one CTA |
| Press release | `press-release-writing.md` | dateline · inverted pyramid · quote · boilerplate |
| Crisis PR statement | `crisis-pr-writing.md` | acknowledge · action taken · no speculation |
| Newsletter | `newsletter-writing.md` | one primary story · scannable · single CTA |
| Op-ed | `op-ed-writing.md` | thesis · argument · call to action |
| Release blog / launch narrative | `release-blog-and-launch-narrative.md` | what's new · why it matters · how to get it |
| Audio script (podcast / VO) | `audio-script-writing.md` | spoken rhythm · no unpronounceables · cue marks |
| Chatbot / conversational | `chatbot-conversation-writing.md` | turn-taking · persona consistency · fallback path |
| NPS / CSAT / review response | `nps-response-writing.md` | acknowledge · specific · action |
| Support / escalation reply | `support-ticket-writing.md` | restate issue · next step · owner · timeline |

### `career-and-formal-writing` — career / academic / legal / policy

| Document type | Reference | Stable structural checks (skeleton) |
|---|---|---|
| Resume / CV | `resume-and-cv-writing.md` | quantified achievements · action verbs · ATS-clean |
| Cover letter | `cover-letter-writing.md` | fit · evidence · ask |
| Job description | `job-description-writing.md` | responsibilities · requirements · compliant pay range |
| Performance / self / peer review | `performance-review-writing.md` | specific behaviors · impact · growth, not labels |
| Academic / scholarly | `academic-and-citation-writing.md` | thesis · evidence · consistent citation style |
| Legal-adjacent (NDA / ToS / privacy) | `legal-adjacent-writing.md` | defined terms · unambiguous obligations · counsel-review flag |
| Policy / governance | `policy-and-governance-writing.md` | scope · roles · enforcement · RFC 2119 terms |
| Survey questions | `survey-question-writing.md` | unbiased · single-barreled · balanced scale |

### `writing-expert` — general prose + the voice layer

| Document type | Reference | Stable structural checks (skeleton) |
|---|---|---|
| General prose / unknown | `editing-and-revision.md` | topic-first paragraphs · active voice · parallel lists · cohesion |
| Email | `email-craft.md` | BLUF · single ask · subject-line discipline |
| Slack / real-time / status post | `realtime-writing-under-pressure.md` | lead with the answer · brevity · timestamp |
| **Any (voice layer — always load)** | `kill-the-AI-ism.md` | Tier-1 ban list · Tier-2/3 structural robot-tells |

## Loading rules (summary)

1. Reference a folded topic by `<hub-id> (references/<name>.md)`; the bare topic
   name is **not** an installed skill id and will not resolve on its own.
2. The loaded reference is authoritative for volatile checks (ban lists, tone,
   conventions). The skeleton columns above are routing hints only.
3. Always co-load the `kill-the-AI-ism` reference for the Pass-13 voice layer.
4. For a document type not listed, use the SKILL.md Step 2.5 dynamic fallback
   (score the hubs, then `Read` the closest reference in the winning hub).
