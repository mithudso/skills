<!-- hub-reference-banner -->
> **Reference file — part of the `da-3-data-acquisition-sampling` hub.** Formerly the standalone `da-3-2-5-web-scraping-crawling` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-3-2-5-web-scraping-crawling
description: |
  Web scraping and crawling as a primary data collection method in the data analysis lifecycle.
  Covers tool selection (Scrapy, Playwright, Puppeteer, BeautifulSoup, Firecrawl), handling
  JavaScript-rendered pages, anti-bot countermeasures, robots.txt compliance, legal/ethical
  frameworks (hiQ v LinkedIn, CFAA, GDPR, ToS), and data quality challenges (DOM drift,
  schema brittleness, selector fragility).

  TRIGGER: Use when the task involves collecting data from websites programmatically,
  building or evaluating a scraping pipeline, assessing legal risk of scraping a specific
  site, dealing with CAPTCHA or rate-limiting defenses, designing crawl schedules, parsing
  HTML/JSON responses, or diagnosing scraped data quality issues (missing fields, stale
  selectors, encoding errors).

  SKIP: Defer to `dom-scraping-resilience` for engineering resilience patterns (selector
  strategies, retry logic, graceful degradation, selector maintenance pipelines). Defer to
  `da-3-1-1-primary-vs-secondary` for the upstream question of whether web data is
  primary or secondary for a given study. Defer to `mongodb-kb` for storage of collected
  corpus data.

related_skills:
  - dom-scraping-resilience
  - da-3-2-collection-methods
  - da-3-data-collection-acquisition
  - da-3-1-1-primary-vs-secondary
  - da-11-ethics-and-privacy
---

# Web Scraping & Crawling (DA 3.2.5)

Web scraping is the automated extraction of data from web pages; crawling is the systematic
traversal of hyperlinks to discover pages at scale. Together they form a primary data
collection method when no official API exists, data must be continuously refreshed, or the
needed dataset spans many pages or domains [Apify legal guide, 2026].

---

## 1. Conceptual Framing

In the OSEMN and CRISP-DM data-lifecycle models, scraping sits in the **Obtain / Data
Understanding** phase — before cleaning, exploration, or modelling. Two distinct activities
are often conflated:

| Activity | Definition |
|---|---|
| **Scraping** | Extracting structured data from a known page or set of pages |
| **Crawling** | Discovering and queuing new URLs by following links (the "spider" pattern) |

Most production pipelines combine both: a crawler discovers URLs, a scraper extracts data
from each. The distinction matters for rate-limiting strategy and for legal analysis —
crawlers generate more server load and more ToS surface area than single-page scrapers.

---

## 2. Tool Landscape

### 2.1 Static HTML parsers

**BeautifulSoup** (Python) parses downloaded HTML using lxml or html.parser. Suitable when
pages are fully server-rendered (no JS hydration required). Fast and low-dependency.
Canonical use: `soup.select_one("div.price")`. No browser overhead; cannot interact with
SPAs.

### 2.2 Full-framework crawlers

**Scrapy** (Python) is the dominant open-source crawl framework for large-scale jobs. It
provides a request queue, middleware pipeline (custom headers, proxy rotation, deduplication),
item pipelines (validation, storage), and a scheduling model. Scrapy is stateless by default;
session/cookie handling requires custom middleware. Scrapy's `DOWNLOAD_DELAY` and
`AUTOTHROTTLE_ENABLED` settings are the primary levers for respectful rate control [Scrapy
docs]. Not suitable for JS-heavy pages without an integrated Playwright/Splash middleware.

### 2.3 Browser automation (JS-rendered pages)

**Playwright** (Python/Node) and **Puppeteer** (Node) launch real Chromium instances and
expose a full DOM after JavaScript execution. Both can handle:

- Single-page applications (React, Vue, Next.js) whose data is injected after initial load
- Infinite scroll: trigger scroll events and wait for network idle
- Auth flows: fill form inputs and persist cookies
- Waiting strategies: `wait_for_selector`, `wait_for_load_state("networkidle")`,
  `wait_for_timeout` as fallback for async data injection

