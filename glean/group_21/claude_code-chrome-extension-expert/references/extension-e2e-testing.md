<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `extension-e2e-testing` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: extension-e2e-testing
version: 1.1.0
category: developer
tags: [chrome-extension, testing, playwright, puppeteer, e2e, ci-cd, mv3, service-worker, content-script]
description: >-
  Chrome extension E2E testing with Playwright — launchPersistentContext, unpacked extension loading,
  popup/options/dashboard page testing, service worker inspection, content script verification,
  message interception fixtures, shadow DOM piercing, storage seeding/reset, and CI/CD with
  GitHub Actions + xvfb.
  TRIGGER: writing or debugging extension E2E tests, setting up Playwright for a Chrome extension,
  testing popup/options/dashboard surfaces, launchPersistentContext setup, extension CI/CD pipeline,
  service worker suspension in tests, content script injection testing.
  SKIP: unit tests for extension helpers (use testing-and-vitest-expert); non-extension web app E2E;
  Firefox/Safari extension testing (Playwright only supports Chrome extensions).
triggers:
  - "Write Playwright E2E tests for my Chrome extension"
  - "How do I load an unpacked extension in Playwright?"
  - "launchPersistentContext setup for extension testing"
  - "Test the popup page with Playwright"
  - "Resolve extension ID dynamically in tests"
  - "Seed chrome.storage before an E2E test"
  - "Service worker suspension causing flaky tests"
  - "Content script injection verification in Playwright"
  - "Chrome extension CI/CD with GitHub Actions"
  - "Intercept chrome.runtime.sendMessage in a test"
related_skills:
  - mv3-service-worker-expert
  - extension-message-bridge
  - chrome-dev
  - testing-and-vitest-expert
---

# Extension E2E Testing (Playwright)

Playwright is the standard tool for Chrome extension E2E testing. Extensions require a persistent browser context, Chromium-only execution, and special handling for the multi-surface architecture (popup, options, dashboard, content scripts, service worker).

**When not to use:** unit tests for pure JS helpers (use Vitest/jest-chrome); Firefox or Safari extension testing; non-extension web app testing.

## Key constraints

- Extensions only load via `launchPersistentContext` with `--load-extension` — `browser.newContext()` will not work.
- `channel: 'chromium'` enables extensions in headless mode; without it, headed mode is required.
- Extension IDs are dynamic — resolve at runtime from the service worker URL.
- MV3 service workers suspend after ~30s idle, affecting test timing.

## Extension surfaces

| Surface | Lifecycle | Access method | Key challenge |
|---------|-----------|---------------|---------------|
| Popup | Ephemeral | `chrome-extension://{id}/popup.html` | Closes on focus shift |
| Options page | Persistent tab | `chrome-extension://{id}/options.html` | Easiest to test |
| Dashboard/newtab | Persistent tab | `chrome-extension://{id}/dashboard.html` | May depend on storage |
| Content script | Injected into host page | Navigate to matching URL | Wait for injection |
| Service worker | Ephemeral, event-driven | `context.serviceWorkers()` | Suspends after ~30s |

## Setup

### playwright.config.js

```javascript
import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: false,       // extensions share state — do not run parallel
  timeout: 120_000,           // extension startup is slow
  projects: [{ name: 'chromium-extension' }],
  reporter: 'list',
  use: { trace: 'retain-on-failure' },
});
```

### Install

```bash
npm install -D @playwright/test
npx playwright install chromium
```

## Loading unpacked extensions

```javascript
import { chromium } from '@playwright/test';
import path from 'path';
import fs from 'fs';
import os from 'os';

const pathToExtension = path.resolve(__dirname, '../..');  // repo root
const userDataDir = fs.mkdtempSync(path.join(os.tmpdir(), 'ext-test-'));

const context = await chromium.launchPersistentContext(userDataDir, {
  channel: 'chromium',   // required for headless extension support
  args: [
    `--disable-extensions-except=${pathToExtension}`,
    `--load-extension=${pathToExtension}`,
    '--no-first-run',
    '--no-default-browser-check',
  ],
});
```

### Resolve extension ID from service worker

