#!/bin/bash
#
# Resolves merge conflicts from the "Removed all Non-SpiderMonkey Files" cherry-pick.
#
# The first cherry-pick deletes thousands of files, creating conflicts where git
# leaves "Version HEAD of <file> left in tree" messages. This script parses those
# messages, deletes the conflicting files, stages the changes, and continues the
# cherry-pick.
#
# Usage: ./resolve_delete_non_spidermonkey_files_cherry_pick.sh <input_file> <output_file>
#   input_file:  the captured output of `git cherry-pick <hash>`
#   output_file: scratch file where extracted filenames are written

set -euo pipefail

if [ $# -ne 2 ]; then
    echo "Usage: $0 <cherry_pick_output_file> <output_file>"
    exit 1
fi

input_file="$1"
output_file="$2"

if [ ! -f "$input_file" ]; then
    echo "ERROR: Input file '$input_file' does not exist."
    exit 1
fi

>"$output_file"

# Extract filenames from "Version HEAD of <path> left in tree" messages.
count=0
while read -r line; do
    result=$(echo "$line" | sed -n 's/.*Version HEAD of \(.*\) left in tree.*/\1/p')
    if [ -n "$result" ]; then
        echo "$result" >>"$output_file"
        count=$((count + 1))
    fi
done <"$input_file"

echo "Found ${count} files to remove."

# Delete each file.
removed=0
while read -r fileToRemove; do
    if [ -e "$fileToRemove" ]; then
        rm "$fileToRemove"
        removed=$((removed + 1))
    fi
done <"$output_file"

echo "Removed ${removed} files."

# Stage all deletions.
git add -u

# Continue the cherry-pick (auto-accept the commit message).
GIT_EDITOR=true git cherry-pick --continue
