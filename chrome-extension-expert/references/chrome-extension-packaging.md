<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `chrome-extension-packaging` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: chrome-extension-packaging
description: Chrome extension packaging and distribution -- build pipelines (webpack/vite/CRXJS), manifest validation, Chrome Web Store publishing, version management, CI/CD workflows, code signing, privacy disclosure, and update mechanisms. Use when building, reviewing, or automating Chrome extension build/publish workflows.
version: 1.1.0
category: developer
tags: [chrome-extension, build, packaging, cws, ci-cd, publishing, manifest-v3]
whenToUse:
  - User asks how to package a Chrome extension for the Web Store
  - User needs a CI/CD pipeline for extension build/test/publish
  - User asks about manifest.json validation or required fields
  - User wants to automate CWS uploads via GitHub Actions
  - User asks about Chrome extension version bumping strategies
  - User needs to set up code signing or manage .pem keys
  - User asks about CWS review process, rejection reasons, or privacy disclosure
  - User asks about self-hosted extension updates vs CWS-managed updates
  - User is choosing a build tool (vite/webpack/esbuild) for an extension project
whenNotToUse:
  - User asks about extension runtime behavior (use mv3-service-worker-expert)
  - User asks about content script injection or DOM manipulation (use chrome-dev)
  - User asks about extension security review (use chrome-extension-security-reviewer)
  - User asks about E2E testing for extensions (use extension-e2e-testing)
  - General npm/library packaging unrelated to Chrome extensions (use code-packaging)
relatedSkills: [chrome-dev, mv3-service-worker-expert, chrome-extension-security-reviewer, extension-e2e-testing, code-packaging]
---

# Chrome Extension Packaging

## Overview

This skill covers the full lifecycle of packaging and distributing a Chrome extension:
from local build pipelines through manifest validation, ZIP packaging, Chrome Web
Store (CWS) submission, code signing, privacy disclosure, version management, and
CI/CD automation. All guidance targets Manifest V3 (MV2 was deprecated June 2025).

Key principle: the `manifest.json` is the single source of truth for the extension
runtime. Every build, validation, and publish step revolves around it.

### When to use this skill

- Packaging an extension for CWS upload or self-hosted distribution
- Setting up or reviewing a build pipeline (vite, webpack, esbuild) for an extension
- Validating manifest.json fields, permissions, or CSP
- Automating build/test/publish via CI/CD (GitHub Actions)
- Managing version bumps, code signing keys, or .pem files
- Preparing privacy disclosures and permission justifications for CWS review

### When NOT to use this skill

- Extension runtime behavior, service worker lifecycle --> use **mv3-service-worker-expert**
- Content script injection, DOM extraction, overlay patterns --> use **chrome-dev**
- Security audit of extension permissions or code --> use **chrome-extension-security-reviewer**
- E2E testing with Playwright for extensions --> use **extension-e2e-testing**
- General npm/library packaging unrelated to extensions --> use **code-packaging**

### Quick-start decision tree

```
Need to package a Chrome extension?
  |
  +-- Do you have a build step (bundler, framework)?
  |     |
  |     +-- Yes, new project --> Vite + CRXJS (see Build Pipeline #2)
  |     +-- Yes, legacy/complex --> Webpack (see Build Pipeline #4)
  |     +-- Yes, minimal (JS only) --> esbuild (see Build Pipeline #5)
  |     +-- Yes, custom layout --> Vite manual (see Build Pipeline #3)
  |
  +-- No build step (vanilla JS) --> Skip to ZIP Packaging
  |
  +-- Ready to publish?
        |
        +-- CWS public --> See CWS Publishing Process
        +-- CWS enterprise --> See Enterprise Publishing
        +-- Self-hosted --> See Code Signing & Update Mechanism
```

---

## Core Concepts

### Extension artifact types

