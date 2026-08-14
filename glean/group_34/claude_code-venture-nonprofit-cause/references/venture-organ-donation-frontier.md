<!-- hub-reference-banner -->
> **Reference file — part of the `venture-nonprofit-cause` hub.** Formerly the standalone `venture-organ-donation-frontier` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-organ-donation-frontier
description: >-
  Cutting-edge, advocacy, and adjacent-domain companion for a North Carolina organ/eye/tissue donation nonprofit founder — the frontier beyond the core donation system. Covers (1) xenotransplantation & bioengineered organs (gene-edited pig-organ trials through 2025-26, the FDA IND/BLA pathway, hype vs reality); (2) donation policy & legislation tracking (OPTN Modernization, the US opt-in NOT opt-out reality, the Living Donor Protection Act, CMS OPO rule litigation, and how an advocate tracks bills); (3) tissue & eye banking regulation (FDA HCT/P, AATB, EBAA, the for-profit-tissue controversy); (4) the National Donate Life Month campaign playbook. Fast-moving — verify before relying. TRIGGER: xenotransplantation/pig-organ trials & FDA status; tracking donation legislation (LDPA, OPTN Act, CMS OPO rule); opt-out/presumed-consent claims about the US; tissue/eye-banking regulation or the for-profit-tissue controversy; National Donate Life Month / April campaign planning; OPTN Modernization status; "what's new in donation." SKIP: the core donation system, registration, allocation, OPOs, consent law, NC Anatomical Gift Act, waitlist, myths → venture-organ-donation-system; campaign strategy & storytelling → venture-cause-nonprofit-marketing; Ad Grants/CRM/email ops → venture-nonprofit-fundraising-ops; multicultural/faith outreach, volunteers, story consent → venture-donor-community-engagement; 501(c)(3) filing → venture-nc-nonprofit-formation; clinical transplant medicine.
category: personal-venture
version: 1.1.0
updated: 2026-06-16
tags: [venture, organ-donation, xenotransplantation, policy-advocacy, tissue-banking]
whenToUse:
  - Explaining the state of xenotransplantation (gene-edited pig organs) and what is hype vs. real as of 2025-26
  - Getting the FDA regulatory pathway for xenotransplants right (expanded access vs. IND clinical trial vs. BLA)
  - Tracking organ-donation legislation — a specific federal bill, a state bill, or the policy landscape
  - Correcting or fact-checking claims that the US has (or is adopting) opt-out / presumed-consent donation
  - Explaining the OPTN Modernization rollout and the CMS OPO final-rule litigation to a board or partner
  - Understanding tissue and eye banking as their own regulated worlds (AATB, EBAA, FDA HCT/P) and their controversies
  - Planning April / National Donate Life Month programming and plugging a small NC org into national toolkits
  - Deciding what is stable enough to put in published messaging vs. what is too fast-moving to assert
triggers:
  - xenotransplantation / pig kidney / gene-edited organ trial / FDA xeno status
  - track organ donation legislation / Living Donor Protection Act / OPTN Act / bill status
  - opt-out or presumed consent in the US (claim to verify)
  - CMS OPO final rule / OPO decertification / lawsuit
  - tissue banking / eye banking / AATB / EBAA / HCT/P regulation
  - for-profit tissue processing controversy
  - National Donate Life Month / Blue & Green Day / April campaign toolkit
  - OPTN Modernization status / multi-vendor / independent board
---

# Organ Donation: The Frontier, Advocacy & Adjacent Worlds

The forward-looking, advocacy, and adjacent-domain layer for a North Carolina organ/eye/tissue donation nonprofit. The companion **venture-organ-donation-system** skill covers the *core* system (HRSA/OPTN/UNOS, OPOs, registration, allocation, consent law, NC Anatomical Gift Act, myths, statistics) — read that for foundations. This skill covers four things that move fast and that a founder doing advocacy and awareness needs to handle well:

1. **Xenotransplantation & bioengineered organs** — where the science actually is.
2. **Policy & legislation tracking** — what's pending and *how to follow it yourself*.
3. **Tissue & eye banking** — separate regulatory worlds with their own controversies.
4. **The National Donate Life Month playbook** — how a small org plugs into the April machine.

> **This is a fast-moving frontier.** Trial results, approval status, and bill progress change month to month. Every dated claim below carries an as-of date. **Re-verify before you publish, present, or fundraise on any of it.** See the Disclaimer. This is not medical or legal advice.

