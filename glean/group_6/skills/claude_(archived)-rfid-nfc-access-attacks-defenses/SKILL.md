---
name: rfid-nfc-access-attacks-defenses
description: >-
  Threat model, attacks, and defenses for card-based physical access control
  (RFID/NFC badges, readers, controllers), for DEFENSE and authorized
  assessment, not offensive how-to. TRIGGER: is a badge deployment weak;
  eavesdrop/skim/clone/replay/relay/downgrade/DoS threats; MIFARE Classic
  Crypto-1 (nested/darkside/hardnested), default-key and UID-clone risk;
  Wiegand wiretap/replay (ESPKey) and OSDP Secure Channel hardening; HID
  Prox/iCLASS legacy vs iCLASS SE/SEOS/DESFire EV2/EV3 AES diversified keys;
  relay/mafia-fraud and UWB/distance-bounding; attacker tooling (Proxmark3,
  Flipper, ChameleonMini) awareness; reader tamper detection;
  weak-deployment checklist. SKIP: card-chip command/key detail ->
  contactless-smartcards-mifare-desfire; RFID bands/physics ->
  rfid-fundamentals-bands-standards; NFC NDEF/tag-format -> the NFC-format
  sibling; reader/controller/PACS architecture -> the PACS sibling; web/app
  security -> security-review; Web Crypto/vault code -> webcrypto-vault-reviewer.
---

# RFID / NFC access security — attacks & defenses

The security/attacks-and-defenses spoke of the physical-access family. Use it to
**reason about how card-based physical access gets defeated and how to defend it**,
and to run an authorized "is our deployment weak?" assessment.

**Defensive framing (read first).** This skill describes attack *classes* and their
*mitigations* at the level needed to defend a system or scope an authorized
penetration test. It deliberately does **not** give step-by-step offensive
exploitation recipes (no key-recovery command lines, wiring diagrams, or clone
walkthroughs). Physical-access testing against systems you do not own or lack
written authorization for is illegal in most jurisdictions. Keep the goal
*hardening*.

**Scope boundary.** This spoke owns the *adversary model and countermeasures*. It
defers the underlying tech to siblings:
- Card chips, MIFARE/DESFire command sets, key structures, product matrix →
  **contactless-smartcards-mifare-desfire**
- LF/HF/UHF bands, coupling, antennas, range physics →
  **rfid-fundamentals-bands-standards**
- NDEF records, tag types, NFC data format → the **NFC-format** sibling
- Readers, controllers, panels, door hardware, system architecture → the **PACS**
  sibling

---

## Threat model

**Assets:** the credential secret (or, for weak tech, just the badge ID), the
reader↔controller channel, the controller's access decision, and ultimately the
door/turnstile. **Trust assumption that usually fails:** "possession of a badge
that reads correctly == an authorized person is present, now, at this door."
Every attack below breaks one clause of that sentence (*authentic*, *present*,
*now*, *this door*).

**Attacker positions:** in RF range of a card or reader (a few cm to ~1 m with a
better antenna); briefly near a victim's wallet/pocket; with momentary physical
access to the wiring behind a reader (often reachable by removing the reader from
the wall — the reader sits on the *unsecured* side of the door); or a remote pair
cooperating over the internet (relay).

**Attack classes (STRIDE-flavored for PACS):**

| Class | What it defeats | Plain description |
| --- | --- | --- |
| Eavesdrop / sniff | Confidentiality of the exchange | Capture the RF or the reader→controller wire and read credential data in transit. |
| Skim | "needs the real card" | Energize and read a card at a distance without the holder's knowledge. |
| Clone | Credential authenticity | Reproduce a credential (trivial for ID-only tags; needs key recovery for crypto tags). |
| Replay | "now / live" | Re-present captured data (RF or Wiegand bits) later to gain entry. |
| Relay (mafia fraud) | "present / this door" | Forward a live exchange in real time between a far-away genuine card and the target reader. |
| Downgrade | Strong-tech assumption | Force the system onto a weak credential/protocol it still accepts (e.g., a multi-tech card's open LF track; OSDP→clear). |
| Crypto break | Key secrecy | Recover the on-card/reader key via cipher or implementation weakness (Crypto-1, iCLASS master key). |
| Tamper / DoS | Availability & integrity | Pull/jam a reader, jam the RF, or fail the door open/closed. |
| Default/known keys & enrollment abuse | Key secrecy / provisioning | Read or write tags protected only by factory/published keys; abuse install/enrollment modes. |

---

## Attacks → affected tech → tooling → mitigation

A defender's quick-reference. Tools are named **for awareness** (so you know what a
realistic adversary carries), not as instructions. "Affected tech" precision
matters — most attacks are tech-specific.

