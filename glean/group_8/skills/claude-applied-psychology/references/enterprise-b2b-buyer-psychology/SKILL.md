---
name: enterprise-b2b-buyer-psychology
description: >-
  How enterprises (especially risk-averse big banks/FSI) buy B2B technology: the psychology and
  structure of the enterprise buying decision. Applied B2B-buying spoke of applied-psychology that
  APPLIES sibling-owned mechanisms and cites them. TRIGGER: who really decides / buying committee /
  buying center / decision-making unit (DMU); a deal stalled in committee or gone "no decision";
  Gartner buying group & buyer "jobs"; champion/mobilizer ID & coaching; perceived purchase/vendor
  risk, "nobody got fired for buying IBM", vendor viability; switching-cost & vendor-lock-in fear;
  FSI/bank procurement & vendor-risk regulation. SKIP (mechanism owners): loss-aversion/status-quo/
  sunk-cost → behavioral-decision-making; trust → trust-and-psychological-safety; ELM/Cialdini →
  persuasion-and-influence-psychology; adoption/SDT/Fogg → behavior-change-psychology; deal-execution/
  MEDDIC → tam-operations; FSI regulatory map → fsi-banking-regulatory-context; customer message →
  content-and-marketing-writing.
version: 1.1.0
updated: 2026-06-18
category: applied-psychology
whenToUse: >-
  Load when reasoning about how an enterprise customer (especially a regulated FSI / large bank) makes
  a B2B technology purchase or renewal: who is in the buying group and what role each plays, why a deal
  is stuck or heading to "no decision", how to find/develop/arm an internal champion, how perceived
  risk and switching-cost/lock-in fear shape a competitive displacement, and how heavyweight FSI
  procurement and vendor-risk regulation gate the decision. This is the applied B2B-buying layer; for
  the underlying psychological mechanisms, route to the sibling references named in SKIP.
keywords:
  - B2B buying committee
  - buying center
  - decision-making unit (DMU)
  - organizational buying behavior
  - Webster and Wind
  - Sheth model
  - BUYGRID buyclass buyphase
  - Gartner buying group 6-10 stakeholders
  - buyer enablement six jobs
  - no decision dysfunctional consensus
  - Challenger Customer mobilizer talker blocker
  - JOLT effect customer indecision
  - perceived risk risk reduction
  - nobody got fired for buying IBM
  - vendor viability vendor risk criterion
  - switching costs Burnham procedural financial relational
  - vendor lock-in fear competitive displacement
  - champion mobilizer MEDDIC MEDDPICC
  - Miller Heiman buying influences coach
  - single-threaded risk multi-threading
  - FSI big bank procurement vendor risk
  - third-party risk DORA SR 11-7 concentration risk exit plan
tags:
  - applied-psychology
  - b2b-buying
  - enterprise-sales
  - buying-committee
  - fsi
  - tam
  - vendor-risk
---

# Enterprise B2B Buyer Psychology & Buying Committees

> Provenance: spoke of the **applied-psychology** hub. Created 2026-06-18 via `/dr` deep research (4 parallel research agents, 40+ independent sources). This is the **B2B-buying synthesis** layer: it *applies* psychological mechanisms owned by sibling hub references and cites them — it does not re-teach the base mechanisms.
>
> `verified-as-of: 2026-06-18` for all volatile claims (vendor-research statistics, the regulatory landscape, the 2026 model-risk revision). Re-verify before relying on a dated number.

This reference explains **how large, risk-averse enterprises, especially big banks and other financial-services institutions (FSI), actually buy B2B technology**: who decides, why deals stall, and what shapes the decision psychologically and structurally. Written for a MongoDB TAM working a flagship FSI account, but general to any complex enterprise purchase.

## What this spoke owns vs. what it cites

This is the **applied B2B-buying layer**. The underlying mechanisms live in sibling references, load them for the mechanism, load this for the B2B application:

