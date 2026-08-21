<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `debt-collectors-and-fdcpa-rights` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: debt-collectors-and-fdcpa-rights
description: 'Third-party debt collectors and your FDCPA / Regulation F rights. Legal INFORMATION, not advice (as of 2026). Spoke of consumer-credit-and-debt. TRIGGER: a collector is calling, texting, or emailing; validate/dispute a debt with the collector (FDCPA validation); make a collector stop; collector harassment or threats; being sued by a collector / answer a summons or complaint; time-barred "zombie" debt and whether to pay; does a payment restart the clock; my rights with collections; which CFPB sample letter to send. SKIP: settling, charge-off, pay-for-delete, 1099-C, re-aging -> charge-offs-collections-and-debt-resolution; any NC statute or NC garnishment -> north-carolina-credit-and-debt-law; federal statute/rule text (FCRA, FDCPA, Reg F, ECOA) -> us-consumer-credit-and-debt-law; disputing a furnisher with a bureau, or how a report/score is built -> credit-reports-and-scores; raising a score -> improving-and-rebuilding-credit; medical-collection specifics -> medical-debt-and-billing.'
metadata:
  changelog:
  - 2026-06-16 sko --meta --no-sync — structural pass vs new consumer-finance family; 0 High 0 Med (meta-validate clean). Added 1 SKIP edge (-> medical-debt-and-billing); trimmed description to 997 chars within 1000 cap; no content passes.
  - '2026-06-16 sko create->v1.0.0 — 2 iterations, CLEAN; fixed 3 High (description 1184->993 chars; sample-letter numbering aligned to CFPB filenames; letter-5 download added + landing-page ref corrected) + 6 Medium (validation written/electronic; em-dash density 1.49->0.98/100; 2 trigger-eval false-positives closed; reference headings ###; letter-4 URL); Pass H predicted 10/10 pos, 10/10 neg'
---

# Debt collectors and your FDCPA rights

> **This is general legal information, not legal advice.** It summarizes the
> federal Fair Debt Collection Practices Act (FDCPA) and its implementing
> rule, **Regulation F (12 CFR part 1006)**, **as of 2026**. Laws and amounts
> change, and **state law often adds protections** (and some states regulate
> the original creditor — see the NC skill). **Deadlines when you are sued are
> strict and set by state court rules.** If a collector sues you, or you are
> unsure, **consult a licensed attorney or your local legal aid office
> promptly** — do not rely on this page alone. Verify current law via the
> primary sources at the end.

**Parent hub:** spoke of `consumer-credit-and-debt`. This skill is the
*collector-interaction* spoke — what to say, what to send, and what to do when
a collector contacts you or sues. For statute text see
`us-consumer-credit-and-debt-law`; for settling the debt see
`charge-offs-collections-and-debt-resolution`; for North Carolina specifics see
`north-carolina-credit-and-debt-law`.

---

## Who the FDCPA covers (and who it doesn't)

The federal FDCPA applies to **third-party debt collectors** collecting debts
that are **primarily personal, family, or household** (not business debts):

- **Collection agencies** collecting for someone else.
- **Debt buyers**, companies that bought the debt and collect it as their own.
- **Collection law firms / attorneys** regularly collecting debts.

**The federal FDCPA generally does NOT cover the original creditor** collecting
its own debt in its own name (e.g., your bank or card issuer's in-house
collections). It also excludes most business/commercial debts.

> **State-law flag:** Several states extend collection rules to **original
> creditors** too. **North Carolina** is one (its Collection Agency Act and the
> Prohibited Practices by Collectors statute reach in-house/first-party
> conduct). **Defer NC specifics to `north-carolina-credit-and-debt-law`.** If
> the caller is the original creditor and you're outside the FDCPA, check your
> state's law.

**If you don't recognize the debt at all,** treat it as possible scam or
identity theft until validated; see "Validate first" below and NC DOJ's
scam-collector guidance.

---

## Your rights — quick checklist

- [ ] **Validation notice.** An FDCPA collector must give you the validation
  information **in the initial communication or within 5 days** of it. The
  **information must be provided in writing or electronically** (a written
  validation notice), even if the first contact was a phone call — an oral-only
  statement does not satisfy the rule. Required details are below. *(Reg F
  § 1006.34)*
- [ ] **30-day dispute / verification right.** You can **dispute the debt in
  writing within 30 days** of getting the validation notice. If you do, the
  collector **must pause collection of the disputed amount until it sends you
  verification.** *(Reg F § 1006.34)*
