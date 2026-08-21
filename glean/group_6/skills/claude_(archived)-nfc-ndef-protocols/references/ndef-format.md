# NDEF format — byte-level detail

Companion to the SKILL.md "NDEF" section. Grounded in the NFC Forum NDEF and
RTD/Text/URI/Smart Poster specs; verify byte-exact conformance against the
purchased specs.

## Record header byte (bit 7 → bit 0)

| Bits | Field | Meaning |
| --- | --- | --- |
| 7 | MB | Message Begin — first record of the message |
| 6 | ME | Message End — last record of the message |
| 5 | CF | Chunk Flag — this record is one chunk of a chunked payload |
| 4 | SR | Short Record — SR=1 ⇒ PAYLOAD LENGTH is 1 byte; SR=0 ⇒ 4 bytes (big-endian) |
| 3 | IL | ID Length present — an ID LENGTH byte and ID field are included |
| 2–0 | TNF | Type Name Format (3-bit) |

A single record sets both MB and ME. A chunked payload: first chunk carries the
real TNF/TYPE with CF=1; middle/last chunks use **TNF 0x06 (Unchanged)** with
TYPE LENGTH 0; the last chunk clears CF.

## Field order

`header` → `TYPE LENGTH (1B)` → `PAYLOAD LENGTH (1B if SR=1 else 4B big-endian)`
→ `ID LENGTH (1B, only if IL=1)` → `TYPE (TYPE LENGTH bytes)` →
`ID (ID LENGTH bytes, only if IL=1)` → `PAYLOAD (PAYLOAD LENGTH bytes)`.

## TNF values (0x00–0x07)

| TNF | Name | TYPE field |
| --- | --- | --- |
| 0x00 | Empty | none (TYPE/ID/PAYLOAD lengths all 0) |
| 0x01 | NFC Forum Well-Known Type | RTD name, e.g. `T`, `U`, `Sp` |
| 0x02 | MIME media-type (RFC 2046) | e.g. `text/vcard`, `application/json` |
| 0x03 | Absolute URI (RFC 3986) | a URI that names the type |
| 0x04 | NFC Forum External Type | namespaced, e.g. `android.com:pkg` |
| 0x05 | Unknown | TYPE LENGTH must be 0 |
| 0x06 | Unchanged | chunk continuation only |
| 0x07 | Reserved | — |

## Text RTD (`T`, 0x54)

Payload = `status byte` + `language code` + `text`:
- **Status byte:** bit 7 = encoding (0 = UTF-8, 1 = UTF-16); bit 6 reserved (0);
  bits 5–0 = length of the language tag.
- **Language code:** IETF BCP 47 tag (e.g. `en`, `en-US`), ASCII, length from
  the status byte.
- **Text:** the remaining payload bytes in the declared encoding. No length
  prefix — it runs to the end of the record payload.

## URI RTD (`U`, 0x55)

Payload = `1 identifier-code byte` + `URI field`. The identifier code prepends a
common prefix so it need not be stored:

| Code | Prefix | Code | Prefix |
| --- | --- | --- | --- |
| 0x00 | (none) | 0x05 | `tel:` |
| 0x01 | `http://www.` | 0x06 | `mailto:` |
| 0x02 | `https://www.` | 0x07 | `ftp://anonymous:anonymous@` |
| 0x03 | `http://` | 0x08 | `ftp://ftp.` |
| 0x04 | `https://` | 0x1D | `file://` |

The table continues to 0x23 (codes 0x09–0x23 cover ftps/sftp, smb, nfs, various
`urn:` schemes, `tel`/`sms` variants, etc.); 0x24–0xFF are RFU. Consult the NFC
Forum URI RTD spec for the complete list.

Example: `https://www.example.com` → identifier code `0x02`, URI field
`example.com`.

## Smart Poster RTD (`Sp`)

A container record whose **payload is itself an NDEF message**. Components:
- **URI record** — exactly one, mandatory.
- **Title (Text) records** — zero or more, one per language (a language MUST NOT
  repeat).
- **Action record** — optional; suggests do-the-action / save-for-later /
  open-for-editing.
- **Icon (MIME image), Type, Size** records — optional metadata.

## Storing NDEF on a tag

NDEF on a tag is wrapped in TLV structures (Type 1/2/5) or files (Type 4) — see
`tag-types-1-5.md`.
