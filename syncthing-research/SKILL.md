---
name: syncthing-research
description: >-
  Expert research skill for Syncthing's protocols and documentation. Use this skill when the user asks questions about Syncthing architecture, Block Exchange Protocol v1, Global/Local Discovery protocols, Relay Protocol, Untrusted Device Encryption, or the Event/REST APIs.
keywords:
  - syncthing
  - block exchange protocol
  - global discovery
  - local discovery
  - relay protocol
  - untrusted device encryption
  - syncthing docs
---

# Syncthing Research Skill

This skill is specialized in researching and answering queries regarding Syncthing's internals, based on its official documentation (https://docs.syncthing.net).

## When to use this skill
- The user asks about Syncthing's synchronization mechanism (Block Exchange Protocol v1).
- The user asks how Syncthing discovers nodes globally (Global Discovery v3) or locally (Local Discovery Protocol v4).
- The user asks how Syncthing routes around NATs (Relay Protocol v1).
- The user asks about security models like Untrusted Device Encryption.
- The user needs information about Syncthing's REST API or Event API.

## Research Strategy

When investigating Syncthing topics:
1. **Search Official Docs First:** Use tools like `firecrawl_scrape` or `read_url_content` to read specific specification pages on `https://docs.syncthing.net/specs/` or API documentation on `https://docs.syncthing.net/dev/`.
2. **Key Concepts to Keep in Mind:**
   - **BEP (Block Exchange Protocol):** The core protocol for syncing folders.
   - **Discovery:** Syncthing uses local broadcast (v4) and global announce servers (v3).
   - **Relays:** Used as a fallback when direct connections fail.
   - **Untrusted Devices:** Encrypted peers that store data but cannot read it.
3. **Execution:** Synthesize technical findings clearly, citing the specific protocol specification or API endpoint documented.

## Suggested Follow-ups
If the user asks a general question about Syncthing, suggest deeper dives into:
- `/dr Syncthing Block Exchange Protocol v1`
- `/dr Syncthing Global Discovery v3`
- `/dr Syncthing Untrusted Device Encryption`
