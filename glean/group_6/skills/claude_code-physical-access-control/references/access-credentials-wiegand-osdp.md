<!-- hub-reference-banner -->
> **Reference file — part of the `physical-access-control` hub.** Formerly the standalone `access-credentials-wiegand-osdp` skill.
> Sibling topics in this family are now reference files under the hubs (`physical-access-control`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: access-credentials-wiegand-osdp
description: >-
  Access-credential data formats and reader-to-controller protocols for physical access
  control. TRIGGER: Wiegand formats (26-bit H10301 = facility code + card number + 2 parity;
  35-bit Corporate 1000; 37-bit H10302/H10304; format collisions); the Wiegand WIRE interface
  (D0/D1, one-way, unsupervised, unencrypted, ~150 m) and why it's insecure; OSDP (SIA, IEC
  60839-11-5) — RS-485 multidrop, CP/PD, polling, command/reply set, LED/buzzer/tamper,
  biometrics, Secure Channel (SCBK, SCBK-D default-key risk, AES-128, MAC); SIA
  Wiegand-to-OSDP migration; HID FORMAT families (Prox 125 kHz; iCLASS/SE/SEOS 13.56 MHz;
  SIO). SKIP: attacks (ESPKey, OSDP SCBK/downgrade) -> rfid-nfc-access-attacks-defenses;
  chip/key internals -> contactless-smartcards-mifare-desfire; RF bands/physics ->
  rfid-fundamentals-bands-standards; NDEF/NFC -> nfc-ndef-protocols; panels/topology ->
  pacs-access-control-architecture; mobile/UWB -> mobile-credentials-aliro-uwb; convergence ->
  physical-security-convergence-standards.
version: 1.1.0
updated: 2026-06-18
metadata:
  changelog:
    - "2026-06-18 greenfield build — Wiegand 26-bit layout, Wiegand-vs-OSDP table, OSDP Secure Channel, HID format families; references/osdp-protocol-detail.md for command/reply set + handshake"
    - "2026-06-18 sko v1.0.0->v1.1.0 — description 1677->998 chars (under 1000 Glean cap, parsed-scalar length); canonicalized H10301 parity wording + capacity reconciliation; explicit CHLNG/CCRYPT-SCRYPT/RMAC_I pairing + distinct-code-space note; moved Prox clonability/tooling to attacks sibling; deduped capacity caveat"
---

# Access Credential Formats & Reader-to-Controller Protocols (Wiegand / OSDP)

The format-and-wire layer of physical access control (PACS): **what bits are on the badge**
(card-data formats) and **how the reader hands those bits to the controller** (the
reader↔controller protocol). This skill owns the credential *data format* and the *Wiegand*
and *OSDP* interfaces. It does **not** cover the chip internals, the radio physics, the panel
hardware, or attack execution — see the SKIP list in the frontmatter and the cross-links below.

Two distinct things are often conflated under the word "Wiegand":

1. **The Wiegand data format** (e.g. 26-bit H10301) — the bit layout of the credential
   number. A "26-bit card" describes the *format*, independent of the wire.
2. **The Wiegand interface** — the physical D0/D1 two-line signalling that carries that bit
   stream from reader to controller. OSDP replaces the *interface* while still transporting
   the same credential *formats* (and more) over an encrypted channel.

---

## 1. Wiegand card-data formats

### 1.1 The 26-bit standard format — H10301

H10301 is the open, non-proprietary, industry-baseline format. 26 bits total:

```
 bit:  1        2 ........... 9    10 ......................... 25   26
      [EP]     [   facility code  ][        card number          ]  [OP]
       │        8 bits (0–255)      16 bits (0–65,535)              │
       │                                                            │
  even parity                                                  odd parity
  over leading 12 data bits                            over trailing 12 data bits
```

| Field | Bits | Width | Range |
| --- | --- | --- | --- |
| Even parity (P) | 1 | 1 | computed (see below) |
| Facility code (FC) / site code | 2–9 | 8 | 0–255 |
| Card number (CN) | 10–25 | 16 | 0–65,535 |
| Odd parity (P) | 26 | 1 | computed (see below) |

