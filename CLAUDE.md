# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run a single test (and its subtests)
go test -run TestCORSPreflightHeaderValidation ./...

# Run one subtest specifically
go test -run TestCORSPreflightHeaderValidation/rejects_disallowed_headers ./...

# Race detector + verbose
go test -race -v ./...

# Benchmarks
go test -bench=. -benchmem ./...

# Static checks
go vet ./...

# Build everything (incl. examples/)
go build ./...
```

Requires Go 1.24 (see `go.mod`). Module path is `github.com/dbubel/intake/v2` — v2 is in the import path even though the repo has no v2/ subdirectory.

## Architecture

Intake is a thin layer over Go 1.22+ `http.ServeMux`. Everything composes around four pieces:

**`Intake` (intake.go)** — the router. Owns the `*http.ServeMux`, a slice of global `MiddleWare`, an optional panic handler, and a `registeredRoutes map[string][]string` registry used by `GetRoutes()` and `AddOptionsEndpoints()`. `AddEndpoint` registers a route using Go 1.22's `"VERB /path"` mux pattern (`fmt.Sprintf("%s %s", verb, path)`), so path parameters like `/users/{id}` work via `r.PathValue`.

**Middleware chain order (intake.go `AddEndpoint`)** — non-obvious and worth preserving:
1. Route-specific middleware wraps the final handler, applied in reverse so the *first* middleware passed is the outermost.
2. Global middleware then wraps that, also reverse-applied — so global middleware always runs *before* route-specific middleware at request time, and the first global middleware added is the outermost.
3. Panic recovery wraps everything last, so it catches panics from both global and route middleware. The `PanicHandler != nil` check is done at request time, so `SetPanicHandler` can be called before or after `AddEndpoint`.

This means: **call `AddGlobalMiddleware` before `AddEndpoint`.** Global middleware added after a route is registered won't apply to that route — it's captured into the chain at registration time.

**`endpoint` / `Endpoints` (endpoint.go, endpoints.go)** — declarative alternative to `AddEndpoint`. `endpoint` is unexported but is returned by exported constructors (`GET`, `POST`, `NewEndpoint`, etc.) — callers use it by value, not by name. `Endpoints` is `[]endpoint` with `Use`/`Append` (add middleware to the tail of each endpoint's chain) and `Prepend` (add to the head). `AddEndpoints` then feeds them into `AddEndpoint` one at a time.

**`CORS` (cors.go)** — middleware factory. At construction time `buildPolicy` precomputes lookup sets and joined header strings into a `corsPolicy` so request-time work is just map lookups and header writes. Origin matching supports three forms: exact match, `*` wildcard (disabled automatically if `AllowCredentials` is true), and subdomain patterns like `https://*.example.com` (apex domain is intentionally not matched). `CORS()` fills in defaults for empty `AllowedOrigins` / `AllowedMethods` / `AllowedHeaders` to match `DefaultCORSConfig`, so a zero `CORSConfig{}` is usable.

**`AddOptionsEndpoints` + CORS pairing** — calling order matters and is the typical preflight pattern:
1. `AddGlobalMiddleware(CORS(...))`
2. Register all real routes via `AddEndpoint` / `AddEndpoints`.
3. `AddOptionsEndpoints()` — synthesizes a no-op 204 OPTIONS handler for every path that doesn't already have one. The CORS middleware (added in step 1) wraps these and short-circuits preflight requests before the 204 handler runs.

**`Run` (intake.go)** — blocks on `ListenAndServe` in a goroutine, waits on SIGINT/SIGTERM, then `Shutdown` with a 60s timeout, falling back to `Close`. Returns the resulting error; `http.ErrServerClosed` is normalized to `nil`.

## Conventions

- Conventional Commits format for commit messages (e.g. `fix: ...`, `feat: ...`, `docs: ...`) — see recent `git log`.
- The project CLAUDE.md at `/Users/debubel/code/CLAUDE.md` defines a separate developer-agent team workflow (orchestrator/coder/jira-agent/github-agent). That's a global config; it isn't required for routine changes to this repo.
