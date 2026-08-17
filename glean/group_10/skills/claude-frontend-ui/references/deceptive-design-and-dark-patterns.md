<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** The design-ETHICS critique axis: *does this
> interface deceive, manipulate, or coerce the user — and what is the honest, compliant alternative?*
> The lens is **DETECT (the UI tell) → RATE (user harm + legal/compliance exposure) → FIX (the honest
> pattern)**. Sibling topics in this family are reference files under the hubs (`frontend-ui`) — **not**
> standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a
> bare sibling skill; load that topic's `references/<name>.md` from the owning hub (see the hub's
> "Cross-hub map").
>
> **Scope boundary — do not duplicate siblings.** This file owns the *deception/manipulation/coercion*
> critique axis and its 2024-26 legal grounding. It is the *ethics* complement to: **persuasion done
> honestly + the heuristics/Laws of UX that dark patterns weaponise**, owned by
> `usability-heuristics-laws-of-ux` (cross-reference, never re-document — esp. Hick, Jakob, the
> aesthetic-usability effect, Von Restorff); **how a design erodes vs earns first-impression TRUST**,
> owned by `emotional-design-and-visual-trust` (cross-reference); **visual-composition principles**
> (hierarchy, gestalt, C.R.A.P.), owned by `visual-design-principles-and-critique`; **how to run/
> facilitate a crit**, owned by `visual-design-critique-methodology`; **objective machine-computable
> metrics**, owned by `computational-aesthetics-ui-metrics`; **WCAG/accessibility**, owned by
> `accessibility-ux-reviewer`. The general decision/persuasion *psychology* (nudge vs sludge as theory,
> reactance, Cialdini, prospect theory) lives in the `applied-psychology` hub. Drafting honest customer
> copy → `content-and-marketing-writing`. **This file is a critique lens, not legal advice** — it rates
> *exposure* and names the governing instrument; it does not adjudicate a specific case.

---
name: deceptive-design-and-dark-patterns
title: Deceptive Design & Dark Patterns (the regulatory generation)
description: >-
  The design-ETHICS critique axis of graphic/UI work: DETECT a deceptive, manipulative, or coercive
  pattern (Brignull, Gray, Mathur dark-pattern taxonomies) → RATE its user harm and 2024-26 legal
  exposure (DSA Art. 25, GDPR/ePrivacy, FTC + Negative Option/click-to-cancel, CCPA/CPRA) → FIX it
  with the honest pattern, across high-risk flows (sign-up, consent/cookie banners, subscription &
  cancellation, checkout, defaults). TRIGGER: audit/critique a UI for dark patterns, deceptive design,
  confirmshaming, roach motel, drip pricing, fake urgency/scarcity, asymmetric cookie banner,
  dark-pattern GDPR/DSA/FTC/CCPA exposure. SKIP: DESIGNING honest urgency/scarcity/persuasion to lift
  conversions (not auditing for deception), or the heuristics/Laws of UX →
  usability-heuristics-laws-of-ux; first-impression trust → emotional-design-and-visual-trust; WCAG →
  accessibility-ux-reviewer; persuasion theory (nudge, Cialdini) → applied-psychology; drafting copy →
  content-and-marketing-writing.
category: developer
version: "1.0.1"
updated: "2026-06-16"
metadata:
  verified-as-of: "2026-06-16"
  changelog:
    - "2026-06-16 sko v1.0.0->v1.0.1 — Pass H 10/10 pos, 9/10->10/10 neg (predicted; build-vs-audit SKIP carve-out added); fixed 1 High (CNIL fine §4.2 inconsistency: corrected to Google €150M/Facebook €60M Dec 2021, matching §5.3 + source), 1 High (description 1439->985 chars, under Glean 1000 cap), 1 Medium (§7 step 3 '5 attribute tests'->6, restored Disparate Treatment) + 2 Low ('leverage'->'most useful'). Pass I clean (all sibling concepts cross-ref, not duplicated); Pass J on-demand reference, length justified; Pass L clean."
    - "2026-06-16 v1.0.0: created via /dr deep-research (FILL-IN round, design-ETHICS critique lens). 4 parallel research agents; ~80 deduped independent sources across taxonomies, EU law, US law, and per-flow application. Volatile legal claims (FTC click-to-cancel vacatur, enforcement amounts) stamped verified-as-of 2026-06-16."
keywords:
  - dark patterns
  - deceptive design
  - deceptive patterns
  - design ethics
  - confirmshaming
  - roach motel
  - drip pricing
  - fake urgency
  - fake scarcity
  - cookie banner consent
  - DSA Article 25
  - GDPR consent dark patterns
  - FTC click to cancel
  - negative option rule
  - CCPA CPRA dark pattern
  - sludge
  - choice architecture
  - manipulation vs persuasion
tags:
  - developer
  - frontend
  - ux
  - design
  - ethics
  - critique
  - privacy
  - compliance
whenToUse:
  - "auditing or critiquing a UI/mockup/flow for deceptive, manipulative, or coercive patterns"
  - "naming a UI tell against the Brignull / Gray / Mathur dark-pattern taxonomies"
  - "rating a design's legal/compliance exposure under DSA Art. 25, GDPR/ePrivacy, FTC Act/ROSCA, or CCPA/CPRA"
  - "fixing a cookie banner, sign-up, subscription/cancellation, checkout, or default to the honest pattern"
  - "deciding whether an urgency/scarcity/social-proof cue is honest persuasion or a deceptive pattern"
whenNotToUse:
  - "honest persuasion technique + the usability heuristics / Laws of UX themselves — use usability-heuristics-laws-of-ux"
  - "first-impression trust, credibility, brand-tone critique — use emotional-design-and-visual-trust"
  - "WCAG / accessibility conformance review — use accessibility-ux-reviewer"
  - "decision/persuasion psychology as theory (nudge, reactance, Cialdini, prospect theory) — use applied-psychology"
  - "drafting the honest customer-facing copy itself — use content-and-marketing-writing"
related_skills:
  - usability-heuristics-laws-of-ux
  - emotional-design-and-visual-trust
  - accessibility-ux-reviewer
  - applied-psychology
---

# Deceptive Design & Dark Patterns — the regulatory generation

The **design-ethics critique lens**. Where the other `frontend-ui` critique references ask "is this
*good*?" (composition, usability, trust, aesthetics), this one asks: **"is this *honest*?"** — does the
interface deceive, manipulate, or coerce the user into a choice they would not otherwise make, against
their own interest and in favour of the business? And in 2024-26, that question is no longer only
ethical: the DSA, GDPR, the FTC, and a dozen US state privacy laws have made many dark patterns
**illegal**, with fines reaching billions.

Use this file to run a **DETECT → RATE → FIX** pass over a UI, mockup, flow, or spec:

- **DETECT** — spot the concrete UI tell and name it against a taxonomy (§2, §4).
- **RATE** — score user-harm severity *and* compliance/legal exposure (§3, §5, §6).
- **FIX** — prescribe the documented honest/compliant alternative (§4, §7).

## Contents

