<!-- hub-reference-banner -->
> **Reference file — part of the `security-review` hub.** Spoke of the "Okta Security Management"
> concept family (sibling of `okta-expert`). Sibling topics in this family are reference files under
> the hubs (`security-review`) — **not** standalone skills. Load sibling topics from the owning hub's
> `references/<name>.md` (see the hub's routing table).

---

---
name: okta-phishing-resistant-auth
version: 1.0.0
updated: 2026-06-11
description: >
  Phishing-resistant authentication rollout with Okta — FastPass architecture
  and device binding, passkeys/WebAuthn (synced vs device-bound, AAGUID
  policy), NIST AAL mapping, authentication policies and authenticator method
  chains, rollout sequencing for mixed populations, recovery without
  phishable factors, FIDO2 enterprise attestation, rollout failure modes.
  TRIGGER: phishing-resistant MFA / passwordless rollout on Okta; FastPass vs
  passkeys vs security keys; app sign-in policies with phishing-resistant
  possession constraints; synced vs device-bound passkey policy or AAGUID
  allowlists; enrollment/recovery without phishable factors; FastPass
  failures (WebView, loopback, SSO extension).
  SKIP: general Okta platform/API/Terraform (okta-expert); admin-tenant
  hardening and device posture depth (sibling spokes — see the hub routing
  table); non-Okta IdP passkey rollouts (Entra ID, Google).
category: developer
tags: [okta, fastpass, passkeys, webauthn, fido2, phishing-resistant, mfa, authentication-policy, aaguid, rollout]
related_skills: [security-review, software-engineering-patterns]
whenToUse:
  - "Plan a phishing-resistant MFA rollout on Okta for a mixed workforce"
  - "Decide between Okta FastPass, passkeys, and FIDO2 security keys"
  - "Configure an app sign-in policy with the phishing-resistant constraint"
  - "Build an authenticator method chain for high-assurance apps"
  - "Set synced vs device-bound passkey policy and AAGUID allowlists"
  - "Design enrollment and account recovery without phishable factors"
  - "Map Okta authenticators to NIST 800-63B AAL2/AAL3"
  - "Debug FastPass Access denied / WebView / loopback failures"
whenNotToUse:
  - "General Okta API, OAuth, hooks, or Terraform work: use okta-expert"
  - "Device assurance / posture policy design: use the hub's device-posture sibling spoke"
origin: local
---

# Okta Phishing-Resistant Authentication Rollout

Expert reference for rolling out phishing-resistant authentication (Okta
FastPass, passkeys/WebAuthn, FIDO2 security keys) across a real workforce:
the architecture that makes each factor phishing-resistant, the policy
machinery that enforces it, the sequencing that makes the rollout survive
contact with mixed user populations, and the enrollment/recovery flows
attackers actually target once sign-in is hardened.

## Why "phishing-resistant" is a protocol property, not a product label

NIST SP 800-63B defines phishing resistance as the authentication protocol's
ability to prevent disclosure of authentication secrets and valid
authenticator outputs to an impostor verifier *without relying on the
claimant's vigilance*. Two mechanisms qualify: **channel binding**
(client-authenticated TLS, PIV/CAC) and **verifier name binding** (WebAuthn —
the credential is scoped to the authenticated RP domain). Anything requiring
manual entry of an authenticator output (OTP, out-of-band codes) is
explicitly *not* phishing-resistant, because the output isn't bound to the
session being authenticated.[^1] CISA calls FIDO/WebAuthn "the only widely
available phishing-resistant authentication" and ranks SMS/voice as
last-resort.[^2]

Assurance-level mapping (NIST 800-63B rev 4, verified-as-of: 2026-06-11):

| Level | Requirement | Okta implication |
| --- | --- | --- |
| AAL2 | Two factors; verifier SHALL offer at least one phishing-resistant option | Synced passkeys acceptable; FastPass, passkeys, or security keys all qualify[^1] |
| AAL3 | Phishing-resistant authenticator with **non-exportable** key, hardware-protected; **syncable authenticators SHALL NOT be used** | Device-bound only: FastPass on TPM/Secure Enclave hardware, or hardware FIDO2 keys; Okta states FastPass can support AAL3 with FIPS 140-3 Level 2 devices[^1][^3] |

NIST's syncable-authenticator appendix additionally requires the sync fabric
itself be protected by AAL2-equivalent MFA and that exported keys be
encrypted — the security boundary moves from the device to the passkey
provider account.[^1]

## Okta FastPass architecture and device binding

FastPass is a device-bound cryptographic authenticator delivered through the
Okta Verify app (Windows, macOS, iOS, Android). At enrollment, Okta Verify
generates key pairs in the device's hardware keystore (TPM / Secure Enclave
where available): a **proof-of-possession** key, plus a separate **user
verification** key if biometrics are enabled. Private keys never leave the
device; the public keys plus device-context metadata register with the Okta
org. The credential is cryptographically bound to one specific Okta tenant
and one device.[^4]

