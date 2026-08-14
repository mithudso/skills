#!/usr/bin/env python3
"""Estimate the cost of a prompt across Claude models and effort levels.

Token counting: uses the Anthropic count_tokens API when credentials are
available (exact), otherwise falls back to an offline heuristic
(~1 token per 3.5 characters). The source used is reported in the output.

Cost model:
  input cost  = input_tokens x input_price
  output cost = (expected_output_tokens x effort_multiplier) x output_price
  total       = input + output

Effort multipliers approximate how much extra thinking+output each effort
level spends relative to the base output estimate. They are heuristics,
not API-reported numbers — tune with --multipliers if you have measured data.

Usage:
  estimate_cost.py --file prompt.txt
  estimate_cost.py --text "Summarize this repo"
  estimate_cost.py --file prompt.txt --output-tokens 4000
  estimate_cost.py --file prompt.txt --models claude-opus-4-8,claude-sonnet-5
  estimate_cost.py --file prompt.txt --json
  estimate_cost.py --file prompt.txt --calls 1000        # scale to N calls
  estimate_cost.py --file prompt.txt --cached            # assume cache-read pricing on input
  estimate_cost.py --file prompt.txt --quick             # fast, free, no API — one-line rough number
  estimate_cost.py --file prompt.txt --confirm-exec       # after the estimate, ask y/N to actually run it
  estimate_cost.py --file prompt.txt --exec               # skip the ask, run it immediately
  estimate_cost.py --file prompt.txt --recommend          # cheap heuristic: what model/effort should this run on
  estimate_cost.py --file prompt.txt --exec --target cheapest   # --exec against pure cost-min instead
"""

import argparse
import json
import sys
from datetime import date

# ---------------------------------------------------------------------------
# Pricing table — per 1M tokens, USD. Cached 2026-06. Verify against
# https://platform.claude.com/docs/en/pricing.md if precision matters.
# ---------------------------------------------------------------------------
MODELS = {
    "claude-fable-5": {
        "display": "Fable 5",
        "input": 10.00,
        "output": 50.00,
        "efforts": ["low", "medium", "high", "xhigh", "max"],
        "notes": "thinking always on",
    },
    "claude-opus-4-8": {
        "display": "Opus 4.8",
        "input": 5.00,
        "output": 25.00,
        "efforts": ["low", "medium", "high", "xhigh", "max"],
    },
    "claude-opus-4-7": {
        "display": "Opus 4.7",
        "input": 5.00,
        "output": 25.00,
        "efforts": ["low", "medium", "high", "xhigh", "max"],
    },
    "claude-opus-4-6": {
        "display": "Opus 4.6",
        "input": 5.00,
        "output": 25.00,
        "efforts": ["low", "medium", "high", "max"],
    },
    "claude-sonnet-5": {
        "display": "Sonnet 5",
        "input": 3.00,
        "output": 15.00,
        "efforts": ["low", "medium", "high", "xhigh", "max"],
        "intro": {"input": 2.00, "output": 10.00, "until": "2026-08-31"},
    },
    "claude-sonnet-4-6": {
        "display": "Sonnet 4.6",
        "input": 3.00,
        "output": 15.00,
        "efforts": ["low", "medium", "high", "max"],
    },
    "claude-haiku-4-5": {
        "display": "Haiku 4.5",
        "input": 1.00,
        "output": 5.00,
        "efforts": ["-"],  # no effort parameter
    },
}

# Heuristic: total output tokens (thinking + response) as a multiple of the
# base expected-output estimate, per effort level.
DEFAULT_MULTIPLIERS = {
    "-": 1.0,       # no effort control (Haiku)
    "low": 1.0,
    "medium": 1.5,
    "high": 2.5,
    "xhigh": 4.0,
    "max": 6.0,
}

CACHE_READ_FACTOR = 0.1  # cache-read input tokens bill at ~0.1x base input

CHARS_PER_TOKEN = 3.5  # offline heuristic

# ---------------------------------------------------------------------------
# Model/effort recommendation — a cheap, free, zero-network heuristic (pure
# keyword + length matching on the prompt text) for "what should I actually
# run this on", separate from the cost matrix's pure cost-minimization view.
# Tier vocabulary mirrors skill-optimizer's Step 4.6 model-recommendation
# table for consistency across the ecosystem, recalibrated for task prompts
# rather than skill files.
# ---------------------------------------------------------------------------
EFFORT_LADDER = ["-", "low", "medium", "high", "xhigh", "max"]