- **Card number is transmitted MSB-first (big-endian).**
- **Capacity:** 256 facility codes × 65,535 card numbers = **16,711,425** addressable
  credentials (per identisource; confirmed by Get Safe and Sound). This figure excludes
  card number 0, which issuance conventionally skips; the raw combinatorial maximum is
  256 × 65,536 = 16,711,936.

**Parity calculation — the canonical convention.** The two parity bits are computed over the
24 data bits, each as a simple even/odd 1-count check:

- **Even parity (bit 1):** set so the **leading 12 data bits** (bits 2–13: the 8 facility
  bits plus the top 4 card-number bits) plus bit 1 itself contain an even number of 1s.
- **Odd parity (bit 26):** set so the **trailing 12 data bits** (bits 14–25: the low 12
  card-number bits) plus bit 26 itself contain an odd number of 1s.

The controller validates both parities on every read; a parity failure means a corrupted
read (long cable run, noise, marginal card) and the read is rejected. **Parity is integrity
checking, not security** — it does nothing against deliberate capture/replay.

> Cross-source caution: some references say "even over first 12, odd over last 12" while
> PageMac states "even over bits 1–13, odd over bits 13–25." Both describe the same H10301
> behaviour; the difference is whether the parity bit's own position is counted in the span.
> Validate against the specific controller's format definition before computing by hand.

### 1.2 The facility/site code

The facility code (FC), also called *site code*, partitions the card-number space so two
sites can both issue card #1234 without colliding *as long as their facility codes differ*.
It is the weakest part of the security model: an 8-bit FC has only 256 values, is frequently
left at a vendor default, is printed/known per site, and is trivially included in any
captured bitstream. It is an *organizational* namespace, not an authentication factor.

### 1.3 Other common formats

| Format | Bits | Layout (summary) | Notes |
| --- | --- | --- | --- |
| **H10301 (standard 26-bit)** | 26 | 8-bit FC + 16-bit CN + 2 parity | Open, non-proprietary; smallest namespace; most cloned. |
| **HID Corporate 1000 (35-bit)** | 35 | 12-bit FC ("company ID") + 20-bit CN + 3 parity | Proprietary/managed by HID; far larger namespace. See §1.4. |
| **HID Corporate 1000 (48-bit)** | 48 | larger FC + larger CN + parity | Even larger managed namespace variant. |
| **H10302** | 37 | **no facility code** — 35-bit CN + 2 parity | Pure 35-bit card number, ~34.4 billion values. |
| **H10304** | 37 | 16-bit FC + 19-bit CN + 2 parity | 37-bit *with* facility code. |

### 1.4 HID Corporate 1000 (35-bit) bit layout

Corporate 1000 is a **managed** format: HID controls company-ID assignment so two
organizations never share the same code space. It is **35 bits**: a 12-bit facility code
(HID terms it the *company ID code*), a 20-bit card number, and **three** parity bits.

```
 bit:  1     2     3 ........... 14    15 ....................... 34    35
      [OP]  [EP]  [  company ID    ]  [        card number         ]  [OP]
              │    12 bits (0–4095)    20 bits (0–1,048,575)           │
   odd parity over all 35 bits;  even parity (bit 2) and odd parity (bit 35) over interleaved bit sets
```

Per the PageMac authoritative format reference, the three parity bits are computed in this
order — **bit 2, then bit 35, then bit 1**:

- **Bit 2 (even):** over bits 3,4,6,7,9,10,12,13,15,16,18,19,21,22,24,25,27,28,30,31,33,34
- **Bit 35 (odd):** over bits 2,3,5,6,8,9,11,12,14,15,17,18,20,21,23,24,26,27,29,30,32,33
- **Bit 1 (odd):** over **all** of bits 2–35

The 12-bit company ID (0–4,095) and 20-bit card number (0–1,048,575) give Corporate 1000 a
much larger, centrally-deconflicted namespace than 26-bit. Do **not** hand-derive these
parity masks for production tooling without confirming against HID's own format spec — the
interleaving differs between sources' transcriptions.

### 1.5 Format collisions and ambiguity

