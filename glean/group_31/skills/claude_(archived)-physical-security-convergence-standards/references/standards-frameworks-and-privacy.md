# Standards, compliance frameworks & privacy — detail

Companion to `physical-security-convergence-standards` §3–§6. The SKILL.md gives
the crosswalk table and the headline privacy items; this file expands each
framework's physical-access scope and the privacy/regulatory landscape.

> **Legal/compliance disclaimer (carries from the parent).** This is **general
> information as of June 2026, not legal or compliance advice**, and is **heavily
> jurisdiction- and date-dependent**. Compliance control IDs are summarized from
> **secondary sources**; the **licensed standard text governs** — confirm clause
> numbers/wording before citing in any audit. Privacy statutes and case law change
> by legislative session and court ruling — re-verify per jurisdiction with
> counsel.

---

## 1. Standards bodies & interop standards

### SIA (Security Industry Association)
- US security-industry trade association and an **ANSI-accredited Standards
  Developing Organization**; also engages with ISO/IEC for global standards.
- **Outputs relevant here:** **OSDP** (Open Supervised Device Protocol, the
  reader↔controller standard — published internationally as **IEC 60839-11-5**,
  2020; current SIA OSDP rev as of late 2024), the legacy **26-bit Wiegand reader
  interface** spec, and CCTV-to-access-control integration standards.
- **Governance:** the **Security Industry Standards Council (SISC)** is the
  balanced-consensus body that votes on proposed standards; SIA also runs the
  **OSDP Verified** conformance program.
- *OSDP wire/secure-channel detail is owned by the `access-credentials-wiegand-
  osdp` sibling — this skill only places SIA as the standards body.*

### ONVIF (interoperability profiles)
- Open standards for IP-based physical-security devices so equipment from
  different vendors interoperates (reduces vendor lock-in). Conformance is by
  **Profile**.
- **Access-control-relevant profiles:**
  - **Profile A** — access-control **configuration**: managing credentials,
    schedules, access rules, and door configuration on networked access-control
    devices.
  - **Profile C** — access-control **door control + events**: door
    monitoring/control and event/alarm handling at the system level.
  - **Profile D** — access-control **peripherals**: readers, locks, sensors, PINpads
    (the device-edge of the access system).
- **Video profiles (named, not in scope):** **S** (streaming), **T** (advanced
  streaming/analytics), **G** (recording/storage), **M** (metadata/analytics).
- Practical read: a converged shop specifies ONVIF profiles to keep cameras,
  readers, and controllers swappable across vendors and to let the PACS and VMS
  exchange events.

### ISO/IEC 60839 series (product/system standards)
- "Alarm and electronic security systems" — the **-11-x** parts cover **electronic
  access control systems (EACS)**:
  - **60839-11-1** — system and component **requirements** (minimum functionality,
    performance, test methods for physical entry/exit control).
  - **60839-11-2** — **application guidelines** (planning, installation,
    commissioning, maintenance, documentation).
  - **60839-11-5** — **OSDP** (the SIA protocol, internationalized).
- These are what the installed *system* is built and tested to, complementary to
  the management/governance frameworks below.

---

## 2. Compliance / governance frameworks — physical-access scope

### ISO/IEC 27001:2022 (ISMS certification)
- The international **Information Security Management System** standard; Annex A
  lists the controls. The **2022** revision reorganized Annex A into **4 themes**;
  **clause 7 = Physical controls** (14 controls, A.7.1–A.7.14).
- **Physical-access-relevant controls:**
  - **A.7.1 Physical security perimeters** — define/protect perimeters (walls,
    barriers, controlled boundaries) around information assets.
  - **A.7.2 Physical entry** — entry controls at secure-area access points
    (locks/keys, card readers, biometrics, guards, monitoring); only authorized
    people enter. *This is the core "access control" control.*
  - **A.7.3** securing offices/rooms/facilities; **A.7.4** physical security
    **monitoring** (surveillance/alarms); **A.7.5** protecting against physical &
    environmental threats; **A.7.6** working in secure areas; A.7.7 clear
    desk/screen; A.7.8–A.7.14 equipment siting, utilities, cabling, maintenance,
    off-site assets, secure disposal.
