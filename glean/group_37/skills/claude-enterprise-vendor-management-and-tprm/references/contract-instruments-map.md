<!-- Provenance: reference under the `enterprise-vendor-management-and-tprm` skill. Created 2026-06-18 via /dr deep-research. Buyer-side governance lens for a B2B/SaaS seller (MongoDB TAM) selling into banks. Educational context, not legal/compliance advice. -->

# Enterprise Commercial Contract Instruments — High-Level Map

> **Scope guard (read first).** This reference **maps** the contract-instrument landscape — what each
> document IS and where it sits in the stack. **Contract-clause DRAFTING — indemnity wording,
> limitation-of-liability construction, precedence-clause text, force-majeure language — is OUT OF SCOPE and
> belongs to the `legal-adjacent-writing` skill** (a reference under the `career-and-formal-writing` hub).
> Everything below defines instruments and their place; it never tells you how to write a clause. This is
> educational context, not legal advice.

Claims dated **as of 2026** where volatile; cited inline as `[^n]`.

## Contents
- MSA — the durable legal framework
- SOW — project-specific execution detail
- Order form / ordering document
- DPA — and how it differs from SCCs
- NDA and BAA (adjacent instruments)
- The document stack & order of precedence
- Vendor-paper vs customer-paper (banks impose their own)
- Disconfirming nuance
- References

---

## MSA — the durable legal framework

A **Master Service(s) Agreement** (sometimes Master *Subscription* Agreement in SaaS) is the **parent record**
that sets the long-term "rules of the game": confidentiality, indemnification, limitation of liability, IP
ownership, governing law, dispute resolution, termination.[^msa1] SOWs and order forms sit underneath it as
child records; **one MSA can have many.** Enterprises prefer the **MSA + downstream orders** structure
specifically to **avoid renegotiating foundational terms for every purchase** — the MSA is negotiated once,
and subsequent transactions (upsells, renewals, new products) execute through simpler documents that reference
it, "dramatically reducing negotiation time and contract friction."[^msa2]

## SOW — project-specific execution detail

A **Statement of Work** defines a particular engagement *under* the MSA — best suited to custom / project-based
work (configuration, integration, professional services) requiring detailed **scope, deliverables, milestones,
and acceptance criteria.** The SOW carries the "what / when / how-much" for that engagement; the MSA carries
the legal terms.[^sow1]

## Order form / ordering document

In SaaS/subscription deals, the **order form** is the **commercial record**: it specifies **product tier,
user/seat count, contract term, and price**, referencing the MSA for all other (legal) terms. For largely
"off-the-shelf" SaaS, an order form alone is often enough to onboard or expand — quicker than an SOW. The clean
split: **order form = "what's being bought, the price, the term"; MSA = "the legal backbone (risk allocation,
liability, IP, governing law)."**[^order1]

## DPA — and how it differs from SCCs

A **Data Processing Agreement (DPA)** is required under **GDPR Art. 28** whenever the vendor processes personal
data as a **"processor"** for the customer **"controller."** Art. 28(3) mandates a binding contract covering
minimum areas (subject-matter, duration, nature/purpose, data types, controller obligations); Art. 28(4)
requires those terms to **flow down to any sub-processors.** The DPA is typically incorporated as an exhibit to
the MSA. Most vendor DPAs use **general sub-processor authorization** — the vendor publishes a sub-processor
list and gives notice of changes with an objection window (commonly ~30 days).[^dpa1]

**Keep the two SCC families straight:**[^scc1]
- **Transfer SCCs** — the **Standard Contractual Clauses** that satisfy GDPR Chapter V for transfers of
  personal data *outside the EEA*. These address *cross-border transfer*, not the controller-processor
  relationship.
- **Art. 28(7) controller-processor SCCs** — an EU-approved *template for the Art. 28 contract itself*
  (Commission Implementing Decision (EU) 2021/915, adopted 4 June 2021), separate from the transfer SCCs.

So a DPA governs the *processing relationship*; SCCs (transfer flavor) govern *international transfer*. They
coexist and are frequently bundled, but they are distinct instruments. *(The two-SCC-family distinction is
`QUALIFIED` — handle with care.)*[^scc1]

## NDA and BAA (adjacent instruments)

- **NDA (Non-Disclosure Agreement)** — general confidentiality protection (trade secrets, business info),
  typically signed *before* substantive negotiation/evaluation. It *precedes* the deal.[^nda1]
- **BAA (Business Associate Agreement)** — the **HIPAA-specific** analog to the DPA, protecting Protected
  Health Information (PHI). It carries "the full weight of HIPAA regulations and penalties" and fires *only*
  when HIPAA PHI is in play — adjacent to most bank deals, relevant to healthcare.[^baa1]

NDA / DPA / BAA coexist as **parallel** privacy/confidentiality instruments rather than a sequence.[^baa1]

## The document stack & order of precedence

A common physical structure: **MSA at the base**, with **SLA, security addendum, and DPA incorporated as
exhibits/addenda**, and **order form / SOW** carrying the deal-specific commercial terms.[^stack1] An SLA
"incorporated by reference into the MSA" is enforceable as a critical exhibit.[^stack1]

**Precedence clauses decide which document wins on conflict — and the default is usually "MSA controls."** The
typical rule: MSA terms prevail over a conflicting SOW *unless* the SOW expressly states it is overriding a
specific MSA section. A common carve-out: an **order form's terms take precedence over the MSA, but only as to
the specific products/services that order form governs.**[^precedence1] Sellers deliberately keep critical
risk areas (warranties, limitation of liability, termination) controlled by the MSA itself.[^precedence1]

## Vendor-paper vs customer-paper (banks impose their own)

