# distributed-systems-consensus

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude Code
**Original Path:** claude-code/distributed-systems-consensus

## Description
Distributed-systems & consensus theory plus blockchain consensus mechanisms — agreement across unreliable nodes. TRIGGER: CAP/PACELC, FLP, linearizability vs causal/eventual consistency, Lamport/vector clocks, state-machine replication, quorum intersection, crash/Byzantine failure models, safety vs liveness; crash-fault consensus (Paxos/Multi-Paxos, Raft, Viewstamped Replication, Zab); BFT (PBFT, Tendermint/CometBFT, HotStuff, 3f+1); gossip/epidemic, CRDTs; blockchain consensus (Nakamoto/longest-chain PoW, Proof of Stake, Ethereum Gasper = Casper FFG + LMD-GHOST, Cardano Ouroboros, Solana Tower BFT, fork-choice, finality, Sybil resistance, the trilemma, long-range/nothing-at-stake/selfish-mining attacks). SKIP: architecture/microservices/sagas → software-engineering-patterns; whole-chain designs + chain params (Bitcoin halving/21M/retarget) → per-chain skills (bitcoin-protocol-expert, ethereum-protocol-expert); crypto-primitive math (signatures/VRF/hashing) → blockchain-crypto-primitives; GPU "distributed training" → distributed-training; MongoDB replica-set ops → mongodb-expert.

---

# Distributed Systems & Consensus

Expert reference for **distributed consensus** — how independent nodes agree on a single value or an ordered log despite unreliable networks, slow processes, and faulty (even malicious) participants. Spans three layers that share one substrate: (1) the **theory** (impossibility results, consistency models, failure models), (2) **classical consensus algorithms** (crash-fault Paxos/Raft family; Byzantine PBFT/HotStuff/Tendermint family), and (3) **blockchain consensus mechanisms** (PoW, PoS, fork-choice, finality, Sybil resistance).

**The unifying idea:** consensus = three properties — **agreement** (no two correct nodes decide differently), **validity** (a decided value was proposed), **termination** (correct nodes eventually decide). Almost every result below is about which of these you can guarantee, under which failure model, at what node count. Two numbers recur: **2f+1** nodes to tolerate **f crash** faults, **3f+1** to tolerate **f Byzantine** faults — both driven by *quorum intersection*.

## How to use this skill

Read this page for the conceptual map and the load-bearing facts (bounds, theorem statements, which-algorithm-when). Load the matching reference file for depth:

| Reference | Covers |
|---|---|
| `references/foundations-theory.md` | CAP/PACELC, FLP impossibility, logical time (Lamport/vector clocks), state-machine replication, quorums & quorum intersection, replication/partitioning, failure models, safety vs liveness, consistency models (linearizability → eventual), CALM theorem |
| `references/crash-fault-consensus.md` | Paxos & Multi-Paxos, Raft (election/log/safety/membership), Viewstamped Replication, ZooKeeper/Zab, and how they relate |
| `references/byzantine-gossip-crdt.md` | Byzantine Generals & 3f+1, PBFT, Tendermint/CometBFT, HotStuff (linear view-change, chained/pipelined), gossip/epidemic protocols, SWIM, CRDTs |
| `references/blockchain-consensus.md` | Nakamoto/longest-chain PoW, Proof of Stake, Ethereum Gasper (Casper FFG + LMD-GHOST), Ouroboros, Tower BFT, finality (probabilistic vs economic), fork-choice rules, Sybil resistance, the trilemma, consensus-layer attacks |

## Core map (the 60-second version)

