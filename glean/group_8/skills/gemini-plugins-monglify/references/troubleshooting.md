# Monglify Troubleshooting

## `kanopy-oidc` Not Found

If you get a "command not found" error for `kanopy-oidc`, tell the user to:

1. Download the latest release from https://github.com/kanopy-platform/kanopy-oidc/releases
2. Place it somewhere on their `$PATH` (e.g., `/usr/local/bin/`)
3. Run `kanopy-oidc login` from a normal terminal to complete the initial login flow

Then retry the Monglify command.

## 401 / Unauthorized Errors

If you get a 401 or other Unauthorized error:

1. Check if `~/.kanopy/token-oidclogin.json` contains `"expiry": "0001-01-01T00:00:00Z"` — this is an invalid expiry caused by a bug in older versions of `kanopy-oidc`.
2. If so, explain to the user that an older version of `kanopy-oidc` wrote an invalid token file, and they need to download a newer version from https://github.com/kanopy-platform/kanopy-oidc/releases.
3. Ask the user for permission before deleting `~/.kanopy/token-oidclogin.json` (use a question tool if available). Do not delete without permission.
4. After the file is deleted, ask the user to run `kanopy-oidc login` from a normal terminal to re-authenticate.
