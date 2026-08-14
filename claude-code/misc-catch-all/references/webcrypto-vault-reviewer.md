<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `webcrypto-vault-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: webcrypto-vault-reviewer
version: 1.2.0
updated: 2026-05-29
description: >-
  Web Crypto API security review and implementation patterns — SubtleCrypto, PBKDF2 key derivation
  (iterations, salts), AES-GCM (IV generation, tag length, nonce reuse), HKDF for high-entropy
  secrets, key wrapping & vault envelope patterns (multi-unlock architecture), passkey/WebAuthn
  PRF vault unlock, chrome.storage session vs local key custody, security checklist for vaults in
  Chrome extensions and web apps. TRIGGER: review/write code calling crypto.subtle.*; auditing a
  vault, password manager, or encrypted-storage feature; PBKDF2 iterations, IV generation,
  key-wrapping flows; WebAuthn PRF passkey unlock; anti-patterns (IV reuse, weak iterations,
  extractable keys, Math.random() secrets); AES-GCM, PBKDF2, HKDF, wrapKey/unwrapKey,
  key-commitment. SKIP: TLS/mTLS config → http-security-headers; OAuth/OIDC flows →
  web-auth-patterns; MongoDB encryption at rest → mongodb-encryption; security audits without
  crypto code → security-reviewer.
category: developer
tags:
  - cryptography
  - web-crypto
  - aes-gcm
  - pbkdf2
  - passkeys
  - security
  - key-wrapping
  - hkdf
  - vault-envelope
  - webauthn
triggers:
  - crypto.subtle
  - SubtleCrypto
  - AES-GCM
  - PBKDF2
  - deriveKey
  - wrapKey
  - unwrapKey
  - importKey
  - vault encrypt
  - vault decrypt
  - WebAuthn PRF
  - passkey vault
  - IV reuse
  - nonce reuse
  - key derivation browser
related_skills:
  - software-engineering-patterns
  - security-review
  - chrome-extension-expert
  - mongodb-operations-expert
---

# WebCrypto Vault Reviewer

Security review and implementation patterns for browser-side cryptography built on the Web Crypto API (`crypto.subtle`). Targets vault envelope patterns commonly used in Chrome extensions and web apps where sensitive data must be encrypted at rest without a server round-trip.

## When to use this skill

- Reviewing or writing code that calls `crypto.subtle.*` methods
- Auditing a vault, password manager, or encrypted-storage feature in a web app or extension
- Evaluating PBKDF2 iteration counts, IV generation, or key-wrapping flows
- Implementing passkey-based vault unlock via the WebAuthn PRF extension
- Checking for common cryptographic anti-patterns (IV reuse, weak iterations, extractable keys)

**Do not use** for TLS/mTLS config (`http-security-headers`), OAuth/OIDC flows (`web-auth-patterns`), MongoDB encryption at rest (`mongodb-encryption`), or general security audits without crypto code (`security-reviewer`).

---

## Core Concepts

### SubtleCrypto API surface

| Method | Purpose | Vault role |
|---|---|---|
| `generateKey` | Create a new CryptoKey | Generate AES-GCM master key |
| `importKey` | Import raw bytes as CryptoKey | Import PBKDF2 base material, import raw AES key |
| `deriveKey` | Derive a CryptoKey from base material | PBKDF2 password → AES wrapping key |
| `encrypt` / `decrypt` | AES-GCM encrypt/decrypt | Encrypt vault payload, wrap master key |
| `wrapKey` / `unwrapKey` | Export+encrypt / decrypt+import a key | Envelope wrap (alternative to manual encrypt) |
| `digest` | SHA-256 hash | Hash PRF output before import |

### CryptoKey properties

```js
const key = await crypto.subtle.generateKey(
  { name: 'AES-GCM', length: 256 },
  false,   // extractable — KEEP FALSE for operational keys
  ['encrypt', 'decrypt']
);
// key.extractable → false  (raw bytes cannot be read from JS)
```

Mark keys as non-extractable (`extractable: false`) unless you specifically need `exportKey()` or `wrapKey()`. Non-extractable keys prevent JavaScript from reading raw key material, limiting exposure to XSS.

### Randomness

Always use `crypto.getRandomValues()`. Never use `Math.random()` for cryptographic values.

**Quota limit:** A single call can fill at most 65,536 bytes. For larger buffers, call in a loop.

```js
function randomBytes(length = 32) {
  return crypto.getRandomValues(new Uint8Array(length));
}
```

---

## PBKDF2 Key Derivation

### Recommended parameters (2025–2026)

| Parameter | Minimum | Recommended | Notes |
|---|---|---|---|
| Hash | SHA-256 | SHA-256 | SHA-512 also acceptable |
| Iterations | 310,000 | 600,000 | OWASP 2025: 600K for FIPS-140; 310K minimum general |
| Salt length | 16 bytes | 16–64 bytes | Unique per envelope, from CSPRNG |
| Output | AES-GCM 256-bit key | AES-GCM 256-bit key | Via `deriveKey`, not `deriveBits` when possible |

### Reference implementation

```js
async function derivePasswordKey(password, saltBase64, iterations = 600_000) {
  const baseKey = await crypto.subtle.importKey(
    'raw',
    new TextEncoder().encode(String(password)),
    'PBKDF2',
    false,
    ['deriveKey']
  );
  return crypto.subtle.deriveKey(
    { name: 'PBKDF2', hash: 'SHA-256', salt: fromBase64(saltBase64), iterations },
    baseKey,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt']
  );
}
```

### Review checklist for PBKDF2

- [ ] Iterations >= 310,000 (prefer 600,000 for new implementations)
- [ ] Salt is at least 16 bytes from `crypto.getRandomValues()`
- [ ] Salt is unique per envelope (never shared across users/envelopes)
- [ ] Salt stored alongside the wrapped payload (it is not secret)
- [ ] Iteration count stored in envelope metadata for future upgradability
- [ ] Base key `extractable: false`
- [ ] Derived key `extractable: false`
- [ ] Usages minimally scoped (`['deriveKey']` for base, `['encrypt','decrypt']` for derived)

---

## AES-GCM Encryption

### IV / Nonce requirements

| Property | Requirement |
|---|---|
| Length | 96 bits (12 bytes) |
| Generation | `crypto.getRandomValues(new Uint8Array(12))` |
| Uniqueness | MUST be unique per encryption with the same key |
| Storage | Stored alongside ciphertext (not secret) |

**Critical:** IV reuse in AES-GCM is catastrophic — it leaks `P1 XOR P2` between any two plaintexts and recovers the authentication key, enabling forgery of arbitrary ciphertexts.

**Nonce-misuse-resistant alternative:** AES-GCM-SIV (RFC 8452) tolerates accidental nonce reuse without catastrophic failure, but is not available in the Web Crypto API as of 2026. Userland libraries (e.g., `@noble/ciphers`) are required, trading off hardware-backed performance and non-extractable key guarantees.

### Reference implementation

```js
async function encryptPayload(keyBytes, plaintext) {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const key = await crypto.subtle.importKey('raw', keyBytes, 'AES-GCM', false, ['encrypt']);
  const ciphertext = await crypto.subtle.encrypt(
    { name: 'AES-GCM', iv },   // tagLength defaults to 128 — keep it
    key,
    new TextEncoder().encode(JSON.stringify(plaintext))
  );
  return { iv: toBase64(iv), cipherText: toBase64(ciphertext) };
}

