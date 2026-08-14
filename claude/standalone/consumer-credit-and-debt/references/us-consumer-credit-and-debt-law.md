<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `us-consumer-credit-and-debt-law` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: us-consumer-credit-and-debt-law
description: >-
  US FEDERAL consumer-credit & debt statute reference — what each law does,
  consumer rights, enforcers (CFPB, FTC). General legal INFORMATION, not legal
  advice. TRIGGER: FCRA (accuracy, dispute/reinvestigation, permissible purpose,
  adverse action), FDCPA + Regulation F (third-party
  collectors, validation, 7-in-7 calls, time-barred debt), ECOA/Reg B (credit
  discrimination), TILA/Reg Z (APR, disclosures, rescission), FCBA (card billing
  errors), CROA (advance-fee ban), EFTA/Reg E (EFTs); the CONTESTED 2024 CFPB
  medical-debt rule; AI/algorithmic fair lending. SKIP: NC statutes ->
  north-carolina-credit-and-debt-law; dispute STEPS / score tactics ->
  credit-reports-and-scores / improving-and-rebuilding-credit; settling
  charge-offs & collections -> charge-offs-collections-and-debt-resolution;
  collector-interaction tactics -> debt-collectors-and-fdcpa-rights; coding
  credit-law compliance -> software-engineering-patterns; medical bills / billing
  disputes / charity care -> medical-debt-and-billing.
metadata:
  changelog:
    - "2026-06-16 sko (initial author + optimize) — Pass H 10/10 pos, 1/10 neg (predicted); fixed 3 High (desc >1000 chars; FDCPA 30-day clock now runs from receipt; ECOA 'CCPA' disambiguated to federal Consumer Credit Protection Act) + 11 Medium (per-statute + deadlines tables, EFTA $500 condition, medical-debt vacatur wording, output-rule block, code SKIP carve-out, em-dash density 2.13->0.91/100w, citations); 0 banned terms; blind-audit CLEAN"
    - "2026-06-16 sko (--meta --no-sync) — 1 Medium structural: added cross-ref SKIP -> medical-debt-and-billing (this skill owns the vacated 2024 CFPB medical-debt rule's statute side; practical medical bills/billing/charity-care layer routes there). Cross-topic, no circular SKIP. Desc trimmed to <=1000 (1000). consumer-credit-and-debt anchor unchanged. Linter 0/0/0. No content passes; --no-sync (registry left stale)."
---

# US Consumer Credit & Debt Law (Federal Framework)

> **Spoke** of the **`consumer-credit-and-debt`** hub. This is the federal
> statutory backbone the other spokes reference for the *law itself*.

## LEGAL FRAMING — read first

This skill is **general legal information, NOT legal advice.** It does not create
an attorney-client relationship. Consumer-credit law is **fact-specific** and
**changes frequently** (statutes are amended; agencies issue, revise, and have
rules **vacated** by courts; dollar thresholds are inflation-adjusted). Content
is **current as of 2026** and may be stale by the time you read it. Statutory
dollar amounts and damage caps cited here are baselines that may have been
adjusted. **Verify the current text** against the primary sources in the last
section, and for any actual dispute, claim, deadline, or lawsuit, **consult a
licensed attorney** in the relevant jurisdiction or contact the **CFPB**
(consumerfinance.gov) or **FTC** (consumer.ftc.gov). State law often adds
stronger protections; federal law is generally a floor, not a ceiling.

Citations: **USC** = U.S. Code (the statute, via Cornell LII). **CFR** = Code of
Federal Regulations (the implementing regulation, via eCFR / CFPB).

---

## The enforcers — CFPB and FTC

- **CFPB (Consumer Financial Protection Bureau)** — created by the Dodd-Frank Act
  (2010). Primary federal **rule-writer** for most of these statutes (it now
  issues Regulations B, E, F, V, Z, etc.), supervises large banks and many
  nonbanks, **enforces**, and runs the public **complaint database** at
  consumerfinance.gov/complaint. NOTE (2026 volatility): the CFPB's funding,
  staffing, and rulemaking posture have been the subject of litigation and
  policy change; verify its current authority/activity before relying on a
  recent rule or guidance.
