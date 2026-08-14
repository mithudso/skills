<!-- Provenance: reference under the standalone `distributed-systems-consensus` skill. Created 2026-06-16 via /dr deep-research. Synthesized from primary papers + authoritative secondary sources; treat all content as synthesized domain knowledge. -->

# Distributed Systems Foundations — Theory

The shared theoretical substrate beneath every consensus algorithm. Per-claim confidence: **[fact]** = 3+ independent sources; **[qualified]** = 2; **[tentative]** = 1.

## Contents
- [CAP theorem and PACELC](#cap-theorem-and-pacelc)
- [FLP impossibility](#flp-impossibility)
- [Logical time: Lamport and vector clocks](#logical-time-lamport-and-vector-clocks)
- [State-machine replication and the consensus ↔ atomic-broadcast equivalence](#state-machine-replication-and-the-consensus--atomic-broadcast-equivalence)
- [Quorums and quorum intersection](#quorums-and-quorum-intersection)
- [Replication and partitioning](#replication-and-partitioning)
- [Failure models](#failure-models)
- [Safety vs liveness](#safety-vs-liveness)
- [Consistency models](#consistency-models)
- [Common confusions](#common-confusions-disconfirming)
- [Sources](#sources)

## CAP theorem and PACELC

CAP began as a conjecture by Eric Brewer at PODC 2000 and was proven by Seth Gilbert and Nancy Lynch in 2002 ("Brewer's Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web Services").[^gl] The precise result is narrow: in the **asynchronous network model**, it is impossible to implement a read/write register guaranteeing both **availability** and **atomic (linearizable) consistency** in all executions *including those where messages are lost* (Theorem 1).[^gl] The proof is a two-line partition argument: split nodes into G1/G2, drop all messages between them; a write of v1 completes at G1 (availability), a subsequent read at G2 completes (availability) but cannot have seen the write, so it returns stale v0 — violating atomicity.[^gl] They also prove the partially-synchronous analogue (Theorem 2), but there a partially-synchronous system *can* return atomic data when no messages are lost, degrading only during partitions.[^gl] **[fact]**

Definitions that matter: "consistency" in CAP means **linearizability** (= atomic / strong consistency) and has **nothing to do with the C in ACID**.[^gl][^kleppmann] "Availability" means *every* request to a *non-failing* node returns a non-error response (an unusually strong definition).[^gl] "Partition tolerance" is not a choice — it means the network may drop/delay messages, which all real networks do.[^kleppmann] **[fact]**

The popular **"pick 2 of 3"** framing is rejected even by its authors. Brewer (2012, "CAP Twelve Years Later") wrote that "2 of 3" "was always misleading" because partitions are rare (so there's no reason to forfeit C or A when un-partitioned), the C-vs-A choice recurs at fine granularity, and all three are "more continuous than binary"; CAP "prohibits only a tiny part of the design space."[^brewer] Kleppmann ("Please stop calling databases CP or AP," 2015; arXiv:1509.05393) argues CAP should be retired, showing real systems are often **neither CP nor AP — just "P"** under CAP's strict definitions, and that "2 of 3" mathematically permits picking 1 or 0 of 3.[^kleppmann] **[fact]**

**PACELC** (Daniel Abadi, 2010 blog, IEEE Computer 2012) fixes CAP's biggest blind spot — latency: **if Partition (P), trade Availability vs Consistency (A/C); Else (E), trade Latency vs Consistency (L/C)**.[^pacelc] Abadi's thesis: the consistency/latency tradeoff "is present at all times during system operation, whereas CAP is only relevant in the arguably rare case of a network partition."[^pacelc] The ELC tradeoff applies only to *replicated* systems; an atomically consistent system pays at least one extra message delay on reads+writes.[^pacelc] Systems classify as PA/EL, PA/EC, PC/EL, PC/EC.[^pacelc] PACELC received a formal proof in 2018 ("Proving PACELC," SIGACT News): a partition-free system cannot always guarantee both consistency and operation latency below half the minimum message delay.[^provingpacelc] **[fact]**

## FLP impossibility

The FLP result (Fischer, Lynch, Paterson; presented 1983, published JACM 1985) is the foundational impossibility theorem for consensus. **Theorem: no consensus protocol is totally correct in spite of one fault.**[^flp] Concretely: in a *completely asynchronous* message-passing system (no bound on process speed or message delay, no synchronized clocks, no failure detection), no *deterministic* protocol can guarantee that all non-faulty processes agree on a binary value even if **only one process may crash** (fail by stopping, undetectably).[^flp] The model assumes a *reliable* message system and no Byzantine faults — yet "the stopping of a single process at an inopportune time" can prevent agreement.[^flp] A correct protocol must satisfy agreement, validity (both 0 and 1 are possible decisions), and termination.[^flp] **[fact]**

The proof is the **bivalence argument**. A configuration is **bivalent** if both decision values are still reachable; **univalent** if only one is.[^flp] Lemma: every such protocol has a bivalent initial configuration. Lemma: from any bivalent configuration with an applicable event e, one can reach another bivalent configuration in which e was applied. Combining these constructs an infinite **non-deciding run** that stays bivalent forever by delaying a single decider message — every process takes infinitely many steps and receives every message, yet no decision commits.[^flp] The engine is *indistinguishability*: a crashed process cannot be told apart from a merely slow one.[^flp] **[fact]**

Real systems **circumvent (not violate)** FLP by weakening a model assumption:[^randconsensus][^flpcap] **(a) Randomization** — Ben-Or (1983) terminates with probability 1 if a majority is correct; **(b) Partial synchrony** — Dwork-Lynch-Stockmeyer (1988): consensus is solvable if the network is eventually synchronous; **(c) Failure detectors** — Chandra-Toueg (1996) identify the weakest detector (◇W / Ω) sufficient for consensus. Crucially, **FLP is a liveness result, not a safety result**: protocols stay safe (never decide inconsistently); they may only fail to terminate.[^flpcap] Timeout-based detection does **not** "beat" FLP — it changes the model, and FLP applies to the whole system including its detectors.[^flpcap] **[fact]**

## Logical time: Lamport and vector clocks

Lamport's "Time, Clocks, and the Ordering of Events in a Distributed System" (CACM 1978) defines the **happens-before** relation (→): (1) same-process precedence; (2) a message send → its receipt; (3) transitivity. Two events are **concurrent** if neither precedes the other — an irreflexive **partial** order.[^lamport] **Lamport (scalar) clocks** satisfy the **Clock Condition: if a→b then C(a) < C(b)** (increment between events; on receipt set clock > received timestamp); tie-break on process ID for a total order.[^lamport] The crucial limitation: the implication is **one-directional** — `a→b ⟹ C(a)<C(b)` but **NOT the converse**, so scalar timestamps **cannot detect concurrency**.[^lamport] **[fact]**

**Vector clocks** (independently Fidge 1988 and Mattern 1989) fix this: each process keeps a length-n vector, increments its own component per event, attaches its vector to messages, and takes the componentwise maximum on receipt.[^vc] The defining property is an isomorphism of causal and temporal structure: **e → e′ ⟺ V(e) < V(e′)**, and **e ‖ e′ ⟺ V(e) ‖ V(e′)** — vector clocks **can** detect both causality and concurrency.[^vc] Charron-Bost proved capturing causality this way needs vectors of size n in the worst case.[^vc] **[fact]**

## State-machine replication and the consensus ↔ atomic-broadcast equivalence

SMR builds a fault-tolerant service by modeling it as a **deterministic state machine** (state + atomic, deterministic commands), replicating it, and feeding every replica **the same command sequence in the same order** — replicas then stay identical.[^schneider][^lamport] Lamport (1978) first observed that totally ordering requests yields an algorithm to implement any state machine.[^lamport] Schneider's 1990 ACM Computing Surveys tutorial is canonical: Replica Coordination decomposes into **Agreement** (every non-faulty replica receives every request) and **Order** (replicas process requests in the same relative order).[^schneider] Determinism is essential — non-deterministic operations would let replicas diverge.[^schneider] The tutorial gives protocols for both fail-stop and Byzantine models.[^schneider] **[fact]**

The "Order" requirement is precisely **total-order (atomic) broadcast**. The deep result is the **equivalence of consensus and atomic broadcast**, proven by Chandra & Toueg (1996, "Unreliable Failure Detectors for Reliable Distributed Systems"): "Consensus and Atomic Broadcast are reducible to each other in asynchronous systems with crash failures."[^ct] One builds each from the other, so they are the same problem with reliable channels — and a corollary is that FLP applies to atomic broadcast too.[^ct] (Terminology note: Défago et al. prefer "total-order broadcast" over "atomic broadcast," since "atomicity" suggests agreement rather than ordering.[^tob]) **[fact]**

## Quorums and quorum intersection

A **quorum system** over node set V is a family of subsets such that **every two quorums intersect** (Q1 ∩ Q2 ≠ ∅).[^ethzquorum][^aspnesquorum] Intersection guarantees a read quorum sees at least one node that participated in the most recent write.[^ethzquorum] The canonical **Majority** quorum system gives every quorum ⌊n/2⌋+1 nodes; in Dynamo-style leaderless replication the practical rule is **R + W > N**.[^ddiarepl] **[fact]**

**Why 2f+1 tolerates f crash faults:** with n = 2f+1, any two majorities (each size f+1) must intersect, and losing any f nodes still leaves a majority (f+1) — making Majority **f-resilient**.[^ethzquorum][^aspnesquorum] **Byzantine** quorum systems need larger intersections: an **f-masking** system requires intersections of 2f+1 (so ≥f+1 correct nodes outvote ≤f Byzantine ones), needing n > 4f for non-self-verifying data; **f-disseminating** (self-verifying/signed data) needs intersection f+1 and n > 3f.[^ethzquorum] **[fact]**

**Flexible Paxos** (Howard, Malkhi, Spiegelman; OPODIS 2016) is the key refinement: majority quorums are **not necessary** — intersection is required only **across the two phases** (leader-election quorum Q1 must intersect every replication quorum Q2), but Q1 and Q2 need not be majorities and need not self-intersect.[^fpaxos] This enables weighted/flexible quorums (e.g. shrink phase-2 quorums to raise throughput, enlarge phase-1 to compensate): "Paxos is just a single point on a broad spectrum of possibilities for safely reaching distributed consensus."[^fpaxos] **[fact]**

## Replication and partitioning

Kleppmann's *Designing Data-Intensive Applications* is the standard secondary treatment.[^ddiarepl] **Replication** keeps copies of the same data on multiple nodes; the difficulty is propagating *changes*. Three architectures:[^ddiarepl] **single-leader** (one leader takes writes, streams a replication log to read-only followers — PostgreSQL, MySQL, MongoDB per-shard); **multi-leader** (multiple write-accepting nodes, needs conflict resolution — LWW or CRDTs); **leaderless** (Dynamo-style: writes/reads sent to all replicas, convergence via quorum intersection + read-repair). **[fact]**

**Synchronous** replication waits for follower ack before confirming (an up-to-date copy, but one slow node blocks writes); **asynchronous** does not wait (low latency, risks data loss on leader failure → replication lag → eventual consistency); **semi-synchronous** makes one follower synchronous as a compromise.[^ddiarepl] **Partitioning (sharding)** splits data so each record belongs to one shard; the two strategies are **key-range** (efficient range scans, hot-spot risk) and **hash** (even distribution, loses range queries).[^ddiarepl] **Sharding is normally combined with replication** — each shard is itself replicated, so a record belongs to one shard but lives on several nodes; cross-partition causality is hard (independent partitions apply writes at different rates).[^ddiarepl] **[fact]**

## Failure models

A severity hierarchy:[^failures] **Crash-stop (= fail-stop)** — a process halts and never acts again (simplest, most-assumed). **Crash-recovery** — halts but may reboot and rejoin (needs durable state). **Omission** — a process/channel fails to send or receive a message. **Byzantine / arbitrary** — a process behaves arbitrarily: omitting, taking incorrect steps, sending conflicting or malicious messages.[^flp][^failures] The **fail-stop model** is formalized by Schlichting & Schneider (1983): a fail-stop processor halts on failure, its failed state is detectable by others, and volatile storage is lost while stable storage survives.[^failstop] There is a translation hierarchy: crash-tolerant algorithms can be transformed to tolerate omission then arbitrary failures (Neiger-Toueg), often by masking corruption with checksums/authentication.[^failures] **Important: a network partition is NOT a node failure** — in a partition the nodes stay up but the network drops messages; that is the "P" of CAP, a property of the system model, not a process.[^gl][^kleppmann] **[fact]**

## Safety vs liveness

Named by Lamport (1977), formalized by Alpern & Schneider, "Defining Liveness" (IPL 1985).[^alpern] A **safety** property says **"nothing bad ever happens"** (mutual exclusion, agreement, "no two processes decide differently"); a **liveness** property says **"something good eventually happens"** (termination, starvation-freedom).[^alpern] Formalization: a property is a set of infinite execution sequences; **safety** is one whose violation is detectable in a *finite prefix* (proved via invariance arguments); **liveness** is one where *every finite prefix can be extended* to a satisfying execution (proved via well-foundedness).[^alpern] Central theorem: **every property is the intersection of a safety property and a liveness property.**[^alpern] **[fact]**

This frames the earlier results: **FLP is a liveness impossibility** — consensus protocols preserve *safety* (agreement always holds) but cannot guarantee *liveness* (termination) under async + one crash.[^flpcap] Likewise in CAP, **consistency (linearizability) is a safety property and availability is a liveness property** — a framing Kleppmann uses to argue Brewer's "continuous" reinterpretation is incoherent (safety/liveness properties hold or do not; there is no degree of continuity).[^kleppmann] **[fact]**

## Consistency models

Models form a hierarchy; Jepsen defines a model as a set of legal histories, where A implies B iff A ⊆ B.[^jepsen] **[fact]**

**Linearizability** (Herlihy & Wing, ACM TOPLAS 1990) is the strongest single-object model: each operation appears to take effect **atomically at a single instant between its invocation and response**, consistent with **real-time** (if A completes before B begins, A precedes B).[^hw] It is a **local (composable)** property — if every object is individually linearizable, the whole system is, *unlike* sequential consistency or serializability.[^hw] Linearizability = "atomic consistency" = the CAP "C," and cannot be totally available under partition.[^jepsen][^gl] **[fact]**

**Sequential consistency** (Lamport) requires *some* single total order respecting each process's program order **but imposes no real-time constraint** — linearizability is sequential consistency *plus* the real-time bound.[^jepsen] Attiya & Welch (1994) proved that with imperfect/no clocks, **linearizability is strictly more expensive than sequential consistency**.[^attiya] **[fact]**

**Causal consistency** (Lamport happens-before; Ahamad et al. 1993): causally-related operations appear in the same order everywhere, while concurrent operations may be seen in different orders.[^jepsen][^cops] **Causal+ / convergent causal consistency** (Lloyd et al., COPS, SOSP 2011) adds convergent conflict handling.[^cops] The pivotal result: **causal consistency (specifically Real-Time Causal) is the STRONGEST model achievable in an always-available, partition-tolerant system** (Mahajan-Alvisi-Dahlin).[^mad] Causal consistency alone has no liveness guarantee (satisfiable by never propagating writes), which is why convergence is added.[^mad] **[fact]**

**Eventual consistency** is the weakest: if writes stop, replicas eventually converge — no ordering or recency guarantee.[^mad] Convergent causal consistency is strictly stronger than (and implies) eventual; eventual does not imply causal.[^mad] **[fact]**

**Session guarantees** (Terry et al., Bayou, 1994) are per-client guarantees mitigating replication-lag anomalies:[^session][^ddiarepl] **Read-your-writes** (see your own prior writes), **Monotonic reads** (reads never go backward — typically pin a session to one replica), **Monotonic writes** (a client's writes apply in order), **Writes-follow-reads**. DDIA reframes the user-visible anomalies as read-after-write, monotonic-reads, and consistent-prefix violations.[^ddiarepl] **[fact]**

**The CALM theorem** (Hellerstein & Alvaro; conjectured PODS 2010, proved by Ameloot-Neven-Van den Bussche; CACM 2020 "Keeping CALM"): **a problem has a consistent, coordination-free distributed implementation if and only if it is monotonic** (Consistency As Logical Monotonicity).[^calm] Monotonic problems are safe under missing information and need no coordination; **non-monotonic** problems (where new input retracts a prior output) require coordination.[^calm] CALM is a program/problem-level result enabling static analysis — the constructive positive counterpart to CAP's negative result, contrasted with storage-level consistency that needs expensive runtime enforcement.[^calm] **[fact]**

## Common confusions (disconfirming)

1. **"CAP = pick 2 of 3" is rejected by its own authors** (Brewer 2012; Kleppmann 2015). The CP/AP labeling is "a false dichotomy"; many real systems are just "P."[^brewer][^kleppmann]
2. **Linearizability ≠ serializability.** Linearizability is single-object, real-time recency (the CAP "C") and is composable; serializability is multi-object transaction isolation (the ACID "I") with **no real-time and not even per-process ordering**, and is not composable. The combination is **strict serializability**.[^bailis][^jepsen]
3. **ACID "consistency" is unrelated to CAP/linearizability consistency** — a frequently-confused overload of the word.[^kleppmann][^bailis]
4. **FLP "doesn't matter" in practice — but not because it's wrong.** Production systems achieve agreement with high probability; FLP only denies a *guaranteed-termination deterministic* protocol.[^flpcap]
5. **ZooKeeper is a worked counterexample**: neither CAP-consistent (default reads aren't linearizable without `sync`) nor CAP-available (writes need a majority), yet it provides atomic broadcast + causal session guarantees.[^kleppmann]

## Sources

[^gl]: Gilbert & Lynch, "Brewer's Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web Services," SIGACT News 2002 — https://www.comp.nus.edu.sg/~gilbert/pubs/BrewersConjecture-SigAct.pdf — *paper*
[^brewer]: Brewer, "CAP Twelve Years Later: How the 'Rules' Have Changed," IEEE Computer 2012 — https://www.infoq.com/articles/cap-twelve-years-later-how-the-rules-have-changed/ — *paper*
[^kleppmann]: Kleppmann, "Please stop calling databases CP or AP," 2015 — https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html — *blog*; "A Critique of the CAP Theorem," arXiv:1509.05393 — https://arxiv.org/pdf/1509.05393 — *paper*
[^pacelc]: Abadi, "Consistency Tradeoffs in Modern Distributed Database System Design," IEEE Computer 2012 — https://www.cs.umd.edu/~abadi/papers/abadi-pacelc.pdf — *paper*
[^provingpacelc]: Golab et al., "Proving PACELC," SIGACT News 2018 — https://uwaterloo.ca/distributed-algorithms-systems-lab/sites/default/files/uploads/files/proving_pacelc.pdf — *paper*
[^flp]: Fischer, Lynch, Paterson, "Impossibility of Distributed Consensus with One Faulty Process," JACM 1985 — https://www.cs.princeton.edu/courses/archive/spring24/cos418/papers/flp.pdf — *paper*
[^randconsensus]: Aspnes, "Randomized Consensus" survey (ETH) — https://disco.ethz.ch/courses/ws0304/seminar/papers/randomized_consensus_survey.pdf — *paper*
[^flpcap]: Dinh, "History of the Impossibles — CAP and FLP" — https://dinhtta.github.io/flpcap/ — *blog*; Cornell CS6410 "FLP Impossibility & Weakest Failure Detector" — https://www.cs.cornell.edu/courses/cs6410/2016fa/slides/18-distributed-systems-flp.pdf — *docs*
[^lamport]: Lamport, "Time, Clocks, and the Ordering of Events in a Distributed System," CACM 1978 — https://lamport.azurewebsites.net/pubs/time-clocks.pdf — *paper*
[^vc]: Mattern, "Virtual Time and Global States" (logical-time notes, ETH) — https://vs.inf.ethz.ch/publ/papers/logtime.pdf — *paper*; Sookocheff, "Vector Clocks" (cites Fidge 1988) — https://sookocheff.com/post/time/vector-clocks/ — *blog*
[^schneider]: Schneider, "Implementing Fault-Tolerant Services Using the State Machine Approach," ACM Computing Surveys 1990 — https://www.cs.cornell.edu/fbs/publications/SMSurvey.pdf — *paper*
[^ct]: Chandra & Toueg, "Unreliable Failure Detectors for Reliable Distributed Systems," JACM 1996 — https://www.cs.princeton.edu/courses/archive/fall08/cos597B/papers/unreliable.pdf — *paper*
[^tob]: Défago, Schiper, Urbán, "Total Order Broadcast and Multicast Algorithms: Taxonomy and Survey," ACM Computing Surveys — https://course.ece.cmu.edu/~ece845/sp11/docs/atomic-broadcast.pdf — *paper*
[^ethzquorum]: Wattenhofer et al. (ETH Distributed Systems), "Quorum Systems," Ch. 21 — https://disco.ethz.ch/courses/hs19/distsys/lnotes/chapter21-2on1.pdf — *docs*
[^aspnesquorum]: Aspnes (Yale), "Quorum Systems," PineWiki — https://www.cs.yale.edu/homes/aspnes/pinewiki/QuorumSystems.html — *docs*
[^ddiarepl]: Kleppmann, *Designing Data-Intensive Applications*, Ch. Replication & Sharding — https://www.oreilly.com/library/view/designing-data-intensive-applications/9781098119058/ — *book*
[^fpaxos]: Howard, Malkhi, Spiegelman, "Flexible Paxos: Quorum Intersection Revisited," OPODIS 2016 — https://fpaxos.github.io/ — *paper/docs*
[^failures]: Blanton (Buffalo CSE586), "Failures" (crash/omission/arbitrary taxonomy) — https://cse.buffalo.edu/~eblanton/course/cse586/2026-Spring/07-failures.pdf — *docs*; UIUC CS425 "Failure Models" — https://courses.grainger.illinois.edu/cs425/fa2009/L11tmp.pdf — *docs*
[^failstop]: Schlichting & Schneider, "Fail-Stop Processors: An Approach to Designing Fault-Tolerant Computing Systems," ACM TOCS 1983 — https://pages.cs.wisc.edu/~remzi/Classes/739/Fall2018/Papers/fail-stop.pdf — *paper*
[^alpern]: Alpern & Schneider, "Defining Liveness," IPL 1985 — https://www.cs.cornell.edu/fbs/publications/DefLiveness.pdf — *paper*
[^jepsen]: Jepsen, "Consistency Models" — https://jepsen.io/consistency/models — *docs*
[^hw]: Herlihy & Wing, "Linearizability: A Correctness Condition for Concurrent Objects," ACM TOPLAS 1990 — https://cs.brown.edu/people/mph/HerlihyW90/p463-herlihy.pdf — *paper*
[^attiya]: Attiya & Welch, "Sequential Consistency versus Linearizability," ACM TOCS 1994 — https://dl.acm.org/doi/10.1145/176575.176576 — *paper*
[^cops]: Lloyd, Freedman, Kaminsky, Andersen, "Don't Settle for Eventual… (COPS)," SOSP 2011 — https://www.cs.princeton.edu/~wlloyd/papers/cops-sosp11.pdf — *paper*
[^mad]: Mahajan, Alvisi, Dahlin, "Consistency, Availability, and Convergence" — https://www.cs.cornell.edu/lorenzo/papers/cac-tr.pdf — *paper*
[^session]: Terry et al., "Session Guarantees for Weakly Consistent Replicated Data" (Bayou), PDIS 1994 — https://www.cs.utexas.edu/~dahlin/Classes/GradOS/papers/SessionGuaranteesPDIS.pdf — *paper*
[^calm]: Hellerstein & Alvaro, "Keeping CALM: When Distributed Consistency Is Easy," CACM 2020 — https://cacm.acm.org/research/keeping-calm/ — *paper*; arXiv:1901.01930 — https://arxiv.org/abs/1901.01930 — *paper*
[^bailis]: Bailis, "Linearizability versus Serializability" — https://www.bailis.org/blog/linearizability-versus-serializability/ — *blog*