1. [The honesty boundary: persuasion vs manipulation vs deception](#1-the-honesty-boundary)
2. [The named taxonomies — your detection vocabulary](#2-the-named-taxonomies)
3. [What makes a pattern "dark": the harm/severity model](#3-what-makes-a-pattern-dark)
4. [The high-risk flows: DETECT → RATE → FIX tables](#4-the-high-risk-flows)
5. [The EU regulatory layer (DSA, GDPR, ePrivacy/CNIL, EDPB, UCPD)](#5-the-eu-regulatory-layer)
6. [The US regulatory layer (FTC, ROSCA, CCPA/CPRA, state laws)](#6-the-us-regulatory-layer)
7. [Running the critique: the audit method + severity rubric](#7-running-the-critique)
8. [Anti-patterns of the critique itself (don't over-flag)](#8-anti-patterns-of-the-critique)
9. [Cross-references](#9-cross-references)
10. [References](#10-references)

---

## 1. The honesty boundary

The whole lens turns on one distinction. **Persuasion** engages the user's deliberate reasoning with
accurate information; **manipulation** hijacks heuristics or hides information; **deception** induces a
false belief.[^chou][^conversation] A critique must separate these because **only deception and
coercion are reliably unlawful** — "annoying" or "persuasive" is not, on its own, a violation.[^conversation][^ftcreport]

Two framings ground the lens:

- **Nudge vs sludge** (Thaler & Sunstein). A *nudge* helps people choose "as judged by themselves"; a
  *sludge* is friction that leaves people worse off "by their own lights."[^thaler][^cma-oca] Dark
  patterns are sludge + deception. The governing principle a critique should assert: **defaults and
  friction must serve the user, not the business.**[^cma-oca][^oecd]
- **The technical-indistinguishability problem.** At a code/crawler level, an honest nudge and a
  deceptive pattern often use the *same* technique (a default, a countdown, a scarcity message). The
  differentiator is **whether the influence benefits or harms the user, and whether the underlying fact
  is true.**[^ceur][^ftcreport] This is why the single most useful operational test is **"is it
  true?"** — *if the urgency/scarcity is real, show it; if it is fabricated, it is deceptive.*[^nng-scarcity][^ftcreport]

> **The single most useful test (the "full-disclosure test").** For any conversion-critical
> element, rewrite it as a perfectly accurate plain-language description of what it actually does. **If
> that honest description would change the user's click, it is a dark-pattern candidate.** Persuasion
> survives disclosure; manipulation collapses under it.[^chou][^nng-deceptive]

Empirically, dark patterns *work* and they *harm*: in the first large randomized experiment, mild dark
patterns more than doubled an unwanted-enrollment rate (~11% → ~26%) and aggressive ones roughly
quadrupled it (~42%); aggressive patterns triggered consumer **backlash** while mild ones did not (so
mild-but-effective patterns are the more insidious policy problem), and **less-educated users were
significantly more susceptible.**[^luguri][^nng-deceptive] Most users do not consciously recognise dark
patterns, but recognition improves once they are informed.[^digeronimo]

---

## 2. The named taxonomies

These are your **detection vocabulary**: spot a tell, name it, and you can both reason about harm and
point to the regulator/research that treats it as a problem. Three taxonomies dominate; a 2024 ontology
unifies them.

### 2.1 Brignull / deceptive.design — the canonical surface types

Harry Brignull coined "dark patterns" in 2010 and runs **deceptive.design** (formerly
darkpatterns.org). He has since shifted terminology to **"deceptive patterns"** (his 2023 book title),
and is explicit that his categories were "a rallying cry," **not a rigorous classification.**[^brignull][^dd-book-p1][^dd-book-p3]
The current canonical type list (cite the live `/types` page for current names):[^dd-types]

| Type | Tell (verbatim definition, condensed) |
| --- | --- |
| **Comparison Prevention** | features/prices combined so the user can't compare |
| **Confirmshaming** | the decline option is worded to guilt/shame the user |
| **Disguised Ads** | an ad made to look like content or an interface element |
| **Fake Scarcity** | a fake "limited supply/popularity" signal |
| **Fake Social Proof** | fake reviews/testimonials/activity messages |
| **Fake Urgency** | a fake time limit |
| **Forced Action** | the user must do an undesirable thing to get what they want |
| **Hard to Cancel** (Roach Motel) | easy to sign up, very hard to cancel |
| **Hidden Costs** (drip pricing) | unexpected fees appear only at checkout |
| **Hidden Subscription** | a recurring charge with no clear disclosure/consent |
| **Nagging** | persistent interruptions to do something not in the user's interest |
| **Obstruction** | barriers that make a task or info hard to reach |
| **Preselection** | a default option pre-selected to steer the decision |
| **Sneaking** | a transaction entered on false pretences (info hidden/delayed) |
| **Trick Wording** | confusing/misleading language steering an action |
| **Visual Interference** | info hidden, obscured, or disguised against expectation |

> Note for honest citation: the modern deceptive.design list **absorbed Gray et al.'s 2018 strategies
> in the 2023 refresh** (nagging, obstruction, sneaking, forced action, visual interference) without
> citation or hierarchy. For "the origin," cite the 2010 definition; for "current types," cite the live
> `/types` page.[^gray-ontology][^dd-types]

### 2.2 Gray, Kou, Battles, Hoggatt & Toombs (2018) — the five strategies (by designer motivation)

CHI 2018. Organised by the **designer's strategy/motivation** (not the surface pattern). The five
high-level strategies, with Brignull's surface tricks nested under them:[^gray2018]

1. **Nagging** — a minor, repeated redirection of expected functionality (pop-ups, prompts, audio).
2. **Obstruction** — making a flow harder than it needs to be to dissuade an action. *Sub: Roach Motel;
   Price Comparison Prevention.*
3. **Sneaking** — hiding/disguising/delaying information relevant to the user. *Sub: Forced Continuity;
   Hidden Costs; Sneak into Basket; Bait and Switch.*
4. **Interface Interference** — manipulating the UI to privilege some actions over others (the broadest
   bucket). *Sub: Hidden Information; Preselection; Aesthetic Manipulation → which contains Toying with
   Emotion (confirmshaming), False Hierarchy, Disguised Ad, Trick Questions.*
5. **Forced Action** — requiring a tangential action to access functionality. *Sub: Social Pyramid
   (forced recruitment); forced disclosure.*

These five "were carried forward" across later taxonomies and are the stable shared backbone — **the
FTC uses "coerced action" in place of "forced action."**[^gray-ontology]

### 2.3 Mathur et al. (2019) "Dark Patterns at Scale" — the e-commerce categories + empirics

CSCW 2019. An automated crawl of **~53K product pages across ~11K shopping sites** found **1,818
dark-pattern instances (15 types in 7 categories) on 1,254 sites (~11.1%)** — an explicit **lower
bound**, and **more common on more popular sites.**[^mathur2019][^mathur-princeton] *(An earlier 2019
draft reported 1,841/1,267; use the final 1,818/1,254 figures.)* The 7 categories, with the single
most prevalent tell noted:[^mathur-princeton][^mathur2019]

| Category | Tells (with notable prevalence) |
| --- | --- |
| **Sneaking** | Sneak into Basket; Hidden Costs; Hidden Subscription |
| **Urgency** | **Countdown Timer (most common urgency tell)**; Limited-time Message — Mathur found **>40% of timers were fake/reset** |
| **Misdirection** | Confirmshaming; Visual Interference; Trick Questions; Pressured Selling |
| **Social Proof** | Activity Messages ("N people viewing"); Testimonials of Uncertain Origin |
| **Scarcity** | **Low-stock Message (the single most prevalent type overall)**; High-demand Message |
| **Obstruction** | Hard to Cancel |
| **Forced Action** | Forced Enrollment |

Mathur's **Urgency / Scarcity / Social Proof** categories are the behavioral-economics contribution
that Gray 2018 lacked; the 2024 ontology unifies them under a sixth high-level pattern, **"social
engineering."**[^gray-ontology]

### 2.4 The unifying ontology (Gray et al. 2024) — three levels

CHI 2024 harmonised **ten** academic + regulatory taxonomies into a **three-level ontology** with
standardized definitions for **65 dark-pattern types**:[^gray-ontology]

- **High-level = strategy** (context-agnostic: the six are nagging, obstruction, sneaking, interface
  interference, forced/coerced action, social engineering).
- **Meso-level = "angle of attack."**
- **Low-level = "means of execution"** (the situated, Brignull-style surface trick).

**Practical use:** map any UI tell to a **low-level** Brignull name, roll it up to a **high-level** Gray
strategy, and you have a stable label that cross-walks to every regulator's vocabulary.

### Cross-taxonomy mapping (use to translate a tell across sources)

| Brignull / deceptive.design | Gray 2018 strategy | Mathur 2019 category |
| --- | --- | --- |
| Sneak into Basket / Hidden Costs / Hidden Subscription | Sneaking | Sneaking |
| Roach Motel / Hard to Cancel | Obstruction | Obstruction |
| Comparison Prevention | Obstruction | — |
| Preselection / bad defaults | Interface Interference | Misdirection |
| Confirmshaming | Interface Interference (Aesthetic → Toying w/ Emotion) | Misdirection |
| Visual Interference / False Hierarchy / Misdirection | Interface Interference (Aesthetic Manip.) | Misdirection |
| Trick Wording / Trick Questions | Interface Interference | Misdirection |
| Disguised Ads | Interface Interference | (Misdirection-adjacent) |
| Nagging | Nagging | (added later) |
| Fake Urgency (countdown / limited-time) | — | Urgency |
| Fake Scarcity (low-stock / high-demand) | — | Scarcity |
| Fake Social Proof (activity / testimonials) | — | Social Proof |
| Forced Action / Forced Enrollment / friend-spam | Forced Action | Forced Action |

**Reading:** Gray's **Interface Interference** ≈ Mathur's **Misdirection** (the broadest buckets).
**Nagging, Obstruction, Sneaking, Forced Action** are the shared backbone across all three (FTC: "coerced
action"). **Urgency/Scarcity/Social Proof** are Mathur's addition, later unified as "social engineering."

---

## 3. What makes a pattern "dark"

A tell is not automatically "dark" — *Forced Action, Obstruction,* and *Nagging* can be entirely
truthful. The field's main conceptual paper (Mathur, Kshirsagar & Mayer 2021) shows there is **no
single definition**, only a set of related considerations.[^mathur2021] Use its two frameworks to RATE.

### 3.1 The six "dark" attributes (does the design have any of these?)

Across two themes:[^mathur2021]

- **Modifying the decision space:** **Asymmetric** (unequal burden on the choices — e.g., a bright
  "Accept" vs a buried "Reject"); **Restrictive** (eliminates choices — forced account, roach motel);
  **Disparate Treatment** (disadvantages one group — e.g., pay-to-skip).
- **Manipulating the information flow:** **Covert** (the influence mechanism is hidden — a decoy);
  **Deceptive** (induces a false belief — a fake timer); **Information Hiding** (obscures/delays needed
  info — drip pricing).

**Rating heuristic:** a tell that is **Deceptive or Restrictive** maps most directly to legal exposure;
a tell that is only **Asymmetric or Covert** is the harder (but often still unlawful under autonomy/
"freely given" standards) case. The more attributes a tell trips, the higher the severity.

### 3.2 The four normative lenses (whose welfare is harmed?)

To judge *darkness*, ask which welfare the design sacrifices:[^mathur2021]

1. **Individual welfare** — financial loss, loss of privacy, or cognitive burden to the user.
2. **Collective welfare** — harm to market competition or to trust in the market (proliferation breeds
   skepticism that harms honest sellers too).
3. **Regulatory objectives** — undermining the aims of a regime (e.g., subverting GDPR/CCPA consent).
4. **Individual autonomy** — how far it impairs the user's ability to decide independently.

---

## 4. The high-risk flows

Five flows concentrate dark patterns. Each table is **DETECT → RATE → FIX**. The legal instruments in
the RATE column are exposure pointers — see §5–§6 for the rule text. (A 2022 EU mystery-shopping study
found **97% of popular EU sites/apps used ≥1 dark pattern**; the five most common were hidden info/false
hierarchy, preselection, nagging, difficult cancellation, and forced registration.[^eu-behavioural])

### 4.1 Sign-up / onboarding

| DETECT (UI tell) | RATE (harm + likely law) | FIX (honest pattern) |
| --- | --- | --- |
| **Confirmshaming** — decline worded to shame ("No thanks, I hate saving money") | Emotional manipulation; trust erosion. FTC Act §5; often subtle enough to escape law | Neutral decline of equal weight ("No thanks"); also give X + click-outside exits[^nng-shaming][^uxbooth] |
| **Forced account creation** before purchase/use | Forces unwanted data disclosure; task barrier. GDPR Art. 25 (EDPB calls guest mode the privacy-protective default) | Offer "create account **OR** continue as guest"; make guest path prominent[^edpb-accounts] |
| **Pre-checked marketing / data-sharing opt-ins** | Consent by inertia. GDPR (pre-ticked boxes invalid — *Planet49*); FTC §5 | **Opt-IN** — unchecked boxes; active choice required[^cookiebot] |
| **Friend-spam / address-book harvesting** — bright "Grow your network", tiny text reveals it emails all contacts | Coerced action; reputational harm. FTC §5 (coerced-action category) | Full disclosure before access; equal-weight Skip; let the user pick which contacts[^thinkdesign] |
| **Privacy-zuckering** — vague settings/broad defaults publish more than intended | Over-sharing under a false sense of protection. GDPR Art. 25; FTC §5; DSA Art. 25 | Privacy-protective defaults; unbundle distinct purposes[^cookiebot] |
| **Trick wording / double-negative toggles** | Cognitive burden on the protective choice. FTC §5; UCPD | Plain, single-polarity language; toggle state matches effect[^nng-deceptive] |

### 4.2 Consent / cookie banners (the most-measured flow)

| DETECT | RATE | FIX |
| --- | --- | --- |
| **Asymmetric Accept vs Reject** — bright "Accept all", "Reject" a grey link / absent on layer 1 | *The* canonical tell. Removing "reject all" from the first page raises consent ~22–23 pts.[^nouwens] GDPR consent invalid; CNIL fined **Google €150M and Facebook €60M (Dec 2021)** for exactly this[^cnil-formal][^techcrunch2022] | **"Accept all" and "Reject all" on the same first layer, equal size/colour**, one click each (CNIL's explicit rule)[^cnil-formal] |
| **False hierarchy** — green Accept vs greyed Reject | Visual misdirection of the System-1 click. EDPB Guidelines 03/2022; DSA Recital 67 ("more prominence to certain choices") | Equal visual weight; neutral styling[^cnil-formal] |
| **Pre-toggled "legitimate interest"** switches on by default | Non-essential processing runs unless the user hunts every toggle. GDPR (LI for ad-tracking widely rejected) | Non-essential categories **off by default**; block until consent[^cnil-formal] |
| **Cookie/consent wall** — "Accept to continue", no real reject | Consent not "freely given". GDPR Art. 7(4); DSA Art. 25 | A genuine reject that still grants access[^cnil-formal] |
| **Nagging re-prompts** after a reject; **"Manage" mazes** (reject buried 2+ clicks, ambiguous labels) | Friction wears users into accepting. CNIL & Belgian DPA cited all of these | One-layer symmetric reject; stop re-prompting after a valid choice[^cnil-formal][^belgian-noyb] |
| **Phantom reject** — "Reject all" exists but ad cookies fire anyway / keep running after withdrawal | The reject is cosmetic. CNIL acted on exactly this | Reject must actually suppress non-essential tags; honor withdrawal[^cnil-formal] |

### 4.3 Subscription & cancellation

> **CURRENCY CAVEAT (volatile, verified 2026-06-16):** The FTC **"Click-to-Cancel" Rule was vacated in
> its entirety by the Eighth Circuit in July 2025** on procedural grounds (`Custom Communications v.
> FTC`); the FTC re-opened rulemaking with a new ANPRM in early 2026. **Do NOT cite click-to-cancel as
> binding federal law.** **ROSCA and FTC Act §5 remain fully in force** and are the live basis for
> subscription enforcement.[^8thcir][^ftc-anprm]

| DETECT | RATE | FIX |
| --- | --- | --- |
| **Roach motel** — one-click signup, but cancel is phone-only with a retention agent / "type this phrase" | Users keep paying. ROSCA "simple mechanism" + FTC §5 — `FTC v. Vonage` ("panoply of hurdles… dark patterns", **$100M**); `FTC v. Amazon` Prime ("Iliad" cancel maze, **$2.5B** settlement) | Cancellation **as easy as signup, in the same medium** used to enroll[^vonage][^amazon] |
| **Retention maze / forced "are you sure" + save-offers** delaying cancellation | Sludge to exhaust the user. FTC §5/ROSCA (the *rule* provision restricting save-attempts was vacated, but the principle stands) | A direct cancel path; any retention offer optional and skippable[^ftc-negopt-policy] |
| **False hierarchy at cancel** — "Keep" large/bright vs "Cancel" small/greyed | Interface interference. FTC report names this exact example | Equal visual weight for Keep vs Cancel[^ftcreport] |
| **Hidden subscription / free-trial-to-paid** auto-conversion; recurring charge under a "one-time" guise | Users charged for things they forgot/never bought; older consumers disproportionately harmed. ROSCA + FTC §5 + CFPB Circular 2023-01 | Clear price + renewal date + "charges continue unless you cancel" **before** taking billing info; easy cancel[^cfpb] |
| **"Finished-looking" cancel that isn't** — a confirmation that didn't actually cancel | Deceives the user while billing continues. FTC §5 (deception) | Unambiguous confirmation that charges have **stopped**; immediate halt[^conversation] |

### 4.4 Checkout / e-commerce

| DETECT | RATE | FIX |
| --- | --- | --- |
| **Drip pricing / hidden costs** — mandatory fees revealed only at the final step | Raises price paid by **up to ~10–13%** and users *stay* with the worse option; #1 cause of cart abandonment ("extra costs too high"; ~70% abandon). FTC junk-fees; UCPD; UK CMA | **All-in pricing up front** — show total incl. fees/shipping early[^drip-bu][^baymard] |
| **Sneak-into-basket** — pre-added items/insurance/donations or a pre-checked add-on | Users pay for things they never chose. FTC §5 (sneaking); UCPD aggressive practice | **Opt-IN add-ons only**; nothing in the cart the user didn't add[^mathur2019] |
| **Fake urgency** — a countdown timer that **resets on reload** or is inconsequential | Fabricated deadline → impulse buys. Mathur: **>40% of timers fake**. FTC names reset timers as deceptive | If the deadline is real, show it; if not, don't invent it[^ftcreport][^nng-scarcity] |
| **Fake scarcity** — "Only 2 left!" regardless of real stock; same count across products | Manufactured FOMO; **backfires** when users sense insincerity. FTC §5 if false | Surface **real inventory** dynamically[^nng-scarcity][^scarcity-backfire] |
| **Fake social proof** — "24 people are viewing this", fabricated review counts | Misleads on popularity/credibility. FTC names "N people viewing" as deceptive | Show real, verifiable activity or omit it; disclose paid reviews[^ftcreport] |
| **Disguised ads** — paid content/ranking formatted to look editorial | User can't tell ad from content. FTC §5 + deceptively-formatted-ads ("net impression") | Clear "Ad/Sponsored" labels; disclose pay-for-placement ranking[^ftcreport] |

### 4.5 Defaults & interface interference (the choice-architecture layer)

| DETECT | RATE | FIX |
| --- | --- | --- |
| **Preselection / bad defaults** — the business-favouring option is pre-chosen (max data sharing on; paid shipping pre-selected) | Defaults are powerful (CMA: "substantial evidence of the power of default settings"). GDPR Art. 25 (privacy by default); UCPD; DSA | **User-serving defaults**; a pre-set default is the option a reasonable user would want[^cma-oca][^oecd] |
| **Visual interference / false hierarchy / misdirection** — the business path is visually dominant, the user path is dim/hidden | Covert steering of the System-1 choice. FTC interface-interference; DSA Recital 67 | Neutral visual weight across options[^ftcreport] |
| **Trick questions / double negatives** in forms/toggles | Cognitive burden on the protective choice (Mathur "asymmetric"). FTC §5; UCPD | Plain, single-polarity language; pretest comprehension[^ftcreport] |
| **Bait-and-switch "X"** — the modal's close X opens the offer instead of closing | Subverts a universal convention (violates Jakob's Law → see usability-heuristics-laws-of-ux). FTC §5 | X closes; honor learned conventions[^ftcreport] |
| **Intermediate-currency / decoy obfuscation** — points/credits or a decoy that hides true cost | Covertly distorts cost comparison via the decoy effect. UCPD; FTC §5 if misleading | Show real prices directly; no manufactured decoys[^mathur2021] |

---

## 5. The EU regulatory layer

**Structure to remember:** the DSA's dark-pattern ban is **residual (lex specialis)** — it applies only
where conduct is *not* already caught by the **UCPD** (commercial practices) or **GDPR** (personal-data
processing). So a cookie-consent dark pattern is a **GDPR/ePrivacy** matter; DSA Art. 25 bites on
non-data, non-commercial manipulations (sign-out friction, repeated nagging on a platform UI).[^dsa][^springer-frag]

### 5.1 EU Digital Services Act (Reg. 2022/2065), Article 25

The operative prohibition (Art. 25(1)): providers of online platforms shall **not design, organise or
operate their online interfaces in a way that deceives or manipulates** recipients, **or otherwise
materially distorts or impairs** their ability to make **free and informed decisions.**[^dsa][^dsa-cms]
Art. 25(2) carves out UCPD/GDPR practices. Art. 25(3) names example tells the Commission may issue
guidance on: **(a)** giving more prominence to certain choices; **(b)** repeatedly requesting a choice
already made (interfering pop-ups); **(c)** making termination harder than subscription.[^dsa-cms]
**Recital 67** defines dark patterns and adds: cancelling significantly harder than signing up, making
sign-out unreasonably difficult, and **default settings that are very difficult to change.**[^dsa]

- **Binds:** providers of **online platforms** (social networks, marketplaces, app stores); VLOPs/VLOSEs
  carry the heaviest supervision. **Applies from 17 Feb 2024.**[^dsa-enforce]
- **Penalties:** up to **6% of annual worldwide turnover** (Art. 52); the Commission directly supervises
  VLOPs/VLOSEs, national Digital Services Coordinators the rest.[^dsa-penalty][^dsa-enforce]
- **Honest caveat:** Art. 25 is widely criticised as **vague and largely residual**, has **no in-article
  definition** (read Recital 67), and as of 2026-06-16 has **no headline CJEU ruling or Commission fine
  specifically on Art. 25 dark patterns** — it is young and largely untested.[^taylorwessing][^dsa-observatory]

### 5.2 EDPB Guidelines 03/2022 — deceptive design patterns (the GDPR-facing taxonomy)

**Adopted in final form 14 Feb 2023**; the final version renamed "dark patterns" → **"deceptive design
patterns."** They are **soft law (not legally binding)** but treated as authoritative by national DPAs
and routinely cited.[^edpb-news][^legiscope] Six categories, each tied to GDPR breaches (mainly Art. 5
fairness/transparency, Art. 12, Arts. 4(11)/7 consent, Art. 25 by-design/default):[^edpb-guidelines]
**Overloading, Skipping, Stirring, Obstructing, Fickle, Left in the dark.**

### 5.3 GDPR consent + ePrivacy/CNIL cookie rules

**Consent (Art. 4(11), Art. 7, Recital 32)** must be **freely given, specific, informed, unambiguous**,
by a **clear affirmative act**; **silence, pre-ticked boxes, and inactivity do not constitute
consent**; withdrawal must be **as easy as giving**; consent must not be bundled or conditional on a
service it isn't needed for.[^gdpr-recital32][^edpb-consent] **CJEU `Planet49` (C-673/17, 2019)** held a
**pre-checked cookie checkbox is not valid consent** (and applies regardless of whether the data is
personal).[^planet49]

**ePrivacy + CNIL:** storing/reading non-essential cookies needs prior **opt-in** consent. CNIL's rule
(guidelines Oct 2020, enforced from 1 Apr 2021): **refusing must be as easy as accepting** — an
equally-prominent one-click **"Reject all" at the first layer.**[^cnil-formal] Enforcement (historical,
appeals possible — verify before relying): **CNIL fined Google €150M and Facebook €60M (Dec 2021)** and
**Google €100M + Amazon €35M (Dec 2020)** for cookie consent dark patterns, under **Art. 82 of the
French Data Protection Act** (implementing ePrivacy — note this lets the local DPA act without the GDPR
one-stop-shop).[^techcrunch2022][^hunton] The **Belgian DPA** ordered a **"reject" button on the first
layer** and found **deceptive button colours** unlawful (Mediahuis decisions).[^belgian-noyb]

### 5.4 UCPD + forward-looking

The **UCPD (2005/29/EC)** Commission Guidance (Dec 2021) confirms it **covers dark patterns** — they
can breach Art. 5 (professional diligence), Arts. 6–7 (misleading), or Arts. 8–9 (aggressive / undue
influence).[^ucpd] **DMA Art. 13(6)** bans gatekeepers from making rights "unduly difficult" via
non-neutral interface design.[^edpb-interplay] **Forward-looking (TENTATIVE):** a **Digital Fairness
Act** is on the Commission's 2026 work programme (indicative Q4 2026) and is expected to consolidate
dark-pattern, addictive-design, and personalisation rules — **not yet law.**[^dfa]

---

## 6. The US regulatory layer

### 6.1 FTC — "Bringing Dark Patterns to Light" (Sept 2022) + Section 5

The staff report flags **four categories**: (1) design that **induces false beliefs** (disguised ads,
fake comparison sites, false urgency); (2) design that **hides/delays material info & costs** (drip
pricing/junk fees); (3) design that **leads to unauthorized charges** (deceptive subscription
enrollment, free-trial conversion); (4) design that **obscures/subverts privacy choices** (steering to
max data sharing).[^ftcreport][^ftc-press] The hook is **FTC Act §5** — "unfair or deceptive acts or
practices" (UDAP): the **deception** test (mislead-likely + reasonable interpretation + material) and
the **unfairness** test (substantial injury, not reasonably avoidable, not outweighed by benefits, §
45(n)).[^usc45][^ftc-unfairness] **§5 alone has no civil penalty for a first violation** (the FTC gets
injunctions; `AMG Capital` stripped §13(b) money) — which is **why the FTC pairs dark-pattern cases with
ROSCA**, where penalties attach.[^ftc-negopt-policy]

### 6.2 The Negative Option / "Click-to-Cancel" Rule — VACATED (volatile)

The 2024 final "click-to-cancel" rule (cancellation at least as easy as sign-up, same medium; clear
disclosure; express informed consent) was **vacated in its entirety by the Eighth Circuit on 8 July
2025** (`Custom Communications, Inc. v. FTC`, 142 F.4th 1060) on procedural grounds — the FTC skipped
the **preliminary regulatory analysis** required once the rule's impact exceeded $100M/yr. Vacatur
**reinstated the 1973 prenotification-only rule.** The FTC **re-opened rulemaking** with an ANPRM in
early 2026.[^8thcir][^ftc-anprm] **Status verified 2026-06-16: no click-to-cancel rule is in force.**

### 6.3 ROSCA — the live federal baseline

**ROSCA (15 U.S.C. §§ 8401-8405)** is **in force** regardless of the vacated rule. **§ 8403** makes it
unlawful to charge via an online negative-option feature unless the seller: **(1)** clearly and
conspicuously discloses **all material terms before** obtaining billing info; **(2)** obtains **express
informed consent** before charging; and **(3)** provides a **simple cancellation mechanism.** Failing
any one element violates ROSCA; the FTC's position is that "simple" means **at least as easy as
enrollment, same medium** (the pre-rule source of "click to cancel", which survives the vacatur).
ROSCA violations carry **civil penalties + redress.**[^rosca][^ftc-negopt-policy]

### 6.4 California CCPA/CPRA — the statutory "dark pattern" definition

CPRA (operative 2023) **defines** a dark pattern (Cal. Civ. Code § 1798.140(l)): *"a user interface
designed or manipulated with the substantial effect of subverting or impairing user autonomy,
decisionmaking, or choice."* Critically, **"agreement obtained through use of dark patterns does not
constitute consent"** (§ 1798.140(h)).[^ccpa-civ] The regulations (**11 CCR § 7004**) add the
**symmetry-in-choice** rule — *the path to the more privacy-protective option must not be longer/harder/
more time-consuming than the less-protective path* — and the **"effect, not intent"** test: a UI **is** a
dark pattern if it **has the effect** of substantially subverting autonomy, regardless of intent.[^ccpa-regs][^cppa-advisory]
The **CPPA** and the **California AG** enforce. **`Sephora` (AG Bonta, Aug 2022, $1.2M)** — the first
public CCPA action — established that businesses must honour a **user-enabled universal opt-out signal
(Global Privacy Control)** the same as a manual "Do Not Sell" click.[^sephora][^iapp-sephora]

### 6.5 Other state laws + 2024-25 enforcement

**~20 state privacy laws** track the CPRA "substantial effect / autonomy" definition and **void
dark-pattern consent** — e.g., **Colorado CPA** (C.R.S. § 6-1-1303(9); regs add "commonly used is not,
alone, enough" to clear a pattern) and **Connecticut CTDPA** (§ 42-515(14), which expressly **incorporates
"any practice the FTC refers to as a dark pattern"**).[^colorado][^ctdpa] Headline FTC enforcement
(amounts volatile, verified 2026-06-16): **`FTC v. Amazon` (Prime "Iliad")** — **$2.5B** settlement
(Sept 2025; $1B penalty + $1.5B redress) for dark-pattern enrollment + a labyrinthine cancel flow;
**`FTC v. Adobe`** (filed Jun 2024, **pending**) — hidden early-termination fee + preselected plan +
onerous cancellation; **`FTC v. Vonage`** — **$100M** (2022) for cancellation dark patterns + junk
fees.[^amazon][^adobe][^vonage]

> **Honest disconfirming note for severity-rating.** The "dark pattern" standard is **contested**: the
> U.S. Chamber of Commerce (amicus in `Amazon`) and ITIF argue the FTC's definitions are **vague /
> "fail to provide fair notice"**; scholars note the CPRA definition is hard to apply and may exceed
> the CPPA's "user interface" authority; and post-`Loper Bright` the FTC's expansive §5 reading is more
> legally vulnerable.[^chamber][^itif][^goldman][^promarket] Rate exposure as **likelihood of harm in
> context**, not as a settled verdict.

---

## 7. Running the critique

A repeatable DETECT → RATE → FIX pass:

1. **Inventory conversion-critical surfaces** — every button, modal, default, toggle, banner, and the
   sign-up↔cancel and add↔remove journeys.
2. **DETECT — run the full-disclosure test** (§1) on each, then **name the tell** against the taxonomy
   (§2). ~15 min/surface. The honest rewrite that *would change the click* is your finding.
3. **DETECT — apply Mathur's 6 attribute tests** (§3.1): is it Asymmetric, Restrictive, Covert,
   Deceptive, or Information-Hiding — plus Disparate Treatment when it targets one group of users
   (e.g., pay-to-skip)? Tally — more attributes = higher severity.
4. **RATE — separate "annoying" from "deceptive/coercive."** Only the latter is reliably unlawful
   (§1). Then attach the exposure: which welfare is harmed (§3.2) and which instrument likely applies
   (§5–§6). Use the severity rubric below.
5. **RATE — for journeys, run a sludge/asymmetry audit:** count "substantial actions"/clicks for the
   user-serving path vs the business-serving path (e.g., signup clicks vs cancel clicks). A large
   asymmetry on a privacy or money choice is a high-severity finding.[^cma-oca]
6. **FIX — prescribe the documented honest pattern** (§4) and, where law applies, the compliant form
   (symmetric reject-all; opt-IN defaults; same-medium cancel; all-in pricing; honor GPC).

### Severity rubric (adapt the `frontend-ui` 0–4 critique scale to ethics)

| Severity | Criterion |
| --- | --- |
| **4 — Blocker / unlawful** | Deceptive or Restrictive tell on a money/privacy/consent choice with a clear governing rule (e.g., pre-ticked consent, phantom reject, drip pricing, roach motel). Likely a §5/GDPR/DSA/CCPA violation — fix before ship. |
| **3 — Major** | Strong manipulation (asymmetric + covert) impairing autonomy on an important choice (false hierarchy at cancel, sneak-into-basket). Strong exposure even if not a named per-se violation. |
| **2 — Moderate** | A manipulative tell that engages a heuristic but the underlying fact is true and a choice exists (mild nagging, a real-but-pushy urgency). Trust risk; the "mild but effective" category — flag it. |
| **1 — Minor** | Borderline; arguably persuasion. Note it and apply the full-disclosure test; recommend a neutral alternative. |
| **0 — Not a dark pattern** | Honest persuasion: true information, real choice, survives disclosure. Do **not** flag (see §8). |

**Detection tooling (assistive, not authoritative):** deceptive.design's **types catalog + "hall of
shame"** as a vocabulary/example bank; the **Mathur Princeton crawler** as the empirical baseline;
emerging **LLM-driven audits** are promising at interaction-level identification but audits "remain
largely manual" — treat AI detectors as a first pass a human verifies.[^dd-types][^mathur-princeton][^llm-audit]

---

## 8. Anti-patterns of the critique

- **Over-flagging honest persuasion.** A true scarcity message ("8 left" when 8 are left), a genuine
  sale countdown, or a clearly-labeled upsell is **not** a dark pattern. The differentiator is **truth +
  real choice**, not the technique.[^nng-scarcity][^ceur] Run the full-disclosure test before flagging.
- **Asserting a universal "effect size."** There is **no single effect size of dark patterns**; many
  users notice and resist, with large demographic variance. Rate *likelihood of harm in context*, not
  universal manipulation.[^chou][^mathur2021]
- **Citing vacated/volatile law as settled.** The FTC click-to-cancel rule is **vacated** (§6.2); DSA
  Art. 25 is **largely untested**; enforcement amounts and pending cases change. Stamp legal claims and
  re-verify (this file: verified-as-of 2026-06-16).[^8thcir][^taylorwessing]
- **Confusing "deceptive design pattern" scopes.** The EDPB term is **GDPR-specific** (privacy
  processing on social media); the FTC's is **consumer-protection**; the CPRA's is **statutory**. Use
  the right vocabulary for the right regime (§5.2, §6.1, §6.4).
- **Treating EDPB/soft-law guidance as binding.** EDPB guidelines and CNIL recommendations are
  authoritative but **not directly binding** — frame them as strong regulatory expectation/risk, not
  black-letter law.[^legiscope]
- **Adjudicating instead of advising.** This lens rates **exposure** and names the instrument; it is
  **not legal advice.** Recommend counsel for a live compliance question.

---

## 9. Cross-references

- **`usability-heuristics-laws-of-ux`** — the honest twin of this lens. Dark patterns weaponise the
  same heuristics it documents: **Jakob's Law** (the bait-and-switch X violates a learned convention),
  **Hick's Law** (overloading/"too many options"), **Von Restorff** (false hierarchy via salience), the
  **aesthetic-usability effect** (a polished skin that masks a manipulative flow). Read it for the
  legitimate use; read this file for the abuse.
- **`emotional-design-and-visual-trust`** — the trust consequence. Dark patterns are a primary way a
  design **erodes** the first-impression trust/credibility that reference teaches you to build; a
  manipulative pattern detected here predicts a trust defect there.
- **`visual-design-principles-and-critique`** / **`visual-design-critique-methodology`** — composition
  principles and how to run the crit; compose this ethics pass alongside them.
- **`accessibility-ux-reviewer`** — WCAG conformance (a separate axis; obstruction can also be an
  accessibility failure).
- **`applied-psychology`** (hub) — the decision/persuasion theory beneath the tells: nudge vs sludge,
  reactance, Cialdini's principles, prospect theory/loss aversion.
- **`content-and-marketing-writing`** (hub) — drafting the honest replacement copy (neutral declines,
  clear disclosures, truthful urgency).

---

## 10. References

[^chou]: Yu-kai Chou — "Brignull's Manipulative Design (Dark Patterns) in UX." https://yukaichou.com/gamification-analysis/dark-patterns-brignull-manipulative-design-ux/ — practitioner/expert — persuasion-vs-manipulation boundary; the full-disclosure test; no single effect size.
[^conversation]: G. Dickinson (legal scholar), "Dark patterns on the web — why aren't they all illegal?" The Conversation. https://theconversation.com/dark-patterns-on-the-web-279961 — explainer — only deception/coercion is reliably unlawful; "finished-looking" fake cancel is deception.
[^thaler]: R. Thaler, "Nudge, not sludge," *Science* 361:431 (2018). https://www.science.org/doi/10.1126/science.aau9241 — research — nudge vs sludge; "as judged by themselves."
[^cma-oca]: UK CMA, "Online Choice Architecture: evidence review of OCA and consumer/competition harm" (2022). https://www.gov.uk/government/publications/online-choice-architecture-how-digital-design-can-harm-competition-and-consumers — regulator — sludge; power of defaults; not-all-friction-is-sludge.
[^oecd]: OECD, "Dark Commercial Patterns" (2022). https://www.oecd.org/content/dam/oecd/en/publications/reports/2022/10/dark-commercial-patterns_9f6169cd/44f5e846-en.pdf — regulator/research — cross-taxonomy overview; sneak-into-basket + preselection as combined sludge.
[^ceur]: "Comparing Nudges and Deceptive Patterns," CEUR Vol-3720 paper 12. https://ceur-ws.org/Vol-3720/paper12.pdf — research — nudges and deceptive patterns are technically indistinguishable; truth/benefit is the differentiator.
[^nng-deceptive]: Nielsen Norman Group, "Deceptive Patterns in UX: How to Recognize and Avoid Them." https://www.nngroup.com/articles/deceptive-patterns/ — practitioner(NNG) — definition; mild patterns unnoticed; audit questions.
[^luguri]: J. Luguri & L. Strahilevitz, "Shining a Light on Dark Patterns," *J. Legal Analysis* 13(1):43 (2021). https://academic.oup.com/jla/article/13/1/43/6180579 — research — mild ~2× / aggressive ~4× unwanted-enrollment; backlash asymmetry; education effect.
[^digeronimo]: Di Geronimo et al., "UI Dark Patterns and Where to Find Them," CHI 2020. https://dl.acm.org/doi/fullHtml/10.1145/3313831.3376600 — research — 240 apps, 95% had ≥1 pattern; uses Gray taxonomy; users mostly don't recognise them.
[^brignull]: Harry Brignull personal site. https://brignull.com/ — practitioner — coined "dark patterns" (2010); expert-witness cases.
[^dd-types]: deceptive.design — Types index. https://www.deceptive.design/types — docs — canonical current 16 type names + definitions; the "hall of shame" example bank.
[^dd-book-p1]: deceptive.design Book, Part 1. https://www.deceptive.design/book/contents/part-1 — book — 2010 definition + the dark→deceptive terminology shift.
[^dd-book-p3]: deceptive.design Book, Part 3. https://www.deceptive.design/book/contents/part-3 — book — Brignull: his naming was a "rallying cry," not a rigorous classification.
[^gray2018]: Gray, Kou, Battles, Hoggatt & Toombs, "The Dark (Patterns) Side of UX Design," CHI 2018, DOI 10.1145/3173574.3174108. https://classes.cs.uchicago.edu/archive/2020/fall/33231-1/readings/2018_Grayetal_CHI_DarkPatternsUXDesign — paper — five strategies + verbatim definitions + nested sub-types.
[^gray-ontology]: Gray, Santos, Bielova & Mildner, "An Ontology of Dark Patterns," CHI 2024. https://colingray.me/wp-content/uploads/2024/02/2024_Grayetal_CHI_OntologyDarkPatterns.pdf — paper — 3-level ontology / 65 types; cross-taxonomy cross-walk; "social engineering" 6th high-level; FTC "coerced action."
[^mathur2019]: Mathur et al., "Dark Patterns at Scale: Findings from a Crawl of 11K Shopping Websites," CSCW 2019, arXiv 1907.07032. https://arxiv.org/pdf/1907.07032 — paper — 1,818 instances / 15 types / 7 categories / 1,254 sites (~11.1%); 5 attributes; >40% of timers fake; lower-bound + popularity correlation.
[^mathur-princeton]: Princeton CITP — "Dark Patterns at Scale" project page. https://webtransparency.cs.princeton.edu/dark-patterns/ — docs/data — the 7 categories, 15 types, per-pattern prevalence, verbatim definitions.
[^mathur2021]: Mathur, Kshirsagar & Mayer, "What Makes a Dark Pattern… Dark?" CHI 2021, arXiv 2101.04843. https://arxiv.org/pdf/2101.04843 — paper — no single definition; 6 attributes / 2 themes (incl. Disparate Treatment); 4 normative lenses; decoy effect; effect heterogeneity.
[^eu-behavioural]: European Commission, "Behavioural study on unfair commercial practices in the digital environment: dark patterns and manipulative personalisation" (2022, DOI 10.2838/859030). https://op.europa.eu/en/publication-detail/-/publication/606365bc-d58b-11ec-a95f-01aa75ed71a1 — study (primary) — 97% of popular EU sites/apps used ≥1 dark pattern; top-5 patterns.
[^nng-shaming]: Nielsen Norman Group, "Stop Shaming Your Users for Micro Conversions." https://www.nngroup.com/articles/shaming-users/ — practitioner(NNG) — confirmshaming/manipulinks; acceptance-link bias; fixes.
[^uxbooth]: UX Booth, "UX Dark Patterns: Manipulinks and Confirmshaming." https://uxbooth.com/articles/ux-dark-patterns-manipulinks-and-confirmshaming/ — blog — negative opt-out; before/after neutral-decline fix.
[^edpb-accounts]: EDPB, "Recommendations 2/2025 on mandatory user accounts" (2025). https://www.edpb.europa.eu/system/files/2025-12/edpb-recommendations-202502-mandatory-user-accounts_en.pdf — regulator — guest mode is the privacy-protective lawful default (GDPR Art. 25).
[^cookiebot]: Cookiebot/Usercentrics, "What are dark patterns." https://www.cookiebot.com/en/dark-patterns/ — vendor/explainer — privacy-zuckering; opt-in-not-opt-out fix; DSA/CCPA exposure list.
[^thinkdesign]: think.design, "Friend Spam / Contact Harvesting." https://think.design/blog/friend-spam-or-contact-harvesting/ — blog — the OAuth address-book trick + full-disclosure fix.
[^nouwens]: Nouwens, Liccardi, Veale, Karger & Kagal, "Dark Patterns after the GDPR," CHI 2020 (consent-banner field study). https://arxiv.org/abs/2001.02479 — research — removing "reject all" from the first page raises consent ~22–23 percentage points.
[^cnil-formal]: CNIL, "Dark patterns in cookie banners: CNIL issues formal notice." https://www.cnil.fr/en/dark-patterns-cookie-banners-cnil-issues-formal-notice-website-publishers — regulator — reject must be as easy as accept; the specific banner tells; equally-prominent first-layer reject.
[^techcrunch2022]: TechCrunch, "France spanks Google $170M, Facebook $68M over cookie consent dark patterns" (Jan 2022). https://techcrunch.com/2022/01/06/cnil-facebook-google-cookie-consent-eprivacy-breaches/ — news — €150M/€60M cookie fines, reject-not-as-easy-as-accept.
[^hunton]: Hunton, "CNIL Fines Big Tech €210M for Cookie Violations." https://www.hunton.com/privacy-and-information-security-law/cnil-fines-big-tech-companies-210-million-euros-for-cookie-violations — news/legal — Art. 82 French DPA basis; fine breakdown; injunctions.
[^belgian-noyb]: noyb, "Belgian DPA settlement turned into proper legal orders on deceptive cookie banners" (Sep 2024). https://noyb.eu/en/noyb-win-belgian-dpa-settlement-turned-proper-legal-orders-deceptive-cookie-banners — advocacy/news — reject-on-first-layer order; deceptive button colours; analytics cookies need consent.
[^ftc-negopt-policy]: FTC, "Enforcement Policy Statement Regarding Negative Option Marketing" (Oct 2021). https://www.ftc.gov/system/files/documents/public_statements/1598063/negative_option_policy_statement-10-22-2021-tobureau.pdf — agency-doc — ROSCA elements; "simple cancellation = as easy as enrollment, same medium"; why §5 pairs with ROSCA.
[^vonage]: FTC, "FTC Action Against Vonage Results in $100 Million…" (Nov 2022). https://www.ftc.gov/news-events/news/press-releases/2022/11/ftc-action-against-vonage-results-100-million-refunds-illegally-charging-consumers-junk-fees — agency-doc — "panoply of hurdles… dark patterns"; $100M; ROSCA §8403.
[^amazon]: FTC, "FTC Secures Historic $2.5 Billion Settlement Against Amazon" (Sep 2025). https://www.ftc.gov/news-events/news/press-releases/2025/09/ftc-secures-historic-25-billion-settlement-against-amazon — agency-doc — $1B penalty + $1.5B redress; Prime "Iliad" cancel maze; ROSCA + §5. VOLATILE.
[^cfpb]: CFPB, "Circular 2023-01: Unlawful Negative Option Marketing Practices." https://files.consumerfinance.gov/f/documents/cfpb_unlawful-negative-option-marketing-practices-circular_2023-01.pdf — regulator — free-trial-to-pay disclosure duties; older-consumer harm.
[^drip-bu]: BU Law / SSRN drip-pricing experiments. https://scholarship.law.bu.edu/cgi/viewcontent.cgi?article=5098&context=faculty_scholarship — research — drip pricing raises price paid ~10–13%; users stick with the worse option.
[^baymard]: Baymard Institute — "Providing a Total Order Cost" (#825) + cart-abandonment research. https://baymard.com/guidelines/825-providing-a-total-order-cost — practitioner(Baymard) — all-in pricing up front; ~70% abandon, "extra costs too high" the #1 reason.
[^nng-scarcity]: NN/G-cited honest-vs-fake scarcity practitioner synthesis. https://www.nngroup.com/articles/deceptive-patterns/ (and scarcity practitioner pieces) — practitioner — operational honest-scarcity line; FTC names reset timers as deceptive.
[^scarcity-backfire]: "When product scarcity backfires," *Journal of Retailing* (2025). https://www.sciencedirect.com/science/article/abs/pii/S0022435925001022 — research — "scarce-insincere inference": fake scarcity reduces purchase intent.
[^8thcir]: U.S. Court of Appeals for the 8th Cir., *Custom Communications, Inc. v. FTC*, No. 24-3137, 142 F.4th 1060 (8th Cir., 8 Jul 2025). https://ecf.ca8.uscourts.gov/opndir/25/07/243137P.pdf — court-opinion (primary) — vacatur of the 2024 Negative Option Rule on § 57b-3(b)(1) procedural grounds; reinstated the 1973 rule. VOLATILE.
[^ftc-anprm]: FTC ANPRM on the Negative Option Rule (2026) + Gibson Dunn analysis. https://www.gibsondunn.com/ftc-restarts-negative-option-rulemaking-after-eighth-circuit-vacatur/ — agency-doc/legal — FTC re-opened rulemaking; ROSCA continues to govern. VOLATILE.
[^usc45]: 15 U.S.C. § 45 (FTC Act §5). https://uscode.house.gov/view.xhtml?req=(title:15%20section:45%20edition:prelim) — legal-text (primary) — UDAP + § 45(n) unfairness standard.
[^ftc-unfairness]: FTC Policy Statement on Unfairness (1980) + Federal Reserve UDAP reprint. https://www.ftc.gov/legal-library/browse/ftc-policy-statement-unfairness — agency-doc — three-part deception + unfairness tests.
[^ftcreport]: FTC Staff Report, "Bringing Dark Patterns to Light" (Sep 2022, P214800). https://www.ftc.gov/system/files/ftc_gov/pdf/P214800%20Dark%20Patterns%20Report%209.14.2022%20-%20FINAL.pdf — agency-doc (primary) — four categories; "N people viewing" + reset-timer + disguised-ad + false-hierarchy examples; "not dark if true."
[^ftc-press]: FTC press release, "FTC Report Shows Rise in Sophisticated Dark Patterns" (Sep 15, 2022). https://www.ftc.gov/news-events/news/press-releases/2022/09/ftc-report-shows-rise-sophisticated-dark-patterns-designed-trick-trap-consumers — agency-doc — four categories + ABCmouse/LendingClub/Vizio examples.
[^rosca]: 15 U.S.C. § 8403 (ROSCA). https://uscode.house.gov/view.xhtml?req=granuleid:USC-prelim-title15-section8403 — legal-text (primary) — the three online negative-option requirements (verbatim).
[^ccpa-civ]: Cal. Civ. Code § 1798.140 (leginfo). https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=CIV&sectionNum=1798.140 — legal-text (primary) — § (l) dark-pattern definition + § (h) "dark patterns ≠ consent" (verbatim).
[^ccpa-regs]: CPPA Final Regulations, 11 CCR § 7004. https://cppa.ca.gov/regulations/pdf/20230329_final_regs_text.pdf — legal-text (primary) — symmetry-in-choice; "effect regardless of intent"; common-use not a defense.
[^cppa-advisory]: CPPA Enforcement Advisory No. 2024-02 (Sep 2024). https://cppa.ca.gov/pdf/enfadvisory202402.pdf — agency-doc — "dark patterns aren't about intent, they're about effect"; symmetrical choices.
[^sephora]: California AG (Bonta), Sephora settlement press release (Aug 24, 2022). https://oag.ca.gov/news/press-releases/attorney-general-bonta-announces-settlement-sephora-part-ongoing-enforcement — agency-doc (primary) — $1.2M; failure to honor Global Privacy Control. VOLATILE (amount).
[^iapp-sephora]: IAPP, "California AG announces first CCPA enforcement action." https://iapp.org/news/a/california-attorney-general-announces-first-ccpa-enforcement-action — news/legal — "first" CCPA action; GPC emphasis.
[^colorado]: Colorado Privacy Act, C.R.S. § 6-1-1303(9) (signed SB 21-190) + 4 CCR 904-3 regs. https://leg.colorado.gov/sites/default/files/2021a_190_signed.pdf — legal-text (primary) — dark-pattern definition; consent via dark patterns void; "commonly used is not, alone, enough."
[^ctdpa]: Connecticut Data Privacy Act, Conn. Gen. Stat. ch. 743jj § 42-515(14). https://www.cga.ct.gov/current/pub/chap_743jj.htm — legal-text (primary) — definition incl. "any practice the FTC refers to as a dark pattern"; consent exclusion.
[^adobe]: FTC, "FTC Takes Action Against Adobe… for Hiding Fees" (Jun 17, 2024). https://www.ftc.gov/news-events/news/press-releases/2024/06/ftc-takes-action-against-adobe-executives-hiding-fees-preventing-consumers-easily-cancelling-software — agency-doc (primary) — hidden ETF, preselected plan, cancellation hurdles; ROSCA + §5. Status: pending (verified 2026-06-16).
[^dsa]: EUR-Lex, Regulation (EU) 2022/2065 (DSA), Art. 25 + Recital 67. https://eur-lex.europa.eu/eli/reg/2022/2065/oj/eng — legal-text (primary) — the dark-pattern prohibition; UCPD/GDPR carve-out; Recital 67 definition + examples.
[^dsa-cms]: CMS DigitalLaws — DSA Art. 25 (verbatim mirror). https://www.cms-digitallaws.com/en/dsa/article-25/ — legal-text (mirror) — Art. 25(1)-(3) incl. the lettered examples.
[^dsa-penalty]: DSA Art. 52 penalties (CMS/eu-digital-services-act mirror). https://www.eu-digital-services-act.com/Digital_Services_Act_Article_52.html — legal-text (mirror) — up to 6% worldwide turnover.
[^dsa-enforce]: EPRS, "Enforcing the DSA: state of play" (2024). https://epthinktank.eu/2024/11/21/enforcing-the-digital-services-act-state-of-play/ — regulator-adjacent — Commission vs DSC enforcement split; 17 Feb 2024 application.
[^springer-frag]: "Fragmenting Consumer Law: DMA, DSA, GDPR, EU Consumer Law," *J. Consumer Policy* (2025). https://link.springer.com/article/10.1007/s10603-025-09584-3 — study — the Art. 25(2) lex specialis carve-out analysis.
[^taylorwessing]: Taylor Wessing, "DSA: Dark Patterns and other current issues." https://www.taylorwessing.com/en/interface/2023/predictions-2023-part-2/digital-services-act-current-issues — legal — Art. 25 vagueness; no in-article definition.
[^dsa-observatory]: DSA Observatory, "Platforms still use manipulative design despite DSA rules" (Aug 2025). https://dsa-observatory.eu/2025/08/07/investigation-platforms-still-use-manipulative-design-despite-dsa-rules/ — study/advocacy — Art. 25 under-enforced (disconfirming).
[^edpb-guidelines]: EDPB Guidelines 03/2022 on Deceptive design patterns in social media platform interfaces (v2.0 final). https://www.edpb.europa.eu/our-work-tools/our-documents/guidelines/guidelines-032022-deceptive-design-patterns-social-media_en — regulator-guidance — six categories + GDPR mapping + worked examples.
[^edpb-news]: EDPB news, "EDPB publishes three guidelines following public consultation" (Feb 2023). https://www.edpb.europa.eu/news/news/2023/edpb-publishes-three-guidelines-following-public-consultation_en — regulator — 14 Feb 2023 final adoption; "dark pattern"→"deceptive design pattern."
[^legiscope]: Legiscope, "European Data Protection Board: Role, Guidelines, Decisions." https://www.legiscope.com/blog/european-data-protection-board-edpb.html — regulator-adjacent — EDPB guidelines non-binding-but-authoritative; Art. 65 decisions binding.
[^gdpr-recital32]: GDPR Recital 32 + Arts. 4(11)/7. https://gdpr-info.eu/recitals/no-32/ — legal-text (primary) — consent by clear affirmative act; silence/pre-ticked boxes/inactivity ≠ consent; withdrawal symmetry.
[^edpb-consent]: EDPB Guidelines 05/2020 on consent. https://www.edpb.europa.eu/sites/default/files/files/file1/edpb_guidelines_202005_consent_en.pdf — regulator-guidance — pre-ticked invalid; scrolling ≠ consent; unbundling.
[^planet49]: CJEU C-673/17 *Planet49* (1 Oct 2019). https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:62017CJ0673 — legal-text (case law) — pre-checked cookie box is not valid consent.
[^ucpd]: EU UCPD Commission Notice, OJ C 526/2021 (Dec 2021) + Covington analysis. https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=oj:JOC_2021_526_R_0001 — legal-text (soft law) — UCPD covers dark patterns; Arts. 5/6-7/8-9 mapping.
[^edpb-interplay]: EDPB Guidelines 3/2025 on DSA–GDPR interplay (2025). https://www.edpb.europa.eu/system/files/2025-09/edpb_guidelines_202503_interplay-dsa-gdpr_v1_en.pdf — regulator-guidance — cross-instrument coordination; DMA Art. 13 anti-dark-pattern clause.
[^dfa]: European Parliament Legislative Train — Digital Fairness Act. https://www.europarl.europa.eu/legislative-train/theme-protecting-our-democracy-upholding-our-values/file-digital-fairness-act — legal-text (pipeline) — DFA on the 2026 work programme; not yet law. TENTATIVE/VOLATILE.
[^chamber]: U.S. Chamber of Commerce amicus, *FTC v. Amazon* (CourtListener). https://storage.courtlistener.com/recap/gov.uscourts.wawd.323520/gov.uscourts.wawd.323520.98.1.pdf — court (adversarial position) — FTC dark-pattern definition "far too indeterminate," "fail to provide fair notice."
[^itif]: ITIF (D. Castro), "The FTC's Efforts to Label Practices 'Dark Patterns'… Regulatory Overreach" (Jan 2023). https://itif.org/publications/2023/01/04/the-ftcs-efforts-to-label-practices-dark-patterns-are-an-attempt-at-regulatory-overreach/ — think-tank (disconfirming) — evidence "thin"; conflates bad design with illegality.
[^goldman]: E. Goldman, CPPA rulemaking comments. https://digitalcommons.law.scu.edu/cgi/viewcontent.cgi?article=3705&context=historical — academic (disconfirming) — "reconsider the definition of 'dark pattern'"; 7004 may exceed "user interface" authority.
[^promarket]: ProMarket, "How Loper Bright and the End to Chevron Impact the FTC" (Sep 2024). https://www.promarket.org/2024/09/05/how-loper-bright-and-the-end-to-the-chevron-doctrine-impact-the-ftc/ — academic/news (disconfirming) — expansive §5 reading more legally vulnerable post-Chevron.
[^llm-audit]: "LLM-Driven Agents for Dark Pattern Audits," arXiv 2603.03881. https://www.arxiv.org/pdf/2603.03881 — research — automated audits still mostly manual; LLM detectors assistive, not authoritative.
