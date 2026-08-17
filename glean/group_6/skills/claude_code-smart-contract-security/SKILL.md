---
name: smart-contract-security
description: >-
  Smart-contract security & auditing expert — finding, classifying, and preventing on-chain vulnerabilities on EVM (and Solana). TRIGGER: reviewing/auditing a Solidity/Vyper/Rust contract for exploits; reentrancy (single-function, cross-function, cross-contract, read-only); access-control failures, unprotected init/selfdestruct, tx.origin auth; integer over/underflow (pre-0.8 vs checked vs unchecked), rounding/precision, ERC-4626 inflation; oracle/price manipulation, flash-loan attacks; front-running/MEV-aware contract bugs; unchecked external calls, delegatecall/storage collision, signature replay/malleability, uninitialized proxies (UUPS), DoS/gas-griefing; the audit process (scoping, threat modeling, writing findings), SWC/SCWE/OWASP-SC-Top-10 taxonomy, severity rubrics (Immunefi/Code4rena/Sherlock); security tooling (Slither, Mythril, Foundry invariant testing, Echidna, Medusa, Certora, Halmos, Ethernaut, Damn Vulnerable DeFi); Solana-program pitfalls (missing signer/owner checks, account confusion, arbitrary CPI, Anchor constraints); landmark exploits (The DAO, Parity, bridge hacks). SKIP: Solidity/Vyper LANGUAGE syntax & general dev tooling → ethereum-smart-contract-development; EVM/protocol internals (opcodes, gas mechanics, state trie) → ethereum-protocol-expert; MEV economics/markets (PBS, builders, order-flow auctions) → blockchain-economics; consensus-layer attacks (51%, long-range, selfish mining) → distributed-systems-consensus; Rust/Anchor/Solana LANGUAGE surface (not security) → lang-rust; general web/app/infra appsec (OWASP web Top 10, CSP, SSRF) → security-review/security-reviewer.
version: 1.0.0
updated: 2026-06-16
category: developer
whenToUse:
  - Auditing or security-reviewing a smart contract (Solidity, Vyper, or Solana/Rust) for vulnerabilities
  - Identifying or explaining a contract vulnerability class (reentrancy, access control, oracle manipulation, delegatecall, signature replay, DoS, integer/precision bugs)
  - Choosing or applying a secure coding pattern (checks-effects-interactions, reentrancy guard, pull-over-push, safe external calls)
  - Selecting or running smart-contract security tooling (Slither, Mythril, Foundry invariant tests, Echidna, Medusa, Certora, Halmos) or practice platforms (Ethernaut, Damn Vulnerable DeFi)
  - Structuring a contract audit, classifying finding severity, or writing an audit finding
  - Reviewing a Solana program for signer/owner/account-confusion/arbitrary-CPI pitfalls
  - Learning from landmark exploits (The DAO, Parity, Wormhole, Nomad, Ronin, Poly Network, Mango Markets) for pattern recognition
keywords:
  - smart contract security
  - smart contract audit
  - reentrancy
  - access control vulnerability
  - oracle manipulation
  - flash loan attack
  - delegatecall storage collision
  - reentrancy guard
  - checks-effects-interactions
  - Slither
  - Mythril
  - Foundry invariant testing
  - Echidna
  - Medusa
  - Certora
  - Halmos
  - Ethernaut
  - Damn Vulnerable DeFi
  - SWC registry
  - OWASP smart contract top 10
  - uninitialized proxy
  - UUPS
  - signature replay
  - tx.origin
  - integer overflow
  - ERC-4626 inflation attack
  - Solana program security
  - Anchor constraints
  - arbitrary CPI
  - The DAO hack
  - Parity wallet hack
  - bridge hack
tags:
  - blockchain
  - smart-contracts
  - security
  - auditing
  - ethereum
  - solana
  - defi
  - solidity
---

# Smart-Contract Security & Auditing

Expert reference for **finding, classifying, and preventing vulnerabilities in smart contracts** — the exploit-and-audit layer above contract development. EVM-first (Solidity/Vyper) with a dedicated Solana/Anchor section so the skill is not EVM-only. This skill is for **security review**: identifying bug classes, applying secure patterns, running security tools, structuring an audit, and recognizing exploit patterns from history.

> **Scope & peer routing.** This skill owns *vulnerabilities, exploits, secure patterns, audit methodology, and security tooling*. It defers:
> - **Solidity/Vyper language syntax, contract dev workflow, deployment** → `ethereum-smart-contract-development` (parallel sibling).
> - **EVM/protocol internals** (opcode-level gas, state trie, mempool mechanics) → `ethereum-protocol-expert`.
> - **MEV economics & markets** (PBS, builders/searchers, order-flow auctions, the supply chain) → `blockchain-economics`. Here we cover MEV *only* as it creates contract bugs and the coding defenses.
> - **Consensus-layer attacks** (51%, long-range, nothing-at-stake, selfish mining) → `distributed-systems-consensus`.
> - **Rust/Anchor/Solana language & dev surface** (PDAs, CPI mechanics, Borsh, compute budget as *development* topics) → `lang-rust`. Here we cover the Solana *security pitfalls* and how Anchor constraints mitigate them.
> - **General web/app/infrastructure appsec** (OWASP web Top 10, CSP, SSRF, secrets) → `security-review` / `security-reviewer`. Smart-contract security is a distinct discipline (on-chain, adversarial-by-default, immutable, value-bearing) — do not conflate the two.

