<!-- hub-reference-banner -->
> **Reference file — part of the `devops-linux-admin` hub.** Formerly the standalone `shell-scripting` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: shell-scripting
title: Shell Scripting (Bash & Zsh)
description: >
  Write, debug, and optimize Bash and Zsh scripts with production-grade error handling, parameter expansion, glob qualifiers, text-processing pipelines, dotfile management, CLI productivity techniques, and the most dangerous scripting anti-patterns. Covers POSIX fundamentals through zsh-specific extensions, signal handling, parallel execution, and modern toolchain replacements (ripgrep, fd, fzf, jq, yq).
  TRIGGER: user writes, reviews, or debugs a Bash or Zsh script; asks about parameter expansion, quoting, arrays, getopts, traps, signal handling, here-docs, process substitution, glob qualifiers, dotfile startup order, or text-processing pipelines (grep/sed/awk/jq/yq); migrating from bash to zsh; ShellCheck findings.
  SKIP: Python/Ruby/Node scripting (use the language skill); Windows batch or PowerShell scripts; Kubernetes YAML (use kubernetes-networking); complex JSON/YAML generation better handled in a real language.
version: "1.1"
category: developer
updated: "2026-05-29"
tags:
  - bash
  - zsh
  - shell
  - scripting
  - cli
  - terminal
  - unix
  - linux
  - macos
  - text-processing
  - dotfiles
  - productivity
related_skills:
  - linux-sysadmin
  - python-patterns
  - javascript-nodejs
---

# Shell Scripting (Bash & Zsh)

## When to Use

- Writing or reviewing shell scripts for automation, CI/CD, or system administration
- Debugging unexpected variable expansion, quoting, or pipeline exit-code behavior
- Configuring dotfiles (.bashrc, .zshrc, .zshenv) or a prompt framework
- Building CLI tools with argument parsing, signal handling, and proper cleanup
- Constructing text-processing pipelines with grep/sed/awk/jq/yq
- Migrating from bash to zsh or adopting modern CLI replacements (ripgrep, fd, fzf)
- Identifying and fixing common anti-patterns flagged by ShellCheck

## When Not to Use

- Complex data manipulation better handled by Python, Ruby, or Node.js
- Performance-critical code that needs real data structures (use a real language)
- Cross-platform scripts that must run on Windows without WSL
- Tasks requiring JSON/YAML generation with complex nesting (use jq/yq or a language with native types)

---

## 1. Bash Fundamentals

### 1.1 Variables and Quoting

```bash
name="Alice"
count=42

echo "Hello, $name"   # Hello, Alice (double-quotes expand variables)
echo 'Hello, $name'   # Hello, $name (single quotes are literal)

# Always quote variable expansions to prevent word splitting and globbing
cp "$source" "$dest"   # safe
cp $source $dest       # DANGEROUS if paths contain spaces
```

### 1.2 Parameter Expansion

```bash
# Default value — use $val if set and non-null, otherwise "default"
echo "${val:-default}"

# Assign default — also writes back to $val
echo "${val:=default}"

# Error if unset or null
echo "${val:?'val is required'}"

# Use alternate value — return "set" only if $val is non-empty
echo "${val:+"set"}"

# String length
echo "${#name}"          # 5

# Substring: ${var:offset:length}
echo "${name:1:3}"       # lic

# Remove shortest prefix matching pattern
path="/usr/local/bin/bash"
echo "${path#*/}"        # usr/local/bin/bash

# Remove longest prefix (basename equivalent)
echo "${path##*/}"       # bash

# Remove shortest suffix (dirname equivalent)
echo "${path%/*}"        # /usr/local/bin

# Remove longest suffix
echo "${path%%/*}"       # (empty)

# Search and replace (first / all matches)
echo "${name/l/L}"       # ALice
echo "${name//l/L}"      # ALice (global)

# Case conversion (bash 4+)
echo "${name^^}"         # ALICE
echo "${name,,}"         # alice
echo "${name^}"          # Alice (first char only)
```

### 1.3 Arrays

