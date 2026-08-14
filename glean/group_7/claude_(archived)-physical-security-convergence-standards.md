# physical-security-convergence-standards

**Category:** Science, Biology & Medicine
**Platform:** Claude (Archived)
**Original Path:** claude/archived-treearch/physical-security-convergence-standards

## Description
Physical-logical security convergence + the standards/compliance/governance layer above a PACS. TRIGGER: unifying PACS with enterprise IT identity (IdP/directory/PIAM, one credential for door + network, CSO/CISO convergence); US federal identity for physical access (HSPD-12, FIPS 201/PIV, PIV-I, CAC, FICAM, OMB M-19-17, SP 800-116 facility-area to auth-factor model, APL, Federal PKI/OCSP, SP 800-73/800-78 by name); standards bodies & frameworks (SIA; ONVIF Profiles A/C/D; ISO/IEC 27001 Annex A.7; SOC 2 CC6; PCI DSS req 9; ISO/IEC 60839 naming); PSIM / command center; visitor management & identity proofing; biometric-privacy law (Illinois BIPA) + surveillance-notice (jurisdiction- and date-dependent, not advice). SKIP: RF bands/physics -> rfid-fundamentals-bands-standards; NDEF/NFC -> nfc-ndef-protocols; chip/key -> contactless-smartcards-mifare-desfire; attacks -> rfid-nfc-access-attacks-defenses; Wiegand/OSDP wire -> access-credentials-wiegand-osdp; panels/door-hardware/topology -> pacs-access-control-architecture; mobile/Aliro/UWB/mDL -> mobile-credentials-aliro-uwb.

---

# Physical-Security Convergence & Standards

This skill owns the **governance, identity-integration, standards, and
compliance layer that sits ABOVE a physical access control system (PACS)** — how
physical security stops being a separate silo and joins enterprise IT identity,
which standards and frameworks govern it, and the privacy law that constrains it.

It does **not** cover the credential technology, the wire protocols, attacks, the
PACS hardware/topology, or mobile/UWB credentials — those are sibling skills (see
the SKIP list in the frontmatter and the cross-links inline). Think of this as
the "org chart + rulebook" spoke, where `pacs-access-control-architecture` is the
"how the door works" spoke.

> **Legal / regulatory disclaimer (read first).** The privacy and compliance
> material in §6 and the references is **general information, not legal advice**,
> and is **highly jurisdiction- and date-dependent**. Biometric-privacy statutes
> (e.g. Illinois BIPA), workplace surveillance-notice rules, and even the federal
> standards cited here change by legislative session, court ruling, and standards
> revision. Figures and provisions are stated **"as of June 2026"**; verify the
> currently adopted text with counsel and the relevant authority before relying
> on any of it. Compliance-framework control IDs (ISO/SOC 2/PCI) are summarized
> from secondary sources — confirm against the licensed standard text.

---

## 1. Physical-logical security convergence (the core idea)

**Convergence** is the unification of *physical* security (doors, badges,
cameras, guards) with *logical / IT* security (network login, IdP, directory,
data access) into one identity, one policy model, and often one org. Historically
these were two separate worlds — a **facilities/security** team ran badge access
on an isolated system, and a separate **IT/infosec** team ran network identity —
with separate credentials, separate enrollment, and no shared record of who a
person is. Convergence collapses that.

### 1.1 Why it matters

- **One identity, provisioned once.** When physical access derives from the same
  authoritative identity as IT access, an HR "hire / move / leave" event drives
  *both* the network account *and* the badge. The classic failure convergence
  fixes is the **terminated employee whose network account is disabled the same
  day but whose badge still opens the building for weeks** (the single most common
  audit exception in access reviews — see SOC 2 CC6 in §4).
- **The attack surface already merged.** Badge readers, cameras, and building
  systems are now IP devices on the network (often joined to Active Directory);
  an IoT/smart-building compromise can pivot into IT and vice-versa. The boundary
  between "physical" and "cyber" no longer matches the org chart, so splitting
  ownership leaves seams.
- **One credential for door + network.** The end-state many orgs want: the *same*
  credential a person uses to log into the network also opens the door (the US
  federal government mandates exactly this with the PIV card — §2). It removes a
  parallel provisioning pipeline, cuts cost, and means revocation is single-point.
