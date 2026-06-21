<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `sticky-notes-local-data` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: sticky-notes-local-data
version: 1.1.0
updated: 2026-05-29
description: >
  Chrome extension sticky notes with local data persistence — chrome.storage.local note CRUD
  (flat ID-keyed schema, single top-level key for atomic bulk ops), popout windows via
  chrome.windows.create with reuse detection, contenteditable rich editing with paste
  sanitization, debounced auto-save (800ms default) with beforeunload flush, note organization
  (tags/folders/full-text search over plainText), and cross-context sync via chrome.storage.onChanged
  with conflict-avoidance (pause-on-blur, resume-on-focus).
  TRIGGER: building or reviewing a Chrome extension note-taking surface, sticky note widget,
  contenteditable auto-save pattern, popout window from an extension, chrome.storage note CRUD,
  cross-context sync between popup/popout/background/content scripts, debounce save in an
  extension, or note organization with tags/folders/search.
  SKIP: server-side note storage or sync (use `backend-patterns`), IndexedDB for large corpora
  (use `indexeddb-patterns` or `dexie-indexeddb-local-first-reviewer`), general chrome.storage
  patterns unrelated to notes (use `chrome-storage-patterns`), extension architecture (use
  `chrome-mv3-advanced`).
category: developer
tags:
  - chrome-extension
  - sticky-notes
  - contenteditable
  - chrome-storage
  - popout-window
  - auto-save
  - debounce
  - cross-context-sync
triggers:
  - sticky note
  - note-taking extension
  - chrome.storage notes
  - popout window
  - contenteditable auto-save
  - note persistence
  - extension notepad
  - chrome.windows.create
  - debounce save
  - cross-context sync
related:
  - chrome-storage-patterns
  - chrome-mv3-advanced
  - extension-message-bridge
  - vanilla-js-ui-reviewer
  - indexeddb-patterns
---

# Sticky Notes with Local Data Persistence

Patterns for building a Chrome extension note-taking surface backed by `chrome.storage.local`, with popout windows, contenteditable editing, debounced auto-save, organizational features, and cross-context sync.

---

## 1. Note Data Model

Design a flat, ID-keyed schema that serializes cleanly to chrome.storage.

```js
// note-schema.js — canonical note shape
function createNote(overrides = {}) {
  return {
    id: crypto.randomUUID(),
    title: '',
    body: '',                   // HTML string from contenteditable
    plainText: '',              // stripped text for search indexing
    tags: [],                   // string[] for categorization
    folder: 'default',          // folder slug
    color: '#fff9c4',           // sticky background hex
    pinned: false,
    createdAt: Date.now(),
    updatedAt: Date.now(),
    ...overrides,
  };
}
```

**Storage key convention:** store all notes under a single top-level key (`stickyNotes`) to allow atomic bulk reads/writes and avoid key collisions.

```js
// { stickyNotes: { [noteId]: NoteObject, ... }, stickyMeta: { nextOrder: 5 } }
```

---

## 2. chrome.storage CRUD Layer

```js
const STORAGE_KEY = 'stickyNotes';

async function getAllNotes() {
  const result = await chrome.storage.local.get(STORAGE_KEY);
  return result[STORAGE_KEY] || {};
}

async function getNote(id) {
  const all = await getAllNotes();
  return all[id] ?? null;
}

async function saveNote(note) {
  const all = await getAllNotes();
  note.updatedAt = Date.now();
  all[note.id] = note;
  await chrome.storage.local.set({ [STORAGE_KEY]: all });
  return note;
}

async function deleteNote(id) {
  const all = await getAllNotes();
  delete all[id];
  await chrome.storage.local.set({ [STORAGE_KEY]: all });
}

async function saveAllNotes(notesObj) {
  await chrome.storage.local.set({ [STORAGE_KEY]: notesObj });
}
```

**Sharing across contexts:** Extension pages (popup, popout, options, dashboard) share the same origin — import as an ES module. Service workers use `importScripts`. Content scripts must route through `chrome.runtime.sendMessage` to the service worker.

### Error handling

```js
async function saveNoteSafe(note) {
  try {
    return await saveNote(note);
  } catch (err) {
    if (err.message?.includes('QUOTA_BYTES')) {
      console.error('[notes] Storage quota exceeded', err);
      // Surface a user-visible warning or offer to prune old notes
    }
    throw err;
  }
}
```

### Storage limits

| Limit | Value | Notes |
|---|---|---|
| Default `chrome.storage.local` quota | 10 MB | Was 5 MB before Chrome 114 |
| Per-item limit | None for local | Unlike sync (8 KB per item) |
| Unlimited storage | Add `"unlimitedStorage"` permission | For heavy collections with embedded images |