## How this skill is organized

The deep material lives in four reference files. Read the one that matches the task:

| Reference | Covers |
| --- | --- |
| `references/vulnerability-classes.md` | The full bug catalog: reentrancy family, access control, delegatecall/storage collision, uninitialized proxies (UUPS), oracle/price manipulation, flash loans, integer over/underflow, rounding/precision (ERC-4626), front-running/MEV-aware bugs, unchecked calls, signature replay & malleability, tx.origin, DoS. Each with mechanism + minimal code sketch + fix. |
| `references/secure-patterns.md` | The defensive patterns: checks-effects-interactions, reentrancy guards (incl. transient-storage), pull-over-push, safe external calls, and why `.transfer`/`.send` are deprecated as a reentrancy defense. |
| `references/audit-methodology.md` | The audit process (scoping → threat modeling → manual review → reporting → fix-review), engagement models (private vs competitive contest vs bug bounty), the taxonomy landscape (SWC legacy, SCWE/SCSVS, EEA EthTrust, OWASP SC Top 10), severity rubrics, and how to write a finding. |
| `references/tooling-and-cases.md` | The tool catalog (Slither, Mythril, Foundry invariant testing, Echidna, Medusa, Certora, Halmos, Ethernaut, Damn Vulnerable DeFi) + the recommended workflow; the Solana-program security section; and landmark exploit case studies. |

## The mental model (read this first)

Smart contracts are **adversarial-by-default, immutable, and value-bearing**: code is public, anyone can call any function in any order with any input, deployed code usually cannot be patched in place, and bugs are directly monetizable. Three consequences drive everything below:

1. **Every external call is a yield point.** When your contract calls out (ETH send, token transfer with a hook, a call to another contract), you hand control to code that may call *back* into you before you finish. This is the root of reentrancy.
2. **Every input and every account is attacker-chosen.** On the EVM, function arguments and call ordering are adversarial. On Solana, the *accounts* passed to an instruction are adversarial and the runtime does **not** verify their relationships — your program must.
3. **Tools find the shallow bugs; humans find the deep ones.** Static analysis and fuzzing are a cheap, high-value first pass, but ~half of real audit findings (business logic, economic/oracle attacks, front-running) require human understanding of the protocol's intent. Tools *augment*, never *replace*, manual review. (Trail of Bits' review of 246 findings: ~49% were "almost impossible to detect with a tool.")[^tob246]

## Vulnerability-class map

The authoritative awareness ranking is the **OWASP Smart Contract Top 10** (the 2025 edition derived from 2024 incidents; a 2026 edition exists).[^owaspsc] The classes, with the reference section that deep-dives each:

| Class | One-line | Detail |
| --- | --- | --- |
| **Access control** | Privileged function reachable by any caller; missing `onlyOwner`/role checks; unprotected init/`selfdestruct`/`mint` | `vulnerability-classes.md` |
| **Reentrancy** | External call hands control to attacker before state is finalized (single-function, cross-function, cross-contract, read-only) | `vulnerability-classes.md` + `secure-patterns.md` |
| **Oracle / price manipulation** | Using a manipulable spot price (DEX reserves, LP-share price) as a trusted feed | `vulnerability-classes.md` |
| **Flash-loan-facilitated** | Atomic uncollateralized capital amplifies an oracle or governance attack within one transaction | `vulnerability-classes.md` |
| **Unchecked external calls** | Low-level `call`/`send` returns `false` on failure instead of reverting — silent failure | `vulnerability-classes.md` |
| **Arithmetic — overflow/underflow** | Pre-0.8 silent wraparound; the re-introduced risk in `unchecked {}` blocks and downcasts | `vulnerability-classes.md` |
| **Arithmetic — rounding/precision** | Integer-division truncation, wrong rounding direction, the ERC-4626 first-depositor inflation attack | `vulnerability-classes.md` |
| **Delegatecall / storage collision** | Callee code runs against caller storage; proxy/implementation layout mismatch; selector clash; arbitrary target | `vulnerability-classes.md` |
| **Uninitialized proxy** | Implementation left uninitialized; attacker seizes it; UUPS self-destruct bricks the proxy | `vulnerability-classes.md` |
| **Signature replay & malleability** | Missing nonce/chainId/domain separator; ECDSA s-malleability; `ecrecover` returning `address(0)` | `vulnerability-classes.md` |
| **tx.origin auth** | `require(tx.origin == owner)` is phishable | `vulnerability-classes.md` |
| **Front-running / MEV-aware** | Transaction-ordering dependence, sandwiching, the `approve()` race; slippage/deadline/commit-reveal defenses | `vulnerability-classes.md` |
| **Denial of service** | Gas griefing, unbounded loops, DoS-by-revert in a push-payment loop | `vulnerability-classes.md` |
| **Solana-specific** | Missing signer/owner checks, account confusion/type-cosplay, arbitrary CPI, PDA bump, `init_if_needed` | `tooling-and-cases.md` |

