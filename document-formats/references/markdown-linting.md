# Markdown Linting & Quality Gates (markdownlint · remark-lint · Vale · link-check)

Enforcing **consistent, valid, well-styled markdown** in editors and CI. This reference owns markdown
linting/quality tooling. For markdown syntax itself load `references/markdown-authoring.md`; for the
remark/unified plumbing that remark-lint runs on load `references/markdown-processing.md`; for the
docs-as-code CI context load `references/markdown-docs-as-code.md`.

---

## 1. Two different jobs: structural lint vs prose lint

- **Structural/style lint** (markdownlint, remark-lint) — markdown *correctness & consistency*: heading
  increments, list marker style, blank-line rules, trailing spaces, fenced-code languages, line length.
- **Prose lint** (Vale, write-good, alex, proselint) — *writing quality*: style-guide terms, passive
  voice, inclusive language, banned words. Complementary, not a substitute. Run both in docs CI.

## 2. markdownlint (the de-facto structural linter)

Rules are `MD001`–`MD0xx` (e.g. `MD013` line-length, `MD025` single-H1, `MD040` fenced-code language,
`MD033` no-inline-HTML). Two main runners:

- **markdownlint-cli2** (current preferred CLI) and **markdownlint-cli**; also the `markdownlint` library
  (Node) and the VS Code extension (David Anson). Ports exist (`pymarkdownlnt` in Python, `mdl` in Ruby).

```jsonc
// .markdownlint-cli2.jsonc  (or .markdownlint.json / .yaml)
{
  "config": {
    "default": true,
    "MD013": { "line_length": 100, "code_blocks": false, "tables": false }, // relax line length
    "MD033": { "allowed_elements": ["details", "summary", "br"] },           // allow some inline HTML
    "MD041": false                                                            // first-line-H1 off (frontmatter)
  },
  "globs": ["**/*.md"],
  "ignores": ["node_modules", "CHANGELOG.md"]
}
```

```bash
npx markdownlint-cli2 "**/*.md"            # lint
npx markdownlint-cli2 --fix "**/*.md"      # auto-fix the fixable rules
```

Inline control: `<!-- markdownlint-disable MD013 -->` … `<!-- markdownlint-enable MD013 -->`, or
`<!-- markdownlint-disable-next-line MD033 -->`.

## 3. remark-lint (lint inside the unified pipeline)

If you already use remark (`references/markdown-processing.md`), lint as plugins on the mdast:

```js
import { remark } from 'remark'
import remarkPresetLintRecommended from 'remark-preset-lint-recommended'
import remarkPresetLintConsistent from 'remark-preset-lint-consistent'

await remark()
  .use(remarkPresetLintRecommended)
  .use(remarkPresetLintConsistent)        // enforce whatever style the doc uses first
  .process(file)                          // messages attach to file.messages
```

Or via `remark-cli` + `.remarkrc`. Presets: `recommended` (catch real problems), `consistent` (match the
file's own first choice), or `remark-preset-lint-markdown-style-guide`. Hundreds of individual
`remark-lint-*` rules compose. **markdownlint vs remark-lint:** markdownlint is simpler/standalone and the
common choice; remark-lint is right when markdown is already flowing through a remark/unified build.

## 4. Link checking & prose linting

```bash
lychee --no-progress "**/*.md"             # fast Rust link checker (internal + external, anchors)
npx markdown-link-check ./README.md        # per-file link check (Node)
vale ./docs                                # prose style-guide linter (Vale + a style like Microsoft/Google)
```

Vale uses a `.vale.ini` + style packages (Microsoft, Google, write-good) and is the standard for
enforcing a *writing* style guide; pair it with markdownlint for full coverage.

## 5. Wire it into the workflow

- **pre-commit** (the framework): hooks for `markdownlint-cli2`, `lychee`, `vale` so issues are caught before commit.
- **CI** (`references/markdown-docs-as-code.md` §4): run lint + link-check + prose-lint as gates; fail the build on errors; many teams run `--fix` locally but check (not fix) in CI.
- **Editor**: the markdownlint VS Code extension surfaces rules live; format-on-save fixes the easy ones.
- **Monorepos/docs**: scope globs to the docs dir and ignore generated files (CHANGELOG, vendored content).

## Sources
- [markdownlint (rules)](https://github.com/DavidAnson/markdownlint) · [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2) · [rule reference](https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md)
- [remark-lint](https://github.com/remarkjs/remark-lint) · [remark-preset-lint-recommended](https://github.com/remarkjs/remark-lint/tree/main/packages/remark-preset-lint-recommended)
- [Vale](https://vale.sh/) · [lychee link checker](https://github.com/lycheeverse/lychee) · [pre-commit](https://pre-commit.com/)
