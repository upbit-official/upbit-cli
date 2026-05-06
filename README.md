English | [한국어](./README_KR.md)

# Upbit CLI

The official CLI for the [Upbit REST API](https://docs.upbit.com).

<!-- x-release-please-start-version -->

## Installation

### Installing with npm

```sh
npm install -g @upbit-official/upbit-cli
```

### Installing with Go

To test or install the CLI locally, you need [Go](https://go.dev/doc/install) version 1.22 or later installed.

```sh
go install 'github.com/upbit-official/upbit-cli/cmd/upbit@latest'
```

Once you have run `go install`, the binary is placed in your Go bin directory:

- **Default location**: `$HOME/go/bin` (or `$GOPATH/bin` if GOPATH is set)
- **Check your path**: Run `go env GOPATH` to see the base directory

If commands aren't found after installation, add the Go bin directory to your PATH:

```sh
# Add to your shell profile (.zshrc, .bashrc, etc.)
export PATH="$PATH:$(go env GOPATH)/bin"
```

<!-- x-release-please-end -->

### Running Locally

After cloning the git repository for this project, you can use the
`scripts/run` script to run the tool locally:

```sh
./scripts/run args...
```

## Usage

The CLI follows a resource-based command structure:

```sh
upbit [resource] <command> [flags...]
```

```sh
upbit accounts list \
--access-key "$UPBIT_ACCESS_KEY" \
--secret-key "$UPBIT_SECRET_KEY"
```

For details about specific commands, use the `--help` flag.

For more runnable examples, see the scripts in [`examples/`](examples/).

### Environment variables

| Environment variable | Description | Required | Default value |
| -------------------- | ----------- | -------- | ------------- |
| `UPBIT_ACCESS_KEY` | The access key provided by Upbit for API authentication. For more details, please refer to https://docs.upbit.com/reference/auth. | no | `null` |
| `UPBIT_SECRET_KEY` | The secret key used to sign API requests for secure verification. For more details, please refer to https://docs.upbit.com/reference/auth. | no | `null` |

### Global flags

- `--access-key` - The access key provided by Upbit for API authentication.
  For more details, please refer to https://docs.upbit.com/reference/auth.
  (can also be set with `UPBIT_ACCESS_KEY` env var)
- `--secret-key` - The secret key used to sign API requests for secure verification.
  For more details, please refer to https://docs.upbit.com/reference/auth.
  (can also be set with `UPBIT_SECRET_KEY` env var)
- `--help` - Show command line usage
- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version
- `--base-url` - Use a custom API backend URL
- `--environment` - Select API environment (`kr`, `sg`, `id`, `th`)
- `--format` - Change the output format (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--format-error` - Change the output format for errors (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--transform` - Transform the data output using [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
- `--transform-error` - Transform the error output using [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)

## Linking different Go SDK versions

You can link the CLI against a different version of the Upbit Go SDK
for development purposes using the `./scripts/link` script.

To link to a specific version from a repository (version can be a branch,
git tag, or commit hash):

```bash
./scripts/link github.com/org/repo@version
```

To link to a local copy of the SDK:

```bash
./scripts/link ../path/to/upbit-go
```

If you run the link script without any arguments, it will default to `../upbit-go`.

## Versioning

This package generally follows [SemVer](https://semver.org/spec/v2.0.0.html) conventions, though certain backwards-incompatible changes may be released as minor versions:

1. Changes to CLI internals which are technically public but not intended or documented for external use.
2. Changes that we do not expect to impact the vast majority of users in practice.

We take backwards-compatibility seriously and work hard to ensure you can rely on a smooth upgrade experience.

We are keen for your feedback; please contact us at open-api@upbit.com with questions, bugs, or suggestions.

## Contributing

The Upbit CLI is in its initial release phase, and public contributions (Issues/PRs) are currently closed.
For bug reports and feedback, please email open-api@upbit.com.
We are considering opening external contribution channels in phases as this project becomes more stable.

© 2026 Dunamu Inc. All rights reserved.