A founder's mental model: your levers are still **education, trust-building, and registry sign-ups** (per the core skill). The frontier topics here are mostly things you *talk about* to build credibility and hope — not things you operate. The advocacy topics are things you can actually *act on* (contact a legislator, join an action alert, sign onto a coalition letter).

---

## 1. Xenotransplantation & Bioengineered Organs

**Xenotransplantation** = transplanting living cells, tissues, or organs from one species into another — here, **gene-edited pig organs into humans**. It is the most-hyped potential answer to the organ shortage (~103,000+ on the US waitlist per the core skill). As of mid-2026 it is **early-stage clinical research**, not standard care. No xenotransplant product is FDA-approved.

### Why pigs, and what "gene-edited" means

Pigs are used for organ size/physiology compatibility and breeding practicality. The barrier is rejection and cross-species incompatibility, addressed by **CRISPR gene editing** of a source pig. Two product strategies dominate (as of 2025-26):

- **United Therapeutics "UKidney"** — a pig with **10 gene edits**: 6 human genes inserted (immune compatibility) + 4 pig genes inactivated — 3 that remove the carbohydrate/"sugar" antigens humans react to (GTKO, CMAH, B4GALNT2) plus 1 (growth-hormone receptor, GHR) to stop the organ from over-growing (United Therapeutics IR, 2025-02-03; EXPAND study, PMC 2025; NYU Langone). **Note: the UKidney does NOT inactivate the porcine retrovirus (PERV)** — PERV stays a monitored safety concern rather than an edit.
- **eGenesis "EGEN-2784"** — a separately engineered multi-gene-edited pig kidney whose signature is **inactivating the porcine endogenous retrovirus (PERV)** across the genome (in addition to antigen knockouts and human transgenes) — the main design contrast with United Therapeutics' pig (eGenesis press releases; CRISPR Medicine News, 2025).

A founder's plain-language version: scientists edit a pig's DNA so its organ looks "less foreign" to the human immune system; one program (eGenesis) also edits out a pig virus, while the other (United Therapeutics) monitors for it instead. **Don't over-claim**: editing reduces but does not eliminate rejection — every human case so far has eventually faced rejection or a related complication, and no design has been proven safe long-term.

### The FDA regulatory pathway (get this distinction right)

There are **three different legal routes**, and conflating them is the most common error:

| Route | What it is | Status for xeno |
|---|---|---|
| **Expanded Access / "compassionate use"** | One desperately ill patient, no other options; case-by-case FDA + hospital IRB authorization. NOT a trial. | How nearly all human cases happened **2022–early 2025**. |
| **IND clinical trial** | An Investigational New Drug application clears the FDA to run a structured, enrolling study with safety endpoints. | **Now underway** — the real milestone (see below). |
| **BLA (Biologics License Application)** | Full FDA marketing approval. A xeno organ is regulated as a **biologic** by **CBER**. | **None granted.** Years away. The trials aim to *support* a future BLA. |

So "FDA cleared a trial" ≠ "FDA approved pig organs." It means an IND was allowed to proceed.

### Where it realistically stands (as of mid-2026)

**Compassionate-use / expanded-access cases (single patients, 2022–2025):**

- **Hearts (Univ. of Maryland):** Two pig-heart transplants into living patients (Jan 2022; Sept 2023). Both patients died (~60 and ~40 days), largely from antibody-mediated rejection (Scientific American; UMaryland; *Circulation*, 2024). No US pig-heart *clinical trial* is enrolling as of mid-2026 — heart xeno lags kidney. [Verify before asserting a third heart case.]
- **Kidneys:** Several expanded-access cases, including **Towana Looney** (NYU Langone, transplanted **Nov 25, 2024**), who set the **longest-survival record at ~130 days** before the kidney was removed after rejection on **April 4, 2025**, returning her to dialysis (NYU Langone News; American Kidney Fund; *Science*; NPR, 2025-04-11).
- **Liver (first-ever, living recipient):** A team in **China (Anhui Medical University, led by Beicheng Sun)** performed the **world's first pig-to-living-human (auxiliary) liver xenotransplant** — a 71-yr-old with unresectable liver cancer. No hyperacute rejection through day 31; the graft was removed on day 38 for **xenotransplantation-associated thrombotic microangiopathy (xTMA)**; the patient died day 171 (*Journal of Hepatology*, 2025; EASL; PMC). Proof-of-concept that a pig liver can perform metabolic functions, with serious limits.

