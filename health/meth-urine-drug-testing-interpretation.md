---
name: meth-urine-drug-testing-interpretation
description: Interpret UDS results in meth recovery: distinguish illicit meth from prescribed amphetamine/methylphenidate via immunoassay cross-reactivity and GC-MS confirmation
---

# UDS Interpretation — Meth Recovery + Prescribed Stimulants

## Immunoassay Screen (IA): What It Detects

Standard cutoff: **1000 ng/mL** (SAMHSA); some labs use 500 ng/mL.

The amphetamine IA detects the **phenethylamine scaffold** — amphetamine, methamphetamine, MDA, MDMA, and structurally similar compounds. Cross-reactants that can cause false positives include:

- **OTC sympathomimetics**: pseudoephedrine, ephedrine, phenylephrine
- **Psychiatric medications**: bupropion, trazodone, phenothiazines (chlorpromazine, promethazine, thioridazine), TCAs (desipramine, doxepin)
- **Selegiline**: metabolizes to l-amphetamine and l-methamphetamine — can produce true positive
- **Benzphetamine** (high-dose): metabolizes to amphetamine

**Methylphenidate does NOT cross-react.** It is a piperidine derivative, not a phenethylamine. A patient on methylphenidate who screens positive on the amphetamine IA requires explanation other than their prescription — either a cross-reactant, illicit use, or lab error.

---

## Confirmatory Testing (GC-MS / LC-MS/MS)

Order confirmatory testing for **all positives** when stimulants are prescribed. Confirmatory mass spec identifies specific analytes and quantities; immunoassay cannot.

### Methylphenidate (Ritalin, Concerta, Focalin)
- Confirmatory shows: **ritalinic acid** (primary metabolite) ± methylphenidate parent
- Shows: **no amphetamine, no methamphetamine**
- A positive IA + confirmatory showing only ritalinic acid = false-positive IA; document accordingly

### Amphetamine salts (Adderall, Adderall XR)
- Confirmatory shows: **d-amphetamine + l-amphetamine**
- Adderall is 75% d-amphetamine, 25% l-amphetamine by formulation; urine reflects this mix
- **No methamphetamine** expected — patients compliant with Adderall should not produce detectable methamphetamine

### Lisdexamfetamine (Vyvanse)
- Prodrug; activated in GI tract → releases **d-amphetamine only**
- Confirmatory shows: **d-amphetamine, no l-amphetamine, no methamphetamine**
- Distinguishable from Adderall by absence of l-amphetamine component

### Illicit methamphetamine
- Confirmatory shows: **d-methamphetamine AND d-amphetamine** (meth metabolizes ~15–30% to amphetamine)
- Presence of methamphetamine itself is the key — Adderall/Vyvanse do not produce methamphetamine
- If methamphetamine is present, source is not a prescribed amphetamine salt

---

## d- vs l-Methamphetamine: The Nasal Inhaler Problem

Standard GC-MS cannot distinguish d- from l-methamphetamine without a **chiral column**.

- **Illicit street meth**: d-methamphetamine (pharmacologically active)
- **Vicks VapoInhaler / L-Desoxyephedrine inhalers**: l-methamphetamine (pharmacologically weak; OTC legal)

When a patient claims nasal inhaler use as the source of a meth-positive, **chiral stereospecific GC-MS** is the only way to confirm or refute. Standard LC-MS/MS will confirm methamphetamine present but cannot assign chirality. Note: controlled studies show Vicks VapoInhaler does NOT produce detectable d-amphetamine, so d-amphetamine co-presence supports illicit meth.

---

## Specimen Validity — Rule Out Adulteration First

Before interpreting any result:

| Finding | Interpretation |
|---|---|
| Creatinine < 20 mg/dL | Dilute specimen |
| Specific gravity < 1.003 | Dilute |
| Nitrite > 500 µg/mL | Adulterated (oxidizer added) |
| pH < 3 or > 11 | Adulterated |
| Temperature out of range | Possible substitution |

Dilution can produce a false-negative. Adulteration invalidates the specimen. Document and re-collect under observation if validity is in question.

---

## Clinical Decision Framework

```
IA positive for amphetamines?
│
├─ Patient prescribed methylphenidate only?
│   └─ Run confirmatory: ritalinic acid only → false-positive IA; document
│      Confirmatory shows amphetamine/meth → illicit use or co-prescribed stimulant
│
├─ Patient prescribed Adderall or Vyvanse?
│   └─ Run confirmatory:
│      Amphetamine only (no meth) → consistent with prescription
│      Both meth + amphetamine → illicit meth use
│      Only d-amphetamine, no l → consistent with Vyvanse
│      d+l amphetamine mix → consistent with Adderall
│
└─ No stimulant prescribed?
    └─ Run confirmatory to rule out cross-reactant
       Meth + amphetamine confirmed → illicit meth use
       Amphetamine only → consider selegiline, OTC cross-reactants
       Negative on confirmatory → false-positive IA
```

---

## Documentation Protocol

1. **Baseline UDS** before initiating any stimulant in meth recovery
2. **Document expected positives** from prescription at each monitoring visit
3. **Order confirmatory (GC-MS/LC-MS)** on all positive screens when stimulants are prescribed — never act on IA alone
4. **Note cross-reactants** in the medical record when applicable
5. Monitoring frequency: weekly or biweekly early in treatment; monthly once stable (risk-stratify)
6. Chain-of-custody documentation if results will be used for legal/regulatory purposes

---

## References

- [Mayo Clinic Proceedings: Urine Drug Screening — Practical Guide for Clinicians](https://www.mayoclinicproceedings.org/article/S0025-6196(11)61120-8/fulltext)
- [PMC: A Practical Guide to Urine Drug Monitoring](https://pmc.ncbi.nlm.nih.gov/articles/PMC6368048/)
- [PMC: Urine methamphetamine-to-amphetamine ratio by LC-MS/MS](https://pmc.ncbi.nlm.nih.gov/articles/PMC12280402/)
- [PMC: Methamphetamine and Amphetamine Isomer Concentrations Following Vicks VapoInhaler Administration](https://pmc.ncbi.nlm.nih.gov/articles/PMC4168291/)
- [PMC: UDS Considerations for the Psychiatric Pharmacist](https://pmc.ncbi.nlm.nih.gov/articles/PMC6009242/)
- [US Pharmacist: Minimizing False-Positives and False-Negatives in UDS](https://www.uspharmacist.com/article/urine-drug-screening-minimizing-false-positives-and-false-negatives-to-optimize-patient-care)
