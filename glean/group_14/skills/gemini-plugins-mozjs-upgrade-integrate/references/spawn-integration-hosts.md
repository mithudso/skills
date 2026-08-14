# Spawning Integration Hosts (Linux x86_64 and Windows)

## Spawn the Hosts

```bash
evergreen host create -d amazon2023-latest-large -k <your-key-name> --no-expire
evergreen host create -d windows-2022-large      -k <your-key-name> --no-expire
```

Check your registered key names with `evergreen keys list`.

`--no-expire` is required because the default expiry (24 h) is too short for a multi-platform
upgrade. Unexpirable hosts are not automatically terminated — remember to destroy them when done:
`evergreen host terminate -h <host-id>`. You can also add a sleep schedule via the Evergreen UI
(Hosts → Edit) to stop the host on nights/weekends and reduce costs.

## Wait for Hosts to Be Running

```bash
evergreen host list --mine
```

Repeat until both show `Status: running`. Usually takes 2-3 minutes.

## Get SSH Connection Details

The host list output has the format:

```
ID: <id>; Name: ; Distro: <distro>; Status: running; Host name: <hostname>; User: <user>; ...
```

SSH commands:

```bash
# Linux x86_64
ssh -A ec2-user@<hostname>

# Windows x86_64
ssh -A Administrator@<hostname>
```

Always use `-A` (SSH agent forwarding) so GitHub access works on the remote host.

## Install Claude Code

### Linux

The workstation setup installs Node v8, which is too old. Install Node 22 via nvm first:

```bash
curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
export NVM_DIR="$HOME/.nvm" && source "$NVM_DIR/nvm.sh"
nvm install 22
nvm use --delete-prefix v22   # clears the ~/.npmrc prefix set by workstation setup
npm install -g @anthropic-ai/claude-code
```

Then authenticate. The easiest approach on a headless server is to copy credentials from your local
machine:

```bash
# Run this locally, not on the remote host
scp -r ~/.claude/ ec2-user@<hostname>:~/.claude/
```

Alternatively, set `ANTHROPIC_API_KEY` in the remote session before running `claude`.

### Windows

From the SSH session:

```powershell
# Install Node.js via winget if not present
winget install OpenJS.NodeJS.LTS

# Then install Claude Code
npm install -g @anthropic-ai/claude-code
```

Copy credentials the same way as Linux (adjust the user path):

```bash
# Run locally
scp -r ~/.claude/ Administrator@<hostname>:/Users/Administrator/.claude/
```

## After the Script: Manual Workstation Setup

Once the script completes and you SSH in, run workstation setup manually so you can complete the
EngFlow authentication interactively (it opens a browser link):

```bash
cd ~/mongo
bash etc/set_up_workstation.sh
source ~/.bashrc
source python3-venv/bin/activate
```

The EngFlow auth is optional for local builds but speeds up subsequent bazel invocations via remote
cache.

## One-Shot Script

Spawns hosts, waits until running, clones the mongo repo, installs Claude Code and LLVM_OBJDUMP,
then prints the SSH commands. Workstation setup is
intentionally excluded — run it manually after SSH-ing in so the EngFlow auth flow can complete
interactively.