```bash
# Indexed array
fruits=("apple" "banana" "cherry")

echo "${fruits[0]}"          # apple
echo "${fruits[@]}"          # all elements (always prefer @)
echo "${#fruits[@]}"         # 3 (length)
echo "${fruits[@]:1:2}"      # banana cherry (slice)

fruits+=("date")             # append

for fruit in "${fruits[@]}"; do
  echo "$fruit"
done

# Associative array (bash 4+)
declare -A ages
ages["alice"]=30
ages["bob"]=25
echo "${ages["alice"]}"      # 30
echo "${!ages[@]}"           # keys: alice bob
echo "${ages[@]}"            # values: 30 25
```

### 1.4 Arithmetic

```bash
result=$(( 3 * 4 + 2 ))      # 14

(( count++ ))
(( count += 5 ))

if (( count > 10 )); then echo "big"; fi

# Float arithmetic via bc or awk
echo "scale=2; 22/7" | bc
awk 'BEGIN { printf "%.2f\n", 22/7 }'
```

### 1.5 Here-Docs and Here-Strings

```bash
# Quoted delimiter — no variable expansion
cat <<'EOF'
Literal $text here.
EOF

# Unquoted delimiter — variables and commands expand
name="World"
cat <<EOF
Hello, $name!
Today is $(date +%F).
EOF

# Strip leading tabs with <<-
if true; then
  cat <<-EOF
	This line has its leading tab stripped.
	EOF
fi

# Passing here-doc to a command
mysql -u root <<SQL
  CREATE DATABASE mydb;
SQL

# Here-string — single-line redirect without creating a subshell
grep "pattern" <<< "some text to search"
read -r first rest <<< "one two three"
```

### 1.6 Process Substitution

```bash
# <(cmd) — treat command output as a readable file descriptor
diff <(sort file1.txt) <(sort file2.txt)
comm <(sort listA) <(sort listB)

# >(cmd) — treat command input as a writable file descriptor
tee >(gzip > log.gz) > log.txt

# Avoid subshell in while-read loop (variable mutations visible outside)
while IFS= read -r line; do
  echo "line: $line"
done < <(find . -name "*.txt")
```

---

## 2. Zsh Enhancements

### 2.1 Extended Globbing and Glob Qualifiers

```zsh
setopt extended_glob   # enable advanced patterns

ls ^*.txt              # everything except .txt files (negate with ^)
ls **/*.js             # recursive glob

# Glob qualifiers appended in parentheses after the pattern
ls *(.)                # regular files only
ls *(/)                # directories only
ls *(@)                # symbolic links only
ls *(*)                # executables only
ls *(r)                # owner-readable
ls *(m-7)              # modified within last 7 days
ls *(Lm+1)             # larger than 1 MB
ls *(om)               # sorted by modification time (newest first)
ls *(Lk+100.)          # regular files > 100 KB
ls -d *(/om[1])        # most recently modified directory

# Numeric ranges
ls file<1-100>.txt     # file1.txt through file100.txt
```

### 2.2 Zsh Associative Arrays

```zsh
typeset -A colors
colors=(red "#FF0000" green "#00FF00" blue "#0000FF")

echo $colors[red]            # #FF0000
echo ${(k)colors}            # keys: red green blue
echo ${(v)colors}            # values: #FF0000 #00FF00 #0000FF

for key val in "${(@kv)colors}"; do
  echo "$key => $val"
done

(( ${+colors[red]} )) && echo "red exists"
```

### 2.3 Zsh-Specific Parameter Expansion Flags

```zsh
name="hello world"
echo ${(U)name}              # HELLO WORLD
echo ${(C)name}              # Hello World

words=(one two three)
echo ${(j:, :)words}         # one, two, three (join with separator)

str="a:b:c"
parts=(${(s/:/)str})         # split into array on ":"
echo $parts[2]               # b

echo ${(F)words}             # join with newlines

list=(c b a b)
echo ${(ou)list}             # a b c (sort + unique)

file="my file with spaces"
echo ${(q)file}              # my\ file\ with\ spaces (shell-quoted)
```

### 2.4 precmd, preexec, and chpwd Hooks

```zsh
# precmd — runs before each prompt is drawn
precmd() {
  local branch
  branch=$(git symbolic-ref --short HEAD 2>/dev/null)
  RPROMPT="${branch:+%F{yellow}[$branch]%f}"
}

# preexec — runs after Enter, before command executes; receives command as $1
preexec() {
  echo "Running: $1"
}

# chpwd — runs whenever the working directory changes
chpwd() {
  ls --color=auto
}

# Register multiple hooks without overwriting
autoload -U add-zsh-hook
add-zsh-hook precmd  my_precmd_function
add-zsh-hook preexec my_preexec_function

RPROMPT='%F{green}%~%f'    # right-side prompt in green
ZLE_RPROMPT_INDENT=0        # remove trailing space at right margin
```

