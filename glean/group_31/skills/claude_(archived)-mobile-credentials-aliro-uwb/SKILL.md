---
name: mobile-credentials-aliro-uwb
description: >-
  Frontier of phone-as-credential physical access: mobile credentials, the Aliro
  open standard, and UWB secure ranging. TRIGGER: phone/wearable as badge or key
  via NFC (HCE) vs BLE vs UWB; Apple Wallet employee badge + the NFC & SE
  Platform / "Apple access" program, Google Wallet corporate badge, Samsung
  Wallet; provisioning/issuance to a wallet, eSE/secure element, applet
  attestation; BLE access UX (twist-and-go, tap) and range/UX/security
  trade-offs; UWB (IEEE 802.15.4z) secure ranging, HRP/STS distance bounding,
  anti-relay, hands-free access, CCC Digital Key; Aliro (CSA) goals/members/
  status + tie to wallets and existing readers; mDL (ISO/IEC 18013-5/-7) and
  verifiable digital credentials as access-adjacent; FIDO2/passkey <-> physical
  convergence. SKIP: RF bands/physics -> rfid-fundamentals-bands-standards;
  NDEF/NFC format -> nfc-ndef-protocols; card chip/key internals ->
  contactless-smartcards-mifare-desfire; attack methodology ->
  rfid-nfc-access-attacks-defenses; Wiegand/OSDP wiring ->
  access-credentials-wiegand-osdp; panels/door hardware/topology ->
  pacs-access-control-architecture; PIV/FICAM + federal-identity/compliance ->
  physical-security-convergence-standards.
---

# Mobile Credentials, Aliro & the UWB Access Frontier

The frontier spoke of the physical-access family: the phone (and wearable)
becoming the credential. Fast-moving — **claims are date-stamped; verify
date-sensitive items before relying on them (knowledge as of June 2026).**

This skill owns the *emerging* layer: wallet-based credentials, the **Aliro**
open standard, **UWB secure ranging**, **mDL**, and the FIDO2↔physical-access
convergence. It defers all the underlying card/RF/reader/panel mechanics to its
siblings (see SKIP in the frontmatter).

---

## 1. Mobile-credential models: NFC-HCE vs BLE vs UWB

A "mobile credential" replaces a plastic badge or metal key with a
cryptographic credential held on a phone/wearable. Three radio transports carry
it, each with a different UX/security profile:

| Transport | Range | UX | Security posture | Typical role |
|---|---|---|---|---|
| **NFC** (13.56 MHz) | ~≤4 cm, intentional tap | Tap-to-enter; deterministic, one door | Proximity is the control; relay-bounded by physics but **link-layer relay is still possible** | Badge tap, transit, payment, key |
| **BLE** | ~up to tens of m | "Twist-and-go," seamless/hands-light, can pre-authenticate at a distance | Convenient but **proximity is *inferred* from RSSI/timing → spoofable; relay-prone** | Long-range/seamless access, parking, garage |
| **UWB** (IEEE 802.15.4z) | ~up to ~200 m LoS; cm-accurate | Truly hands-free, intent-by-location | **Cryptographic distance bounding (STS)** resists relay | Hands-free access, car keys, anti-relay layer |

Key UX/security trade-off (the heart of the BLE-vs-NFC question): **NFC makes
the user's *intent* explicit (you must tap), so range is the security boundary.
BLE buys convenience (the door can sense you approaching) but it *infers*
proximity from signal strength and challenge-response timing, which an attacker
can relay or spoof.** UWB resolves the tension — it keeps the seamless BLE-style
UX *and* adds a cryptographically secured time-of-flight distance measurement so
the reader knows the credential is *physically* within bounds. In practice
modern phone-access stacks **combine** them: BLE for discovery/wake + UWB for
secure ranging + NFC as the universal tap fallback.

### Where the credential lives: HCE vs Secure Element (eSE)
- **HCE (Host Card Emulation)** — the credential/app logic runs in the OS/app,
  keys protected by the OS keystore (TEE/StrongBox on Android). Cheaper, no SE
  partner needed, but a softer trust boundary. Common for BLE-first systems.
- **Embedded Secure Element (eSE)** — a tamper-resistant chip; applet + keys
  live in dedicated hardware. This is what Apple Wallet and high-assurance
  badges use; it also enables **express/power-reserve** taps (works when the
  phone is off or battery is low) because the SE can transact without the OS.
