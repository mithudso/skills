<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** The affective / emotional-resonance and
> visual-trust/credibility axis of design critique: *does this design feel right, trustworthy, and
> emotionally on-tone at first glance* — and how to detect → rate → fix problems on that axis.
> Sibling topics in this family are reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").
>
> **Scope boundary — do not duplicate siblings.** This file owns the *affective* critique axis
> (first impression, emotional tone, perceived trust/credibility, brand personality conveyed by
> visuals). It is the *subjective-emotional* complement to: the **aesthetic-usability effect**, owned by
> `usability-heuristics-laws-of-ux` (cross-reference, never re-document); **visual-composition
> principles** (hierarchy, gestalt, C.R.A.P., balance), owned by `visual-design-principles-and-critique`;
> **how to run/facilitate a crit**, owned by `visual-design-critique-methodology`; **objective
> machine-computable metrics**, owned by `computational-aesthetics-ui-metrics`; **WCAG/accessibility**,
> owned by `accessibility-ux-reviewer`. General-psychology trust/emotion theory (Mayer ABI, the Trust
> Equation, emotion regulation) lives in the `applied-psychology` hub. Prompting a VLM to critique a
> screenshot → `ai-mcp-sdk-prompting`.

---
name: emotional-design-and-visual-trust
title: Emotional Design & Visual Trust/Credibility
description: >-
  The affective / first-impression / trust axis of graphic-and-UI critique — "does this design feel
  right, trustworthy, on-tone." Norman's visceral/behavioral/reflective levels applied to an asset; the
  first-impression aesthetic judgment, halo effect, and aesthetic→trust pathway; Stanford/Fogg web
  credibility and the visual trust signals that build vs erode it; brand personality from
  type/color/imagery/shape (Aaker, warmth×competence); and the CRITIQUE pass — 5-second/squint tests,
  an emotional audit, severity rating, fixes. TRIGGER: critique a design's first impression /
  emotional resonance / "vibe" / trustworthiness / credibility / brand tone; "does this look
  trustworthy", "emotional tone audit". SKIP: the aesthetic-usability effect →
  usability-heuristics-laws-of-ux; visual-composition or running a crit →
  visual-design-principles-and-critique / -critique-methodology; WCAG → accessibility-ux-reviewer;
  general trust/emotion psychology → applied-psychology; photographic/AI-image critique.
category: developer
version: "1.0.2"
updated: "2026-06-16"
metadata:
  changelog:
    - "2026-06-16 v1.0.1->v1.0.2: /dr reciprocal cross-ref added (composition bullet routing deception-driven trust erosion to deceptive-design-and-dark-patterns); related_skills updated. <5% length delta."
    - "2026-06-16 sko v1.0.0->v1.0.1: Pass H 10/10 pos, 0/10 neg (predicted); fixed 4 Medium (Pass I anti-pattern dup -> cross-ref visual-design-critique-methodology; Pass I keyword 'aesthetic usability effect' removed; Pass J section-7 prose -> table; Pass K em-dash density 2.40->0.85/100) + 3 Low (sibling owner pointers); Pass J length retained (on-demand hub reference, no remaining duplication)"
keywords:
  - emotional design
  - visual trust
  - web credibility
  - first impression
  - Norman three levels
  - halo effect
  - brand personality
  - trust signals
  - aesthetic to trust pathway
  - emotional tone audit
  - visceral behavioral reflective
  - warmth competence
tags:
  - developer
  - frontend
  - design
  - ux
  - critique
  - trust
whenToUse:
  - "critiquing a design's first impression, emotional resonance, 'vibe', or feel"
  - "judging whether a UI screen or brand asset looks trustworthy or credible"
  - "auditing the emotional tone / brand personality a design conveys (warmth, competence)"
  - "flagging visual trust signals that are missing, weak, or actively eroding credibility"
  - "running a visceral/behavioral/reflective (Norman) emotional-design audit pass"
  - "rating the severity of an affective or trust defect and turning it into a fix"
whenNotToUse:
  - "the aesthetic-usability effect itself — use usability-heuristics-laws-of-ux"
  - "named visual-composition principles (hierarchy, gestalt, C.R.A.P., balance) — use visual-design-principles-and-critique"
  - "how to run or facilitate a critique session — use visual-design-critique-methodology"
  - "accessibility-only WCAG audit — use accessibility-ux-reviewer"
  - "general trust/emotion psychology (Mayer ABI, Trust Equation, emotion regulation) — use applied-psychology"
  - "critiquing photographs or AI-generated images — out of scope"
related_skills:
  - frontend-ui
  - visual-design-principles-and-critique
  - usability-heuristics-laws-of-ux
  - accessibility-ux-reviewer
  - applied-psychology
  - deceptive-design-and-dark-patterns
---

# Emotional Design & Visual Trust / Credibility

The **affective critique axis**: when you evaluate a graphic or UI design, this is the lens that asks
*does it feel right, trustworthy, and emotionally on-tone — in the first glance, before anyone reads a
word?* It is the subjective-emotional complement to the mechanical principles and heuristics that other
references in this hub own. Run it **first and early** in a multi-lens critique, on a static comp or
hi-fi mockup, because the first-impression read frames every downstream judgment (see §3, the halo
mechanism) and because the first-glance tests require a static artifact.[^lind][^fogg-pi]

This file is **for critiquing** a design on this axis: what visual signals to inspect, how to rate
them, and how to fix them. The theory is here only as far as it makes the critique defensible.

## Contents

