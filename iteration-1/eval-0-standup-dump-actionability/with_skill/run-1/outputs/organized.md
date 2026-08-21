# Standup / notes — week of 6/29 — organized

**Diagnosis:** Operational dump — a week of standup + pricing-meeting capture (Mon–Thu). It's actions, decisions, a couple of ideas and questions, plus incident and reference facts. You asked for "what I owe people" and "what got decided," so I sorted by **actionability** (Actions → Decisions → Ideas → Questions → Reference → Archive), not by topic. Nothing added, nothing dropped. Two things to know up front: (1) the Friday latency rollback was only **partial** — the same analyzer change also shipped to the EU cluster and was never rolled back there, so EU p99 is still degraded; (2) both questions raised this week already got answered (see Questions). No conflicts, nothing parked.

## Actions — what's outstanding

**Reliability / incident**
- **File the incident write-up** for Friday night's search-latency spike (p99 → 800ms). *(owner not stated; appeared twice in the notes)*
- **Carlos — file a ticket** for the EU cluster, whose p99 is still degraded. *(Thu)*

**Pricing page / launch**
- **Sam — ping marketing to sign off on the tier names** (builder / scale / enterprise). This is the current unblocker for the pricing page. *(Sam owns)*
- **Sam — update the pricing-page draft.** ⛔ Blocked on final tier names → marketing sign-off — and the team decided to hold the page until then (see Decisions). *(Sam has the draft)*

## Decisions — what got decided
- **API pricing → usage-based tier**, confirmed for the **September release**. The flat tier is killed. *(decided Mon, reconfirmed Wed)*
- **Hold the pricing page** until marketing signs off on the tier names. *(Thu)*

## Ideas
- **Customer-facing "trust page"** exposing the latency dashboard — and possibly uptime / status too. *(cf. status.io)*

## Questions — both resolved this week
- **Does the Sept release slip if usage metering isn't done?** → **No.** Metering is estimated at ~3 weeks (JT); the release is ~5 weeks out. *(resolved)*
- **Is the EU price-change notice period 60 days?** → **Yes — 60 days**, confirmed by Dana; it's in the DPA addendum. This also closes the "check with Dana" to-do. *(resolved)*

## Reference — facts worth keeping
- **Friday-night search-latency incident:** p99 spiked to 800ms Friday night, traced to an analyzer change (Kavya, Mon). Rolled back → p99 back to 120ms on the main cluster. **But** the same change had also shipped to the **EU cluster and was never rolled back there** (Carlos, Thu), so **EU p99 is still degraded** — the fix was partial. (Drives the two reliability actions above.)
- **Price-change notice periods:** 30 days standard (per legal); **60 days for EU customers**, confirmed by Dana, per the DPA addendum.
- **Proposed tier names:** builder / scale / enterprise — marketing has **not** signed off yet.
- **Metering estimate:** ~3 weeks (JT); September release is ~5 weeks out.

## Archive — stale / superseded
- Carlos was out until Thursday; he's now back. *(Mon → Thu)*
- Pricing meeting was moved to Wednesday — it happened (notes captured above). *(Mon)*

---
**Counts.** 21 note lines in (line 1 is the date header). The four compound lines — the Mon standup line, the Wed "30-day / 60-day / check-Dana" line, the Thu "Carlos back + EU spread" line, and the "EU degraded + ticket" line — split into multiples, giving **26 substantive items**. These map to **15 distinct entries** — 4 actions, 2 decisions, 1 idea, 2 (resolved) questions, 4 reference facts, 2 archived — with the other **11 items merged or folded in**: **3 duplicate/near-duplicate merges** (incident write-up ×2, usage-based-tier decision ×2, trust-page idea ×2); the Carlos out→back status merged into one archive line; the pricing-page blocker folded into its action; four incident fragments (spike, rollback, EU spread, EU still-degraded) folded into one incident record; and Dana's 60-day confirmation folded into the notice-period reference — which also answers the EU-60 question and closes the "check with Dana" to-do. The Sept-slip question is answered by the metering estimate in Reference. **0 flags (⚑ unclear / ⚡ conflict), 0 parked.** All 26 items accounted for.
