<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Created 2026-07-01 as the **build-automation
> & task-runner** reference — `make`/Makefile semantics and the modern task-runner landscape. For CI/CD
> pipeline design use `references/cicd-pipelines.md`; for git-triggered release automation use
> `references/git-workflows.md`. Sibling topics are reference files under the devops hubs — **not**
> standalone skills. Ignore any "use the X skill" pointer that names a bare sibling; load that topic's
> `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: build-automation-make
title: Build Automation (Make & Task Runners)
description: >
  GNU Make semantics and the task-runner landscape that surrounds it. TRIGGER: how make decides what to
  rebuild (target/prerequisite timestamps, the dependency DAG); .PHONY targets and why they matter;
  pattern rules, automatic variables ($@, $<, $^), and variable flavors (= vs := vs ?= vs +=); tab-vs-space
  recipe rule; parallel builds (make -j) and .NOTPARALLEL; common Makefile anti-patterns; when make is the
  wrong tool; modern task runners and build systems (just, Task/Taskfile, npm/pnpm scripts, invoke, Bazel,
  CMake, Ninja) and how to choose.
  SKIP: language-specific build tooling internals (webpack/vite, cargo, gradle — use the language skill);
  CI/CD pipeline orchestration and runners (use cicd-pipelines.md); container image build layering (use
  docker-containers.md); shell-script authoring inside recipes (use shell-scripting.md).
triggers:
  - make
  - makefile
  - build automation
  - phony target
  - pattern rule
  - automatic variables
  - make -j parallel
  - task runner
  - justfile
  - taskfile
  - ninja
  - cmake
  - when to use make
version: "1.0"
updated: "2026-07-01"
category: developer
tags:
  - make
  - build
  - automation
  - task-runner
  - just
  - ninja
  - devops
whenToUse:
  - Explaining why make rebuilds (or fails to rebuild) a target
  - Fixing a Makefile that "works sometimes" (missing .PHONY, tab/space, missing prereqs)
  - Choosing between make, just, Taskfile, or language-native scripts for a repo
  - Setting up incremental builds keyed on file timestamps and a dependency DAG
---

# Build Automation — Make and its modern successors

`make` is two things at once: an **incremental build engine** (only rebuild what changed) and a **task
runner** (named commands). Most confusion comes from using it as the latter while it behaves as the former.

## How make decides what to rebuild

- A rule is `target: prerequisites` + a **recipe** (the indented commands). make rebuilds `target` when it
  is **missing** or **older** than any prerequisite (mtime comparison). This forms a **dependency DAG** that
  make walks bottom-up.
- **`.PHONY`** marks targets that are *not* files (`clean`, `test`, `build`). Without it, if a file named
  `test` exists, make thinks the target is up to date and **silently does nothing** — the single most common
  Makefile bug. Declare `.PHONY: clean test build all`.
- Recipes **must** be indented with a **real tab**, not spaces (`*** missing separator` error otherwise).

## Variables and automatic variables

| Syntax | Meaning |
| --- | --- |
| `VAR = x` | **recursive** — expanded each use (can reference later-defined vars; can cause surprises/loops) |
| `VAR := x` | **simple** — expanded once, at definition (usually what you want) |
| `VAR ?= x` | set only if not already set |
| `VAR += x` | append |

Automatic variables inside a recipe: `$@` = target, `$<` = first prerequisite, `$^` = all prerequisites
(deduped), `$*` = the stem in a pattern rule. **Pattern rules** generalize: `%.o: %.c` with recipe
`cc -c $< -o $@`. Escape a literal `$` as `$$`.

## Parallelism and correctness

- `make -j N` runs independent DAG branches in parallel — but **only correct if prerequisites are fully and
  accurately declared**. Missing prereqs that "worked" serially cause nondeterministic failures under `-j`.
  This is the real reason to keep dependencies honest.
- `.NOTPARALLEL` (whole file) or ordering via prereqs forces serialization where needed.

## When make is the wrong tool — and the alternatives

Reach past make when you want a pure task runner without file-timestamp semantics, or a big multi-language
build graph:

| Tool | Sweet spot |
| --- | --- |
| **just** (`justfile`) | Pure command runner — no phony/tab traps, per-recipe shebangs, args. Best "make but just tasks". |
| **Task** (`Taskfile.yml`) | YAML task runner with declarative deps and up-to-date checks; cross-platform. |
| **npm/pnpm scripts** | JS/TS repos already using package.json. |
| **CMake + Ninja** | Native C/C++: CMake *generates* Ninja/Make files; Ninja is a fast low-level build executor. |
| **Bazel** | Very large hermetic, multi-language, cached/remote builds (high setup cost). |

Rule of thumb: **file-based incremental build → make/Ninja; named tasks for humans → just/Task/package
scripts; huge hermetic monorepo build → Bazel.**
