<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-finance` hub.** Formerly the standalone `medical-debt-and-billing` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-finance`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: medical-debt-and-billing
description: "US medical bills & medical debt: understand a bill, fix errors, resolve it. consumer-finance spoke (sibling hub consumer-credit-and-debt). Educational, not advice; credit-reporting rules volatile as of 2026, verify. TRIGGER: dispute a bill; itemized bill vs EOB; billing errors / CPT codes; surprise/balance billing & the No Surprises Act; good-faith estimate for uninsured/self-pay; hospital charity care / financial assistance / IRS 501(r) / presumptive eligibility; negotiate or settle a bill (prompt-pay discount, payment plan, lump sum); does medical debt hurt my credit; medical bill in collections; CareCredit / medical-card deferred interest; HSA/FSA to pay bills. SKIP: collector conduct / sued / validation → debt-collectors-and-fdcpa-rights; general settlement & 1099-C tax → charge-offs-collections-and-debt-resolution; credit-reporting statute/FCRA → us-consumer-credit-and-debt-law; how health PLANS/coverage/denied-claim appeals work → health-insurance-and-coverage."
version: 1.1.0
updated: 2026-06-16
metadata:
  changelog:
    - "2026-06-16 sko (--meta --no-sync) v1.1.0: structural reconciliation — Pass A/M description hub framing fixed 'hub: consumer-credit-and-debt' -> 'consumer-finance spoke (sibling hub consumer-credit-and-debt)' + TRIGGER trimmed, 999 -> 981 chars (under 1000 cap), all 4 SKIP edges retained; Pass L restored missing closing frontmatter '---' (High, parse break introduced mid-run, fixed same iteration); Pass A body anchor updated consumer-credit-and-debt 'current installed anchor' -> consumer-finance hub + sibling-hub cross-link; added version/updated/metadata frontmatter (were absent); A'/N resolvability closed; Pass I reciprocal peer edges into out-of-set siblings REPORTED not written"
---

# Medical Debt & Billing (US)

> **Framing (read first).** This is **general educational information, NOT legal, financial, tax, or medical advice.** Medical-billing and medical-debt rules vary by state and change often; the **credit-reporting** rules here are especially volatile. Everything is stated **as of 2026**. Verify against the primary sources in *References* before relying on it, and consult a licensed professional (attorney, tax pro, financial counselor, or a hospital financial counselor) for your situation.

**Where this sits.** Spoke of the **`consumer-finance`** hub. Its **sibling hub, `consumer-credit-and-debt`,** owns generic credit/debt/collections and the governing statutes — start there once a bill becomes a collection, settlement, or credit-reporting matter. This skill owns the *medical-specific* layer (bills, the No Surprises Act, hospital charity care/501(r), medical-debt credit reporting); it hands off generic debt mechanics and statute text to the cross-referenced spokes below.

---

## The "I got a big medical bill" playbook

Work the steps in order. Do **not** pay a medical bill (or put it on a credit card) the moment it arrives; errors are common and discounts are available.

### 1. Don't pay yet; get the two source documents and compare them
- **The bill** is what the provider says you owe. **The EOB (Explanation of Benefits)** is what your **insurer** sends: a line-by-line statement of the billed charge, the plan-negotiated (allowed) amount, what insurance paid, and **what you actually owe**. *An EOB is not a bill.* Match the bill against the EOB; if the bill asks for more than the EOB's patient-responsibility line, that is a red flag.
- **Request an itemized bill** (a.k.a. an itemized statement). The first bill is usually a summary. The itemized bill lists every **CPT/HCPCS procedure code** and revenue code, each charge, and the amount insurance paid. You generally cannot find errors without it, so ask the billing department **in writing**. If they stall or refuse, escalate: many states give you a right to an itemized statement, the federal price-transparency rules require posted charges, and you can put the dispute in a CFPB complaint. (A *superbill* is a different document: an itemized receipt the provider gives **you** to submit to your own insurer for out-of-network reimbursement, not the charge breakdown you want for an audit.)

### 2. Hunt for errors on the itemized bill
Common, real errors to look for:
- **Duplicate charges** (same service billed twice) and **unbundling** (charging separately for items that should be one bundled code).
- **Upcoding** (billed for a more expensive service than you received) and **wrong quantity/units**.
- **Services, meds, or supplies you never received**; wrong dates of service.
- **Charges already paid** by you or your insurer; **balance billing** that should be barred (see §4).
- **Wrong patient demographics / insurance not applied**: a coverage/eligibility mismatch the provider can re-file.

