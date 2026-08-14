# APDU framing, status words & the byte-level authentication flow

Deep reference for the **`contactless-smartcards-mifare-desfire`** skill. The
SKILL.md body covers the model at working depth; this file holds the byte-level
APDU framing, status-word tables, the activation sequence, and the DESFire
command codes you need when actually parsing or building card traffic. Values are
from NXP datasheets/app notes and ISO/IEC 14443 & 7816-4 (sources are listed in
the parent SKILL.md). Items flagged `[verify]` need confirmation against the
primary PDF before being quoted verbatim.

---

## 1. ISO/IEC 14443 activation sequence (Type A)

Order of operations from field-on to APDU exchange on a Type A 14443-4 PICC
(e.g., DESFire / Plus):

1. **REQA** (`0x26`, 7-bit short frame) or **WUPA** (`0x52`, wakes a HALTed card)
   — reader polls.
2. **ATQA** (2 bytes): PICC answers. Bits b7–b8 *hint* UID size
   (`00`=4-byte/single, `01`=7-byte/double, `10`=10-byte/triple) but **do not
   rely on ATQA** for chip/UID-size identification (collisions corrupt it; use
   SAK — NXP AN10833).
3. **Anticollision loop**, per **cascade level**:
   - **SEL** byte = `0x93` (CL1) / `0x95` (CL2) / `0x97` (CL3).
   - **Cascade Tag CT = `0x88`** prefixes the partial UID when more bytes remain
     (a 7-byte UID = `CT ‖ UID0..2` at CL1, then `UID3..6` at CL2).
   - PICC returns 4 UID bytes + **BCC** (block check char = XOR of the 4).
   - Reader sends **SELECT** (`SEL ‖ NVB=0x70 ‖ full-UID ‖ BCC`) → PICC returns
     **SAK** (1 byte).
4. **SAK** decode (ISO 14443-3 Table 8 / AN10833):
   - **bit b3 = cascade bit** — set → UID *not complete*, raise cascade level and
     loop; clear → UID complete.
   - **bit b6** — set → PICC is **ISO/IEC 14443-4 compliant** (supports RATS/T=CL
     → speaks APDUs); clear (with b3=0) → not 14443-4 (e.g., MIFARE Classic).
   - Common SAK: `0x20` = UID complete + 14443-4; `0x08` = MIFARE Classic 1K
     (not 14443-4); `0x04` = UID incomplete/cascade `[verify per chip]`.
5. **RATS** (only if b6=1): start byte `0xE0`, param = `(FSDI<<4) | CID`. PICC
   replies **ATS**.
6. **ATS** fields: **TL** (length) → **T0** (format; low nibble = **FSCI** frame
   size, flags which interface bytes follow) → optional **TA(1)** (bit-rate
   capability), **TB(1)** (timing: **FWI** frame-waiting integer + **SFGI**
   start-up guard), **TC(1)** (protocol: NAD-supported, CID-supported) →
   historical bytes (≤15).

**Type B** differs: **REQB/WUPB** (carry AFI + slot count) → **ATQB** (= **PUPI**
4 B + app data 4 B + protocol info 3 B) → **ATTRIB** (selects one PICC by PUPI,
assigns params + **CID**). There is **no RATS/ATS on Type B** — ATQB+ATTRIB carry
the equivalent. **PUPI** is the Type B anticollision identifier (the analogue of
the Type A UID).

**Frame-size code mapping** (FSDI→FSD reader-receive; FSCI→FSC card-receive —
identical table):

| code | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9–15 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| **bytes** | 16 | 24 | **32** | 40 | 48 | 64 | 96 | 128 | **256** | RFU(→256) |

Default **FSCI = 2 → 32 bytes**. **FWT** (frame waiting time) `= 256 × 16/fc ×
2^FWI`, FWI 0–14; an **S(WTX)** block temporarily extends it.

---

## 2. T=CL block protocol (ISO/IEC 14443-4)

Half-duplex block protocol carrying the APDUs. Each block opens with a **PCB**
(Protocol Control Byte); three block types:

- **I-block** (Information): carries the application payload (the APDU). Holds a
  1-bit **block number** (toggles 0/1 for sequencing) and a **chaining bit** (set
  → more I-blocks follow). PCB pattern `0b000x xy1z` `[verify exact bits]`.
