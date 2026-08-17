<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-33-prescriptive-analytics` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-33-prescriptive-analytics
title: Prescriptive Analytics and Optimization (Decision Science)
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
related_skills:
  - da-15-forecasting
  - da-7-machine-learning
  - da-6-statistical-modeling
  - da-12-ab-testing-causal-inference
  - da-25-bayesian-data-analysis
description: >
  Prescriptive analytics / decision science — turning predictions and data into
  the recommended action ("what should we do?"), distinct from predictive
  analytics ("what will happen?"). Covers mathematical optimization (linear
  programming / LP, mixed-integer programming / MILP, convex / QP / SOCP,
  constraint programming / CP-SAT), optimization under uncertainty (stochastic
  programming with recourse, robust and distributionally-robust optimization,
  chance constraints), multi-objective optimization (Pareto fronts, weighted-sum,
  epsilon-constraint, NSGA-II), decision analysis (decision trees, EMV, EVPI/EVSI,
  utility theory, risk attitude), simulation for decisions (discrete-event
  simulation, Monte Carlo, simulation-optimization, queueing theory), decision
  intelligence (Gartner DI / decision modeling), and OR application archetypes
  (routing/VRP, scheduling, assignment, knapsack, inventory/newsvendor, blending,
  facility location). Tooling: OR-Tools (CP-SAT, routing), PuLP, Pyomo, Gurobi,
  CVXPY, SciPy.optimize (linprog/milp), pymoo, SimPy.
  TRIGGER: user wants to RECOMMEND/CHOOSE an optimal action or allocation, asks
  "what should we do", needs to maximize/minimize an objective subject to
  constraints, mentions LP/MILP/ILP, integer/binary decision variables,
  scheduling/routing/assignment/knapsack/blending/inventory/staffing
  optimization, Pareto/multi-objective tradeoffs, optimization under uncertainty
  (stochastic/robust/chance constraints), decision trees / EVPI / utility /
  expected value of information, discrete-event simulation or queueing for
  capacity/staffing decisions, decision intelligence, or which solver/modeling
  library to use (OR-Tools, CP-SAT, PuLP, Pyomo, Gurobi, CVXPY, scipy.optimize,
  pymoo, SimPy).
  SKIP: forecasting a future value with no decision/action attached
  (da-15-forecasting); fitting a predictive model (da-7-machine-learning) or a
  statistical model (da-6-statistical-modeling); estimating a causal effect
  (da-12-ab-testing-causal-inference); pure EDA/visualization (da-5, da-8);
  generic hyperparameter tuning of an ML model (that is model search, not a
  business decision-optimization problem) unless framed as a constrained
  resource-allocation decision.
triggers:
  - prescriptive analytics
  - decision science
  - what should we do
  - optimization model
  - linear programming
  - LP
  - mixed-integer programming
  - MILP
  - integer programming
  - binary decision variable
  - convex optimization
  - quadratic programming
  - constraint programming
  - CP-SAT
  - stochastic programming
  - recourse
  - robust optimization
  - chance constraint
  - distributionally robust
  - multi-objective optimization
  - Pareto front
  - epsilon-constraint
  - weighted sum
  - NSGA-II
  - decision tree analysis
  - expected monetary value
  - EVPI
  - EVSI
  - utility theory
  - decision intelligence
  - discrete-event simulation
  - Monte Carlo decision
  - simulation optimization
  - queueing theory
  - vehicle routing
  - VRP
  - scheduling optimization
  - assignment problem
  - knapsack
  - inventory optimization
  - newsvendor
  - blending problem
  - facility location
  - OR-Tools
  - PuLP
  - Pyomo
  - Gurobi
  - CVXPY
  - scipy.optimize
  - linprog
  - milp
  - pymoo
  - SimPy
---

# Prescriptive Analytics and Optimization (Decision Science)

## Overview

Prescriptive analytics is the fourth and highest rung of Gartner's analytics
maturity ladder (descriptive → diagnostic → predictive → prescriptive). It
answers **"what should be done?"** rather than "what happened?" or "what will
happen?" by recommending (or automating) a specific action. Gartner defines it
as advanced analytics that examines data to answer "what should be done?" using
techniques such as **optimization, simulation, complex event processing, graph
analysis, heuristics, recommendation engines, and machine learning**
([Gartner, Data & Analytics](https://www.gartner.com/en/topics/data-and-analytics)).

The mental model: **predictive feeds prescriptive.** A demand forecast (da-15)
or a propensity model (da-7) produces *parameters*; prescriptive analytics wraps
those parameters in a **decision model** — an objective to optimize, decision
variables you control, and constraints you must respect — and returns the action
that best trades off the objective against the constraints. This skill is the
**optimization + decision-science** layer of the curriculum; da-6/da-7 supply
the predictions it consumes.

A useful framing is the **decision = objective + decision variables +
constraints + uncertainty** quadruple. Choosing a method is mostly about which
of those four is hard: linear and continuous → LP; discrete choices → MILP/CP;
nonlinear-but-convex → convex/QP; uncertainty dominates → stochastic/robust; many
competing objectives → multi-objective; analytical model intractable → simulation.

## Core Concepts

### 1. Predictive → prescriptive distinction
Descriptive/diagnostic give hindsight; predictive/prescriptive give foresight,
and **human involvement decreases** as you move toward prescriptive (which can
drive automated action). Prescriptive consumes a prediction and adds a *decision
rule or optimization* on top
([EAG, 4 types of analytics](https://eaginc.com/understanding-data-analytics/);
[Qlik, descriptive→prescriptive](https://www.qlik.com/blog/embrace-the-future-moving-from-descriptive-to-prescriptive-analytics);
[Gartner glossary](https://www.gartner.com/en/topics/data-and-analytics)).

### 2. Linear programming (LP)
Continuous variables, linear objective and constraints. Solved at a polytope
vertex by simplex or interior-point. The canonical teaching cases are
**blending** (min-cost mix meeting nutritional/chemical specs) and
**product-mix** (max profit subject to resource limits). LP is the substrate
everything else extends ([SciPy linprog](https://docs.scipy.org/doc/scipy/reference/generated/scipy.optimize.linprog.html);
[PuLP docs](https://coin-or.github.io/pulp/);
[Real Python LP](https://realpython.com/linear-programming-python/)).

### 3. Mixed-integer programming (MILP / ILP)
Some/all variables are integer or **binary** (yes/no decisions: open a facility,
assign a job, select an item). NP-hard; solved by **branch-and-bound /
branch-and-cut** with LP relaxations. Binary variables unlock assignment,
knapsack, facility location, scheduling, and routing. For pure integer problems
OR-Tools recommends **CP-SAT**; for mixed continuous+integer it recommends SCIP
or a commercial solver ([OR-Tools MIP](https://developers.google.com/optimization/mip);
[Gurobi MIP basics](https://www.gurobi.com/resources/);
[SciPy milp](https://docs.scipy.org/doc/scipy/reference/generated/scipy.optimize.milp.html)).

### 4. Convex optimization (QP / SOCP / SDP)
Nonlinear but **convex** objective/constraints → any local optimum is global and
solvers are reliable and fast. Includes least-squares, quadratic programming
(e.g., **Markowitz portfolio**), second-order cone, and semidefinite programs.
**Disciplined Convex Programming (DCP)** is the rule system CVXPY uses to *verify*
a model is convex before solving — build expressions from a library of functions
with known curvature ([CVXPY DCP](https://www.cvxpy.org/tutorial/dcp/index.html);
[CVXPY intro](https://www.cvxpy.org/tutorial/intro/index.html);
[Boyd & Vandenberghe, Convex Optimization](https://web.stanford.edu/~boyd/cvxbook/)).

### 5. Constraint programming (CP / CP-SAT)
Declarative: state variables, domains, and combinatorial **constraints**
(`AllDifferent`, no-overlap, cumulative); the solver searches via propagation +
SAT/backtracking. Excels at feasibility-heavy, highly combinatorial problems —
**scheduling, rostering, timetabling**. OR-Tools **CP-SAT** is the flagship and
has repeatedly won the MiniZinc Challenge
([OR-Tools CP](https://developers.google.com/optimization/cp);
[OR-Tools CP-SAT solver](https://developers.google.com/optimization/cp/cp_solver);
[OR-Tools, Wikipedia](https://en.wikipedia.org/wiki/OR-Tools)).

### 6. Optimization under uncertainty: stochastic & robust
- **Stochastic programming with recourse**: split into first-stage (here-and-now,
  before uncertainty) and second-stage **recourse** (wait-and-see corrective)
  decisions; optimize *expected* cost over scenarios. Classic example:
  newsvendor / two-stage capacity-then-adjust.
- **Robust optimization**: optimize the **worst case** over an uncertainty set
  (no probability distribution needed) — more conservative.
- **Distributionally robust (DRO)**: hedge over a *set of distributions*; sits
  between stochastic and robust.
- **Chance constraints**: require a constraint to hold with probability ≥ 1−α
  ([NEOS Guide, Stochastic Programming](https://neos-guide.org/guide/types/stochastic/);
  [SIAM J. Optimization, Distributionally Robust Two-Stage SP](https://epubs.siam.org/doi/10.1137/20M1370227);
  [Birge & Louveaux, Intro to Stochastic Programming](https://link.springer.com/book/10.1007/978-1-4614-0237-4)).

### 7. Multi-objective optimization (Pareto)
Competing objectives (cost vs. service, risk vs. return) have no single optimum
but a **Pareto front** of non-dominated tradeoffs. Scalarization methods:
**weighted-sum** (simple, but misses non-convex regions of the front) and
**epsilon-constraint** (optimize one objective, bound the others — recovers
non-convex fronts). Population methods like **NSGA-II** (non-dominated sorting +
crowding distance) approximate the whole front in one run; `pymoo` implements
both scalarization and evolutionary approaches
([pymoo](https://pymoo.org/);
[pymoo NSGA-II](https://pymoo.org/algorithms/moo/nsga2.html);
[Blank & Deb, pymoo paper](https://arxiv.org/pdf/2002.04504)).

### 8. Decision analysis (trees, EVPI, utility)
For discrete decisions under uncertainty with few alternatives:
- **Decision trees** alternate decision nodes and chance nodes; fold back by
  **expected monetary value (EMV)** to pick the best branch.
- **EVPI** = (expected value *with* perfect information) − (best EMV *without*) —
  the max you'd rationally pay for perfect info; **EVSI** is its sample-info
  analogue (imperfect info / Bayesian update).
- **Utility theory**: replace dollars with a **utility function** to encode risk
  attitude (concave = risk-averse) — maximize *expected utility*, not EMV
  ([Wikipedia, EVPI](https://en.wikipedia.org/wiki/Expected_value_of_perfect_information);
  [Analytica, EVI/EVPI/ESVI](https://docs.analytica.com/index.php/Expected_value_of_information_--_EVI,_EVPI,_and_ESVI);
  [TreeAge, EVPI](https://www.treeage.com/help/Content/31-Analyzing-Decision-Trees/9-Expected-value-perfect-information-EVPI.htm)).

### 9. Simulation for decisions
When the system is too complex for a closed-form model:
- **Discrete-event simulation (DES)**: model entities flowing through
  resources/queues over event-driven time (`SimPy` in Python; Arena/AnyLogic
  commercially) — staffing, throughput, capacity decisions.
- **Monte Carlo**: propagate input distributions through a model to get an
  output *distribution* and risk metrics (P10/P50/P90) for a decision.
- **Simulation-optimization**: wrap a simulation as the objective for an
  optimizer (heuristics, surrogate/Bayesian opt) when no analytic form exists.
- **Queueing theory** (M/M/1, M/M/c, Little's Law `L = λW`) gives analytic
  baselines for waiting-line/staffing decisions
  ([SimPy / DES with SimPy, TDS](https://towardsdatascience.com/object-oriented-discrete-event-simulation-with-simpy-53ad82f5f6e2/);
  [Practical optimization by OR and simulation, ScienceDirect](https://www.sciencedirect.com/science/article/abs/pii/S1569190X08000439);
  [SimLLM (DES + queueing), arXiv 2026](https://arxiv.org/html/2601.06543v1)).

### 10. Decision intelligence (DI)
Gartner's operationalizing umbrella: a discipline that **explicitly models
decisions as reusable assets**, linking data → analytics → action and closing the
loop with outcome feedback. **Decision Intelligence Platforms (DIPs)** compose
data, analytics, decision modeling, and AI to support/augment/automate decisions.
Per the 2024 Gartner Market Guide, ~33% of surveyed organizations had already
deployed DI ([Gartner DI glossary](https://www.gartner.com/en/information-technology/glossary/decision-intelligence);
[Gartner, Market Guide for DI Platforms](https://www.gartner.com/en/documents/5599159);
[FICO, what is DI software](https://www.fico.com/blogs/what-decision-intelligence-software-and-should-you-invest-it)).

## Tools / Frameworks

| Tool | Layer | Best for | Notes |
| --- | --- | --- | --- |
| **SciPy.optimize** (`linprog`, `milp`) | low-level | small LP/MILP, embedding in NumPy pipelines | `milp` (HiGHS) accepts ≤, ≥, =; `linprog` is ≤ only ([SciPy milp](https://docs.scipy.org/doc/scipy/reference/generated/scipy.optimize.milp.html)) |
| **PuLP** | modeling | quick LP/MILP, teaching, readable API | writes LP/MPS; calls CBC, GLPK, HiGHS, CPLEX, Gurobi, OR-Tools ([PuLP](https://coin-or.github.io/pulp/)) |
| **Pyomo** | modeling | large/structured LP, MILP, NLP, MINLP, stochastic | algebraic; NEOS + many solvers; supports nonlinear ([Pyomo docs](https://www.pyomo.org/documentation)) |
| **Google OR-Tools** | modeling+solver | CP-SAT scheduling, **vehicle routing**, assignment | CP-SAT is best-in-class for combinatorial/ILP ([OR-Tools](https://developers.google.com/optimization)) |
| **CVXPY** | modeling | convex (QP/SOCP/SDP), DCP-verified, portfolio | also DQCP, MICP, geometric programming ([CVXPY](https://www.cvxpy.org/)) |
| **Gurobi / CPLEX** | solver | large commercial LP/MILP/QP, speed | licensed (free academic); industry standard ([Gurobi](https://www.gurobi.com/)) |
| **pymoo** | framework | multi-objective / Pareto, NSGA-II, evolutionary | scalarization + EA + visualization ([pymoo](https://pymoo.org/)) |
| **SimPy** | simulation | discrete-event / queueing in pure Python | process-based DES; pair with Monte Carlo ([SimPy](https://simpy.readthedocs.io/)) |

**Selection heuristic:** convex/continuous nonlinear → CVXPY; combinatorial /
scheduling / routing → OR-Tools CP-SAT; plain LP/MILP prototyping → PuLP; large
or nonlinear/stochastic algebraic models → Pyomo; performance at scale → hand the
model to Gurobi; many objectives → pymoo; no analytic model → SimPy + an outer
optimizer.

## Methodology

1. **Frame the decision, not the prediction.** Name the *objective* (what you
   maximize/minimize, in one unit), the *decision variables* (what you control),
   the *constraints* (what limits you), and the *uncertainty* (what you don't
   know). If you can't write these four, it isn't yet an optimization problem.
2. **Classify the problem** → pick the method family (LP / MILP / convex / CP /
   stochastic / robust / multi-objective / simulation) using the quadruple above.
3. **Source parameters** from predictive/statistical models (da-6/da-7/da-15) and
   data; keep the parameter pipeline separate from the decision model.
4. **Build small, validate, scale.** Prototype in PuLP/CVXPY/SciPy on a tiny
   instance; check feasibility and sanity of the dual/shadow prices before
   scaling to a production solver.
5. **Quantify the value of certainty** before buying data: compute **EVPI/EVSI**
   for decision-analysis problems; run **sensitivity / shadow-price** analysis on
   LP/MILP to see which constraints bind.
6. **Stress-test under uncertainty.** Re-solve across scenarios (stochastic) or
   over an uncertainty set (robust); report the *distribution* of outcomes via
   Monte Carlo, not a single point.
7. **Close the loop (DI).** Deploy the recommendation, capture realized outcomes,
   and feed them back to refit parameters and re-tune the model.

## Practical Patterns

- **Blending / diet**: LP, min cost s.t. composition specs → PuLP or SciPy.
- **Product mix / capacity**: LP, max margin s.t. resource limits; read shadow
  prices to find the bottleneck resource.
- **Assignment / matching**: binary MILP (Hungarian for the pure case) → OR-Tools.
- **Knapsack / selection**: binary MILP, max value s.t. budget → CP-SAT.
- **Scheduling / rostering**: CP-SAT with `NoOverlap` / cumulative + interval
  vars — usually beats raw MILP on combinatorial structure.
- **Vehicle routing (VRP/CVRP/VRPTW)**: OR-Tools routing library with capacity
  and time-window dimensions.
- **Inventory / newsvendor**: stochastic — balance overage vs. underage cost; the
  critical-ratio quantile of the demand distribution is the optimal order.
- **Portfolio**: convex QP (Markowitz mean-variance) in CVXPY; multi-objective
  (return vs. risk) yields a Pareto/efficient frontier.
- **Staffing / call center**: queueing baseline (Erlang-C) → refine with DES
  (SimPy) → simulation-optimization for shift design.

## Anti-Patterns

- **Optimizing a forecast instead of a decision.** A great prediction with no
  objective/constraints/action is still just predictive analytics (→ da-15/da-7).
- **Forcing nonlinearity into LP** (or ignoring that a model is non-convex).
  Linearize deliberately (piecewise, big-M) or move to convex/MINLP — don't
  pretend a non-convex problem is solved to global optimality by a local solver.
- **Big-M chosen too large.** Loose big-M constants wreck MILP relaxations and
  blow up solve time and numerics; pick the tightest valid bound.
- **Weighted-sum for non-convex Pareto fronts.** It silently misses entire
  regions; use epsilon-constraint or NSGA-II when the front may be non-convex.
- **Single-scenario "optimal" plans.** Deterministic optimization on a point
  forecast is brittle; use stochastic/robust or at least Monte Carlo stress.
- **Trusting a local optimum as global** on a non-convex/MINLP problem without
  saying so. Report the optimality gap; a MIP solver gives you one — use it.
- **Decision tree with made-up probabilities** and no EVPI. If the recommendation
  flips under plausible probabilities, you need more info, not more precision.
- **Ignoring solver status.** "Optimal" vs. "feasible/time-limit" vs. "infeasible"
  vs. "unbounded" are different answers — always check the status code.

## Troubleshooting

- **Model infeasible.** Relax/soften constraints (add slack with penalty),
  use solver IIS/conflict refiner (Gurobi, CP-SAT) to find the conflicting set;
  a frequent culprit is over-tight equalities or unit mismatches.
- **Unbounded.** Missing an upper bound or a sign error in the objective; add
  realistic variable bounds.
- **MILP too slow.** Tighten big-M, add valid cuts/symmetry-breaking, provide a
  warm-start incumbent, set a MIP gap tolerance, or switch a combinatorial model
  from MILP to CP-SAT.
- **CVXPY rejects the model ("not DCP").** An expression has unknown/ wrong
  curvature; rewrite using DCP-atoms (e.g., `cp.quad_form`, `cp.norm`,
  `cp.log_sum_exp`) per the DCP rules ([CVXPY DCP](https://www.cvxpy.org/tutorial/dcp/index.html)).
- **Numerical issues / ill-conditioning.** Rescale variables and coefficients to
  similar magnitudes; avoid mixing 1e-6 and 1e9 in the same matrix.
- **Simulation results too noisy.** Increase replications, use common random
  numbers across compared scenarios, and report confidence intervals — a single
  run is not an answer.
- **Stochastic model explodes in size.** Reduce scenarios via scenario reduction
  / sample average approximation (SAA) before solving.

## References

1. Gartner — Data & Analytics topic + prescriptive analytics definition, 2025. https://www.gartner.com/en/topics/data-and-analytics
2. Gartner — Decision Intelligence glossary, 2024–2025. https://www.gartner.com/en/information-technology/glossary/decision-intelligence
3. Gartner — Market Guide for Decision Intelligence Platforms, 2024. https://www.gartner.com/en/documents/5599159
4. Google OR-Tools — MIP, CP, CP-SAT, Routing docs, 2024. https://developers.google.com/optimization
5. CVXPY — DCP tutorial + intro, 2024–2025. https://www.cvxpy.org/tutorial/dcp/index.html
6. Boyd & Vandenberghe — Convex Optimization (Cambridge, 2004, repr.). https://web.stanford.edu/~boyd/cvxbook/
7. Pyomo — official documentation. https://www.pyomo.org/documentation
8. PuLP — COIN-OR docs, 2024. https://coin-or.github.io/pulp/
9. SciPy — `linprog` / `milp` (HiGHS) reference, v1.17, 2025. https://docs.scipy.org/doc/scipy/reference/generated/scipy.optimize.milp.html
10. Gurobi — Mathematical Optimization resources, 2024–2025. https://www.gurobi.com/resources/
11. NEOS Guide — Stochastic Programming, 2024. https://neos-guide.org/guide/types/stochastic/
12. SIAM J. on Optimization — Distributionally Robust Two-Stage Stochastic Programming. https://epubs.siam.org/doi/10.1137/20M1370227
13. Birge & Louveaux — Introduction to Stochastic Programming (Springer, 2nd ed., 2011). https://link.springer.com/book/10.1007/978-1-4614-0237-4
14. pymoo — NSGA-II + framework docs; Blank & Deb, 2020. https://pymoo.org/ , https://arxiv.org/pdf/2002.04504
15. Wikipedia / Analytica / TreeAge — EVPI, EVSI, decision-tree analysis, 2024. https://en.wikipedia.org/wiki/Expected_value_of_perfect_information
16. Discrete-event simulation with SimPy (TDS) + simulation-optimization (ScienceDirect) + SimLLM (arXiv 2026). https://towardsdatascience.com/object-oriented-discrete-event-simulation-with-simpy-53ad82f5f6e2/ , https://arxiv.org/html/2601.06543v1
17. Qlik / EAG — descriptive→prescriptive maturity, 2024. https://www.qlik.com/blog/embrace-the-future-moving-from-descriptive-to-prescriptive-analytics