async function decryptPayload(keyBytes, envelope) {
  const key = await crypto.subtle.importKey('raw', keyBytes, 'AES-GCM', false, ['decrypt']);
  const plainBytes = await crypto.subtle.decrypt(
    { name: 'AES-GCM', iv: fromBase64(envelope.iv) },
    key,
    fromBase64(envelope.cipherText)
  );
  return JSON.parse(new TextDecoder().decode(plainBytes));
}
```

### Review checklist for AES-GCM

- [ ] IV is exactly 12 bytes (96 bits)
- [ ] IV generated from `crypto.getRandomValues()`, not a counter or timestamp
- [ ] IV stored with the ciphertext, never discarded
- [ ] Fresh IV generated for every `encrypt()` call — no reuse
- [ ] `tagLength` is 128 (the default) — never reduced below 96
- [ ] Key length is 256 bits
- [ ] Decryption failures (`OperationError`) caught and treated as authentication failures
- [ ] `additionalData` (AAD) used where context-binding is needed (optional but recommended)

---

## Vault Envelope Pattern

Separates the master key from the data it protects, allowing multiple unlock methods without re-encrypting the entire vault.

### Architecture

```
User password  → PBKDF2 → wrapping key → AES-GCM-encrypt(masterKey) → passwordEnvelope
Auth file      → HKDF   → wrapping key → AES-GCM-encrypt(masterKey) → authFileEnvelope
Passkey PRF    → SHA-256 → wrapping key → AES-GCM-encrypt(masterKey) → passkeyEnvelope

