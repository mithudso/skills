<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-canva-founder-brand-stack` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-canva-founder-brand-stack
description: >-
  The DIY-founder design-tooling layer for a cause/nonprofit venture: running Canva as a one-person brand and
  marketing studio. Covers Canva tiers and the free Canva-for-Nonprofits program (501(c)(3) eligibility); the Brand
  Hub / Brand Kit; Magic Studio AI (as of 2026); content types (social, decks, docs, websites, video, PDFs); print
  collateral and the print-production fundamentals a non-designer must know (RGB vs CMYK, bleed/trim/safe margin,
  300 DPI, PDF Print export); the brand-to-templates-to-collateral-to-channel workflow; consistency systems;
  collaboration; and when Canva is the wrong tool. Design theory is cross-referenced, not duplicated. TRIGGER:
  questions about Canva, Brand Kit/Brand Hub, building or exporting marketing collateral or print files (flyers,
  cards, banners), Magic Studio AI, Canva-for-Nonprofits eligibility, print export settings (CMYK/bleed/DPI), or a
  DIY brand-asset workflow. SKIP: marketing strategy / local SEO -> venture-marketing-strategy-local-seo; nonprofit
  campaign or donor-journey specifics -> venture-cause-nonprofit-marketing; product-UI / app design -> frontend-ui;
  sentence-level copywriting -> content-and-marketing-writing; 501(c)(3) formation -> venture-nc-nonprofit-formation.
category: personal-venture
tags: [venture, canva, design, branding, marketing-collateral]
whenToUse:
  - Setting up or organizing a Canva account, folders, or Brand Kit for the venture
  - Deciding between Canva Free, Pro, Teams, or applying for Canva for Nonprofits
  - Building marketing collateral (social graphics, flyers, decks, one-pagers, PDFs) from a brand
  - Exporting a print-ready file (business card, postcard, banner) without color or bleed problems
  - Using Magic Studio AI features (Magic Design, Magic Write, Magic Switch, Background Remover) responsibly
  - Creating a reusable template/consistency system so every asset looks on-brand
  - Deciding whether a job belongs in Canva at all vs Figma or Adobe
  - Onboarding a volunteer or contractor into the venture's Canva workspace
triggers:
  - canva
  - brand kit
  - canva for nonprofits
  - print ready export
  - bleed and CMYK
  - magic studio
  - flyer or business card design
  - marketing collateral templates
---

# Canva and the DIY Founder Brand/Marketing Stack

This skill is the **design-tooling layer** for a solo founder running a North Carolina organ-donation-awareness
nonprofit who produces their own branding and marketing collateral. It treats Canva as a one-person design studio:
how to set it up, get it free as a nonprofit, build a reusable brand system, produce on-brand assets fast, export
files that actually print correctly, and know when to reach for a different tool.

**Scope guardrails.** This skill covers *the tool and the production craft*, not strategy or theory:

- **Marketing strategy, channels, local SEO** -> `venture-marketing-strategy-local-seo`.
- **Donor-journey design, cause-storytelling ethics, campaign specifics** -> `venture-cause-nonprofit-marketing`.
- **Color theory, typographic theory, layout/visual-hierarchy theory** -> the pro library's `frontend-ui`
  (ui-ux / web-design references). This skill gives you *practical defaults and Canva mechanics*, and points there
  for the why.
- **Sentence-level copy** (headlines, body, CTAs) -> `content-and-marketing-writing`. Canva's Magic Write drafts
  copy, but the craft lives there.
- **Product/app UI design** -> `frontend-ui`. **501(c)(3) formation** -> `venture-nc-nonprofit-formation`.

> **Verify before relying.** Canva ships features and changes pricing/eligibility constantly. Every feature, price,
> seat count, and eligibility rule below has an as-of date; treat it as a snapshot and confirm against `canva.com`
> before making a decision. See the Disclaimer.

---

## 1. Canva tiers + the nonprofit program (the most important decision)

As of 2026 Canva sells (consumer/SMB-facing): **Free ($0)**, **Pro (~US$15/mo or ~$120/yr for one person)**, and
**Teams (~$10/seat/mo, minimum seat counts apply)**, plus Enterprise. Pricing and plan names shift; verify on the
[Canva pricing page](https://www.canva.com/pricing/). [as of 2026-06]

| Capability | Free | Pro / Teams | Why it matters to you |
|---|---|---|---|
| Templates, stock assets, real-time collab | Yes (1.6M+ templates, 5GB) | Yes (100M+ premium assets, 1TB+) | Free is genuinely usable |
| **Brand Kit / Brand Hub** | Limited | **Yes** | The core of a consistent brand |
| **Background Remover** | No | **Yes** | Clean cut-outs of people/logos |
| **Magic Switch / Magic Resize** | No | **Yes** | Reformat one design to every channel |
| **CMYK color selection + print proofing** | No | **Yes** | Real print-production control (see Sec. 5) |
| Magic Studio AI (Magic Design, Magic Write, etc.) | Basic/limited | Full (usage limits apply) | See Sec. 3 |
| Brand Controls (lock colors/fonts, approval) | No | Teams/Business/Enterprise | Useful once volunteers join |

Sources: [Canva pricing](https://www.canva.com/pricing/), tier comparisons cross-checked across third-party
breakdowns. [as of 2026-06]

### Canva for Nonprofits — get Pro-tier features free

This is the headline win for a registered nonprofit. **Canva for Nonprofits gives an eligible organization free
access to all premium Canva features (Pro plus team/collaboration tools) for one team of up to 50 members.**
[Canva for Nonprofits](https://www.canva.com/canva-for-nonprofits/) and Canva's newsroom announcement of the
[50 free Canva for Teams seats](https://www.canva.com/newsroom/news/canva-for-nonprofits-seat-increase/). [as of 2026-06]

**Eligibility (US):**

- Your org must be a **recognized charitable nonprofit** that operates not-for-profit and exists for public benefit;
  in the US that means a **501(c)(3)** (corporation, trust, or other federally tax-exempt organization).
- The org must be **independent from government** and **not operating for commercial profit**.
- Source: [Canva for Nonprofits Help Center](https://www.canva.com/help/canva-for-nonprofits/). [as of 2026-06]

**Commonly *ineligible* org types** (even if charitably registered): organizations focused on legislation/politics,
government entities, education providers (K-12 schools, universities — those route to *Canva for Education* instead),
professional sports, financial services, business development, professional societies, job-training orgs, mutual
organizations, fraternities/sororities, and employee/membership-benefit groups. Confirm your mission isn't on the
exclusion list before applying. [as of 2026-06]

**Documentation Canva's vetting typically asks for** (have these ready — they map directly to artifacts produced in
`venture-nc-nonprofit-formation`):

1. IRS **EIN assignment letter**.
2. **Bank letter or voided deposit slip** in the organization's name.
3. **Signed/stamped governing documents** (articles of incorporation / bylaws) — including a statement of mission.
4. Proof you **work for or are authorized to represent** the org.
5. **Government-issued ID**.

[UNVERIFIED — exact document list] The specific documents requested can vary by country and over time and are
collected via Canva's third-party nonprofit-verification partner; the list above reflects what applicants commonly
report. Confirm the current requirements in the application flow at
[canva.com/canva-for-nonprofits](https://www.canva.com/canva-for-nonprofits/). [as of 2026-06]

> **Action for this venture:** apply for Canva for Nonprofits *before* paying for Pro. If approval lags, a single Pro
> seat (~$15/mo) unblocks Brand Kit + CMYK export in the meantime. Don't buy a full Teams plan you can get free.

---

## 2. Canva fundamentals (the editor and your file system)

**The editor.** A Canva design is a stack of **pages**, each holding **elements** (text, shapes, images, frames,
videos). Core moving parts:

- **Left panel:** Design (templates), Elements (shapes/lines/graphics/stickers/charts), Text, Brand, Uploads,
  Tools, Projects, Apps.
- **Top bar:** resize/Magic Switch, position/layers, animate, the **Share** button (where export lives), and
  Present/preview.
- **Frames vs. elements:** drop a photo into a **frame** to mask it to a shape; this keeps a clean, on-brand crop
  you can swap later. Use frames instead of manually cropping when you'll reuse the layout.
- **Grids and "set as background":** snap multiple images into a layout, or pin one image as a true background.

**Templates.** Start from a template, then *make it yours* with Brand Kit assets — never ship a stock template
unchanged (it will look generic and may collide with another org's collateral). Save your customized version as a
**Brand Template** so future assets start on-brand.

**Projects, folders, and naming.** Canva's **Projects** area holds all designs, folders, uploads, and Brand assets.
For a one-person venture, impose structure early so you (and any volunteer) can find things in a year:

```
Projects/
  00_BRAND/                 (logos, color/font references, brand guide PDF)
  01_TEMPLATES/             (master Brand Templates: social, flyer, deck, one-pager)
  02_SOCIAL/                (dated, by campaign: 2026-04_DonateLifeMonth/)
  03_PRINT/                 (flyers, cards, banners — print-ready masters + export PDFs)
  04_PRESENTATIONS_DOCS/
  05_WEB/                   (Canva websites, landing assets)
  99_ARCHIVE/               (retired collateral; keep for reference, out of the way)
