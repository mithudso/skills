# Python CLI & TUI Application Development

`lang-python` hub reference. Building command-line tools (argparse / Click /
Typer) and full terminal UIs (Rich / Textual).

**Scope:** the app-building libraries and their selection. Packaging the result
as a runnable command (entry points / `console_scripts`) is in
`uv-python-toolchain`; distributing it as a standalone binary is in
`python-app-packaging`.

---

## 1. `argparse` — the standard library

Zero dependencies, on every Python since 3.2. Reach for it when you can't add a
dependency or the tool is small.

```python
import argparse

p = argparse.ArgumentParser(prog="mytool", description="…")
p.add_argument("path", type=str)                       # positional
p.add_argument("-n", "--count", type=int, default=1)   # option
p.add_argument("--verbose", action="store_true")
sub = p.add_subparsers(dest="cmd", required=True)       # subcommands
run = sub.add_parser("run"); run.add_argument("--fast", action="store_true")
args = p.parse_args()
```

The standard dispatch idiom for subcommands is `set_defaults(func=handler)` so
`args.func(args)` routes without an if-chain:

```python
run.set_defaults(func=lambda args: print("run", args.fast))
args = p.parse_args()
args.func(args)
```

Strengths: stdlib, predictable. Weaknesses: verbose, manual subcommand wiring,
no completion/colors, awkward type coercion. For anything with nested
subcommands, prefer Click/Typer.

---

## 2. Click — the mature decorator framework

The most widely adopted CLI library (dominant in the ecosystem). Decorator-based,
composable groups, rich parameter types, shell completion, testing harness.

```python
import click

@click.group()
@click.option("--verbose", is_flag=True)
@click.pass_context
def cli(ctx, verbose):
    ctx.obj = {"verbose": verbose}

@cli.command()
@click.argument("path", type=click.Path(exists=True))   # validates for you
@click.option("--count", type=int, default=1)
def run(path, count):
    click.echo(f"running {path} x{count}")
```

- **Parameter types** do validation: `Path(exists=True)`, `Choice([...])`,
  `IntRange`, `File("r")`, `DateTime`. Less hand-rolled checking than argparse.
- **Context** (`ctx.obj`, `pass_context`, `pass_obj`) threads shared state through
  the command tree.
- **Testing**: `click.testing.CliRunner().invoke(cli, [...])` returns
  `result.exit_code` / `result.output` — first-class CLI testing.
- `prompt=`, `confirmation_prompt=`, `password=` for interactive input;
  `click.progressbar`, `click.style`/`secho` for color.

---

## 3. Typer — type-hints as the API

Built **on top of Click** by the FastAPI author. You write a normal typed
function; Typer derives the parser. Least boilerplate, automatic `--help`, plays
with the type system.

```python
import typer
from typing import Annotated

app = typer.Typer()

@app.command()
def run(
    path: Annotated[str, typer.Argument(help="input path")],
    count: Annotated[int, typer.Option(help="repeats")] = 1,
    verbose: bool = False,
) -> None:
    typer.echo(f"running {path} x{count}")

if __name__ == "__main__":
    app()
```

- Types drive parsing: `int`/`float`/`bool`/`Enum`/`Path`/`datetime` → options
  and validation, no extra declarations. `Annotated[...]` is the current idiom
  (over the older default-value `typer.Option(...)` style).
- Inherits Click's ecosystem (completion, `CliRunner`-style testing via
  `typer.testing.CliRunner`). Startup overhead is comparable to Click (Typer
  wraps it), so the choice rarely turns on performance.
- Pairs naturally with Rich (Typer renders help and tracebacks through Rich).

---

## 4. Choosing

| Need | Pick |
| --- | --- |
| No dependencies / tiny script / locked-down env | **argparse** |
| Production CLI, subcommands, mature ecosystem, fine control | **Click** |
| Modern typed codebase, minimal boilerplate, auto docs | **Typer** |
| Full interactive terminal UI (not just commands) | **Textual** (§6) |

Rule of thumb: **Typer** for new typed projects, **Click** when you need its
depth or already use it, **argparse** when you can't add a dependency.

---

## 5. Rich — terminal rendering

The output layer (also the engine under Textual). Pretty printing, tables,
progress bars, syntax highlighting, markdown, and tracebacks.