**The real milestone — first regulated clinical TRIALS (IND-cleared):**

- **United Therapeutics UKidney / EXPAND study:** FDA cleared the IND **Feb 3, 2025**. Design: a "phaseless" combined phase 1/2/3, initial cohort **6 patients** across 2 centers, expandable to **up to 50**; ESRD patients ~ages 55–70 on dialysis ≥6 months; 24-week primary follow-up; 12-week safety gap between the first two transplants; aims to support a BLA (United Therapeutics IR, 2025-02-03). **First transplant in the trial performed Nov 3, 2025 at NYU Langone Health** — the first transplant in a *regulated clinical xenotransplantation study* (United Therapeutics IR, 2025-11-03; NYU Langone News). [Outcomes of trial participants are early — verify current status before citing survival.]
- **eGenesis EGEN-2784:** FDA IND cleared (2025) for a phase 1/2/3 kidney trial in dialysis-dependent ESRD patients age 50+ on the waitlist; 24-week endpoint (eGenesis; CRISPR Medicine News; Inside Precision Medicine, 2025).
- **Other INDs:** Between Feb 2024 and 2025 the FDA granted multiple xeno INDs; a **pig liver** IND was also reported (HCPLive; National Kidney Foundation). [Verify the current count/sponsors — this list grows.]

**Bottom line for messaging:** Xenotransplantation is a genuine, FDA-regulated clinical reality as of late 2025, but it is **experimental, kidney-led, measured in months-not-years of graft survival, and not approved**. Safe framing: "Promising research is underway in FDA-cleared trials; it is not yet an approved treatment, and registering as a donor today is what saves lives now." Do **not** tell the public pig organs are "available" or "approved."

### Adjacent frontier (one sentence each, [UNVERIFIED] for specifics)

- **Bioengineered / lab-grown & 3D-bioprinted organs** and **decellularized scaffolds** — earlier-stage than xeno; no human whole-organ products. [UNVERIFIED for any near-term clinical claim.]
- **Organ preservation tech** — normothermic/machine perfusion ("organ in a box," e.g., warm perfusion devices) is **already in clinical use** and expands the usable donor pool by reconditioning organs; a real, current story worth telling (cross-check current device/approval names before naming them).
- **Zoonosis / PERV risk** — the lifelong-monitoring and infectious-disease-surveillance requirement is *why* xeno trials are so cautious; it's the honest counterweight to hype.

---

## 2. Donation Policy & Legislation Tracking

This is where a founder can genuinely *act*. Two halves: **the current landscape** and **how to track it yourself**.

### 2a. Critical fact-check: the US is OPT-IN, not opt-out

> **Watch for misinformation.** Low-quality/AI-generated articles circulate claiming specific US states (e.g., "California, New York, Virginia") "adopted opt-out / presumed consent starting January 1, 2026." **This is false.** As of mid-2026 the **United States uses an opt-in, first-person-authorization system in every state** — you must affirmatively register. **No US state has presumed consent / opt-out.** Authoritative confirmations: New England Donor Services ("Why Opt-out… Is Not a Good Idea for the U.S."), Donate Life California FAQ, the National Academies, and HRSA. Opt-out is a *foreign* model (e.g., UK, Spain, much of Europe). If you see an opt-out-in-the-US claim, treat it as wrong until a primary source (a state statute or HRSA) proves otherwise.

Why this matters: a founder repeating that claim destroys credibility and can spook the public (privacy/consent fears actually *reduce* registration). Cross-reference the consent-models section in **venture-organ-donation-system** for the opt-in/opt-out debate and why the US hasn't switched.

### 2b. The federal landscape (as of mid-2026 — verify status before relying)

- **OPTN Modernization Initiative (ongoing).** HRSA is breaking the decades-long single-contractor model. Congress passed the **Securing the U.S. OPTN Act (Public Law 118-14, Sept. 2023)**, letting HRSA award **multiple contracts** for separate OPTN functions. **For the first time in ~40 years, HRSA has made multi-vendor awards** (HRSA OPTN Modernization, Sept 2024). An **independent, elected OPTN Board of Directors** — separated from the contractor — has been established. **Patient-safety and committee-support functions were being competed out of the UNOS contract in early 2026**; HRSA exercised short (3-month) continuity options to avoid a gap when the operations contract period ended **Dec 29, 2025** (HRSA OPTN Modernization updates, Nov 2025 / Jan 2026). **Takeaway:** the governance/contractor layer is mid-restructuring — say "the OPTN, overseen by HRSA," not "UNOS runs the system." (More detail in the core skill.)