The party whose template forms the **base draft** controls the default risk allocation ("vendor-paper" vs
"customer-paper"). **In enterprise/bank deals the customer usually imposes its own paper.** Bank-side
literature is consistent that banks/financial institutions run their own professionally-negotiated terms, SLAs,
and financial-penalty schedules, and "discuss contract requirements and modifications with vendors and
negotiate provisions that facilitate effective risk management" — i.e., the bank drives from its
paper.[^paper1] *(Direct head-to-head "whose-paper-wins" sourcing is thinner; the bank-imposes-own-paper
conclusion is inferred from bank-side vendor-management literature — `QUALIFIED`, not `FACT`.)*

**Seller implication:** when selling into a bank, expect to negotiate *off the bank's MSA template*, not your
own — which front-loads the bank's audit-rights, exit, sub-processor, and resilience demands (see
`references/operational-resilience-regulation.md` for *why* those demands exist).

## Disconfirming nuance

- **MSA / order-form roles are widely confused in practice.** Teams "rush to the MSA and miss that most
  commercial disagreements start with an incomplete or vague order form" — the commercial terms (price,
  quantity, term) live in the order form, and getting them vague is a frequent failure mode even when the MSA
  is solid.[^nuance-confuse]
- **Precedence is not always "MSA wins."** Some contracts flip it so the order form (or SOW) controls *for its
  specific scope* — so "the MSA always governs" is an over-generalization; you must read the precedence
  clause.[^precedence1]
- **Terminology drift:** "MSA" is sometimes Master *Subscription* Agreement (SaaS) and sometimes Master
  *Service* Agreement (services); "order form" overlaps with "work order" / "sales order." The *function*
  (durable terms vs deal-specific commercials) matters more than the label.[^msa1]

---

## References

[^msa1]: MSA = durable parent record (confidentiality, indemnity, LoL, IP, governing law); terminology drift. PandaDoc; Ironclad; CloudNuro; OwnData — practitioner/CLM. https://www.pandadoc.com/blog/master-services-agreement-vs-statement-of-work/ ; https://ironcladapp.com/journal/contracts/what-is-an-msa ; https://www.owndata.com/legal/msa
[^msa2]: Master + downstream orders to avoid renegotiating foundational terms. Adams on Contract Drafting; Agiloft; Aline — practitioner. https://www.adamsdrafting.com/using-a-master-agreement-structure/ ; https://www.agiloft.com/blog/master-service-agreements-a-comprehensive-field-guide/
[^sow1]: SOW = project-specific scope/deliverables/milestones/acceptance under the MSA. PandaDoc; Adventum Legal; Ironclad — practitioner/law-firm. https://www.pandadoc.com/blog/master-services-agreement-vs-statement-of-work/ ; https://ironcladapp.com/journal/contracts/msas-and-sows-managing-the-contract-relationship
[^order1]: Order form = commercial record (product/quantity/price/term) under MSA. CloudNuro; AMST Legal; Promise Legal; Sirion — practitioner/law-firm. https://www.cloudnuro.ai/blog/saas-agreement-terms ; https://amstlegal.com/what-is-an-order-form-and-why-important/
[^dpa1]: DPA required under GDPR Art. 28; Art. 28(3) minimum terms; Art. 28(4) sub-processor flow-down; general authorization + ~30-day objection. GDPR-info; ComplianceStack — official-text/practitioner. https://gdpr-info.eu/art-28-gdpr/ ; https://compliancestack.ai/guides/gdpr-data-processing-agreement
[^scc1]: Two SCC families — transfer SCCs (Chapter V) vs Art. 28(7) controller-processor SCCs (Decision (EU) 2021/915). EU Commission; Bird & Bird; Fieldfisher — official/law-firm (QUALIFIED). https://commission.europa.eu/law/law-topic/data-protection/international-dimension-data-protection/new-standard-contractual-clauses-questions-and-answers-overview_en ; https://www.twobirds.com/en/insights/2021/uk/the-other-sccs-new-art-28-sccs-published
[^nda1]: NDA = general confidentiality, precedes the deal. Sirion — practitioner. https://www.sirion.ai/library/contracts/msa-vs-nda/
[^baa1]: BAA = HIPAA-specific PHI instrument; NDA/DPA/BAA parallel not sequential. Ironclad; Accountable; CommonPaper — practitioner/law-firm. https://ironcladapp.com/journal/contracts/business-associate-agreement ; https://www.accountablehq.com/post/are-business-associate-agreements-necessary-when-hipaa-requires-a-baa-and-when-it-doesn-t
[^stack1]: Document stack (MSA base; SLA/DPA/security addendum as exhibits); SLA-by-reference enforceable. CloudNuro; Ironclad; ContractNerds; Sirion — practitioner/CLM. https://www.cloudnuro.ai/blog/saas-agreement-terms ; https://ironcladapp.com/journal/contracts/your-complete-guide-to-service-level-agreements-slas
[^precedence1]: Order-of-precedence defaults (often MSA-wins, sometimes order-form-wins for its scope); critical risk kept in MSA. Law Insider; fynk; Sirion; Macronet — practitioner/CLM. https://www.lawinsider.com/clause/order-of-precedent ; https://fynk.com/en/clauses/order-of-precedence/
[^paper1]: Banks impose their own paper/terms/penalty schedules and negotiate from it. Ncontracts; Advantage-FI; Cornerstone Advisors — practitioner/advisory (QUALIFIED). https://www.ncontracts.com/nsight-blog/interagency-guidance-on-third-party-relationships/ ; https://info.advantage-fi.com/en-us/contract-negotiations
[^nuance-confuse]: MSA/order-form roles confused; vague order forms a frequent failure mode. CloudNuro — practitioner (disconfirming). https://www.cloudnuro.ai/blog/saas-agreement-terms
