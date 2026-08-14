---
name: nfc-ndef-protocols
description: >-
  NFC (a 13.56 MHz HF-RFID subset) and the NDEF data format. TRIGGER: NFC vs HF
  RFID and its tie to ISO/IEC 14443, 18092 (NFCIP-1), FeliCa / JIS X 6319-4,
  15693; NFC Forum specs (Digital Protocol, NCI, LLCP, SNEP); the three modes
  (reader/writer, peer-to-peer + deprecation, card emulation); HCE vs Secure
  Element / eSE and the SE role in payments/access; the NDEF format (records,
  header flags MB/ME/CF/SR/IL, TNF, RTD Text/URI/Smart Poster); NFC Forum Tag
  Types 1-5, chip mappings, capability container, NDEF TLV; phone NFC APIs
  (Android HostApduService / reader mode, Apple Core NFC / background reading);
  tap-to-pair / handover, smart posters, access taps, tag provisioning. SKIP:
  RFID bands & RF physics -> rfid-fundamentals-bands-standards; MIFARE/DESFire
  commands/keys -> contactless-smartcards-mifare-desfire; skim/clone/relay
  attacks -> rfid-nfc-access-attacks-defenses; Wiegand/OSDP wiring ->
  access-credentials-wiegand-osdp; controller/panel PACS architecture; EMV
  payment kernels.
metadata:
  version: 1.0.0
  updated: 2026-06-18
---

# NFC Technology & NDEF

Near Field Communication (NFC) is a **short-range (typically ≤ ~4 cm, design
"operating volume" measured in mm) subset of 13.56 MHz HF RFID** that adds
two-way interaction and a standardized application-data format (NDEF) on top of
contactless-card air interfaces. This skill covers the protocol stack, the three
operating modes, card emulation (HCE vs SE), the NDEF format, the five NFC Forum
Tag Types, and phone NFC APIs at the concept level.

> Scope edge: this is the **NFC/NDEF format & protocol** skill. For RFID band
> physics and the raw air interfaces, the card chip command sets, the attack
> surface, or the wiring/controller layer of access control, defer to the
> sibling skills named in the frontmatter `SKIP:` clause (the PACS-architecture
> sibling there is a planned family member — route by topic if it is absent).

## Core Concepts

### NFC as a slice of HF RFID
- NFC operates at **13.56 MHz**, the same HF carrier as ISO/IEC 14443
  (proximity contactless smartcards) and ISO/IEC 15693 (vicinity cards). It uses
  **inductive (near-field) coupling**, not far-field backscatter (that's UHF
  RFID — see `rfid-fundamentals-bands-standards`).
- What makes NFC *NFC* (vs "just HF RFID"): a converged protocol stack
  (NFC Forum Digital Protocol harmonizing the underlying ISO standards), a
  standard data format (**NDEF**), and **peer/active modes** where a device can
  be initiator *or* target. A passive HF RFID tag is read; an NFC device can
  read, be read, *and* talk peer-to-peer.

### The underlying standards (and the NFC Forum's "technology" names)
The NFC Forum re-labels the legacy air interfaces as **NFC-A / NFC-B / NFC-F /
NFC-V** and unifies them in its **Digital Protocol** spec:

| NFC Forum tech | Underlying standard | Notes |
| --- | --- | --- |
| **NFC-A** | ISO/IEC 14443 Type A | Modified-Miller / 100% ASK; Topaz/MIFARE/NTAG lineage |
| **NFC-B** | ISO/IEC 14443 Type B | NRZ-L / 10% ASK |
| **NFC-F** | JIS X 6319-4 (**FeliCa**), in ISO/IEC 18092 | Manchester / 212–424 kbps; Sony FeliCa |
| **NFC-V** | ISO/IEC 15693 | "Vicinity," longer range, lower data rate |
| **ISO-DEP** | ISO/IEC 14443-4 | Half-duplex block transmission over NFC-A/B; carries ISO/IEC 7816-4 APDUs |
| **NFCIP-1** | ISO/IEC 18092 | Peer-to-peer (active/passive) data exchange protocol |

- **ISO/IEC 14443** = proximity cards (the access-badge / payment-card air
  interface). **ISO/IEC 18092 (NFCIP-1)** = the NFC peer protocol that also
  references FeliCa. **ISO/IEC 15693** = vicinity / NFC-V. Treat exact
  modulation/timing as belonging to `rfid-fundamentals-bands-standards`; here
  we care about which standard each Tag Type / mode maps to.

