<!-- hub-reference-banner -->
> **Reference file — part of the `venture-nonprofit-cause` hub.** Formerly the standalone `venture-nc-nonprofit-formation` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-nc-nonprofit-formation
description: >-
  Guide to forming a North Carolina nonprofit and keeping it compliant: structure choice (GS Ch. 55A
  nonprofit corporation, unincorporated association, or fiscal sponsorship), incorporating with IRS-required
  501(c)(3) purpose and dissolution clauses, governance, EIN, federal 501(c)(3) recognition, NC tax
  exemptions via NCDOR, NC Charitable Solicitation License (GS Ch. 131F), and Form 990 compliance. General to
  any NC charitable nonprofit; tuned for an organ-donation-awareness venture. TRIGGER: forming or maintaining
  a NC nonprofit; getting or keeping 501(c)(3) status; NC nonprofit tax exemption; charitable-solicitation or
  fundraising registration; Form 1023/1023-EZ; public charity vs. private foundation; Form 990 filing;
  nonprofit bylaws/board; must we register before fundraising. SKIP: FOR-PROFIT NC entity formation,
  LLC/C-corp/S-corp choice, NC business tax (taxable entity), or DBA/assumed-name filings → venture-nc-business-formation-tax.
category: personal-venture
version: 1.0.1
updated: 2026-06-16
tags: [venture, nonprofit, north-carolina, 501c3, charitable-solicitation, legal]
metadata:
  changelog:
    - "2026-06-16 sko v1.0.0->v1.0.1 — Pass H 10/10 pos, 0/10 neg (predicted); fixed 3 Medium (Pass M: description 1501->989 chars under Glean cap; Pass N: removed dangling SKIP target; Pass K: em-dash density 1.32->0.97/100), added version/updated; 0 banned terms; first-time hub registration"
whenToUse:
  - Deciding whether to incorporate as an NC nonprofit, stay an unincorporated association, or use a fiscal sponsor
  - Drafting Articles of Incorporation with the IRS-required 501(c)(3) purpose and dissolution clauses
  - Setting up board, bylaws, conflict-of-interest policy, and the organizational meeting
  - Applying for an EIN and federal 501(c)(3) recognition (Form 1023 vs. 1023-EZ)
  - Getting NC corporate income/franchise exemption and the NC sales-and-use tax refund
  - Determining whether the venture must obtain a NC Charitable Solicitation License before fundraising
  - Planning ongoing compliance: Form 990 series, recordkeeping, and license renewals
  - Figuring out registration obligations when soliciting donations in multiple states
triggers:
  - NC nonprofit formation
  - 501(c)(3) application
  - charitable solicitation license
  - Form 1023 vs 1023-EZ
  - nonprofit bylaws and board
  - NC nonprofit tax exemption
  - Form 990 filing
  - fundraising registration
---

# Forming and Maintaining a North Carolina Nonprofit (501(c)(3) + Charitable-Solicitation Compliance)

This skill walks a founder through legally standing up a North Carolina charitable nonprofit and keeping
it compliant across three regulators: the **NC Secretary of State** (incorporation + charitable-solicitation
licensing), the **NC Department of Revenue / NCDOR** (state tax exemptions), and the **IRS** (federal
501(c)(3) recognition + annual Form 990 filing). It is written for a cause venture (e.g., organ-donation
awareness) but applies to any NC charitable nonprofit.

There are three legally distinct things people lump together as "starting a nonprofit," and they happen in
order:

1. **State entity formation** — becoming a legal nonprofit corporation in NC (or choosing not to incorporate).
2. **Federal tax-exempt recognition** — getting the IRS to recognize you as a 501(c)(3) so donations are
   deductible and you are exempt from federal income tax.
3. **Fundraising authorization + state tax treatment** — getting a NC Charitable Solicitation License (if
   required) before asking the public for money, and securing NC tax exemptions.

You can be a corporation without being tax-exempt, and you must usually be licensed to solicit *separately*
from being tax-exempt. Keep the three lanes mentally separate.