| Artifact | Purpose | When used |
|----------|---------|-----------|
| Unpacked directory | Local development, `chrome://extensions` load | Dev loop |
| `.zip` | CWS upload, CI artifact | Publishing |
| `.crx` | Signed package for self-hosted distribution | Enterprise/sideload |
| `.pem` | RSA private key that determines extension ID | Signing, key continuity |

### Manifest V3 constraints that affect packaging

- No remote code execution -- all JS must be bundled in the package.
- Service workers replace persistent background pages (idle timeout ~30s).
- `content_security_policy.extension_pages` cannot allow `unsafe-eval` or remote script sources.
- `host_permissions` are declared separately from `permissions`.
- `web_accessible_resources` require explicit `matches` patterns.

---

## Build Pipeline Options

### 1. No build (vanilla JS)

Suitable for extensions with plain JS, HTML, and CSS -- no transpilation needed.
The repo root IS the extension directory. Load unpacked directly.

```
my-extension/
  manifest.json
  src/background/service-worker.js
  src/content/content-script.js
  src/popup/popup.html
  src/popup/popup.js
```

Validation: `node --check src/background/service-worker.js` for syntax errors.
Packaging: zip the directory excluding dev files (see ZIP Packaging below).

### 2. Vite + CRXJS (recommended for new projects)

CRXJS (`@crxjs/vite-plugin`) reads `manifest.json` as the entry point and
auto-discovers background scripts, content scripts, popup, options, and other pages.

```js
// vite.config.js
import { defineConfig } from 'vite';
import { crx } from '@crxjs/vite-plugin';
import manifest from './manifest.json';

export default defineConfig({
  plugins: [crx({ manifest })],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
});
```

Features: true HMR for popup/options/content scripts, automatic
`web_accessible_resources` generation, static asset handling.

Install: `npm i -D @crxjs/vite-plugin`

### 3. Vite (manual config, no CRXJS)

When you need full control or CRXJS does not fit your project shape.

```js
// vite.config.js
import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    rollupOptions: {
      input: {
        'service-worker': resolve(__dirname, 'src/background/service-worker.js'),
        popup: resolve(__dirname, 'src/popup/popup.html'),
        options: resolve(__dirname, 'src/options/options.html'),
        'content-script': resolve(__dirname, 'src/content/content-script.js'),
      },
      output: {
        entryFileNames: '[name].js',
        chunkFileNames: 'chunks/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash][extname]',
      },
    },
  },
});
```

Copy `manifest.json` into `dist/` during build (use `vite-plugin-static-copy`
or a custom plugin).

### 4. Webpack

Still viable for complex legacy setups. Use `copy-webpack-plugin` to bring
`manifest.json` and static assets into the output directory.

```js
// webpack.config.js
const CopyPlugin = require('copy-webpack-plugin');
const path = require('path');

module.exports = {
  mode: 'production',
  entry: {
    'service-worker': './src/background/service-worker.js',
    'content-script': './src/content/content-script.js',
    popup: './src/popup/popup.js',
  },
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: '[name].js',
    clean: true,
  },
  plugins: [
    new CopyPlugin({
      patterns: [
        { from: 'manifest.json' },
        { from: 'src/popup/popup.html', to: 'popup.html' },
        { from: 'icons', to: 'icons' },
      ],
    }),
  ],
};
```

### 5. esbuild (fast, minimal)

Good for service worker and content script bundling when you do not need
HTML processing.

```js
// build.mjs
import * as esbuild from 'esbuild';

await esbuild.build({
  entryPoints: [
    'src/background/service-worker.js',
    'src/content/content-script.js',
  ],
  bundle: true,
  outdir: 'dist',
  format: 'esm',
  target: 'chrome120',
  minify: true,
});
```

### Build tool selection guide

| Criterion | Vite+CRXJS | Vite manual | Webpack | esbuild |
|-----------|-----------|-------------|---------|---------|
| Zero config | Yes | No | No | No |
| HMR for extensions | Native | Manual | Plugin | No |
| HTML entry points | Auto | Plugin | Plugin | No |
| Build speed (2026) | Fast (Rolldown) | Fast | Slow | Fastest |
| Ecosystem maturity | Good | Excellent | Excellent | Good |
| Best for | New projects | Custom layouts | Legacy | Libraries/CLIs |

