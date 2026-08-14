# pacs-access-control-architecture

**Category:** Science, Biology & Medicine
**Platform:** Claude (Archived)
**Original Path:** claude/archived-treearch/pacs-access-control-architecture

## Description
Physical Access Control System (PACS) architecture end to end: credential -> reader -> controller/panel -> head-end server -> door hardware. TRIGGER: PACS stack & topology (centralized/distributed/edge-PoE, offline/degraded mode); door hardware (electric strikes, maglocks, electrified mortise/cylindrical locks, latch retraction, REX, DPS); fail-safe vs fail-secure + egress/life-safety code (NFPA 101/80, IBC, free egress, maglock release, stair re-entry); head-end (cardholder DB, access levels/schedules/zones, anti-passback, mustering, audit log); integration (elevators, turnstiles/mantraps, visitor mgmt, intrusion, video). SKIP: Wiegand/OSDP -> access-credentials-wiegand-osdp; RF -> rfid-fundamentals-bands-standards; NFC -> nfc-ndef-protocols; chip/key -> contactless-smartcards-mifare-desfire; attacks -> rfid-nfc-access-attacks-defenses; mobile/UWB -> mobile-credentials-aliro-uwb (planned); convergence -> physical-security-convergence-standards (planned).

---

# Physical Access Control System (PACS) Architecture

The end-to-end electronic access-control stack: **credential → reader → door
controller/panel → head-end (access-control server) → door hardware**, plus the
sensors and power that make a door work and the life-safety code that governs how
it must release. This skill owns the **system architecture and door hardware**.
It does **not** cover the credential RF/chip internals, the reader↔controller
*frame format*, attacks, mobile credentials, or the standards/compliance
deep-dive; see the SKIP list in the frontmatter and the cross-links inline.

> **Life-safety disclaimer.** Egress, fire-door, and locking-arrangement rules in
> §3 are summarized from secondary code commentary (idighardware/Allegion, SIA,
> manufacturer literature), not the bound code text. **Code adoption varies by
> AHJ (authority having jurisdiction), edition, and occupancy type.** Always
> verify against the *adopted* edition of NFPA 101, NFPA 80, and the IBC for the
> specific jurisdiction and have the AHJ sign off. Do not design egress from this
> file alone.

---

## 1. The PACS stack (five layers)

A PACS is the layered system that decides, at an opening, whether to unlock.

```
[ CREDENTIAL ]   card / fob / PIN / mobile / biometric  ── what you present
      │ (RF / NFC / BLE-UWB / keypad / sensor)
      ▼
[ READER ]       reads credential, sends data up        ── capture
      │ (Wiegand or OSDP wire — see access-credentials-wiegand-osdp)
      ▼
[ DOOR CONTROLLER / PANEL ]  the local "brain"          ── DECISION
      │  - holds a cached copy of the cardholder DB + access rules
      │  - decides grant/deny locally (works during head-end/network outage)
      │  - drives the lock relay; reads REX + door position; raises alarms
      │ (TCP/IP, usually Ethernet, often over UPS-backed switches)
      ▼
[ HEAD-END / ACCESS-CONTROL SERVER (+ client / cloud) ] ── POLICY + RECORD
      │  - master cardholder DB, access levels, schedules, zones
      │  - pushes enrollments/revocations/schedules DOWN to controllers
      │  - receives + logs every event (grant/deny/alarm/door state) UP
      ▼
[ DOOR HARDWARE ]  lock + REX + door contact + power     ── the physical opening
```

Key architectural facts that recur in every deployment:

- **The decision is made at the controller, not the server.** Controllers hold a
  **local copy of the cardholder database** and the access rules, so grant/deny,
  door-monitoring, and alarms **keep working when the head-end or network is
  down** (this is *offline / degraded / standalone* mode). Buffered events upload
  when comms are restored. This is the single most important reliability property
  of a real PACS; never describe the server as making the live unlock decision.
- **The server's job is policy distribution and the record of truth.** It pushes
  new enrollments, revocations, and schedule changes *down*; it receives and logs
  grants, denials, alarms, and door-state changes *up*.