### Theory substrate (see `foundations-theory.md`)
- **FLP impossibility** (Fischer-Lynch-Paterson 1985): no *deterministic* protocol guarantees consensus in a *fully asynchronous* system if even *one* process can crash. It is a **liveness** result — real systems stay safe but may not terminate. Circumvented (not violated) by randomization, partial synchrony, or failure detectors.[^flp]
- **CAP** (Gilbert-Lynch 2002 proof of Brewer's conjecture): under a network **partition**, you cannot have both **linearizable consistency** and **availability**. "Pick 2 of 3" is a misreading — partition tolerance is not optional, and Brewer himself called the framing misleading. **PACELC** adds: *else* (no partition) you still trade **latency vs consistency**.[^cap][^pacelc]
- **Consistency models** form a hierarchy (Jepsen): **linearizable** (real-time, single-object, composable) > sequential > **causal** (strongest model achievable while staying available under partition) > **eventual** (replicas converge if writes stop). Linearizability ≠ serializability (serializability is multi-object transaction isolation, no real-time order).[^jepsen][^lin]
- **Logical time:** Lamport clocks give a total order consistent with happens-before but **cannot detect concurrency**; **vector clocks** can (causality ⟺ vector order), at O(n) size.[^lamport]
- **State-machine replication (SMR):** replicate a deterministic state machine + feed every replica the same command order = a fault-tolerant service. The "order" requirement is **total-order (atomic) broadcast**, which is **equivalent to consensus** (Chandra-Toueg).[^smr]
- **Quorum intersection:** any two majorities of 2f+1 nodes share ≥1 node — that overlap is what carries a committed value across leader changes. **Flexible Paxos**: only the two *phases'* quorums must intersect, not every quorum with itself.[^quorum]

### Crash-fault-tolerant consensus (see `crash-fault-consensus.md`)
Leader-based replicated log, majority quorum (f+1 of 2f+1), non-Byzantine (nodes crash, never lie).
- **Paxos** (Lamport): two phases (prepare/promise, accept/accepted), safe by majority intersection, no guaranteed liveness. **Multi-Paxos**: stable leader runs Phase 1 once, then skips to Phase 2 per log slot.[^paxos]
- **Raft** (Ongaro-Ousterhout 2014): decomposed for understandability — terms + randomized election timeouts, AppendEntries log replication, the Log-Matching / Leader-Completeness / State-Machine-Safety properties. Whether it's actually "simpler" than Paxos is contested.[^raft]
- **Viewstamped Replication** (Oki-Liskov 1988): independently discovered consensus; views, primary/backup, view-change. **Zab** (ZooKeeper): atomic broadcast with primary-order/FIFO and zxid epochs.[^vr]

### Byzantine fault tolerance (see `byzantine-gossip-crdt.md`)
Nodes may behave arbitrarily/maliciously. **n ≥ 3f+1** for unauthenticated Byzantine agreement (the tight bound under partial synchrony); 2f+1-size quorum certificates.
- **PBFT** (Castro-Liskov 1999): three phases (pre-prepare/prepare/commit), O(n²) normal-case messages, O(n³) view-change. First *practical* BFT.[^pbft]
- **HotStuff** (2019): **linear O(n)** communication and **linear view-change** via threshold signatures + leader-collects-votes; three-chain commit, pipelinable. Basis of Diem; HotStuff-2/Jolteon lineage.[^hotstuff]
- **Tendermint/CometBFT**: round-based propose/prevote/precommit, **>2/3 voting power**, locking → instant deterministic finality (no forks). The BFT core under Cosmos.[^tendermint]
- **Gossip/epidemic** (Demers 1987): anti-entropy + rumor-mongering, O(log n) dissemination; **SWIM** for membership/failure detection. **CRDTs** (Shapiro 2011): state-based (join-semilattice merge) or op-based (commutative), giving Strong Eventual Consistency.[^gossip][^crdt]

### Blockchain consensus (see `blockchain-consensus.md`)
Open/permissionless → needs **Sybil resistance** (classical BFT assumes a known validator set; here identity must be made costly).
- **Nakamoto PoW** (Bitcoin 2008): hash-based leader election + longest/heaviest-chain rule, **probabilistic finality** (reversal probability decays exponentially with confirmations), honest-majority-of-hashpower assumption. Sybil resistance via compute cost.[^pow]
- **Proof of Stake:** voting power ∝ bonded stake, **slashing** penalizes provable misbehavior; must handle **nothing-at-stake**, **long-range attacks** (→ weak subjectivity checkpoints).[^pos]
- **Ethereum Gasper = Casper FFG + LMD-GHOST**: FFG is a **finality gadget** (checkpoints, 2/3 supermajority, two slashing "commandments", economic finality); LMD-GHOST is the **fork-choice** rule (heaviest subtree by latest attestations). The gadget constrains the fork-choice to descend from the last justified checkpoint.[^gasper]
- **Ouroboros** (Cardano): provably secure PoS, VRF private leader election, common-prefix/chain-quality/chain-growth proofs. **Tower BFT** (Solana): PBFT-style voting clocked by Proof-of-History.[^ouroboros]
- **Finality:** probabilistic (PoW) vs deterministic/**economic** (BFT gadgets); a **fork-choice rule** picks a revisable head, a **finality rule** makes blocks irreversible.[^finality]
- **Trilemma** (Buterin): scalability/security/decentralization — a heuristic, not a theorem (and formally contested).[^trilemma]

## Decision guidance

- **Choosing a consensus algorithm:** known/trusted members + only crashes → **Raft** (new systems: etcd, Consul, CockroachDB) or **Multi-Paxos** (Chubby, Spanner). Mutually-distrusting members / adversarial → **BFT**: small fixed set → **PBFT**; large rotating validator set / blockchain → **HotStuff** or **Tendermint**. Open/permissionless → a **blockchain** Sybil-resistant mechanism (PoW/PoS).
- **Fault budget math:** to survive **f crash** failures you need **2f+1** replicas (majority quorum f+1). To survive **f Byzantine** failures you need **3f+1** (quorum 2f+1). A 5-node cluster tolerates 2 crashes or 1 Byzantine fault.
- **Consistency target:** need real-time global order across one object → **linearizability** (costs a consensus round / can't be totally available under partition). Want always-available + multi-master → settle for **causal+** or **eventual** consistency (CRDTs converge without coordination, but only for monotonic/commutative operations — CALM theorem).
- **Blockchain finality expectation:** PoW → wait N confirmations for *probabilistic* finality (size confirmations against a ~25% attacker, not 51%, per selfish-mining results). BFT-gadget PoS → *economic* finality after the supermajority commit, but liveness can still stall under partition/attack.

## Anti-patterns and common confusions

- **"CAP says pick 2 of 3."** No. Partitions are a fact of networks, not a dial; the real choice is C-vs-A *during* a partition, and L-vs-C otherwise (PACELC). Both Brewer and Kleppmann reject the 2-of-3 framing.[^cap][^pacelc]
- **Confusing "consistency" senses.** CAP/linearizability consistency (real-time recency of one object) is unrelated to ACID "C" (a transaction preserves invariants). And linearizability ≠ serializability.[^lin]
- **"Timeouts beat FLP."** They don't — adding failure detectors/timeouts *changes the model* (to partial synchrony); FLP still applies to the asynchronous model including its detectors.[^flp]
- **Treating Lamport timestamps as causality oracles.** `C(a) < C(b)` does **not** imply `a → b`; only vector clocks detect concurrency.[^lamport]
- **Assuming Raft/Paxos handle malicious nodes.** They tolerate *crashes only*. A single lying node breaks them — you need BFT (3f+1).[^raft][^pbft]
- **Assuming "honest majority (>50%)" is the PoW safety bar.** Selfish mining (Eyal-Sirer) shows Bitcoin isn't incentive-compatible above ~25-33% adversarial hashpower; design for that, not 51%.[^attacks]
- **Calling blockchain "finality" absolute.** PoW finality is probabilistic forever; PoS "deterministic" finality is *economic and conditional* (holds while <1/3 stake is adversarial) and can degrade to a liveness stall.[^finality]
- **Conflating fork-choice with finality.** A fork-choice rule (longest-chain, GHOST, LMD-GHOST) picks a *revisable* head; a finality gadget makes checkpoints *irreversible*. Different jobs.[^gasper]

## Peer skills (deferral)

- **Language-agnostic architecture, microservices decomposition, saga/CQRS/event-sourcing patterns** → `software-engineering-patterns` (and its `microservices-patterns` reference). This skill owns the *consensus algorithm*; that one owns *service architecture*.
- **Full per-chain protocol designs and chain-specific parameters** (the whole Bitcoin/Ethereum/Solana system: execution, mempool, gas, tokenomics, EVM; plus chain-specific consensus *parameters* — Bitcoin's difficulty-retarget cadence, the halving, the 21M cap; per-chain block/tx formats) → the per-chain skills (e.g. `bitcoin-protocol-expert`). This skill covers the consensus *mechanism as a class* (longest-chain PoW, PoS, BFT, finality), reciprocal with those skills' own SKIP clauses.
- **Cryptographic primitive internals** (digital signatures, VRFs, threshold/BLS signatures, hash functions, zk-proofs — the math) → `blockchain-crypto-primitives` (a planned sibling, may not be installed yet). Here they are named as black-box building blocks; if that skill is absent, the math is simply out of scope.
- **GPU-parallel "distributed training"** of LLMs (FSDP/ZeRO/tensor-parallel) → the `distributed-training` reference in the `ai-llm-model-layer` hub. Unrelated to consensus despite the word "distributed."
- **MongoDB replica-set elections / oplog / read-concern operations** → `mongodb-expert`. MongoDB's elections are a *Raft-derived* protocol; the *theory* is here, the *product operations* are there.
- **Self-healing / autonomic control loops (MAPE-K), Kubernetes reconciliation** → `self-healing-systems-autonomic-computing` (registered under the `devops-containers-cicd` hub). Control loops, not agreement protocols.

## References

[^flp]: Fischer, Lynch, Paterson, "Impossibility of Distributed Consensus with One Faulty Process," JACM 1985. https://www.cs.princeton.edu/courses/archive/spring24/cos418/papers/flp.pdf — *paper*
[^cap]: Gilbert & Lynch, "Brewer's Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web Services," SIGACT News 2002. https://www.comp.nus.edu.sg/~gilbert/pubs/BrewersConjecture-SigAct.pdf — *paper*. Brewer, "CAP Twelve Years Later," IEEE Computer 2012. https://www.infoq.com/articles/cap-twelve-years-later-how-the-rules-have-changed/ — *paper*
[^pacelc]: Abadi, "Consistency Tradeoffs in Modern Distributed Database System Design," IEEE Computer 2012. https://www.cs.umd.edu/~abadi/papers/abadi-pacelc.pdf — *paper*
[^jepsen]: Jepsen, "Consistency Models." https://jepsen.io/consistency/models — *docs*
[^lin]: Herlihy & Wing, "Linearizability: A Correctness Condition for Concurrent Objects," ACM TOPLAS 1990. https://cs.brown.edu/people/mph/HerlihyW90/p463-herlihy.pdf — *paper*. Bailis, "Linearizability versus Serializability." https://www.bailis.org/blog/linearizability-versus-serializability/ — *blog*
[^lamport]: Lamport, "Time, Clocks, and the Ordering of Events in a Distributed System," CACM 1978. https://lamport.azurewebsites.net/pubs/time-clocks.pdf — *paper*
[^smr]: Schneider, "Implementing Fault-Tolerant Services Using the State Machine Approach," ACM Computing Surveys 1990. https://www.cs.cornell.edu/fbs/publications/SMSurvey.pdf — *paper*. Chandra & Toueg, "Unreliable Failure Detectors for Reliable Distributed Systems," JACM 1996. https://www.cs.princeton.edu/courses/archive/fall08/cos597B/papers/unreliable.pdf — *paper*
[^quorum]: Howard, Malkhi, Spiegelman, "Flexible Paxos: Quorum Intersection Revisited," OPODIS 2016. https://fpaxos.github.io/ — *paper/docs*
[^paxos]: Lamport, "Paxos Made Simple," SIGACT News 2001. https://lamport.azurewebsites.net/pubs/paxos-simple.pdf — *paper*
[^raft]: Ongaro & Ousterhout, "In Search of an Understandable Consensus Algorithm (Extended)," USENIX ATC 2014. https://web.stanford.edu/~ouster/cgi-bin/papers/raft-extended.pdf — *paper*. Howard & Mortier, "Paxos vs Raft," PaPoC 2020. https://arxiv.org/abs/2004.05074 — *paper*
[^vr]: Liskov & Cowling, "Viewstamped Replication Revisited," MIT-CSAIL-TR-2012-021. https://pmg.csail.mit.edu/papers/vr-revisited.pdf — *paper*. Junqueira, Reed, Serafini, "Zab," DSN 2011. https://marcoserafini.github.io/assets/pdf/zab.pdf — *paper*
[^pbft]: Castro & Liskov, "Practical Byzantine Fault Tolerance," OSDI 1999. https://pmg.csail.mit.edu/papers/bft-tocs.pdf — *paper*
[^hotstuff]: Yin, Malkhi, Reiter, Golan-Gueta, Abraham, "HotStuff: BFT Consensus in the Lens of Blockchain," arXiv 1803.05069 (PODC 2019). https://arxiv.org/pdf/1803.05069 — *paper*
[^tendermint]: Buchman, Kwon, Milosevic, "The latest gossip on BFT consensus," arXiv 1807.04938. https://arxiv.org/abs/1807.04938 — *paper*
[^gossip]: Demers et al., "Epidemic Algorithms for Replicated Database Maintenance," PODC 1987. https://www.bitsavers.org/pdf/xerox/parc/techReports/CSL-89-1_Epidemic_Algorithms_for_Replicated_Database_Maintenance.pdf — *paper*. Das, Gupta, Motivala, "SWIM," DSN 2002. https://www.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf — *paper*
[^crdt]: Shapiro, Preguiça, Baquero, Zawirski, "Conflict-free Replicated Data Types," SSS 2011. https://perso.lip6.fr/Marc.Shapiro/papers/2011/CRDTs_SSS-2011.pdf — *paper*
[^pow]: Nakamoto, "Bitcoin: A Peer-to-Peer Electronic Cash System," 2008. https://bitcoin.org/bitcoin.pdf — *paper*
[^pos]: Buterin & Griffith, "Casper the Friendly Finality Gadget," arXiv 1710.09437. https://arxiv.org/pdf/1710.09437 — *paper*. ethereum.org, "Proof-of-stake (PoS)." https://ethereum.org/developers/docs/consensus-mechanisms/pos/ — *docs* (verified-as-of 2026-06-16)
[^gasper]: Buterin et al., "Combining GHOST and Casper" (Gasper), arXiv 2003.03052. https://arxiv.org/abs/2003.03052 — *paper*. ethereum.org, "Gasper." https://ethereum.org/developers/docs/consensus-mechanisms/pos/gasper/ — *docs* (verified-as-of 2026-06-16)
[^ouroboros]: Kiayias, Russell, David, Oliynykov, "Ouroboros: A Provably Secure PoS Blockchain Protocol," CRYPTO 2017. https://eprint.iacr.org/2016/889 — *paper*
[^finality]: Anceaume et al., "On Finality in Blockchains," OPODIS 2021. https://drops.dagstuhl.de/storage/00lipics/lipics-vol217-opodis2021/LIPIcs.OPODIS.2021.6/LIPIcs.OPODIS.2021.6.pdf — *paper*. Trail of Bits, "The Engineer's Guide to Blockchain Finality." https://blog.trailofbits.com/2023/08/23/the-engineers-guide-to-blockchain-finality/ — *blog*
[^trilemma]: Buterin, "Why sharding is great: demystifying the technical properties." https://vitalik.eth.limo/general/2021/04/07/sharding.html — *blog* (verified-as-of 2026-06-16)
[^attacks]: Eyal & Sirer, "Majority Is Not Enough: Bitcoin Mining Is Vulnerable," CACM (arXiv 1311.0243). https://cacm.acm.org/research/majority-is-not-enough/ — *paper*