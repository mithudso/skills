<!-- Provenance: reference under the `splunk-platform-spl` standalone skill. Created 2026-06-18 via /dr deep-research (multi-source web research, ≥3 independent sources/concept). Overview level; RBA covered in depth. Volatile claims stamped verified-as-of: 2026-06-18. -->

# Splunk for SIEM — Enterprise Security (ES)

`verified-as-of: 2026-06-18` for ES 8.x terminology, editions, the SOAR/UBA portfolio, and the Cisco direction (all VOLATILE — ES is on a fast 8.0→8.5 cadence; verify against the running version). This is an overview-level reference; **risk-based alerting (RBA)** gets the most depth. ES is sourced mostly to Splunk Docs + .conf security talks + Splunk Lantern.

## Contents
- What ES is (and its CIM/data-model dependency)
- Correlation searches
- Notable events / findings + the analyst triage workflow (the ES 7.x→8.x terminology shift)
- Risk-based alerting (RBA) — in depth
- Threat intel, assets & identities, MITRE ATT&CK
- The Splunk security portfolio (ES, SOAR, UBA/UEBA, Cisco)
- Disconfirming findings / gotchas

## What ES is
Splunk Enterprise Security is the **premium SIEM app that runs on top of Splunk core** (Splunk Enterprise or Splunk Cloud Platform) — not standalone software. It requires an underlying Splunk platform plus a Daily Indexing Volume or vCPU/SVC license; buying ES grants security use of ingested data but **no extra ingest capacity.**[^3][^2] It reuses Splunk's search and correlation capabilities across access, endpoint, and network domains.[^2]

ES is **functionally inseparable from the CIM and accelerated data models.** The CIM (the `Splunk_SA_CIM` add-on, bundled with ES) provides normalized field names and tags per data domain; ES dashboards and detections are **populated from CIM data-model acceleration** and query the summaries via **`tstats`.**[^5][^25][^6] ES also installs its own data models (Risk, Threat Intelligence, Incident Management).[^1] **Practical consequence:** if the right CIM data models are not accelerated (or source data is not CIM-normalized), ES dashboards and correlation searches return empty/incomplete results.[^24][^25] (See `spl-language-and-commands.md` for CIM mechanics.)

## Correlation searches
A **correlation search scans multiple normalized data sources for a defined pattern; when it matches, it performs an adaptive response action** — most commonly creating a notable event.[^8][^9] They can span any domain plus asset/identity lists and threat intel, aggregating an initial search with SPL (counts, thresholds) and acting when results meet conditions — e.g., detect brute-force by counting authentication failures above a threshold followed by a success.[^9]

Three control layers:[^7][^9]
- **Search + thresholds:** the SPL logic plus threshold conditions; admins can modify thresholds, response actions, and schedule/cron.
- **Trigger conditions:** evaluate whether returned results match before an action fires.
- **Throttling:** applied *after* trigger conditions, prevents creating more than one alert of a given type within a window (per throttle fields). **Throttling occurs before notable-event suppression.** Modifying a correlation search does **not** retroactively change existing notables.[^7]

**Out-of-the-box content:** ES ships with **most of its core correlation searches disabled by default**, so teams enable only what is relevant.[^7] The large body of prebuilt detection content comes via **Enterprise Security Content Update (ESCU)** — a free app maintained by the **Splunk Threat Research Team (STRT)** delivering "analytic stories" mapped to MITRE ATT&CK / Kill Chain / CIS, bundled with detection searches (and SOAR playbooks where available); the open-source mirror is `splunk/security_content` / research.splunk.com.[^18][^20] (Note the two are different objects: ESCU detections of the *Correlation* analytic **type** are enabled by default on a cron schedule, whereas the core ES correlation searches are mostly shipped disabled.[^22]) (ES 8.x renames "correlation search" → **detection**; see below.)

## Notable events / findings + the analyst triage workflow — VOLATILE
A **notable event is the alert object a correlation search creates** when a pattern is detected; notables are written to the dedicated **`notable` index** and surfaced on the **Incident Review dashboard.**[^23][^30] Classic triage workflow (ES ≤7.x):[^27][^30] monitor Incident Review → assign an **owner** (status New → In Progress) → research/record findings, possibly escalate to an **Investigation** → resolve (Resolved/Closed). Triage fields: **Status, Urgency, Owner, Disposition.** **Urgency is auto-calculated** by `urgency_lookup` from the notable's **severity** (set on the correlation search) combined with the **priority** of the matched asset/identity (higher of asset vs identity wins; analysts can override).[^27][^29]