- **Living Donor Protection Act (LDPA) of 2025** — the marquee living-donor bill. **S. 1552** (Sen. Gillibrand et al.) and House companions **H.R. 4582 / H.R. 4583**, introduced ~May 1, 2025 (119th Congress). It would (a) bar life/disability/long-term-care insurers from discriminating against people for donating an organ/tissue, and (b) clarify that **FMLA** covers living-donor recovery time. **Milestone: it reached/cleared Senate HELP Committee markup** (a first for the bill), reported out of committee ~Feb 2026, awaiting full-chamber action — **not yet enacted** as of this writing (Congress.gov S.1552; National Kidney Foundation; GovTrack). [Verify current stage — this is live.]

- **CMS OPO Final Rule litigation.** The **2020 CMS Final Rule** redefined OPO performance metrics and recertification tiers (donation-rate and transplantation-rate outcome measures, with a ranking that puts the bottom tier at risk of **decertification**), with consequences landing in **2026**. **Seven OPOs sued** — ***LifeLink Foundation et al. v. Kennedy et al.***, filed **Aug 1, 2025**, **M.D. Fla.** (Case No. 25-cv-02042). Core APA arguments: CMS omitted statutorily-required "process" measures, relied on flawed state **death-certificate** data, used two highly-correlated outcome measures + a confidence interval biased against larger OPOs, and built a relative-ranking system where some OPOs must always fall to the bottom (AOPO statement; Crowell & Moring; Mondaq). CMS issued **additional guidance on the organ-donation process in early 2026** and signaled decertification proceedings would follow rule finalization (Holland & Knight, Mar 2026). **Takeaway:** OPO accountability is real and contested; this is a "system is being held to outcomes, and the metrics are disputed" story — handle even-handedly. [Verify case posture and any 2026 rule changes before asserting outcomes.]

- **Other living-donor bills to know:** the **Expanding Support for Living Donors Act** (reimburses lost wages/childcare/travel for living donors) reappears across Congresses (Rep. Miller et al.; American Transplant Foundation). [Verify current bill number/status.]

### 2c. HOW to track legislation (the advocate's toolkit)

Build this into a simple recurring habit — don't rely on memory or social media:

**Federal:**
- **Congress.gov** — the authoritative source. Search a bill number (e.g., "S. 1552") or keyword ("organ donation," "living donor"); open the bill page for **status/actions, cosponsors, committees, and full text**. Create a **free Congress.gov account → save searches → email alerts** on a bill or keyword.
- **GovTrack.us** and **LegiScan** — friendlier trackers with a "% progression" view, prognosis, and email alerts; good for a quick read but **confirm specifics on Congress.gov**.
- **Federal rules (CMS/HRSA):** **Regulations.gov** for proposed/final rules and to **submit a public comment** during open comment periods (a real advocacy action a small org can take); the **Federal Register** for the rule text.

**State (NC and others):**
- The **NC General Assembly** site (ncleg.gov) for NC bills; most states have an equivalent. **LegiScan** and **Quorum/FastDemocracy**-type trackers aggregate multi-state bills if you watch a national issue.
- Watch your **state donation coalition** (Donate Life NC) and the OPOs serving NC (**HonorBridge**, **LifeShare Carolinas** — see core skill) for NC-specific legislative news.

**Curated advocacy feeds (let the experts surface what matters):**
- **Donate Life America** — national observances + a **"Take Action / Get Involved Locally"** hub and **state-contact directory** (donatelife.net). Join their list for action alerts.
- **National Kidney Foundation**, **American Transplant Foundation**, **AOPO**, **AST/ASTS**, and regional OPOs publish **legislative action alerts** (often one-click "email your representative" tools). Subscribing to 2–3 of these is the highest-leverage move for a one-person advocacy shop.

**Practical cadence:** set Congress.gov + one state tracker email alerts on your 2–3 priority bills; skim a Donate Life / NKF action alert email weekly; act (call/email/comment) only on the few moments that matter (a markup, a floor vote, an open comment period). Cross-reference **venture-cause-nonprofit-marketing** for turning an alert into a supporter campaign, and note **501(c)(3) lobbying limits** (insubstantial-lobbying / 501(h) election) — see **venture-nc-nonprofit-formation** before running anything that looks like lobbying.

