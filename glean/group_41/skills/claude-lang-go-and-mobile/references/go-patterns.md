<!-- hub-reference-banner -->
> **Reference file — part of the `lang-go-and-mobile` hub.** Formerly the standalone `go-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: go-patterns
version: "1.0.0"
updated: "2026-05-29"
category: developer
tags: [go, golang, patterns, concurrency, goroutines, channels, interfaces, error-handling, testing, modules, context, http, grpc]
description: "Go (Golang) production patterns expert — idiomatic Go code, concurrency primitives, error handling, interfaces, testing, HTTP servers, and module management. TRIGGER: user is writing or reviewing Go code; asking about goroutines, channels, mutexes, or context propagation; implementing Go interfaces or embedding; handling errors idiomatically in Go; writing Go tests with table-driven patterns or testify; setting up Go modules or workspace; building HTTP servers or gRPC services in Go; profiling or benchmarking Go code; migrating from Go generics syntax. SKIP: user is asking about a different language (Rust, C++, Java) even if the concepts are similar; pure algorithm questions with no Go-specific concern; infrastructure tooling written in Go but where the question is about the tool's behavior (kubectl, Terraform), not the Go code itself."
related_skills: [javascript-nodejs, python-patterns, typescript-expert, docker-containers, cicd-pipelines]
---

# Go Production Patterns

## When not to use

Skip this skill when the user is asking about a different language, or about a Go-based tool's behavior (kubectl, Terraform, etc.) rather than Go code they are writing.

## Key decisions

### Error handling

Go errors are values. The idiomatic pattern:

```go
// Return errors; never use panic for expected failure paths
func fetchUser(id string) (*User, error) {
    if id == "" {
        return nil, fmt.Errorf("fetchUser: id must not be empty")
    }
    // ...
}

// Wrap with context using %w (Go 1.13+) so callers can unwrap
if err != nil {
    return fmt.Errorf("fetchUser %s: %w", id, err)
}

// Sentinel errors for switch/Is checks
var ErrNotFound = errors.New("not found")

// Custom error types for structured inspection
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}

// Check with errors.Is / errors.As (not == on wrapped errors)
if errors.Is(err, ErrNotFound) { ... }
var ve *ValidationError
if errors.As(err, &ve) { ... }
```

### Concurrency

**Goroutines + channels** for pipeline and fan-out patterns; **sync primitives** for shared state.

```go
// Fan-out: spread work across N workers
func fanOut(jobs <-chan Job, results chan<- Result, n int) {
    var wg sync.WaitGroup
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    go func() { wg.Wait(); close(results) }()
}

// Context cancellation — always propagate ctx
func doWork(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-longOperation():
        _ = result
        return nil
    }
}

// Mutex for shared state (prefer RWMutex when reads dominate)
type SafeCounter struct {
    mu sync.RWMutex
    v  map[string]int
}
func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.v[key]++
}
func (c *SafeCounter) Value(key string) int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.v[key]
}
```

**Goroutine leak prevention:** every goroutine needs a termination path — context cancellation, channel close, or explicit stop signal.

### Interfaces

Keep interfaces small (1–3 methods). Define interfaces at the point of use (consumer side), not the provider side.

```go
// Good: small, consumer-defined interface
type Storer interface {
    Store(ctx context.Context, key string, val []byte) error
}

// Bad: large interface on the provider
type Database interface { // 20 methods — hard to mock, couples callers
    ...
}

// Embedding to compose
type ReadWriter interface {
    io.Reader
    io.Writer
}
```

### HTTP servers

```go
func NewServer(addr string, store Storer) *http.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /items/{id}", handleGetItem(store))
    mux.HandleFunc("POST /items", handleCreateItem(store))
    return &http.Server{
        Addr:         addr,
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
}

// Always set timeouts — default http.Server has none
// Go 1.22+: use method+path pattern routing ("GET /items/{id}")
```

### Graceful shutdown

```go
srv := NewServer(":8080", store)
go srv.ListenAndServe()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("shutdown: %v", err)
}
```

### Testing — table-driven pattern

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name    string
        a, b    int
        want    int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, 1, 0},
        {"zeros", 0, 0, 0},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := Add(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}

// Subtests run in parallel with t.Parallel()
// Use testify/require for assertion helpers when the stdlib is verbose
```

### Generics (Go 1.18+)

Use generics to eliminate type switch boilerplate and write reusable containers, not to abstract over unrelated types.

```go
// Typed collection with constraints
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

// Constraint with union
type Number interface { ~int | ~int64 | ~float64 }

func Sum[T Number](s []T) T {
    var total T
    for _, v := range s {
        total += v
    }
    return total
}
```

### Context rules

- Always accept `ctx context.Context` as the **first parameter** of functions that do I/O or spawn goroutines.
- Never store context in a struct — pass it through function arguments.
- Use `context.WithValue` only for request-scoped metadata (trace IDs, auth principals), never for optional parameters.
- Cancel contexts when you create them: `ctx, cancel := context.WithTimeout(...); defer cancel()`.

## Common pitfalls and fixes

| Symptom | Cause | Fix |
|---------|-------|-----|
| Goroutine leak | No termination path | Add ctx cancellation or channel close in every goroutine |
| Data race on map | Concurrent reads+writes without lock | Use `sync.RWMutex` or `sync.Map` |
| Loop variable captured in goroutine (pre-Go 1.22) | `i` refers to loop variable by reference | Pass `i` as argument: `go func(i int) { ... }(i)` |
| `nil` panic on interface | Concrete `nil` pointer assigned to interface | Return typed nil: `return (*T)(nil), err` is still non-nil interface; return `nil, err` instead |
| HTTP response body not closed | Missing `defer resp.Body.Close()` | Always `defer resp.Body.Close()` immediately after checking err |
| `json.Unmarshal` silently ignores unknown fields | Default behavior | Use `json.Decoder` with `DisallowUnknownFields()` when strict parsing needed |
| Test helpers pollute failure output | Helper not calling `t.Helper()` | Add `t.Helper()` at the top of every test helper function |
| Build constraint ignored | Old `// +build` syntax | Use `//go:build` (Go 1.17+); both forms for backward compat if supporting older Go |

## Module and toolchain quick reference

```bash
go mod init example.com/myproject   # create module
go get github.com/pkg/errors@v0.9.1 # add dependency at specific version
go mod tidy                          # prune unused, add missing
go mod vendor                        # vendor deps for reproducible builds

go test ./...                        # run all tests
go test -race ./...                  # run with race detector (CI must pass this)
go test -bench=. -benchmem ./...     # benchmarks with memory allocation stats
go build -trimpath -ldflags="-s -w" # stripped production binary

go vet ./...                         # static analysis (run before committing)
staticcheck ./...                    # extended linting (install separately)
```

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go specification](https://go.dev/ref/spec)
- [Go standard library](https://pkg.go.dev/std)
- [Go concurrency patterns (Pike, 2012)](https://go.dev/talks/2012/concurrency.slide)
- [Go blog: error handling](https://go.dev/blog/error-handling-and-go)
- [Go generics proposal / tutorial](https://go.dev/doc/tutorial/generics)
- [testify](https://github.com/stretchr/testify)
- [staticcheck](https://staticcheck.dev/)
