<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `improving-and-rebuilding-credit` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: improving-and-rebuilding-credit
description: 'Playbook to build, improve, or rebuild a US consumer credit score: on-time-payment systems, lowering utilization, account age, inquiries, credit mix, disputing errors, goodwill letters, and rebuilding tools, with realistic timelines. Educational, not advice. Spoke of consumer-credit-and-debt. TRIGGER: raise/build/rebuild a score, lower utilization, AZEO, dispute credit-report errors, goodwill letter, secured card, credit-builder loan, authorized user, Experian Boost, build a thin/no file, rebuild after bankruptcy or charge-offs, "is credit repair a scam", CROA. SKIP: how scoring works, factor weights, FICO vs VantageScore, inquiries, freezes/fraud alerts -> credit-reports-and-scores; statute text (FCRA/CROA/FDCPA, suing) -> us-consumer-credit-and-debt-law; settling charge-offs & collections (pay-for-delete, 1099-C) -> charge-offs-collections-and-debt-resolution; stopping collectors -> debt-collectors-and-fdcpa-rights; budgeting / debt-payoff planning -> budgeting-and-saving.'
metadata:
  changelog:
  - 2026-06-16 sko --meta --no-sync — structural pass vs new consumer-finance family; 0 High 0 Med (meta-validate clean). Added 1 SKIP edge (-> budgeting-and-saving); trimmed description to 989 chars within 1000 cap; no content passes (already CLEAN prior).
---

# Improving and Rebuilding Credit

General educational information, **not financial advice**. Credit outcomes depend
on your full file and the scoring model a lender uses; **results vary and no point
gain is guaranteed.** Current **as of 2026**, score factors, products, and rules
change, so verify against the primary sources at the end. This is a spoke of the
**`consumer-credit-and-debt`** hub.

Everything actionable here ladders back to the five score factors (covered in
depth in `credit-reports-and-scores`). The approximate FICO weights are the map of
where your effort pays off:

| Factor | ~Weight | Lever in this skill |
|---|---|---|
| Payment history | 35% | On-time payment systems / autopay |
| Amounts owed (mostly utilization) | 30% | Lower utilization; AZEO |
| Length of credit history | 15% | Keep old accounts open |
| New credit (applications/inquiries) | 10% | Limit new applications |
| Credit mix | 10% | Add an installment loan only if it fits |

---

## The prioritized action playbook

Work top-down. The top two factors are ~65% of a FICO score, so they get the
effort first.

### 1. Never miss a payment again (highest leverage, 35%)
- **Get current and stay current.** A single 30-day late mark can drop a good
  score sharply and stays ~7 years. Nothing else matters as much.
- **Automate.** Set **autopay for at least the minimum** on every account, then
  pay more manually. Autopay-minimum is the safety net that prevents a missed
  payment; it is not a strategy to carry a balance.
- **Add reminders** a few days before each due date as a backstop. CFPB: "on
  time" means the payment **reaches the company by the due date** — if mailing,
  send several days early.
- If money is tight, pay the account whose lateness would be **newly reported**
  first (a payment <30 days late is usually not reported; 30+ days is).

### 2. Lower your utilization (30%, and it updates fast)
Utilization = revolving balances ÷ revolving credit limits. It is recalculated
every cycle, so this is the **fastest** lever to move a score.
- **Target <30%, ideally <10%**, both overall and on each individual card. (CFPB:
  some experts say under 30%, others under 10%.)
- Levers, fastest first: **pay before the statement closes** (the statement
  balance is usually what gets reported, not the due-date balance); **pay
  multiple times a month**; **request a credit-limit increase** (prefer a soft-pull
  increase); **spread balances** across cards instead of maxing one.
- **AZEO ("All Zero Except One")** — a short-term optimization for the ~30-45 days
  **before a mortgage or major loan application**: let **every card report a $0
  balance except one**, which reports a small balance (commonly 1-9%). Mortgage-era
  FICO models (FICO 2/4/5) are sensitive to the *number of cards carrying a
  balance*, so one-card-with-a-small-balance often scores best. Caveats: it keys
  off the **statement/reporting date**, not the due date; it is a temporary tactic,
  **not** something you must do every month; and never let utilization hit a hard 0%
  across all cards long-term.
- Note: installment loans (auto, mortgage, student, credit-builder) factor into
  "amounts owed" too, but **revolving utilization dominates** and is what to manage.