- Conceptual SE-vs-HCE detail and the NDEF/HCE API surface live in the **sibling
  `nfc-ndef-protocols`** and SE/key-custody detail in
  **`contactless-smartcards-mifare-desfire`** — defer chip/key internals there.

### Provisioning / issuance to a wallet (generic flow)
1. Employer/service provider runs an app (or MDM) that authenticates the user.
2. The wallet platform's server downloads a **signed applet** for the card
   scheme and **carves a memory partition (applet instance)** on the device SE.
3. Control passes to the **service-provider/TSM ("partner") servers** to
   **personalize** the instance (write the credential/keys).
4. The badge appears in the wallet; the user taps to enter.
- **Attestation:** before deployment the applet must be **reviewed and validated
  by an independent accredited third-party lab** (Apple's requirement) — proof
  the applet is non-harmful and follows platform-security guidance. Device-level
  attestation (Secure Enclave / hardware keystore) underpins trust that the
  credential lives on genuine hardware.

---

## 2. The three wallets (employee/access badge programs)

> Coverage, eligibility, and OS-version gates change frequently — **verify
> current state for any specific deployment.**

### Apple — Wallet + the NFC & SE Platform ("Apple access")
- **iOS 18.1 (late 2024)** introduced the **NFC & SE Platform** APIs that let
  *authorized* developers add/store/present a contactless card from inside their
  own iOS app — the framework that replaced the older "you need PassKit/closed
  program" gate. Supported use cases: **in-store payments, car keys, closed-loop
  transit, corporate badges, student ID, home keys, hotel keys, loyalty, event
  tickets**, and **government ID (iOS 26.4+)**.
- **Corporate Badge eligibility:** you must operate an office building (or have a
  binding agreement to offer virtual corporate badges) in eligible territories;
  then you request access to the iOS APIs to develop/distribute a badge app.
- **Express Mode / power reserve:** an SE feature — tap without unlocking the
  phone, and (on supported devices) for several hours after the battery dies.
- **Employee Badge in Apple Wallet** has shipped via partners (HID, Kastle,
  Safetrust, Sharry, rf IDEAS, Wavelynx/SwiftConnect, etc.) since 2022.
- The broader **Apple Wallet credential ecosystem** (access-adjacent):
  **Home Key** (HomeKit/Matter locks — Schlage, Yale, Aqara, Level, Lockly,
  ULTRALOQ), **Car Key** (Car Connectivity Consortium **Digital Key**, NFC + BLE
  + UWB), **Hotel Key**, and **IDs in Wallet** (state mDLs + the new Digital ID,
  see §4).

### Google — Wallet corporate badge
- **Corporate Badges** are NFC-enabled digital cards in Google Wallet for
  building access. Requirements: Android 9.0+, NFC on, latest Service Provider +
  Google Wallet apps; **one phone + one watch** at a time; workplace needs
  **NFC readers that accept a Wallet credential via the service provider's
  protocol**. Has a developer program (`developers.google.com/wallet/access`),
  managed-device (Android Enterprise / MDM) provisioning, and **Wear OS** badge
  support (Dec 2024). Google Android also carries **CCC Digital Key** car keys
  and Android's identity-credential mDL stack.

### Samsung — Wallet (Company ID / access)
- Samsung Wallet holds a **Company ID / access key**; once added, the user need
  not wake or unlock the phone to tap. Deployed via partners (e.g.
  Wavelynx/SwiftConnect deploying *corporate badge in Google Wallet + company ID
  in Samsung Wallet*). Samsung devices were also early UWB phones.

**Reader reality:** wallet badges still talk to physical readers. Vendors
(HID **Signo**/iCLASS SE with **Seos**/Origo, Wavelynx, Kastle, Nedap,
Allegion/Schlage) bridge wallet credentials to existing reader/panel infra — so
a wallet rollout is usually a *credential* change, not a full reader rip-and-
replace. Reader↔controller wiring (OSDP/Wiegand) and panel topology are the
**`access-credentials-wiegand-osdp`** and **`pacs-access-control-architecture`**
siblings.

---

## 3. Aliro — the open mobile-access standard (CSA)

**What it is.** Aliro is the **Connectivity Standards Alliance (CSA)** open
standard for mobile/wearable physical access — a common credential + communication
protocol so *any* compliant device unlocks *any* compliant reader, across mobile
wallets. Think "the Matter of door access." (CSA is the same body that runs
Matter and Zigbee.)

**Why it exists.** Mobile access today is fragmented and vendor-locked (HID
Origo/Seos, Apple's program, per-vendor BLE apps). Aliro aims to **reduce
implementation barriers, cut complexity across the value chain, and give one
interoperable, secure digital-access experience** — one device, many doors,
many vendors.

**Who.** Announced **Nov 9, 2023**; founding/contributing members include
**Apple, Google, Samsung, ASSA ABLOY, Allegion, HID, Infineon, NXP,
STMicroelectronics, Qualcomm, Kastle Systems, Last Lock**. By the 1.0 release
**220+ member companies** had contributed. The Aliro Working Group has org
roles drawn from ASSA ABLOY/Allegion/etc.

**Status (DATE-SENSITIVE).**
- **Aliro 1.0 specification released ~Nov 8, 2024**, *with* a formal
  **certification program**, as a **living standard**.
- **Commercial products / first certifications** tracked into **early–mid 2026**
  (industry coverage Jan–Feb 2026 described products "launching shortly" and
  the 1.0 standard "launching"; **Google Wallet confirmed Aliro support**).
  First-to-certify names floated: **Apple, Allegion, Aqara, Google, HID,
  Kastle, Kwikset, Last Lock, Nordic Semiconductor, Nuki, NXP, Qorvo, Samsung,
  STMicroelectronics.** *Verify what has actually shipped/certified now.*

**How it works (conceptual).**
- **Transports:** NFC, BLE, and **BLE + UWB** — the same three models in §1, so
  Aliro is the standardization of the whole transport stack, *including* UWB
  secure ranging.
- **Crypto:** **asymmetric / public-key, certificate-based authentication** to
  establish trust, then **symmetric encryption** for the session. Aligns with
  the wallet ecosystems (Apple/Google/Samsung) and defines credential data,
  NFC/BLE/UWB experiences, and the asymmetric-crypto model.
- **Relationship to readers/Apple's ecosystem:** an Aliro reader reads
  credentials *across* wallets (tap your office in the morning, your front door
  at night — same device, different credentials). It complements rather than
  erases existing programs; vendors map it onto installed reader infra.

**Adjacent open standards (don't confuse them).** **PKOC** (Public Key Open
Credential — a simpler open public-key card/credential format) and **LEAF**
are sibling open-credential efforts; Aliro is the broader mobile-centric
interoperability standard. The industry framing in 2026 is "Aliro, PKOC, LEAF —
same direction (open, public-key, interoperable), different scope."

---

## 4. UWB secure ranging (the anti-relay layer)

See `references/uwb-secure-ranging.md` for the deep dive. Essentials:

- **UWB** = **IEEE 802.15.4z** (the "-4z" amendment, 2020) ultra-wideband PHY.
  It measures distance by **time-of-flight** of very short RF pulses →
  **centimeter-class** ranging, robust in multipath (garages, lobbies).
- **Why it matters for access:** BLE/NFC proximity can be **relayed** (attacker
  forwards the radio exchange to make a far credential look near — the
  car-key/relay attack). UWB does **cryptographic distance bounding**: the
  reader learns a *trustworthy upper bound* on the real physical distance, so a
  relayed signal fails the range check.
- **HRP + STS:** secure ranging uses the **High Rate Pulse-repetition-frequency
  (HRP)** PHY with the **Scrambled Timestamp Sequence (STS)** — an *encrypted*,
  pseudo-random waveform embedded in the UWB frame that only the two ranging
  devices know. Because the timestamp marker is unpredictable to an attacker,
  they can't forge an *earlier* arrival to shorten the measured distance. HRP is
  the mode certified by the **FiRa Consortium** program.
- **FiRa Consortium** profiles 802.15.4z features and runs interoperability +
  **certification** for UWB ranging.
- **Ecosystem:** this is the same tech behind **iPhone/Apple Watch precise
  finding, AirTag, and Car Key**; the **Car Connectivity Consortium (CCC)
  Digital Key** spec uses BLE+UWB for hands-free, relay-resistant car entry
  (CCC certifications jumped from ~2 in 2024 to ~115 in 2025 — momentum signal).
- **Caveat (date-sensitive):** academic work (e.g. *Secure Ranging with IEEE
  802.15.4z HRP UWB*, 2023–24) shows STS, while strong, is **not unconditionally
  unbreakable** under some adversarial PHY attacks — treat UWB as "raises the
  bar dramatically vs BLE/NFC," not "provably unspoofable."

---

## 5. mDL & verifiable digital credentials (access-adjacent)

A **mobile driver's license (mDL)** is a government ID on the phone; it's
access-*adjacent* — increasingly the same wallet, sometimes used for
high-assurance entry/age/identity checks.

- **ISO/IEC 18013-5:2021** — the core mDL standard: data model + **proximity**
  (in-person) presentation (device retrieval over NFC/BLE/QR), with
  **selective disclosure** (share only "over 21," not your address),
  **issuer-signed** data, and **anti-cloning**.
- **ISO/IEC TS 18013-7:2025** — adds **online / unattended (remote)**
  presentation, including via **OpenID for Verifiable Presentations (OID4VP)**.
- **Verifiable Credentials (W3C VC) / OID4VC** — the broader open framework for
  cryptographically verifiable digital credentials; the mDL world and the VC
  world are converging on OID4VP for remote flows.
- **Deployment (DATE-SENSITIVE):** state **mDLs in Apple Wallet** (12+ states +
  Puerto Rico, plus Japan My Number Card by late 2025); **Apple Digital ID**
  (built from a **U.S. passport**) launched **Nov 2025**, accepted in beta at
  **TSA checkpoints in 250+ airports** (also via Google/Samsung wallets and
  state apps). *Not* a passport replacement; coverage expanding — verify.
- **The full federal-identity / PIV / FICAM / compliance landscape is the
  `physical-security-convergence-standards` sibling — only the Aliro/mDL
  frontier lives here.**

---

## 6. FIDO2 / passkeys ↔ physical-access convergence (brief)

The end-state many vendors point at: **one credential for logical *and*
physical** access.

- **FIDO2 / passkeys** (FIDO Alliance: WebAuthn + CTAP) are phishing-resistant,
  public-key login credentials — the same **asymmetric, certificate/key-based**
  trust model Aliro and mDL use, which is *why* convergence is plausible.
- **Converged badge / converged credential:** smart cards (and increasingly
  phones) that carry **PKI/FIDO for IT login + a physical-access credential** on
  one token (e.g. Thales, HID converged offerings) — tap to enter the building,
  tap/auth to log into the laptop.
- The shared direction: public-key credentials, device-bound keys, hardware
  attestation, and selective disclosure — converging the badge, the car/home
  key, the government ID, and the login passkey onto the same secure-element +
  wallet substrate. **Convergence is directional/aspirational, not a single
  shipped product** — flag it as such.

---

## Quick decision/triage cues
- "Phone unlocks door from a distance / twist-and-go" → BLE UX, watch for relay
  → recommend **UWB secure ranging** for anti-relay; NFC tap as fallback.
- "Make any phone open any reader, vendor-neutral" → **Aliro** (+ note PKOC/LEAF).
- "Put the badge in Apple/Google/Samsung Wallet" → §2 program + eligibility +
  SE/HCE + provisioning/attestation in §1.
- "Relay/car-key attack on phone entry" → §4 UWB/STS distance bounding (attack
  *methodology* itself → `rfid-nfc-access-attacks-defenses`).
- "Driver's license / ID on the phone for entry or age check" → mDL §5
  (federal PIV/FICAM/compliance → `physical-security-convergence-standards`).
- "One credential for login + door" → §6 FIDO/passkey convergence.

---

## Sources
All accessed June 2026; mark date-sensitive items before relying on them.

**Aliro / CSA**
- CSA newsroom — *Introducing Aliro 1.0* (1.0 release, members, transports, cert program): https://csa-iot.org/newsroom/introducing-aliro-1-0-a-unified-standard-to-transform-the-access-control-ecosystem/
- CSA / PR Newswire — *CSA Announces Aliro* (Nov 9, 2023 founding, founding members): https://www.prnewswire.com/news-releases/the-connectivity-standards-alliance-announces-aliro-a-new-effort-to-make-mobile-devices--wearables-central-to-a-digital-access-future-301982416.html
- SecurityInfoWatch — *CSA Releases Aliro 1.0 Access Control Standard*: https://www.securityinfowatch.com/access-identity/access-control/news/55360417/connectivity-standards-alliance-releases-aliro-10-access-control-standard
- Allegion — *Aliro* founding-member blog: https://www.allegion.com/corp/en/news/blog/2023/Aliro-announcement.html
- Security Industry Association — *Aliro, PKOC, LEAF… What Do They All Mean?* (2026 landscape): https://www.securityindustry.org/2026/03/18/aliro-pkoc-leaf-what-do-they-all-mean/
- SecurityInfoWatch — *Aliro vs. PKOC: Two Standards, One Direction*: https://www.securityinfowatch.com/access-identity/access-control/article/55378999/aliro-vs-pkoc-two-standards-one-direction
- AppleInsider — *Apple-backed smart lock standard Aliro is launching shortly* (Jan 2026, date-sensitive status): https://appleinsider.com/articles/26/01/05/apple-backed-smart-lock-standard-aliro-is-launching-shortly
- kormax/aliro — Aliro access credential protocol research (community reverse-engineering): https://github.com/kormax/aliro

**Wallets (Apple/Google/Samsung)**
- Apple Developer — *NFC & SE Platform* (iOS 18.1 APIs, provisioning, applet attestation, corporate-badge eligibility, use cases): https://developer.apple.com/support/nfc-se-platform/
- Apple Support — *Add an employee badge to Apple Wallet*: https://support.apple.com/en-us/119901
- Apple Newsroom — *Apple introduces Digital ID* (Nov 2025, TSA 250+ airports): https://www.apple.com/newsroom/2025/11/apple-introduces-digital-id-a-new-way-to-create-and-present-an-id-in-apple-wallet/
- Google for Developers — *Corporate Badge Overview*: https://developers.google.com/wallet/access/corporate-badge/get-started/overview
- Google Wallet Help — *Save your corporate badge to Google Wallet*: https://support.google.com/wallet/answer/14085932
- SecurityInfoWatch — *SwiftConnect, Wavelynx deploy corporate badge in Google Wallet and company ID in Samsung Wallet*: https://www.securityinfowatch.com/access-identity/access-control/cards-tokens/press-release/53099266/

**UWB / secure ranging / CCC**
- FiRa Consortium — *What UWB Does* (ranging, hands-free access, 802.15.4z, certification): https://www.firaconsortium.org/discover/what-uwb-does
- FiRa Consortium — *UWB Secure Ranging: Revolutionizing Security Technology*: https://www.firaconsortium.org/resource-hub/blog/uwb-secure-ranging-revolutionizing-security-technology
- *Secure Ranging with IEEE 802.15.4z HRP UWB* (arXiv 2312.03964 — HRP/STS, limits): https://arxiv.org/abs/2312.03964
- Car Connectivity Consortium — *Digital Key* / Advanced Digital Key initiative: https://carconnectivity.org/ccc-members-advanced-digital-key-initative/
- NCC Group — *BLE relay vulnerability* (link-layer relay, ~8ms, why BLE proximity is spoofable): https://www.nccgroup.com/newsroom/ncc-group-uncovers-bluetooth-low-energy-ble-vulnerability-that-puts-millions-of-cars-mobile-devices-and-locking-systems-at-risk/

**mDL / verifiable credentials**
- Dock — *ISO 18013-5 Standard: What It Is And How It Works*: https://www.dock.io/post/iso-18013-5
- ID Tech Wire — *Newly Published ISO/IEC Standard Enables Remote IDV via mDL* (18013-7:2025, OID4VP): https://idtechwire.com/newly-published-iso-iec-standard-enables-remote-idv-via-mobile-drivers-license/
- TSA — *Digital ID* (mDL/Digital ID acceptance): https://www.tsa.gov/digital-id

**FIDO2 / convergence**
- FIDO Alliance — *Passkeys*: https://fidoalliance.org/passkeys/
- Thales — *FIDO Devices / Converged Badge* (one token for physical + PKI/FIDO): https://cpl.thalesgroup.com/access-management/authenticators/fido-devices
