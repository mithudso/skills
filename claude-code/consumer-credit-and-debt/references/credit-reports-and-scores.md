<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `credit-reports-and-scores` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: credit-reports-and-scores
description: >-
  US consumer credit reporting & scoring: bureaus/furnishers, what's on a
  report, FICO vs VantageScore models, inquiries, freezes/locks, and
  negative-item aging. Spoke of consumer-credit-and-debt. Educational, not
  advice.
  TRIGGER: how a report or score is built; what the bureaus collect; FICO factor
  weights/ranges; FICO 10T vs VantageScore 4.0; hard vs soft inquiries; free
  reports (AnnualCreditReport.com), freeze vs lock vs fraud alert;
  how long a late, charge-off, bankruptcy, or inquiry stays on a report;
  mortgage score transition (FHFA); BNPL/cashflow scoring.
  SKIP: raising/rebuilding a score, dispute tactics, secured cards ->
  improving-and-rebuilding-credit; settling a charge-off/collection ->
  charge-offs-collections-and-debt-resolution; statute text/disputes/suing ->
  us-consumer-credit-and-debt-law; stopping collectors, validation ->
  debt-collectors-and-fdcpa-rights; how medical debt reports/ages ->
  medical-debt-and-billing; credit-based insurance scores ->
  personal-insurance.
metadata:
  changelog:
    - "2026-06-16 sko --meta --no-sync — structural pass vs new consumer-finance family; 0 High 0 Med (meta-validate clean). Added 2 SKIP edges (-> medical-debt-and-billing, -> personal-insurance); trimmed description to 993 chars within 1000 cap; no content passes (already CLEAN prior)."
---

# Credit Reports & Scores (US)

**General educational information only — not financial or legal advice.** Credit
scoring models, bureau policies, and federal rules change; numbers below are
**as of 2026** and approximate. Verify current details with the bureau, the
score provider (FICO/VantageScore), or the CFPB before relying on them. This is
a **spoke of the `consumer-credit-and-debt` hub**; see the SKIP list in the
description for sibling skills that own dispute tactics, debt resolution, FCRA
legal mechanics, and FDCPA rights.

---

## 1. The three nationwide bureaus & how data gets there

The three **nationwide consumer reporting agencies (NCRAs / CRAs)** are
**Equifax, Experian, and TransUnion**. They are private companies, not
government agencies, and they do **not** share a single file; each maintains
its own database, so your three reports (and scores) often differ.

- **Furnishers** are the lenders, card issuers, collectors, and other
  creditors that voluntarily **report** account data to the bureaus, typically
  monthly. Furnishing is voluntary; a given creditor may report to all three,
  some, or none. This is why an account can appear on one report but not another.
- Bureaus aggregate furnished data plus public records into a **credit file**;
  a **credit report** is a snapshot of that file. **Credit scores** are
  separate analytic products (FICO, VantageScore) computed *from* the report
  data, sold by the bureaus and by FICO.
- The data ecosystem is largely self-correcting through **disputes** under the
  Fair Credit Reporting Act (FCRA). The *mechanics* of disputes live in the
  `us-consumer-credit-and-debt-law` sibling skill, not here.

---

## 2. What's on a credit report

Four broad sections:

1. **Personal / identifying information** — name(s), current and former
   addresses, date of birth, Social Security number (often partially masked),
   employers. *Not used in scoring*; used to match the file to you.
2. **Tradelines (accounts)** — one line per credit account: creditor, account
   type (revolving / installment / open / mortgage), date opened, credit limit
   or original loan amount, current balance, monthly payment, and a
   month-by-month **payment history** (on-time / 30/60/90/120+ days late /
   charge-off). The single most influential section.
3. **Inquiries** — records of who pulled the file (see §5). Hard inquiries are
   visible to others; soft inquiries are visible only to you.
4. **Public records & collections** — now narrow. After the **National Consumer
   Assistance Plan (2017–2018)**, the bureaus **stopped including civil
   judgments and tax liens** that lack full identifying data; in practice
   **bankruptcies** are the main remaining public record. **Collection
   accounts** (debts sold/assigned to a collector) appear as their own
   tradelines.

---

## 3. Credit scores — FICO factors & weights