Playwright is preferred over Puppeteer in new projects due to multi-browser support
(Chromium, Firefox, WebKit) and better async Python bindings [Firecrawl, 2025].

Browser automation is 10–50× slower than static parsing and triggers more detection signals.
Reserve it for pages where static parsing cannot reach the data [groupbwt, 2025].

### 2.4 Managed extraction APIs

**Firecrawl** is a hosted API that renders pages, handles anti-bot layers, and returns clean
markdown or structured JSON. It accepts a schema definition and uses AI extraction to map
DOM content to typed fields, eliminating hand-written CSS selectors. When a site redesigns,
Firecrawl's AI adapts automatically, reducing selector maintenance by up to 90% in reported
benchmarks [Firecrawl blog, 2025]. Useful when:

- Target site changes frequently (schema brittleness is the dominant cost driver)
- Team lacks scraping infrastructure
- LLM-ready markdown output is the desired format

Trade-offs: per-page API cost, third-party data handling, less control over crawl logic.

### 2.5 Tool selection matrix

| Criteria | BeautifulSoup | Scrapy | Playwright | Firecrawl |
|---|---|---|---|---|
| JS rendering | No | No (+ plugin) | Yes | Yes |
| Scale (pages/day) | Low | High | Medium | Medium |
| Selector maintenance | Manual | Manual | Manual | AI-assisted |
| Cost | Free | Free | Free | Per-page fee |
| Anti-bot bypass | None | Middleware | TLS fingerprint risk | Managed |

---

## 3. JavaScript-Rendered Pages

The majority of modern web pages use client-side rendering (CSR) via React, Next.js, or
Vue. The initial HTTP response contains a near-empty HTML shell; actual data is injected
after JS execution. Static scrapers see only the shell.

**Detection pattern:** Fetch the raw URL with `requests`; if the target data fields are
absent in the response but visible in the browser, JS rendering is required.

**Practical approaches:**

1. **Look for a backing API first.** React/Vue apps almost always call a REST or GraphQL
   endpoint that returns the same data as JSON. Open DevTools, go to the Network tab, and
   filter by XHR/Fetch. Scraping the API directly is faster, more stable, and legally
   equivalent to scraping the rendered page [Apify guide].

2. **Intercept XHR/fetch responses with Playwright `page.route()`** when you need the raw
   JSON before it reaches the DOM. Example: `page.route("**/api/products**", handler)`
   captures the JSON payload at the network layer, bypassing DOM parsing entirely.

3. **Use Playwright with `networkidle` wait** when no clean API exists and data is injected
   after multiple async requests.

---

## 4. Anti-Bot Countermeasures

### 4.1 Detection mechanisms

Modern defenses operate in layers [groupbwt, 2025; aimultiple, 2026]:

| Layer | Mechanism |
|---|---|
| **Network** | IP reputation, datacenter CIDR blocks, ASN blocklists |
| **TLS fingerprint** | JA3/JA4 hash of TLS handshake parameters; differs between `requests` and a real browser |
| **HTTP headers** | Missing or inconsistent `Accept-Language`, `Sec-Fetch-*`, `User-Agent` header profiles |
| **Browser fingerprint** | Canvas hash, WebGL renderer, font enumeration, navigator.webdriver flag |
| **Behavioral** | Mouse movement entropy, click timing, scroll velocity, inter-request cadence |
| **Honeypots** | CSS-hidden links or fields; only bots follow them. Triggering flags the session |
| **CAPTCHA** | hCaptcha, reCAPTCHA v3 (score-based, frictionless), Cloudflare Turnstile |

In July 2025, Cloudflare began blocking AI-agent crawlers by default, requiring explicit
opt-in from site operators [Cloudflare, 2025].

### 4.2 Mitigation strategies

- **Rotate residential proxies** to avoid datacenter IP blocks. Residential IPs are tied to
  real ISP subscribers and pass most IP-reputation checks.