- [ ] **Request the original-creditor name.** You can ask (within 30 days) for
  the name and address of the original creditor.
- [ ] **Limits on when/where they contact you.** No contact at times/places
  they know are inconvenient — **presumed inconvenient = before 8:00 a.m. or
  after 9:00 p.m. your local time.** No contact **at work if they know your
  employer prohibits it.** *(Reg F § 1006.6)*
- [ ] **Restrict the channel/time.** You can tell them how to reach you and to
  stop using a particular medium (e.g., "no calls at work," "email only").
- [ ] **Talk to your lawyer instead.** If you're represented and they can reach
  your attorney, they generally must contact the **attorney**, not you.
- [ ] **Cease all communication.** A **written** "stop contacting me" request
  means they must stop (narrow exceptions below). *(Reg F § 1006.6)*
- [ ] **Call-frequency cap.** Presumed-OK is **no more than 7 calls in 7
  consecutive days** about the same debt, and **no call within 7 days after a
  phone conversation** about it. More than that is presumed harassment. *(Reg F
  § 1006.14)*
- [ ] **No harassment.** No threats of violence, obscene language, repeated
  calls to annoy, publishing debtor lists, or hiding their identity. *(§ 1006.14)*
- [ ] **No false or misleading statements.** They can't misstate the amount,
  falsely claim to be an attorney/government, falsely threaten arrest, or
  threaten action they can't or won't take. *(§ 1006.18)*
- [ ] **Limited reach to third parties.** They generally can't tell your
  family, friends, neighbors, or employer that you owe a debt; contacts to
  *locate* you are tightly limited. *(§ 1006.6)*
- [ ] **Complain.** You can file a complaint with the **CFPB**
  (consumerfinance.gov/complaint), the **FTC** (reportfraud.ftc.gov), and your
  **state AG** (in NC, the NCDOJ Consumer Protection Division). FDCPA also lets
  you sue a violating collector.

### What must be in the validation notice (Reg F § 1006.34)

1. A statement that the communication is from a debt collector.
2. The collector's and your **name and mailing information**.
3. The **name of the creditor** you owe (the current creditor; original too if
   different) and the **account number**, if any.
4. An **itemization of the current amount**, showing interest, fees, payments,
   and credits since an itemization date.
5. The **current amount** of the debt as of the notice date.
6. **How to reply** (e.g., "this isn't mine," "the amount is wrong").
7. An **end date for the 30-day dispute period**.
8. For consumer-finance debts, a pointer that more help is at
   `cfpb.gov/debt-collection`.

> A bare phone demand with none of this is a red flag — ask for it in writing
> before you discuss or pay anything.

---

## Interaction playbook (do this, in order)

1. **Don't confirm or pay anything on the first call.** Don't give bank or card
   numbers, your SSN, or "yes I'll pay" language until the debt is validated —
   especially if it might be **time-barred** (a payment can restart the clock,
   below).
2. **Get it in writing. Communicate in writing.** Ask for the validation
   notice. From here on, **prefer mail/email over phone** so there's a record.
   For mailed disputes/cease letters, consider **certified mail, return receipt
   requested**, and keep a copy.
3. **Keep records.** Log every call (date, time, caller, what was said); save
   every letter, voicemail, and text. This evidence is what wins an FDCPA
   complaint or lawsuit, and what proves you disputed on time.
4. **Validate / dispute within 30 days.** If anything is off — wrong person,
   wrong amount, already paid, too old — **send a written dispute within the
   30-day window.** That forces the collector to **stop collecting the disputed
   amount until it verifies.** (You can still dispute after 30 days, but you
   lose the automatic stop-and-verify.)
5. **Set contact rules or cut off contact.** Use the right CFPB sample letter
   (next section): restrict the channel/time, or stop contact entirely.
   - **Cease-communication consequence:** a stop-contact letter ends the
     calls/letters, **but it does not erase the debt and does not stop a
     lawsuit.** Cutting off contact can make a collector *more* likely to sue,
     because that's the only remaining avenue. Decide accordingly.
6. **Confirm it's not a scam.** Legitimate collectors will provide validation
   in writing. **Pressure to pay immediately by gift card / wire / crypto,
   refusal to identify themselves, or threats of arrest/deportation = scam or
   illegal.** Hang up and report it (NC DOJ: 1-877-5-NO-SCAM; FTC).
7. **If the debt may be time-barred, get advice before paying.** See below.

---

## CFPB sample letters — which one to send