Always store a `plainText` field alongside `body` so full-text search never parses HTML at query time.

---

## 3. Popout Windows via chrome.windows.create

```js
async function openNotePopout(noteId) {
  const url = chrome.runtime.getURL(
    `popout.html?noteId=${encodeURIComponent(noteId)}`
  );

  // Reuse an existing popout for the same note if one is already open
  const allWindows = await chrome.windows.getAll({ populate: true });
  for (const win of allWindows) {
    for (const tab of win.tabs ?? []) {
      if (tab.url === url) {
        await chrome.windows.update(win.id, { focused: true });
        return win;
      }
    }
  }

  return chrome.windows.create({
    url,
    type: 'popup',
    width: 420,
    height: 520,
    top: 80,
    left: Math.round(screen.availWidth - 460),
  });
}
```

### Manifest entry

```jsonc
{
  "manifest_version": 3,
  "permissions": ["storage"],
  "action": { "default_popup": "popup.html" }
  // chrome.windows.create requires no extra permission
}
```

---

## 4. contenteditable Editing Surface

```js
function initEditor(el) {
  el.setAttribute('contenteditable', 'true');
  el.setAttribute('spellcheck', 'true');
  el.setAttribute('role', 'textbox');
  el.setAttribute('aria-multiline', 'true');

  // Paste as plain text — prevents style injection and XSS
  el.addEventListener('paste', (e) => {
    e.preventDefault();
    const text = e.clipboardData.getData('text/plain');
    document.execCommand('insertText', false, text);
  });

  // Basic formatting shortcuts
  el.addEventListener('keydown', (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'b') {
      e.preventDefault();
      document.execCommand('bold');
    }
    if ((e.metaKey || e.ctrlKey) && e.key === 'i') {
      e.preventDefault();
      document.execCommand('italic');
    }
  });
}
```

### Extracting plain text for search

```js
function stripHtml(html) {
  const tmp = document.createElement('div');
  tmp.textContent = html;   // assign as text, not markup — safe
  return tmp.textContent || '';
}
```

### Safe DOM content loading

When loading saved note HTML back into the editor, sanitize first to prevent stored-XSS.

```js
// Option A: DOMPurify (recommended when bundling is an option)
// insertAdjacentHTML is an XSS sink — ONLY use after DOMPurify sanitization:
// const clean = DOMPurify.sanitize(note.body, { ALLOWED_TAGS: ['b','i','br','ul','ol','li'] });
// editor.insertAdjacentHTML('afterbegin', clean);  // safe: DOMPurify output only

// Option B: manual allowlist via DOM parsing
function loadNoteSafe(editor, html) {
  const ALLOWED = new Set(['B', 'I', 'BR', 'UL', 'OL', 'LI', 'P']);
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, 'text/html');

  function walk(node, parent) {
    if (node.nodeType === Node.TEXT_NODE) {
      parent.appendChild(document.createTextNode(node.textContent));
    } else if (node.nodeType === Node.ELEMENT_NODE && ALLOWED.has(node.tagName)) {
      const el = document.createElement(node.tagName);
      for (const child of node.childNodes) walk(child, el);
      parent.appendChild(el);
    }
    // Disallowed elements are silently dropped
  }

  const frag = document.createDocumentFragment();
  for (const child of doc.body.childNodes) walk(child, frag);
  editor.textContent = '';
  editor.appendChild(frag);
}
```

---

## 5. Auto-Save with Debounce

```js
function debounce(fn, delayMs) {
  let timerId = null;
  function debounced(...args) {
    clearTimeout(timerId);
    timerId = setTimeout(() => fn(...args), delayMs);
  }
  debounced.cancel = () => clearTimeout(timerId);
  debounced.flush = (...args) => { clearTimeout(timerId); fn(...args); };
  return debounced;
}
```

### Wiring auto-save to the editor

```js
const editor = document.getElementById('note-editor');
const noteId = new URLSearchParams(location.search).get('noteId');
let localLastSave = 0;

// Load note on open
(async () => {
  const note = await getNote(noteId);
  if (note) loadNoteSafe(editor, note.body);
})();

const autoSave = debounce(async () => {
  const body = editor.innerHTML;
  const plainText = stripHtml(body);
  const note = await getNote(noteId);
  if (!note) return;
  note.body = body;
  note.plainText = plainText;
  await saveNote(note);
  localLastSave = note.updatedAt;
  showSaveIndicator();
}, 800);

editor.addEventListener('input', autoSave);
window.addEventListener('beforeunload', () => autoSave.flush());
```

### Choosing the debounce delay