A **FICO Score** (the dominant model, used in the large majority of lending
decisions) is built from five categories. Approximate weights for the general
population (**individual weighting varies**):

| Factor | ~Weight | What it captures |
|---|---|---|
| **Payment history** | **~35%** | Late/on-time record, severity, recency, public records/collections |
| **Amounts owed / utilization** | **~30%** | Balances vs. limits; **revolving utilization** matters most; installment loans weigh less |
| **Length of credit history** | **~15%** | Age of oldest account, average account age, age of newest |
| **New credit** | **~10%** | Recent account openings and hard inquiries |
| **Credit mix** | **~10%** | Variety: revolving (cards) + installment (auto/mortgage/student) |

**Utilization** = balance ÷ credit limit on revolving accounts. Lower is better;
it is recalculated each time balances are reported and has **no memory** (paying
down a card can move the score the next cycle, unlike payment-history damage).

### Score ranges & tiers (base FICO, 300–850)

| Tier | Range |
|---|---|
| Poor | 300–579 |
| Fair | 580–669 |
| Good | 670–739 |
| Very Good | 740–799 |
| Exceptional | 800–850 |

Note: **industry-specific FICO scores** (e.g., FICO Auto Score, Bankcard Score)
use a **250–900** range but keep the same tier groupings in the middle. There are
also **many FICO versions in use at once** (FICO 8 is most common in general
lending; FICO 9, 10, and 10T also exist), so the "FICO score" a lender sees can
differ from a free score you check.

### VantageScore

**VantageScore** is the joint model from the three bureaus and the main
competitor to FICO. Also a **300–850** scale (current versions). Instead of
fixed percentages it ranks factors in influence bands:

| Influence | Factor |
|---|---|
| Most influential | Payment history |
| Highly influential | Age and type of credit; credit utilization |
| Moderately influential | Total balances / debt |
| Less influential | Recent credit behavior; available credit |

Other practical differences vs. FICO:

- Can score **"thin-file" / unscoreable** consumers FICO can't, partly by using
  a longer look-back and (in 4.0) **alternative data**.
- Free scores from many consumer apps and card issuers are frequently
  VantageScore, which is why an app score may not match a lender's FICO.

### Why your three scores differ

A common surprise: you have many scores, not one. They differ because of three
stacked variables, so do not expect a single number:

1. **Different bureau data** — furnishing is voluntary, so each bureau may hold
   a different set of accounts (see §1).
2. **Different model and version** — FICO 8 vs FICO 9/10/10T vs VantageScore,
   each weighting the same data differently.
3. **Different pull date** — balances and inquiries change week to week.

The free score you check is often a different model or version than the FICO a
given lender pulls, so a gap between "my app score" and "what the lender saw"
is normal.

---

## 4. Frontier — emerging, as of 2026

> These developments are real and in progress, but their **dates, sequencing,
> and adoption levels keep moving**. Treat specific timing below as fluid as of
> 2026 and verify with the Federal Housing Finance Agency (FHFA), the bureaus,
> and the score providers.

- **FICO 10 / 10T & VantageScore 4.0** are the newest mainstream models.
  **FICO 10T and VantageScore 4.0 both use "trended data"**: roughly **24+
  months** of balance/limit/payment trajectory, not just the latest snapshot,
  so whether balances are rising or falling matters, not only today's
  utilization. **VantageScore 4.0** also incorporates **alternative data** (rent,
  utility, telecom) where available. Adoption is gradual; FICO 8 still dominates
  most non-mortgage lending.
- **Mortgage credit-score transition (FHFA / Fannie Mae & Freddie Mac).** Fannie
  and Freddie (the government-sponsored enterprises, or GSEs, that FHFA
  oversees) historically required the older **"Classic FICO"** plus a
  **tri-merge** (all three bureaus). FHFA has **validated both FICO 10T and
  VantageScore 4.0** and is moving the GSE market toward (a) **accepting
  VantageScore 4.0** alongside or instead of Classic FICO, and (b) a **bi-merge**
  (two-bureau) requirement instead of tri-merge, introduced through an
  **interim phase** letting approved lenders deliver loans scored by *either*
  Classic FICO *or* VantageScore 4.0. Confirm the current sequencing on fhfa.gov.