### 2.5 Useful Zsh Options

```zsh
setopt AUTO_CD            # cd by typing directory name alone
setopt CORRECT            # suggest corrections for mistyped commands
setopt HIST_IGNORE_DUPS   # skip duplicate history entries
setopt SHARE_HISTORY      # share history across sessions in real time
setopt EXTENDED_GLOB      # enable advanced glob patterns
setopt NO_CASE_GLOB       # case-insensitive globbing
setopt COMPLETE_ALIASES   # complete through aliases
setopt PUSHD_IGNORE_DUPS  # no duplicate entries in dirstack
```

---

## 3. Shell Scripting Patterns

### 3.1 Strict Mode

```bash
#!/usr/bin/env bash
set -euo pipefail
# -e          exit immediately on command failure
# -u          error on unset variables (use ${var:-} for optional vars)
# -o pipefail pipeline returns exit code of first failed command

IFS=$'\n\t'   # prevent accidental word splitting in loops

# Allow a command to fail intentionally
grep -q "pattern" file || true
```

### 3.2 Cleanup with trap

```bash
#!/usr/bin/env bash
set -euo pipefail

TMPDIR=""
LOCKFILE=""

cleanup() {
  local exit_code=$?
  [[ -d "$TMPDIR" ]] && rm -rf "$TMPDIR"
  [[ -f "$LOCKFILE" ]] && rm -f "$LOCKFILE"
  exit "$exit_code"
}

trap cleanup EXIT INT TERM

TMPDIR=$(mktemp -d)
LOCKFILE="/tmp/myscript.lock"
```

### 3.3 Lock Files with flock

```bash
#!/usr/bin/env bash
LOCKFILE="/var/run/myscript.lock"

exec 9>"$LOCKFILE"
if ! flock --nonblock 9; then
  echo "Another instance is running. Exiting." >&2
  exit 1
fi
trap 'flock --unlock 9' EXIT
```

### 3.4 Argument Parsing with getopts

```bash
#!/usr/bin/env bash
usage() {
  echo "Usage: $0 [-v] [-o output_file] [-n count] [args...]" >&2
  exit 1
}

VERBOSE=false
OUTPUT=""
COUNT=10

while getopts ":vo:n:" opt; do
  case $opt in
    v) VERBOSE=true ;;
    o) OUTPUT="$OPTARG" ;;
    n) COUNT="$OPTARG" ;;
    :) echo "Option -$OPTARG requires an argument." >&2; usage ;;
    \?) echo "Unknown option: -$OPTARG" >&2; usage ;;
  esac
done
shift $(( OPTIND - 1 ))   # remove parsed options; $@ = remaining positional args
```

### 3.5 Temporary Files

```bash
# Always use mktemp — never hand-craft /tmp/myscript.$$
tmpfile=$(mktemp)
tmpdir=$(mktemp -d)
tmpfile=$(mktemp --suffix=.json)

# macOS mktemp requires Xs in template
tmpfile=$(mktemp /tmp/myscript.XXXXXX)

trap 'rm -f "$tmpfile"; rm -rf "$tmpdir"' EXIT
```

### 3.6 Parallel Execution

```bash
# xargs -P — run up to N jobs in parallel
find . -name "*.log" | xargs -P4 -I{} gzip {}

# GNU parallel — more control
parallel -j4 gzip ::: *.log
parallel --jobs 8 --progress 'process_file {1}' ::: *.csv

# Background jobs with wait
for file in *.png; do
  convert "$file" "${file%.png}.jpg" &
done
wait

# Capture exit codes from background jobs
pids=()
for item in "${items[@]}"; do
  process "$item" & pids+=($!)
done
for pid in "${pids[@]}"; do
  wait "$pid" || echo "Job $pid failed"
done
```

### 3.7 Signal Handling

```bash
#!/usr/bin/env bash

handle_sigterm() {
  echo "SIGTERM received — cleaning up..." >&2
  exit 0
}

handle_sigint() {
  echo "" >&2
  echo "Interrupted by user" >&2
  exit 130   # convention: 128 + signal number (SIGINT=2)
}

trap handle_sigterm TERM
trap handle_sigint  INT
```