- **Same bit length ≠ same format.** A reader/controller must be told which *format* a 26-
  or 37-bit stream is. Two different 37-bit formats (e.g. H10302 vs H10304) map the same 37
  bits to entirely different facility/card values.
- **26-bit is small and openly published**, so card-number/facility-code collisions across
  unrelated deployments are common; the open format is also why 26-bit is the easiest to
  clone and the most often duplicated.
- **Custom/OEM formats** (variable length, embedded issue codes, non-standard parity) exist
  to enlarge the namespace and add a weak "proprietary" hurdle, but length alone never
  guarantees uniqueness — only a managed program (Corporate 1000) or cryptographic
  credential (Seos/SIO) does.
- The credential *number* that the format encodes is what the access-control software stores
  and matches; the format is just the agreed encoding of that number on the card and wire.

---

## 2. The Wiegand interface (the wire) — and why it is a problem

The legacy Wiegand interface predates and is separate from the data format. It carries the
credential bitstream from reader to controller on **two data lines plus ground/power**:

- **D0 (Data-0 / "DATA LOW"):** pulsed low to signal a binary **0**.
- **D1 (Data-1 / "DATA HIGH"):** pulsed low to signal a binary **1**.
- Idle: both lines held high. Each bit = one short pulse (~50 µs) on the appropriate line,
  with a pulse interval (~1–2 ms) between bits; the controller reconstructs the bit stream
  from the pulse train. (Confirm exact pulse/interval timing against the controller spec —
  values vary slightly by vendor and are not part of a single ratified standard.)

Why this is a security and operational problem (the SIA "hole in the boat" position):

1. **One-way / unidirectional.** Data flows reader → controller only. The controller cannot
   query the reader, so it **cannot know the reader's state** — whether it is healthy,
   replaced, or has been physically opened.
2. **Unsupervised.** No line supervision means a cut wire or a swapped reader is not
   detected by the protocol itself.
3. **Unencrypted / "in the clear."** Bits are plain pulses = raw ones and zeros. Anyone who
   taps D0/D1 reads the full credential (facility code + card number) and can replay it.
   (The concrete wiretap/replay attack and the ESPKey implant belong to the security sibling
   — see `rfid-nfc-access-attacks-defenses`; do not reproduce that methodology here.)
4. **Distance-limited.** Practical reader-to-controller runs are capped at roughly **150 m
   / 500 ft**; longer runs corrupt the unbalanced pulse signalling.
5. **Wire-heavy.** Buzzer, LED(s), tamper, hold, etc. each need their own conductors — a
   typical Wiegand reader run is "12+ wires," not 2.
6. **No interoperability guarantee** beyond the basic D0/D1 convention; behaviour around
   LED/beeper control is vendor-specific.

These limitations are exactly what OSDP was designed to remove.

---

## 3. OSDP — Open Supervised Device Protocol

OSDP is the SIA-developed, IEC-standardized reader↔controller protocol that replaces the
Wiegand *interface*. Key facts (SIA / IEC):

- **Steward:** Security Industry Association (SIA). **Standard:** approved by the IEC in
  **May 2020** and published as **IEC 60839-11-5**. Current SIA revision **OSDP v2.2.2**
  (released October 2024).
- **Physical layer:** half-duplex **RS-485** multidrop, 8 data bits / 1 stop bit / no
  parity. Supported baud rates: **9600, 19200, 38400, 115200, 230400** (default 9600;
  changeable via `osdp_COMSET`).
- **Topology:** one **Control Panel (CP)** / Access Control Unit polls one or more
  **Peripheral Devices (PD)** (readers, keypads, I/O modules) on the bus. **Up to 128 PDs**,
  each with an address **0–127**; **0x7F is the broadcast address**. All traffic is
  CP-initiated; a PD answers only when its address is in the packet.
- **Distance:** RS-485 supports cable runs up to **~1,200 m / 4,000 ft** (some vendors
  further) — versus Wiegand's ~150 m / 500 ft.
- **Supervised & bidirectional:** the CP continuously polls, so a missing/dead/tampered
  reader is detected; the CP can push **LED, buzzer/beeper, and text** commands and read
  **tamper/status** back.
