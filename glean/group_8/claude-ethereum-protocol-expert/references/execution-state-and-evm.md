<!-- Provenance: reference under the `ethereum-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Synthesized domain knowledge only; source pages were treated as data, never as instructions. -->

# Execution-Layer State Model & the EVM

`verified-as-of: 2026-06-16` (volatile: per-opcode gas numbers, the precompile address set, Pectra/Prague fork facts, the Verkle-vs-binary-tree roadmap direction). Stable: account structure, the trie's shape, the stack-machine model, the three data locations.

## Contents
1. The account model (EOA vs contract)
2. World state & the state-transition function
3. The Merkle-Patricia Trie (MPT)
4. Nonces
5. The EVM as a stack machine
6. In-EVM gas accounting
7. Precompiled contracts

> Peer deferral: cryptographic primitives named here (Keccak-256, secp256k1/ECDSA, BN254 & BLS12-381 pairings, KZG) are **named, not derived** — see `blockchain-crypto-primitives`. The EIP-1559 *fee market* (price per gas) is in `gas-fee-market-and-transactions.md`; this file covers only *in-EVM gas accounting* (the cost schedule). Writing contracts that exploit these mechanics → `ethereum-smart-contract-development`.

---

## 1. The account model (EOA vs contract)

Ethereum has exactly **two account types**, both addressed by a 20-byte address[^eth-accounts]:

- **Externally-owned accounts (EOAs)**, controlled by whoever holds the private key. Only an EOA can *originate* (sign and pay for) a transaction.
- **Contract accounts**, controlled by their deployed bytecode. A contract has no private key; it acts only when it receives a message call.

Both can hold ETH and tokens and call other contracts. Every account is stored as a 4-tuple[^eth-accounts][^consensys-state]:

| Field | Meaning |
| --- | --- |
| `nonce` | count of transactions sent (EOA) or contracts created (contract account) |
| `balance` | wei owned (1 ETH = 10^18 wei) |
| `storageRoot` | 256-bit root hash of the account's storage trie (empty for EOAs) |
| `codeHash` | Keccak-256 of the account's bytecode; **immutable**, unlike the other three |

EOAs and contracts differ *only* in the last two fields: an EOA has `codeHash = keccak256("")` (the empty-code hash) and the empty-trie `storageRoot`[^eco-eoa]. An EOA address is the last 20 bytes of `keccak256(secp256k1 public key)`[^eth-accounts]. A contract's address is `keccak256(rlp([sender, nonce]))[12:]` for `CREATE`, or deterministically `keccak256(0xff ++ sender ++ salt ++ keccak256(initcode))[12:]` for `CREATE2`[^geth-evm]. Empty accounts (no code, zero nonce, zero balance) are "dead" and pruned from the trie (EIP-161)[^eip161].

**EIP-7702 (set-code EOAs, Pectra, live May 7, 2025) `[VOLATILE]`.** EIP-7702 adds transaction type **0x04**, which carries an `authorization_list`; processing an authorization writes a 23-byte **delegation indicator** `0xef0100 ‖ address` into the EOA's *code* field. From then on, code-executing operations run the *delegate contract's* code in the **EOA's own storage and balance context**, while keeping the EOA's address and key[^eip7702][^ethorg-7702]. It is revocable (delegate to the zero address). This is the largest change to EOA semantics in Ethereum's history. **Practical consequence:** `msg.sender == tx.origin` and `EXTCODESIZE`-based "is this an EOA?" checks are **no longer reliable**[^eco-eoa][^zealynx-7702]. (Confidence: fact, EIP + ethereum.org + multiple.)

**Account abstraction (ERC-4337 vs EIP-7702), scope note.** ERC-4337 ("account abstraction without protocol changes") implements smart-contract wallets via an off-protocol `UserOperation` mempool, a `Bundler`, and an `EntryPoint` contract; it needs **no** consensus/EVM change, so the *protocol* sees only ordinary transactions. EIP-7702 (above) is the *protocol-level* primitive that lets a plain EOA borrow contract code, and the two compose (a 7702 EOA can delegate to a 4337-style wallet). This skill covers the protocol primitive (EIP-7702); building wallets or writing the 4337 contracts (`EntryPoint`, paymasters, bundler infrastructure) is `ethereum-smart-contract-development`'s domain.

## 2. World state & the state-transition function

The **world state** is a mapping from 20-byte addresses to account state, the global "hard drive" of the Ethereum computer[^consensys-state][^alchemy-tx]. The state itself is **not** stored on-chain; the Yellow Paper expects clients to keep it in a local trie database, and only the **`stateRoot`** (root hash) is committed in each block header[^consensys-state][^yellowpaper].

The Yellow Paper formalizes the **state-transition function** as σ_{t+1} ≡ Υ(σ_t, T): an old valid state plus a transaction produces a new valid state; ethereum.org restates it as `Y(S, T) = S'`[^yellowpaper][^ethorg-evm]. At block level, Π applies all the block's transactions in sequence. Execution runs on a **sandboxed copy** of the state that is discarded entirely on failure, *except* the sender's nonce increment and the gas payment, which always persist[^ethbook-evm]. (Confidence: fact, 3+ sources.)