- **Cashflow / alternative-data underwriting.** Lenders increasingly supplement
  bureau scores with **bank-transaction (cashflow) data** and consumer-permissioned
  income/utility/rent data to score thin-file applicants. This generally rides
  *alongside* the bureau file rather than replacing it.
- **Buy Now, Pay Later (BNPL).** The bureaus are building BNPL into their
  systems, but treatment is **inconsistent and still evolving**: as of 2026 much
  furnished BNPL data is **visible on the report but not yet factored into legacy
  FICO/VantageScore models**, and large providers (e.g., Affirm) have begun
  reporting pay-over-time loans to Experian and TransUnion. Expect this to feed
  newer models over time. Verify per-bureau and per-lender.

---

## 5. Hard vs. soft inquiries

| | **Hard inquiry (hard pull)** | **Soft inquiry (soft pull)** |
|---|---|---|
| When | You **apply** for credit (card, mortgage, auto, personal loan) | Pre-approval offers, account reviews, **checking your own** report/score, many employer/insurance checks |
| Affects score? | **Yes**, usually under ~5 points for most people (more with few accounts or a short history) | **No** |
| Visible to others? | Yes | **Only you** |
| Stays on report | **~2 years** (but FICO only *considers* the last **12 months**) | Varies; not scored |

- **Rate-shopping protection:** in FICO, multiple hard inquiries of the **same
  type** (mortgage, auto, or student loan) within a short window count as
  **one**, so shop rates in a tight burst. The window is **45 days** in newer
  FICO versions (**14 days** in older ones). VantageScore uses a rolling ~14-day
  window that dedupes across loan types, not just the same type.
- **Credit cards and personal loans are generally not bundled** this way: each
  application is its own inquiry, so the burst trick does not protect
  card shopping.
- **Checking your own credit is always a soft pull** and never lowers your
  score, including AnnualCreditReport.com and most free-score apps.

---

## 6. Getting reports, monitoring, and protection tools

### Obtaining your reports

- **AnnualCreditReport.com is the only federally authorized source** for your
  free reports. As of 2026 all three bureaus provide **free reports weekly**
  (a COVID-era expansion that was made **permanent**), on top of the statutory
  free annual report. Beware sound-alike sites that charge or phish.
- A free *report* does not always include a free *score*; scores can be bought
  from the bureaus/FICO or obtained free via many card issuers and apps (often
  VantageScore or a specific FICO version, so check which).

### Credit monitoring

Paid or free services that alert you to file changes (new inquiries, new
accounts, balance/score moves). Useful for early fraud detection, but monitoring
is **detective, not preventive**: it tells you after something happens.

### Security freeze vs. credit lock vs. fraud alert

