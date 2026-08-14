# UWB secure ranging — deep reference

Conceptual depth for §4 of SKILL.md. **Knowledge as of June 2026; the security
literature here is active — re-verify before citing a hard guarantee.**

## What UWB is
- **Ultra-wideband**: radio that spreads energy over a very wide band (≥500 MHz)
  using extremely short pulses (sub-nanosecond). The wide bandwidth → very fine
  time resolution → precise time-of-flight distance measurement.
- Standardized as **IEEE 802.15.4z-2020** ("4z"), an amendment to 802.15.4 that
  adds **enhanced UWB PHYs and associated ranging techniques**, explicitly for
  *secure* ranging. (The earlier 802.15.4a UWB lacked the security extension.)
- Operates with line of sight at **up to ~200 m**; robust in **congested
  multipath** environments (parking structures, lobbies, airports) where
  RSSI-based BLE proximity is unreliable.

## How ranging works
- **Two-Way Ranging (TWR)** / time-of-flight: device A sends a pulse, device B
  replies; the round-trip time (minus a known processing delay) × speed of light
  / 2 = distance. UWB's pulse sharpness yields **centimeter-class** accuracy.
- Variants: single-sided TWR, double-sided TWR (cancels clock drift), and
  TDoA for positioning infrastructure. For access, TWR between phone and reader
  is the relevant case.

## Why ranging must be *secure* (the relay/distance problem)
- The threat is the **relay attack** (a.k.a. mafia-fraud / distance-fraud): an
  attacker forwards the legitimate radio exchange between a far-away credential
  and the reader so the reader believes the credential is adjacent. This is the
  classic **car key fob relay theft**, and it applies to BLE phone-as-key too.
- BLE/NFC cannot truly prevent this: BLE infers proximity from **RSSI and
  challenge-response latency**, both of which a fast relay defeats (NCC Group
  demonstrated a **link-layer BLE relay adding only ~8 ms** round-trip — under
  the tolerance of typical proximity checks). NFC's short range raises the cost
  but link-layer relay is still demonstrated.
- **Distance bounding** is the cryptographic countermeasure: bound the *physical*
  distance by measuring signal time-of-flight in a way an attacker cannot
  shorten. UWB's time resolution is what makes a *trustworthy* distance bound
  practical at the PHY.

## HRP + STS (the secure-ranging mechanism, conceptually)
- 802.15.4z defines two PHYs: **HRP (High Rate Pulse-repetition-frequency)** and
  LRP (Low Rate). **Secure ranging is specified in HRP**, and **HRP is the mode
  the FiRa certification program certifies.**
- **STS = Scrambled Timestamp Sequence**: an **encrypted, pseudo-random pulse
  sequence** inserted into the UWB frame, generated from a key + nonce that
  **only the two ranging devices share** (typically via a DRBG, e.g. AES-based).
- The ranging timestamp is derived from the STS. Because the STS is
  cryptographically unpredictable to an attacker, the attacker **cannot guess or
  forge the waveform early** to make the pulse appear to arrive *sooner* (which
  is what shortening the measured distance requires). So a relayed/replayed
  signal fails the distance check.
- Conceptually: STS turns "measure arrival time of a known preamble" (forgeable)
  into "measure arrival time of a secret waveform" (not forgeable without the
  key) → a **cryptographically enforced upper bound** on real distance.

## Limits & honest caveats (DATE-SENSITIVE, active research)
- STS is **strong but not unconditionally unbreakable.** Academic work
  (*Secure Ranging with IEEE 802.15.4z HRP UWB*, arXiv 2312.03964, 2023–24)
  shows PHY-layer attacks can still, under some conditions, manipulate the
  measured distance against an HRP-STS implementation — i.e. the security
  depends on implementation choices (STS configuration, detection thresholds),
  not the spec alone.
- Correct framing for advice: **UWB secure ranging raises the relay/distance bar
  by orders of magnitude vs BLE/NFC and is the current best-practice anti-relay
  layer — but treat it as "very hard to relay," not "provably impossible."**

## Ecosystem touchpoints
- Same UWB silicon and ranging power **iPhone/Apple Watch Precision Finding,
  AirTag, and Apple Car Key**; Android UWB (Pixel, Samsung Galaxy) exposes
  ranging APIs; chip vendors: **NXP, Qorvo, STMicroelectronics, Infineon,
  Nordic** (several are also Aliro members).
- **Car Connectivity Consortium (CCC) Digital Key** uses **BLE for wake/discovery
  + UWB for secure ranging** (plus NFC tap fallback) — the reference design for
  relay-resistant hands-free entry. CCC Digital Key certifications rose from
  ~2 (2024) to ~115 (2025) — an adoption-momentum signal.
- **Aliro** standardizes the **BLE + UWB** experience for door access, bringing
  this same secure-ranging model from cars to buildings/homes.
- **FiRa Consortium** profiles 802.15.4z, defines interoperability, and runs the
  UWB certification program (HRP).

## Cross-links
- Relay/skim/clone **attack methodology and tooling** → sibling
  `rfid-nfc-access-attacks-defenses` (this file covers the *defense mechanism*,
  not how to run the attack).
- RF bands/coupling physics → `rfid-fundamentals-bands-standards`.
- Reader↔controller wiring that carries the unlock decision →
  `access-credentials-wiegand-osdp`; door hardware/topology →
  `pacs-access-control-architecture`.
