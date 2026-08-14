<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `template-config-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: template-config-patterns
version: 1.1.0
updated: 2026-05-29
description: >
  Configurable report and prompt template patterns for JavaScript operator tools, Chrome
  extensions, and browser-based UIs — Mustache-style {{var}} interpolation with a zero-dependency
  regex renderer, template definition/content separation, token builder with fallbacks, catalog
  and map construction for UI dropdowns, user-customizable template libraries with normalization
  and ID deduplication, prompt override pattern, chrome.storage integration with save/load,
  placeholder validation, conditional sections, schema versioning, and live preview wiring.
  TRIGGER: building or reviewing template interpolation, variable substitution, configurable
  report/prompt templates, template rendering in a Chrome extension or browser UI, token builders,
  template catalog dropdowns, user-editable template libraries, template versioning/migration,
  or live preview with debounce.
  SKIP: backend templating with Handlebars/Nunjucks/EJS (this skill covers browser/CSP-safe
  patterns only), markdown rendering (use `markdown-rendering-browser`), LLM prompt engineering
  (use `prompt-engineering`), report generation for TAM accounts (use `tam-account-reports`).
category: developer
tags:
  - template-engine
  - interpolation
  - chrome-extension
  - operator-tools
  - mustache
  - user-customization
  - csp-safe
triggers:
  - template interpolation
  - variable substitution
  - configurable templates
  - report templates
  - prompt templates
  - template rendering
  - template tokens
  - token builder
  - template catalog
  - custom template library
  - template preview
  - template validation
  - Mustache pattern
  - template versioning
  - placeholder
related_skills:
  - chrome-storage-patterns
  - operator-report-generator
  - case-tracker
  - llm-integration-reviewer
  - vanilla-js-ui-reviewer
---

# Template Configuration Patterns

Expert reference for building configurable report and prompt templates with variable interpolation in JavaScript operator tools, Chrome extensions, and browser-based UIs.

**Canonical implementation:** `src/shared/template-config.js` in the mdb-case-assistant repo. Tests: `tests/unit/template-config.test.js`.

## 1. Template Engine Selection

### When to use which syntax

| Syntax | Use when | Avoid when |
|---|---|---|
| `{{var}}` (Mustache/Handlebars) | User-facing templates stored as plain strings; CSP-restricted environments (Chrome extensions); templates that travel across contexts (storage, postMessage, service worker) | You need expression evaluation inside the template |
| `${var}` (JS template literals) | Compile-time interpolation inside JS source code; developer-facing strings that never leave the module | Templates are user-editable strings loaded at runtime (template literals require `eval` or `new Function` to render dynamic strings) |
| Tagged template literals | You own both the template and the renderer; you need custom escaping or transforms | Templates are persisted to storage or edited by non-developers |

**Rule of thumb:** If the template is stored as a string (in chrome.storage, a JSON file, a database, or a textarea), use `{{var}}` with a regex renderer. If the template is inline JS code, use native template literals.

### Micro-Mustache pattern (zero dependencies)

```js
// Core renderer — ~4 lines, no dependencies, CSP-safe
function renderTemplate(template = '', tokens = {}) {
  return String(template || '').replace(
    /\{\{\s*([a-zA-Z0-9_]+)\s*\}\}/g,
    (_match, key) => {
      const value = tokens?.[key];
      return value == null ? '' : String(value);
    }
  );
}
```

Key design decisions:
- Whitespace tolerance: `{{ var }}` and `{{var}}` both work
- Key charset limited to `[a-zA-Z0-9_]` — prevents injection via crafted keys
- Null/undefined collapse to empty string rather than leaving the placeholder
- No nested property access (intentional — keeps the attack surface minimal)

## 2. Template Definition Architecture

### Definition vs content separation

Keep template metadata (key, label, description) separate from template content. This enables catalog UIs, dropdown selectors, and storage-efficient updates.

```js
// Template definitions — metadata only
const TEMPLATE_DEFINITIONS = Object.freeze([
  { key: 'ack',      label: 'Acknowledge + summarize',
    description: 'Short acknowledgement with working hypothesis.' },
  { key: 'diag',     label: 'Ask for diagnostics',
    description: 'Requests the next needed evidence.' },
  { key: 'escalate', label: 'Escalation update',
    description: 'Internal escalation summary.' },
  { key: 'followup', label: 'Follow-up nudge',
    description: 'Follow-up for pending actions.' },
]);