| Tool | What it does | Cost | Duration | Set up |
|---|---|---|---|---|
| **Security freeze** | Blocks new creditors from pulling your file, so new accounts can't be opened in your name. Your existing creditors and you still have access. | **Free by federal law** | Until **you** lift it (lift/relift is free; must be done **per bureau**) | Contact **each of the three bureaus** |
| **Credit lock** | Similar blocking effect, controlled via an app | Often **bundled into a paid product** | While locked | Per bureau/product; **contractual, not statutory** |
| **Fraud alert** | Tells creditors to **verify your identity** before opening credit (doesn't block access) | **Free** | **Initial: 1 year** (renewable); **extended: 7 years** with an identity-theft report; **active-duty** alert: 1 year | Tell **one** bureau — it must notify the other two |

Practical guidance the FTC/CFPB give: a **freeze is the strongest free
protection** and is generally preferred over a paid lock for that reason; a
freeze must be placed/lifted at **all three** bureaus separately, whereas a
**fraud alert** propagates from one bureau to the others.

---

## 7. Negative-item aging timelines

How long adverse items remain (FCRA-driven; **general rule = 7 years**). Always
**verify the controlling date** for a specific item:

| Item | Stays on report |
|---|---|
| Late payments / most negative info | **~7 years** |
| **Charge-offs & collections** | **~7 years from the original delinquency date** that led to it, in practice **~7½ years** (the FCRA's ~180-day allowance before the clock starts). Crucially, the clock runs from the **original delinquency**, *not* from when a debt is sold or a collector re-reports it. |
| **Chapter 7 bankruptcy** | **10 years** from the filing date |
| **Chapter 13 bankruptcy** | typically **7 years** from filing (vs. 10 for Chapter 7) |
| **Hard inquiries** | **~2 years** (scored only ~12 months) |
| Civil judgments / tax liens | Largely **no longer reported** by the bureaus post-2017–18 (see §2) |
| Closed accounts **in good standing** / positive history | Can remain up to **~10 years**, which *helps* your average account age |

Note: removing an item at the 7-year mark is **automatic**; you don't have to do
anything. "Re-aging" a debt to restart the clock is **prohibited**.

---

## 8. How-to quick steps

1. **Pull all three reports** free at AnnualCreditReport.com (stagger across the
   year, or pull all three at once now that weekly is free). Read each section in §2.
2. **Check accuracy** of every tradeline, the personal-info section, and the
   inquiry list. (For the formal *dispute* process and your FCRA rights, go to
   the `us-consumer-credit-and-debt-law` sibling skill.)
3. **Know your scores' provenance** — which model/version and which bureau the
   number came from, so you compare like-for-like.
4. **Protect the file:** place a **free security freeze** at all three bureaus
   if you're not actively applying for credit; add a **fraud alert** if you
   suspect identity theft.
5. **Shop loans in a tight window** (≤14 days to be safe) so rate-shopping
   inquiries dedupe.
6. **Re-verify volatile facts** (model adoption, FHFA timing, BNPL treatment,
   freeze rules) against the primary sources below before acting.

---

## References / verify current

Primary, authoritative sources (re-check; these update):

- **CFPB, Credit reports & scores hub:** https://www.consumerfinance.gov/consumer-tools/credit-reports-and-scores/
- **CFPB, How long does information stay on my credit report?:** https://www.consumerfinance.gov/ask-cfpb/how-long-does-information-stay-on-my-credit-report-en-323/
- **CFPB, How long does a bankruptcy appear?:** https://www.consumerfinance.gov/ask-cfpb/how-long-does-a-bankruptcy-appear-on-credit-reports-en-325/
- **CFPB, What is a credit inquiry?:** https://www.consumerfinance.gov/ask-cfpb/what-is-a-credit-inquiry-en-1317/
- **CFPB, What is a credit freeze / security freeze?:** https://www.consumerfinance.gov/ask-cfpb/what-is-a-credit-freeze-or-security-freeze-on-my-credit-report-en-1341/
- **CFPB, BNPL and credit reporting:** https://www.consumerfinance.gov/about-us/blog/by-now-pay-later-and-credit-reporting/
- **FTC Consumer Advice, Credit freezes & fraud alerts:** https://consumer.ftc.gov/articles/credit-freezes-and-fraud-alerts
- **FTC Consumer Advice, Free credit reports / permanent free weekly:** https://consumer.ftc.gov/articles/free-credit-reports and https://consumer.ftc.gov/consumer-alerts/2023/10/you-now-have-permanent-access-free-weekly-credit-reports
- **AnnualCreditReport.com** (the only authorized source): https://www.annualcreditreport.com/
- **myFICO, What's in your FICO score (factors & weights):** https://www.myfico.com/credit-education/whats-in-your-credit-score
- **myFICO, What is a credit score (ranges & tiers):** https://www.myfico.com/credit-education/credit-scores
- **myFICO, Managing credit inquiries (rate-shopping window):** https://www.myfico.com/credit-education/credit-reports/manage-credit-inquiries
- **myFICO, FICO score versions (8/9/10/10T):** https://www.myfico.com/credit-education/credit-scores/fico-score-versions
- **VantageScore, knowledge center (4.0, model comparisons):** https://vantagescore.com/resources/knowledge-center/
- **FHFA, Credit Scores policy hub & news (FICO 10T / VantageScore 4.0 validation, bi-merge):** https://www.fhfa.gov/policy/credit-scores and https://www.fhfa.gov/news/news-release/fhfa-announces-validation-of-fico-10t-and-vantagescore-4.0-for-use-by-fannie-mae-and-freddie-mac
- **Experian / TransUnion BNPL pages** (per-bureau treatment, evolving): https://www.experian.com/consumer-education-content/buy-now-pay-later-faq and https://www.transunion.com/buy-now-pay-later