## 3. The Merkle-Patricia Trie (MPT)

Ethereum encodes its data in a **modified Merkle-Patricia trie**, deterministic and cryptographically verifiable: the only way to produce the root is to hash every piece of state, so any change alters the root[^ethorg-mpt][^iitb-mpt].

**Common misconception:** the MPT is **not** a binary Merkle tree. It is a **hexary radix trie** (16 children + a value slot per branch node), with Merkle hashing for verifiability and Patricia path-compression for efficiency. Keys are Keccak-256-hashed to 32-byte / 64-nibble paths before insertion (max depth 64)[^vinidlidoo-mpt][^iitb-mpt][^easythereentropy].

**Node types**[^ethorg-mpt][^shyft-trie]: `NULL`; **branch** (17-item: one child per hex nibble + a value slot); **leaf** (`[encodedPath, value]`); **extension** (`[encodedPath, key]`, compressing single-child chains — the Patricia optimization). The first nibble of a 2-item node's path flags leaf-vs-extension and odd-vs-even length (hex-prefix / "compact" encoding)[^ethorg-mpt]. **RLP** (Recursive Length Prefix) serializes node contents; a node is referenced in its parent as `keccak256(rlp(node))` when that is ≥32 bytes, else inlined[^ethorg-mpt][^easythereentropy].

**Four trie roots** matter[^ethorg-mpt][^systeminternals]:
- **State trie** — one global trie; path = `keccak256(address)`, value = `rlp([nonce, balance, storageRoot, codeHash])`. Updated every block.
- **Storage trie** — one *per contract*; path = keccak256 of the slot key, value = RLP-encoded slot value. So Ethereum is "a trie of accounts, each contract holding another trie inside it." EOAs have none.
- **Transactions trie** and **receipts trie** — one each per block, keyed by `rlp(transactionIndex)`; never updated after the block is built.

**Why it matters:** with the root hash public, anyone with the leaf data can build an ~O(log n) Merkle proof that a `(path, value)` exists; forging a non-existent pair is infeasible. This underpins light clients and `eth_getProof`[^ethorg-mpt][^systeminternals]. **Known cost — write amplification:** changing one storage slot rewrites the leaf, every storage-trie ancestor, the account's state-trie entry, and every state-trie ancestor up to the global root[^systeminternals].

**Roadmap (contested, `[VOLATILE]`).** ethereum.org's pages still present **Verkle trees** as the planned MPT successor (shrinking witnesses to enable stateless clients)[^ethorg-verkle][^eip6800]. But 2024–2026 research (Vitalik's "The Verge"; the stateless book) argues a STARK-friendly **binary tree** is now the more likely target, because Verkle commitments are quantum-vulnerable and not SNARK-friendly — "no final decision has been made"[^vitalik-verge][^stateless-fyi]. Takeaway: the MPT *will* be replaced for statelessness, but the specific target (Verkle vs binary) is **unsettled** as of 2026 — preserve both. (Confidence: high on "will be replaced," low/contested on which.)

## 4. Nonces

