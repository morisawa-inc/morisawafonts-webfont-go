# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Test/Lint Commands

This project uses mise as the task runner:

```bash
mise tasks run lint     # Run golangci-lint
mise tasks run test     # Run tests with gotestsum
mise tasks run ci       # Run both lint and test
```

## Architecture

Go API client library for the Morisawa Fonts web font API. Requires Go 1.24+.

### Package Structure

- **`client/`** - Low-level HTTP client using resty.dev/v3. Handles GET/POST/DELETE requests, authentication, retries, and error handling.
- **`option/`** - Functional options pattern for client configuration (API token, base URL, timeout, retry count, custom HTTP client).
- **`pager/`** - Generic pagination with cursor-based paging and Go 1.22+ iterator support via `Iter()`.
- **`resource/domain/`** - Domain management operations (List, Add, Delete).
- **`internal/clienttest/`** - Test utilities providing pre-configured test clients.

### Key Patterns

- **Functional options**: Configuration via `option.WithAPIToken()`, `option.WithTimeout()`, etc.
- **Context-first APIs**: All operations take `context.Context` as first parameter.
- **Generic pagination**: `pager.Pager[T, M]` supports both manual pagination and iterator pattern.
- **Composition**: Main client embeds `client.Client` and exposes sub-clients (e.g., `Domains`).

## Testing

- Uses `github.com/jarcoal/httpmock` for HTTP mocking
- Uses `github.com/stretchr/testify/assert` for assertions
- Test client helper: `clienttest.NewClient(t)` creates pre-configured test client with cleanup