How the browser hands the challenge to Okta Verify determines whether the
flow is actually phishing-resistant:

| Probing scheme | Type | Phishing-resistant? |
| --- | --- | --- |
| Loopback server (Okta Verify binds localhost port 8769; the browser probes a fixed port list) | Silent | **Yes** — Okta Verify captures the browser's `Origin` header, signs it into the response; Okta rejects domain mismatches and logs `Okta FastPass declined phishing attempt`[^4][^5] |
| Credential SSO extension (managed Apple devices, Safari) | Silent | **Yes** — Apple Extensible SSO framework signs headers for the configured domain[^4] |
| Custom URI / AppLink / Universal Link | Interactive fallback | **No** — no origin tracking; exists to launch Okta Verify, not to prove origin[^4][^21] |

The single most load-bearing configuration fact: **FastPass only guarantees
phishing resistance when the app sign-in policy's possession constraint
explicitly requires "Phishing resistant."** Without that constraint, a failed
loopback (port conflict, DNS-rebind-protecting router, AiTM-modified
loopback address) silently falls back to the URI scheme and AiTM phishing
works again. Okta's own security team confirms that modern "FastPass bypass"
claims almost always reduce to orgs that never enforced the constraint in
policy.[^6][^21]

## Authentication policies and authenticator method chains

Okta Identity Engine enforcement lives in **app sign-in policies** (formerly
authentication policies). The possession-factor constraints that matter:

- **Phishing resistant** — requires possession factors that cryptographically
  verify the login origin (FastPass, Passkeys/FIDO2 WebAuthn). Selecting it
  auto-selects Device bound.[^7]
- **Hardware protected** — keys must live in TPM/Secure Enclave; implies
  Device bound.[^7]
- **Device bound (excludes phone and email)** — keys not transferable
  without re-enrollment; SMS and email are the only non-device-bound
  possession factors.[^7]
- **User presence/interaction** — require user verification (biometric/PIN)
  so Okta Verify can't satisfy possession silently via certificate-based
  auth; set User verification to Required on each authenticator so new
  enrollments satisfy it.[^7]

In Terraform these map to `okta_app_signon_policy_rule` constraint JSON:
`possession.phishingResistant`, `possession.hardwareProtection`,
`possession.deviceBound`, each `OPTIONAL` (default) or `REQUIRED` — codify
the constraint in IaC so a console edit can't silently drop it.[^8]

**Authenticator method chains** (verified-as-of: 2026-06-11) let a rule
specify ordered sequences instead of "any N factors": e.g., password →
passkey for one chain and Okta Verify → TOTP for another, per-step factor
constraints (phishing-resistant FastPass, hardware-protected smart card),
and even two phishing-resistant steps for the most sensitive apps (WebAuthn
then FastPass).[^9] Use chains where "any 2 factor types" would let a
phishable pair (password + push) satisfy a rule you intended to be
phishing-resistant — the chain is the policy-level antidote to downgrade.