- **Richer transport:** carries raw card data, **keypad/PIN** data, and **biometric
  templates**, and supports remote configuration and (in Secure Channel) encrypted smart-
  card data between reader and controller.
- **Wiring:** two wires (RS-485 A/B) instead of 12+, enabling true multidrop.

### 3.1 Command / reply model (essentials)

The CP sends a **command (CMD)**; the addressed PD returns a **reply (REPLY)**. Core set
(full table with hex codes in `references/osdp-protocol-detail.md`):

- **Status/identity:** `osdp_POLL` (periodic heartbeat), `osdp_ID` → `osdp_PDID`,
  `osdp_CAP` → `osdp_PDCAP`, `osdp_LSTAT`/`osdp_ISTAT`/`osdp_OSTAT`/`osdp_RSTAT`.
- **Credential data (PD→CP):** `osdp_RAW` (raw card bit array — this is where the Wiegand
  *format* still lives, now carried inside OSDP), `osdp_KEYPAD` (PIN/keypress data).
- **Output control (CP→PD):** `osdp_OUT`, `osdp_LED`, `osdp_BUZ` (buzzer), `osdp_TEXT`,
  `osdp_COMSET` (set address/baud).
- **Acknowledgements:** `osdp_ACK` (nothing to report), `osdp_NAK` (error), `osdp_BUSY`.
- **Security:** `osdp_KEYSET`; the Secure Channel handshake pairs CP commands with PD
  replies — `osdp_CHLNG` → `osdp_CCRYPT`, `osdp_SCRYPT` → `osdp_RMAC_I` (§3.2).
- **Extension:** `osdp_MFG` (manufacturer-specific).

Command and reply code spaces are distinct, so a command and a reply can share a hex value
(e.g. `osdp_CHLNG` and `osdp_CCRYPT` are both 0x76) without ambiguity — the full hex table is
in `references/osdp-protocol-detail.md`.

### 3.2 OSDP Secure Channel (SC)

Secure Channel is the OSDP feature that closes the Wiegand "in the clear" gap. **AES-128**,
mutual authentication, encrypted+authenticated messages. AES-128 Secure Channel is
**required for U.S. federal government applications**.

- **SCBK — Secure Channel Base Key:** the 128-bit per-PD base key. From it the link derives
  three session keys on each SC setup: **S-ENC** (encryption) and **S-MAC1 / S-MAC2**
  (message authentication).
- **SCBK-D — the default base key:** a key published in the OSDP spec, intended **only** for
  install/provisioning. **Leaving SCBK-D in place in production is the headline OSDP
  misconfiguration** — the key is public, so any party can establish a "secure" channel and
  the encryption gives no protection. Provisioning must replace it with a unique site SCBK.
- **Handshake (mutual challenge-response):**
  1. `osdp_CHLNG` — CP → PD with a CP random number.
  2. `osdp_CCRYPT` — PD → CP with a PD random number **and the client cryptogram** (lets the
     CP authenticate the PD).
  3. `osdp_SCRYPT` — CP → PD with the **server cryptogram** (lets the PD authenticate the CP).
  4. `osdp_RMAC_I` — PD → CP with the initial reply-MAC; both ends now hold matching session
     keys and the channel is up. The reply-MAC chains forward as the IV/seed for
     authenticating subsequent traffic (MAC chaining gives sequencing + tamper-evidence).
- **Install mode & `osdp_KEYSET`:** when no SCBK is set, the PD enters install mode and will
  accept an `osdp_KEYSET` carrying a fresh site SCBK; on success it exits install mode and
  refuses further key-set. The documented residual risk — a physical attacker forcing
  install mode (e.g. a pinhole switch) to inject their own key — is **attack methodology and
  belongs to `rfid-nfc-access-attacks-defenses`**; here, the defensive takeaway is: never
  ship SCBK-D, provision unique per-reader keys, and prefer secure provisioning (e.g. a
  config card) over physically-triggered install mode.

> Deeper protocol detail — full CMD/REPLY hex table, packet framing, addressing/broadcast,
> baud negotiation, and the SC key-derivation/handshake steps — is in
> **`references/osdp-protocol-detail.md`**.

---

## 4. Wiegand interface vs OSDP — comparison

