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

If it is missing, install it (see the xenv README) and confirm with the user before changing their system. If `xenv` is unavailable or errors, stop and explain — do not silently fall back to reading `.env` files by hand unless the user explicitly authorizes it.

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

Pass the secret straight into a child process instead of exporting it into the shell:

```bash
xenv inject API_KEY -- curl -H "Authorization: Bearer $API_KEY" https://api.example.com
xenv inject GITHUB_TOKEN AWS_ACCESS_KEY_ID -- ./deploy.sh
```

The value exists only in the child process and disappears when it exits. `inject` fails before running the command if any requested key is missing, so a typo cannot silently run with an empty variable. The `--` separator is optional when the command has no flags (`xenv inject API_KEY printenv API_KEY`).

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

If the given file does not exist, `xenv` falls back to the process environment.

## Exit codes

- `check` exits `1` when the key does not exist.
- `inject` exits `1` when a requested key is missing or the command cannot be executed.

Use these instead of parsing output when scripting.

## Fallback

If `xenv` fails or is not installed:

1. Stop and report the failure and what you were trying to do.
2. Explain why reading the secret file directly is risky.
3. Only proceed another way if the user explicitly authorizes it.