```

Name designs `YYYY-MM_channel_purpose_vN` (e.g., `2026-04_IG_DonateLifeMonth_storyset_v2`). Folders + a naming
convention are the cheapest "design ops" you can buy.

---

## 3. Magic Studio (Canva's AI) as of 2026

Canva's AI is branded **Magic Studio**. As of 2026 it is *brand-aware* (it can pull from your Brand Kit) and the
core tools are bundled into Pro/Teams (and therefore into Canva for Nonprofits) rather than sold à la carte; **usage
limits/credits apply** on lower tiers. [as of 2026-06] Treat the specific tool roster as fluid — Canva renames and
reshuffles these often.

- **Magic Design** — describe what you want ("a flyer for Donate Life Month, teal, with a registration QR space")
  via text or voice and get fully editable layout options in seconds. Great for a fast first draft to then refine.
- **Magic Write** — in-canvas copy generation with tone control and 100+ languages. Use it to *draft* headlines/body,
  then edit for voice and accuracy (copy craft -> `content-and-marketing-writing`; cause-messaging ethics ->
  `venture-cause-nonprofit-marketing`).
- **Magic Media** — text-to-image and text-to-video generation (multiple visual styles). **Caution for a cause org:**
  do **not** use AI-generated images to depict real patients, recipients, donors, or medical situations as if they
  were real — it's misleading and erodes trust. Prefer authentic photography with consent for anything implying a
  real person/story.
- **Magic Switch / Magic Resize** — reformat one design into other sizes/formats (Instagram post -> Story -> Facebook
  -> flyer) and even transform a deck into a doc or summary. AI now relayouts intelligently rather than just cropping.
  This is your single biggest time-saver for multi-channel collateral (Pro feature). [as of 2026-06]
- **Background Remover** — one-click cut-out of subject from background for photos (and video, via Effects). Pro
  feature; the everyday workhorse for clean portraits and logo placement.
- Adjacent helpers you'll see: **Magic Edit/Grab/Eraser** (object-level photo edits), **Bulk Create** (generate many
  variants from a data list), **Magic Charts / Canva Sheets**, and **Brand Voice** (keeps generated copy on-tone).

Sources: Canva product pages for [Magic Resize](https://www.canva.com/pro/magic-resize/) and
[Brand Kit](https://www.canva.com/pro/brand-kit/), cross-checked against 2026 Magic Studio overviews. [as of 2026-06]

> **AI-honesty rule for a donation cause:** AI is fine for layout, drafting, resizing, and background cleanup. It is
> **not** fine for fabricating testimonials, faces of "recipients," or statistics. Donation messaging lives or dies on
> trust — see `venture-organ-donation-system` for the real facts and `venture-cause-nonprofit-marketing` for ethics.

---

## 4. Content types you can produce in Canva

One tool covers most of a small nonprofit's output:

- **Social graphics** — posts, Stories, Reels covers, carousels, profile/banner art. Use channel-correct sizes (or
  let Magic Switch generate them). There's a built-in **Content Planner** to schedule posts (Pro/Teams). [as of 2026-06]
- **Presentations** — board decks, partner pitches, volunteer training; present from Canva or export to PDF/PPTX.
- **Docs (Canva Docs)** — lightweight, visual documents (one-pagers, simple reports) as an alternative to Word.
- **Websites** — simple one-page sites / landing pages published on a Canva subdomain or a custom domain. Fine for a
  microsite or event page; **not** a substitute for a real CMS for an ongoing content site (see Sec. 8).
- **Video / Reels** — timeline editor with transitions, captions, stock clips, and background removal. Good for short
  social video; not a replacement for a pro NLE on complex edits.
- **PDFs** — two export flavors that matter: **PDF Standard** (screen/email, smaller) and **PDF Print** (300 DPI,
  CMYK, crop marks/bleed — see Sec. 5). Choosing the wrong one is the #1 print mistake.
- **Print products** — business cards, flyers, postcards, brochures, posters, banners, stickers, apparel — design in
  Canva and either **export a print-ready PDF for a local printer** or order through **Canva Print** fulfillment.

---

## 5. Print collateral + the print-production fundamentals (don't skip this)

A non-designer's print jobs go wrong in predictable ways. These are the fundamentals that prevent a wasted print run.

### RGB vs CMYK (color models)

Screens mix **RGB** (light). Presses mix **CMYK** (ink). Bright RGB blues, greens, and neons often **can't be
reproduced in CMYK** and shift duller/darker in print — this is **color drift**, and it's normal, not a Canva bug.

- **Canva designs in RGB.** On the **PDF Print** export, Canva performs an **automatic RGB -> CMYK conversion**, and
  the resulting PDF contains CMYK color data.
- **CMYK color *selection*** (picking a specific CMYK value, and the print-proofing view) is a **Pro / Teams /
  Education / Nonprofits feature** — another reason to get on the nonprofit plan before doing print.
- Source: Canva Help Center, *Proof your designs for print*
  (https://www.canva.com/help/proof-designs-print/). [as of 2026-06]

**Practical defaults:** avoid pure RGB neon for anything that must print; if a brand color is critical, ask your
printer for the target **CMYK** (or **Pantone/PMS**) values and set them, or order a physical proof. Pure black large
fills can look washed — printers often want a "rich black" build; ask them.

### Bleed, trim, and safe margin

Three nested boundaries on any cut piece:

- **Bleed** — artwork that extends *past* the cut line so a slightly off cut doesn't leave a white sliver. **Canva
  adds 0.125 in (3.175 mm) of bleed on all sides** for products that get trimmed. Backgrounds and edge-to-edge
  images **must extend into the bleed**, not stop at the edge.
- **Trim** — where the paper is actually cut (the finished size).
- **Safe zone / margin** — keep all critical text, logos, and the QR code at least **~5 mm (≈0.2 in)** *inside* the
  trim so nothing important gets clipped.

In Canva: turn on **File -> View settings -> Show margins** and **Show print bleed** (faint dashed lines). On export,
tick **"Crop marks and bleed."** Source: Canva Help Center, *Use margins, bleed, rulers, and crop marks*
(https://www.canva.com/help/margins-bleed-crop-marks/). [as of 2026-06]

### Resolution / DPI

Print needs **≥300 DPI at final size**. Two consequences:

1. **PDF Print exports at 300 DPI** — use it (not PNG/JPG, not PDF Standard) for anything going to a press.
2. **Don't blow up small images.** A 600 px logo dragged to fill a banner will look soft/pixelated. Source the
   highest-resolution logo and photos you have; for large-format (banners), 150 DPI can be acceptable because of
   viewing distance — confirm with the printer.

### The print-ready export checklist (run every time)

1. Canvas set to the **exact finished size** (or a Canva print template for that product).
2. Backgrounds/photos **bleed past the trim**; nothing critical outside the **safe margin**.
3. Images are **≥300 DPI at print size** (no "low resolution" warnings).
4. **Share -> Download -> PDF Print**, with **Crop marks and bleed** checked.
5. (Pro/Nonprofit) Use the **print proofing / CMYK** view; expect some color drift vs. screen.
6. Order a **physical proof** for important runs, or send the PDF to your printer and ask them to confirm before
   running. Fix in Canva and re-export — never edit the PDF by hand.

---

## 6. The brand-identity -> templates -> collateral -> channel workflow

This is the backbone: do it once, reuse forever.

1. **Define the brand identity** (decide *before* touching Canva). Logo, 2-4 brand colors with HEX values, 1-2 fonts
   (a heading + a body), and a one-line voice description. If you don't have these yet, *theory and decisions* live in
   `frontend-ui` (color/type) and your visual identity should reflect the cause tone in
   `venture-cause-nonprofit-marketing`. Canva's **Brand Kit Builder** can even extract colors/fonts/logos from an
   existing website or PDF to bootstrap the kit.
2. **Build the Brand Kit / Brand Hub** (Sec. 7). Load logos, the exact palette, fonts. This makes every future design
   one click from on-brand.
3. **Create master Brand Templates** for your recurring formats: one social-post template, one Story template, one
   flyer, one deck, one one-pager. Lock the look; vary only content. Save them in `01_TEMPLATES/`.
4. **Produce collateral** from those templates per campaign. Draft copy (Magic Write -> edit), drop in approved
   photos, keep type/colors from the kit.
5. **Adapt to every channel with Magic Switch**, then **export per destination** (PDF Print for press, PNG/MP4 for
   social, PDF Standard for email).

### Worked example: a "Donate Life Month" flyer + matching social set

Goal: an April awareness push — a printed flyer for a DMV/clinic counter **and** a coordinated social set, all on-brand.

1. **Brand first.** Brand Kit holds the logo, palette (e.g., teal + white + a warm accent), heading + body fonts.
   (Donate Life's own brand cues and the campaign timing come from `venture-organ-donation-system` /
   `venture-cause-nonprofit-marketing`.)
2. **Flyer master (print).** New design at the printer's flyer size (e.g., 5×7 in or A5). Turn on margins + bleed.
   Headline via Magic Write, edited for voice; one authentic, high-res, consented photo (Background Remover if you
   need a clean subject); a **QR code** (Canva's QR app) pointing to the donor-registration page — placed **inside
   the safe margin**; logo + URL in the footer. Background color bleeds to the edge.
3. **Export the flyer.** Share -> Download -> **PDF Print** + **crop marks and bleed**; check the CMYK proof; expect
   the teal to shift slightly; order a proof or confirm with the local printer.
4. **Social set via Magic Switch.** From the flyer (or a dedicated social master), generate an **Instagram/Facebook
   post (1080×1080)**, a **Story (1080×1920)**, and a **LinkedIn** size. Magic Switch relayouts; you nudge text and
   re-crop the photo. Keep the same headline, colors, and logo so print and social read as one campaign.
5. **Export social** as PNG (static) or MP4 (if animated); schedule via the Content Planner. Email version: **PDF
   Standard** or a Canva-built email graphic.
6. **Archive.** Save masters in `01_TEMPLATES/`, the campaign assets in `02_SOCIAL/2026-04_DonateLifeMonth/` and
   `03_PRINT/`. Next April, duplicate and update.

Result: one brand decision -> one set of templates -> a print piece + four social sizes, consistent and fast.

---

## 7. Brand Kit / Brand Hub — your consistency engine

The **Brand Kit** (in the newer **Brand Hub**) stores, in one place: **logos**, **brand colors** (with exact HEX
values), **brand fonts** (upload your own or pick from Canva's library), plus **photos, graphics, icons, chart
colors**, and a **Brand Voice**. The 2026 Brand System adds **custom categories, sections, and components** so a
larger asset set stays organized. Brand Kits are available on **Pro, Teams, Business, Enterprise, Education, and
Nonprofits**. Sources: Canva Help Center [Brand Kit](https://www.canva.com/help/brand-kit/) and
[Brand Kit best practices](https://www.canva.com/help/brand-kit-best-practices/); product page
[canva.com/pro/brand-kit](https://www.canva.com/pro/brand-kit/). [as of 2026-06]

**Set it up well:**

- **Upload logo variants:** full-color, white/reverse (for dark backgrounds), and a mark-only/icon version. Use PNGs
  with transparent backgrounds (or SVG if you have it).
- **Add the *exact* palette by HEX.** Don't eyeball colors per design — pull them from the kit every time. Note your
  print **CMYK/Pantone** equivalents somewhere (a notes page or your brand guide PDF) since presses need those.
- **Set heading + body fonts** so applying brand type is one click.
- **Brand Controls (Teams/Business/Enterprise):** lock designs to approved colors/fonts and require approval before
  publishing. Useful *once volunteers or a contractor* are editing — it stops off-brand drift. [as of 2026-06]
- **Brand Templates:** publish your master layouts as Brand Templates so collaborators start on-brand and can't
  wreck the structure.

---

## 8. Consistency systems, collaboration, and when *not* to use Canva

### Consistency systems (design ops for one person)

- **Brand Kit + Brand Templates** are the system. Everything starts from a template; content varies, look doesn't.
- **One source of truth per asset:** edit the master, then re-export — don't fork five near-identical copies.
- **Folder + naming convention** (Sec. 2) so you and future-you can find and reuse.
- A short **one-page brand guide** (logo usage, HEX + CMYK/Pantone, fonts, do/don'ts) stored in `00_BRAND/` —
  invaluable when a printer, partner, or volunteer asks "what's your blue?"

### Collaboration

Canva is multiplayer: **share a design or folder** with view/comment/edit roles, leave **comments/@mentions**, and on
Teams use **approval workflows + Brand Controls**. For a solo founder bringing in a volunteer or designer: share the
specific **folder**, give **edit** only where needed, and lean on Brand Templates + Brand Controls so help doesn't
become cleanup.

### When Canva is the *wrong* tool

Canva's strength is speed and accessibility for marketing collateral. Reach for something else when:

- **Product / app / web UI design and prototyping** -> **Figma**. Components, design systems, interactive prototypes,
  developer handoff. (Skill: `frontend-ui`.) Don't build a real app interface in Canva.
- **Advanced print / pro production** -> **Adobe InDesign** (multi-page, precise typography, imposition) or
  **Illustrator** (precise vector, logo *design* from scratch, scalable line art). When a commercial printer demands
  exact spot colors, complex prepress, or true vector, Canva's auto-conversion may not be enough.
- **Heavy photo retouching / complex vector** -> **Photoshop / Illustrator / Affinity** (Canva's nonprofit plan
  includes Affinity access). Canva's editing is light.
- **An ongoing content website / blog with a CMS, SEO control, forms, analytics** -> a real website platform; Canva
  websites are fine for a microsite or one-pager, not a growing site. (Strategy/SEO -> `venture-marketing-strategy-local-seo`.)
- **Long-form documents with rigorous formatting / citations** -> a word processor or layout tool, not Canva Docs.

A healthy pattern for many nonprofits: **Canva for 90% of day-to-day collateral**, a designer in Adobe/Figma for the
occasional high-stakes piece (the logo itself, a complex annual report, a billboard).

---

## Sources

Official Canva pages (authoritative for features/eligibility; verify live — Canva's help pages are JS-rendered and
may not load in simple fetchers):

- Canva for Nonprofits — landing + program: https://www.canva.com/canva-for-nonprofits/ [as of 2026-06]
- Canva for Nonprofits — Help Center (eligibility): https://www.canva.com/help/canva-for-nonprofits/ [as of 2026-06]
- Canva newsroom — 50 free Canva for Teams seats for nonprofits: https://www.canva.com/newsroom/news/canva-for-nonprofits-seat-increase/ [as of 2026-06]
- Canva nonprofit features: https://www.canva.com/nonprofits/features/ [as of 2026-06]
- Canva pricing: https://www.canva.com/pricing/ [as of 2026-06]
- Brand Kit — Help Center: https://www.canva.com/help/brand-kit/ ; best practices: https://www.canva.com/help/brand-kit-best-practices/ [as of 2026-06]
- Brand Kit — product page: https://www.canva.com/pro/brand-kit/ [as of 2026-06]
- Magic Resize / Magic Switch — product page: https://www.canva.com/pro/magic-resize/ [as of 2026-06]
- Proof your designs for print (RGB->CMYK conversion, CMYK is a Pro/Teams/Education/Nonprofits feature): https://www.canva.com/help/proof-designs-print/ [as of 2026-06]
- Margins, bleed, rulers, and crop marks (0.125 in / 3.175 mm bleed): https://www.canva.com/help/margins-bleed-crop-marks/ [as of 2026-06]
- Resize designs and size limits: https://www.canva.com/help/resize/ [as of 2026-06]

Cross-referenced third-party (used to corroborate, not as primary truth): Zeffy, Style Factory, and multiple 2026
Canva pricing/Magic-Studio breakdowns.

Related skills: `venture-cause-nonprofit-marketing`, `venture-organ-donation-system`,
`venture-nc-nonprofit-formation`, `venture-marketing-strategy-local-seo`, `frontend-ui`,
`content-and-marketing-writing`.

## Disclaimer

Canva changes features, plan names, pricing, AI-tool branding, and the Canva-for-Nonprofits eligibility rules and
required documents **frequently**. Everything here is a snapshot **as of 2026-06** and is for practical guidance, not
a guarantee. Before you apply for the nonprofit program, pay for a plan, or send a file to a printer, **verify the
current details on canva.com** (the pages in Sources) and, for print, **confirm specs with your printer and order a
physical proof**. Items marked **[UNVERIFIED]** could not be confirmed against a live official source at authoring
time. This skill does not provide legal, tax, or compliance advice — for 501(c)(3) status and charitable-solicitation
matters see `venture-nc-nonprofit-formation`.
