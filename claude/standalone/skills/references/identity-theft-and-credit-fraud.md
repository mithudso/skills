<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `identity-theft-and-credit-fraud` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: identity-theft-and-credit-fraud
description: >-
  Identity theft and credit/financial fraud — prevention, detection, and the
  step-by-step FTC IdentityTheft.gov recovery plan, plus FCRA recovery rights
  and tax-related ID theft. Spoke of consumer-credit-and-debt. Educational
  info, NOT legal advice (as of 2026).
  TRIGGER: identity theft; new-account fraud; account takeover; synthetic
  identity; freeze my credit after fraud; fraud alert;
  IdentityTheft.gov / FTC Identity Theft Report; FCRA 605B / 609(e); tax ID
  theft / IRS IP PIN / Form 14039; data breach; child, medical, or criminal
  identity theft; collectors on a fraudulent debt.
  SKIP: how freezes/scores GENERALLY work → credit-reports-and-scores;
  FCRA/FDCPA statute TEXT → us-consumer-credit-and-debt-law; collector tactics on
  a real debt → debt-collectors-and-fdcpa-rights; settling/charge-off, 1099-C →
  charge-offs-collections-and-debt-resolution; investment-scam red flags
  (Ponzi / guaranteed returns) → investing-and-retirement; NC ID Theft Act
  detail → north-carolina-credit-and-debt-law.
metadata:
  version: 1.1.0
  updated: 2026-06-16
  changelog:
    - "2026-06-16 sko v1.0.0->v1.1.0 (--meta --no-sync) — 1 Medium structural: Pass I collision vs newly-expanded consumer-finance sibling set; +SKIP edge investment-scam red flags -> investing-and-retirement (reciprocal; investing owns scam-spotting, this skill owns post-scam identity recovery) in desc + body; desc trimmed to 998 (<=1000). No content passes."
---

# Identity Theft & Credit / Financial Fraud

General **educational information, NOT legal, financial, or tax advice.** Identity
theft is stressful and time-sensitive; this is a map, not a substitute for the
official tools or a professional. **Processes, links, phone numbers, and dollar/time
limits change; verify against the primary sources below. Current as of 2026.**

Spoke of the **consumer-credit-and-debt** hub. This skill is the *fraud / "it wasn't
me"* lane: prevention, detection, the full recovery plan, **data-breach response (§7)**,
and a brief **North Carolina overlay (§8)**. For how credit reporting, freezes,
monitoring, and scores work in the normal (non-fraud) case, see **credit-reports-and-scores**.
For the literal text of the FCRA/FDCPA, see **us-consumer-credit-and-debt-law**.

**The one thing to remember:** the federal government's free, official one-stop site
is **IdentityTheft.gov**. It generates your recovery plan, your FTC Identity Theft
Report, and pre-filled dispute letters. Start there.

---

## 1. Types of identity theft (and the warning signs)

Identity theft = someone uses your personal information (name, SSN, card/account
numbers, medical or insurance IDs) without permission, usually for financial gain.

| Type | What it is |
|---|---|
| **New-account fraud** | Thief opens *new* credit cards, loans, utilities, or accounts in your name. |
| **Account takeover (ATO)** | Thief gains control of your *existing* account (changes login, address, adds an authorized user) and drains or charges it. |
| **Synthetic identity** | Thief combines *real* data (often a real SSN, frequently a child's or someone with a thin file) with a *made-up* name/DOB to create a brand-new "person" and build credit. Hard to detect because it may not map cleanly to one victim's file. The fastest-growing form of identity fraud. |
| **Tax-related ID theft** | Thief uses your SSN to file a fraudulent tax return and grab your refund, or to get a job (employment-related). See §6. |
| **Medical identity theft** | Thief uses your identity to get care, drugs, or to bill insurance. Dangerous: corrupts your *medical record* (wrong blood type, allergies) and is harder to clean up than a credit file. |
| **Criminal identity theft** | Someone gives *your* name to police at arrest/citation. Surfaces as a bogus warrant, failed background check, or denied job/clearance. |
| **Child identity theft** | A minor's SSN is exploited (children have clean, unmonitored files). Often discovered years later. |

**Warning signs (investigate if you notice):**
- Bills for things you didn't buy; withdrawals you didn't make.
- Expected bills or statements *stop arriving* (address may have been changed).
- Accounts, inquiries, or addresses on your credit report you don't recognize.
- A debt collector calls about a debt that isn't yours.
- A health plan rejects a claim because you've "reached your limit"; a medical bill
  or Explanation of Benefits for care you never got.
- The **IRS** rejects your e-filed return as a duplicate, or sends a notice about a
  return/wages/employer you don't recognize.
- **For a child:** a bill/collection notice in the child's name, a denied government
  benefit, an IRS notice, or a pre-approved credit offer addressed to your minor.

---

## 2. Prevention checklist

**Protect the keys (especially the SSN):**
- Don't carry your Social Security card; store it and other sensitive docs in a safe
  place. Give your SSN only when truly required; ask why, and how it's protected.
- Never email your SSN or full account/card numbers. Don't print SSN/DL on checks.

**Lock down accounts:**
- Strong, unique passwords via a **password manager**; never reuse the password you
  used at a breached company.
- Turn on **multi-factor authentication (MFA)** everywhere it's offered (prefer an
  authenticator app or hardware security key over SMS).