The CFPB publishes **five free, fill-in-the-blank "action" letters** for
replying to a collector. Pick by situation:

| # | Letter (CFPB name) | Use it when… | Effect |
|---|---|---|---|
| 1 | **"I do not owe this debt"** (dispute + verification) | You believe it isn't yours or the amount is wrong | Disputes the debt and tells them to **stop collecting until they prove** you owe it |
| 2 | **"I need more information"** | You're not sure the debt is yours/valid | States you **dispute** it until they answer specific questions about what's owed |
| 3 | **"Stop contacting me"** (cease communication) | You want all contact to stop | Federal law requires them to **stop after a written request** (narrow exceptions) |
| 4 | **"Contact only my lawyer"** | You've hired an attorney | Gives the attorney's info; collector contacts the **lawyer**, not you |
| 5 | **"Limit how/when you contact me"** (contact restriction) | Calls come at bad times/places | Tells them your **preferred method/time**; bars inconvenient contact |

The numbering follows the CFPB's own letter set (letter 1 = "not my debt"). All
five are downloadable; the direct links are in **References**.

Best practice: **fill in the date and your reference info, keep a copy, and
send by certified mail (return receipt) or a tracked channel** so you can prove
delivery and timing, which is critical for the 30-day dispute and for any later
dispute about whether you sent a cease letter. Download links are in
**References** below.

---

## Time-barred ("zombie") debt — handle with care

A debt is **time-barred** when the **statute of limitations (SOL)** to *sue*
on it has expired. Key points:

- **The SOL varies by state and debt type** — commonly **3–6 years**, sometimes
  longer, and the choice-of-law clause in your contract can matter. *(Federal
  student loans generally have no SOL.)* **Look up your state's SOL** — for NC,
  see `north-carolina-credit-and-debt-law`.
- **It doesn't delete the debt.** The debt still exists; collectors may still
  **ask** you to pay, and (depending on your state) may still **contact** you.
- **They can't sue (or threaten to sue) on a time-barred debt.** If they do,
  that itself can be an FDCPA violation, and the **SOL is an affirmative
  defense** you raise in court.
- **DANGER — restarting the clock.** In many states, **making a payment (even a
  small one) or acknowledging/promising in writing that you owe** an old debt
  can **revive the debt and restart the SOL**, exposing you to a fresh lawsuit
  for the full balance plus interest and fees.
- **So:** before you pay, settle, or sign *anything* on an old debt, **find out
  whether it's time-barred and whether a payment would restart the clock** —
  and consider getting legal advice first. Collectors must make certain
  time-barred disclosures, but **don't rely on them to tell you the clock will
  reset.**

> Related but different: settling/negotiating a charge-off or collection (and
> the 1099-C tax angle) lives in
> `charge-offs-collections-and-debt-resolution`.

---

## "You've been sued by a collector" — step list

> **The single most important rule: do NOT ignore a lawsuit or summons.**
> Failing to respond by the deadline lets the court enter a **default judgment**
> for the full claimed amount plus fees and interest — and a judgment is very
> hard to undo. You have a far better chance fighting *before* judgment.

1. **Read the papers and find the deadline.** The summons/complaint states **how
   many days you have to respond** (often ~20–30, but **state-specific**) and
   *where*. Calendar it immediately.
2. **Get help fast.** Contact a **licensed attorney or legal aid**; many help
   with FDCPA/collection defense and some are free or low-cost. Do this early —
   deadlines are unforgiving.
3. **File a written Answer with the court by the deadline.** Respond to each
   allegation (admit/deny/lack-knowledge). Filing an Answer is what prevents the
   default judgment. **Don't just call the collector — file with the court.**
4. **Make them prove standing / ownership.** Especially against **debt buyers**,
   demand proof they **own this debt** and that the **amount is correct** (chain
   of assignment, the account contract, a complete itemization). Many cases
   stall when the plaintiff can't produce documentation.
5. **Raise affirmative defenses.** Common ones: **statute of limitations**
   (time-barred), wrong defendant / identity theft, already paid/settled, no
   proof of ownership. **You must assert SOL — the court won't apply it for
   you.** Counterclaims for FDCPA violations may be available.
6. **Show up to every hearing.** Bring your records.
7. **Know what a judgment enables** (so you understand the stakes): **wage
   garnishment, bank-account levy/garnishment, and property liens.** *State law
   limits these and provides exemptions —* and some states **sharply restrict
   wage garnishment for ordinary consumer debts** (North Carolina is notably
   protective). **Defer the NC garnishment/exemption detail to
   `north-carolina-credit-and-debt-law`.**

