---
name: enrollment-eligibility-medicaid-pregnancy-postpartum
description: Medicaid eligibility during pregnancy and the postpartum period — pregnancy as a Medicaid eligibility category, income limits for pregnant women (typically 185–215% FPL), CHIP unborn child option, 12-month postpartum coverage (ARPA state option now adopted by 49 states as of Feb 2026), presumptive eligibility for pregnant women, coverage during labor and delivery, and interaction with other Medicaid eligibility categories after pregnancy. TRIGGER: Medicaid pregnancy; pregnancy Medicaid eligibility; Medicaid while pregnant; Medicaid postpartum; 12-month postpartum Medicaid; Medicaid maternity coverage; CHIP pregnancy; pregnant Medicaid income limit; Medicaid prenatal coverage; Medicaid labor delivery; postpartum coverage extension; Medicaid after baby born; Medicaid for pregnant women income; NC Medicaid pregnancy. SKIP: presumptive eligibility broadly → enrollment-eligibility-presumptive-eligibility; CHIP enrollment → enrollment-eligibility-chip-enrollment; general Medicaid eligibility determination → enrollment-eligibility-medicaid-eligibility-determination; retroactive Medicaid → enrollment-eligibility-retroactive-medicaid.
version: "1.0.0"
updated: "2026-06-22"
category: health/enrollment-eligibility
tags: [medicaid, pregnancy, postpartum, maternity, CHIP, ARPA, presumptive-eligibility, 12-month-postpartum]
---

# Medicaid Eligibility During Pregnancy and Postpartum

## Pregnancy as a Medicaid Eligibility Category

Pregnancy is an independent Medicaid eligibility category — a person can qualify for Medicaid based on pregnancy alone even if they would not otherwise qualify. Key features:

- **Higher income limits**: Federal minimum is 133% FPL for pregnant women; most states go significantly higher (many at 185%–215% FPL)
- **No asset test** under MAGI methodology (pregnancy is a MAGI category)
- **Broader coverage**: All pregnancy-related and related medical services required
- Coverage begins **from the date of application** (or retroactively up to 3 months prior if eligible)

## Income Limits by State (Common Examples, 2025)

| State | Pregnant Women Medicaid Income Limit |
|---|---|
| North Carolina | 196% FPL |
| Texas | 198% FPL |
| California | 213% FPL |
| New York | 223% FPL |
| Federal minimum | 133% FPL |

## CHIP Unborn Child Option

States may elect to cover the unborn child through CHIP (rather than the pregnant individual), allowing coverage for prenatal care under CHIP at higher income levels than Medicaid. This is a separate election from covering the pregnant person.

## Presumptive Eligibility for Pregnant Women

Qualified entities (hospitals, FQHCs, OB/GYN practices, Navigators) can make same-day temporary Medicaid eligibility determinations for pregnant women while the full application is processed. Coverage begins immediately, allowing prenatal appointments without delay. See `enrollment-eligibility-presumptive-eligibility` for full mechanics.

## Coverage During Pregnancy

Medicaid for pregnant women covers:
- Prenatal care (all visits, labs, ultrasounds)
- Labor and delivery (inpatient hospital)
- Postpartum follow-up
- Treatment for conditions that may complicate the pregnancy
- Dental and vision in many states (optional benefits)

## The Postpartum Extension: From 60 Days to 12 Months

**Original rule:** Medicaid coverage for pregnant women ended 60 days after delivery.

**ARPA 2021 state option:** American Rescue Plan Act gave states the option to extend postpartum coverage to 12 months.

**Consolidated Appropriations Act 2023:** Made the 12-month postpartum option **permanent** (states can elect it without time limit).

**Current adoption (as of February 2026):**
- **49 states + DC** have adopted 12-month postpartum coverage (Wisconsin passed it February 2026)
- Only **Arkansas** had not adopted as of early 2026

## What 12-Month Postpartum Means

- Full Medicaid coverage continues for **12 months** after the end of pregnancy (delivery, miscarriage, or termination)
- Coverage is not limited to pregnancy-related services — it covers all medically necessary care
- Critical for maternal mortality prevention: ~50% of pregnancy-related deaths occur in the postpartum year
- Applies to all individuals who were eligible for pregnancy Medicaid, regardless of income changes after delivery

## Transition After 12 Months Postpartum

At 12 months postpartum:
1. State redetermines eligibility under the applicable coverage group
2. If income is at or below 138% FPL (in expansion states like NC): likely eligible for adult expansion Medicaid
3. If above expansion income limits: may qualify for ACA marketplace coverage with APTCs
4. Children born during the pregnancy are eligible for Medicaid/CHIP in their own right

## NC-Specific Notes

- NC Medicaid covers pregnant women at up to **196% FPL**
- NC adopted 12-month postpartum coverage
- Apply through NC FAST (epass.nc.gov) or in-person at local DSS
- NC provides presumptive eligibility for pregnant women through hospitals and FQHCs
- NC Tailored Plans cover pregnancy services; Standard Plans cover most pregnant enrollees

## Implementation Concerns (2025–2026)

The One Big Beautiful Bill (H.R.1, 2025) creates new administrative burdens that could affect postpartum enrollees:
- Semi-annual redeterminations for expansion adults could disrupt postpartum coverage transitions
- New documentation requirements may create gaps for recently-postpartum individuals moving to the expansion group
- States must build systems to identify and protect the postpartum period from inappropriate redeterminations