// Template content — keyed by definition key
const DEFAULT_TEMPLATES = Object.freeze({
  ack: [
    'Hi team,',
    '',
    'Thanks for the update on case {{caseNumber}}.',
    'Current hypothesis: {{hypothesis}}',
    '',
    '{{nextActionsBullets}}',
  ].join('\n'),
});
```

### Placeholder registry

Declare all supported placeholders in one place. This powers the options UI reference list, validation, and documentation.

```js
const TEMPLATE_PLACEHOLDERS = Object.freeze([
  '{{caseNumber}}',
  '{{accountName}}',
  '{{severity}}',
  '{{status}}',
  '{{owner}}',
  '{{hypothesis}}',
  '{{executiveSummary}}',
  '{{nextActionsBullets}}',
  '{{blockersBullets}}',
  '{{currentTimestamp}}',
  '{{caseUrl}}',
]);
```

## 3. Token Building

Tokens are the runtime values that replace placeholders. Build them from structured context with explicit fallbacks.

```js
function buildTemplateTokens(context = {}) {
  return {
    caseNumber:    context.caseNumber || '[Case number]',
    accountName:   context.accountName || '[Customer name]',
    severity:      context.severity || '[Severity]',
    status:        context.status || '[Status]',
    owner:         context.owner || '[Name]',
    hypothesis:    context.hypothesis || '[root cause hypothesis]',

    // Derived tokens — computed from multiple fields
    executiveSummary: firstNonEmpty(
      context.executiveSummary,
      context.hypothesis,
      '[2-4 sentences]'
    ),

    // List tokens — arrays rendered as bullet lists
    nextActionsBullets: bulletList(
      context.nextActions,
      '[next action placeholder]'
    ),
    blockersBullets: bulletList(
      context.blockers,
      '[waiting on customer / engineering / internal team]'
    ),

    currentTimestamp: formatTimestamp(),
    caseUrl: context.pageUrl || '[link]',
  };
}
```

### Helper functions

```js
function firstNonEmpty(...values) {
  return values.map(v => String(v || '').trim()).find(Boolean) || '';
}

function bulletList(items = [], fallback = '[placeholder]') {
  return (Array.isArray(items) && items.length ? items : [fallback])
    .map(item => `- ${item}`)
    .join('\n');
}

function formatTimestamp(value = new Date()) {
  const date = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(date.getTime())) return '[YYYY-MM-DD HH:MM TZ]';
  return `${date.toISOString().slice(0, 16).replace('T', ' ')} UTC`;
}
```

## 4. Template Catalog and Map

### Catalog (for UI dropdowns)

```js
function buildTemplateCatalog(options = {}) {
  const builtIns = TEMPLATE_DEFINITIONS.map(def => ({
    ...def, builtin: true,
  }));
  const custom = normalizeCustomTemplates(options.customTemplates || [])
    .map(entry => ({
      key:         entry.id,
      label:       entry.name,
      description: 'User-defined template.',
      builtin:     false,
    }));
  return builtIns.concat(custom);
}
```

### Map (for rendering)

```js
function buildTemplateMap(options = {}) {
  const builtIns = normalizeTemplates(options.templates || {});
  const custom = normalizeCustomTemplates(options.customTemplates || []);
  return {
    ...builtIns,
    ...Object.fromEntries(custom.map(e => [e.id, e.content])),
  };
}
```

### End-to-end rendering

```js
function insertTemplate(templateKey, options, context) {
  const templates = buildTemplateMap(options);
  if (!templateKey || !templates[templateKey]) return '';
  const tokens = buildTemplateTokens(context);
  return renderTemplate(templates[templateKey], tokens);
}
```

## 5. User-Customizable Template Libraries

### Normalization pipeline

```js
function normalizeTemplateText(value = '') {
  return typeof value === 'string'
    ? value.replace(/\r\n/g, '\n').trim()
    : '';
}

function normalizeLibraryId(value = '', prefix = 'entry') {
  const raw = String(value || '').trim().toLowerCase()
    .replace(/[^a-z0-9_-]+/g, '-')
    .replace(/^-+|-+$/g, '');
  if (!raw) return '';
  return raw.startsWith(`${prefix}-`) ? raw : `${prefix}-${raw}`;
}

function normalizeLibraryEntries(value = [], prefix = 'entry') {
  const entries = Array.isArray(value) ? value : [];
  const usedIds = new Set();
  const result = [];

  for (const entry of entries) {
    const name = normalizeTemplateText(
      entry?.name || entry?.label || entry?.title || ''
    );
    const content = normalizeTemplateText(
      entry?.content || entry?.body || entry?.prompt || entry?.template || ''
    );
    if (!name || !content) continue;

    let id = normalizeLibraryId(entry?.id, prefix)
          || normalizeLibraryId(name, prefix);
    let suffix = 2;
    const baseId = id;
    while (usedIds.has(id)) { id = `${baseId}-${suffix++}`; }
    usedIds.add(id);

    result.push({ id, name, content });
  }
  return result;
}
```

### Dual library pattern

```js
function normalizeCustomPromptLibrary(value = []) {
  return normalizeLibraryEntries(value, 'custom-prompt');
}

