<!-- hub-reference-banner -->
> **Reference file — part of the `career-and-formal-writing` hub.** Formerly the standalone `academic-and-citation-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: academic-and-citation-writing
description: Academic prose, scholarly argument, and citation discipline. Covers the five major citation systems (APA 7th, Chicago 17th notes-and-bibliography and author-date, MLA 9th, IEEE, Vancouver), in-text-vs-reference-list mechanics, signal phrases ("Smith argues that..."), avoidance of pseudo-citation ("studies show"), primary-vs-secondary source discipline, persistent identifiers (DOI, ORCID, ISBN, ISSN), reference-management tools (Zotero, BibTeX, Mendeley, EndNote), Toulmin argument structure (claim, grounds, warrant, backing, qualifier, rebuttal), the literature-review funnel pattern, and the abstract-vs-introduction distinction. TRIGGER: user asks to write, draft, edit, or review a paper, thesis, dissertation, journal article, conference paper, literature review, abstract, or systematic review; user says "cite this in APA", "Chicago footnote", "MLA works cited", "IEEE reference", "Vancouver numeric style", "DOI", "ORCID", "Zotero", "BibTeX", "signal phrase", "Toulmin", "literature review", "primary vs secondary source", "in-text citation", "reference list", "annotated bibliography". SKIP: writing a policy or governance document (use policy-and-governance-writing); writing an RFC, ADR, or technical design doc (use rfc-and-design-docs); writing a research proposal, grant application, or RFP response (use proposal-and-grant-writing); blog posts or marketing copy (use writing-expert); changelog or release-note entries (use changelog-and-release-notes); a runbook or procedure (use runbook-craft); executive memos (use executive-comms); rhetorical-argument analysis without citation discipline (use rhetorical-frameworks-deep).
---

# Academic and Citation Writing

## Overview

Academic writing is argument with receipts. Every non-trivial claim ties to a source the reader can independently retrieve and verify. The receipts come in two parts — an in-text marker that points to a fully specified entry in a reference list — and the rules for both parts are defined by a citation style.

The four problems this skill solves: choosing the right citation style for the venue, mechanically formatting in-text citations and reference entries in that style, integrating sources into prose with signal phrases (instead of dropping quote-bombs), and structuring an argument so that evidence supports a clearly stated claim rather than orbiting around it.

The five citation systems covered (APA 7th, Chicago 17th, MLA 9th, IEEE, Vancouver) cover the vast majority of academic and technical venues. Each makes different tradeoffs — author-date emphasizes recency, numeric systems emphasize compactness, notes-and-bibliography accommodates exotic sources — and choosing the wrong one is a desk-reject risk.

## Core Concepts

### 1. The five major citation systems at a glance

| Style | In-text format | Reference list | Field |
|---|---|---|---|
| **APA 7th** | (Smith, 2024, p. 42) | References — author-date | Social sciences, psychology, education, nursing |
| **Chicago 17th — Author-Date** | (Smith 2024, 42) | References — author-date | Sciences, social sciences |
| **Chicago 17th — Notes-Bibliography** | ¹ footnote with full cite, then short form | Bibliography — alphabetical | Humanities, history, literature, arts |
| **MLA 9th** | (Smith 42) | Works Cited — alphabetical, "core elements" template | Languages, literature, humanities |
| **IEEE** | [1] | References — numeric, by appearance | Engineering, computer science |
| **Vancouver** | (1) or superscript ¹ | References — numeric, by appearance | Biomedical, medical |

**Choose by venue, not preference.** Journals and conferences specify their style; theses follow institutional rules; submitting in the wrong style is grounds for return without review. When the venue offers a choice, default to APA in the social sciences, Chicago author-date in mixed STEM-humanities, MLA in literature/languages, IEEE in computing, and Vancouver in biomedical.

### 2. APA 7th edition — author-date system

In-text citation pairs author surname with publication year (and page number for direct quotes). Two forms:

- **Parenthetical:** "Working memory degrades under cognitive load (Smith, 2024)."
- **Narrative:** "Smith (2024) found that working memory degrades under cognitive load."

Direct quotes require a page or paragraph number: "(Smith, 2024, p. 42)" or "(Smith, 2024, para. 7)" for unpaginated sources.

**Reference list entry — journal article:**

> Smith, J. R., & Doe, A. B. (2024). Working memory and cognitive load. *Journal of Cognitive Science*, *48*(3), 215–238. https://doi.org/10.1234/jcs.2024.0042

Key features: hanging indent, italicized journal title and volume, alphabetical order by first author surname, all named authors listed up to 20 before truncation, DOI as a clickable URL.

APA 7 recommends **past tense or present perfect** for citing prior research ("Smith argued..." or "Smith has argued..."), reserving present tense for established findings ("Working memory is finite").

### 3. Chicago Manual of Style 17th edition — two systems

Chicago is unusual in that it offers two complete, equally official systems. The choice is by field convention.

**Notes-Bibliography system** — for humanities. In-text is a superscript number; the matching footnote (or endnote) contains the full citation on first appearance, then a short form on subsequent appearances.

First note:

> 1. John R. Smith, *Cognitive Load Theory in Practice* (Chicago: University of Chicago Press, 2024), 42.

Subsequent note:

> 2. Smith, *Cognitive Load*, 67.

A separate Bibliography lists every work cited, alphabetical by author surname.

**Author-Date system** — for sciences and social sciences. Looks much like APA but with no comma before the year:

> "Working memory degrades under cognitive load (Smith 2024, 42)."

Reference list entry:

> Smith, John R. 2024. *Cognitive Load Theory in Practice*. Chicago: University of Chicago Press.

Both systems are equally Chicago; pick one and use it consistently throughout the document.

### 4. MLA 9th edition — core-elements template

MLA 9 (2021) introduced a flexible **core-elements template** that handles any source type by filling slots in order: Author. "Title of Source." *Title of Container*, Other contributors, Version, Number, Publisher, Publication date, Location.

In-text: author and page only, no year. "(Smith 42)" or "Smith argues that... (42)."

Works Cited entry — journal article:

> Smith, John R. "Working Memory and Cognitive Load." *Journal of Cognitive Science*, vol. 48, no. 3, 2024, pp. 215–38.

Works Cited entry — webpage:

> Smith, John R. "Cognitive Load in 2024." *Cognitive Today*, 14 Mar. 2024, www.cognitivetoday.org/articles/load-2024.

MLA's strength is handling unusual sources (TikToks, podcasts, video-game cutscenes) through the same template. Format: title the page **Works Cited**, alphabetical, hanging indent, double-spaced.

### 5. IEEE and Vancouver — numeric systems

Both number references in the order they first appear in the text. Both keep in-text citations short, which suits dense technical prose with many citations.

**IEEE** — sources cited as bracketed numbers: "Working memory degrades [1], [2]." The reference list is numbered, in the order of first appearance, not alphabetical.

> [1] J. R. Smith, "Working memory and cognitive load," *J. Cognitive Sci.*, vol. 48, no. 3, pp. 215–238, 2024, doi: 10.1234/jcs.2024.0042.

IEEE abbreviates journal titles per its own list, uses initials-first author names, and italicizes the journal title.

**Vancouver** — used in biomedical writing, defined by the International Committee of Medical Journal Editors (ICMJE). In-text uses parenthesized or superscript numbers: "(1)" or "¹". Reference list uses **Index Medicus / NLM abbreviations** for journal names:

> 1. Smith JR, Doe AB. Working memory and cognitive load. J Cogn Sci. 2024;48(3):215–38.

Vancouver lists the first six authors then "et al." (varies by journal), uses no italics, and is the dominant style in PubMed-indexed journals.

### 6. Signal phrases — integrating sources into prose

A signal phrase introduces a source in your own sentence, so the citation feels woven in rather than tacked on. Compare:

**Quote-bomb (weak):**
> Working memory has a limited capacity. "Most adults can hold about four items in working memory at once" (Cowan 2010, p. 87).

**Signal-phrase integration (strong):**
> As Cowan (2010) argued, "most adults can hold about four items in working memory at once" (p. 87), a finding that has shaped subsequent cognitive-load research.

Signal-phrase verbs carry stance. They tell the reader how confident the cited author is, and how confident you are in repeating them:

- **Neutral:** notes, states, reports, observes, describes, writes
- **Argumentative:** argues, contends, claims, asserts, maintains
- **Evidence-based:** found, demonstrated, showed, established, documented
- **Hedged:** suggests, proposes, hypothesizes, speculates, implies
- **Critical:** disputes, challenges, questions, complicates, refutes

APA 7 prefers **past tense or present perfect** with signal phrases ("Cowan argued" or "Cowan has argued"), reflecting that the research is completed; reserve present tense for findings treated as ongoing fact ("Working memory is limited").

### 7. Pseudo-citation — the failure mode

**"Studies show..."** with no citation is pseudo-citation. So are "research suggests", "it is well known that", "experts agree", and any other appeal to unnamed authority. Pseudo-citation is rejected at every level of academic review because it cannot be verified.

Three repairs:

1. **Name the source.** "Studies show" → "Cowan (2010) and Baddeley (2012) showed".
2. **Drop the unsupportable claim.** If you cannot find a citation, you do not yet have evidence.
3. **Mark it as your own argument.** "I argue that..." is honest. "Studies show..." with no studies is not.

A related failure: citing a secondary source as if it were primary. If you read Cowan being quoted in a 2024 review, cite the 2024 review — or go read Cowan directly. Never pretend you read what you only read about.

### 8. Primary, secondary, tertiary sources

- **Primary** — original research, raw data, firsthand accounts. Empirical studies, clinical trials, archival documents, interviews, statutes, original artworks.
- **Secondary** — analysis or interpretation of primary sources. Literature reviews, meta-analyses, textbooks, biographies, scholarly monographs.
- **Tertiary** — compilations of secondary sources. Encyclopedias, Wikipedia, dictionaries, handbook chapters that survey a field.

Discipline-specific norms:

- **Sciences:** lean heavily on primary research articles; secondary sources for context only; tertiary almost never cited.
- **Humanities:** primary texts (the novel, the statute, the artwork) anchor the argument; secondary scholarship situates it; tertiary occasionally tolerated for historical orientation.
- **Always:** when citing a finding, cite the primary source whenever possible. Cite the review only when discussing the review itself.

### 9. Persistent identifiers — DOI, ORCID, ISBN, ISSN

Modern citations rely on persistent identifiers that survive URL rot:

- **DOI (Digital Object Identifier)** — a permanent link to a digital object. Format `10.xxxx/yyyy`. Include in references as `https://doi.org/10.xxxx/yyyy`. Required by APA 7, IEEE, and most journals.
- **ORCID** — a persistent identifier for an author (format `0000-0000-0000-0000`). Used in submissions metadata; not typically inserted into prose citations but increasingly required in author bylines.
- **ISBN** — books (10-digit pre-2007 or 13-digit). Optional in most styles; required by some publishers in works-cited entries.
- **ISSN** — journals/serials. Rarely cited in-text; used in cataloging.
- **PMID / PMCID** — biomedical, indexed in PubMed. Vancouver style often includes these.
- **arXiv ID** — preprints in physics, CS, math. Include in IEEE-style references as `arXiv:2401.01234`.

