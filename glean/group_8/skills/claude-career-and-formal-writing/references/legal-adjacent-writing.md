<!-- hub-reference-banner -->
> **Reference file — part of the `career-and-formal-writing` hub.** Formerly the standalone `legal-adjacent-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: legal-adjacent-writing
description: "Craft drafting for legal-adjacent prose that will be reviewed by counsel — disclaimers, MSAs, NDAs, terms of service, privacy notices, security-incident disclosures (8-K cyber, GDPR Article 33, US state breach-notification statutes), warranty language, limitation-of-liability clauses, indemnification clauses with carve-outs, force-majeure clauses, the 'AS IS' disclaimer pattern, GDPR/CCPA/CPRA disclosure conventions, and SOC2/HIPAA notice patterns. Covers the architecture of risk-allocating contract clauses, the four-business-day SEC 8-K Item 1.05 cyber disclosure clock, the 72-hour GDPR notification window, and the language patterns that hold up in adversarial reading. TRIGGER: user asks to draft, edit, or critique a disclaimer, MSA, NDA, terms of service, privacy notice, breach-notification email, warranty, limitation-of-liability clause, indemnification clause, force-majeure clause, AS-IS disclaimer, 8-K cyber disclosure, GDPR Article 33 notice, CCPA/CPRA disclosure, SOC2/HIPAA notice; user mentions 'legal review', 'goes to counsel', 'redline', 'risk allocation', or any contract-clause name. SKIP: actual legal advice or jurisdictional interpretation (always recommend counsel); plain-language rewrites for accessibility (use plain-language); policy and governance writing that is internally facing and not contractually binding (use policy-and-governance-writing); UI microcopy that happens to include a 'Terms apply' line (use microcopy-and-ui-writing)."
origin: local
version: "1.0.0"
updated: "2026-05-29"
keywords:
  - legal-adjacent writing
  - disclaimer
  - MSA
  - master service agreement
  - NDA
  - terms of service
  - privacy notice
  - breach notification
  - 8-K cyber disclosure
  - GDPR Article 33
  - limitation of liability
  - indemnification
  - force majeure
  - AS IS disclaimer
  - warranty disclaimer
  - CCPA
  - CPRA
  - SOC2 notice
  - HIPAA notice
tags:
  - writing
  - legal-writing
  - compliance
  - contracts
  - privacy
  - security-disclosure
whenToUse:
  - User asks to draft, edit, or critique a contract clause (limitation of liability, indemnification, force majeure, warranty, IP, confidentiality)
  - User is preparing a draft that will go to in-house or outside counsel for review
  - User must publish a breach notification under GDPR Article 33, a US state breach-notification statute, or HIPAA
  - User must draft a Form 8-K Item 1.05 disclosure for a material cybersecurity incident
  - User is writing a public-facing privacy notice under GDPR, CCPA, CPRA, UK DPA 2018, or similar
  - User is writing customer-facing terms of service or an end-user license agreement (EULA)
  - User is writing an NDA (mutual or one-way) or confidentiality clause inside a larger agreement
  - User is writing a SOC2, HIPAA, or ISO 27001-related notice, attestation, or customer-facing security statement
  - User needs to draft an "AS IS" warranty disclaimer for a SaaS, open-source, or API product
whenNotToUse:
  - User is asking for legal advice or jurisdictional interpretation — always recommend qualified counsel
  - User wants UI microcopy that happens to mention legal terms (use microcopy-and-ui-writing)
  - User wants plain-language rewriting of a legal doc for accessibility (use plain-language)
  - User wants an internal policy or governance document with no contractual or regulatory binding force (use policy-and-governance-writing)
  - User wants academic writing about law, history of law, or comparative law
related_skills:
  - plain-language
  - policy-and-governance-writing
  - writing-expert
  - http-security-headers
  - mongodb-security-architecture
---

# Legal-Adjacent Writing

Reference for drafting legal-adjacent prose that will go to counsel: contracts, disclaimers, privacy notices, breach disclosures, and regulator-facing statements. This skill is **craft for drafts**, not legal advice. Every output should carry a "counsel must review before execution" footer when produced for the user.

## When to use this skill

Activate when the user needs to:

- Draft or edit a contract clause that allocates risk (LoL, indemnification, warranty, IP, confidentiality, force majeure, termination, governing law)
- Draft a privacy notice, terms of service, EULA, or DPA (data processing addendum)
- Draft a breach-notification email, regulator filing, or 8-K Item 1.05 disclosure
- Draft an "AS IS" disclaimer for a SaaS product, open-source release, or API
- Draft a SOC2, HIPAA, ISO 27001, or similar customer-facing security attestation
- Critique an existing clause for clarity, internal consistency, or alignment with a known market-standard pattern

## When NOT to use this skill

- User is asking what the law says — recommend qualified counsel
- User wants to interpret a specific clause's enforceability in a specific jurisdiction — recommend counsel
- User wants a microcopy line ("By signing up you agree…") — use microcopy-and-ui-writing
- User wants an internal-only policy document — use policy-and-governance-writing

## The five-point legal-adjacent writing test

Every legal-adjacent draft should pass these five tests before going to counsel. The skill's job is to get the draft to "ready for counsel review" — not to replace counsel.

1. **Is the risk-allocating verb correct?** "Shall," "will," "must," and "may" are not synonyms. "Shall" creates an obligation. "May" creates a right or permission. "Will" is forward-looking and sometimes ambiguous. The Plain Writing trend in legal drafting now prefers "must" over "shall" because "shall" has been litigated into ambiguity (it has been read as both mandatory and permissive in different cases). Pick one register per document and stay consistent.

2. **Are the defined terms actually defined?** Every Capitalized Term should appear once in a Definitions section, indented and alphabetical, and every subsequent use should match exactly. "Confidential Information," "Services," "Effective Date," "Term," "Party," and "Affiliate" are the most-abused. If a term is used three times without definition, define it. If it is used once, lowercase it.

3. **Does the carve-out language survive a hostile read?** A limitation-of-liability clause that caps liability "except for breach of confidentiality" is read narrowly by some courts. "Except for" should be paired with a non-exhaustive list ("including but not limited to") only when you want breadth. Pair with an exhaustive list ("the following exceptions and no others") when you want narrowness. Pick deliberately.

4. **Is the temporal scope explicit?** "In the 12 months preceding the event giving rise to the claim" is unambiguous. "In the prior year" is ambiguous (calendar year? rolling year? contract year?). Every duration, deadline, and lookback window must specify (a) the start trigger, (b) the unit (business days vs calendar days), and (c) the end trigger.

5. **Is the notice-and-cure mechanism workable?** If a clause requires "written notice," specify the delivery channel (email, certified mail, registered post, the contract's Notices section), the recipient (named individual, role, or address), and the cure window. Vague notice clauses are the most common source of post-execution disputes.

## Core concepts

### 1. The risk-allocation architecture

A commercial contract is, structurally, a series of nested risk allocations. The skill of legal-adjacent drafting is mapping each risk to the right clause type:

| Risk | Clause type | What it does |
|---|---|---|
| Product fails to do what you said | Warranty | Promises the product works; allocates remedy for breach |
| Product does what you said but causes harm | Limitation of liability | Caps the dollar exposure |
| You get sued by a third party because of your product | Indemnification | Other party covers the defense and judgment |
| You can't perform because of an Act of God | Force majeure | Excuses performance without penalty |
| Either party wants out | Termination | Defines exit ramps and notice |
| Something goes wrong with confidential data | Confidentiality + data security | Defines what's protected, for how long, and what happens on breach |
| The deal falls apart in court | Governing law + venue + dispute resolution | Defines where the fight happens |

Drafting from this map prevents the most common error: stuffing risk allocation into the wrong clause (e.g., trying to cap IP-infringement liability in the warranty section instead of carving it out of the LoL cap).

### 2. The limitation-of-liability triangle

LoL clauses have three dials. Every market-standard LoL sets all three:

1. **Cap amount.** Most common SaaS form: "fees paid by Customer in the 12 months preceding the event giving rise to the claim." Variants: fixed dollar cap ($1M, $5M), multiple of fees (2x, 3x), insurance-policy-limit cap. Higher-value deals trend toward higher caps. Sources note that SaaS deals around $500K typically negotiate to ~12 months of fees as the cap.

2. **Damages exclusion.** The "no indirect, incidental, special, consequential, or punitive damages, including lost profits, lost revenue, lost data, or lost business opportunity" exclusion. The 5–7 word list is market-standard. The trick is the "including but not limited to" qualifier: it expands the exclusion. Sophisticated counter-parties will push to delete it.

3. **Carve-outs.** What the cap and exclusion don't apply to. Standard market carve-outs (uncapped or super-capped):
   - Breach of confidentiality obligations
   - Breach of IP indemnification (typically uncapped — IP suits are existential)
   - Payment obligations (customer can't cap what they owe)
   - Gross negligence and willful misconduct
   - Indemnification obligations more broadly
   - Death or personal injury (legally unwaivable in many jurisdictions)

The carve-out list is where the negotiation happens. A clause that caps liability "for any reason" without carve-outs is unenforceable in many states.

### 3. Indemnification: mutual vs one-way, capped vs uncapped

Indemnification means "if you get sued because of something I did, I'll pay for it." The clause has four levers:

1. **Direction.** Mutual (both sides indemnify) or one-way (one party indemnifies the other). Modern B2B norm is mutual for everything except IP, where the vendor indemnifies the customer for vendor-IP claims and the customer indemnifies the vendor for customer-content claims.
2. **Trigger.** Third-party claims only (standard) or any claim including first-party (broader, less common). The "third-party claim" framing is more enforceable.
3. **Procedure.** Notice (prompt written), control of defense (indemnifying party usually controls), consent to settlement (typically required from indemnified party for any settlement admitting liability).
4. **Carve-outs from indemnification.** Common: claims arising from the indemnified party's modification of the product, combination with third-party tech not provided by indemnifying party, use outside the documented scope.

### 4. The "AS IS" warranty disclaimer

In US contract law, the implied warranties of merchantability and fitness for a particular purpose attach by default to most commercial sales. To disclaim them, the contract must do so **conspicuously** — meaning all caps, bold, or a heading that announces the disclaimer. The standard pattern:

```
EXCEPT AS EXPRESSLY SET FORTH IN THIS AGREEMENT, THE SERVICES ARE PROVIDED "AS IS"
AND "AS AVAILABLE," AND PROVIDER MAKES NO REPRESENTATIONS OR WARRANTIES OF ANY KIND,
WHETHER EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, INCLUDING WITHOUT LIMITATION ANY
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, TITLE, AND
NON-INFRINGEMENT. PROVIDER DOES NOT WARRANT THAT THE SERVICES WILL BE UNINTERRUPTED,
ERROR-FREE, OR FREE OF HARMFUL COMPONENTS, OR THAT ANY DATA WILL BE SECURE OR NOT
OTHERWISE LOST OR DAMAGED.
```

The all-caps formatting is not stylistic. It is a UCC § 2-316 "conspicuousness" requirement. A non-conspicuous disclaimer can be ignored.

### 5. The 8-K Item 1.05 cyber disclosure

The SEC's 2023 cybersecurity rules require public companies to file a Form 8-K within **four business days** of determining that a cybersecurity incident is **material**. The materiality determination is itself a separate disclosure obligation — the clock starts on determination, not on the incident itself.

What must be disclosed:
- The **material aspects** of the nature, scope, and timing of the incident
- The **material impact** or reasonably likely material impact on the registrant, including on financial condition and results of operations

What is **not** required (and should be omitted to preserve incident response):
- Specific technical detail about the attack vector
- Specific detail about the response plan
- Specific detail about cybersecurity systems, networks, or device configurations
- Any detail that would impede ongoing remediation

The carve-out for omitting technical detail is your friend. Use it. The drafting move is to describe **impact** (downtime, revenue effect, customer notification) rather than **mechanism** (the CVE, the malware family, the affected service). If information is not yet known at filing time, say so explicitly and amend later.

A delay is permitted only if the US Attorney General determines that immediate disclosure would pose a substantial risk to national security or public safety and notifies the Commission in writing.

### 6. GDPR Article 33 — the 72-hour clock

GDPR Article 33 requires controllers to notify the supervisory authority of a personal data breach "without undue delay and, where feasible, not later than 72 hours after having become aware of it." Key drafting points:

- The clock starts on **awareness**, not on the breach itself. "Awareness" does not require certainty. If you reasonably believe personal data has been compromised, the clock starts.
- **Phased notification is explicitly permitted.** Article 33(4) allows you to provide information "in phases without undue further delay" if not all information is available at 72 hours.
- The notification must include:
  - Nature of the breach, including categories and approximate number of data subjects and records
  - Name and contact details of the DPO or other contact point
  - Likely consequences of the breach
  - Measures taken or proposed to address the breach and mitigate adverse effects
- **High-risk breaches require notice to data subjects** "without undue delay" (Article 34), in "clear and plain language." Article 34 has a higher threshold than Article 33 — "high risk to rights and freedoms" rather than just "risk."

### 7. Privacy notice architecture (GDPR / CCPA / CPRA)

A modern privacy notice is a structured disclosure document, not marketing copy. The required components under GDPR Article 13/14:

1. Identity and contact details of the controller (and DPO if applicable)
2. Purposes of processing and legal basis for each
3. Legitimate interests pursued (if relied on)
4. Recipients or categories of recipients
5. Transfers to third countries and safeguards
6. Storage period or criteria
7. Data subject rights (access, rectification, erasure, restriction, portability, objection)
8. Right to withdraw consent
9. Right to lodge a complaint with a supervisory authority
10. Whether provision is statutory/contractual and consequences of not providing
11. Existence of automated decision-making, including profiling

CCPA/CPRA adds: categories of personal information collected, sold, or disclosed; rights to delete, correct, and opt out of sale/share; "Do Not Sell or Share My Personal Information" link; "Limit the Use of My Sensitive Personal Information" link; financial incentives disclosure if applicable.

The drafting trap: writing a single notice that tries to satisfy both regimes by union-merging the requirements. This produces a notice that is hard to read and easy to challenge. Better: use one schema (typically GDPR's, which is stricter) and add a CCPA/CPRA-specific section for California-resident-only rights.

## Templates

### Template 1: Mutual limitation of liability (market-standard SaaS)

```
LIMITATION OF LIABILITY.

