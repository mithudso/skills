---
name: nc-health-services-nccare360
description: NCCARE360 — North Carolina's statewide closed-loop social needs referral network; first statewide coordinated care network in the US; available all 100 NC counties; connects health providers and community organizations electronically to address social determinants; powers Healthy Opportunities Pilots; NCDHHS-FHLI public-private partnership; Unite Us technology platform; 2019 launch, ongoing. TRIGGER: NCCARE360, NC social needs referral, NC closed-loop referral, NC statewide care coordination network, NC social determinants referral system, NC community resource referral, NC health and human services coordination, Unite Us North Carolina, NC FHLI Foundation for Health Leadership Innovation, NC SDOH referral platform, NC non-medical needs referral, NCCARE360 Healthy Opportunities, NC 211 care coordination. SKIP: NC Healthy Opportunities Pilots clinical program → nc-health-services-healthy-opportunities-pilots; NC 988 crisis line → nc-health-services-988-crisis; NC Medicaid behavioral health → nc-health-services-behavioral-health-sud.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [nc-health-services, sdoh, care-coordination, social-needs, referral-network]
keywords:
  - NCCARE360
  - closed-loop referral
  - NC social needs referral
  - NC care coordination network
  - Unite Us North Carolina
  - social determinants NC
  - NC community resource network
  - NC FHLI
  - NC statewide referral platform
  - NC non-medical needs
whenToUse:
  - "NCCARE360 how it works"
  - "NC statewide social needs referral platform"
  - "closed-loop referral North Carolina"
  - "NC care coordination for social determinants"
  - "connecting patients to community resources in NC"
whenNotToUse:
  - "NC Healthy Opportunities Pilots program mechanics → nc-health-services-healthy-opportunities-pilots"
  - "NC 988 crisis referrals → nc-health-services-988-crisis"
  - "NC FQHCs as referral destinations → nc-health-services-fqhcs"
---

# NCCARE360: NC's Statewide Social Needs Referral Network

## Overview

**NCCARE360** is the first statewide coordinated care network in the United States to electronically connect individuals with identified social needs to community resources — and critically, to **track whether those needs were actually met** (closing the loop). It launched in 2019 through a public-private partnership between NCDHHS and the Foundation for Health Leadership and Innovation (FHLI) and is available in all 100 NC counties as of June 2020.

NCCARE360 addresses a fundamental gap in health and human services: providers could refer patients to social services, but had no mechanism to know if the referral was acted on or if the need was resolved.

## How It Works

The platform enables **electronic, trackable referrals** across organizational boundaries:

1. A provider (hospital, clinic, FQHC, health plan, school) identifies a patient's social need (food insecurity, housing instability, transportation, utility assistance, etc.)
2. The provider submits an electronic referral through NCCARE360 to a community-based organization (CBO) that can address the need
3. The CBO receives the referral, contacts the patient, and updates the referral status
4. The provider is notified of the outcome — the "closed loop"
5. Data flows into a statewide repository, enabling population-level tracking of social needs and service gaps

### Key Technical Features

- **Shared technology platform** via Unite Us (the technology partner) — integrates resource directory, referral workflow, and outcome tracking
- **NC 211 integration** — United Way of NC's 211 information/referral line feeds into and from NCCARE360
- **Secure electronic communication** — HIPAA-compliant referral transmission between organizations
- **Self-serve resource navigation** — a public-facing resource directory tool launched October 2025 allows individuals to search for resources directly
- **HIE integration** — NCCARE360 social determinants data integrates with NC HealthConnex (the state HIE) for event-based notifications (e.g., hospital discharge triggers a housing referral)
- **Standardized data format** — uses open standards (HSDS — Human Services Data Specification) for resource interoperability

## Organizational Partners

| Partner | Role |
|---|---|
| NC DHHS | State government lead; policy, funding, oversight |
| Foundation for Health Leadership & Innovation (FHLI) | Network operator and convener |
| United Way of NC / 211 | Resource directory; call center; navigator support |
| Unite Us | Technology platform provider |
| Expound Decision Systems | Data verification and quality |
| Health systems, FQHCs, CBOs | Network participants (senders and receivers of referrals) |
| K-12 school districts | Pilot participants (school-based SDOH identification) |

## Social Needs Addressed

NCCARE360 addresses the core social determinants of health:

- **Food security** — food pantries, SNAP enrollment, meal programs
- **Housing stability** — emergency shelter, rental assistance, utility programs
- **Transportation** — medical transportation, ride programs
- **Employment and income** — workforce development, benefits enrollment
- **Safety** — domestic violence, safety planning resources
- **Social isolation and stress** — peer support, mental health navigation

## Relationship to Healthy Opportunities Pilots (HOP)

NCCARE360 is the **technology backbone** of NC's Healthy Opportunities Pilots (HOP), the first Medicaid Section 1115 waiver to test payment for non-medical social services. HOP providers use NCCARE360 to:
- Document identified social needs of Medicaid beneficiaries
- Refer to HOP-approved CBOs for food, housing, transportation, and interpersonal violence services
- Track service delivery and outcomes for HOP billing purposes

**Important note**: HOP's resource navigation services were **paused as of January 22, 2026** due to the suspension of the HOP program. NCCARE360 itself continues operating beyond HOP; it is a broader infrastructure program not tied solely to HOP.

## Current Status (2025-2026)

- **All 100 counties activated**: Full statewide coverage since June 2020
- **Self-serve tool launched**: October 2025 — public-facing resource navigation
- **Resource navigation paused**: January 22, 2026 (HOP suspension impacts)
- Platform continues operating for general closed-loop referral use outside HOP

## Research and Evidence

A 2024 study published in the NC Medical Journal (PubMed PMID 39412323) examined NCCARE360 implementation during the COVID-19 Support Services Program and found it effectively facilitated closed-loop coordination between health and social care sectors. The platform has been cited as a nationally recognized model for statewide resource navigation.

Duke Department of Psychiatry and Behavioral Sciences has also documented NCCARE360's role in helping North Carolinians address non-medical social needs.

## Implications for NC Medicaid

Under NC Medicaid managed care, Standard Plans and Tailored Plans are expected to address health-related social needs (HRSN) of their members. NCCARE360 provides the referral infrastructure to:
- Document SDOH screening results
- Route referrals to community resources
- Support ILOS (In-Lieu-of-Service) documentation for social needs interventions
- Satisfy reporting requirements under managed care contracts

## References

[^1]: NCCARE360. "Building Connections for a Healthier NC." [https://nccare360.org/](https://nccare360.org/) verified-as-of: 2026-06-22

[^2]: NCDHHS. "NCCARE360." [https://www.ncdhhs.gov/about/department-initiatives/healthy-opportunities/nccare360](https://www.ncdhhs.gov/about/department-initiatives/healthy-opportunities/nccare360)

[^3]: NC Medical Journal. "Implementation of NCCARE360, a Digital Statewide Closed-Loop Referral Platform." 2024. PubMed: [https://pubmed.ncbi.nlm.nih.gov/39412323/](https://pubmed.ncbi.nlm.nih.gov/39412323/)

[^4]: Unite Us. "NCCARE360 – North Carolina Closed-Loop Referral Network." [https://uniteus.com/networks/north-carolina/](https://uniteus.com/networks/north-carolina/)

[^5]: Open Referral. "Introducing NCCARE360: a coordinated statewide resource referral platform." [https://openreferral.org/introducing-nccare360-a-coordinated-statewide-resource-referral-platform/](https://openreferral.org/introducing-nccare360-a-coordinated-statewide-resource-referral-platform/)
