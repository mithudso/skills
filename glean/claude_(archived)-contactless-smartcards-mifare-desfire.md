# contactless-smartcards-mifare-desfire

**Category:** Science, Biology & Medicine
**Platform:** Claude (Archived)
**Original Path:** claude/archived-treearch/contactless-smartcards-mifare-desfire

## Description
Contactless smart cards and on-card crypto: ISO/IEC 14443 Type A/B PICCs, the ISO/IEC 7816-4 APDU model, and the NXP MIFARE family (Classic, Ultralight/C/EV1, Plus, DESFire EV1/EV2/EV3). TRIGGER: MIFARE/DESFire tier choice & capabilities; card memory model (sectors/blocks, key A/B + access bits; DESFire apps/files/keys); on-card symmetric crypto (DES/2K3DES/3K3DES/AES-128), 3-pass mutual auth & session keys; EV2 SecureMessaging/CMAC, Transaction MAC, EV3 LRP & SUN/SDM; native-vs-7816-4 APDUs (CLA 0x90, 91xx); key diversification (AN10922); SAM/HSM key custody; random UID; perso & default keys. SKIP: skim/clone/relay/Crypto-1 attacks -> rfid-nfc-access-attacks-defenses; RF bands/physics -> rfid-fundamentals-bands-standards; NFC NDEF -> nfc-ndef-protocols; Wiegand/OSDP wire -> access-credentials-wiegand-osdp; reader/panel topology -> pacs-access-control-architecture; mobile/UWB & convergence -> mobile-credentials-aliro-uwb / physical-security-convergence-standards; EMV applets.

---

# Contactless Smart Cards (MIFARE / DESFire) + On-Card Cryptography

The **chip and key layer** of 13.56 MHz contactless credentials: how a card is
addressed (ISO/IEC 14443), how applications talk to it (ISO/IEC 7816-4 APDUs),
what each MIFARE/DESFire tier can do, and how the symmetric keys are
authenticated, diversified, and custodied. This is the "what's inside the card
and how do I trust it" skill: not RF physics, not NDEF, not attacks, not panels
(see SKIP routing in the frontmatter).

> **One-line security caveat:** MIFARE Classic's **Crypto-1** (a proprietary
> 48-bit stream cipher) is **cryptographically broken**: reverse-engineered and
> key-recoverable in seconds (Nohl & Plötz 2007; Garcia et al., *Dismantling
> MIFARE Classic*, ESORICS 2008), so treat MIFARE Classic as having **no
> security value** and migrate to AES (DESFire EV1/EV2/EV3 or MIFARE Plus SL3).
> The attack *methodology and tooling* live in
> **`rfid-nfc-access-attacks-defenses`**; this skill states only the fact.

---

## Core Concepts