---

## Manifest Validation

### Required fields (extension will not load without these)

```json
{
  "manifest_version": 3,
  "name": "My Extension",
  "version": "1.0.0"
}
```

### Recommended fields for CWS submission

```json
{
  "description": "One-sentence summary (max 132 chars for store listing).",
  "icons": {
    "16": "icons/icon-16.png",
    "48": "icons/icon-48.png",
    "128": "icons/icon-128.png"
  },
  "permissions": ["storage", "activeTab"],
  "host_permissions": ["https://example.com/*"],
  "background": {
    "service_worker": "src/background/service-worker.js",
    "type": "module"
  },
  "content_security_policy": {
    "extension_pages": "script-src 'self'; object-src 'self'"
  },
  "minimum_chrome_version": "120"
}
```

### Validation script (CI-friendly)

```bash
#!/usr/bin/env bash
set -euo pipefail

MANIFEST="${1:-manifest.json}"

# JSON syntax check
python3 -m json.tool "$MANIFEST" > /dev/null

# Required fields
for field in manifest_version name version; do
  jq -e ".$field" "$MANIFEST" > /dev/null || {
    echo "FAIL: missing required field '$field'" >&2; exit 1
  }
done

# MV3 check
MV=$(jq '.manifest_version' "$MANIFEST")
[[ "$MV" == "3" ]] || { echo "FAIL: manifest_version must be 3, got $MV" >&2; exit 1; }

# Version format: 1-4 dot-separated integers
VERSION=$(jq -r '.version' "$MANIFEST")
[[ "$VERSION" =~ ^[0-9]+(\.[0-9]+){0,3}$ ]] || {
  echo "FAIL: invalid version format '$VERSION'" >&2; exit 1
}

# Icons present
for size in 16 48 128; do
  ICON=$(jq -r ".icons.\"$size\" // empty" "$MANIFEST")
  if [[ -n "$ICON" && ! -f "$ICON" ]]; then
    echo "WARN: icon file '$ICON' declared but missing" >&2
  fi
done

# Permissions audit -- flag high-risk permissions
HIGH_RISK=("debugger" "cookies" "webRequestBlocking" "<all_urls>")
for perm in "${HIGH_RISK[@]}"; do
  if jq -e ".permissions // [] | index(\"$perm\")" "$MANIFEST" > /dev/null 2>&1; then
    echo "WARN: high-risk permission '$perm' -- ensure CWS justification is ready"
  fi
done

echo "OK: manifest validation passed"
```

### Permissions audit checklist

- Request only the minimum permissions needed for core functionality.
- Move optional features behind `chrome.permissions.request()` (optional_permissions).
- Prefer `activeTab` over broad host patterns when click-triggered access suffices.
- Every permission declared in manifest.json must have a written justification
  ready for the CWS privacy practices tab.
- `<all_urls>` and `webRequest`/`webRequestBlocking` trigger manual review.

### Content Security Policy rules (MV3)

- Default: `script-src 'self'; object-src 'self'` -- sufficient for most extensions.
- Cannot use `unsafe-eval`, `unsafe-inline`, or remote script sources.
- Sandbox pages can use a more relaxed policy via `content_security_policy.sandbox`.
- Wasm: add `script-src 'self' 'wasm-unsafe-eval'` if you ship WebAssembly.

---

## CWS Publishing Process

### Prerequisites

1. **Developer account**: one-time $5 registration fee at https://chrome.google.com/webstore/devconsole
2. **Extension ZIP**: see ZIP Packaging below.
3. **Store listing assets**: 128x128 icon, at least one 1280x800 screenshot, promotional tile (440x280).
4. **Privacy policy URL**: required if extension uses any permissions or collects data.
5. **Permissions justification**: text explanation for every declared permission.

### Submission workflow

