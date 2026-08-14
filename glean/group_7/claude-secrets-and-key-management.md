# secrets-and-key-management

**Category:** Science, Biology & Medicine
**Platform:** Claude
**Original Path:** claude/standalone/secrets-and-key-management

## Description
Foundational spine for the secure custody of secrets — the generate → store → transfer → use → rotate → destroy lifecycle and envelope encryption (a KEK wraps DEKs), stated independently of any one implementation. The organizing question: where do the keys live relative to the data they protect? TRIGGER: envelope encryption / key hierarchy (KEK vs DEK, why a key stored next to its ciphertext is theater); KDF selection & parameters (Argon2id vs scrypt vs PBKDF2 vs bcrypt; password-hashing vs key-derivation vs encryption); AES-GCM nonce/IV uniqueness (NIST SP 800-38D) and key-commitment; KMS vs HSM vs software keystore, key ceremonies, Shamir/quorum unseal; secure key transfer & wrapping (AES-KW RFC 3394/5649, HPKE RFC 9180, TLS-in-transit); crypto-periods, rotation & crypto-shredding; PCI split-knowledge & dual control (PCI DSS 4.0.1 Req 3.6/3.7); FIPS 140-3; at-rest vs in-transit vs in-use. SKIP: MongoDB CSFLE/QE field-encryption setup -> mongodb-encryption; MongoDB ops-level encryption & KMS config -> mongodb-operations-expert; browser/WebCrypto vault code review (SubtleCrypto, wrapKey, PBKDF2 iters, WebAuthn PRF) -> webcrypto-vault-reviewer; Atlas tenant/PII isolation for vector search -> atlas-vector-search-pii-isolation; embedding/vector-store leakage -> embedding-inversion-threat-model; bank AI/model-risk governance -> bank-genai-model-risk-governance; FSI regulation map -> fsi-banking-regulatory-context; app/web security review -> security-review.

---

# Secrets & Key Management (Foundational Spine)

The one pattern every credible secret-storage design resolves to, stated
independently of MongoDB, browsers, or any one cloud. **Core thesis: every
secret — a private key, a password, a connection string, a data-encryption key —
moves through the same lifecycle (generate → store at rest → transfer → use →
rotate → destroy), and "security" means controlling the trust boundary at each
stage. The organizing question is: *where do the keys live, relative to the data
they protect?* Encryption with a key stored next to its ciphertext is theater.**

This is the **trunk**. Implementation leaves — `mongodb-encryption` (CSFLE/QE),
`webcrypto-vault-reviewer` (browser vaults), `atlas-vector-search-pii-isolation`
— all *assume* this pattern; this skill *teaches* it. Route implementation
specifics to those leaves (see SKIP).

## 1. The spine: the key hierarchy (envelope encryption)

Never encrypt bulk data directly with a long-lived master key. Wrap short-lived
data keys with it instead.

| Tier | Name | Where it lives | Rule |
|---|---|---|---|
| 1 | **Master key / CMK / KEK** (Key-Encryption Key) | HSM or KMS — **never leaves in cleartext** | Only ever *wraps/unwraps* other keys; never touches plaintext data |
| 2 | **DEK** (Data-Encryption Key) | Stored **wrapped** next to the data; unwrapped only in memory at use | Encrypts the actual data; short-lived; one per object/tenant/field where practical |
| 3 | **Plaintext data / secret** | Never at rest unencrypted | Exists in cleartext only transiently in memory |

Why it wins: rotating the KEK only re-wraps DEKs (cheap) instead of re-encrypting
all data; a single HSM/KMS boundary protects everything; the blast radius of a
leaked DEK is one object, not the estate. NIST SP 800-57 Part 1 (Rev. 5 is the
current published edition; a Rev. 6 is in progress — ⚠️ verify before citing a
revision number) is the reference for key types and crypto-periods.

## 2. The secret lifecycle (the trust boundary at each stage)

