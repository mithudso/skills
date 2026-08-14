#!/usr/bin/env bash
# Usage: verify-threads.sh <owner> <repo> <pr-number>
# Prints the count of unresolved threads. Should be 0 when done.
set -euo pipefail

owner=$1
repo=$2
pr_number=$3

gh api graphql -f query='
query($owner: String!, $repo: String!, $number: Int!) {
  repository(owner: $owner, name: $repo) {
    pullRequest(number: $number) {
      reviewThreads(first: 100) {
        nodes {
          isResolved
        }
      }
    }
  }
}' -f owner="$owner" -f repo="$repo" -F number="$pr_number" \
  --jq '[.data.repository.pullRequest.reviewThreads.nodes[] | select(.isResolved == false)] | length'