1. Upload ZIP via Developer Dashboard or API.
2. Fill out listing: description, category, screenshots, promo images.
3. Fill out privacy practices: data usage disclosures, permission justifications, Limited Use certification.
4. Submit for review.

### Review tracks

| Track | Mechanism | Typical duration |
|-------|-----------|-----------------|
| Fast (automated) | Malware scan, policy pattern matching | Under 1 hour |
| Manual | Human reviewer examines code and listing | 1-3 business days |
| Extended | Complex permissions, prior violations, surges | Up to 7+ days |

### Common rejection reasons (>70% of rejections)

1. **Obfuscated or minified code without source maps** -- provide readable source or source map.
2. **Over-broad permissions** without justification.
3. **Incomplete listing** -- missing screenshots, description, or privacy policy.
4. **Privacy policy mismatch** -- listing says "collects browsing data" but policy omits it.
5. **Unused permissions** declared in manifest.
6. **Remote code loading** -- MV3 forbids this entirely.

### Enterprise publishing (2026)

Organization publishing allows approval-gated private distribution across organizations.
Extensions still pass standard CWS review. Useful for internal enterprise tools.

---

## Version Management

### Chrome version format

```
major.minor.patch.build
```

- 1 to 4 dot-separated integers (no text, no semver prerelease labels).
- Each integer: 0-65535.
- CWS rejects uploads with a version <= the currently published version.

### Strategies

#### Manual bump (small projects)

Edit `manifest.json` version directly. Keep `package.json` version in sync.

#### Script-based auto-increment

```bash
#!/usr/bin/env bash
# bump-version.sh -- increment patch in manifest.json and package.json
set -euo pipefail

MANIFEST="manifest.json"
CURRENT=$(jq -r '.version' "$MANIFEST")
IFS='.' read -r MAJOR MINOR PATCH BUILD <<< "$CURRENT"
PATCH=$((PATCH + 1))
NEW_VER="${MAJOR}.${MINOR}.${PATCH}"

jq --arg v "$NEW_VER" '.version = $v' "$MANIFEST" > tmp.json && mv tmp.json "$MANIFEST"

if [[ -f package.json ]]; then
  jq --arg v "$NEW_VER" '.version = $v' package.json > tmp.json && mv tmp.json package.json
fi

echo "Version bumped: $CURRENT -> $NEW_VER"
```

#### Semantic release automation

Use `semantic-release` with a Chrome extension plugin:

```json
{
  "branches": ["main"],
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    ["@nicedoc/semantic-release-chrome", {
      "asset": "dist.zip",
      "distFolder": "dist"
    }],
    "@semantic-release/github"
  ]
}
```

#### Convention for 4-segment versions

```
major.minor.patch.build
  |      |      |     |
  |      |      |     +-- CI build number or channel indicator
  |      |      +-------- bug fixes
  |      +--------------- features / UX enhancements
  +---------------------- breaking changes / min Chrome version bump
```

---

## ZIP Packaging for CWS Upload

### What to include

- `manifest.json` (must be at ZIP root, not inside a subdirectory)
- All files referenced by manifest: service worker, content scripts, popup HTML/JS/CSS,
  options page, icons, `_locales/`, `web_accessible_resources`
- Any runtime dependencies bundled into the output

### What to exclude

```
.git/
.github/
.vscode/
.claude/
node_modules/
tests/
docs/
scripts/
*.md
*.pem
*.map (unless you want source maps for debugging)
.env
.eslintrc*
.prettierrc*
tsconfig.json
vite.config.*
webpack.config.*
package.json
package-lock.json
```

### Packaging script