- **R-block** (Receive-ready) — ACK/NAK + block number; no INF field. Drives
  chaining flow-control and error recovery. PCB pattern `0b101x xy0z` `[verify]`.
- **S-block** (Supervisory) — control. **S(WTX)** (PICC asks for more time on a
  long op, e.g. a crypto compute) and **S(DESELECT)** (drop to HALT). e.g.
  `0xC2` = S(WTX), `0xC0`/`0xCA` = S(DESELECT) (the `0x08` bit sets when a CID is
  appended) `[verify]`.

Optional prologue fields: **CID** (Card IDentifier 0–14 — lets a PCD juggle
several active PICCs) and **NAD** (Node ADdress for logical connections).

---

## 3. ISO/IEC 7816-4 APDU structure

**Command APDU (C-APDU)** = mandatory header `CLA INS P1 P2` + conditional body:

| Field | Bytes | Meaning |
| --- | --- | --- |
| **CLA** | 1 | Class (interindustry vs proprietary; SM bits b4–b3; logical-channel bits b1–b2). |
| **INS** | 1 | Instruction code. |
| **P1 P2** | 2 | Parameters (offset/selection mode; `00` if unused). |
| **Lc** | 0/1/3 | Length of command data (Nc). 3-byte form = `00` marker + 2-byte len (extended). |
| **Data** | Nc | Command data. |
| **Le** | 0/1/2/3 | Expected response length (Ne); short `00`→256, extended 2-byte after `00` marker. |

**The four ISO cases:**

| Case | Cmd data (Lc) | Resp data (Le) | Body |
| --- | --- | --- | --- |
| 1 | — | — | `CLA INS P1 P2` only |
| 2 | — | yes | `Le` only |
| 3 | yes | — | `Lc` + data |
| 4 | yes | yes | `Lc` + data + `Le` |

**Short** length: max 256 B each way. **Extended**: C-APDU data up to 65,535 B,
R-APDU data up to 65,536 B (contactless cards commonly support extended; contact
T=0 needs ENVELOPE/GET RESPONSE chaining).

**Response APDU (R-APDU)** = `[response data] SW1 SW2` (2 mandatory status bytes).

### Status words (SW1 SW2)

| SW | Meaning |
| --- | --- |
| `9000` | Success / normal completion |
| `61xx` | Success; `xx` more bytes available → issue **GET RESPONSE** (T=0) |
| `6Cxx` | Wrong Le; `xx` = exact length → re-issue with Le=xx |
| `6700` | Wrong length (Lc/Le) |
| `6A82` | File or application **not found** |
| `6A86` | Incorrect P1-P2 |
| `6982` | **Security status not satisfied** (need prior auth) |
| `6983` | Authentication method blocked |
| `6985` | Conditions of use not satisfied |
| `6986` | Command not allowed (no current EF) |
| `63Cx` | Warning — verification failed; `x` = remaining retries |
| `6D00` | INS not supported |
| `6E00` | CLA not supported |
| `6F00` | No precise diagnosis |

Ranges: `61xx`/`9xxx` normal; `62xx`/`63xx` warning; `64xx`/`65xx` execution
error; `67`–`6F` checking error.

### Interindustry commands (INS)

| Command | INS | Notes |
| --- | --- | --- |
| **SELECT** | `A4` | P1=`04` by AID/DF name; P1=`00` by file ID; P2=`0C` no-response/first, `00` return FCI. |
| **READ BINARY** | `B0` | Transparent EF; offset in P1-P2 (or short-EF id in P1 if b8=1). |
| **UPDATE BINARY** | `D6` | Write transparent EF. |
| **READ RECORD** | `B2` | Record-structured EF. |
| **UPDATE RECORD** | `DC` | Record-structured EF. |
| **GET RESPONSE** | `C0` | Retrieve pending data after `61xx` (T=0). |
| **VERIFY** | `20` | PIN/CHV verification (→`63Cx` on fail). |
| **GET CHALLENGE** | `84` | Random for auth. |
| **INTERNAL / EXTERNAL AUTHENTICATE** | `88` / `82` | Mutual auth. |
| **ENVELOPE** | `C2` | Wrap long/SM APDUs over T=0. |

