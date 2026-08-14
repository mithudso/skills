# Phone NFC APIs — Android vs Apple capability matrix

Companion to the SKILL.md "Phone NFC APIs" section. Grounded in the Android
`android.nfc` docs and Apple Core NFC docs. Conceptual; verify current
availability/entitlement details against the live developer docs.

## Capability matrix

| Capability | Android | Apple (Core NFC) |
| --- | --- | --- |
| Read NDEF | `Ndef` / `NdefMessage` / `NdefRecord` via `NfcAdapter` foreground dispatch or Reader Mode | `NFCNDEFReaderSession` |
| Write NDEF | yes (`Ndef.writeNdefMessage`) | yes, on supported devices (`NFCNDEFReaderSession` write) |
| Raw tag tech | `NfcA/NfcB/NfcF/NfcV`, `IsoDep`, `MifareUltralight`, `MifareClassic`, `NdefFormatable` | `NFCTagReaderSession` → ISO7816, ISO15693, FeliCa, MIFARE |
| Reader trigger | foreground dispatch or `enableReaderMode` | explicit session (user-initiated) + background tag reading |
| Background read | foreground/Reader Mode while app active; some OEM background scan | OS background read of an NDEF URL record (Universal Link), no app open |
| Card emulation | `HostApduService` (HCE) + `OffHostApduService` (SE) | reserved for Apple Pay (SE); third-party via gated entitlement |
| P2P | deprecated/removed (Android Beam) | never shipped |

## Android specifics

- **Reader/Writer:** `NfcAdapter` + foreground dispatch, or `enableReaderMode`
  with flags to select technologies and skip platform NDEF checks. `Tag`
  dispatched via intent (`ACTION_NDEF_DISCOVERED` etc.).
- **Card emulation:** see `card-emulation-hce-se.md` — `HostApduService` (HCE),
  `OffHostApduService` (SE), `CardEmulation` for defaults.
- **P2P:** `NdefMessage` push (Android Beam) — **deprecated in Android 10,
  removed in Android 14**. Do not build new features on it.

## Apple specifics

- **`NFCNDEFReaderSession`** — the easy path; reads (and on supported devices
  writes) NDEF from Tag Types 1–5. Requires the NFC entitlement and a usage
  string; sessions are foreground and user-initiated.
- **`NFCTagReaderSession`** — lower-level raw access for ISO 7816, ISO 15693,
  FeliCa, MIFARE.
- **Background tag reading** (iPhone XS+, iOS 12+): the OS reads an NDEF tag in
  the background **with no app open** if the first NDEF record is a **URL record
  that is a valid Apple Universal Link**; only the first URL record is processed,
  and it routes the URL to the matching app. Suppressed while a reader session,
  Wallet/Apple Pay, or the camera is active, in Airplane Mode, or before the
  first unlock after boot.
- **Card emulation:** historically Apple Pay only (Secure Element). Third-party
  HCE-style emulation opened via a gated entitlement on recent iOS — treat as
  region/program-gated and verify Apple's current NFC & SE / PassKit docs before
  relying on it.

## Reading availability caveats

NFC entitlement, device model, and iOS/Android version all gate these APIs.
Background reading needs iPhone XS+; third-party iOS card emulation is
entitlement-gated and has shifted across releases. Confirm against the live docs.