## Passkeys: synced vs device-bound, AAGUID policy, attestation

Two passkey families with different security properties:

- **Device-bound passkeys** — private key created and stored on one
  authenticator (security key, or platform authenticator with syncing
  unavailable/blocked). Supports attestation; can satisfy hardware-protected
  constraints and AAL3.[^1][^10]
- **Synced passkeys** — private key encrypted and synced through a provider
  cloud (iCloud Keychain, Google Password Manager, 1Password, Bitwarden).
  Phishing-resistant (still WebAuthn origin-bound) but exportable across
  devices, **no attestation**, and excluded from AAL3.[^1][^10]

The **AAGUID** (Authenticator Attestation GUID) is the 128-bit make/model
identifier each FIDO2 authenticator presents at registration. Critical
caveat: **without enforced attestation, the AAGUID is self-reported and
cannot be trusted as a security control** — synced-passkey providers can't
sign with a hardware batch certificate, so AAGUID allow/blocklists for
synced passkeys are organizational policy guides, not controls.[^10]

Okta's controls (verified-as-of: 2026-06-11):

- **FIDO MDS + custom AAGUID list.** Okta resolves authenticator
  make/model and `keyProtection` (hardware) status against the FIDO Alliance
  Metadata Service; authenticators absent from MDS can be added to a custom
  AAGUID list (custom entries override MDS entries). Manage via Admin
  Console (Security > Authenticators > Passkeys) or the Authenticators API
  (`/api/v1/authenticators/{id}/aaguids`).[^7][^11][^23]
- **Authenticator groups.** Group recognized passkey authenticators and
  reference the groups in enrollment policies, so different populations are
  restricted to different hardware (note: FIDO U2F-era keys can't enroll
  under authenticator groups).[^12]
- **Allowlist failure mode.** When an AAGUID allowlist is active, a new key
  model (including newer YubiKey firmware with a *different AAGUID*) fails
  login with a System Log error naming the rejected AAGUID; fix by adding
  the AAGUID, no re-enrollment needed. Budget for this at every hardware
  refresh.[^13]
- **Block synced passkeys.** An org-wide toggle blocks *new* enrollments of
  authenticators known to create syncable credentials. It is explicitly
  **not an attestation-based check** — it blocklists known-syncable
  authenticators — and it has UX casualties: Safari Touch ID on macOS
  Monterey and iOS 16 platform enrollments are blocked too (use FastPass or
  NFC/USB-C keys for those users). Existing synced enrollments are not
  revoked — audit and clean them up separately.[^11][^14]

**FIDO2 enterprise attestation (EA):** by default FIDO prohibits uniquely
identifying information in attestation; EA removes that limit by binding the
key pair to a serial number (or equivalent) so an authorized RP can track a
specific issued authenticator. EA only works for RP IDs pre-provisioned on
the authenticator by the vendor (vendor-facilitated) or authorized by an
enterprise-managed platform/browser (platform-managed).[^15][^16] Okta's
WebAuthn enrollment API exposes the standard attestation conveyance enum
including `enterprise` (verified-as-of: 2026-06-11); treat EA as a
procurement decision — it must be ordered into the keys (e.g., custom
YubiKey programs) and is justified mainly where per-device audit trails are
a compliance requirement.[^16][^17]

## Okta Verify (FastPass) vs platform authenticators — choosing per population

