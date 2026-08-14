# Source Basis

This skill was built from the exported Okta markdown corpus plus the existing MongoDB and 10gen skills.

## Highest-signal source documents

- Training Manual and Playbook for MongoDB Incident Responders: Okta
- Okta Tier-0 Incident Responder Training Guide
- Okta Tier-0 IR AI Assistant Playbook
- Okta Incident Fire-Drill Plan and Playbook (MongoDB × Okta)
- Okta Incident Fire-Drill Process Runbook
- Okta Fire-Drill — IR Card (On-Shift Responder)
- Okta IR Operational Checklist
- Okta IR Consolidated Handbook
- Okta AMER IR Role Brief
- Enhanced Support S1 Outage Playbook – Okta
- Okta IR Readiness Gap Assessment
- Case Analysis Report: Okta 5/19/26

## Existing skill inputs folded into routing

- `atlas-diagnostics-expert`
- `mongodb-expert`
- `mongodb-atlas-expert`
- `mongodb-operations-expert`
- `mongodb-kb`
- `10gen`

## Extracted method

The repeatable method extracted from the corpus is:

1. Treat credible Okta S1/S2 signals as incident command work, not normal support handling.
2. Establish command and a scenario hypothesis inside the first 10–15 minutes.
3. Separate current symptom and customer impact from root-cause explanation.
4. Route deep technical work through the correct MongoDB/Atlas/10gen skill rather than improvising across all domains at once.
5. Keep a clear distinction between live-incident action, fire-drill simulation, and readiness work.
6. Escalate early when restoration may depend on Atlas or cloud-provider intervention.

## Boundaries

This method is tuned for Okta Tier-0 Atlas incidents, drills, and high-severity case reviews. It is not a general replacement for normal low-severity support handling.

## Validation status

Status: PROVISIONAL

The extracted workflow was cross-checked against independent source types inside the corpus:
- role brief
- responder card
- fire-drill plan
- consolidated handbook
- outage scenario index

Those sources reinforced the same core method: early command, scenario-led routing, restoration-first posture, and explicit drill safety boundaries. It has not yet been independently validated on a fresh live incident outside this corpus.