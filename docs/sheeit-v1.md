# `sheeit`: The spreadsheet CLI

[Neeraj Kashyap (zomglings)](mailto:nkashy1@gmail.com)

`sheeit` is a command-line tool which allows its users to drive spreadsheets from their terminal or, more importantly, through their AI agents (via the terminal).

Currently, `sheeit` can only be used to drive spreadsheets hosted on [Google Sheets](https://workspace.google.com/products/sheets/).

`sheeit` is inspired by [`clacks`](https://github.com/zomglings/clacks). It is free software (MIT licensed), hosted at [https://github.com/zomglings/sheeit](https://github.com/zomglings/sheeit)

`sheeit` is built primarily for use in AI environments such as Claude Code, Codex, and Cursor. It is meant to be easy for the AIs present in those environments to use while also costing very little in context tokens.

## CLI structure

`sheeit` is implemented in Go using [spf13/cobra](https://github.com/spf13/cobra).

```
sheeit version              # Print version information
sheeit completions <shell>  # Generate shell completions
sheeit auth <command>       # Authentication management
sheeit config <command>     # Configuration management
sheeit <operation>          # Spreadsheet operations
```

* **`sheeit version`** - Prints version information
* **`sheeit completions`** - Generates shell completion scripts. Supported shells: `bash`, `zsh`, `fish`, `powershell`
* **`sheeit auth`** - Authentication commands (`login`, `status`, `logout`)
* **`sheeit config`** - Configuration commands (`contexts`, `switch`)
* **Operational commands** - Live at the top level (e.g., `sheeit read`, not `sheeit values read`)

Operational commands are specified individually in separate GitHub issues. Each command corresponds to one or more capabilities from the "What users can do with Google Sheets" section.

## Functionality

### Authorization and authentication

`sheeit` authenticates to Google Sheets using OAuth 2.0, [as required by the Google Sheets API](https://developers.google.com/identity/protocols/oauth2?utm_source=chatgpt.com).

`sheeit` is a write-capable tool. It requires the [`https://www.googleapis.com/auth/spreadsheets`](https://www.googleapis.com/auth/spreadsheets) OAuth scope for Google Sheets ([reference](https://developers.google.com/workspace/sheets/api/scopes)).

The `sheeit` authorization process is interactive and will require the user to submit consent in a browser.

`sheeit` stores and reuses OAuth tokens locally to avoid repeated authorization prompts and to keep workflows fast in AI-driven environments. These tokens are stored in a platform-dependent credentials file, a [single TOML file](https://toml.io/en/) located at:

* Linux: `$XDG_CONFIG_HOME/sheeit/credentials.toml` or `~/.config/sheeit/credentials.toml`  
* macOS: `~/Library/Application Support/sheeit/credentials.toml`  
* Windows: `%APPDATA%\sheeit\credentials.toml`

The file must contain valid TOML. There is a single, required global field:

* `schema_version` (integer). This is the version of the credentials file schema.  The current supported value is `1`.

All other top-level tables are represented as *contexts* in the sense of `kubectl` or `aws` CLI. Each context is identified by the name of its top-level table. Context names:

* MUST be unique  
* MUST match `[A-Za-z0-9_-]+` (i.e. must be a valid TOML table name)  
* Are case-sensitive

Each context table must contain the following keys:

* `scopes` (array of strings). OAuth scopes granted for this context. For `sheeit`, this MUST include [`https://www.googleapis.com/auth/spreadsheets`](https://www.googleapis.com/auth/spreadsheets)  
* `token_uri` (string). OAuth 2.0 token endpoint. For Google, this MUST be [`https://oauth2.googleapis.com/token`](https://oauth2.googleapis.com/token)  
* `access_token` (string). OAuth 2.0 access token used for API requests.  
* `refresh_token` (string). OAuth 2.0 refresh token used to obtain new access tokens.  
* `expiry` (string, RFC 3339 timestamp). Expiration time of the access token.  
* `token_type` (string). Token type returned by the OAuth server (typically `Bearer`).

Unknown keys present in a credentials file MUST be ignored by `sheeit`.

### Authentication commands

#### Login

```
sheeit auth login -c CONTEXT [--overwrite]
```

Creates a new context by initiating the Google OAuth 2.0 browser flow. The user must grant consent in their browser.

* `-c, --context CONTEXT` (required). Name for the new context.
* `--overwrite`. Replace an existing context with the same name. Without this flag, `sheeit` MUST error if the context already exists.

On success, the new context is written to `credentials.toml`.

#### Status

```
sheeit auth status
```

Displays the current authentication status, including the current context name and token expiry.

#### Logout

```
sheeit auth logout [-c CONTEXT]
```

Removes a context from `credentials.toml`.

* `-c, --context CONTEXT`. The context to remove. Defaults to the current context if not specified.

If the removed context was the current context, `current_context` in `state.toml` is cleared.

### What users can do with Google Sheets

- Read cell values from a range
  [`spreadsheets.values.get`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets.values/get)

- Read cell values from multiple ranges in a single request
  [`spreadsheets.values.batchGet`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets.values/batchGet)

- Write or update cell values (including formulas)
  [`spreadsheets.values.update`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets.values/update)

- Write or update values across multiple ranges in a single request
  [`spreadsheets.values.batchUpdate`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets.values/batchUpdate)

- Append rows of data to a sheet
  [`spreadsheets.values.append`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets.values/append)

- Clear values from a range while preserving formatting
  [`spreadsheets.values.clear`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets.values/clear)

- Retrieve spreadsheet metadata and structure (sheets, properties, dimensions)
  [`spreadsheets.get`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets/get)

- Perform structural and formatting changes via batch update requests
  [`spreadsheets.batchUpdate`](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets/batchUpdate)

- Create, duplicate, rename, and delete sheets (tabs)
  [Batch update requests](https://developers.google.com/workspace/sheets/api/guides/batchupdate)

- Insert, delete, move, and resize rows and columns
  [Batch update requests](https://developers.google.com/workspace/sheets/api/guides/batchupdate)

- Apply cell formatting (number formats, text styles, alignment, colors)
  [Cell formatting](https://developers.google.com/workspace/sheets/api/guides/formats)

- Apply and manage conditional formatting rules
  [Conditional formatting](https://developers.google.com/workspace/sheets/api/guides/conditional-format)

- Define, update, and remove named and protected ranges
  [Named & protected ranges](https://developers.google.com/workspace/sheets/api/samples/ranges)

- Create and update pivot tables
  [Pivot tables](https://developers.google.com/workspace/sheets/api/samples/pivot-tables)

- Create and update embedded charts
  [Charts](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets/charts)

- Apply basic filters and filter views
  [Filters](https://developers.google.com/workspace/sheets/api/guides/filters)

- Set data validation rules on cells
  [Data operations](https://developers.google.com/workspace/sheets/api/samples/data)

- Add, update, and remove notes on cells
  [CellData](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets/cells)

- Merge and unmerge cells
  [Basic formatting](https://developers.google.com/sheets/api/samples/formatting)

- Find and replace text across ranges, sheets, or entire spreadsheets
  [Requests](https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets/request)

- Sort data in a range by one or more columns
  [Data operations](https://developers.google.com/workspace/sheets/api/samples/data)

- Attach developer metadata to rows, columns, sheets, or the spreadsheet
  [Developer metadata](https://developers.google.com/workspace/sheets/api/guides/metadata)

### CLI conventions

#### Output format

All commands support:

```
--format plain|json
```

* `plain` (default). Human-readable output.
* `json`. Machine-parseable JSON output.

#### Exit codes

* `0` - Success
* `1` - Any error

#### Error output

All errors are written to stderr. Data output goes to stdout. This applies regardless of `--format`.

#### Spreadsheet identification

Operational commands that act on a spreadsheet accept:

```
-s, --spreadsheet SPREADSHEET
```

This accepts:
* A spreadsheet ID (e.g., `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms`)
* A full Google Sheets URL
* Any other identifier that can be resolved to a spreadsheet ID

`sheeit` normalizes the input internally, which may involve API calls to resolve the identifier.

#### Range notation

Commands that operate on cell ranges accept any of the following notations:

* **A1 notation** - e.g., `Sheet1!A1:B10`, `A1:C5`
* **R1C1 notation** - e.g., `R1C1:R10C2`
* **Named ranges** - ranges defined in the spreadsheet by name

`sheeit` auto-detects the notation format.

### State file

`sheeit` stores operational state separately from credentials in a `state.toml` file located in the same configuration directory:

* Linux: `$XDG_CONFIG_HOME/sheeit/state.toml` or `~/.config/sheeit/state.toml`
* macOS: `~/Library/Application Support/sheeit/state.toml`
* Windows: `%APPDATA%\sheeit\state.toml`

The state file contains:

* `current_context` (string, optional). The name of the currently active context. If not set, `sheeit` MUST error when an operation requires a context.

Unknown keys present in the state file MUST be ignored by `sheeit`.

If the state file does not exist, `sheeit` MUST create it when a context is first switched.

### Context selection

`sheeit` supports multiple contexts, allowing users to manage credentials for different Google accounts or projects.

#### Listing contexts

```
sheeit config contexts [--limit LIMIT] [--offset OFFSET]
```

Lists all available contexts from `credentials.toml`. Supports pagination via `--limit` (default: no limit, shows all) and `--offset` (default: 0).

#### Switching contexts

```
sheeit config switch -C CONTEXT
```

Sets `current_context` in `state.toml` to the specified context name. The context MUST already exist in `credentials.toml`.

#### Context usage in commands

Operational commands (e.g., reading/writing spreadsheets) use the current context by default. All operational commands accept:

```
-c, --context CONTEXT
```

This overrides the current context for that command only, without modifying `state.toml`. The specified context MUST exist in `credentials.toml`.

#### Error handling

Context validation occurs when an operational command is executed, not at startup. `sheeit` MUST error if:

* The `current_context` (from `state.toml` or `-c` override) does not exist in `credentials.toml`
* The context's tokens are stale or expired and cannot be refreshed

Non-operational commands (e.g., `sheeit config contexts`) do not require a valid context and MUST NOT error due to context issues.

#### Configuration directory override

All `sheeit` subcommands accept:

```
-C, --config-dir CONFIG_DIR
```

This overrides the default platform-specific configuration directory. When specified, `sheeit` reads from `CONFIG_DIR/credentials.toml` and `CONFIG_DIR/state.toml` instead of the default locations.

