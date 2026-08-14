<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `north-carolina-credit-and-debt-law` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: north-carolina-credit-and-debt-law
description: >-
  North Carolina consumer credit & debt-law overlays: NC rules that change
  outcomes vs. the federal baseline. General legal INFORMATION, not legal advice.
  Spoke of consumer-credit-and-debt.
  TRIGGER: NC debt/credit law; NC Debt Collection Act (G.S. 75-50 to 75-56,
  reaches ORIGINAL creditors); NC SOL (3 yrs, G.S. 1-52); "can my wages be
  garnished in NC" (generally no); NC homestead/property exemptions (G.S.
  1C-1601); NC repossession & deficiency; NC foreclosure (power-of-sale, clerk
  hearing, Ch. 45); NC payday loans (banned) / small-loan caps (G.S. Ch. 53);
  NC Collection Agency Act (G.S. Ch. 58); NC AG complaint.
  SKIP: federal statute text (FCRA/FDCPA/ECOA/TILA) -> us-consumer-credit-and-debt-law;
  a collector calling, or debt-validation/collector tactics (even in NC) ->
  debt-collectors-and-fdcpa-rights; settling a charge-off, pay-for-delete, 1099-C
  -> charge-offs-collections-and-debt-resolution; the executor's NC probate /
  estate-administration process -> estate-planning-and-wills.
metadata:
  changelog:
    - "2026-06-16 sko (--meta --no-sync) — 1 Medium structural: added reciprocal SKIP -> estate-planning-and-wills (the executor's NC probate / estate-administration process; this skill keeps who-owes / creditor-reach / SOL). Resolved circular-SKIP [creditor] flagged by meta-validate into a one-way gradient by wording the edge around the probate-process slice (no 'creditor'/'debts' token). Desc trimmed to <=1000 (992). consumer-credit-and-debt anchor unchanged. Linter 0/0/0. No content passes; --no-sync."
---

# North Carolina Credit & Debt Law

NC-specific overlays on US consumer credit & debt law. This is a **spoke of the
`consumer-credit-and-debt` hub**: the hub and `us-consumer-credit-and-debt-law`
hold the federal baseline (FCRA, FDCPA/Reg F, ECOA, TILA); this skill holds only
the **North Carolina differences that change outcomes**.

> **This is general legal INFORMATION, NOT legal advice.** It does not create an
> attorney-client relationship. NC statutes, dollar amounts, and case law change;
> everything here is **as of 2026** and must be re-verified against ncleg.gov
> before relied on (see *References* below). For a specific situation, consult a
> **NC-licensed attorney** or **Legal Aid of North Carolina** (legalaidnc.org,
> 1-866-219-5262).

## How to use this skill (output rules)

When answering from this skill:

1. **Lead or close with the disclaimer.** State "general legal information, not
   legal advice" in any substantive answer.
2. **Never declare a specific debt time-barred.** Frame it as "this *may* be
   time-barred under G.S. 1-52; confirm with counsel," and flag the
   restart/revival trap (§2) before the person pays or acknowledges an old debt.
3. **On any live dispute, lawsuit, or foreclosure**, surface the **attorney /
   Legal Aid of NC** referral and the deadline that applies (e.g., the 10-day
   foreclosure appeal in §7, or a lawsuit answer deadline).