| Attack | Affected tech (precise) | Tool an attacker might use | Primary mitigation |
| --- | --- | --- | --- |
| UID/ID clone | LF 125 kHz ID-only: EM4100, HID Prox, AWID, Indala; HF UID-only use of MIFARE Classic/Ultralight where only the UID is checked | Flipper Zero, Proxmark3, ChameleonMini/Ultra, NFC phone | Stop authenticating on UID/ID alone; move to a credential with mutual auth (DESFire EV2/EV3 AES, iCLASS SE/SEOS, PIV PKI-AUTH/SYM-CAK) |
| Crypto-1 key recovery (darkside / nested / hardnested) | **MIFARE Classic** (1K/4K) and Classic-compatible; *not* DESFire/SEOS | Proxmark3; Flipper does nested-with-dictionary only | Retire MIFARE Classic for access; use AES-based cards with per-card diversified keys |
| Default / published keys | Any tag left on factory/known keys (Classic default keys; NTAG/Ultralight no-auth; misconfigured DESFire apps) | Any reader/dictionary tool | Provision unique site keys; never ship factory keys; diversify keys per card |
| iCLASS legacy master-key clone | **HID iCLASS legacy / "standard security"** (shared global key); *not* iCLASS SE/Elite/SEOS | Proxmark3 (iCLASS support) | Migrate to iCLASS SE/SEOS or DESFire EV2/EV3; use Elite/custom keys, not the global key |
| Eavesdrop the RF | Any contactless during an exchange; worst on cards with no encrypted session | SDR / Proxmark3 sniff | Encrypted secure messaging (DESFire EV2/EV3, SEOS); short range alone is *not* a control |
| Wiegand wiretap + replay | **Any** card tech behind a reader using the legacy **Wiegand** reader→controller protocol (cleartext) | ESPKey-class implant on the reader wires | **OSDP with Secure Channel** instead of Wiegand; reader tamper detection; conduit/run protection |
| Wiegand brute force / enumeration | Short formats (e.g., 26-bit) with one known card | Wiegand injector behind reader | Larger/custom encoded formats, but treat as defense-in-depth; the real fix is OSDP + tamper alarms |
| OSDP attacks (install-mode key request, daisy-chain key capture, downgrade-to-clear, cleartext command byte) | **OSDP** deployments that don't enforce Secure Channel, leave install mode on, or use the default key SCBK-D | BishopFox `mellon`-class line tool | Enforce Secure Channel; remove install mode after pairing; ban SCBK-D; key per reader; monitor for clear sessions |
| Replay (RF) | Static-data credentials (ID-only tags) | Flipper/Proxmark emulation | Challenge–response / nonce-based mutual auth; anti-replay counters |
| Relay / mafia fraud | **Any** proximity credential, incl. crypto cards — relaying forwards a *valid* live session | Two NFC phones / custom relay rig over IP | **Distance bounding / UWB secure ranging** (802.15.4z STS); short reader range; transaction velocity & geofencing; require a second factor |
| Reader tamper / removal | Physical reader (gives wire access) | Screwdriver | Tamper switch + alarm; mounting that detects removal; keep wiring on the secure side where possible |
| Jamming / DoS | The RF channel and/or reader | RF noise source / blocker tag | Jam-detection alerting; fail-secure vs fail-safe policy chosen per door + life-safety code |

Deep mechanism notes (still defensive, no recipes) live in
**references/attack-classes.md**.

---

## Defenses & countermeasures

Layered. No single control is sufficient; relay in particular defeats *all*
cryptographic strength because it forwards a genuine session.

