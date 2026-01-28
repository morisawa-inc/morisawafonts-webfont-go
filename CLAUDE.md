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
- **`resource/`** - API resource implementations. Each sub-package corresponds to an API endpoint group (e.g., `domain/`, `stats/`).
- **`internal/clienttest/`** - Test utilities providing pre-configured test clients.

### Key Patterns

- **Functional options**: Configuration via `option.WithAPIToken()`, `option.WithTimeout()`, etc.
- **Context-first APIs**: All operations take `context.Context` as first parameter.
- **Generic pagination**: `pager.Pager[T, M]` supports both manual pagination and iterator pattern.
- **Composition**: Main client embeds `client.Client` and exposes resource-specific sub-clients.

## Testing

### Test Structure

- Test files follow Go convention (`*_test.go`) and are co-located with implementation files
- Test functions follow the pattern `TestTypeName_MethodName` (e.g., `TestDomains_List`, `TestPV_Get`)

### Test Client Helper

Use `clienttest.NewClient(t)` from `internal/clienttest/` to create pre-configured test clients:
- Pre-sets `http.DefaultClient` and a test API token (`"test-token"`)
- Accepts additional `option.Option` arguments to override defaults
- Automatically registers cleanup via `t.Cleanup()` for proper resource management

### HTTP Mocking with httpmock

Tests use `github.com/jarcoal/httpmock` for HTTP mocking:
- Call `httpmock.Activate(t)` at the start of each test (cleanup is automatic via `testing.T`)
- Register responders with `httpmock.RegisterResponder(method, url, handler)`
- Use `httpmock.NewJsonResponse()` for JSON responses and `httpmock.NewBytesResponse()` for non-content responses (e.g., 204 No Content)
- Use `httpmock.ConnectionFailure` to simulate network errors

### Pagination Testing

- Handlers use cursor query parameters to determine which page to return (typically via switch statements)
- Tests iterate through all pages using Go 1.22+ iterator syntax: `for item, err := range pager.Iter(ctx)`
- Validate pagination metadata: `HasNext` flag and `NextCursor` values for each page

### Assertions

- Uses `github.com/stretchr/testify/assert` for assertions
- Common assertions: `NoError`, `Equal`, `ErrorIs`, `ErrorAs`, `JSONEq`, `Len`
- Uses `github.com/samber/lo` utilities: `lo.ToPtr()` for pointer creation, `lo.Ternary()` for conditional values

### Test Data

- Test data is constructed inline within test functions (no separate fixture files)
- Uses realistic domain names (e.g., `1.example.com`, `2.example.com`) for readability