### File model (MF/DF/EF)

- **MF** = Master File (root DF), file ID **`3F00`**.
- **DF** = Dedicated File (directory; an application = a DF selectable by **AID**).
- **EF** = Elementary File (data): **transparent** (binary, byte-addressable),
  **linear fixed/variable**, **cyclic fixed** (record-oriented).
- File ID = 2 bytes; short-EF ID = 5 bits (1–30); record numbers `01`–`FE`.
- **AID** = up to **16 bytes** = **RID** (5-byte registered provider ID) + **PIX**
  (0–11 byte proprietary extension).

### Secure messaging & logical channels (CLA byte)

- **Logical channels:** CLA b1–b2 = channels 0–3 (basic); channel 0 always open.
- **Secure Messaging (SM):** CLA b4–b3 signal SM (b4=1 = SM). SM wraps the body in
  **BER-TLV** objects — cryptographic checksum/MAC (tag `8E`), encrypted data
  (`81`/`87`), protected Le (`97`) — integrity-protecting and/or encrypting the
  APDU so the air gap can't be sniffed/tampered. (7816-4 Annex F.)

---

## 4. ATR vs ATS vs ATQA/SAK

| Term | Layer | What it is |
| --- | --- | --- |
| **ATR** (Answer To Reset) | contact, 7816-3 | Byte string a *contact* card emits on electrical reset. **No ATR on the air interface.** |
| **ATQA / ATQB** | contactless, 14443-**3** | Polling response; presence + UID-size hint (A) / PUPI+protocol (B). |
| **SAK** | contactless, 14443-**3** | Post-SELECT byte; UID-complete (b3) + 14443-4 compliant (b6). |
| **ATS** (Answer To Select) | contactless, 14443-**4** | The contactless analogue of the ATR; frame size, bit rates, timing, CID/NAD support, historical bytes. (Type B carries the equivalent in ATQB+ATTRIB.) |

PC/SC middleware synthesizes a **pseudo-ATR** for non-14443-4 chips
(Classic/Ultralight) so host software sees a uniform ATR-like descriptor.

---

## 5. DESFire native-wrapped APDU framing

DESFire commands ride over T=CL in one of two framings — **native-wrapped** or
**true ISO 7816-4** — and **the first command of a session fixes the framing**
(styles cannot be mixed).

**Native-wrapped layout:**

```
Command:   CLA   INS   P1   P2   Lc    [Data]   Le
           0x90  cc    00   00   len   data     00
             │    │
   class 0x90 │   native DESFire command code as INS (e.g. 0x60 GetVersion)
   marks a    │
   wrapped    │
   native APDU┘

Response:  [Data]  SW1  SW2
                   0x91  ss        ss = DESFire native status
```

- **CLA = `0x90`**; **INS = DESFire command code**; **P1=P2=`0x00`**; **Lc** =
  data length; **Le = `0x00`** for wrapped commands.
- **Status = `91 ss`**: **`9100` = OPERATION_OK**, **`91AF` = ADDITIONAL_FRAME**
  (more data follows — chaining). Other `91xx` are specific errors `[verify the
  full table against the datasheet — only 00/AF are high-confidence]`.
- **Chaining example (GetVersion):** `90 60 00 00 00` → `… 91 AF`, then
  `90 AF 00 00 00` repeatedly to fetch each subsequent frame until `9100`.

**True ISO 7816-4 framing:** standard INS bytes (`A4` SELECT, `B0` READ
BINARY…), ISO file model, ISO status words (`9000`/`61xx`/…).

### DESFire command codes (native)

| Command | Code | Use |
| --- | --- | --- |
| `Authenticate` (DES / 2K3DES) | `0x0A` | legacy auth |
| `AuthenticateISO` (3K3DES) | `0x1A` | ISO-mode auth |
| `AuthenticateAES` (AES-128) | `0xAA` | AES auth |
| `AuthenticateEV2First` | `0x71` | EV2/EV3 secure-messaging start (also LRP-mode on EV3) |
| `AuthenticateEV2NonFirst` | `0x77` | subsequent auth in an EV2/EV3 session |
| `GetVersion` | `0x60` | chip/version (multi-frame) |
| `FormatPICC` | `0xFC` | wipe all apps/files (master-key-gated) |
| `Read_Sig` | `0x3C` | read ECC originality signature `[verify code]` |

