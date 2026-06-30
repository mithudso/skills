---
name: medicaid-nc-dma-administration
description: >-
  NC Division of Medical Assistance (DMA) / Division of Health Benefits (DHB) — the state Medicaid agency; organizational structure, responsibilities, NCDHHS hierarchy, NC Medicaid policy development, Clinical Coverage Policies (CCPs), provider enrollment and appeals, NC Medicaid Direct (fee-for-service), NC Medicaid recipient rights, NC HealthConnex HIE mandate. TRIGGER: NC DMA, NC Division of Medical Assistance, NC Division of Health Benefits, NC Medicaid agency, NC Medicaid Director, NC Medicaid fee-for-service, NC Medicaid Direct, NC Medicaid Clinical Coverage Policy, NC Medicaid provider enrollment, NC Medicaid appeals, NC HealthConnex HIE, NC health information exchange. SKIP: NC Medicaid Managed Care → nc-health-services-medicaid-managed-care; NC LME-MCOs → nc-health-services-lme-mcos; NC HCBS waivers → nc-health-services-medicaid-waivers; NC Medicaid behavioral health → nc-drug-rehabilitation-treatment.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [medicaid, NC-DMA, NC-DHB, NCDHHS, fee-for-service, NC-HealthConnex, HIE, administration, North-Carolina]
related_skills:
  - nc-health-services-medicaid-managed-care
  - nc-health-services-lme-mcos
  - nc-health-services-medicaid-waivers
  - medicaid-fmap-federal-financing
  - medicaid-fraud-waste-abuse
---

# NC Division of Medical Assistance / Division of Health Benefits

> **Educational only — NOT legal or medical advice.** Administrative structures and policies change; verify current information at medicaid.ncdhhs.gov.

**verified-as-of: 2026-06-22**

---

## 1. Organizational Structure

### NC DHHS Hierarchy
**NC Department of Health and Human Services (DHHS)** is the state-level umbrella agency. Within DHHS, Medicaid is administered by the **Division of Health Benefits (DHB)**, formerly known as the **Division of Medical Assistance (DMA)**. Both terms are used; official current name is DHB.

```
NC DHHS Secretary
└── Division of Health Benefits (DHB)
    ├── NC Medicaid Managed Care (Standard Plans, Tailored Plans, CFSP)
    ├── NC Medicaid Direct (fee-for-service, exempt populations)
    ├── Policy & Clinical Coverage (CCPs, benefits)
    ├── Provider Enrollment & Credentialing
    ├── Program Integrity
    ├── Member Services & Recipient Rights
    └── NC HealthConnex (HIE)
```

### NC Medicaid Director
The DHB Director serves as NC's Medicaid Director, responsible for CMS compliance, state plan amendments (SPAs), waiver applications, and managed care contract oversight.

---

## 2. NC Medicaid Direct (Fee-for-Service)

While most NC Medicaid enrollees are in managed care (Standard Plans or Tailored Plans), **NC Medicaid Direct** (fee-for-service/FFS) continues for exempt populations:

**Populations remaining in NC Medicaid Direct (as of 2026)**:
- Individuals enrolled in **PACE** programs
- **Foster care youth** (ages 0–25 — NC adopted extended foster care)
- **American Indian/Alaska Native** Medicaid recipients (exempt from mandatory managed care)
- Members of **federally recognized tribes**
- Individuals with **spend-down** status
- Some **dual eligibles** not enrolled in managed care
- **Emergency Medicaid** (non-citizen emergency services only)
- Certain waiver populations receiving LTSS not yet transitioned to managed care

NC Medicaid Direct reimburses providers directly on a fee-for-service basis using NC Medicaid fee schedules.

---

## 3. Clinical Coverage Policies (CCPs)

**Clinical Coverage Policies (CCPs)** are NC Medicaid's equivalent of coverage determinations — they define what services are covered, what medical necessity criteria apply, and what prior authorization requirements exist for each service category.

CCPs are developed by DHB's Clinical Policy team, informed by:
- CMS Medicaid statute and regulations
- Evidence-based medicine reviews
- MHPAEA parity compliance requirements
- 1115 waiver terms
- Input from provider community and stakeholders

**How providers use CCPs**:
- CCPs are publicly available at medicaid.ncdhhs.gov
- Prior authorization requests must document medical necessity per the CCP
- Appeal of denied services must reference the applicable CCP criteria

**2025–2026 CCP developments**: Major SUD service CCPs were updated effective January 1, 2026, aligned with the 1115 SUD Demonstration Waiver renewal (see `nc-behavioral-health-medicaid.md`). MHPAEA parity-compliant CCPs for MH/SUD launched January 1, 2025.

---

## 4. Provider Enrollment