(a) Excluded Damages. EXCEPT FOR EXCLUDED CLAIMS (DEFINED BELOW), IN NO EVENT WILL
EITHER PARTY BE LIABLE TO THE OTHER FOR ANY INDIRECT, INCIDENTAL, SPECIAL,
CONSEQUENTIAL, OR PUNITIVE DAMAGES, OR ANY LOSS OF PROFITS, REVENUE, DATA, OR
BUSINESS OPPORTUNITY, ARISING OUT OF OR RELATED TO THIS AGREEMENT, WHETHER IN
CONTRACT, TORT (INCLUDING NEGLIGENCE), STRICT LIABILITY, OR ANY OTHER LEGAL
THEORY, AND WHETHER OR NOT THE PARTY HAS BEEN ADVISED OF THE POSSIBILITY OF SUCH
DAMAGES.

(b) Liability Cap. EXCEPT FOR EXCLUDED CLAIMS, EACH PARTY'S TOTAL CUMULATIVE
LIABILITY ARISING OUT OF OR RELATED TO THIS AGREEMENT WILL NOT EXCEED THE TOTAL
FEES PAID OR PAYABLE BY CUSTOMER TO PROVIDER UNDER THIS AGREEMENT IN THE TWELVE
(12) MONTHS PRECEDING THE EVENT GIVING RISE TO THE CLAIM.

