<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `chrome-native-messaging` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: chrome-native-messaging
description: >
  Chrome extension native messaging — writing a native host in Python or Node.js,
  registering the host manifest on macOS/Linux/Windows, wiring
  chrome.runtime.connectNative or sendNativeMessage in an MV3 service worker,
  debugging "host not found" or protocol errors, and understanding the 4-byte
  framing protocol and 1 MB message size limit.
  TRIGGER: implementing or debugging chrome.runtime.connectNative or
  sendNativeMessage, writing a Python/Node.js native host, registering the host
  manifest, troubleshooting "Specified native messaging host not found" or
  "Access to the specified native messaging host is forbidden", or understanding
  the 4-byte framing protocol.
  SKIP: general IPC between browser tabs (use chrome.runtime.sendMessage),
  WebSocket or HTTP localhost communication (no native host needed),
  extension-to-extension messaging, or browser automation outside Chrome.
version: 1.1.0
category: developer
tags: [chrome-extension, native-messaging, mv3, service-worker, python, nodejs, framing, protocol, manifest]
related_skills: [mv3-service-worker-expert, chrome-mv3-advanced, chrome-extension-security-reviewer]
updated: 2026-05-29
---

# Chrome Native Messaging

## Overview

Chrome extensions communicate with local native applications via **stdin/stdout using a 4-byte length prefix** (native byte order, little-endian on all common platforms) followed by a UTF-8 JSON payload. The native host is a standalone executable registered via a manifest file placed in an OS-specific directory — it is not bundled inside the extension.

## Protocol — 4-Byte Framing

Every message in both directions uses the same wire format:

```
[4 bytes: uint32, native byte order][N bytes: UTF-8 JSON payload]
```

- Length prefix is a **32-bit unsigned integer in native byte order** (little-endian on x86/x64/ARM; use `@I` in Python struct, `readUInt32LE`/`writeUInt32LE` in Node.js Buffer)
- **Max message native host → Chrome:** 1 MB (1,048,576 bytes) — Chrome rejects larger responses
- **Max message Chrome → native host:** 64 MiB
- Debug output **must** go to **stderr only** — any bytes on stdout corrupt the framing

## Host Manifest

Place a JSON file named `com.company.app.json` in the OS-specific directory below.

```json
{
  "name": "com.my_company.my_app",
  "description": "My native host",
  "path": "/usr/local/bin/my-native-host",
  "type": "stdio",
  "allowed_origins": ["chrome-extension://EXTENSION_ID_HERE/"]
}
```

| Field | Notes |
|---|---|
| `name` | Lowercase alphanumeric, underscores, dots only. No leading/trailing dots. |
| `path` | Absolute on macOS/Linux. May be relative on Windows. |
| `type` | Always `"stdio"`. |
| `allowed_origins` | Extension origin(s), no wildcards. Format: `chrome-extension://ID/` |

## OS Registration Paths

### macOS

| Scope | Path |
|---|---|
| System-wide | `/Library/Google/Chrome/NativeMessagingHosts/com.my_company.my_app.json` |
| User | `~/Library/Application Support/Google/Chrome/NativeMessagingHosts/com.my_company.my_app.json` |
| Chromium (user) | `~/Library/Application Support/Chromium/NativeMessagingHosts/com.my_company.my_app.json` |

### Linux

| Scope | Chrome | Chromium |
|---|---|---|
| System-wide | `/etc/opt/chrome/native-messaging-hosts/com.my_company.my_app.json` | `/etc/chromium/native-messaging-hosts/com.my_company.my_app.json` |
| User | `~/.config/google-chrome/NativeMessagingHosts/com.my_company.my_app.json` | `~/.config/chromium/NativeMessagingHosts/com.my_company.my_app.json` |

### Windows

Registry key value points to the manifest JSON file path:

```
HKEY_CURRENT_USER\SOFTWARE\Google\Chrome\NativeMessagingHosts\com.my_company.my_app
  (Default) = "C:\path\to\com.my_company.my_app.json"
```

Also register under `HKEY_LOCAL_MACHINE` for system-wide access. On 64-bit Windows, Chrome checks the 32-bit registry hive (`HKLM\SOFTWARE\WOW6432Node\Google\Chrome\NativeMessagingHosts\...`) first.

