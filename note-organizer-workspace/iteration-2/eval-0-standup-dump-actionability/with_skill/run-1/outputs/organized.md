# Week of 6/29 — organized

**Diagnosis:** Operational dump — standup and meeting captures full of actions,
decisions, blockers, and one running incident. Your stated goal ("what I owe
people and what got decided") confirms the shape: actionability sort. Assumed
the un-owned incident-writeup TODO is yours, since the notes assign it to no
one — flagged rather than assigned.

## Actions

**Latency incident**
- File the incident writeup for Friday night's search-latency spike — no owner
  stated in the notes *(dedup: captured Mon and again Weds)*
- EU cluster ticket for the still-degraded p99 — Carlos said Thursday he's
  filing it (worth confirming it landed)

**Pricing**
- Ping marketing re: tier-name sign-off — owner: Sam *(Thurs)*
- Update the pricing page draft — Sam has it; currently on hold per the
  Decisions section until marketing signs off tier names *(blocked — see
  Problems)*

## Problems

- **EU p99 is still degraded.** The analyzer change shipped to the EU cluster
  too and was never rolled back there *(Carlos, Thurs)*; Carlos is filing a
  ticket — addressed by the EU-ticket action above.
  ⚡ conflict: Monday's notes say the analyzer change was "rolled back" and
  that the rollback fixed p99 — Thursday's info shows that was true of the
  main cluster only.
- **Pricing page draft is blocked on final tier names** *(Sam, Weds)* —
  marketing hasn't signed off; addressed by the ping-marketing action and the
  hold decision below.

## Decisions

- **Usage-based tier for the API; flat tier killed** *(Mon)*. Confirmed Weds:
  it goes in the September release.
- **Hold the pricing page until marketing signs off the tier names** *(Thurs)*.

## Wins

- Rollback fixed p99 on the main cluster: 800ms → back to 120ms *(Kavya, Mon)*
  — main cluster only; EU still degraded, see Problems.
- EU price-change notice question resolved: Dana confirmed 60 days, per the
  DPA addendum *(Thurs)* — closes Wednesday's "check with Dana".
- September-release timing risk resolved: metering estimate is 3 weeks (JT)
  and the release is 5 weeks out, so the date holds — closes the open question
  below.

## Ideas

- Customer-facing trust page: expose the latency dashboard to customers; could
  also show uptime — status.io does this *(merged from two notes, Mon + Weds)*.

## Questions

- Does the Sept release date slip if usage metering isn't done? — **answered
  in the notes:** metering is estimated at 3 weeks (JT) vs. 5 weeks to
  release, so no slip expected.

## Reference

- Price-change notice periods: 30 days standard (legal), 60 days for EU
  customers (Dana confirmed; source: DPA addendum).
- Candidate tier names: Builder / Scale / Enterprise — pending marketing
  sign-off.
- Usage-metering estimate: 3 weeks (JT); release is 5 weeks out.

## Archive

- Carlos out until Thursday — superseded; he was back Thursday.
- Pricing meeting moved to Wednesday — happened; notes captured above.

---

**Counts:** 22 items in → 16 appear above, 3 merged as duplicates (second
incident-writeup TODO; trust-page/uptime idea into one; "EU p99 still
degraded + filing ticket" into the EU rollback problem), 3 folded into other
items (Sept-release confirmation into the usage-tier decision; Dana's 60-day
confirmation into the notice question; the metering estimate into the
release-slip question), 1 flag (⚡ rollback conflict), 0 parked. 22 = 16 + 3
+ 3, all accounted for; Parking-lot section omitted because it's empty.