(c) Excluded Claims. "Excluded Claims" means: (i) either party's indemnification
obligations under Section [X]; (ii) either party's breach of its confidentiality
obligations under Section [Y]; (iii) Customer's payment obligations; (iv) either
party's gross negligence, willful misconduct, or fraud; and (v) any liability
that cannot be limited under applicable law.

(d) Basis of the Bargain. The parties acknowledge that the limitations in this
Section are an essential element of the bargain between them and that the fees
reflect this allocation of risk.
```

### Template 2: GDPR Article 33 supervisory authority notification

```
Subject: Personal Data Breach Notification — [Controller Name] — [Reference]

To: [Supervisory Authority Name and address]

Date and time of awareness of breach: [YYYY-MM-DD HH:MM UTC]
This notification is submitted within 72 hours of awareness pursuant to Article
33 of the GDPR / UK GDPR.

1. Controller details
   Name: [Legal entity name]
   Registered address: [Address]
   Representative (if applicable): [Name and address]
   DPO contact: [Name, email, phone]
   Reference number assigned by controller: [REF-YYYY-NNNN]

2. Nature of the breach
   Type of breach: [Confidentiality / Integrity / Availability — select all that apply]
   Date and time of the breach: [Known/estimated]
   Date and time of discovery: [Known/estimated]
   How the breach was discovered: [Brief factual description, no speculation]
   Status: [Ongoing / Contained / Resolved]