## Python Host Implementation

```python
#!/usr/bin/env python3
# Run as: python3 -u host.py
# -u flag required on Windows for binary stdin/stdout (or use the .bat wrapper below)

import sys
import json
import struct

def read_message():
    raw_len = sys.stdin.buffer.read(4)
    if len(raw_len) == 0:
        sys.exit(0)
    msg_len = struct.unpack('@I', raw_len)[0]  # '@' = native byte order
    return json.loads(sys.stdin.buffer.read(msg_len).decode('utf-8'))

def send_message(msg):
    payload = json.dumps(msg, separators=(',', ':')).encode('utf-8')
    sys.stdout.buffer.write(struct.pack('@I', len(payload)))
    sys.stdout.buffer.write(payload)
    sys.stdout.buffer.flush()

while True:
    try:
        msg = read_message()
        send_message({'echo': msg})
    except Exception as e:
        sys.stderr.write(f'Error: {e}\n')  # stderr only — never print() to stdout
        sys.exit(1)
```

**Windows wrapper** (`host.bat` — point manifest `path` here):
```batch
@echo off
python -u "C:\path\to\host.py"
```

## Node.js Host Implementation

```javascript
#!/usr/bin/env node
// CommonJS. Shebang not supported on Windows — invoke with: node host.js

process.stdin.resume();
process.stdin.once('end', () => process.exit(0));

function readExact(n) {
  // Accumulates chunks until exactly n bytes are collected — handles partial reads.
  return new Promise((resolve) => {
    const parts = [];
    let remaining = n;
    const tryRead = () => {
      while (remaining > 0) {
        const chunk = process.stdin.read(remaining);
        if (chunk === null) { process.stdin.once('readable', tryRead); return; }
        parts.push(chunk);
        remaining -= chunk.length;
      }
      resolve(Buffer.concat(parts));
    };
    tryRead();
  });
}

async function readMessage() {
  const header = await readExact(4);
  const len = header.readUInt32LE(0);  // native byte order (LE on x86/ARM)
  const body = await readExact(len);
  return JSON.parse(body.toString('utf-8'));
}

function sendMessage(msg) {
  const payload = Buffer.from(JSON.stringify(msg), 'utf-8');
  const header = Buffer.alloc(4);
  header.writeUInt32LE(payload.length, 0);
  process.stdout.write(header);
  process.stdout.write(payload);
}

(async () => {
  while (true) {
    try {
      const msg = await readMessage();
      sendMessage({ echo: msg });
    } catch (e) {
      process.stderr.write(`Error: ${e}\n`);  // stderr only — never console.log to stdout
      process.exit(1);
    }
  }
})();
```

## install.sh Registration Script (macOS/Linux)

A shell script is the standard way to register the host manifest with the correct extension ID and absolute path:

```bash
#!/bin/bash
# install.sh — registers the native messaging host manifest
# Usage: ./install.sh <extension-id>

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOST_NAME="com.my_company.my_app"
EXTENSION_ID="${1:?Usage: install.sh <extension-id>}"

# Write manifest with correct extension ID and absolute path to host binary
cat > "${SCRIPT_DIR}/manifest.json" <<EOF
{
  "name": "${HOST_NAME}",
  "description": "My native messaging host",
  "path": "${SCRIPT_DIR}/host.py",
  "type": "stdio",
  "allowed_origins": ["chrome-extension://${EXTENSION_ID}/"]
}
EOF

# Copy to OS registration path
if [[ "$OSTYPE" == "darwin"* ]]; then
  HOSTS_DIR=~/Library/Application\ Support/Google/Chrome/NativeMessagingHosts
else
  HOSTS_DIR=~/.config/google-chrome/NativeMessagingHosts
fi

mkdir -p "$HOSTS_DIR"
cp "${SCRIPT_DIR}/manifest.json" "${HOSTS_DIR}/${HOST_NAME}.json"
echo "Installed: ${HOSTS_DIR}/${HOST_NAME}.json"
echo "IMPORTANT: Restart Chrome after installation."
```

Run this script after any change to host code or path:
1. Re-run `install.sh <extension-id>` if the path or `allowed_origins` changed
2. **Restart Chrome** — native host manifests are read at Chrome startup, not on extension reload
3. Re-run the install step when Chrome assigns a new unpacked extension ID