---

## 3. Tissue & Eye Banking — Separate Regulated Worlds

Organ, **eye**, and **tissue** donation are one phrase but **three different regulatory systems**. Donor-registration usually authorizes all three, but recovery, processing, and oversight diverge sharply. A founder should know the map so they speak accurately and don't accidentally describe tissue like organs.

### The big structural difference

- **Organs** (heart, kidney, liver, lung, pancreas, intestine) — allocated by **OPTN/UNOS**, recovered by **OPOs**, can't be bought/sold (NOTA), transplanted quickly. (Core skill.)
- **Tissue** (skin, bone, tendons, heart valves, corneas-as-eye-tissue, veins) — recovered by tissue banks/OPOs, **processed, stored long-term, and distributed** — often through **for-profit processors** — and used in tens of thousands more procedures than organs. **Not allocated by OPTN.**
- **Eyes/corneas** — their own world with their own accreditor (EBAA). Cornea transplant (keratoplasty) is one of the most common and successful transplants.

### Who regulates what

- **FDA (CBER)** regulates tissue as **HCT/Ps — Human Cells, Tissues, and cellular/tissue-based Products** under **21 CFR Part 1271**. The key legal split:
  - **"361" HCT/Ps** (named for §361 of the Public Health Service Act): tissue that is **minimally manipulated** and used for its **same/homologous function** → regulated as tissue only; the bank must **register & list with FDA, screen/test donors, follow Good Tissue Practice, keep records, and report adverse events** — but no premarket approval.
  - **"351" products**: if cells are **more-than-minimally manipulated** or given a **non-homologous use** → regulated as a **biologic/drug** requiring a **BLA/IND**.
  - **The controversy lives in that line.** FDA historically let *manufacturers self-classify*, which enabled aggressive marketing of poorly-evidenced products (see below) (AABB; ProPublica; FDA 21 CFR 1271).
- **AATB (American Association of Tissue Banks)** — the voluntary **accreditation** body for tissue banks. Publishes the **Standards for Tissue Banking** (first 1984; accreditation program since 1986; **re-accreditation ~every 3 years**, with the possibility of unannounced inspections). The AATB standards are the most comprehensive private tissue-banking standards in the world (AATB; Wikipedia/AATB). Look for **AATB-accredited** when vetting a tissue partner.
- **EBAA (Eye Bank Association of America)** — the analogous body for **eye banks**. Publishes **EBAA Medical Standards** (set by its Medical Advisory Board with FDA input), runs **accreditation with on-site inspection ≥ every 3 years** by an eye banker + corneal surgeon, and is recognized by the CDC/federal health agencies. Eye banks must forward FDA inspection findings (e.g., **Form 483**) to EBAA within 10 business days and report serious adverse reactions to FDA within 15 days (EBAA "restoresight.org"; *Eye Banking and Corneal Transplantation* journal, 2024–25 standards).

### The for-profit tissue-processing controversy (handle carefully)

A real, recurring scandal-zone a founder should understand but discuss responsibly:

- **The profit tension.** Donated tissue is given altruistically (NOTA bars *organ* sale and limits what can be charged), but **for-profit processors** turn donated tissue into products generating substantial revenue. Critics argue a profit motive conflicts with altruistic donation and outpaces oversight; defenders note processing genuinely costs money and saves the products' usefulness (PMC, international ethics review; ICIJ).
- **Specific abuses reported over the years:** the **"birth-tissue / amniotic" profiteers** marketing unproven stem-cell-style products by exploiting the 361/351 self-classification gap (**ProPublica, "The Birth-Tissue Profiteers"**); historic **illegal body-broker/body-parts trade** cases that prompted federal legislation (Rep. Pallone et al.); and earlier safety/sterilization concerns with soft-tissue allografts (PMC).
- **Why it matters to a donation-awareness org:** families sometimes fear "my loved one's tissue will be sold for profit." The honest, careful answer: **donation itself is a gift and isn't sold; reputable recovery is done by accredited (AATB/EBAA) nonprofits and OPOs; processing fees are regulated; and the abuses that exist are at the downstream for-profit-product end, which FDA and Congress have repeatedly acted on.** Don't dismiss the fear and don't sensationalize it. Cross-reference the **myths/trust** content in **venture-organ-donation-system** and the trusted-messenger / historical-mistrust guidance in **venture-donor-community-engagement**.