- **Use real browser profiles** (Playwright with `--user-data-dir` and genuine UA strings)
  to match TLS and header fingerprints.
- **Randomize timing** between requests. Deterministic cadence (e.g., exactly 1.0s every
  request) is a strong bot signal.
- **Respect Retry-After headers** on HTTP 429 responses rather than exponential backoff
  alone.
- **Never automate CAPTCHA solving** without explicit site permission. Automated CAPTCHA
  solvers violate ToS universally and may trigger CFAA exposure in the US.

### 4.3 Rate limiting as an ethical control

Even when technically possible, sending hundreds of requests per second imposes server
cost on the target. Scrapy's `AUTOTHROTTLE_ENABLED` dynamically adjusts delay based on
server response latency, making it the closest approximation to a self-regulating polite crawler.
Manual `DOWNLOAD_DELAY = 1.0` is a baseline floor.

---

## 5. robots.txt

`robots.txt` is a plain-text file at `https://example.com/robots.txt` that specifies which
user-agents may crawl which paths. Syntax:

```
User-agent: *
Disallow: /private/

User-agent: Googlebot
Allow: /
Crawl-delay: 2
```

**Legal weight:** robots.txt has no binding legal force under US federal law. However:

- Courts have treated deliberate ignoring of robots.txt as evidence of bad faith in ToS
  breach and trespass-to-chattels claims [Apify legal guide, 2026].
- Under GDPR, scraping personal data from pages marked `Disallow` on a public-sector site
  could weaken a "legitimate interest" basis defense.
- Industry standard (Robots Exclusion Protocol, RFC 9309) requires all compliant crawlers
  to honor `Disallow` directives for their user-agent.

**Crawl-delay:** Honor this directive even when your scraper otherwise ignores robots.txt
restrictions; it signals the minimum server tolerance for automated traffic.

**Per-page signals:** `<meta name="robots" content="noindex,nofollow">` tags and
`X-Robots-Tag` HTTP response headers are per-page equivalents of robots.txt and carry the
same good-faith weight. Check both when deciding whether to index or follow links from a
page.

---

## 6. Legal and Ethical Frameworks

### 6.1 Computer Fraud and Abuse Act (CFAA) — US

The CFAA criminalizes access to a computer "without authorization" or "exceeding authorized
access." The landmark **hiQ Labs v. LinkedIn** litigation (2017–2022) clarified scope:

- **2019 (9th Circuit):** Accessing publicly available web pages does not constitute
  "unauthorized access" under the CFAA. hiQ was granted a preliminary injunction allowing
  continued scraping of public LinkedIn profiles [Ninth Circuit, 2019].
- **2022 (9th Circuit, on remand):** Reaffirmed the narrow CFAA interpretation. Public data
  access cannot trigger CFAA liability regardless of cease-and-desist [Jenner & Block, 2022].
- **Final outcome:** LinkedIn obtained summary judgment on a **breach of contract** theory
  (ToS violation), not CFAA. The parties settled with a consent judgment in late 2022
  [Privacy World, 2022].

Key practical takeaways:
- Scraping **public** pages: low CFAA risk, but ToS breach remains a civil liability vector.
- Scraping **behind a login**: authentication bypass triggers CFAA even for data you could
  see when logged in manually.
- Scraping after receiving a **cease-and-desist letter** and continuing to access the site
  creates strong "exceeding authorized access" arguments [Fordham IPLJ].

### 6.2 Terms of Service

ToS are contracts between the site operator and the user. Automated access typically
violates "no bots / automated access" clauses. Consequences:

- Civil litigation for breach of contract (injunction, damages)
- Account termination
- IP bans

ToS violation is not a CFAA crime for public pages, but is a valid civil cause of action
[hiQ v. LinkedIn final ruling]. Always review the ToS of target sites before production
scraping. If the ToS explicitly prohibits scraping, consider requesting data access or an
official API key instead.

### 6.3 GDPR and privacy law