- **FTC (Federal Trade Commission)** — enforces the FDCPA, FCRA, CROA, and
  others against entities outside CFPB's reach; brings UDAP actions under FTC Act
  §5; consumer education at consumer.ftc.gov; report fraud at reportfraud.ftc.gov.
- **State Attorneys General** and **private rights of action** also enforce most
  of these statutes (the FCRA, FDCPA, TILA, ECOA, and EFTA all allow consumers
  to sue, with statutory damages and attorney's fees in many cases).

---

## How to answer with this skill

When answering from this material:

- **Cite the section.** Pair every rule with its USC/CFR anchor (e.g. "FDCPA
  §1692g", "Reg Z §1026.23").
- **Flag staleness on every number.** Any deadline, dollar amount, or damage cap
  you surface must carry the caveat that it may have been amended or
  inflation-adjusted; point the reader to the matching primary source in the
  References section.
- **Never assert the volatile items as settled.** Do not state the 2024 CFPB
  medical-debt rule is in force, and do not present disparate-impact doctrine or
  AI-underwriting guidance as fixed law; flag them as contested and tell the
  reader to verify.
- **Route, don't guess.** Send state-specific, how-to/dispute-steps,
  debt-settlement, collector-interaction, and software-implementation questions to
  the spokes named in the closing Cross-references section.
- **Repeat the "not legal advice" framing** when the question is about a real
  situation, and point to a licensed attorney or the CFPB.

## Statute quick reference

One row per core statute (citation, implementing regulation, who enforces,
private-right-of-action remedy). Dollar caps are statutory baselines that may be
inflation-adjusted; verify before relying.

| Statute | USC | Implementing reg | Enforcers | Private right / remedy |
|---|---|---|---|---|
| FCRA | 15 USC §1681 | Reg V (12 CFR 1022) | CFPB, FTC, state AGs | Yes — §1681n willful / §1681o negligent; actual + statutory + fees |
| FDCPA | 15 USC §1692 | Reg F (12 CFR 1006) | CFPB, FTC, state AGs | Yes — §1692k; actual + up to $1,000 statutory + fees; 1-yr SOL |
| ECOA | 15 USC §1691 | Reg B (12 CFR 1002) | CFPB, FTC, DOJ, prudential regulators | Yes — §1691e; actual + punitive (≤ $10,000 individual; capped class) + fees |
| TILA | 15 USC §1601 | Reg Z (12 CFR 1026) | CFPB, FTC | Yes — §1640; actual + statutory + fees; rescission (§1635) |
| FCBA | 15 USC §1666 | Reg Z (12 CFR 1026) | CFPB, FTC | Yes — via TILA §1640; forfeiture of disputed amount on noncompliance |
| CROA | 15 USC §1679 | (statute; no separate reg) | CFPB, FTC | Yes — §1679g; actual + punitive + fees; contracts voidable |
| EFTA | 15 USC §1693 | Reg E (12 CFR 1005) | CFPB, FTC | Yes — §1693m; actual + statutory + fees |

## Key deadlines & dollar amounts

Highest-searched figures, gathered. Each is subject to the staleness caveat
above; the substantive rules and conditions are in the per-statute sections.

| Item | Figure | Source |
|---|---|---|
| FCRA reinvestigation | 30 days (+15 = 45 if consumer adds info) | §1681i(a)(1) |
| FCRA results notice to consumer | 5 business days after completion | §1681i(a)(6) |
| FDCPA validation notice sent | within 5 days of initial communication | §1692g(a) |
| FDCPA dispute window | 30 days after **receiving** the notice | §1692g(a)(3) |
| Reg F call-frequency presumption | > 7 calls / 7 days, or within 7 days of a conversation | §1006.14(b) |
| FDCPA statutory damages | up to $1,000 | §1692k |
| TILA right of rescission | 3 business days (up to 3 years if disclosures omitted) | §1635 |
| FCBA billing-error notice | 60 days after statement sent | §1666(a) |
| FCBA resolution | acknowledge ≤ 30 days; resolve ≤ 2 billing cycles / 90 days | §1666(a) |
| ECOA adverse-action notice | generally 30 days | §1691(d) |
| ECOA damages | punitive ≤ $10,000 individual (class capped) | §1691e |
| EFTA unauthorized-transfer liability | $50 / $500 / unlimited (time-tiered) | §1693g; Reg E §1005.6(b) |

---

## FCRA — Fair Credit Reporting Act (15 USC §1681 et seq.)

Governs **consumer reporting agencies (CRAs / "credit bureaus")**, the
**furnishers** who report data to them, and the **users** who pull reports.
Implemented by **Regulation V (12 CFR 1022)**. Enforced by CFPB and FTC; private
right of action (§1681n willful, §1681o negligent).

**What it does / key consumer rights:**

- **Permissible purpose (§1681b).** A CRA may furnish a report only for a listed
  purpose (credit, insurance, employment *with written consent*, a court order,
  legitimate business need tied to a consumer-initiated transaction, etc.).
  Pulling a report without permissible purpose is unlawful.
- **Accuracy & the dispute / reinvestigation process (§1681i).** On a consumer
  dispute, the CRA must conduct a **reasonable reinvestigation** free of charge
  and complete it **within 30 days** (extendable to **45 days** if the consumer
  submits additional info during the initial period). It must forward relevant
  info to the furnisher, **review all relevant information** the consumer
  submits, and **delete or modify** any item that is inaccurate, incomplete, or
  cannot be verified. Results notice to the consumer **within 5 business days**
  of completion. Frivolous-dispute terminations require notice within 5 business
  days.
- **Furnisher duties (§1681s-2).** Furnishers must not report information they
  know or have reasonable cause to believe is inaccurate, must correct and update,
  and must investigate disputes forwarded by a CRA. Note: the *direct* private
  right of action against furnishers is limited; most consumer suits run through
  the CRA-forwarded-dispute path of §1681s-2(b). (See CFPB Circular 2022-07 on
  what a "reasonable investigation" requires.)
- **Adverse-action notice (§1681m).** A user who takes adverse action based on a
  consumer report must tell the consumer, name the CRA used, and state the right
  to a free report and to dispute. (Credit-score-based denials also trigger a
  risk-based pricing / credit-score disclosure under §1681m(h) and Reg V
  subpart H.)
- **Free disclosures (§1681g, §1681j).** Consumers get file disclosures and at
  least one free report per year (AnnualCreditReport.com is the federally
  mandated source; weekly free access has been available in recent years; verify
  the current cadence against the FTC/CFPB sources in the References section).
- **Other rights:** security freezes and one-call fraud alerts (added by later
  amendments / FACTA), identity-theft blocking, and limits on reporting of stale
  negative items (generally **7 years**; bankruptcies **10 years**) under §1681c.

**Consumer-rights checklist (FCRA):**
- [ ] Get your reports (AnnualCreditReport.com) and check each bureau separately.
- [ ] Dispute errors **in writing** to the **CRA** (creates the §1681i duty);
      disputing to the furnisher can also trigger duties but the CRA path
      preserves the strongest remedy.
- [ ] Keep proof of mailing/dates; the 30/45-day clock matters.
- [ ] Demand the adverse-action notice / CRA name after any credit denial.
- [ ] Consider a security freeze (free) if at risk of identity theft.

---

## FDCPA — Fair Debt Collection Practices Act (15 USC §1692 et seq.)

**Scope:** applies to **third-party debt collectors** and debt buyers collecting
**consumer** (personal/household) debts; *not*, in general, to the **original
creditor** collecting its own debt in its own name (§1692a(6) definition of "debt
collector"). Enforced by CFPB and FTC; private right of action (§1692k:
actual + up to **$1,000** statutory damages + attorney's fees; 1-year statute of
limitations).

**Key prohibited practices:**
- **§1692c — communication restrictions:** no contact at unusual/inconvenient
  times or places (presumed 8 a.m.–9 p.m. local), at work if the collector knows
  the employer prohibits it, or with the consumer **after a written cease
  request** (except limited notices); must go through counsel if represented.
- **§1692d — harassment or abuse:** no threats of violence, obscene language,
  repeated phone calls to annoy, publishing "deadbeat" lists.
- **§1692e — false or misleading representations:** no falsely implying
  attorney/government affiliation, misstating the amount or legal status of a
  debt, or threatening action that can't legally be taken or isn't intended.
- **§1692f — unfair practices:** no collecting amounts not authorized by the
  agreement or law, no deceptive postcards.
- **§1692g — debt validation:** within **5 days** after the initial
  communication the collector must send a written validation notice (amount,
  creditor, the consumer's right to dispute), **unless that information was in
  the initial communication** or the debt is already paid. The consumer then has
  **30 days after receiving the notice** (§1692g(a)(3); Reg F §1006.34(b)(5)
  lets the collector assume receipt 5 business days after sending) to dispute
  **in writing**; if they do, the collector must **cease collection until it
  mails verification.**

### Regulation F — 12 CFR Part 1006 (effective Nov. 30, 2021)

The CFPB's modern rule **implementing the FDCPA**. Key additions:

- **Validation notice content (§1006.34):** standardized disclosures —
  itemization of the debt as of an "itemization date," the current amount,
  creditor name, a tear-off dispute form, and a plain statement of consumer
  rights. A **model validation notice** in Appendix B provides a safe harbor.
- **"7-in-7" call-frequency presumption (§1006.14(b)).** A collector is
  **presumed to violate** the harassment prohibition if it places telephone
  calls **more than seven (7) times within seven (7) consecutive days**, *or*
  **within seven (7) consecutive days after having had a telephone conversation**
  with the person about that debt. (Per-debt; the presumption can be rebutted.)
- **Limited-content message (§1006.2(j)).** Defines a voicemail/message a
  collector may leave that is **not** treated as a "communication" (so it can be
  overheard without violating third-party-disclosure rules); it may contain only
  specified content (consumer name, a request to reply, business name that does
  not indicate it's a debt collector, etc.).
- **Time-barred debt (§1006.26).** A collector may **not sue or threaten to sue**
  on a debt it knows or should know is **time-barred** (past the state statute of
  limitations, usually 3–6 years, varying by state and debt type). The rule also
  governs any required time-barred-debt disclosure on the validation notice.
- **Electronic communications:** Reg F sets opt-out and convenience rules for
  email and text, and addresses social-media contact.

**Consumer-rights checklist (FDCPA / Reg F):**
- [ ] You can demand validation **in writing within 30 days** — collection must
      pause until they verify.
- [ ] You can tell a collector to **stop contacting you** in writing (they may
      still notify you of specific actions).
- [ ] Excessive calls (the 7-in-7 line) and any threat to sue on old
      (time-barred) debt are red flags; document dates/times.
- [ ] Never assume a payment or new promise won't **restart** a state limitations
      clock; verify under your state's law before paying old debt.

---

## ECOA — Equal Credit Opportunity Act (15 USC §1691 et seq.) / Reg B (12 CFR 1002)

**Anti-discrimination in *any* aspect of a credit transaction.** Implemented by
**Regulation B (12 CFR Part 1002)**. (Older citations to 12 CFR 202 are the
superseded Federal Reserve version; rulemaking moved to the CFPB.) Enforced by
CFPB, FTC, DOJ; private right of action (actual + punitive up to $10,000
individual / capped class damages + fees).

- **Prohibited bases (§1691(a)):** **race, color, religion, national origin,
  sex** (incl. sexual orientation and gender identity per CFPB interpretation),
  **marital status, age** (if the applicant can contract); because income comes
  from a **public-assistance program**; or because the applicant in good faith
  **exercised any right under the Consumer Credit Protection Act** (§1691(a)(3) —
  the federal umbrella statute that includes ECOA/TILA/FDCPA, *not* the
  California privacy law of the same initials).
- **Adverse-action notice (§1691(d)).** A creditor must notify the applicant of
  action taken and give **specific reasons** for adverse action (or notice of the
  right to request reasons). Generally **within 30 days**.
- **Disparate treatment vs. disparate impact.** ECOA reaches both intentional
  discrimination **and** facially neutral policies with a disproportionate
  adverse effect on a protected group that aren't justified by business necessity
  (with a less-discriminatory-alternative analysis). (See frontier note below on
  AI underwriting.)

---

## TILA — Truth in Lending Act (15 USC §1601 et seq.) / Reg Z (12 CFR 1026)

**Cost-of-credit transparency.** Implemented by **Regulation Z (12 CFR Part
1026)**. Enforced by CFPB/FTC; private right of action.

- **Disclosures.** Creditors must disclose the **finance charge** and the
  **Annual Percentage Rate (APR)**. The APR is the standardized "total cost of
  credit" figure (not the same as the nominal interest rate), plus payment terms,
  total of payments, etc.
- **Right of rescission (§1635).** For a credit transaction secured by the
  consumer's **principal dwelling** (e.g., a refinance or home-equity loan, not
  a purchase-money mortgage), the consumer may **rescind until midnight of the
  third business day** after consummation or delivery of the required disclosures
  and rescission notice, whichever is later. If material disclosures are never
  given, the rescission window can extend up to **3 years**.
- TILA also houses the **TILA-RESPA Integrated Disclosure (TRID)** mortgage
  forms, credit-card protections from the **CARD Act**, and ability-to-repay /
  qualified-mortgage rules.

---

## FCBA — Fair Credit Billing Act (15 USC §1666 et seq.)

A 1974 **amendment to TILA** (so it's enforced through Reg Z) governing
**billing-error disputes on open-end / credit-card accounts.**

- The consumer must send **written notice within 60 days** after the first
  statement containing the error was sent.
- The creditor must **acknowledge in writing within 30 days** and **resolve
  within two complete billing cycles (max 90 days)**, either correcting the
  account or explaining why it's correct.
- During the dispute the creditor **may not** try to collect the disputed amount,
  report it as delinquent to a CRA, or close the account for nonpayment of it.
- "Billing error" includes unauthorized charges, wrong amounts/dates, goods/
  services not delivered or not as agreed, and accounting/posting failures.
- Distinct from the **EFTA/Reg E** error process, which covers *debit*/electronic
  transfers (different deadlines and liability rules; see below).

---

## CROA — Credit Repair Organizations Act (15 USC §1679 et seq.)

Regulates for-profit **credit-repair companies.** Enforced by CFPB/FTC; private
right of action; contracts violating CROA are voidable.

- **Prohibited conduct (§1679b(a)):** may not make or advise consumers to make
  **untrue or misleading statements** to a CRA or creditor, and may not engage in
  fraud or deception.
- **Advance-fee ban (§1679b(b)) — the core rule:** *"No credit repair
  organization may charge or receive any money or other valuable consideration
  for the performance of any service which the credit repair organization has
  agreed to perform for any consumer before such service is fully performed."*
  → They **cannot charge you up front.**
- **Required disclosures & contract terms (§§1679c–1679e):** a written contract,
  a "Consumer Credit File Rights Under State and Federal Law" disclosure, and a
  **3-day right to cancel.**
- Reality check: a credit-repair firm can do nothing a consumer can't do for free
  themselves; accurate negative information that's within the reporting period
  cannot be lawfully "erased." Up-front-fee demands are a classic illegal-operator
  tell.

---

## EFTA — Electronic Fund Transfer Act (15 USC §1693 et seq.) / Reg E (12 CFR 1005) — brief

Consumer rights for **electronic fund transfers** (debit cards, ATM, ACH, P2P,
direct deposit; also remittance transfers since 2010). Implemented by
**Regulation E (12 CFR Part 1005).**

- **Unauthorized-transfer liability (§1693g):** tiered and time-sensitive —
  generally capped at **$50** if the consumer reports a lost/stolen access device
  within **2 business days** of learning of the loss; up to **$500** if they fail
  to report within that 2-day window; and potentially **unlimited** for
  unauthorized transfers appearing on a statement that go unreported for more
  than **60 days**. (Reg E §1005.6(b) details; card-network "zero liability"
  policies often add more protection.)
- **Error-resolution (§1693f):** the consumer reports an error (generally within
  **60 days** of the statement); the institution must investigate, often
  providing **provisional credit** within 10 business days while it does.
- Key contrast: **credit-card** disputes run under **FCBA** (above); **debit/
  electronic** disputes run under **EFTA/Reg E** — different deadlines and
  liability exposure, which is why payment method matters.

---

## FRONTIER / VOLATILE AREAS — verify current status before relying on these

> These are the parts most likely to have **changed** since this skill was
> written. Treat every statement here as "as of 2026, subject to change —
> confirm."

### 1. CFPB 2024 medical-debt rule (amends Regulation V) — CONTESTED / likely vacated

- **What it was:** a CFPB final rule (published in the Federal Register Jan. 14,
  2025) that would **remove most medical debt from consumer credit reports** and
  bar creditors from considering medical debt in many credit decisions, by
  deleting a regulatory exception in **Regulation V (FCRA)**.
- **Status — DO NOT assert it is in force.** It was reported **vacated on July
  11, 2025** by the **U.S. District Court for the Eastern District of Texas** (in
  *Cornerstone Credit Union League v. CFPB*, on a joint request of the Bureau and
  plaintiffs), on the ground that it exceeded the CFPB's FCRA authority. As of
  2026 that vacatur is the last reported posture. **Independently verify the
  current status** (any appeal, re-proposal, or successor rule) before stating
  the law; do not represent the rule as effective.
- **Practical nuance to confirm separately:** independent of the CFPB rule, the
  three nationwide bureaus voluntarily stopped reporting **paid** medical
  collections and those **under $500 / under 12 months old** (a 2022–2023
  industry change), and **several states** have enacted their own medical-debt
  credit-reporting bans. Those are distinct from the vacated federal rule and have
  their own current status; check each.

### 2. Fair lending, disparate impact & AI / algorithmic underwriting (ECOA / Reg B)

- ECOA's disparate-impact reach to **automated / "black-box" underwriting** is an
  active enforcement and policy area. Per **CFPB Circular 2023-03**, a creditor
  using **AI/ML or complex algorithms** is **not excused** from ECOA's
  adverse-action obligations: it must give the applicant **specific and accurate
  reasons** for a denial; it cannot fall back on generic checklist reasons that
  don't reflect the model's actual basis, and "the model is too complex to
  explain" is not a defense.
- The CFPB has also signaled that lenders should **test models for prohibited-
  basis disparities** and search for **less-discriminatory alternatives.**
- **Volatility:** the legal status of **disparate-impact liability** generally,
  and agency guidance on AI underwriting specifically, have been subject to
  executive-branch policy shifts and litigation. Confirm whether the circulars,
  guidance, and disparate-impact doctrine cited here remain operative before
  relying on them.

---

## References / verify current law

Primary sources (statute = Cornell LII; regulation = eCFR / CFPB; consumer
guidance = FTC/CFPB). **Re-check these; text and rule status change.**

**Statutes (U.S. Code, Cornell LII):**
- FCRA — 15 USC §1681 et seq.: https://www.law.cornell.edu/uscode/text/15/1681
  (disputes §1681i: https://www.law.cornell.edu/uscode/text/15/1681i ;
  furnishers §1681s-2: https://www.law.cornell.edu/uscode/text/15/1681s-2 )
- FDCPA — 15 USC §1692 et seq.: https://www.law.cornell.edu/uscode/text/15/1692
- ECOA — 15 USC §1691 et seq.: https://www.law.cornell.edu/uscode/text/15/1691
- TILA — 15 USC §1601 et seq.: https://www.law.cornell.edu/uscode/text/15/1601
  (rescission §1635: https://www.law.cornell.edu/uscode/text/15/1635 )
- FCBA — 15 USC §1666 et seq.: https://www.law.cornell.edu/uscode/text/15/1666
- CROA — 15 USC §1679 et seq.: https://www.law.cornell.edu/uscode/text/15/1679b
- EFTA — 15 USC §1693 et seq.: https://www.law.cornell.edu/uscode/text/15/1693

**Regulations (eCFR / CFPB):**
- Reg B (ECOA) — 12 CFR 1002: https://www.consumerfinance.gov/rules-policy/regulations/1002/
- Reg E (EFTA) — 12 CFR 1005: https://www.consumerfinance.gov/rules-policy/regulations/1005/
- Reg F (FDCPA) — 12 CFR 1006: https://www.ecfr.gov/current/title-12/chapter-X/part-1006
  (validation §1006.34: https://www.consumerfinance.gov/rules-policy/regulations/1006/34/ ;
  call frequency §1006.14: https://www.consumerfinance.gov/rules-policy/regulations/1006/14/ )
- Reg V (FCRA) — 12 CFR 1022: https://www.consumerfinance.gov/rules-policy/regulations/1022/
- Reg Z (TILA) — 12 CFR 1026: https://www.consumerfinance.gov/rules-policy/regulations/1026/

**Frontier / contested:**
- CFPB medical-debt final rule (Reg V), Federal Register (2025-01-14):
  https://www.federalregister.gov/documents/2025/01/14/2024-30824/prohibition-on-creditors-and-consumer-reporting-agencies-concerning-medical-information-regulation-v
- CFPB Circular 2023-03 (adverse-action notices, incl. AI/ML models):
  https://www.consumerfinance.gov/compliance/circulars/circular-2023-03-adverse-action-notification-requirements-and-the-proper-use-of-the-cfpbs-sample-forms-provided-in-regulation-b/
- CFPB on AI credit denials:
  https://www.consumerfinance.gov/about-us/newsroom/cfpb-issues-guidance-on-credit-denials-by-lenders-using-artificial-intelligence/

**Agencies / help:**
- CFPB: https://www.consumerfinance.gov  (complaints: https://www.consumerfinance.gov/complaint/ )
- FTC consumer advice: https://consumer.ftc.gov  (report fraud: https://reportfraud.ftc.gov )

---

## Cross-references (within the `consumer-credit-and-debt` hub)

- **How credit reports/scores are built** (bureaus, FICO/VantageScore, inquiries,
  aging timelines) → `credit-reports-and-scores`
- **Tactics to raise/rebuild a score** and DIY dispute *steps* / goodwill letters
  / secured cards → `improving-and-rebuilding-credit`
- **Settling or negotiating charge-offs & collections** (pay-for-delete,
  settlement %, 1099-C) → `charge-offs-collections-and-debt-resolution`
- **Practical interaction with debt collectors** (validation letters, stopping
  calls, harassment responses) → `debt-collectors-and-fdcpa-rights`
- **North Carolina-specific statutes** (state collection / usury / debt law) →
  `north-carolina-credit-and-debt-law`

This skill stays on the **federal statutory framework**; route how-to and
state-specific questions to the spokes above.