TIER_DEFAULTS = {
    "mechanical": ("claude-haiku-4-5", "-"),
    "routine": ("claude-sonnet-5", "medium"),
    "analytical": ("claude-opus-4-8", "high"),
    "long_horizon": ("claude-opus-4-8", "xhigh"),
    "frontier": ("claude-fable-5", "xhigh"),
}

TIER_KEYWORDS = {
    "frontier": ["frontier reasoning", "novel research", "prove that", "cutting-edge reasoning"],
    "long_horizon": ["multi-step", "multi step", "pipeline", "orchestrat", "end-to-end",
                      "end to end", "agent loop", "autonomous", "multi-agent"],
    "analytical": ["review", "audit", "diagnose", "debug", "analyze", "analyse", "evaluate",
                   "compare", "optimize", "optimise", "refactor", "design", "architecture",
                   "root cause"],
    "mechanical": ["classify", "convert", "rename", "extract", "lookup", "look up",
                   "one word", "single word", "yes or no", "true or false", "format as"],
}

TIER_RATIONALE = {
    "mechanical": "short/lookup-style task, no judgment required",
    "routine": "default: templated or moderate-judgment task, no strong signal either way",
    "analytical": "review/design/diagnosis-style task needs real judgment",
    "long_horizon": "multi-step or agentic task benefits from deeper reasoning budget",
    "frontier": "explicitly frontier-reasoning task",
}


def classify_tier(text: str) -> str:
    lower = text.lower()
    for tier in ("frontier", "long_horizon", "analytical", "mechanical"):
        if any(kw in lower for kw in TIER_KEYWORDS[tier]):
            return tier
    return "routine"


def closest_effort(model_id: str, desired: str) -> str:
    """Nearest effort the model actually supports, at or below desired on the ladder."""
    available = MODELS[model_id]["efforts"]
    if desired in available:
        return desired
    want_idx = EFFORT_LADDER.index(desired) if desired in EFFORT_LADDER else 0
    ranked = sorted(available, key=lambda e: abs(EFFORT_LADDER.index(e) - want_idx)
                     if e in EFFORT_LADDER else 99)
    return ranked[0]


def recommend(text: str, model_ids: list) -> dict:
    """Cheap, free, no-network recommendation of which model/effort to actually run on."""
    tier = classify_tier(text)
    model_id, effort = TIER_DEFAULTS[tier]
    substituted = False
    if model_id not in model_ids:
        # Requested scope doesn't include the tier's default model — fall back to
        # the cheapest-output model actually in scope and note the substitution.
        model_id = min(model_ids, key=lambda m: MODELS[m]["output"])
        substituted = True
    effort = closest_effort(model_id, effort)
    rationale = TIER_RATIONALE[tier]
    if substituted:
        rationale += f"; {TIER_DEFAULTS[classify_tier(text)][0]} not in --models scope, using {model_id}"
    return {"tier": tier, "model": model_id, "effort": effort, "rationale": rationale}


def count_tokens_api(text: str, model: str = "claude-opus-4-8"):
    """Exact count via the Anthropic API. Returns None on any failure."""
    try:
        import anthropic
        client = anthropic.Anthropic()
        resp = client.messages.count_tokens(
            model=model,
            messages=[{"role": "user", "content": text}],
        )
        return resp.input_tokens
    except Exception:
        return None


def count_tokens_heuristic(text: str) -> int:
    return max(1, round(len(text) / CHARS_PER_TOKEN))


def execute_prompt(model_id: str, effort: str, text: str) -> dict:
    """Actually send the prompt to the model. Returns {ok, text|error, usage}."""
    try:
        import anthropic
    except ImportError:
        return {"ok": False, "error": "anthropic package not installed (pip install anthropic)"}
    try:
        client = anthropic.Anthropic()
        kwargs = {
            "model": model_id,
            "max_tokens": 4096,
            "messages": [{"role": "user", "content": text}],
        }
        if effort not in ("-", None):
            kwargs["effort"] = effort
        resp = client.messages.create(**kwargs)
        out_text = "".join(
            block.text for block in resp.content if getattr(block, "type", None) == "text"
        )
        usage = getattr(resp, "usage", None)
        return {
            "ok": True,
            "text": out_text,
            "input_tokens": getattr(usage, "input_tokens", None),
            "output_tokens": getattr(usage, "output_tokens", None),
        }
    except Exception as e:
        return {"ok": False, "error": str(e)}