- **Generate** — use a CSPRNG (`/dev/urandom`, `crypto.getRandomValues`,
  `secrets`/`os.urandom`). **Never** `Math.random`, PID/timestamp seeds, or
  user-chosen bytes for key material.
- **Store at rest** — wrap the DEK with a KEK held elsewhere (envelope). Keys and
  ciphertext must not share a trust boundary.
- **Transfer** — TLS for in-transit; **wrap** (not re-encrypt) for at-rest handoff
  (§5). Never place a secret in argv, environment dumps, URLs, logs, tickets, or
  source control.
- **Use** — unwrap into memory for the shortest possible window; zeroize where the
  language allows.
- **Rotate** — on a crypto-period schedule and on suspected compromise (§6).
- **Destroy** — crypto-shredding: destroy the key and the ciphertext is
  irrecoverable (§6).

## 3. KDFs: derivation vs hashing vs encryption (do not conflate)

Three distinct jobs; picking the wrong primitive is a top finding.

- **Password hashing** (verify a human password): use a **memory-hard, slow**
  function. Preference order **Argon2id > scrypt > bcrypt > PBKDF2**.
- **Key derivation** (turn a password/shared secret into key bytes): same slow KDF
  when the input is a password; **HKDF** only when the input is already
  high-entropy (HKDF is fast and is *not* a password stretcher).
- **Encryption** (confidentiality): AEAD (§4). A hash or KDF is not encryption.

⚠️ **Parameters drift — treat as version-gated; confirm against the current OWASP
Password Storage Cheat Sheet before quoting.** Representative 2024–2025 minimums:

| Primitive | Representative minimum | Notes |
|---|---|---|
| **Argon2id** | m=19 MiB, t=2, p=1 (or m=12 MiB,t=3 / m=7 MiB,t=5) | Preferred; memory-hard resists GPU/ASIC |
| **scrypt** | N=2^17, r=8, p=1 | Memory-hard; good Argon2 alternative |
| **bcrypt** | work factor ≥ 10 | ⚠️ silently truncates input at 72 bytes; pre-hash long inputs |
| **PBKDF2** | PBKDF2-HMAC-SHA256 ≥ 600,000 iters | **FIPS-approved** — use when FIPS 140 compliance is required; weakest here |

Always use a unique per-secret **salt** (stored alongside); add a **pepper** (a
secret stored separately, ideally in a KMS/HSM) for defense in depth.

## 4. Symmetric encryption at rest: AES-GCM done safely

AES-GCM is the default AEAD, but it is **nonce-fragile**:

- **IV/nonce uniqueness is mandatory (NIST SP 800-38D).** For a given key, an IV
  must **never** repeat. Nonce reuse under GCM is catastrophic — it leaks the GHASH
  authentication key and can reveal plaintext XORs. Use a 96-bit random IV and cap
  invocations per key (≈2^32) or use a deterministic counter you can prove is
  unique; **rotate the DEK** before the limit.
- **GCM is not key-committing.** A single ciphertext can decrypt under two keys
  without an auth error — dangerous in multi-recipient/multi-key designs. Add a
  key-commitment step or use a committing scheme where that matters.
- **Nonce-misuse resistance:** AES-GCM-SIV degrades gracefully on accidental nonce
  reuse (still not committing). Prefer it when unique nonces are hard to guarantee.
- **Always AEAD** (authenticated). Never CBC/CTR without a separate MAC (padding-
  oracle / bit-flip exposure). Bind context with the AAD field.

## 5. Where keys live, and how they move

Ordered by assurance and cost:

| Option | Boundary | Use when |
|---|---|---|
| **Software keystore** (file, OS keyring) | Same host as app | Low assurance; dev/small scale |
| **Cloud KMS** (AWS KMS, Azure Key Vault, GCP KMS, KMIP) | Provider HSM-backed; key never exported | Default for cloud; wraps DEKs, audit-logged |
| **Dedicated HSM / CloudHSM** (FIPS 140-2/3 L3) | Tamper-resistant hardware you control | Regulated/root-of-trust; key ceremonies |

