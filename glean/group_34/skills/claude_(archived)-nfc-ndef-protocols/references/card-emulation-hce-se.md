# Card emulation — HCE vs Secure Element, AID routing, APDU flow

Companion to the SKILL.md "Card Emulation" section. Grounded in the Android HCE
developer docs, the Secure Technology Alliance "HCE 101" whitepaper, and SE-vs-HCE
payment-security analyses. Concepts, not production code.

## The reader-side handshake

In card-emulation mode the phone is the **target**. It comes up over ISO-DEP
(ISO/IEC 14443-4) and speaks **ISO/IEC 7816-4 APDUs**. The exchange:

1. Reader powers the field, anticollision/activation selects the phone as a card.
2. Reader sends a **SELECT (by name / AID)** APDU:
   `00 A4 04 00 <Lc> <AID> 00`, where `<AID>` is up to 16 bytes (ISO/IEC 7816-5).
3. The NFC controller matches the AID against its **routing table** and forwards
   subsequent APDUs to the matched destination (host service or secure element).
4. Application-specific command/response APDUs flow until a new SELECT (different
   AID) or the RF link drops.

## AID routing table

Each entry maps an **AID** (or AID prefix) to a **destination**: the host CPU
(HCE) or an off-host secure element (eSE / UICC).

- **AID groups** are routed **atomically**: every AID in a group goes to the same
  destination, or none do — never a partial split.
- **Categories:** `CATEGORY_PAYMENT` (only one payment AID group active
  system-wide, tied to the default wallet) and `CATEGORY_OTHER` (may be always
  active; loyalty, transit, access).
- **Conflict resolution:** default wallet wins → else the sole registered service
  → else the OS prompts the user to choose.

## HCE service (host path)

`HostApduService` (Android 4.4+):
- `processCommandApdu(commandApdu, extras)` — called on the main thread with each
  APDU; return the response bytes, or `null` and call `sendResponseApdu()` later
  for async work.
- `onDeactivated(reason)` — RF link lost or a different AID selected.
- One logical channel.
- Declared in the manifest with action
  `android.nfc.cardemulation.action.HOST_APDU_SERVICE`, permission
  `BIND_NFC_SERVICE`, and an `apduservice` XML resource declaring the AID groups.

## Off-host service (SE path)

`OffHostApduService` declares which AIDs route to the secure element
(`android:secureElementName="eSE"` or a UICC). Android never starts/binds it —
the transaction runs entirely on the SE; the app can query status afterward.
`requireDeviceUnlock` does not apply to off-host AIDs.

## HCE vs SE — the security model

| | HCE | Secure Element |
| --- | --- | --- |
| Credential custody | host / app / cloud | tamper-resistant certified chip |
| Trust assumption | host is compromisable | hardware is the trust anchor |
| Mitigations | tokenization, **limited-use keys**, device fingerprinting, transaction risk analysis | EMV-grade, payment-network evaluated |
| Suited for | loyalty, transit, lower-assurance access, fast time-to-market | open-loop payments, identity, high-assurance access |
| Examples | Google Pay, Samsung Pay (tokenized cloud creds) | Apple Pay (tokens/keys in the SE) |

**The SE's role in payments/access:** it stores the keys and (for SE wallets) the
tokenized PAN, computes the EMV cryptogram in hardware, and is the trust anchor
that lets the phone substitute for a plastic EMV card or a high-assurance access
credential. EMV payment-kernel internals are out of scope for this skill.