A citation with a DOI is verifiable years later; a citation with only a URL may not be. Always prefer DOI.

### 10. Reference-management tools

Manual citation formatting is error-prone and time-consuming. Use a reference manager.

- **Zotero** — free, open source, browser-integrated. Excellent for capturing sources from web pages and JSTOR. Exports to Word, Google Docs, LaTeX (via Better BibTeX plugin), and Markdown.
- **Mendeley** — Elsevier-owned, PDF-centric, social/sharing features. Free; some features behind paywall.
- **EndNote** — commercial; heavy in life sciences and corporate research libraries.
- **BibTeX / BibLaTeX** — plaintext bibliography format for LaTeX. The lingua franca for computer science and physics. A `.bib` file is git-friendly and tool-agnostic.
- **CSL (Citation Style Language)** — XML-based style files used by Zotero, Mendeley, Pandoc. The Zotero Style Repository has 10,000+ journal styles ready to install.
- **Pandoc + CSL** — for Markdown-to-anything authoring, Pandoc renders `[@smith2024]` citations against a `.bib` file and a CSL style. The native pipeline for many computational researchers.

Discipline shorthand: humanities → Zotero, social sciences → Zotero or Mendeley, life sciences → EndNote or Zotero, computing/physics → BibTeX, anyone using Pandoc → CSL.