### ES 7.x → 8.x terminology shift (VOLATILE — current as of ES 8.0–8.5)
ES 8.0 (announced .conf24, June 2024) embedded **Mission Control natively in ES** and renamed core objects to align with the **Open Cybersecurity Schema Framework (OCSF):**[^21][^16][^14]

| ES ≤7.3 | ES 8.x |
|---|---|
| Notable event, risk notable | **Finding** |
| Risk event | **Intermediate finding** |
| Correlation search / risk rule | **Detection** (event-based or finding-based) |
| MC incident, ES investigation | **Investigation** |
| Risk object | **Entity** |
| Incident Review page | **Analyst Queue** (under Mission Control) |

A **finding combines what were separately "notable events" and "risk events" into one object** carrying tactics, techniques, confidence, impact, risk score, threat objects, and impacted entity.[^14] Findings are triaged in the **Analyst Queue on the Mission Control page** (Owner / Status / Urgency / Disposition).[^15] **Finding groups** auto-aggregate related findings by rules (similar findings, cumulative risk score, MITRE ATT&CK thresholds).[^21][^14] **Intermediate findings** (formerly risk events) are *not* triageable in the queue — they appear nested inside finding groups.[^14] Mission Control adds native **Splunk SOAR** playbook/action integration from the queue.[^21]

## Risk-based alerting (RBA) — in depth
**What it is.** RBA **attributes incremental "risk" to entities (assets/identities) from many lower-confidence detections, accumulates that risk in a dedicated index, and only raises a notable when accumulated risk crosses a threshold** — converting a flood of discrete alerts into a small number of high-fidelity, context-rich notables.[^10][^12] Two stages:[^33][^10]
1. **Risk rules** (correlation searches with a Risk Analysis adaptive-response action) observe anomalies and **write "risk events" / "risk modifiers" to the `risk` index** instead of alerting. Each risk event carries at minimum `risk_score`, `risk_object`, and `risk_object_type`.[^11][^35]
2. **Risk Incident Rules (RIRs)** scan the risk index, **aggregate risk by object**, and create a **risk notable** when the object's summed risk crosses a threshold over a window.[^10][^33]

**Key objects:**
- **Risk object / risk_object_type** (= "Entity" in ES 8.x): the system, host, user, credential, or device a detection reports on; `risk_object` references a search field like `src`/`dest`. Only `risk_object` + `risk_object_type` are strictly required to create a risk notable.[^11]
- **Risk modifiers** are the risk-index events carrying the score; **risk factors** are *multipliers/adjusters* applied based on entity characteristics (privileged user, critical-priority asset, time-of-day, geolocation). ES 8.3 ships **seven default risk factors** (e.g., critical-priority destination ×1.5, high-priority user ×1.25). Scoring: **base score** (from the detection) → **calculated score** (after risk factors) → **total score** = sum of calculated scores for an entity in a window.[^34][^36][^12]
- **Risk index** (`risk` data model `All_Risk`): all risk events.[^24]
- **Threat objects** (domain, command line, IP, registry key, filename): best practice recommends aggregating/curating RIRs by *threat object* rather than risk object, because threat objects are stronger behavioral indicators.[^bestpractices]

**Default RIRs (load-bearing examples):**[^33]
- **"Risk Threshold Exceeded for Object Over 24 Hour Period"** — sums `calculated_risk_score` by object via `tstats … from datamodel=Risk.All_Risk`; fires when score **> 100 over 24 hours.**
- **"ATT&CK Tactic Threshold Exceeded for Object Over Previous 7 Days"** — fires when distinct MITRE **tactics exceed 3 over 7 days** for an object.

These thresholds are tunable. ES 8.x also adds an **Entity Risk Score (ERS)** — a 0–100 normalized, weighted-average score over 7 days, recalculated every 20 min, used by native UEBA.[^38] *(ERS specifics are single-source Splunk Docs — qualified.)*