---

## 4. National Donate Life Month — The Awareness-Campaign Playbook

**April is National Donate Life Month (NDLM)**, established by **Donate Life America (DLA)** in 2003. It is the single biggest annual moment for a donation-awareness org, and DLA does most of the heavy lifting — a small NC org's job is to **plug in**, not reinvent. (Campaign *mechanics* — channels, donor journey, content calendar — live in **venture-cause-nonprofit-marketing**; this section is the donation-specific calendar and assets.)

### The April calendar (2026 dates — verify on donatelife.net each year)

- **National Donate Life Month** — all of **April**.
- **2026 theme: "Leave a Legacy,"** using **trees** as the symbol (a forest of interconnected lives — donors, recipients, families) (donatelife.net; Donate Life California; Lifeline of Ohio toolkit, 2026).
- **Donate Life Living Donor Day** — **April 1, 2026** (honors living donors).
- **Blue & Green Spirit Week** — **April 4–11, 2026**, with themed days (appreciation, outreach, education — e.g., "Thank Your Healthcare Team," "Write a Message of Hope," "Dress Up Your Pet," "Make Blue & Green Treats").
- **National Donate Life Blue & Green Day** — **April 10, 2026** — wear **blue and green** and promote registration (donatelife.net; Donor Alliance; multiple partner sites). [One outlier source listed Apr 11; **donatelife.net says Apr 10** — verify the year you run it.]
- **National Pediatric Transplant Week** — **April 19–25, 2026** (focus: ending the pediatric waitlist).

### Other observances beyond April (round out the year)