### 11. Toulmin argument structure

Stephen Toulmin's argument model gives a standardized vocabulary for breaking an argument into parts and inspecting whether each is supported. Six elements:

1. **Claim** — the position being defended. ("Working memory training transfers to general cognitive ability.")
2. **Grounds** (data, evidence) — the facts that support the claim. ("Jaeggi et al. (2008) showed transfer in n-back training.")
3. **Warrant** — the underlying principle that lets the grounds support the claim. ("Transfer of cognitive skill requires shared cognitive substrates.")
4. **Backing** — support for the warrant itself, often theoretical. ("Cognitive load theory (Sweller, 1988) predicts substrate-shared transfer.")
5. **Qualifier** — limits on the claim. ("In healthy young adults, under controlled training conditions.")
6. **Rebuttal** — counterclaims and exceptions acknowledged. ("Melby-Lervåg & Hulme (2013) failed to replicate.")

Toulmin's value in writing: every paragraph in an academic argument should be analyzable into these slots. A paragraph with grounds but no warrant feels like data without conclusion; a paragraph with claim but no grounds feels like assertion.

### 12. The literature-review funnel

A literature review opens broad, narrows, and lands on the specific gap your work addresses. The funnel shape:

1. **Field opening** — the broad domain (1–2 paragraphs). "Cognitive load research has shaped instructional design for four decades."
2. **Sub-area** — the specific corner of the field. "Within cognitive load theory, the distinction between intrinsic, extraneous, and germane load remains contested."
3. **Active debate** — the specific scholarly conversation. "Recent work disputes whether germane load is empirically separable from intrinsic load (Kalyuga, 2011; Sweller, 2010)."
4. **The gap** — what is missing. "No study has tested this distinction in real-time multimedia learning environments."
5. **Your contribution** — what your work does. "The present study addresses this gap by..."

