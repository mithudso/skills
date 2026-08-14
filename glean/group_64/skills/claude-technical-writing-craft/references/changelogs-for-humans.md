<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `changelogs-for-humans` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: changelogs-for-humans
description: "User-facing changelog craft — the end-user counterpart to the developer changelog. Benefit-led language (not feature-list language), screenshots and GIFs for UI changes, the 'skip the version numbers' approach (hide or downplay semver), grouping by audience benefit not file-tree, the Linear/Stripe/Vercel/Notion changelog pattern with images, ChangeKeep readable conventions, RSS-feed friendliness, monthly digest patterns, and Mailchimp 'What's New' email format. SIBLING to changelog-and-release-notes, NOT a replacement — that skill is developer-facing (semver, breaking changes, migration guides); this skill is end-user-facing (benefits, screenshots, tone). TRIGGER: 'user-facing changelog', 'what's new page', 'product update post', 'monthly digest email', 'release highlights for users', 'changelog with screenshots', 'in-product release notes', 'feature announcement digest', 'how do I make my changelog readable', 'Linear-style changelog', 'Vercel-style changelog'. SKIP: developer-facing semver changelog, breaking-change announcement, migration guide (use changelog-and-release-notes); a single big-bang product launch (use release-blog-and-launch-narrative); a conversion-focused landing page (use sales-and-marketing-copy); the prose-craft pass (use writing-expert)."
version: "1.0.0"
updated: "2026-05-29"
category: custom
tags:
  - writing
  - release-notes
  - user-comms
  - product-marketing
parent_concept: "Writing and Documentation"
whenToUse:
  - "Write a user-facing changelog entry"
  - "Draft a 'What's New' page or email"
  - "Translate engineering release notes into user benefits"
  - "Make our changelog look like Linear's or Vercel's"
  - "Add screenshots to a changelog entry"
  - "Plan a monthly product-update digest"
  - "Write release highlights for the in-product widget"
  - "Should the changelog show version numbers"
  - "Group these changes by user benefit, not by service"
  - "Write a Mailchimp 'What's New' style email"
related_skills:
  - changelog-and-release-notes
  - release-blog-and-launch-narrative
  - sales-and-marketing-copy
  - writing-expert
  - technical-writing-craft
  - executive-comms
triggers:
  - user-facing changelog
  - what's new page
  - product update post
  - monthly digest email
  - release highlights for users
  - changelog with screenshots
  - in-product release notes
  - feature announcement digest
  - Linear-style changelog
  - Vercel-style changelog
  - Stripe changelog format
  - Notion changelog
  - benefit-led changelog
  - hide version numbers in changelog
---

# Changelogs for Humans

Reference for **user-facing changelog craft**: the public-product update feed read by end users, not the engineering changelog read by integrators. The sibling skill (`changelog-and-release-notes`) covers developer-facing release notes (semver, breaking changes, migration guides, Keep a Changelog spec). This skill covers the visual, benefit-led, screenshot-rich form pioneered by Linear, Vercel, Stripe, Notion, GitHub, and Mailchimp.

Sources: Linear changelog (linear.app/changelog), Vercel changelog (vercel.com/changelog), Stripe Blog: Changelog (stripe.com/blog/changelog), GitHub changelog (github.blog/changelog), Notion changelog (notion.so/releases), Mailchimp "What's New" (mailchimp.com/whats-new), keepachangelog.com (Olivier Lacan), Mintlify "Five changelog principles from best-in-class developer brands" (mintlify.com/blog), Usersnap "10 Inspiring Changelog Examples" (usersnap.com/blog), Beamer / AnnounceKit / Frill / Quickhunt / Screenhance guides, Microsoft Writing Style Guide, Google developer-doc style guide.

SKIP: a developer-facing semver changelog, breaking-change announcement, or migration guide (use `changelog-and-release-notes` — the sibling, not the replacement); a one-shot product-launch narrative arc (use `release-blog-and-launch-narrative`); a conversion-focused landing page (use `sales-and-marketing-copy`); the prose-craft polish pass (use `writing-expert`).

---

## How to use this skill

When invoked, follow this sequence:

1. **Confirm the audience is the end user, not the developer/integrator.** If the change is breaking, requires a migration, or touches an API contract, the *developer* changelog (use `changelog-and-release-notes`) is the primary artifact and this skill produces an optional companion entry.
2. **Identify the surface.** Is this for a public changelog feed (Linear-style)? A monthly digest email (Mailchimp-style)? An in-product "What's New" panel (Notion-style)? Surface dictates length, tone, and visual treatment.
3. **Group by user benefit, not by file tree.** See Section 2.
4. **Lead with the outcome, not the mechanism.** See Section 3.
5. **Decide on visual treatment.** Major UI changes get a screenshot or GIF; routine fixes stay text-only. See Section 4.
6. **Decide on version-number visibility.** See Section 5.

If the input is a raw engineering changelog or a list of merged PRs, ask once: "Is this for the user-facing feed (Linear/Stripe-style) or the developer changelog (semver, migration notes)?" Don't produce both. Recommend the developer skill if the input is API-shaped.

---

## 1. The shape of a user-facing changelog

The format Linear, Vercel, Stripe, Notion, GitHub, and Mailchimp converged on:

```text
┌────────────────────────────────────────────────────────────────┐
│                                                                │
│   [Date]                                                       │
│   ## [Headline framed as a user benefit, not a feature name]   │
│                                                                │
│   [One- to three-paragraph narrative explaining what changed   │
│   and what it means for the user. Author voice. Real human.]   │
│                                                                │
│   [Optional screenshot, GIF, or short video — inline, large,   │
│   high contrast. Captioned only if not self-explanatory.]      │
│                                                                │
│   ─── Smaller items below the fold ──────────────────────────  │
│                                                                │
│   • Bug fix: [user-visible symptom, not stack trace]           │
│   • Improvement: [outcome verb]                                │
│   • [Optional: "We also..."]                                   │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

What this shape gets right:

- **Date is the strongest ordering signal**, not version number.
- **Headline is a benefit**, not a feature name. "Faster search" beats "Search v2.3."
- **Author voice** — written like a product person talking, not like CI output.
- **Visual hierarchy** — one big thing per entry gets the image; everything else is a bullet.
- **No labels like FIXED / ADDED / CHANGED** unless the audience is integrators. Those belong in the developer changelog.

---

## 2. Group by user benefit, not by file tree

The developer changelog can group by service, package, or commit type. The user changelog must group by **what the user can now do** (or no longer has to do).

| Engineering grouping (wrong here) | User-benefit grouping (right) |
|---|---|
| `auth-service`, `billing-api`, `webhook-worker` | "Faster sign-in," "Clearer billing," "Webhooks you can trust" |
| `feat:`, `fix:`, `chore:`, `refactor:` | "New," "Improved," "Fixed" (audience-readable verbs) |
| Version `1.42.0`, `1.42.1`, `1.42.2` | "May 14," "May 21," "June 3" — one month, three entries |

If your input arrives grouped by service, **re-section it** before drafting. The user does not know which service owns what.

**Heuristic.** If a user reading the changelog can't say "oh, that affects me, here's what I do differently now" within five seconds, the entry is grouped wrong, written wrong, or both.

---

## 3. Benefit-led language (the translation table)

Every user-facing entry rewrites engineering language into outcomes. Apply this translation reflexively.

| Engineering wording | Benefit-led rewrite |
|---|---|
| "Migrated sync layer to WebSockets" | "Changes now appear instantly on every device" |
| "Improved performance" | "Dashboard loads in under one second on slow connections" |
| "Refactored permission engine" | "You can now share a project without granting full account access" |
| "Added pagination to API endpoint" | "Lists with more than 100 items now load without timing out" |
| "Upgraded to Postgres 16" | (Omit. Internal change, no user-visible effect.) |
| "Fixed bug in import handler" | "Fixed: CSV imports no longer drop the last row" |
| "Released v1.42.0" | (Omit version number, lead with the date.) |
| "BREAKING: API v1 sunset" | Use the developer changelog (`changelog-and-release-notes`); this skill is wrong tool. |

**Rules.**

1. **Replace internal nouns with user verbs.** "Sync layer" → "changes appear." Verbs the user does or experiences, not subsystems we own.
2. **Quantify when honest.** "3x faster" is honest if measured; "blazing fast" is unfalsifiable. If you don't have a number, describe the felt experience ("now feels instant on common pages").
3. **Cut entries with no user-visible effect.** Postgres upgrades, dependency bumps, refactors, internal tooling — these belong in commit history, not the user changelog.
4. **Name the bug by the symptom, not the cause.** "Fixed: CSV imports drop last row" beats "Fixed off-by-one in `importCsv` line 442."
5. **Keep one sentence per atom.** If an entry needs three sentences, it's probably two changes.

---

## 4. Visuals — when, where, what

The Linear/Vercel pattern made one rule the industry default: **major changes get an image; routine items stay text.** This creates a natural visual hierarchy that scales — readers know that a screenshot means "look at this." Teams that add visuals to user-facing changelogs report 2–3x higher engagement than text-only changelogs (per Screenhance and Beamer analyses).

| Change type | Visual | Why |
|---|---|---|
| New UI surface (panel, page, modal) | Static screenshot, full-bleed | "Here's what the new thing looks like." |
| Interaction change (drag, hover, multi-select) | Short GIF, 4–8 seconds, autoplay-friendly | Stills can't show interaction. |
| Workflow change spanning multiple screens | Captioned screenshot pair, before/after | One image is wrong; show the transition. |
| Performance improvement | Optional: chart, gauge, or "before/after" timer | Only if you have the data. |
| Bug fix | None | Don't screenshot the absence of a bug. |
| Backend/infra change with no UI surface | None | If the user can't see it, don't pretend. |

**Image craft checklist:**

- **Real product data** wherever the user's eye lands. Lorem ipsum and "Acme Corp" make the product look unfinished.
- **Dark mode if your product has dark mode** — match the default user experience. Pick one mode per changelog so the page doesn't strobe.
- **Crop tightly.** Don't show the browser chrome or the whole app shell when the change lives in one panel.
- **Caption only when the image isn't self-explanatory** — overrunning every image with a caption is noise.
- **Alt text** for screen readers describes the change, not the chrome. "Settings panel now shows a 'Default workspace' dropdown" — not "screenshot."
- **GIFs**: under 5 MB, autoplay-on, loop. If it's longer than 8 seconds, it's a video, not a GIF.

---

## 5. Version numbers — hide, downplay, or omit

Linear, Vercel, GitHub, Notion: no version numbers in the user changelog. Stripe: dated. Mailchimp "What's New": dated, no semver. The pattern is consistent — **the date is the version**.

| Version number visibility | Use when |
|---|---|
| **Hidden / omitted** | Continuous-deploy SaaS, no install step, no client-side caching the user controls. (Linear, Vercel, Notion default.) |
| **Downplayed** (footer, faint) | Same as above but you want changelog readers to be able to cite a specific deploy. |
| **Shown** | Self-hosted, downloadable client, mobile/desktop app with explicit updates, anything where the user might be on an older version. (Use the dev changelog instead — `changelog-and-release-notes`.) |

**Heuristic.** If the user has no way to roll back, no way to be on an older version, and no way to be affected by semver — strip the version number. The date is enough.

When the product *does* require version awareness (mobile apps, electron apps, self-hosted), the user changelog can include a small "build 4.12.0 · May 28" line at the top of each entry, but the entry itself still leads with the benefit.

---

## 6. Channels and distribution patterns

A user changelog usually lives on more than one surface. Pattern each so they reinforce, not duplicate.

| Channel | Length | Visual | Cadence |
|---|---|---|---|
| **Public changelog page** (Linear/Vercel/Stripe) | Long form per entry, 1–5 paragraphs | Inline screenshot/GIF per major item | Per-ship — could be daily |
| **RSS feed** (off the public page) | Same content; ensure absolute image URLs, no relative links, clean HTML | Yes (must use absolute URLs) | Auto from the page |
| **In-product "What's New" panel** | 1–2 sentences per item, link to full entry | Small thumbnail or icon | Weekly to monthly |
| **Monthly digest email** (Mailchimp pattern) | "Top 5 of the month" + link to changelog page | Hero image + per-section image | Monthly, fixed date |
| **Tweet/X thread** | One line + link, one tweet per major item | Yes | Per-ship for headline items only |
| **Slack/Discord update channel** | One line + link to changelog entry | Optional inline image | Per-ship |

**Rules.**

1. **The changelog page is the canonical artifact.** Every other channel links back to it.
2. **The RSS feed must work.** Absolute image URLs, valid XML, stable item GUIDs, dates in RFC 822 or ISO 8601, no orphaned closing tags. If your RSS feed is broken, integration tools (Zapier, Slack/RSS, IFTTT, Beamer) silently drop you and you don't notice.
3. **The monthly digest is not the changelog repeated.** It is an editorial curation — top three to five items, each rewritten for an audience that won't read the long-form entries.
4. **The in-product panel is the shortest form.** The reader is mid-task; respect that. One line, one link, dismissable.

---

## 7. Monthly digest pattern (Mailchimp / Notion style)

A monthly digest answers: "What if I haven't checked the changelog in 30 days?"

Structure:

```text
Subject: What's new in [Product] — [Month Year]
Preview text: [One headline benefit, no "we're excited"]