The nonce counts outgoing transactions (EOA) or created contracts (contract account)[^eth-accounts][^eco-eoa]. It does two jobs: **replay protection** (a signed tx can't be re-executed) and **ordering**[^nethereum-nonce]. Enforcement is **strict ascending, no gaps**: a tx with nonce N can't be included before N−1. Too-low nonces are rejected immediately; too-high nonces sit in the mempool *queue* (un-mined, not propagated) until the gap fills — so one stuck/underpriced tx freezes all later txs from that account[^chainstack-nonce][^se-highnonce]. Since EIP-161, a created contract's nonce starts at 1, not 0[^eip161]. (Account-abstraction proposals introduce 2-D nonces that relax sequentiality, but that is not current EOA behavior.) (Confidence: fact, 3+ sources; the "no gaps" rule is the most common misconception.)

## 5. The EVM as a stack machine

The EVM is a **stack machine** with a maximum depth of **1024 items**, each a **256-bit (32-byte) word** (chosen for convenience with Keccak-256 and secp256k1)[^ethorg-evm][^cornell-evm]. The stack is LIFO and local to a message call. Stack ops: POP, PUSH1–PUSH32, DUP1–DUP16, SWAP1–SWAP16[^yellowpaper-evm][^coinculture].

**Quasi-Turing-complete, not Turing-complete:** it can run any program but only for a finite number of steps bounded by gas, so infinite loops can't stall the chain[^yellowpaper-evm][^ethbook-evm].

**The three data locations** (the developer-critical distinction)[^ethorg-evm][^calnix-evm][^alchemy-data]:

| Location | Lifetime | Addressing | Cost | Opcodes |
| --- | --- | --- | --- | --- |
| **Storage** | persistent (in the state trie) | 256-bit key→256-bit value | most expensive | SLOAD/SSTORE |
| **Memory** | per message call, zero-initialized | byte-addressed linear array | cheap | MLOAD/MSTORE/MSTORE8 |
| **Calldata** | read-only tx input | byte-addressed | cheapest | CALLDATALOAD/COPY/SIZE |
| **Transient storage** | per *transaction*, then cleared (EIP-1153) | 256-bit key→value | cheap | TLOAD/TSTORE |

The classic error is treating these as interchangeable: storage persists on-chain (and is by far the costliest), memory vanishes after the call, calldata is read-only and cheapest for unmodified inputs[^alchemy-data][^cyfrin-data].

**Machine state** μ tracks gas remaining, the program counter (starts at 0; JUMP/JUMPI may set it only to a valid JUMPDEST), memory, the stack, and output[^yellowpaper-evm][^cornell-evm]. **Opcode categories** (1 byte each, grouped in 16s): arithmetic (0x01+), comparison/bitwise (0x10+), Keccak (0x20), environment/block (0x30–0x4f), stack-memory-storage-flow (0x50–0x5f), PUSH/DUP/SWAP (0x60–0x9f), logging (0xa0–0xa4), system/call (0xf0–0xff)[^ethorg-opcodes][^coinculture].

**The message-call model** — five system opcodes with different contexts[^tevm-system][^zealynx-delegatecall][^se-calltypes]:

| Opcode | `msg.sender` | Storage used | `msg.value` | State changes |
| --- | --- | --- | --- | --- |
| `CALL` | the caller | callee's | sent | allowed |
| `DELEGATECALL` | **preserved from parent** | **caller's** | **preserved** | allowed (on caller's storage) |
| `STATICCALL` | the caller | callee's | 0 | **forbidden** (reverts on state change) |
| `CALLCODE` (deprecated) | the caller | caller's | sent | allowed |