- **Audit, regulation, and the board push it.** Frameworks increasingly treat
  physical and logical access as **one control family** (SOC 2 puts them in the
  same criterion, CC6; ISO 27001 has a physical-controls annex; NIST CSF 2.0
  spans both). Boards want single-point accountability for "security," full stop.

### 1.2 The org model: CSO / CISO convergence

- **CISO** (Chief Information Security Officer) traditionally owns *information /
  cyber* security. **CSO** (Chief Security Officer) is the broader role that can
  own *all* security — physical, personnel, investigations, and sometimes cyber.
- Convergence at the org level means **breaking down the CSO/CISO silo**: a single
  accountable security leader (or a deliberately integrated reporting line) over
  both physical and logical security, with unified policy, shared risk
  registers, and joint incident response. Organizations that converge report
  being more resilient and faster to detect/respond, because threats that cross
  the physical/cyber boundary hit one team, not two that have to coordinate.
- Caveat: titles are **not standardized**. In some firms the CSO reports to the
  CISO; in others the reverse; in many "CSO" and "CISO" are used interchangeably.
  The substance is *who owns the converged risk*, not the title.

### 1.3 PIAM — the technical glue

**PIAM (Physical Identity & Access Management)** is the software layer that makes
convergence real on the physical side. It is the *authoritative system for
physical access rights*, sitting between the enterprise identity sources and the
PACS:

- **Integrates with the HR/IT system of record** (Workday, SAP SuccessFactors,
  Oracle HCM) and the **IdP / directory** (Entra ID/Active Directory, Okta) so
  physical identities live in lockstep with the enterprise identity.
- **Automates the lifecycle:** **birthright / role-based access** provisions a new
  hire's badge rights on day one from their role; role changes re-provision;
  **deprovisioning** on termination revokes building access the moment the HR
  event fires (vendors claim large reductions in provisioning time and the
  closing of the "still-has-a-badge" gap).
- **Policy + workflow:** approval routing, area-owner attestation, **periodic
  access reviews / certification** of who can go where (the physical analogue of
  IT access-recertification), and audit evidence for the frameworks in §4.
- **Spans systems and sites:** one identity model over multiple PACS head-ends,
  buildings, and the **visitor management system** (§5).

> Distinguish PIAM from the PACS head-end: the **PACS head-end** (see
> `pacs-access-control-architecture`) decides and records door events; **PIAM**
> governs *identities and entitlements* across PACS, IdP, and HR. PIAM is the
> convergence/governance plane; the head-end is the access-control plane.

The deep federal instantiation of "one credential for door + network" is the PIV
card — that is §2 and `references/federal-identity-piv-ficam.md`.

---

## 2. US federal identity for physical access (HSPD-12 / PIV / FICAM)

The US federal government is the canonical converged-identity program: a single,
standardized, PKI-backed credential — the **PIV card** — that authenticates a
person to **both** facilities **and** information systems government-wide. This
section names the framework; the mechanism detail (authentication-mechanism
acronyms, PKI path validation, the data model) is in
`references/federal-identity-piv-ficam.md`.

### 2.1 The policy + standards stack

| Layer | What it is | Owns |
| --- | --- | --- |
| **HSPD-12** (2004) | Homeland Security Presidential Directive 12, "Policy for a Common Identification Standard for Federal Employees and Contractors" | The *mandate*: one secure, reliable, interoperable government-wide ID for physical + logical access |
| **FIPS 201** (current: **201-3**, Jan 2022) | NIST Federal Information Processing Standard implementing HSPD-12 | The *standard* for the **PIV** (Personal Identity Verification) credential: identity proofing, card, data model, auth mechanisms, issuer accreditation |
| **OMB M-19-17** (2019) | OMB ICAM policy ("Enabling Mission Delivery through Improved ICAM") | The modern *implementation policy*; reinforces HSPD-12/FIPS 201, expands permitted authenticators, mandates agency **ICAM** programs (updates the older M-05-24 era) |
| **FICAM** | Federal Identity, Credential, and Access Management — GSA program + architecture + playbooks at idmanagement.gov | The government-wide *architecture and guidance* (incl. **PACS 101**, the FICAM-compliant-PACS rules) |
| **NIST SP 800-116 Rev. 1** (2018) | "Guidelines for the Use of PIV Credentials in Facility Access" (Rev 1 retitled; original 800-116 withdrawn) | *How to use PIV in a PACS*: the facility-area → authentication-factor risk model (§2.3) |