Under GDPR, personal data (names, emails, professional profiles) collected by scraping is
subject to the same obligations as any other collected personal data:

- **Lawful basis required.** "Legitimate interests" is the most common basis claimed for
  public-data scraping, but it requires a balancing test against data subjects' interests.
- **Data minimisation.** Collect only fields necessary for the stated purpose.
- **Retention limits.** Don't store scraped personal data indefinitely.
- **Right to erasure.** Even if data was scraped from a public source, subjects may still
  exercise erasure rights under GDPR Article 17.

The CCPA in California applies similar principles to personal data of California residents.
Scraping at scale without a lawful basis has resulted in DPA enforcement actions in the EU
[tendem.ai, 2025].

### 6.4 Ethical principles beyond legal compliance

- **Do not scrape for the purpose of building a competing product** using a site's
  proprietary dataset — even if technically lawful, this is widely considered unethical and
  increases litigation risk.
- **Attribute data sources** in published work.
- **Disclose automated data collection** in research methodology sections.
- **Do not republish scraped personal data** unless a clear consent or public-interest basis
  exists.

---

## 7. Data Quality Challenges

### 7.1 DOM drift and selector fragility

HTML structure changes without versioning or notice. CSS class names generated by modern
bundlers (e.g., `class="sc-abc123"`) are arbitrary hashes that change on each deploy.
XPath selectors break when parent element hierarchy changes. Symptoms:

- Fields return `None` silently
- Scraped data is structurally valid but semantically wrong (selector matched wrong element
  after DOM shift)
- Volume drops after a site redesign without error logs

**Mitigation:** Use semantic selectors (`[data-testid="price"]`, ARIA labels, heading
text) where available. Instrument extraction with field-level null-rate monitoring; alert
when >5% of records have null in a required field. Consider AI-extraction APIs (Firecrawl)
when field drift is the dominant maintenance cost.

### 7.2 Schema brittleness

The implicit schema of scraped data (field names, types, nesting) is owned by the target
site and changes unilaterally. Downstream pipelines break when a field is renamed, moved
from a text node to an attribute, or removed.

**Mitigation:**
- Validate extracted records against a schema (Pydantic, JSON Schema, Zod) at ingestion.
- Store raw HTML alongside extracted records during development; reparse from raw when the
  extractor logic changes.
- Version extraction logic alongside the data — tag which extractor version produced each
  record.

### 7.3 Encoding and encoding errors

Non-UTF-8 pages, double-encoded entities, and mixed-encoding responses are common sources
of corrupted scraped data. Always specify `response.encoding` explicitly in `requests`
rather than relying on auto-detection. BeautifulSoup's `from_encoding` parameter handles
legacy charsets (ISO-8859-1, Windows-1252).

### 7.4 Pagination and completeness

Infinite scroll and cursor-based pagination require state management across requests.
Missing the termination condition results in duplicate records or incomplete datasets.
Track seen URLs in a set; validate final record count against a known total when the site
exposes one (e.g., "Showing 1–25 of 1,432 results").

### 7.5 Freshness and staleness

Scraped data captures a point-in-time snapshot. For time-sensitive analysis (pricing,
inventory, social metrics), establish a crawl cadence matched to the expected change rate
of the data, and stamp every record with a `scraped_at` timestamp. Consider delta scraping
(re-scraping only changed pages detected via `Last-Modified` or `ETag` headers) to reduce
volume and server load.

---

## 8. Crawl Architecture Patterns

### 8.1 Breadth-first vs depth-first

| Pattern | Use case |
|---|---|
| **BFS** | Site-wide coverage; discover all pages before going deep on any branch |
| **DFS** | Category pages with known depth; follow pagination chains fully before moving to next category |
| **Focused crawl** | Start from a seed set and only follow links that match URL patterns or page-content heuristics |

### 8.2 Sitemap-driven seed discovery

Before building a custom link-following crawler, check for `sitemap.xml` (typically at
`https://example.com/sitemap.xml` or declared in `robots.txt` as `Sitemap: <url>`). A
sitemap lists canonical URLs with optional `<lastmod>` timestamps, enabling:

