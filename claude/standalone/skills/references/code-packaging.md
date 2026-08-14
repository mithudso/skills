<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Formerly the standalone `code-packaging` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: code-packaging
description: >
  Library packaging and distribution expert for JavaScript/TypeScript, Python, and Rust.
  Covers ESM vs CJS module format selection and dual-package configuration; monorepo setup
  with Turborepo, Nx, pnpm, or Cargo workspaces; build tool selection (tsup, tsdown, Rollup,
  esbuild, Vite library mode); semantic versioning and automated releases; tree-shaking and
  sideEffects configuration; peer dependency and lockfile strategy; plugin system and
  extension-point design; supply-chain security with provenance attestation; and API
  documentation generation (TypeDoc, JSDoc, Sphinx, rustdoc).
  TRIGGER: user is designing, building, or publishing a reusable library or shared module
  (npm, PyPI, crates.io, GitHub Packages); asking about ESM/CJS/dual-package, package.json
  exports field, tree-shaking, monorepo tooling, semver, provenance, peer dependencies,
  plugin architecture, or API documentation.
  SKIP: application-level bundling (Vite/webpack app builds); framework-specific app setup
  (React, Next.js, Django); runtime debugging of already-published packages (use
  programming-languages); general CI/CD pipeline design beyond the publish step.
version: "1.1.1"
updated: "2026-05-31"
category: developer
tags:
  - npm
  - esm
  - commonjs
  - monorepo
  - library
  - versioning
  - build-tools
  - rust
  - python
  - supply-chain
keywords:
  - npm
  - esm
  - commonjs
  - dual-package
  - conditional-exports
  - tree-shaking
  - sideEffects
  - semver
  - semantic-versioning
  - monorepo
  - turborepo
  - nx
  - pnpm
  - cargo
  - rust
  - rollup
  - esbuild
  - tsup
  - tsdown
  - unbuild
  - pkgroll
  - vite
  - library
  - publish
  - crates
  - pypi
  - pip
  - python
  - peer-dependencies
  - lockfile
  - barrel-exports
  - typedoc
  - jsdoc
  - rustdoc
  - provenance
  - supply-chain
  - plugin-system
  - versioning
whenToUse:
  - Designing the public API surface and export shape of a new library
  - Publishing or releasing a package to npm, PyPI, or crates.io
  - Choosing between ESM-only, CJS-only, or dual ESM+CJS output
  - Configuring package.json exports field and conditional exports
  - Selecting a build tool for a TypeScript or JavaScript library (tsup, Rollup, esbuild, Vite)
  - Setting up a monorepo with Turborepo, Nx, pnpm workspaces, or Cargo workspaces
  - Applying semantic versioning rules and automating changelog generation
  - Configuring tree-shaking via sideEffects field and deep imports
  - Diagnosing or resolving peer dependency conflicts and lockfile drift
  - Designing a plugin system, hook registry, or extension-point architecture
  - Generating API documentation with TypeDoc, JSDoc, Sphinx, or rustdoc
  - Hardening publish pipelines with provenance attestation and trusted publishing
whenNotToUse:
  - Application-level bundling (Vite or webpack app builds, not library mode)
  - Framework-specific app setup (React, Next.js, Django)
  - Runtime debugging of already-published packages — use programming-languages
  - General CI/CD pipeline design beyond the publish step
  - Scanning/auditing a Python project's dependencies for CVEs, generating an SBOM, bandit SAST, or sigstore/PEP 740 attestation verification (the consumer/security-audit angle, not the library-publish angle) — use programming-languages references/python-supply-chain-security.md
related_skills:
  - programming-languages
  - devops-infra
  - software-engineering-patterns
---

# Code Packaging and Library Design

Full lifecycle of packaging reusable code: from module system selection and API surface design,
through build tooling and monorepo orchestration, to publishing, versioning, and documentation.

**Full reference:** `references/code-packaging-context.md`

## Navigation

| Topic | Section in reference |
|---|---|
| ESM vs CJS, dual packages, conditional exports | Module Systems |
| npm / pnpm / yarn / pip / cargo workspaces | Package Management |
| Tree-shaking, barrel exports, sideEffects | Library Design Patterns |
| semver rules, breaking changes, prerelease | Semantic Versioning |
| Turborepo, Nx, Lerna, pnpm workspaces | Monorepo Tooling |
| tsup, Rollup, esbuild, Vite library mode | Build Systems |
| npm publish, PyPI, crates.io, provenance | Code Distribution |
| Lockfiles, peer deps, dependency injection | Dependency Management |
| Plugin systems, IoC, extension points | Modular Architecture |
| TypeDoc, JSDoc, Sphinx, rustdoc | Documentation Generation |

## Language Routing

| Language | Relevant sections |
|---|---|
| JavaScript / TypeScript | Module Systems, Build Systems, Package Management (npm/pnpm/yarn) |
| Python | Package Management (pip/PyPI), Semantic Versioning, Documentation (Sphinx) |
| Rust | Package Management (Cargo workspaces), Code Distribution (crates.io), Documentation (rustdoc) |

## Example Inputs and Expected Outputs

**"I'm publishing a TypeScript utility library. Should I ship ESM, CJS, or both?"**
Decision table from Module Systems + recommended `package.json` exports config + build tool recommendation (tsup for zero-config dual output).

**"My monorepo has 8 packages and CI takes 12 minutes. How do I speed it up?"**
Turborepo/Nx comparison table, caching strategy, `turbo.json` pipeline config skeleton.

**"Users keep getting duplicate React instances from my library."**
Peer dependency configuration pattern, `peerDependencies` vs `dependencies` distinction, lockfile deduplication steps.

**"How do I publish a crate to crates.io from GitHub Actions without storing tokens?"**
Trusted Publishing / OIDC setup steps for crates.io, `crates-io-auth-action` config snippet.

## Known Edge Cases and Failure Modes

- **Dual-package hazard:** When both ESM and CJS versions load in the same runtime (mixed `require`/`import`), module-level state is duplicated. Guard with `instanceof` checks or a singleton registry pattern.
- **Barrel files break tree-shaking:** `index.ts` re-exporting everything prevents bundlers from eliminating unused code unless `"sideEffects": false` is set and the barrel itself has no side effects. Prefer deep imports for large libraries.
- **`tsc` declaration bottleneck:** On large TypeScript projects, `--dts` generation dominates build time. Use `isolatedDeclarations` (TS 5.5+) or `dts-bundle-generator` to parallelize.
- **Cargo workspace + binary crate:** Binary crates in a workspace cannot be published to crates.io as workspace members — publish independently with `cargo publish -p <name>`.
- **pnpm phantom dependencies:** pnpm's strict isolation means packages can only import what they declare. Code that imports a hoisted transitive dep will break in pnpm — declare all direct dependencies explicitly.