### 3. Dispute errors in writing
- Send a **written** dispute to the provider's billing office; include **copies** (never originals) of the EOB, itemized bill, and any records. Keep a dated log of every call (name, date, what was said).
- If the problem is a **denied or misapplied insurance claim**, you also have appeal rights with your **health plan**: an **internal appeal** and an **external review**. (How plan appeals work → cross-ref `health-insurance-and-coverage`.)
- Unresolved? You can file a complaint with the **CFPB** (financial-product/credit-reporting/collections angle), your **state attorney general / insurance commissioner**, and, for No Surprises Act violations, the federal **No Surprises Help Desk (1-800-985-3059)**.

### 4. Check No Surprises Act protections (see full section below)
If any out-of-network surprise billing is involved (ER, an out-of-network provider at an in-network facility, air ambulance), the balance bill may be **illegal**. If you are **uninsured/self-pay**, check your **good-faith estimate** against the bill.

### 5. Ask about hospital financial assistance / charity care BEFORE negotiating (see 501(r) section)
At a nonprofit hospital this can wipe out or steeply cut the bill, and it can be **retroactive** to a bill you already have. Always screen for this first.

### 6. Then negotiate or set up a plan (see negotiation section)
Ask for the **self-pay / prompt-pay (cash) discount**, propose a **realistic interest-free payment plan**, or offer a **lump-sum settlement**. Get any agreement **in writing** before paying.

### 7. Be careful how you finance it
**Avoid CareCredit-style deferred-interest medical cards** unless you are certain you can pay in full before the promo ends (see that section). If you have one, **HSA/FSA funds** pay qualified medical expenses tax-advantaged.

> **The 60-second version:** get the itemized bill + EOB, find errors, screen for charity care (nonprofit hospitals must offer it), check No Surprises Act protections, then negotiate a cash discount or interest-free plan before you pay or borrow.

---

## Itemized bill vs. EOB vs. CPT codes (the documents)

| Document | Who sends it | What it tells you |
|---|---|---|
| **EOB (Explanation of Benefits)** | Your **insurer** | Billed charge, plan-allowed amount, plan paid, **your responsibility**. *Not a bill.* |
| **Summary bill** | Provider | High-level "amount due." Often the only thing you get first. |
| **Itemized bill / statement** | Provider (on request) | Every **CPT/HCPCS code**, units, charge, and insurer payment; needed to audit. |
| **Superbill** | Provider (on request) | An itemized receipt **you** submit to your own insurer for out-of-network reimbursement. Not the same as a billing-office itemized statement. |

**CPT/HCPCS codes** are the standardized 5-character procedure codes (e.g., physician services); **ICD-10** codes are diagnoses; **revenue codes** group hospital charges. Mismatches between these and what actually happened are where errors hide.

---

## The No Surprises Act (NSA): surprise & balance billing

Federal law effective **January 1, 2022**. **Balance billing** means an out-of-network provider bills you the gap between their charge and what your plan paid. The NSA **bans balance billing** in the common surprise scenarios and limits your cost-share to in-network amounts:

- **Emergency services** at any hospital/ER (in- or out-of-network), including post-stabilization care in many cases.
- **Non-emergency care by out-of-network providers at an in-network facility** (e.g., an out-of-network anesthesiologist, radiologist, or assistant surgeon at your in-network hospital): you generally cannot be balance-billed and did not knowingly waive protection.
- **Air ambulance** (out-of-network). *Ground ambulance is largely NOT covered by the NSA, a known gap.*
- You **cannot be asked to waive** these protections in the situations the law treats as non-waivable; any "consent to waive" form has strict limits and never applies to emergencies.

**Good-faith estimate (GFE) for the uninsured / self-pay.** If you don't have insurance (or won't use it), the provider must give a **written good-faith estimate** of expected charges before a scheduled service. If your final bill is **at least $400 more** than the GFE for a provider/facility, you can use the **patient-provider dispute resolution (PPDR)** process. *(This $400 PPDR threshold is distinct from the $400-and-up "high-balance medical collection" credit-reporting context below.)*

**If you think the NSA was violated:** call the **No Surprises Help Desk at 1-800-985-3059** and don't pay the disputed balance-billed amount while it's contested.

---

## Hospital financial assistance / charity care — IRS 501(r)

Every **nonprofit (501(c)(3)) hospital** must, under **IRS Section 501(r)**, run a financial-assistance program. This is the single most overlooked way to cut a hospital bill. Key requirements:

- **§501(r)(4), Written Financial Assistance Policy (FAP).** The hospital must have a written FAP covering all emergency and medically necessary care, stating the **eligibility criteria** and the discounts/free care available, plus a plain-language summary. **Ask for the FAP by name.** Eligibility is usually income-based (often a multiple of the **Federal Poverty Level**), and many hospitals offer **free care** below one threshold and **discounts** on a sliding scale above it.
- **Presumptive eligibility.** A hospital **may** determine you're eligible using outside data (e.g., enrollment in means-tested programs, prior FAP approval) **without** a full application, and must tell you the basis and how to apply for more generous aid.
- **§501(r)(5), Limit on charges (AGB).** A FAP-eligible patient cannot be charged **more than the "amounts generally billed" (AGB)** to insured patients, never full "gross/chargemaster" charges.
- **§501(r)(6), Billing & collections.** Before "extraordinary collection actions" (e.g., reporting to a bureau, suing, wage garnishment, selling the debt), the hospital must make **reasonable efforts to determine FAP eligibility**, generally including a **240-day application period** from the first post-discharge bill.

**Practical:** ask for the FAP and application *before* paying or negotiating. Assistance can apply **retroactively** to bills already incurred. **Government / for-profit hospitals** aren't bound by 501(r), but many have their own charity-care policies and some states mandate them, so always ask.

---

## Negotiating & settling a medical bill

- **Self-pay / prompt-pay (cash) discount.** Uninsured or paying out of pocket? Ask billing for the self-pay or prompt-pay rate; chargemaster "sticker" prices are routinely discounted. Use the hospital's **price-transparency** tools (federal rules require hospitals to post standard charges, including payer-negotiated and cash rates) and **fair-price benchmarks** to anchor your ask.
- **Payment plan.** Many providers offer **interest-free** in-house plans. Propose a monthly amount you can actually sustain; get the **no-interest** term in writing. (Beware plans routed through a third-party lender; those may carry interest. See medical cards below.)
- **Lump-sum settlement.** Offering a one-time payment for less than the full balance can work, especially once a bill is old or with a collector. Get **"paid in full"** in writing before paying.
- General settlement mechanics, "pay for delete," and the **1099-C canceled-debt tax** consequence (forgiven debt over $600 can be taxable income) live in **`charge-offs-collections-and-debt-resolution`**; cross-ref before settling a large balance.

---

## Medical debt & your credit report (volatile; verify)

**As of 2026, medical debt CAN still appear on credit reports.** Two separate developments; keep them straight:

1. **The nationwide bureaus' voluntary changes (2022–2023), still in effect:**
   - **Paid** medical collections were removed (effective **July 1, 2022**).
   - The wait before an **unpaid** medical collection can appear was extended from 6 months to **1 year** (July 2022), giving you time to sort out insurance/disputes first.
   - Medical collections with an **initial balance under $500** were removed (effective **~April 2023**), which took roughly half of affected consumers' medical collections off their files; **balances of $500+ can still report.**

   These three are a **voluntary industry policy**, not a statute. They are how the bureaus say they handle medical collections, so if a paid or under-$500 medical collection is **still showing** despite the policy, you dispute it as a reporting error (FCRA reinvestigation) rather than enforce a "right" against the furnisher.

2. **The 2024 CFPB medical-debt rule was VACATED in 2025.** The CFPB finalized a rule (Jan 7, 2025) to bar most medical debt from credit reports and stop creditors from considering it. A federal court (**E.D. Texas**) **vacated that rule on July 11, 2025** (*Cornerstone Credit Union League v. CFPB*) as exceeding the Bureau's FCRA authority. **So the federal ban is NOT in force.** Some **states** have enacted their own medical-debt reporting limits, so check your state.