function normalizeCustomDraftTemplates(value = []) {
  return normalizeLibraryEntries(value, 'custom-template');
}
```

## 6. Prompt Override Pattern

```js
const DEFAULT_ANALYSIS_PROMPT = [
  'You are analyzing a live customer support case.',
  'Use only the supplied case text and metadata.',
  '',
  'Rules:',
  '- Keep the tone flat and technical.',
  '- If the case does not provide enough info, use empty strings.',
].join('\n');

function normalizePromptOverride(value = '', fallback = '') {
  return normalizeTemplateText(value) || fallback;
}

// Usage
const prompt = normalizePromptOverride(
  userProvidedOverride,
  DEFAULT_ANALYSIS_PROMPT
);
```

## 7. Storage Integration (Chrome Extensions)

### Durable storage layout

```js
const STORAGE_DEFAULTS = {
  draftTemplates:       DEFAULT_TEMPLATES,
  customPromptLibrary:  [],
  customDraftTemplates: [],
  caseAnalysisPromptOverride: '',
  researchPromptOverride:     '',
  deepDivePromptOverride:     '',
};
```

### Save with normalization

```js
async function saveTemplateOptions(payload) {
  const next = {};
  next.caseAnalysisPromptOverride = normalizePromptOverride(
    payload.caseAnalysisPromptOverride,
    DEFAULT_ANALYSIS_PROMPT
  );
  next.customPromptLibrary = normalizeCustomPromptLibrary(
    payload.customPromptLibrary || []
  );
  next.customDraftTemplates = normalizeCustomDraftTemplates(
    payload.customDraftTemplates || []
  );
  next.draftTemplates = normalizeTemplates(payload.draftTemplates || {});
  await chrome.storage.local.set(next);
  return next;
}
```

### Load with defaults

```js
async function loadTemplateOptions() {
  const stored = await chrome.storage.local.get(STORAGE_DEFAULTS);
  return {
    draftTemplates: normalizeTemplates(stored.draftTemplates),
    customPromptLibrary: normalizeCustomPromptLibrary(stored.customPromptLibrary),
    customDraftTemplates: normalizeCustomDraftTemplates(stored.customDraftTemplates),
  };
}
```

## 8. Template Validation

```js
function extractPlaceholders(template = '') {
  const matches = [];
  const re = /\{\{\s*([a-zA-Z0-9_]+)\s*\}\}/g;
  let match;
  while ((match = re.exec(template)) !== null) {
    matches.push(match[1]);
  }
  return [...new Set(matches)];
}

function validateTemplate(template, knownPlaceholders) {
  const used = extractPlaceholders(template);
  const known = new Set(
    knownPlaceholders.map(p => p.replace(/^\{\{|\}\}$/g, ''))
  );
  const unknown = used.filter(p => !known.has(p));
  return {
    valid:   unknown.length === 0,
    used,
    unknown,
    message: unknown.length
      ? `Unknown placeholders: ${unknown.map(p => '{{' + p + '}}').join(', ')}`
      : 'All placeholders are valid.',
  };
}

function validateTemplateEntry(entry) {
  const errors = [];
  if (!entry?.name?.trim()) errors.push('Name is required.');
  if (!entry?.content?.trim()) errors.push('Content is required.');
  if (entry?.name && entry.name.length > 100)
    errors.push('Name must be under 100 characters.');
  if (entry?.content && entry.content.length > 50_000)
    errors.push('Content must be under 50,000 characters.');
  return { valid: errors.length === 0, errors };
}
```

## 9. Conditional Sections

```js
function renderWithConditionals(template = '', tokens = {}) {
  let result = template.replace(
    /\{\{#if\s+(\w+)\}\}\n?([\s\S]*?)\{\{\/if\}\}\n?/g,
    (_match, key, body) => {
      const value = tokens?.[key];
      const hasValue = value != null
        && String(value).trim() !== ''
        && String(value).trim() !== '- [placeholder]';
      return hasValue ? body : '';
    }
  );
  return renderTemplate(result, tokens);
}

// Ternary inline: {{severity|default:Unknown}}
function renderWithDefaults(template = '', tokens = {}) {
  return String(template || '').replace(
    /\{\{\s*([a-zA-Z0-9_]+)\s*\|\s*default\s*:\s*([^}]*)\}\}/g,
    (_match, key, defaultValue) => {
      const value = tokens?.[key];
      return (value != null && String(value).trim())
        ? String(value)
        : defaultValue.trim();
    }
  );
}
```

## 10. Template Versioning

```js
const TEMPLATE_SCHEMA_VERSION = 2;

function migrateTemplateConfig(stored) {
  const version = stored?._schemaVersion ?? 1;

  if (version < 2) {
    // v1 → v2: rename "body" field to "content"
    if (Array.isArray(stored.customTemplates)) {
      stored.customTemplates = stored.customTemplates.map(t => ({
        ...t,
        content: t.content || t.body || '',
      }));
    }
  }

  stored._schemaVersion = TEMPLATE_SCHEMA_VERSION;
  return stored;
}

function detectBuiltinChanges(storedTemplates, currentDefaults) {
  const changes = [];
  for (const [key, defaultContent] of Object.entries(currentDefaults)) {
    const stored = storedTemplates[key];
    if (stored && stored !== defaultContent && stored !== '') {
      changes.push({ key, hasUserOverride: true, defaultChanged: true });
    }
  }
  return changes;
}
```

## 11. Template Preview and Options UI

```js
// Live preview with 300ms debounce
function setupTemplatePreview(textarea, previewEl, getTokens) {
  let debounceTimer;
  textarea.addEventListener('input', () => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      previewEl.textContent = renderTemplate(textarea.value, getTokens());
    }, 300);
  });
}