4. **Cite the G.S. section** and date-stamp dollar amounts ("as of 2026; verify
   at ncleg.gov") so the reader can check current law.
5. **Route, don't answer, out-of-scope asks** per the SKIP list and the
   Cross-references section.

---

## Why North Carolina is different (the headline)

NC is one of the **most debtor-protective states** for ordinary consumer debt.
Three things drive most real-world outcomes:

1. **Wages generally CANNOT be garnished** for ordinary consumer/commercial debt.
2. The **NC Debt Collection Act reaches ORIGINAL creditors**, not just third-party
   collectors (broader than the federal FDCPA), and carries **statutory damages
   + treble damages** via Chapter 75 UDAP.
3. A **short 3-year statute of limitations** on most consumer debt.

---

## 1. NC Debt Collection Act: G.S. Ch. 75, Art. 2 (§§ 75-50 to 75-56)

> **How NC differs from federal:** The federal **FDCPA applies only to
> third-party "debt collectors"** and exempts a creditor collecting its **own**
> debt. The **NC Debt Collection Act has no original-creditor exemption.** Under
> **G.S. 75-50(3)**, a "debt collector" is *"any person engaging, **directly or
> indirectly**, in debt collection from a consumer"* — the **only** carve-out is
> persons subject to **Article 70 of Chapter 58** (licensed collection agencies,
> who are regulated separately, see §5). NC appellate courts have therefore
> applied Article 2's prohibited-practices provisions to **original creditors
> collecting their own debts**. So a NC consumer may have a state claim against a
> bank, hospital, utility, or landlord even where the FDCPA would not apply.

**Definitions (G.S. 75-50):** "consumer" = a natural person who incurred a debt
for **personal, family, household, or agricultural** purposes; "debt" = any
obligation owed or alleged owed by a consumer.

**Prohibited practices** (mirror, but are broader in reach than, the FDCPA):
- **§ 75-51 — Threats and coercion** (violence, harm to reputation/property,
  threatening illegal acts, falsely threatening arrest or criminal prosecution).
- **§ 75-52 — Harassment** (oppressive/abusive conduct, profane language,
  repeated/continuous calls intended to annoy, calls without disclosing identity).
- **§ 75-53 — Unreasonable publication** (improper disclosure of the debt to
  third parties / employers, with narrow exceptions for credit bureaus and
  location info).
- **§ 75-54 — Deceptive representation** (false statement of the debt's character/
  amount/legal status, simulating legal/judicial process, false implication of
  attorney or government affiliation).
- **§ 75-55 — Unconscionable means** (collecting unauthorized fees/charges,
  improper venue, seeking unconscionable waivers).

**Remedies — this is where NC has teeth (G.S. 75-56 + G.S. 75-16):**
- **§ 75-56** makes Article 2 violations the **exclusive UDAP standard** under
  **G.S. 75-1.1** for debt collection, and provides **civil penalties the court
  may allow of not less than $500 nor more than $4,000 per violation**.
- Because a violation **is** an unfair/deceptive act under **G.S. 75-1.1**, the
  consumer gets a **private right of action** with **automatic TREBLE (3×)
  damages** under **G.S. 75-16**, plus **attorneys' fees** under **G.S. 75-16.1**.
- Treble damages + fee-shifting + the per-violation penalty make NC debt-collection
  claims unusually valuable; this is a frequent counterclaim when a consumer is
  sued on a debt.

> Routing: for *how to interact with a collector* (validation letters, cease
> contact, recording, what to say) -> `debt-collectors-and-fdcpa-rights`. For the
> *federal* FDCPA/Reg F text -> `us-consumer-credit-and-debt-law`.

---

## 2. NC statute of limitations on debt: G.S. 1-52 (general 3 years)

- **G.S. 1-52(1): three (3) years** for an action *"upon a contract, obligation
  or liability arising out of a contract, express or implied."* This is the
  default for **most consumer debt** — credit cards, medical bills, open
  accounts, personal loans, auto-loan deficiencies.
- **G.S. 1-47: ten (10) years** for an action on a **sealed instrument** (a
  contract "under seal") or an instrument conveying an interest in real property,
  against the principal. Some older promissory notes are under seal.
- **G.S. 25-3-118 (UCC):** negotiable instruments / notes payable at a definite
  time generally carry a **6-year** limitations period.

> **How NC differs / why it matters:** NC's 3-year clock is **shorter than many
> states'** (and far shorter than the 6 years some other states use for credit
> cards). The SOL is an **affirmative defense** — it does **not** erase the debt;
> it bars the lawsuit if properly raised. Suing on, or threatening suit on, a
> **time-barred** debt is itself actionable: **G.S. 58-70-115** makes it a
> prohibited practice for a debt buyer (or one acting for it) to sue when it
> "knows, or reasonably should know," collection is barred by the SOL.
> **Caution:** a new written promise or (in some circumstances) a payment can
> **restart/revive** the clock — see `charge-offs-collections-and-debt-resolution`
> before paying or acknowledging an old debt.