---

## 4. CLI Productivity

### 4.1 Readline Keybindings (Emacs Mode)

| Key | Action |
|-----|--------|
| Ctrl+A | Move to beginning of line |
| Ctrl+E | Move to end of line |
| Alt+F / Alt+B | Forward / backward one word |
| Ctrl+K | Kill (cut) to end of line |
| Ctrl+U | Kill to beginning of line |
| Ctrl+W | Kill previous word |
| Ctrl+Y | Yank (paste) killed text |
| Ctrl+R | Reverse incremental history search |
| Ctrl+L | Clear screen |
| Alt+. | Insert last argument of previous command |
| Ctrl+X Ctrl+E | Open current line in $EDITOR |

```bash
# Switch to Vi mode
set -o vi          # bash
bindkey -v         # zsh

# Custom ~/.inputrc
"\C-p": history-search-backward
"\C-n": history-search-forward
set completion-ignore-case on
set show-all-if-ambiguous on
```

### 4.2 History Expansion

```bash
!!          # repeat last command
!$          # last argument of previous command
!^          # first argument of previous command
!*          # all arguments of previous command
!n          # command number n from history
!string     # most recent command starting with "string"
^old^new    # quick substitution: replace "old" with "new" in last command
!!:p        # print last command without executing
!!:h        # head (directory part of last arg)
!!:t        # tail (filename of last arg)

# Practical: re-run with sudo
sudo !!
```

### 4.3 Brace Expansion

```bash
echo {jpg,png,gif}                   # jpg png gif
cp file.{txt,bak}                    # cp file.txt file.bak
mkdir -p project/{src,test,docs}

echo {1..10}                         # 1 2 3 4 5 6 7 8 9 10
echo {a..z}
echo {01..10}                        # zero-padded
echo {1..20..3}                      # step=3: 1 4 7 10 13 16 19

echo {A,B}{1,2}                      # A1 A2 B1 B2 (nested)
```

### 4.4 Process Management

```bash
long_command &         # run in background
jobs -l                # list jobs with PIDs
fg %1                  # bring job 1 to foreground
bg %1                  # resume suspended job in background
Ctrl+Z                 # suspend foreground job
disown %1              # remove from job table (keeps running)
nohup cmd > out.log 2>&1 &  # immune to SIGHUP from start
```

### 4.5 tmux Session Management

```bash
tmux new-session -s work
tmux new-session -s deploy -n logs
tmux list-sessions
tmux attach -t work
# Inside tmux (prefix Ctrl+B):
# c = new window, % = split vertical, " = split horizontal
# z = zoom pane, d = detach, n/p = next/prev window, arrows = move pane
tmux send-keys -t work:0 "npm run dev" Enter
tmux kill-session -t work
```

---

## 5. Dotfile Management

### 5.1 Shell Startup File Order

**Bash** interactive login shell: `/etc/profile` then first of `~/.bash_profile`, `~/.bash_login`, `~/.profile`
**Bash** interactive non-login: `/etc/bash.bashrc` then `~/.bashrc`

**Zsh** load order:
1. `~/.zshenv` — always sourced (env vars, PATH); keep it lean
2. `~/.zprofile` — login shells only
3. `~/.zshrc` — interactive shells (aliases, functions, completions)
4. `~/.zlogin` — login shells (after zshrc)

```bash
# ~/.bash_profile — keep minimal; source .bashrc for interactive config
[[ -f ~/.bashrc ]] && source ~/.bashrc
```

### 5.2 PATH Management

```bash
# Bash: helper functions to avoid duplicates
path_prepend() {
  [[ ":$PATH:" != *":$1:"* ]] && export PATH="$1:$PATH"
}
path_prepend "$HOME/.local/bin"
path_prepend "$HOME/go/bin"

# Zsh: use $path array (tied bidirectionally to $PATH)
typeset -U path          # -U enforces uniqueness automatically
path=("$HOME/.local/bin" $path)
```

### 5.3 Oh My Zsh

```bash
# Install
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

# .zshrc key settings
ZSH_THEME="robbyrussell"
plugins=(git node npm docker kubectl fzf z)

omz update
```

### 5.4 Starship Prompt