```bash
#!/bin/bash
set -euo pipefail

KEY=${1:-pubKey2}
BRANCH=${2:-}        # optional: branch to fetch and checkout after setup
GIT_NAME=${3:-""}    # your git name, e.g. "Jane Smith"
GIT_EMAIL=${4:-""}   # your git email, e.g. "jane.smith@mongodb.com"

# Snapshot existing instance IDs before spawning so we can identify the new
# ones by diff. `host create` returns an Evergreen spawn ID but `host list`
# uses AWS instance IDs — they cannot be correlated directly.
BEFORE=$(evergreen host list --mine 2>/dev/null | grep -oP "ID: \K[^;]+")

echo "Spawning hosts..."
evergreen host create -d amazon2023-latest-large -k "$KEY" --no-expire 2>&1
evergreen host create -d windows-2022-large      -k "$KEY" --no-expire 2>&1

echo "Waiting for hosts to be running..."
while true; do
    LIST=$(evergreen host list --mine 2>/dev/null)
    AFTER=$(echo "$LIST" | grep -oP "ID: \K[^;]+")
    NEW_IDS=$(comm -13 <(echo "$BEFORE" | sort) <(echo "$AFTER" | sort))

    if [ -z "$NEW_IDS" ]; then
        echo "  waiting for new hosts to appear..."
        sleep 15
        continue
    fi

    ALL_RUNNING=true
    while IFS= read -r id; do
        STATUS=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Status: \K[^;]+")
        DISTRO=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Distro: \K[^;]+")
        echo "  $DISTRO ($id): ${STATUS:-pending}"
        [ "$STATUS" != "running" ] && ALL_RUNNING=false
    done <<< "$NEW_IDS"

    $ALL_RUNNING && break
    sleep 15
done

# Give SSH a moment to become available after host reports running
sleep 10

LIST=$(evergreen host list --mine 2>/dev/null)
AFTER=$(echo "$LIST" | grep -oP "ID: \K[^;]+")
NEW_IDS=$(comm -13 <(echo "$BEFORE" | sort) <(echo "$AFTER" | sort))

# Run setup on each new host via SSH with agent forwarding.
# --setup flag is not used because it runs before SSH and can't forward the
# GitHub SSH key needed to clone the private 10gen/mongo repo.
while IFS= read -r id; do
    HOST=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Host name: \K[^;]+")
    REMOTE_USER=$(echo "$LIST" | grep "ID: $id;" | grep -oP "User: \K[^;]+")
    DISTRO=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Distro: \K[^;]+")
    echo ""
    echo "Setting up $DISTRO ($HOST) — this takes a few minutes..."
    # Retry SSH until the daemon is ready (may take 10-20s after status=running)
    for _i in $(seq 1 10); do
        ssh -A -o StrictHostKeyChecking=no -o ConnectTimeout=10 "$REMOTE_USER@$HOST" true 2>/dev/null && break
        echo "  SSH not ready yet, retrying in 15s..."
        sleep 15
    done
    ssh -A -o StrictHostKeyChecking=no "$REMOTE_USER@$HOST" bash -s << ENDSSH
set -eo pipefail
ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null
# On Windows (Cygwin), MinGW git uses its own /bin/ssh path which doesn't exist.
# Fix: set core.sshCommand to the Cygwin ssh binary using a Windows-style path.
# The cygpath -m approach produced /bin/ssh in practice; hardcoding C:/cygwin/bin/ssh
# is reliable because Evergreen Windows hosts always install Cygwin to C:\cygwin.
command -v cygpath &>/dev/null && git config --global core.sshCommand "C:/cygwin/bin/ssh"
git clone git@github.com:10gen/mongo.git
cd mongo
git config user.name "$GIT_NAME"
git config user.email "$GIT_EMAIL"

# Workstation setup is intentionally skipped here — it blocks on EngFlow auth
# which requires interactive browser login. Run it manually after SSH-ing in.

# Node v8 is pre-installed but too old for Claude Code; install Node 22 via nvm.
curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
export NVM_DIR="\$HOME/.nvm" && source "\$NVM_DIR/nvm.sh"
nvm install 22
nvm use --delete-prefix v22  # clears the conflicting ~/.npmrc prefix set earlier
npm install -g @anthropic-ai/claude-code

# llvm-objdump is needed by extract.sh but is not on PATH by default.
# Do NOT add the full toolchain v5 bin to PATH — it switches to v5 g++ which
# statically links libstdc++, and moz.configure rejects that.
echo 'export LLVM_OBJDUMP=/opt/mongodbtoolchain/v5/bin/llvm-objdump' >> ~/.bashrc
$( [ -n "$BRANCH" ] && echo "git fetch origin '$BRANCH' && git checkout '$BRANCH'" )
echo "Setup complete on $DISTRO."
ENDSSH
done <<< "$NEW_IDS"

echo ""
echo "=== Ready — SSH Commands ==="
while IFS= read -r id; do
    HOST=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Host name: \K[^;]+")
    REMOTE_USER=$(echo "$LIST" | grep "ID: $id;" | grep -oP "User: \K[^;]+")
    DISTRO=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Distro: \K[^;]+")
    echo "$DISTRO:  ssh -A ${REMOTE_USER}@${HOST}"
done <<< "$NEW_IDS"
```

Usage:

```bash
bash spawn-integration-hosts.sh [key-name] [branch] [git-name] [git-email]
# e.g.:
bash spawn-integration-hosts.sh pubKey2 your-branch/SERVER-XXXXX "Jane Smith" "jane.smith@mongodb.com"
```