---

## 3. NC wage garnishment: generally PROHIBITED (a major NC protection)

> **How NC differs from federal:** Federal law (Title III, CCPA) merely **caps**
> wage garnishment. **North Carolina goes further and generally does NOT permit
> wage garnishment for ordinary consumer or commercial debts at all.** A
> judgment creditor on a credit card, medical bill, or personal loan in NC
> **cannot** garnish wages for the debt itself. NC's exemption of *"earnings of
> the debtor for personal services rendered within 60 days"* needed for family
> support (**G.S. 1C-1601(a)... and G.S. 1-362**) underpins this protection.

**Garnishment IS allowed in NC only for specific, statutory exceptions:**
- **Child support / spousal support** — **G.S. 110-136** and related family-law
  provisions.
- **Unpaid state/local & federal taxes** — e.g., NCDOR attachment & garnishment
  under **G.S. 105-368**; IRS levies under federal law.
- **Federal/federally-guaranteed student loans** — administrative wage
  garnishment under federal law (income-based limits apply).
- **Court-ordered restitution**, certain **public-assistance overpayments**, and
  **ambulance / EMS and other public-hospital service** debts via specific
  statutes.
- A creditor who has reduced a debt to judgment may still pursue **other**
  enforcement (bank-account attachment, liens on non-exempt property,
  examinations of judgment debtors) — **wages** are the protected piece.

> **Out-of-state wrinkle:** A creditor with a judgment from **another state**
> where the debtor now lives in NC generally **cannot** use that judgment to
> garnish NC wages; NC's anti-garnishment policy controls wages of NC residents.
> Verify with NC counsel for cross-border facts.

---

## 4. NC property exemptions: G.S. 1C-1601 (date-stamped amounts)

What a judgment debtor can protect from execution/levy. **Amounts and the
subsection numbering below are as of 2026** and reflect the **post-Sept-2025**
version of 1C-1601 (verify at ncleg.gov; amounts are periodically increased):

| Exemption | Amount | Cite |
|---|---|---|
| **Homestead / residence** (real or personal property used as a residence by debtor or a dependent) + a burial plot | **$35,000**; **$60,000** if the debtor is **65 or older** and the property was previously co-owned with a now-deceased spouse/co-owner | G.S. 1C-1601(a)(1) |
| **"Wildcard"** — any unused portion of the homestead exemption, applied to **any** property | up to **$5,000** | G.S. 1C-1601(a)(2) |
| **Motor vehicle** (one) | up to **$3,500** | G.S. 1C-1601(a)(3) |
| **Household goods** — furnishings, clothing, appliances, books, animals, crops held primarily for personal/family/household use | **$5,000** + **$1,000 per dependent** (max **$4,000** added) | G.S. 1C-1601(a)(4) |
| **Tools / implements of trade** of the debtor or a dependent | up to **$2,000** | G.S. 1C-1601(a)(5) |
| **Life insurance** (per the NC Constitution) | as provided | G.S. 1C-1601(a)(6) |
| **Professionally prescribed health aids** | fully exempt | G.S. 1C-1601(a)(7) |
| **Compensation for personal injury / wrongful death** of a person the debtor depended on | exempt | G.S. 1C-1601(a)(8) |

Also exempt by **other** statutes: most **retirement accounts** (IRAs/401(k)s),
**Social Security**, **unemployment** and **workers' comp** benefits, and (as
above) recent **earnings** for family support. A former dollar-capped subsection
was **repealed effective Sept. 1, 2025**, which is why the list above ends at
(a)(8); confirm the current numbering at ncleg.gov before citing a subsection.

> **How NC differs:** Homestead/vehicle amounts are **modest** vs. some states,
> but combined with the **no-wage-garnishment** rule, many NC consumers are
> effectively "**judgment-proof**" for unsecured debt — a judgment exists on
> paper but the creditor has little to collect. NC also lets debtors choose the
> **state** exemptions (NC has **opted out** of the federal bankruptcy exemption
> set, so NC filers use 1C-1601 in bankruptcy).

