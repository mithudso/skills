#!/usr/bin/env bash
# setup.sh — clone-and-go bootstrapper for the skills repo (github.com/mithudso/skills).
# Installs Homebrew+Git, clones the repo to ~/.claude/skills, then: prompt keys
# (~/.claude/.env) -> lay down ~/.claude config (non-destructive) -> Ollama+model ->
# build index -> embed -> boot refresh agent -> merge MCP servers -> nightly cron.
# Idempotent and safe to re-run. Secrets are prompted, never committed.
# Flags: --repo <url> --dir <path> --skip-ollama --skip-service --skip-mcp
#        --skip-keys --skip-config --skip-cron --help
set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"   # repo root
SC_DIR="$DIR/skill-consolidation"                      # vendored index system
CC_DIR="$DIR/claude-config"                            # sanitized ~/.claude config
ENV_FILE="$HOME/.claude/.env"                          # secrets live here (git-ignored)
MODEL="${SKILLS_EMBED_MODEL:-qwen3-embedding:4b}"
OS="$(uname -s)"
DO_OLLAMA=1; DO_SERVICE=1; DO_MCP=1; DO_KEYS=1; DO_CONFIG=1; DO_CRON=1
REPO_URL="${REPO_URL:-https://github.com/mithudso/skills.git}"
TARGET_DIR="${TARGET_DIR:-$HOME/.claude/skills}"

while [ $# -gt 0 ]; do case "$1" in
  --skip-ollama) DO_OLLAMA=0 ;;
  --skip-service) DO_SERVICE=0 ;;
  --skip-mcp) DO_MCP=0 ;;
  --skip-keys) DO_KEYS=0 ;;
  --skip-config) DO_CONFIG=0 ;;
  --skip-cron) DO_CRON=0 ;;
  --repo) shift; REPO_URL="${1:-}" ;;
  --repo=*) REPO_URL="${1#--repo=}" ;;
  --dir) shift; TARGET_DIR="${1:-}" ;;
  --dir=*) TARGET_DIR="${1#--dir=}" ;;
  --help|-h) sed -n '2,8p' "$0" 2>/dev/null || true
    printf '\nUsage: ./setup.sh [--repo <git-url>] [--dir <path>] [--skip-ollama]\n'
    printf '       [--skip-service] [--skip-mcp] [--skip-keys] [--skip-config] [--skip-cron]\n'
    printf '\nStandalone bootstrap (no checkout yet):\n'
    printf '  bash setup.sh        # clones %s into %s\n' "$REPO_URL" "$TARGET_DIR"
    exit 0 ;;
  *) echo "unknown flag: $1 (try --help)"; exit 2 ;;
esac; shift; done

say(){ printf '\n\033[1m== %s\033[0m\n' "$1"; }
ok(){ printf '   \033[32mok\033[0m %s\n' "$1"; }
warn(){ printf '   \033[33m!!\033[0m %s\n' "$1"; }

brew_path(){ for b in /opt/homebrew/bin/brew /usr/local/bin/brew /home/linuxbrew/.linuxbrew/bin/brew; do [ -x "$b" ] && { echo "$b"; return 0; }; done; command -v brew 2>/dev/null || true; }
ensure_brew(){
  command -v brew >/dev/null 2>&1 && { ok "Homebrew present"; return 0; }
  local b; b="$(brew_path)"
  if [ -n "$b" ]; then eval "$("$b" shellenv)"; ok "Homebrew present"; return 0; fi
  warn "Homebrew not found — installing (non-interactive)"
  local script; script="$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" \
    || { warn "could not download the Homebrew installer"; return 1; }
  NONINTERACTIVE=1 /bin/bash -c "$script" || { warn "Homebrew install failed"; return 1; }
  b="$(brew_path)"; [ -n "$b" ] && eval "$("$b" shellenv)"
  command -v brew >/dev/null 2>&1 && ok "Homebrew installed" || { warn "Homebrew still not on PATH"; return 1; }
}
ensure_git(){
  command -v git >/dev/null 2>&1 && { ok "git $(git --version 2>/dev/null | awk '{print $3}')"; return 0; }
  warn "git not found — installing"
  if   command -v brew    >/dev/null 2>&1; then brew install git
  elif command -v apt-get >/dev/null 2>&1; then sudo apt-get update -y && sudo apt-get install -y git
  elif command -v dnf     >/dev/null 2>&1; then sudo dnf install -y git
  elif command -v yum     >/dev/null 2>&1; then sudo yum install -y git
  elif command -v pacman  >/dev/null 2>&1; then sudo pacman -S --noconfirm git
  else warn "no package manager available — install git manually, then re-run"; return 1; fi
  command -v git >/dev/null 2>&1 && ok "git installed"
}