```bash
#!/usr/bin/env bash
set -euo pipefail

NAME=$(jq -r '.name' manifest.json | tr ' ' '-' | tr '[:upper:]' '[:lower:]')
VERSION=$(jq -r '.version' manifest.json)
OUTFILE="${NAME}-v${VERSION}.zip"

# If using a build pipeline, zip from dist/
# If no build, zip from repo root with exclusions
if [[ -d dist ]]; then
  cd dist
  zip -r "../$OUTFILE" . -x "*.map"
  cd ..
else
  zip -r "$OUTFILE" . \
    -x ".git/*" ".github/*" ".vscode/*" ".claude/*" \
    -x "node_modules/*" "tests/*" "docs/*" "scripts/*" \
    -x "*.md" "*.pem" ".env" ".eslintrc*" ".prettierrc*" \
    -x "tsconfig.json" "vite.config.*" "webpack.config.*" \
    -x "package.json" "package-lock.json" \
    -x "*.swp" ".DS_Store"
fi

echo "Packaged: $OUTFILE ($(du -h "$OUTFILE" | cut -f1))"
```

### Size considerations

- No hard CWS size limit, but Chrome enforces a 4 GB per-file upload limit.
- Smaller packages review faster and install faster.
- Minify JS/CSS in production builds. Strip source maps unless needed for debugging.

---

## Code Signing & Keys

### How CRX signing works

1. First packaging generates an RSA key pair and saves the private key as `.pem`.
2. The extension ID is derived from a hash of the public key -- same key = same ID.
3. Every subsequent version must be signed with the same `.pem` to retain the ID.
4. If the `.pem` is lost, the extension ID changes and users must reinstall.

### Packaging a CRX from command line

```bash
# First time (generates key)
google-chrome --pack-extension=./dist

# Subsequent (reuse key)
google-chrome --pack-extension=./dist --pack-extension-key=./my-extension.pem
```

### GitHub Actions CRX signing

```yaml
- name: Sign CRX
  uses: nicedoc/sign-chrome-extension@v1
  with:
    pem: ${{ secrets.CRX_PEM_KEY }}
    src: dist
    out: my-extension.crx
```

### Key management best practices

- **Never commit `.pem` files** to version control.
- Store the PEM as a GitHub Actions secret or in a vault (1Password, AWS Secrets Manager).
- Back up the PEM in at least two secure locations -- loss means a different extension ID.
- For CWS-only distribution, you do not need a `.pem` -- CWS manages signing.
- For self-hosted distribution, the `.pem` is essential for update continuity.

---

## CI/CD Workflows

