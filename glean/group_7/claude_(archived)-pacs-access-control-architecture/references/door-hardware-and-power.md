# Door Hardware, REX/DPS Wiring & Power — Detail

Companion to `pacs-access-control-architecture/SKILL.md` §2, §4, §7. Electrical and
sensor-wiring depth that doesn't belong in the main body. **Life-safety release
behavior is governed by SKILL.md §3 and the adopted code — this file is about the
electrical/wiring reality, not a substitute for code.**

---

## 1. Locking devices — electrical behavior

| Device | Voltage (typical) | Current draw character | Fail mode | Notes |
| --- | --- | --- | --- | --- |
| Electric strike | 12 / 24 VDC (some AC) | Steady holding current when energized; modest | Fail-safe or fail-secure (often field-selectable by jumper/orientation) | Continuous-duty vs intermittent-duty ratings matter; a fail-secure strike on a fire door must be **continuous-duty + fire-listed** so it can sit energized-locked indefinitely. |
| Magnetic lock (maglock) | 12 / 24 VDC | **Continuous** draw (the magnet must stay energized to hold) — often 250–500 mA per lock | Fail-safe only | The continuous draw + the fact it must drop on fire/power is why maglocks usually sit on a **dedicated, fire-interfaced supply**, not a PoE port. Holding force commonly 600 or 1200 lbf. |
| Electrified mortise / cylindrical | 12 / 24 VDC | Solenoid energized while in the non-default state | Fail-safe or fail-secure | "Electrically locked" (fail-safe: powered = locked) vs "electrically unlocked"/storeroom (fail-secure: powered = unlocked, latch always retractable from inside). |
| Electric latch retraction (EL/QEL) on exit device | 12 / 24 VDC | **High inrush** to retract the motor/solenoid, then low hold; "QEL" = quiet/low-current variant | Fail-secure | Size the supply for the **inrush** (can be several amps for a moment) or staggered/soft-start if multiple devices fire together. |

**Inrush vs hold** is the recurring power-sizing gotcha: a supply rated for the
*holding* current of N latch-retraction devices can still brown out if they all
**retract simultaneously** (e.g. a scheduled mass-unlock or a fire-alarm release).
Use a supply rated for aggregate inrush, or devices with sequenced/soft start.

---

## 2. PoE power budgets vs lock draw

PoE delivers power + data on one Cat-5/6 drop, with the budget set by the standard:

| Standard | Marketing name | Power *to the device* (approx) |
| --- | --- | --- |
| IEEE 802.3af | PoE | ~12.95 W |
| IEEE 802.3at | PoE+ | ~25.5 W |
| IEEE 802.3bt (Type 3/4) | PoE++ / 4PPoE | ~51 W / ~71 W |

Implications for an edge/IP door:
- A **reader + electric strike + door contact + REX** on a low-draw fail-secure
  strike can live within PoE+ comfortably — this is the clean "one Cat-6 per door,
  no separate lock supply" case.
- A **continuous-draw maglock** or a **high-inrush EL panic device** can exceed a
  modest PoE budget (continuous load + inrush headroom), so those designs commonly
  **keep a local lock power supply** even when the controller/reader is PoE — i.e.
  "IP at the door" does not automatically mean "lock powered by PoE."
- **Centralizing UPS at the switch/IDF** is the big edge-topology reliability win:
  one managed battery plant in the closet instead of a battery above every door.

**Rule:** size PoE against the **specific device's hold + inrush**, never against a
nominal "it's PoE+." When in doubt, power the lock from a dedicated supply.

---

## 3. Power supply architecture + battery backup

- **Listed access power supply** (UL 294; UL 1076 if intrusion-capable) converts
  mains to 12/24 VDC for locks and sometimes the controller.
- **Multi-output power-controller boards** (e.g. Altronix/LifeSafety) provide N
  independently-configurable outputs, each settable **fail-safe or fail-secure**,
  each able to be **tied to the fire-alarm interface**.
