# NFC Forum specification inventory

Companion to the SKILL.md "NFC Forum governance" section. Sourced from the NFC
Forum Technical Specifications index (versions as listed there, mid-2026);
versions advance — re-check the index before citing a version as current.

## The protocol stack (top → bottom)

| Spec | Version | Role |
| --- | --- | --- |
| NFC Controller Interface (NCI) | 2.3 | Standard interface between the NFC controller and the host application processor (2.3 adds NFC WLC + removal detection) |
| NFC Digital Protocol | 2.4 | Implementation layer above ISO/IEC 18092 & 14443; harmonizes NFC-A/B/F (2.0 added NFC-V + Active Comm Mode; 2.4 aligns EGT to EMVCo, T5T special frames to ISO/IEC 15693) |
| NFC Activity | 2.3 | How the Digital Protocol is sequenced to set up communication |
| NFC Analog | 3.0 | RF/analog characteristics of the interface (3.0 introduced the 20 mm operating volume); covers NFC-A/B/F |
| Profiles | 1.1 | How Activities combine for a use case |

## Data format & tags

| Spec | Role |
| --- | --- |
| Data Exchange Format (NDEF) | The common record/message data format |
| Record Type Definition (RTD) | Rules for building standard record types |
| Text RTD / URI RTD / Smart Poster RTD | Well-known record types (see `ndef-format.md`) |
| Multiple URI RTD | Efficiently encodes multiple URIs in one record |
| Signature RTD | Signs single or multiple NDEF records |
| Device Information RTD | Conveys model/identity info |
| Type 1–5 Tag specs | NDEF mapping per air interface (see `tag-types-1-5.md`) |

## Peer-to-peer & tag exchange

| Spec | Version | Role |
| --- | --- | --- |
| Logical Link Control Protocol (LLCP) | 1.4 | OSI L2 peer-to-peer transport, based on IEEE 802.2 (1.4 uses NAP 1.0 for secured transfer) |
| Simple NDEF Exchange Protocol (SNEP) | — | Exchanges NDEF messages over LLCP connection-oriented mode |
| Tag NDEF Exchange Protocol (TNEP) | 1.0 | Bidirectional NDEF exchange on Type 2/3/4/5 tags |

> P2P / LLCP-SNEP is effectively obsolete on phones (Android Beam removed in
> Android 14; Apple never shipped it). NFC's surviving "tap to share" role is
> **Connection Handover**, not bulk P2P.

## Application / advanced specs

| Spec | Version | Role |
| --- | --- | --- |
| Connection Handover | 1.5 | Bootstraps a Bluetooth/Wi-Fi link via a tap (first to use TNEP) |
| NFC Authentication Protocol | 1.0 (Dec 2022) | Authentication and secured data transfer |
| Wireless Charging (WLC) | 2.0 | Wireless charging of small battery devices |
| Bluetooth Secure Simple Pairing | 1.3 | NFC ↔ Bluetooth pairing handover detail |
| Personal Health Device Communication (PHDC) | 1.2 | Transport for ISO/IEEE 11073-20601 health devices |

## Underlying ISO/JIS standards the specs reference

- **ISO/IEC 14443** (Type A and B) — proximity contactless cards → NFC-A / NFC-B.
- **ISO/IEC 18092 (NFCIP-1)** — the NFC peer protocol; also references FeliCa.
- **ISO/IEC 15693** — vicinity cards → NFC-V.
- **JIS X 6319-4 (FeliCa)** — Sony's standard → NFC-F / Type 3 Tag.
- **ISO/IEC 7816-4** — APDU command/response set, carried over ISO-DEP in CE/Type 4.
- **IEEE 802.2** — the basis for LLCP.

## Governance

The **NFC Forum** (founded 2004 by Nokia, Philips/NXP, and Sony) owns the
interoperability specs above and runs the certification program. Specs are
purchased/downloaded from the NFC Forum; the index page is the authority for
current version numbers.