---

## 6. The three-pass mutual authentication, byte-level

Conceptual flow for DESFire-style auth with shared key **K** (the same pattern
underlies ISO/IEC 9798-2 three-pass mutual authentication — neither side ever
transmits K):

1. **Reader → card:** `Authenticate(keyNo)`. **Card → reader:** card picks nonce
   **RndB**, returns **`E_K(RndB)`** (enciphered — unlike MIFARE Classic, whose
   first challenge is plaintext).
2. **Reader → card:** reader decrypts to recover RndB, picks nonce **RndA**,
   rotates RndB left by one byte → **RndB'**, returns **`E_K(RndA ‖ RndB')`**. The
   card decrypts and checks RndB' → **reader proven** to know K.
3. **Card → reader:** card rotates RndA → **RndA'**, returns **`E_K(RndA')`**. The
   reader checks RndA' → **card proven** to know K.

**Session-key derivation:** both sides now hold **RndA and RndB** and concatenate
selected bytes of each (the exact byte selection differs by crypto type — DES vs
AES vs LRP) into the **session key**, which then keys **CMAC** (integrity) and
**CBC** (encryption) for the rest of the session. Per-file communication mode
(**Plain / MACed / Fully-Enciphered**) decides which of those apply to each
file's data.

**EV2 SecureMessaging** (entered via `AuthenticateEV2First`/`NonFirst`) uses
CMAC-chained command/response counters + AES-CBC encryption with the session
keys. **EV3 LRP** (Leakage-Resilient Primitive, NXP AN12304) is an alternative
secure-messaging mode built only from standard AES-128 (no proprietary crypto),
designed to resist side-channel/fault analysis; it is selectable in place of the
AES-CBC mode and is driven by the same `0x71`/`0x77` codes reinterpreted in LRP
mode.

---

## 7. Key diversification (NXP AN10922) — byte-level

CMAC-based one-way derivation of a **unique per-card key** from a master key.

- **Diversification input M** = unique data, recommended **card UID ‖ AID ‖
  system identifier**. Worked AN10922 example:
  `UID=04782E21801D80`, `AID=3042F5`, `SystemIdentifier=4E585020416275` →
  `M = 04 78 2E 21 80 1D 80 30 42 F5 4E 58 50 20 41 62 75`. Max **M = 31 bytes**
  for AES.
- **Compute:** `CMAC(K, DIVconst ‖ M)` under the master key K, using the CMAC
  sub-key K1 (M fills the block) or **K2** (M padded with `0x80` then `0x00…`).
  Leading **DIV constant** by key type:

| Key type | Master | Leading DIV constant(s) | Output |
| --- | --- | --- | --- |
| AES-128 | 16 B | `0x01` | 16 B (1 key load, 3 AES ops) |
| AES-192 | 24 B | `0x11`, `0x12` | 24 B (two CMACs, recombined) |
| 2K3DES (2TDEA) | 16 B | `0x21`, `0x22` | 16 B |
| 3K3DES (3TDEA) | 24 B | `0x31`, `0x32`, `0x33` | 24 B |

- For 2TDEA/3TDEA, the **LSB of each 8-byte block is the DES key-version bit** —
  copy the master key's version bits into the diversified key.
- **AN10922 Rev. 2.2 (2019)** extends the same construction to **AES-256**
  (`[verify the AES-256 DIV constant in Rev. 2.2 — not present in the Rev. 2.0
  text]`).
- The terminal (ideally a **SAM**) re-derives the per-card key at transaction
  time; the **master key never leaves the SAM/HSM**.

---

## Open items to confirm against primary PDFs

1. Exact PCB bit layouts for I/R/S blocks and the chip-specific SAK value `0x04`.
2. The full DESFire `0x91xx` status-code table (only `9100`/`91AF` are
   high-confidence here) and `Read_Sig` code `0x3C`.
3. The AES-256 leading DIV constant in AN10922 Rev. 2.2.
4. EV2/EV3 exact maximum-application wording (EV1's 28-app cap is lifted to
   memory-limited).
