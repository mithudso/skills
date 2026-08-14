# Attack classes: defensive mechanism notes

Companion to SKILL.md. Mechanism-level detail so a defender understands *why* a
control works, written at a defensive/awareness level, **no operational
exploitation recipes** (no command sequences, key-recovery procedures, wiring,
or clone walkthroughs). For card-chip internals see
**contactless-smartcards-mifare-desfire**; for band/coupling physics see
**rfid-fundamentals-bands-standards**.

## 1. Eavesdropping / sniffing

A passive attacker captures either the **RF** between card and reader or the
**electrical signal** on the reader→controller wire. RF capture range is small
for HF (13.56 MHz) but a tuned antenna extends it beyond the nominal few cm.
Cards that send credential data in the clear (any ID-only tag, and the cleartext
phases of weak protocols) leak directly; cards using encrypted secure messaging
(DESFire EV2/EV3, SEOS) leak only ciphertext. **Why short range is not a
control:** range is an attacker-budget parameter, not a security boundary.
**Defense:** encrypted session (secure messaging) so captured traffic is useless,
plus channel protection on the wire (OSDP Secure Channel).

## 2. Skimming

Active version of eavesdropping: the attacker's reader energizes a nearby card
(in a pocket/wallet) and completes enough of an exchange to harvest usable data -
the badge ID for ID-only tags. Crypto cards resist this because the attacker
cannot complete mutual authentication without the key, so they get at most a UID.
**Defense:** don't trust the UID; require mutual auth. Shielding (RF-blocking
sleeves) is a user-side mitigation, not a system control.

## 3. Cloning

- **ID-only tags (LF Prox/EM/AWID/Indala; HF tags checked by UID only):** the
  "secret" is just an identifier that the card broadcasts, so a functional copy
  needs only that identifier on a writable blank or an emulator. This is the
  single most common real-world weakness.
- **Crypto tags (Crypto-1, iCLASS legacy):** cloning first requires *recovering a
  key* (sections 7–8), then the card can be reproduced. DESFire EV2/EV3 and
  SEOS with diversified keys resist cloning because each card's key is unique and
  never exposed.

**Defense:** AES mutual-auth credentials with **per-card diversified keys** so a
single recovered card key cannot forge others.

## 4. Replay

Re-presenting previously captured data. Two distinct surfaces:
- **RF replay** against static-data credentials, the captured response is always
  valid, so emulating it later opens the door.
- **Wiegand replay**, replaying captured cleartext Wiegand bits onto the
  reader→controller wire (section 6).

**Defense:** *freshness*. Challenge–response with a fresh nonce each transaction
and monotonic anti-replay counters mean a captured exchange does not validate
again.

## 5. Relay (mafia fraud): the one crypto can't fix

Two cooperating devices bridge distance in real time: a "leech" near the genuine
card and a "ghost" at the target reader, forwarded over any fast link (often the
internet). The reader runs a **completely valid** session, with the real card,
just far away. Because the cryptography is genuine, *credential strength does not
help*. Named "mafia fraud" (Desmedt/Bengio/Goutier, 1987). Demonstrated
practically on ISO 14443 proximity cards (Hancke) and on contactless transactions
with NFC phones.

**Defense:** bound *distance* or *time*, not just identity -
- **Distance-bounding protocols** measure round-trip time so a relayed (longer)
  path is detectable (Hancke & Kuhn).
- **UWB secure ranging (IEEE 802.15.4z HRP with Scrambled Timestamp Sequence)**
  gives a physically grounded distance estimate; deployed for car entry (Apple
  U1, NXP, Qorvo). **Caveat:** research (WiSec 2021; arXiv 2312.03964) shows
  physical-layer distance-reduction attacks against energy-detection receivers -
  STS narrows but does not eliminate this; receiver design matters. Treat UWB as
  a strong layer, not a guarantee.
- Operationally: minimal reader range, transaction-velocity limits, geofencing,
  and a **second factor** (PIN/biometric/phone presence) on sensitive doors.

## 6. Wiegand wiretapping, replay, brute force