def actual_cost(model_id: str, effort: str, in_tok: int, out_tok: int) -> float:
    info = MODELS[model_id]
    today = date.today().isoformat()
    use_intro = "intro" in info and today <= info["intro"]["until"]
    in_price = effective_input_price(info, use_intro)
    out_price = effective_output_price(info, use_intro)
    return (in_tok or 0) * in_price / 1_000_000 + (out_tok or 0) * out_price / 1_000_000


def effective_input_price(model_info: dict, use_intro: bool) -> float:
    if use_intro and "intro" in model_info:
        return model_info["intro"]["input"]
    return model_info["input"]


def effective_output_price(model_info: dict, use_intro: bool) -> float:
    if use_intro and "intro" in model_info:
        return model_info["intro"]["output"]
    return model_info["output"]


def build_estimates(input_tokens, output_tokens, multipliers, model_ids,
                    cached=False, calls=1):
    today = date.today().isoformat()
    rows = []
    for mid in model_ids:
        info = MODELS[mid]
        use_intro = "intro" in info and today <= info["intro"]["until"]
        in_price = effective_input_price(info, use_intro)
        out_price = effective_output_price(info, use_intro)
        in_factor = CACHE_READ_FACTOR if cached else 1.0
        input_cost = input_tokens * in_price * in_factor / 1_000_000
        for effort in info["efforts"]:
            mult = multipliers.get(effort, 1.0)
            est_out = round(output_tokens * mult)
            output_cost = est_out * out_price / 1_000_000
            total = (input_cost + output_cost) * calls
            rows.append({
                "model": mid,
                "display": info["display"] + (" (intro)" if use_intro else ""),
                "effort": effort,
                "input_tokens": input_tokens,
                "est_output_tokens": est_out,
                "input_cost": round(input_cost * calls, 6),
                "output_cost": round(output_cost * calls, 6),
                "total_cost": round(total, 6),
            })
    return rows


def fmt_usd(x):
    if x >= 0.01:
        return f"${x:,.4f}".rstrip("0").rstrip(".")
    return f"${x:.6f}"


def print_table(rows, input_tokens, token_source, output_tokens, calls, cached):
    print(f"Input tokens: {input_tokens:,} ({token_source})")
    print(f"Base output estimate: {output_tokens:,} tokens "
          f"(effort multiplies this; override with --output-tokens)")
    if cached:
        print("Input priced at cache-read rate (~0.1x). First call pays full/write price.")
    if calls > 1:
        print(f"Costs shown for {calls:,} calls.")
    print()
    header = f"{'Model':<20} {'Effort':<8} {'Est out tok':>12} {'Input $':>12} {'Output $':>12} {'Total $':>12}"
    print(header)
    print("-" * len(header))
    for r in rows:
        print(f"{r['display']:<20} {r['effort']:<8} {r['est_output_tokens']:>12,} "
              f"{fmt_usd(r['input_cost']):>12} {fmt_usd(r['output_cost']):>12} "
              f"{fmt_usd(r['total_cost']):>12}")
    print()
    cheapest = min(rows, key=lambda r: r["total_cost"])
    priciest = max(rows, key=lambda r: r["total_cost"])
    print(f"Cheapest: {cheapest['display']} @ {cheapest['effort']} = {fmt_usd(cheapest['total_cost'])}")
    print(f"Most expensive: {priciest['display']} @ {priciest['effort']} = {fmt_usd(priciest['total_cost'])}")


def print_quick(rows, input_tokens, token_source, calls):
    """One-line rough estimate. Assumes rows already narrowed to one row."""
    r = rows[0]
    scale = f" x{calls:,} calls" if calls > 1 else ""
    eff = "" if r["effort"] in ("-", None) else f" @ {r['effort']}"
    print(f"~{fmt_usd(r['total_cost'])}  ({r['display']}{eff}{scale})")
    print(f"  {input_tokens:,} in + ~{r['est_output_tokens']:,} out tok  "
          f"[{token_source}; rough — use full mode for a real quote]")


def print_recommendation(rec: dict):
    display = MODELS[rec["model"]]["display"]
    eff = "" if rec["effort"] in ("-", None) else f" @ {rec['effort']}"
    print(f"\nRecommended: {display}{eff}  [tier: {rec['tier']}]")
    print(f"  why: {rec['rationale']}")
    print(f"  to switch this session: /model {rec['model']}"
          + (f"  /effort {rec['effort']}" if rec["effort"] not in ("-", None) else ""))