- **Door capacity is per controller.** A common intelligent controller handles up
  to ~**8 doors** directly; modular/distributed lines scale by adding downstream
  door modules (e.g. a Mercury LP1501 is edge-capable for one opening but expands
  to ~17 openings via downstream SIO / MR modules). Big head-end platforms scale
  to tens of thousands of doors and hundreds of thousands of cardholders.

Deeper topology and head-end-feature detail: `references/head-end-and-topology.md`.

---

## 2. Door hardware

The physical opening is a small system of its own: a **locking device**, a
**request-to-exit** sensor, a **door position switch**, and a **power supply**
with battery backup. Get the lock *type* and its *fail mode* right or you create
either a security hole or a life-safety violation (§3).

### 2.1 Locking devices

| Device | What it is | Fail mode availability | Free mechanical egress? | Typical use / notes |
| --- | --- | --- | --- | --- |
| **Electric strike** | Replaces the strike plate in the *frame*; releases the latch/keeper while the lockset stays mechanically intact. Pairs with a cylindrical/mortise lock, rim lock, or panic bar. | **Fail-safe OR fail-secure** (model/orientation dependent) | Yes — turning the lever/pushing the panic bar still retracts the door's own latch | Cheapest, least invasive (frame-only); latch keeper is exposed = lower attack resistance. On a **fire door the strike MUST be fail-secure** and fire-listed (NFPA 80). |
| **Magnetic lock (maglock)** | Electromagnet on the frame + armature plate on the door; energized magnet holds the door. | **Fail-safe ONLY** — no holding force without power | **NO** — a maglock does not retract anything; it provides *no* free egress on its own | Requires a code-compliant **release path** (REX sensor or door-mounted switch) **and** fire-alarm + power-loss release to be legal egress (§3.3). |
| **Electrified mortise lock** | A mortise lockset with a built-in solenoid; "storeroom function" = outside lever locked, **inside lever always free**. | **Fail-safe OR fail-secure** (and "electrically locked" vs "electrically unlocked" variants) | **Yes** — inside lever always retracts the latch | More secure than a strike (deadlatch protected); fire-listed versions exist. Common high-security choice. |
| **Electrified cylindrical lock** | Same idea, cylindrical (bored) lockset. Storeroom function, inside lever free. | **Fail-safe OR fail-secure** | **Yes** — inside lever free | Lower cost than mortise; watch deadlatch alignment if combined with a strike. |
| **Electric latch retraction (EL/QEL) on a panic device** | Motor/solenoid retracts the exit-device latch on command (e.g. for scheduled unlock or remote buzz-in). | **Fail-secure** (latch projects on power loss) | **Yes — always** (it's a panic bar: pushing the bar mechanically retracts the latch) | Pairs access control with code-required panic hardware; high inrush current — size the supply. Fire-exit-hardware versions exist. |
| **Electromagnetic shear lock** | Maglock variant that engages in shear (top of door). | **Fail-safe** (some have a delay) | No — same egress caveat as a maglock | Niche; same release-path requirements as a maglock. |

**The decisive distinctions:**
- **Strike vs lock:** a *strike* lives in the frame and only frees the latch; a
  *lock* lives in the door and changes the lever function. Locks are more secure
  (no exposed keeper, deadlatch protected); strikes are cheaper and lockset-agnostic.
- **Maglock is the special case:** it is the *only* common device that gives **no
  inherent free egress** and is **fail-safe-only**, which is exactly why it draws
  the heaviest code requirements (§3.3). Every other device on the list lets you
  exit mechanically by operating the inside lever or panic bar.

### 2.2 Request-to-exit (REX / RTE)

REX tells the system that an exit is *intended/legitimate* so the controller can
**shunt** (suppress) the door-forced-open alarm and, for a maglock, **release**
the lock. Two main forms:

- **PIR motion sensor (ceiling/over-door):** detects a person approaching from
  the secure side; releases/shunts hands-free. Cheap and common, but it fires on
  *any* motion (passersby, HVAC airflow), so it is a weaker egress *authorization*
  than a switch.
- **Switch in the hardware:** a **REX switch inside the lever set / mortise lock**
  (fires when the lever is turned) or a **push-bar / panic-device switch** (fires
  when the bar is pressed). This is the preferred form because the event is tied
  to an actual hand on the egress hardware, and it shunts the alarm *without*
  needlessly energizing the lock.
- **Push-to-exit button (RTE button):** a manual button releasing the door; when
  it is the code release for a maglock it has strict requirements (§3.3).

> A REX is *not* itself the life-safety egress for a maglock — a PIR-only REX can
> fail (it's electronics on the same power). Code requires the release to also
> occur on **power loss** and **fire alarm** regardless of the REX (§3.3).

### 2.3 Door position switch (DPS) / door contact + monitoring

A **magnetic door contact (DPS / DCS)** reports door **open vs closed**. With it,
the controller generates the two core door alarms:

- **Door Forced Open (DFO):** door opened with **no valid grant and no REX** — a
  break-in or propped-with-force condition.
- **Door Held Open (DHO):** door stayed open **longer than the allowed time** after
  a valid grant (tailgating/propping risk).

A separate **latch-bolt / latch monitor** can confirm the door actually *latched*
(closed ≠ secured if the latch didn't catch). On a valid grant or REX the
controller **shunts** the DPS for a programmed interval so a normal entry/exit
doesn't raise DFO/DHO.

### 2.4 Power supply + battery backup

- **Listed power supply** feeds the locks and (often) the controller. UL 294 is
  the access-control system standard; where the panel also does intrusion it adds
  UL 1076. Lock voltage is typically **12 or 24 VDC**.
- **Battery / standby backup** holds power to **fail-secure** locks and the
  controller through an AC outage (so doors don't unlock just because the mains
  blinked). Switchover should be glitch-free. A *fail-safe* lock **releases** on
  power loss by definition, so battery-backing it keeps it *locked* through an
  outage — which is correct for a fail-safe lock you want held, but **do not
  battery-back a fail-safe maglock that serves as code egress**: its §3.3 release
  *requires* it to drop on loss of power, so that output is deliberately left
  un-backed (or fire/power-drop-interfaced).
- **Fire-alarm interface (FACP):** a relay/input that **drops power to fail-safe
  egress locks (e.g. maglocks) on fire alarm**, and can trigger **elevator
  recall / HVAC shutdown** via Form-C dry contacts. This interface is what makes a
  maglock legal (§3.3).
- **Per-output fail mode:** modern multi-output supplies let each lock output be
  set **fail-safe or fail-secure** and tie selected outputs to the fire interface.

Electrical detail (PoE power budgets, inrush, separate lock supplies, REX/DPS
wiring) is in `references/door-hardware-and-power.md`.

---

## 3. Fail-safe vs fail-secure and egress / life-safety (READ CAREFULLY)

This is the section that must be correct. The governing principle in US life
safety is **free egress: a person must always be able to leave** without a key,
tool, special knowledge, or effort, and (generally) with **one motion/one hand**.
Electrified locking must never defeat that.

### 3.1 The two terms — power vs lock state

| Term | On loss of power the door becomes… | Power is applied to… | Mnemonic |
| --- | --- | --- | --- |
| **Fail-SAFE** | **UNLOCKED** (on the access/outside side) | **LOCK** the door | "safe" for *people* — power off ⇒ open |
| **Fail-SECURE** | **LOCKED** (on the access/outside side) | **UNLOCK** the door | "secure" for *property* — power off ⇒ stays locked |

"Access side" matters: it describes the **outside/entry** behavior. **Free egress
from the inside is a separate property of the hardware** (see §3.2).

### 3.2 Fail-safe/fail-secure is about *ingress* — most hardware still gives free *egress*

The most common point of confusion: **fail-secure does NOT mean people are
trapped inside.** On a normal lockset or panic device, the **inside lever / panic
bar always retracts the latch mechanically**, so occupants exit at all times —
*regardless of whether the device is fail-safe or fail-secure*, and *regardless of
power*. Fail-safe vs fail-secure only changes what happens on the **outside** when
power is lost.

- "Most products that are not special locking arrangements provide free egress
  whether they are fail safe or fail secure." (idighardware/Allegion)
- The **exception is the maglock** (and shear lock): it has **no mechanical
  release at all**, so it cannot provide free egress by itself and is treated as a
  *special locking arrangement* (§3.3).

### 3.3 The maglock problem and its two code paths

Because a maglock is fail-safe-only with **no free egress**, code permits it only
under one of two named arrangements (apply **one complete set**, not a mix). Both
sit in NFPA 101 (Life Safety Code) and the IBC (these are the common 2015→2024
provisions; **verify the adopted edition with the AHJ**):

**Path A — Sensor release ("Access-Controlled Egress Doors" / IBC "Sensor Release
of Electrically Locked Egress Doors").** The maglock must unlock on **all** of:
1. **Loss of power** to the sensor *or* to the lock/locking system;
2. **Activation of the fire alarm / sprinkler** (stays unlocked until the system
   is reset);
3. A **manual push button** that is: mounted **40–48 in** above the floor, **within
   5 ft** of the door, marked **"PUSH TO EXIT,"** **directly interrupts lock
   power independent of other electronics**, and holds the door unlocked for a
   **minimum of 30 seconds**;
4. a **REX sensor** that detects an approaching occupant and releases the lock.

**Path B — Door-hardware release ("Electrically Controlled Egress Door Assemblies"
/ IBC "Door Hardware Release of Electrically Locked Egress Doors").** Release comes
from a **switch built into the door-mounted hardware** (lever or panic bar). It
must: operate with **one hand / obvious operation under all lighting**, **directly
interrupt power** to unlock immediately, and **unlock on loss of power** to the
lock. **Note:** in Path B, fire-alarm activation does **not** itself have to
release the lock (the hand-on-hardware switch is the always-available release).

> Practical read: **Path A** is the "maglock + PIR + Push-to-Exit + fire-alarm
> drop" recipe; **Path B** is "maglock released by a switch in the panic bar/lever."
> The 30-second / 40–48 in / within-5-ft / "PUSH TO EXIT" numbers are the
> classic, frequently-tested Path-A specifics, but **confirm against the adopted
> edition**, since wording and details have shifted across IBC/NFPA cycles.

There are also tighter **special locking arrangements**: *delayed-egress* (door
releases after a 15-second nuisance-delay countdown, used to deter theft/elopement)
and *controlled-egress* (e.g. healthcare units holding patients; staff can release
for evacuation). These have their own occupancy limits and are out of scope here
beyond naming them; treat them as AHJ-driven.

### 3.4 Fire doors — fail-SECURE is mandatory

On a **fire-rated door assembly**, an **electric strike (and the locking hardware)
must be fail-secure and fire-listed** per **NFPA 80**. A fire door must **stay
latched** to hold back fire/smoke; a fail-safe device that drops the latch on
power loss (or on the fire alarm) would defeat the rated assembly. So: fire door ⇒
fail-secure latch, but the inside still gives **free egress** mechanically (§3.2).
Maglocks on fire doors are generally not used for this reason.

### 3.5 Stairwell re-entry — fail-SAFE is the requirement

Stairwell doors are the opposite case. Occupants must be able to **leave the
stair** (re-enter a floor) if that floor's landing is blocked. Stair-side
electrified locks are therefore **fail-safe** and must **unlock on fire alarm, on
a signal from the fire command center, and on power loss** (e.g. the 2024 IBC
stair re-entry provision). Elevator-lobby egress and some healthcare doors carry
similar fail-safe-on-alarm + two-way-communication requirements.

### 3.6 Decision summary

| Situation | Required/typical fail mode | Why |
| --- | --- | --- |
| General secured perimeter/interior door | Fail-**secure** (default) | Door stays locked through a power blip; inside still egresses mechanically |
| Fire-rated door assembly | Fail-**secure** + fire-listed (NFPA 80) | Door must stay latched to hold the rating |
| Maglock anywhere | Fail-**safe** (only option) + Path A or B release + fire/power drop | No mechanical egress; must release for life safety |
| Stairwell re-entry door | Fail-**safe**, unlock on alarm/FCC/power loss | Occupant must be able to leave the stair |
| Where you want doors to **open** automatically in a fire (e.g. maglocked egress) | Fail-**safe** tied to FACP | Free egress during the emergency |
| Where you want doors to **stay locked** in a power loss | Fail-**secure** + battery backup | Property protection + buffered controller |

---

## 4. Controllers / panels — topologies

How the decision-makers are arranged. All three coexist in the field; the trade is
**wiring/labor vs. door cost vs. failure domain**.

| Topology | Where the brain lives | Wiring pattern | Pros | Cons |
| --- | --- | --- | --- | --- |
| **Centralized (panel)** | Multi-door **panels** in an IDF/closet; readers + locks home-run back to the panel | Long home-runs from each door to the closet | Fewer network nodes; panel on one UPS; mature, high-density | Heavy copper/labor; a panel failure can drop **all** its doors |
| **Distributed** | **Intelligent controller** + **downstream door/IO modules** on a local bus (e.g. Mercury LP + MR/SIO) | Controller central, short runs to per-door modules | Scales by adding modules; less head-run cable than pure centralized | Still a controller failure domain; bus design needed |
| **Edge / PoE-IP** | A **single-door controller (or reader-controller) at each door**, powered + networked over **PoE** | One Cat-5/6 PoE drop per door; no separate lock supply needed for low-draw locks | Minimal failure domain (one door); UPS/standby centralized at the **switch/IDF**; least home-run copper; easy adds | Network-dependent; PoE power budget limits high-draw locks (maglocks/EL) — see refs; more managed endpoints |

Notes:
- **Edge does not remove the controller; it relocates it to the door.** The
  grant/deny logic and offline caching still live in the edge device.
- **PoE budgets (power *to the device*):** 802.3af ≈ 12.95 W (15.4 W at the port),
  802.3at (PoE+) ≈ 25.5 W, 802.3bt (PoE++) higher. A continuous-draw **maglock**
  or high-inrush **electric latch retraction** can exceed a modest PoE budget, so
  those often keep a **dedicated lock power supply** even in an "IP" design.
  Detail in `references/door-hardware-and-power.md`.
- **I/O on a controller:** reader port(s) (Wiegand/OSDP), lock relay output(s),
  REX input, DPS/door-contact input, tamper input, aux inputs/outputs (for FACP
  drop, elevator recall, HVAC, alarm shunts).
- **Cloud / ACaaS** is a head-end deployment choice layered on any of the above:
  the controller still decides locally; the cloud is the management/record plane.

---

## 5. Head-end (access-control software) — what it does

The management server + client (on-prem or cloud) is the **policy source and the
record of truth**. Core feature set:

- **Cardholder database:** identities, credential bindings (card/fob/mobile/PIN/
  biometric), status (active/suspended), and the credential→person mapping the
  controllers cache.
- **Access levels / privileges:** the named bundles of *(door or door-group) ×
  (time schedule)* that say *who* may go *where* and *when*. Assigned to
  cardholders or groups.
- **Time schedules / time zones + holidays:** the *when*: business hours,
  shift windows, auto-unlock/auto-lock of public doors, holiday overrides.
- **Zones / areas:** logical groupings of doors used for anti-passback, mustering,
  occupancy, and area-based rules.
- **Anti-passback (APB):** prevents a credential entering an area twice without an
  exit (badge-sharing / tailgating control). Modes: **hard** (violation = deny),
  **soft** (allow but log/flag), **timed** (same reader can't re-fire for N
  seconds), **area/zone** (must traverse areas in order), and **global** APB
  across all controllers. Requires **in *and* out readers** on the area boundary.
- **Mustering / muster reports:** in an evacuation, a real-time roll of *who is
  still inside* a zone (and who has badged to a muster point), derived from the
  same in/out tracking as APB. Used for emergency accountability.
- **Occupancy / area control:** live count per area (max/min occupancy rules,
  two-person rule, first-in/last-out unlock).
- **Audit log / event history:** time-stamped record of **every** grant, denial,
  REX, door-state change, and alarm — buffered at the controller during outages,
  uploaded on reconnect. The compliance/forensic backbone.
- **Alarm monitoring:** live operator view of **DFO, DHO**, tamper, comm-loss,
  forced/held, and (when integrated) intrusion points — with acknowledge/clear
  workflow, maps, and linked video.

Deeper definitions and worked examples in `references/head-end-and-topology.md`.

---

## 6. Integration points (named, not deep)

A PACS rarely stands alone; the controller's I/O and the head-end's interfaces tie
it to neighboring systems:

- **Elevators / destination dispatch:** floor-level permissions gate which floors
  a credential can select; with **destination-dispatch** the PACS passes the
  authorized floor(s) to the elevator controller (often API-based) so the car
  assignment respects access rights. The PACS also drives **fire-service elevator
  recall** via a relay.
- **Turnstiles / mantraps (interlocks):** turnstiles enforce **one-person-per-
  credential** anti-tailgating at lobbies; a **mantrap / access vestibule** uses
  **two interlocked doors** (the first must close/lock before the second opens),
  often with APB and occupancy rules, for high-security one-at-a-time passage.
- **Visitor management (VMS):** issues **temporary, time-bound credentials** on
  host approval and auto-revokes at checkout; writes the visitor into the same
  cardholder/audit model.
- **Intrusion detection:** arm/disarm by valid badge, share door/zone state, and
  cross-trigger (e.g. a DFO raising an intrusion alarm); panels doing both invoke
  UL 1076.
- **Video / VMS:** access events **bookmark and call up video** (a badge-in or a
  DFO jumps the operator to that camera), and analytics add **anti-tailgating** at
  the portal. The two systems negotiate **which is "master" for alarms**.

These are integration *surfaces*; the protocol/standards depth (and convergence
of physical + logical security) belongs to
`physical-security-convergence-standards`.

---

## 7. Availability & reliability (failure modes)

- **Decision survives the network:** controllers run on **cached data** in
  offline/degraded mode — doors keep granting/denying and logging (buffered) when
  the head-end or LAN is down. Design so an IDF/switch outage degrades to *local*
  control, not *no* control.
- **Power:** **battery/standby backup** rides through AC loss; **fail-secure +
  battery** keeps a door locked through an outage, while **fail-safe** releases by
  design (so you choose the fail mode per the §3 life-safety requirement, *then*
  size backup accordingly).
- **Fire wins over security:** the **FACP interface drops fail-safe egress locks**
  on alarm; life safety overrides access control, always.
- **Failure domain by topology:** centralized panel = many doors lost on one
  failure; **edge/PoE** = one door lost per failure but more nodes and a network
  dependency (mitigate with PoE UPS at the switch). Pick the blast radius you can
  live with.
- **Monitoring the system itself:** comm-loss, tamper, power-fault, and low-battery
  are themselves **alarms** in the head-end; a PACS should report its own health.

---

## Sources

Authoritative and reputable-integrator references consulted (accessed June 2026).
Code summaries are secondary commentary; verify against the adopted edition with
the AHJ.

1. **idighardware (Lori Greene / Allegion) — "Decoded: Fail Safe vs. Fail Secure
   – When and Where?" (2023)**: fail-safe/fail-secure definitions; maglock is
   fail-safe-only; NFPA 80 fail-secure on fire doors; stair re-entry fail-safe;
   most locksets give free egress either way:
   https://idighardware.com/2023/10/decoded-fail-safe-vs-fail-secure-when-and-where/
2. **idighardware — "Code Requirements for Electromagnetic Locks" (2017)**: the
   two maglock code paths (sensor-release vs door-hardware-release), the 40–48 in
   / within-5-ft / "PUSH TO EXIT" / 30-second / power-loss / fire-alarm specifics,
   and the NFPA 101 vs IBC section naming:
   https://idighardware.com/2017/06/code-requirements-for-electromagnetic-locks/
3. **Allegion — "Article Decoded: Fail Safe vs. Fail Secure – When and Where?"
   (Lori Greene, PDF)**: same fail-safe/secure framework, free-egress principle:
   https://us.allegion.com/content/dam/allegion-us-2/web-documents-2/Article/Allegion_Fail_Safe_Fail_Secure_Article_112140.pdf
4. **NIH ORF Technical Bulletin — "Fail Safe vs. Fail Secure Electronic Locksets"
   (Issue 101, 2020)**: government facilities guidance on fail-safe vs fail-secure
   selection and fire/egress interplay:
   https://orf.od.nih.gov/TechnicalResources/Documents/Technical%20Bulletins/20TB/Fail%20Safe%20vs.%20Fail%20Secure%20Electronic%20Locksets%20June%202020%20-%20Technical%20Bulletin%20UPDATED_508.pdf
5. **Avigilon — "Physical Access Control System (PACS): Components + Examples"** —
   five-layer component model; controller makes the live decision; server stores
   users/privileges/audit logs; cloud vs on-prem; video integration:
   https://www.avigilon.com/blog/physical-access-control
6. **ButterflyMX — "Physical Access Control System (PACS): Your Complete Guide"** —
   control panel as the brain; local memory so it runs during comms loss;
   controllers ~8 doors; event flow up to the server:
   https://butterflymx.com/blog/physical-access-control/
7. **GSA — "Physical Access Control Systems (PACS) Customer Ordering Guide" (2025)**
   federal/FICAM PACS architecture: head-end server, workstations, controllers
   on a high-speed network; component taxonomy:
   https://buy.gsa.gov/api/system/files/documents/physical-access-control-systems-pacs-customer-ordering-guide-508-compliant-12-01-2025-final.pdf
8. **Mercury Security — Intelligent Controllers / LP1501 / MR-SIO modules** —
   distributed open-architecture controllers; edge-capable; expansion via
   downstream serial I/O modules; embedded Linux + crypto memory:
   https://www.mercury-security.com/controllers/ and https://www.mercury-security.com/products/lp1501/
9. **IPVM — "Access: IP Readers vs. Control Panels"**: centralized panel vs
   edge/IP-reader topology trade-offs; failure domain; cabling:
   https://ipvm.com/reports/access-ip-readers-vs-control-panels
10. **Kisi — "PoE Door Access Control"** and **Anixter — "Powering Access Control
    over PoE"**: PoE delivery (802.3af/at/bt budgets), edge controllers at the
    door, centralized UPS at the switch, power-budget cautions:
    https://www.getkisi.com/guides/poe and https://www.anixter.com/en_us/resources/literature/techbriefs/powering-access-control-over-poe.html
11. **Altronix / LifeSafety Power — access power supplies (datasheets +
    Locksmith Ledger "Powering Up Access Control")**: per-output fail-safe/
    fail-secure, fire-alarm interface to cut power to egress locks, standby
    battery switchover, Form-C relays for elevator recall/HVAC:
    https://www.altronix.com/products/AL300ULM and https://www.locksmithledger.com/electronics-access-control/control-panels-accessories/article/21151325/powering-up-access-control
12. **SDC — "Code Compliant Access Controlled Egress Door" (application guide)** —
    maglock egress arrangement wiring: REX, push-to-exit, fire-alarm release:
    https://sdcsecurity.com/docs/solution6.pdf
13. **Locksmith Ledger / Park Avenue Locks / Morrison-Maierle — electrified
    hardware comparisons**: electric strike vs electrified mortise/cylindrical
    vs electric latch retraction; storeroom function; deadlatch caution:
    https://www.locksmithledger.com/electronics-access-control/article/12164378/failsafe-va-failsecyre and https://www.parkavenuelocks.com/blog/post/electric-strikes-vs-electric-locks-which-one-should-you-choose
14. **Silva Consultants / Sielox glossary / Kisi — anti-passback & door
    monitoring**: hard/soft/timed/area/global APB; mustering; DFO/DHO; REX:
    https://www.silvaconsultants.com/anti-passback-in-access-control and https://www.getkisi.com/guides/anti-passback
15. **UL 294 (Access Control System Units) / UL 1076 (Proprietary Burglar Alarm
    Units)**: the listing standards a PACS and intrusion-capable panel must meet:
    https://www.shopulstandards.com/ProductDetail.aspx?productId=UL294

**Uncertainty / verify-before-use notes:**
- All §3 egress, fire-door, and maglock numbers (30 s, 40–48 in, within 5 ft,
  stair re-entry triggers) are summarized from **secondary code commentary** and
  have **changed across IBC/NFPA editions**. They are correct as commonly cited
  but are **not a substitute for the adopted code text and AHJ approval**.
- Door capacity per controller (~8 doors; ~17 with expansion) is vendor-typical,
  not universal; confirm per product line.
- PoE power sufficiency for a given lock depends on the specific device's hold/
  inrush current vs the switch's per-port budget; size it, don't assume.