<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `webapp-testing` skill.
> Sibling topics in this family are now reference files under the hubs (`software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: webapp-testing
description: >
  Playwright-based local web application testing and UI automation. Write Python Playwright
  scripts to verify frontend functionality, capture screenshots, inspect the DOM, and view
  browser console logs. Handles static HTML, dynamic single-page apps, and multi-server setups
  (backend + frontend). Includes the with_server.py helper for server lifecycle management.
  TRIGGER: user wants to test, verify, or automate a local web app; asking about Playwright
  scripts, browser automation, DOM inspection, screenshot capture, console log capture, or
  UI behavior verification for a locally running app.
  SKIP: testing libraries (Vitest, Jest, pytest) with no browser/Playwright component; remote
  or production web apps; non-Python Playwright usage (use playwright docs directly).
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags:
  - playwright
  - testing
  - web-automation
  - browser
  - python
  - ui-testing
keywords:
  - playwright
  - browser automation
  - DOM inspection
  - screenshot
  - networkidle
  - with_server.py
  - headless
  - web app testing
  - console logs
  - selectors
whenToUse:
  - Verifying that a local web app feature works correctly by running a Playwright script
  - Capturing a screenshot of a local web page for visual inspection
  - Inspecting the rendered DOM of a dynamic single-page app
  - Capturing browser console logs while interacting with a page
  - Testing a workflow that requires a backend server and a frontend dev server running together
  - Automating form submissions or button clicks on a local app
whenNotToUse:
  - Unit or integration testing with Vitest, Jest, or pytest (no browser) — use testing-and-vitest-expert
  - Remote or production web app automation — consider security and ToS implications first
  - E2E testing of the Chrome extension itself — use extension-e2e-testing
related_skills:
  - testing-and-vitest-expert
  - chrome-extension-expert
  - programming-languages
---

# Web Application Testing

Test local web applications using native Python Playwright scripts.

**Always run helper scripts with `--help` first** before reading source — these scripts are
large and will pollute your context window. Use them as black boxes.

## Decision Tree

```
User task → Is it static HTML?
    ├─ Yes → Read HTML file to identify selectors
    │         └─ Write Playwright script using discovered selectors
    │
    └─ No (dynamic app) → Is the server already running?
        ├─ No → python scripts/with_server.py --help
        │        Then use the helper + write Playwright script
        │
        └─ Yes → Reconnaissance-then-action:
            1. Navigate + wait for networkidle
            2. Take screenshot or inspect DOM
            3. Identify selectors from rendered state
            4. Execute actions with discovered selectors
```

## Starting a Server with with_server.py

**Single server:**
```bash
python scripts/with_server.py --server "npm run dev" --port 5173 -- python your_automation.py
```

**Multiple servers (backend + frontend):**
```bash
python scripts/with_server.py \
  --server "cd backend && python server.py" --port 3000 \
  --server "cd frontend && npm run dev" --port 5173 \
  -- python your_automation.py
```

**Automation script template** (server managed externally — include only Playwright logic):
```python
from playwright.sync_api import sync_playwright

with sync_playwright() as p:
    browser = p.chromium.launch(headless=True)  # Always headless
    page = browser.new_page()
    page.goto('http://localhost:5173')
    page.wait_for_load_state('networkidle')  # CRITICAL: wait for JS before inspection
    # ... your automation logic
    browser.close()
```

## Reconnaissance-Then-Action Pattern

1. **Inspect rendered DOM:**
   ```python
   page.screenshot(path='/tmp/inspect.png', full_page=True)
   content = page.content()
   page.locator('button').all()
   ```

2. **Identify selectors** from inspection results

3. **Execute actions** using discovered selectors

## Common Pitfall

Do NOT inspect the DOM before waiting for `networkidle` on dynamic apps — the DOM will be incomplete and selectors will fail.

## Best Practices

- Use `sync_playwright()` for all scripts
- Always close the browser when done
- Prefer descriptive selectors: `text=`, `role=`, CSS selectors, or IDs
- Add explicit waits: `page.wait_for_selector()` rather than fixed `page.wait_for_timeout()`

## Reference Files

- `examples/element_discovery.py` — discovering buttons, links, and inputs
- `examples/static_html_automation.py` — using `file://` URLs for local HTML
- `examples/console_logging.py` — capturing console logs during automation
- `scripts/with_server.py` — server lifecycle management (run with `--help` first)