// Populate built-in template editors from storage
function populateTemplateEditors(draftTemplates) {
  TEMPLATE_DEFINITIONS.forEach(({ key }) => {
    const field = document.getElementById(`template-draft-${key}`);
    if (field) field.value = draftTemplates[key] || DEFAULT_TEMPLATES[key];
  });
}

// Dropdown selector for panel UI
function renderTemplateOptions(selectEl, options) {
  const catalog = buildTemplateCatalog(options);
  const currentValue = selectEl.value;
  selectEl.innerHTML = ['<option value="">Template</option>']
    .concat(catalog.map(entry =>
      `<option value="${escapeHtml(entry.key)}">${escapeHtml(entry.label)}</option>`
    ))
    .join('');
  if (currentValue) {
    const opt = selectEl.querySelector(`option[value="${CSS.escape(currentValue)}"]`);
    if (opt) selectEl.value = currentValue;
  }
}
```

## 12. Anti-Patterns

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| `eval()` or `new Function()` to render templates | CSP violation in Chrome extensions; XSS in web apps | Use regex-based `{{var}}` rendering |
| Storing rendered output instead of template + tokens | Cannot re-render when tokens change; wastes storage | Store template string, re-render on demand |
| Arbitrary JS expressions in placeholders | Sandbox escape; injection risk | Limit keys to `[a-zA-Z0-9_]` |
| No normalization on load | CRLF drift, whitespace, null entries accumulate | Always normalize on read and on write |
| Hardcoding template content in UI code | Cannot update templates without code changes | Separate definition from content; load from storage |
| Skipping ID deduplication in custom libraries | Duplicate IDs cause silent overwrites in the template map | Use suffix-incrementing deduplication |
| No fallback for missing tokens | Rendered output contains raw placeholder markers | Collapse missing tokens to empty string or bracket placeholder |

## 13. Testing Templates

```js
import { describe, expect, it } from 'vitest';

describe('template rendering', () => {
  it('renders tokens and collapses missing values', () => {
    const result = renderTemplate(
      'Case {{caseNumber}} -- {{severity}} -- {{missing}}',
      { caseNumber: '01234', severity: 'Sev2' }
    );
    expect(result).toBe('Case 01234 -- Sev2 -- ');
  });

  it('normalizes custom entries and deduplicates IDs', () => {
    const entries = normalizeLibraryEntries([
      { name: 'Handoff', content: 'Template A' },
      { name: 'Handoff', content: 'Template B' },
    ], 'tpl');
    expect(entries[0].id).toBe('tpl-handoff');
    expect(entries[1].id).toBe('tpl-handoff-2');
  });

  it('strips conditional blocks when token is empty', () => {
    const tpl = 'A\n{{#if blockers}}\nBlockers: {{blockersBullets}}\n{{/if}}\nB';
    const result = renderWithConditionals(tpl, { blockersBullets: '' });
    expect(result.trim()).toBe('A\nB');
  });
});
```

## 14. Quick-Reference Checklist

- [ ] Templates stored as plain strings, not code
- [ ] Mustache-style syntax with regex renderer (no eval)
- [ ] Placeholder registry declared in one file
- [ ] Token builder with explicit fallbacks for every placeholder
- [ ] Normalization on every read and write path
- [ ] ID deduplication for user-submitted custom templates
- [ ] Built-in vs custom templates separated in catalog and map
- [ ] Prompt overrides use `normalizePromptOverride(value, fallback)`
- [ ] Conditional sections use lightweight block stripping
- [ ] Schema version field for storage migrations
- [ ] Live preview with 300ms debounce
- [ ] Validation reports unknown placeholders before save
- [ ] Tests cover rendering, normalization, deduplication, and conditionals