masterKey → AES-GCM-encrypt(vaultData) → encryptedPayload

Stored envelope:
{
  "version": 1,
  "unlockMethods": [
    { "method": "password", "salt": "...", "iterations": 600000, "iv": "...", "wrappedKey": "..." },
    { "method": "auth_file", "salt": "...", "iv": "...", "wrappedKey": "..." },
    { "method": "passkey", "credentialId": "...", "prfSalt": "...", "iv": "...", "wrappedKey": "..." }
  ],
  "payload": { "iv": "...", "cipherText": "..." }
}
```

### Key wrapping

```js
async function wrapMasterKey(wrappingKey, masterKeyBytes) {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const cipher = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, wrappingKey, masterKeyBytes);
  return { iv: toBase64(iv), wrappedKey: toBase64(cipher) };
}

async function unwrapMasterKey(wrappingKey, envelope) {
  const plain = await crypto.subtle.decrypt(
    { name: 'AES-GCM', iv: fromBase64(envelope.iv) },
    wrappingKey,
    fromBase64(envelope.wrappedKey)
  );
  return new Uint8Array(plain);
}
```

### HKDF for auth-file key derivation

When the secret source is already high-entropy (e.g., a 32-byte random auth file), use HKDF — not PBKDF2. PBKDF2's iteration cost adds latency with no security benefit for already-random input.

```js
async function deriveAuthFileKey(fileSecretBytes, saltBase64) {
  const material = await crypto.subtle.importKey('raw', fileSecretBytes, 'HKDF', false, ['deriveKey']);
  return crypto.subtle.deriveKey(
    {
      name: 'HKDF',
      hash: 'SHA-256',
      salt: fromBase64(saltBase64),
      info: new TextEncoder().encode('my-app-auth-file-wrap-v1')
    },
    material,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt']
  );
}
```

### Review checklist for vault envelope

- [ ] Master key is 32 bytes from `crypto.getRandomValues()`
- [ ] Master key raw bytes never persisted to storage in cleartext
- [ ] Master key held in memory only while vault is unlocked
- [ ] Each unlock method wraps the same master key independently
- [ ] Each wrapping envelope has its own unique IV and salt
- [ ] Envelope includes a version field
- [ ] HKDF (not PBKDF2) used for high-entropy secrets
- [ ] HKDF `info` field includes a unique domain-separation string
- [ ] On vault lock, master key bytes zeroed from memory where possible

### Master key rotation

1. Unlock vault with current master key.
2. Generate new 32-byte master key via `crypto.getRandomValues()`.
3. Re-encrypt vault payload with new master key.
4. Re-wrap new master key under every active unlock method.
5. Write the updated envelope atomically (`chrome.storage.local.set` single call).
6. Zero the old master key bytes from memory.

---

## Passkey / WebAuthn PRF Extension

The WebAuthn PRF extension allows a passkey authenticator to deterministically derive a symmetric secret from a relying-party-provided salt, enabling passwordless vault unlock.

### How it works

1. **Registration:** Create passkey with `navigator.credentials.create()`, store `credentialId`.
2. **Authentication:** Call `navigator.credentials.get()` with `prf` extension and a salt.
3. **Key derivation:** Authenticator computes `HMAC-SHA-256(internalKey, salt)`, returns 32-byte secret.
4. **Vault unlock:** Hash or import PRF output, use it to unwrap the master key.

### PRF feature detection at registration

```js
if (typeof PublicKeyCredential?.getClientCapabilities === 'function') {
  const caps = await PublicKeyCredential.getClientCapabilities();
  if (!caps?.['extension:prf']) {
    throw new Error('This browser/authenticator does not support the PRF extension.');
  }
}
```

### Authentication with PRF

```js
async function derivePasskeySecret(credentialIdBase64, prfSaltBase64) {
  const assertion = await navigator.credentials.get({
    publicKey: {
      challenge: crypto.getRandomValues(new Uint8Array(32)),
      allowCredentials: [{ id: fromBase64(credentialIdBase64), type: 'public-key' }],
      userVerification: 'required',
      timeout: 60000,
      extensions: { prf: { eval: { first: fromBase64(prfSaltBase64) } } }
    }
  });
  const prf = assertion.getClientExtensionResults()?.prf;
  const secret = prf?.results?.first ?? prf?.first ?? null;
  if (!secret) throw new Error('PRF not supported by this authenticator');
  return new Uint8Array(secret);
}
```

**base64 vs base64url:** WebAuthn uses base64url (RFC 4648 §5) for credential IDs. Mixing base64 and base64url causes lookup failures. Use consistent encoding throughout.

### Browser support (2025–2026)

| Browser | PRF support |
|---|---|
| Chrome 120+ | Yes (platform + roaming authenticators) |
| Edge 120+ | Yes (follows Chromium) |
| Safari 18+ (macOS 15, iOS 18) | Partial (platform authenticators) |
| Firefox | Not yet supported |

### Review checklist for passkey PRF

- [ ] PRF salt is at least 32 bytes from `crypto.getRandomValues()`
- [ ] PRF salt unique per credential enrollment and stored in the envelope
- [ ] PRF output hashed (SHA-256) or run through HKDF before use as AES key
- [ ] `userVerification: 'required'` set on both create and get
- [ ] Code checks for PRF support and falls back gracefully
- [ ] `credentialId` stored alongside the wrapped envelope
- [ ] Error handling covers timeout, user cancellation, and missing PRF support

---

## Security Checklist (Full Review)

### Randomness and entropy
- [ ] All random values use `crypto.getRandomValues()` — never `Math.random()`
- [ ] Master key is 32 bytes of CSPRNG output
- [ ] IVs are 12 bytes of CSPRNG output, generated fresh per encrypt call
- [ ] Salts are at least 16 bytes of CSPRNG output, unique per derivation

### Key derivation
- [ ] PBKDF2 iterations >= 310,000 (prefer >= 600,000 for new code)
- [ ] PBKDF2 uses SHA-256 or SHA-512
- [ ] Iteration count stored in envelope metadata for upgradability
- [ ] High-entropy sources (auth files, PRF output) use HKDF, not PBKDF2
- [ ] HKDF `info` field includes a domain-separation string

### Key management
- [ ] All CryptoKeys created with `extractable: false` unless wrapping requires otherwise
- [ ] Key usages minimally scoped
- [ ] Master key bytes zeroed on vault lock where possible
- [ ] Master key never written to persistent storage in cleartext
- [ ] Unlocked key material lives in `chrome.storage.session` or in-memory only
- [ ] `chrome.storage.session.setAccessLevel` defaults to `TRUSTED_CONTEXTS` unless content scripts need session data

### Encryption
- [ ] AES-GCM with 256-bit keys
- [ ] IV is 12 bytes, unique per encryption, stored with ciphertext
- [ ] Tag length is 128 bits (default — never reduced)
- [ ] `OperationError` from `decrypt()` treated as authentication failure

### Vault envelope structure
- [ ] Envelope includes a version field
- [ ] Multiple unlock methods each independently wrap the same master key
- [ ] Each unlock method has its own salt and IV
- [ ] Envelope stored in durable storage in encrypted form only

### Timing and side channels
- [ ] Password comparison uses constant-time HMAC, not `===`
- [ ] Token/HMAC verification uses `crypto.subtle.verify()` or double-HMAC pattern

### Error handling
- [ ] Decryption errors do not leak whether password/key was wrong vs. data was corrupted
- [ ] Error messages do not include key material, IVs, or salt values
- [ ] Failed unlock attempts are rate-limited in the UI
- [ ] Passkey PRF fallback path exists for browsers without PRF support

---

## Anti-Patterns and Vulnerabilities

| Anti-pattern | Severity | Impact | Fix |
|---|---|---|---|
| IV reuse | Critical | Leaks `P1 XOR P2`, recovers auth key — complete cryptographic break | Fresh `getRandomValues(12)` per encrypt call |
| Weak PBKDF2 iterations (< 310K) | High | Offline brute-force on weak passwords | Use 600,000 iterations |
| Extractable master key | High | XSS can call `exportKey()` and exfiltrate raw bytes | `extractable: false` |
| `Math.random()` for crypto values | High | Deterministic, predictable output | `crypto.getRandomValues()` |
| Master key in durable storage cleartext | High | Persists across restarts, readable by any extension code | Session storage only; durable only when wrapped |
| Reduced tag length (tagLength: 32) | High | 32-bit tag brute-forced with ~2^32 attempts | Use default 128-bit tag |
| Missing PRF support check | Medium | Runtime error breaks unlock flow on Firefox/older Safari | Check `getClientExtensionResults()?.prf` before accessing |
| PBKDF2 for high-entropy secrets | Medium | ~100ms latency penalty with no security gain | Use HKDF for auth files and PRF output |
| Key-commitment vulnerability | Low (single-user vault) | Attacker controlling two wrapping keys could craft ambiguous ciphertext | Add HMAC commitment or key fingerprint in multi-key scenarios |

---

## Troubleshooting

| Error | Likely cause | Fix |
|---|---|---|
| `OperationError` on decrypt | IV mismatch, wrong key, tampered ciphertext, or truncated storage | Check base64 round-trips; verify salt/iterations match |
| `InvalidAccessError` on importKey | Wrong `usages` array or wrong algorithm name | Match usages to intended operation |
| `NotSupportedError` on deriveKey | Unsupported algorithm in older WebViews; wrong key type | PBKDF2 base material must be imported as `'PBKDF2'`, not `'AES-GCM'` |
| PRF returns null | Browser or authenticator doesn't support PRF; result path varies | Check both `prf.results.first` and `prf.first` |
| PBKDF2 too slow (>500ms) | Blocking main thread | Run derivation in a Web Worker; do not reduce iterations |

**PBKDF2 in a Web Worker:**
```js
// worker.js — perform the full operation inside the worker
self.onmessage = async ({ data }) => {
  const key = await derivePasswordKey(data.password, data.salt, data.iterations);
  const masterKeyBytes = await unwrapMasterKey(key, data.envelope);
  // ArrayBuffer is transferable
  self.postMessage({ masterKeyBytes }, [masterKeyBytes.buffer]);
};
```

---

## References

- [MDN: SubtleCrypto](https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto)
- [OWASP: Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [NIST SP 800-132: Password-Based Key Derivation](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-132.pdf)
- [Yubico: PRF Extension Developer Guide](https://developers.yubico.com/WebAuthn/Concepts/PRF_Extension/Developers_Guide_to_PRF.html)
- [RFC 8452: AES-GCM-SIV](https://datatracker.ietf.org/doc/rfc8452/)
- [elttam: Attacks on GCM with Repeated Nonces](https://www.elttam.com/blog/key-recovery-attacks-on-gcm)