3. Affected data
   Categories of personal data: [e.g., names, email addresses, hashed passwords,
   IP addresses, partial payment-card data — list each category]
   Categories of data subjects: [e.g., end-users, employees, contractors]
   Approximate number of data subjects affected: [Number or range; state if unknown]
   Approximate number of records affected: [Number or range; state if unknown]
   Special categories of personal data involved: [Yes/No — if yes, specify]

4. Likely consequences
   [Concrete, evidence-based assessment of harm. Avoid speculation. State what
   is known, what is being investigated, and what is unknown.]

5. Measures taken or proposed
   Containment measures already taken: [List with timestamps]
   Mitigation measures planned: [List]
   Communication to data subjects: [Planned/completed/not required — with reasoning]
   Engagement with law enforcement: [Yes/No]

6. Cross-border processing
   Is this cross-border processing? [Yes/No]
   Lead supervisory authority: [If applicable]
   Other supervisory authorities notified: [List]

7. Further information
   This notification is provided based on currently available information. Where
   information was not available at the time of this notification, the controller
   will provide it in phases as it becomes available, in accordance with Article
   33(4) GDPR.

[Signature and contact details]
```

### Template 3: Form 8-K Item 1.05 cyber-incident disclosure

```
Item 1.05 — Material Cybersecurity Incidents