- **§1: Norman's three levels** (visceral / behavioral / reflective) as the audit spine
- **§2: First impressions & the halo** (the 50ms finding, aesthetic→trust pathway, what's real vs overstated)
- **§3: Visual web credibility** (Stanford/Fogg: the 10 guidelines, "design look = 46%", surface credibility, trust signals)
- **§4: Brand personality from visuals** (Aaker's five dimensions, warmth×competence, color/type/shape/imagery → tone)
- **§5: The critique pass** (first-glance protocols, audit structure, severity rating, fix vocabulary, composition, anti-patterns)
- **§6: Quick-reference checklist**
- **§7: Per-claim confidence & contested claims**
- **## References**

---

## §1: Norman's three levels of emotional design (the audit spine)

Don Norman's *Emotional Design* (2004) frames human response to designed things at three levels of
processing — **visceral, behavioral, reflective** — derived from formal affect theory (Ortony, Norman
& Revelle).[^norman-jnd][^norman-ortony][^nng-ed] For a critic, the value is that each level is a
*different, inspectable facet* of the asset and a *different kind of user reaction*. Read each level off
the artifact:

| Level | What it is | What you inspect on the asset | The critic's diagnostic question |
| --- | --- | --- | --- |
| **Visceral** | Pre-conscious gut reaction; "pure style, pure surface"; automatic, fast, largely universal; the "wow effect".[^norman-ortony][^norman-jnd] | Color & contrast, imagery, typography *at first glance*, polish/craft, motion, white space, overall organization — the impression before any reading.[^cheesecake][^ixdf-3levels][^jarosz] | "What's my immediate, pre-conscious reaction — attractive/safe/inviting vs ugly/cluttered/forbidding? Does the look telegraph the right feeling and the product's purpose *before* I read anything?" |
| **Behavioral** | Pleasure & effectiveness of *use*; expectation-induced; "good tools are invisible"; violated expectations → frustration/anger.[^norman-ortony][^norman-jnd] | Perceived ease/clarity, affordance & signifier clarity, consistency, feedback after actions, error handling, perceived control.[^jarosz][^cheesecake] | "Does it *look like* it will be easy and let me feel in control?" — **then defer the rigorous audit to the usability lens.** |
| **Reflective** | Meaning, self-image, story, brand; conscious; pride/shame/admiration; most culturally variable.[^norman-ortony][^norman-jnd] | Brand story/identity coherence, status/prestige cues, perceived value & trust, self-expression, cultural references, memorability.[^jarosz][^ixdf-3levels][^cheesecake] | "What does using or displaying this say about *me*? What story, status, or meaning does it carry afterward — attachment, or just a finished task?" |

**Ownership inside this hub.** The affective axis owns mostly the **visceral** layer (first impression,
appearance) and the **reflective** layer (trust, brand image, "what it says about me"). The
**behavioral** layer overlaps heavily with classic usability, i.e. Nielsen-heuristics territory:
*cross-reference `usability-heuristics-laws-of-ux`, do not re-own it.*[^cheesecake][^nng-au]

**"Attractive things work better": the engine behind why emotion belongs in a critique.** Norman
argues appeal and usability are *complementary*: positive affect broadens thinking and makes users
**more tolerant of minor difficulties**; anxiety narrows it (grounded in Isen's affect research,
Fredrickson's broaden-and-build, Zajonc's affective primacy).[^norman-au][^mueller] His own caveat
is the load-bearing one for a critic: he warns against **"façade design"**: "to be truly beautiful…
the product has to fulfill a useful function, work well, and be usable."[^norman-au]

**Tensions a critic should name** (Norman: "real products involve continual conflicts among the three
levels"):[^mueller]
- **Viscerally beautiful but behaviorally broken**: the classic failure (Norman's own unusable-but-pretty Jensen clock). Do not let a high visceral read launder a severe usability defect (§5).[^norman-jnd][^nng-au]
- **Reflective meaning buying tolerance for weak behavioral showing**: a luxury/status product wins on reflective grounds and earns patience for friction.[^ixdf-3levels][^mueller]
- **Norman's own later self-critique** ("Where Emotional Design Fails"): award-winning, gorgeous designs can be "art-centered, prize-centered, object-centered… the one thing they are not is human-centered."[^norman-fails]

> **Calibration:** the three levels are *different axes*, not one quality score. Designers have most
> control over visceral, least over reflective (and even there only via "emotional affordances" whose
> uptake is "beyond the designer's control") — so frame reflective judgments as *intended/probable*
> meaning, hedged for audience and culture.[^norman-ortony] The relative weight of each level is
> product- and context-dependent (an alarm clock is function-first; a wall clock may be bought for
> visceral/reflective appeal).[^norman-jnd]

---

## §2: First impressions & the halo effect

### The "first 50 ms" finding: stated precisely, then bounded

Lindgaard, Fernandes, Dudek & Brown (2006), *"Attention web designers: You have 50 milliseconds to
make a good first impression!"* found that **average visual-appeal ratings of homepages shown for
50 ms correlate ~0.95 with the same pages shown for 500 ms** (r = .947, r² ≈ .90, aggregated across
pages).[^lind] The mechanism the authors invoke is a chain a critic must understand because it is *why
first impressions are high-severity*:[^lind]
1. **Mere exposure** (Zajonc): affective preference forms pre-cognitively, "before the organism has had a chance cognitively to analyse" the stimulus.
2. **Halo effect**: the first impression "carries over… to the evaluation of other attributes."
3. **Confirmation bias**: once formed, users "search exclusively for confirmatory evidence": a good first impression makes later flaws "generously overlooked"; a bad one resists revision "even in the presence of strong disconfirmatory evidence."

> **Honest caveat: preserve this contradiction.** The *stability* of fast aesthetic judgments is
> well-replicated (Tractinsky 2006: 500 ms ↔ 10 s; Tuch 2012 pushed effects to **17 ms**).[^tract06][^tuch-pt]
> The **literal "50 ms" deadline is overstated**: it rests on a single between-subjects experiment
> (n=40) with deliberately *polarized* stimuli (only very-appealing and very-unappealing pages); the
> authors explicitly say their aim was "not to determine an accurate threshold"; individual-level
> reliability was weak (40% of participants < r .60 at 50 ms); within-subject designs *inflate* these
> correlations (Thielsch & Hirschfeld); and the research does **not** say users *decide to leave* in
> 50 ms — only that an appeal *rating* is reliable.[^lind][^papachristos][^thielsch] Use "first
> impressions form in well under a second and stick" as the defensible claim; treat "50 ms" as a
> memorable gloss, not a measured deadline.

**What predicts the first-glance appeal.** Reinecke et al. (2013, CHI) found computational
**colorfulness + visual complexity, plus demographics, explain ≈48% of the variance** in
first-impression appeal (after 500 ms).[^reinecke13] The nuance (more precise than "low + low wins"):
**visual complexity dominates; colorfulness is a smaller, secondary factor**; higher complexity
generally lowers appeal (Tuch found a *linear* "lower = better"), though Berlyne's inverted-U
(moderate complexity peaks) competes; and the effect is **not universal** (age × complexity, education
× colorfulness interactions).[^reinecke13][^tuch-pt] Tuch's deeper finding: **low visual complexity +
high prototypicality** (how *typical* the layout is for its genre) yields the highest beauty ratings —
familiar, uncluttered layouts win the first glance.[^tuch-pt]

### The halo effect and the aesthetic → perceived-trust pathway

The halo effect (Thorndike 1920; "what is beautiful is good", Dion/Berscheid/Walster 1972) means a
rating on one attribute drags the others in the same direction.[^nng-halo][^dion] In interfaces, **an
attractive first impression biases judgments of trustworthiness, credibility, content quality, and
usability**, before content is read.[^nng-halo][^reinecke13] The trust-specific evidence:
- **Robins & Holmes (2008):** presenting **identical content** in higher- vs lower-aesthetic treatments, the higher-aesthetic version was judged **more credible in 19 of 21 cases (90%)** — an "amelioration effect" operating "within the first few seconds," "probably limited to the visceral level."[^robins]
- **Lindgaard et al. (2011):** appeal, trust, and perceived usability are *all* largely driven by visual appeal, but trust is **anchored to appeal, not identical to it** (it relies on somewhat different visual attributes and is processed more deliberately).[^lind11]
- Corroborated by Kwak et al. (aesthetics the strongest system feature for perceived quality/reputation via the attractiveness halo) and "What is beautiful is secure" (appeal drove *perceived security* over real identity verification).[^kwak][^secure]

> **Boundary conditions on the halo (preserve the contradiction).** The halo is **moderate and
> conditional, not universal.** Eagly et al.'s (1991) meta-analysis ("…But…") found the effect
> **content-specific**: largest for perceived *competence*, **near zero for integrity / concern for
> others**, and **weakened by individuating information**. So a pretty UI inflates perceived
> *competence/quality* but does little for perceived *honesty*, and real content can override
> it.[^eagly] First-impression preferences also vary **substantially by culture, gender, age, and
> education** (Reinecke & Gajos: 2.4M ratings / ~40k people; 80,901 sites across 44 countries) — there
> is no single universal "appealing"; appeal is target-audience-relative.[^gajos14][^gajos18]

> **Cross-reference, do not duplicate:** the **aesthetic-usability effect** (users tolerate/overlook
> *usability* problems in attractive UIs) is the usability-specific special case of the halo and is
> owned by `usability-heuristics-laws-of-ux`. Note only that its forgiveness is **bounded** — "a
> pretty design can make users forgiving of *minor* usability problems, but not of *large*
> ones."[^nng-au]

---

## §3: Visual web credibility & trust signals (Stanford / Fogg)

The **Stanford Web Credibility Project** (B.J. Fogg's Persuasive Technology Lab) is the canonical
research on what makes a site *believed*. Two results anchor the visual-trust axis.

**"Design look" is the single most-mentioned credibility factor.** In Fogg et al. (2003), 2,684
participants evaluated live sites; coding their free-text comments, **"Design Look" appeared in 46.1%
of comments**, far ahead of information design (28.5%), information focus (25.1%), and everything
else.[^fogg-46][^fogg-acm] Fogg's conclusion: "people do judge a Web site by how it looks. That's the
first test… if it doesn't look credible… they go elsewhere."[^wiki-scwc]

> **Calibrate the famous number.** 46.1% is the share of *coded comments that mentioned design* (what
> people most **notice and remark on**) — **not** "46% of credibility is statistically explained by
> design," and **not** a causal coefficient. State it as *"design look was the single most-mentioned
> credibility factor (46.1% of comments)."* Domain experts are far less swayed by surface than general
> consumers.[^fogg-46][^bokardo]

**Prominence-Interpretation Theory** explains *how*: credibility impact = **Prominence** (a user
notices an element) × **Interpretation** (the user judges it). If either is zero, that element has no
credibility effect.[^fogg-pi][^nng-pi] This resolves the apparent paradox that surveys rank
contact-info/privacy-policies as important, yet free evaluation is dominated by "design look": design
is simply the most **prominent** thing.[^fogg-pi]

**Fogg's four credibility types**: and which one visual design drives:[^fogg-elements][^fogg-pi]
- **Presumed**: from the user's assumptions/stereotypes (e.g., `.org` reads as more trustworthy).
- **Reputed**: from third-party endorsement: seals, awards, certifications, inbound links. *(Fogg also used "referred" for the link/recommendation variant — same third-party family.)*
- **Surface**: from simple inspection / first impression ("judge a book by its cover"). **This is the type visual design directly drives**, and the home of the affective-trust axis.
- **Earned**: from firsthand experience over time.

**The Stanford Guidelines for Web Credibility (10), abbreviated for inspection** (Fogg 2002):[^stanford-guidelines][^wiki-scwc]
1. Make it easy to **verify the accuracy** of your information (cite/link sources).
2. Show there's a **real organization** behind the site (physical address, office photo).
3. **Highlight expertise** (credentials, respected affiliations; don't link to non-credible sites).
4. Show **honest, trustworthy people** stand behind it (real people, bios).
5. Make it **easy to contact you** (phone, address, email).
6. **Design the site to look professional** (and *appropriate for its purpose*) — "people quickly evaluate a site by visual design alone."
7. Make the site **easy to use and useful**.
8. **Update content often** (or show it's recently reviewed).
9. **Use restraint with promotional content** (avoid pop-ups; distinguish ads from content).
10. **Avoid errors of all types**, however small (typos, broken links — they "hurt more than most people imagine").

At least six of the ten are visual/design-surfaced; guideline 6 is explicitly about look.

### Trust signals: BUILDERS vs ERODERS (the inspection list)

**Builders (raise perceived credibility):**[^stanford-guidelines][^nng-trust][^nng-comm][^baymard-seal][^baymard-sec]
- **Professional design matched to the genre/purpose** (the strongest surface cue; Jakob's Law — users expect your site to look like the category they know).
- **Real photos of real people / evidence of a real organization** (named experts, office, team bios).
- **Visible contact info**: phone, physical address, email (not just a form).
- **Specific social proof & third-party reviews**: testimonials with full name + photo + company + concrete result; external reviews (G2/Trustpilot/Google) outweigh on-site quotes users discount as cherry-picked.
- **Recognized trust/security cues at the point of risk**: SSL/padlock, recognized payment icons, "your information is encrypted" *next to the form*; Baymard: users judge payment-form security *visually* ("trust their gut"), and **the brand on a seal matters more than what it certifies**.
- **Fresh, comprehensive, correct content; fast, working, error-free site.**

**Eroders (damage credibility):**[^stanford-guidelines][^fogg-chi01][^nng-comm][^pixxen]
- Typos, broken links, visual errors.
- Aggressive, deceptive, or pop-up ads.
- Dated or **amateurish** design (Fogg's CHI 2001 study: the two factors that *hurt* credibility were **commercial implications** and **amateurism**).
- Stock-photo overload; haphazard/incomplete content (missing product photos read as "brittle").
- Inconsistent branding / wrong look for the genre.
- **Over-slickness** that reads as marketing gloss (the "too slick" backlash in the Stanford data).
- **Dark patterns**, hidden costs, premature login walls (NN/g "Hierarchy of Trust": don't make high-commitment demands before lower-level trust is established).

> **Durability of the old study.** The core Stanford findings are from **2002–2003** (pre-mobile,
> pre-modern-dark-patterns; the original page is partly link-rotted) — flag the age. But NN/g's 2016
> re-test of Nielsen's trust factors found them still holding ("design patterns change… human behavior
> does not"), and Alberts & van der Geest (2011) replicated "visual cues over content." Treat the
> **principles as durable, the surface specifics as dated.**[^nng-trust][^makinggood]

> **Boundary, not duplicated here:** general interpersonal/organizational trust models — **Mayer's ABI
> (ability-benevolence-integrity)** and the **Trust Equation** — explain *why humans extend trust* and
> live in `applied-psychology`. The Stanford/Fogg work is specifically about **web/visual surface
> credibility cues.** Cross-reference; don't merge.

---

## §4: Brand personality conveyed by visual choices (the emotional-tone audit)

### The two frameworks

**Aaker's five dimensions of brand personality** (1997, the canonical framework, "the set of human
characteristics associated with a brand," serving a symbolic/self-expressive function):[^aaker]
**Sincerity** (honest, wholesome, down-to-earth), **Excitement** (daring, spirited, trendy),
**Competence** (reliable, intelligent, successful), **Sophistication** (glamorous, upper-class,
charming), **Ruggedness** (tough, outdoorsy). This is the vocabulary for naming *what personality a
design projects.*

**Warmth × Competence** (the Stereotype Content Model — Fiske, Cuddy, Glick — applied to brands as the
**Brands as Intentional Agents Framework**, Kervyn, Fiske & Malone 2012). People judge brands, like
people, on **warmth** ("what are its intentions toward me?") and **competence** ("can it act on
them?"), forming four quadrants each with a distinct emotion: warm+competent → **admiration**;
cold+competent → **envy** (luxury brands); warm+incompetent → **pity**; cold+incompetent →
**contempt**. Both dimensions independently predict purchase intent and loyalty.[^biaf] A telling
manipulation: switching a page's domain from `.com` to `.org` made the *same* company read as **warmer
but less competent**: small visual/contextual cues move the warmth-competence position.[^biaf]

### How each visual element pushes tone

| Element | Pushes *warm / soft / approachable* | Pushes *competent / cold / serious* | Evidence & caveat |
| --- | --- | --- | --- |
| **Color** | warmer/softer palettes; (red → *excitement*) | blue → *competence* (most common logo color) | Labrecque & Milne "Exciting red and competent blue"; saturation & value matter, not just hue. **Heavily caveated** — see below.[^labrecque] |
| **Type** | script/handwriting (casual, personal, elegant); rounded sans | serif (stable, mature, formal); precise geometric sans | Fonts carry consistent personalities *by family*; **congruence** is the real effect — an incongruent typeface damages perceived message *and author credibility*.[^shaikh][^fox] |
| **Shape/contour** | rounded/curved (friendly, safe) | angular/sharp (strong, capable) | Bar & Neta: humans prefer curves; sharp contours read as *threat* (amygdala). Galli & Chattopadhyay: circular logos → "softness", angular → "hardness", transferring to product judgments. Frontiers: angular → competence, rounded → warmth.[^barneta][^galli][^frontiers-shape] |
| **Imagery** | authentic/candid human photography; illustration | polished studio, restrained | **Authentic beats stock and (currently) AI** for trust (Getty: 98% say authentic imagery pivotal to trust; Cornell: real photos beat stock). Over-modified images *reduce* authenticity → lower trust.[^getty][^cornell] |

> **The honest caveats: this is where pop "color psychology" fails.** The authoritative review
> (Elliot & Maier 2014, *Annual Review of Psychology*) concludes color *can* carry meaning but the
> literature is "at a nascent stage" and **not yet ready for strong applied recommendations**; the
> governing principle is **Color-in-Context** — the same color carries different, even opposite,
> meanings by context (red = failure in achievement contexts, attraction in mating contexts).[^elliot]
> Color–emotion associations are **both universal and culture-specific** (Jonauskaite: 4,598 people /
> 30 nations; red is strongly *positive* in Chinese culture), and color alone explains only **~5–10% of
> emotional-response variance**.[^jonauskaite][^colorarchive] **Do not assert single-color "meanings"
> as universal laws.** Likewise: **"serif is more readable than sans" is largely a myth** (Lund's
> review of 28 studies; a 42-study review found no reading-performance difference from serifs alone —
> familiarity dominates) — use typeface *tone/personality* (robust), not readability superiority.[^lund][^hosang]
> And distrust the famous stats: **"90 seconds / 62–90% of assessment from color"** and **"consistent
> color → +80% recognition / +33% revenue"** are widely repeated but **untraceable pop-stats** — assert
> only the *direction* (consistency aids recognition; color is mainly a fit/differentiation asset).[^singh][^colorarchive]

### Congruence is the load-bearing finding

Across color, type, and shape, **fit beats intrinsic meaning**: matching element semantics to the
brand/message semantics improves perception; *mismatch degrades it.*[^bottomley][^fox][^galli] The
emotional-tone audit's core question: **does the visual language (type + color + imagery + shape +
layout) match the intended brand personality AND the audience's category expectations?** Canonical
mismatch findings: "looks like a fintech but it's a children's charity," a luxury brand "undermined
by stock photography," a bank using red that triggers anxiety where trust/competence is
expected.[^advergize][^adessi] And **consistency** across touchpoints reinforces personality (reduces
cognitive friction, builds the brand halo); shifting the visual language "resets the recognition
clock."[^cohesion][^adessi]

---

## §5: The critique pass: run → rate → fix

This is the operational synthesis. The pass produces findings as a **triple**:

> **Observation → Severity → Fix** — or, for an explicitly affective audit, **Observation → Change →
> Expected emotion** (name the affective delta the change should produce).

A finding like *"User feels frustrated"* is not actionable; *"Value prop illegible at 5 s; users can't
say what the product is"* is.[^uxcrush][^nng-journey]

### 5a. First-glance protocols (run these first)

**The 5-second / "blink" test** probes *comprehension + recalled impression*. Show a static design for
exactly 5 s, hide it, then ask **3–5 open-text-first** questions — **without priming** (don't name the
brand or say "first impressions"):[^koji][^maze-5s][^cleverx][^lyssna]
- *Comprehension:* "What is this about? What product/service is offered?"
- *Audience:* "Who is this designed for?"
- *Dominant element / recall:* "What stood out? What do you remember most?"
- *Sentiment / trust:* "What was your first impression? How trustworthy did it feel (1–5)?"
- *Brand-tone:* "Three words that describe this design/brand."
- *CTA recognition:* "What did it want you to do next?"

*Failure looks like:* can't state what it is / who it's for; recall fixes on a decoration; low trust
rating; brand words contradict the intended personality.[^cleverx][^lyssna]

**The squint / blur test** probes *visual hierarchy*. Blur the design (squint + step back, or 8–15px
Gaussian blur) until only shapes/color/contrast remain — what survives is what users process
pre-attentively.[^nng-squint][^lukew][^polypane] Check: is the **primary focus still recognizable**?
Are there **competing focal points** (everything equal weight = eyes bounce)? Are **background
elements more prominent** than the primary one? Still on-brand? Wroblewski's move turns it strategic —
"which elements *should* stand out" forces the purpose/audience question.[^lukew]

> **Validity caveat: directional signal, not measurement (preserve the contradiction).** Peer-reviewed
> work questions first-impression *timing*: Kuric et al. (2023) call fixed-5s testing "pseudoscience"
> in its current state (cognitive ability × visual complexity sway early impressions — pilot the
> exposure); Gronier (2016) found **no significant difference** between a 5-second condition and a full
> usability test and warns it is **not valid for finding usability problems** and must use a static
> hi-fi artifact.[^kuric][^gronier] Use 5-second/squint outputs as *directional signal about what
> communicates visually*, keep exposure consistent within a study, and **never use them to certify
> usability.**

### 5b. Structuring the emotional audit

Inspect across Norman's **visceral → behavioral → reflective** spine (§1), owning the visceral and
reflective layers and deferring behavioral to the usability lens.[^jarosz][^cheesecake] To ground the
affective column in *data* rather than the critic's taste, recruit a qualitative emotion probe:
- **Microsoft Desirability Toolkit / Product Reaction Cards** (Benedek & Miner): 118 word-cards (~60% positive / 40% negative — the skew is deliberate); show the design, have participants pick the words, **narrow to 5**, then **interview "why" for each** — the debrief is where the insight lives.[^nng-cards][^uxfirm]
- **PrEmo** (non-verbal cartoon emotions, cross-cultural), **SAM** (pictorial pleasure/arousal/dominance), **Geneva Emotion Wheel** (20 emotion families × valence/control) — non-verbal instruments exist precisely because emotion words don't translate.[^bluehair][^gew]
- **Emotion / emotional journey maps**: one emotion curve across phases; **dips = highest-priority pain points**, each annotated with the specific friction thought.[^nng-journey][^uxcrush]

> **Validity caveat (multi-source).** These are **qualitative probes, not validated metrics.**
> MeasuringU: the Desirability Toolkit has **no published benchmarks** and "**no evidence** that [it]
> actually measures desirability." Comparative studies disagree on which instrument best matches
> pre-classified emotions — none is ground truth. Use them to *structure and triangulate* the affective
> column (combine with a validated behavioral metric like SUS and direct observation), never to assert
> a validated emotion score.[^measuringu-des][^bluehair]

### 5c. Severity rating for affective / trust findings

Borrow the **Nielsen 0–4 severity scale** verbatim and the **frequency × impact × persistence**
model, reframed for emotional/trust defects (the scale itself is owned by
`usability-heuristics-laws-of-ux`; borrowed here and re-anchored to affective findings):[^nng-severity]

| Rating | Meaning | Typical affective/trust example |
| --- | --- | --- |
| **0** | Not a problem | — |
| **1** | Cosmetic only | A single slightly-off-tone icon in a non-prominent spot |
| **2** | Minor | Mildly generic stock hero; one weak testimonial |
| **3** | Major (high priority) | First-impression hierarchy fails the squint test; brand tone mismatches the audience; visible typos |
| **4** | Catastrophe (fix before release) | Value prop illegible at 5 s; the design reads as untrustworthy/scammy; a dark pattern causing user harm |

**Why affective findings punch above their apparent weight.** A weak first impression is **high
severity** because of the prominence-interpretation + halo chain: design look is the most *prominent*
factor (§3) and the halo + confirmation bias make it *tax every downstream judgment* (§2). Map a
value-prop-illegible-at-5s or CTA-vanishes-on-squint finding to **3–4** (frequent: every visitor,
high-impact, persistent).[^fogg-pi][^fogg-46][^nng-au] A brand-tone mismatch is typically **2–3**
(persistent, brand-eroding, usually overcome-able).[^advergize][^adessi]

**Severity ≠ priority.** Severity is the *degree of experiential harm* (set by the critic); priority
is the *business decision of when to fix* (deadlines, visibility, revenue). The relationship is
non-linear: the canonical case is the **misspelled logo: low severity, high priority** (every visitor
sees it; brand trust damaged → fix now).[^betterqa][^contextqa] This is exactly why affective findings
are **often low functional severity but high business priority**: a typo, a stock photo, a tone clash
break *nothing* yet are seen by everyone and erode the trust that gates conversion. Most are also
**quick wins** (high impact, low effort: a real photo, a visible phone number, and fixing typos are free
and fast), so they sort to "do first" on a value×effort matrix.[^contextqa][^altersquare]

> **Reliability caveat (multi-source).** Single-rater severity is **unreliable**, and affective
> severity is *more* subjective than functional. Nielsen's own finding: single-evaluator ratings are
> "very unreliable," but the **mean of ~4 evaluators** lands within half a point of true severity. The
> broader "evaluator effect" ("you say disaster, I say no problem") shows inter-rater agreement often
> well under acceptable thresholds. **Mitigations: use multiple independent raters, a written
> definitions/codebook, report rater backgrounds, and compute inter-rater agreement.**[^hertzum][^nng-severity]

### 5d. Fix vocabulary (remediation moves)

A critique that only names problems is half-done on this axis, where fixes are well-catalogued and
mostly cheap. Group by the defect treated:
- **(a) Raise visceral polish**: increase contrast, add white space, replace low-quality/stock imagery, tighten typographic craft, enforce consistency.[^saleshub][^cheesecake]
- **(b) Fix first-impression hierarchy**: re-weight so the dominant element carries the core message; one primary action per screen; match visual weight to the business goal (the squint-test fix).[^lukew][^polypane]
- **(c) Add / strengthen trust signals**: real team/office/product photos (never stock), visible contact info, *specific* social proof (name + photo + company + result) and external reviews, recognized security cues *at the friction point*, named author bios/credentials, error-free copy.[^stanford-guidelines][^saleshub][^baymard-sec]
- **(d) Correct brand-tone mismatch**: align type/color/imagery/shape to the intended Aaker personality and warmth-competence position; resolve element-vs-element clashes (e.g., sophisticated palette + stock photos).[^aaker][^advergize][^fox]
- **(e) Remove trust eroders**: replace dated/amateurish design, cut stock-photo overload and aggressive/disguised ads, fix typos, and **flag dark patterns as an ethics finding** (harm pathway: manipulation → emotional distress → trust/autonomy erosion), not just a severity number. Use *real* badges only — fake/aspirational ones "erode trust faster than no badges."[^stanford-guidelines][^pixxen][^darkpattern]

### 5e. How the affective lens composes with other lenses

Affective critique is **one axis** alongside usability heuristics (behavioral level), visual-composition
principles, accessibility, and objective/computational metrics — run them as parallel, separately-scoped
domains and merge into one severity-ranked report.[^uiux-toolkit] Discipline to avoid double-counting:
- **Assign each finding to the lens that *owns* it, and count it once.** Low color contrast is an **accessibility** finding (`accessibility-ux-reviewer`) even though it also reads as low visceral polish: file it under a11y and cross-reference. "Decorative elements compete with content" lives in the visual-composition lens (`visual-design-principles-and-critique`); objective scoring (colorfulness, clutter) lives in `computational-aesthetics-ui-metrics`. Don't re-score a sibling's finding here.[^uiux-toolkit][^deque]
- **Consolidate & de-duplicate before rating**: merge findings across evaluators and lenses, remove duplicates, average severity.[^uiux-toolkit]
- **Triangulate, don't duplicate**: pair affective signal with a *validated* behavioral metric (SUS) and direct observation.[^nng-cards]
- **Trust *erosion* via deception is a separate, owned lens.** When low perceived trust traces to a *deceptive or manipulative* pattern (a fake-urgency timer, an asymmetric cookie banner, a roach-motel cancel, drip pricing), that is a design-**ethics** finding — route it to `deceptive-design-and-dark-patterns`, which rates user harm **plus** legal exposure (DSA Art. 25, GDPR/CNIL, FTC, CCPA/CPRA) and prescribes the honest fix. This affective lens owns the *first-impression trust signal*; the dark-patterns lens owns *why a manipulative interface forfeits it*.

> **The decoupling rule: don't let a high visceral score launder a real defect.** Polish inflates
> *perceived* credibility independent of actual quality (Robins & Holmes' amelioration effect), so a
> slick shell can mask weak content — and **perceived trust ≠ trustworthiness.** Score actual
> content/usability/accessibility *separately* and flag explicitly when polish and substance diverge.
> A high visceral score is a *reason to keep looking, never a reason to stop.* (This holds under the
> live contradiction in the literature: Tuch et al. 2012 and Grishin & Sauro 2019 **failed to find**
> the aesthetics→perceived-usability effect, with poor usability instead *lowering* post-use aesthetics
> — "usable is beautiful.")[^robins][^graticle][^tuch-chb][^grishin]

### 5f. Anti-patterns specific to the affective critique

The general crit anti-patterns (vague praise / Nielsen's "Make-it-Pop syndrome," the compliment sandwich, HiPPO and bikeshedding, and "diagnose, don't prescribe") are owned by `visual-design-critique-methodology` and apply unchanged here. The affective lens is the *most* vulnerable to subjectivity, so a few anti-patterns bite hardest on this axis and deserve calling out:[^nielsen-crit][^uxtigers][^prototypr]

- **Subjective taste presented as emotional fact.** "I hate blue" or "this feels cheap to me" is a sample size of one. Ban "I like / I don't like"; anchor every affective note to the intended persona, brand goal, or a cited finding. "Blue doesn't match the trust-building goal" is feedback; "I don't like blue" is not.
- **Ignoring audience and cultural context.** First-impression and emotion reads are audience-, culture-, age-, and gender-dependent (§2, §4). Anchor to *the intended* audience, never to personal or Western defaults; an "appealing" or "trustworthy" verdict is target-relative.
- **Over-indexing on first impression while ignoring behavioral substance.** A beautiful 5-second read says nothing about whether tasks complete. The mirror of the decoupling rule (§5e): "pay attention not only to what users say, but what they do."[^cleverx]

---

## §6: Quick-reference checklist

Run top to bottom on a static comp:

1. **Squint test**: does the dominant element (after blur) carry the core message? Competing focal points? Background louder than foreground? *(fail → severity 3–4)*
2. **5-second test**: can a viewer state what it is, who it's for, and what to do next? Trust rating? Brand words match intent? *(fail → 3–4)*
3. **Visceral read (Norman §1)**: polish/craft, color/contrast, imagery quality, typographic care; "wow" vs "cheap/dated." Does the look telegraph the right feeling + purpose?
4. **Trust signals (§3)**: real photos? visible contact info? specific social proof / external reviews? security cues at the friction point? error-free copy? Any eroders (typos, stock overload, aggressive ads, dark patterns, dated design)?
5. **Emotional tone / brand personality (§4)**: name the Aaker personality and the warmth×competence quadrant; does the visual language match the intended personality *and* the audience's category expectations? Element-vs-element clashes? Consistency across touchpoints?
6. **Reflective read (Norman §1)**: brand story coherence; what does choosing this say about the user; memorability.
7. **Rate** each finding 0–4 (frequency × impact × persistence), separating **severity from priority**; average across ≥1 other rater where possible.
8. **Decouple**: score actual content/usability/accessibility separately; never let visceral polish launder a real defect.
9. **Write findings** as Observation → Severity → Fix (or → Expected emotion); diagnose, don't prescribe; specific not vague.

---

## §7: Per-claim confidence & contested claims

Confidence per load-bearing claim (the inline blockquotes in §1–§5 carry the detail; this table is the index). "High" = 3+ independent sources agree; "Qualified" = real but single-primary or bounded; "Contested" = preserve the contradiction.

| Claim | Confidence | Note / boundary |
|---|---|---|
| Norman's three levels & their design signals[^norman-jnd][^norman-ortony][^ixdf-3levels][^jarosz] | High | author + IxDF + NN/g + peer-reviewed audit converge |
| "Attractive things work better" / positive-affect-broadens[^norman-au][^mueller] | High (as thesis) | Norman's synthesis of Isen/Fredrickson/Zajonc, not one replicated experiment |
| 50ms↔500ms appeal correlation ~.95 (aggregate)[^lind][^tract06] | High | stability of sub-second appeal is the robust claim |
| The literal "50 ms deadline"[^lind][^papachristos][^thielsch] | Contested | single n=40 study, polarized stimuli, authors disclaim it; within-subject inflation |
| Halo / aesthetic→perceived-trust pathway[^robins][^lind11][^kwak][^secure] | High | it is *perceived* trust, not trustworthiness |
| Halo strength[^eagly][^gajos14] | Bounded | moderate, content-specific (~0 for integrity), weakened by individuating info, culture/age/gender-variable |
| "Design look = 46.1% of credibility comments"[^fogg-46][^fogg-acm][^bokardo] | High (figure) | *most-mentioned*, not "% of credibility explained" |
| The 10 Stanford guidelines & four credibility types[^stanford-guidelines][^fogg-elements] | High | reputed-vs-"referred" is a terminology nuance |
| Aaker's five dimensions; warmth×competence / BIAF[^aaker][^biaf] | High | primary papers |
| Color → personality (red=excitement / blue=competence)[^labrecque] | Qualified | real effect, single-primary mapping |
| Color meaning is context/culture-dependent; pop color psychology overstated[^elliot][^jonauskaite] | High | the better-supported claim; do not assert single-color "meanings" |
| Typeface *tone* attribution & congruence[^shaikh][^fox] | High | use for tone |
| Serif-vs-sans *readability* superiority[^lund][^hosang] | Contested / debunked | use tone, not readability |
| Rounded→warm / angular→competent (incl. threat/amygdala)[^barneta][^galli][^frontiers-shape] | High | power/self-construal is a boundary condition |
| Authentic imagery > stock/AI for trust[^getty][^cornell] | High | authenticity may be a baseline expectation, not a bonus |
| 5-second / squint *mechanics*[^koji][^nng-squint] | High | — |
| First-impression *timing* validity[^kuric][^gronier] | Contested | directional signal, never a usability certifier |
| Emotion-measurement tools (PrEmo/SAM/GEW/Desirability)[^measuringu-des][^bluehair] | Qualified | qualitative probes without established validation |
| Nielsen 0–4 scale + freq×impact×persistence[^nng-severity] | High (verbatim) | scale owned by `usability-heuristics-laws-of-ux`; its *application to affective findings* is a flagged synthesis |
| Severity-rating reliability (average ≥4 raters)[^hertzum] | High | — |
| Untraceable pop-stats ("90 sec / 62–90% from color"; "+80% recognition / +33% revenue")[^singh][^colorarchive] | Not citable | assert only the direction |

---

## References

[^norman-jnd]: https://jnd.org/emotional-design-people-and-things/ — Norman's concise visceral/behavioral/reflective definitions; product/context weighting; the unusable-but-pretty Jensen clock. | tier: org-docs (jnd.org, primary author).
[^norman-ortony]: https://projectsfinal.interactionivrea.org/2004-2005/SYMPOSIUM%202005/communication%20material/DESIGNERS%20AND%20USERS_Norman.pdf — Norman & Ortony, "Designers and Users" (2003): the reactive/behavioral/reflective → visceral/behavioral/reflective derivation; "emotional affordances"; designers' limited reflective control. | tier: paper.
[^nng-ed]: https://www.nngroup.com/books/emotional-design/ — NN/g page for *Emotional Design*; the three levels + "Attractive Things Work Better." | tier: org-docs (NN/g).
[^ixdf-3levels]: https://ixdf.org/literature/article/norman-s-three-levels-of-design — IxDF canonical article: per-level signals; smartwatch behavioral-vs-reflective conflict. | tier: org-docs (Interaction Design Foundation).
[^jarosz]: https://dbc.wroc.pl/Content/139298/Jarosz_Methodological_Concept_for_User_Experience_Research.pdf — Jarosz (2024): peer-reviewed UI audit operationalizing Norman's model (≈9 visceral / 8 behavioral / 7 reflective criteria; good-vs-poor rating; 7 stages). | tier: paper.
[^cheesecake]: https://cheesecakelabs.com/blog/applying-emotional-design-and-ux-heuristics-in-real-projects/ — maps each level to concrete UI signals; pairs the model with Nielsen's heuristics. | tier: blog.
[^mueller]: https://journals.library.ualberta.ca/iejll/index.php/iejll/article/download/707/367 — Mueller review essay: the three-aspects table; affective primacy (Zajonc), positive-affect-broadens (Isen/Fredrickson); "continual conflicts among the three levels." | tier: paper.
[^norman-au]: https://jnd.org/emotion-design-attractive-things-work-better/ — Norman, "Attractive Things Work Better": positive affect → breadth/tolerance; the anti-"façade design" caveat. | tier: org-docs (primary; also *Interactions* 2002).
[^norman-fails]: https://jnd.org/where-emotional-design-fails/ — Norman's self-critique: award-winning emotional design can be object-centered, not human-centered. | tier: org-docs (primary).
[^nng-au]: https://www.nngroup.com/articles/aesthetic-usability-effect/ — NN/g, aesthetic-usability effect: perceived vs actual usability; the "minor not large problems" bound; beauty can mask usability problems. (Owned by usability-heuristics-laws-of-ux; cross-ref only.) | tier: org-docs (NN/g).
[^lind]: https://www.makinggood.ac.nz/media/1265/lindegaardetal_2006_attention.pdf — Lindgaard et al. (2006), full text: 50↔500 ms r=.947; methodology; polarized stimuli; 50 ms reliability caveats; mere-exposure/halo/confirmation-bias mechanism. | tier: paper (primary).
[^tract06]: https://www.sciencedirect.com/science/article/abs/pii/S1071581906000863 — Tractinsky et al. (2006): 500 ms ↔ 10 s stability; extreme judgments made faster. | tier: paper.
[^tuch-pt]: https://static.googleusercontent.com/media/research.google.com/en//pubs/archive/38315.pdf — Tuch et al. (2012): visual complexity + prototypicality at 17 ms; low VC + high PT = highest appeal; complexity processed before prototypicality. | tier: paper.
[^reinecke13]: https://wildlab.cs.washington.edu/Publications_files/Reinecke_CHI2013.pdf — Reinecke et al. (2013, CHI): colorfulness + complexity + demographics ≈48% of appeal variance; complexity dominates colorfulness; "not universal." | tier: paper.
[^papachristos]: https://dl.ifip.org/db/conf/interact/interact2011-1/PapachristosA11.pdf — Papachristos & Avouris (2011): replicates aggregate consistency but within-rater appeal avg r=.521; flags Lindgaard's polarized sample. | tier: paper (negation).
[^thielsch]: http://www.thielsch.org/download/paper/Thielsch_Hirschfeld_2012.pdf — Thielsch & Hirschfeld (2012): within-subject designs overestimate first-impression correlations; between-subjects 50 ms "mediocre but robust." | tier: paper (negation).
[^nng-halo]: https://www.nngroup.com/articles/halo-effect/ — NN/g, Halo Effect: definition; Thorndike 1920; Lindgaard & Dudek 2002 (>50% task failure, high satisfaction); masking risk. | tier: org-docs (NN/g).
[^dion]: https://psycnet.apa.org/record/1973-09160-001 — Dion, Berscheid & Walster (1972) "What Is Beautiful Is Good." | tier: paper (primary).
[^robins]: https://www.makinggood.ac.nz/media/1276/robinsholmes_2008_aestheticsandcredibilityinwebdesign.pdf — Robins & Holmes (2008): identical content + higher aesthetics → more credible in 19/21 (90%); "amelioration effect," visceral-level, "within the first few seconds." | tier: paper (primary).
[^lind11]: https://dl.acm.org/doi/10.1145/1959022.1959023 — Lindgaard et al. (2011, ACM TOCHI): appeal/trust/perceived-usability largely driven by appeal; trust anchored to but not identical to appeal. | tier: paper (primary).
[^kwak]: https://aisel.aisnet.org/cgi/viewcontent.cgi?article=1917&context=jais — Kwak, Ramamurthy & Nazareth, "Beautiful Is Good and Good Is Reputable" (JAIS): aesthetics strongest system feature for quality/reputation via halo. | tier: paper (primary).
[^secure]: https://dl.acm.org/doi/fullHtml/10.1145/3533047 — Stojmenović et al., "What Is Beautiful Is Secure" (ACM TOPS 2022): appeal drove perceived security over real identity verification. | tier: paper (primary).
[^eagly]: https://psycnet.apa.org/doiLanding?doi=10.1037%2F0033-2909.110.1.109 — Eagly, Ashmore, Makhijani & Longo (1991), "…But…" meta-analysis: beauty-is-good moderate, content-specific, ~0 for integrity, weakened by individuating info. | tier: paper (primary meta-analysis, negation).
[^gajos14]: https://eecs.harvard.edu/~kgajos/papers/2014/reinecke14visual.shtml — Reinecke & Gajos (2014, CHI), "Quantifying Visual Preferences Around the World": 2.4M ratings / ~40k people; peak appeal varies by gender/education/country. | tier: paper (primary, boundary).
[^gajos18]: https://dl.acm.org/doi/10.1145/3173574.3173911 — Reinecke & Gajos (2018), "A Case for Design Localization": 80,901 sites / 44 countries; systematic cross-country design differences. | tier: paper (primary, boundary).
[^fogg-46]: https://credibility.stanford.edu/pdf/How_Do_People_Evaluate_a_Web_Site's_Credibility_v37.pdf — Fogg et al. report: ranked credibility factors with Design Look 46.1%; coding definition. | tier: org-docs (Stanford, primary).
[^fogg-acm]: https://dl.acm.org/doi/10.1145/997078.997097 — ACM record, Fogg et al. (2003) DUX study, 2,684 participants; "present in 46.1% of the comments." | tier: paper (primary).
[^wiki-scwc]: https://en.wikipedia.org/wiki/Stanford_Web_Credibility_Project — independent confirmation of the 10 guidelines, participant figures, 46.1%, per-category breakdown, Fogg quote. | tier: encyclopedia (corroborating).
[^bokardo]: http://bokardo.com/archives/web_credibility/ — Joshua Porter critique: don't over-read the 46% study as "visuals matter most"; "only just a start." | tier: blog (named-expert critique).
[^fogg-pi]: https://credibility.stanford.edu/pdf/p-iTheory_Fogg_Oct02.pdf — Fogg, Prominence-Interpretation Theory: Prominence × Interpretation; the prominence/interpretation factors; mapping of the four credibility types. | tier: org-docs (Stanford, primary).
[^nng-pi]: https://www.nngroup.com/articles/prominence-interpretation-theory/ — NN/g restatement of PIT. | tier: org-docs (NN/g).
[^fogg-elements]: https://credibility.stanford.edu/pdf/p80-fogg.pdf — Fogg & Tseng, "The Elements of Computer Credibility" (CHI 1999): presumed/reputed/surface/earned taxonomy. | tier: paper (primary).
[^stanford-guidelines]: https://credibility.stanford.edu/guidelines/ — the verbatim 10 Stanford Guidelines for Web Credibility + per-guideline commentary (Fogg, 2002). | tier: org-docs (Stanford, primary).
[^fogg-chi01]: https://dl.acm.org/doi/10.1145/365024.365037 — Fogg et al. (CHI 2001): 7 factors; +real-world-feel/ease/expertise/trustworthiness/tailoring, −commercial implications/amateurism. | tier: paper (primary).
[^nng-trust]: https://www.nngroup.com/articles/trustworthy-design/ — NN/g (2016): 4 credibility factors; durability-of-guidelines argument with cross-cultural re-test. | tier: org-docs (NN/g).
[^nng-comm]: https://www.nngroup.com/articles/communicating-trustworthiness/ — Nielsen's original 4 trust factors; brittle-content and link-out examples; Hierarchy-of-Trust lineage. | tier: org-docs (NN/g).
[^baymard-seal]: https://baymard.com/blog/site-seal-trust — Baymard site-seal survey: brand recognition on a seal > what it technically certifies. | tier: org-docs (e-commerce research).
[^baymard-sec]: https://baymard.com/blog/perceived-security-of-payment-form — Baymard: users judge payment-form security *visually*; visual encapsulation + recognized marks reduce anxiety. | tier: org-docs (research).
[^pixxen]: https://pixxen.com/blog/website-trust-signals/ — trust-signal placement near friction; fake/aspirational badges backfire. | tier: blog (corroborating consensus).
[^makinggood]: https://www.makinggood.ac.nz/methods/web-credibility/ — WCP summary + the "too slick" backlash (Fogg 2002) and Alberts & van der Geest (2011) replication. | tier: blog (synthesis citing primaries).
[^aaker]: https://gsb-courses.stanford.edu/building-innovative-brands/wp-content/uploads/sites/25/2022/04/dimensions_of_brand_personality.pdf — Aaker (1997), full text: brand-personality definition; the five dimensions (Sincerity/Excitement/Competence/Sophistication/Ruggedness). | tier: paper (primary).
[^biaf]: https://pmc.ncbi.nlm.nih.gov/articles/PMC3882007/ — Kervyn, Fiske & Malone (2012), Brands as Intentional Agents Framework: warmth/competence quadrants + emotions; both predict purchase/loyalty; the .org/.com warmth-competence manipulation. | tier: paper (primary).
[^labrecque]: https://link.springer.com/article/10.1007/s11747-010-0245-y — Labrecque & Milne (2012), "Exciting red and competent blue": hue → personality; saturation/value amplify; red-excitement / blue-competence. | tier: paper (primary).
[^bottomley]: https://journals.sagepub.com/doi/10.1177/1470593106061263 — Bottomley & Doyle (2006): color-brand *congruity* (functional vs sensory-social colors); appropriateness > color in isolation. | tier: paper (primary).
[^elliot]: https://www.annualreviews.org/content/journals/10.1146/annurev-psych-010213-115035 — Elliot & Maier (2014), Annual Review of Psychology: authoritative color-psychology review; "nascent stage"; Color-in-Context (context-specific, even opposite meanings). | tier: paper (primary review, negation anchor).
[^jonauskaite]: https://journals.sagepub.com/doi/10.1177/0956797620948810 — Jonauskaite et al. (2020, Psych Science): 4,598 people / 30 nations; color–emotion associations both universal (r≈.88) and culture-specific. | tier: paper (primary).
[^colorarchive]: https://colorarchive.org/guides/color-psychology-branding/ — synthesis: color alone ≈5–10% of emotional variance; "+80% recognition" untraceable; consistency (not hue) drives recognition. | tier: blog (synthesis of research; magnitudes flagged unreliable).
[^singh]: https://newion.uwinnipeg.ca/~ssingh5/x/color.pdf — Singh (2006), "Impact of color on marketing": source of the "90 sec / 62–90% from color" claim AND its own admission of "absence of conclusive scientific results." | tier: paper (the headline stat is a soft claim — flag).
[^shaikh]: https://soma.sbcc.edu/users/russotti/113/personality_Shaikh.pdf — Shaikh, Chaparro & Fox (2006), "Perception of Fonts": fonts get consistent personalities by family (serif=stable/mature/formal; script=casual). | tier: paper.
[^fox]: https://journals.sagepub.com/doi/10.1177/154193120705100508 — Fox, Shaikh & Chaparro (2007): incongruent typeface → more negative document + author/ethos perception; cites Brumberger 2003 typeface persona. | tier: paper.
[^lund]: https://mjtsai.com/blog/2013/01/23/the-serif-readability-myth/ — summary of Ole Lund's review of 28 legibility studies (1896–1997): serif-readability is a myth; legibility ≠ readability ≠ reading speed. | tier: blog (summarizing a dissertation, negation).
[^hosang]: https://www.visible-language.org/Issue-59-3/VL-59-3-Ho-Sang-Petraca-8855.pdf — Ho Sang & Petraca review of 42 peer-reviewed studies: no significant reading-performance difference from serifs alone; familiarity dominant. | tier: paper (review, negation).
[^barneta]: https://psychology.unl.edu/can-lab/pubs/BarNetaPsychSci.pdf — Bar & Neta (2006), "Humans Prefer Curved Visual Objects": sharp contour → threat → negative bias (+ amygdala follow-up). | tier: paper (primary).
[^galli]: https://ira.lib.polyu.edu.hk/handle/10397/64985 — Galli & Chattopadhyay (2016, J. Consumer Research): circular logos → "softness," angular → "hardness," transferring to product attribute judgments. | tier: paper (primary).
[^frontiers-shape]: https://www.frontiersin.org/journals/psychology/articles/10.3389/fpsyg.2021.615647/full — Shape–Trait Consistency: angular → competence, rounded → warmth; power-state moderates preference for angular. | tier: paper (primary).
[^getty]: https://newsroom.gettyimages.com/en/getty-images/nearly-90-of-consumers-want-transparency-on-ai-images-finds-getty-images-report — Getty VisualGPS (30k/25 countries): 98% say authentic imagery pivotal to trust; want AI disclosure; less favorable to AI people/products. | tier: org-docs (large industry survey).
[^cornell]: https://news.cornell.edu/stories/2019/01/seeing-believing-depends-photo-quality-study-says — Cornell Tech (2019): high-quality user photos beat stock on trust ("too good to be true"); higher sales. | tier: org-docs (university research).
[^advergize]: https://www.advergize.com/color-psychology-marketing/ — appropriateness > universal rules; bank-red-anxiety and toy-brand-black mismatch examples (sources Bottomley & Doyle). | tier: blog (practitioner, sources primaries).
[^adessi]: https://adessi.io/consistency-in-visual-language/ — brand-consistency elements; "resets the recognition clock"; repeats the untraceable +80%/+33% stats (flagged). | tier: blog (direction only; magnitudes unreliable).
[^cohesion]: https://link.springer.com/article/10.1057/s41262-025-00410-2 — *J. Brand Management* (2025): "design coherence" (consistent logos/typefaces/color) as a precursor to brand-halo and loyalty. | tier: paper (primary).
[^koji]: https://www.koji.so/docs/5-second-test-guide — 5-second test mechanics, question taxonomy, timing / no-priming rules. | tier: blog (vendor guide).
[^maze-5s]: https://maze.co/collections/user-research/five-second-test/ — how to run a five-second test; impression questions. | tier: blog (vendor, standard method).
[^cleverx]: https://cleverx.com/blog/five-second-testing-how-to-measure-first-impressions-in-ux-research/ — 4 core questions, ≤5 open-text, what it does/doesn't measure; "does not assess usability/tasks." | tier: blog (vendor).
[^lyssna]: https://www.lyssna.com/guides/five-second-testing/ — step-by-step; enforce the timer; ask the most-important question first. | tier: blog (vendor).
[^nng-squint]: https://www.nngroup.com/videos/squint-test/ — NN/g: squint test to find attention-grabbing elements & assess hierarchy. | tier: org-docs (NN/g).
[^lukew]: https://www.lukew.com/ff/entry.asp?2013= — Wroblewski coining the squint test; "which elements *should* stand out" strategic framing. | tier: org-docs (origin source).
[^polypane]: https://polypane.app/blog/debug-your-visual-hierarchy-with-the-squint-test/ — squint inspection checklist (primary focus, relationships, foreground/background). | tier: blog.
[^kuric]: https://www.tandfonline.com/doi/full/10.1080/0144929X.2023.2272747 — Kuric et al. (2023): "pseudoscience" critique of fixed-5s testing; cognitive ability × visual complexity affect first impressions. | tier: paper (peer-reviewed, negation).
[^gronier]: https://uxpajournal.org/measuring-testing-validity-5-second-test/ — Gronier (2016, JUS): 5-second-test validity; null differences vs full usability test; static-page-only; not for usability problems. | tier: paper (peer-reviewed, negation).
[^nng-cards]: https://www.nngroup.com/articles/microsoft-desirability-toolkit/ — NN/g: 118 cards, ≥40% negative, select → narrow-to-5, controlled vocabulary; triangulate. | tier: org-docs (NN/g).
[^uxfirm]: https://www.uxfirm.com/microsofts-product-reaction-cards-unlock-user-satisfaction — 60/40 positivity rationale; triangulate with SUS. | tier: blog (corroborating).
[^bluehair]: https://www.bluehair.co/wp-content/uploads/2021/03/assessment-of-existing-tools-for-the-measurement-of-emotions.pdf — descriptions of SAM, PrEmo, Geneva Emotion Wheel; comparison criteria. | tier: blog (report).
[^gew]: https://www.unige.ch/cisa/gew — Geneva Emotion Wheel: 20 emotion families, valence × control axes, 5 intensities. | tier: org-docs (instrument authors).
[^nng-journey]: https://www.nngroup.com/articles/journey-mapping-101/ — NN/g journey-map anatomy; emotion as a single up/down line. | tier: org-docs (NN/g).
[^uxcrush]: https://uxcrush.com/user-journey-map — emotion-curve dips = highest-priority pain points; "confused by step 3" is actionable. | tier: blog.
[^measuringu-des]: https://measuringu.com/microsoft-desirability/ — Sauro/MeasuringU: how to quantify reaction cards; **no benchmarks; no evidence it measures desirability.** | tier: org-docs (MeasuringU, negation).
[^nng-severity]: https://www.nngroup.com/articles/how-to-rate-the-severity-of-usability-problems/ — Nielsen 0–4 scale; frequency/impact/persistence; combine into one rating; single-rater unreliable, mean-of-4 close to true. | tier: org-docs (NN/g, primary).
[^betterqa]: https://betterqa.co/bug-priority-vs-severity-levels/ — severity (technical, set by evaluator) vs priority (business, set by PM); non-linear. | tier: blog.
[^contextqa]: https://contextqa.com/blog/difference-between-severity-and-priority/ — ISTQB defs; 2×2 quadrant; misspelled-logo = low severity / high priority. | tier: blog.
[^altersquare]: https://altersquare.io/how-to-prioritize-bugs-during-product-iteration/ — value × effort matrix; quick-wins; high-severity/low-priority cases. | tier: blog.
[^hertzum]: https://mortenhertzum.dk/publ/HFES1998.pdf — Hertzum & Jacobsen, the "evaluator effect": large inter-rater disagreement in usability evaluation; severity ratings unreliable from one rater. | tier: paper (negation).
[^saleshub]: https://www.saleshubhq.com/strong-website-trust-signals-that-drive-more-leads/ — trust-signal catalog: real photos, contact info, specific social proof, security cues at the form, error-free copy. | tier: blog (multi-source consensus).
[^darkpattern]: https://link.springer.com/article/10.1007/s13520-026-00254-2 — dark-pattern harm pathway: manipulation → emotional distress → trust/autonomy erosion. | tier: paper (peer-reviewed).
[^uiux-toolkit]: https://heurilens.com/heuristic-analysis — multi-domain audit composition; "merge findings, remove duplicates, average severity"; priority tiers. | tier: blog (practitioner meta-method).
[^deque]: https://www.deque.com/blog/supporting-the-design-phase-with-accessibility-heuristics-evaluations/ — accessibility-as-heuristic-evaluation; 3–5 evaluators; scope → debrief → prioritize. | tier: org-docs (Deque).
[^graticle]: https://graticle.com/blog/why-people-trust-bad-websites-more-than-slick-ones/ — slick design can bury honesty signals; "match polish to depth of proof." | tier: blog (corroborates the decoupling rule).
[^tuch-chb]: https://www.sciencedirect.com/science/article/abs/pii/S0747563212000908 — Tuch et al. (2012, CHB), "Is beautiful really usable?": aesthetics did NOT raise perceived usability; poor usability lowered post-use aesthetics ("usable is beautiful"). | tier: paper (primary, negation).
[^grishin]: https://uxpajournal.org/wp-content/uploads/sites/7/pdf/JUS_Grishin_Feb2019.pdf — Grishin & Sauro (2019, JUS): failed to replicate aesthetics→perceived-usability; SUS driven by actual usability. | tier: paper (primary, negation).
[^nielsen-crit]: https://jakobnielsenphd.substack.com/p/design-crit — Nielsen on crit anti-patterns: bikeshedding, "Make-it-Pop" vague feedback, HiPPO + silent-critique fix. | tier: org-docs (Nielsen).
[^uxtigers]: https://www.uxtigers.com/post/design-crit — diagnose-don't-prescribe; ban "I like / I don't like"; the "this may cause trouble because…" stem. | tier: blog (named expert).
[^prototypr]: https://blog.prototypr.io/6-anti-behaviours-in-design-critique-and-how-to-deal-with-them-9e1c3f4d2f1a — HiPPO, compliment-sandwich / "Feedback Pancakes," empty praise; ask "what would make it better." | tier: blog.
