#!/bin/bash
#
# refresh-mcp-auth.sh
# -------------------
# Re-establish the devprod-mcp-gateway (kanopy-oidc) session WITHOUT a local
# browser, then seed the proxy's token file so the running proxy uses it.
#
# Run this when the proxy needs a fresh session (e.g. each weekday morning, or
# after a >~12h gap such as a weekend). It prints a URL + device code; approve
# it on any browser/phone. If you DON'T approve, the device code expires after a
# few minutes and this script aborts WITHOUT touching your existing token.
#
# After it succeeds: reconnect the gateway with `/mcp` in Claude Code (or
# restart the session) so the proxy picks up the new token immediately.
#
# Usage:  ./refresh-mcp-auth.sh
#
# Notes:
#  - The proxy reads its OWN file (PROXY below) first and only falls back to the
#    CLI's native file if its own is MISSING. A stale-but-present file is NOT
#    replaced automatically, which is why this script overwrites it explicitly.
#  - Override the target env's token file with PROXY=... if you ever point at
#    staging, e.g. PROXY=~/.kanopy/token-devprod-mcp-proxy-staging.json ./refresh-mcp-auth.sh

set -euo pipefail

KANOPY_DIR="$HOME/.kanopy"
NATIVE="${NATIVE:-$KANOPY_DIR/token-oidclogin.json}"                 # kanopy-oidc CLI's own token cache
PROXY="${PROXY:-$KANOPY_DIR/token-devprod-mcp-proxy-prod.json}"      # the file the proxy reads

# --- preflight ---------------------------------------------------------------
command -v kanopy-oidc >/dev/null 2>&1 || {
  echo "ERROR: 'kanopy-oidc' is not on PATH. Install the darwin-arm64 build into ~/.local/bin." >&2
  exit 127
}
command -v python3 >/dev/null 2>&1 || {
  echo "ERROR: 'python3' is required (used to validate the token JSON)." >&2
  exit 127
}

# Remember the native file's mtime so we can tell whether login actually wrote
# a fresh token (vs. a silent no-op / stale cache).
before_mtime=$(stat -f %m "$NATIVE" 2>/dev/null || echo 0)

# --- interactive login (blocks until you approve OR the code expires) ---------
echo ">>> Starting kanopy-oidc device-flow login (no local browser)."
echo ">>> A verification URL + code will print below — approve it (phone is fine)."
echo ">>> If you don't approve, the device code expires in a few minutes and this aborts."
echo

# 'if !' keeps 'set -e' from firing so we can give a clean message on failure.
# We do NOT redirect stdout: the device code prints there, and you need to see it.
if ! kanopy-oidc login --flow device --no-browser; then
  echo >&2
  echo "ERROR: login did not complete (no approval, code expired, or denied)." >&2
  echo "       Your existing proxy token was left UNTOUCHED." >&2
  exit 1
fi

# --- validate what login produced --------------------------------------------
if [ ! -s "$NATIVE" ] || ! python3 -m json.tool "$NATIVE" >/dev/null 2>&1; then
  echo "ERROR: '$NATIVE' is missing or not valid JSON after login — NOT seeding the proxy." >&2
  exit 1
fi

after_mtime=$(stat -f %m "$NATIVE" 2>/dev/null || echo 0)
if [ "$after_mtime" = "$before_mtime" ]; then
  echo "WARN: '$NATIVE' was not updated by login (mtime unchanged); seeding it anyway." >&2
fi

if ! python3 -c "import json,sys; sys.exit(0 if json.load(open('$NATIVE')).get('refresh_token') else 1)"; then
  echo "WARN: the new token has NO refresh_token — silent refresh won't work and you'll re-auth sooner." >&2
fi

# --- seed the proxy file (back up the old one first) -------------------------
cp -p "$PROXY" "$PROXY.bak" 2>/dev/null || true
cp -p "$NATIVE" "$PROXY"

echo
echo "OK: seeded fresh token -> $PROXY"
echo "    backup of previous  -> $PROXY.bak"
echo ">>> Next: run '/mcp' in Claude Code (or restart) so the proxy adopts it now."
