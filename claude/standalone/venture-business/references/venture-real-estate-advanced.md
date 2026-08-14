<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-real-estate-advanced` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-real-estate-advanced
description: >-
  Advanced NC real-estate operations beyond intro investing and licensing law: NC property
  management law, short-term/vacation rentals, and creative/seller financing. Educational only,
  not legal advice.
  TRIGGER: managing others' NC rentals for pay and whether a broker license is needed; NC
  security deposits / caps / trust account / accounting; NC eviction / summary ejectment /
  self-help; broker trust accounting; short-term / vacation-rental underwriting (ADR / RevPAR /
  occupancy / seasonality); NC Vacation Rental Act; STR zoning / occupancy tax; subject-to /
  wraparound / lease-option / owner financing; due-on-sale; Dodd-Frank seller-financing limits.
  SKIP: NC licensing / agency / brokerage law -> venture-nc-real-estate-law; intro RE marketing
  + investing metrics (cap rate, BRRRR, flip, 1031) -> venture-real-estate-marketing-investing;
  marketing compliance -> venture-marketing-strategy-local-seo; brokerage entity / tax ->
  venture-nc-business-formation-tax.
category: personal-venture
tags: [venture, real-estate, property-management, short-term-rental, creative-financing, north-carolina]
whenToUse:
  - You want to manage rental property for other owners in NC and need to know whether an NC broker license is required and what the exemptions are.
  - You are handling residential security deposits in NC and need the deposit caps, trust-account/bond rule, and the 30/60-day accounting deadlines.
  - You need the basic NC eviction (summary ejectment) sequence and want to avoid an illegal self-help eviction.
  - You are underwriting a short-term / vacation rental and need occupancy, ADR, RevPAR, and seasonality math, and why long-term cap-rate math misleads here.
  - You are renting property for under 90 days in NC and need the Vacation Rental Act (Ch. 42A) rules plus local zoning/permit/occupancy-tax realities.
  - You are evaluating subject-to, wraparound, lease-option, or seller/owner financing and need an honest read on due-on-sale and Dodd-Frank/SAFE Act risk.
  - You manage trust money as a broker and need the 21 NCAC 58A .0116 handling rules.
triggers:
  - "Do I need a license to manage someone else's rental in NC?"
  - "What's the max security deposit I can charge in North Carolina?"
  - "How long do I have to return a tenant's deposit in NC?"
  - "How do I evict a tenant in North Carolina / what is summary ejectment?"
  - "How do I underwrite an Airbnb, what's a good occupancy and RevPAR?"
  - "What is the NC Vacation Rental Act and does it apply to my Airbnb?"
  - "Is subject-to / a wraparound mortgage legal, what about the due-on-sale clause?"
  - "Can I do seller financing on a house without becoming a mortgage loan originator?"
version: 1.1.0
updated: 2026-06-16
metadata:
  changelog:
    - "2026-06-16 sko v1.0.0->v1.1.0: fixed 1 High (Pass J body ~10.2k->~8.6k tok, extracted worked example + cheat sheet + sources to references/) + 2 Medium (Pass M description 2238->967 chars; Pass K em-dash density 1.44->0.79); 0 banned terms; blind re-audit gate cleared after 1 extra iteration (en-dash/em-dash count reconciled)"
---

# Advanced Real Estate: Property Management, Short-Term Rentals & Creative Financing (NC-anchored)

> **Educational only — not legal, tax, or financial advice.** Property-management licensing,
> the Tenant Security Deposit Act, the Vacation Rental Act, and creative/seller financing are
> all legally sensitive and fast-moving. Statutes and local ordinances change. Creative
> financing in particular can cross into territory that triggers federal mortgage-licensing
> rules or lets a lender call a loan. **Verify with a licensed NC real-estate attorney before
> acting**, and route licensing/agency questions to the sibling skill.

## What this skill is (and what it isn't)

This is the **advanced operations** layer of the real-estate venture stack. It assumes you
already know the intro material and the licensing baseline, and goes deep on three things a
NC operator actually has to get right:

1. **Running rentals for other people** (property management), and the licensing + trust +
   deposit + eviction rules that govern it in NC.
2. **Short-term / vacation rentals**, which underwrite and regulate differently from
   long-term holds.
3. **Creative & seller financing**: powerful, and the most legally dangerous thing in this
   skill.

**Cross-references (do not duplicate these; go there):**

- **`venture-nc-real-estate-law`** — the NCREC, GS Ch. 93A, getting/keeping a license,
  provisional broker, **broker-in-charge (BIC)**, firm license, agency law, RPOADS, the WWREA
  disclosure, fair housing, advertising rules, discipline. **Any "who is legally allowed to
  do X / how do I get licensed / what disclosure is required" question goes there.**
- **`venture-real-estate-marketing-investing`**: listing marketing, lead gen, the NAR
  settlement, and **intro investing metrics** (cap rate, cash-on-cash, NOI, 1% rule, GRM,
  DSCR, BRRRR, flip/70% rule, wholesaling, REITs, financing types, 1031/depreciation basics).
  **Cap-rate fundamentals and standard buy-and-hold analysis live there**; this skill builds
  on them for STR.
- **`venture-nc-business-formation-tax`**: forming the LLC/PLLC/entity that owns the
  business and its tax filings.

---

## PILLAR 1 — North Carolina Property Management

### 1.1 Does managing other people's property require an NC broker license?

**Short answer: generally yes.** North Carolina treats most property-management activity
performed *for another* *for compensation* as **real estate brokerage**, which requires a
license issued by the **North Carolina Real Estate Commission (NCREC)** under **GS Chapter
93A**. The classic licensed activities include leasing/renting others' property, negotiating
leases, showing units, collecting rent, and handling tenant security deposits on an owner's
behalf. (NCREC, "Can I Get Paid?", bulletins.ncrec.gov, as of 2026.)

The test NCREC applies is roughly: are you (a) performing an act of brokerage, (b) for
another person, (c) for compensation or the expectation of it? If yes, you need a license.
(NCREC bulletins, as of 2026 — confirm the current list of brokerage acts on ncrec.gov.)

**The narrow exemption that matters for a small operator:** a **W-2 salaried employee** who
manages property for a **single employer/owner** generally does not need a license to do so.
The exemption is narrow: it breaks if you manage for **multiple owners**, or if you are paid
on a **commission/transaction basis** rather than salary, or if you are an independent
contractor rather than an employee. (NCREC, as of 2026; **[UNVERIFIED]** as to every edge of
the salaried-employee exemption — verify the exact wording of GS 93A-2 exclusions and current
NCREC guidance before relying on it.)

**You can always manage your own property.** Owning and renting out property you hold title to
is not brokerage. The license question only arises when you act **for someone else for pay**.

**If you do form a property-management company:** if it is an entity (LLC/corp), the entity
itself needs a **firm license** and must designate a **broker-in-charge (BIC)** to supervise
brokerage activity and trust money. (NCREC firm-license guidance, as of 2026.) **Firm
license, BIC qualification, and supervision rules belong to `venture-nc-real-estate-law` —
go there for the how.**

> **Routing rule:** "Do I need a license / how do I get the firm license / who can be BIC?"
> → these are *licensing-law* questions. Read this section for the property-management
> *context*, then go to **`venture-nc-real-estate-law`** for the licensing mechanics.

### 1.2 The NC Tenant Security Deposit Act (TSDA) — G.S. 42-50 et seq.

This is the single most-litigated thing a NC landlord/PM touches. The TSDA is **Article 6 of
GS Chapter 42** (§§ 42-50 through 42-56). It applies to **residential** rentals.

**Deposit caps (by lease term).** A residential security deposit may not exceed
(G.S. 42-50/42-51, as of 2026):

| Tenancy type | Max security deposit |
|---|---|
| Week-to-week | **2 weeks' rent** |
| Month-to-month | **1.5 months' rent** |
| Term greater than month-to-month (e.g., a 1-year lease) | **2 months' rent** |

(Statutory text: "two weeks' rent if a tenancy is week to week, one and one-half months' rent
if a tenancy is month to month, and two months' rent for terms greater than month to month."
The cap language appears in the deposit statute; some secondary sources cite the cap at
G.S. 42-51 — **read the current statute**.) A separate **reasonable, nonrefundable pet fee**
is generally allowed on top of the deposit and is not counted against the cap. **[UNVERIFIED]**
as to the pet-fee treatment — confirm against current law.

**Trust account OR bond, mandatory.** Residential security deposits must either be:
- **Deposited in a trust account** with a **licensed and federally insured depository
  institution** (or a trust institution authorized to do business in NC); **or**
- Secured by a **bond** from an insurance company licensed in NC (the landlord's option).
- Deposits may be held in a trust account **outside NC only if** the landlord provides the
  tenant an adequate bond in the amount of the deposits.

(G.S. 42-50, as of 2026.) **Note:** the bond *option* exists for ordinary long-term
residential deposits — but **NOT** for vacation rentals under Ch. 42A (see Pillar 2; Ch. 42A
forces a trust account).

**30-day location notice.** The landlord/agent must **notify the tenant, within 30 days after
the lease term begins,** of the **name and address of the bank/institution** holding the
deposit (or the insurance company if a bond is used). (G.S. 42-50, as of 2026.)

**Permitted deductions.** A landlord may apply the deposit to (G.S. 42-51, as of 2026):
unpaid rent; damage **beyond normal wear and tear**; unpaid bills (e.g., utilities/water/sewer)
that become a lien; costs of **re-renting** after the tenant's breach; **removal and storage**
of the tenant's property after a lawful eviction; court costs; and cleaning costs needed to
restore the unit. Deductions for **normal wear and tear are not allowed.**

**Accounting deadlines, 30 / 60 days.** After the tenancy ends and the tenant surrenders
possession (G.S. 42-52, as of 2026):
- The landlord must **itemize and refund** the balance within **30 days**.
- If the amount of the claim **cannot be determined within 30 days**, the landlord must send
  an **interim accounting within 30 days** and a **final accounting within 60 days** of
  termination + delivery of possession.
- The itemized statement is sent by **first-class mail to the tenant's last known address**
  (a forwarding address the tenant provides).

**Penalty for getting it wrong.** A landlord who **fails to provide the required itemized
accounting** within the deadlines **forfeits the right to retain any portion of the deposit**
and may be liable for the deposit plus **damages and attorney's fees**. (G.S. 42-55/42-56, as
of 2026.) This is why deposit accounting discipline matters: a sloppy or late accounting can
turn a legitimate deduction into a total loss plus the tenant's legal fees.

### 1.3 Trust accounting for property managers (broker rules — 21 NCAC 58A .0116)

If you are a **licensed broker** handling other people's money (rent, deposits, association
funds), NCREC trust-money rules apply on top of the TSDA. Key points from **21 NCAC 58A .0116
("Handling of Trust Money")** (as of 2026):

- **Deposit deadline:** trust money must be placed in a **trust/escrow account within three
  banking days** of receipt. (Earnest money / tenant security deposits paid by non-currency
  means must be deposited **no later than three days following acceptance** of the offer to
  purchase/lease.)
- **No conversion / no commingling:** a broker may not convert others' money to the broker's
  own use or apply it to a purpose other than intended; trust money is kept separate from the
  broker's operating funds.
- **Tenant security deposits** held by a broker are disposed of **per G.S. 42-50 through 42-56
  and G.S. 42A-18** (i.e., the TSDA and Vacation Rental Act control the substance).
- **Property owners' association (POA) funds** controlled by a broker are **trust money** and
  must sit in an account **dedicated exclusively to a single POA**, never commingled with
  other associations' or other parties' funds.
- **Records:** brokers must keep records tracing ownership of trust money from deposit through
  final disbursement, sufficient to verify accuracy and proper use.

(Account-labeling specifics — "Trust Account"/"Escrow Account" designation, withdrawal-on-
demand, reconciliation cadence — are detailed in the companion rules **.0117/.0118**; confirm
the current text. **[UNVERIFIED]** as to the exact current labeling/reconciliation wording.)

> **Routing rule:** the *mechanics* of broker trust accounts and reconciliation are NCREC
> compliance — for the licensing-side detail also see **`venture-nc-real-estate-law`** (trust/
> escrow account section). This skill covers how it intersects with deposits and PM operations.

### 1.4 Eviction in NC — summary ejectment basics

NC eviction is called **summary ejectment** (GS Ch. 42, Article 3, §§ 42-26 et seq.). The
operator-critical points (as of 2026):

1. **Grounds (G.S. 42-26):** the main grounds are (a) **nonpayment of rent**, (b) **breach of
   a lease condition** for which the lease provides forfeiture/re-entry, and (c) **holding
   over** after the lease term ends.
2. **Notice/demand first.** For nonpayment, NC requires a **10-day demand for rent** before
   filing (G.S. 42-3, the statutory grace/demand for residential leases, subject to lease
   terms). For holdover/other breaches, give the notice the lease and statute require.
   **[UNVERIFIED]** as to the exact notice period for each ground in your situation — confirm
   the current statute and your lease.
3. **File a Complaint in Summary Ejectment** with the **magistrate (small claims) court** —
   the clerk of superior court issues a summons. **No self-help.**
4. **Fast hearing.** The summons sets a hearing **not more than 7 days** out (excluding
   weekends/holidays); service must be achieved at least **2 days before** the hearing.
   A magistrate may continue the case only on good cause, and not more than **5 days** unless
   the parties consent.
5. **Judgment + appeal + writ.** If the landlord wins, the tenant has an **appeal window**
   (to district court), and only after that does the landlord obtain a **writ of possession**
   executed by the **sheriff** — the sheriff, not the landlord, removes the tenant.

**Self-help eviction is illegal** for residential property. Changing the locks, shutting off
utilities, or removing a tenant's belongings to force them out exposes you to liability (and
potentially criminal exposure). **Always go through the court + sheriff.** (G.S. 42-25.6 et
seq.; UNC School of Government summary-ejectment materials, as of 2026.)

---

## PILLAR 2 — Short-Term / Vacation Rentals (Airbnb / VRBO)

### 2.1 STR underwriting — why it is NOT long-term cap-rate math

Long-term rentals (LTR) underwrite on **stable monthly rent → NOI → cap rate / cash-on-cash**
(that's the sibling skill, `venture-real-estate-marketing-investing`). STRs are a small
**hospitality business**, and they underwrite on hotel-style metrics:

- **ADR (Average Daily Rate)** = total room/night revenue ÷ nights booked. The price per
  booked night. A high ADR alone can *hide* weak demand if you're booking few nights.
- **Occupancy** = nights booked ÷ nights available. (US STR occupancy has hovered around the
  mid-50s% (about **54.9%** stabilizing in 2025 per industry data), but this is **wildly
  market- and seasonality-dependent**; a beach town can be 85% in summer and 20% in winter.)
- **RevPAR (Revenue Per Available Rental/Room)** = **ADR × Occupancy**. This is the single
  cleanest read on whether pricing *and* demand are working together. Use RevPAR, not ADR, to
  compare properties.
- **RevPAN** (per available night) is the same idea for whole-unit rentals.

**Why you can't just slap an LTR cap rate on an STR pro forma:**
- STRs carry **higher operating costs**: management (often **18–25%**, up to **30–40%** in
  luxury/resort markets), **cleaning/turnover** (US average cleaning fee ~$160/turn; many
  turns/month), supplies, utilities (you pay them), higher insurance, furnishing/replacement,
  platform fees, and **occupancy/sales taxes**.
- STRs carry **revenue volatility and management dependency**, i.e., **operational risk** an
  LTR doesn't have. Hospitality assets trade at **cap rates ~100–300 bps wider** than
  comparable multifamily for exactly this reason; using a multifamily cap rate on a hotel-like
  pro forma can **overstate value by ~15–30%**. (TILT Analytics hotel-underwriting guidance,
  as of 2026.) Treat STR the same way: discount for risk, don't capitalize peak season as if
  it were year-round.
- **Seasonality** must be modeled month-by-month. Underwriting off a single peak month or a
  blended annual average that ignores the trough will badly mislead you.

**Tooling:** AirDNA (Rentalizer/MarketMinder), Mashvisor, and similar pull historical ADR/
occupancy by market and submarket. Treat their revenue estimates as a **starting point**;
they often include cleaning fees in "revenue" and assume professional management/pricing.
Sanity-check against actual comps and your own expense stack. (AirDNA, Mashvisor, as of 2026.)

### 2.2 Worked example — STR vs LTR on the same NC house

Same 3BR house (about $320,000 all-in). As a **long-term rental** at $2,200/mo with about 40%
expenses, NOI lands near $15,840 (cap rate about 4.95%), stable and low-touch. As a
**short-term rental** at ADR $240 and 55% occupancy, gross runs about $48,240, but STR opex is
far heavier (management about 20%, cleaning, utilities you pay, STR insurance, furnishings),
roughly 64% of gross, so NOI lands near $17,290 (cap rate about 5.4%).

The lesson: the STR grosses nearly double, yet after the heavier expense load the *NOI* gap is
small (about $17.3k vs $15.8k) while the STR carries far more operational risk, volatility,
regulatory exposure, and your time. Underwrite the STR at a **wider required cap rate / lower
price**, **stress-test occupancy and ADR downward**, and model the **low season** explicitly.
An STR that only pencils at peak-season numbers is a bad deal. **Full line-by-line math:
`references/worked-examples-and-tables.md`.**

### 2.3 The NC Vacation Rental Act — G.S. Chapter 42A (rentals < 90 days)

NC carves vacation rentals out of the ordinary landlord-tenant chapter into **GS Chapter 42A,
the Vacation Rental Act (VRA)**. Coverage and key rules (as of 2026):

- **What it covers:** the rental of residential property for **vacation, leisure, or
  recreation** for **fewer than 90 days** by a person who **has a permanent residence to
  return to**. (G.S. 42A-4.) This is the legal definition that separates an STR from a normal
  lease.
- **Written agreement required.** A vacation rental must be under a **written vacation rental
  agreement** with statutorily required terms. (G.S. 42A-10/42A-11.)
- **Security deposits: trust account, NO bond option.** Vacation-rental security deposits
  must go into a **trust account per G.S. 42-50**, and unlike ordinary residential deposits,
  the landlord/broker **does NOT get the bond-in-lieu option**. (G.S. 42A-18.)
- **Faster deposit accounting, 45 days.** The landlord/broker must **apply, account for, or
  refund** the deposit within **45 days** after the tenancy ends (vs the 30/60-day TSDA
  timeline for long-term). (G.S. 42A-18.)
- **Advance payments → trust account within 3 banking days.** Any **advance payment** the
  tenant makes (rent or otherwise) must be deposited into a trust account at a federally
  insured depository/trust institution **no later than three banking days** after receipt.
  Those funds **don't earn interest** unless the agreement says so and states who gets it.
  (G.S. 42A-16.)
- **Tenant eviction under 42A.** Ch. 42A has its own **expedited eviction** provisions for
  vacation tenants (G.S. 42A-24 et seq.) distinct from summary ejectment. **[UNVERIFIED]** as
  to the exact 42A expedited-eviction procedure — confirm before relying.

> **Key distinction for an operator:** if your Airbnb guest is a vacationer staying < 90 days,
> **Ch. 42A — not the ordinary TSDA/Article-6 rules — governs the agreement and deposits.** A
> longer-term furnished rental (>90 days, or to someone with no other permanent home) may fall
> back under ordinary landlord-tenant law. Get the characterization right.

### 2.4 Local zoning, permits & occupancy tax — the part that actually kills deals

The VRA is *state* contract law. **Whether you may operate an STR at all is decided locally**,
and NC municipalities vary enormously. **Check the specific city/county before you buy.**

- **Zoning & permits vary widely.** Examples (verify current ordinances — these move):
  - **Asheville** is notably **restrictive**: whole-home STRs are largely **prohibited in
    residential districts**; owner-present **"homestays"** are allowed under conditions, and
    violations carry steep penalties. Some uses need **conditional zoning** approval.
  - **Charlotte** regulates via **zoning + a noise ordinance** rather than a standalone STR
    registry (per its UDO); operators still must comply with zoning.
  - **Wilmington and the beach towns** (Wrightsville Beach, Carolina Beach, OBX communities)
    generally **permit STRs with registration**, and some require **commercial general
    liability insurance** (e.g., **$500,000+ per occurrence** for whole-house lodging in some
    jurisdictions).
  - (Awning / UNC School of Government STR land-use materials, as of 2026 — **[UNVERIFIED]**
    for any specific current threshold; ordinances change frequently and are litigated.)
- **NC limits how far cities can go.** NC has a state law constraining local STR **registration/
  permitting** in some respects (the legislature has restricted certain local STR registries),
  so local power is real but not unlimited. **[UNVERIFIED]** — confirm the current state-vs-
  local balance with NC counsel, as this area has active litigation and legislation.
- **Occupancy + sales tax.** STRs (< 90 continuous days) are taxable accommodations:
  - **State sales tax on accommodations** applies (NC general state rate **4.75%** + local
    sales tax, commonly **~6.75–7.5%** combined; G.S. 105-164.4(a)(3)). **[UNVERIFIED]** exact
    combined rate by county — confirm.
  - **Local room/occupancy tax** is layered on top, set by county/municipality, **commonly
    3%–8%** (e.g., many counties at **6%**; some add tourism/convention-center fees).
  - **Marketplace facilitators (Airbnb/VRBO) collect and remit** the occupancy/sales tax for
    bookings made and paid **through the platform** in NC. If **all** your bookings run through
    such a facilitator, you generally **do not need to register or file** lodging-tax returns
    yourself. But if guests **book/pay you directly** (off-platform), **you** are the retailer
    responsible for collecting and remitting. (Avalara MyLodgeTax NC guide; NC county
    occupancy-tax FAQs, as of 2026.)

> **Practical workflow:** before buying for STR — (1) confirm the **specific local ordinance**
> allows your intended use (whole-home vs owner-occupied homestay) and get the **permit/
> registration** requirements; (2) confirm **HOA/restrictive covenants** don't ban STRs
> (covenants can be stricter than zoning); (3) line up **STR-appropriate insurance**; (4) set
> up **occupancy/sales-tax** handling (rely on the platform only if 100% of bookings are
> on-platform); (5) underwrite per 2.1–2.2 with seasonality and the heavier expense load.

---

## PILLAR 3 — Creative & Seller Financing (the legally dangerous pillar)

> **Read this first.** Creative financing can be legitimate and powerful, but several of these
> structures sit on top of legal trip-wires: a lender can **call the loan** (due-on-sale), or
> the deal can trigger **federal mortgage-licensing/consumer-protection rules** (Dodd-Frank /
> SAFE Act) when the buyer is an **owner-occupant**. **Treat every owner-occupied
> seller-financed deal as presumptively regulated until an attorney confirms otherwise.** The
> intro skill (`venture-real-estate-marketing-investing`) lists financing *types*; this section
> is about the **legal sensitivity** of the creative ones.

### 3.1 The structures, in plain terms

- **Seller carryback / owner financing.** The seller acts as the bank: buyer signs a
  **promissory note** secured by a **deed of trust**, pays the seller over time. Used when a
  buyer can't get (or doesn't want) bank financing, or to spread the seller's gain.
- **Subject-to ("sub-to").** Buyer takes **title** but the **seller's existing mortgage stays
  in place** in the seller's name; buyer just makes the payments. No new loan. **Triggers the
  due-on-sale clause** (see 3.2). The seller stays legally on the hook to the lender.
- **Wraparound ("wrap") mortgage.** Buyer signs a **new note to the seller** that **wraps
  around** the seller's existing loan; the seller keeps paying the original lender and pockets
  the spread. Also a transfer that **triggers due-on-sale**, with added risk that the seller
  must keep paying the underlying loan from the buyer's payments.
- **Lease-option / lease-purchase.** Tenant **leases** with an **option (or obligation) to
  buy** later at a set price; sometimes a portion of rent is **credited** to the purchase.
  In NC, residential lease-with-option contracts are governed by **GS Chapter 47G** (see 3.4)
  — and if rent credits build equity, the deal can be treated as **disguised financing** and
  drop into Dodd-Frank's orbit (see 3.3).

### 3.2 Due-on-sale risk & the Garn-St Germain Act

Almost every conventional mortgage from the last ~40 years has a **due-on-sale (acceleration)
clause**: if the borrower **transfers title**, the lender may **call the entire balance due**.
Federal law, the **Garn-St Germain Depository Institutions Act, 12 U.S.C. § 1701j-3**,
**preempts state law and lets lenders enforce** those clauses, with a list of **exceptions**
where the lender **may not** call the loan.

**The exceptions are about death, divorce, and family estate-planning transfers, not
investor deals.** They include things like transfer on the borrower's death to a relative,
transfer to a spouse/child, transfer incident to divorce, and a **transfer into an inter
vivos trust in which the borrower is and remains a beneficiary and which does not relate to a
transfer of rights of occupancy** (§ 1701j-3(d)(8); 12 C.F.R. Part 191). (Adam Leitman
Bailey / Miller Miller & Canby summaries, as of 2026.)

**Subject-to and wraparound transfers are NOT on the exception list.** That means:
- The lender **legally may accelerate** the loan when it learns of the transfer.
- "Most lenders don't call the loan as long as payments stay current" is a **business
  observation, not a legal defense.** In a rising-rate environment a lender has *more* reason
  to call a cheap old loan. Plan for the risk; don't assume it away.
- The seller in a sub-to/wrap **remains personally liable** on the original note and the loan
  stays on their credit — a real harm if the buyer stops paying.

**Do not** try to dodge due-on-sale by mischaracterizing an investor purchase as an
"inter vivos trust transfer." Lenders and courts see through land-trust gimmicks, and the
trust exception requires the **borrower** to remain beneficiary and **occupancy not to
transfer** — which an investor flip doesn't satisfy. **[UNVERIFIED]** as to any specific
trust structure's effectiveness — this is squarely an attorney question.

### 3.3 Dodd-Frank & the SAFE Act — financing **owner-occupied** homes

This is the federal layer that catches a lot of would-be seller-financers. The **SAFE Act
(2008)** requires people **in the business of originating residential mortgage loans** to be
**licensed Mortgage Loan Originators (MLOs)**. **Dodd-Frank's Loan Originator Rule** (CFPB,
Regulation Z) extended consumer-protection requirements (incl. **ability-to-repay**) to
residential-mortgage origination, and **seller financing of a dwelling can count** unless an
**exclusion** applies.

**Critical scoping:** these rules bite when the financed property is a **dwelling** and the
buyer is a **consumer / owner-occupant** in a **1–4 unit residential** property. **Seller
financing to a business entity, or for purely investment/commercial property, generally falls
outside** the consumer-mortgage rules. (Frascona; Barnes Walker; NAR "The SAFE Act: Seller
Financing"; as of 2026 — confirm current CFPB rule text.)

Two CFPB **exclusions** let a non-MLO seller carry a note on an owner-occupied home **without
becoming a licensed MLO** (as of 2026 — verify the current Regulation Z text and any state
overlay):

| | **One-property exclusion** | **Three-property exclusion** |
|---|---|---|
| Who can be the seller | **Natural person, estate, or trust** (NOT an LLC/corp/partnership) | Seller may be a **person OR an entity** |
| Max financed sales / 12 months | **1** | **3** |
| Balloon payment | **Allowed** (no negative amortization) | **Not allowed — must be fully amortizing** |
| Ability-to-repay | You **do not have to document** it (good faith still expected) | Must **determine in good faith and document** reasonable ability to repay |
| Rate | (no specific structural rate test beyond no-neg-am) | **Fixed**, or **adjustable resetting after ≥5 years** with reasonable caps |
| Did not build/construct the home | required | required |

(Sources: ORT/First American "CFPB and Seller Financing"; NAR; Barnes Walker; Frascona —
as of 2026.) **If you exceed these limits, or the seller is an entity wanting balloon terms,
you likely need a licensed MLO to originate the loan (or use a Residential Mortgage Loan
Origination ("RMLO") service to underwrite the ability-to-repay).** **[UNVERIFIED]** as to the
exact current thresholds and whether NC imposes any additional state MLO requirement — confirm
with counsel and the NC Commissioner of Banks.

**Bottom line:** seller-financing **investment property to an entity** is the low-risk lane;
seller-financing an **owner-occupied home to a consumer** is the regulated lane: keep it to
**1 (or ≤3) deals/year within an exclusion**, document ability-to-repay, **avoid balloons**
on owner-occupied consumer deals, and get the note/deed-of-trust and disclosures drafted by an
attorney/RMLO.

### 3.4 NC lease-options — GS Chapter 47G

NC specifically regulates **residential lease-with-option-to-purchase** deals under **GS
Chapter 47G** ("Option to Purchase Contracts Executed with Lease Agreements for Residential
Property"). Key requirements (as of 2026):

- The **option contract must be in writing** and contain minimum terms: **full names/addresses
  of the parties, legal description of the property, the sales price, and the option fee** (and
  any other fees/payments between the parties). (G.S. 47G-2.)
- **Recordation and consumer-protection features** apply (the chapter exists precisely because
  these deals were abused). NCREC has **disciplined firms** for offering **illegal lease-option/
  "rent-to-own" + credit-repair** schemes — so this is an enforcement-active area. (NCREC
  bulletins, as of 2026.)
- **Federal overlay:** if rent payments **credit toward the price or build equity**, the deal
  can be recharacterized as a **disguised installment sale / seller financing** and fall under
  **Dodd-Frank** (3.3). Structure (and label) carefully, with counsel. **[UNVERIFIED]** as to
  exactly when a given NC lease-option tips into a regulated financing — attorney question.

> **Routing rule:** doing a lease-option, sub-to, or wrap usually means you are **performing
> brokerage and/or originating financing** — both of which raise **licensing** questions that
> belong to **`venture-nc-real-estate-law`** (real-estate license) and the **NC Commissioner
> of Banks** (MLO license). Use this section for the deal mechanics + risk; route the
> "do I need a license to do this" question accordingly.

---

## Quick-reference cheat sheet, sources & child concepts

The **one-line-per-rule cheat sheet** (deposit caps, deadlines, trust-money rule, eviction,
VRA, STR taxes, due-on-sale, seller-financing exclusions, lease-options), the **full
primary/secondary source list** (with as-of dates), and the deeper **child concepts** for
further research all live in **`references/worked-examples-and-tables.md`**. Read that file when
you need the at-a-glance table, a citation, or a pointer to an adjacent topic to research next.

Key numbers to remember without opening it: deposit caps **2 weeks / 1.5 months / 2 months**;
deposit accounting **30 days** (interim 30 / final 60); broker trust money in **3 banking
days**; vacation-rental deposit **45 days, no bond option**; STR local occupancy tax **about 3
to 8%**; **RevPAR = ADR x Occupancy**; seller financing an owner-occupied home stays within the
**1- or 3-property CFPB exclusion** or needs an MLO.

---

## Disclaimer

This skill is **educational information, not legal, tax, or financial advice**, and does not
create an attorney-client or advisory relationship. North Carolina property-management
licensing, the Tenant Security Deposit Act, the Vacation Rental Act, broker trust-accounting
rules, local STR ordinances/occupancy taxes, and federal due-on-sale and Dodd-Frank/SAFE Act
rules **change over time and are fact-specific**. **Creative and seller financing are
especially legally sensitive** — subject-to and wraparound deals can let a lender accelerate
the loan, and financing an owner-occupied home can trigger federal mortgage-licensing and
ability-to-repay requirements. Items flagged **[UNVERIFIED]** were not confirmed against
primary text during authoring and must be checked. **Before acting on anything here, verify
the current statute/rule and consult a licensed North Carolina real-estate attorney** (and,
for financing, a mortgage/lending attorney or licensed MLO, and a CPA for tax). Route
licensing/agency/brokerage-law questions to `venture-nc-real-estate-law`.
