#!/usr/bin/env bash
# Usage: resolve-thread.sh <thread-id>
# Thread IDs look like PRRT_kwDO...
set -euo pipefail

thread_id=$1

gh api graphql -f query='
mutation($threadId: ID!) {
  resolveReviewThread(input: {threadId: $threadId}) {
    thread { isResolved }
  }
}' -f threadId="$thread_id" --jq '.data.resolveReviewThread.thread.isResolved'