Supporting NIST SP 800-series, **named here, detail in the reference**: **SP
800-73** (PIV card interfaces / data model), **SP 800-76** (biometric data
specs), **SP 800-78** (cryptographic algorithms & key sizes for PIV), **SP
800-79** (accreditation of PIV card issuers), **SP 800-156** (derived-credential
data model), **SP 800-157** (derived PIV credentials, e.g. on mobile).

### 2.2 PIV vs PIV-I vs CAC (don't conflate)

| Credential | Who | Issued under | Key point |
| --- | --- | --- | --- |
| **PIV** | Federal civilian employees & contractors | FIPS 201, Federal PKI | The government's own HSPD-12 credential |
| **CAC** (Common Access Card) | DoD military / civilian / eligible contractors | DoD's implementation of FIPS 201 | The DoD form of the PIV card (same standard, DoD issuance/PKI) |
| **PIV-I** (PIV-**Interoperable**) | **Non-federal** issuers (state/local govt, contractors, critical infrastructure) | Issued by non-federal entities to be *technically interoperable* with federal PIV, but chained to a **separate (non-Federal) PKI** trust hierarchy | Lets non-federal credentials be *trusted and interoperable* without being federally issued; widely adopted voluntarily |

All three follow the **same FIPS 201 technical model** (smart card, PKI certs,
biometrics, in-person proofing); they differ in **who issues** and **which PKI
trust chain** they anchor to.

### 2.3 The SP 800-116 facility-area → authentication-factor model (key concept)

This is the single most reusable idea from the federal model, and it applies
beyond government: **match the strength of authentication to the risk of the
area.** SP 800-116 Rev 1 / FICAM map facility **security areas** to required
**authentication factors** (factors = *have* the card / *know* the PIN / *are* the
biometric):

| Security area | Risk | Factors required | Typical PIV mechanism |
| --- | --- | --- | --- |
| **Unrestricted** | (public) | (no PIV auth) | — |
| **Controlled** | Lower | **1 factor** (have) | **PKI-CAK** (card auth key) or **SM-AUTH** |
| **Limited** | Moderate | **2 factors** (have + know) | **PKI-AUTH** + PIN |
| **Exclusion** | Highest | **3 factors** (have + know + are) | PKI-AUTH + PIN + **biometric** |

A **FICAM-compliant PACS** must, per OMB M-19-17 + SP 800-116 Rev 1: use
**high-assurance** PIV credentials, use **non-deprecated** auth mechanisms (FIPS
201-3 **removed CHUID** and **deprecated SYM-CAK and VIS** — never authenticate on
card-serial or by-eye alone), **validate credential status and authenticity**
(PKI path validation chaining to the **Federal Common Policy** root, with **OCSP**
preferred or **CRL** for revocation), **interoperate** with PIV from other
agencies, and use products on the **GSA FIPS 201 Approved Products List (APL)**.
Mechanism acronyms, PKI validation, and the E-PACS validation system are detailed
in `references/federal-identity-piv-ficam.md`.