## MV3 Service Worker Integration

Declare permission in `manifest.json`:
```json
{ "permissions": ["nativeMessaging"] }
```

**One-shot message** (spawns new host process per call):
```javascript
const response = await chrome.runtime.sendNativeMessage('com.my_company.my_app', { cmd: 'ping' });
```

**Persistent connection** (host lives until port is destroyed — keeps SW alive):
```javascript
const port = chrome.runtime.connectNative('com.my_company.my_app');

port.onMessage.addListener((msg) => {
  console.log('Host sent:', msg);
});

port.onDisconnect.addListener(() => {
  if (chrome.runtime.lastError) {
    console.error('Disconnected:', chrome.runtime.lastError.message);
  }
});

port.postMessage({ cmd: 'start' });
// port.disconnect() when done
```

**Key constraints:**
- `connectNative` is only available in **extension pages and service worker** — not in content scripts
- An open `connectNative` port keeps the MV3 service worker alive as long as the port remains connected
- `sendNativeMessage` starts a new host process per call; only the **first response** is delivered — the host should write one response then exit cleanly
- To find the extension ID for `allowed_origins` during development: open `chrome://extensions`, enable Developer mode, copy the ID shown under the extension name

## Error Reference

| Error message | Cause | Fix |
|---|---|---|
| "Specified native messaging host not found" | Wrong manifest filename or location | Filename must match `name` field exactly; check OS path |
| "Access to the specified native messaging host is forbidden" | Extension ID not in `allowed_origins` | Add `chrome-extension://YOUR_ID/` to manifest |
| "Error when communicating with the native messaging host" | Protocol violation — bad stdout output | Remove all `print()`/`console.log()` to stdout; use stderr |
| "Failed to start native messaging host" | Missing execute permission or bad path | `chmod +x host.py`; verify path is absolute on macOS/Linux |
| "Invalid native messaging host name" | Name contains invalid chars | Only lowercase alphanum, `_`, `.` allowed |
| Host crashes silently | Unhandled exception or wrong binary mode | Add try/except; use `-u` flag on Python/Windows |

## Debugging

1. Check `chrome://extensions` → click "Errors" on the extension card for native messaging errors
2. Start Chrome from terminal on macOS/Linux to see stderr output from the host in the shell
3. **Windows**: launch Chrome with `--enable-logging`; logs in `%LOCALAPPDATA%\Google\Chrome\User Data\chrome_debug.log`
4. Validate manifest JSON: `python3 -m json.tool com.my_company.my_app.json`
5. Test host in isolation: `echo -ne '\x09\x00\x00\x00{"cmd":1}' | python3 host.py`

## Security Considerations

1. **Extension ID validation**: Chrome enforces `allowed_origins` at the OS level — the host only starts if the calling extension is listed. No wildcards allowed.
2. **Input sanitization**: Validate all JSON fields received from the extension. Content scripts are less trusted than the service worker and can be compromised by malicious pages.
3. **Minimal `allowed_origins`**: one entry per legitimate extension.
4. **Host executable permissions**: `chmod 755`; do not make writable by untrusted users.
5. **No localhost server needed**: native messaging is a direct process pipe — do not expose an HTTP port just for this purpose.
6. **First argument is the caller origin**: Chrome passes `chrome-extension://ID` as the first CLI argument — use it to double-check the caller when `allowed_origins` contains multiple entries.

## Common Mistakes

- Printing debug text to **stdout** breaks framing — always use stderr
- Omitting **`sys.stdout.buffer.flush()`** after each write causes Chrome to hang
- Using Python **text mode** instead of binary mode on Windows (missing `-u` flag)
- Setting `path` to a **relative path** on macOS/Linux (must be absolute)
- Manifest filename **not matching** the `name` field (e.g., `com.foo.bar.json` for name `com.foo.bar`)
- Calling `connectNative` from a **content script** (not allowed — route through service worker)
- Exceeding **1 MB** response size — use compact JSON `separators=(',', ':')`

## Sources

- [Chrome for Developers: Native Messaging](https://developer.chrome.com/docs/extensions/develop/concepts/native-messaging)
- [MDN: Native messaging](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Native_messaging)
- [chrome.runtime API reference](https://developer.chrome.com/docs/extensions/reference/api/runtime)
