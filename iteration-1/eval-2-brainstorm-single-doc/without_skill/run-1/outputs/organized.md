# Trail Conditions App — Concept Doc

## The one-liner
**Waze for trails.** Recent, crowdsourced condition reports from people who are actually out on the trail right now.

## The problem
Hikers never really know the *current* state of a trail before they head out. The best existing option — AllTrails reviews — runs weeks behind. A comment from three weeks ago tells you nothing about whether the creek is running high today or whether that section is still buried in snow.

## The core idea
A lightweight app where hikers file quick, real-time condition reports as they go, and other hikers see fresh, trustworthy info before they set out. The value is *freshness* — knowing what a trail is like today, not last month.

### Report types
- Mud
- Snow / ice
- Downed trees
- Water crossings
- Crowding
- Trail closure

## What makes it different: freshness done right
The standout idea is **reports that expire by type.** A snow report from three weeks ago is just noise, but a downed tree stays relevant until someone clears it. So reports auto-decay on a schedule that fits the condition:
- **Mud** expires fast.
- **Downed trees** persist until cleared.
- (Snow, water crossings, etc. each get their own sensible lifespan.)

Nobody handles freshness well today. This expire-by-type mechanic is probably the real differentiator.

## Key features

**Offline-first (non-negotiable).**
There's usually no signal on the trail. The app has to queue reports locally and sync them automatically once the phone is back in range.

**Photo with every report.**
Each report includes a photo, with the timestamp and GPS pulled automatically from the EXIF data — so reports are anchored to a real time and place without extra effort.

**Confidence scoring / verification.**
Multiple reports that agree raise a confidence score. A single report is flagged as "unconfirmed." This keeps one-off or mistaken reports from misleading people.

## Open questions and things to decide

**Gamification?**
Points for filing reports could drive participation — but it feels a bit wrong for the outdoors crowd. Flagged as an open question, not a commitment.

**The cold-start problem.**
An empty app is useless; there's no reason to open it until there's data in it. One idea: seed it by scraping recent AllTrails / Reddit trail mentions. That's a legal gray area, so it needs more thought before relying on it.

**Scope creep to avoid: offline maps.**
Tempting, but out of scope. Better to link out to Gaia or CalTopo rather than build mapping ourselves.

## Business / revenue ideas
- **Premium tier** with a weather overlay.
- **Parks department partnerships.** They may actively *want* this data for trail-maintenance prioritization — sell them a dashboard.
- **Crowding data** as a side benefit: crowding reports could help smooth out trailhead parking chaos.

## Liability and safety
Real risk: someone trusts a "water crossing passable" report, attempts it, and drowns. Mitigation:
- Strong terms of service.
- "Conditions change" disclaimers throughout.
- **Never** use the word "safe" — report conditions, don't certify safety.

## Naming
Candidates: **TrailPulse**, **FreshTracks**, **TrailCheck**.
(FreshTracks is probably already taken by a ski-related product.)

## MVP scope
Keep the first version deliberately small:
- **iOS only**
- **3 report types:** mud, snow, blockage
- **One metro area** to start
- **No accounts** — device ID only

**Launch city: Seattle.** Chosen over Denver because it has more year-round variance in conditions, which gives the freshness angle more to work with.
