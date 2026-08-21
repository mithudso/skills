# Defenses and countermeasures: design detail

Companion to SKILL.md. Expands the defense layers and gives a credential-strength
ladder and a configuration-hardening view. Defensive framing throughout. Card-chip
specifics (exact key types, command flows) belong to
**contactless-smartcards-mifare-desfire**; this file owns the *security-design*
reasoning.

## Credential-strength ladder (worst → best for access control)

Pick the strongest the deployment can support; the goal is to leave nothing on the
weak rungs for anything that matters.

1. **LF 125 kHz ID-only**, EM4100, HID Prox, AWID, Indala. No crypto; the ID is
   broadcast. *Trivially cloned.* Avoid for access control.
2. **HF UID-only usage**, MIFARE Classic/Ultralight where only the UID is
   checked. The UID is not a secret. *Trivially cloned.* Avoid.
3. **MIFARE Classic with Crypto-1**, broken cipher (nested/darkside/hardnested).
   *Do not rely on it.* Retire.
4. **HID iCLASS legacy ("standard security")**, global shared master key,
   published. *Do not rely on it.* Migrate.
5. **AES mutual-auth with diversified keys**, **MIFARE DESFire EV2/EV3** (AES-128,
   3-pass mutual auth, EV2 secure messaging), **HID iCLASS SE/SEOS**. Card proves
   key possession without exposing it; per-card keys. *Current good practice.*
6. **PIV / high-assurance**, NIST SP 800-116 mechanisms (**PKI-AUTH**,
   **SYM-CAK**), optionally with on-card biometric (**OCC-AUTH**) and PIN, for
   multi-factor, high-assurance facilities.

Rungs 5–6 close cloning, skimming, RF-replay, and crypto-break, but **not relay**
(needs the ranging/time layer) and **not Wiegand wiretap** (needs the channel
layer).

## Layer 1: strong credentials

- **Mutual authentication:** both card and reader prove themselves; a skimmer
  without the key gets at most a UID.
- **Secure messaging:** the post-auth session is encrypted and integrity-protected
  (DESFire EV2/EV3, SEOS), defeating eavesdropping.
- **Per-card key diversification:** each card's key derives from a master key plus
  the card serial, so extracting one card key does not compromise the fleet. This
  is the structural fix for both the iCLASS-legacy and MIFARE-default failures.

## Layer 2: secure the reader-to-controller channel (OSDP)

**Replace Wiegand with OSDP (SIA standard) + Secure Channel (AES-128).** OSDP adds
encryption, bidirectional supervised communication (the controller knows if a
reader is removed/replaced), and is the SIA-recommended successor to Wiegand.

**But OSDP must be configured correctly**, researchers (Bishop Fox, Petro &
Vargas, *Badge of Shame*) showed five weaknesses in real deployments, all
configuration/protocol issues rather than AES breaks:
- **Secure Channel is off by default** → *enforce it; reject readers that won't
  negotiate it.*
- **Install mode** lets a reader request the base key from the controller and is
  often left on → *disable install mode after pairing.*
- **SCBK-D**, a hardcoded default key, provides no security → *never use it; set a
  unique key per reader.*
- **Key delivery over the daisy-chain** can expose the base key to an inline
  listener → *protect wiring; treat the bus as untrusted during pairing.*
- **The command byte stays cleartext** even inside a Secure Channel session →
  accept this as a residual and lean on tamper detection + monitoring.

Net: OSDP Secure Channel correctly configured defeats Wiegand wiretap/replay; OSDP
misconfigured re-opens it. (Tool awareness: BishopFox `mellon` exercises these.)

## Layer 3: anti-replay (freshness)

Challenge–response with a fresh nonce per transaction and monotonic anti-replay
counters so captured traffic never validates twice. This is a property of the
credential/protocol choice (present in AES mutual-auth schemes), not an add-on.

## Layer 4: anti-relay (distance and time)

Relay forwards a genuine session, so identity strength alone fails. Bound the
physical relationship:
- **Distance-bounding protocols** (Hancke & Kuhn) measure round-trip time; a
  relayed path is longer and detectable.
- **UWB secure ranging, IEEE 802.15.4z HRP with STS** (Scrambled Timestamp
  Sequence): physically grounded distance; deployed for vehicle entry (Apple U1,
  NXP, Qorvo). **Caveat:** distance-reduction attacks exist against
  energy-detection receivers (WiSec 2021; arXiv 2312.03964), STS reduces but
  does not remove the risk; receiver design is decisive. A strong layer, not a
  guarantee.
- **Operational anti-relay:** minimal reader range; transaction-velocity limits;
  geofencing/time-of-day rules; a **second factor** on sensitive doors so a relay
  alone is insufficient.

## Layer 5: physical, monitoring, operational

- **Reader tamper detection** with alarms on removal; mount so removal is
  detectable; keep wiring on the secure side where door design allows.
- **Door alarms:** door-forced, door-held-open, **anti-passback** (a credential
  can't re-enter without exiting, also blunts casual cloning/lending).
- **Jam/DoS detection:** alert when in-band noise exceeds a threshold; OSDP's
  supervision flags lost readers.
- **Fail-mode policy:** decide **fail-secure vs fail-safe** per door under
  life-safety code; document it; don't leave it to default.
- **Defense in depth:** zoning so one cloned badge ≠ whole building; mantraps and
  turnstiles against tailgating; escorts and second factor on crown-jewel areas;
  guard/CCTV correlation with access events.

## Layer 6: provisioning and lifecycle

- **Custom site keys**, never factory/default; **per-card diversification**.
- **Patch the whole chain**, cards, readers, *and* encoders/issuance hardware
  (CISA ICSA-24-037-01 shows encoders can leak keys); rotate keys after exposure.
- **Revocation that reaches the reader/controller** promptly; no orphaned
  credentials; audited issuance and return.

## Standards & references anchors

- **SIA OSDP**, successor to Wiegand; Secure Channel.
- **NIST SP 800-116 Rev. 1**, PIV in PACS; PKI-AUTH / SYM-CAK; risk-based
  mechanism selection; CHUID deprecated.
- **NXP DESFire EV3 (MF3Dx3) data sheet**, AES-128, 3-pass mutual auth, EV2
  secure messaging.
- **CISA ICSA-24-037-01**, HID encoder key-extraction advisory (verify exact
  CVE/CVSS against the live page).

Full link list: SKILL.md "Sources".