---

## 5. NC Collection Agency Act: G.S. Ch. 58, Art. 70 (third-party agencies & debt buyers)

Third-party collection agencies and **debt buyers** sit outside **Article 2 of
Chapter 75** (the Debt Collection Act's prohibited-practices scheme: they are the
entities **excluded** from 75-50(3)); instead they fall under **Chapter 58,
Article 70**, administered by the **NC Department of Insurance / Commissioner of
Insurance**. (Their violations are still **UDAP-exposed under G.S. 75-1.1 / 75-16**,
see below.)
- **§ 58-70-1 et seq.** — scope/definitions (includes **"debt buyer"** = a person
  in the business of buying delinquent/charged-off consumer debt for collection).
- **§ 58-70-5** — must apply to the Commissioner for a **permit**; **§ 58-70-20**
  requires a **surety bond ($10,000** for the initial permit). Operating without
  a permit is unlawful.
- **§§ 58-70-90 to 58-70-130 (Part 5)** — **prohibited practices** mirroring (and
  in places exceeding) the FDCPA, applied to agencies/debt buyers. Notably for
  debt buyers: **§ 58-70-115** bars **suing on a time-barred debt** and bars suit
  **without valid documentation of ownership** and **reasonable verification of
  the amount**; **§ 58-70-110/-115** require a written notice with the original
  creditor's name, the account number, a copy of the contract, and an itemized
  accounting before/with suit.

> **How NC differs:** NC's **debt-buyer documentation and SOL rules are stricter
> than the bare FDCPA**, giving consumers concrete defenses when a junk-debt buyer
> sues without a clean chain of title or after the 3-year clock. Violations are
> again **UDAP/treble-damages** exposed via Chapter 75.

---

## 6. NC vehicle repossession & deficiency: UCC Article 9 (G.S. Ch. 25)

NC has adopted **UCC Article 9** at **G.S. 25-9-601 et seq.** Core points:
- **Self-help repossession** is allowed once in default **without a court order**
  *if it can be done **without a breach of the peace*** (G.S. 25-9-609). Breaking
  into a locked garage or repossessing over the debtor's physical objection can
  be a breach of the peace.
- After repossession the secured party must send a **reasonable notification of
  disposition** and conduct a sale that is **commercially reasonable in every
  aspect** (G.S. 25-9-610 to 25-9-614).
- A **deficiency** (loan balance minus sale proceeds) may be pursued, but if the
  creditor **fails to comply** with the notice/commercial-reasonableness rules,
  the deficiency can be **reduced or barred** (G.S. 25-9-625 to 25-9-626;
  consumer "rebuttable presumption" rule).
- A consumer also has the **right to redeem** the vehicle before sale by paying
  the full obligation + costs (G.S. 25-9-623).

> **How NC differs:** The mechanics are standard Article 9, but a resulting
> **deficiency judgment is collectible in NC only against non-exempt property,
> NOT wages** (see §3), and a deficiency claim is itself a **contract action
> subject to the 3-year SOL** (§2).

---

## 7. NC foreclosure: predominantly non-judicial "power of sale": G.S. Ch. 45

Most NC mortgages are **deeds of trust** with a **power-of-sale** clause, so NC
foreclosures are typically **non-judicial** but **supervised by the Clerk of
Superior Court** under **G.S. Ch. 45, Art. 2A** (esp. **G.S. 45-21.16**).

**Process:**
1. Trustee files a **notice of hearing** as a **special proceeding** with the
   **Clerk of Superior Court** in the county where the property sits; the notice
   is served on the borrower/record owner (sheriff or certified mail).
2. At the hearing the **clerk must find six things** to allow the sale
   (**G.S. 45-21.16(d)**): (i) a **valid debt** of which the foreclosing party is
   the holder, (ii) **default**, (iii) **right to foreclose** under the
   instrument, (iv) proper **notice** to those entitled, (v) whether the debt is a
   **"home loan"** and, if so, that the **pre-foreclosure notice** to the borrower
   was given, and (vi) that the **sale is not barred** under **G.S. 45-21.12A**.
