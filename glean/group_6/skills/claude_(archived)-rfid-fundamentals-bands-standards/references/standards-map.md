# RFID Air-Interface Standards Map — Deep Reference

Companion to `../SKILL.md`. Part-by-part scope of the core air-interface standards, the
EPC/ISO convergence at UHF, and the NFC-vs-HF-RFID relationship. Figures are grounded in the
sources in SKILL.md; mark uncertain items as such and verify against the ISO/GS1 text before
relying on exact clause numbers (the official standards are paywalled — do not invent clause
detail).

## The big picture: who owns what

- **ISO/IEC JTC 1/SC 31** maintains the vendor-neutral air-interface standards: the
  **ISO/IEC 18000 series** (item management, organized by band) plus **ISO/IEC 14443** and
  **ISO/IEC 15693** (HF cards).
- **GS1 / EPCglobal** owns the **EPC** data standard (the *identifier* scheme — what the bits
  mean) and the **Gen2** UHF air interface. GS1 contributed Gen2 to ISO, where it became part
  of 18000-6/-63.
- **RAIN RFID** is an industry alliance and the marketing umbrella for the passive-UHF Gen2 /
  18000-63 ecosystem (analogous to how "Wi-Fi" brands IEEE 802.11). It is *not* a separate
  air interface.
- **NFC Forum** maintains NFC device and tag-type specifications on top of the HF (13.56 MHz)
  standards.

## ISO/IEC 18000 series (item management, by band)

| Part | Band / freq | Scope |
|---|---|---|
| 18000-1 | — | Reference architecture / parameter definitions for the series. |
| 18000-2 | LF < 135 kHz | LF air interface. |
| 18000-3 | HF 13.56 MHz | HF air interface. **Mode 1** aligns with ISO/IEC 15693; later editions define additional modes (e.g. Mode 2, Mode 3). Modes are *not* interoperable. |
| 18000-4 | Microwave 2.45 GHz | 2.45 GHz air interface. |
| 18000-6 | UHF 860–960 MHz | UHF air interface; historically had types A/B/C. **Type C = EPC Gen2** (a.k.a. 18000-6C). The 2013 edition's Types A/B/C/D were later split into the standalone parts 18000-61/-62/-63/-64. |
| **18000-63** | UHF 860–960 MHz | The current passive-UHF part (formerly 18000-6 Type C); corresponds to **EPC Gen2 V2** (evolved from 18000-6C). The dominant supply-chain protocol. |
| 18000-64 | UHF 860–960 MHz | Passive-UHF **Type D** (formerly 18000-6 Type D) — a minor, little-used UHF type. *Not* the active band. |
| **18000-7** | 433 MHz | Active-tag air interface (battery-powered, long range, hundreds of metres). A separate part — never renumbered. |

> Numbering note: the four UHF types of the old 18000-6 (Types A/B/C/D) were split into the
> standalone parts **18000-61 / -62 / -63 / -64** — all UHF 860–960 MHz. The familiar
> EPC Gen2 air interface is **18000-6 Type C ≈ 18000-63** (and Gen2 V2 ≈ 18000-63). The
> active 433 MHz standard is and remains **ISO/IEC 18000-7** — it was *not* renumbered to
> 18000-64 (18000-64 is the UHF Type D part). You will see both "18000-6C" and "18000-63" for
> the passive-UHF protocol in the wild; treat 18000-6C ≈ Gen2 v1 and 18000-63 ≈ Gen2 V2 as the
> working mapping (verify the precise edition against ISO if a contract depends on it).

## EPC Gen2 (the UHF workhorse)

- **Identity:** EPC Gen2 is GS1/EPCglobal's UHF Class-1 Generation-2 air-interface protocol.
  Ratified into ISO as **18000-6C** (Gen2 v1); **Gen2 V2** became **ISO/IEC 18000-63**.
- **Version history (working summary):**
  - **Gen2 v1 (2004–2008):** original Class-1 Gen-2; → ISO 18000-6C.
  - **Gen2 V2 (≈2013–2015):** added optional security/privacy features — untraceability,
    richer access control, anti-counterfeiting (crypto-suite hooks), better loss-prevention. →
    ISO/IEC 18000-63.
  - Later maintenance revisions exist (point releases); confirm the exact current edition on
    ref.gs1.org rather than asserting a specific "v2.1/v3" here — *uncertain, verify*.