| Dimension | FastPass | Platform passkeys (synced) | Hardware FIDO2 keys |
| --- | --- | --- | --- |
| Phishing-resistant | Yes (when enforced in policy) | Yes | Yes |
| Device-bound | Yes | No (synced across provider account) | Yes |
| Device posture signals at auth time | Yes — OS version, encryption, jailbreak, EDR signals, silent re-checks mid-session | No | No |
| Attestation / AAGUID trust | N/A (Okta-proprietary binding) | No attestation | Yes (incl. EA option) |
| Works without installing software | No (Okta Verify app) | Yes (built into OS/browser) | Yes (any WebAuthn browser) |
| Cross-ecosystem consistency | Same UX on Windows/macOS/iOS/Android | Fragmented by vendor ecosystem | Consistent but physical-token logistics |
| AAL ceiling | AAL3-capable on FIPS-validated hardware[^3] | AAL2[^1] | AAL3 |

Default pattern: FastPass as the primary workforce authenticator (it doubles
as the device-context source), passkeys/security keys as the
phishing-resistant *backup and browser-independent* factor, hardware keys
with attestation for admins and for populations that can't run Okta
Verify.[^3][^10]

## Rollout sequencing for mixed populations

Okta's deployment guidance and CISA's prioritization converge on the same
shape (phases, not big bang)[^2][^3][^18][^19]:

1. **Pre-deployment.** OIE required. Baseline help-desk password-reset
   volume (the ROI metric). Decide the factor mix per population. Configure
   a global session policy; enable FastPass and Passkeys authenticators;
   set User verification to Required.
2. **Pilot.** Small early-adopter group; deploy Okta Verify via MDM where
   available; validate real environments (managed/BYOD/VDI), not just IT's
   test machines.
3. **Admins and high-value targets first.** CISA: start with system
   administrators, help desk, and other high-value accounts (attorneys, HR);
   Okta: prioritize high-risk users (IT, execs) and sensitive apps. Give
   this tier device-bound-only policy (hardware keys or FastPass with
   hardware protection, AAGUID allowlist, no synced passkeys).[^2][^3]
4. **Policy ratchet, app by app.** Apply phishing-resistant constraints to
   low-traffic, non-sensitive apps first to validate policy mechanics, then
   expand. Edit authentication policies in small increments and test —
   policy mistakes lock out populations at scale.[^3]
5. **General workforce.** Opt-in/silent enrollment prompts first, then
   mandatory. Synced passkeys are an acceptable AAL2 tier for the broad
   workforce when device-bound coverage is impractical.[^1][^10]
6. **Contractors and unmanaged/BYOD.** FastPass works on unmanaged devices
   (registered via Okta Verify without MDM) and still enforces device
   assurance; where Okta Verify can't be installed, issue security keys or
   accept synced passkeys with a tighter app scope.[^3]
7. **Mobile-only and frontline users.** Verify the actual sign-in
   environments: shared kiosks and retail floors break assumptions —
   cross-device passkey flows need Bluetooth proximity, which is often
   disabled on kiosks; shared devices break 1:1 user-device binding
   entirely. Plan NFC security keys or dedicated flows for this tier.[^19]
8. **Deprecate phishable factors.** The rollout only delivers its security
   value at this step: remove SMS, voice, email OTP, and bare push from
   policies *and from enrolled factors per user* — an enrolled-but-unused
   phishable factor is still a downgrade path. Track FastPass adoption via
   the MFA Usage report (`Okta Verify - signed_nonce`).[^3][^20]

## Fallback and recovery without phishable factors

Once sign-in is phishing-resistant, attackers move to enrollment, recovery,
and help-desk flows — Okta's product security team describes this as the
"chicken and egg" problem: the factors that verify identity *before*
enrolling a phishing-resistant factor are often not themselves
phishing-resistant.[^6] Demonstrated attacks: AiTM-proxying the Okta Verify
enrollment redirect to capture the authorization code and bind an
*attacker's* device to the victim account; and downgrade attacks where the
AiTM page simply hides the passkey option and presents the phishable backup
factor.[^20][^21]

Controls, in priority order:

1. **Require a phishing-resistant authenticator to enroll additional
   authenticators** (Early Access setting, verified-as-of: 2026-06-11) —
   closes the enrollment hole for everyone already holding a
   phishing-resistant factor.[^6][^21]
