---
name: coding-standards
description: >
  Baseline cross-project coding conventions for naming, readability, immutability, error handling, and code-quality review. Covers JavaScript, TypeScript, and general principles (KISS, DRY, YAGNI). Use when enforcing consistency, reviewing code for quality, onboarding contributors, or applying naming and formatting rules.
  TRIGGER: "naming convention", "code quality", "immutability", "error handling pattern", "variable naming", "function naming", "code smells", "magic numbers", "early returns", "DRY", "YAGNI", "KISS".
  SKIP: React component composition or hooks (use frontend-design); backend service architecture or API endpoint design (use backend-patterns or api-design-patterns); design pattern selection (use coding-patterns); language-specific framework details when a narrower skill exists.
origin: ECC
version: 1.1.0
updated: "2026-05-29"
category: developer
tags: [coding-standards, naming, readability, immutability, error-handling, javascript, typescript]
related_skills:
  - software-engineering-patterns
  - frontend-ui
  - lang-js-ts
whenToUse:
  - enforcing naming conventions across a project
  - reviewing code for readability and maintainability
  - setting up linting, formatting, or type-checking rules
  - applying DRY, KISS, or YAGNI to reduce code complexity
  - identifying code smells (long functions, deep nesting, magic numbers)
  - onboarding contributors to team conventions
whenNotToUse:
  - React component composition, hooks, or rendering patterns (use frontend-design)
  - backend architecture, API design, or database layering (use backend-patterns)
  - design pattern selection for structural decisions (use coding-patterns)
  - domain-specific framework guidance when a narrower skill already exists
---

# Coding Standards & Best Practices

Baseline coding conventions applicable across JavaScript/TypeScript projects.

**Scope:** This skill is the shared floor — naming, readability, immutability, error handling, and code-smell detection. It is not the framework playbook.

- Use `frontend-design` for React, state, forms, rendering, and UI architecture.
- Use `backend-patterns` or `api-design-patterns` for repository/service layers, endpoint design, and server-specific concerns.
- Use `coding-patterns` for structural design decisions (factory, observer, state machine, etc.).

---

## Code Quality Principles

| Principle | Rule |
|-----------|------|
| **Readability first** | Code is read more than written. Clear names and consistent formatting beat clever code. |
| **KISS** | Use the simplest solution that works. Optimize only when you have evidence of a bottleneck. |
| **DRY** | Extract logic used in 3+ places into a shared function. Avoid copy-paste programming. |
| **YAGNI** | Don't build features before they're needed. Add complexity only when a second concrete case requires it. |

---

## Naming

### Variables

```typescript
// GOOD: Descriptive names
const marketSearchQuery = 'election'
const isUserAuthenticated = true
const totalRevenue = 1000

// BAD: Unclear names
const q = 'election'
const flag = true
const x = 1000
```

### Functions

```typescript
// GOOD: Verb-noun pattern
async function fetchMarketData(marketId: string) { }
function calculateSimilarity(a: number[], b: number[]) { }
function isValidEmail(email: string): boolean { }

// BAD: Noun-only or unclear
async function market(id: string) { }
function similarity(a, b) { }
```

---

## Immutability

Never mutate objects or arrays directly. Always return new copies.

```typescript
// GOOD: spread operator
const updatedUser = { ...user, name: 'New Name' }
const updatedArray = [...items, newItem]

// BAD: direct mutation
user.name = 'New Name'
items.push(newItem)
```

---

## Error Handling

```typescript
// GOOD: Explicit HTTP status check + typed error re-throw
async function fetchData(url: string) {
  try {
    const response = await fetch(url)
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Fetch failed:', error)
    throw new Error('Failed to fetch data')
  }
}

// BAD: No error handling
async function fetchData(url) {
  const response = await fetch(url)
  return response.json()
}
```

---

## Async / Await

Run independent async calls in parallel with `Promise.all`.

```typescript
// GOOD: Parallel execution
const [users, markets, stats] = await Promise.all([
  fetchUsers(),
  fetchMarkets(),
  fetchStats()
])

// BAD: Sequential when there is no dependency between calls
const users = await fetchUsers()
const markets = await fetchMarkets()
const stats = await fetchStats()
```

---

## Type Safety (TypeScript)

Avoid `any`. Define interfaces for domain objects.

```typescript
// GOOD
interface Market {
  id: string
  name: string
  status: 'active' | 'resolved' | 'closed'
  created_at: Date
}

function getMarket(id: string): Promise<Market> { /* ... */ }

// BAD
function getMarket(id: any): Promise<any> { /* ... */ }
```

---

## Comments

Comment *why*, not *what*. The code states what; the comment explains the non-obvious reasoning.

```typescript
// GOOD: Explains non-obvious reasoning
// Use exponential backoff to avoid overwhelming the API during outages
const delay = Math.min(1000 * Math.pow(2, retryCount), 30000)

// BAD: Restates the obvious
// Increment counter by 1
count++
```

JSDoc for public APIs:

```typescript
/**
 * Searches markets using semantic similarity.
 * @param query - Natural language search query
 * @param limit - Maximum results (default: 10)
 * @returns Markets sorted by similarity score
 * @throws {Error} If the search API fails
 */
export async function searchMarkets(query: string, limit = 10): Promise<Market[]> { }
```

---

## Code Smell Detection

### Long functions
Split functions over ~50 lines into focused sub-functions. Each function should do one thing.

```typescript
// BAD: 100-line monolith
function processMarketData() { /* 100 lines */ }

// GOOD: Composed sub-functions
function processMarketData() {
  const validated = validateData()
  const transformed = transformData(validated)
  return saveData(transformed)
}
```

### Deep nesting
Use early returns to flatten guard conditions.

```typescript
// BAD: 5 levels deep
if (user) {
  if (user.isAdmin) {
    if (market) {
      if (market.isActive) {
        if (hasPermission) { /* do something */ }
      }
    }
  }
}

// GOOD: Early returns
if (!user) return
if (!user.isAdmin) return
if (!market) return
if (!market.isActive) return
if (!hasPermission) return
// do something
```

### Magic numbers
Name every non-obvious constant.

```typescript
// BAD
if (retryCount > 3) { }
setTimeout(callback, 500)

// GOOD
const MAX_RETRIES = 3
const DEBOUNCE_DELAY_MS = 500
if (retryCount > MAX_RETRIES) { }
setTimeout(callback, DEBOUNCE_DELAY_MS)
```

---

## Testing Standards

### AAA pattern

```typescript
test('calculates similarity correctly', () => {
  // Arrange
  const vector1 = [1, 0, 0]
  const vector2 = [0, 1, 0]

  // Act
  const similarity = calculateCosineSimilarity(vector1, vector2)

  // Assert
  expect(similarity).toBe(0)
})
```

### Test naming
Name tests as behavioral specifications.

```typescript
// GOOD
test('returns empty array when no markets match query', () => { })
test('throws when OpenAI API key is missing', () => { })

// BAD
test('works', () => { })
test('test search', () => { })
```

---

## File Organization (Next.js / Node reference)

```
src/
├── app/          # Next.js App Router pages + API routes
├── components/   # React components (ui/, forms/, layouts/)
├── hooks/        # Custom React hooks (use* prefix)
├── lib/          # Utilities, API clients, constants
├── types/        # TypeScript type definitions
└── styles/       # Global styles
```

File naming: `Button.tsx` (PascalCase for components), `useAuth.ts` (camelCase with `use` prefix for hooks), `formatDate.ts` (camelCase for utilities).

---

> Code quality is not optional. Clear, maintainable code is what enables fast, confident refactoring.