3. The clerk's ruling may be **appealed to a Superior/District Court judge within
   10 days**, heard **de novo**.
4. After the sale there is an **upset-bid period: 10 days**, with **successive**
   10-day periods each time a higher (upset) bid is filed (**G.S. 45-21.27**);
   the sale is not final until the upset-bid window closes.

**Deficiency after foreclosure:**
- **G.S. 45-21.36** gives the borrower a **"fair value" defense/offset**: in a
  deficiency suit where the **lender bought** the property at its own sale, the
  borrower may show the property was **worth the debt** (or that the bid was
  substantially below true value) to **defeat or offset** the deficiency.
- **G.S. 45-21.38** bars any deficiency on a **purchase-money** mortgage (seller-
  financed) on the secured property; **G.S. 45-21.38A** addresses deficiency
  limits in certain residential contexts.

> **How NC differs:** NC is a **non-judicial / power-of-sale** state but adds a
> **mandatory clerk hearing with six findings** and a **borrower-favorable
> fair-value deficiency defense**, more borrower protection than a pure
> non-judicial sale. Free help: **Legal Aid of NC** and the **NC Housing Finance
> Agency / State Home Foreclosure Prevention Project**.

---

## 8. NC high-cost & payday lending: effectively banned; small-loan caps: G.S. Ch. 53

> **How NC differs:** **North Carolina effectively prohibits payday lending.**
> NC let its payday-loan authorization expire and the **NC Attorney General /
> DOJ has litigated payday and high-rate lenders out of the state**, including
> out-of-state and "rent-a-bank" schemes that try to evade NC rate caps. There is
> no legal storefront payday loan in NC.

- **NC Consumer Finance Act — G.S. Ch. 53, Art. 15** licenses small consumer
  lenders and **caps rates** on small installment loans (tiered, well below
  payday APRs; see **G.S. 53-176** for rates/maturities/amounts and **G.S.
  53-180** for prohibited practices/limitations). Unlicensed lending above the
  cap is **usurious and unenforceable**.
- General **usury** limits live in **G.S. Ch. 24** (e.g., the legal rate and
  contract-rate ceilings under **G.S. 24-1, 24-1.1**).
- A loan made in violation of NC rate caps can expose the lender to **UDAP /
  Chapter 75** liability and the consumer may avoid the unlawful charges.

> The **NC DOJ "Payday Loans"** page warns consumers that payday/quick-cash and
> online high-rate loans are **illegal in NC** and urges reporting them.

---

## 9. NC Attorney General / Department of Justice: consumer-protection role & complaints

- The **NC DOJ (ncdoj.gov)** Consumer Protection Division enforces **Chapter 75
  UDAP**, sues abusive collectors and predatory/payday lenders, and offers
  **informal mediation** of consumer complaints. It does **not** act as your
  private lawyer or give individual legal advice.
