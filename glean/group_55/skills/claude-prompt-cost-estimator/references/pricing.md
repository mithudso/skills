# Claude model pricing snapshot

Cached 2026-06 from the claude-api skill / platform.claude.com. Per 1M tokens, USD, first-party API. The script's `MODELS` dict in `scripts/estimate_cost.py` mirrors this table — update both together.

| Model | ID | Input $/1M | Output $/1M | Effort levels | Notes |
| --- | --- | --- | --- | --- | --- |
| Fable 5 | `claude-fable-5` | 10.00 | 50.00 | low–max | thinking always on |
| Opus 4.8 | `claude-opus-4-8` | 5.00 | 25.00 | low–max | |
| Opus 4.7 | `claude-opus-4-7` | 5.00 | 25.00 | low–max | |
| Opus 4.6 | `claude-opus-4-6` | 5.00 | 25.00 | low, medium, high, max (no xhigh) | |
| Sonnet 5 | `claude-sonnet-5` | 3.00 | 15.00 | low–max | intro $2.00/$10.00 through 2026-08-31 |
| Sonnet 4.6 | `claude-sonnet-4-6` | 3.00 | 15.00 | low, medium, high, max (no xhigh) | |
| Haiku 4.5 | `claude-haiku-4-5` | 1.00 | 5.00 | none (effort param unsupported) | 200K context |

## Pricing modifiers

| Modifier | Factor | Applies to |
| --- | --- | --- |
| Cache read | ~0.1x | input tokens served from cache |
| Cache write (5m TTL) | ~1.25x | input tokens written to cache |
| Cache write (1h TTL) | ~2x | input tokens written to cache |
| Batch API | 0.5x | all token usage, non-latency-sensitive |

## Effort multiplier heuristics

Output-spend multipliers relative to the base output estimate (script defaults; not API-reported):

| Effort | Multiplier |
| --- | --- |
| low | 1.0 |
| medium | 1.5 |
| high | 2.5 |
| xhigh | 4.0 |
| max | 6.0 |

## Refreshing

Verify against https://platform.claude.com/docs/en/pricing.md (or the claude-api skill's model table) when the user needs current numbers. Update this file, the `MODELS` dict, and any intro-pricing end dates together.
