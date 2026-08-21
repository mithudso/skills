---
name: add-warm-handoff-care-transitions
description: Warm handoff and care transitions for SUD treatment — ED to outpatient, detox to long-term care, peer specialist roles, meth-specific challenges, Bridge Clinic model, reducing dropout at transition points.
---

# Warm Handoff & Care Transitions — Meth/SUD

## Definition & Why It Matters

Warm handoff: direct, real-time provider-to-provider (or provider-to-peer) connection at the moment of discharge — patient is introduced, not sent. Contrast with passive referral (give patient a phone number): engagement rate <20%. Warm handoff achieves 60–80% engagement.

Every care transition is a dropout event waiting to happen. Crisis → treatment, ED → outpatient, detox → long-term care = highest-risk windows for meth SUD patients where treatment retention is already low (~40% at 30 days without active support).

---

## ED → Outpatient: The Bridge Clinic Model

**Core components:**
- Same-day or next-day SUD evaluation appointment scheduled *before* ED discharge
- Peer Recovery Specialist (PRS) present at discharge to escort/accompany to first visit
- Bridge Clinic: rapid-access outpatient site (UCSF model) accepting walk-ins, no appointment barriers, harm-reduction framing

**Evidence:**
- Brown et al. 2021 (*Annals Emergency Medicine*): Bridge Clinic patients 4.3× more likely to engage SUD treatment at 30 days vs. standard referral
- Kinnard et al. 2021: ED-initiated buprenorphine + Bridge Clinic = >60% 30-day treatment engagement
- SUN intervention (Philadelphia, 3 hospitals, 2023–2024): 50.4% treatment engagement at 30 days vs. 15.9% controls — aOR 3.7 (ScienceDirect, 2022 implementation study)

Bridge clinics provide MOUD initiation, harm reduction services, stabilization during high-risk transitions, and linkage to long-term providers. Patients report value in welcoming environments, no-appointment access, and compassionate staff.

---

## Detox → Long-Term Treatment

Without active follow-up, 30-day treatment engagement post-detox drops to <30%.

**Highest-ROI interventions:**
1. **PRS contact within 24–48 hours** post-discharge — single most effective lever
2. **Patient navigation** — help with insurance enrollment, transportation, childcare removes administrative barriers that silently cause dropout
3. **Discharge coordinator protocol**: contacts receiving program *before* patient leaves, confirms bed/appointment, arranges transportation; patient never holds a number alone

---

## MAT Continuity Across Settings

Primarily OUD-relevant but applies when meth patients have co-occurring OUD (common):
- Buprenorphine/methadone gaps between hospital/crisis and outpatient prescriber = rapid relapse + overdose death
- ED bridge buprenorphine prescriptions (now standard at many EDs) require explicit warm handoff to SUD prescriber — the bridge is incomplete without the handoff
- Methadone pathways (ED → community OTP): require peer support + navigation + bridge dosing to close gaps

---

## Peer Recovery Specialists as Handoff Agents

PRS = trained person with lived SUD experience. Most effective mediator at every transition point.

**Functions:**
- Accompaniment to first appointment (removes the "I'll go tomorrow" dropout point)
- Text/call check-ins between appointments
- Motivational support, stigma reduction
- Administrative navigation (insurance, childcare, transport)
- Linkage to housing, legal, social services

**Evidence:** PRS services reduce substance re-use, improve treatment retention and patient satisfaction, reduce costly acute care utilization (NCBI Bookshelf, Ch. 4; CHCS 2024).

---

## Meth-Specific Challenges

| Challenge | Implication for Handoff |
|---|---|
| No FDA-approved pharmacotherapy | No medication to bridge; behavioral engagement is the only lever |
| Low baseline retention (~40% at 30 days) | PRS accompaniment + contingency management at first visit critical |
| Meth psychosis during transition | Patient may be too disorganized to navigate — delay handoff attempt until acute psychosis resolves; family/support person can hold the connection |
| Rural distribution of meth use | Telehealth for follow-up visits reduces transportation barrier |

**Contingency management (CM):** Start CM at the first visit — tangible reinforcement for attendance creates immediate reason to return. Especially critical for meth SUD given absence of MAT.

---

## Technology-Supported Transitions

- SMS/text check-ins increase 30-day retention
- Automated appointment reminders reduce no-shows
- Telehealth removes rural transportation barrier
- JITAI (Just-In-Time Adaptive Interventions) — see `nat-jitai-digital-phenotyping` skill for meth-specific behavioral sensing + intervention timing

---

## Protocol Summary (Practitioner Reference)

```
Transition point identified →
  1. Confirm receiving provider/program before patient leaves
  2. Schedule appointment (same-day or <48h) before discharge
  3. Assign PRS for accompaniment
  4. Patient navigation: insurance, transport, childcare
  5. PRS check-in call 24–48h post-discharge
  6. Telehealth fallback if in-person not feasible
  7. CM incentive at first attended appointment
```

---

## Key Sources

- Brown et al. 2021, *Annals Emergency Medicine* — Bridge Clinic 4.3× engagement
- Kinnard et al. 2021 — ED buprenorphine + Bridge >60% 30-day engagement
- SUN intervention implementation study, ScienceDirect 2022 — aOR 3.7 for navigation vs. control
- NCBI Bookshelf NBK596261/NBK596269 — Peer Support in SUD Treatment (SAMHSA)
- PMC10101823 — Bridge clinic models, evidence, future directions (2023)
- PMC11043281 — Post-hospitalization care transition strategies, narrative review (2024)
- CHCS 2024 — Peer Recovery Support Services in SUD treatment