```bash
# Install (cross-shell, written in Rust)
curl -sS https://starship.rs/install.sh | sh

# Enable in .zshrc / .bashrc
eval "$(starship init zsh)"
eval "$(starship init bash)"

# ~/.config/starship.toml
# [character]
# success_symbol = "[>](bold green)"
# error_symbol   = "[x](bold red)"
#
# [git_branch]
# truncation_length = 20
#
# [nodejs]
# format = "via [ $version](bold green) "
```

### 5.5 Shell Completion

```bash
# Bash: install bash-completion
[[ -r "/opt/homebrew/etc/profile.d/bash_completion.sh" ]] && \
  source "/opt/homebrew/etc/profile.d/bash_completion.sh"

source <(kubectl completion bash)
source <(helm completion bash)

# Zsh: compinit
autoload -Uz compinit && compinit

zstyle ':completion:*' menu select
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Z}'   # case-insensitive

fpath=(~/.zsh/completions $fpath)
```

---

## 6. Text Processing Pipelines

### 6.1 grep and ripgrep

```bash
grep "pattern" file.txt
grep -i "pattern" file.txt        # case-insensitive
grep -r "pattern" ./src/          # recursive
grep -n "pattern" file.txt        # line numbers
grep -l "pattern" *.txt           # filenames only
grep -v "pattern" file.txt        # invert
grep -E "regex+" file.txt         # extended regex
grep -A3 -B1 "error" log.txt      # 3 lines after, 1 before

# ripgrep (rg) — faster, .gitignore-aware
rg "pattern"
rg -i "pattern"
rg -t js "pattern"                # only .js files
rg -l "pattern"                   # filenames only
rg "pattern" --json | jq '.'      # machine-readable output
```

### 6.2 sed

```bash
sed 's/old/new/' file.txt          # first match per line
sed 's/old/new/g' file.txt         # all matches per line
sed -i 's/foo/bar/g' file.txt      # in-place (GNU; macOS: -i '')
sed -i.bak 's/foo/bar/g' file.txt  # in-place with backup
sed '/pattern/d' file.txt          # delete matching lines
sed '5d' file.txt                  # delete line 5
sed -n '10,20p' file.txt           # print lines 10-20
sed -n '/start/,/end/p' file.txt   # print between patterns
sed -e 's/foo/bar/' -e 's/baz/qux/' file.txt   # multiple expressions
```

### 6.3 awk

```bash
awk '{print $1, $3}' file.txt
awk -F: '{print $1}' /etc/passwd    # custom delimiter

awk '$3 > 100 {print $1, $3}' data.txt

awk '{sum += $2} END {print sum}' file.txt

awk '{count[$1]++} END {for (k in count) print k, count[k]}' file.txt

awk 'BEGIN {FS=","; OFS="\t"} {print $2, $1} END {print "done"}' csv

awk '/ERROR/ {printf "[%s] %s\n", $1, $NF}' app.log
```

### 6.4 cut, sort, uniq, tr

```bash
cut -d: -f1,3 /etc/passwd          # fields 1 and 3
cut -c1-10 file.txt                # characters 1-10

sort file.txt                      # lexicographic
sort -n file.txt                   # numeric
sort -rn file.txt                  # reverse numeric
sort -k2,2n file.txt               # sort by field 2 numerically
sort -t: -k3,3n /etc/passwd

sort file.txt | uniq               # deduplicate
sort file.txt | uniq -c            # count occurrences
sort file.txt | uniq -d            # duplicates only

# Canonical "top N" pattern
sort | uniq -c | sort -rn | head -10

tr 'a-z' 'A-Z' <<< "hello"        # HELLO
tr -d '\r' < windows.txt           # remove carriage returns
tr -s ' ' < file.txt               # squeeze multiple spaces
```

### 6.5 jq and yq

```bash
jq '.'                   data.json  # pretty-print
jq '.name'               data.json  # extract field
jq '.users[].email'      data.json  # iterate array
jq 'select(.age > 30)'  data.json   # filter
jq '[.users[] | {name, email}]' data.json   # reshape
jq -r '.name'            data.json  # raw string (no quotes)
jq -c '.'                data.json  # compact output
jq --arg k "v" '.[$k]'  data.json   # variable injection
jq -s '.[0] * .[1]' a.json b.json  # merge two objects

# Build JSON from shell variables
jq -n --arg name "$name" --argjson count "$count" \
  '{name: $name, count: $count}'

# yq (YAML processor, jq-compatible API)
yq '.spec.replicas' deployment.yaml
yq '.metadata.name = "new-name"' deployment.yaml
yq -o=json '.' config.yaml         # YAML to JSON
yq 'del(.status)' deployment.yaml
```