- Targeted seed sets without crawling every link on every page
- Delta crawling: re-scrape only URLs whose `<lastmod>` changed since the last run
- Coverage validation: compare scraped URL count against the sitemap count to detect gaps

Scrapy's `SitemapSpider` base class parses sitemaps automatically and filters by URL
pattern via the `sitemap_rules` attribute.

### 8.3 Scrapy + Playwright hybrid

For large-scale crawls with mixed static/dynamic pages: Scrapy handles request queuing,
deduplication, and item pipelines; `scrapy-playwright` middleware dispatches to a Playwright
browser only for pages that require JS rendering (detected by URL pattern or response
content). This hybrid keeps browser instances to a minimum while retaining Scrapy's
scalability.

### 8.3 Crawl politeness checklist

- [ ] Honor `robots.txt` Disallow and Crawl-delay
- [ ] Set a descriptive `User-Agent` identifying your project and contact email
- [ ] Implement `AUTOTHROTTLE` or fixed `DOWNLOAD_DELAY`
- [ ] Respect `Retry-After` on 429/503 responses
- [ ] Limit concurrent requests per domain
- [ ] Log all requests with timestamps for post-hoc audit

---

## 9. Relationship to dom-scraping-resilience

This skill covers scraping as a **data collection method** — conceptual grounding, tool
selection, legal/ethical framing, and data quality principles. The companion skill
`dom-scraping-resilience` covers the **engineering resilience** angle: selector strategies,
retry and fallback logic, graceful degradation, selector health monitoring pipelines, and
automated re-extraction from raw HTML archives. Use both skills together for production
scraping system design.

---

## Sources

- [Apify Web Scraping Legal Guide (2026)](https://use-apify.com/blog/web-scraping-legal-guide)
- [Facing the Real Web Scraping Challenges in 2025 — GroupBWT](https://groupbwt.com/blog/challenges-in-web-scraping/)
- [hiQ v. LinkedIn Case Law — Apify Blog](https://blog.apify.com/hiq-v-linkedin/)
- [Ninth Circuit Holds Data Scraping is Legal in hiQ v. LinkedIn — California Lawyers Association](https://calawyers.org/privacy-law/ninth-circuit-holds-data-scraping-is-legal-in-hiq-v-linkedin/)
- [LinkedIn's Data Scraping Battle with hiQ Labs Ends — Privacy World](https://www.privacyworld.blog/2022/12/linkedins-data-scraping-battle-with-hiq-labs-ends-with-proposed-judgment/)
- [Jenner & Block: hiQ v. LinkedIn Ninth Circuit Reaffirms Narrow CFAA Interpretation](https://www.jenner.com/en/news-insights/publications/client-alert-data-scraping-in-hiq-v-linkedin-the-ninth-circuit-reaffirms-narrow-interpretation-of-cfaa)
- [Firecrawl vs Playwright for Web Scraping (2025)](https://www.firecrawl.dev/blog/playwright-vs-firecrawl)
- [Best Open-Source Web Scraping Libraries in 2026 — Firecrawl](https://www.firecrawl.dev/blog/best-open-source-web-scraping-libraries)
- [The Death of the Brittle Scraper: How Firecrawl Solves Schema Drift — Medium](https://medium.com/@raisrujan/the-death-of-the-brittle-scraper-how-firecrawl-is-solving-the-webs-hardest-data-problems-04b6f70341fa)
- [Is Web Scraping Legal? GDPR, CCPA & CFAA Frameworks Explained — tendem.ai](https://tendem.ai/blog/is-web-scraping-legal-compliance-overview)
- [Most Common Web Scraping Challenges in 2026 — AIM Multiple](https://aimultiple.com/web-scraping-challenges)
- [Data Scraping as a Cause of Action: Limiting Use of the CFAA — Fordham IPLJ](https://ir.lawnet.fordham.edu/cgi/viewcontent.cgi?article=1705&context=iplj)