| Dimension | Wiegand interface | OSDP (Secure Channel) |
| --- | --- | --- |
| Steward / standard | De-facto convention; no single ratified protocol | SIA; **IEC 60839-11-5** (2020); SIA OSDP **v2.2.2** (2024) |
| Physical layer | D0/D1 two data lines + ground (pulse signalling) | **RS-485** half-duplex, 2-wire (A/B) |
| Direction | **One-way** reader → controller | **Bidirectional** CP ↔ PD |
| Supervision / tamper | **None** (unsupervised) | Continuous **polling**, status + tamper reporting |
| Encryption | **None** — credential sent in the clear | **AES-128** Secure Channel (mutual auth, MAC) |
| Max distance | **~150 m / 500 ft** | **~1,200 m / 4,000 ft** (vendor-dependent further) |
| Wiring | 12+ conductors typical (data + LED + buzzer + tamper…) | **2 wires**, multidrop |
| Devices per run | One reader per controller port | **Up to 128 PDs**, addresses 0–127 (0x7F broadcast) |
| Baud / data rate | Pulse-rate fixed by reader; no negotiation | **9600 / 19200 / 38400 / 115200 / 230400** |
| LED / beeper / text | Vendor-specific, extra wires | In-protocol (`osdp_LED` / `osdp_BUZ` / `osdp_TEXT`) |
| Carries | Card bitstream only | Card data, **PIN/keypad**, **biometric templates**, config |
| Interoperability | Weak beyond basic D0/D1 | **SIA OSDP Verified** conformance program |
| Replay resistance | None (capture = clone) | Encrypted + MAC'd + session-fresh |

**SIA position:** Wiegand is "a hole in the boat" — its lack of supervision and encryption
makes it unsuitable for new installations; SIA recommends specifying OSDP (and, where a rip-
and-replace is not yet feasible, Wiegand-to-OSDP converters as an interim step).

---

## 5. HID credential families — at the format level

These are the credential *technologies* a TAM/integrator will name; this section is the
*format/security* tier (chip command sets and key diversification belong to
`contactless-smartcards-mifare-desfire`; band/coupling physics to
`rfid-fundamentals-bands-standards`).

| Family | Frequency | Security tier | Format / notes |
| --- | --- | --- | --- |
| **HID Prox** | **125 kHz** (LF) | **Lowest** — no encryption, no mutual auth | Legacy 125 kHz proximity. Carries a Wiegand format (often 26-bit) with no on-card cryptography. Still ubiquitous; the prime migration target. (Clonability and tooling are attack detail → `rfid-nfc-access-attacks-defenses`.) |
| **iCLASS (legacy)** | **13.56 MHz** (HF) | Low-moderate | First-gen HID smart card. Encrypted storage but legacy keying; superseded by SE/Seos for new high-security work. |
| **iCLASS SE** | **13.56 MHz** | Moderate-high | Built on **SIO** (Secure Identity Object) + HID's Trusted Identity Platform; mutual authentication, larger storage; SIO makes the credential more portable/secure than plain iCLASS. |
| **HID Seos (SEOS)** | **13.56 MHz** | **Highest (HID classic line)** | Standards-based secure messaging; the recommended high-security HID credential and the basis for HID mobile credentials. |
| **HID Corporate 1000** | format (any compatible carrier) | namespace control | A **35-bit (or 48-bit) managed FORMAT** (§1.4), not a chip — HID centrally assigns company IDs to prevent cross-customer collisions. Often issued on iCLASS/Seos carriers. |
| **SIO (Secure Identity Object)** | data model | data-protection layer | HID's portable, signed/encrypted credential **data structure** ("PACS data" container) used by iCLASS SE / Seos — a credential methodology, not a frequency or chip. |

**Multi-technology cards** (e.g. iCLASS SE + Prox) deliberately put 125 kHz Prox *and* a
13.56 MHz Seos/SE credential on one card so a site can run readers in mixed mode during a
**Prox → Seos/SE migration** without re-badging everyone at once.