# Already inside the cloned skills repo?
IN_REPO=0; [ -f "$SC_DIR/gen-skills-index.mjs" ] && IN_REPO=1

# ── 1. Dependencies: Homebrew (macOS) + Git ──────────────────────────────────
say "1/11 Dependencies (Homebrew + Git)"
if [ "$OS" = "Darwin" ]; then ensure_brew || warn "continuing without Homebrew (later steps may be skipped)"
else command -v brew >/dev/null 2>&1 && ok "Homebrew present" || warn "no Homebrew (Linux uses native package managers below)"; fi
ensure_git || { [ "$IN_REPO" -eq 1 ] || { echo "   git is required to clone the repo"; exit 1; }; }

# ── 2. Clone / update the repo, then re-exec from inside it ───────────────────
say "2/11 Repository"
if [ "$IN_REPO" -eq 1 ]; then
  ok "running from the repo ($DIR)"
elif [ "${SKILLS_BOOTSTRAPPED:-0}" = "1" ]; then
  echo "   bootstrap re-exec did not land inside the repo at $TARGET_DIR"; exit 1
else
  [ -n "$REPO_URL" ] || { echo "   need a repo URL: ./setup.sh --repo <git-url> [--dir <path>]"; exit 2; }
  if [ -d "$TARGET_DIR/.git" ]; then
    git -C "$TARGET_DIR" pull --ff-only && ok "updated existing checkout at $TARGET_DIR" || warn "pull failed; using existing checkout"
  else
    mkdir -p "$(dirname "$TARGET_DIR")"
    git clone "$REPO_URL" "$TARGET_DIR" && ok "cloned into $TARGET_DIR"
  fi
  REEXEC=""
  [ "$DO_OLLAMA"  -eq 0 ] && REEXEC="$REEXEC --skip-ollama"
  [ "$DO_SERVICE" -eq 0 ] && REEXEC="$REEXEC --skip-service"
  [ "$DO_MCP"     -eq 0 ] && REEXEC="$REEXEC --skip-mcp"
  [ "$DO_KEYS"    -eq 0 ] && REEXEC="$REEXEC --skip-keys"
  [ "$DO_CONFIG"  -eq 0 ] && REEXEC="$REEXEC --skip-config"
  [ "$DO_CRON"    -eq 0 ] && REEXEC="$REEXEC --skip-cron"
  say "Re-running from the clone"
  exec env SKILLS_BOOTSTRAPPED=1 bash "$TARGET_DIR/setup.sh" $REEXEC
fi

# ── 3. Node ≥ 18 ─────────────────────────────────────────────────────────────
say "3/11 Node.js"
if ! command -v node >/dev/null 2>&1 && command -v brew >/dev/null 2>&1; then brew install node || true; fi
command -v node >/dev/null 2>&1 || { echo "   node not found. Install Node ≥18 (e.g. 'brew install node') and re-run."; exit 1; }
NODE_MAJOR="$(node -e 'process.stdout.write(String(process.versions.node.split(".")[0]))')"
[ "$NODE_MAJOR" -ge 18 ] || { echo "   need Node ≥18 (have $(node -v))"; exit 1; }
ok "$(node -v)"

# ── 4. Configure API keys (~/.claude/.env) ───────────────────────────────────
say "4/11 Configure API keys ($ENV_FILE)"
mkdir -p "$HOME/.claude"
if [ ! -f "$ENV_FILE" ]; then
  if [ -f "$SC_DIR/.env.example" ]; then cp "$SC_DIR/.env.example" "$ENV_FILE"; else : > "$ENV_FILE"; fi
  ok ".env created"
