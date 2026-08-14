# NFC Forum Tag Types 1–5 — capability containers & TLVs

Companion to the SKILL.md "NFC Forum Tag Types 1–5" section. Grounded in the NFC
Forum Type 1–5 Tag specs and the nfcpy tag documentation; confirm per-chip
memory against the specific datasheet.

## Type → technology → chip

| Type | NFC tech / standard | Representative chips | Spec (current) |
| --- | --- | --- | --- |
| 1 | NFC-A; ISO/IEC 14443A frame + Jewel/Topaz commands | Broadcom (Innovision) Topaz/Jewel | withdrawn in TS Release 2021 |
| 2 | NFC-A; ISO/IEC 14443A | NXP MIFARE Ultralight, NTAG21x | v1.3 |
| 3 | NFC-F; FeliCa / JIS X 6319-4 (ISO/IEC 18092) | Sony FeliCa Lite | v1.1 |
| 4 | NFC-A **and** NFC-B; ISO-DEP (ISO/IEC 14443-4); ISO/IEC 7816-4 APDUs | NXP DESFire, SmartMX; phone HCE | v1.2 |
| 5 | NFC-V; ISO/IEC 15693 | NXP ICODE, ST ST25TV | v1.3 |

## Capability Container (CC) per type

The CC tells a reader the tag's NDEF capability, version, size, and access rights.

- **Type 1 / Type 2 / Type 5:** a small fixed CC block plus **TLV**-wrapped data.
  - **Type 2 CC** is 4 bytes at block 3: byte 0 magic `0xE1`; byte 1 version
    (high/low nibble) + read/write access; byte 2 memory size in units of 8
    bytes; byte 3 read/write access conditions.
- **Type 3 (FeliCa):** an **Attribute Information Block** plays the CC role
  (NDEF version, number of blocks readable/writable in one command, max size,
  write flag, RW flag, checksum).
- **Type 4:** a **Capability Container file** (file ID `E103`), selected by APDU,
  whose contents include the NDEF File Control TLV pointing at the **NDEF file**.

## TLV structures (Type 1/2/5)

Data area is a sequence of Type-Length-Value blocks:

| Tag | TLV | Purpose |
| --- | --- | --- |
| 0x00 | NULL | padding, ignored |
| 0x01 | Lock Control | location/size of dynamic lock bits |
| 0x02 | Memory Control | reserved memory areas |
| 0x03 | NDEF Message | the NDEF message (Length then the message bytes) |
| 0xFD | Proprietary | vendor-defined |
| 0xFE | Terminator | end of the TLV area (1 byte, no length/value) |

Length encoding: one byte for 0–254; for ≥255, a 3-byte form `0xFF` + 2-byte
big-endian length.

## Reading NDEF

- **Type 1/2/5:** read the CC, walk the TLVs to the **NDEF Message TLV (0x03)**,
  read its value.
- **Type 3:** read the Attribute Information Block, then the data blocks via
  FeliCa Check/Update commands.
- **Type 4:** `SELECT` the NDEF application/CC file, `ReadBinary` the CC to learn
  the NDEF file ID and max size, `SELECT` the NDEF file, `ReadBinary` (first 2
  bytes are the NDEF length).

## Formatting a blank tag

A freshly formatted Type 2 tag gets a CC, optional Lock Control / Memory Control
TLVs, and an **empty NDEF Message TLV** (`0x03 0x00`) followed by a Terminator
TLV (`0xFE`) — the same shape nfcpy writes when it formats a Topaz/Type 2 tag.