- Read card/bank statements promptly; note when bills are due and chase missing ones.

**Use the free legal tools (federal law):**
- **Security freeze** on all three bureaus (see §4): the strongest free protection
  against *new-account* fraud.
- Pull and review your **free reports** at AnnualCreditReport.com (the official site).

**Reduce exposure:** shred old financial/medical paperwork and labels; be alert to
phishing (don't act on unsolicited "verify your account" emails/texts/calls); secure
your mailbox and devices.

---

## 3. Detection / monitoring

- **AnnualCreditReport.com**: the only federally authorized free source for your
  Equifax, Experian, and TransUnion reports. Review for unknown accounts, inquiries,
  or addresses. (How scores/reports are built → **credit-reports-and-scores**.)
- **Statement review**: your own card and bank statements catch fraud the credit
  report won't (small "test" charges, micro-deposits).
- **Free credit monitoring** if offered after a breach: accept it; it's an early-warning
  signal, not a fix.
- **Tax:** opt into the IRS **IP PIN** (§6), the key control against tax-refund fraud.

---

## 4. Freezes vs. fraud alerts (both free under federal law)

Free nationwide under federal law since 2018. These are different tools; most ID-theft
victims should do **both**. (Conceptual freeze-vs-lock detail → **credit-reports-and-scores**.)

| | **Security freeze** | **Fraud alert** |
|---|---|---|
| What it does | Blocks new creditors from pulling your report, so new accounts can't be opened. | Tells creditors to verify your identity before extending new credit. |
| Strength | Strongest barrier to *new-account* fraud. | A speed bump; you can still get credit normally. |
| How to place | Contact **each** of the 3 bureaus individually (one does **not** notify the others). | Contact **one** bureau; it must notify the other two. |
| Timing | Bureau must place online/phone request within **1 business day**; must lift within **1 hour** online/phone. | — |
| Duration | Until you lift it. | **Initial: 1 year** (renewable). **Extended: 7 years.** **Active-duty: 1 year.** |
| Effect on score | **None.** | None. |

- **Initial (1-yr) alert:** anyone who suspects fraud can place it; entitles you to a
  free credit report. Placing an extended alert entitles you to **2 free reports** in
  the following 12 months (on top of your AnnualCreditReport.com pulls).
- **Extended (7-yr) alert:** only for confirmed victims who have an **FTC Identity Theft
  Report or police report**; also removes you from prescreened-offer marketing lists for
  5 years.
- **Active-duty alert:** for deployed servicemembers (1 year); removes you from
  prescreened-offer marketing lists for 2 years.
- **Protected Consumer / minor freeze:** a parent/guardian can freeze the file of a child
  or an incapacitated adult (strong defense against child identity theft).

**Recommended order for a victim:** place the **fraud alert first** (it's instant at one
bureau and that bureau notifies the other two), then place the **freeze** at all three (the
stronger lock against new-account fraud). Do both; the alert buys time while you set freezes.

Bureau contacts (verify, these change): Equifax 800-685-1111; Experian 888-397-3742;
TransUnion 888-909-8872.

---

## 5. Recovery — the FTC IdentityTheft.gov plan

This is the core workflow. **IdentityTheft.gov** walks you through it, tracks progress,
and pre-fills the letters.

**Documents to have ready** (the recovery flow and the §605B / §609(e) letters all need
these): a government photo ID, proof of address (utility bill or lease), your **SSN**, a
list of the fraudulent accounts/charges, and (once you have it) your **FTC Identity Theft
Report number**. IdentityTheft.gov pre-fills the letters from what you enter.

### Step 1 — Act immediately
- **Call the fraud department** of each company where fraud happened. Ask them to
  **close or freeze** the fraudulent accounts. Get confirmation in writing.
- **Change** logins, passwords, and PINs on affected (and similar) accounts.
- **Place a fraud alert** with one bureau, then a **freeze** at all three (free; §4);
  do both (order and rationale in §4).
- **Get your credit reports** (AnnualCreditReport.com) and mark every fraudulent item.

### Step 2 — Report to the FTC and create your FTC Identity Theft Report
- Report at **IdentityTheft.gov**. You receive a **personal recovery plan** and an
  **FTC Identity Theft Report**.
- The **FTC Identity Theft Report** is an official report that *proves to businesses and
  bureaus that you're a victim*. For most companies and all three bureaus, it is
  sufficient documentation on its own to fix fraud-related problems.

### Step 3 — File a police report (optional, situational)
- Often optional, but **do it** if you know the thief, a creditor/collector or the DMV
  demands it, or for criminal/medical/synthetic cases. Bring your FTC Identity Theft
  Report, a government ID, proof of address, and proof of the theft.
- An **FTC report + a police report** together is the classic **"Identity Theft Report"**
  that unlocks the strongest FCRA rights below.

### Step 4 — Use your FCRA recovery rights
These federal rights are *triggered by* having an Identity Theft Report:

- **§605B block of fraudulent information.** Send the bureau your Identity Theft Report,
  proof of identity, and a statement identifying which items are fraudulent. The credit
  bureau must **block** that information from your report **within 4 business days** of
  receiving the request and must notify the furnisher. Result: the fraudulent tradeline
  stops showing and stops being reported.
- **Disputes / reinvestigation (§611).** Separately dispute fraudulent items with each
  bureau; the furnisher must investigate. (Dispute *mechanics/letters* and statute text
  → **us-consumer-credit-and-debt-law**.)
- **§609(e) right to business records.** You can demand, **in writing**, that any
  business that opened a fraudulent account or transacted with the thief give you the
  **application and transaction records**, **free of charge, within 30 days.** Powerful
  for proving the fraud and for a police investigation.
- **Free freezes and fraud alerts**, including the **7-year extended alert** (§4).

### Step 5 — Handle creditors and collectors on fraudulent debts
- Tell the creditor/collector, in writing, that the debt resulted from identity theft and
  **you don't owe it**; enclose your **Identity Theft Report**.
- Ask the collector to **stop collection** and **stop reporting** the debt to the bureaus
  (FDCPA / Regulation F). Use **§609(e)** to demand the underlying records.
- **Scope note:** that's the *fraud* angle. For a collector's *conduct/tactics*, debt
  validation on a (possibly real) debt, or being sued, see **debt-collectors-and-fdcpa-rights**.
  Time-barred "zombie" debt, settlement, 1099-C → **charge-offs-collections-and-debt-resolution**.

### Step 6 — Special cleanups
- **Medical:** request copies of and corrections to your medical records and your
  insurer's Explanation of Benefits; flag the file with the provider and plan.
- **Criminal:** work with the arresting jurisdiction's court/police to clear the record;
  some states issue an identity-theft "clearance"/passport.
- **Deceased-relative ("ghosting") theft:** send a copy of the death certificate to each
  bureau, the SSA, and the IRS to flag the file, and request a **deceased alert**.
- **New SSN (last resort):** if fraud persists after you've done everything above, the
  **SSA** can issue a new Social Security number in limited cases; it's disruptive (it
  disconnects your existing credit history), so it's a last resort, not a first step.