- Generic **loss aversion, status-quo bias, sunk-cost, framing, prospect theory, nudge** → `behavioral-decision-making` (`../behavioral-decision-making/SKILL.md`). This spoke shows how they manifest as *vendor switching inertia*, not how they work.
- **Trust formation & repair, ABI model, Trust Equation, psychological safety** → `trust-and-psychological-safety` (`../trust-and-psychological-safety/SKILL.md`). Vendor viability and reference-checking are trust signals; the trust mechanism is owned there.
- **Persuasion: ELM/HSM, reactance, Cialdini, inoculation, norms** → `persuasion-and-influence-psychology` (`../persuasion-and-influence-psychology/SKILL.md`). Arming a champion is applied persuasion; the routes/principles are owned there.
- **Adoption & motivation: SDT, Fogg B=MAP, habits, stages of change** → `behavior-change-psychology` (`../behavior-change-psychology/SKILL.md`). Post-purchase adoption is owned there.
- Drafting the actual customer message → `content-and-marketing-writing`; TAM deliverables / QBR / churn frameworks → `tam-operations`.
- **The vendor-risk/TPRM *mechanics* themselves** (TPRM lifecycle, the security-questionnaire gauntlet, DORA/OCC/FFIEC the-rule-says detail, GRC tooling) → `enterprise-vendor-management-and-tprm`; **FSI seller fluency & the regulatory map** (how bank segments differ as buyers, Basel/SOX/GLBA/SR 11-7 pointers) → `fsi-banking-regulatory-context`. This spoke owns the *psychology* of why those gates make a bank a risk-averse buyer and how to clear them via a champion; it does not own the regulatory mechanics. (Both sibling skills may be registry-resolved before they are installed locally.)
- **When NOT to over-apply this:** straight rebuys, single-owner SMB deals, and routine renewals compress or skip the buying center, do not force DMU mapping, multi-threading, or a consensus toolkit onto a 2-person decision. The machinery here assumes a large, multi-stakeholder, high-perceived-risk *new-task* purchase (see Section 1's BUYGRID). At **renewal** you are the incumbent: the Section 3–4 inertia and "nobody got fired" mechanics now protect your install base, and the job inverts to defending against a challenger.

## Contents

This reference is split for token economy. Read the section you need:

1. **Core models & the buying center (DMU)** — Webster & Wind, Sheth, BUYGRID, the six roles. *(below)*
2. **The modern Gartner buying group & "no decision"** — 6–10+ stakeholders, the six buyer jobs, buyer enablement, the no-decision/dysfunctional-consensus outcome. *(below)*
3. **The three layers: rational + political + emotional** — perceived risk, "nobody got fired for buying IBM", incumbent advantage, vendor viability as a criterion. *(below)*
4. **Switching-cost psychology & lock-in fear** — Burnham's taxonomy; how lock-in fear shapes competitive displacement. *(below)*
5. **Champions & mobilizers** — identify, develop, arm; single-threaded risk and multi-threading. *(below)*
6. **FSI / big-bank buying behavior** — the deep regulatory + procurement detail, and vendor implications → **`references/fsi-bank-buying.md`** (load for the FSI deep-dive).

A one-screen field playbook for a TAM is at the end.

---

## 1. Core models & the buying center (DMU)

Three foundational academic models still anchor how B2B buying is taught. They are old but their **core dimensions remain validated**; some original *assumptions* did not hold up (flagged below).[^ww][^sheth][^buygrid50]

### Webster & Wind (1972)
"A General Model for Understanding Organizational Buying Behavior" (*Journal of Marketing* 36(2):12–19) defines organizational buying as the decision process by which organizations establish a need and identify, evaluate, and choose among suppliers.[^ww] Two enduring contributions:

- **Four nested sets of variables** shaping the decision: **environmental** (economic, technological, legal, competitive), **organizational** (objectives, structure, purchasing policy, centralization, reward systems), **interpersonal/social** (authority, status, persuasion among participants), and **individual** (the decision-maker's own goals, perceived risk, personal needs).[^ww][^mbaknol] Each variable has a **task** dimension (tied to the buying problem, e.g., lowest price) and a **nontask** dimension (personal values, career needs), most have both, with one predominant.[^ww]
- **The "buying center" construct** — only a subset of organizational actors is involved in any given purchase, motivated by "a complex interaction of individual and organizational goals."[^ww] This is the origin of the buying-center / decision-making-unit (DMU) idea.[^dmu]

Their dictum to remember: **"in the final analysis, all organizational buying behavior is individual behavior."**[^mbaknol]

### Sheth (1973)
"A Model of Industrial Buyer Behavior" (*Journal of Marketing* 37(4):50–56) focuses on the **group process** where Webster & Wind focused on context.[^sheth][^uky] Three aspects:

1. **The psychological world of the decision-makers** — the model centers on the *expectations* of participants (classically the purchasing agent, the engineer, the user) about suppliers. Five processes make those expectations *differ*: background, information sources, active search, perceptual distortion, and satisfaction with past purchases. Specialists are rewarded for specialized viewpoints, so buying motives diverge systematically (cost vs. quality vs. usability).[^sheth][^jagsheth]
2. **Conditions for joint vs. autonomous decisions** — six factors: three **product-specific** (perceived risk, type of purchase, time pressure) and three **company-specific** (company orientation, size, centralization). High perceived risk, capital/novel purchases, low time pressure, large size, and low centralization push toward **joint** (multi-party) decisions; the inverse delegates to one person.[^sheth][^jagsheth]
3. **Conflict resolution in joint decisions** — conflict is the normal consequence of deciding jointly across different goals. Sheth names four modes: **problem-solving** and **persuasion** (rational, conflict over criteria, resolved by information or appeal to objectives → better decisions) and **bargaining** and **politicking** (conflict over fundamental goals → irrational, possibly detrimental supplier choices).[^jagsheth] Sheth also allows that some industrial decisions are driven by **ad hoc situational factors**, not any systematic process.[^jagsheth]

> **TAM use:** when a deal is being decided by problem-solving/persuasion, *more/better information wins*. When it has degraded to bargaining/politicking, you have a **goal-conflict** problem in the buying center that information alone will not fix — you need a mobilizer to broker it (Section 5).

### The six buying-center / DMU roles
Webster & Wind's original paper lists **five** roles verbatim; the widely taught **sixth (initiator)** is the Kotler/textbook addition.[^ww][^dmu] Some texts add a seventh, **approver**.[^yourarticle]

| Role | Definition | Typical in a B2B tech / database purchase |
| --- | --- | --- |
| **Initiator** | First sees the need and starts the process (Kotler addition).[^dmu] | An SRE/eng lead who flags the current datastore can't scale. |
| **User** | Will actually use the product; shapes specs.[^ww] | Developers, DBAs, analysts who run the workload daily. |
| **Influencer** | Shapes evaluation criteria; can veto on technical grounds.[^ww] | Staff/principal architect; security engineer (SOC 2, latency SLOs). |
| **Decider** | Authority to choose among alternatives.[^ww] | VP Engineering / CTO. |
| **Buyer** | Formal authority to contract & negotiate.[^ww] | Procurement / vendor management. |
| **Gatekeeper** | Controls the flow of information (and access) into the center.[^ww] | Exec assistant; an IT/security reviewer who decides which vendors reach the deciders, "many promising conversations die here."[^dmu] |
| *(Approver)* | Authorizes the deciders'/buyers' actions (7-role variant).[^yourarticle] | CFO sign-off above a budget threshold. |

One person can hold several roles; several people can share one.[^ww] **Mapping the DMU early, naming who plays each role, is the single most useful structural move on an enterprise deal.**

### BUYGRID & why "new task" buys are the biggest
Robinson, Faris & Wind (1967) crossed **3 buyclasses × 8 buyphases**.[^buygrid][^buygrid50]

- **Buyclasses:** **New task** (first-time, high uncertainty, max information & alternatives sought), **modified rebuy** (changing specs/price/supplier on something already bought, smaller, faster center), **straight rebuy** (routine reorder from an approved supplier, often just a buyer / automated).[^buygrid]
- **Buyphases:** problem recognition → general need description → product specification → supplier search → proposal solicitation → supplier selection → order-routine specification → performance review.[^buygrid]
- **Why new task = largest buying center:** the full eight phases run only for a new task; rebuys compress or skip stages. Higher novelty/uncertainty/perceived risk draws in more participants and more functions, which is exactly Sheth's prediction (high perceived risk + capital/novel → joint decisions).[^buygrid][^sheth] A first-ever NoSQL or cloud-database adoption at a bank is a *new task* and will recruit the largest, most cross-functional committee.

**Disconfirming note (don't over-claim these models).** The 50-year and 25/30-year retrospectives find the *dimensions* (buying situation, buying process, buying center) survived, but Wind & Thomas concede some original *assumptions* were not supported, and Ferguson (1979) found BUYGRID didn't fit services rebuys; a 2022 *Industrial Marketing Management* review judged the buying-center concept had "partially lost its vigor."[^buygrid50][^buyclassretro] Treat them as durable scaffolding, not laws.

### B2B vs. consumer buying (the honest caveat)
Organizational buying differs from consumer buying in **structure** (derived, more volatile/inelastic demand; fewer, larger, concentrated buyers), **who decides** (a multi-person buying center), and **how** (more formal, professional, rational; RFPs, written POs, longer cycles).[^kotler][^lumen] **But it remains fundamentally human** — businesspeople buy solutions to *two* problems at once: the organization's economic/strategic problem *and* their own need for achievement, reward, and safety. So B2B decisions are **both rational and emotional**, exactly as Webster & Wind's task/nontask split and Sheth's "psychological world" imply.[^kotler] This dual nature is the bridge to Section 3.

---

## 2. The modern Gartner buying group & "no decision"

> **Provenance flag:** CEB (Corporate Executive Board) was **acquired by Gartner in 2017**. "CEB research," "Challenger research," and "Gartner research" are often the *same lineage* under different brands and years; Challenger Inc. is now a separate training company holding the IP. Statistics get cross-attributed and the headline numbers have **drifted upward over time**. Treat brand + year carefully.[^challengerinc]

### The buying group is large and self-directed
Gartner's own published language (current as of 2023–2024): **"The typical buying group for a complex B2B solution now involves six to 10 decision makers, each armed with four or five pieces of information they've gathered independently."**[^gartnerjourney][^advweek] `verified-as-of: 2026-06-18`.

- The **"~11" figure** is the *same lineage, enterprise/complex end* — Gartner elsewhere states **"five to 11 stakeholders, who represent an average of five distinct business functions"**, and Challenger Inc. cites ~11 "and in complex deals, as many as 20."[^challengerinc][^gartnerreport] Don't present "11" and "6–10" as the same measurement.
- **Historical drift (cite it for credibility):** CEB's figure was **5.4 (2015) → ~6.8 (2017) → 10.2 (2018)**, now "6–10" with enterprise up to ~20.[^geniusdrive] Buying-group size grew materially over the decade.
- Companion stat: **77% of B2B buyers say their latest purchase was very complex or difficult.**[^advweek]
- **Disconfirming nuance, use a *range by deal size*, not one number.** The single headline figure is criticized for hiding variance: SMB ≈ 2–4 people, mid-market ≈ 4–7, enterprise ≈ 8–14, strategic platform deals 15–25.[^abmatic] Likewise the famous "57% of the journey complete before contacting sales" is a *complex-deal average*, not a universal law, even CEB's own people walked it back.[^57myth]

### Buyers spend most of their time *away* from vendors
Gartner's 750-buyer study breaks the journey time as **27% researching independently online, 22% meeting the internal buying group, 18% researching offline, 17% meeting *all* potential suppliers combined, 16% other.**[^leveragepoint][^b2bmktg] Implications:

- ~**45% is independent research** and ~**two-thirds is spent away from suppliers entirely.**[^leveragepoint]
- With multiple vendors competing for the 17%, **any one sales rep gets only ~5–6%** of the buyer's time.[^sechub]
- **~75% of B2B buyers say they prefer a "rep-free" sales experience** (Gartner; lower figures of 67–72% appear in other waves).[^gartnerjourney] `verified-as-of: 2026-06-18`.

### Buyer enablement and the six "jobs"
Gartner reframes buying not as a funnel but as a set of **"jobs" buyers must complete, and they loop, revisiting each at least once**, not in a linear order.[^gartnerjourney][^gartnerbe] The six jobs (Gartner's buyer-voice framing):

1. **Problem identification** — "We need to do something."
2. **Solution exploration** — "What's out there to solve our problem?"
3. **Requirements building** — "What exactly do we need the purchase to do?"
4. **Supplier selection** — "Does this do what we want it to do?"
5. **Validation** — "We think we know the right answer, but we need to be sure."
6. **Consensus creation** — "We need to get everyone on board."

*(Gartner sometimes presents the first four as the "core jobs" with Validation and Consensus Creation as cross-cutting jobs that recur, the six terms are stable, the grouping is presented inconsistently even within Gartner.)*[^gartnerjourney]

**Buyer enablement** = giving buyers the information and tools to complete those jobs. Gartner's load-bearing finding: **"information, not interactions with salespeople,"** is what most helps buyers progress, so good content is **shareable, channel-agnostic, consistent, and prescriptive** (mapped to a specific job).[^gartnerbe][^gartnerreport] Buyers are **1.8× more likely** to close a high-quality, low-regret deal when they use supplier digital tools *with* a rep, and **2.8× more likely** when they perceive high information consistency between website and rep; self-service-only buyers are **1.65× more likely to regret** the purchase.[^gartnerreport] `verified-as-of: 2026-06-18`.

### The "no decision" / dysfunctional-consensus outcome
This is the single most-misattributed area. **Two distinct research bodies, two distinct numbers, keep them separate:**

**(a) "40–60% of *lost* deals end in no decision" = *The JOLT Effect* (Dixon & McKenna, 2022), NOT *The Challenger Customer*.**[^jolt][^challengerjolt] From a study of **2.5 million recorded sales conversations**: 40–60% of deals are lost to customers who express intent to buy but **fail to act**. JOLT attributes this primarily to **individual buyer indecision / "Fear Of Messing Up" (FOMU)** — and partly *corrects* the older "just beat the status quo" framing. The JOLT playbook: **J**udge indecision, **O**ffer a recommendation, **L**imit the exploration, **T**ake risk off the table.[^jolt][^challengerjolt] `verified-as-of: 2026-06-18`.

**(b) Buying-group-dysfunction stats = *The Challenger Customer* (CEB, 2015), different numbers.**[^challengercust][^challengerinc] In the average B2B group of **5.4 stakeholders**, the likelihood of *any* purchase drops to **~30%**; going from 1→2 decision-makers drops purchase likelihood from **81% to 55%**.[^challengercust][^consensus] **46% of customers** say agreeing as a group is highly difficult, and **conflict peaks ~37% of the way** through the journey.[^challengercust] Without alignment, a buying group either **settles for a "good enough" lowest-common-denominator** or, *more commonly* — **does nothing at all.**[^challengercust] Challenger Inc.'s own figure: **~38% of purchase attempts end in no decision.**[^challengerinc] *"A collection of yeses is not the same as a collective yes."*[^challengercust]

> **The single most common error in this topic** is attributing "40–60% no decision" to The Challenger Customer. It is JOLT (2022). Also watch the **denominator**: 40–60% is the share of *lost* deals; "no decision" as a share of *all* outcomes is more like 20–35%.[^umbrex]

> **TAM use:** "No decision" is your real competitor more often than a rival vendor. Both bodies of research point to the same job: **de-risk the individual buyer's fear of being wrong (JOLT) *and* help the group actually reach a collective yes (Challenger)** — which is what champions/mobilizers do (Section 5).

---

## 3. The three layers: rational + political + emotional

Enterprise buying runs on three simultaneous layers. The rational layer (features, price, fit) is the visible one; the **political** layer (whose budget, whose turf, whose career) and the **emotional** layer (fear of being wrong) decide most stalled deals. The base mechanisms here, loss aversion, status-quo bias, sunk cost, are owned by `behavioral-decision-making`; this section is their **B2B manifestation**.

### Perceived risk (the emotional engine)
Perceived risk originates with Bauer (1960): buying is risk-taking because outcomes are uncertain and some are unpleasant, and it is **subjective/perceived**, not objective, risk that drives behavior.[^bauer][^mitchell] It is modeled as **uncertainty × consequences** (Cunningham 1967).[^perceivedrisk] The classic risk-type taxonomy is **Jacoby & Kaplan (1972): functional/performance, physical, financial, social, psychological** risk, with **time/convenience** commonly added.[^jacoby][^perceivedrisk]

**In B2B, add political/career risk to the individual buyer.** Organizational purchase risk = **(a) economic/performance risk to the organization + (b) psychosocial/personal risk to the buyer** — fear of looking bad to peers or hurting promotion prospects after a poor choice.[^orgbrisk] Empirically, B2B buyers experience **product-performance, personal-psychological, and personal-financial** risk.[^b2brisk] And **risk perception rises with purchase importance and uncertainty** — experiments with real supply managers show that when a category is important/difficult, buyers prefer the **more certain** supplier *even at a lower expected payoff.*[^supplierselect] This career-risk logic is formalized in **Anderson & Chambers' (1985) Reward/Measurement model**: buyers favor vendors that raise the probability they're *seen* as high performers.[^rewardmodel]

*(Disconfirming: the multiplicative uncertainty×consequences model is critiqued as psychometrically "suspect", multiplying ordinal scales, and there's no consensual definition of perceived risk across disciplines.[^riskcritique] The *construct* holds up; its *measurement* is contested.)*

### Risk-reduction strategies buyers use
Buyers deploy identifiable **risk relievers**: stay with the **incumbent** (the most common), buy from **large, established vendors**, demand **references + independent verification**, run **pilots/POCs** (a controlled experiment with success metrics defined *before* the demo so the vendor can't redefine success after), and **multi-source** to avoid single-point failure.[^riskrelievers][^vendorref][^poc][^dualsource]

**"Nobody ever got fired for buying IBM."** The phrase appears in print by **1978** and spread from **1984**; it was never an IBM ad (IBM alluded to it, "we sell a good night's sleep").[^ibm][^ipglossary] Its meaning is **career-risk minimization and herd safety**: bureaucracies are risk-averse and want the solution their peers bought, because being "in the herd" obviates criticism, forcing new entrants to overcome the incumbent's standing reputation.[^ipglossary] It mutated to "buying Microsoft" (1990s) and today "buying the market leader."[^ibm] `verified-as-of: 2026-06-18` (phrase usage, not the underlying mechanism). **For a challenger vendor (e.g., MongoDB displacing a legacy RDBMS at a bank), this heuristic IS the competitor** — you must give the buyer enough cover (peer-bank references, proof, viability) that choosing you is the *safe* career move, not the bold one.

### Incumbent advantage as switching inertia
The base mechanisms (status-quo bias, Samuelson & Zeckhauser 1988; endowment effect / loss aversion — Kahneman, Knetsch & Thaler 1991; sunk cost) are owned by `behavioral-decision-making`.[^statusquo][^endowment] Their **B2B manifestation is documented switching inertia**: Polites & Karahanna ("Shackled to the Status Quo," *MIS Quarterly*) show incumbent-system habit, transition-cost rationalization, and sunk-cost commitment depress new-system adoption *above and beyond* rational perceptions; demand-side work formalizes a numeric **"inertia penalty"** — a challenger must be *more than marginally* better to win.[^shackled][^inertia] In sales terms, this shows up as **elongated evaluations, pilot fatigue, deferred decisions**, i.e., "No Decision," which Forrester put at **43% of lost IT-sales opportunities.**[^pedowitz][^forrester] Incumbents deliberately exploit it (reminding the customer of sunk evaluation cost, the "burn-victim" effect of a prior bad switch).[^incumbent]

### Vendor viability as an explicit buying criterion (the rational layer formalizing the fear)
"Will this vendor still exist, and still support this, in three years?" is a *scored* criterion, not a vibe. The recurring framework asks three independent questions: **financial viability**, **operational viability** (SLAs, support, security *today*), and **roadmap fit**: a vendor strong on two legs and weak on the third is "a timed risk."[^viability] The concrete fear: a distressed vendor may **discontinue products, be acquired, or shut down**, causing missed roadmap, support delays, pricing changes, or sudden service shutdown.[^vendordisappear][^vendorfin] Procurement/risk functions formalize it with **financial-health metrics** (cash runway, revenue trend, debt-to-equity, quick ratio, Altman Z-score, customer concentration), **deprecation-policy commitments**, **severity-weighted scorecards with go/no-go gates**, and **M&A-specific diligence** (12–24-month bridge commitments, escrow, migration assistance).[^vendorrisk1][^vendorrisk2] This matters most in fast-consolidating categories. *(The "67% of venture-backed SaaS fail before Series C" and "single customer >40% of revenue" figures are vendor/consultancy estimates, directional, not peer-reviewed.)*[^viability] `verified-as-of: 2026-06-18`.

---

## 4. Switching-cost psychology & lock-in fear

The authoritative taxonomy is **Burnham, Frels & Mahajan (2003), "Consumer Switching Costs: A Typology, Antecedents, and Consequences"** (*Journal of the Academy of Marketing Science* 31(2):109–126).[^burnham] Its headline empirical finding: **all three switching-cost types significantly drive intention to stay, explaining more variance than satisfaction does.** Three higher-order types, eight first-order facets:[^burnham][^mofokeng][^jsalt][^blut]

**1. Procedural switching costs** (lost time and effort):
- **Economic-risk costs** — accepting uncertainty about an unknown new provider's performance.
- **Evaluation costs** — time/effort to search and analyze alternatives.
- **Learning costs** — time/effort to acquire the skills to use the new product effectively.
- **Setup costs** — time/effort to initiate the relationship / configure for first use.

**2. Financial switching costs** (quantifiable resource loss):
- **Benefit-loss costs** — forfeited contractual/accumulated benefits (loyalty status, accrued value).
- **Monetary-loss costs** — one-time out-of-pocket switching outlays (not the new product's price); often equated with sunk costs.

**3. Relational switching costs** (psychological/emotional discomfort):
- **Personal-relationship-loss costs** — emotional loss from severing bonds with the people/account team.
- **Brand-relationship-loss costs** — emotional loss from breaking identification with the provider's brand.

### How this becomes lock-in fear in a competitive displacement
High switching costs are literally labeled **"being locked-in"**, producing **calculative commitment**; research finds customers perceive even *small procedural* switching costs as painful and **stay even when dissatisfied.**[^mofokeng] In B2B specifically, switching costs predict **actual** switching behavior, not just intentions.[^ha] The buyer contemplating a displacement feels it as **migration & retraining anxiety** and fear of disruption, countered with rollback plans, exit ramps, and phased deployment.[^pedowitz]

**Technology/data mapping (sound, but don't over-claim the dollar figures):**
- **Data gravity** — large datasets in proprietary formats are expensive and slow to move; in regulated industries migration adds compliance risk and downtime (≈ monetary-loss + economic-risk costs).[^datagravity]
- **Retraining** — migrating may require retraining whole teams, temporarily cutting productivity (≈ learning costs).[^datagravity]
- **Integration rework** — every proprietary integration point (queues, triggers, event rules) adds weeks/months to leave (≈ setup costs).[^datagravity]
- Practitioners reframe this as a measurable **"Lock-In Liability"** (engineering hours × loaded cost to migrate away). *Specific figures ($500K–$2M after 2–3 years; migration ≈ 2–4× first-year license savings) are illustrative practitioner/vendor estimates, TENTATIVE.*[^datagravity] `verified-as-of: 2026-06-18`.

### Disconfirming evidence (important for honesty and for displacement strategy)
- **Vendor lock-in fear is frequently *overstated*.** Multiple sources argue the probability an enterprise actually switches clouds/platforms is very low absent extreme economic/geopolitical pressure, and that *multi-cloud-to-avoid-lock-in* can **increase** cost and lock-in (one case: +35% opex, slower deployments); the discussion "happens less and less" as buyers prioritize speed over future portability.[^lockinoverstated][^fowlerlockin] So the *fear* is real and shapes rhetoric, but its actual influence on switching *behavior* is contested.
- **The three types are NOT interchangeable in B2B.** Blut/Evanschitzky et al. found **financial switching costs *insignificant*** for switching behavior and share-of-wallet, only **relational** switching costs affected all outcomes; a separate study found the same.[^blut][^ha][^jsalt] **Implication for displacing an incumbent:** the binding cost is usually **relational and procedural** (the account team relationship, retraining, integration rework), *not* the headline dollar figure. Win the relational/procedural argument, and note that in a *regulated FSI*, lock-in flips from a psychological cost to a **regulatory requirement** (concentration risk + tested exit plan, see Section 6).

---

## 5. Champions & mobilizers

> **Honesty flag.** MEDDIC/MEDDPICC and Miller Heiman Strategic Selling are **practitioner/commercial sales methodologies, not peer-reviewed science** — their evidence base is internal win/loss analysis and vendor-published research. The Challenger/Gartner Mobilizer line is closer to formal research but is also commercially published. Treat all as **expert heuristics**, and treat vendor close-rate statistics (Gong/Chorus, below) as **TENTATIVE** unless independently verified.[^prospeo][^spotlightbroken]

### What a champion is (and is not)
A **champion** is an internal advocate who has **power/influence and access** and who **actively sells on your behalf when you're not in the room.** The MEDDIC test requires three traits *simultaneously*: (1) power/political influence, (2) a **personal win** tied to the deal, (3) **active internal selling**; many add (4) **access to the economic buyer**.[^meddiccchampion][^prospeo] The distinguishing word is **"actively."** *"A champion is not someone who likes your product. A champion is someone who fights for it internally."*[^spotlightval]

| Role | Key difference |
| --- | --- |
| **Coach** | Gives you intel and wants you to win, but **lacks influence or won't spend political capital.** "A Coach tells you what's happening; a Champion makes things happen." A coach is a champion missing one of the three traits.[^accountmap] |
| **Champion** | Has power *and* acts. |
| **Economic Buyer** | Final budget authority + veto; one per deal. The champion is usually *not* the EB — part of the champion's job is getting you *access* to the EB.[^millerheiman] |
| **Sponsor** | Informal vocabulary (often an exec giving "air cover"); **not a formally defined role** in MEDDIC or Miller Heiman, don't treat it as one.[^prospeo] |

**Test a real champion vs. a "cheerleader": give them a task that costs social capital** — set up an EB meeting, share internal evaluation criteria, present your business case in a leadership review. *"If they will not act, they are a coach, not a champion."*[^salesmotion][^b2bchampion] A common rubric: Grey (no champion) → Yellow (personal win known but EB access/influence unvalidated) → Green (booking EB meetings, shaping criteria, co-building the ROI case).[^prospeo]

**Miller Heiman Strategic Selling** (Miller & Heiman, 1985) names **four buying influences**: **Economic Buyer** (one; final yes/no), **User Buyer(s)** (judge job-performance impact), **Technical Buyer(s)** (gatekeepers/screeners, IT, security, legal, procurement, who **"can only say no, not yes"**), and **Coach** (insider who gives access).[^millerheiman][^avoma] **Terminology clash to avoid:** Miller Heiman's "Coach" ≈ the modern "Champion," whereas MEDDIC draws a *sharp* Coach-vs-Champion distinction. Don't conflate them.[^avoma]

### Identify and develop (arm) the champion
**Observable signs of a true champion** (behavior, not self-report): gives you the real objections and who's blocking; navigates politics and debriefs you after internal meetings; has **skin in the game** (their promotion/credibility rides on your success); can **articulate your value without your slides**; *requests ammunition*; opens doors to executives; tells you when things go badly.[^accountmap][^spotlightval]

**Arm them with a "consensus toolkit" / champion enablement** (strikingly uniform across sources):[^reworkchampion][^gtmplaybook][^kayvon]
- A **business-case/ROI template** pre-populated with your data, structured to answer the EB's questions (problem, cost of the problem, solution, ROI, risk of *doing nothing* vs. adopting).
- An **executive-ready narrative** (3-slide story / short deck: problem → why now → solution → proof → ROI → implementation → next steps).
- **Stakeholder-specific talking points** (CFO vs. CTO vs. VP Ops) + an **FAQ answering the objections you know are coming** (security, integration, pricing).
- **Peer proof** — references/case studies matched to their size, industry, and tech stack.
- **Objection rehearsal:** brief them before every internal meeting "like a witness before a deposition", role-play the EB's questions and give the exact answers.[^b2bchampion]
- A **Mutual Action Plan** that positions the champion as the project lead/hero. *"Help them look like heroes; position this as their initiative."*[^reworkstakeholder]

### Single-threaded risk and multi-threading
The most heavily corroborated practical claim in the corpus: a **single-threaded** deal dies if the champion **leaves/gets reorged, loses a political battle, or simply lacks the influence to build consensus.**[^multithread][^allston] *"When scrutiny increases, belief without political cover collapses."*[^meddmondays]

- **Multi-thread by ROLE, not headcount** — "one champion plus five end users is still single-threaded." Minimum role coverage: economic buyer, champion, technical/security evaluator, end-user rep, + blocker/exec sponsor; enterprise best practice **4–6 contacts across roles.**[^multithread][^prospectory]
- **Multi-thread *with* the champion, not around them** — frame it as reducing *their* political risk ("I want to make sure you have air cover from Finance and IT"). Going around the champion signals distrust and triggers backlash.[^reworkmulti]
- *Vendor close-rate stats, 71% of single-threaded deals end in no-decision/loss; deals with 3+ engaged contacts close at 2.4×; champion turnover ~33% per cycle, are repeatedly attributed to Gong/Chorus but were not independently verified here; treat as **TENTATIVE.***[^accountmapsingle] `verified-as-of: 2026-06-18`.

### Champion ⇄ Mobilizer (the cross-link)
The champion you develop should ideally be a **Mobilizer**, not a **Talker.** *The Challenger Customer* (CEB/Gartner) segments seven stakeholder profiles into three groups:[^challengercust][^stakeholderpdf]

- **Mobilizers** (drive consensus and organizational change): **Go-Getter, Teacher, Skeptic.**
- **Talkers** (agreeable but don't drive action): **Guide, Friend, Climber.**
- **Blocker** (favors the status quo): its own anti-category.

*(The three Mobilizers are Go-Getter, Teacher, Skeptic per Challenger Inc., Demand Gen Report quoting Spenner, and Salesforce; a minority of low-tier blogs swap Guide↔Teacher, none authoritative.)*[^challengerinc][^demandgen] Targeting Mobilizers correlates with high-performer status (Gartner: ~**31% more likely** to be a high performer).[^challengerinc] **Key warning:** the easiest, friendliest contact is often a **Talker** (≈ the MEDDIC "cheerleader"); a Mobilizer is frequently the skeptic who asks hard questions. Gartner's prescription, **"Commercial Coaching": arm the Mobilizer to build consensus** — is the same enablement motion as arming a champion.[^challengerinc] **Connection rule:** a validated MEDDIC Champion ≈ an engaged Gartner Mobilizer; a "cheerleader/fan" ≈ a Talker. Develop Mobilizers; don't be seduced by Talkers. *(Deep treatment of the seven profiles is summarized here only, the mechanism of consensus/social influence is owned by `persuasion-and-influence-psychology`.)*

---

## 6. FSI / big-bank buying behavior

Everything above intensifies in a large, regulated financial institution: more stakeholders with independent vetoes, longer cycles, and, crucially, **risk and lock-in become *regulatory* obligations, not just psychology.** The deep regulatory and procurement detail, the source list for it, and the vendor-action implications are in **`references/fsi-bank-buying.md`**. Load that file for an FSI deal. The essentials:

- **Bank buying is assurance-led, not demo-led.** Cycles run **9–18 months**; the **vendor risk assessment alone takes 3–6 months**, and security/InfoSec review is the dominant time sink. *"The sponsor is your champion; procurement is your buyer"* — the business sponsor cannot approve until risk/InfoSec/compliance/legal gates clear.[^fsi]
- **Regulation makes risk and exit mandatory.** US: the **Interagency Guidance on Third-Party Relationships (2023)** (a third-party risk life cycle, supervisory guidance, "not law" but examined), **OCC Heightened Standards** (covered banks ≥ $50B; new/unproven tech is itself a supervised strategic risk), **FFIEC** outsourcing guidance, and **model-risk guidance**. Note: **SR 11-7 was *revised* in April 2026 (SR 26-2 / OCC 2026-13)**, still covering vendor/third-party models. EU: **DORA** (binding since 17 Jan 2025) *legally requires* an ICT third-party risk strategy, a Register of Information, **concentration-risk assessment (Art. 29)**, mandatory contractual clauses (Art. 30), and a **credible, tested exit plan** for critical functions. `verified-as-of: 2026-06-18` — these specifics change; verify against primary regulators.[^fsi]
- **Lock-in is a *regulatory* concern here.** Inability to leave a vendor (no tested exit plan; excessive concentration) is a **supervised deficiency** under DORA, fundamentally different from the psychological switching-cost framing in Section 4. *"An untested exit is not credible."* For a **database vendor**, this flips lock-in into a **selling point**: demonstrable data egress, no one-way proprietary doors, multi-region/data-residency controls, and a rehearsed migration path satisfy the regulator's "credible exit plan" mandate.[^fsi]
- **Buyers buy "defensible execution, not innovation hype"** — referenceable outcomes from comparable institutions are decisive; the "reference-able peer bank" instinct counters institutional conservatism. *(Disconfirming: FSIs are increasing tech/cloud spend and self-identify as innovators, and marketplace standard contracts compress legal cycles, but governance gates remain the binding constraint; deals got bigger and more deliberate, not faster.)*[^fsi]

---

## TAM field playbook (one screen)

1. **Map the DMU first.** Name the initiator, users, influencers, decider, buyer, gatekeeper(s), approver. New-task buys → expect the biggest committee (Section 1).
2. **Assume "No Decision" is your main competitor.** Work *both* the individual's fear of being wrong (JOLT) and the group's inability to reach a collective yes (Challenger), not just the rival vendor (Section 2).
3. **Find/validate a champion by testing for *action*, not warmth.** A friendly contact who won't spend political capital is a Coach/Talker, not a Champion/Mobilizer (Section 5).
4. **Arm the champion with a risk-tilted consensus toolkit**: business case, peer-bank proof, stakeholder talking points, objection rehearsal. Position them as the hero (Section 5).
5. **Multi-thread by role, with the champion** (air cover, not betrayal). Cover the vetoes (Section 5).
6. **De-risk the decision on all three layers**: rational (fit/ROI), political (whose budget/turf/career), emotional (fear of being wrong). For a challenger displacing an incumbent, make choosing you the *safe* career move (Section 3).
7. **In FSI: engage risk/InfoSec/procurement in week 1, in parallel; lead with a "vendor trust pack" and a credible, tested exit story.** Turn DORA's exit/concentration mandate into a differentiator (Section 6 + `fsi-bank-buying.md`).
8. **Stay honest about the evidence.** Frameworks are expert heuristics; vendor stats are often unverified; date the volatile claims. Cite the mechanism references for the *why*.

**What to capture (the deliverable).** Turn this analysis into a living account artifact, refreshed each touchpoint: (a) a **DMU role table** (person → role → stance), (b) a **champion status** flag per the Grey/Yellow/Green rubric (Section 5), (c) a **role-coverage checklist** for multi-threading (which of EB / champion / technical-security / compliance / procurement / end-user are covered), and (d) the **open risk gates** (security review, exit-plan/concentration, model validation) with owners. This plugs into the account-review and QBR deliverables owned by `tam-operations`.

---

## References

Full citation list (every footnote marker above) lives in **`references/sources.md`**: ~85 sources with tiers, extracted to keep this on-demand reference lean. FSI-specific sources are in **`references/fsi-bank-buying.md`**.