2. **Restrict enrollment by network zone** — deny authenticator enrollment
   from IPs outside trusted zones; pairs well with on-site onboarding.[^21]
3. **Exclude low-assurance factors from enrollment-verification flows**
   (enrollment policy options) so SMS/email can't bootstrap a new
   authenticator.[^6]
4. **Pre-enrolled hardware keys for day-one and recovery.** Yubico FIDO
   Pre-reg with Okta (GA since November 2024; orchestrated through Okta
   Workflows templates) ships factory-pre-registered YubiKeys directly to users with
   the PIN delivered separately — no password or phishable factor ever
   exists in the flow, and the key doubles as the recovery factor when a
   FastPass phone is lost, without a help-desk call.[^22][^24]
5. **Proximity-based FastPass re-enrollment** — new-device enrollment
   authorized by physical proximity to an already-registered device.[^6]
6. **Harden the help desk itself.** Recovery via help desk is the
   social-engineering target (the pattern behind the 2023 casino breaches);
   require phishing-resistant verification or manager/ID-proofing for
   resets, never SMS-to-caller. A passkey rollout is not complete until
   recovery is also phishing-resistant — otherwise recovery becomes the new
   front door.[^19][^20]

## Common rollout failure modes

| Failure mode | Mechanism | Mitigation |
| --- | --- | --- |
| Phishing-resistant constraint never enforced | FastPass silently falls back to Custom URI scheme (no origin check) on loopback failure; AiTM works | Require "Phishing resistant" in every relevant app sign-in policy rule; treat it as the control, not the authenticator[^6][^21] |
| WebView-based apps | Embedded WebViews can't do the loopback/extension dance; auth fails with `Access denied` under phishing-resistant policy | Inventory WebView apps in the pilot; exempt or fix before broad enforcement[^5] |
| macOS Safari without SSO extension | Safari FastPass not phishing-resistant unless the managed-device SSO extension profile is pushed | MDM-deploy the Credential SSO extension profile[^5] |
| DNS Rebind Protection on routers | Blocks browser→localhost loopback; phishing-resistance checks fail (home/SOHO networks) | Document for remote workers; expect and triage these tickets[^5] |
| Loopback port conflict (8769) | Okta Verify fails to bind, no error, no alternative port; silent downgrade | Enforced phishing-resistant policy turns silent downgrade into visible failure — by design[^21] |
| Malicious browser extensions | `declarativeNetRequest` + host permissions can spoof `Origin` to the loopback server, breaking origin trust | Extension allowlisting on managed browsers; require user verification so relays still prompt the user[^21] |
| Downgrade via leftover factors | AiTM kits hide the passkey button and harvest the OTP/push backup | Remove phishable enrolled factors per user; method chains; conditional access to managed devices[^9][^20] |
| AAGUID allowlist drift | New key firmware = new AAGUID = login failures | Watch System Log for allowlist rejections; update list at procurement time[^13] |
| Block-synced-passkeys side effects | iOS 16 / macOS Monterey Safari platform enrollments blocked along with synced passkeys | Offer FastPass or NFC/USB-C keys to affected users before flipping the toggle[^11] |
| Single-passkey lockouts | One enrollment in one browser = lockout on browser block or device loss | Encourage multiple enrollments across browsers/devices; keep a hardware-key recovery tier[^3][^19] |
| Backend metrics overstate success | Client-side failures (cancelled biometric prompts, missing credential manager, wrong credential manager) never reach the IdP | Track adoption per population, support tickets by OS/browser, registered-but-unused passkeys — not just auth success rate[^19] |
| "Enabled ≠ enforced" adoption stall | Passkeys enabled but <5% usage; users default back to passwords | Staged rollout with nudges, make the phishing-resistant factor the easiest path, then retire phishable methods per cohort[^19] |

## Review checklist (security-review hub)