**The two layers (don't conflate them).**
- **ISO/IEC 14443** is the *air interface and transport* for **proximity** cards
  at **13.56 MHz**, operating distance **~up to 10 cm**. A card is a **PICC**
  (Proximity IC Card); a reader is a **PCD** (Proximity Coupling Device).
  Half-duplex; base bit rate **106 kbit/s**, optionally 212/424/848 kbit/s. It
  has four parts: Part 1 physical, Part 2 RF/modulation (**Type A vs Type B**),
  Part 3 init/anticollision (UID selection, ATQA/SAK), Part 4 the **T=CL** block
  transmission protocol (RATS/ATS, I/R/S blocks).
- **ISO/IEC 7816-4** is the *application protocol* — the **APDU**
  command/response model (CLA/INS/P1/P2, status words, the MF/DF/EF file model).
  It is the same APDU language used over contact pins; ISO 14443-4 just carries
  it over the air. **14443 = the pipe; 7816-4 = the language.**

**Type A vs Type B** (ISO 14443-2 — both at 13.56 MHz):

| | Type A | Type B |
| --- | --- | --- |
| Reader→card | **100% ASK** (OOK), **modified Miller** | **10% ASK**, **NRZ** |
| Card→reader | **OOK**/Manchester on 847.5 kHz subcarrier | **BPSK**/NRZ on subcarrier |
| Anticollision ID | **UID** (4 / 7 / 10 byte) | **PUPI** (4 byte) |
| Activation | REQA→anticollision→SELECT→**SAK**→RATS→**ATS** | REQB→**ATQB**→**ATTRIB** |

MIFARE products are **Type A**. The **SAK** byte after SELECT tells the reader
whether the UID is complete (bit b3) and whether the card is **ISO 14443-4
compliant** (bit b6 — i.e., speaks T=CL/APDUs). MIFARE **Classic/Ultralight**
have b6=0 (proprietary command set, *not* 7816-4); **DESFire / Plus EV1+** are
14443-4 and speak APDUs.

**Activation vocabulary** (a frequent point of confusion): **ATR** is the
contact-card reset answer (ISO 7816-3) — **there is no ATR on the air
interface**. On contactless, **ATQA/SAK** (Part 3) are the "are you there / what
are you" handshake, and **ATS** (Part 4) is the capability descriptor that plays
the ATR's role. PC/SC middleware often synthesizes a *pseudo-ATR* for non-4
chips (Classic/Ultralight) so the host sees a uniform interface.

**The card data model differs sharply by family:**
- **MIFARE Classic** = a flat **sector/block** memory (16-byte blocks), each
  sector gated by **Key A / Key B + access bits**. No applications, no file
  system, no AES — just Crypto-1.
- **DESFire** = a real **on-card file system**: a PICC master application plus
  **applications (3-byte AID)**, each holding **files** (standard, backup,
  value, linear/cyclic record, + Transaction-MAC on EV2/EV3) with **per-file
  access rights** and **up to 14 application keys**. This is what makes DESFire
  multi-application and AES-grade.

---

## MIFARE / DESFire Tier Comparison

Pick the **lowest tier that meets the security model**; for any access-control or
value application, that is **AES-based** (DESFire EVx or Plus SL3). Numbers below
are from NXP datasheets/fact sheets, cross-checked against the ANSSI EV2 security
target and Wikipedia; a few items are flagged `[verify]`.

| Family / tier | Memory | App / file model | On-card crypto | Mutual auth | Notes |
| --- | --- | --- | --- | --- | --- |
| **Ultralight** (MF0ICU1) | 64 B (16×4 B pages), ~48 B user | none (flat pages + OTP lock bits) | **none** | none | Disposable tickets; no security. |
| **Ultralight C** (MF0ICU2) | 192 B, 7-byte UID | flat pages | **3DES**, 16-byte key | 3DES | Single 16-bit one-way counter. |
| **Ultralight EV1** (MF0ULx1) | 48 B or 128 B user | flat pages + OTP | **32-bit password** (PWD) + 16-bit **PACK**; ECC originality sig | password only | **Three 24-bit one-way counters**; PWD/PACK never readable. Disposable/limited-use. |
| **MIFARE Classic 1K** (EV1) | 1 KB = **16 sectors × 4 blocks** (64 blocks); ~752 B user | flat sectors, no apps | **Crypto-1 (48-bit) — BROKEN** | Crypto-1 | Key A/B 6 bytes each; access bits per block. |
| **MIFARE Classic 4K** (EV1) | 4 KB = **40 sectors** (0–31 ×4 blocks, 32–39 ×16 blocks); ~3,440 B user | flat sectors | **Crypto-1 — BROKEN** | Crypto-1 | EV1 adds ECC originality sig, optional 7-byte UID. |
| **MIFARE Plus** (S/X/SE; EV1/EV2) | 1K/2K/4K class | flat sectors (Classic-compatible) | **AES-128** (+ Crypto-1 for back-compat) | AES (SL3) | **Security levels SL0–SL3** (below); staged Classic→AES migration. |
| **DESFire EV1** (MF3ICDx21/41/81) | **2K / 4K / 8K** | **28 applications**, up to **32 files/app** | **DES, 2K3DES, 3K3DES, +AES-128** | 3-pass | **CC EAL4+**. AES introduced at EV1. |
| **DESFire EV2** (MF3Dx2) | **2K / 4K / 8K** | apps **memory-limited** (EV1's 28-app cap removed `[verify exact wording]`), ≤32 files/app | DES/2K3DES/3K3DES/AES-128 | 3-pass, **EV2 SecureMessaging** | **CC EAL5+**. Announced Mar 2016. **+Transaction MAC, Proximity Check, multiple key sets / key rolling, Virtual Card Architecture.** |
| **DESFire EV3** (MF3D(H)x3) | **2K / 4K / 8K** | flexible file structure, memory-limited apps | DES/2K3DES/3K3DES/AES-128 (+ **LRP** AES mode) | 3-pass, AES **or LRP** SecureMessaging | **CC EAL5+**. Announced 2 Jun 2020. **+LRP (Leakage-Resilient Primitive), SUN/SDM (Secure Dynamic Messaging → dynamic NDEF), Transaction Timer.** Backward-compatible with EV2/EV1. |

**MIFARE Plus security levels** (field-switchable, for staged migration):
- **SL0**: blank/perso state (pre-configuration).
- **SL1**: **MIFARE Classic-compatible** mode using **Crypto-1** (drop-in for
  existing Classic readers; optional AES auth alongside).
- **SL2**: Crypto-1 for data, **AES authentication** layered in.
  `[verify: SL2 is not implemented on all Plus SKUs; some go SL1→SL3.]`
- **SL3**: **full AES-128**; all auth + messaging AES-secured.

**DESFire crypto key sizes:** DES = 8 B; **2K3DES (2TDEA) = 16 B (112-bit)**;
**3K3DES (3TDEA) = 24 B (168-bit)**; **AES-128 = 16 B**. AES is the only choice
to design with new; DES/2K3DES exist for legacy interop only.

---

## The APDU / Authentication Model

DESFire (and other 14443-4 cards) expose commands as **APDUs** over the T=CL
protocol. DESFire offers **two framings** of the same operations:

- **Native-wrapped APDU**: proprietary command in a 7816-style envelope:
  `CLA=0x90  INS=<native cmd code>  P1=0x00  P2=0x00  Lc=<len>  [data]  Le=0x00`.
  Status returns as **`91xx`** — `9100`=OK, **`91AF`=ADDITIONAL_FRAME** (chaining;
  re-issue `90 AF 00 00 00` to fetch the next frame). Other `91xx` are specific
  errors `[verify the full table against the datasheet]`.
- **True ISO 7816-4 APDU**: standard INS bytes (`A4` SELECT, `B0` READ BINARY,
  `D6` UPDATE BINARY, `C0` GET RESPONSE…), ISO file model, status **`9000` /
  `61xx` / `6Cxx` / `6A82` / `6982` …**.

**Rule:** the **first command fixes the session mode** — native, native-wrapped,
and ISO styles **cannot be mixed** within one session.

**ISO/IEC 7816-4 C-APDU** = mandatory header `CLA INS P1 P2` + conditional
`[Lc][data][Le]`, in **four cases**: (1) no data/no response, (2) response only
(`Le`), (3) command data only (`Lc`+data), (4) both. **Short** length uses 1-byte
Lc/Le (max 256 B); **extended** uses a `0x00` marker + 2-byte length (command
data up to 65,535 B, response data up to 65,536 B). **R-APDU** = `[data] SW1
SW2`. Memorize a few status words: **`9000`**
success, **`61xx`** more data (GET RESPONSE), **`6A82`** file not found, **`6982`**
security status not satisfied, **`6700`** wrong length, **`63Cx`** auth failed with
*x* retries left.

**Three-pass mutual authentication (ISO/IEC 9798-2 pattern).** Neither side ever
transmits the shared key K — each proves *knowledge* of it, and both derive a
fresh **session key** from the two nonces:

1. Card → reader: card picks nonce **RndB**, sends it **enciphered** `E_K(RndB)`.
2. Reader → card: reader recovers RndB, picks **RndA**, returns
   `E_K(RndA ‖ rot(RndB))`. Card checks rot(RndB) → reader proven.
3. Card → reader: card returns `E_K(rot(RndA))`. Reader checks → card proven.

Both now hold **RndA and RndB** and combine selected bytes into the **session
key**, used with **CMAC** (integrity) and **CBC** (encryption) for the rest of
the session. DESFire command codes: `Authenticate`=`0x0A` (DES/2K3DES),
`AuthenticateISO`=`0x1A` (3K3DES), `AuthenticateAES`=`0xAA`,
`AuthenticateEV2First`=`0x71` / `AuthenticateEV2NonFirst`=`0x77` (EV2/EV3 secure
messaging; the same `0x71/0x77` codes drive **LRP** mode on EV3).

**Per-file communication modes** select how each file is protected after auth:
**Plain**, **MACed** (CMAC), or **Fully Enciphered** (CBC). EV2 SecureMessaging
adds CMAC-chained, AES-CBC encrypted sessions; the **Transaction MAC** file type
lets a back-end cryptographically verify *that a transaction occurred on the
card*; EV3 **LRP** is a side-channel/fault-resistant secure-messaging mode that
wraps only standard AES-128 (no proprietary crypto). For the full APDU framing,
status-word tables, and the byte-level auth flow, see
**`references/apdu-and-authentication.md`**.

---

## Key Diversification & SAM/HSM (Key Management)

The single most important operational rule: **never put one master key on every
card**, and **never let the master key reach reader software**.

**Key diversification: per-card derived keys (NXP AN10922).** Instead of a
shared key, each card gets a **unique key derived from a master/base key + a
unique input** (typically the **card UID ‖ AID ‖ system identifier**). If one
card's key is extracted, only *that* card is exposed, not the whole estate. The
scheme is **CMAC-based** (a one-way MAC, not encryption), per key type:

| Key type | Master size | Div-input (M) | Leading DIV constant(s) | Output |
| --- | --- | --- | --- | --- |
| **AES-128** | 16 B | 1–31 B | `0x01` | 16 B |
| AES-192 | 24 B | 1–31 B | `0x11`, `0x12` | 24 B |
| 2K3DES (2TDEA) | 16 B | 1–15 B | `0x21`, `0x22` | 16 B |
| 3K3DES (3TDEA) | 24 B | 1–15 B | `0x31`, `0x32`, `0x33` | 24 B |

(CMAC sub-key K1/K2 selection by whether M is padded with `0x80 00…`; AN10922
Rev. 2.2 (2019) extends the same construction to **AES-256**.) The terminal
re-derives the card's unique key at transaction time; the **master key stays in
the SAM/HSM**.

**SAM (Secure Access Module): key custody at the reader edge.** A SAM is itself
a **secure-element smart chip placed *inside the reader*** (removable SIM-style or
soldered). It **stores the master keys, performs the card authentication +
AN10922 diversification + session-key generation on-module**, and runs a secure
host↔SAM channel, so **master keys never appear in reader application
memory/software.**
- **NXP MIFARE SAM AV2**: AV1/AV2 modes; AV2 adds host authentication + three
  messaging modes (**Plain / MAC / Full**) and **key classes** (Host / PICC /
  OfflineChange / OfflineCrypto). **S-mode** = SAM↔host-MCU only; **X-mode** =
  SAM also wired directly to the contactless reader IC for a faster data path.
- **NXP MIFARE SAM AV3**: on **SmartMX2 P60**, hardware **CC EAL6+**. Crypto:
  Crypto-1; TDEA 56/112/168; **AES 128/192/256**; SHA-1/224/256; **RSA; ECC**.
  ~**128 symmetric key entries** (+ RSA/ECC/EMV-CA entries `[verify counts]`);
  host↔SAM auth via AES *or* RSA PKI; S-mode + X-mode; **Programmable Logic**
  for custom diversification/secure-messaging.

**HSM vs SAM — the key hierarchy (backend issuance vs reader runtime).** An
**HSM** is the tamper-resistant device in the **issuance/personalization
backend** whose master keys are flagged **non-exportable** (they never leave in
the clear — a stronger guarantee than a SAM, at higher cost). The chain:

```
Backend HSM  (issuance/perso; master keys, non-exportable)
     │  derives + injects diversified keys (AN10922 CMAC) at personalization
     ▼
Diversified per-card keys ──loaded into──▶  Card (PICC)
     ▲
Reader-side SAM  (holds the SAME master keys; re-derives the per-card key
                  at the door so the master key never enters reader software)
```

**Card UID & random UID (privacy).** DESFire's fixed UID is **7 bytes**. In
**random-ID mode** the card emits a **random 4-byte UID with leading byte
`0x08`** at each anticollision; the real UID is retrievable only **after
authentication** (`Get_Card_UID`). (`NUID`, a 4-byte *non-unique* ID, is a
separate MIFARE Classic EV1 concept to extend the exhausted 4-byte space.) An
**ECC-based NXP originality signature** (readable via `Read_Sig`) lets a reader
verify genuine NXP silicon (anti-clone, probabilistic).

**Personalization lifecycle.** Cards ship **blank with all-zero default keys**
(DES `8×00`, AES `16×00`, 3K3DES `24×00`); the PICC is fully open until
personalized. Perso flow: authenticate with the default/transport key → set the
PICC key-settings → `CreateApplication` (AID, key settings, #keys, crypto type) →
`CreateStdDataFile`/`CreateValueFile`/`CreateRecordFile` with per-file access
rights + comm mode → `ChangeKey` to install the **diversified** keys → optionally
freeze config. **`FormatPICC` (`0xFC`)** wipes all apps/files (master-key-gated).
**Leaving cards on factory-default/known keys is a classic deployment failure** —
the keys *must* be replaced at personalization.

---

## Sources

1. NXP **MF1S50YYX_V1** — MIFARE Classic EV1 **1K** datasheet (sector/block layout, Key A/B, access bits, UID). https://www.nxp.com/docs/en/data-sheet/MF1S50YYX_V1.pdf
2. NXP **MF1S70YYX_V1** — MIFARE Classic EV1 **4K** datasheet (40-sector layout). https://www.nxp.com/docs/en/data-sheet/MF1S70YYX_V1.pdf
3. NXP **MF0ULX1** — MIFARE Ultralight EV1 datasheet (PWD/PACK, three 24-bit counters, ECC originality). https://www.nxp.com/docs/en/data-sheet/MF0ULX1.pdf
4. NXP **MF0ICU2** — MIFARE Ultralight C datasheet (3DES, 192-byte EEPROM). https://www.nxp.com/docs/en/data-sheet/MF0ICU2.pdf
5. NXP **MF3ICDx21_41_81** — MIFARE DESFire **EV1** datasheet (apps/files/keys, DES/2K3DES/3K3DES/AES, EAL4+). https://www.nxp.com/docs/en/data-sheet/MF3ICDX21_41_81_SDS.pdf
6. NXP **MF3Dx2/MF3DHx2** — MIFARE DESFire **EV2** datasheet (EV2 SecureMessaging, Transaction MAC, multiple key sets). https://www.nxp.com/docs/en/data-sheet/MF3DX2_MF3DHX2_SDS.pdf
7. NXP **MF3D(H)x3** — MIFARE DESFire **EV3** datasheet (file types, SUN/SDM, EAL5+, LRP). https://www.nxp.com/docs/en/data-sheet/MF3D_H_X3_SDS.pdf
8. NXP **DESFire EV2 Fact Sheet** (EAL5+, TMAC, key rolling, Virtual Card, Proximity Check). https://www.nxp.com/docs/en/fact-sheet/MIFARE-DESFIRE-EV2-FS.pdf
9. ANSSI — **MF3Dx2 DESFire EV2 Security Target Lite** (independent CC **EAL5+**). https://cyber.gouv.fr/sites/default/files/2016/05/mf3dx2_mifare_desfireev2_securitytargetlite_v1.5.pdf
10. NXP **AN10922** — *Symmetric key diversifications* (CMAC scheme; the canonical per-card key derivation). Mirror: https://www.cardlogix.com/downloads/support/AN10922-symmetric-key-diversifications-nxp-sam-av2.pdf
11. NXP **AN12304** — *Leakage-Resilient Primitive (LRP) Specification* (EV3 secure messaging over AES-128). https://www.nxp.com/docs/en/application-note/AN12304.pdf
12. NXP **MIFARE SAM AV3** datasheet (MF4SAM3X_S) + product page (EAL6+, key entries, S/X-mode, Programmable Logic). https://www.nxp.com/products/MIFSAMAV3
13. NXP **MIFARE SAM AV2** datasheet (P5DF081_SDS) (AV1/AV2 modes, key classes, S/X-mode). https://www.nxp.com/docs/en/data-sheet/P5DF081_SDS.pdf
14. NXP **AN10833** — *MIFARE Type Identification Procedure* (Type A activation/anticollision, SAK bit coding, cascade tag `0x88`). https://www.nxp.com/docs/en/application-note/AN10833.pdf
15. **ISO/IEC 14443** overview + Type A/B modulation/coding, four parts, 13.56 MHz/106 kbit/s. https://en.wikipedia.org/wiki/ISO/IEC_14443
16. **ISO/IEC 7816-4** APDU model — C-APDU/R-APDU, four cases, file model (MF/DF/EF), status words. https://cardwerk.com/smart-card-standard-iso7816-4-section-5-basic-organizations/ ; https://en.wikipedia.org/wiki/Smart_card_application_protocol_data_unit
17. Garcia, de Koning Gans, Muijrers, van Rossum, Verdult, Wichers Schreur, Jacobs — **"Dismantling MIFARE Classic," ESORICS 2008** (canonical Crypto-1 cryptanalysis, Radboud University). https://link.springer.com/chapter/10.1007/978-3-540-88313-5_7
18. **Crypto-1** structure + reverse-engineering (Nohl & Plötz) + NXP migration recommendation. https://en.wikipedia.org/wiki/Crypto-1
19. **Secure access module** (SAM) definition, functions, form factors. https://en.wikipedia.org/wiki/Secure_access_module

> **Uncertainties to confirm against primary PDFs before quoting verbatim:** the
> exact EV2/EV3 maximum-application wording (the EV1 28-app cap is lifted to
> "memory-limited"); MIFARE Plus SL2 availability per SKU; the full DESFire
> `0x91xx` status-code table (only `9100`/`91AF` are high-confidence here); and
> the SAM AV3 RSA/ECC/EMV-CA key-entry counts.