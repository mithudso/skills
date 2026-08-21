#!/usr/bin/env bash
# Usage: fetch-threads.sh <owner> <repo> <pr-number>
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
          id
          isResolved
          comments(first: 10) {
            nodes {
              author { login }
              body
              path
              line
              outdated
            }
          }
        }
      }
    }
  }
}' -f owner="$owner" -f repo="$repo" -F number="$pr_number" \
  --jq '.data.repository.pullRequest.reviewThreads.nodes[]
        | select(.isResolved == false)
        | {
            id: .id,
            path: .comments.nodes[0].path,
            line: .comments.nodes[0].line,
            outdated: .comments.nodes[0].outdated,
            author: .comments.nodes[0].author.login,
            comments: [.comments.nodes[].body]
          }'
