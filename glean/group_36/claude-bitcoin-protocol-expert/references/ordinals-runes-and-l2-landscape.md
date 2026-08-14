<!-- Provenance: reference under the standalone `bitcoin-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. SURVEY DEPTH — emerging/experimental/contested domain. All volatile claims stamped verified-as-of 2026-06-16. Peer deferral: crypto math (Schnorr/adaptor/zk) → blockchain-crypto-primitives; consensus theory → distributed-systems-consensus; token valuation → blockchain-economics. -->

# Ordinals, Runes & the Bitcoin L2 / Restaking Landscape

**SURVEY DEPTH** — a "know it exists and roughly how it works" map of Bitcoin's metaprotocols and emerging L2/restaking frontier.

> **Read this as emerging, experimental, and heavily contested.** Adoption numbers, mainnet/testnet status, and project maturity move fast. **All volatile claims are stamped verified-as-of 2026-06-16.** Base-layer mechanics (Taproot/witness/OP_RETURN) are assumed known — see the other references. Lightning (the payments L2) is covered in `references/lightning-network.md`. Be conservative: this domain is hype-prone.

## Contents
- [Ordinals & inscriptions](#ordinals--inscriptions)
- [BRC-20 & Runes](#brc-20--runes)
- [BitVM](#bitvm)
- [Bitcoin L2 / staking landscape](#bitcoin-l2--staking-landscape)
- [The covenant constraint](#the-covenant-constraint)
- [Babylon](#babylon)
- [The honest meta-point](#the-honest-meta-point)
- [References](#references)

## Ordinals & inscriptions

**Ordinal Theory** numbers every satoshi by **mining order** (genesis-block first sat = ordinal 0) and tracks it across transactions by a **FIFO** input→output rule. It needs **no sidechain, no token, no consensus change** (fact).[^1][^2][^3]

**Inscriptions** attach arbitrary content (images, text) to a sat, stored **entirely on-chain in the Taproot script-path witness** — which has few content restrictions and **gets the witness discount** (fact).[^1][^4] Mechanically: content is wrapped in an **"envelope"** — an `OP_FALSE OP_IF … OP_ENDIF` block (a no-op, so script semantics are unchanged) — and revealed by a **two-phase commit/reveal** (a commit tx creates a Taproot output committing to the script; a reveal tx spends it, exposing the content). The framing is "Bitcoin-native digital artifacts / NFTs."[^1] The **`ord` software** (github.com/ordinals/ord, by **Casey Rodarmor**) is the reference index + explorer + CLI wallet — explicitly "experimental software with no warranty," requiring `bitcoind -txindex`; it warns Bitcoin Core is **not inscription-aware** (using `bitcoin-cli` on an `ord` wallet can lose inscriptions).[^5]

**Inscriptions vs OP_RETURN** (fact): OP_RETURN is Bitcoin's sanctioned small data-carrier (historically ~80 bytes) in a provably-unspendable output; inscriptions hide larger data in the **witness** for the discount. A nuance: large inscribed images sit in a **non-executed** witness part (no signatures to verify), so per unit of weight they can be *cheaper* to validate than ordinary txs and don't bloat the UTXO set.[^6]

**The "spam vs legitimate use" controversy** (fact — ongoing). Critics (notably **Luke Dashjr**) call Ordinals a "spam attack"/"bug" to be filtered (his 2023–24 proposal was closed; CVE-2023-50428 framed witness data-carrying as a `datacarriersize` bypass).[^7][^8] It reignited in 2025 as the "spam wars": Peter Todd's PR to **remove the OP_RETURN 80-byte limit** was merged into Bitcoin Core (~Oct 2025), on the logic that filters are futile when Taproot already carries arbitrary data — drawing fierce opposition. The deeper fault line: Bitcoin as **"peer-to-peer cash only"** vs **fee-market-decides**.[^9][^10][^11]

## BRC-20 & Runes

**BRC-20** (pseudonymous dev **domo, March 2023**) — an **experimental inscription-based fungible-token standard** embedding **JSON in inscriptions** with `deploy`/`mint`/`transfer` ops (fact). Bitcoin has **no VM**, so there is **no on-chain enforcement**: balances exist only when an **off-chain indexer** reads inscriptions and applies the ruleset — a full node alone is insufficient; wallets/markets must run or trust an indexer.[^12][^13][^14] Criticisms: **indexer dependence** (outages/bugs/divergence → corrupt balances, centralization) and **UTXO-set bloat** (BRC-20 txs resemble ordinary txs, gaining little from the witness discount but **proliferating UTXOs** — BitMEX Research ties much of the UTXO-set growth from ~84M to ~169M outputs, Dec 2022–Sep 2025, largely to BRC-20). It hit >$1B market cap within ~2 months. (fact; the BitMEX attribution is one research house — qualify.)[^6][^14][^15]

**Runes** (**Casey Rodarmor**, launched at the **2024 halving — block 840,000, ~Apr 20, 2024**) — a **UTXO-based** fungible-token protocol (fact).[^16][^17][^18] Design (from Rodarmor's own blog): a simple, ~2,000-line, **OP_RETURN-based** protocol with a **minimal on-chain footprint, no off-chain data, no native token**, positioned as **harm reduction vs BRC-20's UTXO proliferation**. Balances **live in UTXOs**; a "**runestone**" is an output whose scriptPubKey is `OP_RETURN OP_13` carrying **edicts** (output #, rune ID, amount) for **etching/minting/transferring**. Because OP_RETURN outputs are provably unspendable, they **don't congest the UTXO set**. Rodarmor candidly called it "built for degens and memecoins." Its launch spiked halving-day fees to records (~$20M in fees in the first ~18 post-halving blocks).[^16][^17][^18][^19]

**Current state (verified-as-of 2026-06-16):** Runes has generally **dominated BRC-20/Ordinals in on-chain tx share** since launch (cumulative ~2,500 BTC in associated fees reported 2024–25). Precise mid-2026 figures and any decline are **not well-attested** — treat current adoption magnitude as **tentative**.[^20] (Protocol design & launch = fact; current dominance = moderate; precise 2026 metrics = low/flag.)

## BitVM

**BitVM** (**Robin Linus, late 2023**) expresses **Turing-complete / arbitrary computation verifiable on Bitcoin with NO consensus change** (fact). The inversion: computations are **verified, not executed** — a prover claims `f(inputs)=output`; if false, a verifier runs a **succinct fraud proof** and punishes the prover, using only **hashlocks, timelocks, large Taproot trees**, and pre-signed **challenge-response** transactions. If both cooperate, complex computation happens off-chain with minimal on-chain footprint; on-chain action is needed **only in a dispute**. The original design was **two-party only**.[^21][^22][^23] **BitVM2** generalizes to **permissionless verification**: a one-time setup with a **1-of-n honesty assumption**, but at runtime **anyone** can challenge an invalid assertion (~2 rounds / 3 on-chain txs), splitting SNARK verification into Bitcoin-tx sub-programs. Its headline use is the **BitVM Bridge**, lowering deposit-safety trust from honest-majority (t-of-n) to **1-of-n at setup**, with liveness needing only one rational operator. Developed under **ZeroSync / the BitVM Alliance**.[^24][^25][^26]

**Honest "not production-proven" framing** (fact — primary BitVM docs admit the limitation): a production-level BitVM2 implementation exists and a full challenge verification was **executed on Bitcoin mainnet** (Linus et al., ePrint 2025/1158), and Citrea's "Clementine" bridge deployed on **testnet** (Sept 2024) — so it has moved from whitepaper to working demos. **But** it is **trust-minimized, not trustless**, and not mainstream-ready: a massive on-chain footprint, extreme complexity, **weeks-long dispute resolution**, and **operator liveness/capital requirements** (the "bridge" is really a **reimbursement system** — operators must **front their own BTC**, liquidity scaling linearly with TVL). **"ZK rollups on Bitcoin do not yet work"** — only optimistic rollups via BitVM are possible today; and **if all operators go offline, funds freeze** (no working escape hatch under Bitcoin's UTXO model without future covenant opcodes), enabling a ransom attack.[^24][^26][^27][^28][^29]

## Bitcoin L2 / staking landscape

**There is no agreed definition of "Bitcoin L2."** A common test: (a) security **derived from Bitcoin** and (b) **exit to L1** without fully trusting a third party; the *strict* version (unilateral exit enforced by Bitcoin consensus) is met by essentially **only rollups — arguably none fully today** (fact).[^30][^31] Rough taxonomy:
- **Sidechains** — own consensus, two-way peg, **do not inherit Bitcoin security**. **Liquid** (Blockstream, **federated**, ~65 functionaries, live 2018); **Rootstock/RSK** (**merge-mined**, EVM-compatible, PowPeg federation); **Stacks** (Proof-of-Transfer).[^30][^31][^32]
- **State chains / channels & hybrids** — statechains, **Ark**, **Spark**; **RGB** (client-side validation).[^31][^33]
- **Rollups (attempts)** — genuinely Bitcoin-secured rollups are **nascent**: by credible accounts only 3–4 teams build true Bitcoin rollups and **none are live in production**; most "Bitcoin rollups" are bridged EVM-stack forks.[^29][^34]
- **Drivechains (BIP300/301)** — a **proposal, NOT deployed**. **BIP-300 (Hashrate Escrow)** + **BIP-301 (Blind Merged Mining)** by **Paul Sztorc** would let anyone launch miner-secured sidechains with a ~3–6-month miner-voted withdrawal peg. **Both remain "Draft," not activated** (no soft-fork consensus); contemplated via a CUSF "enforcer" client. Long-running, polarizing (miner-centralization concerns); **not on a near-term roadmap as of 2026**.[^35][^36][^37][^38]

## The covenant constraint

**The gating constraint: Bitcoin has no native covenant opcode.** Trust-minimized L2 peg designs largely wait on new Script capabilities; the debate is live and unresolved (verified-as-of 2026-06-16) (fact for proposal status; "what activates first" is analyst speculation — flag):[^39][^40][^41][^42]
- **OP_CTV / BIP-119 (CHECKTEMPLATEVERIFY)** — Rubin (2020); commits a UTXO to a **fixed tx template**; narrow, non-recursive (vaults, congestion control); tested on signet ~2 years.
- **OP_CAT / BIP-347** — Heilman/Sabouri; stack concatenation enabling (via Schnorr tricks) **introspection and recursive covenants**; BIP merged, Core PR open, on signet. More powerful, higher risk-surface.
- **APO / SIGHASH_ANYPREVOUT / BIP-118** — Decker/Towns; a sighash flag mainly for **LN-Symmetry/eltoo**; Draft.
- Analysts guess **CTV (+CSFS)** activates first and **OP_CAT** stays contested — speculation, not consensus.

*(The cryptographic mechanics → `blockchain-crypto-primitives`.)*

## Babylon

**Babylon** is a **Bitcoin staking / timestamping protocol** letting **native BTC be staked without bridging or wrapping** to secure PoS systems (fact).[^43][^44][^45] BTC is locked in a **self-custodial Bitcoin Script** (Taproot output with timelock-staking, on-demand-unbonding, and slashing paths); the staker delegates voting power to a **finality provider**. Slashing uses **Extractable One-Time Signatures (EOTS)** — a finality provider that **double-signs** exposes its EOTS key, letting **pre-signed slashing transactions** burn a portion of stake; a **covenant committee** (threshold multisig) co-signs unbonding/slashing because Bitcoin lacks native covenants. Co-founded by **David Tse**; raised $70M led by Paradigm.[^44][^46]

**Current state (verified-as-of 2026-06-16):** **Phase-1 mainnet** launched **Aug 22, 2024** — BTC locking only, **no slashing**. **Phase-2 / "Babylon Genesis"** (a Cosmos-SDK L1, token **BABY**) launched **Apr 10, 2025**, **activating slashing**, permissionless from Apr 24, 2025; Genesis is billed as the first "Bitcoin Supercharged Network." Scale: peak ~**57,000 BTC** staked (~$6.1B; ~0.3% of BTC), TVL **volatile** — Messari ~45,600 BTC ($4.9B) Q2 2025 (−12.6% QoQ); DefiLlama ~$3.2–4.8B through 2026. (Mechanics & timeline = fact; live TVL = moderate/volatile.)[^46][^47][^48][^49][^50][^51]

**Babylon risks** (fact — independently reported; dollar figures volatile): (1) **EOTS native-BTC slashing is a brand-new primitive with no battle-tested precedent** — an honest bug could **irreversibly burn** staked BTC; (2) the **covenant committee is a trusted multisig** — compromise/censorship could **freeze unbonding** (a temporary centralization vector meant to be removed *if* Bitcoin gets covenant opcodes — not guaranteed); (3) **systemic/cascade risk** from a large LRT/restaking stack built on Babylon; (4) real bugs found (a BLS vote-extension consensus-bypass class patched in v4.2.0; OpenZeppelin-disclosed delegation/jailing issues); (5) BABY emission-driven yields. Independent analysis stresses it is **"trust-minimized, not trustless."**[^45][^52][^53][^54]

## The honest meta-point

Across the stack, the consistent, independently-corroborated takeaway: **most of this is early, experimental, and contested, and much "Bitcoin L2" branding is marketing** (fact). Analysts counted **70+** self-described "Bitcoin L2s" in ~9 months, many straining any definition. The recurring failure mode: **the bridge is the weak point** — the majority of Bitcoin-L2 TVL is secured by **5–10 federated multisig signers**, reintroducing the trusted third party Bitcoin was built to remove (bridge exploits have caused >$2B in losses crypto-wide). **"There is no such thing as a trustless L2"** — every L2 doing something Bitcoin can't adds a new category of trust. BitVM and covenant opcodes are the credible *technical* paths to shrinking that trust, but the former is complex/not-yet-mainstream and the latter is stuck in an unresolved activation debate.[^29][^31][^41][^55][^56][^57][^59]

## Child sub-concepts (future research)
The OP_RETURN "spam wars" & datacarrier policy (Knots vs Core, "what is Bitcoin for"); BitVM bridge implementations in depth (Citrea/Clementine, operator economics, BitVM3); the covenant activation debate as governance (LNHANCE bundle); BTC restaking / LRT stack & systemic risk; client-side-validation protocols (RGB, Taproot Assets); Stacks/sBTC & Proof-of-Transfer.

## References

[^1]: https://docs.ordinals.com/inscriptions.html — docs (Ordinal Theory Handbook): inscription mechanics, witness storage/discount, commit/reveal, envelopes, numbering.
[^2]: https://docs.ordinals.com/ — docs: Ordinal Theory intro; no sidechain/token; authors.
[^3]: https://docs.ordinals.com/overview.html — docs: mine-order numbering, FIFO transfer, ordinal ranges.
[^4]: https://github.com/ordinals/ord/blob/master/docs/src/inscriptions.md — docs (primary repo): inscription mechanics, "NFTs," witness discount.
[^5]: https://github.com/ordinals/ord — docs (primary repo): `ord` = experimental index/explorer/wallet; requires `bitcoind -txindex`; Core not inscription-aware.
[^6]: https://www.bitmex.com/blog/ordinals-impact-on-node-runners — research blog: inscriptions in non-executed witness (cheap, no UTXO bloat) vs BRC-20 UTXO-set growth (84M→169M).
[^7]: https://github.com/bitcoin/bitcoin/issues/29187 — forum (Luke Dashjr): CVE-2023-50428; OP_FALSE OP_IF witness data framed as spam.
[^8]: https://www.coindesk.com/tech/2024/01/09/bitcoin-developers-proposal-to-stop-spam-nfts-gets-shut-down — news: Dashjr inscription-filter proposal closed.
[^9]: https://www.coindesk.com/tech/2025/04/30/bitcoin-debate-on-looser-data-limits — news: Peter Todd OP_RETURN-limit-removal proposal; how inscriptions exploit Taproot.
[^10]: https://protos.com/bitcoin-devs-continue-fight-over-arbitrary-data-storage/ — news: PR #32359 "spam wars," mixed dev reactions.
[^11]: https://bitcoinmagazine.com/technical/op_return-limits-bitcoins-battle-over-arbitrary-data — blog: "Filterors" position; Core merged removing the limit (~Oct 2025).
[^12]: https://chain.link/education-hub/brc-20-token — docs/edu: BRC-20 by domo Mar 2023; JSON deploy/mint/transfer; off-chain indexer required.
[^13]: https://ar5iv.labs.arxiv.org/html/2310.10652 — paper ("BRC-20: Hope or Hype"): inscription/JSON mechanics; off-chain-indexer reconciliation; domo "fun experimental standard."
[^14]: https://daic.capital/blog/brc20-bitcoin-token-standard — blog: no VM; indexer dependence and failure modes.
[^15]: https://research.nansen.ai/articles/a-primer-on-bitcoin-s-wild-west-ordinals-and-brc-20 — research blog: >$1B mcap in ~2 months; off-chain indexing centralization critique.
[^16]: https://rodarmor.com/blog/runes/ — blog (primary, Casey Rodarmor): Runes spec/intent — UTXO-held balances, OP_RETURN tag, no off-chain data, no native token, harm reduction vs BRC-20.
[^17]: https://sovryn.com/all-things-sovryn/bitcoin-runes-tokens — blog: launched 4th halving Apr 2024; runestone/edicts/etching; OP_RETURN unspendable → no UTXO congestion.
[^18]: https://www.coindesk.com/tech/2024/04/17/runes-casey-rodarmors-protocol-for-shtcoins — news: Rodarmor "built for degens and memecoins," ~2,000-line OP_RETURN protocol.
[^19]: https://insights.blockbase.co/what-is-runes-protocol/ — blog: block 840,000 / Apr 20 2024; record halving-day fees (~$20M in first ~18 blocks).
[^20]: https://decrypt.co/ (aggregated news digest) — news (aggregated): Runes leading BRC-20/Ordinals in tx share since Apr 2024; 2026-specific metrics weakly attested.
[^21]: https://bitvm.org/bitvm.pdf — whitepaper (primary, Robin Linus): BitVM — verify-not-execute, fraud proofs, hashlocks/timelocks/Taproot trees, 2-party.
[^22]: https://bitvm.org/ — docs (primary): project overview; "anyone can fraud-proof"; ZeroSync/BitVM Alliance.
[^23]: https://eprint.iacr.org/2024/1995 — paper (peer-reviewed): "BitVM: Quasi-Turing Complete Computation"; 3 txs optimistic; no consensus change.
[^24]: https://bitvm.org/bitvm2.html — docs (primary): BitVM2 permissionless verification; 1-of-n setup; honest-operator limitation → funds unspendable / ransom attack.
[^25]: https://bitvm.org/bitvm_bridge.pdf — whitepaper (primary): BitVM2 bridge; trust t-of-n → 1-of-n; liveness w/ 1 rational operator.
[^26]: https://eprint.iacr.org/2025/1158 — paper (primary): production-level BitVM2; full challenge verification executed on Bitcoin mainnet.
[^27]: https://blockspace.media/insight/why-you-might-want-to-pay-attention-to-bitvm-again/ — research blog: Citrea "Clementine" testnet (Sep 2024); weeks-long disputes; operators front BTC ("reimbursement system"); trust-minimized not trustless.
[^28]: https://onchainbitcoin.substack.com/p/is-bitvm-going-to-enable-bitcoin — blog: operator capital fronting, liveness/collusion risk; "niche, not mainstream DeFi."
[^29]: https://bitvm-acc.org/assets/BitVM_Status_Report_Feb_2025.pdf — status report: "ZK rollups on Bitcoin do not yet work"; no working escape hatch; all-operators-offline → funds frozen.
[^30]: https://www.dwf-labs.com/research/402-bitcoin-layer-2s-classification-an-update — research blog: sidechain taxonomy; Liquid federated vs Rootstock merge-mined; sidechains don't fully inherit Bitcoin security.
[^31]: https://www.spark.money/research/bitcoin-layer-2-comparison — research blog: no agreed L2 definition; strict vs broad; Liquid functionaries, Rootstock PowPeg; hybrids.
[^32]: https://www.nervos.org/knowledge-base/ultimate_guide_bitcoin_layer_two_part_1 — blog: strict ETH definition (only rollups); Rootstock merge-mined; Stacks PoX.
[^33]: https://chain.link/education-hub/bitcoin-layer-2 — docs/edu: state channels / sidechains / rollups.
[^34]: https://januszgrze.substack.com/p/validiums-volitions-and-superstitions — blog: only 3–4 teams building true Bitcoin rollups; none live; most are bridged EVM forks.
[^35]: https://github.com/bitcoin/bips/blob/master/bip-0300.mediawiki — BIP-300 (primary): Hashrate Escrow, Status Draft, miner-voted withdrawals.
[^36]: https://github.com/bitcoin/bips/blob/master/bip-0301.mediawiki — BIP-301 (primary): Blind Merged Mining, Status Draft.
[^37]: https://www.spark.money/research/bitcoin-drivechain-bip-300 — research blog: drivechain framework; "remains in Draft, not activated" as of 2026.
[^38]: https://www.learnbitcoin.com/glossary/bip-300-drivechains — edu: not activated; "not on near-term roadmap as of 2026."
[^39]: https://www.spark.money/tools/bitcoin-covenant-proposals-compared — research blog: CTV/APO/OP_CAT comparison; recursive-covenant debate.
[^40]: https://www.galaxy.com/insights/research/bitcoins-next-major-upgrade-op-cat-and-op-ctv — research blog: CTV (BIP119) covenant mechanics; uncertain activation timeline.
[^41]: https://www.spark.money/research/bitcoin-op-cat-covenant-debate — research blog: proposal/BIP/status table; "CTV+CSFS likely first, OP_CAT contested."
[^42]: https://bitcoinops.org/en/topics/op_checktemplateverify/ — Optech: CTV field commitments; OP_CAT+CSFS can emulate CTV; signet usage.
[^43]: https://docs.babylonlabs.io/guides/overview/bitcoin_staking/ — docs (primary): staking on Bitcoin UTXO/Script, self-custody, no wrap/bridge, EOTS + covenant-committee slashing.
[^44]: https://babylonlabs.io/blog/babylon-bitcoin-staking-mainnet-launch-phase-1 — blog (primary): Phase-1 mechanics; finality-provider delegation; no slashing in Phase-1.
[^45]: https://www.openzeppelin.com/news/state-changes-at-the-boundary-lessons-from-security-research-on-babylon — security research: Taproot script-tree paths; covenant-committee sigs; pre-signed slashing; four on-chain bugs found.
[^46]: https://www.prweb.com/releases/trustless-bitcoin-staking-protocol-babylon-launches-on-mainnet.html — press release (primary): Phase-1 mainnet Aug 22 2024; $70M raise led by Paradigm.
[^47]: https://github.com/babylonlabs-io/babylon/blob/main/docs/register-bitcoin-stake.md — docs (primary repo): pre-signed slashing for double-signing; burn-address portion + timelock.
[^48]: https://babylonlabs.io/blog/phase-2-launch-round-up — blog (primary): Genesis launched Apr 10 2025; permissionless Apr 24; >50,000 BTC.
[^49]: https://docs.babylonlabs.io/guides/overview/babylon_genesis/ — docs (primary): Genesis = first BSN; BABY token; dual quorum; >57,000 BTC (~$6.1B, ~0.3% of BTC).
[^50]: https://messari.io/report/state-of-babylon-q2-2025 — research: Genesis Apr 10 2025 activated slashing; TVL 45,600 BTC ($4.9B), −12.6% QoQ.
[^51]: https://defillama.com/protocol/babylon-protocol — data: TVL ~$3.26B, BABY mcap ~$58M, as of 2026.
[^52]: https://hindenrank.com/protocol/babylon — risk analysis: novel EOTS slashing (honest bug can burn BTC), trusted covenant multisig (freeze risk), LRT cascade, BABY emission yields.
[^53]: https://onchainbitcoin.substack.com/p/how-trustless-are-babylons-bitcoin — blog: "trust-minimized, not trustless."
[^54]: https://hindenrank.com/blog/how-does-babylon-protocol-work — risk analysis: BLS vote-extension vuln patched v4.2.0; whale-concentration unstaking event.
[^55]: https://blockspace.media/insight/most-bitcoin-layer-2s-are-marketing-scams/ — research blog: 70+ "Bitcoin L2s" in ~9 months, many falsely marketed.
[^56]: https://chainscorelabs.com/blog/where-bitcoin-layer-2s-can-fail — blog: "bridge is the weakest link"; >90% of L2 TVL on 5–10 federated signers; >$2B bridge hacks.
[^57]: https://cointelegraph-magazine.com/bitcoin-layer2-sidechains-not-really-bitcoin-l2s/ — blog: most "L2s" are sidechains w/ multisig bridges; "L2" = marketing.
[^59]: https://www.somethinginteresting.news/p/there-are-only-three-kinds-of-bitcoin — blog: "There is no such thing as a trustless L2."