| Delay | Tradeoff |
|---|---|
| 300 ms | Near-instant save; higher storage write frequency |
| 800 ms | Good balance for typical typing speed (recommended) |
| 1500 ms | Minimal writes; user may notice delay in indicator |

---

## 6. Note Organization — Tags, Folders, Search

```js
// Tags
function addTag(note, tag) {
  const normalized = tag.trim().toLowerCase();
  if (!note.tags.includes(normalized)) note.tags.push(normalized);
  return note;
}

function removeTag(note, tag) {
  note.tags = note.tags.filter(t => t !== tag);
  return note;
}

// Folders — derived from the set of distinct folder values
async function getFolders() {
  const all = await getAllNotes();
  const folders = new Set(Object.values(all).map(n => n.folder));
  folders.add('default');
  return [...folders].sort();
}

// Full-text search over plainText
async function searchNotes(query) {
  const q = query.toLowerCase();
  const all = await getAllNotes();
  return Object.values(all).filter(n =>
    n.plainText.toLowerCase().includes(q) ||
    n.title.toLowerCase().includes(q) ||
    n.tags.some(t => t.includes(q))
  );
}
```

---

## 7. Cross-Context Sync

Use `chrome.storage.onChanged` as the single sync primitive across popup, popout, dashboard, and content script overlay.

```js
// sync-listener.js — imported by every UI context
function onNotesChanged(callback) {
  chrome.storage.onChanged.addListener((changes, area) => {
    if (area !== 'local') return;
    if (!changes.stickyNotes) return;

    const oldNotes = changes.stickyNotes.oldValue || {};
    const newNotes = changes.stickyNotes.newValue || {};

    const changedIds = new Set();
    for (const id of Object.keys(newNotes)) {
      if (!oldNotes[id] || oldNotes[id].updatedAt !== newNotes[id].updatedAt)
        changedIds.add(id);
    }
    for (const id of Object.keys(oldNotes)) {
      if (!newNotes[id]) changedIds.add(id);   // deleted
    }

    if (changedIds.size > 0) callback(newNotes, changedIds);
  });
}
```

### Applying sync in an open popout editor

```js
onNotesChanged((allNotes, changedIds) => {
  if (!changedIds.has(noteId)) return;
  const updated = allNotes[noteId];
  if (!updated) { window.close(); return; }  // deleted in another context
  if (updated.updatedAt > localLastSave) {
    loadNoteSafe(editor, updated.body);
  }
});
```

### Conflict-avoidance rule

Because chrome.storage writes are last-writer-wins, only the *focused* editor should auto-save. Unfocused popouts pause their debounce timer and re-read on focus.

```js
window.addEventListener('blur', () => autoSave.cancel());
window.addEventListener('focus', () => {
  getNote(noteId).then(n => {
    if (n) loadNoteSafe(editor, n.body);
  });
});
```

---

## 8. Manifest Permissions Checklist

```jsonc
{
  "manifest_version": 3,
  "permissions": [
    "storage"
    // "unlimitedStorage"  — add only if note corpus exceeds 10 MB
  ]
  // chrome.windows.create and contenteditable require no extra permissions
}
```

---

## 9. Anti-Patterns

| Anti-pattern | Why it fails | Better approach |
|---|---|---|
| One storage key per note | Hundreds of keys thrash quota checks; atomic bulk ops impossible | Single `stickyNotes` object key |
| Saving on every `input` event | Write per keystroke degrades performance | Debounce at 800 ms |
| `localStorage` instead of `chrome.storage` | Not accessible from service worker; no `onChanged` sync | Always use `chrome.storage.local` |
| Unsanitized HTML paste or load | Stored/pasted HTML can inject scripts or break layout | Intercept paste as plain text; sanitize loads with DOMPurify or allowlist |
| No `beforeunload` flush | Closing popout before debounce fires loses last edit | Call `autoSave.flush()` on unload |
| Writing from blurred windows | Last-writer-wins overwrites newer edits from focused context | Pause auto-save on blur; re-read on focus |
| Storing HTML without plainText | Full-text search must parse HTML at query time | Always maintain `plainText` alongside `body` |

---

## 10. Testing Considerations

- **Unit-test the CRUD layer** by mocking `chrome.storage.local` with a plain object and verifying get/set/delete round-trips.
- **Unit-test the debounce** with fake timers (`vi.useFakeTimers()` in Vitest) — assert callback fires once after the delay and not on every input.
- **E2E-test the popout** with Playwright's `launchPersistentContext` to load the extension, open the popout via `chrome.windows.create`, type into the contenteditable, and verify storage was updated.
- **Cross-context sync test:** open two extension pages, edit in one, and assert the other receives the `onChanged` update within one second.