### NFC Forum governance
The **NFC Forum** is the industry consortium (founded 2004 by Nokia, Philips/NXP,
and Sony) that owns the interoperability specs and certification. Its core stack,
top-down (current versions as listed on the NFC Forum spec page, June 2026):

- **NFC Controller Interface (NCI) 2.3** — host CPU ↔ NFC controller interface.
- **NFC Digital Protocol 2.4** — implementation layer above ISO/IEC 18092 &
  14443; harmonizes NFC-A/B/F and (from 2.0) NFC-V and Active Comm Mode.
- **NFC Activity 2.3** + **Analog 3.0** — how the digital protocol is sequenced;
  RF/analog characteristics (3.0 introduced the 20 mm operating volume).
- **Data Exchange Format (NDEF)** + **RTD** + the **Type 1–5 Tag** specs (see
  below).
- Peer-to-peer link layer: **LLCP 1.4** (OSI L2, based on IEEE 802.2) and
  **SNEP** (exchanges NDEF messages over LLCP). **TNEP 1.0** is the
  tag-based bidirectional NDEF exchange protocol for Type 2/3/4/5 tags.
- **Connection Handover 1.5** (tap-to-pair) and **NFC Authentication Protocol
  1.0** (adopted Dec 2022).

→ Full spec inventory with versions in `references/nfc-forum-specs.md`.

## The Three Operating Modes

NFC devices support three modes, built on ISO/IEC 18092 (NFCIP-1) and ISO/IEC
14443:

1. **Reader/Writer mode:** the device acts as an active reader of NFC Forum
   tags (and writes NDEF to them). This is the "tap a smart poster," "read a
   product tag," "provision a tag" mode. The device powers the passive tag via
   its RF field.

2. **Peer-to-Peer (P2P) mode:** two NFC devices exchange data directly,
   standardized in **ISO/IEC 18092 (NFCIP-1)** with the NFC Forum **LLCP**
   transport and **SNEP** for NDEF exchange (this was Android Beam's substrate).
   - **Deprecation note:** P2P / LLCP-SNEP is effectively **obsolete on phones**.
     **Android Beam was deprecated in Android 10 (2019) and removed in
     Android 14**; **Apple never shipped NFC P2P**. Device-to-device sharing
     migrated to Wi-Fi/Bluetooth-based transfers (e.g., Nearby Share / AirDrop).
     NFC's role in tap-to-share is now mainly **Connection Handover** (NFC
     bootstraps a Bluetooth/Wi-Fi link — "tap-to-pair"), not bulk P2P transfer.

3. **Card Emulation (CE) mode:** the device behaves like a contactless
   smartcard/tag toward an external reader (it is the *target*). This is the
   mode behind tap-to-pay wallets, transit, and mobile access credentials. CE
   splits into **HCE** and **Secure Element** emulation (next section).

## Card Emulation: HCE vs Secure Element (SE / eSE)

In card emulation the phone presents itself over **ISO-DEP (ISO/IEC 14443-4)**
and answers **ISO/IEC 7816-4 APDUs**. The reader begins with a **SELECT AID**
(Application Identifier, up to 16 bytes per ISO/IEC 7816-5) APDU; an **AID
routing table** in the NFC controller decides where subsequent APDUs go.

| Dimension | **Host Card Emulation (HCE)** | **Secure Element (SE / eSE)** |
| --- | --- | --- |
| Where credentials live | Host CPU / app / cloud (no dedicated secure chip) | Dedicated tamper-resistant chip (embedded eSE, UICC/SIM, or microSD) |
| Who handles the APDUs | An Android **`HostApduService`** (or `OffHostApduService` points at the SE) | The secure element itself; the OS/app is not in the transaction path |
| Introduced | **Android 4.4 "KitKat" (2013)** | Pre-HCE; the original NFC payment model |
| Security model | Assume host is compromisable → mitigate with **tokenization, limited-use keys, device fingerprinting, transaction risk analysis** | EMV-grade certified hardware; payment-network evaluated; stores real PAN/keys |
| Trade-off | Faster time-to-market, no SE-issuer/TSM commercial deals, no secure-element dependency | Highest assurance; needed for high-security use (open-loop payments, identity, high-security access) |
| Real-world | **Samsung Pay, Google Pay** use HCE + tokenized cloud credentials | **Apple Pay** stores tokens/keys in the device **Secure Element** |

**Android AID routing specifics** (concepts, not full code):
- Services register **AID groups**; routing is **atomic** per group (all AIDs in
  a group route to the same destination, never a split state).
- Two categories: **`CATEGORY_PAYMENT`** (only one payment AID group active
  system-wide, tied to the default wallet) and **`CATEGORY_OTHER`** (can be
  always active; e.g., loyalty, access, transit).