The format/security through-line: **Prox (clonable, plaintext format on the card) → iCLASS →
iCLASS SE / Seos (mutual-auth, SIO-protected)** mirrors the wire-side **Wiegand → OSDP
Secure Channel** progression. A genuinely hardened deployment upgrades *both* the credential
(to Seos/SE) *and* the reader↔controller link (to OSDP Secure Channel) — fixing one while
leaving the other (e.g. Seos cards over a Wiegand wire) leaves an exploitable gap.

---

## Sources

Authoritative and reputable-integrator references consulted (accessed June 2026):

1. **Security Industry Association — Open Supervised Device Protocol (OSDP)** (standard
   overview; IEC 60839-11-5; OSDP v2.2.2; OSDP Verified; AES-128 required for federal):
   https://www.securityindustry.org/industry-standards/open-supervised-device-protocol/
2. **SIA — "There Is a Hole in the Boat: Why Access Control Professionals Need to Move From
   Wiegand to OSDP"** (Wiegand one-way/unsupervised/in-the-clear; 500 ft; OSDP 4,000 ft;
   migration position): https://www.securityindustry.org/2021/11/09/there-is-a-hole-in-the-boat-why-access-control-professionals-need-to-move-from-wiegand-to-osdp/
3. **Axis Communications — "Why choose OSDP over Wiegand in access control" (white paper,
   May 2025)** (Wiegand vs OSDP comparison; 150 m vs 1,200 m; biometric/mobile transport;
   IEC 60839-11-5): https://whitepapers.axis.com/en-us/osdp-protocol-in-access-control
4. **PageMac — HID RFID data formats** (exact 26-bit H10301 and 35-bit Corporate 1000 bit
   layouts and parity masks): https://www.pagemac.com/projects/rfid/hid_data_formats
5. **identisource — The Standard 26-bit Format** (H10301 field widths, FC/CN ranges,
   capacity): https://www.identisource.net/26_bit_format_layout.cfm
6. **Get Safe and Sound — 26-Bit Wiegand Format** (parity halves, 16,711,425 capacity,
   facility-code role): https://getsafeandsound.com/blog/26-bit-wiegand-format/
7. **libosdp / doc.osdp.dev — OSDP protocol docs** (CP/PD model, RS-485, AES-128, command/
   reply codes, Secure Channel base key & enforce-secure flag):
   https://doc.osdp.dev/index.html and https://doc.osdp.dev/protocol/commands-and-replies.html
8. **libosdp / doc.osdp.dev — Secure Channel** (SCBK / SCBK-D default-key risk; S-ENC,
   S-MAC1/2 session keys; CHLNG/CCRYPT/SCRYPT/RMAC_I handshake; install mode + KEYSET):
   https://doc.osdp.dev/libosdp/secure-channel.html
9. **libosdp introduction / OSDP World / Suprema BioStar** (baud rates 9600–230400; up to
   128 PDs; address 0–127; 0x7F broadcast; default addr 0 / 9600):
   https://libosdp.gotomain.io/protocol/introduction.html and https://osdpworld.com/2021/04/03/osdp-device-communication/
10. **HID Global — iCLASS SE / Seos / SIO and Corporate 1000** (13.56 MHz SIO-based families;
    Prox 125 kHz; multi-tech migration; managed Corporate 1000 namespace):
    https://www.hidglobal.com/product-mix/seos , https://www3.hidglobal.com/products/cards-and-credentials/iclass-se , and https://www.hidglobal.com/solutions/corporate-1000
11. **Telaeris / Get Kisi — HID security card types** (Prox vs iCLASS vs SE vs Seos security
    tiers; clonability of 125 kHz): https://telaeris.com/understanding-hid-security-card-types
    and https://www.getkisi.com/blog/hid-card-types

**Uncertainty / verify-before-use notes:**
- Wiegand pulse width / inter-pulse interval (~50 µs / ~1–2 ms) is a common practical
  convention, **not** a single ratified figure — confirm against the specific controller.
- The H10301 parity *span* is described two ways across sources ("12 data bits" vs "13 bits
  incl. parity"); the resulting parity bit is the same (see §1.1 and the §1.1 cross-source
  caution). The Corporate 1000 parity masks are transcribed from PageMac — **do not
  hand-implement** a format calculator without confirming against HID's own published format
  definition.