### GitHub Actions: build, test, validate

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm

      - run: npm ci

      # Manifest validation
      - name: Validate manifest.json
        run: |
          python3 -m json.tool manifest.json > /dev/null
          node -e "
            const m = require('./manifest.json');
            if (m.manifest_version !== 3) throw 'Must be MV3';
            if (!m.version.match(/^\d+(\.\d+){0,3}$/)) throw 'Bad version';
            console.log('Manifest OK: v' + m.version);
          "

      # Syntax check critical files
      - name: Syntax check
        run: |
          for f in src/background/*.js src/content/*.js; do
            node --check "$f" || exit 1
          done

      # Unit tests
      - name: Run tests
        run: npm test

      # Package ZIP artifact
      - name: Package extension
        run: |
          VERSION=$(node -p "require('./manifest.json').version")
          zip -r "extension-v${VERSION}.zip" . \
            -x ".git/*" "node_modules/*" "tests/*" ".github/*" \
            -x "*.md" "*.pem" ".env" ".claude/*"

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: extension-zip
          path: extension-v*.zip
```

### GitHub Actions: publish to CWS on tag

```yaml
# .github/workflows/publish.yml
name: Publish to Chrome Web Store

on:
  push:
    tags: ['v*']

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm

      - run: npm ci

      # Build (if applicable)
      - run: npm run build --if-present

      # Package
      - name: Create ZIP
        run: |
          cd dist
          zip -r ../extension.zip .

      # Upload and publish
      - name: Publish to CWS
        uses: mnao305/chrome-extension-upload@v5
        with:
          file-path: extension.zip
          extension-id: ${{ secrets.CWS_EXTENSION_ID }}
          client-id: ${{ secrets.CWS_CLIENT_ID }}
          client-secret: ${{ secrets.CWS_CLIENT_SECRET }}
          refresh-token: ${{ secrets.CWS_REFRESH_TOKEN }}
          publish: true
```

### CWS API credentials setup

1. Create a project in Google Cloud Console.
2. Enable the Chrome Web Store API.
3. Create OAuth 2.0 credentials (Desktop app type).
4. Generate a refresh token using the OAuth consent flow.
5. Store `client_id`, `client_secret`, `refresh_token`, and `extension_id` as GitHub Actions secrets.

### Version-gated publish safety

```yaml
# Verify version was actually bumped before publishing
- name: Check version bump
  run: |
    TAG_VERSION="${GITHUB_REF#refs/tags/v}"
    MANIFEST_VERSION=$(jq -r '.version' manifest.json)
    if [[ "$TAG_VERSION" != "$MANIFEST_VERSION" ]]; then
      echo "ERROR: tag v${TAG_VERSION} does not match manifest version ${MANIFEST_VERSION}"
      exit 1
    fi
```

---

## Privacy Practices Disclosure

### Data usage disclosure (CWS dashboard)

You must declare each data type your extension collects or accesses:

| Data type | Example triggers |
|-----------|-----------------|
| Personally identifiable info | User login, email collection |
| Health info | Health-related content access |
| Financial info | Payment processing |
| Authentication info | OAuth tokens, passwords |
| Personal communications | Email reading, message access |
| Location | Geolocation API usage |
| Web history | `tabs`, `webNavigation`, `history` permissions |
| User activity | Click tracking, usage analytics |
| Website content | DOM reading via content scripts |

### Limited Use certification

A binding declaration that you will:
- Only use data for the extension's user-facing features.
- Not transfer data for advertising or to data brokers.
- Not sell user data.
- Restrict human review of user data except for security/compliance/legal.

### Permissions justification template

```
Permission: storage
Justification: Persists user preferences (theme, enabled features) locally.
No data is transmitted externally.

Permission: activeTab
Justification: Reads the current tab URL to determine if the extension
should activate its overlay on supported pages. No browsing history is
stored or transmitted.

Permission: host_permissions (https://example.com/*)
Justification: Required to inject content scripts that extract case data
from the support portal. Data stays local to the extension.
```

---

## Update Mechanism

### CWS-managed updates (default)

- Extensions installed from CWS auto-update via Chrome's built-in mechanism.
- Chrome checks for updates every few hours.
- Developer has no control over update timing or rollout speed.
- No `update_url` field needed in manifest -- CWS handles it.

### Self-hosted updates (enterprise/sideload)

Add to manifest.json:
```json
{
  "update_url": "https://myhost.example.com/updates.xml"
}
```

Host an update manifest XML:
```xml
<?xml version='1.0' encoding='UTF-8'?>
<gupdate xmlns='http://www.google.com/update2/response' protocol='2.0'>
  <app appid='EXTENSION_ID_HERE'>
    <updatecheck codebase='https://myhost.example.com/my-extension.crx'
                 version='1.2.3' />
  </app>
</gupdate>
```

Chrome polls this URL every few hours. If the XML version > installed version,
Chrome downloads and installs the new CRX.

### Choosing an update strategy

| Factor | CWS-managed | Self-hosted |
|--------|------------|-------------|
| Review required | Yes | No |
| User trust | Higher (CWS badge) | Lower |
| Update speed | Hours (after review) | Minutes |
| Enterprise control | Via CWS policies | Full control |
| Key management | CWS handles | You manage .pem |
| Auto-install via policy | Supported | Supported |

---

## Anti-Patterns

1. **Committing `.pem` files** -- exposes signing key; attacker can publish malicious updates.
2. **Using `<all_urls>` when `activeTab` suffices** -- triggers manual review, erodes user trust.
3. **Minifying without source maps in dev** -- impossible to debug CWS rejection feedback.
4. **Hardcoding version** in multiple places -- manifest.json should be the single source; derive others.
5. **Skipping manifest validation in CI** -- broken manifests waste CWS review cycles.
6. **Uploading node_modules in the ZIP** -- massively inflated package, guaranteed slow review.
7. **Using dynamic code evaluation** -- MV3 CSP blocks eval and dynamic code construction; extension will not load.
8. **Forgetting to bump version before CWS upload** -- API rejects same-or-lower version.
9. **No privacy policy for permission-using extensions** -- immediate CWS rejection.
10. **Declaring permissions you do not use** -- reviewers flag unused permissions.

---

## Troubleshooting

### "Invalid package" on CWS upload
- Ensure `manifest.json` is at the ZIP root, not nested in a subdirectory.
- Validate JSON syntax: `python3 -m json.tool manifest.json`.
- Confirm `manifest_version` is `3`.

### Extension not loading after unpack
- Check `chrome://extensions` for error banners.
- `node --check <file>` every JS file referenced by manifest.
- Verify all icon files exist at declared paths.
- Confirm service worker path matches `background.service_worker`.

### CWS review taking >3 days
- Check for high-risk permissions (`debugger`, `cookies`, broad host patterns).
- Ensure no obfuscated code (minified is fine if readable patterns remain).
- Verify privacy policy URL resolves and covers declared data types.
- Check the developer dashboard for specific rejection feedback.

### Version rejected by CWS API
- Published version must be strictly less than uploaded version.
- Use the version check script above before uploading.
- If testing, use the `.build` segment (e.g., `1.0.0.1`, `1.0.0.2`).

### Self-hosted update not applying
- Verify `update_url` in manifest points to a reachable XML file.
- Confirm XML `version` attribute is higher than the installed version.
- Check that `appid` matches the extension's actual ID.
- Ensure the CRX file at `codebase` URL is downloadable and correctly signed.

---

## Related Skills

| Skill | Use for |
|-------|---------|
| `mv3-service-worker-expert` | Service worker lifecycle, idle termination, event-driven activation |
| `chrome-dev` | Content scripts, DOM interaction, extension runtime APIs |
| `chrome-extension-security-reviewer` | Security audit of permissions, CSP, code patterns |
| `extension-e2e-testing` | Playwright-based E2E testing for extension pages |
| `code-packaging` | General library/module packaging (npm, PyPI) unrelated to extensions |
| `code-reviewer` | Code review patterns applicable to extension PRs |

---

## References

- [Chrome Extensions Manifest V3 docs](https://developer.chrome.com/docs/extensions/mv3/)
- [CWS Publishing guide](https://developer.chrome.com/docs/webstore/publish)
- [CWS Review process](https://developer.chrome.com/docs/webstore/review-process)
- [CWS Program policies](https://developer.chrome.com/docs/webstore/program-policies)
- [CWS Privacy practices tab](https://developer.chrome.com/docs/webstore/cws-dashboard-privacy)
- [Manifest version field](https://developer.chrome.com/docs/extensions/reference/manifest/version)
- [Manifest CSP field](https://developer.chrome.com/docs/extensions/reference/manifest/content-security-policy)
- [Declare permissions](https://developer.chrome.com/docs/extensions/develop/concepts/declare-permissions)
- [CRXJS Vite plugin](https://github.com/crxjs/chrome-extension-tools)
- [chrome-extension-upload GitHub Action](https://github.com/marketplace/actions/chrome-extension-upload-action)
- [mnao305/chrome-extension-upload](https://github.com/marketplace/actions/publish-chrome-extension-to-chrome-web-store)
- [Extension versioning strategies](https://davestewart.co.uk/blog/extension-versioning/)
- [CWS Enterprise publishing (2026)](https://www.adwaitx.com/chrome-web-store-enterprise-publishing/)
- [Automating Chrome Extension Releases with CI/CD](https://zenn.dev/atani/articles/chrome-extension-auto-publish-guide)
- [CWS Policy updates 2025](https://developer.chrome.com/blog/cws-policy-updates-2025)
