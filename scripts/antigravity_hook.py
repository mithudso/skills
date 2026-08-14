#!/usr/bin/env python3
import sys
import json
import subprocess
import os

def main():
    try:
        # Read hook payload from stdin
        raw = sys.stdin.read()
        if raw.strip():
            payload = json.loads(raw)
        else:
            payload = {}
    except Exception:
        payload = {}

    # Trigger background sync script
    trigger_script = "/Users/mitch.hudson/dev/skills/scripts/trigger_sync.sh"
    try:
        subprocess.Popen([trigger_script], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, start_new_session=True)
    except Exception:
        pass

    # Always output valid empty JSON object
    sys.stdout.write("{}\n")
    sys.stdout.flush()

if __name__ == "__main__":
    main()
