# GitHub Pending Review — API Reference

The GitHub review API has two distinct workflows depending on what you're posting. Mixing them up causes 404 or 422 errors.

## New inline diff comments → pending review

Use a pending review when posting *new* inline comments on diff lines (not replies to existing threads). The entire review is invisible until the user publishes it.

```bash
# Build payload in a file — body often contains em-dashes, backticks, angle
# brackets that break shell quoting if passed inline
cat > /tmp/review_payload.json <<'EOF'
{
  "body": "<PR-level summary>",
  "comments": [
    {
      "path": "src/mongo/db/query/planner.cpp",
      "line": 42,
      "side": "RIGHT",
      "body": "**[tag]** comment text"
    }
  ]
}
EOF

gh api repos/10gen/mongo/pulls/<number>/reviews \
  --method POST \
  --input /tmp/review_payload.json
```

**Quirks that cause failures if ignored:**

| Quirk | What breaks |
|---|---|
| Do NOT set `"event": "PENDING"` | GitHub rejects it as an invalid enum value — omit the field entirely |
| `"body"` MUST be non-empty | The PATCH endpoint cannot add a body later; returns "Could not edit a review with a missing body" |
| Always set `"side": "RIGHT"` on inline comments | Targets the new-file side of the diff |
| Use `"line"` not `"position"` | `position` is deprecated diff-relative offset; `line` is the actual line number in the new file |
| Write payload to a file; use `--input` | Heredoc piped to `--input -` breaks on special characters in comment bodies |

**Verify the review is pending before reporting:**

```bash
gh api repos/10gen/mongo/pulls/<number>/reviews/<review_id> \
  --jq '{state: .state, body_chars: (.body | length), submitted_at: .submitted_at}'
# Expected: "state": "PENDING"

gh api repos/10gen/mongo/pulls/<number>/reviews/<review_id>/comments \
  --jq 'length'
# Expected: matches the number of inline comments you submitted
```

If `state` is anything other than `PENDING`, the review was published immediately — stop and tell the user before proceeding.

---

## Replies to existing threads → direct POST

`in_reply_to` does **not** work through the pending review endpoint. Use the pull request comments endpoint directly instead. These replies publish immediately.

```bash
gh api --method POST repos/10gen/mongo/pulls/<number>/comments \
  --field body="<reply text>" \
  --field in_reply_to=<root_comment_id>
```

When `in_reply_to` is set, `path`, `line`, and `commit_id` are all inherited from the parent comment — do not include them.

---

## Mixing both in one session

If you need to stage new inline comments (pending) AND reply to existing threads in the same session:

1. Post the pending review with new inline comments first.
2. **Delete the pending review** before posting `in_reply_to` replies. Having an open pending review causes `POST .../comments` with `in_reply_to` to fail with 422 "user_id can only have one pending review per pull request."
3. Post replies via the direct `POST .../comments` endpoint.

```bash
# Delete pending review after verifying it looks correct
gh api --method DELETE repos/10gen/mongo/pulls/<number>/reviews/<review_id>

# Then post thread replies
gh api --method POST repos/10gen/mongo/pulls/<number>/comments \
  --field body="<reply>" \
  --field in_reply_to=<root_comment_id>
```

The pending review (new inline comments) remains deleted at this point — the user will need to publish those through the GitHub UI. If you need those comments to survive, submit the pending review as `COMMENT` event before deleting it, then post the thread replies.

---

## Publishing a pending review

```bash
gh api --method POST repos/10gen/mongo/pulls/<number>/reviews/<review_id>/events \
  --field event="COMMENT"
# Use APPROVE or REQUEST_CHANGES instead of COMMENT when warranted
# Never submit on behalf of the user without explicit confirmation
```
