# Federal identity for physical access — PIV / FICAM mechanism detail

Companion to `physical-security-convergence-standards` §2. The SKILL.md names the
HSPD-12 / FIPS 201 / FICAM stack and the facility-area → factor model; this file
holds the mechanism-level detail: the PIV authentication-mechanism acronyms, the
PKI path-validation / E-PACS validation chain, and the data-model / crypto
companion specs by name.

> Same disclaimer as the parent skill: this is **general information as of June
> 2026, not legal/compliance advice**. NIST revises FIPS 201 and the SP 800-series
> on its own cadence — verify against the current published text before relying on
> any acronym status or numeric below. The original **SP 800-116** PDF was
> image-encoded and not machine-readable in this build, so the mechanism details
> here are grounded in the **idmanagement.gov FICAM PACS-101 page + NIST FIPS
> 201-3 announcements**, cross-checked, rather than transcribed from the bound SP
> 800-116 text. Treat as a study map, not an audit citation.

---

## 1. PIV authentication mechanisms (acronym map)

A PIV card can be authenticated several ways; each provides a different number of
**factors** (*have* the card / *know* the PIN / *are* the biometric) and runs over
the **contact** or **contactless** interface. The PACS designer picks the
mechanism to match the area's risk (§3). Names and status **as of FIPS 201-3 /
June 2026**:

| Mechanism | Full name | Factors | Interface | Status (FIPS 201-3) |
| --- | --- | --- | --- | --- |
| **VIS** | Visual (human guard reads the card / photo) | ~0–1 (visual) | n/a (by eye) | **Deprecated** — not a cryptographic authentication |
| **CHUID** | Cardholder Unique Identifier | 1 (have, but **not** cryptographically strong) | contact/contactless | **Removed** in 201-3 (was deprecated in 201-2) — do not use |
| **PKI-CAK** | PKI using the **Card Authentication Key** (asymmetric) | 1 (have) | contact **or contactless** | Current — the standard **1-factor** door mechanism (works contactless, no PIN) |
| **SYM-CAK** / SM-AUTH | Symmetric Card Auth Key / **Secure Messaging** authentication | 1 (have) | contactless | **SYM-CAK deprecated** in 201-3 (use not prohibited yet, slated for removal); **SM-AUTH** added in 201-3 as a current secure-messaging mechanism |
| **BIO** | Biometric (fingerprint) match, **PIN-released** | 2 (have + are; PIN releases the template) | contact | Current |
| **BIO-A** | Biometric match **with attended/Attendant** supervision | 3 (have + know + are, attended) | contact | Current — highest assurance with a guard present |
| **OCC-AUTH** | **On-Card Comparison** (match-on-card biometric) | 2 (have + are) | contact (and 201-3 broadens) | Current — biometric compared on the card itself |
| **PKI-AUTH** | PKI using the **PIV Authentication Key** (asymmetric, **PIN-protected**) | 2 (have + know) | contact (PIN required) | Current — the standard **2-factor** mechanism |

Mnemonic for the common three a PACS uses:
- **PKI-CAK** = "tap, 1 factor" (contactless, no PIN) → **Controlled** area.
- **PKI-AUTH** = "insert + PIN, 2 factors" → **Limited** area.
- **PKI-AUTH + PIN + biometric** = 3 factors → **Exclusion** area.

The decisive split: the **CAK** lives on the card and needs **no PIN** (so it can
run contactless for throughput), while the **PIV Auth key** is **PIN-protected**
(so it inherently adds the *know* factor). Both are real asymmetric PKI keys whose
certificates are validated against the Federal PKI — which is the point of §2: a
PIV reader doesn't trust the card *number*, it cryptographically proves the card
holds a valid, unrevoked key that chains to the federal trust root.

---

## 2. PKI path validation & the E-PACS validation chain

FICAM compliance hinges on **validating the credential**, not reading an
identifier. The flow at a high-assurance door:

1. **Read the certificate** from the card (CAK cert for 1-factor; PIV-Auth cert
   for 2-factor, which also requires the PIN).
2. **Cryptographic challenge** — the reader/controller has the card prove
   possession of the private key (signs a nonce), so a cloned identifier alone
   fails.
3. **Path validation** — build and validate the certificate chain up to the
   **Federal Common Policy CA** (the Federal PKI / FCPCA trust anchor). A cert
   that does not chain to the federal root is not trusted.
4. **Revocation check** — confirm the cert is not revoked via **OCSP** (Online
   Certificate Status Protocol — preferred, real-time) or a **CRL** (Certificate
   Revocation List — loaded periodically, e.g. daily, for offline validation).
5. **Authorization** — only *after* the credential is authenticated does the PACS
   check whether *this* identity is **registered** and **entitled** to this door
   (the local access-control decision — see `pacs-access-control-architecture`).