```python
from rich.console import Console
from rich.table import Table
from rich.progress import track

console = Console()
console.print("[bold green]done[/]  :rocket:")          # console markup
t = Table(title="results"); t.add_column("name"); t.add_row("alpha")
console.print(t)
for item in track(items, description="processing…"):    # progress bar
    work(item)
```

- `rich.traceback.install()` — colorized, source-context tracebacks.
- `RichHandler` plugs Rich into the stdlib `logging` module (see
  `python-logging-observability`).
- `Console(record=True)` + `export_text/html/svg` to capture output.

---

## 6. Textual — full TUI framework

A Rapid-Application-Development framework for terminal UIs that can **also run in
a browser**. Built on Rich's rendering; a reactive, DOM-like architecture often
described as "React for the terminal." Renders via delta-updates of dirty regions
(high frame rates vs. `curses`).

```python
from textual.app import App, ComposeResult
from textual.widgets import Header, Footer, Button
from textual.reactive import reactive

class Counter(App):
    count: reactive[int] = reactive(0)               # reactive attribute

    def compose(self) -> ComposeResult:               # build the widget tree
        yield Header(); yield Button("inc", id="go"); yield Footer()

    def on_button_pressed(self, _) -> None:           # message handler
        self.count += 1

    def watch_count(self, value: int) -> None:        # auto-called on change
        self.title = f"count={value}"

if __name__ == "__main__":
    Counter().run()
```

Core concepts:
- **Widgets** — composable building blocks (`Button`, `Input`, `DataTable`,
  `Tree`, `TextArea`, `Switch`, …). Compose them in `compose()`.
- **Reactive attributes** (`reactive(...)`) — assignment is observed; a
  `watch_<name>` method fires on change (no manual re-render).
- **Messages & events** — widgets emit messages (`on_button_pressed`,
  `post_message`); the app dispatches them. Event-driven, not polling.
- **TCSS** — Textual CSS in `.tcss` files / `CSS` strings styles the layout and
  theme separately from logic.
- **`@work` / workers** — run async or threaded background work without blocking
  the UI; integrates with the asyncio model (see `python-concurrency`).
- Run in a terminal (`app.run()`), headless for tests (`run_test()`), or serve to
  a browser (`textual serve`).

---

## 7. Best practices & pitfalls

- **Make the package runnable.** Declare a `[project.scripts]` entry point
  (`mytool = "mypkg.cli:app"`) so `uv run` / pipx / pip install exposes the
  command — don't ship a bare `python script.py`. See `uv-python-toolchain`.
- **Exit codes matter.** Return/raise to set non-zero exit on failure
  (`raise typer.Exit(1)`, `sys.exit(2)`, Click's `ctx.exit`). Scripts and CI
  depend on them.
- **Separate I/O from logic.** Keep parsing thin; put behavior in plain functions
  you can unit-test without invoking the CLI.
- **Test through the runner**, not by shelling out — `CliRunner` is fast and
  asserts on exit code + output.
- **Don't reinvent Rich.** Use it for tables/progress/color rather than manual
  ANSI codes; respect `NO_COLOR` and non-TTY output (Rich auto-detects).
- **Async CLIs:** `asyncio.run(...)` inside the command, or `anyio.run`; Textual
  is already async — use its workers, never block the loop.
- **Pitfall:** heavy imports at module top slow every `--help`; defer expensive
  imports into the command body for snappy startup.

---

## Sources

- [Typer docs — Alternatives, Inspiration & Comparisons](https://typer.tiangolo.com/alternatives/)
- [Comparing argparse, Click, and Typer — CodeCut](https://codecut.ai/comparing-python-command-line-interface-tools-argparse-click-and-typer/)
- [Building CLI Tools with Python: Click, Typer, argparse (2025)](https://dasroot.net/posts/2025/12/building-cli-tools-python-click-typer-argparse/)
- [Textual — Home](https://textual.textualize.io/) and [Tutorial](https://textual.textualize.io/tutorial/) / [Widgets guide](https://textual.textualize.io/guide/widgets/)
- [Python Textual: Build Beautiful UIs in the Terminal — Real Python](https://realpython.com/python-textual/)
- [argparse — Python stdlib docs](https://docs.python.org/3/library/argparse.html)