- Conflict resolution: default wallet wins → else the sole registered service →
  else the OS prompts the user to pick.
- `OffHostApduService` declares AIDs that the controller routes to the SE
  (`android:secureElementName="eSE"` / UICC); Android never binds that service.

**The SE's role in payments/access:** it holds the cryptographic keys and (for
SE-based wallets) the tokenized card, performs the EMV cryptogram in hardware,
and is the trust anchor that lets a phone substitute for a plastic EMV card or a
high-assurance access credential. EMV payment-kernel internals are **out of
scope** — defer; this skill stops at "the SE answers the APDUs / holds the keys."

→ HCE/SE flow, APDU exchange, and the routing table in
`references/card-emulation-hce-se.md`.

## NDEF: the NFC Data Exchange Format

**NDEF** is a lightweight binary **message** format. An **NDEF message** is an
ordered sequence of **NDEF records**; each record carries a typed payload.

### Record header byte (bit 7 → bit 0)

| Bits | Field | Meaning |
| --- | --- | --- |
| 7 | **MB** | Message Begin — set on the first record of a message |
| 6 | **ME** | Message End — set on the last record |
| 5 | **CF** | Chunk Flag — record is one chunk of a chunked payload |
| 4 | **SR** | Short Record — SR=1 ⇒ PAYLOAD LENGTH encoded as **1 byte** (max 255); SR=0 ⇒ **4 bytes**, big-endian |
| 3 | **IL** | ID Length present — an ID field (and its length byte) is included |
| 2–0 | **TNF** | Type Name Format (3 bits) — how to interpret the TYPE field |

### Field order within a record
`header byte` → `TYPE LENGTH (1 byte)` → `PAYLOAD LENGTH (1 byte when SR=1, else
4 bytes big-endian)` → `ID LENGTH (1 byte, only if IL=1)` → `TYPE (var)` →
`ID (var, only if IL)` → `PAYLOAD (var)`.

### TNF values (3-bit field)

| TNF | Name | TYPE field holds |
| --- | --- | --- |
| 0x00 | **Empty** | no type / no payload |
| 0x01 | **NFC Forum Well-Known Type** | an RTD name (e.g., `T`, `U`, `Sp`) |
| 0x02 | **MIME media-type** | a media type per RFC 2046 (e.g., `text/vcard`) |
| 0x03 | **Absolute URI** | an absolute URI identifying the type |
| 0x04 | **NFC Forum External Type** | a namespaced external type (e.g., `android.com:pkg`) |
| 0x05 | **Unknown** | type is unknown (TYPE length 0) |
| 0x06 | **Unchanged** | used for middle/terminating **chunks** |
| 0x07 | **Reserved** | reserved |

### RTD — Record Type Definitions (the well-known types, TNF 0x01)
RTD defines standard record formats. The common ones:

- **Text RTD (`T`, 0x54):** payload = **status byte + language code + text**.
  Status byte: MSB = encoding (0 = UTF-8, 1 = UTF-16); low 6 bits = length of the
  IETF **BCP 47** language tag that follows; remaining bytes are the text.
- **URI RTD (`U`, 0x55):** payload = **1 identifier-code byte + URI string**.
  The first byte is a prefix abbreviation that expands a common scheme to save
  space: `0x00` = no prefix, `0x01` = `http://www.`, `0x02` = `https://www.`,
  `0x03` = `http://`, `0x04` = `https://`, `0x05` = `tel:`, `0x06` = `mailto:`,
  `0x1D` = `file://`. (The full identifier-code table runs to 0x23; see the NFC
  Forum URI RTD spec.)
- **Smart Poster RTD (`Sp`):** a *container* record whose payload is **itself an
  NDEF message**. It bundles exactly **one mandatory URI record** plus optional
  **Title (Text)** records (one per language, no repeats), an **Action** record,
  Icon/Type/Size records. The canonical "tap this poster to open a URL / dial /
  SMS" experience.
- Others: **Multiple URI**, **Signature RTD** (signs records), **Device
  Information RTD**, **Smart/Connection Handover** records.

→ Byte-level walkthroughs of Text/URI/Smart Poster, the full URI prefix table,
and chunking in `references/ndef-format.md`.

## NFC Forum Tag Types 1–5

NDEF can live on a passive tag. The NFC Forum defines a Tag Type as a mapping of
NDEF onto a specific air interface + chip behavior. The NFC Forum **withdrew the
Type 1 Tag specification in Technical Specification Release 2021**, so **Types
2–5 are the actively maintained set**; Type 1 is still widely referenced as the
historical first member and appears below for completeness.