### 3. Keep old accounts open (15% — age)
- **Don't close your oldest cards.** Length of history uses the age of your oldest
  account, newest account, and the **average age** of all accounts. Closing a card
  also **removes its limit**, which can spike utilization (a double hit).
- Keep no-annual-fee cards alive with a tiny recurring charge + autopay.
- For an annual-fee card you don't want, ask to **product-change/downgrade** to a
  no-fee card from the same issuer to preserve the account's age rather than
  closing it.

### 4. Limit new applications and hard inquiries (10%)
- Each application is usually a **hard inquiry** (small, temporary ding); several
  in a short window compounds and shortens average account age.
- **Rate-shopping exception:** multiple inquiries for the *same* loan type
  (reliably auto and mortgage; some models also student) within a focused window
  (commonly ~14-45 days depending on model) typically count as **one** inquiry. Do
  rate shopping in a tight cluster.
- Space out new credit-card applications by several months. Apply only when the
  account genuinely helps.

### 5. Credit mix (10% — lowest priority)
- A mix of revolving (cards) and installment (loans) can help, but **do not take on
  debt or interest just to "improve mix."** A **credit-builder loan** (below) is the
  sane way to add installment history if you have none.

---

## Fixing what's wrong: disputes, goodwill, negotiation

### Disputing errors (the practical FCRA process)
Errors are common and you have the right to fix them **for free**. Studies and CFPB
guidance: dispute with **both** the credit reporting company **and** the furnisher
(the bank/lender/landlord that supplied the data).
1. **Get your reports.** Free weekly reports from all three bureaus at
   **AnnualCreditReport.com** (the only federally authorized free source).
2. **Dispute with each bureau** (Equifax, Experian, TransUnion) that shows the
   error — online, by phone, or by mail. Include your contact info, the disputed
   item(s) and account numbers, a clear explanation, a marked-up copy of the report,
   and **copies** (never originals) of supporting documents.
3. **Also dispute with the furnisher** directly, in writing.
4. **Use certified mail, return receipt requested** for a paper trail (online keeps
   its own log).
5. **Timeline:** the investigation generally must complete within **~30 days** (up
   to **45 days** if you dispute after receiving your free annual report or add
   information mid-investigation). If the item is inaccurate or unverifiable, it must
   be **corrected or deleted** and the furnisher must notify the bureaus. If it
   stays, you can add a brief statement of dispute to your file.
- **If it comes back "verified" but you believe it is wrong:** re-dispute with new
  supporting documents, escalate to the furnisher, file a complaint with the CFPB
  (consumerfinance.gov/complaint), and for a serious or repeated error consult a
  consumer attorney (`us-consumer-credit-and-debt-law`).
- **Watch the 7-year clock and reinsertion.** The ~7-year reporting window runs from
  the **date of first delinquency** (FCRA), which should not move. If a deleted item
  reappears or the delinquency date is reset to look newer ("re-aging"), dispute it.
- Reality check: **accurate** negative info cannot be disputed away. Disputing
  correct items as a tactic is what scam "credit repair" sells (see below).

### Goodwill letters (for accurate, one-off lates)
- A **goodwill letter** asks a creditor to remove an isolated late mark **as a
  courtesy** — most effective when you have **one slip** on an otherwise clean,
  long history and have since paid on time.
- Write a **personal** note (not a template), briefly explain the cause, what you
  fixed, and ask politely; addressing a higher-level contact can help.
- It is **discretionary**: creditors are not required to agree, and many decline by
  policy. Free to try; manage expectations.

### Negotiating with creditors
- For the specifics of settling charge-offs and collections (settlement
  percentages, **pay-for-delete**, lump-sum vs. payment plans, 1099-C tax on
  forgiven debt), use **`charge-offs-collections-and-debt-resolution`** (planned
  sibling). For stopping collector contact and your rights, use
  **`debt-collectors-and-fdcpa-rights`** (planned sibling).
- Score-relevant note: **paying a collection does not automatically delete it.**
  Newer models (FICO 9/10, VantageScore 3/4) ignore **paid** collections and weigh
  **medical** collections less, but older models a lender may still use can keep
  counting it. Get any delete agreement **in writing before paying.**

---

## Rebuilding tools (build a thin/no file or recover after damage)