---

## When to route elsewhere

- **Negotiating a settlement, pay-for-delete, lump-sum vs. payment plan,
  charge-off mechanics, or 1099-C cancellation-of-debt income** →
  `charge-offs-collections-and-debt-resolution`.
- **North Carolina statutes** (NC Collection Agency Act, Prohibited Practices
  by Collectors, NC SOL, NC wage-garnishment limits, NCDOJ enforcement) →
  `north-carolina-credit-and-debt-law`.
- **Full federal statute/regulation text and enforcement structure** (FCRA,
  FDCPA, every Reg F section, ECOA/TILA) → `us-consumer-credit-and-debt-law`.
- **Removing a collection from your report / improving your score** →
  `improving-and-rebuilding-credit`; **how reporting & scoring work** →
  `credit-reports-and-scores`.

---

## References / verify current law (as of 2026)

Laws and dollar/time thresholds change — **confirm against these primary
sources** before acting, and get an attorney for anything litigation-related.

### CFPB Regulation F (12 CFR part 1006), the operative rules
- Validation notice contents & 30-day dispute, § 1006.34:
  https://www.consumerfinance.gov/rules-policy/regulations/1006/34/
- Communications, inconvenient times, cease-communication, § 1006.6:
  https://www.consumerfinance.gov/rules-policy/regulations/1006/6/
- Harassment/abuse & the 7-calls-in-7-days presumption, § 1006.14:
  https://www.consumerfinance.gov/rules-policy/regulations/1006/14/
- Full Part 1006 (eCFR): https://www.ecfr.gov/current/title-12/chapter-X/part-1006

### CFPB consumer guidance
- Debt collection resource hub: https://www.consumerfinance.gov/consumer-tools/debt-collection/
- What to do when a collector contacts you:
  https://www.consumerfinance.gov/ask-cfpb/what-should-i-do-when-a-debt-collector-contacts-me-en-1695/
- What info a collector must give you:
  https://www.consumerfinance.gov/ask-cfpb/what-information-does-a-debt-collector-have-to-give-me-about-the-debt-en-331/
- How to get a collector to stop contacting you:
  https://www.consumerfinance.gov/ask-cfpb/how-do-i-get-a-debt-collector-to-stop-contacting-me-en-1411/
- Statute of limitations / time-barred debt:
  https://www.consumerfinance.gov/ask-cfpb/what-is-a-statute-of-limitations-on-a-debt-en-1389/
- If you're sued by a collector or creditor:
  https://www.consumerfinance.gov/ask-cfpb/what-should-i-do-if-im-sued-by-a-debt-collector-or-creditor-en-334/
- What a judgment is:
  https://www.consumerfinance.gov/ask-cfpb/what-is-a-judgment-en-1381/

### CFPB sample letters (download .doc)
All five letters and CFPB's usage instructions are on the "what should I do when
a debt collector contacts me" page (listed under consumer guidance above).
- 1, "Not my debt" (dispute + verification):
  https://files.consumerfinance.gov/f/documents/cfpb_debt-collection-letter_1-not-my-debt.doc
- 2, "More information":
  https://files.consumerfinance.gov/f/documents/cfpb_debt-collection-letter-2_more-information.doc
- 3, "Stop contacting me":
  https://files.consumerfinance.gov/f/documents/cfpb_debt-collection-letter-3_stop-contacting.doc
- 4, "Contact my lawyer":
  https://files.consumerfinance.gov/f/documents/cfpb_debt-collection-letter-4_contact-my-lawyer.doc
- 5, "Here's how to contact me" (contact restriction):
  https://files.consumerfinance.gov/f/documents/cfpb_debt-collection-letter-5_heres-how-to-contact-me.doc

### FTC
- Debt Collection FAQs (rights, time-barred debt, being sued):
  https://consumer.ftc.gov/articles/debt-collection-faqs-0
- FDCPA statute (legal library):
  https://www.ftc.gov/legal-library/browse/statutes/fair-debt-collection-practices-act

### North Carolina (defer NC detail to the NC skill)

- NCDOJ, Debt Collectors:
  https://ncdoj.gov/protecting-consumers/credit-and-debt/debt-collectors/
- NCDOJ, Watch out for debt collection scams:
  https://ncdoj.gov/watch-out-for-debt-collection-scams/
- NCDOJ, File a complaint (1-877-5-NO-SCAM):
  https://ncdoj.gov/file-a-complaint/