| Type | NFC tech / standard | Representative chips | Typical memory | Notes |
| --- | --- | --- | --- | --- |
| **Type 1** | NFC-A; ISO/IEC 14443A frame + Jewel/**Topaz** commands | Broadcom (Innovision) **Topaz / Jewel** | ~96 B–2 KB | Simple, low-cost; collision-free single-tag; largely legacy (removed in 2021 release) |
| **Type 2** | NFC-A; ISO/IEC 14443A | NXP **MIFARE Ultralight**, **NTAG21x** | ~48 B–~888 B | The dominant cheap NDEF tag (stickers, posters); spec current v1.3 |
| **Type 3** | NFC-F; **FeliCa** / JIS X 6319-4 (ISO/IEC 18092) | Sony **FeliCa Lite** | up to ~1 MB (variable) | Higher data rate; common in Japan; spec current v1.1 |
| **Type 4** | NFC-A **and** NFC-B; **ISO-DEP (ISO/IEC 14443-4)**, **ISO/IEC 7816-4 APDUs** | NXP **DESFire**, SmartMX; HCE on phones | up to ~64 KB+ | APDU-based; the type a phone in CE mode emulates; spec current v1.2 |
| **Type 5** | NFC-V; **ISO/IEC 15693** | NXP **ICODE**, ST **ST25TV** | ~few hundred B–few KB | Vicinity range (longer read distance), lower data rate; spec current v1.3 |

**Capability Container (CC) & NDEF storage.** A tag doesn't just hold raw NDEF —
it advertises how to read/write it via a **Capability Container**:
- **Type 1 / Type 2 / Type 5:** the CC is a small fixed block (e.g., the
  Type 2 4-byte CC: magic `0xE1`, version, memory-size unit, read/write access),
  and NDEF is wrapped in **TLV** structures — the **NDEF Message TLV** (tag
  `0x03`), plus optional **Lock Control TLV** and **Memory Control TLV**, with a
  **Terminator TLV** (`0xFE`).
- **Type 3 (FeliCa):** an **Attribute Information Block** plays the CC role.
- **Type 4:** a **Capability Container file** (selected by APDU) points at the
  **NDEF file**; reading is `SELECT` → `ReadBinary` APDUs.

→ Per-type CC layout, the Type 2 TLV walkthrough, and how a phone formats a blank
tag in `references/tag-types-1-5.md`.

## Phone NFC APIs (conceptual)

**Android** (`android.nfc` / `android.nfc.cardemulation`):
- **Reader/Writer:** `NfcAdapter` + foreground dispatch or **Reader Mode**
  (`enableReaderMode`); `Tag`, `Ndef`, `NdefMessage`, `NdefRecord` to read/write
  NDEF; tech classes `NfcA/NfcB/NfcF/NfcV`, `IsoDep`, `MifareUltralight`, etc.
- **Card emulation:** **`HostApduService`** (HCE; override
  `processCommandApdu` / `onDeactivated`) and **`OffHostApduService`** (routes to
  the SE). AID groups + categories declared in an `apduservice` XML resource;
  `CardEmulation` manages defaults.
- **P2P:** `NdefMessage` push (Android Beam) — **deprecated/removed** (see above);
  do not design new features on it.

**Apple Core NFC** (`CoreNFC`, iPhone 7+; entitlement-gated):
- **`NFCNDEFReaderSession`:** the easy path; read (and on supported devices,
  write) **NDEF** messages from Tag Types 1–5.
- **`NFCTagReaderSession`:** lower-level raw tag access for **ISO 7816**,
  **ISO 15693**, **FeliCa**, and **MIFARE** tags.
- **Background tag reading** (iPhone XS+): the OS reads NDEF tags in the
  background with **no app open** *if* the first record is a **URL record that is
  a valid Apple Universal Link**; it routes the URL to the matching app. Disabled
  while a reader session, Wallet/Apple Pay, or the camera is active, in Airplane
  Mode, or before the first unlock after boot.
- **Card emulation:** Apple historically reserved CE for Apple Pay (Secure
  Element). Third-party HCE-style access opened only via the gated **PassKit
  Semantic Tags / NFC & SE entitlement** on recent iOS — treat availability as
  region/program-gated and verify current entitlement docs before relying on it.

→ Side-by-side API capability matrix in `references/phone-nfc-apis.md`.

## Practical Uses

- **Tap-to-pair (Connection Handover 1.5):** NFC carries the Bluetooth/Wi-Fi
  credentials so a tap bootstraps a higher-bandwidth link (speakers, headphones,
  printers). This — not P2P — is NFC's surviving "tap to share" role.
- **Smart posters:** a `Sp` record (URI + localized Title + Action) behind a
  poster/sticker so a tap opens a URL, dials, or sends an SMS.
- **Access taps:** a phone in **card emulation** (HCE or SE) presents a mobile
  credential to a reader (ISO-DEP/APDU). The *credential data format*, the
  *reader→controller* link, and the *controller/panel architecture* are sibling
  concerns (see SKIP list); here the relevant layer is "CE mode + AID + APDU
  exchange."
- **Tag provisioning:** writing NDEF (often a URI or External/MIME record) onto
  Type 2/5 tags at scale for product authentication, configuration, asset IDs,
  Wi-Fi onboarding, and "tap-to-launch-app" stickers.

**When NFC is the wrong layer.** Bulk or streaming device-to-device transfer →
do *not* use NFC for the payload; use NFC only to carry the pairing credentials
via Connection Handover, then move the data over the Bluetooth/Wi-Fi link it
bootstrapped. No physical tap, or range beyond a few centimetres → NFC cannot
reach; use BLE or far-field UHF RFID instead (band/range trade-offs →
`rfid-fundamentals-bands-standards`). High-assurance credential storage on the
phone → card emulation belongs on the Secure Element, not HCE.

## Uncertainties / verify-before-relying

- **Exact per-chip memory sizes** vary by SKU (e.g., NTAG213/215/216 differ);
  the table gives ranges — confirm the datasheet for a specific part.
- **iOS third-party card emulation / HCE** availability is **entitlement- and
  region-gated** and has changed across iOS releases; verify current Apple
  developer docs (PassKit / NFC & SE entitlement) before designing for it.
- **Type 1 Tag** was de-emphasized in the 2021 spec release; new designs rarely
  target it. Treat it as legacy.
- Sources here are vendor/spec summaries plus the NFC Forum spec index; for
  byte-exact conformance, the **purchased NFC Forum specs** and chip datasheets
  are authoritative over any third-party page cited below.

## Sources

- NFC Forum — Technical Specifications index (NDEF, RTD, Text/URI/Smart Poster,
  Type 1–5 Tag, Digital Protocol 2.4, Activity 2.3, Analog 3.0, NCI 2.3,
  LLCP 1.4, SNEP, TNEP 1.0, Connection Handover 1.5):
  https://nfc-forum.org/build/specifications/
- NFC Data Exchange Format (NDEF) Technical Specification (record/message
  structure, TNF): https://nfc-forum.org/build/specifications/data-exchange-format-ndef-technical-specification/
- NDEF record binary structure — header flags (MB/ME/CF/SR/IL), TNF table, URI
  prefix codes, Text record format (TLV/NDEF breakdown):
  https://tucker.the-twomeys.com/blog/posts/ndef-tlv/
- Android Developers — Host-based card emulation (HCE) overview
  (`HostApduService`, `OffHostApduService`, AID routing, ISO-DEP/7816-4,
  KitKat 4.4): https://developer.android.com/develop/connectivity/nfc/hce
- Apple Developer — Core NFC (`NFCNDEFReaderSession`, `NFCTagReaderSession`,
  Tag Types 1–5, ISO7816/ISO15693/FeliCa/MIFARE):
  https://developer.apple.com/documentation/corenfc
- Apple Developer — Adding Support for Background Tag Reading (Universal Link
  URL record, conditions): https://developer.apple.com/documentation/corenfc/adding-support-for-background-tag-reading
- Secure Technology Alliance — "Host Card Emulation (HCE) 101" whitepaper
  (HCE vs SE, TSM, four security pillars):
  https://www.securetechalliance.org/wp-content/uploads/HCE-101-WP-FINAL-081114-clean.pdf
- Sony — FeliCa and the NFC Forum (NFC-F, ISO/IEC 18092, JIS X 6319-4,
  Type 3 Tag, three modes): https://www.sony.net/Products/felica/NFC/forum.html
- nfcpy documentation — `nfc.tag` (Tag Types 1–5, capability container, NDEF
  TLV, Topaz format): https://nfcpy.readthedocs.io/en/latest/modules/tag.html
- Finextra / FinTech Futures — Secure Element vs cloud-based HCE for NFC
  payments (EMV-grade SE, HCE tokenization, Apple-SE vs Samsung-HCE):
  https://www.finextra.com/blogposting/10087/se-vs-hce-what-is-more-secure-for-nfc-mobile-payments