### Lost or stolen wallet (a common precipitating event)
If a wallet/purse with cards or ID is lost or stolen, act fast: inventory what was in it;
call each **card/bank issuer** to freeze or reissue; report a stolen **driver's license**
to the DMV; if your **SSN card** was in it, treat it as an SSN exposure (place a freeze +
fraud alert, §4); then watch for downstream fraud and run this §5 plan if any appears.

---

## 6. Tax-related identity theft (IRS)

- **IP PIN (Identity Protection PIN):** a 6-digit number that must appear on your federal
  return for it to be accepted; it blocks anyone else from filing as you. **Any taxpayer**
  can opt in via the **Get an IP PIN** tool at IRS.gov (it's reissued yearly; the IRS sends
  it on a **CP01A** notice). The single best defense against refund fraud.
- **Form 14039, Identity Theft Affidavit:** file (online or paper) **only if** you're a
  victim of tax-related ID theft *and* you haven't already gotten an IRS letter about it.
  **Do NOT file 14039** if the IRS sent you an identity-verification letter (e.g., **5071C**);
  just follow that letter's instructions (verify at the IRS site/phone). The IRS opens an
  identity-theft case and works to resolve it and flag your account.
- Still complete the **IdentityTheft.gov** plan and the §5 steps (tax ID theft often
  travels with broader identity theft).

---

## 7. Data-breach response

Getting a breach notice ≠ you've been defrauded, but act:
1. **Change the password** at the breached company and anywhere you reused it; turn on **MFA**.
2. Check **what was exposed**: visit **IdentityTheft.gov/databreach** for tailored steps.
3. **If your SSN was exposed:** pull your free reports, watch for unknown accounts, and
   place a **freeze** (and/or fraud alert).
4. **Accept** any free credit monitoring / ID-theft insurance the company offers.
5. **If fraud appears:** run the §5 recovery plan at IdentityTheft.gov.