- **National Minority Donor Awareness Month — August** (equity focus; high-value for NC's under-registered communities — pair with the disparities content in the core skill and the outreach playbook in **venture-donor-community-engagement**). [Verify exact 2026 framing.]
- **National Donor Sabbath** — mid-**November** (the Fri–Sun ~2 weeks before Thanksgiving) — faith-community engagement; pairs with the faith-leader-partnership guidance in **venture-donor-community-engagement**.
- **National DMV / first-person-authorization moments** and **Donate Life Rose Parade Float** (Jan 1) — national visibility hooks. [Verify current-year specifics.]

### The Donate Life signature assets a small org uses

- **Blue & Green** — the official Donate Life colors (use them in April collateral).
- **The Donate Life flag** — raised at hospitals/OPOs/government buildings in April; a simple, photogenic local action (ask your OPO or city about a flag-raising).
- **Donor-family & living-donor recognition** — quilts, "Legacy" tributes, honor walks, and *Donate Life* recognition events — emotionally central and partner-supported.
- **DLA partner toolkits** — **free downloadable resources**: social graphics ("square posts" sized for Facebook/Instagram/LinkedIn/X), Instagram/Facebook **Stories**, **Spirit Week** graphics, posters, and **English + Spanish** assets; sample social copy and activity ideas (donatelife.net; state affiliates like Donate Life NC, plus OPO toolkits e.g., Lifeline of Ohio, NJ Sharing Network).

### How a small NC org plugs in (practical)

1. **Affiliate / align with Donate Life NC** and your OPO (**HonorBridge**, **LifeShare Carolinas**) — use the **DLA state-contact directory** (donatelife.net "Get Involved Locally") to find the right people; co-brand rather than compete (core skill has the NC registry/OPO map).
2. **Adopt the national theme & dates** so your messaging rides the bigger wave (use "Leave a Legacy" / trees in 2026).
3. **Download and localize DLA toolkits** — don't design from scratch; add your NC org's logo and a local call-to-action (register at the NC registry / National Donate Life Registry — see core skill).
4. **Do one or two photogenic local actions** — a Blue & Green Day photo push, a flag-raising, a donor-family recognition moment — and amplify with the partner graphics.
5. **Run the year-round beats** (Minority Donor Awareness in August, Donor Sabbath in November) so April isn't your only moment.
6. **Mind compliance** — charitable-solicitation registration before fundraising (**venture-nc-nonprofit-formation**) and ad/email/deliverability operations (**venture-nonprofit-fundraising-ops**).

---

## Cross-references

- **venture-organ-donation-system** — the core US/NC donation system, OPOs, registration, allocation, consent law, NC Anatomical Gift Act, myths, statistics, the OPTN/UNOS/HRSA structure, and the opt-in/opt-out debate. **Read first for foundations.**
- **venture-cause-nonprofit-marketing** — campaign strategy and donor-journey design, cause-storytelling ethics, and turning an action alert into a supporter campaign.
- **venture-nonprofit-fundraising-ops** — the build-and-run layer: Google Ad Grants account structure, email lifecycle automation, donor CRM selection, and registry-signup conversion measurement.
- **venture-donor-community-engagement** — multicultural & faith-community outreach where registration lags, volunteer/ambassador programs, and ethical (consented) donation-story collection.
- **venture-nc-nonprofit-formation** — 501(c)(3), charitable-solicitation registration, and **lobbying limits** (501(h)) before any advocacy that looks like lobbying.

---

## Sources

Fast-moving topics; each carries an as-of date inline. Re-verify before relying.

**Xenotransplantation:**
- United Therapeutics IR — IND clearance for UKidney (2025-02-03); first transplant in EXPAND trial (2025-11-03), ir.unither.com
- NYU Langone News — Towana Looney pig-kidney case; first regulated trial transplant (2025), nyulangone.org
- *Science*, NPR (2025-04-11) — longest pig-kidney case (~130 days) failure/removal
- eGenesis press releases; CRISPR Medicine News; Inside Precision Medicine (2025) — EGEN-2784 IND
- *Journal of Hepatology* / EASL / PMC (2025) — first pig-to-living-human liver xenotransplant (China, Anhui Medical Univ.)
- Scientific American; Univ. of Maryland; *Circulation* (2024) — pig-heart cases
- National Kidney Foundation; American Kidney Fund; HCPLive — FDA xeno trial overview

**Policy & legislation:**
- HRSA OPTN Modernization updates (Sept 2024; Nov 2025; Jan 2026), hrsa.gov; Public Law 118-14 (Securing the US OPTN Act, 2023)
- Congress.gov — S. 1552 / H.R. 4582 / H.R. 4583 (Living Donor Protection Act of 2025); GovTrack; LegiScan
- National Kidney Foundation — LDPA Senate HELP markup; American Transplant Foundation — Expanding Support for Living Donors Act
- AOPO statement; Crowell & Moring; Mondaq; Holland & Knight (Mar 2026) — *LifeLink Foundation v. Kennedy* (CMS OPO Final Rule litigation, M.D. Fla. 25-cv-02042)
- New England Donor Services; Donate Life California; National Academies; HRSA — US is opt-in (NO US opt-out/presumed consent)
- Tracking tools: Congress.gov, GovTrack, LegiScan, Regulations.gov, Federal Register, ncleg.gov, Donate Life America "Take Action"

**Tissue & eye banking:**
- FDA 21 CFR Part 1271 (HCT/Ps; 361 vs 351); AABB regulatory pages
- AATB — Standards for Tissue Banking; accreditation program (aatb.org)
- EBAA — Medical Standards & accreditation (restoresight.org); *Eye Banking and Corneal Transplantation* (2024–25 standards)
- ProPublica "The Birth-Tissue Profiteers"; ICIJ; PMC international ethics review; Rep. Pallone (body-parts-trade legislation) — for-profit-tissue controversy

**Donate Life Month:**
- Donate Life America (donatelife.net) — National Observances; Blue & Green Day; NDLM; toolkits
- Donate Life California; Donor Alliance; Lifeline of Ohio; NJ Sharing Network; Donate Life NY — 2026 dates, "Leave a Legacy" theme, partner toolkits

---

## Disclaimer

This skill describes a **fast-moving frontier**. Clinical-trial results, FDA approval status, bill progress, court rulings, regulatory rules, and campaign dates **change frequently** and some details above carry an as-of date of mid-2026 or are marked **[UNVERIFIED]**. **This is educational information, not medical or legal advice.** Nothing here states that any xenotransplant product is FDA-approved (none is), that any bill has become law (verify on Congress.gov), or that the US has adopted opt-out donation (it has not). **Always re-verify trial status, approval status, legislation status, and dates against primary sources** (FDA, HRSA/OPTN, Congress.gov, court dockets, AATB/EBAA, donatelife.net) **before publishing, presenting, fundraising, or advising anyone.** For the core donation system and NC law, defer to **venture-organ-donation-system**; for nonprofit legal/lobbying questions, defer to **venture-nc-nonprofit-formation** and qualified counsel.
