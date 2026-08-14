<!-- hub-reference-banner -->
> **Reference file — part of the `physical-access-control` hub.** Formerly the standalone `rfid-fundamentals-bands-standards` skill.
> Sibling topics in this family are now reference files under the hubs (`physical-access-control`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: rfid-fundamentals-bands-standards
description: >-
  RFID fundamentals — system anatomy, frequency bands, coupling physics, and core
  air-interface standards. TRIGGER: choosing/comparing an RFID band (LF 125-134 kHz,
  HF 13.56 MHz, UHF 860-960 MHz, microwave 2.45 GHz); read range / data rate / penetration
  / cost trade-offs; passive vs active vs semi-passive (BAP) tags; inductive coupling vs
  backscatter, near- vs far-field; tag/reader/antenna/backend anatomy; air-interface
  standards (ISO/IEC 14443, 15693, 18000 series, 18000-63 / EPC Gen2 V2, EPCglobal/GS1/RAIN);
  EPC/TID/user memory model; anti-collision (Q-algorithm / slotted ALOHA); which band fits
  access-control/ID (HF) vs supply-chain/inventory (UHF). SKIP: NFC NDEF/record format ->
  nfc-format sibling; MIFARE/DESFire card command & key detail ->
  contactless-smartcards-mifare-desfire; attacks/cloning/skimming/relay ->
  rfid-nfc-access-attacks-defenses; controller/Wiegand-OSDP access-control architecture ->
  PACS sibling; asset-tracking/RTLS depth -> asset-tracking sibling.
version: 1.0.1
updated: 2026-06-18
metadata:
  changelog:
    - "2026-06-18 sko v1.0.0->v1.0.1 — fixed 18000-64/18000-7 mislabel (4 locations), description 1041->985 chars under cap, A/m unit gloss, deduped UHF regional table"
---

# RFID Fundamentals, Frequency Bands & Standards

The vendor-neutral physics-and-standards layer for radio-frequency identification. Use it
to reason about *which band, which coupling mode, and which air-interface standard* a job
needs — before diving into card-product, NFC-format, PACS, or attack detail (those are
sibling skills; see SKIP in the frontmatter).

## Core concepts

### System anatomy (four parts)
1. **Tag / transponder** — a microchip (IC) bonded to an antenna, on a label, card, or
   hard housing. Holds an identifier and optional memory. May be passive, active, or
   semi-passive.
2. **Reader / interrogator** — generates the RF field, runs the air-interface protocol,
   resolves tag collisions, and decodes responses. Can be fixed (portal/gateway) or handheld.
3. **Antenna** — couples energy between reader and tag. On LF/HF it is a wound coil
   (magnetic, near-field); on UHF/microwave it is a radiating element (electric, far-field).
   Often a separate component from the reader on UHF systems.
4. **Backend / host / middleware** — software that filters, deduplicates, and routes tag
   events into business systems (WMS, access control, ERP). The reader "publishes events";
   the backend gives them meaning.

### Power & coupling: inductive (near-field) vs backscatter (far-field)
The dividing line between near-field and far-field is **λ/2π** (wavelength ÷ 2π).

- **Inductive coupling (LF & HF, near-field).** Reader and tag coils share a magnetic
  field; a changing current in the reader coil induces current in the tag coil (mutual
  inductance). The tag answers by **load modulation** — switching a load across its coil on
  and off, which is detectable as sideband frequencies on the reader's antenna. In the
  near-field, magnetic field strength falls as **1/d³** and *available power* as **1/d⁶**, so
  range is short and steeply distance-limited. At 13.56 MHz, λ ≈ 22.1 m and λ/2π ≈ 3.5 m
  (theoretical near-field edge); practical read range is well under 1 m.
- **Backscatter coupling (UHF & microwave, far-field).** The reader radiates an EM wave; the
  tag rectifies part of it to power its chip, then transmits by **modulated backscatter** —
  toggling a load to make its antenna a better or worse reflector, varying the reflected
  signal the reader sees. Far-field power falls only as **1/d²**, enabling multi-metre range.
  A resonant tag antenna reflects strongly when its dimension is ≥ ~½ wavelength.

(Full physics, including the 1/d³ vs 1/d⁶ vs 1/d² derivations and orientation/material
effects, is in `references/frequency-bands.md`.)

### Tag power classes
- **Passive** — no battery; harvests all power from the reader field. Cheapest, smallest,
  effectively unlimited lifetime; shortest range. Dominant in access cards, retail/supply-chain
  labels.
- **Active** — onboard battery powers a transmitter; the tag *initiates* transmission.
  Longest range (tens to hundreds of metres), highest cost, battery-limited life. Typical at
  433 MHz (ISO/IEC 18000-7) and 2.45 GHz; used for real-time location, container/yard tracking.
- **Semi-passive / Battery-Assisted Passive (BAP)** — battery powers the chip (and often
  sensors) but the tag still answers by **backscatter**, not active transmit. Better
  range/sensitivity than passive; used for cold-chain and sensor logging (e.g. temperature).

## Frequency-band comparison

| | **LF** | **HF** | **UHF (passive)** | **Microwave** |
|---|---|---|---|---|
| Frequency | 125 / 134.2 kHz | 13.56 MHz | 860–960 MHz | 2.45 GHz (also 5.8 GHz) |
| Coupling | Inductive (near-field) | Inductive (near-field) | Backscatter (far-field) | Backscatter (far-field) |
| Typical read range | ~≤10 cm (a few in.) | ≤1 m (NFC ~≤10 cm) | ~2–10 m (to ~15–20 ft) | ~1–2 m+ (active much more) |
| Data rate | Lowest | Moderate (106–848 kbit/s on 14443) | High (reads 100s of tags/s) | High |
| Liquid / tissue | Penetrates well | Penetrates well | Absorbed (detunes near water) | Strongly absorbed |
| Metal | Tolerant | Sensitive (detunes) | Detunes (needs on-metal tags) | Very sensitive |
| Tag cost | Higher than HF | Lowest of the three | Low, falling | Higher |
| Common power class | Passive | Passive | Passive / semi-passive | Active / semi-passive |
| Typical use | Animal ID, immobilizer, **access control** | **Access/ID cards**, payment, transit, library, NFC | **Supply chain, inventory, retail** | RTLS, toll (e.g. ETC), yard |

Rule of thumb: **lower frequency → shorter range, slower data, better liquid/metal
penetration, more orientation tolerance; higher frequency → longer range, faster, but
absorbed by water and detuned by metal.** **Access control & ID skew HF; supply chain &
inventory skew UHF.**

UHF is **region-regulated** — the band is not globally identical (US/FCC 902–928 MHz,
Europe/ETSI 865–868 MHz, Japan ~916–921 MHz, each with its own power/channel rules), so design
for the deployment region or buy globally-certified hardware. The full per-region table (power
limits, Listen-Before-Talk) is in `references/frequency-bands.md`. LF and HF (13.56 MHz ISM)
are effectively global.

## Standards map

**Two standards worlds meet at UHF:** GS1/EPCglobal's *EPC Gen2* air interface and ISO's
*18000-63* are the **same protocol** — Gen2 was ratified into ISO as 18000-6C, and Gen2 V2
became ISO/IEC 18000-63. "RAIN RFID" is the industry alliance/brand for this UHF Gen2 family.

| Standard | Band | Role |
|---|---|---|
| **ISO/IEC 14443** (A/B, parts 1–4) | HF 13.56 MHz | **Proximity** (~≤10 cm). 106/212/424/848 kbit/s. Payments (EMVCo), transit, secure access/eID. Basis of NFC Forum Type 2 & Type 4 tags. |
| **ISO/IEC 15693** (parts 1–3) | HF 13.56 MHz | **Vicinity** (~≤1–1.5 m), block-oriented, lower data rate, weaker field (0.15–5 A/m vs 1.5–7.5 A/m for proximity — A/m = magnetic field strength H; higher = stronger). Library, light asset ID. Basis of NFC Forum Type 5 tags. |
| **ISO/IEC 18000-3** | HF 13.56 MHz | HF item-management air interface; Mode 1 aligns with 15693. (Modes are not interoperable.) |
| **ISO/IEC 18000-2** | LF <135 kHz | LF item-management air interface. |
| **ISO/IEC 18000-63** (= EPC Gen2 V2 / 18000-6C) | UHF 860–960 MHz | **The** dominant passive-UHF supply-chain protocol; Q-algorithm anti-collision. |
| **ISO/IEC 18000-7** | 433 MHz | Active-tag air interface (battery-powered, long range). |
| **ISO/IEC 14223** | LF 134.2 kHz | Animal identification (advanced transponders). |
| **EPCglobal / GS1 / RAIN RFID** | UHF | Governs the EPC data standard (the identifier scheme) and the Gen2 air interface; GS1 maintains EPC Tag Data Standard. |
| **NFC Forum / ISO/IEC 18092** | HF 13.56 MHz | NFC is a **subset of HF RFID** for phones (peer-to-peer + card emulation + reader). NFC-A/B = 14443, NFC-V = 15693, NFC-F = FeliCa. (NDEF/record format → NFC-format sibling.) |

(Deeper standards detail — part-by-part scope, version history, NFC technology mapping — is
in `references/standards-map.md`.)

### EPC Gen2 (UHF) tag memory model — four banks
- **Bank 00 — Reserved:** 32-bit **Kill** password (silences the tag permanently) + 32-bit
  **Access** password. The only bank that can be both read- and write-locked.
- **Bank 01 — EPC:** 16-bit CRC-16 + 16-bit Protocol-Control (PC) word + the **EPC** itself
  (96–496 bits) — the item identifier the user programs.
- **Bank 10 — TID:** factory-set, perma-locked **Tag Identifier**. Its leading bits encode
  the **chip mask-designer (vendor) ID** and tag model number, so the TID identifies the
  *chip*, not the item. Globally unique on serialized chips.
- **Bank 11 — User:** optional extra memory (commonly 512 bits, up to a few KB) for sensor
  logs, expiry dates, maintenance data.

(HF cards — 14443/15693 — use a block-oriented memory layout instead; UID + data blocks.
Their sector/key/command detail is the smart-card sibling's domain.)

### Anti-collision (reading many tags at once)
- **UHF Gen2 / 18000-63** uses a probabilistic **slotted-ALOHA "Q-algorithm":** the reader
  broadcasts a Q value; each tag picks a random slot counter in 0…2^Q−1 and replies only when
  its counter hits 0, decrementing otherwise. The reader adapts Q to the estimated tag count —
  too-high Q wastes empty slots, too-low Q causes collisions. This singulates hundreds of tags
  per second.
- **HF 14443** uses a deterministic **bit-collision / binary-tree (anticollision loop)** walk
  over UIDs; **15693** uses a slotted (16-slot) anticollision scheme.

### Chip vendors (high level)
UHF Gen2 ICs are dominated by **Impinj** (Monza/M-series), **NXP** (UCODE), and **Alien**
(Higgs). **NXP** also leads HF/NFC (NTAG, MIFARE, ICODE) and offers LF; **STMicroelectronics**
(ST25), **Infineon**, **EM Microelectronic**, **Texas Instruments**, and **LEGIC** are other
significant players. The TID's mask-designer field maps back to these vendors.

## Common misconceptions
- **"RFID and NFC are different technologies."** NFC is a *subset of 13.56 MHz HF RFID* with
  added phone modes (peer-to-peer, card emulation). All NFC is HF RFID; not all HF RFID is NFC.
- **"EPC Gen2 and ISO 18000-63 are competing standards."** They are the *same* UHF air
  interface; Gen2 V2 *is* ISO/IEC 18000-63 (and Gen2 v1 was 18000-6C). "RAIN RFID" brands it.
- **"UHF reads further, so it's always better."** UHF is absorbed by water/tissue and detuned
  by metal, and its band is region-locked. For tap-to-authenticate access control, HF's short,
  controlled range is a *feature*, not a limitation.
- **"Passive tags have a battery."** Passive tags have no battery — they harvest reader energy.
  Battery-Assisted Passive (semi-passive) tags have a battery but still answer by backscatter,
  not active transmit.
- **"The barcode-style number lives in the chip's serial."** On UHF the user-programmed item
  ID is the **EPC** (Bank 01); the factory **TID** (Bank 10) identifies the *chip/vendor*, not
  the item — they are different banks.
- **"125 kHz prox cards are secure ID."** Band/coupling choice says nothing about security;
  the security model lives in the card product and crypto (smart-card sibling) — and legacy LF
  prox is widely considered weak (attack/defense sibling).
- **"More tag memory = better."** Most deployments only need the ID; user memory adds cost and
  read time. Add it only when on-tag data (sensor logs, etc.) is genuinely required.

## Sources
- GS1 / EPCglobal — EPC Radio-Frequency Identity Gen2 UHF RFID Standard: https://ref.gs1.org/standards/gen2/
- GS1 — Overview of UHF frequency allocations for RAIN RFID: https://www.gs1.org/docs/epc/uhf_regulations.pdf
- ISO/IEC 18000-7:2014 (active 433 MHz air interface): https://www.iso.org/standard/57336.html
- ISO/IEC 15693 (vicinity cards, parts 1–3) — Wikipedia (cites the standard): https://en.wikipedia.org/wiki/ISO/IEC_15693
- Near-field communication (NFC) — Wikipedia: https://en.wikipedia.org/wiki/Near-field_communication
- RFID4U — Inductive and Backscatter Coupling (near/far-field, λ/2π, 1/d³ vs 1/d⁶ vs 1/d²): https://rfid4u.com/inductive-and-backscatter-coupling/
- RFID4U — RFID EPC Gen2 Memory bank layout (Reserved/EPC/TID/User): https://rfid4u.com/rfid-epc-gen2-memory/
- RFID4U — RFID frequency selection: https://rfid4u.com/rfid-frequency/
- Syncotek — HF RFID Standards (14443 / 15693 / 18000-3 / NFC / EMV): https://syncotek.com/hf-rfid-standards/
- Syncotek — What is Backscatter in RFID (modulated backscatter): https://syncotek.com/what-is-backscatter-in-rfid/
- AssetPulse — RFID Frequency Ranges (LF/HF/NFC/UHF specs): https://www.assetpulse.com/blog/rfid-frequency-ranges
- RFIDLabel — What is RAIN RFID? GS1 EPC Gen2 & ISO 18000-63: https://www.rfidlabel.com/what-is-rain-rfid-demystifying-gs1-epc-gen2-and-iso-18000-63-global-standards/
- RFIDLabel — Navigating US (FCC) vs EU (ETSI) UHF frequencies: https://www.rfidlabel.com/navigating-the-global-rfid-maze-a-guide-to-us-vs-eu-uhf-frequencies/
- RFID Journal — Leading RFID companies (chip vendors): https://www.rfidjournal.com/ask-the-experts/what-are-the-leading-rfid-companies/

See `references/frequency-bands.md` and `references/standards-map.md` for depth.
