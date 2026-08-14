# OSDP protocol detail — command/reply set, framing, addressing, Secure Channel handshake

Companion to `SKILL.md` §3. Deeper reference for the OSDP (Open Supervised Device Protocol /
IEC 60839-11-5) wire protocol. Scope here is the protocol mechanics; **attack execution
against Secure Channel (install-mode key injection, SCBK-D exploitation, downgrade) belongs
to `rfid-nfc-access-attacks-defenses`** — this file documents the protocol so defenders can
configure and verify it, not how to break it.

## Bus & addressing

- **Physical:** half-duplex RS-485, asynchronous serial, **8 data bits, 1 stop bit, no
  parity**.
- **Baud rates:** 9600, 19200, 38400, 115200, 230400. **Default 9600.**
- **Roles:** one **CP** (Control Panel / ACU) is the bus master; one or more **PD**
  (Peripheral Device — reader, keypad, I/O) respond. **All communication is CP-initiated.**
- **Addressing:** up to **128 PDs**, address **0–127 (0x00–0x7E)**; **0x7F is the broadcast
  address** (all PDs act on it). **Default PD address 0.** A PD replies only when its address
  is in the request.
- **Reconfiguration:** address and baud are changed at runtime with **`osdp_COMSET`**.
- **Polling cadence:** the CP polls each PD frequently (libosdp recommends servicing the
  state machine at least every ~50 ms; an idle Secure Channel session times out around
  400 ms without traffic). Continuous polling is what provides *supervision*.

## Packet framing (overview)

Each OSDP packet carries: a start-of-message marker (`0x53`), the PD address byte, a length
field, a message-control byte (sequence number + Secure Channel block usage), the
**command/reply code**, the data payload, and a trailing checksum/CRC. When Secure Channel
is active, security control blocks and a per-message **MAC** are added and the payload is
AES-128 encrypted. (Exact byte offsets are in the SIA OSDP spec / IEC 60839-11-5; treat the
above as the conceptual layout, not a byte map.)

## Command codes (CP → PD)

| Command | Hex | Purpose |
| --- | --- | --- |
| `osdp_POLL` | 0x60 | Periodic poll / heartbeat (drives supervision). |
| `osdp_ID` | 0x61 | ID report request → `osdp_PDID`. |
| `osdp_CAP` | 0x62 | Capabilities request → `osdp_PDCAP`. |
| `osdp_LSTAT` | 0x64 | Local status report request. |
| `osdp_ISTAT` | 0x65 | Input status report request. |
| `osdp_OSTAT` | 0x66 | Output status report request. |
| `osdp_RSTAT` | 0x67 | Reader tamper/status report request. |
| `osdp_OUT` | 0x68 | Output control (relays, etc.). |
| `osdp_LED` | 0x69 | Reader LED control. |
| `osdp_BUZ` | 0x6A | Reader buzzer/beeper control. |
| `osdp_TEXT` | 0x6B | Reader text-display output. |
| `osdp_COMSET` | 0x6E | Set PD communication config (address, baud). |
| `osdp_KEYSET` | 0x75 | Encryption key set (provision/replace SCBK). |
| `osdp_CHLNG` | 0x76 | Challenge & Secure Channel session init (CP random). |
| `osdp_SCRYPT` | 0x77 | Server random + server cryptogram. |
| `osdp_MFG` | 0x80 | Manufacturer-specific command. |

## Reply codes (PD → CP)

| Reply | Hex | Purpose |
| --- | --- | --- |
| `osdp_ACK` | 0x40 | General acknowledge, nothing to report. |
| `osdp_NAK` | 0x41 | Negative acknowledge (error/unsupported). |
| `osdp_PDID` | 0x45 | Device identification report. |
| `osdp_PDCAP` | 0x46 | Device capabilities report. |
| `osdp_RAW` | 0x50 | **Card data report — raw bit array** (the Wiegand *format* rides here). |
| `osdp_KEYPAD` | 0x53 | Keypad / PIN data report. |
| `osdp_CCRYPT` | 0x76 | Client ID + client random number + client cryptogram. |
| `osdp_RMAC_I` | 0x78 | Client cryptogram packet + initial R-MAC. |
| `osdp_BUSY` | 0x79 | PD busy, retry later. |

> Hex values above are from the libosdp command/reply reference and match the SIA OSDP
> code assignments; verify against the version of the spec your devices implement
> (`osdp_KEYPAD` and `osdp_CCRYPT` both appearing as 0x53/0x76 in the reply vs command
> spaces is expected — command and reply code spaces are distinct).

## Secure Channel (SC) — keys, handshake, MAC

### Keys
- **SCBK** — Secure Channel Base Key, **128-bit per-PD** AES key. The root from which session
  keys derive.
- **SCBK-D** — the **default base key published in the OSDP spec**, for install/provisioning
  only. **Using SCBK-D in production removes all protection** (public key → anyone can
  negotiate the channel). Provisioning must `osdp_KEYSET` a unique site SCBK.
- **Session keys** derived per SC setup: **S-ENC** (AES-128 payload encryption), **S-MAC1**
  and **S-MAC2** (message authentication). They are functions of the SCBK and the handshake
  random numbers, so each session is fresh (replay-resistant).

### Handshake (mutual challenge-response, AES-128)
1. **CP → PD `osdp_CHLNG`** — carries the CP's random number (RND.A).
2. **PD → CP `osdp_CCRYPT`** — carries the PD's random number (RND.B) and the **client
   cryptogram**; the CP verifies it to **authenticate the PD**.
3. **CP → PD `osdp_SCRYPT`** — carries the **server cryptogram**; the PD verifies it to
   **authenticate the CP**.
4. **PD → CP `osdp_RMAC_I`** — the initial reply-MAC. Both ends now hold matching S-ENC /
   S-MAC keys; the channel is established.

After setup, every secured message carries a MAC; the reply-MAC chains forward (used as the
IV/seed for the next message's authentication), giving **message sequencing + tamper
evidence**. An `OSDP_FLAG_ENFORCE_SECURE`-style policy makes the CP refuse to fall back to
cleartext — important so an attacker cannot force a downgrade (downgrade *exploitation* is
covered in the security sibling).

### Install mode & key provisioning
- A PD with **no SCBK set enters install mode** and will accept an `osdp_KEYSET` carrying a
  fresh SCBK; on success it **exits install mode and refuses further key-set** in that state.
- Defensive guidance (the *what to do*, not the attack): **never deploy with SCBK-D**,
  provision a **unique SCBK per reader**, prefer secure provisioning (e.g. a configuration
  card presented to the reader) over a physically-triggered install switch, and enforce
  secure-only operation. The residual physical-attack path on install mode is documented in
  `rfid-nfc-access-attacks-defenses`.

## Sources
- libosdp / doc.osdp.dev — index, commands & replies, Secure Channel:
  https://doc.osdp.dev/index.html ,
  https://doc.osdp.dev/protocol/commands-and-replies.html ,
  https://doc.osdp.dev/libosdp/secure-channel.html
- libosdp introduction (bus parameters, baud rates):
  https://libosdp.gotomain.io/protocol/introduction.html
- OSDP World — OSDP device communication (addressing, broadcast 0x7F, defaults):
  https://osdpworld.com/2021/04/03/osdp-device-communication/
- Security Industry Association — OSDP standard (IEC 60839-11-5; v2.2.2; AES-128):
  https://www.securityindustry.org/industry-standards/open-supervised-device-protocol/