**Why RBA exists — alert-fatigue reduction (the origin story).** RBA was introduced to the Splunk community by **Stuart McIntosh and Jim Apger** at .conf (2018 "Say Goodbye to Your Big Alert Pipeline…"; 2019 "Getting Started with RBA and MITRE").[^E3] The problem: the **one-to-one alert model** (one detection = one ticket) suffers **alert fatigue, doesn't scale, over-zealous exclusions, little correlation, and no alert narrative.** The **risk-based model** decouples detection count from investigation count — you can run a *large* inventory of high- and low-confidence detections, but **more detections ≠ more investigations** because conditions must aggregate past a threshold.[^E3] Splunk's tech brief frames the canonical SOC outcome: ~250–500 risk events/day collapse to ~60–100 actual tickets.[^E5] RBA began as a community add-on (`apger/SA-RBA`, now deprecated) and became **natively supported in ES 6.6.**[^E6][^10]

### RBA tuning pain / gotchas (disconfirming — the real-world traps)
- **RBA is not automatic noise reduction.** Uncurated RBA can generate as many notables as traditional alerting and produce *duplicate* risk notables; "excessive/duplicate risk notables from normal business activity confuse analysts." Tuning is mandatory.[^bestpractices]
- **Keep the threshold constant; tune the *scores* around it.** Splunk's RBA FAQ stresses scores are meant to be dynamic — lower scores for benign/regular traffic via the `risk_score` field; keep the threshold (e.g., 100) stable.[^RBA-FAQ]
- **Tuning levers:** filter known false positives *before* writing to the risk index; **throttle** on risk_object/score/message to stop duplicate notables; **dedup** similar events; **weight** noisy sources (count a noisy rule as 0.5×); identify high-benign-rate searches via a `benign_rate` calculation.[^rba-gh][^33]
- **RBA is a "journey, not one-and-done."** The hardest parts are **stakeholder/leadership buy-in, SPL familiarity, getting the SOC involved early, asset-&-identity data, and continuous tuning.** A common rollout tip: force risk notables to "Informational" so you don't wreck the SOC's metrics while tuning.[^RBA-FAQ][^37]

## Threat intel, assets & identities, MITRE ATT&CK (overview)
**Threat intelligence:** ES ingests indicators into KV-store collections; **threat-matching searches monitor the CIM data models** for those entities and enrich investigations. Admins add feeds (download, custom CSV, or from events) with a **Weight** that raises the risk score of assets/identities interacting with high-confidence indicators. Bundled sources include Emerging Threats, PhishTank, the SANS blocklist; generic enrichment lists include MITRE ATT&CK. **VOLATILE:** ES 8.x has two parallel systems — the native **Threat Intelligence Framework** (KV-store) and **Threat Intelligence Management (Cloud / TIM)**; post-Cisco, **Cisco Talos** intel is being integrated.[^intel][^available-intel]

**Assets & identities:** when correlation is on, ES compares indexed events against merged asset/identity lists at search time via automatic lookups — matching `src`/`dest`/`dvc` (asset) and `user`/`src_user` (identity). Matches add fields, feed **urgency** (asset/identity priority drives the urgency calc), and feed **risk factors.** Entity zones (`cim_entity_zone`) disambiguate the same IP/name across zones.[^A&I-manage][^A&I-config]

**MITRE ATT&CK mapping:** ES detections carry **annotations** mapping them to frameworks (managed: MITRE ATT&CK, Kill Chain, CIS, NIST, PCI + Confidence/Impact). Annotations live in `savedsearches.conf` as `action.correlationsearch.annotations` JSON; the `mitre_attack_enrichment` automatic lookup expands a technique ID into tactic/technique context. ATT&CK annotations are what make the "ATT&CK Tactic Threshold Exceeded" RIR and finding-group MITRE thresholds work.[^annotations][^33]