A literature review that stays broad is a textbook chapter; one that opens narrow is a research note. The funnel keeps the reader oriented from familiar territory to the specific question.

### 13. Abstract vs. introduction — different jobs

The two are routinely confused.

**Abstract** — a self-contained summary readable without the paper. Usually 150–300 words. Four parts (IMRD compressed):

1. Context / problem (1 sentence)
2. Method (1–2 sentences)
3. Results (2–3 sentences with the actual numbers)
4. Implication (1 sentence)

The abstract gets indexed; it has to stand alone, which means no citations (in most styles), no abbreviations on first use, and the key result must appear.

**Introduction** — the opening section of the paper proper. Several pages. Functions:

1. Establishes the broader context
2. Reviews relevant prior work (citations everywhere)
3. Identifies the gap
4. States the research question and hypothesis
5. Previews the paper's structure

Confusing them: an abstract that reads like the first paragraph of the intro buries the result; an intro that reads like an expanded abstract skips the literature.

## Templates and Examples

### Template: APA reference entries (the four most common)

```
Journal article:
Author, A. A., & Author, B. B. (Year). Title of article. Journal Name,
  Volume(Issue), Pages. https://doi.org/xx.xxxx/yyyy

Book:
Author, A. A. (Year). Title of book (Nth ed.). Publisher.

Book chapter:
Author, A. A. (Year). Title of chapter. In E. Editor (Ed.), Title of book
  (pp. xx–xx). Publisher. https://doi.org/...

Web page:
Author, A. A. (Year, Month Day). Title of page. Site Name. URL
```