Medicaid providers in NC must enroll with DHB. The enrollment process:
1. Apply via NC Tracks (NC Medicaid's claims and provider management system)
2. Background screening (OIG LEIE, SAM, NC criminal background — varies by provider type)
3. Site visits (required for home health, personal care, certain behavioral health providers)
4. Credentialing (licensed providers must have valid NC license)
5. Re-enrollment every 5 years (or when ownership/control changes)

**High-risk providers** (home health, DME, behavioral health) face enhanced screening. Providers excluded from Medicare/federal programs are ineligible for Medicaid enrollment.

---

## 5. State Plan Amendments (SPAs) and Waivers

NC Medicaid's basic coverage rules are set in the **NC Medicaid State Plan** — the formal agreement between NC and CMS. Changes require:
- **State Plan Amendment (SPA)**: For most changes to covered services, eligibility, or provider rates — requires CMS approval
- **1915(c) waiver**: For HCBS waiver programs (CAP/DA, CAP/C, Innovations)
- **1115 demonstration waiver**: For major demonstration programs (SUD waiver, Healthy Opportunities Pilots, managed care itself)
- **1915(b) waiver**: Used to mandate managed care enrollment for certain populations

---

## 6. Member Services and Recipient Rights

**NC Medicaid member rights** include:
- Right to receive covered services without discrimination
- Right to receive information in an understandable format (language access, ADA)
- Right to appeal coverage denials
- Right to disenroll from a managed care plan with cause

**Appeals process**:
1. Level 1: Internal appeal to managed care plan (Standard Plan, Tailored Plan, LME-MCO) within 60 days of adverse action
2. Level 2: State Fair Hearing through NC OAH (Office of Administrative Hearings) — must be requested within 90 days
3. Federal standard: Continuation of benefits during appeal (if requested timely)

---

## 7. NC HealthConnex — Statewide Health Information Exchange

**NC HealthConnex** (formerly NC HIE) is North Carolina's statewide **Health Information Exchange (HIE)**, operated by the NC Health Information Exchange Authority (HIEA), a state agency within NCDHHS.

### Purpose
Enables real-time, secure sharing of patient health records across providers, payers, and public health agencies to support care coordination and reduce duplicative testing.

### Mandatory Participation (GS § 90-414.4)
**NC law mandates** participation in NC HealthConnex for:
- **Medicaid-enrolled providers** in NC must connect to NC HealthConnex and submit/query data as a condition of Medicaid participation
- Facilities licensed under GS 131E and 122C
- Specific implementation deadlines by provider type (hospitals first, then practices)

This makes NC unusual nationally — it has a statutory mandate for Medicaid provider HIE participation.

### What NC HealthConnex Contains
- Clinical documents (discharge summaries, visit notes, lab results, imaging)
- Medication history
- Immunization records
- Patient demographic data
- Event notifications (ADT — Admission, Discharge, Transfer alerts)

### Relevance to NC Medicaid
- DHB uses NC HealthConnex data for care coordination (especially for high-utilizers)
- Tailored Plan MCOs and Standard Plan MCOs are required to utilize HIE data for care management
- NC HealthConnex generates real-time alerts to care managers when Medicaid patients are admitted to hospitals or ERs

### 2025 Updates
- NC HealthConnex expanded API access for providers using EHR systems
- New behavioral health data elements added for SUD treatment coordination
- NCCARE360 social determinants platform integrated with HIE event notifications

---

## 8. NC Medicaid Data and Analytics

DHB publishes Medicaid data through:
- **NC Medicaid Data Dashboard**: Enrollment, expenditure, managed care quality metrics
- **NCTracks**: The state claims processing system (providers submit claims; public data available for enrollment/provider lookup)
- **MACPAC and KFF**: State Medicaid profiles including NC

---

## 9. Anti-Patterns

**"DHB and DMA are different agencies"** — Same agency; DMA was renamed DHB. Legacy documents use DMA; current official name is DHB.

**"NC Medicaid fee schedules are the same as Medicare"** — No. NC Medicaid rates are set by NC, subject to CMS adequacy standards and federal upper payment limits. NC Medicaid rates are generally lower than Medicare for many service categories.

**"Providers don't need to connect to NC HealthConnex if they're small"** — NC statute mandates participation for Medicaid-enrolled providers; size does not exempt. However, implementation timelines varied by provider type.

---

## 10. References

- NCDHHS Division of Health Benefits: [medicaid.ncdhhs.gov](https://medicaid.ncdhhs.gov)
- NC HealthConnex / HIEA: [hiea.nc.gov](https://hiea.nc.gov)
- GS § 90-414.4 — NC HealthConnex mandatory participation
- NCTracks (provider enrollment portal): [nctracks.nc.gov](https://nctracks.nc.gov)
- NC Medicaid State Plan (current version): [medicaid.ncdhhs.gov/about-nc-medicaid/nc-state-plan](https://medicaid.ncdhhs.gov/about-nc-medicaid/nc-state-plan)
- MACPAC, NC Medicaid Profile: [macpac.gov/state/north-carolina](https://www.macpac.gov/state/north-carolina/)
- KFF, NC Medicaid facts: [kff.org/medicaid/state-indicator](https://kff.org/medicaid/state-indicator/)
- NC OAH (State Fair Hearings): [oah.nc.gov](https://www.oah.nc.gov)