## The Splunk security portfolio (high level) — VOLATILE
- **Enterprise Security (ES)** — the flagship SIEM, repositioned as a unified **TDIR (Threat Detection, Investigation, Response)** platform; sold in **Essentials** (SIEM + AI Assistant + Detection Studio) and **Premier** (adds native UEBA + the converged Mission Control + SOAR experience) editions.[^2][^Portfolio]
- **Splunk SOAR** (Security Orchestration, Automation, Response) — **formerly Phantom** (acquired 2018); playbook-driven automation across many third-party tools; now a **native capability within ES** (Mission Control runs SOAR from the analyst queue) but a **separate subscription.**[^SOAR][^21]
- **UBA vs UEBA:** **Splunk UBA** (User Behavior Analytics) is the **legacy standalone** ML/anomaly product, reportedly reaching **End of Sale December 2025**; its successor **UEBA** is **built natively into ES Premier** (not sold standalone). Treat "Splunk UBA" as sunsetting.[^UEBA]
- **Cisco acquisition:** **Cisco completed its acquisition of Splunk on March 18, 2024** (~$157/share ≈ $28 billion). Talos intel is being woven across the portfolio; products are branded "Splunk Enterprise by Cisco" in some listings.[^Cisco-IR][^32]

## Disconfirming findings / gotchas
1. **RBA is not automatic noise reduction** — uncurated RBA generates as many (and duplicate) notables; tuning (throttle/dedup/weight/filter) is mandatory; keep threshold constant and tune scores.[^bestpractices][^RBA-FAQ]
2. **ES 8.x terminology churn has teeth** — ES 8.x **enforces previously-optional fields as required** (`risk_message`, `risk_object`, `rule_description`, `investigation_type`, `description`); if missing, **findings/notables silently fail to be created after upgrade**, and admins see migration errors. Post-upgrade breakage (filters, saved views, SAML/Mission-Control redirects) is documented in Known Issues. *(Strong but sourced to Splunk Known Issues + a partner blog — qualified.)*[^E1][^E2]
3. **RIRs were not auto-converted** to finding-based detections in 8.0; legacy + new coexist.[^16]
4. **CIM dependency is a hard requirement** — un-accelerated/un-normalized data ⇒ empty ES dashboards and non-firing correlation searches.[^24][^25]
5. **Cost & complexity** are the dominant independent criticisms — ES is among the priciest SIEMs, with multi-month implementations and a steep SPL learning curve; ingestion-based licensing creates "perverse incentives" to under-ingest (coverage gaps). *(Comparison figures are Tier-3 consultant/review blogs — directionally consistent, treat $ as illustrative.)*[^Zendikt][^decryptiondigest]
6. **Native RBA scaling ceiling** reported ~100k risk events/week before specialized engineering is needed. *(Single practitioner source with a commercial interest — tentative.)*[^Outpost]

## Adjacent / frontier concepts
Splunk Mission Control (the unified analyst work-surface); Splunk Attack Analyzer (automated phishing/malware analysis); Federated Analytics / Federated Search for security ("analyze data where it's stored"); **Detection-as-Code / Splunk Security Content / ESCU + `contentctl`** (the YAML-defined, Git-managed, MITRE-mapped detection pipeline — a strong candidate for its own skill); OCSF (the cross-vendor schema ES 8.x aligns to); finding-based detections & Finding Groups; native UEBA vs legacy UBA; Cisco convergence (XDR, Secure Network Analytics, Talos).