On [Date], [Registrant] (the "Company") determined that a cybersecurity
incident affecting [scope — e.g., a specific business segment, system, or set
of services] is material to the Company.

The Company first detected unauthorized activity on [Date]. Upon detection, the
Company [activated its incident response plan / engaged third-party cybersecurity
experts / notified law enforcement / began containment activities].

Based on the Company's investigation to date, the incident [describe nature and
scope in material aspects without disclosing specific technical detail that
could impede ongoing remediation — e.g., "involved unauthorized access to
certain Company systems and the exfiltration of data including [categories]"].

The Company has [contained / is in the process of containing] the incident.
[Describe operational impact in material terms — service disruptions, customer
notifications, recovery activities.]

The Company is continuing to investigate the incident, including the full
nature and scope of impacted data and the financial and operational impact.
The Company has not yet determined [whether the incident has had or is
reasonably likely to have a material impact on the Company's financial
condition or results of operations / the full extent of remediation costs and
related expenses / [other unknowns]]. The Company will amend this Form 8-K to
disclose this information when it becomes available.

The Company is providing this disclosure within four business days after
determining that the incident is material, in accordance with Item 1.05 of
Form 8-K.

Forward-Looking Statements
[Standard forward-looking statements disclaimer covering the cyber incident
disclosure.]
```

### Template 4: Mutual NDA (one-page, for early-stage discussions)

```
MUTUAL NONDISCLOSURE AGREEMENT

This Mutual Nondisclosure Agreement ("Agreement") is entered into as of [Date]
("Effective Date") between [Party A], a [State] [entity type], and [Party B],
a [State] [entity type] (each a "Party" and together the "Parties").

1. Purpose. The Parties wish to explore a potential business relationship
   ("Purpose") and, in connection with the Purpose, each Party may disclose
   Confidential Information to the other.

2. Confidential Information. "Confidential Information" means any non-public
   information disclosed by one Party (the "Discloser") to the other (the
   "Recipient"), whether orally, visually, in writing, or in any other form,
   that is identified as confidential at the time of disclosure or that a
   reasonable person would understand to be confidential given the nature of
   the information and the circumstances of disclosure.