**1. Use credentials with mutual authentication and secrecy.** Move off ID-only
tags. Prefer **MIFARE DESFire EV2/EV3 (AES-128, 3-pass mutual auth, EV2 secure
messaging)** or **HID iCLASS SE/SEOS / PIV (PKI-AUTH or SYM-CAK)**. The card
proves it holds a secret without exposing it; the reader proves itself too.

**2. Diversify keys per card.** Derive each card's key from a master key + the
card UID/serial so one extracted card key cannot read or forge others
(DESFire/SEOS support this). Never run a fleet on one shared key (the iCLASS
legacy and MIFARE-default failures).

**3. Secure the reader→controller channel.** Replace **Wiegand** (cleartext, open
to wiretap/replay/brute-force) with **OSDP (SIA) + Secure Channel (AES-128)**.
Then *configure it correctly*: enforce Secure Channel (it is off by default),
disable install mode after pairing, never use the default key **SCBK-D**, and use
a unique key per reader — OSDP misconfiguration re-opens the door it was meant to
close.

**4. Defeat replay with freshness.** Challenge–response with nonces and
monotonic anti-replay counters so captured traffic cannot be re-presented.

**5. Defeat relay with distance/time, not just crypto.** Add **distance bounding
/ UWB secure ranging (IEEE 802.15.4z HRP with STS)** where supported; keep reader
range minimal; enforce transaction velocity, geofencing, and a second factor
(PIN/biometric/phone presence) for high-value doors. Note UWB ranging itself has
shown physical-layer distance-reduction weaknesses against energy-detection
receivers — treat it as a strong layer, not a guarantee.

**6. Physical & monitoring layers.** Reader **tamper detection** with alarms;
protect/conceal wiring runs; **anti-passback** and door-held/forced alarms;
**jam-detection** alerting; choose **fail-secure vs fail-safe** per door against
life-safety code; log and review reader/controller events.

**7. Provisioning & lifecycle hygiene.** Custom site keys (never factory),
revocation that actually reaches readers, no orphaned credentials, secure
issuance, and firmware patching (e.g., the **HID iCLASS SE encoder key-extraction
advisory, CISA ICSA-24-037-01**, shows encoders/readers themselves can leak
keys — patch and rotate).

More design detail and a tech-selection ladder in **references/defenses.md**.

---

## "Is our deployment weak?" checklist

Run this for an authorized self-assessment. Any **Yes** in the first block is a
material weakness.

**Credential technology**
- [ ] Do any doors authenticate on **UID / card-ID only** (no on-card crypto)?
- [ ] Are we using **125 kHz Prox / EM / AWID / Indala** for anything that matters?
- [ ] Are we using **MIFARE Classic** or **HID iCLASS legacy** (Crypto-1 /
  shared-master-key tech)?
- [ ] Are any cards still on **factory/default keys** or a single shared key
  (no per-card diversification)?

**Reader-to-controller channel**
- [ ] Is any reader wired with **Wiegand** (not OSDP)?
- [ ] If OSDP: is **Secure Channel actually enforced**, **install mode disabled**,
  and **SCBK-D / default keys removed**?
- [ ] Are reader wiring runs reachable from the **unsecured** side without a
  tamper alarm?

**Anti-relay / anti-replay**
- [ ] Could a **relay** succeed (no distance bounding / velocity / second factor
  on sensitive doors)?
- [ ] Do credentials carry **anti-replay** freshness, or is the exchange static?

**Physical / operational**
- [ ] Do readers have **tamper detection** that alarms on removal?
- [ ] Are **door-forced / door-held / anti-passback** alarms enabled and watched?
- [ ] Is there **jam/DoS detection** and a deliberate **fail-secure/fail-safe**
  decision per door?
- [ ] Is credential **revocation** prompt and verified at the reader/controller?
- [ ] Is reader/controller/encoder **firmware patched** (e.g., ICSA-24-037-01)?

**Process**
- [ ] Is there **defense in depth** (a single cloned badge ≠ full building access:
  zoning, escorts, mantraps, second factor on crown-jewel areas)?
- [ ] Is physical-access testing **authorized in writing** and scoped before any
  assessment?

A deployment that uses AES mutual-auth diversified-key credentials, OSDP Secure
Channel correctly configured, reader tamper alarms, anti-replay, and a second
factor on sensitive doors has closed the common attack classes — except relay,
which needs the ranging/velocity/second-factor layer specifically.