### Template: Toulmin paragraph skeleton

```
[Claim] — opening sentence stating the position.
[Grounds] — evidence: "Smith (2024) found that..."
[Warrant] — connective principle: "This matters because..."
[Backing, if needed] — support for the warrant.
[Qualifier] — limits: "In the population studied..."
[Rebuttal] — counter-evidence acknowledged: "Though Jones (2023) reported..."
```

### Template: Lit-review funnel paragraph progression

```
¶1: Field — broad domain
¶2-3: Sub-area — specific corner
¶4-6: Debate — competing positions, cited
¶7: Gap — what is missing
¶8: Bridge — "The present study addresses this gap by..."
```

### Example: Signal-phrase integration in three styles

**APA 7th:**
> Cowan (2010) argued that "most adults can hold about four items in working memory at once" (p. 87), challenging Miller's earlier seven-item estimate.

**Chicago Author-Date:**
> Cowan (2010, 87) argued that "most adults can hold about four items in working memory at once," challenging Miller's earlier seven-item estimate.

**IEEE:**
> Cowan argued that "most adults can hold about four items in working memory at once" [1], challenging Miller's earlier seven-item estimate.

## Anti-Patterns

1. **Pseudo-citation** — "Studies show...", "research suggests...", "it is widely accepted..." with no source. Fix: name the source or drop the claim.

2. **Quote-bombing** — paragraph after paragraph of standalone quotations with no integration or analysis. Fix: signal phrases plus your own interpretation between every quoted sentence.

3. **Secondary cited as primary** — "Cowan (2010) found that working memory is limited" when you read it in a 2024 review and never opened Cowan. Fix: read the primary or cite the secondary honestly ("Cowan, as cited in Brown, 2024").

4. **Style mixing** — APA in-text with MLA works-cited entries, or footnotes alongside parenthetical author-date. Fix: pick one style at the start and stay consistent.

5. **Missing DOIs** — references with URLs that will rot, when DOIs exist. Fix: look up the DOI on crossref.org and include it.

6. **Et al. abuse** — collapsing "Smith, Jones, & Brown (2024)" to "Smith et al. (2024)" on first mention in styles that require the full author list. APA 7 in particular uses et al. from the first citation only when 3+ authors; older style had a different rule. Fix: check the current edition.

7. **Wikipedia as primary citation** — citing Wikipedia (a tertiary source) for an empirical claim. Fix: follow the Wikipedia citation to its primary source and cite that.

8. **Inflated reference list** — citing sources you skimmed but did not read, to look thorough. A reviewer will catch it when your characterization of the source is wrong. Fix: cite only what you have actually read and can defend.

9. **Abstract that announces ("This paper will discuss...")** instead of summarizing the result. Fix: state what you found, not what you intend to do.

10. **The literature review with no gap** — describes the field for ten pages without ever saying what is missing. Reader leaves not knowing why the present work exists. Fix: every lit review ends with the gap and your contribution.

11. **Year-only in narrative ("Smith showed")** without page number for a direct quotation. APA requires page or paragraph for quotes. Fix: add (p. xx) or (para. xx).

12. **Inconsistent hanging-indent or alphabetization** in the reference list. Reviewers and editors flag this on sight. Fix: use a reference manager and let it format.

## Decision Heuristics

**Which citation style do I use?**

- Check the venue's submission guidelines first — they specify the style.
- Thesis or dissertation → institutional graduate-school style guide.
- Social sciences / psychology / education → APA 7.
- Humanities, history, literature → Chicago Notes-Bibliography or MLA 9.
- Sciences, mixed STEM → Chicago Author-Date or APA 7.
- Engineering, CS → IEEE.
- Biomedical, clinical → Vancouver (ICMJE).

