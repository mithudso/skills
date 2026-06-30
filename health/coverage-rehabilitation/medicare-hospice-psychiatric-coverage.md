---
name: medicare-hospice-psychiatric-coverage
description: >-
  Medicare hospice benefit as it applies to individuals with psychiatric
  diagnoses and behavioral health needs: eligibility for hospice with terminal
  psychiatric illness (schizophrenia, late-stage dementia with behavioral
  symptoms, anorexia nervosa), the 6-month prognosis requirement for BH
  conditions, Medicare hospice interdisciplinary team (IDT) requirements
  including social work and chaplaincy for BH, psychiatric medication coverage
  under the hospice per-diem, Medicare hospice coverage of psychotherapy and
  behavioral symptom management, concurrent care, respite, and the tension
  between curative BH treatment and hospice election. Educational only, NOT
  medical/insurance advice; verified-as-of 2026-06-22.
  TRIGGER: Medicare hospice mental illness, Medicare hospice psychiatric,
  Medicare hospice dementia behavioral symptoms, hospice anorexia nervosa,
  hospice schizophrenia prognosis, Medicare hospice 6-month rule psychiatric,
  hospice psychotropic medications, hospice behavioral health services, hospice
  social work mental health, Medicare concurrent care hospice mental health,
  hospice election mental health treatment, Medicare hospice interdisciplinary
  team behavioral health.
  SKIP: Medicare hospice overall → see consumer-finance (not built in health);
  inpatient psychiatric separate from hospice → medicare-mental-health-sud-coverage;
  Medicaid LTSS nursing facility BH → medicaid-ltss-hcbs-rehab; discharge
  planning from hospice → mhpaea-prior-auth-discharge.
version: "1.1.0"
updated: "2026-06-22"
category: health
tags: [medicare, hospice, behavioral-health, psychiatric, coverage, end-of-life]
hub: health
---

# Medicare Hospice Benefit and Behavioral Health

## Overview

Medicare hospice covers individuals with terminal illnesses who forgo curative treatment for the terminal diagnosis and have a prognosis of 6 months or less if the illness follows its normal course. While hospice is most associated with cancer and organ failure, individuals with terminal psychiatric conditions — including severe dementia, treatment-refractory anorexia nervosa, and end-stage conditions complicated by psychiatric illness — may qualify and benefit from the BH components of the hospice interdisciplinary team.

**Statutory basis:** Medicare Part A hospice benefit, 42 U.S.C. § 1395d; **Conditions of Participation:** 42 CFR Part 418 (verified current at eCFR.gov, May 14, 2026).

## Hospice Eligibility with Psychiatric Diagnoses

### Terminal Illness Requirement (42 CFR 418.22)
- Prognosis of 6 months or less with normal disease course
- Certifying physicians: attending physician (if any) + hospice medical director must both certify
- Benefit periods: two 90-day periods, then unlimited 60-day periods (42 CFR 418.21)
- Face-to-face encounter by hospice physician or NP required before 3rd benefit period and each recertification thereafter (42 CFR 418.22(a)(4), effective Jan. 1, 2011)

### Psychiatric Conditions That May Meet Criteria

