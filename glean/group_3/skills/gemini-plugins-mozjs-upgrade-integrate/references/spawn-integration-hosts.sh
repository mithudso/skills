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
evergreen host create -d windows-2022-large -k "$KEY" --no-expire 2>&1

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
    done <<<"$NEW_IDS"

    $ALL_RUNNING && break
    sleep 15
done

# Give SSH a moment to become available after host reports running
sleep 10

LIST=$(evergreen host list --mine 2>/dev/null)
AFTER=$(echo "$LIST" | grep -oP "ID: \K[^;]+")
NEW_IDS=$(comm -13 <(echo "$BEFORE" | sort) <(echo "$AFTER" | sort))

# Run setup on each new host via SSH with agent forwarding.
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
    ssh -A -o StrictHostKeyChecking=no "$REMOTE_USER@$HOST" bash -s <<ENDSSH
set -eo pipefail
ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null
command -v cygpath &>/dev/null && git config --global core.sshCommand "C:/cygwin/bin/ssh"
git clone git@github.com:10gen/mongo.git
cd mongo
git config user.name "$GIT_NAME"
git config user.email "$GIT_EMAIL"

curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
export NVM_DIR="\$HOME/.nvm" && source "\$NVM_DIR/nvm.sh"
nvm install 22
nvm use --delete-prefix v22
npm install -g @anthropic-ai/claude-code

echo 'export LLVM_OBJDUMP=/opt/mongodbtoolchain/v5/bin/llvm-objdump' >> ~/.bashrc
$([ -n "$BRANCH" ] && echo "git fetch origin '$BRANCH' && git checkout '$BRANCH'")
echo "Setup complete on $DISTRO."
ENDSSH
done <<<"$NEW_IDS"

echo ""
echo "=== Ready — SSH Commands ==="
while IFS= read -r id; do
    HOST=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Host name: \K[^;]+")
    REMOTE_USER=$(echo "$LIST" | grep "ID: $id;" | grep -oP "User: \K[^;]+")
    DISTRO=$(echo "$LIST" | grep "ID: $id;" | grep -oP "Distro: \K[^;]+")
    echo "$DISTRO:  ssh -A ${REMOTE_USER}@${HOST}"
done <<<"$NEW_IDS"