- Physical access is **one control theme** inside a broad ISMS — certification
  audits the whole system, with §7 as the physical slice.

### SOC 2 (AICPA Trust Services Criteria)
- A US **attestation report** (common for SaaS/service orgs), not a certification.
  Built on the **Trust Services Criteria**; the mandatory **Security / Common
  Criteria** are CC1–CC9.
- **CC6 = "Logical and Physical Access Controls"** — notably, **physical and
  logical access are the SAME criterion** (CC6.1–CC6.8). This is the convergence
  point *in audit terms*: the framework already treats them as one control family.
  CC6 covers restricting access, granting/modifying/removing access, and physical
  access to facilities/server rooms (badge readers, etc.).
- **Most common audit exceptions land in CC6**, especially **terminated personnel
  who still have access** and missing periodic **access reviews** — exactly the gap
  convergence/PIAM closes (parent §1.1, §1.3).

### PCI DSS v4.0.1 (payment-card mandate)
- **Mandatory** for entities that store/process/transmit cardholder data
  (contractual, not a certification body). **Requirement 9 = "Restrict physical
  access to cardholder data."**
- **Sub-requirements:**
  - **9.2** — facility **entry controls** into areas with the cardholder-data
    environment.
  - **9.3** — **personnel + visitor** management: authorize/identify staff,
    **escort visitors, badge them, log them**, revoke on termination.
  - **9.4** — **media** handling: classify, secure, log off-site movement, and
    **destroy** media with cardholder data.
  - **9.5** — protect **POI / POS devices** from **tampering and substitution**
    (skimmer defense) — a physical-security control unique to PCI.
- Narrow scope (one data type) but very prescriptive on the visitor/media/device
  mechanics that a converged VMS + PACS must produce evidence for.

### How the frameworks relate (one-paragraph mental model)
**ISO 27001 and SOC 2** are the broad "is your security program sound" regimes —
physical access is one control family in each (ISO Annex A §7; SOC 2 CC6).
**PCI DSS** is a narrow, prescriptive mandate for cardholder data (Req 9). **ISO
60839 and ONVIF** are the product/interop standards the installed system is built
and tested to. **FICAM/NIST (FIPS 201, SP 800-116)** is the US-federal regime
(parent §2). A real converged program answers to **several at once** and wants one
PIAM + PACS that emits evidence satisfying all of them — that single-evidence-
source goal is itself a convergence driver.

---

## 3. Visitor management & identity-proofing assurance vocabulary

Parent §6.1 covers the VMS workflow (pre-registration, ID-scan proofing, watchlist
screening, temporary zone/time-bound credentials). The assurance vocabulary worth
naming:

- **Identity proofing vs authentication** — *proofing* establishes who someone is
  at enrollment (the hard problem for one-time visitors); *authentication*
  re-verifies a known identity later at the door.
- **NIST SP 800-63 assurance levels** (the digital-identity guideline, named for
  orientation): **IAL** (Identity Assurance Level — strength of proofing), **AAL**
  (Authenticator Assurance Level — strength of authentication), **FAL**
  (Federation Assurance Level). Government visitor proofing targets higher IAL than
  commercial sign-in. *Full SP 800-63 mechanics are out of scope here; named so the
  IAL/AAL terms are recognizable when they appear in a converged-identity spec.*
- **Convergence point:** a mature VMS writes the visitor into the **same cardholder
  + audit model** as employees, so the temporary credential is governed,
  time-bound, zone-scoped, and auto-revoked — satisfying PCI 9.3 and ISO A.7.2 in
  one system.

---

## 4. Privacy & regulatory landscape (US-centric; jurisdiction- & date-dependent; NOT legal advice)

Re-read the disclaimer at the top. Converged systems collect **biometrics, video,
and movement logs** — heavily regulated data. **As of June 2026:**

### 4.1 Biometric-privacy law

