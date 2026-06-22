---
name: nat-jitai-digital-phenotyping
description: JITAI & digital phenotyping — passive sensing, ML-driven relapse triggers, moment-level intervention delivery. TRIGGER: just-in-time adaptive intervention, digital phenotyping addiction, passive sensing relapse prediction. SKIP: EMA survey design → ecological-momentary-assessment; mHealth app general → digital-therapeutics.
---

# Just-in-Time Adaptive Interventions and Digital Phenotyping

## Core Concepts
- **JITAI**: Intervention framework delivering support precisely when an individual is most vulnerable, adapting in real-time to changing internal states (craving, stress) and external contexts (location, social environment)
- **Digital phenotyping**: Quantitative characterization of individual behavior and psychophysiology via continuous passive smartphone/wearable sensor data — GPS, accelerometry, call/text logs, screen time, voice acoustics
- **Decision points**: Moments when an individual's state or context triggers intervention delivery; distinguished from time-contingent (fixed-schedule) delivery
- **Tailoring variables**: Observable signals (sensor features, EMA responses) used to personalize intervention type and timing; may be moderated by stable traits (impulsivity) or dynamic states (stress level)
- **Micro-randomized trials (MRT)**: Gold-standard JITAI study design randomizing individual decision points (not persons) to intervention vs. no-intervention to estimate proximal causal effects

## Mechanisms / Neuroscience
- Craving and relapse risk are state-dependent and fluctuate across minutes-to-hours; fixed-schedule interventions miss peak vulnerability windows
- GPS mobility patterns signal social isolation, high-risk location proximity (bars, dealers), and reduced goal-directed activity — each predicts craving onset
- Accelerometry distinguishes purposive movement from aimless locomotion; combined with GPS provides behavioral context for ML feature engineering
- Voice acoustic features (pitch variability, speaking rate, jitter/shimmer) index autonomic arousal and negative affect without self-report; sensitive to stress preceding lapse
- Ecological Momentary Assessment (EMA) surveys embedded in JITAI pipelines capture subjective craving/affect at decision points; passive sensing reduces EMA burden via automated state detection
- Machine learning classifiers (random forest, LSTM, gradient boosting) trained on sensor time-series achieve AUC 0.70–0.85 for predicting lapse in smoking and alcohol studies

## Clinical Picture
- Opioid use disorder: Epstein et al. (NIDA) combined passive GPS with ML to predict opioid craving/stress 90 minutes in the future among buprenorphine/methadone patients; positive predictive value 0.93
- Binge drinking: Smartphone sensor ML model predicted same-day binge events 1–6 hours prior with acceptable-to-good accuracy; call logs and social behavior were top features
- Smoking cessation: Sense2Stop MRT used wearable biosensors to optimize stress-management JITAI for relapse prevention; passive detection of internal stress states achieved <50% sensitivity — current bottleneck
- Alcohol relapse: ML models combining demographics, prior use, and psychological factors reached ~77% predictive accuracy; passive sensing alone underperforms self-report fusion
- Engagement attrition: Sensor-only (no prompts) systems show highest retention; EMA-heavy designs see 30–40% non-compliance within 2 weeks

## Evidence & Treatment
- 2025 systematic review (Frontiers Digital Health): 23 JITAI trials across mental health/addiction; majority published ≥2024; common delivery modalities are push notifications, brief CBT micro-sessions, and crisis helpline links
- 2024 Annual Reviews meta-analysis found JITAI effect sizes modest (d ≈ 0.20–0.35) for substance use outcomes vs. static digital interventions; larger effects when triggering algorithm outperforms random timing
- Actionable digital phenotyping framework (Mohr et al.) distinguishes passive monitoring → state classification → intervention selection pipeline; each stage introduces error that compounds
- Privacy-preserving ML (federated learning, on-device inference) is actively developed to address consent and data-minimization concerns; no large addiction trials yet use federated design
- Key unresolved question: burden-efficacy tradeoff — more sensors improve trigger accuracy but increase dropout; optimal sensor battery remains undetermined

## Key Facts
- Passive internal-state detection accuracy: <50% sensitivity in current systems (smoking/stress domain)
- Opioid craving prediction PPV: 0.93 (Epstein, NIDA; GPS + ML, 90-min horizon)
- Binge-drinking prediction window: 1–6 hours prior to event using smartphone sensors alone
- JITAI effect size vs. static mHealth: d ≈ 0.20–0.35 advantage when triggering algorithm is optimized
- EMA non-compliance onset: ~30–40% within 2 weeks in sensor-heavy protocols

## Related Concepts
- nat-real-time-fmri-neurofeedback
- nat-contingency-management-stimulant-sud
- add-recovery-capital-measurement
- add-prefrontal-cortex-impairment
- add-paws-sleep-disruption
