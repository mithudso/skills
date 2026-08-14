<!-- hub-reference-banner -->
> **Reference file — part of the `lang-js-ts` hub.** Formerly the standalone `typescript-expert` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: typescript-expert
description: >
  TypeScript general expert: project setup, tsconfig strategy, module resolution, everyday types, generics, narrowing, utility types, declaration files, and coding standards. Covers strict mode, interface vs type alias, function signatures, null/undefined handling, control-flow narrowing, and TSConfig policy choices.
  TRIGGER: user writes, reviews, or debugs TypeScript code; asks about tsconfig options, module boundaries, declaration files, any/unknown tradeoffs, generic APIs, or utility types.
  SKIP: advanced type-system internals (conditional types, mapped types, branded types) — use typescript-advanced-types; runtime validation with Zod — use zod-schema-validation; framework-specific typing (React props, Express handlers) — use the framework skill.
version: "1.1"
category: developer
updated: "2026-05-29"
tags:
  - typescript
  - tsconfig
  - generics
  - modules
  - narrowing
  - declaration-files
  - utility-types
related_skills:
  - typescript-advanced-types
  - zod-schema-validation
  - javascript-nodejs
---

# TypeScript Expert

Practical TypeScript reference for everyday types, generics, narrowing, utility types, modules, declaration files, and TSConfig strategy. Treat the official [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html) and [TSConfig reference](https://www.typescriptlang.org/tsconfig/) as the source of truth.

**Version note:** based on TypeScript rolling docs as accessed 2026-05-10. The "Creating Types from Types" page was last updated 2026-05-04.

## When to use this skill

- Setting up a TypeScript project or tuning `tsconfig.json`
- Reviewing or writing type annotations, interfaces, or generics
- Understanding module boundaries (ES Modules vs CommonJS, script vs module files)
- Writing or consuming `.d.ts` declaration files
- Debugging `any`/`unknown` tradeoffs or type widening
- Choosing between `interface` and `type` alias
- Handling `null`/`undefined` optionality safely

## When NOT to use this skill

- Advanced type-system patterns (conditional types, mapped types, branded types) — use `typescript-advanced-types`
- Runtime validation with Zod or io-ts — use `zod-schema-validation`
- React, Express, or other framework-specific typing — use the framework skill
- Vitest test file types — use `testing-and-vitest-expert`

---

## Quick rules

1. Prefer `string`, `number`, `boolean` over boxed `String`, `Number`, `Boolean`.
2. Avoid `any` — use `unknown` and narrow it. `any` discards type guarantees for that path.
3. Narrow unions with runtime checks (`typeof`, `in`, `instanceof`, equality) before accessing members.
4. Use generics when input and output types are related — `any` throws that relationship away.
5. Derive types from types with `keyof`, `typeof`, indexed access, conditional types, and mapped types instead of copying shapes by hand.
6. A file with a top-level `import` or `export` is a module; without one it is a global script.
7. Use `visudo`-style discipline for `.d.ts` files — follow official templates, not ad hoc ambient declarations.
8. Treat `tsconfig.json` as codebase policy, not just transpilation plumbing.

---

## Core type-system model

### Primitive and array types

- Primitives: `string`, `number`, `boolean`, `bigint`, `symbol`, `null`, `undefined`.
- Arrays: `T[]` and `Array<T>` are equivalent.
- `any` — escape hatch; disables checking. Use `unknown` instead when type is genuinely unknown, then narrow.

### Object, interface, and type alias

- Object types can be anonymous, or named via `interface` or `type` alias.
- `interface` supports declaration merging; `type` alias does not. Both express structural types.
- Properties can be optional (`prop?: T`) or readonly (`readonly prop: T`).
- Treat optional properties as genuinely absent-capable — never assume they are present.

### Function types

- Function type expressions: `(a: string) => void`. Parameter name is required syntactically.
- Call signatures in object types model callable values that also carry properties or overloads.
- `(string) => void` means a parameter literally named `string` with implicit `any` — not a string-typed parameter.

---

## Control-flow narrowing

TypeScript tracks reachable control flow and refines union types within branches.

| Check | Example |
|-------|---------|
| `typeof` | `typeof x === "string"` |
| Truthiness | `if (x)` narrows out `null`/`undefined`/`0`/`""` |
| Equality | `x === null`, `x === "literal"` |
| `in` operator | `"prop" in obj` narrows to types that have `prop` |
| `instanceof` | `x instanceof Date` |
| Type predicate | `function isString(v: unknown): v is string` |

Write runtime checks that reflect real invariants; the type system becomes more precise automatically.

---

## Generics

- Use generics when input and output types are related: `function identity<T>(x: T): T`.
- Generic constraints: `<T extends { length: number }>` — `T` must have `length`.
- Good generic API: capture the type the caller already has, return it with minimal widening.
- Do not add a generic just because you can — a simple concrete type is clearer when the relationship does not matter.

---

## Creating types from types

| Operator | Purpose | Example |
|----------|---------|---------|
| `keyof T` | Union of property names | `keyof { a: 1; b: 2 }` → `"a" \| "b"` |
| `typeof val` (type position) | Derive type from value | `type Cfg = typeof defaultConfig` |
| `T["key"]` | Indexed access | `User["email"]` → `string` |
| Conditional `T extends U ? X : Y` | Type-level branching | `IsArray<T>` |
| Mapped `{ [K in keyof T]: ... }` | Transform properties | `Readonly<T>`, `Partial<T>` |
| Template literal `` `${A}-${B}` `` | String type patterns | `EventName` |

---

## Utility types reference

| Utility | Effect |
|---------|--------|
| `Partial<T>` | All properties optional |
| `Required<T>` | All properties required |
| `Readonly<T>` | All properties readonly |
| `Pick<T, K>` | Keep only keys K |
| `Omit<T, K>` | Drop keys K |
| `Record<K, V>` | Object type with keys K and values V |
| `Exclude<T, U>` | Remove U from union T |
| `Extract<T, U>` | Keep members of T assignable to U |
| `NonNullable<T>` | Remove `null` and `undefined` |
| `Parameters<T>` | Tuple of function parameter types |
| `ReturnType<T>` | Return type of a function |
| `Awaited<T>` | Recursively unwrap Promise (TS 4.5+) |

Use built-in utility types as the first stop for common transformations. Custom mapped/conditional types are for cases built-ins cannot cover.

---

## Modules and declaration files

### Module vs script

- File with top-level `import`/`export` = **module** (scoped namespace).
- File without either = **script** (global namespace, shared scope).
- Accidentally leaving a file as a script leaks names globally. Add `export {}` to force module mode when needed.

### ES Modules vs CommonJS

- TypeScript centers guidance on ES Modules; CommonJS is documented for ecosystem compatibility.
- `moduleResolution: "bundler"` or `"node16"` affects how imports resolve. Match the runtime.
- For Node.js: use `"module": "nodenext"` and `"moduleResolution": "nodenext"` together.

### Declaration files

- `.d.ts` files describe JS/library API surfaces without runtime code.
- Most common use: typing npm packages that lack built-in types.
- Follow official [declaration file templates](https://www.typescriptlang.org/docs/handbook/declaration-files/templates.html) instead of ad hoc ambient declarations.
- Do not sprinkle `declare module "x"` across the codebase — use a single `.d.ts` per package.

---

## TSConfig strategy

TSConfig is codebase policy. Key decisions:

| Option | Recommendation |
|--------|---------------|
| `strict` | Enable. Turns on `strictNullChecks`, `noImplicitAny`, and more. |
| `noUncheckedIndexedAccess` | Enable for safer array/object indexing. |
| `exactOptionalPropertyTypes` | Enable to distinguish `undefined` assignment from property absence. |
| `moduleResolution` | Match the runtime: `"bundler"` for Vite/webpack; `"nodenext"` for Node.js. |
| `skipLibCheck` | Enable in development to avoid type-checking `node_modules` .d.ts files. |
| `allowUnreachableCode: false` | Surface provably unreachable code as an error. |
| `allowUnusedLabels: false` | Catch accidental label-like mistakes. |

**Strict mode checklist:**
- `strictNullChecks` — `null` and `undefined` are not assignable to every type.
- `noImplicitAny` — implicit `any` is an error.
- `strictFunctionTypes` — function parameter types checked contravariantly.
- `useUnknownInCatchVariables` — catch variables are `unknown`, not `any`.

---

## Coding standards

### Public API typing

- Model parameter and return relationships with generics rather than `any`.
- Use interfaces and function type expressions to make public contracts readable.
- Annotate function signatures; let the compiler infer local variable types.

### Null/undefined handling

- Treat `?` properties as genuinely absent-capable.
- Use narrowing and explicit checks instead of non-null assertion (`!`) as a first reflex.
- Prefer `??` (nullish coalescing) over `||` when `0` and `""` are valid values.

### Type reuse

- Derive types from existing types and values. Duplication creates sync bugs.
- Utility types first; custom mapped/conditional types only when built-ins fall short.

### Declaration file hygiene

- Follow official templates for library typing.
- Do not use declaration merging to patch third-party types in application code — use module augmentation in a dedicated `.d.ts` file.

---

## Version-sensitive notes

- `Awaited<T>` introduced in TS 4.5.
- `satisfies` operator introduced in TS 4.9 (see `typescript-advanced-types` for usage).
- `NoInfer<T>` introduced in TS 5.4.
- `const` type parameters introduced in TS 5.0.
- Docs are rolling — TSConfig options and utility types often note their release version.

---

## References

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Everyday Types](https://www.typescriptlang.org/docs/handbook/2/everyday-types.html)
- [Narrowing](https://www.typescriptlang.org/docs/handbook/2/narrowing.html)
- [Generics](https://www.typescriptlang.org/docs/handbook/2/generics.html)
- [Creating Types from Types](https://www.typescriptlang.org/docs/handbook/2/types-from-types.html)
- [Utility Types](https://www.typescriptlang.org/docs/handbook/utility-types.html)
- [Modules](https://www.typescriptlang.org/docs/handbook/2/modules.html)
- [Declaration Files](https://www.typescriptlang.org/docs/handbook/declaration-files/introduction.html)
- [TSConfig Reference](https://www.typescriptlang.org/tsconfig/)