- **Illinois BIPA (740 ILCS 14)** — the landmark US biometric statute and the
  dominant litigation risk because it has a **private right of action**.
  - **Covers** private entities collecting **biometric identifiers** — fingerprint,
    face geometry, hand geometry, iris/retina scan, voiceprint — which directly
    implicates **biometric door readers** and **facial-recognition access**.
  - **Core duties:** **informed written consent before collection**; a published
    **retention + destruction schedule** (destroy when purpose satisfied or within
    a set period); prohibition on **selling/profiting** from biometrics; limits on
    **disclosure**; reasonable **safeguards**.
  - **Damages:** historically **$1,000 per negligent** violation / **$5,000 per
    intentional or reckless** violation (or actual damages) — which, multiplied
    across employees scanning a reader daily, drove large class actions.
  - **2024 amendment SB 2979** narrowed exposure after *Cothron v. White Castle*:
    repeated collection of the **same** biometric from the **same** person is now a
    **single** violation (not per-scan), and a valid **"written release" includes an
    electronic signature**. Litigation remained heavy through 2024–25.
- **Other state biometric statutes:** **Texas CUBI** and **Washington**'s
  biometric law impose similar duties but are **AG-enforced** (no private right of
  action), so BIPA dominates risk. A growing patchwork of state comprehensive
  privacy laws also reaches biometrics.
- **Bottom line:** if a converged deployment uses biometrics, **biometric-privacy
  law is a first-order design constraint** (consent flow, retention policy, vendor
  contracts) — and it is **changing fast**; check the specific state + latest
  amendments/rulings with counsel.

### 4.2 Video surveillance & notice

- **No single US federal workplace-CCTV privacy regime.** The governing pieces:
  - **Federal wiretap law (ECPA / Wiretap Act)** strongly restricts **audio**
    recording (**one-party consent** federally; **two-party/all-party** in some
    states — e.g. **California**, **Massachusetts**) — far more than silent video.
    **Never assume a camera that also records audio is treated like a silent one.**
  - **State electronic-monitoring notice laws** — several states require **written
    notice** of electronic monitoring (incl. video), e.g. **New York** (2021 law,
    effective May 2022 — notice at hire + conspicuous posting + acknowledgment),
    **Connecticut**, **Delaware**.
  - **Reasonable-expectation-of-privacy bar** — **all** states prohibit cameras in
    restrooms, locker/changing rooms, and similar private spaces.
- **Practical baseline (not advice):** post **conspicuous signage**, give **written
  notice** where required, avoid private areas, separate audio from video
  decisions, and document retention.

### 4.3 General privacy regimes that reach physical-security data

- **State comprehensive privacy laws** (e.g. California **CCPA/CPRA**, and the
  growing roster of other US state laws) treat biometrics/video/access logs as
  personal (often **sensitive**) data with notice, purpose-limitation, retention,
  and data-subject-rights obligations. Applicability thresholds vary.
- **EU GDPR** (for multinationals) treats biometrics as **special-category** data
  (Art. 9) and video/access logs as personal data, with lawful-basis, DPIA,
  retention, and rights obligations.
- These are **flagged, not covered in depth** — scope and applicability are
  fact-specific; engage privacy counsel.

> **Every item in §4 is general, jurisdiction-specific, and time-sensitive
> information — NOT legal advice.** Statutes (BIPA damages, SB 2979, NY/CT/DE
> notice rules, two-party-consent lists) and the case law interpreting them change
> by session and ruling. Verify the current text per jurisdiction before relying
> on anything here.

---

## 5. Sources

See the parent `SKILL.md` **Sources** section — entries **8–9** (SIA, ONVIF), **10**
(ISO 27001 A.7), **11** (SOC 2 CC6), **12** (PCI DSS Req 9), **13** (ISO 60839),
**17** (BIPA + SB 2979), **18** (video-surveillance/notice law), and **19**
(visitor management) ground this file. NIST SP 800-63 (IAL/AAL/FAL) is named for
orientation only. All compliance control IDs are secondary summaries; the licensed
standard text governs.