- **Secured credit card:** you post a refundable deposit that equals your (small)
  limit; it reports like a normal card. The fastest first step for a thin file or
  post-derogatory rebuild. Pick one that **reports to all three bureaus**, has
  **low/no annual fee**, and ideally **graduates** to unsecured. Keep utilization
  low and pay in full. Watch high fees/APRs on subprime cards.
- **Credit-builder loan:** a small loan whose proceeds sit in a locked savings
  account; **you make the fixed payments first** and receive the saved principal at
  the end (you pay interest/fees on top of the payments, and some lenders rebate
  part of the interest). Builds **installment** and payment history with little
  risk. Offered mainly by credit unions and **CDFIs (Community Development
  Financial Institutions)**; confirm it reports to all three bureaus.
- **Authorized user (AU) / tradelines:** being added to a **trusted person's**
  old, low-utilization, perfectly-paid card can import that history to your file.
  Cautions: it only helps if **that account stays clean** (their late payment or
  high balance can hurt you); some scoring models discount AU tradelines; and you
  can be removed anytime.
  - **Warning — paid "tradeline" schemes:** buying AU spots on strangers' accounts
    from a broker is a **manipulation tactic** lenders actively detect and discount,
    can look like fraud, and wastes money. Avoid.
- **Rent / utility / subscription reporting:** **Experian Boost** (and rent-reporting
  services) can add on-time **utility, telecom, rent, and some streaming** payments
  to your file. Caveats: Boost affects **only your Experian** file and only certain
  FICO/VantageScore versions; rent reporting only helps if the **landlord/service
  actually reports** and the bureau/model counts it. Useful for **thin files**;
  modest for already-established files.

---

## Realistic timelines

Educational ranges, not promises:
- **Utilization paydown:** next reporting cycle (**~30-60 days**) once lower balances
  report — the fastest visible change.
- **New positive account (secured card / builder loan):** on-time history starts
  helping in **~3-6 months**; a meaningful track record takes **6-12 months**.
- **Hard inquiry:** minor effect fades within **~12 months**; drops off the report at
  **~24 months**.
- **Late payment:** most damaging early, **decreasing impact over time**, off the
  report at **~7 years**.
- **Collections / charged-off accounts:** stay **~7 years**; impact lessens as they
  age.
- **Chapter 7 bankruptcy:** reports **~10 years**; **Chapter 13: ~7 years.** Scores
  can begin recovering well before that with active rebuilding.
- **Thin / no file → first score:** typically need **~3-6 months** of at least one
  reporting account before a score can even be generated.

**Rebuilding after bankruptcy or charge-offs:** confirm the discharged/settled debts
report correctly ($0 balance, "included in bankruptcy"/"settled"); dispute anything
wrong; open **one** secured card and/or a credit-builder loan; keep utilization low;
automate payments; and **be patient**: recovery is driven by *new* positive history
accumulating, not by paying anyone to erase the past.

---

## Legit help vs. "credit repair" scams (CROA)

There are three legitimate paths and one trap.

**Legitimate:**
1. **DIY:** everything in this skill. Disputes, goodwill letters, and rebuilding
   are things **you can do yourself for free.**
2. **Nonprofit credit counseling:** e.g., **NFCC** member agencies (nfcc.org) and
   **HUD-approved housing counselors**. The **initial budget session is typically
   free**; they offer budgeting and, for unaffordable debt, a **Debt Management Plan
   (DMP)** (one monthly payment, often reduced interest). Find one via the **NFCC
   Agency Finder**.
3. **A consumer attorney** for genuine FCRA violations (see
   `us-consumer-credit-and-debt-law`).

**The trap — advance-fee "credit repair":** Under the **Credit Repair Organizations
Act (CROA)** it is **illegal** for a credit-repair company to:
- **Charge any fee before services are fully performed** (no advance/upfront fees);
- **Lie about what they can do** or promise to remove **accurate, timely** negative
  information (no one can);
- Fail to give you a **written contract**, a description of services, and your
  **3-day right to cancel.**

**Red flags = walk away:** demands payment up front; guarantees a specific point
jump or to remove accurate items; tells you **not to contact the bureaus
yourself**; or pushes a **"new credit identity" / CPN (credit privacy/profile
number)**. Using a **CPN** in place of your SSN to apply for credit is **fraud** and
can carry **fines or prison** — these "numbers" are often stolen SSNs.

---

## Myth list (common false beliefs)