North Carolina note: breached entities must notify affected NC residents (see §8).

---

## 8. North Carolina note (brief — cross-ref the NC skill)

The **NC Identity Theft Protection Act** (G.S. **Chapter 75, Article 2A**) is North
Carolina's overlay on the federal baseline. In brief:
- **Free security freezes** for NC consumers, including a **Protected Consumer freeze**
  for a child or an incapacitated adult.
- **Security-breach notification:** a business or government agency that suffers a breach
  of NC residents' personal information must notify the affected individuals **and** the
  **NC Attorney General's Consumer Protection Division** without unreasonable delay (the AG
  view is no later than 30 days). A large breach also requires notifying the nationwide
  credit bureaus; for the resident-count threshold and exact deadlines, see the NC skill.
- **SSN protections:** restrictions on collecting, displaying, and disposing of SSNs.
- You can file a consumer complaint with the **NC DOJ** (ncdoj.gov).

For statutory detail, deadlines, and how NC differs from federal law → **north-carolina-credit-and-debt-law**.

---

## When to route elsewhere

- How freezes/locks, monitoring, inquiries, or scores work in the *non-fraud* case →
  **credit-reports-and-scores**
- Literal **FCRA / FDCPA / Regulation F** statute text → **us-consumer-credit-and-debt-law**
- A collector's **conduct/tactics**, validation on a (real) debt, or being **sued** →
  **debt-collectors-and-fdcpa-rights**
- Settling, charge-offs, "zombie" debt, 1099-C canceled debt →
  **charge-offs-collections-and-debt-resolution**
- Rebuilding a score after the fraud is cleared → **improving-and-rebuilding-credit**
- **Spotting investment-scam red flags** *before* losing money (guaranteed returns,
  Ponzi, crypto/"pig-butchering") → **investing-and-retirement** (this skill owns the
  *post-scam* identity/financial-fraud recovery once your data or accounts are compromised)
- NC Identity Theft Protection Act statutory detail → **north-carolina-credit-and-debt-law**

---

## References / verify current (as of 2026)

Primary, authoritative sources; links and details change, confirm before relying.

- **FTC IdentityTheft.gov** (report + recovery plan + Identity Theft Report + letters):
  https://www.identitytheft.gov  ·  steps: https://www.identitytheft.gov/Steps
- **FTC "What To Know About Identity Theft"** (types, warning signs, prevention):
  https://consumer.ftc.gov/articles/what-know-about-identity-theft
- **FTC Credit Freezes and Fraud Alerts** (free freezes/alerts, durations, bureau contacts):
  https://consumer.ftc.gov/articles/credit-freezes-and-fraud-alerts
- **FTC "Identity Theft: A Recovery Plan"** (booklet, ordered steps):
  https://www.bulkorder.ftc.gov/system/files/publications/501a_idt_a_recovery_plan_508.pdf
- **FTC Data breach: what to do** + **IdentityTheft.gov/databreach**:
  https://consumer.ftc.gov/media/79862  ·  https://www.identitytheft.gov/databreach
- **FTC Protect your child from identity theft:**
  https://consumer.ftc.gov/articles/how-protect-your-child-identity-theft
- **FCRA §605B** (15 U.S.C. §1681c-2), block of fraudulent info (4 business days):
  https://consumer.ftc.gov/articles/pdf-0089-fcra-605b.pdf
- **FCRA §609(e)** (15 U.S.C. §1681g(e)), victim's right to business records (30 days, free):
  https://consumer.ftc.gov/system/files/consumer_ftc_gov/pdf/fcra-609e.pdf
- **CFPB Credit freeze / security freeze:**
  https://www.consumerfinance.gov/ask-cfpb/what-is-a-credit-freeze-or-security-freeze-on-my-credit-report-en-1341/
- **CFPB "I've been a victim of identity theft":**
  https://www.consumerfinance.gov/ask-cfpb/what-do-i-do-if-i-think-i-have-been-a-victim-of-identity-theft-en-31/
- **IRS Identity theft guide for individuals:**
  https://www.irs.gov/identity-theft-central/identity-theft-guide-for-individuals
- **IRS Get an IP PIN:** https://www.irs.gov/identity-theft-fraud-scams/get-an-identity-protection-pin
- **IRS Form 14039 (Identity Theft Affidavit)** + when to file:
  https://www.irs.gov/pub/irs-pdf/f14039.pdf  ·  https://www.irs.gov/newsroom/when-to-file-an-identity-theft-affidavit
- **NC DOJ Identity theft / Free security freeze:**
  https://ncdoj.gov/protecting-consumers/identity-theft/  ·  https://ncdoj.gov/protecting-consumers/protecting-your-identity/free-security-freeze/
- **AnnualCreditReport.com**, official free reports: https://www.annualcreditreport.com