**Parenthetical or narrative citation?**

- Emphasis on the *finding* → parenthetical: "Working memory is limited (Cowan, 2010)."
- Emphasis on the *author* or comparing authors → narrative: "Cowan (2010) found... whereas Miller (1956) argued..."
- Use both throughout — varying the form keeps prose readable.

**MUST I include a page number?**

- Direct quotation → always.
- Paraphrase of a specific argument in a long work → recommended.
- General reference to a whole work → not required.

**Primary or secondary?**

- Discussing a finding → primary source.
- Discussing how the field has interpreted a finding → secondary source.
- Cannot get the primary → cite the secondary honestly with "as cited in".

**Signal-phrase verb choice?**

- The cited author *asserted* without empirical support → "argues", "claims", "contends".
- The cited author *demonstrated* with data → "found", "showed", "established".
- The cited author *speculated* or proposed → "suggests", "proposes", "hypothesizes".
- The cited author was *refuted* later → "challenged", "complicated", "disputed".

**Reference manager?**

- New to academic writing → Zotero (free, easy, exports everywhere).
- LaTeX user → BibTeX (with Better BibTeX in Zotero for sync).
- Lab uses EndNote → use EndNote.
- Markdown / Pandoc workflow → CSL + .bib.

## References

1. **Publication Manual of the American Psychological Association, 7th ed.** (2020). APA. <https://apastyle.apa.org/>
2. **The Chicago Manual of Style, 17th ed.** (2017). University of Chicago Press. <https://www.chicagomanualofstyle.org/>
3. **MLA Handbook, 9th ed.** (2021). Modern Language Association. <https://style.mla.org/>
4. **IEEE Reference Guide** (current). IEEE. <https://ieeeauthorcenter.ieee.org/wp-content/uploads/IEEE-Reference-Guide.pdf>
5. **ICMJE Recommendations** (Vancouver style). International Committee of Medical Journal Editors. <https://www.icmje.org/recommendations/>
6. **Purdue OWL — APA, Chicago, MLA guides.** <https://owl.purdue.edu/owl/research_and_citation/resources.html>
7. **Toulmin, S. E.** (2003). *The Uses of Argument* (Updated ed.). Cambridge University Press.
8. **Zotero documentation.** <https://www.zotero.org/support/>

## Cross-references

- **policy-and-governance-writing** — for documents that prescribe rather than argue
- **rfc-and-design-docs** — for technical proposals (not academic papers)
- **writing-expert** — for prose-level craft inside the academic structure
- **rhetorical-frameworks-deep** — for argumentation theory beyond Toulmin
- **plain-language** — when academic prose needs to reach a public audience

---

## Footnote vs endnote vs inline citation — choosing the form

**Rule.** Three citation forms coexist; the choice depends on *what the reader needs to do with the reference*. None is universally better; using the wrong form for the task is the common failure.

| Form | What the reader sees | Best for | Worst for |
|---|---|---|---|
| Inline parenthetical | "(Williams, 1990, p. 50)" or "[1]" | Frequent in-text citation; reader who only needs author/year to triangulate | Bibliographic detail; commentary on the cited source |
| Footnote (bottom of page) | superscript ¹ with full citation at page bottom | Reader-on-the-page who wants the source without breaking flow; commentary, asides, side arguments | Documents with very dense citation; very long bibliographic notes |
| Endnote (end of document or chapter) | superscript ¹ with full citation at the end | Documents with heavy citation that would clutter every page; print volumes where notes are gathered for typesetting | Web docs where the reader cannot easily jump back and forth; commentary that should be read inline |

### Style-system defaults

