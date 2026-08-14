# Head-End Features & Controller Topology — Detail

Companion to `pacs-access-control-architecture/SKILL.md` §1, §4, §5. Deeper
definitions of the access-management software model and the controller-topology
trade-offs.

---

## 1. The data model the controller caches

The head-end is the master; each controller caches the subset it needs to decide
locally. The objects:

- **Cardholder / identity:** the person. Holds bindings to one or more
  **credentials** (card number/format, mobile credential ID, PIN, biometric
  template ref) and a **status** (active / suspended / expired).
- **Credential:** the presentable token; maps to exactly one cardholder. The
  reader captures it (as a Wiegand/OSDP bitstream — see
  `access-credentials-wiegand-osdp`), the controller matches it against the cache.
- **Access level / access group:** the reusable bundle of **(door or door-group)
  × (time schedule)**. "Lab Techs" = {Lab doors} × {Mon–Fri 06:00–20:00}. Assign
  the *level* to people, not doors to people one-by-one.
- **Time schedule / time zone (+ holidays):** named time windows. Drive both
  access (when a level is valid) and **auto-unlock/auto-lock** of public doors,
  with **holiday tables** overriding normal days.
- **Door / reader / door-group:** the controlled openings and their logical sets.
- **Area / zone:** doors grouped by the space they bound; the unit for
  anti-passback, mustering, and occupancy.

A grant decision = *does this credential's cardholder hold an access level that
includes this door at this moment (and pass any APB/occupancy rule)?* It is
computed **at the controller** against cached data.

---

## 2. Anti-passback (APB) — modes

APB stops a credential from entering an area twice without an intervening exit
(badge-sharing, tailgating, hand-back-the-card). Requires **in and out readers**
on the area boundary so the system knows which side you're on.

| Mode | Behavior on violation | Use |
| --- | --- | --- |
| **Hard** | **Deny** the second entry (or the exit if no matching entry) | Strict areas; enforced occupancy/parking |
| **Soft** | **Allow** but **log/flag** for the operator | Where denying would cause operational pain but you still want the audit signal |
| **Timed** | Same reader/credential **can't re-fire for N seconds** | Stop pass-back of one card to a second person right behind |
| **Area / zone** | Must traverse areas **in order** (can't appear in zone 3 without zone 1→2) | Layered/nested security; prevents "skipping" to an inner zone |
| **Global** | APB enforced **across all controllers/host**, not just one panel | Multi-panel campuses; needs host (or peer) coordination |

APB state can get "stuck" (someone tailgates out, so the system thinks they're
still in), so most systems provide an **APB reset / forgive** and a
**first-card-in** exception. APB requires reliable in/out reads; a single-reader
(entry-only) door can't do true APB.

---

## 3. Mustering & occupancy

- **Mustering / muster report:** during an evacuation, a **live roll of who is
  still inside** an area, plus who has badged to a **muster point** reader outside.
  Built from the same in/out tracking as APB. The emergency-accountability
  deliverable: "everyone accounted for" / "these N are unaccounted."
- **Occupancy counting / area control:** running count per area; supports
  **max-occupancy** (deny entry when full: labs, cleanrooms, parking) and
  **min-occupancy / two-person rule** (no single person alone in a sensitive
  area; both must badge).
- **Two-person (dual-authorization) rule:** two valid credentials required to
  unlock (vaults, datacenters, pharma). Distinct from occupancy but uses the same
  area model.

These features all **depend on in/out reader pairs and APB-grade tracking**: you
cannot muster accurately if entries/exits aren't both read.

---

## 4. Audit log & alarm monitoring

- **Audit log / event history:** time-stamped record of **every** transaction:
  grant, denial (with reason: no access, wrong schedule, APB, expired), REX,
  door-state change, operator action, and alarm. **Buffered at the controller**
  during a head-end/network outage and **uploaded on reconnect**, so the record is
  complete even through downtime. This is the forensic/compliance backbone.
- **Alarm monitoring:** the live operator surface. Core alarm types:
  - **Door Forced Open (DFO):** opened with no grant and no REX.
  - **Door Held Open (DHO):** open past the allowed time after a grant.
  - **Tamper:** reader/enclosure opened.
  - **Comm loss:** controller or downstream module unreachable.
  - **Power fault / low battery / AC fail:** the power plant reporting itself.
  - **Intrusion points** (when the panel/integration covers intrusion).
  Operators **acknowledge/clear** with notes; alarms can drive **maps**, **linked
  video** (call up the camera on the event), and notifications.

---

## 5. Controller topology — deeper trade

| Axis | Centralized panel | Distributed (controller + modules) | Edge / PoE-IP |
| --- | --- | --- | --- |
| Brain location | Closet/IDF panel, many doors | Controller + downstream door/IO modules on a bus | One controller **at each door** |
| Cabling | Long home-runs (reader + lock) per door to closet | Controller central; short runs to per-door modules | One PoE drop per door; minimal home-run copper |
| Failure domain | **All doors on the panel** drop together | Controller (or a bus segment) is the domain | **One door** per failure (but many nodes) |
| Scale pattern | Add panels | Add downstream modules to a controller | Add doors = add network endpoints |
| Power/backup | One UPS per panel in the closet | Controller + supply in closet | **Centralized UPS at the switch/IDF** |
| Network dependence | Panel↔host; doors work offline regardless | Same | More network-dependent; mitigate with switch UPS |
| Best fit | High door density in a building core | Mixed density, modular growth | Distributed doors, retrofits, easy adds |

Key truths:
- **Edge relocates the controller to the door; it does not delete it.** The
  offline grant/deny + caching still happen in the edge device. "Reader-only +
  cloud-decides" designs that *don't* cache locally trade away the core PACS
  reliability property — scrutinize any product that claims the cloud makes the
  live unlock decision.
- **Open-architecture controllers** (e.g. the Mercury platform used under many
  head-end brands) let the *head-end software* and the *controller hardware* be
  chosen independently, and expand via **downstream serial I/O (SIO) modules** —
  the practical face of "distributed."
- **Cloud / ACaaS** is a **head-end deployment model** (management + record in the
  cloud), orthogonal to centralized/distributed/edge at the door. You can run a
  cloud head-end over edge controllers, or an on-prem head-end over centralized
  panels, etc.

---

## 6. Controller I/O recap

A door point on a controller typically exposes: reader port(s) (Wiegand/OSDP),
lock relay output(s), REX input, DPS/door-contact input, optional latch/bond
input, tamper input, and aux inputs/outputs (FACP drop, elevator recall, HVAC
shutdown, alarm shunt). Door capacity per controller is vendor-typical (~8 doors
direct; more via expansion modules) — confirm per product line.
