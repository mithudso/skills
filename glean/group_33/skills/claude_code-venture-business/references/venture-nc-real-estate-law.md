<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-nc-real-estate-law` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-nc-real-estate-law
description: >-
  North Carolina real-estate LICENSING and brokerage-law reference for a NC founder, grounded in
  ncrec.gov and GS Chapter 93A: the NCREC, the single-license/provisional-broker system, getting
  licensed (prelicensing, PSI exam, postlicensing, CE), firm license and the BIC role,
  trust/escrow accounts, agency law and disclosures, the RPOADS and material-fact duty, fair
  housing, advertising, and discipline.
  TRIGGER: NC real-estate licensing, NCREC, provisional broker, pre/postlicensing/CE hours, firm
  license, broker-in-charge (BIC), trust/escrow accounts, the WWREA disclosure, dual/designated
  dual agency, buyer-agency agreements, RPOADS or material-fact disclosure, fair housing, NCREC
  advertising, or license discipline (GS 93A).
  SKIP: real-estate MARKETING/lead-gen or INVESTING -> venture-real-estate-marketing-investing;
  the brokerage ENTITY or TAX -> venture-nc-business-formation-tax; general marketing compliance
  (CAN-SPAM, TCPA, GBP/local SEO) -> venture-marketing-strategy-local-seo.
category: personal-venture
tags: [venture, north-carolina, real-estate, licensing, ncrec, agency-law]
version: 1.0.1
updated: 2026-06-16
whenToUse:
  - "How do I get a real estate license in North Carolina"
  - "What is a provisional broker vs a broker vs a broker-in-charge in NC"
  - "How many prelicensing, postlicensing, or continuing-education hours does NC require"
  - "Do I need a firm license and who can be the broker-in-charge"
  - "What are NC trust/escrow account and recordkeeping rules"
  - "When must I give the Working With Real Estate Agents disclosure; how does dual agency work"
  - "When must buyer-agency agreements be in writing after the NAR settlement"
  - "What must a broker disclose on the RPOADS or as a material fact, even in an as-is sale"
  - "What are the NC fair-housing and NCREC advertising (blind-ad, team-name) rules"
triggers:
  - "NC real estate license"
  - "NCREC"
  - "provisional broker"
  - "broker-in-charge BIC"
  - "real estate trust account NC"
  - "Working With Real Estate Agents disclosure"
  - "RPOADS material fact"
  - "NC fair housing real estate"
---

# North Carolina Real Estate Licensing & Brokerage Law

Reference for the licensing path and brokerage-law duties of a North Carolina real-estate
broker. The single authoritative regulator is the **North Carolina Real Estate Commission
(NCREC)**, and the governing law is **General Statutes (GS) Chapter 93A** plus the Commission's
administrative rules in **Title 21, Chapter 58 of the NC Administrative Code (21 NCAC 58)**.
Cite both inline. Cross-references: real-estate **marketing/investing** ->
`venture-real-estate-marketing-investing`; forming the brokerage **entity/tax** ->
`venture-nc-business-formation-tax`; general **marketing compliance** ->
`venture-marketing-strategy-local-seo`.

> **This is general information, NOT legal advice.** Hours, fees, rules, and forms change.
> Verify everything against ncrec.gov and current GS Chapter 93A before relying on it, and
> consult a NC attorney for legal questions. See the Disclaimer at the end.

---

## 1. The NCREC and what it regulates

The NCREC is an independent state regulatory agency that licenses and regulates real-estate
brokers and brokerage firms in NC. It is created by **GS 93A-3**, which sets a **nine-member
Commission** (seven appointed by the Governor and two by the General Assembly, one each on
Senate and House recommendation); at least three members must be licensed brokers and at least
two must have no real-estate/appraisal business involvement (GS 93A-3, as of ncleg.gov 2026-06).

What NCREC does (GS Chapter 93A; 21 NCAC 58):

- **Issues and renews broker and firm licenses** and sets education requirements.
- **Approves prelicensing/postlicensing/CE education** and the schools/providers that deliver it
  (GS 93A Article 3).
- **Adopts the rules** brokers must follow (agency, trust accounts, advertising, disclosure).
- **Investigates complaints and disciplines licensees** (GS 93A-6).
- **Administers the Real Estate Education and Recovery Fund** (GS 93A Article 2), which
  reimburses consumers for certain losses caused by licensee misconduct.

**Who must be licensed (GS 93A-1, GS 93A-2).** It is unlawful to act as a real-estate broker
(listing, selling, buying, leasing, renting, auctioning, or negotiating those for *others* for
compensation) without a NCREC license (GS 93A-1, as of ncleg.gov 2026-06). **GS 93A-2** lists
exemptions, including owners dealing in their *own* property, attorneys practicing law (Chapter
84), and court-appointed receivers/trustees/guardians/executors. *Note:* dealing in your own
property is exempt, but a *licensed* broker selling their own property is still bound by
material-fact disclosure duties (see Section 8).

---

## 2. The single-license / single-level system (the provisional-broker model)

This is the most-misunderstood feature of NC real estate, especially for people coming from
other states. **North Carolina is a single-license state: it issues only one kind of license —
a BROKER license. There is no separate "salesperson" or "sales agent" license.** (NCREC,
ncrec.gov "Apply for a License" / "Real Estate Licensing in NC", as of 2026-06.)

> Historical wrinkle to flag: GS Chapter 93A **Article 1 is still titled "Real Estate Brokers
> and Salespersons,"** and GS 93A-2 still references "salesperson" definitions, but in practice
> NCREC issues only broker licenses (the statutory "salesperson" tier is no longer issued). When
> reading the raw statute, treat "salesperson" as legacy terminology. (GS 93A Article 1 heading,
> ncleg.gov, as of 2026-06.)

Instead of two license *types*, NC uses one license with **status levels** for the same broker
license:

1. **Provisional Broker (PB).** The entry level. Every brand-new licensee starts here. A PB may
   perform brokerage only while **supervised by a Broker-in-Charge (BIC)** and only while
   on **active** status. The "provisional" label is removed once the PB completes the 90-hour
   Postlicensing education (Section 5).
2. **Broker (full / non-provisional).** After Postlicensing is complete, the "provisional"
   status is removed and the licensee is a full Broker. A full broker still needs to be
   affiliated with a BIC to operate actively, unless the broker is themselves a BIC or sole
   proprietor BIC.
3. **Broker-in-Charge (BIC) — a *designation*, not a separate license.** A BIC is a broker who
   has met extra eligibility requirements and been designated to supervise an office and the
   provisional/affiliated brokers in it (Section 6).

So the progression is: **Provisional Broker -> Broker -> (optionally) Broker-in-Charge
designation.** All three are the same underlying broker license at different stages.

---

## 3. License path — overview (numbered steps)

The standard path for a NC resident with no prior real-estate license (NCREC "Apply for a
License" and "Real Estate Licensing in NC", as of 2026-06):

1. **Meet the basic eligibility.** Be at least **18 years old** and be a U.S. citizen, a
   non-citizen national, a qualified alien, or otherwise lawfully present and authorized to work
   (NCREC, as of 2026-06; GS 93A-4 sets the 18-year minimum). You must also satisfy a
   **character/background check** assessing honesty, trustworthiness, and integrity.
2. **Complete the 75-hour Broker Prelicensing Course** from an approved provider (Section 4).
3. **Apply to NCREC** with the application and fee; the Commission reviews eligibility (including
   the background review) and, if you qualify, issues a **Notice of Examination Eligibility**
   with a **180-day examination-eligibility window**.
4. **Pass the license examination** (administered by PSI; Section 4).
5. **Receive a Provisional Broker license** — initially on **inactive** status (you cannot
   practice yet).
6. **Activate the license under a Broker-in-Charge** by affiliating with a firm/BIC and filing
   the **License Activation and Broker Affiliation form (REC 2.08)** — the BIC is responsible
   for its submission per Rule 58A .0506(b) (NCREC, as of 2026-06).
7. **Complete the 90-hour Postlicensing education within 18 months** of initial licensure to
   keep active eligibility and remove provisional status (Section 5; Rule 58A .1902).
8. **Complete 8 hours of Continuing Education each year** by June 10 to renew (Section 5; Rule
   58A .1702).

Steps 7 and 8 are independent: Postlicensing (one-time, 90 hours, 18-month deadline) and CE
(recurring, 8 hours/year) are *both* required and do not substitute for each other.

---

## 4. Getting licensed — prelicensing and the exam

**75-hour Broker Prelicensing Course.** You must complete a NCREC-approved **75-hour** Broker
Prelicensing Course before applying (NCREC "Prelicensing FAQ" / "Real Estate Licensing in NC",
as of 2026-06; the 75-hour figure tracks GS 93A-4(a), which requires "at least 75 hours of
instruction" within three years before application, or equivalent experience). NC prelicensing
is delivered by **live instruction**, but approved **distance-education providers** offer
online, self-paced versions. The course ends in a provider final exam you must pass.

**Out-of-state waiver.** An applicant holding a *current* license in another state that has been
on active status within the previous three years and is equivalent to NC's provisional or full
broker license may be able to **waive the 75-hour Prelicensing Course** (NCREC, as of 2026-06).
Equivalency and any state-portion exam requirements are determined by NCREC.

**The license examination (PSI).** The exam is administered by **PSI Services LLC** and has two
sections (NCREC "New NC Real Estate Broker License Examination" / State Exam Study Guidelines, as
of 2026-06):

- **National section** — pass with **71% or higher**.
- **State section** — **40 questions**, pass with **72.5% or higher**.

You schedule and sit the exam during the **180-day examination-eligibility period**. *[UNVERIFIED:
the exact current National-section question count — historically ~100-120 — should be confirmed on
the current PSI/NCREC candidate handbook.]*

**Fees.** Effective **April 1, 2026**, the original broker license application fee is **$105.00**
(NCREC, as of 2026-06). GS 93A-4 caps the application fee in a $100-$120 range; the PSI exam fee
and fingerprint/background-check fees are separate. **Verify all current fees on ncrec.gov before
budgeting.**

**After you pass:** the license is issued as a **Provisional Broker** on **inactive** status. You
cannot practice brokerage until you activate it under a BIC (Section 3, step 6). On inactive
status a broker must immediately cease all brokerage activity (NCREC "Active versus Inactive
License Status", as of 2026-06).

---

## 5. Postlicensing education and Continuing Education

These are two distinct, ongoing education obligations. Do not confuse them.

### 5a. Postlicensing — 90 hours within 18 months (one-time)

A Provisional Broker must complete **all 90 hours of Postlicensing education within 18 months of
the initial license date** under **Rule 58A .1902(b)** (NCREC, as of 2026-06). The 90 hours are
three separate **30-hour courses** (commonly "Post 301/302/303"). Key points:

- The **18-month clock runs from the license issuance date regardless of activation status**; it
  does not pause if you stay inactive.
- **No extensions are granted, for any reason** (NCREC, citing Rule 58A .1902).
- Each Postlicensing course must be completed (including passing the end-of-course exam) within
  **180 days of enrollment**.
- **Consequence of missing the deadline:** the license is placed on **inactive** status until
  activation criteria are met (**Rule 58A .1902(c)**). To later reactivate, all three courses must
  have been completed within the **two years immediately before** the activation request; courses
  older than two years must be retaken.
- **Completing all three courses removes "provisional" status**, so the licensee becomes a full
  Broker.

> **Important correction for anyone who heard "90 hours over 3 years":** that framing is wrong.
> The verified rule is **90 hours within 18 months** of initial licensure (Rule 58A .1902(b)). The
> "three years" number that floats around usually refers to GS 93A-4's rule that *Prelicensing*
> instruction must have been taken within three years before *application* — a different thing.

### 5b. Continuing Education — 8 hours per license year (recurring)

After the first renewal cycle, **every NC broker (resident or not) must complete 8 hours of CE
each license year**, due by **June 10**, under **Rule 58A .1702** (NCREC, as of 2026-06). The 8
hours are:

- **A 4-hour Commission-mandated Update course**, plus
- **4 hours of a Commission-approved elective.**

Which Update course depends on status:

- **General Update (GENUP)** — for brokers who are *not* BIC or BIC-Eligible.
- **Broker-in-Charge Update (BICUP) / Broker-in-Charge Annual Review (BICAR)** — required for
  brokers who are **BIC** or **BIC-Eligible**. (NCREC has updated the BIC update-course branding
  over time; recent material uses **BICAR**. Confirm the current course name each license year.)

The license year ends June 30, and CE for that year is due by **June 10**. A new licensee's first
CE applies for the renewal *after* they have held the license, working alongside the separate
Postlicensing requirement. Failing to complete CE by June 10 results in the license going
**inactive** until CE is made up (NCREC "Continuing Education Information", as of 2026-06).

---

## 6. Firm licensing and the Broker-in-Charge (BIC)

### 6a. Firm license

**Every business entity other than a sole proprietorship that operates as a real-estate broker
for compensation must obtain its own firm license** before doing business, under **Rule 58A
.0502** (NCREC, as of 2026-06). This covers corporations, LLCs/PLLCs, partnerships, and similar
entities.

- The firm must designate a **Qualifying Broker (QB)** — a principal of the firm who **holds an
  active broker license in good standing** (Rule 58A .0502). The QB is responsible for ensuring
  there is a **BIC for each office location** of the firm.
- **LLC applications** must include the operating agreement (or a written description of the
  managers' rights and duties) and the name of each manager; the application is marked incomplete
  without it (NCREC, as of 2026-06).
- **Narrow no-BIC exception (Rule 58A .0110(c)):** a licensed firm need not designate a BIC if it
  exists *solely* to receive compensation, is an IRS pass-through, has no principal/branch office,
  and has no licensee other than its QB. (A common structure for a broker's personal
  compensation-only entity.) Entity choice/tax for such an entity -> see
  `venture-nc-business-formation-tax`.

### 6b. Broker-in-Charge (BIC) eligibility and designation

The BIC is the supervising broker for an office. BIC is a **designation on a broker license**,
not a separate license (NCREC "BIC Eligible Status and Declaration", as of 2026-06).

**Eligibility (Rule 58A .0110).** To be **BIC-Eligible**, a broker must:

1. Hold a **non-provisional (full) broker license on active status**.
2. Have at least **two years of full-time (or equivalent) real-estate brokerage experience**
   within the previous five years (NCREC describes this as roughly 40 hours/week-equivalent for
   two years within the prior five).
3. **Complete the 12-hour Broker-in-Charge (BIC) Course** — no earlier than **one year before**,
   and no later than **120 days after**, submitting the BIC request form.
4. Submit the **Request for BIC Eligible Status and/or BIC Designation (Form REC 2.25)**.

**Maintaining BIC eligibility:** renew the broker license on time each year **and** complete the
required **Update course (BICUP/BICAR)** plus elective by **June 10** annually (Section 5b). Lapsing
either can cause loss of BIC-Eligible status.

**BIC responsibilities (Rule 58A .0110, .0506; BIC Best Practices Guide).** A BIC is responsible
for supervising the office's brokerage practice, including:

- Submitting **License Activation/Affiliation forms (REC 2.08)** for affiliated provisional
  brokers (Rule 58A .0506(b)).
- Ensuring affiliated brokers have **active licenses and completed their CE/Postlicensing**.
- **Overseeing trust-account handling**, advertising, agency agreements, and disclosure
  compliance for the office.
- Maintaining required records and supervising the conduct of affiliated brokers and unlicensed
  staff.

---

## 7. Trust / escrow accounts and recordkeeping

Brokers who hold money belonging to others (earnest-money deposits, rents, POA funds, etc.) are
**fiduciaries** and must handle that money under **Rules 58A .0116, .0117, and .0118** (NCREC
"Opening a Real Estate Broker Trust Account" / NC Real Estate License Law and Rules, as of
2026-06). Core requirements:

- **Use a designated trust/escrow account (Rule 58A .0117).** It must be a **demand-deposit
  account** at a **federally insured depository institution doing business in NC** that agrees to
  make the records available to NCREC for inspection. The account, **bank statements, deposit
  tickets, and checks must bear the words "Trust Account" or "Escrow Account."**
- **Deposit promptly (Rule 58A .0116).** Trust money must be deposited **no later than three (3)
  banking days following receipt**, unless a different time frame applies under Rule 58A .0116(b).
- **No commingling.** A broker may not commingle clients' trust money with the broker's own funds
  (commingling and failure to deposit trust money are explicit disciplinary grounds; GS 93A-6).
- **Recordkeeping (Rule 58A .0117).** Maintain records that show the **ownership of every dollar
  of trust money from deposit through final disbursement**, sufficient for NCREC to verify
  accuracy and proper use; all records must be available for NCREC inspection.
- **Property Owners' Association (POA) funds (Rule 58A .0118).** When a broker collects/controls
  POA funds, those funds are **trust money** and must be kept in a trust/escrow account
  **dedicated exclusively to a single POA**, not commingled with other associations' or persons'
  funds.
- **Disputed funds.** When parties dispute entitlement to trust money (e.g., contested earnest
  money), a broker generally should not unilaterally disburse; NC guidance allows surrendering
  disputed funds to the **Clerk of Court** in appropriate cases (NCREC, as of 2026-06).

Trust-account violations are among the most serious and most-disciplined license-law issues. The
BIC bears supervisory responsibility for the office's trust accounting (Section 6b).

---

## 8. NC agency law and disclosures

NC agency law is built on the **Working With Real Estate Agents (WWREA) disclosure**, the rule
that **all agency agreements be in writing** (in the end), and the dual/designated-agency
framework — all anchored in **Rule 58A .0104**.

### 8a. The Working With Real Estate Agents (WWREA) disclosure — timing

Under **Rule 58A .0104(c)**, in every real-estate **sales** transaction a broker must, **at first
substantial contact** with a prospective buyer or seller, provide that person a copy of the
NCREC publication **"Working With Real Estate Agents."** (NCREC "Reminder: Working with Real
Estate Agents Disclosure", as of 2026-06.)

- **"First substantial contact"** is the point at which discussion moves beyond casual greeting
  toward specifics (the consumer begins to disclose needs/motivation, or the broker begins to
  advise), i.e., *before* the consumer reveals confidential information. NCREC does not give a
  rigid definition; brokers are expected to err on the early side.
- **Format (revised 2021):** the disclosure is a **single page, double-sided — one side for
  sellers, one side for buyers** — accompanied by a longer explanatory brochure.
- The disclosure **educates the consumer about agency options and is not itself an agency
  contract.** After it is delivered, a proper agency relationship is formed orally or in writing
  per Rule 58A .0104.

### 8b. Types of agency

- **Seller's agent (listing agent):** represents the seller.
- **Buyer's agent:** represents the buyer.
- **Seller subagency:** a cooperating broker represents the *seller* (not the buyer); now
  uncommon but still permitted with proper disclosure/authority.
- **Dual agency:** the firm represents both buyer and seller in the *same* transaction (Section
  8c).
- **Designated dual agency:** a refinement of dual agency (Section 8c).

### 8c. Dual agency and designated dual agency

- **Dual agency** is allowed in NC **only with the written authorization of both parties** in
  their agency agreements (Rule 58A .0104). In pure dual agency the firm (and potentially the same
  broker) represents both sides; the broker must stay neutral and **may not disclose one party's
  price, terms, motivation, or other confidential information to the other.**
- **Designated dual agency** is NC's preferred refinement: within one firm, **one broker is
  designated to represent only the seller and a different broker only the buyer**, so each client
  gets an individual advocate while the firm represents both. Each designated agent acts only
  in their own client's interest and may not share that client's confidential information (NCREC
  "Case Study: Designated Dual Agency", as of 2026-06).
- **BIC conflict rules:** a BIC may serve as a designated dual agent **only if a non-provisional
  broker is the designated agent for the competing party**, and a **BIC may not be a designated
  dual agent when a provisional broker they supervise is the designated agent for the other side**
  (supervisory conflict). When the BIC is a designated dual agent, they must set up
  policies/procedures (e.g., a separate broker reviewing records) to wall off the other party's
  confidential information.

### 8d. Buyer-agency agreements and the 2024 NAR settlement

- **NC has long required agency agreements to be in writing** under **Rule 58A .0104**. A buyer
  agency agreement may begin **orally**, but it **must be reduced to writing no later than the time
  any offer is made/written** (Rule 58A .0104). Required contents include a **definite term, the
  broker's license number, a termination provision, and nondiscrimination (Fair Housing)
  language**.
- **2024 NAR settlement overlay:** effective **August 17, 2024**, NAR's settlement requires MLS
  participants working with a buyer to enter a **written buyer agreement *before touring* a home**,
  and **prohibits offers of buyer-broker compensation in the MLS** (compensation may still be
  negotiated and offered *outside* the MLS). (NCREC "Has the World Exploded? The NAR Settlement,
  Commission Law and Rules", Oct 2024 eBulletin, as of 2026-06.)
- **Net effect in NC:** NC rule already required *written* agency agreements (by offer time); the
  settlement pushes the *written buyer agreement earlier* (before showings) and **decouples buyer-
  agent compensation from the MLS.** Practitioners should now use a **signed written buyer-agency
  agreement that addresses compensation before showing property.** *[UNVERIFIED: whether NCREC has
  amended Rule 58A .0104 itself to codify the "before touring" timing, vs. it being an
  NAR/MLS-membership obligation — confirm the current rule text on ncrec.gov.]*

---

## 9. RPOADS and the broker's material-fact disclosure duty

### 9a. The Residential Property and Owners' Association Disclosure Statement (RPOADS)

Under **NC General Statutes Chapter 47E** (the Residential Property Disclosure Act), a **seller**
of residential property of **one to four units** (by sale, exchange, installment land contract, or
option) must give the buyer the **Residential Property and Owners' Association Disclosure Statement
(RPOADS)** — NCREC form REC 4.22 (NCREC "Sellers Required by Law to Provide Two Disclosure
Statements to Buyers", as of 2026-06). Key features:

- For each item the seller may answer **"Yes," "No," or "No Representation."** Because of the **"No
  Representation"** option, the form is mandatory to *deliver* but does not force the seller to
  *reveal* anything specific.
- A seller who answers **"No" when the truthful answer is "Yes"** can face **civil liability for
  misrepresentation**; "No Representation" carries no such representation risk.
- GS Chapter 47E also requires a separate **Mineral and Oil and Gas Rights Mandatory Disclosure
  Statement** in covered residential transactions (the "two disclosure statements").

### 9b. The broker's independent material-fact duty (the critical point)

**Regardless of what the seller does on the RPOADS, a broker must disclose all known material
facts to a prospective buyer.** This is a license-law duty separate from the seller's GS 47E
form (NCREC "Did You Know? Brokers must disclose material facts on their own properties" and
"Sellers' Obligation to Disclose Latent Defects", as of 2026-06):

- A **material fact** is anything that, if known, might affect the buyer's decision or the price
  (physical defects, latent defects, and certain off-site or transaction facts).
- A **listing agent has an affirmative duty to discover and disclose** material facts; NCREC
  holds brokers responsible for facts they **know or reasonably should know**. Even if the seller
  omits a defect from the RPOADS (or answers "No Representation"), the **broker must still disclose
  a known material fact to buyers/their agents before a sales contract is formed.**
- **"As-is" does NOT excuse disclosure.** Selling "as-is" lets the seller decline to *repair*; it
  does **not** permit the seller or broker to **conceal** material facts. A seller can answer "No
  Representation," but a broker — **including a broker selling their own property** — must always
  disclose known material facts.
- **Failing or willfully/negligently omitting a material fact is an explicit disciplinary ground**
  under **GS 93A-6(a)(1)** (see Section 10).

---

## 10. Fair housing

Two layers apply, and NCREC enforces compliance against licensees in addition to the dedicated
fair-housing agencies.

**Federal — Fair Housing Act** (Title VIII of the Civil Rights Act of 1968, as amended by the
1988 Fair Housing Amendments Act). It prohibits discrimination in the sale, rental, financing, and
**advertising** of housing based on the federal **protected classes: race, color, religion, sex,
national origin, familial status, and disability (handicap)** (NCREC Fair Housing Brochure /
FH Laws, as of 2026-06). *Note:* HUD enforces "sex" to include sexual orientation and gender
identity under current guidance — confirm current HUD/agency posture, as enforcement
interpretations can shift.

**State — NC State Fair Housing Act, GS Chapter 41A.** NC's act parallels the federal protected
classes and is **enforced by the NC Human Relations Commission** (NCREC, as of 2026-06).

**Advertising and practice rules for brokers.** All real-estate advertising and conduct must be
**free of language or images suggesting a preference, limitation, or exclusion** based on a
protected class:

- Describe **the property, not who should live there.** Avoid steering language such as "perfect
  for young professionals," "great for a Christian family," or "no children."
- Fair-housing obligations extend across listings, social media, property management, and tenant
  screening. Brokers practicing **property management** should re-check fair-housing rules on
  occupancy, assistance animals/reasonable accommodation, and screening.
- NCREC can **discipline licensees** for discriminatory advertising or conduct (Section 10), in
  addition to federal/state fair-housing enforcement.

---

## 11. NCREC advertising rules (blind ads, team names)

Advertising of real estate and brokerage services is governed by **Rule 58A .0105** (with assumed-
name rules in **Rule 58A .0103**) (NCREC "Proper Use of Names in Advertising", as of 2026-06):

- **No blind ads.** A **"blind ad"** is one that does not disclose that the advertiser is a real-
  estate **broker** (it makes the offer look like it comes from a private party/principal only).
  Brokers' ads **must disclose the broker/firm's involvement.**
- **Firm name required in the ad.** A broker's advertisement must include **the name of the firm
  (or sole proprietorship) the broker is affiliated with.** On **social media**, the firm name
  must appear in **the actual ad content**, not merely in a bio/profile.
- **Authority to advertise.** A broker generally may advertise a specific property **only with the
  authority of the listing firm/owner** (don't advertise another firm's listing without
  permission).
- **Team and assumed (DBA) names (Rule 58A .0103).** A licensed firm may operate under one or more
  **assumed/team names** (e.g., "Team One Realty"), but it must **file an assumed-business-name
  certificate** in compliance with NC law **and notify NCREC in writing** of the assumed name.
  Team-name advertising must still make the affiliated **firm** identifiable.
- **Special identifiers.** Nonresident commercial licensees must conspicuously identify themselves
  as a **"Limited Nonresident Commercial Real Estate Broker"** in advertising. *[UNVERIFIED: the
  exact current sub-section lettering within Rule 58A .0105 for blind-ad vs. firm-name provisions —
  confirm against the current NC Real Estate License Law and Rules booklet.]*

General marketing-compliance rules that aren't NCREC-specific (CAN-SPAM, TCPA texting/calling
consent, FTC endorsement/fake-review rules, Google Business Profile) -> see
`venture-marketing-strategy-local-seo`.

---

## 12. Discipline and license-law violations

NCREC can discipline any licensee who violates GS Chapter 93A or Commission rules, after a
hearing, under **GS 93A-6** (ncleg.gov, as of 2026-06).

**Sanctions (GS 93A-6).**

- **Reprimand / censure**: formal discipline short of loss of license.
- **Suspension**: temporary; may be **active** (no brokerage allowed during it) or **stayed**
  (brokerage allowed subject to conditions), or a combination.
- **Revocation**: the Commission takes the license away.
- NCREC may also require **education** or impose conditions.

**Common grounds (GS 93A-6(a)).** Among the enumerated grounds:

- **Willful or negligent misrepresentation, or willful/negligent omission of a material fact**
  (GS 93A-6(a)(1)) — ties directly to the disclosure duty in Section 9.
- Making false promises or pursuing a course of misrepresentation/false promises.
- Acting for more than one party without the knowledge/consent of all (improper undisclosed dual
  agency).
- **Commingling** others' funds with the broker's own, or **failing to deposit** trust money in a
  trust account; failing within a reasonable time to **account for or remit** money belonging to
  others (Section 7).
- Performing **unauthorized legal services** (practicing law).
- **Conviction of a crime** involving moral turpitude, fraud, theft, embezzlement, or forgery; or
  obtaining a license by fraud.
- **Improper, fraudulent, or dishonest dealing / incompetence** endangering the public.
- **Violating any Commission rule** (21 NCAC 58).

**Real Estate Education and Recovery Fund (GS 93A Article 2).** A consumer who suffers monetary
loss from a licensee's fraud, misrepresentation, conversion of trust funds, or similar misconduct,
and who obtains an unsatisfied judgment, may seek **reimbursement from the Recovery Fund** (per-
transaction/per-licensee statutory caps apply; verify current limits in GS 93A Article 2). The
Fund protects consumers; it is **not** insurance for the broker, and NCREC may pursue the broker
for reimbursement.

**If a complaint is filed against you:** respond timely and truthfully to NCREC; failing to
respond to a Commission inquiry can itself be a violation. Consider counsel for any matter that
could threaten the license.

---

## Sources

Primary sources, all consulted as of **2026-06**:

1. **NCREC — Apply for a License / Real Estate Licensing in North Carolina** (single broker
   license; 75-hour Prelicensing; out-of-state waiver; $105 fee eff. 2026-04-01):
   https://www.ncrec.gov/Licensing/ApplyLicense ; https://www.ncrec.gov/Brochures/general.pdf
2. **NCREC — Prelicensing FAQ / Course info:** https://www.ncrec.gov/Education/PreFAQ
3. **NCREC — License examination changes / State Exam Study Guidelines** (PSI; National 71%,
   State 72.5% / 40 questions; 180-day eligibility):
   https://bulletins.ncrec.gov/new-nc-real-estate-broker-license-examination/ ;
   https://www.ncrec.gov/Pdfs/Licensing/StateExamStudyGuidelines.pdf
4. **NCREC — Postlicensing education** (90 hrs / 3x30; 18-month deadline; Rule 58A .1902):
   https://bulletins.ncrec.gov/what-do-pbs-and-their-bics-need-to-know-about-postlicensing-education/ ;
   https://bulletins.ncrec.gov/reminder-have-you-completed-your-postlicensing-education/
5. **NCREC — Continuing Education Information** (8 hrs/yr; GENUP/BICUP/BICAR + elective; June 10;
   Rule 58A .1702): https://www.ncrec.gov/Education/General
6. **NCREC — BIC Eligible Status and Declaration / All about BIC Eligibility / 12-hour BIC Course
   FAQ / BIC Best Practices Guide** (Rule 58A .0110; Form REC 2.25; Rule 58A .0506):
   https://www.ncrec.gov/BICEligibleReqAndDesignation ;
   https://bulletins.ncrec.gov/all-about-bic-eligibility-rules/ ;
   https://www.ncrec.gov/Education/BICCRSFAQ ; https://www.ncrec.gov/Pdfs/bicguide.pdf
7. **NCREC — Firm license** (Rule 58A .0502 QB; Rule 58A .0110(c) no-BIC exception):
   https://bulletins.ncrec.gov/do-i-need-a-firm-license/ ;
   https://bulletins.ncrec.gov/how-to-apply-for-a-firm-license/
8. **NCREC — Trust accounts** (Rules 58A .0116/.0117/.0118; 3 banking days; "Trust Account"
   labeling; POA funds): https://bulletins.ncrec.gov/opening-a-real-estate-broker-trust-account/ ;
   https://www.ncrec.gov/Pdfs/Rules/NCRECLawAndRules.pdf
9. **NCREC — WWREA disclosure** (first substantial contact; one-page two-sided; Rule 58A .0104(c)):
   https://bulletins.ncrec.gov/working-with-real-estate-agents-disclosure/ ;
   https://www.ncrec.gov/Forms/WWREA/WWREADisclosureForm.pdf
10. **NCREC — NAR settlement and NC agency law** (Rule 58A .0104 written agreements; Aug 17 2024;
    MLS compensation decoupling):
    https://bulletins.ncrec.gov/has-the-world-exploded-the-nar-settlement-commission-law-and-rules/ ;
    https://bulletins.ncrec.gov/nar-settlement/
11. **NCREC — Designated dual agency** (one firm/two designated agents; BIC conflict rules):
    https://bulletins.ncrec.gov/case-study-designated-dual-agency/ ;
    https://www.ncrec.gov/Pdfs/genupdate/2022-2023%20Section%202%20Dual%20Agency.pdf
12. **NCREC — RPOADS and material-fact duty** (GS 47E; "No Representation"; broker affirmative
    duty; as-is): https://bulletins.ncrec.gov/sellers-required-by-law-to-provide-two-disclosure-statements-to-buyers/ ;
    https://bulletins.ncrec.gov/did-you-know-brokers-must-disclose-material-facts-on-their-own-properties/ ;
    https://www.ncrec.gov/Forms/Consumer/rec422.pdf
13. **NCREC — Fair housing** (federal FHA classes; NC GS Ch. 41A; advertising):
    https://www.ncrec.gov/Brochures/FairHousingBrochure.pdf ; https://www.ncrec.gov/FHResources/FHLaws
14. **NCREC — Advertising rules** (Rule 58A .0105 blind ads/firm name; Rule 58A .0103 assumed/team
    names): https://bulletins.ncrec.gov/proper-use-of-names-in-advertising/ ;
    http://www.ncrec.gov/Pdfs/bicar/Advertising.pdf
15. **NCREC — Discipline / Recovery Fund** (GS 93A-6 sanctions; Article 2 Recovery Fund):
    https://bulletins.ncrec.gov/disciplinary-action/ ;
    https://bulletins.ncrec.gov/real-estate-education-and-recovery-fund-reimburses-victims/
16. **NC General Statutes Chapter 93A** (Articles 1-4; GS 93A-1/-2/-3/-4/-6), ncleg.gov:
    https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/ByChapter/Chapter_93A.html
17. **NC Real Estate License Law and Commission Rules booklet** (21 NCAC 58):
    https://www.ncrec.gov/Pdfs/Rules/NCRECLawAndRules.pdf

---

## Disclaimer

This skill is **general information, not legal advice**, and does not create an
attorney-client relationship. Real-estate licensing requirements, education hours, fees, forms,
rules, and statutes **change** — sometimes mid-year — and several specific figures here are
flagged **[UNVERIFIED]** pending confirmation against the current source. Before relying on
anything in this skill: (1) verify current rules, hours, fees, and forms directly at
**ncrec.gov** and in the current **GS Chapter 93A** (ncleg.gov) and **21 NCAC 58** rules booklet;
(2) for any disciplinary matter, transaction question, or entity/contract issue, consult a
**licensed North Carolina attorney**; and (3) treat NCREC as the controlling authority on all
licensing and brokerage-law questions. Entity formation and tax are out of scope (see
`venture-nc-business-formation-tax`), as are marketing/investing strategy (see
`venture-real-estate-marketing-investing`) and general marketing compliance (see
`venture-marketing-strategy-local-seo`).