```javascript
let [serviceWorker] = context.serviceWorkers();
if (!serviceWorker) {
  serviceWorker = await context.waitForEvent('serviceworker');
}
const extensionId = new URL(serviceWorker.url()).host;
```

## Testing extension pages

### Popup

```javascript
const popup = await context.newPage();
await popup.goto(`chrome-extension://${extensionId}/src/popup/popup.html`);
await expect(popup.getByRole('heading', { name: 'Quick status' })).toBeVisible();
await expect(popup.getByRole('button', { name: 'Toggle Overlay' })).toBeDisabled();
```

### Options page

```javascript
const options = await context.newPage();
await options.goto(`chrome-extension://${extensionId}/src/options/options.html`);
await options.locator('#tracking-enabled').check();
await options.getByRole('button', { name: 'Save Changes' }).click();
await expect(options.locator('#options-status')).toContainText('Options saved.');
```

### Dashboard (with storage pre-seeding)

```javascript
const dashboard = await context.newPage();
await dashboard.goto(`chrome-extension://${extensionId}/src/dashboard/dashboard.html`);

await dashboard.evaluate(async () => {
  await chrome.storage.local.set({
    mca_tracking_state_v1: {
      trackedCases: { '01234567': { caseNumber: '01234567', status: 'Open' } },
    },
  });
});
await dashboard.reload();
await expect(dashboard.getByText('01234567')).toBeVisible();
```

### Cross-page navigation (popup → dashboard)

```javascript
const dashboardPagePromise = context.waitForEvent('page');
await popup.getByRole('button', { name: 'Open Dashboard' }).click();
const dashboardPage = await dashboardPagePromise;
await dashboardPage.waitForLoadState('domcontentloaded');
await expect(dashboardPage).toHaveURL(
  new RegExp(`chrome-extension://${extensionId}/src/dashboard/dashboard\\.html`)
);
```

## Service worker testing

### Wake before assertion

MV3 service workers suspend after ~30s. Wake before asserting background state:

```javascript
const page = await context.newPage();
await page.goto(`chrome-extension://${extensionId}/src/popup/popup.html`);
await page.evaluate(async () => {
  await chrome.runtime.sendMessage({ type: 'MCA_GET_AUTH_STATUS' });
});
// Now safe to evaluate on the service worker
const version = await serviceWorker.evaluate(() =>
  chrome.runtime.getManifest().version
);
```

### Handle "Service worker restarted" errors

```javascript
async function resilientEvaluate(sw, fn) {
  try {
    return await sw.evaluate(fn);
  } catch (err) {
    if (err.message.includes('Service worker restarted')) {
      return await sw.evaluate(fn);  // retry once
    }
    throw err;
  }
}
```

### Keep-alive for long tests

```javascript
// In beforeEach
const keepAlive = setInterval(async () => {
  try { await serviceWorker.evaluate(() => true); } catch { /* restarting */ }
}, 15_000);
// In afterEach
clearInterval(keepAlive);
```

## Intercepting runtime messages

### Stub sendMessage responses

```javascript
await page.addInitScript(() => {
  const orig = chrome.runtime.sendMessage.bind(chrome.runtime);
  chrome.runtime.sendMessage = async (message, ...args) => {
    if (message?.type === 'MCA_GET_AUTH_STATUS') {
      return { ok: true, result: { support: { ok: true }, glean: { ok: true } } };
    }
    return orig(message, ...args);
  };
});
```

### Simulate push messages (onMessage)

```javascript
await page.addInitScript(() => {
  const listeners = [];
  const origAdd = chrome.runtime.onMessage.addListener.bind(chrome.runtime.onMessage);
  chrome.runtime.onMessage.addListener = (fn) => { listeners.push(fn); return origAdd(fn); };
  window.__dispatchRuntimeMessage = (msg) => {
    for (const fn of listeners) fn(msg, {}, () => {});
  };
});

await page.evaluate((msg) => window.__dispatchRuntimeMessage(msg), {
  type: 'MCA_TRACKED_CASE_ANALYSIS_STATUS',
  payload: { caseNumber: '01234567', status: 'completed' },
});
```

## Content script testing

```javascript
const page = await context.newPage();
await page.goto('https://hub.corp.mongodb.com/case/01234567');
await page.waitForSelector('#mca-overlay-shell', { timeout: 10_000 });
await expect(page.locator('#mca-overlay-shell')).toBeVisible();

