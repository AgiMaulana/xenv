---
name: xenv
description: Discover and use project-managed environment variables and secrets safely with the xenv CLI, without ever exposing secret values. Use this skill whenever a task touches environment variables, .env files, API keys, access tokens, credentials, database URLs, or connection strings — including reading config, checking whether a key exists, or running a command that needs a secret. Reach for it before opening a .env or secrets file by hand, and prefer it over guessing or echoing values.
---

# xenv — safe secret discovery and use

`xenv` is a CLI that lets an agent see *which* secrets exist and *use* them, while keeping the values out of stdout, logs, patches, and chat. It reads from the process environment by default, or from a secret file with `-i`.

## Why this matters

Secret values leak through the most ordinary actions: echoing them to inspect, `cat`-ing a `.env`, printing a command before running it, or pasting output into a summary. Once a value lands in a transcript, log, or commit, it is effectively exposed. `xenv` splits "does this key exist / what is its name" from "run something that needs the value," so you can do real work without the value ever touching the message stream.

## When to apply

Use `xenv` any time the work involves:

- Listing, checking, or using environment variables or secrets
- Reading configuration that may contain credentials (`.env`, `.env.production`, `*.secrets`, etc.)
- Running a command, script, test, or server that needs a secret in its environment
- Remapping one secret into a differently named variable

## Preconditions

Confirm `xenv` is available before relying on it:

```bash
xenv version
```

If it is missing, install it (see the xenv README) and confirm with the user before changing their system. If it is unavailable when a secret is actually needed, follow the fallback guidance below rather than reading secret files directly.

## Core rules

These keep values out of the transcript. Follow them even when it seems faster to ignore them:

- Never `cat`, `grep`, or otherwise read a secret file directly when `xenv` can answer the question.
- Never echo, print, log, or hard-code a secret value.
- Never place a secret value in chat responses, commit messages, patches, screenshots, or test output — only key names and existence states.
- Prefer existence checks and redacted output over retrieving values at all.
- Ask the user before creating, changing, rotating, or deleting secret state.

## Workflow

### 1. Discover which keys exist

List every key with values redacted:

```bash
xenv read --json
# {"API_KEY":"<redacted>","DATABASE_URL":"<redacted>"}
```

`--json` (`-j`) is the machine-friendly form; use it when parsing. Plain `xenv read` is fine for a quick human-facing look.

### 2. Check a single key

```bash
xenv check DATABASE_URL --json
# {"key":"DATABASE_URL","exists":true,"state":"available"}
```

`check` reports `<available>` or `<not-exist>` and exits `1` when the key is missing, so it works directly in shell conditionals.

### 3. Run a command that needs a secret

`inject` puts the secret into the child process's environment, not the parent's. That distinction matters, because the shell expands `$VAR` in the command line *before* `xenv` runs — using the parent environment, where the injected value does not exist yet. So a bare command line silently drops the secret:

```bash
# WRONG: $API_KEY expands in the parent shell (usually empty) before inject runs
xenv inject API_KEY -- curl -H "Authorization: Bearer $API_KEY" https://api.example.com
```

Make the child do the expansion by wrapping the command in a shell with single quotes, so the parent shell passes `$API_KEY` through literally:

```bash
xenv inject API_KEY -- sh -c 'curl -H "Authorization: Bearer $API_KEY" https://api.example.com'
```

For anything non-trivial, put the command in a script that reads the variable and invoke the script:

```bash
xenv inject API_KEY -- ./call-api.sh   # call-api.sh uses "$API_KEY"
```

Multiple keys are allowed. Use the `--` separator to separate the keys from the command; it is optional only for a single key:

```bash
xenv inject GITHUB_TOKEN AWS_ACCESS_KEY_ID -- ./deploy.sh
xenv inject API_KEY printenv API_KEY
```

The value exists only in the child process and disappears when it exits. `inject` fails before running the command if any requested key is missing, so a typo cannot silently run with an empty variable.

### 4. Remap a key into another variable

When a tool expects a differently named variable, derive an `export` command and `eval` it. Never print the generated command — it contains the value:

```bash
eval "$(xenv export --derive DATABASE_URL --key APP_DATABASE_URL)"
```

`xenv export` refuses to print when stdout is an interactive terminal, which is a guardrail, not a bug. Always capture it through `eval`, as above.

### 5. Read from a secret file

Add `-i` / `--input` before the subcommand:

```bash
xenv -i .env.production read --json
xenv -i .env.production check DATABASE_URL
eval "$(xenv -i .env.production export -d DATABASE_URL -k APP_DATABASE_URL)"
```

If the given file does not exist, `xenv` silently falls back to the process environment. That fallback is worth calling out, because it undercuts what `-i` looks like it guarantees: `xenv -i .env.production check DATABASE_URL` returning `available` does **not** prove the value came from `.env.production` — it may have come from the process environment because the file was missing. When the source must be unambiguous, confirm the file exists before trusting the result, and say so rather than assuming.

## Exit codes

- `check` exits `1` when the key does not exist.
- `inject` exits `1` when a requested key is missing or the command cannot be executed.

Use these instead of parsing output when scripting.

## Fallback

If a task needs a secret and `xenv` is unavailable or errors, stop rather than working around it:

1. Stop before doing anything that would expose the value.
2. Report the failure and what you were trying to do.
3. Explain why reading the secret file directly is risky.
4. Only proceed another way if the user explicitly authorizes it.