- **File a complaint:** **ncdoj.gov/file-a-complaint/consumer-complaint/**
  (online), or call **1-877-5-NO-SCAM (1-877-566-7226)** in-state /
  **919-716-6000**. Filing a DOJ complaint is **free** and is a common first step
  alongside (not a substitute for) a private Ch. 75/Article 2 claim.
- **Self-help / courts:** **NC Judicial Branch — nccourts.gov** has self-help on
  **foreclosures**, **small claims** (money judgments up to the statutory limit),
  and forms. **Legal Aid of North Carolina (legalaidnc.org, 1-866-219-5262)**
  serves low-income consumers.

---

## Cross-references

- **Federal statute text & rights** (FCRA, FDCPA/Reg F, ECOA/Reg B, TILA/Reg Z,
  the contested CFPB medical-debt rule) -> **`us-consumer-credit-and-debt-law`**.
- **Dealing with a collector** — validation letters, cease-contact, recording,
  being sued, FDCPA tactics -> **`debt-collectors-and-fdcpa-rights`**.
- **Resolving defaulted debt** — charge-offs, settlement, pay-for-delete, debt
  buyers, **1099-C** cancellation-of-debt tax -> **`charge-offs-collections-and-debt-resolution`**.
- **How a score/report works or is rebuilt** -> **`credit-reports-and-scores`**,
  **`improving-and-rebuilding-credit`**.
- Parent hub: **`consumer-credit-and-debt`** family. The federal baseline lives
  in **`us-consumer-credit-and-debt-law`**; this spoke holds only the NC deltas.

---

## References / verify current law (re-check before relying: as of 2026)

NC statutes and dollar amounts change; **confirm every citation at ncleg.gov**.

**NC Debt Collection Act (Ch. 75, Art. 2):**
- G.S. 75-50 (definitions): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/BySection/Chapter_75/GS_75-50.html
- Article 2 full text (§§ 75-50 to 75-56): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/ByArticle/Chapter_75/Article_2.html
- G.S. 75-56 (enforcement; UDAP standard; $500–$4,000 penalty): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/ByArticle/Chapter_75/Article_2.html
- G.S. 75-1.1 / 75-16 / 75-16.1 (UDAP, treble damages, attorneys' fees): https://www.ncleg.gov/Laws/GeneralStatuteSections/Chapter75

**Statute of limitations:**
- G.S. 1-52 (three years): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_1/GS_1-52.pdf
- G.S. 1-47 (ten years, sealed instruments): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/BySection/Chapter_1/GS_1-47.html
- G.S. 25-3-118 (UCC, negotiable instruments): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_25/GS_25-3-118.pdf

**Wage garnishment & exemptions:**
- G.S. 1C-1601 (exempt property / homestead): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/BySection/Chapter_1C/GS_1C-1601.html
- Chapter 1C, Article 16 (enforcement of judgments): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/ByArticle/Chapter_1C/Article_16.pdf
- G.S. 110-136 (child-support garnishment): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_110/GS_110-136.pdf
- G.S. 105-368 (tax attachment & garnishment): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_105/GS_105-368.pdf

**Collection Agency Act (Ch. 58, Art. 70):**
- Article 70 (collection agencies / debt buyers): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/ByArticle/Chapter_58/Article_70.html
- G.S. 58-70-5 (permit application): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_58/GS_58-70-5.pdf
- G.S. 58-70-115 (debt-buyer prohibited acts; time-barred suits): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/BySection/Chapter_58/GS_58-70-115.html

**Repossession (UCC Art. 9):**
- Chapter 25 (UCC): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/ByChapter/Chapter_25.html

**Foreclosure (Ch. 45):**
- Chapter 45, Article 2A (power of sale): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/ByArticle/Chapter_45/Article_2A.html
- G.S. 45-21.16 (notice & hearing; clerk findings): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_45/GS_45-21.16.pdf
- G.S. 45-21.27 (upset bids): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/BySection/Chapter_45/GS_45-21.27.html
- G.S. 45-21.36 (deficiency / fair-value defense): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/BySection/Chapter_45/GS_45-21.36.html
- NC Judicial Branch — Foreclosures (self-help): https://www.nccourts.gov/help-topics/housing/foreclosures

**Payday / high-cost / small loans (Ch. 53, Ch. 24):**
- G.S. 53-176 (Consumer Finance Act — rates/maturities/amounts): https://www.ncleg.gov/EnactedLegislation/Statutes/PDF/BySection/Chapter_53/GS_53-176.pdf
- Chapter 24 (interest / usury): https://www.ncleg.gov/EnactedLegislation/Statutes/HTML/ByChapter/Chapter_24.html
- NC DOJ — Payday Loans: https://ncdoj.gov/protecting-consumers/credit-and-debt/payday-loans/

**NC DOJ / Attorney General:**
- Consumer credit & debt hub: https://ncdoj.gov/protecting-consumers/credit-and-debt/
- File a consumer complaint: https://ncdoj.gov/file-a-complaint/consumer-complaint/

*Statutes and amounts stated as of 2026. Verify against ncleg.gov before relying.
This is general information, not legal advice — consult a NC-licensed attorney or
Legal Aid of North Carolina for your situation.*