Hero image (one feature, the month's flagship change)

[Month] in [Product]
[One opening paragraph — the editorial framing of the month.]

The big one
  [Flagship change, 2-3 sentences, screenshot or GIF, "Learn more" link]

Smaller wins
  • [Item, one line, link]
  • [Item, one line, link]
  • [Item, one line, link]

Coming soon
  [One or two items teased, if any. Optional. Only if you'll actually ship them.]

Quietly under the hood
  [One line about non-user-visible stability/perf work, optional. Builds trust without padding.]

[Footer: link to full changelog, unsubscribe, RSS link]
```

**What this format gets right:**

- **Editorial voice** — a human picked the order. The changelog page is chronological; the digest is curated.
- **One hero, several supporting items.** Mirrors the launch-post arc (one big thing per release) at monthly scale.
- **"Coming soon" is optional and conditional.** Only include things you will actually ship within the next 30–60 days. Tease something you cancel and the next digest loses trust.
- **"Under the hood" entry** acknowledges the invisible work. One sentence. Don't pad.

---

## 8. RSS-feed friendliness

Many users will never visit the changelog page. They subscribe via RSS, Slack/RSS integration, Beamer, AnnounceKit, or aggregator tools. To serve them:

1. **Absolute URLs for all images.** Relative URLs break the moment the feed is rendered outside your domain.
2. **Stable per-item GUIDs.** If you regenerate GUIDs on every deploy, every existing item re-broadcasts and integrations rate-limit you or unsubscribe.
3. **Valid dates** in RFC 822 (`Tue, 28 May 2026 14:00:00 +0000`) or ISO 8601 (`2026-05-28T14:00:00Z`). Ambiguous date strings are the most common reason feeds disappear from readers.
4. **Strip CSS, keep semantic HTML.** RSS readers strip styles anyway, and unusual elements often render as `<unknown>` blocks. Stick to `<p>`, `<h2>`–`<h4>`, `<ul>`, `<img>`, `<a>`, `<code>`, `<pre>`.
5. **Test in three readers.** Feedly (web), NetNewsWire (Mac), and a Slack `/feed` subscription cover most of the wild. Visual differences are surprisingly large.
6. **Don't backfill old items into the feed when redesigning.** A redesign that re-emits the last year of changelog as "new" will mass-resend every subscriber's notifications.

---

## 9. Templates

### 9.1 A single user-facing changelog entry

```markdown
## [DATE — e.g., May 28, 2026]

### [Benefit-led headline — what the user can now do]

[One- to three-paragraph narrative in a human voice. Lead with the
outcome. End with how to find or use the new thing if it's not
obvious from the screenshot.]

![Alt text describing the change, not the screenshot.](/changelog/2026-05-28-headline.png)

#### Also this week
- **Improved:** [User-visible improvement, one line.]
- **Improved:** [User-visible improvement, one line.]
- **Fixed:** [Symptom, not stack trace.]
- **Fixed:** [Symptom, not stack trace.]
```

### 9.2 In-product "What's New" panel item

```markdown
**[Benefit headline, six words or fewer]**

[One sentence, what's new for you.] [Learn more →]
```

### 9.3 Monthly digest email skeleton

```markdown
Subject: What's new in [Product] — [Month Year]
Preview: [Headline benefit. No "we're excited."]

# [Month] in [Product]

[One-paragraph editorial framing of the month.]

---

## The big one

[Flagship feature headline]

[2–3 sentences. Screenshot or GIF.]

[Try it →]

---

## Smaller wins

- **[Benefit]** — [one-line description] · [Link]
- **[Benefit]** — [one-line description] · [Link]
- **[Benefit]** — [one-line description] · [Link]

---

## Coming soon

[Optional. One or two items you will ship. Skip if uncertain.]

---

## Quietly under the hood

[One sentence about reliability/performance/security work.]

---

[See the full changelog →] · [Unsubscribe] · [RSS]
```

### 9.4 Translation worksheet — engineering input → user output

| Engineering input (from PR / commit) | User-facing output |
|---|---|
| `feat(api): add cursor pagination to /v1/projects` | (Internal-only — omit) |
| `feat(ui): keyboard shortcut Cmd+K opens command palette` | "Press Cmd-K from anywhere to jump to a project, file, or setting. No more clicking through nested menus." (Screenshot of the palette.) |
| `fix(billing): correct timezone display on invoice PDFs` | "Fixed: invoice timestamps now match your account timezone instead of UTC." |
| `chore(deps): bump react to 18.3` | (Omit) |
| `perf(search): debounce input by 120ms, switch to fuzzy match` | "Search now feels instant on long lists and finds results even when you mistype. Hold ⌘ to switch back to exact match." (GIF.) |
| `feat(integrations): Slack message-link unfurling` | "Paste a [Product] link into Slack and you'll see a rich preview with project name, owner, and status — no extra clicks." (Screenshot of an unfurled message.) |

### 9.5 Visual hierarchy decision tree

```text
Is this a user-visible change?
├── No  → Don't write it. Belongs in commit history.
└── Yes
    ├── Does it change how something looks or behaves on screen?
    │   ├── Big visual change         → Screenshot, large, above the fold.
    │   ├── Interaction change        → GIF or short video.
    │   └── Subtle tweak              → Optional thumbnail, below the fold.
    └── Backend with measurable effect?
        ├── Performance metric        → Optional chart or "before/after."
        └── Reliability fix           → Text-only.
```

---

## 10. Anti-patterns

**The version-number-only changelog.** "v1.42.0 — bug fixes and improvements." This is not a changelog; it's an empty timestamp. Either say what changed for the user, or skip the entry.

**The "we're excited to announce" opener.** Every entry. Cut. Open with the change.

**The commit-log dump.** Pasting `git log --oneline` into a webpage and calling it a changelog. Users do not know what `auth-svc` is; they want to know if they should change anything they do.

**The "improved performance" non-entry.** Either quantify it ("dashboard loads in under one second on slow connections") or omit it.

**The "FIXED / ADDED / CHANGED" labels on a user feed.** Those are Keep a Changelog labels — correct for developer changelogs, wrong for user changelogs. The user audience wants "New," "Improved," "Fixed" at most, and even those are optional when the headline is benefit-led.

**The 30-bullet wall of patch notes.** A user-facing changelog has a hero entry, three to seven supporting items, and stops. Game-patch-notes culture is a genre; if you're not shipping a game, don't import it.

**The screenshot of nothing.** A screenshot of an empty modal, a default state with no data, or a UI element with no context. Always shoot with realistic data, and crop tight.

**The "internal upgrade" entry.** "We upgraded our database to version 16." The user does not care. If the upgrade has a user-visible effect, write *that* — "lists with more than 10,000 items now load 5x faster." If it doesn't, omit it.

**The changelog with no author voice.** Bland, passive-voice, committee-written entries. A user changelog is a small product blog. Sound like a person.

**The relative-URL image in the RSS feed.** Works on the page, breaks the moment it leaves your domain.

**The "Coming Soon" entry that becomes a graveyard.** Promising features that slip quarter after quarter. Either ship or stop teasing.

**The "patch notes" tone in a B2B product.** Hyperbolic, slang-heavy game-patch tone in a serious productivity tool. Calibrate to the product's voice.

**The changelog buried four clicks deep.** If nobody can find it, nobody reads it. Link from the marketing nav, the footer, and the in-product help menu.

---

## 11. Decision heuristics

- **Audience is the integrator/developer** → wrong skill. Use `changelog-and-release-notes`.
- **Change has no user-visible effect** → don't write the entry.
- **Entry needs three sentences** → it's probably two entries.
- **Headline reads as a feature name** ("Search v2") → rewrite as a benefit ("Faster search across projects").
- **Big change, no image** → add one. The reader's eye expects a visual cue for "this is worth reading."
- **Small change, big image** → cut the image. Reserve images for the hero of the week/month.
- **Version number is in the headline** → move to footer or remove. The date is the version.
- **RSS feed item has a relative URL** → fix it now. Subscribers see broken images.
- **Three months in a row of "improvements and bug fixes"** → your changelog is dead. Either ship something or pause the changelog.
- **Monthly digest "Coming Soon" item already appeared last month** → cut it. You're padding.
- **Image has placeholder data (Lorem ipsum, "Acme Corp")** → reshoot with realistic data.

---

## 12. References

- Linear changelog — linear.app/changelog (gold standard: dated, illustrated, narrative).
- Vercel changelog — vercel.com/changelog (high-cadence, tag-grouped, image-heavy).
- Stripe Blog: Changelog — stripe.com/blog/changelog (curated cadence, product-area tags).
- Stripe Newsroom — stripe.com/newsroom (launch posts that pair with the changelog).
- GitHub Changelog — github.blog/changelog (post-per-feature, screenshot-rich).
- Notion changelog — notion.so/releases (mobile/desktop versioning, dated cadence).
- Mailchimp "What's New" — mailchimp.com/whats-new (consumer-product tone, monthly-digest pattern).
- Mintlify, "Five Changelog Principles from Best-in-Class Developer Brands" — mintlify.com/blog.
- Usersnap, "10 Inspiring Changelog Examples" — usersnap.com/blog/changelog-examples.
- Beamer, "Changelog or Release Notes?" — getbeamer.com/blog/changelog-or-release-notes.
- AnnounceKit, "Release Notes Best Practices (2026)" — announcekit.app/guides/release-notes-best-practices.
- Appcues, "Changelog vs. release notes" — appcues.com/blog/changelog-vs-release-notes.
- Screenhance, "How to Create Visual Changelogs That Users Actually Read" — screenhance.com/blog/changelog-visual-guide.
- Frill, "15 Changelog Formatting Best Practices" — frill.co/blog/posts/changelog-format.
- Roadmap.ai, "Product Changelog Best Practices" — theroadmapai.com/blog/product-changelog-best-practices.
- Easydesk, "Product Changelog Guide For 2026" — easydesk.app/blog/product-changelog.
- Olivier Lacan, *Keep a Changelog* — keepachangelog.com (sibling spec, developer-facing).
- Microsoft Writing Style Guide; Google developer documentation style guide — tone/voice references.

---

## 13. Cross-skill boundaries

- **`changelog-and-release-notes`** — the SIBLING skill. Developer-facing semver changelog, breaking-change announcements, migration guides, Keep a Changelog spec, Conventional Commits mapping. Pair with this skill: a release that needs both audiences gets *both* artifacts — a developer changelog entry *and* a user-facing entry. They are not the same document.
- **`release-blog-and-launch-narrative`** — the one-shot launch post. The user changelog calls back to the launch post (the hero item links to the launch blog); the launch post is not itself a changelog entry.
- **`sales-and-marketing-copy`** — conversion copy for landing pages and ads. A user changelog entry is not a sales pitch — it's a journal entry. Keep promotional language out of the changelog.
- **`executive-comms`** — the internal exec memo about a launch or release. Different audience, different artifact.
- **`writing-expert`** — the prose-craft pass (anti-AI-isms, nominalizations, sentence flow). Run after this skill produces the structural draft.
- **`technical-writing-craft`** — deeper reference on technical-writing principles. Use when an individual changelog entry needs developer-doc rigor.
