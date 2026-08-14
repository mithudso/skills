# blockchain

**Category:** Science, Biology & Medicine
**Platform:** Claude Code
**Original Path:** claude-code/blockchain

## Description
Blockchain & crypto-protocol hub — how chains work beneath the app layer, their economics & on-chain security; routes to spokes. TRIGGER: Bitcoin protocol (UTXO, Script, Taproot, Lightning, halving/21M, mempool/fees, ordinals/Runes); Ethereum protocol (account model, EVM/gas, EIP-1559, Proof of Stake/Gasper, EIP-4844 blobs, L2 rollups); Solana (account model, Sealevel/SVM, programs); tokenomics & cryptoeconomics (emission/burns, staking economics, security budget, fee markets, MEV/PBS, governance); smart-contract security & auditing (reentrancy, oracle/flash-loan, EVM & Solana pitfalls, audit process, Slither/Mythril/Foundry/Echidna/Certora). SKIP: consensus THEORY (CAP/FLP/Paxos/Raft/BFT, PoW-vs-PoS as a class) → distributed-systems-consensus; Rust/Solidity LANGUAGE syntax → lang-rust; crypto-primitive math (ECDSA/Schnorr/BLS/KZG) → blockchain-crypto-primitives; crypto trading/investing → trading-and-investing.

---

# Blockchain & Crypto Protocols (hub)

Front door for **how blockchains work beneath application layer** — base-layer protocol mechanics, economics that secure them, on-chain security. Hub routes to reference spokes. For agreement-theory (CAP/FLP/BFT/Paxos/Raft and PoW-vs-PoS *as class*) neighbor is **`distributed-systems-consensus`**; for crypto-primitive math, **`blockchain-crypto-primitives`**.

**When activated:** match question to row, then **Read listed `references/` file** before answering at depth — table is router only.

## Routing table

| Spoke | Use when | Reference file |
| --- | --- | --- |
| `bitcoin-protocol-expert` | Bitcoin base layer & ecosystem: UTXO model, Script/addresses, SegWit/Taproot, mining/halving/21M, mempool/fees (RBF/CPFP), Lightning, ordinals/Runes, L2/metaprotocols. | `references/bitcoin-protocol-expert.md` |
| `ethereum-protocol-expert` | Ethereum execution/consensus/scaling: account model, EVM/opcodes/gas, EIP-1559 fee market, Proof of Stake (Beacon/Gasper), EIP-4844 blobs, optimistic vs ZK rollups. | `references/ethereum-protocol-expert.md` |
| `blockchain-economics` | Tokenomics & cryptoeconomics: issuance/emission/burns, staking economics & security budget, fee-market design, MEV & PBS supply chain, governance economics. | `references/blockchain-economics.md` |
| `smart-contract-security` | Finding/preventing on-chain vulns: reentrancy, access control, oracle/flash-loan attacks, EVM & Solana-program pitfalls, audit process, severity rubrics, Slither/Mythril/Foundry/Echidna/Certora. | `references/smart-contract-security.md` |
| `solana-account-model-and-state` | Solana (account model & state): accounts as typed data, rent, program-derived addresses, ownership, account-centric execution model. | `references/solana-account-model-and-state.md` |
| `solana-sealevel-and-svm` | Solana (Sealevel/SVM): parallel runtime and Solana VM — parallel transaction execution, programming model, how differs from EVM. | `references/solana-sealevel-and-svm.md` |

## Cross-cutting notes

- **Consensus theory is neighbor, not here.** Safety/liveness proofs, CAP/FLP, Paxos/Raft/PBFT, PoW-vs-PoS *as class* live in `distributed-systems-consensus`. Hub covers each chain's *concrete* consensus (Nakamoto PoW, Gasper, Tower BFT) as part of protocol.
- **Language vs protocol.** Solidity/Vyper/Foundry and Rust/Anchor *language* surface → `lang-rust` and Ethereum protocol reference's dev section; hub owns protocol mechanics and security, not language syntax.
- **Not investing.** Token *valuation*, trading, portfolio decisions → `trading-and-investing`. Hub is mechanism & economics, not financial advice.

<!-- cross-hub-map -->
## Cross-hub map — where every blockchain topic lives

All blockchain-family topics live as reference files under single hub. If question's deep material **not** in routing table above, out of scope for family — see SKIP / Cross-cutting notes for appropriate neighbor hub.

| Hub | Owns | Reference files |
| --- | --- | --- |
| `blockchain` | Blockchain & Crypto Protocols | `references/bitcoin-protocol-expert.md`, `references/ethereum-protocol-expert.md`, `references/blockchain-economics.md`, `references/smart-contract-security.md`, `references/solana-account-model-and-state.md`, `references/solana-sealevel-and-svm.md` |