`DELEGATECALL` runs the *target's code* against the *caller's storage* — the foundation of proxy/upgradeable and library patterns (its mechanics matter for the protocol; *writing* such contracts is `ethereum-smart-contract-development`'s job). Nested calls receive at most **63/64** of the remaining gas (EIP-150)[^tevm-system][^wolflo-gas]. (Confidence: fact, 3+ sources.)

## 6. In-EVM gas accounting

**Why gas exists:** (1) decouple computational cost from ETH's market price, and (2) **DoS defense** — since Turing-complete programs can't be statically bounded (the halting problem), the only safe approach is to meter execution and halt at the gas limit, forfeiting gas on out-of-gas[^ethbook-gas][^leastauth-gas]. Out-of-gas reverts all state changes except the nonce increment and gas paid[^ethbook-evm].

`[VOLATILE — all gas numbers are fork-dependent]`:
- **Per-opcode base costs**: STOP 0; ADD/SUB 3; MUL/DIV/MOD 5; comparison/bitwise 3; PUSH/DUP/SWAP 3; JUMP 8/JUMPI 10/JUMPDEST 1[^ethorg-opcodes][^evmcodes].
- **Intrinsic gas** (paid by the EOA before any code runs): base **21,000**; **+32,000** if contract-creation; **+4** per zero calldata byte, **+16** per non-zero byte (EIP-2028); plus any access-list cost[^wolflo-gas][^evmcodes][^ethorg-tx].
- **Memory expansion**: `Cmem(a) = 3·a + ⌊a²/512⌋` (a = words) — linear then quadratic, charged on the delta[^yellowpaper-evm].
- **Cold/warm access (EIP-2929, Berlin)**: first ("cold") touch of an account or slot is expensive; repeat ("warm") is cheap. Constants: cold account access **2600**, cold SLOAD **2100**, warm read **100**[^eip2929][^rareskills-2930].
- **SSTORE net metering (EIP-2200)**: `SSTORE_SET` (zero→nonzero) **20,000**; `SSTORE_RESET` (nonzero→nonzero) **5,000**; requires >2,300 gas (the call stipend) available[^eip2200][^geth-gastable].
- **Access lists (EIP-2930, Type-1 tx)**: pre-warm addresses (**2,400**/addr) and slots (**1,900**/slot) at intrinsic-gas time, paying the cold cost upfront at a small discount[^wolflo-gas][^rareskills-2930].
- **Gas refunds (EIP-3529, London)**: removed the SELFDESTRUCT refund, lowered the storage-clearing refund, and **capped total refunds at `gas_used / 5`** — to kill gas-token exploits and reduce block-size variance[^eip3529].

(Confidence: fact for the mechanisms; the *numbers* are post-Berlin/London/Cancun/Pectra and change at any fork — re-verify against the live opcode schedule.)

## 7. Precompiled contracts

Precompiles are special "contracts" bundled with the EVM at fixed low addresses (from 0x01), called like normal contracts with deterministic gas, providing crypto/utility operations too expensive in raw bytecode[^evmcodes-precompiled][^eco-precompile]. The set grows at hard forks. `[VOLATILE — set & gas are fork-dependent]`[^evmcodes-precompiled][^geth-contracts][^tevm-precompiles]:

| Addr | Name | Function | Since |
| --- | --- | --- | --- |
| 0x01 | ecRecover | ECDSA secp256k1 public-key recovery | Frontier |
| 0x02 | SHA2-256 | SHA-256 hash | Frontier |
| 0x03 | RIPEMD-160 | RIPEMD-160 hash | Frontier |
| 0x04 | identity | copy input→output | Frontier |
| 0x05 | modexp | modular exponentiation | Byzantium (repriced EIP-2565) |
| 0x06–0x08 | ecAdd / ecMul / ecPairing (BN254) | EC arithmetic + pairing (zkSNARK verification) | Byzantium |
| 0x09 | blake2f | BLAKE2 compression F | Istanbul |
| 0x0a | point evaluation (KZG) | verify a KZG commitment vs a versioned hash | Cancun (EIP-4844) |
| 0x0b–0x11 | BLS12-381 ops | BLS12-381 curve arithmetic & pairing | **Pectra (EIP-2537)** |

The BN254 trio (0x06–0x08) underpins zkSNARK verification on L1. The 0x0a point-evaluation/KZG precompile (Cancun) is the first precompile added specifically for **L2 scaling** — rollups use it to prove blob contents inside the EVM (see `data-availability-and-danksharding.md`). The crypto inside these (pairings, KZG) is `blockchain-crypto-primitives`' domain. (Confidence: fact for the set; gas values volatile.)

---

## References

[^eth-accounts]: ethereum.org — Ethereum accounts. https://ethereum.org/developers/docs/accounts/ — account types, four fields, address derivation. tier: docs.
[^consensys-state]: Consensys — Ethereum explained: Merkle trees, world state. https://consensys.io/blog/ethereum-explained-merkle-trees-world-state-transactions-and-more — world-state mapping, account state, block roots. tier: blog.
[^eco-eoa]: eco.com — EOA vs Contract Account. https://eco.com/support/en/articles/14796240-eoa-vs-contract-account-architecture — empty-trie/empty-code hashes, EIP-7702 consequences. tier: docs/blog.
[^geth-evm]: go-ethereum core/vm/evm.go. https://github.com/ethereum/go-ethereum/blob/master/core/vm/evm.go — CREATE/CREATE2 address derivation. tier: spec.
[^eip161]: EIP-161 — State trie clearing. https://github.com/ethereum/EIPs/blob/master/EIPS/eip-161.md — empty/dead accounts; CREATE nonce starts at 1. tier: eip.
[^eip7702]: EIP-7702 — Set Code for EOAs. https://eips.ethereum.org/EIPS/eip-7702 — type-0x04 set-code tx, delegation indicator. tier: eip.
[^ethorg-7702]: ethereum.org — Pectra / EIP-7702. https://ethereum.org/roadmap/pectra/7702/ — delegation semantics, revocation. tier: docs.
[^zealynx-7702]: Zealynx glossary — EIP-7702. https://www.zealynx.io/glossary/eip-7702 — 0xef0100 marker, is-EOA-check breakage. tier: blog.
[^yellowpaper]: Ethereum Yellow Paper (Gavin Wood). https://ethereum.github.io/yellowpaper/paper.pdf — formal Υ(σ,T) state transition, σ world state. tier: spec.
[^alchemy-tx]: Alchemy — How Ethereum transactions work. https://www.alchemy.com/docs/how-ethereum-transactions-work — chain of states σt→σt+1. tier: docs.
[^ethorg-evm]: ethereum.org — Ethereum Virtual Machine. https://ethereum.org/developers/docs/evm/ — state machine, Y(S,T), stack machine, data locations, transient storage. tier: docs.
[^ethbook-evm]: Mastering Ethereum, ch.13 (EVM). https://cypherpunks-core.github.io/ethereumbook/13evm.html — quasi-Turing, sandboxed state. tier: docs (book).
[^ethorg-mpt]: ethereum.org — Merkle Patricia Trie. https://ethereum.org/developers/docs/data-structures-and-encoding/patricia-merkle-trie/ — MPT structure, node types, HP encoding, RLP, the tries. tier: docs.
[^iitb-mpt]: IIT-B (Vijayakumaran) — Ethereum Data Structures and Encoding. https://www.ee.iitb.ac.in/~sarva/courses/EE465/2024/slides/EthereumDataEncoding.pdf — node types, HP encoding examples. tier: docs (academic).
[^vinidlidoo-mpt]: vinidlidoo — Ethereum Merkle Patricia Trie. https://vinidlidoo.github.io/blog/ethereum-merkle-patricia-trie/ — hexary-vs-binary clarification, 64-nibble paths. tier: blog.
[^easythereentropy]: easythereentropy — Understanding the Ethereum trie. https://easythereentropy.wordpress.com/2014/06/04/understanding-the-ethereum-trie/ — RLP node references, content addressing. tier: blog.
[^shyft-trie]: Shyft/Medium — Understanding trie databases in Ethereum. https://medium.com/shyft-network/understanding-trie-databases-in-ethereum-9f03d2c3325d — node types, HP vs RLP roles. tier: blog.
[^systeminternals]: systeminternals.dev — Ethereum Internals. https://systeminternals.dev/ethereum/ — account 4-tuple, four trie types, MPT write amplification. tier: blog.
[^ethorg-verkle]: ethereum.org — Verkle trees roadmap. https://ethereum.org/roadmap/verkle-trees/ — Verkle rationale, witness sizes. tier: docs.
[^eip6800]: EIP-6800 — Unified Verkle tree. https://eips.ethereum.org/EIPS/eip-6800 — witness sizes, Verkle-alongside-MPT. tier: eip.
[^vitalik-verge]: Vitalik Buterin — "Possible futures part 4: The Verge". https://vitalikblog.w3eth.io/general/2024/10/23/futures4.html — Verkle quantum vulnerability, binary-tree+STARK alternative. tier: blog (primary).
[^stateless-fyi]: stateless.fyi — Verkle Tree. https://stateless.fyi/trees/vkt-tree.html — Verkle vs binary-tree direction (2026). tier: docs/community.
[^nethereum-nonce]: Nethereum docs — Managing nonces. https://nethereum.readthedocs.io/en/latest/nethereum-managing-nonces/ — ordering + replay functions. tier: docs.
[^chainstack-nonce]: Chainstack — Ethereum nonce management. https://chainstack.com/ethereum-nonce-management/ — no-gaps rule, cascade. tier: blog.
[^se-highnonce]: ETH StackExchange — too-high nonce. https://ethereum.stackexchange.com/questions/2808/ — queue behavior, per-sender cap. tier: forum.
[^cornell-evm]: Cornell CS4998 — EVM Data Structures. https://cs4998.cornellblockchain.org/understanding-the-evm/the-ethereum-virtual-machine/evm-data-structures — stack/memory/storage/PC. tier: docs (academic).
[^yellowpaper-evm]: ethereum.org — Understanding the Yellow Paper's EVM. https://ethereum.org/developers/tutorials/yellow-paper-evm/ — machine state μ, 256-bit words, quasi-Turing, memory cost fn. tier: docs.
[^coinculture]: CoinCulture evm-tools guide. https://github.com/CoinCulture/evm-tools/blob/master/analysis/guide.md — opcode categories, stack/memory/storage model. tier: spec.
[^calnix-evm]: calnix — EVM storage opcodes. https://calnix.gitbook.io/eth-dev/evm-storage-opcodes/evm — memory/storage/calldata semantics. tier: docs/blog.
[^alchemy-data]: Alchemy — storage vs memory vs calldata. https://www.alchemy.com/docs/when-to-use-storage-vs-memory-vs-calldata-in-solidity — data-location costs & mistakes. tier: docs.
[^cyfrin-data]: Cyfrin — data location must be memory or calldata. https://www.cyfrin.io/blog/fixing-data-location-must-be-memory-or-calldata — data-location errors. tier: blog.
[^ethorg-opcodes]: ethereum.org — Opcodes for the EVM. https://ethereum.org/developers/docs/evm/opcodes/ — opcode list + base gas costs. tier: docs.
[^tevm-system]: Tevm — System Instructions. https://tevm.mintlify.app/evm/instructions/system — call-type matrix. tier: docs.
[^zealynx-delegatecall]: Zealynx glossary — Delegatecall. https://www.zealynx.io/glossary/delegatecall — call/delegatecall/staticcall matrix. tier: blog.
[^se-calltypes]: ETH StackExchange — CALL vs CALLCODE vs DELEGATECALL. https://ethereum.stackexchange.com/questions/3667/ — storage/sender semantics. tier: forum.
[^wolflo-gas]: wolflo/evm-opcodes gas.md. https://github.com/wolflo/evm-opcodes/blob/main/gas.md — intrinsic gas, access sets, SLOAD/SSTORE/CALL dynamic gas, 63/64 rule. tier: spec.
[^ethbook-gas]: Mastering Ethereum — gas chapter. https://github.com/merlox/ethereumbook/blob/tech-review/gas.asciidoc — gas as DoS defense, the halting problem. tier: docs (book).
[^leastauth-gas]: LeastAuthority — GasEcon. https://github.com/LeastAuthority/ethereum-analyses/blob/master/GasEcon.md — gas & the halting problem. tier: analysis.
[^evmcodes]: evm.codes — About. https://www.evm.codes/about — intrinsic/calldata gas, opcode reference. tier: docs.
[^ethorg-tx]: ethereum.org — Transactions. https://ethereum.org/developers/docs/transactions/ — 21,000 base gas, tx fields. tier: docs.
[^eip2929]: EIP-2929 — Gas cost increases for state-access opcodes. https://eips.ethereum.org/EIPS/eip-2929 — cold/warm access sets and constants. tier: eip.
[^rareskills-2930]: RareSkills — EIP-2930 access list. https://rareskills.io/post/eip-2930-optional-access-list-ethereum — access-list costs & net savings. tier: blog.
[^eip2200]: EIP-2200 — Structured net gas metering. https://eips.ethereum.org/EIPS/eip-2200 — SSTORE set/reset/clears constants. tier: eip.
[^geth-gastable]: go-ethereum core/vm/gas_table.go. https://github.com/ethereum/go-ethereum/blob/master/core/vm/gas_table.go — memory expansion + EIP-2200 SSTORE. tier: spec.
[^eip3529]: EIP-3529 — Reduction in refunds. https://eips.ethereum.org/EIPS/eip-3529 — SELFDESTRUCT refund removal, 1/5 cap. tier: eip.
[^evmcodes-precompiled]: evm.codes — Precompiled Contracts. https://www.evm.codes/precompiled — precompile list & gas. tier: docs.
[^eco-precompile]: eco.com — What is a Precompile. https://eco.com/support/en/articles/11011385-what-is-a-precompile-built-in-evm-contracts-explained — precompile use cases. tier: docs/blog.
[^geth-contracts]: go-ethereum core/vm/contracts.go. https://github.com/ethereum/go-ethereum/blob/master/core/vm/contracts.go — per-fork precompile address tables (Prague adds 0x0b–0x11). tier: spec.
[^tevm-precompiles]: Tevm — EVM precompiles table. https://voltaire.tevm.sh/evm/precompiles — precompile gas/forks. tier: docs.
