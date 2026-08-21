---
name: medicaid-drug-rebate-340b
description: >-
  Medicaid pharmaceutical pricing and drug access — Medicaid Drug Rebate Program (MDRP, Section 1927), best price rule, AMP, inflation rebates, supplemental rebates, covered outpatient drugs, prior authorization for drugs, and the 340B Drug Pricing Program (safety-net entity discounts, covered entities, 340B ceiling price, duplicate discounts, 340B controversy and CMS-340B tensions). TRIGGER: Medicaid drug rebate, MDRP, Medicaid drug pricing, best price Medicaid, AMP average manufacturer price, Medicaid inflation rebate, supplemental drug rebates, Medicaid drug prior authorization, step therapy Medicaid, preferred drug list PDL, 340B program, 340B covered entity, 340B ceiling price, 340B hospital, 340B discount, 340B controversy, duplicate discount 340B. SKIP: general Medicaid program structure → medicaid-program-structure; Medicaid managed care → nc-health-services-medicaid-managed-care; Medicaid fraud → medicaid-fraud-waste-abuse.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [medicaid, drug-rebate, MDRP, 340B, pharmaceutical-pricing, best-price, AMP, PDL, prior-authorization]
related_skills:
  - medicaid-program-structure
  - medicaid-fmap-federal-financing
  - medicaid-fraud-waste-abuse
  - nc-health-services-medicaid-managed-care
  - nc-health-services-fqhcs
---

# Medicaid Drug Rebate Program and 340B

> **Educational only — NOT legal, medical, or pharmaceutical advice.** Drug pricing rules are complex and change frequently.

**verified-as-of: 2026-06-22** *(compiled from training knowledge; 340B contract pharmacy litigation actively evolving; verify current CMS/HRSA guidance and IRA drug negotiation updates)*

---

## 1. Medicaid Drug Rebate Program (MDRP) Overview

The **Medicaid Drug Rebate Program** (Section 1927 of the Social Security Act) requires pharmaceutical manufacturers to enter rebate agreements with CMS as a condition of having their drugs covered by Medicaid. Without a signed rebate agreement, states may not cover the drug under Medicaid (with narrow exceptions for single-source life-saving drugs).

**Scale:** MDRP generates ~$30–35 billion in annual rebates to states and the federal government, making it one of the largest drug pricing tools in the US.

---

## 2. Key MDRP Pricing Concepts

### Average Manufacturer Price (AMP)
The **average price paid to the manufacturer** by wholesalers/distributors for drugs distributed to retail pharmacies, net of discounts and allowances. Manufacturers report AMP monthly and quarterly to CMS. AMP is used to calculate rebate amounts and (for generics) retail pharmacy reimbursement benchmarks.

### Best Price
Manufacturers must report the **lowest price** they charge to any private-sector customer (excluding 340B, federal programs, and certain other exempt customers). Best price sets a floor for Medicaid rebates: Medicaid cannot pay more than the best commercial price.

**Best Price rule controversies:**
- Manufacturers argue that including value-based contracts, patient assistance programs, and rebates in best price calculations creates perverse disincentives to offering discounts to commercial payers
- CMS has adjusted best price rules to accommodate outcomes-based/value-based arrangements

### Basic Rebate
For **innovator (brand) drugs:** Greater of 23.1% of AMP or AMP minus best price
For **non-innovator (generic) drugs:** 13% of AMP

### Inflation Rebate (Penalty Rebate)
If a brand drug's AMP increases **faster than inflation (CPI-U)** from its base period, the manufacturer owes an additional rebate equal to the excess price increase. This creates a ceiling on price growth for Medicaid — a major cost-control mechanism strengthened by the Inflation Reduction Act (2022) which extended inflation rebates to Medicare Part D.

### Supplemental Rebates
States may negotiate **supplemental rebates** directly with manufacturers on top of the federal rebate, often in exchange for preferred drug list (PDL) placement. Supplemental rebate revenue is split between state and federal government per FMAP formula.

---

## 3. Drug Access Mechanisms Under Medicaid

### Covered Outpatient Drug Definition
MDRP applies to "covered outpatient drugs" — drugs dispensed at retail pharmacies or outpatient settings. Drugs administered in clinical settings (e.g., hospital inpatient, physician-administered infusions) are covered under different Medicaid payment rules and may not receive MDRP rebates the same way.

### Preferred Drug Lists (PDLs)
States maintain PDLs — formularies that give preferred status to certain drugs within a therapeutic class, typically the drug with the best rebate + net cost. Non-preferred drugs may require prior authorization (PA).

**CMS rules on access:** States must ensure PDL restrictions don't create unreasonable access barriers. Certain drug classes must have minimum number of drugs available (e.g., antipsychotics for mental illness, antiretrovirals for HIV — states cannot restrict below established minimums).

### Prior Authorization (PA) for Drugs
Medicaid programs routinely require PA for high-cost, non-preferred, specialty, or newly approved drugs. PA requirements must not unreasonably limit access to medically necessary treatment. Step therapy (try lower-cost drug first) is also common but subject to override for clinical necessity.