> This is general information, **not legal or tax advice.** Fees, forms, and thresholds change — verify every
> number against the primary source before you act, and engage a licensed NC attorney and CPA. See the
> Disclaimer at the end.

---

## 1. Choose the structure

Three viable structures, roughly in order of how most founders proceed:

### A. NC nonprofit corporation (the default, and what 501(c)(3) usually assumes)
Formed under the **North Carolina Nonprofit Corporation Act, GS Chapter 55A** (NC General Statutes Chapter
55A, ncleg.gov, as of 2026-06). This is the standard vehicle for a 501(c)(3). Benefits:

- **Limited liability** for directors, officers, and members; the corporation, not individuals, bears its
  obligations.
- **Perpetual existence** independent of any one founder.
- It is the structure the IRS Form 1023/1023-EZ process is built around (the application asks for your
  organizing document's purpose and dissolution clauses, which a corporation supplies via its Articles).
- Filings under Chapter 55A are administered through the business-entity filing rules of **GS Chapter 55D**
  (ncleg.gov, as of 2026-06).

**Recommended for the organ-donation venture** if it will solicit donations, hold assets, sign contracts,
or seek grants.

### B. Unincorporated nonprofit association
A group can operate as an informal unincorporated association without filing Articles. It *can* still get an
EIN and even apply for 501(c)(3) status. But:

- **No corporate liability shield** in the same clean way — members can carry more personal exposure.
- Harder to open bank accounts, sign leases, and win grants; many funders require incorporation.
- Generally only sensible for a tiny, short-lived, low-risk effort. **Most cause ventures should incorporate.**

### C. Fiscal sponsorship (an alternative to forming your own entity, at least at first)
Instead of forming a new 501(c)(3), you operate under the umbrella of an existing 501(c)(3) "fiscal sponsor."
The sponsor receives tax-deductible donations on your behalf and lets you use its exempt status, typically
for an administrative fee (often a percentage of funds raised). Trade-offs:

- **Fast and cheap to start** — no incorporation, no Form 1023, immediate deductible donations.
- You give up autonomy and pay a cut; the sponsor controls the funds and bears legal responsibility.
- Good as a **bridge** while you decide whether the venture is durable, or as a permanent home for a small
  project. [UNVERIFIED — fiscal-sponsorship terms are private contracts, not statute; vet any sponsor's
  agreement and IRS-recognized status with counsel.]

**Decision heuristic:** durable, donor-facing, asset-holding cause → incorporate (A) and pursue 501(c)(3).
Testing the idea or want speed → start under a fiscal sponsor (C), incorporate later.

---

## 2. Incorporation steps (NC nonprofit corporation)

The incorporation document is the **Articles of Incorporation for a Nonprofit Corporation, Form N-01**, filed
with the NC Secretary of State under **GS §55A-2-02** (sosnc.gov; ncleg.gov, as of 2026-06).

1. **Pick a compliant name.** It must be distinguishable from other NC entities. Check availability on the
   SOS business search (sosnc.gov). Name-style requirements are governed by GS Chapter 55D.
2. **Appoint a registered agent and registered office in NC.** Required — the agent is the entity's official
   point of contact for legal service. Can be an individual NC resident or a qualified business; the office
   must be a NC street address (GS Chapter 55A / 55D, ncleg.gov, as of 2026-06).
3. **Name at least one incorporator.** The incorporator signs and files the Articles.
4. **Write the purpose and dissolution clauses with IRS 501(c)(3) language** (this is the step most DIY
   founders get wrong — see Section 3 below). The bare NC form will form a valid corporation, but **without
   the specific 501(c)(3) language the IRS will reject your exemption application.**
5. **File Form N-01 with the NC Secretary of State and pay the filing fee.** The fee is **$60** for nonprofit
   Articles of Incorporation (sosnc.gov, reported as of 2026-06 — *verify on the current SOS fee schedule
   before filing*). [UNVERIFIED — the $60 figure came from a SOS-sourced search summary; the live fee page
   and the Form N-01 instructions could not be machine-read in this research pass. Confirm the exact current
   amount and any expedited-service surcharge directly with sosnc.gov.]
6. **Effective date.** The corporation exists once the SOS accepts the filing (or on a delayed effective date
   if you request one).

> Tip: Some founders file the IRS-ready language directly in the Articles; others file minimal Articles then
> amend. Filing it right the first time avoids a later **Articles of Amendment (Form N-02)** and re-fee.

---

## 3. The IRS-required 501(c)(3) language (put it in the Articles)

To qualify for 501(c)(3), your **organizing document must pass the IRS "organizational test"**: it must limit
your purposes to exempt purposes and must dedicate your assets to an exempt purpose on dissolution (IRS,
Instructions for Form 1023-EZ, Rev. January 2025, irs.gov; IRS Pub. 557). Two clauses do this:

### Purpose clause
State that the corporation is organized **exclusively for one or more 501(c)(3) purposes** — charitable,
educational, etc. For an organ-donation venture, "charitable and educational purposes, including increasing
public awareness of organ and tissue donation" fits within the charitable/educational categories. The
operative IRS magic words are that the organization is organized exclusively for purposes described in
**section 501(c)(3) of the Internal Revenue Code**.

### Dissolution clause
State that on dissolution, remaining assets will be distributed for an exempt purpose (to another 501(c)(3),
or to a government for a public purpose) — **assets may not go to private individuals, founders, or members.**
This is the "dedication of assets" requirement. NC GS §55A also addresses distribution of assets on
dissolution, but the IRS requires its specific exempt-dedication language in your organizing document.

### Two prohibitions the operational test enforces (bake these into bylaws and practice)
- **No private inurement / private benefit** — net earnings may not benefit insiders.
- **No (or only insubstantial) lobbying, and absolutely no political campaign intervention** — 501(c)(3)s
  may not support or oppose candidates (IRS, Instructions for Form 1023-EZ, Rev. Jan. 2025, irs.gov).

IRS Publication 557 ("Tax-Exempt Status for Your Organization," irs.gov) gives sample purpose and dissolution
clause text — use it as the template.

---

## 4. Governance (board, bylaws, COI policy, organizational meeting)

### Board of directors
- Under **GS Chapter 55A**, a NC nonprofit corporation must have a board of directors; the statute sets a
  **minimum of one director** (ncleg.gov, as of 2026-06). [UNVERIFIED — confirm the exact current minimum and
  any officer requirements in GS §55A-8-03 and the officer provisions; statute text could not be fully
  machine-read this pass.]
- **Practical recommendation, not a legal minimum:** the IRS and good governance favor **at least 3 unrelated
  directors**, because a board dominated by related parties or a single person raises private-benefit and
  governance concerns on Form 1023 and with funders. Most credible 501(c)(3)s run 3+ independent directors.

### Bylaws
Not filed with the state, but essential. Bylaws set how the organization is governed: number and terms of
directors, officer roles, meeting and quorum rules, voting, committees, indemnification, and amendment
procedures. The IRS asks whether you have adopted bylaws.

### Conflict-of-interest (COI) policy
Strongly expected. The IRS Form 1023 process asks whether you have adopted a conflict-of-interest policy, and
**Pub. 557 / the Form 1023 instructions include a sample COI policy** you can adapt (irs.gov). It governs how
directors disclose and recuse from transactions in which they have a financial interest. Adopt it at the
first meeting.

### Organizational meeting
After incorporation, hold the **first (organizational) meeting of the board** to: adopt the bylaws, adopt the
COI policy, elect officers, authorize opening a bank account, authorize applying for an EIN and 501(c)(3)
status, and set the fiscal year. **Keep minutes** — they are part of your permanent records and evidence of
proper formation.

---

## 5. Get an EIN (Employer Identification Number)

Every nonprofit needs an **EIN**, even with no employees; it is required for banking, the 501(c)(3)
application, and tax filings (IRS, "Obtaining an EIN for an exempt organization," irs.gov, as of 2026-06).

- **Apply only after the entity is formed with the state.** The IRS specifically warns that if you do not form
  your entity with the state first, your EIN application may be delayed (irs.gov, as of 2026-06).
- **Apply with Form SS-4** (Rev. December 2025). Online application is **free** and issues the EIN
  immediately for U.S. applicants; fax and mail are also available (irs.gov, as of 2026-06).
- Getting an EIN does **not** make you tax-exempt — exemption requires the separate Form 1023/1023-EZ step.

---

## 6. Federal 501(c)(3) recognition (Form 1023 vs. 1023-EZ)

Apply to the IRS for recognition of exemption. Two forms, chosen by eligibility:

### Form 1023-EZ (Streamlined) — if you qualify
- **User fee: $275** (IRS, "Form 1023 and 1023-EZ: Amount of user fee," irs.gov, page reviewed 09-Oct-2025;
  amounts are subject to change — verify on Pay.gov before filing).
- **Eligibility (high level):** annual gross receipts **not exceeding $50,000** in any of the past three
  years and not projected to exceed $50,000 in any of the next three years, **and** total assets with fair
  market value **not more than $250,000** (per the IRS 1023-EZ eligibility framework, irs.gov, as of 2026-06).
- You must complete the **Eligibility Worksheet** in the Instructions for Form 1023-EZ (Rev. January 2025);
  a "Yes" to any worksheet question disqualifies you and you must file the full Form 1023.
- **Ineligible categories** include churches, schools, and hospitals, private operating foundations, and LLCs
  (Instructions for Form 1023-EZ, Rev. Jan. 2025, irs.gov).
- **Filed electronically only**, via IRS.gov/Form1023EZ on Pay.gov — no paper.

> For a brand-new, lean organ-donation-awareness nonprofit projecting modest early revenue, **1023-EZ is
> often the right path** — but project receipts honestly; understating to qualify is risky.

### Form 1023 (long form) — if you do not qualify for EZ, or want a stronger record
- **User fee: $600** (IRS, "Form 1023 and 1023-EZ: Amount of user fee," irs.gov, reviewed 09-Oct-2025;
  subject to change).
- Required if you expect to exceed the EZ thresholds, or are a church/school/hospital/private foundation/LLC,
  or want the more detailed determination (narrative of activities, financial data, schedules).
- Filed electronically on Pay.gov.

### What 501(c)(3) gives you and requires
- **Exemption from federal income tax** on exempt-purpose activities, and **deductibility of donations** to
  you under IRC §170.
- In exchange you must pass the **organizational test** (Section 3 above) and the **operational test** — be
  operated exclusively for exempt purposes, with no private inurement and no political-campaign activity
  (Instructions for Form 1023-EZ, Rev. Jan. 2025, irs.gov).

### Public charity vs. private foundation (and the public support test)
Every 501(c)(3) is **presumed a private foundation unless it qualifies as a public charity** (IRS,
"Determine your foundation classification," irs.gov, as of 2026-06). Public-charity status is generally
preferable (fewer excise taxes and restrictions, better donor optics). The two common public-charity tests:

- **§509(a)(1) / §170(b)(1)(A)(vi):** generally must receive **at least one-third (33⅓%) of its support from
  the general public / governmental units**, or meet a 10%-plus facts-and-circumstances test.
- **§509(a)(2):** generally must receive **more than one-third of support from contributions and gross
  receipts from exempt-purpose activities**, **and no more than one-third from gross investment income +
  unrelated business taxable income.**
- Both are measured over a **five-year support period** (IRS, "Public charity support test," irs.gov, as of
  2026-06).

A donation-driven awareness nonprofit usually aims for **§509(a)(1)** public-charity status by keeping
broad-based public support. You indicate your requested classification on the exemption application.

---

## 7. NC state tax treatment (NCDOR)

Federal 501(c)(3) status does **not** automatically grant NC tax exemptions — handle these with NCDOR
(ncdor.gov) separately.

### Corporate income & franchise tax exemption
- NCDOR **does not issue "exempt numbers"** for franchise/corporate income tax, but it **issues letters of
  tax exemption** (NCDOR, "Nonprofit Corporate Tax Information," ncdor.gov, as of 2026-06).
- For a nonprofit organized under **GS Chapter 55A**, there is **no separate application or fee** required to
  obtain the state franchise/income exemption — the process flows from incorporating with the NC Secretary of
  State (NCDOR, ncdor.gov, as of 2026-06).
- **Caution:** a nonprofit corporation that does **not** request/obtain the tax-exempt letter is treated as
  **subject to** franchise and corporate income tax. So request the exemption letter from NCDOR rather than
  assuming it. Confirm the current request procedure with NCDOR.

### NC sales-and-use tax — refund mechanism (not an upfront exemption)
NC generally does **not** exempt qualifying nonprofits from *paying* sales tax at purchase; instead it gives a
**refund** of sales-and-use tax paid on direct purchases used in the nonprofit's work (NCDOR, "Nonprofit Sales
and Use Tax Information," ncdor.gov, as of 2026-06):

1. **Register first** with NCDOR using **Form E-585NPA**, Application for Nonprofit Sales and Use Tax Refund
   Account ID, to get a refund-claim account ID (ncdor.gov, as of 2026-06).
2. **File refund claims on Form E-585** ("Nonprofit and Governmental Entity Claim for Refund State, County,
   and Transit Sales and Use Taxes"), under the authority of **GS §105-164.14(b)** (ncdor.gov, as of 2026-06).
3. Eligible nonprofits get a **semiannual refund** of NC sales-and-use tax on qualifying direct purchases.
4. [UNVERIFIED — any annual cap on refundable amounts and the exact claim deadlines should be confirmed on the
   current Form E-585 instructions, ncdor.gov; specific cap figures were not extracted this pass.]

> Practical note: keep clean purchase records — the refund mechanism is documentation-driven.

---

## 8. NC Charitable Solicitation License (GS Chapter 131F)

Administered by the **NC Secretary of State, Charitable Solicitation Licensing Section** (sosnc.gov). This is
**separate from incorporation and from 501(c)(3)** — being a recognized 501(c)(3) does **not** exempt you from
needing this license.

### Who must get a license
Any organization (or person) that intends to **directly solicit charitable contributions in NC, or hires
someone to solicit on its behalf**, must obtain an annual license **before soliciting**, unless an exemption
applies (sosnc.gov; GS Chapter 131F, ncleg.gov, as of 2026-06). For an organ-donation venture that plans to
ask the NC public for donations, **assume you need this license unless you clearly fit an exemption.**

### The key small-organization exemption (likely relevant at launch)
Under **GS §131F-3**, exempt persons include **"any person who receives less than fifty thousand dollars
($50,000) in contributions in any calendar year and does not provide compensation to any officer, trustee,
organizer, incorporator, fund-raiser, or solicitor"** (ncleg.gov; sosnc.gov, as of 2026-06). So an
**all-volunteer organization under $50,000/year in contributions is exempt** — but pay anyone in those roles,
or cross $50,000, and the exemption is lost.

### Other statutory exemptions (GS §131F-3)
Religious institutions; federal/state/local government; accredited educational institutions and their
related foundations; licensed hospitals and related foundations; certain noncommercial radio/TV stations;
qualified community trusts; bona fide volunteers/salaried officers; advising attorneys/investment
counselors/bankers; volunteer fire departments, REACT teams, rescue squads, and EMS; YMCAs/YWCAs; nonprofit
continuing-care facilities; and certain tax-exempt nonprofit fire/EMS organizations (sosnc.gov exemption
list; GS §131F-3, ncleg.gov, as of 2026-06).

### How to claim an exemption
Even when exempt, NC expects you to **submit a written request for exemption with supporting documentation**
to the Secretary of State; on approval you receive a **letter of exemption**, and the Department may review or
cancel exemptions at any time (sosnc.gov, as of 2026-06). Don't assume — file the exemption request.

### License fees (sliding scale by prior-year contributions)
Set by statute on a **sliding scale tied to contributions received per fiscal year**, with a fee waiver in
some cases (sosnc.gov fee schedule, as of 2026-06):

- Contributions **under $50,000 with no compensated** officers/fundraisers → **exempt** (no license required).
- Contributions **under $5,000** (with compensated personnel) → **license required but no fee**.
- **$5,000–$49,999** → **$50**.
- **$50,000–$99,999** → **$50**.
- **$100,000–$199,999** → **$100**.
- **$200,000 and above** → **$200**.
- A **parent organization filing on behalf of its chapters** has a statutory maximum around **$400**
  (sosnc.gov; GS Chapter 131F Article 2, ncleg.gov, as of 2026-06).

[UNVERIFIED — these bracket amounts come from SOS-sourced summaries; confirm the exact current brackets and
fees on the live sosnc.gov fee schedule and the current license application before filing, as the schedule
is statute-driven and periodically amended (e.g., Session Law 2023-119 modified Chapter 131F licensing).]

### Professional fundraisers / solicitors
If you hire a **professional fundraising consultant or professional solicitor**, they must be **separately
licensed** — fee reported as **$200** each, renewed annually (sosnc.gov, as of 2026-06). Their contracts and
solicitation activity are regulated under Chapter 131F; this is a common compliance trap when outsourcing
fundraising.

### Renewal
The license must be **renewed annually**. Renewal is due by the close of each fiscal year in which you
solicited in NC, **or** by the date of any applicable federal information-return (Form 990) filing extension,
whichever is later; the Department may grant an extension up to **60 days** tied to a federal filing extension
(GS Chapter 131F, ncleg.gov; sosnc.gov, as of 2026-06).

> **Sequencing for the venture:** if launch-year contributions will stay under $50,000 and everyone is a
> volunteer, file the **exemption request** now and revisit before you hire anyone or cross $50,000. If you
> will pay staff or expect bigger giving, get the **license before your first public solicitation.**

---

## 9. Ongoing compliance

### IRS — annual Form 990 series (filing keeps your exemption alive)
Tax-exempt organizations file an annual return; **which version depends on size** (IRS, "Annual electronic
filing requirement…Form 990-N (e-Postcard)," irs.gov, as of 2026-06):

- **Form 990-N (e-Postcard):** for small orgs whose **gross receipts are normally $50,000 or less.**
  "Normally $50,000 or less" includes: in existence 1 year or less and received/pledged **$75,000 or less** in
  year one; or at least 3 years old and **averaged $50,000 or less** over the prior 3 years. Electronic only.
- **Form 990-EZ:** if **gross receipts < $200,000 and total assets < $500,000** at year end.
- **Form 990 (full):** above those thresholds.
- **Deadline:** the **15th day of the 5th month after the close of your tax year** (e.g., May 15 for a
  calendar-year org).
- **Critical:** failing to file for **three consecutive years triggers automatic revocation** of tax-exempt
  status (irs.gov, as of 2026-06). This is the single most common way small nonprofits lose their status.

### NC Secretary of State — annual report
**Nonprofit corporations formed under GS Chapter 55A currently do NOT file an annual report** with the NC
Secretary of State (sosnc.gov, as of 2026-06) — unlike business corporations and LLCs. **Watch this:**
legislation has been proposed that would impose a nonprofit annual-report requirement; confirm current status
before relying on the exemption. [UNVERIFIED — proposed bills (e.g., 2025 session) were noted in research but
not enacted as of this pass; re-check sosnc.gov.]

### NC charitable license — annual renewal
Renew the Charitable Solicitation License (or refresh the exemption letter) **every year** as in Section 8.
Many states (including NC) require submitting **financial information** with the renewal.

### Recordkeeping (keep permanently or for the statutory retention period)
- Articles, bylaws, COI policy, board/member meeting minutes, and the IRS determination letter — **permanent**.
- Financial records, Form 990 filings, donation records, and sales-tax-refund documentation — retain per IRS
  and NCDOR guidance.
- **Public disclosure:** 501(c)(3)s must make their **Form 1023 application and recent Form 990s** available
  for public inspection (IRS rules, irs.gov). Build this into your records hygiene.

---

## 10. Multistate charitable registration (if you solicit beyond NC)

Charitable-solicitation laws are **state-by-state**. The moment your organ-donation venture solicits
donations from residents of other states (including via a national website, email appeals, or social
fundraising), those states may require their **own registration**; roughly **40+ states** regulate
charitable solicitation.

- The **Unified Registration Statement (URS)** was a 1990s effort (by NASCO and NAAG) to create a single
  multistate form. In practice its utility is now **limited**: it has not been substantively updated since
  ~2010, most states want **state-specific online filings**, and only a couple of states (e.g., Kentucky,
  Louisiana) treat it as a primary method (multistatefiling.org; industry sources, as of 2026-06).
- **Practical reality:** most growing nonprofits file **state-specific registrations** (often via a
  compliance vendor) rather than relying on the URS.
- **Online-fundraising nuance:** soliciting through a public website can trigger registration in states under
  the **"Charleston Principles"** framework, especially where you specifically target or repeatedly receive
  contributions from a state's residents. [UNVERIFIED — the Charleston Principles are non-binding guidance
  adopted by NASCO, not law in every state; get state-specific advice before a national campaign.]

> For launch, **register/clear NC first.** Add other states deliberately as you begin actively soliciting
> their residents, with counsel or a registration service.

---

## Quick reference (verify all figures before acting; as of 2026-06)

| Item | Value | Primary source |
|---|---|---|
| NC nonprofit Articles of Incorporation form | Form N-01, under GS §55A-2-02 | sosnc.gov; ncleg.gov |
| NC incorporation filing fee | **$60** [UNVERIFIED — confirm] | sosnc.gov fee schedule |
| EIN application (Form SS-4, online) | **Free** | irs.gov |
| IRS Form 1023-EZ user fee | **$275** | irs.gov (rev. 09-Oct-2025) |
| Form 1023-EZ eligibility | ≤ $50,000 gross receipts (3-yr look-back/forward) **and** ≤ $250,000 total assets | irs.gov |
| IRS Form 1023 (long form) user fee | **$600** | irs.gov (rev. 09-Oct-2025) |
| Public-charity support test | ≥ 33⅓% public support (§509(a)(1)) over 5 yrs | irs.gov |
| NC corporate income/franchise exemption | Letter from NCDOR; no fee for Ch. 55A nonprofits | ncdor.gov |
| NC sales/use tax | Refund via Form E-585 (register on E-585NPA), GS §105-164.14(b) | ncdor.gov |
| NC charitable license small-org exemption | < **$50,000**/yr contributions **and** no compensated officers/fundraisers/solicitors | GS §131F-3, ncleg.gov |
| NC charitable license fees | $50 / $100 / $200 sliding scale (≈$400 parent max) | sosnc.gov; GS Ch. 131F Art. 2 |
| Professional solicitor/consultant license | **$200** each, annual | sosnc.gov |
| IRS Form 990-N threshold | gross receipts **normally ≤ $50,000** | irs.gov |
| IRS Form 990-EZ threshold | gross receipts < **$200,000** and assets < **$500,000** | irs.gov |
| Auto-revocation trigger | **3 consecutive years** of non-filing | irs.gov |
| NC SOS nonprofit annual report | **Not required** for Ch. 55A nonprofits (proposal pending) | sosnc.gov |

---

## Sources

Primary sources consulted (as of 2026-06):

- **NC General Statutes, Chapter 55A:** North Carolina Nonprofit Corporation Act (incl. §55A-2-02
  incorporation): https://www.ncleg.gov/Laws/GeneralStatuteSections/Chapter55A and
  https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/ByChapter/Chapter_55a.pdf
- **NC General Statutes, Chapter 131F:** Charitable Solicitation (incl. §131F-3 exemptions, §131F-5):
  https://www.ncleg.gov/enactedlegislation/statutes/html/bychapter/chapter_131f.html and
  https://www.ncleg.gov/enactedlegislation/statutes/pdf/bysection/chapter_131f/gs_131f-3.pdf
- **Session Law 2023-119:** modifications to Chapter 131F charitable-solicitation licensing:
  https://www.ncleg.gov/EnactedLegislation/SessionLaws/PDF/2023-2024/SL2023-119.pdf
- **NC Secretary of State, Charitable Solicitation Licensing** (who must register, fees, renewal, exemptions):
  https://www.sosnc.gov/divisions/charities/licensing ,
  https://www.sosnc.gov/divisions/charities/exemptions ,
  https://www.sosnc.gov/manual/launching_a_business/nonprofit_solicitation ,
  https://www.sosnc.gov/fees/by_title/_charities_charities_Sponsors
- **NC Secretary of State, Nonprofit Articles of Incorporation (Form N-01)** and Articles of Amendment
  (Form N-02): https://www.sosnc.gov/webfiles/documents/forms/Business_Registration/nonprofit_corporations/articles_of_incorporation_for_nonprofit.pdf
- **NC Secretary of State, Annual report applicability** (nonprofits not required to file):
  https://www.sosnc.gov/divisions/business_registration/annual_report_help
- **NCDOR, Nonprofit Corporate Tax Information** (income/franchise exemption letter):
  https://www.ncdor.gov/taxes-forms/corporate-income-franchise-tax/nonprofit-corporate-tax-information
- **NCDOR, Nonprofit Sales and Use Tax Information** and Form E-585 / E-585NPA (refund mechanism, GS
  §105-164.14(b)): https://www.ncdor.gov/taxes-forms/sales-and-use-tax/other-sales-and-use-tax-resources/nonprofit-sales-and-use-tax-information
  and https://www.ncdor.gov/form-e-585npa-application-nonprofit-sales-and-use-tax-refund-account-id/open
- **IRS, Instructions for Form 1023-EZ (Rev. January 2025):** https://www.irs.gov/instructions/i1023ez
- **IRS, Form 1023 and 1023-EZ: Amount of user fee** ($600 / $275; reviewed 09-Oct-2025):
  https://www.irs.gov/charities-non-profits/form-1023-and-1023-ez-amount-of-user-fee
- **IRS, Obtaining an EIN for an exempt organization** and Form SS-4 (Rev. Dec. 2025):
  https://www.irs.gov/charities-non-profits/obtaining-an-employer-identification-number-for-an-exempt-organization
  and https://www.irs.gov/instructions/iss4
- **IRS, Determine your foundation classification** and **Public charity support test** (509(a)(1)/(2)):
  https://www.irs.gov/charities-non-profits/determine-your-foundation-classification and
  https://www.irs.gov/charities-non-profits/charitable-organizations/exempt-organizations-annual-reporting-requirements-form-990-schedules-a-and-b-organization-may-be-publicly-supported-under-one-of-two-tests
- **IRS, Form 990-N (e-Postcard) filing requirement** (thresholds, deadline, auto-revocation):
  https://www.irs.gov/charities-non-profits/annual-electronic-filing-requirement-for-small-exempt-organizations-form-990-n-e-postcard
- **IRS, Publication 557, Tax-Exempt Status for Your Organization** (sample purpose/dissolution/COI language):
  https://www.irs.gov/forms-pubs (search Pub. 557)
- **Unified Registration Statement (multistate charitable registration)**: https://multistatefiling.org/

---

## Disclaimer

This skill is **general information only and is not legal, tax, or accounting advice.** It does not create an
attorney-client or any professional relationship. Statutes, regulations, forms, form numbers, fees, and
dollar thresholds **change frequently**, and several figures here are flagged **[UNVERIFIED]** because they
could not be confirmed against the live primary source in this research pass (notably the exact NC
incorporation filing fee, the precise charitable-license fee brackets, any sales-tax refund cap, the GS §55A
director minimum, and the status of any pending NC nonprofit annual-report legislation). Before acting, the
founder should **verify every fee, form, threshold, and requirement directly** with the **NC Secretary of
State** (sosnc.gov), the **NC Department of Revenue** (ncdor.gov), and the **IRS** (irs.gov), and should
**consult a licensed North Carolina attorney and a CPA** experienced with nonprofits and 501(c)(3) compliance.
Getting the Articles' 501(c)(3) language, the public-charity classification, and the solicitation-licensing
posture right at the outset is far cheaper than fixing them later.
