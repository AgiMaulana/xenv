# xenv

AI optimized secret reader and redactor.

`xenv` lets AI agents and coding tools discover which secrets exist in your environment — key names and existence checks — without ever exposing their values.

## Features

- **Redacted reads** — lists every variable in your secret file with values replaced by `<redacted>`
- **Safe existence checks** — confirms whether a key exists, printing `<available>` or `<missing>` instead of the value
- **JSON output** — machine-friendly output for AI agents (`--json` / `-j`)
- **Custom secret file** — point at any file with `--input` / `-i` (defaults to `secret`)
- **Safe exports** — derives a shell `export` command from a secret without printing the value to an interactive terminal

## Recommended `AGENT.md` instruction

Add the following instruction to your project's `AGENT.md` (or equivalent agent instructions file) to help AI coding agents handle project configuration and secrets safely:

```md
## Environment variables and secrets

Use `xenv` whenever you need to check or use project-managed environment variables or secrets. Do not read `.env` files, secret files, shell history, or secret stores directly when `xenv` can provide the required information.

- Use `xenv check --json KEY` to verify whether a key exists.
- Use `xenv read --json` to discover available keys. Treat all returned values as redacted and never attempt to recover or guess them.
- When a secret must be supplied to a command, use `xenv export` with command substitution or `eval` as documented by `xenv`; do not print, log, copy, or otherwise expose the generated command.
- Never include secret values in chat responses, logs, patches, commits, screenshots, or test output.
- Prefer existence checks and redacted output over retrieving secret values.
- Ask the user before creating, changing, rotating, deleting, or otherwise modifying secret state.
- Do not bypass `xenv` if it is unavailable or fails. Stop and explain the failure, unless the user explicitly authorizes another method.
- Normal runtime environment variables explicitly provided by the current process may be used directly when appropriate; this policy applies to project-managed configuration and secrets.
```

This instruction provides agent guidance; it is not an access-control boundary. Use `xenv`'s permissions and redaction behavior as the enforcement layer.

## Installation

One-line install (downloads the latest release, verifies its signature and installs to `/usr/local/bin`):

```bash
curl -fsSL https://raw.githubusercontent.com/AgiMaulana/xenv/main/install.sh | sh
```

Install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/AgiMaulana/xenv/main/install.sh | sh -s -- v0.0.1
```

Install somewhere else (no sudo needed):

```bash
XENV_INSTALL_DIR="$HOME/.local/bin" sh -c "$(curl -fsSL https://raw.githubusercontent.com/AgiMaulana/xenv/main/install.sh)"
```

> The installer verifies each binary with [cosign](https://github.com/sigstore/cosign) keyless signing when cosign is available; otherwise it skips verification with a warning.

Or with Go:

```bash
go install github.com/AgiMaulana/xenv@latest
```

Or build from source:

```bash
git clone https://github.com/AgiMaulana/xenv
cd xenv
go build
```

## Usage

List all secret keys, with values redacted:

```bash
xenv read
```

Check whether a specific key exists:

```bash
xenv check DATABASE_URL
```

Use JSON output for AI-friendly results:

```bash
xenv read --json
# {"API_KEY":"<redacted>","DATABASE_URL":"<redacted>"}

xenv check DATABASE_URL --json
# {"key":"DATABASE_URL","exists":true}
```

Use a custom secret file:

```bash
xenv -i .env.production read
```

Export a secret into another environment variable. The command refuses to print the generated command when stdout is an interactive terminal, so use it with `eval`:

```bash
eval "$(xenv export --derive DATABASE_URL --key APP_DATABASE_URL)"
```

With a custom secret file:

```bash
eval "$(xenv -i .env.production export -d DATABASE_URL -k APP_DATABASE_URL)"
```

The export value is shell-quoted by `xenv`; avoid logging or otherwise exposing the command substitution output.

### Exit codes

`check` exits with code `1` when the key does not exist, so it can be used directly in shell scripts.

## Secret file format

A simple `KEY=VALUE` format. Blank lines and `#` comments are ignored:

```env
# Database
DATABASE_URL=postgres://localhost:5432/app
API_KEY=super-secret-value
```