- **Standby battery** (sealed lead-acid, sized for a required hold time) with
  **automatic, glitch-free switchover** on AC loss. *Fail-secure* outputs stay
  powered (locked) on battery; *fail-safe* outputs you may deliberately let drop
  per the egress design.
- **Fire-Alarm Control Panel (FACP) interface:** a supervised input that, on
  alarm, **removes power from designated fail-safe egress locks** (maglocks) for
  free egress. Same board often carries **Form-C dry relays** for **elevator
  recall** and **HVAC/smoke shutdown**.
- **Low-battery, AC-fail, and power-fault** should be reported as **alarms** to the
  head-end — the power plant is part of the monitored system.

---

## 4. REX / RTE — forms and wiring nuance

| REX form | Trigger | Releases lock? | Shunts DPS alarm? | Caveat |
| --- | --- | --- | --- | --- |
| **PIR motion sensor** | Person approaches from secure side | Yes (commonly used to drop a maglock hands-free) | Yes | Fires on *any* motion (passersby/airflow); it is electronics on power, so it cannot be the *sole* life-safety release for a maglock (code still requires power-loss + fire-alarm release). |
| **Lever/mortise REX switch** | Inside lever turned | Often configured to shunt only (door's own latch already retracts) | Yes | Best fidelity: a real hand on the egress hardware; avoids needlessly energizing the lock. |
| **Panic-bar / exit-device switch** | Bar pressed | Shunt and/or release | Yes | The bar mechanically frees the door regardless; switch is for *monitoring/shunt*. |
| **Push-to-exit button (RTE)** | Button pressed | Yes (can be the code release for a maglock) | Yes | When it's the maglock release: 40–48 in AFF, within 5 ft of door, "PUSH TO EXIT," **directly** breaks lock power, holds ≥ 30 s (SKILL.md §3.3 Path A). |

**Why a REX is needed even with a strike:** without a REX, a normal *exit* (someone
turning the inside lever and pushing the door open) opens the monitored door with
no grant → the controller raises a **Door Forced Open** false alarm. The REX tells
the controller "this opening is legitimate" and **shunts** the DPS for the egress.

---

## 5. Door position monitoring

- **Door Position Switch (DPS) / door contact:** magnetic reed (recessed or
  surface) reporting **open/closed**. The input the controller uses to compute
  **Door Forced Open (DFO)** and **Door Held Open (DHO)**.
- **Latch-bolt monitor:** confirms the latch actually threw — *closed* (DPS) does
  not prove *latched/secured*. Important on fire doors and high-security openings.
- **Bond-sensor (maglock):** reports whether the magnet is actually bonded to the
  armature (holding force present) — maglock equivalent of a latch monitor.
- **Shunt time:** on a valid grant or REX, the controller suppresses the DPS alarm
  for a programmed window (e.g. unlock time + a few seconds of pass-through) so a
  normal entry/exit doesn't raise DFO/DHO. Too long = tailgating window; too short
  = nuisance DHO.
- **Request-to-enter / tamper / position** inputs round out a fully-monitored door.

---

## 6. A fully-instrumented door (typical I/O)

```
Controller door point:
  reader port      ── credential in (Wiegand/OSDP)
  lock relay out   ── energize/de-energize the locking device (per fail mode)
  REX input        ── legitimate-exit signal (shunt + maybe release)
  DPS input        ── door open/closed (DFO / DHO computation)
  latch/bond input ── (optional) actually-secured confirmation
  tamper input     ── reader/enclosure tamper
  aux out          ── (optional) local sounder / second relay
  FACP tie         ── fire-alarm drop for fail-safe egress locks (at the supply)
```

A door is only as trustworthy as its monitoring: a lock with **no DPS** can be
propped or defeated with no alarm; a **DPS with no latch monitor** can report
"closed" on a door that never latched.