### 6.6 xargs

```bash
find . -name "*.log" | xargs rm

# Null-delimited — safe for filenames with spaces
find . -name "*.log" -print0 | xargs -0 rm

# Inject at specific position
find . -name "*.txt" | xargs -I{} cp {} {}.bak

# Parallel execution
find . -name "*.png" | xargs -P8 -I{} convert {} -resize 800x {}

# Limit items per invocation
echo {1..10} | xargs -n3 echo

# Find files containing a pattern
find . -name "*.js" | xargs grep -l "TODO"
```

---

## 7. Scripting Anti-Patterns

### 7.1 Unquoted Variables (Word Splitting and Globbing)

```bash
# WRONG — word splits on spaces; * globs expand
file="my report.txt"
rm $file              # tries to rm "my" and "report.txt"
for f in $(ls); do .. # broken

# CORRECT
rm "$file"
for f in ./*; do ..   # use glob, not ls
```

### 7.2 Parsing ls Output

```bash
# WRONG — filenames can contain spaces, newlines, special characters
for f in $(ls /tmp); do process "$f"; done

# CORRECT — use globs or find with -print0
for f in /tmp/*; do
  [[ -f "$f" ]] && process "$f"
done

find /tmp -maxdepth 1 -type f -print0 | while IFS= read -r -d '' f; do
  process "$f"
done
```

### 7.3 Useless Use of cat (UUOC)

```bash
# WRONG — cat adds an unnecessary process
cat file.txt | grep "pattern"
cat file.txt | wc -l

# CORRECT — redirect directly
grep "pattern" file.txt
wc -l < file.txt

# cat IS appropriate when concatenating multiple files
cat file1.txt file2.txt > combined.txt
```

### 7.4 Not Checking Exit Codes

```bash
# WRONG — silent failure
mkdir /new/dir
cd /new/dir        # silently goes wrong place if mkdir failed
rm -rf *           # catastrophe

# CORRECT
mkdir -p /new/dir || { echo "mkdir failed" >&2; exit 1; }
cd /new/dir || exit 1
```

### 7.5 Pipeline Exit Code Blindness

```bash
# WRONG — $? only captures exit code of last command (tee)
cat broken.txt | process_data | tee output.txt
echo $?    # always 0 from tee

# CORRECT
set -o pipefail

# Or inspect PIPESTATUS array (bash)
cat broken.txt | process_data | tee output.txt
echo "${PIPESTATUS[@]}"   # exit codes of each pipeline stage
```

### 7.6 Using [ ] Instead of [[ ]] in Bash

```bash
# WRONG — fragile with empty or whitespace values
[ $var == "value" ]    # syntax error if $var is empty

# CORRECT — [[ ]] is a bash keyword; no word splitting inside
[[ $var == "value" ]]
[[ $var =~ ^[0-9]+$ ]]  # regex matching
[[ -z $var ]]           # test empty string
```

### 7.7 Subshell Variable Scoping in Pipelines

```bash
# WRONG — variable mutated inside pipeline subshell is invisible outside
count=0
cat file.txt | while IFS= read -r line; do
  (( count++ ))   # modifies subshell copy
done
echo "$count"     # still 0!

# CORRECT — use process substitution to keep while in the main shell
count=0
while IFS= read -r line; do
  (( count++ ))
done < <(cat file.txt)
echo "$count"     # correct!

# Or use mapfile (bash 4+)
mapfile -t lines < file.txt
count=${#lines[@]}
```

### 7.8 Missing IFS= -r in read Loops

```bash
# WRONG — strips leading/trailing whitespace; interprets backslashes
while read line; do ...

# CORRECT
while IFS= read -r line; do
  echo "$line"
done < file.txt
```

---

## 8. Modern CLI Toolchain