def print_exec_result(result: dict, model_id: str, effort: str, est_total: float):
    if not result["ok"]:
        print(f"\nExecution failed: {result['error']}")
        return
    print(f"\n--- Response ({MODELS[model_id]['display']} @ {effort}) ---")
    print(result["text"])
    real = actual_cost(model_id, effort, result["input_tokens"], result["output_tokens"])
    print(f"\n--- Actual usage ---")
    print(f"{result['input_tokens']:,} in + {result['output_tokens']:,} out tok "
          f"= {fmt_usd(real)} actual  (estimate was {fmt_usd(est_total)})")


def main():
    p = argparse.ArgumentParser(description="Estimate prompt cost across Claude models and efforts")
    src = p.add_mutually_exclusive_group(required=True)
    src.add_argument("--file", help="Path to prompt file")
    src.add_argument("--text", help="Prompt text inline")
    src.add_argument("--tokens", type=int, help="Skip counting; use this input token count")
    p.add_argument("--output-tokens", type=int, default=2000,
                   help="Expected base output tokens at low effort (default 2000)")
    p.add_argument("--models", help="Comma-separated model IDs (default: all)")
    p.add_argument("--multipliers", help='JSON dict overriding effort multipliers, e.g. \'{"high": 3.0}\'')
    p.add_argument("--calls", type=int, default=1, help="Number of calls to scale costs by")
    p.add_argument("--cached", action="store_true",
                   help="Price input at cache-read rate (~0.1x)")
    p.add_argument("--offline", action="store_true",
                   help="Skip the count_tokens API; use the character heuristic")
    p.add_argument("--quick", action="store_true",
                   help="Fast, free rough estimate: forces --offline (no API), picks the "
                        "cheapest model at low effort unless --models is given, and prints "
                        "a single-line number instead of the full table")
    p.add_argument("--json", action="store_true", help="Emit JSON instead of a table")
    p.add_argument("--confirm-exec", action="store_true",
                   help="After showing the estimate, ask y/N whether to actually run the "
                        "prompt against the cheapest-cost model/effort")
    p.add_argument("--exec", action="store_true", dest="exec_now",
                   help="Skip the estimate-only default and immediately run the prompt "
                        "against the target model/effort (no confirmation) — see --target")
    p.add_argument("--recommend", action="store_true",
                   help="Cheap, free, no-network heuristic: gauge the optimal model/effort "
                        "for this prompt and print it (with the /model, /effort commands "
                        "to switch this session to it)")
    p.add_argument("--target", choices=["recommended", "cheapest"], default="recommended",
                   help="Which row --exec/--confirm-exec runs against (default: recommended, "
                        "the --recommend heuristic's pick; 'cheapest' restores pure cost-min). "
                        "Ignored under --quick, which always targets its own cheapest row.")
    p.add_argument("--price-only", action="store_true",
                   help="Print nothing but the dollar figure for --target (e.g. '$0.0031') and "
                        "exit — for scripting, hooks, and status lines. Implies --offline. "
                        "Suppresses --recommend/--exec/--confirm-exec output (price only, no side effects).")
    args = p.parse_args()

    # --quick and --price-only are speed/cost presets: never touch the network.
    if args.quick or args.price_only:
        args.offline = True

    if args.output_tokens <= 0:
        sys.exit("--output-tokens must be a positive integer")
    if args.calls <= 0:
        sys.exit("--calls must be a positive integer")
    if args.tokens is not None and args.tokens <= 0:
        sys.exit("--tokens must be a positive integer")

    # Resolve input tokens
    token_source = None
    prompt_text = None  # only set when we have the actual prompt (needed to --exec)
    if args.tokens is not None:
        input_tokens = args.tokens
        token_source = "user-supplied"
    else:
        if args.text is not None:
            text = args.text
        else:
            try:
                with open(args.file, encoding="utf-8", errors="replace") as f:
                    text = f.read()
            except OSError as e:
                sys.exit(f"Cannot read {args.file}: {e}")
        prompt_text = text
        input_tokens = None
        if not args.offline:
            input_tokens = count_tokens_api(text)
            if input_tokens is not None:
                token_source = "exact, count_tokens API"
        if input_tokens is None:
            input_tokens = count_tokens_heuristic(text)
            token_source = f"heuristic, ~{CHARS_PER_TOKEN} chars/token"

    multipliers = dict(DEFAULT_MULTIPLIERS)
    if args.multipliers:
        try:
            overrides = json.loads(args.multipliers)
        except json.JSONDecodeError as e:
            sys.exit(f"--multipliers is not valid JSON: {e}")
        if not isinstance(overrides, dict) or not all(
            isinstance(v, (int, float)) and v >= 0 for v in overrides.values()
        ):
            sys.exit("--multipliers must be a JSON object of non-negative numbers")
        unknown_keys = set(overrides) - set(DEFAULT_MULTIPLIERS)
        if unknown_keys:
            sys.exit(f"Unknown effort level(s) in --multipliers: "
                     f"{', '.join(sorted(unknown_keys))}. "
                     f"Known: {', '.join(DEFAULT_MULTIPLIERS)}")
        multipliers.update(overrides)

    if args.models:
        model_ids = [m.strip() for m in args.models.split(",")]
        unknown = [m for m in model_ids if m not in MODELS]
        if unknown:
            sys.exit(f"Unknown model(s): {', '.join(unknown)}. Known: {', '.join(MODELS)}")
    elif args.quick:
        # Cheapest model gives the "how little could this cost" number.
        model_ids = [min(MODELS, key=lambda m: MODELS[m]["output"])]
    else:
        model_ids = list(MODELS)

    rows = build_estimates(input_tokens, args.output_tokens, multipliers,
                           model_ids, cached=args.cached, calls=args.calls)

    if args.quick:
        # Narrow to the single cheapest row for a one-line answer.
        rows = [min(rows, key=lambda r: r["total_cost"])]

    cheapest = min(rows, key=lambda r: r["total_cost"])

    # Recommendation is only computable with real prompt text, and doesn't
    # apply under --quick (that mode's whole point is "cheapest, no thinking").
    recommendation = None
    if prompt_text is not None and not args.quick:
        recommendation = recommend(prompt_text, model_ids)

    # Resolve which row --exec/--confirm-exec will act on.
    target_row = cheapest
    if args.target == "recommended" and recommendation is not None:
        target_row = next(
            (r for r in rows if r["model"] == recommendation["model"]
             and r["effort"] == recommendation["effort"]),
            cheapest,  # recommended model/effort not in the current row set
        )

    if args.price_only:
        # Bare number, nothing else — safe for a status line or hook to
        # capture directly with no parsing (no exec, no recommendation printed).
        print(fmt_usd(target_row["total_cost"]))
        return

    if args.json:
        payload = {
            "input_tokens": input_tokens,
            "token_source": token_source,
            "base_output_tokens": args.output_tokens,
            "calls": args.calls,
            "cached_input": args.cached,
            "quick": args.quick,
            "estimates": rows,
        }
        if recommendation is not None:
            payload["recommendation"] = recommendation
        if args.exec_now:
            if prompt_text is None:
                payload["execution"] = {"ok": False, "error": "no prompt text available (--tokens was used)"}
            else:
                result = execute_prompt(target_row["model"], target_row["effort"], prompt_text)
                payload["execution"] = result
        print(json.dumps(payload, indent=2))
        return

    if args.quick:
        print_quick(rows, input_tokens, token_source, args.calls)
    else:
        print_table(rows, input_tokens, token_source, args.output_tokens,
                    args.calls, args.cached)

    if args.recommend:
        if recommendation is None:
            print("\nCannot recommend: no prompt text available, or --quick was used.")
        else:
            print_recommendation(recommendation)

    if not (args.exec_now or args.confirm_exec):
        return

    if prompt_text is None:
        print("\nCannot execute: no prompt text available (--tokens was used without --file/--text).")
        return

    label = "recommended" if (args.target == "recommended" and recommendation is not None) else "cheapest"
    run_it = args.exec_now
    if args.confirm_exec and not run_it:
        prompt = (f"\nExecute now with {target_row['display']} @ {target_row['effort']} "
                  f"({label}, est. {fmt_usd(target_row['total_cost'])})? [y/N]: ")
        if not sys.stdin.isatty():
            print(prompt + "(non-interactive stdin, skipping)")
            return
        run_it = input(prompt).strip().lower() in ("y", "yes")

    if run_it:
        if args.target == "recommended" and recommendation is not None and not args.recommend:
            print_recommendation(recommendation)
        result = execute_prompt(target_row["model"], target_row["effort"], prompt_text)
        print_exec_result(result, target_row["model"], target_row["effort"], target_row["total_cost"])


if __name__ == "__main__":
    main()