- **"Closing a credit card helps my score."** False: closing usually **lowers**
  available credit (raising utilization) and can shrink average account age.
- **"Paying off a collection always removes it / instantly fixes my score."**
  False: it typically **stays ~7 years**; older models still count paid
  collections. Get **pay-for-delete in writing** if offered.
- **"Checking my own credit hurts my score."** False: that's a **soft pull**;
  soft inquiries (including your own checks and prequalifications) **never** affect
  your score. Only **hard** inquiries from applications do.
- **"I need to carry a balance / pay interest to build credit."** False: whether
  you revolve a balance is **not** a scoring factor. **Pay in full**; on-time
  payment and low utilization are what build score (and you avoid interest).
- **"My income / age / race / marital status affects my credit score."** False:
  none of these are in the score; only the data on your credit report is.
- **"A credit-repair company can legally erase accurate bad info."** False: see
  CROA above; no one can remove accurate, timely information.
- **"AZEO must be done every month."** False: it's a **short-term, pre-application**
  optimization, not a maintenance routine.
- **"Disputing accurate items is a smart strategy."** False/risky: it's the core of
  scam credit repair, items return when re-verified, and frivolous disputes can be
  ignored.

---

## Do / Don't

**Do**
- Set autopay-for-minimum on everything, then pay in full.
- Pay balances **before the statement closes** to report low utilization.
- Pull all three reports free at AnnualCreditReport.com and dispute real errors.
- Keep your oldest accounts open and active.
- Rebuild with a **secured card** and/or **credit-builder loan** that report to all
  three bureaus.
- Use **free nonprofit (NFCC/HUD) counseling** when debt is the real problem.

**Don't**
- Don't pay anyone an **advance fee** to "repair" credit, or buy a **CPN/tradeline**.
- Don't close your oldest card or chase a better "mix" with new debt.
- Don't apply for several cards at once.
- Don't dispute information you know is accurate.
- Don't assume paying a collection erases it — get deletion **in writing** first.

---

## References / verify current (as of 2026)

Primary, authoritative sources — re-check, as products and rules change.

- **CFPB — How to rebuild your credit:**
  https://www.consumerfinance.gov/consumer-tools/credit-reports-and-scores/how-to-rebuild-your-credit/
- **CFPB — How do I get and keep a good credit score?:**
  https://www.consumerfinance.gov/ask-cfpb/how-do-i-get-and-keep-a-good-credit-score-en-318/
- **CFPB — How do I dispute an error on my credit report?:**
  https://www.consumerfinance.gov/ask-cfpb/how-do-i-dispute-an-error-on-my-credit-report-en-314/
- **CFPB — Will paying off my balance every month improve my score?:**
  https://www.consumerfinance.gov/ask-cfpb/will-paying-off-my-credit-card-balance-every-month-improve-my-score-en-1293/
- **FTC (consumer.ftc.gov) — Credit Repair / Fixing Your Credit FAQs (incl. CROA, CPN scams):**
  https://consumer.ftc.gov/articles/0225-credit-repair-scams
- **FTC — Credit Repair Organizations Act (statute overview):**
  https://www.ftc.gov/legal-library/browse/statutes/credit-repair-organizations-act
- **MyFICO — What's in your credit score (factor weights):**
  https://www.myfico.com/credit-education/whats-in-your-credit-score
- **MyFICO — Top credit-score myths:**
  https://www.myfico.com/credit-education/blog/5-credit-score-falsehoods
- **NFCC — nonprofit credit counseling / DMP / Agency Finder:**
  https://www.nfcc.org/  •  https://www.nfcc.org/agency-finder/
- **Experian (consumer education) — credit-builder loans vs. secured cards; thin file; goodwill letters:**
  https://www.experian.com/blogs/ask-experian/should-i-get-a-credit-builder-loan-or-a-secured-credit-card/  •  https://www.experian.com/blogs/ask-experian/what-is-goodwill-letter/
- **AnnualCreditReport.com — the only federally authorized free credit reports:**
  https://www.annualcreditreport.com/

**Cross-references:** how scoring works -> `credit-reports-and-scores`; statute/legal
detail (FCRA, CROA, FDCPA, suing) -> `us-consumer-credit-and-debt-law`; settling
charge-offs & collections -> `charge-offs-collections-and-debt-resolution`; dealing
with collectors -> `debt-collectors-and-fdcpa-rights`; hub ->
`consumer-credit-and-debt`.
