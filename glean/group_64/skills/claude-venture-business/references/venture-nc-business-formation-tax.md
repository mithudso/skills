<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-nc-business-formation-tax` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-nc-business-formation-tax
description: >-
  Forming a FOR-PROFIT North Carolina business and keeping it tax- and filing-compliant, for a solo
  founder: entity choice and liability-vs-tax tradeoffs (sole proprietorship, LLC, PLLC, C-corp, the S-corp
  tax election), NC formation steps (Articles of Organization, registered agent, EIN), assumed-name/DBA at
  the county Register of Deeds, licenses (NC has NO general business license), NCDOR tax
  (sales & use, withholding, corporate income, franchise), federal owner taxes (self-employment tax,
  quarterly estimates), the Corporate Transparency Act / FinCEN BOI status, and the annual report.
  TRIGGER: forming a for-profit NC entity; LLC vs S-corp vs C-corp; DBA; NC business
  tax; CTA / BOI / FinCEN; NC annual report; EIN; registered agent. SKIP: nonprofit / 501(c)(3) / Form 990
  -> venture-nc-nonprofit-formation; business plan / financials / funding -> venture-small-business-planning;
  real-estate brokerage / agency licensing -> venture-nc-real-estate-law; consumer-credit / debt law.
category: personal-venture
tags: [venture, north-carolina, business-formation, llc, tax, compliance]
version: 1.1.0
updated: 2026-06-16
metadata:
  changelog:
    - "2026-06-16 sko v1.0.0->v1.1.0 — Pass H 10/10 pos, 10/10 neg (predicted); fixed 1 High (Pass M: description 1372->987 chars under Glean cap) + 3 Medium (Pass A: franchise-tax C-corp vs S-corp precision; Pass A/D: 'minimal records'->'proper records' in pierce-the-veil note; Pass K: em-dash density 2.22->0.94/100); 0 banned terms; first-time hub registration"
whenToUse:
  - Choosing a for-profit entity in NC (sole proprietorship, LLC, PLLC, C-corp) and weighing liability vs. tax
  - Deciding whether to make the S-corp tax election and how reasonable-compensation works
  - Filing Articles of Organization with the NC Secretary of State and setting up a registered agent, operating agreement, and EIN
  - Filing a DBA / Certificate of Assumed Business Name at the county Register of Deeds
  - Figuring out which licenses or permits apply when NC has no general statewide business license
  - Registering with NCDOR for sales & use tax, employer withholding, or corporate income/franchise tax
  - Planning federal owner taxes: self-employment tax and quarterly estimated payments
  - Checking the current Corporate Transparency Act / FinCEN BOI reporting obligation
  - Filing the annual report and opening a business bank account
triggers:
  - NC LLC formation
  - LLC vs S-corp vs C-corp
  - NC assumed business name DBA
  - NC business sales tax registration
  - corporate transparency act BOI FinCEN
  - NC annual report
  - S-corp election reasonable compensation
  - NC franchise tax
---

# Forming a For-Profit North Carolina Business and Keeping It Tax-Compliant

This skill walks a solo or small founder through legally standing up a **taxable** North Carolina business
and keeping it compliant across the regulators that touch a for-profit: the **NC Secretary of State**
(entity formation + annual report), the **county Register of Deeds** (assumed-name/DBA), the **NC
Department of Revenue / NCDOR** (state taxes), the **IRS** (EIN, federal income + self-employment tax, the
S-corp election), and **FinCEN** (Corporate Transparency Act beneficial-ownership reporting).

For a NONPROFIT (501(c)(3), charitable-solicitation license, Form 990), use **venture-nc-nonprofit-formation**
— do not use this skill. For the **business plan, financial model, unit economics, and funding**, use
**venture-small-business-planning**. This skill is purely the entity-formation + tax/compliance ops layer
for a for-profit.

Keep three lanes mentally separate, because people lump them together:

1. **Entity formation**: what legal form the business takes (or whether you bother forming an entity at all).
2. **Tax treatment**: how the IRS and NCDOR tax the business and you. *Tax treatment is a separate choice
   from legal form*; an LLC can be taxed as a sole proprietorship, partnership, S-corp, or C-corp.
3. **Ongoing compliance**: registrations, returns, annual report, and federal BOI status that keep you in
   good standing.

> **One thing to internalize early:** "LLC" is a *legal* status; "S-corp" is a *tax* status. You do not
> "form an S-corp." You form an LLC (or corporation) and then *elect* S-corp taxation. This trips up almost
> every first-time founder.

---

## 1. Entity choice: liability vs. tax tradeoffs

Five practical options for a solo/small NC founder. The first decision is **liability** (does a lawsuit or
debt reach your personal house and savings?); the second is **tax** (how is the income taxed, and can you
cut self-employment tax?).

### The entity-comparison table

| Entity | Liability shield | How it's taxed by default | NC formation filing | Best fit for a solo founder |
|---|---|---|---|---|
| **Sole proprietorship** | **None** — you *are* the business | Pass-through on your personal 1040 (Schedule C); all net profit hit by self-employment tax | None (just a DBA if using a trade name) | Lowest-risk side hustle, testing an idea, no employees, minimal liability exposure |
| **General partnership** | **None** for the partners | Pass-through (Form 1065 + K-1s) | None (DBA if using a trade name) | Two+ owners testing something informally (rare to recommend — no shield) |
| **NC LLC** | **Yes** (limited liability) | Default: single-member = disregarded (Schedule C); multi-member = partnership. *Can elect* S-corp or C-corp | Articles of Organization, **$125** [GS § 57D-1-22] | The default recommendation for most solo/small founders — shield + flexible taxation |
| **NC PLLC** | Yes | Same as LLC | Articles of Organization (Form PLLC-02) + licensing-board certification | A **licensed professional** (doctor, lawyer, CPA, engineer, architect, etc.) — see §1.1 |
| **C-corporation** | Yes | **Entity-level** corporate tax, then dividends taxed again ("double taxation") | Articles of Incorporation | Raising VC/angel money, issuing stock, QSBS, wanting to retain earnings in the entity |

**S-corporation is not in this table because it is a tax election, not an entity** — see §6 and the federal
section.

**The usual answer for a solo NC founder is an LLC**, because it gives a liability shield with minimal
formality and lets you keep simple pass-through taxation now and *add* the S-corp election later once profit
justifies it (§6). A sole proprietorship is cheaper and simpler but gives you **zero** liability protection —
fine for a tiny, low-risk venture, risky once you have customers, contracts, or assets to lose. Reserve the
C-corp for when you genuinely need outside equity investment or the specific C-corp tax features; for a
bootstrapped solo business the double taxation is usually a drawback. (For the funding-strategy side of this
— bootstrapping vs. angel/VC — see **venture-small-business-planning**.)

### 1.1 PLLC: for licensed professions only

If your business practices a **licensed profession** in NC (medicine, law, accounting, engineering,
architecture, etc.), you generally cannot use an ordinary LLC — you must form a **Professional LLC (PLLC)**
under **NC GS Chapter 55B and Chapter 57D** [sosnc.gov, PLLC requirements]. Key differences:

- The name must end in "Professional Limited Liability Company" or **P.L.L.C.** [sosnc.gov].
- You file the **PLLC Articles of Organization (Form PLLC-02)**, not the standard LLC form [sosnc.gov forms].
- You need a **certification from the applicable licensing board** that ownership complies with Ch. 55B/57D,
  filed with or supporting the Articles [sosnc.gov].
- **Ownership rule:** at least **two-thirds** of the ownership interests must be held by licensees of the
  profession, with at least one NC-licensed member/manager per profession offered; up to one-third may be
  owned by non-licensed employees [sosnc.gov, PLLC requirements] — *confirm the exact current ratio with
  your licensing board, as boards set their own rules*.
- Get **licensing-board approval first**, then file with the Secretary of State.

A PLLC shields you from ordinary business liabilities but **does not** shield you from your own professional
malpractice; that's what malpractice insurance is for.

---

## 2. NC formation steps (LLC walkthrough)

The mechanics for the most common case, forming a single-member or small NC LLC. (A corporation is similar
but files **Articles of Incorporation** and adopts bylaws + issues stock instead of an operating agreement.)

1. **Pick and clear a name.** It must be distinguishable from other NC entities and include an LLC
   designator ("LLC," "L.L.C.," "Limited Liability Company"). Search availability on the **NC Secretary of
   State** business-registration search at **sosnc.gov** before filing.
2. **Appoint a registered agent.** Every NC LLC must have a registered agent with a **physical NC street
   address** (not a PO box) available during business hours to receive legal/state notices [sosnc.gov]. You
   can be your own agent if you have an NC street address, or hire a commercial registered-agent service
   (commonly ~$100-150/yr) — useful for privacy and if you work from home or travel.
3. **File the Articles of Organization** with the NC Secretary of State. Filing fee **$125** [GS § 57D-1-22;
   confirm at sosnc.gov]. Online and paper carry the **same base fee**; online just processes faster (paper
   can take weeks). Articles require: LLC name, registered agent + NC address, principal-office address, and
   organizer information [sosnc.gov]. NC LLCs are **member-managed by default**; you can elect
   **manager-managed** structure (members appoint managers to run day-to-day operations).
4. **Get an EIN from the IRS.** Free, directly from **irs.gov** (never pay a third party for one). A
   single-member LLC with **no employees** can technically use the owner's SSN, but get an EIN anyway: you
   need it to open a business bank account and it keeps your SSN off forms. A **multi-member LLC, any LLC with
   employees, and any S-corp/C-corp election requires an EIN** [irs.gov].
5. **Adopt an operating agreement.** **Not required by NC law**, but strongly recommended even for a
   single-member LLC [sosnc.gov]. It documents ownership, capital contributions, management, profit
   distributions, and member-add/-remove procedures, and helps preserve the liability shield by showing the
   LLC is a real, separate entity. **It is NOT filed with the Secretary of State**; keep it in your records.
6. **Handle post-formation items:** open a dedicated business bank account (§9), register with NCDOR for any
   applicable taxes (§5), check licenses/permits (§4), file the DBA if using a trade name (§3), and confirm
   your federal BOI status (§7).

> **Pierce-the-veil warning:** the liability shield only holds if you treat the LLC as separate from
> yourself: a separate bank account, no commingling of personal and business funds, signing contracts in the
> LLC's name, and keeping proper records. Sloppy separation lets a court "pierce the veil" and reach you
> personally.

---

## 3. Assumed business name / DBA (Certificate of Assumed Business Name)

A **DBA** ("doing business as," also called an assumed name, trade name, or fictitious name) lets you operate
under a name **different from your legal name** (sole proprietor) or your **registered entity name** (LLC/corp).
Example: you formed "Hudson Ventures LLC" but want to brand as "Carolina Coffee Roasters."

- **You need a DBA when** a sole proprietor uses anything other than their own legal name, or an LLC/corp
  uses a brand different from its registered legal name. You do **not** need one if you operate under your
  exact legal name or your exact registered entity name.
- **Where to file:** NC DBAs are filed as a **Certificate of Assumed Business Name** with the **county
  Register of Deeds**, *not* the Secretary of State [NC GS Chapter 66, Article 14A].
- **The 2017 centralized system:** Under **Session Law 2017-23 (House Bill 228)**, effective **December 1,
  2017**, NC modernized the assumed-name system. You file **one** Certificate of Assumed Business Name with a
  single county Register of Deeds, and that filing is transmitted into a **statewide central database**
  maintained at the state level — so a single filing covers the whole state (you no longer file in every
  county where you do business) [NC GS § 66-71.4; county Register of Deeds]. One certificate can list up to
  five assumed names and the counties of operation.
- **Cost:** set by the county (commonly around **$26**, plus a small fee per extra page) — *verify the exact
  fee with your county Register of Deeds*. [UNVERIFIED — fee varies by county]
- **A DBA is not an entity and gives no liability protection.** It is only a public name registration.

---

## 4. Licenses & permits: NC has NO general statewide business license

This is one of the most misunderstood points, so state it plainly:

> **North Carolina does not have a single, general, statewide "business license."** There is no one form or
> fee that licenses "a business" in NC. [EDPNC / BLNC; multiple NC sources]

Licensing in NC is **industry-specific and locality-specific** instead:

- **Privilege-license history (why there's no general license):** NC eliminated municipal authority to levy
  local **privilege license taxes** on most businesses effective **July 1, 2015** (Session Law 2015-141)
  [Tax Foundation; NCDOR]. Separately, the statewide **professional privilege license** that NCDOR
  administered under GS § 105-41 (lawyers, physicians, and certain other professions) was **repealed by
  Session Law 2023-134, effective July 1, 2024** — those professionals **no longer** pay that state privilege
  license [NCDOR]. So both the old local and the old statewide professional privilege licenses are gone.
- **What you may still need (industry/locality-specific):**
  - **Occupational/professional licenses** from a state licensing board (contractors, cosmetologists, real-
    estate brokers, electricians, health professionals, etc.) — required to *practice the profession*.
  - **Sales & use tax registration / Certificate of Registration** from NCDOR if you sell taxable goods or
    certain services (§5); this is a tax registration, not a "business license," but it functions like one
    for sellers.
  - **Local permits** — zoning/land-use, a county/city **privilege license is largely gone** but local
    **regulatory permits** (signage, building, fire, health-department food permits, home-occupation permits)
    can still apply. Check your **city and county**.
  - **Industry-specific federal/state permits** — alcohol (NC ABC), firearms, food/restaurants (NC DHHS),
    childcare, transportation, etc.
- **The resource to use:** **Business Link North Carolina (BLNC)**, a free service of the **Economic
  Development Partnership of North Carolina (EDPNC)**, gives **customized license-and-permit lists** and
  free one-on-one counseling. Website **blnc.gov** (and edpnc.com), toll-free **1-800-228-8443**, weekdays
  8 a.m.-5 p.m. [EDPNC / BLNC]. This is the single best starting point for "what licenses do I actually need?"

---

## 5. NC tax registration & ongoing state tax (NCDOR)

Register with the **NC Department of Revenue (NCDOR)** for whatever taxes apply. The main intake is **Form
NC-BR (Business Registration Application)**, which can register you for **sales & use tax** and **income-tax
withholding** in one application; you can also register online at **ncdor.gov** [NCDOR].

### 5.1 Sales & use tax

- **When it applies:** if you sell **taxable tangible goods** or certain **taxable services** in NC, you must
  register, collect, and remit sales tax. Many pure-service businesses are not subject to sales tax — *check
  whether your specific service is taxable* (NC taxes some services, e.g., certain repair, maintenance, and
  installation services).
- **Rate:** the **state rate is 4.75%**, plus **local and transit** rates, so the **combined rate is roughly
  6.75% to 7.5%** depending on the county [NCDOR, Current Sales and Use Tax Rates]. *Look up your county's
  combined rate on NCDOR.*
- **Economic nexus (remote/online sellers):** an out-of-state seller must register and collect once **gross
  sales sourced to NC exceed $100,000** in the current or previous calendar year (NC repealed the prior
  200-transaction prong) [NCDOR / NC sales-tax guidance]. *Verify the current threshold with NCDOR if
  selling across state lines.*
- **Filing:** NCDOR assigns a filing frequency (monthly/quarterly) based on volume; remit on the assigned
  schedule.

### 5.2 Employer withholding (only if you have employees)

If you pay **employees** (W-2), register for **NC income-tax withholding** via **Form NC-BR**, have each
employee complete **Form NC-4** (state withholding allowances), withhold NC income tax from wages, and file
withholding returns (**Form NC-5 / NC-5P**) on the schedule NCDOR assigns [NCDOR]. You'll also have federal
payroll obligations (FICA, FUTA, Form 941) and **NC unemployment insurance** (register with the **NC Division
of Employment Security**). **Paying yourself as an S-corp owner-employee means you become an employer** with
withholding and payroll obligations (§6).

### 5.3 NC corporate income & franchise tax (C-corps; franchise also S-corps)

These apply to entities taxed as **corporations**. Pass-through entities (sole prop, partnership, default
LLC, S-corp) generally do **not** pay NC corporate income tax; their income flows to the owners' personal NC
returns (§5.4). But note S-corps and C-corps both owe **NC franchise tax**.

- **NC corporate income tax — being phased out to zero:** the rate is on a statutory glide path:
  **2.25% (2025) → 2.0% (2026) → 1.0% (2028) → 0% (2030)** [NCDOR; PwC; BDO]. (The 2.5% rate applied through
  2024.) So for tax year **2026** a C-corp pays **2.0%** on NC taxable income, and the tax disappears entirely
  in 2030.
- **NC franchise tax (C-corps and S-corps):** a tax on the corporation's net worth/capital base, **minimum
  $200**, owed even in a no-profit year (hence the $200 floor). The structure differs slightly by entity
  [NCDOR, effective for tax years on/after Jan 1, 2025]:
  - **C-corporation:** **$1.50 per $1,000** of the tax base, but capped at **$500 on the first $1,000,000**;
    the base **above $1,000,000** continues to be taxed at $1.50 per $1,000.
  - **S-corporation:** a flat **$200 on the first $1,000,000** of the tax base, then **$1.50 per $1,000** on
    the base **above $1,000,000**.
  *The franchise tax is NOT being phased out — only the corporate income tax is.*

### 5.4 Pass-through treatment & the NC individual income-tax rate

- **Pass-through entities** (sole proprietorship, partnership, default-taxed LLC, S-corp) pay no separate NC
  *income* tax at the entity level; the profit "passes through" to the owners, who report it on their
  **personal NC income-tax return (Form D-400)** and pay NC individual income tax. (S-corps still owe NC
  franchise tax per §5.3, and NC offers an elective **pass-through-entity (PTET) "Taxed PTE"** option some
  owners use to work around the federal SALT cap — *ask a CPA whether the NC PTET election helps you*.)
- **NC has a FLAT individual income-tax rate, trending down:** **4.25% (2025) → 3.99% (2026)**, then
  scheduled to drop to **3.49% (2027)** and **2.99% (2028)** *if revenue benchmarks are met* [NCDOR, Tax Rate
  Schedules; NC OSBM]. (It was 4.5% in 2024.) So your NC tax on pass-through business profit for **2026** is
  **3.99%** of your NC taxable income, flat, regardless of bracket.

---

## 6. The S-corp tax election: when it fits a small NC founder

This is the single biggest *tax-planning* lever for a profitable solo business, and the most over-applied —
so understand the mechanics, not just the headline.

**What it is:** you keep your LLC (or corporation) as the legal entity but file **IRS Form 2553** to elect
**S-corporation taxation** [irs.gov]. The business stays a pass-through (no entity income tax), but the way
*your* income is taxed changes.

**Why founders do it — the self-employment-tax angle:** as a sole proprietor or default LLC, **all** your net
profit is hit by **15.3% self-employment (SE) tax** (§8). As an S-corp owner-employee, you split the profit
into two buckets:

- **A reasonable salary** (W-2 wages to yourself) — subject to payroll/FICA taxes (the employer+employee
  equivalent of SE tax).
- **Distributions** of the remaining profit are **not** subject to SE/payroll tax [irs.gov; widely cited].

The distributions portion escapes the ~15.3% SE tax, which is the savings. Example pattern (illustrative,
not advice): on $150k profit, a $70k reasonable salary + $80k distributions means the $80k avoids SE tax —
order-of-magnitude ~$12k of potential savings at 15.3%.

**The catch — "reasonable compensation":** the IRS **requires** S-corp owner-employees to pay themselves a
**reasonable salary** for the work they actually do, *before* taking distributions [irs.gov]. The IRS doesn't
publish a magic number — it weighs duties, time, training, and **comparable salaries** for similar work in
your industry and location. **Setting the salary artificially low to dodge payroll tax is a classic audit
trigger**; the IRS can reclassify distributions as wages and add back tax plus penalties.

**When the S-corp election is worth it (rules of thumb, not advice):**

- You have **consistent net profit meaningfully above a reasonable salary** for your role — commonly cited as
  roughly **$40k-80k+ of profit above your salary** before the savings outweigh the costs. Below that, the
  overhead isn't worth it.
- You can **bear the added cost and complexity**: running payroll for yourself, filing **Form 1120-S** plus
  K-1s, NC franchise tax (§5.3), payroll-tax filings, and typically a **CPA** (most S-corp owners pay for
  payroll + tax prep, easily $1-2k+/yr).
- **Deadline:** Form 2553 is generally due within **2 months and 15 days** of the start of the tax year the
  election takes effect (late-election relief exists). *Talk to a CPA before electing.*

> **Bottom line:** form the LLC now for the liability shield and simple pass-through taxes; **add the S-corp
> election later** when your profit is reliably high enough that the SE-tax savings exceed the payroll/CPA
> overhead. Don't elect on day one for a business with little or no profit.

---

## 7. Corporate Transparency Act / FinCEN BOI reporting: CHECK THE CURRENT STATUS

**This area changed repeatedly in 2024-2025 through litigation and rulemaking. Treat the status below as of
its date and VERIFY with FinCEN before relying on it.**

**Background:** The **Corporate Transparency Act (CTA)** required most LLCs and corporations ("reporting
companies") to file **Beneficial Ownership Information (BOI)** reports with the **Financial Crimes Enforcement
Network (FinCEN)** identifying their human owners. Through 2024 and early 2025 the requirement was repeatedly
enjoined, reinstated, and paused by federal courts, creating whiplash for small businesses.

**Current status (as of mid-2026 — VERIFY at fincen.gov/boi):**

- On **March 21, 2025**, FinCEN issued an **interim final rule** (published in the **Federal Register on
  March 26, 2025**) that **removes the BOI reporting requirement for U.S. companies and U.S. persons**
  [FinCEN news release, "FinCEN Removes Beneficial Ownership Reporting Requirements for U.S. Companies and
  U.S. Persons"; 90 Fed. Reg., RIN/doc. 2025-05199].
- The rule **revised the definition of "reporting company"** to mean **only entities formed under the law of
  a foreign country** that have registered to do business in a U.S. state or tribal jurisdiction. **Entities
  created in the United States — all domestic LLCs and corporations — and their U.S. beneficial owners are
  EXEMPT** from filing BOI [FinCEN; Federal Register 2025-05199].
- Foreign reporting companies still had deadlines (e.g., those registered before March 26, 2025 had until
  **April 25, 2025**), and even they were **not required to report U.S.-citizen beneficial owners** [FinCEN].
- FinCEN took public comment on the interim final rule (through **May 2025**) and signaled a possible **final
  rule** later. **A final rule, new litigation, or a change of administration policy could alter this.**

**Practical takeaway for a NC founder (as of this writing):** if you form a **domestic NC LLC or corporation**,
you are currently **NOT required to file a FinCEN BOI report** under the interim final rule. **But this is the
single most volatile compliance item in this skill — confirm the live requirement at
https://www.fincen.gov/boi (and the BOI FAQs) before assuming you have no filing.** [UNVERIFIED — LIVE STATUS;
confirm directly with FinCEN.] Do **not** pay a third-party service that tries to scare you into a paid "BOI
filing" without checking fincen.gov first; filing with FinCEN is free.

---

## 8. Federal owner taxes (self-employment tax & quarterly estimates)

Owners of pass-through businesses (sole prop, partnership, default LLC) pay their own federal taxes; there's
no employer withholding from an owner's draw. Two pieces:

### 8.1 Self-employment (SE) tax

- **Rate: 15.3%**: **12.4% Social Security** (on net earnings up to the annual SS wage base) + **2.9%
  Medicare** (no cap; an extra 0.9% Medicare surtax applies above high-income thresholds) [irs.gov].
- It's charged on your **net self-employment earnings**, computed on **92.35%** of net profit (you multiply
  net profit by 0.9235 first) [irs.gov, Form 1040-ES]. You deduct half the SE tax above the line.
- This is **on top of** regular federal income tax; the part that surprises first-time founders. (The S-corp
  election in §6 is the main lever to reduce the SE-tax base.)

### 8.2 Quarterly estimated taxes (Form 1040-ES)

Because no employer withholds for you, you generally must pay federal tax in **four quarterly estimated
payments** using **Form 1040-ES** [irs.gov]:

- **You must pay estimates if** you expect to owe **$1,000 or more** in federal tax for the year after
  withholding/credits [irs.gov].
- **Due dates:** **April 15, June 15, September 15, and January 15** (of the following year) [irs.gov].
- **Safe harbor (to avoid an underpayment penalty):** pay at least **90% of the current year's tax**, or
  **100% of last year's tax** (**110%** if your AGI exceeded **$150,000**) [irs.gov].
- Don't forget **NC estimated tax** too (**Form NC-40**) — you owe NC individual income tax (§5.4) on the same
  pass-through profit and generally must pay it quarterly as well.

Set aside roughly **25-35%+** of net profit for combined federal income tax + SE tax + NC income tax as you
earn it (the exact rate depends on your total income) so you're not caught short at the deadlines. (For the
cash-flow modeling around this, see **venture-small-business-planning**.)

---

## 9. The annual report (keep the entity alive)

NC LLCs and business corporations must file an **annual report** with the **Secretary of State** to stay in
good standing.

- **Fee:** **$200** by mail/in person, **$203** filed online (a small electronic-filing surcharge) for an LLC
  [sosnc.gov]. *Corporations file an annual report too; confirm the corp fee at sosnc.gov.*
- **Due date:** by **April 15** of each year. For a brand-new LLC, the **first** annual report is due **April
  15 of the year following** the year you formed [sosnc.gov].
- **File at** sosnc.gov (online is fastest; paper can take many weeks to process).
- **Consequence of missing it:** for an LLC, failure to file by the **60th day after** it's due is grounds for
  **administrative dissolution**; the state can shut your LLC down, and you lose the liability shield until
  you reinstate [sosnc.gov; NC GS Ch. 57D]. Set a recurring April reminder.

---

## 10. Business banking basics

Not a legal filing, but operationally essential and tied to your liability shield:

- **Open a dedicated business bank account** as soon as the entity is formed and you have an **EIN**. Banks
  typically want: the **EIN confirmation (IRS letter / CP 575)**, the **filed Articles of Organization**, and
  often the **operating agreement**.
- **Never commingle** personal and business money. Commingling is the fastest way a court pierces the LLC veil
  (§2) and reaches your personal assets. Pay yourself by a clean transfer/draw, not by swiping the business
  card for groceries.
- **Get a business debit/credit card** in the entity's name, keep clean books from day one (even a simple
  accounting tool), and reconcile monthly.
- **Sole proprietors** can use a personal account legally, but a separate account (often under the DBA) is
  still strongly advised for clean records and taxes.

---

## Quick-start checklist (solo NC LLC, the common path)

1. Choose entity (LLC for most; PLLC if licensed profession; C-corp only if raising equity)  [§1]
2. Clear the name on sosnc.gov; appoint a registered agent (NC street address)  [§2]
3. File Articles of Organization, **$125**  [§2]
4. Get a free **EIN** from irs.gov  [§2]
5. Adopt an operating agreement (keep in records; not filed)  [§2]
6. File a **DBA** at the county Register of Deeds if using a trade name  [§3]
7. Get a customized license/permit list from **BLNC (blnc.gov / 1-800-228-8443)**  [§4]
8. Register with **NCDOR** (Form NC-BR) for sales tax and/or withholding if applicable  [§5]
9. Confirm **FinCEN BOI** status at fincen.gov/boi (currently exempt for domestic entities — VERIFY)  [§7]
10. Open a **business bank account**; set aside ~25-35% of profit for taxes; plan **quarterly estimates**  [§8-9]
11. Calendar the **April 15 annual report** ($200/$203)  [§9]
12. Revisit the **S-corp election** with a CPA once profit is reliably high  [§6]

---

## Sources

Primary sources (verify the live figures — fees, rates, and especially the CTA/BOI status change):

- **NC Secretary of State (sosnc.gov)**: LLC and PLLC formation, registered agent, annual report, forms:
  - LLC fees / GS § 57D-1-22 (Articles of Organization $125): https://sosnc.gov
  - PLLC requirements & Form PLLC-02: https://sosnc.gov/manual/launching_a_business/professional_limited_liability_company_requirements
  - Annual report ($200 / $203, due April 15): https://sosnc.gov
- **NC Department of Revenue (NCDOR, ncdor.gov)**: business registration & state taxes:
  - Business registration / Form NC-BR: https://www.ncdor.gov/registration
  - Current Sales and Use Tax Rates (4.75% state + local): https://www.ncdor.gov/taxes-forms/sales-and-use-tax/sales-and-use-tax-rates/current-sales-and-use-tax-rates
  - Corporate Income & Franchise Tax Rates (income phase-out; franchise $1.50/$1,000, $200 min): https://www.ncdor.gov/taxes-forms/corporate-income-franchise-tax/corporate-income-and-franchise-tax-rates
  - Individual Income Tax Rate Schedules (flat 4.25% 2025 / 3.99% 2026): https://www.ncdor.gov/taxes-forms/individual-income-tax/tax-rate-schedules
- **IRS (irs.gov)**: federal owner taxes & S-corp:
  - S corporations / Form 2553 election: https://www.irs.gov/businesses/small-businesses-self-employed/s-corporations
  - Estimated tax / Form 1040-ES (SE tax 15.3%, 92.35%, safe harbors, due dates): https://www.irs.gov/faqs/estimated-tax and https://www.irs.gov/pub/irs-pdf/f1040es.pdf
  - EIN (free): https://www.irs.gov/businesses/small-businesses-self-employed/get-an-employer-identification-number
- **FinCEN (fincen.gov)**: Corporate Transparency Act / BOI (LIVE STATUS — verify):
  - BOI landing page: https://www.fincen.gov/boi  ·  BOI FAQs: https://www.fincen.gov/boi-faqs
  - News release: "FinCEN Removes Beneficial Ownership Reporting Requirements for U.S. Companies and U.S. Persons": https://www.fincen.gov/news/news-releases/fincen-removes-beneficial-ownership-reporting-requirements-us-companies-and-us
  - Federal Register, Interim Final Rule (doc. 2025-05199, published 2025-03-26): https://www.federalregister.gov/documents/2025/03/26/2025-05199/beneficial-ownership-information-reporting-requirement-revision-and-deadline-extension
- **EDPNC / Business Link NC (BLNC)**: licenses & permits, free counseling: https://edpnc.com/get-business-support/start-a-business/ and https://blnc.gov (1-800-228-8443)
- **NC General Assembly (ncleg.gov)**: statutes: Ch. 57D (LLC Act), Ch. 55B (Professional Corporations), Ch. 66 Art. 14A (assumed business names), Ch. 105 (taxes; privilege-tax repeals via SL 2015-141 and SL 2023-134), assumed-name modernization SL 2017-23 (HB 228): https://www.ncleg.gov

Secondary/explanatory sources consulted (law firms, Tax Foundation, accounting guides) corroborated the
above and are not a substitute for the primary sources.

---

## Disclaimer

This skill is **general information, not legal or tax advice**, and does not create an attorney-client or
accountant-client relationship. Fees, tax rates, due dates, statutes, and **especially the Corporate
Transparency Act / FinCEN BOI reporting status** change — sometimes abruptly through litigation or new
rulemaking — and several figures here carry an as-of date or an [UNVERIFIED] flag. Before you file, register,
elect, or rely on anything in this skill, **verify the current requirements directly with the NC Secretary of
State (sosnc.gov), the NC Department of Revenue (ncdor.gov), the IRS (irs.gov), and FinCEN (fincen.gov/boi)**,
and **consult a licensed North Carolina business attorney and a CPA** for advice on your specific situation —
particularly for entity choice, the S-corp election and reasonable compensation, multi-state nexus, and any
BOI obligation. County DBA fees and procedures vary by Register of Deeds. Professional-licensing rules (PLLC
ownership ratios, board approvals) are set by each licensing board — confirm with yours.