3. Exclusions. Confidential Information does not include information that:
   (a) is or becomes publicly known through no breach of this Agreement by
   Recipient; (b) was rightfully known to Recipient without restriction before
   disclosure; (c) is rightfully received from a third party without
   restriction; or (d) is independently developed by Recipient without use of
   or reference to the Discloser's Confidential Information.

4. Obligations. Recipient will: (a) use Confidential Information solely for the
   Purpose; (b) protect Confidential Information using the same degree of care
   it uses to protect its own confidential information of like importance, and
   in no event less than reasonable care; and (c) not disclose Confidential
   Information to any third party except to its employees, contractors, and
   advisors who have a need to know for the Purpose and who are bound by
   confidentiality obligations no less protective than those in this Agreement.

5. Compelled Disclosure. If Recipient is compelled by law to disclose
   Confidential Information, Recipient will give the Discloser prompt prior
   notice (to the extent legally permitted) and reasonable assistance, at
   Discloser's expense, in seeking a protective order or other remedy.

6. Term. This Agreement begins on the Effective Date and continues for two (2)
   years. Recipient's obligations with respect to Confidential Information that
   constitutes a trade secret continue for as long as the information remains a
   trade secret under applicable law.

7. No License. Nothing in this Agreement grants any license or other right in
   or to any Confidential Information except the limited right to use it for
   the Purpose.

8. No Obligation. Nothing in this Agreement obligates either Party to enter
   into any further agreement or business relationship.

9. Remedies. The Parties agree that monetary damages may be inadequate for a
   breach of this Agreement and that the non-breaching Party is entitled to
   seek equitable relief, including injunctive relief, in addition to any
   other remedies.

10. Governing Law. This Agreement is governed by the laws of [State], without
    regard to its conflict-of-laws principles.

[Signature blocks]
```

### Template 5: Force majeure clause (modern pandemic-aware)

```
Force Majeure. Neither party will be liable for any delay or failure to perform
its obligations under this Agreement (other than payment obligations) to the
extent that the delay or failure is caused by an event or circumstance beyond
the reasonable control of the affected party, including acts of God, fire,
flood, earthquake, severe weather, war (declared or undeclared), armed conflict,
acts of terrorism, civil disturbance, pandemic, epidemic, public-health
emergency, governmental order or restriction (including quarantine and
shelter-in-place orders), labor strike not specific to the affected party, and
failure of public utilities or third-party telecommunications networks
("Force Majeure Event"). The affected party will: (a) give prompt written
notice to the other party describing the Force Majeure Event and its expected
duration; (b) use commercially reasonable efforts to mitigate the impact and
resume performance; and (c) keep the other party informed of material changes
in the status of the Force Majeure Event. If a Force Majeure Event continues
for more than thirty (30) consecutive days, either party may terminate the
affected Order Form or Statement of Work upon written notice without liability.
```

### Template 6: Conspicuous warranty disclaimer

```
DISCLAIMER OF WARRANTIES.

EXCEPT AS EXPRESSLY SET FORTH IN SECTION [X] (LIMITED WARRANTY), THE SERVICES
ARE PROVIDED "AS IS" AND "AS AVAILABLE." TO THE MAXIMUM EXTENT PERMITTED BY
APPLICABLE LAW, PROVIDER DISCLAIMS ALL WARRANTIES, REPRESENTATIONS, AND
CONDITIONS, WHETHER EXPRESS, IMPLIED, STATUTORY, OR ARISING FROM COURSE OF
DEALING OR USAGE OF TRADE, INCLUDING WITHOUT LIMITATION ANY WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, TITLE, ACCURACY OF DATA, AND
NON-INFRINGEMENT. PROVIDER DOES NOT WARRANT THAT THE SERVICES WILL MEET
CUSTOMER'S REQUIREMENTS OR EXPECTATIONS, OPERATE UNINTERRUPTED OR ERROR-FREE,
BE FREE FROM HARMFUL COMPONENTS, OR THAT ANY DATA WILL BE COMPLETELY SECURE OR
NEVER LOST OR DAMAGED. CUSTOMER ACKNOWLEDGES THAT THE SERVICES ARE NOT DESIGNED
FOR USE IN HIGH-RISK ACTIVITIES WHERE FAILURE COULD LEAD TO DEATH, PERSONAL
INJURY, OR ENVIRONMENTAL DAMAGE, AND PROVIDER DISCLAIMS ANY EXPRESS OR IMPLIED
WARRANTY OF FITNESS FOR SUCH PURPOSES.
```

### Template 7: Customer breach-notification email (post-detection)

```
Subject: Important Security Notice Regarding Your [Service] Account