### Managed Care and MDRP
When Medicaid beneficiaries are enrolled in managed care, states can require MCOs to use PDLs and collect rebates. Federal rules allow states to require MCOs to participate in supplemental rebate agreements and pass rebate savings into managed care capitation rates (carve-in vs carve-out).

---

## 4. The 340B Drug Pricing Program

### What 340B Is
The **340B Drug Pricing Program** (Section 340B of the Public Health Service Act) requires manufacturers that participate in MDRP to sell outpatient drugs to **covered entities** at or below a federally calculated **ceiling price** — typically 25–50% off commercial prices.

### Who Are Covered Entities?
Eligible covered entities include:
- Federally Qualified Health Centers (FQHCs) and look-alikes
- Ryan White HIV/AIDS program grantees
- Disproportionate Share Hospitals (DSH hospitals serving >11.75% Medicaid inpatients)
- Children's hospitals, critical access hospitals, rural referral centers (with waivers)
- Tribal/Indian health clinics
- Black lung clinics, tuberculosis programs

### 340B Ceiling Price Calculation
Ceiling price = AMP minus MDRP rebate (simplified). For many drugs, ceiling price is substantially below commercial prices, enabling covered entities to stretch limited resources.

### How Covered Entities Use 340B
Entities can purchase 340B drugs and bill payers (including Medicaid, Medicare, private insurance) at normal rates. The **"spread"** between acquisition cost and reimbursement is retained to fund services — often cross-subsidizing care for uninsured patients or expanding clinic services.

### 340B and Duplicate Discounts
A **duplicate discount** occurs if a covered entity receives both a 340B price on a drug AND a Medicaid rebate for the same drug on the same patient encounter. This is **prohibited by law**. States must have systems to prevent Medicaid MCOs and FFS programs from claiming rebates on 340B drugs.

**Carve-in vs carve-out:** To avoid duplicate discounts, states can "carve out" 340B from Medicaid rebates (covered entities can't use 340B for Medicaid patients) or "carve in" (use 340B but track carefully to avoid duplicate rebate claims). NC's approach has evolved with managed care expansion.

### 340B Controversies (Current as of 2026)

**Contract pharmacies:** Covered entities can dispense 340B drugs through contract pharmacies (third-party retail pharmacies). Many manufacturers began restricting contract pharmacy 340B access (2020–2021), leading to litigation. Circuit courts have split on manufacturer authority to restrict (e.g., 3rd Cir. in *Sanofi-Aventis U.S. LLC v. HHS*, 2023); SCOTUS's major 340B decision was *AHA v. Becerra* (2022, Medicare Part B payment cuts — separate issue); contract pharmacy litigation continues in lower courts as of 2026.

**Program expansion concerns:** HRSA OIG reports have found some covered entities used 340B savings for purposes beyond serving low-income patients; Congress has debated stronger accountability/audit requirements.

**Hospital 340B:** Large academic medical centers with DSH status have substantially grown 340B purchasing. Policy debate: Are large hospitals the intended beneficiaries? Proposals exist to cap 340B eligibility at hospitals with demonstrated charity care levels.

**NC FQHCs and 340B:** NC's 42 FQHCs rely heavily on 340B to fund operations serving ~500,000 patients with limited Medicaid and Medicare revenue per visit. 340B restrictions by manufacturers directly threaten FQHC financial stability in NC.

---

## 5. MDRP + 340B Interaction

| Scenario | MDRP rebate? | 340B price? |
|---|---|---|
| FFS Medicaid, non-340B drug | Yes | No |
| FFS Medicaid, 340B drug (carve-out state) | Yes | No (carve-out) |
| FFS Medicaid, 340B drug (carve-in with tracking) | No (tracked) | Yes |
| MCO Medicaid, 340B drug (carve-in MCO) | No | Yes |
| Medicare Part B (physician-administered) | No MDRP | Yes |
| Commercial payer | No MDRP | Yes (if CE) |

---

## Sources

- CMS, *Medicaid Drug Rebate Program* (Section 1927 of the SSA; CMS drug rebate guidance)
- HRSA, *340B Drug Pricing Program* (hrsa.gov/opa)
- MedPAC/MACPAC, *340B Program Overview and Issues* (2021–2025 reports)
- KFF, *Medicaid Drug Pricing* — issue brief series (2024)
- 3rd Cir., *Sanofi-Aventis U.S. LLC v. HHS* (2023) — contract pharmacy restrictions; related circuit court litigation ongoing
- SCOTUS, *AHA v. Becerra*, 596 U.S. 724 (2022) — Medicare Part B 340B payment cuts (separate 340B issue)
- OIG, *340B Program Integrity: Contract Pharmacy Arrangements* (2023)
- GAO, *Drug Pricing: Manufacturer Discounts in the 340B Program Offer Benefits, But Federal Oversight Needs Improvement* (2022)
