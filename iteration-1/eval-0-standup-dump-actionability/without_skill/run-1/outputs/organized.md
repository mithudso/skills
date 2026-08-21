# Standup Notes — Week of June 29

Cleaned up from the week's raw dump. Two things you asked for are up top: **what you owe / who owes what** and **what got decided**. Everything from the notes is preserved below; nothing added.

> **Heads-up:** The Friday latency incident is **not actually closed.** The rollback only fixed the primary cluster — the analyzer change also shipped to the EU cluster and was never rolled back there, so EU p99 is still degraded (found Thursday). This affects your incident writeup.

---

## Action items — who owes what

| Owner | Item | Status |
| --- | --- | --- |
| **You** | File the incident report/writeup for Friday's latency spike | **Open** — noted Mon and again Wed. Must now cover the EU cluster gap (see incident below), not just the primary rollback |
| **Sam** | Update the pricing-page draft (Sam has the doc) | **Blocked / on hold** — waiting on final tier names; decision is to hold until marketing signs off |
| **Sam** | Ping marketing about the tier names | **Open** |
| **Carlos** | File a ticket for the degraded EU p99 | **In progress** — was filing it Thursday |
| **Dana** | Confirm the EU price-change notice period | **Done** — confirmed 60 days (see decisions) |

**Waiting on someone else:** Marketing sign-off on the tier names — this is the single blocker holding up the pricing page.

---

## What got decided

- **Usage-based tier for the API is a go.** Flat tier is killed. **Confirmed** for the September release.
- **Price-change notice periods:** 30 days standard; **60 days for EU customers.** Dana confirmed the EU 60-day figure — it's in the DPA addendum. (This resolves the earlier "EU 60??" open question.)
- **Hold the pricing page** until marketing signs off on the tier names.
- **September release date holds.** Usage-metering work is estimated at 3 weeks (jt) and the release is 5 weeks out, so no slip.

**Decided but not final — needs sign-off:**
- Proposed tier names are **builder / scale / enterprise**. Marketing has **not** signed off yet.

---

## Open incident — Friday latency spike

- **Fri night:** Search latency spiked, p99 hit 800ms (flagged by Kavya).
- **Response:** Rolled back the analyzer change. Kavya confirmed the rollback brought p99 back to 120ms.
- **Thursday update (Carlos, back from being out):** The analyzer change had also shipped to the **EU cluster** and was **never rolled back there.** EU p99 is **still degraded.** Carlos is filing a ticket.
- **Net:** Primary cluster recovered; EU cluster still broken. Incident stays open until EU is rolled back and the writeup reflects both clusters.

---

## Resolved question

- **Does the Sept release slip if usage metering isn't done in time?** → **No.** Metering is a 3-week job (jt) and the release is 5 weeks out.

---

## Ideas / parking lot (not decided)

- Expose the latency dashboard to customers as a **"trust page."**
- The trust page could also show **uptime** — the way status.io does it.

---

## Timeline (context)

- **Mon:** Standup. Carlos out until Thursday. Pricing meeting moved to Wednesday.
- **Wed:** Pricing meeting.
- **Thu:** Carlos back; EU cluster issue surfaced.
