<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-nc-entity-lifecycle` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-nc-entity-lifecycle
description: >-
  Lifecycle changes to an existing FOR-PROFIT North Carolina entity, for a solo founder. (1) DISSOLUTION:
  Articles of Dissolution at the NC SoS (LLC Form L-07; corp B-05/B-06), winding-up, creditor notice and
  the claim-bar, NCDOR final returns + Form NC-BN, owner distributions, LLC-vs-corp differences. (2)
  CONVERSION: sole proprietorship -> LLC, an LLC ELECTING S-corp tax status vs a statutory CONVERSION to a
  corporation, merger mechanics, and the EIN / tax-year fallout. (3) FOREIGN QUALIFICATION: taking a NC
  entity to another state, or registering an out-of-state entity in NC (Certificate of Authority,
  registered agent, "transacting business" thresholds). (4) REINSTATEMENT after administrative dissolution
  (back annual reports + fee). (5) the NC Taxed PTE / PTET election (the SALT-cap workaround — eligibility,
  how to elect, rate, who benefits). Cites the NC SoS, NCDOR, and IRS with as-of dates.
  TRIGGER: dissolving/closing a NC LLC or corp; converting a sole prop to an LLC, or an LLC to a corp;
  S-corp election vs statutory conversion; foreign qualification / Certificate of Authority in or out of
  NC; reinstating an administratively dissolved entity; the NC Taxed PTE / PTET / SALT-cap election.
  SKIP: INITIAL formation, entity choice, EIN setup, BOI/FinCEN, annual report -> venture-nc-business-
  formation-tax; nonprofit dissolution / 501(c)(3) -> venture-nc-nonprofit-formation; business plan ->
  venture-small-business-planning; hiring/payroll -> venture-nc-employer-payroll.
category: personal-venture
tags: [venture, north-carolina, entity-lifecycle, dissolution, conversion, ptet]
version: 1.0.0
updated: 2026-06-16
whenToUse:
  - Dissolving or winding down a NC LLC or business corporation and filing Articles of Dissolution
  - Handling creditors, final NCDOR returns, and member/shareholder distributions when closing a NC entity
  - Converting a sole proprietorship into an LLC, or statutorily converting an LLC into a corporation
  - Deciding between an S-corp TAX election and an actual statutory conversion to a corporation
  - Qualifying a NC entity to do business in another state, or registering a foreign entity in NC
  - Reinstating a NC entity that was administratively dissolved for missed annual reports
  - Evaluating or making the NC Taxed Pass-Through Entity (PTET) election as a SALT-cap workaround
  - Understanding when a restructuring forces a new EIN or a short tax year
triggers:
  - dissolve NC LLC or corporation
  - NC Articles of Dissolution L-07 B-06
  - convert sole proprietorship to LLC
  - LLC convert to corporation S-corp election
  - NC Certificate of Authority foreign qualification
  - reinstate administratively dissolved NC entity
  - NC Taxed PTE PTET election SALT cap
  - new EIN after entity conversion
---

# NC Entity Lifecycle: Dissolution, Conversion, Foreign Qualification, Reinstatement & PTET

Changes to a NC for-profit entity **after** it exists. For the birth of the entity (entity choice,
Articles of Organization, EIN, DBA, BOI/FinCEN, the annual report itself) use
**venture-nc-business-formation-tax**. Nonprofit dissolution and 501(c)(3) matters differ on the asset
side — use **venture-nc-nonprofit-formation**.

Two agencies do most of the work and they do **not** talk to each other for closing or restructuring a
for-profit entity:
- **NC Secretary of State (SoS), Business Registration Division** — the *legal* existence of the
  entity (formation, conversion, dissolution, reinstatement, foreign authority). Forms and fees at
  sosnc.gov.
- **NC Department of Revenue (NCDOR)** — the *tax* accounts (sales & use, withholding, corporate
  income/franchise) and final returns. ncdor.gov.
- **IRS** — federal classification, the EIN, and the S-election.

> **NC has no general "tax clearance certificate" prerequisite to dissolve a for-profit entity.** Unlike
> some states, the NC SoS does **not** require an NCDOR clearance letter before accepting Articles of
> Dissolution for an LLC or business corporation [UNVERIFIED — confirm against current sosnc.gov
> instructions]. You still must file final returns and close your NCDOR tax accounts separately; the two
> tracks just run in parallel. (Nonprofit dissolution is different — it routes through an asset-distribution
> plan; see venture-nc-nonprofit-formation.)

All forms/fees below are **as of 2026-06**; the SoS revises form numbers and fees periodically — verify
the current PDF and fee on sosnc.gov before filing.

---

## 1. Dissolution / winding down a NC entity

Closing a NC business is a **voluntary** act: you file Articles of Dissolution for the entity type, then
the entity may only do "wind-up" work — it does not vanish on filing.
(Source: NC SoS, "Closing a North Carolina Business," sosnc.gov, as of 2026-06.)

### 1a. The forms and fees (NC SoS, as of 2026-06)

| Entity | Form | Fee | Use |
|---|---|---|---|
| **LLC** | Articles of Dissolution (**L-07**) | **$30** | Dissolve an LLC on the SoS record |
| **Business corp — no shares issued** | Articles of Dissolution (**B-05**) | **$30** | Dissolved by incorporators/directors before any shares issued |
| **Business corp — shares issued** | Articles of Dissolution (**B-06**) | **$30** | Dissolved by board + shareholders |

Mail with check to: NC Secretary of State, Business Registration Division, PO Box 29622, Raleigh, NC
27626-0622 (verify current PO Box on the form), or file online via the SoS portal.

### 1b. LLC vs corporation — the real difference is the *approval* step

- **LLC (NCGS Ch. 57D):** dissolution is triggered per the operating agreement or by the members
  (commonly unanimous member consent absent a contrary operating-agreement provision), then file L-07.
  Fewer formal corporate steps.
- **Corporation (NCGS Ch. 55, Article 14):** a **two-track** approval —
  - **No shares issued:** the directors (or a majority of incorporators if there are no directors) may
    dissolve and file B-05 (NCGS 55-14-01).
  - **Shares issued:** the **board recommends** dissolution and the **shareholders approve** it, then
    file B-06 (NCGS 55-14-02/03).
  A corporation is dissolved on the **effective date of the Articles of Dissolution** (NCGS 55-14-03),
  after which it continues to exist *only* to wind up (NCGS 55-14-05).

### 1c. Winding-up steps (do these, in roughly this order)

1. **Approve the dissolution** at the right level (member consent / board + shareholder vote) and
   document it in writing (minutes or written consent).
2. **File Articles of Dissolution** (L-07 / B-05 / B-06) with the NC SoS.
3. **Stop normal business**; do only wind-up acts — collect assets, sell property, discharge
   liabilities (NCGS 55-14-05 for corps; 57D for LLCs).
4. **Give creditor notice** (see 1d) — this is what limits your future exposure.
5. **Pay/settle creditors first**, then distribute what remains (see 1e).
6. **Close NCDOR tax accounts and file final returns** (see 1f).
7. **Close the EIN account with the IRS** (optional but tidy) — send a letter to the IRS asking to close
   the business account; the EIN itself is never reused. (Source: IRS, "Canceling an EIN / closing your
   account," irs.gov.)
8. **Cancel** local privilege licenses/permits, bank accounts, registered-agent service, insurance, and
   any **foreign registrations** in other states (file a withdrawal there — see §3).

### 1d. Creditor notice and the claim-bar (corporations — NCGS 55-14-06/07)

This is the payoff of dissolving *formally*: you can cut off creditor claims.

- **Known claims (NCGS 55-14-06):** the dissolved corporation may deliver **written notice** to each
  known claimant describing the claim, telling them where to send it, and setting a **deadline** by
  which the claim must be received. A known claim is **barred** if the claimant misses that stated
  deadline (the statute sets a minimum notice window — verify the current number of days in NCGS
  55-14-06 [UNVERIFIED — exact day count]).
- **Unknown / contingent claims (NCGS 55-14-07):** the corporation may **publish** a notice of
  dissolution **once** in a newspaper of general circulation in the county of its principal office,
  requesting claims. Claims by anyone who didn't get written notice are then **barred unless suit is
  filed within five (5) years** of the publication date.
- **LLCs:** NC's LLC Act (Ch. 57D) provides for winding up and disposing of known/unknown claims on a
  comparable structure; the practical move (written notice to knowns + published notice for unknowns) is
  the same [UNVERIFIED — confirm the specific 57D claim-bar sections].

### 1e. Distributions to members / shareholders

**Creditors before owners — always.** Distributing assets to owners while leaving creditors unpaid can
expose those owners to clawback. After liabilities are paid or reserved for, distribute the remainder:
per the operating agreement / shareholder rights (preferences first, then common). Liquidating
distributions have **tax consequences** (gain/loss vs basis; possible deemed asset sale for a corp) —
run them past a CPA.

### 1f. NCDOR final returns + closing accounts

- **File final returns** for every tax type the entity was registered for: corporate income & franchise
  (Form CD-405, or CD-401S for an S-corp), partnership (D-403), sales & use (E-500), and withholding.
  Mark them **final** where the form allows.
- **Form NC-BN, "Out-of-Business Notification"** — file this (online or by mail) to tell NCDOR you've
  ceased business and to close your sales-and-use and withholding accounts. (Source: NCDOR, "NC-BN
  Out-of-Business Notification," ncdor.gov, as of 2026-06.)
- **If you had employees:** within **30 days** of closing, file **Form NC-3** plus the required W-2/1099
  statements and furnish copies to workers. (Source: NCDOR Withholding FAQ, ncdor.gov.) Federal final
  Form 941 (check the "final return" box) and Form 940 are also due.
- The annual-report obligation stops once the entity is dissolved on the SoS record (don't keep paying
  it after dissolution).

---

## 2. Conversion & restructuring

"Conversion" gets used loosely. Separate three very different things:

### 2a. Sole proprietorship -> LLC (this is *formation*, not a statutory conversion)

A sole proprietorship has **no separate legal existence**, so there is nothing to "convert" at the SoS.
You simply **form a new LLC** (file Articles of Organization) and move the business into it. See
**venture-nc-business-formation-tax** for the formation steps.

- **EIN:** a single-member LLC taxed as a **disregarded entity** with no employees and no excise tax may
  keep using the owner's SSN/EIN — **no new EIN required**. But if the new LLC will be **taxed as a
  corporation or S-corp**, or will **have employees**, it needs a **new EIN**. (Source: IRS, "When to get
  a new EIN," irs.gov, as of 2026-06.)
- Transfer assets, re-paper contracts/bank accounts in the LLC's name, and update licenses/registrations.

### 2b. LLC **electing** S-corp tax status (a tax election — the entity does NOT change)

This is the common move and it is **purely a federal tax election** — your LLC stays an LLC at the NC
SoS; only its IRS classification changes.

1. The LLC must first be (or elect to be) taxed as a **corporation**, then elect S status — in practice
   you file **IRS Form 2553** (an LLC can use 2553 to elect S-corp treatment and is deemed to have made
   the underlying entity-classification election). (Source: IRS Form 2553 instructions, irs.gov.)
2. **Timing:** generally file Form 2553 by **2 months and 15 days** after the start of the tax year the
   election is to take effect (late-election relief exists under Rev. Proc. 2013-30).
3. **No new EIN** is needed merely to elect S-corp status for an existing entity that already has its own
   EIN; but an LLC that newly elects to be taxed as a corp/S-corp **does** need its own EIN if it was a
   disregarded entity using the owner's SSN. (Source: IRS, "When to get a new EIN.")
4. **NC side:** NC follows the federal S election; the entity files **CD-401S** (NC S-corp return). No
   separate NC S-election form. See venture-nc-business-formation-tax for the reasonable-compensation /
   payroll mechanics of running an S-corp.

> S-corp **election** vs statutory **conversion**: electing S status keeps your LLC and its liability
> shield exactly as-is — you only change tax treatment. A statutory conversion (2c) actually makes the
> entity a corporation under state law (new governance: bylaws, board, stock). Most solo founders want
> the **election**, not the conversion.

### 2c. NC statutory conversion / merger (the entity's *legal type* actually changes)

NC allows a business entity to **convert** into a different type (e.g., LLC -> business corporation) by
adopting a **plan of conversion**, getting it approved by the converting entity, and filing **articles
that include articles of conversion** with the SoS. (Authority: NCGS Ch. 55 Art. 11A and Ch. 57D Art. 9.)

- **LLC -> corporation:** file **"Articles of Incorporation including Articles of Conversion"** (NC SoS
  form, fee **$125** as of 2026-06 — i.e., the corp-formation fee). The resulting corporation is treated
  as the same entity that existed before; it keeps its history but is now governed by corporate law.
- **Mergers** (combining two entities) use Articles of Merger and are governed by the same chapters.
- **Tax/EIN fallout of an actual conversion:**
  - **LLC -> corporation** is, for federal tax, the formation of a corporation; per IRS rules **a new
    EIN is generally required** when an LLC converts to a corporation (or vice-versa). (Source: IRS, "When
    to get a new EIN," irs.gov.)
  - A change of tax classification can create a **short tax year** and require a final return for the old
    classification and a first return for the new one. Confirm with a CPA before the effective date.

> **Pick the cheapest path that gets the outcome you want.** If the goal is the S-corp payroll-tax
> benefit, **elect** S status (2b) and skip the statutory conversion entirely. Reserve a true conversion
> (2c) for when you actually need corporate governance/stock (e.g., taking on investors who want shares).

---

## 3. Foreign qualification (operating across state lines)

"Foreign" = formed in another state, not overseas. Two directions:

### 3a. An out-of-state entity registering **in NC**

A foreign entity that will **transact business** in NC must obtain a **Certificate of Authority** from
the NC SoS. (Source: NC SoS, "Register a Foreign Business," sosnc.gov, as of 2026-06.)

**Application requirements (NC SoS):**
- **Application for Certificate of Authority** for the entity type (foreign LLC / foreign corporation).
- A **registered agent and registered office in NC** (an NC resident or an entity authorized in NC).
- A **Certificate of Existence / Good Standing** from the home state, dated **no more than 6 months** old.
- Entity name, home jurisdiction + formation date, principal office, and officer/manager info. If the
  name isn't available in NC, the entity must use a **fictitious (assumed) name** in NC.

**Fees (NC SoS, as of 2026-06):**

| Foreign entity type | Application fee |
|---|---|
| Foreign **LLC** | **$250** |
| Foreign **business corporation** | **$250** |
| Foreign **limited partnership** | **$50** |

### 3b. "Transacting business" — when you actually have to register

NC (like the model acts) lists activities that **do NOT** by themselves constitute transacting business
(the safe harbor — NCGS 55-15-01 for corporations; 57D-7-01 for LLCs), e.g.: maintaining or defending a
lawsuit; holding manager/member meetings; maintaining bank accounts; selling through independent
contractors; soliciting/obtaining orders that require acceptance outside NC before they become contracts;
isolated transactions completed within 30 days; and owning property. **Crossing the line** usually means
a physical presence, NC employees, a place of business, or repeated in-state contracts.

**Penalty for not qualifying when you should (NCGS 55-15-02 / 57D-7-02):** a foreign entity transacting
business in NC without a Certificate of Authority **may not maintain a lawsuit in NC courts** until it
qualifies, and may owe back fees and penalties. It does **not** void the entity or its contracts.
[UNVERIFIED — confirm current penalty amounts and whether they accrue per year.]

### 3c. Taking a **NC entity into another state**

Mirror image: register your NC entity as a *foreign* entity in **that** state (each state has its own
Certificate of Authority, fee, and registered-agent requirement). You'll typically need a **Certificate
of Existence from the NC SoS** to attach. When you leave a state, file that state's **withdrawal /
cancellation of authority** so you stop owing its annual fees. Research the target state's rules
specifically — they vary widely.

---

## 4. Reinstatement after administrative dissolution

If you stop filing **annual reports** (or lose your registered agent), the NC SoS can **administratively
dissolve** the entity. (Source: NC SoS Administrative Dissolution manual, sosnc.gov.)

**How it happens:**
1. The SoS sends a **Notice of Grounds for Dissolution (NOG)** (a postcard) when an entity is out of
   compliance with the annual-report requirement.
2. If not cured, the entity is **administratively dissolved** — it loses good standing and may not carry
   on business except to wind up, and it can lose exclusive rights to its name.

**Reinstatement (NCGS 57D-6-06 for LLCs; the LLC statute adopts the corporate procedure in NCGS
55-14-22/23/24):**
1. **File all delinquent annual reports** and pay those past-due report fees to cure the ground for
   dissolution.
2. **File the Application for Reinstatement** with the SoS, reciting the entity name and the effective
   date of administrative dissolution and stating that the ground(s) no longer exist.
3. **Pay the reinstatement fee — $100** (as of 2026-06). (Source: NC SoS fee schedule, sosnc.gov.)
4. There is generally **no statutory time limit** to apply for reinstatement, and once granted,
   reinstatement **relates back** to the date of dissolution as if it never occurred (NCGS 55-14-22(c)).

> If the entity's **name** was taken by someone else during the dissolution, you may have to choose a new
> name or use an assumed name as part of reinstatement. Reinstatement is far cheaper than re-forming and
> it restores your operating history — do it rather than starting over.

---

## 5. The NC Pass-Through Entity Tax ("Taxed PTE" / PTET) election

The **SALT-cap workaround.** Since the federal $10,000 cap on the state-and-local-tax (SALT) itemized
deduction hurt owners of pass-throughs, NC (like most states) lets an eligible PTE **elect to pay NC
income tax at the entity level** so the tax becomes a *business* deduction that bypasses the individual
SALT cap. Effective for tax years beginning on/after **Jan 1, 2022**. (Sources: NCDOR "Important Notice
Regarding North Carolina's Recently Enacted Pass-Through Entity Tax" and **Directive TA-23-1**,
ncdor.gov, as of 2026-06; Session Law 2023-134.)

### 5a. Who is eligible

- An **eligible S corporation** and an **eligible partnership** (including an **LLC taxed as a
  partnership** or as an **S-corp**). A single-member LLC that is a **disregarded entity is not** a PTE
  and cannot elect.
- **Owner-composition rule:** the entity was originally ineligible if it had **any owner that was itself
  a corporation or another pass-through**. A **2023 expansion (SL 2023-134)** broadened eligibility so a
  PTE with a **qualifying corporate partner** or a **qualifying trust partner** can now elect — but the
  Taxed PTE's tax base **excludes** the share allocable to a qualifying *corporate* partner, while a
  qualifying *trust* partner's share **is included**. (Source: NCDOR Directive TA-23-1.)
  [UNVERIFIED — confirm whether tiered/multi-PTE ownership disqualifies under the current statute.]

### 5b. How to elect

- The election is made **on the entity's timely-filed annual NC return**, by **checking the "Taxed PTE
  Election" box** on the front of the return:
  - **S corp:** Form **CD-401S**.
  - **Partnership / LLC-as-partnership:** Form **D-403**.
- **Deadline:** by the **due date of the return, including extensions.** An election on a **late-filed**
  return is **not valid**. (Source: NCDOR Important Notice, ncdor.gov.)
- The election is **annual** (made return-by-return). [UNVERIFIED — confirm current-year revocation
  rules; NCDOR has issued guidance allowing certain amended-return elections for specific years.]

### 5c. The rate and how owners are made whole

- The Taxed PTE pays NC income tax on its NC taxable income at the **NC individual income-tax rate** for
  the year: **4.5% (TY2024), 4.25% (TY2025), 3.99% (TY2026 and after)**. (Source: NCDOR tax-rate
  schedules, ncdor.gov.)
- **Owners then deduct their share** of the Taxed PTE's income from their NC taxable income (to the
  extent it was in the PTE's NC taxable income and in the owner's AGI), reported via the **NC-PE**
  schedule — so the income is **not taxed twice** at the NC level. (Source: NCDOR, NC-PE instructions and
  Important Notice.)

### 5d. Who actually benefits (and who shouldn't bother)

- **Benefits:** profitable PTEs whose owners **itemize federally and are over the $10k SALT cap** — the
  entity-level NC tax becomes a federal deduction the owners couldn't otherwise take. Best for
  higher-income owners with meaningful NC PTE income.
- **Often not worth it:** owners taking the federal **standard deduction**; very low NC income; or out-of
  -state owners whose home state won't give a credit for NC entity tax. The election can also complicate
  estimated payments (the PTE may owe NC estimates; see **Form NC-429B PTE** underpayment rules) and
  reduce some owners' NC credits.

> **Moving target.** The federal SALT cap and any PTET workaround are politically active — the cap and
> NC's response have changed before and may change again. **Confirm the current-year rate, the SALT-cap
> status, and the election mechanics with a CPA before relying on this.**

---

## Quick cross-references

- **Initial formation, entity choice, EIN setup, DBA, BOI/FinCEN, the annual report itself** ->
  **venture-nc-business-formation-tax**
- **Nonprofit dissolution, 501(c)(3), asset-distribution-on-dissolution, Form 990** ->
  **venture-nc-nonprofit-formation**
- **Business plan, financial model, funding, valuation** -> **venture-small-business-planning**
- **Hiring your first employee, payroll/withholding/SUTA setup** -> **venture-nc-employer-payroll**

---

## Sources

Primary (verify the live page before filing/electing — fees and forms change):

- **NC Secretary of State — "Closing a North Carolina Business"** (Articles of Dissolution forms L-07,
  B-05, B-06; $30 fees; wind-up): https://www.sosnc.gov/divisions/business_registration/closing_nc_business
  (as of 2026-06)
- **NC Secretary of State — "Register a Foreign Business" / Requirements** (Certificate of Authority,
  registered agent, certificate of existence <6 months, $250 foreign LLC/corp, $50 foreign LP):
  https://www.sosnc.gov/manual/Register_A_Foreign_Business (as of 2026-06)
- **NC Secretary of State — Administrative Dissolution manual** (NOG notice; reinstatement):
  https://www.sosnc.gov/manual/Administrative_Dissolution (as of 2026-06)
- **NC Secretary of State — conversion forms** ("Articles of Incorporation including Articles of
  Conversion," $125): https://www.sosnc.gov/forms (as of 2026-06)
- **NCGS Chapter 55, Article 14 — Dissolution** (corp two-step dissolution, winding up, known/unknown
  claims, 5-year publication bar; §§55-14-01/02/03/05/06/07; reinstatement §§55-14-22/23/24):
  https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/ByArticle/Chapter_55/Article_14.pdf
- **NCGS Chapter 57D — NC Limited Liability Company Act** (LLC dissolution & administrative dissolution
  §57D-6-06; conversion Art. 9; foreign safe harbor §57D-7-01):
  https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_57D/GS_57D-6-06.pdf
- **NCDOR — "Important Notice Regarding North Carolina's Recently Enacted Pass-Through Entity Tax"**
  (Taxed PTE eligibility, election on the return, owner deduction): https://www.ncdor.gov/taxes-forms/information-tax-professionals/tax-bulletins-directives-and-other-important-notices/important-notices-and-frequently-asked-questions-personal-taxes/important-notice-regarding-north-carolinas-recently-enacted-pass-through-entity-tax
  (as of 2026-06)
- **NCDOR — Directive TA-23-1** (2023 expansion of eligible partners; qualifying corporate/trust
  partner): https://www.ncdor.gov/taxes-forms/individual-income-tax/personal-taxes-division-directives/directive-ta-23-1
- **NCDOR — Form NC-BN, Out-of-Business Notification**: https://www.ncdor.gov/nc-bn-out-business-notification
- **NCDOR — Tax Rate Schedules** (individual/PTE rate: 4.5% TY2024, 4.25% TY2025, 3.99% TY2026+):
  https://www.ncdor.gov/taxes-forms/individual-income-tax/tax-rate-schedules
- **NCDOR — NC-429B PTE** (Taxed PTE estimated-tax underpayment):
  https://www.ncdor.gov/nc-429b-pte-underpayment-estimated-tax-taxed-pass-through-entities
- **IRS — "When to get a new EIN"** (sole-prop->LLC, LLC->corp, classification changes):
  https://www.irs.gov/businesses/small-businesses-self-employed/when-to-get-a-new-ein
- **IRS — Form 2553 instructions** (S-corp election timing): https://www.irs.gov/instructions/iss4
  (and irs.gov Form 2553)

Items flagged **[UNVERIFIED]** above were not confirmed against the primary statute/form text in this
pass and should be checked before relying on them: the absence of any NC tax-clearance prerequisite to
dissolve; the exact NCGS 55-14-06 known-claim notice window; the specific NCGS 57D LLC claim-bar
sections; current foreign-qualification penalty amounts; tiered-PTE-ownership disqualification; and
current-year PTET revocation rules.

---

## Disclaimer

This is general educational information, **not legal or tax advice**, and it is **not a substitute for a
North Carolina attorney or CPA**. Forms, fees, statutes, tax rates, and the SALT-cap/PTET landscape
change — everything here is **as of 2026-06** and must be verified against the live **NC Secretary of
State (sosnc.gov)**, **NCDOR (ncdor.gov)**, and **IRS (irs.gov)** sources before you act. Dissolution,
conversion, and the PTET election have lasting legal and tax consequences (creditor liability, gain
recognition, EIN/tax-year effects); consult a qualified NC attorney and CPA for your specific situation.