**Wiegand** is the legacy reader→controller protocol: cleartext, one-directional,
open-collector signaling, with no encryption, authentication, or supervision. The
reader sits on the *unsecured* side of the door, so its wires are often reachable
by removing the reader from the wall.
- **Wiretap:** an inline implant (ESPKey-class) on the data lines logs every
  badge's bits in the clear, works **regardless of how strong the card is**,
  because it taps *after* the card has been validated by the reader.
- **Replay:** captured bits are injected later onto the wire to open the door.
- **Brute force / enumeration:** short formats (e.g., 26-bit: 8 facility-code
  bits + 16 ID bits) have a small namespace; from one known card an attacker can
  guess neighbors. Larger/custom formats raise the bar but are *defense in depth*,
  not a fix.

**Defense:** **OSDP with Secure Channel** (section below); reader **tamper
detection**; protect/conceal the wiring run; keep wiring on the secure side where
the door design allows.

## 7. MIFARE Classic: Crypto-1 weaknesses (tech-specific)

Applies to **MIFARE Classic 1K/4K and Classic-compatible** chips only, *not*
DESFire, SEOS, or PIV. The proprietary **Crypto-1** stream cipher and its weak
PRNG were broken by academic research:
- **Garcia et al. (2008, arXiv 0803.2285):** keystream recovery from PRNG and
  cipher weaknesses.
- **Darkside / nested:** families of attacks exploiting NACK leakage and the
  nonce mechanism to recover keys (one known sector key bootstraps the rest in
  the nested case).
- **Hardnested (Meijer & Verdult, CCS 2015):** ciphertext-only attack against
  *hardened* Classic cards that fixed the earlier implementation bugs, it
  attacks the cipher itself.
- **Static encrypted nonce variant (IACR ePrint 2024/1275):** a newer
  Classic-compatible weakness; confirm affected-product scope against the live
  paper.

**Defense:** **retire MIFARE Classic for access control.** No configuration makes
Crypto-1 safe; migrate to AES-based DESFire EV2/EV3 or SEOS with diversified
keys.

## 8. Default / known keys and legacy shared master keys

- **Default keys:** tags left on factory keys (Classic default keys; NTAG/
  Ultralight sectors with no auth; DESFire apps never re-keyed) are readable by
  anyone with a dictionary.
- **HID iCLASS legacy ("standard security"):** used a **single global master
  key** shared across the install base, recovered and published (Meriac, *Heart
  of Darkness*, 2010), enabling read/clone of legacy iCLASS. **iCLASS SE / Elite
  / SEOS** with site-specific keys are the remediation, *not* the legacy global
  key.
- **Encoder/reader key leakage:** even issuance hardware can leak keys, see
  **CISA ICSA-24-037-01** (HID iCLASS SE CP1000 encoder). Patch and rotate.

**Defense:** unique **site keys** (never factory), **per-card diversification**,
and key rotation/patching across cards, readers, *and* encoders.

## 9. Downgrade

Many credentials and systems accept *multiple* technologies for backward
compatibility, a multi-tech card may carry both a strong HF application and an
open LF Prox track; a reader may accept both. An attacker targets the **weakest
accepted** path (read the open LF track, or coerce OSDP back to a clear/Wiegand
mode). **Defense:** disable weak technologies at the reader/controller; don't
issue multi-tech cards with an open legacy track; enforce the strong path only.

## 10. Tamper / jamming / denial of service

- **Reader tamper/removal:** removing the reader exposes wiring (enables section
  6) and may itself be the goal (disable a door).
- **RF jamming / blocker tags:** noise in-band or a blocker tag stops the reader
  reading, an availability attack.
- **Fail mode:** the door's **fail-secure vs fail-safe** behavior on power/comm
  loss is a deliberate, life-safety-constrained decision; attackers probe for the
  weaker mode.

**Defense:** **tamper switches + alarms**, **jam/DoS detection** alerting,
supervised communication (OSDP reports reader presence), and a per-door
fail-mode decision aligned with life-safety code.

## Cross-references

- Countermeasure design and a tech-selection ladder: **references/defenses.md**.
- Card/chip command and key detail: **contactless-smartcards-mifare-desfire**.
- Source list with links: SKILL.md "Sources".