fi
if [ "$DO_KEYS" -eq 1 ] && [ -t 0 ]; then
  set_env(){ ENV_K="$1" ENV_V="$2" awk '
      BEGIN{k=ENVIRON["ENV_K"]; v=ENVIRON["ENV_V"]; done=0}
      $0 ~ "^"k"=" {print k"="v; done=1; next} {print}
      END{if(!done) print k"="v}' "$ENV_FILE" > "$ENV_FILE.tmp" && mv "$ENV_FILE.tmp" "$ENV_FILE"; }
  prompt_key(){ local k="$1" label="$2" req="${3:-0}" cur val
    cur="$(awk -F= -v k="$k" '$1==k{sub(/^[^=]*=/,"");print;exit}' "$ENV_FILE")"
    if [ -n "$cur" ]; then ok "$label already set (kept)"; return 0; fi
    if [ "$req" = 1 ]; then printf '   %s (required, Enter to skip): ' "$label"
    else printf '   %s (optional, Enter to skip): ' "$label"; fi
    read -r val || val=""
    if   [ -n "$val" ]; then set_env "$k" "$val"; ok "$label saved"
    elif [ "$req" = 1 ]; then warn "$label left blank — its MCP server(s) will be disabled"
    else warn "$label skipped — its MCP server is disabled"; fi; }
  prompt_key MONDAY_API_TOKEN          "Monday.com API token"           1
  prompt_key FIRECRAWL_API_KEY         "Firecrawl API key"              0
  prompt_key EXA_API_KEY               "Exa API key"                    0
  prompt_key MDB_MCP_API_CLIENT_ID     "MongoDB MCP API client id"      0
  prompt_key MDB_MCP_API_CLIENT_SECRET "MongoDB MCP API client secret"  0
  prompt_key DASHBOARD_API_TOKEN       "TAM dashboard API token"        0
  ok "keys captured; servers without keys stay disabled (step 10)"
else
  warn "non-interactive or --skip-keys — edit $ENV_FILE by hand, then re-run"
fi

# ── 5. Lay down ~/.claude config (non-destructive) ───────────────────────────
say "5/11 Lay down ~/.claude config"
if [ "$DO_CONFIG" -eq 0 ]; then
  warn "skipped (--skip-config)"
elif [ ! -d "$CC_DIR" ]; then
  warn "no claude-config/ in the repo — skipping"
else
  install_cfg(){  # src target — render __HOME__->$HOME; back up an existing, differing target
    local src="$1" tgt="$2" tmp; mkdir -p "$(dirname "$tgt")"; tmp="$(mktemp)"
    sed "s#__HOME__#$HOME#g" "$src" > "$tmp"
    if [ -f "$tgt" ]; then cmp -s "$tmp" "$tgt" && { rm -f "$tmp"; return 0; }; cp "$tgt" "$tgt.bak-$(date +%s)"; fi
    mv "$tmp" "$tgt"
  }
  [ -f "$CC_DIR/settings.json.template" ] && install_cfg "$CC_DIR/settings.json.template" "$HOME/.claude/settings.json"
  [ -f "$CC_DIR/mcp.json.template" ]      && install_cfg "$CC_DIR/mcp.json.template"      "$HOME/.claude/mcp.json"
  [ -f "$CC_DIR/CLAUDE.md" ]              && install_cfg "$CC_DIR/CLAUDE.md"              "$HOME/.claude/CLAUDE.md"
  [ -f "$CC_DIR/policy-limits.json" ]     && install_cfg "$CC_DIR/policy-limits.json"     "$HOME/.claude/policy-limits.json"
  for d in agents commands hooks prompts; do
    [ -d "$CC_DIR/$d" ] || continue
    while IFS= read -r f; do rel="${f#"$CC_DIR"/}"; install_cfg "$f" "$HOME/.claude/$rel"; done < <(find "$CC_DIR/$d" -type f)
  done
  ok "config laid down (existing files backed up to *.bak-<ts>; ~/.claude/.env untouched)"
  warn "some hooks/MCP servers reference private repos you may not have — they no-op until present"
fi

# ── 6. Ollama + embedding model (boot service) ───────────────────────────────
say "6/11 Embedding server (Ollama + $MODEL)"
if [ "$DO_OLLAMA" -eq 0 ]; then
  warn "skipped (--skip-ollama); semantic --search will run keyword-only until a server is up"
else
  if ! command -v ollama >/dev/null 2>&1; then
    if [ "$OS" = "Darwin" ] && command -v brew >/dev/null 2>&1; then brew install ollama
    else warn "ollama not found. Install it (https://ollama.com/download) then re-run, or use --skip-ollama"; fi
  fi
  if command -v ollama >/dev/null 2>&1; then
    if [ "$OS" = "Darwin" ] && command -v brew >/dev/null 2>&1; then
      brew services start ollama >/dev/null 2>&1 && ok "ollama running as a brew service (persists across reboot)" \
        || warn "could not start the brew service (a manual 'ollama serve' may own port 11434)"
    else
      warn "start ollama at boot yourself (e.g. 'systemctl --user enable --now ollama' or a login item)"
    fi
    for _ in 1 2 3 4 5 6 7 8; do curl -fsS http://localhost:11434/api/tags >/dev/null 2>&1 && break; sleep 1; done
    ollama pull "$MODEL" >/dev/null 2>&1 && ok "model $MODEL present" || warn "could not pull $MODEL (pull it manually later)"
  fi
fi

# ── 7. Build the regex index (offline, no server) ────────────────────────────
say "7/11 Build SKILLS-INDEX.{json,md}"
node "$SC_DIR/gen-skills-index.mjs" >/dev/null && ok "index regenerated from ~/.claude/skills"

# ── 8. Build / sync the vector corpus (best-effort, fail-open) ───────────────
say "8/11 Embed SKILLS-EMBEDDINGS.json"
if node "$SC_DIR/gen-skills-index.mjs" --embed >/dev/null 2>&1; then
  ok "vectors in sync"
else
  warn "embed skipped — server down or model missing. Re-run 'node skill-consolidation/gen-skills-index.mjs --embed' later"
fi

# ── 9. Boot refresh agent (macOS launchd / Linux systemd) ────────────────────
say "9/11 Auto-refresh on boot + on skill changes"
if [ "$DO_SERVICE" -eq 0 ]; then
  warn "skipped (--skip-service)"
elif [ "$OS" = "Darwin" ]; then
  PLIST="$HOME/Library/LaunchAgents/com.skills-embed.plist"
  mkdir -p "$HOME/Library/LaunchAgents"
  sed -e "s#__DIR__#$SC_DIR#g" -e "s#__HOME__#$HOME#g" "$SC_DIR/com.skills-embed.plist.template" > "$PLIST"
  plutil -lint "$PLIST" >/dev/null
  launchctl bootout "gui/$(id -u)/com.skills-embed" >/dev/null 2>&1 || true
  launchctl bootstrap "gui/$(id -u)" "$PLIST" 2>/dev/null || launchctl load "$PLIST" 2>/dev/null || true
  launchctl enable "gui/$(id -u)/com.skills-embed" 2>/dev/null || true
  ok "launchd agent com.skills-embed loaded (RunAtLoad + WatchPaths + 6h)"
else
  UD="$HOME/.config/systemd/user"; mkdir -p "$UD"
  for u in service timer path; do
    sed -e "s#__DIR__#$SC_DIR#g" -e "s#__HOME__#$HOME#g" "$SC_DIR/skills-embed.$u.template" > "$UD/skills-embed.$u"
  done
  systemctl --user daemon-reload 2>/dev/null || true
  if systemctl --user enable --now skills-embed.timer skills-embed.path 2>/dev/null; then
    ok "systemd user units enabled (timer + path watch)"
  else
    warn "systemd --user unavailable; run the refresh from cron or a login shell instead"
  fi
fi

# ── 10. Merge MCP servers into ~/.claude.json (non-destructive) ──────────────
say "10/11 MCP servers"
if [ "$DO_MCP" -eq 0 ]; then
  warn "skipped (--skip-mcp)"
elif [ ! -f "$ENV_FILE" ]; then
  warn "no $ENV_FILE — fill keys (step 4) then re-run, or pass --skip-mcp"
else
  set -a; . "$ENV_FILE"; set +a   # export every var in .env for the node merge below
  SETUP_DIR="$SC_DIR" node <<'NODE'
const fs=require('fs'), path=require('path'), os=require('os');
const DIR=process.env.SETUP_DIR, HOME=os.homedir();
if(!process.env.SKILLS_RELAY_PATH) process.env.SKILLS_RELAY_PATH=path.join(HOME,'.claude','scripts','skills-relay.js');
const tpl=JSON.parse(fs.readFileSync(path.join(DIR,'mcp-servers.template.json'),'utf8')).mcpServers||{};
const sub=s=>s.replace(/\$\{([A-Z0-9_]+)\}/g,(_,k)=> (process.env[k]&&process.env[k].length)?process.env[k]:'${'+k+'}');
const deep=v=>typeof v==='string'?sub(v):Array.isArray(v)?v.map(deep):(v&&typeof v==='object')?Object.fromEntries(Object.entries(v).map(([k,x])=>[k,deep(x)])):v;
const unresolved=v=>/\$\{[A-Z0-9_]+\}/.test(JSON.stringify(v));
const cfgPath=path.join(HOME,'.claude.json');
let cfg={}; try{cfg=JSON.parse(fs.readFileSync(cfgPath,'utf8'));}catch{}
cfg.mcpServers=cfg.mcpServers||{};
const added=[],kept=[],skipped=[];
for(const [name,raw] of Object.entries(tpl)){
  if(name.startsWith('_')) continue;
  const r=deep(raw);
  if(unresolved(r)){ skipped.push(name); continue; }     // missing creds/path in .env
  if(cfg.mcpServers[name]){ kept.push(name); continue; }  // never overwrite existing
  cfg.mcpServers[name]=r; added.push(name);
}
if(added.length){
  if(fs.existsSync(cfgPath)) fs.copyFileSync(cfgPath, cfgPath+'.bak-'+Date.now());
  fs.writeFileSync(cfgPath, JSON.stringify(cfg,null,2));
}
const p=(label,a)=>a.length?console.log(`   ${label}: ${a.join(', ')}`):null;   // names only, never values
p('added', added); p('kept existing', kept); p('skipped (fill .env)', skipped);
NODE
  ok "~/.claude.json updated (backup written; existing servers untouched)"
fi

# ── 11. Nightly auto-update (cron: git pull + re-embed) ──────────────────────
say "11/11 Nightly auto-update (cron)"
if [ "$DO_CRON" -eq 0 ]; then
  warn "skipped (--skip-cron)"
elif ! command -v crontab >/dev/null 2>&1; then
  warn "crontab not available — skipping nightly auto-update"
elif ! git -C "$DIR" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  warn "$DIR is not a git checkout — nightly git-pull auto-update not installed"
else
  MARK="# skills nightly auto-update"
  CRON_PATH="/opt/homebrew/bin:/usr/local/bin:/home/linuxbrew/.linuxbrew/bin:/usr/bin:/bin"
  LINE="13 3 * * * /bin/bash -lc 'export PATH=\"$CRON_PATH:\$PATH\"; cd \"$DIR\" && git pull --ff-only && \"$SC_DIR/embed-refresh.sh\"' >> \"$DIR/auto-update.log\" 2>&1 $MARK"
  ( crontab -l 2>/dev/null | grep -vF "$MARK"; echo "$LINE" ) | crontab - \
    && ok "nightly cron installed (03:13 daily: git pull + re-embed); your other cron jobs preserved" \
    || warn "could not update crontab"
fi

say "Done"
echo "   Read skill-consolidation/SKILLS-INDEX.md to route by keyword;"
echo "   'node skill-consolidation/gen-skills-index.mjs --search \"...\"' for semantic search."
echo "   Restart Claude Code to pick up MCP server + config changes."