- **Memory model (four banks):**
  - **Bank 00 Reserved** — 32-bit Kill password (bits 00h–1Fh) + 32-bit Access password (bits
    20h–3Fh). Only bank that can be read- *and* write-locked.
  - **Bank 01 EPC** — 16-bit CRC-16, 16-bit Protocol-Control (PC) word, then the EPC
    (96–496 bits). The PC word encodes EPC length, a user-memory indicator, an XPC indicator,
    and a toggle (EPCglobal vs ISO/IEC 15961 data).
  - **Bank 10 TID** — factory-programmed, usually perma-locked. Leading bits identify the
    **chip mask-designer (vendor)** and **tag model number**; serialized chips carry a globally
    unique TID. Identifies the *chip*, not the item.
  - **Bank 11 User** — optional; commonly 512 bits, up to a few KB on some chips; for sensor
    logs, expiry/maintenance data.
- **Anti-collision:** slotted-ALOHA **Q-algorithm** (see SKILL.md). Reader broadcasts Q; tags
  pick a counter in 0…2^Q−1, reply at 0, decrement otherwise; reader adapts Q to tag density.
- **EPC identifier:** the data carried in Bank 01 follows the **GS1 EPC Tag Data Standard**
  (e.g. SGTIN, SSCC encodings) — i.e. a serialized GTIN, not a raw barcode number. (Encoding
  detail is a GS1-data topic, lighter-touch here.)

## HF card standards (13.56 MHz)

### ISO/IEC 14443 — proximity cards (PICC/PCD)
- Range ~≤10 cm; needs a stronger field (1.5–7.5 A/m). Types **A** and **B** differ in
  modulation/coding/anticollision. Data rates **106 / 212 / 424 / 848 kbit/s**.
- Four parts: -1 physical, -2 RF power & signal, -3 initialization & **anticollision**
  (bit-collision binary-tree walk over UIDs), -4 transmission protocol (**ISO-DEP**, the APDU
  transport used by payment/eID).
- Uses: contactless **payments (EMVCo)**, transit fare media, secure access/eID, e-passport
  chips. Card products built on it (MIFARE, DESFire, etc.) → smart-card sibling.

### ISO/IEC 15693 — vicinity cards (VICC/VCD)
- Range ~≤1–1.5 m; weaker field (0.15–5 A/m); block/memory-oriented; lower data rate
  (~26.48 kbit/s high, 6.62/6.67 kbit/s low, down to ~1.65 kbit/s in 1-of-256 mode).
- Three parts: -1 physical (2018), -2 air interface & initialization (2019), -3 anticollision &
  transmission protocol (16-slot scheme; latest edition 2026 per Wikipedia — *verify*).
- Uses: libraries, light asset ID, document/laundry tracking, some access.

## NFC vs HF RFID (the boundary this skill owns)

NFC is a **subset of 13.56 MHz HF RFID** that adds phone-centric modes. This skill covers
*that NFC sits on the HF RFID standards*; the **NDEF record/message format and tag-type
content** are the NFC-format sibling's domain.

- **ISO/IEC 18092 (NFCIP-1)** defines NFC device modes at 13.56 MHz (active + passive comms,
  peer-to-peer). NFC devices also read/emulate 14443 and 15693 tags and FeliCa.
- **NFC technology mapping:**
  - **NFC-A** ← ISO/IEC 14443 Type A
  - **NFC-B** ← ISO/IEC 14443 Type B
  - **NFC-F** ← FeliCa (JIS X 6319-4)
  - **NFC-V** ← ISO/IEC 15693
- **NFC Forum tag types ↔ base standard:** Type 1 (Topaz/ISO 14443A-derived), **Type 2**
  (ISO/IEC 14443A, e.g. NTAG), Type 3 (FeliCa), **Type 4** (ISO/IEC 14443A/B), **Type 5**
  (ISO/IEC 15693). Type 5 was added by the NFC Forum in 2015, letting phones read 15693
  vicinity tags (e.g. NXP ICODE SLIX) — slightly longer phone read range than Type 2.

## Adjacent / related standards (pointers, not depth)
- **ISO 11784 / 11785 & ISO/IEC 14223** — LF animal identification (code structure + advanced
  transponders).
- **ISO/IEC 15961 / 15962 / 15963** — RFID *data* protocol, encoding rules, and unique tag
  identification (data layer above the air interface).
- **EPCIS / GS1** — event/data-sharing layer above the tag (supply-chain visibility); out of
  scope here beyond the EPC identifier note.
- **EMVCo** — contactless payment kernels riding on ISO/IEC 14443 (card-product/payment layer).

## What to verify before quoting clause-level detail
The authoritative ISO and GS1 documents are paywalled. The band/standard *mappings*, memory
banks, frequency allocations, and version lineage above are well-corroborated by the cited
public sources. **Do not invent** clause numbers, exact bit offsets beyond the documented
password/CRC/PC layout, or specific Gen2 point-release version numbers; pull those from
ref.gs1.org or the ISO catalogue when precision matters.