| System | Primary citation | When notes are used |
|---|---|---|
| APA (7th ed.) | Author–date inline parenthetical | Notes are *content* footnotes only; bibliographic references are inline + reference list |
| Chicago notes-and-bibliography | Footnotes (or endnotes) with shortened repeat-citations | Notes carry full bibliographic info; bibliography is alphabetical at end |
| Chicago author–date | Author–date inline | Notes are content only |
| MLA (9th ed.) | Author–page inline parenthetical "(Williams 50)" | Notes are content only; works-cited list at end |
| IEEE | Numbered inline "[1]" | No notes; numbered reference list at end |
| Vancouver / ICMJE | Numbered inline | No notes; numbered reference list at end |

### When to choose each

**Use inline parenthetical when:**
- The discipline expects it (APA, MLA, IEEE, Vancouver — most science, social science, engineering).
- The reader cares about *who said it and when*, not the full bibliographic form, while reading.
- Citation density is moderate to high — full footnotes for every citation would overwhelm the page.

**Use footnotes when:**
- The discipline expects it (Chicago notes-and-bibliography — common in humanities, history, law, theology).
- The reader will want the source *while reading the page* and is reading print or print-like layouts where the eye can drop to the bottom.
- You have asides, commentary, or side arguments that would derail the main text but are worth preserving.
- You are writing for a publication that typesets the footnotes well (most academic books, some journals).

**Use endnotes when:**
- The discipline allows it and the document is long (a print book or monograph) where pages would otherwise be visually dominated by note matter.
- The reader is expected to read straight through and consult notes only when finished, or only when curious.
- The publisher requires it for typesetting reasons.
- *Avoid* endnotes for documents read on screen — the user cannot easily flip back and forth as they could in print.

### Worked example — same citation, three forms

The claim: *Williams argues that style is a matter of clarity, not ornament.*

- APA inline: *Williams argues that style is a matter of clarity, not ornament (Williams, 1990, p. 50).*
- Chicago footnote: *Williams argues that style is a matter of clarity, not ornament.¹* with footnote: *¹ Joseph M. Williams, Style: Toward Clarity and Grace (Chicago: University of Chicago Press, 1990), 50.*
- IEEE numbered: *Williams argues that style is a matter of clarity, not ornament [1].* with reference list: *[1] J. M. Williams, Style: Toward Clarity and Grace. Chicago: Univ. of Chicago Press, 1990, p. 50.*

### Content notes vs bibliographic notes

Both APA and Chicago author–date allow *content footnotes* — notes that contain commentary, qualifications, or tangents rather than citations. Use sparingly. A doctrine-driven test: if the content note is load-bearing for the argument, integrate it into the main text. If it is genuinely a side comment, the footnote is the right home. If you have more than ~3 content footnotes per page, your main text is probably hiding from the reader.

### When to break it

- *Hybrid documents* (a technical report with one or two citations) can use a single inline citation form without a full apparatus.
- *Web articles* often use linked inline citations (a hyperlink on the claim) — a fourth form, increasingly common. It combines inline parenthetical with one-click access, but loses the bibliographic detail unless paired with a references section.
- *Long-form journalism* commonly uses endnotes-only (placed at the end of the article, often hidden behind a "Notes" link). This is fine for a single article but breaks down when the reader needs to verify a specific claim while reading.

### References

- *Publication Manual of the American Psychological Association*, 7th ed. (2020), §2.13 (footnotes) and Chapter 8 (citations).
- *The Chicago Manual of Style*, 17th ed. (2017), Chapter 14 (Notes and Bibliography) and Chapter 15 (Author-Date).
- *MLA Handbook*, 9th ed. (2021), §§6.1–6.5 (in-text citations and notes).
- IEEE Editorial Style Manual — section on citation format. https://journals.ieeeauthorcenter.ieee.org/wp-content/uploads/sites/7/IEEE-Reference-Guide.pdf
- Turabian, K. *A Manual for Writers of Research Papers, Theses, and Dissertations*, 9th ed., University of Chicago Press, 2018 — the student-facing companion to Chicago, especially clear on footnote vs endnote choice.