> **Taxonomy note (verified-as-of 2026-06-16):** The **SWC Registry** (SWC-100…SWC-136) is the historically famous classification but is **unmaintained since ~2020** — useful for cross-referencing legacy findings, not as a current authority. The maintained successors are **OWASP SCS** (the SCSVS Verification Standard + the **SCWE** weakness enumeration, the modern CWE-style successor to SWC) and the **EEA EthTrust Security Levels** (certification). No single taxonomy is universally canonical; cite the OWASP SC Top 10 for awareness and EthTrust/SCSVS for verification checklists.[^swc][^ethtrust] See `audit-methodology.md`.

## Secure-patterns quick reference

The defensive core (full treatment in `secure-patterns.md`):

- **Checks-Effects-Interactions (CEI):** order every function as validate → write all state → external calls *last*. Apply it even to "trusted" callees (a trusted party can still hand control to a third party).
- **Reentrancy guard / mutex:** OpenZeppelin `ReentrancyGuard` (`nonReentrant`); on EIP-1153 chains prefer `ReentrancyGuardTransient` (Solidity ≥0.8.24, Cancun, OZ v5.1+). A single-contract guard does **not** stop *cross-contract* reentrancy.
- **Pull-over-push:** let recipients withdraw rather than pushing funds in a loop — neutralizes both reentrancy and the one-reverting-recipient DoS.
- **Safe external calls:** prefer `(bool ok,) = addr.call{value:v}(""); require(ok);`. **Stop using `.transfer`/`.send` as a reentrancy defense** — their 2300-gas stipend broke after EIP-1884 repriced `SLOAD`; rely on CEI + guards instead.
- **Check every low-level return value.** `call`/`delegatecall`/`staticcall`/`send` return `false` on failure silently.

## Tooling pipeline (one line)

Static analysis (**Slither**, every commit/CI) → fuzzing & invariant testing (**Foundry** `invariant_`, **Echidna**, **Medusa**) → symbolic execution / formal verification for critical properties (**Mythril**, **Halmos**, **Certora**) → **manual review** of what tools can't reach (business logic, economic attacks, front-running). Practice on **Ethernaut** and **Damn Vulnerable DeFi**. Tools augment expert review; they do not replace it. Detail and current status (e.g. MythX is shut down; Echidna is in maintenance with Medusa the recommended path; Certora open-sourced Feb 2025) in `tooling-and-cases.md`.

## When reviewing a contract — a checklist starter

1. **Map the trust boundaries** — which functions are privileged, who can call them, what each external call hands control to.
2. **For every state-changing function with an external call:** is it CEI-ordered? Guarded? Does it check return values?
3. **For every privileged function:** is access control present and correct? Is the upgrade path (`_authorizeUpgrade`) guarded? Are implementations initialized-locked (`_disableInitializers()`)?
4. **For every price/value input:** is the oracle manipulation-resistant, or is it a spot price an attacker can move with a flash loan?
5. **For every arithmetic op:** any `unchecked {}` blocks? Downcasts? Rounding that should favor the protocol? Share-price / first-depositor math?
6. **For every signature check:** nonce + chainId + domain separator (EIP-712)? Low-s enforced? `ecrecover != address(0)`?
7. **For loops and batch payments:** bounded? Pull-over-push? Can one reverting party block everyone?
8. **(Solana)** every account: signer checked? owner checked? type/discriminator checked? every CPI target pinned?
9. **Run the tools, then read the code the tools can't reason about.**

## References (skill-level sources)

These ground the SKILL.md claims above; each reference file carries its own fuller citation list.

[^tob246]: Trail of Bits, "246 Findings From Our Smart Contract Audits: An Executive Summary" — ~49% of findings "almost impossible to detect with a tool"; tools augment, not replace, manual review. https://blog.trailofbits.com/2019/08/08/246-findings-from-our-smart-contract-audits-an-executive-summary/ (blog, verified-as-of 2026-06-16)
[^owaspsc]: OWASP Smart Contract Top 10 (2025, from 2024 incidents; 2026 edition exists). https://scs.owasp.org/sctop10/ (standard, verified-as-of 2026-06-16)
[^swc]: SWC Registry — README states it is no longer actively maintained (no new entries since ~2020); based on EIP-1470/CWE. https://github.com/SmartContractSecurity/SWC-registry (standard, verified-as-of 2026-06-16)
[^ethtrust]: EEA EthTrust Security Levels (v2 published 2023-12-13; v3 in draft); OWASP SCS umbrella houses SCSVS + SCWE. https://entethalliance.org/specs/ethtrust-sl/ and https://scs.owasp.org/ (standard, verified-as-of 2026-06-16)