> Why this matters even off-government: the **risk-tiered, multi-factor-at-the-
> door** pattern (and "validate the credential cryptographically, don't trust the
> card number") is the transferable lesson for any high-security commercial PACS,
> and it ties physical access back to the same PKI/identity assurance concepts IT
> security already uses.

---

## 3. Standards bodies (who writes the rules)

| Body | Scope | Relevant outputs |
| --- | --- | --- |
| **SIA** (Security Industry Association) | US security-industry trade association; **ANSI-accredited standards developer**; engages ISO/IEC | **OSDP** (reader↔controller protocol; now also **IEC 60839-11-5**), the legacy **26-bit Wiegand reader interface** spec, CCTV-to-access-control integration standards; runs the **Security Industry Standards Council (SISC)** consensus body and the **OSDP Verified** program. *(OSDP wire detail → `access-credentials-wiegand-osdp`.)* |
| **ONVIF** | Open standards for IP-based physical security device interoperability (cameras + access control), reduces vendor lock-in | **Profiles** (§4): **A** and **C** = access control; **D** = access-control peripherals; **S/T/G/M** = video/metadata |
| **NIST** | US federal standards (see §2) | FIPS 201, SP 800-series |
| **ISO/IEC** | International standards | **27001** (ISMS, Annex A.7 physical controls); **60839** series (alarm & electronic security systems, incl. access control) |
| **AICPA** | US accounting body that defines **SOC 2** | Trust Services Criteria (incl. CC6 logical+physical access) |
| **PCI SSC** (PCI Security Standards Council) | Payment-card industry | **PCI DSS** (Requirement 9 = physical access) |

---

## 4. Standards & compliance frameworks map (what each governs)

The crosswalk every converged program needs — *which standard says what about
physical access*. Control IDs are summarized from secondary sources (as of June
2026); confirm against the licensed standard text. Deeper notes per row in
`references/standards-frameworks-and-privacy.md`.

| Framework | Type | Physical-access scope | Key clauses / profiles |
| --- | --- | --- | --- |
| **ISO/IEC 27001:2022** (Annex A) | ISMS certification (international) | **Annex A clause 7 = Physical & environmental controls** (A.7.1 physical security perimeters; **A.7.2 physical entry** — entry controls, locks/cards/biometrics; A.7.3 securing offices/rooms; A.7.4 monitoring; A.7.5 protecting against physical/environmental threats; through A.7.14) | A.7.1, **A.7.2**, A.7.4 |
| **SOC 2** (AICPA TSC) | Attestation report (US, common in SaaS) | **CC6 = Logical *and* Physical Access Controls** — physical and logical access live in the **same criterion** (the convergence point in audit terms); covers restricting/granting/removing physical access, badge readers on server rooms | **CC6** (CC6.1–CC6.8) |
| **PCI DSS v4.0.1** | Mandatory for card-data handlers | **Requirement 9 = Restrict physical access to cardholder data** — facility entry controls, **visitor management** (escort, badges, logs), media handling/destruction, and **POI/POS device tamper-protection** (req 9.5) | **Req 9** (9.2 entry, 9.3 personnel/visitors, 9.4 media, 9.5 POI) |
| **ISO/IEC 60839 series** | Product/system standards (international) | **Electronic access control systems** requirements & guidance: **60839-11-1** (system & component requirements), **60839-11-2** (application/installation guidelines), **60839-11-5** (= **OSDP**) | 11-1, 11-2, 11-5 |
| **ONVIF Profiles** | Interoperability profiles | Device-level interop: **Profile A** (access-control configuration — credentials, schedules, door config), **Profile C** (door control + event/alarm at the system level), **Profile D** (access-control peripherals — readers, locks, sensors) | **A, C, D** |
| **NIST FICAM / SP 800-116** | US federal (see §2) | PIV-in-PACS, facility-area → factor model, APL | §2.3 |

How they relate: **ISO 27001 and SOC 2 are the broad "is your security program
sound" frameworks** (physical access is one control family inside each); **PCI DSS
is a narrow mandate** for one data type; **ISO 60839 and ONVIF are the
product/interop standards** the installed system is built to; **FICAM/NIST is the
US-federal regime**. A converged program typically answers to *several at once*
and wants one PIAM/PACS that produces evidence for all of them.

---

## 5. PSIM and the command center

**PSIM (Physical Security Information Management)** is middleware that **integrates
many disparate, unconnected security systems** — video surveillance (VMS), access
control, intrusion, fire, sensors, analytics, building systems — into **one
correlated operating picture** in a **command / control center** (control room,
GSOC, C4I). It is the *situational-awareness and response-orchestration* layer
above the individual systems.

- **What it does:** collects + **correlates events** across systems (e.g. a
  door-forced-open from the PACS + motion on the matching camera + an intrusion
  zone = one *situation*, not three alarms); presents a unified UI/map; and drives
  **SOP-based response workflows** (step-by-step operator guidance, dispatch,
  escalation) so an operator resolves a situation instead of watching panes.
- **Vs a PACS head-end or a VMS:** a head-end manages *access*; a VMS manages
  *video*; **PSIM sits above both and ties everything together** for the command
  center. PSIM is integration/orchestration **middleware**, typically vendor-
  agnostic (it consumes the others via their APIs / ONVIF / OSDP-fed events).
- **Vs PIAM:** **PIAM governs identities/entitlements** (who *should* have access);
  **PSIM governs real-time events/situations** (what is *happening now*). They are
  complementary planes, not competitors.
- **Landscape (examples, not endorsements):** Everbridge (Control Center),
  Advancis (WinGuard), and others; the category overlaps increasingly with
  "converged/unified security platforms." Market labels shift — treat specific
  product claims as time-sensitive.

---

## 6. Visitor management, identity proofing & privacy/regulatory

### 6.1 Visitor management & identity proofing for physical access

A **visitor management system (VMS)** extends the converged identity model to
**non-employees** (guests, contractors, vendors, auditors):

- **Pre-registration:** a host requests/approves a visit before arrival;
  pre-approval cuts wait time and moves the access decision earlier.
- **Identity proofing at check-in:** verifying the visitor *is who they claim* —
  typically **scanning a government ID** (driver's license/passport) and capturing
  a photo; higher-assurance sites add document authentication.
- **Watchlist / screening:** automatic check against internal deny-lists or, for
  defense/aerospace/government, **national or compliance watchlists**, with alerting.
- **Temporary, time-bound, zone-scoped credentials:** the VMS issues a badge that
  is **active only for approved areas and the approved window** and **auto-revokes
  at checkout** — written into the *same* cardholder/audit model as employees
  (this is the convergence point; PCI DSS req 9.3 and ISO A.7.2 expect exactly
  this escort/badge/log discipline).
- **Identity proofing vs authentication:** *proofing* establishes who someone is
  at enrollment (the harder problem for visitors); *authentication* re-verifies a
  known identity at the door later. Government visitor proofing is deliberately
  stronger than commercial. (For the assurance-level vocabulary — NIST SP 800-63
  IAL/AAL — see `references/standards-frameworks-and-privacy.md`.)

### 6.2 Privacy & regulatory (US-centric; jurisdiction- and date-dependent; NOT legal advice)

Converged systems collect **biometrics, video, and movement logs** — among the
most regulated data classes. **Re-read the §-top disclaimer.** Highlights, **as of
June 2026**, verify currency with counsel:

- **Biometric-privacy law — Illinois BIPA** is the landmark US statute (740 ILCS
  14). It governs **private entities** collecting **biometric identifiers**
  (fingerprint, face/hand geometry, iris/retina, voiceprint) — squarely
  implicating biometric door readers and facial-recognition access. Core duties:
  **informed written consent before collection**, a published **retention/
  destruction schedule**, and limits on disclosure/profiting. It carries a
  **private right of action** with statutory damages (historically **$1,000
  negligent / $5,000 intentional or reckless** per violation), which drove heavy
  class-action litigation. **2024 amendment SB 2979** narrowed exposure — repeated
  collection of the *same* biometric from the *same* person is **a single
  violation** (post-*Cothron*), and "written release" now includes an **electronic
  signature**. **Texas (CUBI)** and **Washington** have biometric statutes too,
  but BIPA's private right of action makes it the dominant risk. **State biometric
  law is a patchwork and actively changing — check the specific state and the
  latest amendments/rulings.**
- **Video surveillance & notice.** No single US federal video-privacy regime for
  workplace CCTV; it's a mix of **federal wiretap law (ECPA)** — which restricts
  **audio** recording (one-party consent federally; some states **two-party**,
  e.g. CA/MA) far more than silent video — plus **state notice rules**. Several
  states require **written notice of electronic monitoring** (e.g. **NY**'s 2021
  law, effective May 2022; **CT**, **DE**), and **all** states bar cameras in
  areas with a **reasonable expectation of privacy** (restrooms, changing rooms).
  Practical baseline: post conspicuous signage, give written notice, avoid
  private areas, and **never assume video+audio is treated like video alone**.
- **General privacy regimes that reach physical-security data.** State
  comprehensive privacy laws (e.g. California **CCPA/CPRA**) and, for
  multinationals, the EU **GDPR** treat biometrics/video/access logs as personal
  (often *special-category*) data with notice, purpose-limitation, retention, and
  data-subject-rights obligations. Scope/applicability vary widely — out of depth
  here beyond flagging them.

> Every item in §6.2 is **general information, jurisdiction-specific, and
> time-sensitive**. It is **not legal advice**; engage counsel for any actual
> deployment, and re-verify because these statutes and the case law move fast.

---

## Sources

Authoritative (NIST/idmanagement.gov/standards bodies/AICPA-PCI) and
reputable-secondary references consulted (accessed June 2026). Legal/compliance
items are secondary summaries and are jurisdiction- and date-dependent — verify
against primary text.

1. **NIST FIPS 201-3, "Personal Identity Verification (PIV) of Federal Employees
   and Contractors" (Jan 24, 2022)** — the PIV standard; expanded derived
   credentials/federation, supervised remote identity proofing, **removed CHUID**,
   **deprecated SYM-CAK + VIS**, added SM-AUTH:
   https://csrc.nist.gov/pubs/fips/201-3/final and
   https://csrc.nist.gov/News/2022/fips-201-3-nist-revises-piv-standard
2. **HSPD-12 (Aug 27, 2004)** — the directive mandating a common federal
   identification standard for physical + logical access (referenced via FIPS
   201-3 and NIST PIV project): https://csrc.nist.gov/projects/piv
3. **NIST SP 800-116 Rev. 1, "Guidelines for the Use of PIV Credentials in
   Facility Access" (June 2018)** — PIV-in-PACS; facility-area → authentication-
   factor risk model; original SP 800-116 withdrawn:
   https://csrc.nist.gov/pubs/sp/800/116/r1/final
4. **idmanagement.gov — FICAM program, architecture, and "Physical Access Control
   Systems 101"** — FICAM-compliant-PACS rules (high-assurance creds, non-
   deprecated mechanisms, status validation, interoperability, GSA APL), security-
   area → factor mapping, Federal Common Policy chaining, OCSP/CRL:
   https://www.idmanagement.gov/ficam/ and https://www.idmanagement.gov/university/pacs/
5. **OMB M-19-17, "Enabling Mission Delivery through Improved Identity,
   Credential, and Access Management" (2019)** — federal ICAM policy reinforcing
   HSPD-12/FIPS 201 and mandating agency ICAM:
   https://www.whitehouse.gov/wp-content/uploads/2019/05/M-19-17.pdf
6. **NIST — SP 800-73-5 (PIV interfaces/data model) and SP 800-78-5
   (cryptographic algorithms & key sizes), revised July 2024** — companion PIV
   specs (named-level):
   https://www.nist.gov/news-events/news/2024/07/personal-identity-verification-piv-interfaces-cryptographic-algorithms-and
   and https://csrc.nist.gov/pubs/sp/800/78/5/final
7. **PIV vs CAC vs PIV-I** — non-federal PIV-Interoperable credentials; same FIPS
   201 model, separate PKI hierarchies; DoD CAC as the DoD PIV form:
   https://legalclarity.org/piv-and-cac-card-differences-eligibility-and-issuance/
8. **Security Industry Association (SIA) — Industry Standards & OSDP** — ANSI-
   accredited SDO; OSDP published as **IEC 60839-11-5** (May 2020); SISC consensus
   body; OSDP Verified; 26-bit Wiegand interface standard:
   https://www.securityindustry.org/industry-standards/ and
   https://www.securityindustry.org/industry-standards/open-supervised-device-protocol/
9. **ONVIF — Profiles** — Profile **A** & **C** for access control, **D** for
   access-control peripherals; **S/T/G/M** for video; cross-vendor interop:
   https://www.onvif.org/profiles/
10. **ISO/IEC 27001:2022 Annex A clause 7 (Physical controls)** — A.7.1 perimeters,
    **A.7.2 physical entry**, A.7.4 monitoring, A.7.5 environmental threats
    (secondary summaries): https://www.isms.online/iso-27001/annex-a-2022/7-2-physical-entry-2022/
    and https://www.isms.online/iso-27001/annex-a-2022/7-1-physical-security-perimeters-2022/
11. **SOC 2 Trust Services Criteria — CC6 "Logical and Physical Access Controls"**
    (CC6.1–CC6.8); physical + logical access in one criterion; access-review/
    termination exceptions are the most common findings:
    https://secureframe.com/hub/soc-2/common-criteria
12. **PCI DSS v4.0 Requirement 9 — "Restrict physical access to cardholder data"**
    — 9.2 entry controls, 9.3 personnel/visitor management, 9.4 media, 9.5 POI
    tamper-protection: https://www.isms.online/pci-dss/requirement-9/ and
    https://www.herodevs.com/blog-posts/pci-dss-4-0-requirement-9-how-to-restrict-physical-access-to-cardholder-data
13. **ISO/IEC 60839-11-1 / -11-2 / -11-5 (Alarm and electronic security systems —
    Electronic access control systems)** — system/component requirements,
    application guidelines, and OSDP:
    https://webstore.iec.ch/en/publication/33414 and
    https://standards.globalspec.com/std/1646780/en-60839-11-1
14. **PSIM** — definition (integration middleware correlating disparate security
    systems into one command-center UI/workflow); vendor landscape (Everbridge,
    Advancis, et al.): https://en.wikipedia.org/wiki/Physical_security_information_management
    and https://www.ifsecglobal.com/psim/
15. **PIAM (Physical Identity & Access Management)** — authoritative physical-
    identity lifecycle layer; HR/IdP integration (Workday/SAP/Entra/Okta);
    birthright/role-based provisioning + deprovisioning; access certification:
    https://www.rightcrowd.com/resources/blog/what-is-physical-identity-and-access-management-piam/
    and https://www.hidglobal.com/solutions/physical-identity-access-management
16. **Physical/logical convergence & CSO/CISO** — convergence rationale, merged
    attack surface, single-credential goal, org-silo breakdown, resilience
    benefit: https://www.csoonline.com/article/514481/strategic-planning-erm-physical-and-it-security-convergence-the-basics.html
    and https://www.rightcrowd.com/resources/blog/the-convergence-of-physical-and-logical-access-control/
17. **Illinois BIPA (740 ILCS 14)** — biometric consent/retention/private right of
    action; statutory damages; **2024 SB 2979** (single-violation + electronic
    signature) post-*Cothron*; 2024–25 litigation context:
    https://www.aclu-il.org/campaigns-initiatives/biometric-information-privacy-act-bipa/
    and https://www.gtlaw.com/en/insights/2024/8/bipa-update-illinois-limits-liability-and-clarifies-electronic-consent-for-biometric-data-collection
18. **Workplace video-surveillance / notice law** — ECPA audio vs silent-video
    distinction; two-party-consent states (CA/MA); state written-notice rules
    (NY 2021, CT, DE); reasonable-expectation-of-privacy bar (secondary summaries,
    US-centric, state-dependent):
    https://www.workplacefairness.org/workplace-surveillance/ and
    https://www.rhombus.com/blog/are-you-required-to-notify-employees-or-customers-that-you-have-security-cameras/
19. **Visitor management & identity proofing (government/contractor)** — pre-
    registration, ID-scan proofing, watchlist screening, temporary zone/time-bound
    credentials, integration with PACS:
    https://www.avigilon.com/blog/government-visitor-management

**Uncertainty / verify-before-use notes:**
- **All §6 legal/privacy content is jurisdiction- and date-dependent and is NOT
  legal advice.** BIPA damages, the SB 2979 single-violation rule, surveillance-
  notice requirements, and which states are two-party-consent all change by
  session and court ruling — verify the current statute/case law per jurisdiction.
- **Compliance control IDs (ISO 27001 A.7.x, SOC 2 CC6.x, PCI DSS 9.x, ISO 60839
  parts)** are summarized from **secondary sources**; the licensed standard text
  governs — confirm clause numbers and wording before citing in an audit context.
- **SP 800-116 facility-area → factor mapping** (Controlled=1 / Limited=2 /
  Exclusion=3) and the deprecated-mechanism status (CHUID removed; SYM-CAK, VIS
  deprecated) reflect FIPS 201-3 / SP 800-116 Rev 1 **as of June 2026**; NIST
  revises these — check for a newer FIPS 201 revision before relying.
- **PSIM/PIAM vendor names** are illustrative, not endorsements, and the product
  landscape (and the PSIM-vs-"unified platform" labeling) shifts over time.
- Mechanism-acronym and PKI-validation specifics live in
  `references/federal-identity-piv-ficam.md`; the original SP 800-116 PDF was
  image-encoded and not machine-readable in this build, so those specifics are
  grounded in the FICAM PACS-101 page + NIST announcements rather than the bound
  SP 800-116 text — verify against the PDF for audit-grade use.