When reviewing an Okta org's phishing-resistant posture, verify with
evidence:

1. App sign-in policies for sensitive apps require the **Phishing
   resistant** possession constraint (not merely "2 factor types").
2. User verification set to **Required** on FastPass/Passkeys.
3. Enrollment policy requires phishing-resistant verification (or network
   zone restriction) before new-authenticator enrollment.
4. Phishable factors (SMS, voice, email, bare push) removed from policies
   AND from privileged users' enrolled factors.
5. Admin tier: device-bound-only (attestation/AAGUID-controlled hardware
   keys or hardware-protected FastPass); synced passkeys blocked or scoped.
6. Recovery/help-desk flow cannot be completed with phishable factors alone.
7. System Log monitored for `Okta FastPass declined phishing attempt` and
   AAGUID allowlist rejections, ideally wired to Workflows alerts.[^4]

## References

All sources fetched or surfaced via search during the 2026-06-11 research
run. Access date for all: 2026-06-11. Inline anchors `[^n]` cite the
correspondingly numbered entry below.

1. NIST SP 800-63B rev 4 — Digital Identity Guidelines: Authentication and
   Authenticator Management (AALs, phishing-resistance definition, syncable
   authenticator appendix). Standards body.
   <https://pages.nist.gov/800-63-4/sp800-63b.html>
2. CISA — Implementing Phishing-Resistant MFA fact sheet (Oct 2022;
   strongest-to-weakest MFA table, prioritization guidance). Government.
   <https://www.cisa.gov/sites/default/files/publications/fact-sheet-implementing-phishing-resistant-mfa-508c.pdf>
3. Okta — Beyond MFA: The FastPass Advantage eBook (2025; rollout
   checklists, AAL claims, factor-comparison tables). Vendor docs.
   <https://www.okta.com/sites/default/files/2025-08/Okta-FastPass-Best-Practices-eBook-2025.pdf>
4. Okta blog — A Deep Dive Into Okta FastPass (enrollment key flows, probing
   schemes, loopback origin check). Vendor engineering blog.
   <https://www.okta.com/blog/product-innovation/a-deep-dive-into-okta-fastpass/>
5. Okta docs — Phishing-resistant authentication (Identity Engine;
   restrictions: WebView, macOS SSO extension, UWP, DNS rebind). Vendor docs.
   <https://help.okta.com/oie/en-us/content/topics/identity-engine/authenticators/phishing-resistant-auth.htm>
6. Okta Security — FastPass: The battle-hardened authenticator (hardening
   history, enrollment/recovery attacks, user-verification fallback, trusted
   app filters). Vendor security blog.
   <https://sec.okta.com/articles/fastpasshardening/>
7. Okta docs — Add an app sign-in policy rule (possession constraints,
   FIDO MDS keyProtection, user-presence requirements). Vendor docs.
   <https://help.okta.com/oie/en-us/content/topics/identity-engine/policies/add-app-sign-on-policy-rule.htm>
8. Terraform Registry — okta_app_signon_policy_rule (phishingResistant /
   hardwareProtection / deviceBound constraint schema). Vendor IaC docs.
   <https://registry.terraform.io/providers/okta/okta/6.5.3/docs/resources/app_signon_policy_rule>
9. Okta docs — Authentication method chain (ordered methods, per-step
   constraints, multiple chains). Vendor docs.
   <https://help.okta.com/OIE/en-us/content/topics/identity-engine/policies/authentication-method-chain.htm>
10. Microsoft Learn — Passkeys (FIDO2) in Entra ID: synced vs device-bound,
    attestation, passkey profiles, AAGUID restrictions (cross-vendor
    corroboration of passkey-class properties). Vendor docs.
    <https://learn.microsoft.com/en-us/entra/identity/authentication/concept-authentication-passkeys-fido2>