## References
[^1]: Configure data models for Splunk ES — Splunk Docs (ES 8.3) — https://help.splunk.com/en/splunk-enterprise-security-8/install/8.3/installation/configure-data-models-for-splunk-enterprise-security — ES uses CIM data-model acceleration; CIM bundled; ES installs unique models.
[^2]: About Splunk ES — Splunk Docs (ES 8.4) — https://help.splunk.com/en/splunk-enterprise-security-8/user-guide/8.4/introduction/about-splunk-enterprise-security — ES combines SIEM+SOAR+threat-intel; built on Splunk platform; unified TDIR; OCSF; RBA.
[^3]: Licensing for Splunk ES — Splunk Docs (ES 7.3) — https://help.splunk.com/en/splunk-enterprise-security-7/user-guide/7.3/introduction/licensing-for-splunk-enterprise-security — ES is a premium app on Splunk Enterprise/Cloud; needs DIV or vCPU license; no extra ingest.
[^5]: Overview of the Splunk CIM — Splunk Docs — https://help.splunk.com/en/?resourceId=common_information_model_data_managemenet — CIM = shared semantic model; data models normalized; packaged with ES.
[^6]: tstats command reference — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/search/spl-search-reference/10.4/search-commands/tstats — tstats over tsidx/accelerated data models.
[^7]: Configure correlation searches — Splunk Docs (ES 7.3) — https://help.splunk.com/en/splunk-enterprise-security-7/administer/7.3/correlation-searches/configure-correlation-searches-in-splunk-enterprise-security — thresholds, throttling, most disabled by default, throttling before suppression.
[^8]: Configure correlation searches (PCI/8.x) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise-security-8/splunk-app-for-pci-compliance — correlation search = scan sources for pattern → adaptive response.
[^9]: Correlation search overview — Splunk Docs (ES 7.3.2) — https://docs.splunk.com/Documentation/ES/7.3.2/Admin/Correlationsearchoverview — sources across domains; aggregate with SPL; example brute-force/malware.
[^10]: How RBA works in ES — Splunk Docs (ES 7.2) — https://help.splunk.com/en/splunk-enterprise-security-7/risk-based-alerting/7.2/introduction/how-risk-based-alerting-works-in-splunk-enterprise-security — risk rules→risk index; RIRs aggregate→risk notables; risk modifiers.
[^11]: How risk objects impact risk scores — Splunk Docs (ES 7.2) — https://help.splunk.com/en/splunk-enterprise-security-7/risk-based-alerting/7.2/create-risk-objects/how-risk-objects-impact-risk-scores-in-splunk-enterprise-security — risk_object/risk_object_type; references src/dest; only those two required.
[^12]: About RBA — Splunk Docs (ES 7.2) — https://help.splunk.com/en/splunk-enterprise-security-7/risk-based-alerting — RBA purpose: high-fidelity notables, increase true positives, reduce triage volume.
[^14]: Monitor SOC with findings — Splunk Docs (ES 8.0) — https://help.splunk.com/en/splunk-enterprise-security-8/user-guide/8.0/findings/monitor-your-security-operations-center-with-findings-in-splunk-enterprise-security — finding = merged notable+risk event; intermediate finding; Analyst Queue.
[^15]: Manage analyst workflows using the analyst queue — Splunk Docs (ES 8.0) — https://help.splunk.com/en/splunk-enterprise-security-8/user-guide/8.0/mission-control/manage-analyst-workflows-using-the-analyst-queue-in-splunk-enterprise-security — Analyst Queue actions; intermediate findings not triageable.
[^16]: Installing/upgrading to ES 8.x — Splunk Lantern — https://lantern.splunk.com/Security_Use_Cases/Advanced_Threat_Detection/Installing_and_upgrading_to_Splunk_Enterprise_Security_8x — 7.3→8.x terminology table; RIRs coexist; migration notes.
[^18]: splunk/security_content — GitHub (Splunk OSS) — https://github.com/splunk/security_content — analytic stories mapped to MITRE/Kill Chain/CIS; ships as ESCU; detection-as-code YAML.
[^20]: Splunk Security Content — research.splunk.com — https://research.splunk.com/ — detection catalog; STRT; SOAR playbooks.
[^21]: Introducing ES 8.0 — Splunk blog — https://www.splunk.com/en_us/blog/conf-splunklive/introducing-the-siem-of-the-future-splunk-enterprise-security-8-0.html — unified work surface, native SOAR, Finding Groups, OCSF.
[^22]: ESCU "Types of detection analytics" — Splunk Docs — https://help.splunk.com/en/splunk-cloud-platform/security-content-update — Hunting/TTP/Baseline/Anomaly/Correlation behaviors; Correlation enabled by default.
[^23]: Managing Incident Review — Splunk Docs (ES 7.3) — https://help.splunk.com/en/splunk-enterprise-security-7/user-guide/7.3/incident-review-and-investigations/managing-incident-review-in-splunk-enterprise-security — correlation search → notable; Incident Review categorizes by severity.
[^24]: Configure data models for ES (model list) — Splunk Docs (ES 8.4) — https://help.splunk.com/en/splunk-enterprise-security-8/install/8.4/installation/configure-data-models-for-splunk-enterprise-security — CIM model list + summary ranges (Risk=All Time).
[^25]: Data model reference for dashboards — Splunk Docs (ES 8.5) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.5/dashboards/data-model-reference-for-dashboards-in-splunk-enterprise-security — ES dashboards rely on CIM + accelerations.
[^27]: Triage notables on Incident Review — Splunk Docs (ES 7.3) — https://help.splunk.com/en/splunk-enterprise-security-7/user-guide/7.3/incident-review-and-investigations/triage-notables-on-incident-review-in-splunk-enterprise-security — Status/Urgency/Owner/Disposition fields.
[^29]: How urgency is assigned — Splunk Docs (ES 7.3) — https://help.splunk.com/en/splunk-enterprise-security-7/user-guide/7.3/incident-review-and-investigations/how-urgency-is-assigned-to-notable-events-in-splunk-enterprise-security — urgency_lookup = severity × asset/identity priority; higher priority wins.
[^30]: Overview of Incident Review — Splunk Docs (ES 7.3) — https://help.splunk.com/en/splunk-enterprise-security-7/user-guide/7.3/incident-review-and-investigations/overview-of-incident-review-in-splunk-enterprise-security — Incident Review dashboard; analyst triage workflow.
[^32]: .conf24 Splunk security innovations — Cisco Newsroom — https://newsroom.cisco.com/c/r/newsroom/en/us/a/y2024/m06/conf24-splunk-introduces-new-security-innovations.html — ES 8.0 + Mission Control + Federated Analytics; Talos across ES/SOAR/Attack Analyzer.
[^33]: Default risk incident rules / risk scoring — Splunk Docs (ES 7.3/8.5) — https://help.splunk.com/en/splunk-enterprise-security-7/risk-based-alerting/7.3/risk-incident-rules/default-risk-incident-rules-in-splunk-enterprise-security — "Risk Threshold >100/24h"; "ATT&CK >3 tactics/7d"; tstats RIR SPL.
[^34]: Adjusting risk using risk factors — Splunk Docs (ES 8.3) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.3/risk-based-alerting/adjusting-risk-using-risk-factors-in-splunk-enterprise-security — base/calculated/total score; 7 default risk factors; multipliers.
[^35]: How risk modifiers impact risk scores — Splunk Docs (ES 7.2) — https://help.splunk.com/en/splunk-enterprise-security-7/risk-based-alerting/7.2/modify-risk/how-risk-modifiers-impact-risk-scores-in-splunk-enterprise-security — modifiers carry risk_score/object/type; factors are multipliers.
[^36]: Create risk factors — Splunk Docs (ES 7.2) — https://help.splunk.com/en/splunk-enterprise-security-7/risk-based-alerting/7.2/create-risk-factors — risk-factor multipliers (×1.5/×1.25).
[^38]: Entity risk scoring — Splunk Docs (ES 8.3) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.3/risk-based-alerting/entity-risk-scoring-in-splunk-enterprise-security — ERS 0–100, 7-day, recalced every 20min; feeds UEBA. (single-source — qualified)
[^bestpractices]: Prioritizing threat objects over risk objects — Splunk Docs (ES 7.2) — https://help.splunk.com/en/splunk-enterprise-security-7/risk-based-alerting/7.2/best-practices/prioritizing-threat-objects-over-risk-objects-in-risk-incident-rules — curate RIRs by threat object; excessive/duplicate notables confuse analysts; reduce weight of noisy rules.
[^intel]: Overview of threat intelligence in ES — Splunk Docs (ES 8.4) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.4/threat-intelligence/overview-of-threat-intelligence-in-splunk-enterprise-security — intel ingest→KV store→threat-matching searches; native vs TIM Cloud split.
[^available-intel]: Available threat/generic intel sources — Splunk Docs (ES 8.4) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.4/threat-intelligence/available-threat-intelligence-and-generic-intelligence-sources-in-splunk-enterprise-security — bundled feeds; MITRE ATT&CK as generic source.
[^A&I-manage]: Manage assets & identities — Splunk Docs (ES 8.4) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.4/asset-and-identity-management/manage-assets-and-identities-to-enrich-findings-in-splunk-enterprise-security — correlation matches src/dest/dvc & user/src_user; entity zones.
[^A&I-config]: Configure asset & identity correlation — Splunk Docs (ES 8.5) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.5/asset-and-identity-management/configure-asset-and-identity-correlation-in-splunk-enterprise-security — search-time enrichment; no longer auto-correlates host/orig_host.
[^annotations]: Add annotations to detections — Splunk Docs (ES 8.5) — https://help.splunk.com/en/splunk-enterprise-security-8/administer/8.5/detections/add-annotations-to-detections-in-splunk-enterprise-security — managed (MITRE/Kill Chain/CIS/NIST/PCI) vs unmanaged; savedsearches.conf JSON; mitre_attack_enrichment lookup.
[^Portfolio]: Splunk Enterprise Security product page — Splunk (vendor) — https://www.splunk.com/en_us/products/enterprise-security.html — TDIR platform; Essentials vs Premier; native UEBA in Premier. VOLATILE.
[^SOAR]: Splunk SOAR product page — Splunk (vendor) — https://www.splunk.com/en_us/products/splunk-security-orchestration-and-automation.html — playbook automation; native within ES; separate subscription. VOLATILE.
[^UEBA]: Splunk UEBA product page — Splunk (vendor) — https://www.splunk.com/en_us/products/user-and-entity-behavior-analytics.html — UBA legacy EOS; UEBA native in ES Premier. VOLATILE (single vendor source — qualify).
[^Cisco-IR]: Cisco completes acquisition of Splunk — Cisco IR — https://investor.cisco.com/news/news-details/2024/Cisco-Completes-Acquisition-of-Splunk/default.aspx — closed March 18, 2024; ~$157/share ≈ $28B.
[^E3]: Getting Started with RBA and MITRE (.conf 2019 SEC1538) — Splunk .conf talk — https://conf.splunk.com/files/2019/slides/SEC1538.pdf — Apger/McIntosh; one-to-one vs risk-based model; alert-fatigue rationale.
[^E5]: RBA Helps SOCs Focus tech brief — Splunk (vendor) — https://www.splunk.com/en_us/pdfs/tech-brief/risk-based-alerting-helps-socs-focus-on-what-really-matters.pdf — 250–500 risk events → ~60–100 tickets.
[^E6]: apger/SA-RBA (deprecated add-on) — GitHub — https://github.com/apger/SA-RBA — original RBA mechanics; "built into ES, supported as of 6.6".
[^RBA-FAQ]: Risk-Based Alerting FAQs — Splunk Community — https://community.splunk.com/t5/Splunk-Enterprise-Security/Risk-Based-Alerting-FAQs/m-p/668170 — keep threshold constant/tune scores; implementation challenges; "RBA is a journey".
[^rba-gh]: Splunk RBA essential searches / RIR ideas — splunk.github.io/rba — https://splunk.github.io/rba/ — concrete tuning: throttle/dedup/weight/filter; benign_rate.
[^37]: Planning for Success with RBA — Splunk blog (Mills) — https://www.splunk.com/en_us/blog/security/planning-for-success-with-risk-based-alerting.html — maturity model; stakeholder buy-in.
[^Outpost]: Getting Started with RBA in Splunk ES — Outpost Security (practitioner/vendor) — https://www.outpost-security.com/rba-getting-started — ~100k risk events/week native ceiling; force notables to Informational during rollout. (tentative — commercial interest)
[^E1]: ES 8.0 Known Issues — Splunk Docs — https://help.splunk.com/en/splunk-enterprise-security-8/release-notes-and-resources/8.0/release-notes/known-issues — post-upgrade broken filters/saved views, SAML/Mission-Control redirect bug.
[^E2]: Upgrading ES from 7.x to 8.x — Seynur (partner blog) — https://blog.seynur.com/splunk/2026/02/18/upgrading-es-from-es7x-to-es811.html — ES 8.x now-required fields silently break findings; "Failed to migrate ES detections". (qualify)
[^Zendikt]: Splunk ES review — Zendikt (review aggregator) — https://www.zendikt.com/product/splunk-es — pricing range; multi-month implementation; post-Cisco pricing-complexity complaints. (Tier 3 — illustrative)
[^decryptiondigest]: Splunk vs Sentinel — decryptiondigest (blog) — https://www.decryptiondigest.com/blog/splunk-vs-microsoft-sentinel-siem-comparison — ingestion pricing "perverse incentives"; dual-SIEM/federated search. (Tier 3 — illustrative)