**Authentication vs authorization** is the recurring distinction: FICAM/PIV makes
**authentication** strong and government-wide-interoperable (any agency's PIV can
be *proven genuine* anywhere); **authorization** (does this person get into *this*
building) stays local to the site's PACS / PIAM. Interoperability means agency B's
reader can validate agency A's PIV; it does **not** mean agency A's card
automatically opens agency B's doors.

### E-PACS (Enterprise PACS)

A standalone PACS can't continuously check federal revocation. An **E-PACS**
connects to the enterprise network and a **validation system / validation server**
that performs the path-validation and OCSP/CRL polling centrally and pushes
validated status to controllers — so credentials are checked against current
federal revocation continuously, not just at enrollment. This is what makes
"validate status and authenticity" operationally real at scale, and it is itself
an instance of the convergence theme (the PACS now depends on enterprise network
+ PKI services).

---

## 3. Facility-area → authentication-factor model (expanded)

The SP 800-116 Rev 1 / FICAM risk model (parent §2.3), with the mechanism mapping
from §1 of this file:

| Security area | Risk | Factors | Acceptable PIV mechanisms |
| --- | --- | --- | --- |
| **Unrestricted** | public; no controlled boundary | — | none |
| **Controlled** | lower-risk controlled space | **1** (have) | **PKI-CAK** (or SM-AUTH) — contactless OK |
| **Limited** | moderate risk | **2** (have + know) | **PKI-AUTH** + PIN |
| **Exclusion** | highest risk / most sensitive | **3** (have + know + are) | **PKI-AUTH** + PIN + **biometric** (BIO/BIO-A/OCC-AUTH) |

Design rules a FICAM-compliant PACS follows (per OMB M-19-17 + SP 800-116 Rev 1):
- **No deprecated mechanisms** — never authenticate on **VIS** (by-eye) or
  **CHUID** (card identifier); both are out. Prefer PKI mechanisms.
- **More factors for more sensitive areas** — escalate from CAK → PIV-Auth+PIN →
  +biometric as area risk rises.
- **Validate every time** — path-validate to the Federal Common Policy root and
  check revocation (OCSP/CRL); don't cache "valid" indefinitely.
- **Buy from the GSA FIPS 201 Approved Products List (APL)** — components are
  tested for FIPS 201 conformance + interoperability; agencies are required to
  source from it.
- **Interoperate** — accept and validate PIV from other federal agencies.

---

## 4. Companion NIST SP 800-series (named, for orientation)

These define the *parts* of the PIV system FIPS 201 points to. Named here so the
acronyms are recognizable; consult NIST for the authoritative text.

| Spec | Title (abbrev.) | Governs |
| --- | --- | --- |
| **SP 800-73** (current Parts 1–3, **-5**) | Interfaces for PIV | The **card data model** (Part 1), **card edge / card interface** (Part 2), **client API** (Part 3) — how data is laid out on the card and read |
| **SP 800-76** | Biometric Specifications for PIV | The **biometric data** formats (fingerprint, face, iris) captured/stored |
| **SP 800-78** (current **-5**) | Cryptographic Algorithms & Key Sizes for PIV | Allowed **algorithms/key sizes** (e.g. RSA 2048/3072, ECDSA P-256/P-384; deprecates 3TDEA); aligns crypto with FIPS 201-3 |
| **SP 800-79** | Guidelines for the Accreditation of PIV Card Issuers | Accreditation of the **organizations that issue** PIV cards (the issuer-trust leg) |
| **SP 800-156** | Representation of PIV Chain-of-Trust for Import/Export | Data model for the **chain-of-trust** record (enables derived credentials) |
| **SP 800-157** | Guidelines for **Derived PIV Credentials** | Issuing a **derived** credential (e.g. onto a mobile device) from a valid PIV — the bridge toward mobile credentials |

> Mobile / derived-credential *use* and the broader mobile-credential ecosystem
> (Aliro, UWB, mDL) are **out of scope here** — see the `mobile-credentials-
> aliro-uwb` sibling. SP 800-157/156 are named only to show where the federal
> standard hands off to derived credentials.

---

## 5. Sources

See the parent `SKILL.md` **Sources** section (entries 1–7 cover FIPS 201-3,
HSPD-12, SP 800-116 Rev 1, FICAM/idmanagement.gov PACS-101, OMB M-19-17, SP
800-73-5 / 800-78-5, and PIV/CAC/PIV-I). The mechanism-status and PKI-validation
specifics in this file are grounded primarily in the **idmanagement.gov FICAM
PACS-101 page** and the **NIST FIPS 201-3 announcement** (CHUID removed; SYM-CAK +
VIS deprecated; SM-AUTH added), cross-checked against secondary FIPS-201 summaries
— **not** transcribed from the bound SP 800-116 PDF (image-encoded in this build).
Verify acronym status and the area→factor mapping against the published FIPS 201-3
and SP 800-116 Rev 1 text for any audit-grade use.