11. Okta docs — Configure the Passkeys (FIDO2 WebAuthn) authenticator (MDS
    AAGUID review, custom AAGUIDs, syncable-passkey blocking side effects,
    multi-enrollment guidance). Vendor docs.
    <https://help.okta.com/oie/en-us/Content/Topics/identity-engine/authenticators/configure-passkeys.htm>
12. Okta docs — Configure Passkeys authenticator groups. Vendor docs.
    <https://help.okta.com/oie/en-us/content/topics/identity-engine/authenticators/passkeys-authenticator-groups.htm>
13. Okta Support — Unable to Log In to Okta Using YubiKey (AAGUID allowlist
    rejection symptom and fix). Vendor KB.
    <https://support.okta.com/help/s/article/unable-to-log-in-to-okta-using-yubikey>
14. Okta Support — Passkey Management (Block synced Passkeys behavior:
    new-enrollment scope, non-attestation-based, multi-device blocking).
    Vendor KB.
    <https://support.okta.com/help/s/article/Passkey-Management?language=en_US>
15. FIDO Alliance — Attestation White Paper (2024; enterprise attestation,
    RP ID provisioning, serial-number binding). Standards body.
    <https://fidoalliance.org/wp-content/uploads/2024/06/EDWG_Attestation-White-Paper_2024-1.pdf>
16. Yubico developers — Enterprise Attestation (vendor-facilitated vs
    platform-managed EA, RP ID lists). Vendor docs.
    <https://developers.yubico.com/WebAuthn/Concepts/Enterprise_Attestation/>
17. Okta developer — Start a WebAuthn enrollment (MyAccount API; attestation
    conveyance enum incl. `enterprise`). Vendor API docs.
    <https://developer.okta.com/docs/api/openapi/okta-myaccount/myaccount/webauthn/startwebauthnenrollment>
18. Okta — Step-by-step guide to becoming phishing resistant with Okta
    FastPass (phased deployment starting with pilot). Vendor PDF guide.
    <https://www.okta.com/sites/default/files/2023-10/Step-by-step-guide-to-becoming-phishing-resistant-with-Okta-FastPass.pdf>
19. Entra.News / Corbado — 5 Lessons from Rolling Out Passkeys to Millions
    of Users (May 2026; staged rollout, recovery as stage 4, kiosk/Bluetooth
    pitfalls, telemetry blind spots, multi-passkey guidance). Practitioner.
    <https://entra.news/p/5-lessons-from-rolling-out-passkeys>
20. Push Security — MFA downgrade: How attackers are getting around
    phishing-resistant authentication (Jul 2025; AiTM downgrade mechanics,
    backup-factor risk, conditional-access mitigations). Security research.
    <https://pushsecurity.com/blog/mfa-downgrade-attacks>
21. Obsidian Security — Behind the Shield: Cracking the Limits of Okta
    FastPass (2025; enrollment AiTM, loopback port conflict downgrade,
    URI-scheme relay, extension Origin spoofing, mitigations). Security
    research (adversarial/disconfirming source).
    <https://www.obsidiansecurity.com/blog/behind-the-shield-cracking-the-limits-of-okta-fastpass>
22. Yubico docs — FIDO Pre-reg with Okta (pre-enrolled YubiKey shipment
    workflow, PIN separation, recovery use). Vendor docs.
    <https://docs.yubico.com/cloud-services/fidoprereg-okta/about-fpr-okta.html>
23. Okta developer — Authenticators Administration API (Passkey authenticator
    `webauthn` key, custom AAGUID CRUD endpoints, 2026.04.0 rename note).
    Vendor API docs.
    <https://developer.okta.com/docs/api/openapi/okta-management/management/tags/authenticator>
24. Okta blog — New Okta and Yubico integration delivers phishing-resistant
    authentication at enterprise scale (FIDO Pre-reg GA November 2024).
    Vendor blog.
    <https://www.okta.com/blog/customers-and-partners/new-okta-and-yubico-integration-delivers-phishing-resistant-authentication-at/>