| Classic | Modern | Key advantage |
|---------|--------|---------------|
| grep -r | rg (ripgrep) | .gitignore-aware, 5-10x faster |
| find | fd | intuitive syntax, colorized, respects .gitignore |
| cat | bat | syntax highlighting, Git integration |
| ls | eza / lsd | colors, icons, tree view, Git status |
| cd | zoxide | frecency-based smart directory jumping |
| history search | fzf | fuzzy search across history, files, Git branches |
| man | tldr | practical examples-first pages |

```bash
# fd — find replacement
fd "*.js"
fd -t f -e json
fd --hidden --no-ignore "\.env"

# fzf — integrate into shell (add to .bashrc/.zshrc)
eval "$(fzf --bash)"    # or --zsh
# Ctrl+R = fuzzy history search; Ctrl+T = fuzzy file insert; Alt+C = fuzzy cd
fzf --preview 'bat --color=always {}'
kill $(ps aux | fzf | awk '{print $2}')

# zoxide
z projects/myapp        # jump to frecent match
zi                      # interactive selection with fzf
```

---

## References

1. [GNU Bash Reference Manual — Shell Parameter Expansion](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html)
2. [GNU Bash Reference Manual — Arrays](https://www.gnu.org/software/bash/manual/html_node/Arrays.html)
3. [GNU Bash Reference Manual — Readline Init File Syntax](https://www.gnu.org/software/bash/manual/html_node/Readline-Init-File-Syntax.html)
4. [Bash Hackers Wiki — Parameter Expansion](https://bash-hackers.gabe565.com/syntax/pe/)
5. [Greg's Wiki — BashGuide Parameters](https://mywiki.wooledge.org/BashGuide/Parameters)
6. [Greg's Wiki — Word Splitting](https://mywiki.wooledge.org/WordSplitting)
7. [Zsh Documentation — Expansion](https://zsh.sourceforge.io/Doc/Release/Expansion.html)
8. [Mastering Zsh — Hooks](https://github.com/rothgar/mastering-zsh/blob/master/docs/config/hooks.md)
9. [Zsh ArchWiki](https://wiki.archlinux.org/title/Zsh)
10. [How To Use set and pipefail — How-To Geek](https://www.howtogeek.com/782514/how-to-use-set-and-pipefail-in-bash-scripts-on-linux/)
11. [Shell Scripting Best Practices for Production](https://oneuptime.com/blog/post/2026-02-13-shell-scripting-best-practices/view)
12. [Handling Signals in Bash — Baeldung](https://www.baeldung.com/linux/bash-signal-handling)
13. [Linux Job Control: disown and nohup — Baeldung](https://www.baeldung.com/linux/job-control-disown-nohup)
14. [How To Use Bash Job Control — DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-use-bash-s-job-control-to-manage-foreground-and-background-processes)
15. [Bash Process Substitution — Linux Handbook](https://linuxhandbook.com/bash-process-substitution/)
16. [Bash Heredoc Complete Guide — Linuxize](https://linuxize.com/post/bash-heredoc/)
17. [Moving to zsh Part 2: Configuration Files — Scripting OS X](https://scriptingosx.com/2019/06/moving-to-zsh-part-2-configuration-files/)
18. [Starship Cross-Shell Prompt](https://starship.rs/faq/)
19. [Text Processing in the Shell — Balthazar Rouberol](https://blog.balthazar-rouberol.com/text-processing-in-the-shell)
20. [grep, sed and awk Essentials — mylinux.work](https://mylinux.work/guides/grep-sed-awk-essentials/)
21. [Useless Use of Cat — Tom Ryder](https://blog.sanctum.geek.nz/useless-use-of-cat/)
22. [ShellCheck SC2002 — Useless cat](https://www.shellcheck.net/wiki/SC2002)
23. [ParsingLs — Wooledge Wiki](http://bash.cumulonim.biz/ParsingLs.html)
24. [Modern CLI Tools: fd and ripgrep](https://brtkwr.com/posts/2025-07-21-modern-cli-tools-fd-and-ripgrep/)
25. [CLI Essentials: fd/rg/fzf/bat — Medium](https://addozhang.medium.com/cli-essentials-fd-rg-fzf-bat-6137c34d5d1b)
26. [ZSH Gem: Extended Globbing — Refining Linux](https://www.refining-linux.org/archives/37-ZSH-Gem-2-Extended-globbing-and-expansion.html)
27. [Tmux Session Management — terminal.guide](https://www.terminal.guide/tools/multiplexer/tmux/session-guide/)