---

## Sources

Authoritative and community references used to ground this skill.

**Peer-reviewed cryptanalysis & relay/ranging research**
- de Koning Gans, Hoepman, Garcia, *A Practical Attack on the MIFARE Classic*
  (CARDIS/arXiv 0803.2285) — Crypto-1 keystream recovery.
  https://arxiv.org/pdf/0803.2285
- Meijer & Verdult, *Ciphertext-only Cryptanalysis on Hardened MIFARE Classic
  Cards* (ACM CCS 2015) — the **hardnested** attack.
  http://cs.ru.nl/~rverdult/Ciphertext-only_Cryptanalysis_on_Hardened_Mifare_Classic_Cards-CCS_2015.pdf
- *MIFARE Classic: exposing the static encrypted nonce variant* (IACR ePrint
  2024/1275) — newer Classic-compatible nonce weakness. https://eprint.iacr.org/2024/1275.pdf
- Hancke & Kuhn, *An RFID Distance Bounding Protocol* / "Keep Your Enemies
  Close: Distance Bounding Against Smartcard Relay Attacks" — relay + distance
  bounding foundations. https://www.researchgate.net/publication/250317922
- *A Practical Relay Attack on ISO 14443 Proximity Cards* (Hancke).
  https://www.researchgate.net/publication/228902569
- Singh et al., *Security Analysis of IEEE 802.15.4z/HRP UWB Time-of-Flight
  Distance Measurement* (ACM WiSec 2021) and *Secure Ranging with IEEE
  802.15.4z HRP UWB* (arXiv 2312.03964) — UWB STS ranging and its
  distance-reduction limits. https://arxiv.org/pdf/2312.03964

**Vendor / standards / advisories**
- Security Industry Association (SIA), **OSDP** standard & "move from Wiegand to
  OSDP." https://www.securityindustry.org/industry-standards/open-supervised-device-protocol/
  and https://www.securityindustry.org/2021/11/09/there-is-a-hole-in-the-boat-why-access-control-professionals-need-to-move-from-wiegand-to-osdp/
- NXP, **MIFARE DESFire EV3 (MF3Dx3) data sheet** — AES-128, 3-pass mutual
  auth, EV2 secure messaging. https://www.nxp.com/docs/en/data-sheet/MF3D_H_X3_SDS.pdf
- NIST **SP 800-116 Rev. 1**, *Guidelines for the Use of PIV Credentials in
  Facility Access* (PACS, PKI-AUTH/SYM-CAK, CHUID deprecation).
  https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-116r1.pdf
- **CISA ICS advisory ICSA-24-037-01**, *HID Global Encoders* (iCLASS SE
  CP1000 encoder key-extraction). https://www.cisa.gov/news-events/ics-advisories/icsa-24-037-01

**RFID-hacking / red-team community (defense awareness)**
- Bishop Fox (Petro & Vargas), *Badge of Shame: Breaking Into Secure Facilities
  with OSDP* — five OSDP weaknesses; `mellon` tool.
  https://bishopfox.com/blog/breaking-into-secure-facilities-with-osdp
  and https://github.com/BishopFox/mellon
- Meriac, *Heart of Darkness — exploring the uncharted backwaters of HID iCLASS
  security* (2010) — iCLASS legacy master-key recovery.
- **Proxmark3** (Iceman/RfidResearchGroup) docs — capabilities of the leading
  RFID research tool. https://github.com/RfidResearchGroup/proxmark3
- L. Gavoni, *RFID Exploitation and Countermeasures* (arXiv 2110.00094) — survey
  of attacks + countermeasures. https://arxiv.org/pdf/2110.00094

**Uncertainty notes.** CISA ICSA-24-037-01 was cited from secondary summaries
(the primary page returned 403/timeout at authoring time) — verify the exact
CVE IDs and CVSS scores against the live advisory before quoting them. The 2024
ePrint static-encrypted-nonce result was evolving at authoring time; confirm
current affected-product scope. UWB STS security depends heavily on the
receiver implementation — vendor claims of "relay-proof" should be tested, not
assumed.