[Date]

Dear [Customer Name],

We are writing to inform you of a recent security incident that may have
affected information associated with your [Service] account. We take the
protection of your information seriously and want to provide you with the
facts of what happened, what we have done, and what you can do.

What happened
On [Date], we detected [brief factual description, e.g., "unauthorized access
to a system that stored certain customer information"]. We immediately
[contained the incident / engaged outside cybersecurity experts / notified law
enforcement] and began an investigation.

What information was involved
Based on our investigation to date, the affected information may have included
[specific categories — name, email address, hashed password, IP address, partial
payment-card information, etc.]. [State explicitly what was NOT affected if you
know — e.g., "Full payment-card numbers, Social Security numbers, and account
passwords in usable form were not affected."]

What we are doing
[Concrete list: password resets, additional monitoring, third-party credit
monitoring offered, security improvements, notifications to regulators.]

What you can do
1. [Specific user action — e.g., "Reset your password using the secure link
   below."]
2. Monitor your account activity and report anything unusual.
3. Be alert to phishing attempts that may try to use this incident to obtain
   additional information from you. We will never ask for your password by
   email.
4. [Credit monitoring sign-up, if offered.]

For more information
[FAQ link.] If you have questions, you can contact us at [dedicated incident
response email] or [phone number, staffed hours].

We sincerely regret the concern this incident may cause. Thank you for
trusting [Company].

[Signed]
[Name, role]
```

## Anti-patterns

1. **Mixing "shall" and "must" within the same document.** Pick one register and stay consistent. Modern drafting prefers "must" for clarity.

2. **The undefined "reasonable."** "Reasonable efforts," "reasonable notice," and "commercially reasonable" are all defensible — but each should be defined or paired with a measurable benchmark. "Commercially reasonable efforts" generally means more than "reasonable efforts" and less than "best efforts," but the line is litigated.

3. **The buried carve-out.** Putting an indemnification carve-out only inside the LoL section, or vice versa, creates an interpretation conflict. Carve-outs that limit a cap should live in the LoL clause. Carve-outs that limit an indemnification obligation should live in the indemnification clause. State them in both places only if you intend the same carve-out to apply to both.

4. **Non-conspicuous warranty disclaimers.** A disclaimer of implied warranties that is not in all caps, bold, or under a heading announcing the disclaimer may be ignored as non-conspicuous under UCC § 2-316.

5. **The "perpetual" confidentiality term.** Most courts will not enforce a perpetual confidentiality obligation. Use a defined term (2 years, 5 years, 7 years) with a trade-secret carve-out that survives as long as the information remains a trade secret.

6. **Promising what you can't deliver in a privacy notice.** A privacy notice that says "we will never share your data with anyone" creates a contractual representation that the company can be sued for breaching. Better: "We share your data only as described in this Notice."

7. **The all-caps wall-of-text disclaimer.** All-caps is required for conspicuousness on warranty disclaimers and is conventional for damages exclusions. But all-caps everything is hostile and harder to read. Use all-caps only where conspicuousness is legally required.

8. **The "click here to agree" terms-of-service link.** Click-through must be unambiguous to be enforceable. The button must clearly indicate assent ("I Agree" or "Create Account and Agree to Terms"), and the terms must be available before the click, not after.

9. **Stuffing technical detail into an 8-K cyber disclosure.** Item 1.05 explicitly does not require — and the rule encourages omission of — specific technical detail that would impede remediation. Describe impact, not mechanism.

10. **The privacy notice that "doubles as marketing."** A privacy notice is a disclosure document. Marketing language ("we love your data and treat it like our own") creates representations and undermines the document's clarity.

## Decision heuristics

| Situation | Choice |
|---|---|
| User pushes back on a market-standard 12-month LoL cap and wants $0 cap | Document the push; flag to counsel. A zero cap with no carve-outs is unenforceable in many states. |
| Question: should we disclose a breach now or wait for more facts? | If under GDPR Article 33 and you are "aware" of a likely-risk breach, the 72-hour clock has started. Phased notification is permitted; waiting is not. |
| Question: is this a "material" cybersecurity incident under SEC rules? | The materiality determination is the registrant's, made in consultation with counsel, on a totality-of-the-circumstances basis (qualitative + quantitative). The 4-business-day clock starts on determination, not on incident. |
| Indemnification: third-party claims only, or any claim? | Third-party claims only. This is market-standard and avoids the "indemnification swallows the LoL" trap. |
| Warranty disclaimer format | All caps. Under a heading that says DISCLAIMER. UCC § 2-316 conspicuousness. |
| Force majeure: include pandemic? | Yes — modern force majeure clauses since 2020 typically include "pandemic, epidemic, public-health emergency, governmental order including quarantine." Older form clauses that say only "Act of God" have been litigated. |
| Privacy notice: one notice for global, or per-region? | Use one notice with a global schema (typically GDPR's) and add region-specific sections for CCPA/CPRA, UK DPA, LGPD, etc. |
| NDA: mutual or one-way? | Mutual is the default for exploratory business discussions. One-way is appropriate when only one side is sharing (e.g., a vendor pitching a customer who is sharing nothing). |
| Confidentiality term length | 2–5 years for non-trade-secret information; perpetual for trade secrets (with the carve-out that the obligation continues only as long as the information remains a trade secret). |

## Cross-skill notes

- **Use plain-language alongside this skill** when drafting customer-facing privacy notices, breach-notification emails, or terms-of-service summaries. Plain-language craft (Flesch-Kincaid 8th-grade target, common-word substitution) applies; legal-adjacent craft (precise verbs, defined terms, risk allocation) constrains it.
- **Use policy-and-governance-writing** for internal policies (acceptable use, code of conduct, data classification) that are not contractually binding.
- **Use writing-expert** for the executive summary or cover memo that introduces a legal-adjacent document. The disclaimer body uses this skill; the cover memo summarizing it for an executive uses writing-expert.

## Final reminder

Every output from this skill should carry a footer:

> This is draft language only. It is not legal advice. Qualified counsel must review before execution, filing, or public release. Jurisdiction, deal-specific facts, and recent case law may change what is enforceable.

## References

1. SEC, *Cybersecurity Risk Management, Strategy, Governance, and Incident Disclosure* (Form 8-K Item 1.05 final rule, effective Dec 18, 2023): https://www.sec.gov/newsroom/press-releases/2023-139
2. SEC, *Disclosure of Cybersecurity Incidents Determined To Be Material and Other Cybersecurity Incidents* (Gerding statement, May 2024): https://www.sec.gov/newsroom/speeches-statements/gerding-cybersecurity-incidents-05212024
3. UK ICO, *Personal data breaches: A guide*: https://ico.org.uk/for-organisations/report-a-breach/personal-data-breach/
4. EDPB / GDPR Article 33 and 34 (controller notification + data subject communication).
5. UCC § 2-316 (conspicuous disclaimer of implied warranties).
6. Bloomberg Law, *Limitation of Liability Contract Provision Examples*: https://www.bloomberglaw.com/external/document/XFI6F0SC000000/commercial-drafting-guide-limitation-of-liability-contract-provi
7. ContractNerds, *Five Critical Elements of Limitation of Liability Provisions*: https://contractnerds.com/five-critical-elements-of-limitation-of-liability-provisions/
8. Aaron Hall, *Carveouts to Mutual Indemnity in SaaS Agreements*: https://aaronhall.com/carveouts-to-mutual-indemnity-in-saas-agreements/
9. IAPP resources on GDPR Article 33, CCPA/CPRA, and cross-border breach notification.
10. Practical Law / Thomson Reuters, *General Contract Clauses: Indemnification (Mutual Long Form)*.