**Dementia with behavioral symptoms (most common)**
- Alzheimer's disease and related dementias with significant behavioral symptoms (agitation, aggression, psychosis) in late stages
- Medicare LCD criteria for dementia (L34567; varies by MAC): FAST scale ≥ stage 7 (no consistent meaningful verbal communication; ≤6 intelligible words) plus KPS/PPS < 70%, dependence on 2+ ADLs, and at least one life-threatening complication (aspiration pneumonia, pyelonephritis, Stage 3–4 pressure ulcer, or septicemia) within the prior 12 months [verified-as-of 2026-06-22; MAC contractor LCDs may vary — confirm under the patient's MAC]
- Behavioral symptoms (agitation, combativeness) managed through IDT including social work and behavioral consultation

**Treatment-Refractory Anorexia Nervosa**
- Establishing 6-month terminal prognosis is clinically and procedurally challenging; the hospice 6-month rule was not designed around primary psychiatric diagnoses
- Very low BMI, severe medical complications, exhaustion of evidence-based treatments, and patient declining further treatment are relevant clinical factors — but no Medicare-specific BMI cutoff exists in CMS coverage policy
- Regulatory barrier: many hospices are reluctant to certify a primary psychiatric diagnosis as terminal; some decline on clinical or ethical grounds [source: Hospice News, 2024; AMA Journal of Ethics, 2023]
- Palliative care (non-hospice) is an active alternative model for severe and enduring anorexia nervosa (SE-AN) that does not require the 6-month terminal prognosis and hospice election

**Schizophrenia / End-Stage Psychotic Disorders**
- Rarely meet 6-month terminal prognosis criteria based on psychiatric diagnosis alone
- May qualify via comorbid terminal medical conditions (cancer, heart failure) where psychiatric symptom burden significantly affects end-of-life care
- CMS coverage criteria reference physical functional decline indicators; there is no Medicare-specific hospice eligibility pathway for terminal psychiatric illness absent concurrent terminal medical condition

## Hospice Per-Diem and Covered Services

Medicare pays a prospective daily per-diem for one of four levels of care (42 CFR Part 418 Subpart G); the hospice is responsible for covering ALL services related to the terminal diagnosis and related conditions:

**Four levels of care** (42 CFR 418.202(e)):
1. Routine Home Care (RHC) — standard daily in-home care
2. Continuous Home Care (CHC) — crisis nursing care at home, minimum 8 hrs/day
3. General Inpatient Care (GIP) — short-term inpatient for symptom control unmanageable at home
4. Inpatient Respite Care (IRC) — brief inpatient to relieve family/caregivers

FY 2026 update: 2.6% increase over FY 2025 rates; aggregate cap $35,361.44 per beneficiary [CMS-1835-F, published Aug. 5, 2025, effective Oct. 1, 2025].

### Psychiatric Medications in Hospice (42 CFR 418.202(f))
- Drugs used for pain relief and symptom control are covered in the per-diem under 42 CFR 418.202(f)
- **For dementia**: antipsychotics for behavioral symptoms, antidepressants, benzodiazepines for anxiety — covered as comfort/palliation when related to the terminal diagnosis
- Medications for concurrent non-terminal conditions unrelated to the terminal illness are not hospice per-diem expenses; may be billed separately to Part D (GW modifier used to indicate "not related to terminal condition")
- **MAT note:** Part D opioid safety management interventions (including point-of-sale edits) do not apply to hospice/palliative/EOL patients; beneficiary access to MAT such as buprenorphine should not be impacted by hospice enrollment per CMS guidance, though specific concurrent OUD treatment + hospice payment mechanics are not addressed in a single definitive CMS policy document — consult MAC for billing guidance

## Hospice Interdisciplinary Team (IDT) — BH Components

### Required IDT members (42 CFR 418.56(a)(1), as amended by CAA 2023 §4121, eff. Jan. 1, 2024):
- **(i) Physician** (MD or DO, employee or under contract)
- **(ii) Registered nurse** (must coordinate care and ensure continuous assessment)
- **(iii) Social worker, marriage and family therapist, OR mental health counselor** (CAA 2023 expanded this slot from social worker only; the discipline filling this role must be a direct hospice employee because counseling is a core service under 42 CFR 418.64)
- **(iv) Pastoral or other counselor** (spiritual counseling; broader than "chaplain")

### Core services the hospice MUST directly provide (42 CFR 418.64):
- Nursing services
- Medical social services (qualified social worker, under physician direction; based on psychosocial assessment)
- Counseling services — including:
  - **Bereavement counseling**: organized program supervised by qualified professional with grief/loss education; available up to 1 year post-death (42 CFR 418.64(d)(1))
  - **Dietary counseling** (42 CFR 418.64(d)(2))
  - **Spiritual counseling**: hospice must provide assessment, provide counseling consistent with patient/family beliefs, and make reasonable efforts to facilitate clergy/pastoral visits (42 CFR 418.64(d)(3))

Optional/contracted services that may be relevant to BH:
- Psychologist, LCSW, MFT, MHC: individual psychotherapy for emotional distress in terminal illness — if provided as part of the hospice plan of care, included in the per-diem; NOT separately billable to Medicare Part B
- Music therapist, art therapist: non-pharmacologic behavioral symptom management in dementia (not listed as covered services in 42 CFR 418.202 but may be covered as "other reasonable services" under 42 CFR 418.202(i) when specified in the plan of care)

## Coverage of Psychotherapy in Hospice

- Psychotherapy provided by the hospice as part of the care plan is included in the per-diem — NOT separately billable to Medicare Part B (the hospice election waives Part B billing for services related to the terminal illness and related conditions)
- If a hospice patient receives outpatient psychotherapy for a *separately identifiable, non-terminal condition*, the psychotherapy may be billed to Part B with the GW modifier ("not related to terminal condition") — but this is an exception and requires documentation of unrelatedness
- **Practical reality**: Most BH support in hospice is delivered by the social work/MFT/MHC IDT member and through spiritual counseling; formal psychotherapy referrals within the per-diem are uncommon

## Hospice Election and the "Waiver" Requirement

When a beneficiary elects Medicare hospice (42 CFR 418.24):
- The beneficiary waives all rights to Medicare payments for items, services, and drugs related to the terminal illness and related conditions from any source other than the designated hospice and attending physician
- Services unrelated to the terminal illness remain covered under regular Medicare fee-for-service (Part A, Part B, Part D)
- **BH example**: Patient with terminal cancer elects hospice; antipsychotics prescribed for comfort/agitation management (related to terminal condition) → hospice per-diem; ongoing psychiatric outpatient treatment for a pre-existing, non-terminal schizophrenia condition → potentially still covered via Part B (GW modifier required)

## Tension: Hospice Election vs. BH Treatment

Key dilemma: some BH treatments are inherently disease-modifying and difficult to categorize as "curative" vs. "palliative":
- **MAT (buprenorphine, methadone) for OUD**: hospice may cover as comfort measure if related to terminal condition symptom management; Part D opioid safety edits exempt hospice patients; specific coverage mechanics for concurrent OUD treatment and hospice are not resolved in a single CMS policy — contact MAC
- **Antidepressants**: generally covered in hospice as comfort-oriented symptom management
- **ECT for depression**: typically considered a treatment-oriented (not purely palliative) intervention; hospices generally will not provide ECT under the per-diem for this reason; patient would need to revoke hospice election to receive ECT for depression separately — this reflects common hospice practice, not an explicit CMS prohibition
- **Assisted Outpatient Treatment (AOT) compliance**: patient on hospice who is under an AOT court order faces legal and ethical complexity not addressed in CMS guidance

## Concurrent Care

- **Adults — Medicare hospice**: election of hospice means forgoing curative treatment for the terminal diagnosis; no general Medicare concurrent curative care exception exists for adults
- **Adults — Medicare Care Choices Model (MCCM)**: CMS ran a demonstration (2016–2024) allowing Medicare FFS beneficiaries with certain terminal conditions to receive supportive/hospice-like services without giving up curative treatment; the model has since concluded
- **Children — Medicaid/CHIP**: ACA Section 2302 (Pub. L. 111-148, enacted March 23, 2010) amended SSA §§ 1905(o)(1) and 2110(a)(23) to allow children enrolled in Medicaid or CHIP to receive hospice services without forgoing curative treatment related to the terminal illness; this is a Medicaid/CHIP provision, not Medicare

## Nursing Facility Hospice and PASRR (Psychiatric Residents)

- When a NF resident with serious mental illness (SMI) or intellectual/developmental disability (I/DD) elects hospice:
  - Terminal illness (defined per 42 CFR 418.3) constitutes a basis for a **categorical determination** under 42 CFR 483.130(d)(2) that nursing facility-level services are appropriate — this streamlines the PASRR Level II evaluation process
  - This is not a blanket PASRR waiver; it allows State authorities to make a categorical eligibility determination based on terminal status rather than requiring a full individual evaluation
- 42 CFR 418.112 governs hospice conditions of participation for hospice care provided to residents of SNF/NF/ICF/IID

## References

- **42 CFR Part 418 — Hospice Care** (current through May 14, 2026): https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-B/part-418
- **42 CFR 418.56 — IDT composition** (incl. CAA 2023 MFT/MHC amendment): https://www.law.cornell.edu/cfr/text/42/418.56
- **42 CFR 418.64 — Core services** (nursing, social work, counseling): https://www.law.cornell.edu/cfr/text/42/418.64
- **42 CFR 418.202 — Covered hospice services**: https://www.law.cornell.edu/cfr/text/42/418.202
- **42 CFR 483.130 — PASRR categorical determinations** (incl. terminal illness §(d)(2)): https://www.law.cornell.edu/cfr/text/42/483.130
- **CMS Hospice payment page** (FY 2026 update, levels of care): https://www.cms.gov/medicare/payment/fee-for-service-providers/hospice
- **FY 2026 Hospice Final Rule (CMS-1835-F) Fact Sheet** (2.6% update, $35,361.44 cap, eff. Oct. 1, 2025): https://www.cms.gov/newsroom/fact-sheets/fy-2026-hospice-wage-index-and-payment-rate-update-hospice-quality-reporting-program
- **Federal Register FY 2026 Hospice Final Rule** (Aug. 5, 2025): https://www.federalregister.gov/documents/2025/08/05/2025-14782/medicare-program-fy-2026-hospice-wage-index-and-payment-rate-update-and-hospice-quality-reporting
- **CMS LCD — Hospice Alzheimer's Disease & Related Disorders (L34567)** (FAST criteria): https://www.cms.gov/medicare-coverage-database/view/lcd.aspx?LCDId=34567
- **ACA Section 2302 — Concurrent Care for Children** (Pub. L. 111-148, Mar. 23, 2010): https://www.congress.gov/111/plaws/publ148/PLAW-111publ148.pdf
- **CMS Guidance on Concurrent Care for Children (SMD #10-018)**: https://www.medicaid.gov/federal-policy-guidance/downloads/smd10018.pdf
- **CMS MFT/MHC FAQ — CAA 2023 §4121 hospice IDG amendment**: https://www.cms.gov/files/document/marriage-and-family-therapists-and-mental-health-counselors-faq.pdf
- **Hospice News — Terminal Anorexia Referral/Regulatory Barriers (May 2024)**: https://hospicenews.com/2024/05/31/terminal-anorexia-patients-hitting-referral-regulatory-barriers-to-hospice-care/
- **AMA Journal of Ethics — Palliative Care for SE-AN (Sept. 2023)**: https://journalofethics.ama-assn.org/article/life-affirming-palliative-care-model-severe-and-enduring-anorexia-nervosa/2023-09