// Shadow DOM (Playwright pierces open shadow by default)
const shadowHost = page.locator('#mca-overlay-host');
await expect(shadowHost.locator('.overlay-panel')).toBeVisible();

// Closed shadow DOM — use evaluate
const text = await page.evaluate(() =>
  document.querySelector('#mca-overlay-host')?.shadowRoot?.querySelector('.overlay-title')?.textContent
);
```

## Fixtures

```javascript
// tests/e2e/fixtures.js
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { chromium, test as base } from '@playwright/test';

const repoRoot = path.resolve(__dirname, '../..');

export const test = base.extend({
  context: async ({}, use) => {
    const userDataDir = fs.mkdtempSync(path.join(os.tmpdir(), 'ext-test-'));
    const context = await chromium.launchPersistentContext(userDataDir, {
      channel: 'chromium',
      args: [
        `--disable-extensions-except=${repoRoot}`,
        `--load-extension=${repoRoot}`,
        '--no-first-run',
        '--no-default-browser-check',
      ],
    });
    await use(context);
    await context.close();
    fs.rmSync(userDataDir, { recursive: true, force: true });
  },

  extensionId: async ({ context }, use) => {
    let [sw] = context.serviceWorkers();
    if (!sw) sw = await context.waitForEvent('serviceworker');
    await use(new URL(sw.url()).host);
  },
});

export const expect = test.expect;
```

### Storage seeding helper

```javascript
export async function seedStorage(context, extensionId, data) {
  const page = await context.newPage();
  await page.goto(`chrome-extension://${extensionId}/src/options/options.html`);
  await page.evaluate(async (d) => { await chrome.storage.local.set(d); }, data);
  await page.close();
}
```

### Storage reset between tests

```javascript
test.afterEach(async ({ context, extensionId }) => {
  const page = await context.newPage();
  await page.goto(`chrome-extension://${extensionId}/src/options/options.html`);
  await page.evaluate(async () => {
    await chrome.storage.local.clear();
    await chrome.storage.session.clear();
  });
  await page.close();
});
```

## CI/CD

```yaml
# .github/workflows/e2e.yml
name: Extension E2E Tests
on: [push, pull_request]
jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: 20 }
      - run: npm ci
      - run: npx playwright install chromium --with-deps
      - run: xvfb-run --auto-servernum npx playwright test
      - if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: playwright-report
          path: playwright-report/
          retention-days: 7
```

`xvfb-run` provides a virtual display on Linux CI. Even with `channel: 'chromium'` some extension behaviors need a display server.

## Anti-patterns

| Anti-pattern | Fix |
|-------------|-----|
| Hardcoded extension ID | Resolve from `new URL(serviceWorker.url()).host` |
| `browser.newContext()` | Use `chromium.launchPersistentContext()` |
| `fullyParallel: true` | Extensions share `chrome.storage` — run serially |
| Multi-step popup flows | Test atomically; popup closes on focus shift |
| `waitForTimeout(3000)` | Use `waitForSelector` with explicit timeout |
| No storage reset between tests | Add `afterEach` clear or seed fresh state per test |
| No SW wake before assertion | Send a message to wake SW before `evaluate()` |

## Troubleshooting

**Extension not loaded / timeout waiting for service worker**
1. Verify path contains a valid `manifest.json`
2. Ensure `channel: 'chromium'` is set
3. Check manifest has `background.service_worker`

**Extension ID is undefined** — SW not started yet; use the defensive `waitForEvent('serviceworker')` pattern.

**Tests pass locally, fail in CI** — Add `xvfb-run`; increase timeouts to 120s+; verify Chromium version matches.

**`chrome.storage` empty in evaluate** — Navigate to extension page first; wait for `domcontentloaded` before calling `chrome.storage`.

**Content script not injecting** — Check `content_scripts.matches` covers the test URL; `run_at` value (`document_idle` vs `document_start`).
