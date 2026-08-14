# Live Hub Toolkit — DOM Extraction Helpers

## Core utility functions

```js
function normalizeWhitespace(value) {
  return String(value || '').replace(/\s+/g, ' ').trim();
}

function isElementVisible(element) {
  if (!(element instanceof Element)) return false;
  const style = window.getComputedStyle(element);
  if (style.display === 'none' || style.visibility === 'hidden' || Number(style.opacity) === 0) return false;
  const rect = element.getBoundingClientRect();
  return rect.width > 0 && rect.height > 0;
}

function readElementValue(element) {
  if (!(element instanceof Element)) return '';
  const candidates = [
    element.getAttribute('datetime'),
    element.getAttribute('data-datetime'),
    element.getAttribute('title'),
    element.getAttribute('aria-label'),
    element.getAttribute('content'),
    'value' in element ? element.value : '',
    element.textContent || ''
  ];
  for (const candidate of candidates) {
    const text = normalizeWhitespace(candidate);
    if (text) return text;
  }
  return '';
}

function escapeRegex(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
```

## Universal labeled value extractor

Support portals render key-value pairs in various HTML structures (dt/dd, th/td, label/span). Extractor tries multiple strategies:

```js
function extractLabeledValue(labels = []) {
  const visibleNodes = Array.from(
    document.querySelectorAll('dt,dd,th,td,li,span,div,p,strong,b,h1,h2,h3,h4,a')
  ).filter(isElementVisible);

  for (const label of labels) {
    for (const node of visibleNodes) {
      const text = normalizeWhitespace(node.textContent);

      // Strategy 1: inline "Label: Value" pattern
      const inlineMatch = text.match(
        new RegExp(`^${escapeRegex(label)}\\s*[:\\-]\\s*(.+)$`, 'i')
      );
      if (inlineMatch?.[1]) return normalizeWhitespace(inlineMatch[1]);

      // Strategy 2: sibling element contains the value
      if (text.toLowerCase() === label.toLowerCase()) {
        const sibling = node.nextElementSibling;
        const value = readElementValue(sibling);
        if (value) return value;
      }
    }
  }
  return '';
}
```

## Account case table extraction

Walk visible case links and container rows, build structured records:

```js
function extractAccountCases() {
  const seen = new Map();
  const caseLinks = Array.from(
    document.querySelectorAll('a[href*="/case/"], a[href*="/cases/"]')
  ).filter(isElementVisible);

  const containers = new Set(
    caseLinks.map(link =>
      link.closest('tr, article, li, [role="row"], [role="listitem"]')
    ).filter(Boolean)
  );

  for (const container of containers) {
    const candidate = buildAccountCaseCandidate(container, { accountId, accountName });
    if (!candidate?.caseNumber) continue;

    const existing = seen.get(candidate.caseNumber);
    const score = Object.values(candidate).filter(Boolean).length;
    const existingScore = existing ? Object.values(existing).filter(Boolean).length : -1;

    if (!existing || score > existingScore) {
      seen.set(candidate.caseNumber, candidate);
    }
  }
  return Array.from(seen.values());
}
```

## Account matching helper

Match discovered record to tracked account by ID (preferred) or name (fallback):

```js
function matchesTrackedAccount(record = {}, accountId = '', accountName = '') {
  const recordAccountId   = String(record.accountId   || record.account_id   || '').trim();
  const recordAccountName = String(record.accountName || record.account_name || '').trim().toLowerCase();
  const targetId   = String(accountId   || '').trim();
  const targetName = String(accountName || '').trim().toLowerCase();
  if (targetId   && recordAccountId)   return targetId   === recordAccountId;
  if (targetName && recordAccountName) return targetName === recordAccountName;
  return false;
}
```