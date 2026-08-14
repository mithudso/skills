# physical-access-control

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Claude Code
**Original Path:** claude-code/physical-access-control

## Description
Physical access control hub — badges, readers, controllers, and the credential/RF stack beneath them; routes to spokes. TRIGGER: RFID fundamentals (LF/HF/UHF bands, coupling, EPC Gen2, anti-collision); NFC & NDEF (NFC Forum, tag types, HCE/SE, phone NFC APIs); contactless smart cards (ISO 14443, APDU, MIFARE Classic/Plus/Ultralight, DESFire, on-card crypto); reader-to-controller credentials (Wiegand formats & wire, OSDP Secure Channel, HID Prox/iCLASS/SEOS); PACS architecture (controllers, door hardware, fail-safe/secure, egress/NFPA codes, head-end, anti-passback); attacks & defenses (skim/clone/relay, Crypto-1, ESPKey, Proxmark/Flipper, hardening); mobile credentials (phone-as-badge, Apple/Google Wallet, Aliro, UWB ranging, mDL); physical-logical convergence & standards (PIV/FIPS 201/FICAM, ONVIF, BIPA, PSIM). SKIP: web/app/cloud security → security-review; Web Crypto / vault code → webcrypto-vault-reviewer; blockchain & crypto-primitive math → blockchain.

---

# Physical Access Control (hub)

The front door for **how physical access control works end to end** — the RF and credential stack (RFID → NFC → smart cards → Wiegand/OSDP), the PACS that reads them (controllers, door hardware, head-end), the attack/defense model, and the mobile-credential and physical-logical-convergence frontier. This hub routes to reference spokes.

**When activated:** match the question to a row, then **Read the listed `references/` file** before answering at depth — the table is only a router.

## Routing table

| Spoke | Use when | Reference file |
| --- | --- | --- |
| `rfid-fundamentals-bands-standards` | RFID system anatomy, frequency-band choice (LF/HF/UHF/microwave), coupling physics, EPC Gen2 / ISO 18000, anti-collision; which band fits access vs supply chain. | `references/rfid-fundamentals-bands-standards.md` |
| `nfc-ndef-protocols` | NFC as an HF-RFID subset and the NDEF data format: NFC Forum specs, the three modes, HCE vs secure element, tag types 1–5, phone NFC APIs. | `references/nfc-ndef-protocols.md` |
| `contactless-smartcards-mifare-desfire` | ISO 14443 PICCs, the 7816-4 APDU model, MIFARE (Classic/Ultralight/Plus) and DESFire EV1/2/3 — memory model, keys, on-card crypto, secure messaging. | `references/contactless-smartcards-mifare-desfire.md` |
| `access-credentials-wiegand-osdp` | Credential data formats and reader-to-controller protocols: Wiegand formats, the Wiegand wire, OSDP Secure Channel, HID Prox/iCLASS/SE/SEOS families. | `references/access-credentials-wiegand-osdp.md` |
| `pacs-access-control-architecture` | The PACS stack end to end: controllers/panels, door hardware (strikes/maglocks/REX/DPS), fail-safe vs fail-secure + egress/life-safety codes, head-end, integration. | `references/pacs-access-control-architecture.md` |
| `rfid-nfc-access-attacks-defenses` | Threat model, attacks, and defenses for card-based access (eavesdrop/skim/clone/replay/relay/downgrade), Crypto-1, ESPKey, Proxmark/Flipper, reader hardening — for defense/authorized assessment. | `references/rfid-nfc-access-attacks-defenses.md` |
| `mobile-credentials-aliro-uwb` | Phone-as-credential frontier: NFC/BLE/UWB, Apple/Google/Samsung Wallet badges, Aliro (CSA), UWB secure ranging & anti-relay, CCC Digital Key, mDL. | `references/mobile-credentials-aliro-uwb.md` |
| `physical-security-convergence-standards` | Physical-logical convergence and the standards/compliance layer: PIV/FIPS 201/FICAM, SIA, ONVIF, ISO 27001 A.7, PSIM, visitor management, biometric-privacy law (BIPA). | `references/physical-security-convergence-standards.md` |

## Cross-cutting notes

- **Layer order.** RF/physics (`rfid-fundamentals`) → data format (`nfc-ndef`, `contactless-smartcards`) → reader-to-panel wire (`access-credentials-wiegand-osdp`) → system (`pacs-access-control-architecture`). Attacks (`rfid-nfc-access-attacks-defenses`) cut across all four; mobile/convergence sit above the system.
- **Defense, not offense.** The attacks spoke is for hardening and authorized assessment. App/web/cloud security is a neighbor, not here → `security-review`.
- **Card crypto vs Web Crypto.** On-card symmetric crypto (DES/AES on MIFARE/DESFire) lives in `contactless-smartcards-mifare-desfire`; Web Crypto API / vault implementation code → `misc-catch-all` (references/webcrypto-vault-reviewer.md).

<!-- cross-hub-map -->