- **Key ceremony:** witnessed, scripted generation/loading of a root key into an
  HSM with split custody — required for PCI/FSI root keys.
- **Shamir Secret Sharing / quorum unseal:** the root/unseal key is split into *n*
  shares needing *k* to reconstruct (e.g. HashiCorp Vault unseal) — no single
  person can unseal.

**Secure transfer / wrapping (never move a key in cleartext):**
- **In transit:** TLS 1.2+/1.3 with modern ciphers; mTLS for service-to-service.
- **At-rest handoff / backup:** **wrap** the key under a recipient KEK using
  **AES Key Wrap (RFC 3394, or RFC 5649 with padding)** — the deterministic,
  standards-track way to protect a key with another key.
- **Public-key sealed transport:** **HPKE (RFC 9180)** for hybrid encrypt-to-a-
  public-key delivery.

## 6. Rotation, crypto-periods, and destruction

- **Crypto-period:** every key has a bounded lifetime (NIST SP 800-57 Part 1).
  Rotate on schedule *and* on suspected compromise. Envelope design makes KEK
  rotation cheap: re-wrap DEKs, no bulk re-encryption.
- **Crypto-shredding:** to make data unrecoverable, destroy its DEK. This is the
  practical basis for "right to erasure" and fast tenant offboarding when each
  tenant/object has its own DEK.

## 7. Bank / PCI / FIPS controls (the high-assurance edge)

- **Split knowledge & dual control (PCI DSS v4.0.1 Req 3.6/3.7):** manual
  cleartext key operations require ≥2 people, each holding only part of the key,
  so no individual can reconstruct it. ⚠️ Cite the current PCI DSS version and
  requirement numbering; they shift between releases.
- **FIPS 140-3** is the current US/Canada crypto-module validation standard
  (superseding 140-2, which CMVP has moved toward historical status — ⚠️ verify
  the transition date/status before stating it). "FIPS-approved" constrains
  primitive choice (e.g. PBKDF2 over Argon2id).
- Complements, not covered here: FSI regulatory map → `fsi-banking-regulatory-context`.

## 8. MongoDB instantiation (routing, not mechanics)

MongoDB is this spine at a specific assurance level: **CSFLE and Queryable
Encryption use exactly the KEK/DEK envelope** — a KMS provider (AWS/Azure/GCP
KMS, KMIP, HashiCorp Vault) holds the **CMK/KEK**, which wraps per-field **DEKs**
in the key vault. ⚠️ **Version-gated:** QE **Equality** and **Range** are GA;
**Prefix/Suffix/Substring** were in **preview** — confirm current GA status per
release notes before telling a customer. Setup, driver helpers, key rotation, and
provider config live in `mongodb-encryption` / `mongodb-operations-expert`.

## 9. Anti-pattern checklist (top findings)

- Key stored beside its ciphertext (no separate boundary) — "encryption theater".
- `Math.random`/weak seed for key or IV material.
- AES-GCM nonce reuse; or unauthenticated CBC/CTR without a MAC.
- Password put through a fast hash (SHA-256) or HKDF instead of a slow KDF.
- Under-parameterized KDF (low PBKDF2 iters), or no per-secret salt.
- Secrets in argv, env dumps, URLs, logs, tickets, or committed to source control.
- Long-lived master key used to encrypt bulk data directly (no DEK layer).
- No rotation policy; no crypto-shredding path for erasure/offboarding.
- Non-extractable/exportable expectations violated (exported HSM/KMS root key).

## Verify-first (version-drift, not soft science)

These are specs, so the risk is **staleness**, not replication. Re-confirm before
customer use: OWASP KDF parameters; PCI DSS version & requirement numbers; FIPS
140-2→140-3 transition status; NIST SP 800-57 revision; MongoDB QE GA-vs-preview
matrix. State the primitive and the pattern with confidence; date the parameters.