> **Bottom line (as of 2026):** paid and under-$500 medical collections should already be off your report (per the bureaus' voluntary policy); larger unpaid medical collections **can still report after the 1-year wait.** The federal rule that would have changed this was struck down. This is the most fast-moving item in this skill, so **verify current status** and cross-ref the statute spoke (`us-consumer-credit-and-debt-law`) for FCRA dispute mechanics.

**Dispute a medical item that IS on your report** the same way as any tradeline (FCRA reinvestigation), but first dispute the **bill** with the provider. How reports/scores work and how to dispute → `credit-reports-and-scores`; FCRA text → `us-consumer-credit-and-debt-law`.

---

## Medical debt in collections (brief; cross-ref)

A medical bill sent to a third-party collector triggers your **FDCPA / Regulation F** rights: written **debt validation**, limits on contact, and the right to dispute. **Always validate a medical collection** and reconcile it to your EOB/itemized bill; collectors frequently chase amounts that were miscoded, insurance-eligible, or already paid. Detailed collector conduct, validation letters, time-barred "zombie" debt, and being sued → **`debt-collectors-and-fdcpa-rights`**.

---

## Medical credit cards, HSA/FSA, and price transparency

- **CareCredit / medical-credit-card deferred interest: a trap.** "No interest if paid in full" really means **deferred interest**: interest (commonly **~27% APR**) accrues from day one and, if **any** balance remains when the promo period (often 6–24 months) ends, **all** the back-interest is charged retroactively. The CFPB ordered CareCredit to refund **$34.1M** for deceptive enrollment and has warned that these products are pushed on patients at the point of care. **Don't sign one in the provider's office under pressure;** an interest-free hospital plan or financial assistance is almost always better. If you use one, pay it **in full before the promo ends.**
- **HSA / FSA to pay bills.** If you have a **Health Savings Account** or **Flexible Spending Account**, qualified medical expenses (IRS Publication 502 categories) can be paid with **pre-tax** dollars, cheaper than paying with after-tax cash or a credit card. Keep receipts; HSA funds roll over and are yours to keep.
- **Hospital price transparency.** Federal rules require hospitals to publish a machine-readable file of standard charges and a consumer-friendly display of **shoppable services**, including **payer-negotiated and cash-discount prices**. Use these to benchmark before you negotiate.

---

## References / verify current

Primary sources only. **Re-check these before relying on anything above, especially the credit-reporting status, which changed in 2025.**

**No Surprises Act / surprise billing (CMS/CMS.gov):**
- Understand your rights against surprise medical bills: https://www.cms.gov/newsroom/fact-sheets/no-surprises-understand-your-rights-against-surprise-medical-bills
- No Surprises Act overview, rules & fact sheets: https://www.cms.gov/nosurprises/policies-and-resources/overview-of-rules-fact-sheets
- Good-faith estimate & patient-provider dispute resolution: https://www.cms.gov/nosurprises/providers-payment-resolution-with-patients
- No Surprises Help Desk: **1-800-985-3059**

**Hospital financial assistance / charity care (IRS Section 501(r)):**
- Requirements for 501(c)(3) hospitals under the ACA, Section 501(r): https://www.irs.gov/charities-non-profits/charitable-organizations/requirements-for-501c3-hospitals-under-the-affordable-care-act-section-501r
- Financial Assistance Policy & emergency care, §501(r)(4): https://www.irs.gov/charities-non-profits/financial-assistance-policy-and-emergency-medical-care-policy-section-501r4
- Limitation on charges (AGB), §501(r)(5): https://www.irs.gov/charities-non-profits/limitation-on-charges-section-501r5
- Billing & collections, §501(r)(6): https://www.irs.gov/charities-non-profits/billing-and-collections-section-501r6

**Bills, errors, collections & credit reporting (CFPB):**
- What should I do if I can't pay a medical bill? https://www.consumerfinance.gov/ask-cfpb/what-should-i-do-if-i-cant-pay-a-medical-bill-en-2125/
- Know your rights re: medical bills and collections: https://www.consumerfinance.gov/about-us/blog/know-your-rights-and-protections-when-it-comes-to-medical-bills-and-collections/
- Medical debt rule portal (tracks status; rule vacated 2025): https://www.consumerfinance.gov/rules-policy/medical-debt/
- "Paid or under $500 should no longer be on your report": https://www.consumerfinance.gov/about-us/blog/medical-debt-anything-already-paid-or-under-500-should-no-longer-be-on-your-credit-report/
- Medical credit cards & financing plans: https://www.consumerfinance.gov/ask-cfpb/what-should-i-know-about-medical-credit-cards-and-payment-plans-for-medical-bills-en-1827/
- CareCredit $34.1M enforcement action: https://www.consumerfinance.gov/about-us/newsroom/cfpb-orders-ge-carecredit-to-refund-34-1-million-for-deceptive-health-care-credit-card-enrollment/

**Bureaus' voluntary medical-collection changes (2022–2023):**
- Equifax/Experian/TransUnion remove medical collections under $500: https://investor.equifax.com/news-events/press-releases/detail/1286/equifax-experian-and-transunion-remove-medical-collections

**Tax-advantaged payment (IRS):**
- Publication 502, Medical and Dental Expenses (HSA/FSA-qualified): https://www.irs.gov/publications/p502

**General consumer protection (FTC):**
- FTC consumer advice (medical billing & debt collection): https://consumer.ftc.gov/

**Cross-references (installed skills):**
- `debt-collectors-and-fdcpa-rights`: a collector is contacting/suing you; debt validation.
- `charge-offs-collections-and-debt-resolution`: general settlement mechanics, pay-for-delete, 1099-C canceled-debt tax.
- `us-consumer-credit-and-debt-law`: FCRA dispute/reporting statute detail; the vacated CFPB rule.
- `credit-reports-and-scores`: how reports/scores work; disputing a tradeline.
- `health-insurance-and-coverage`: how plans/coverage and coverage-denial appeals work.
