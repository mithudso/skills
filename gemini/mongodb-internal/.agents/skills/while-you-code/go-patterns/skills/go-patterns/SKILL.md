---
name: go-patterns
description: Common Go code patterns used in the Fern codebase. Use when writing new Go code, adding new CLI commands, working with OTEL tracing, configuration access, logging, or template rendering. Includes Cobra command structure, OpenTelemetry spans, config access, structured logging, Go template patterns, and steps for adding new commands.
user-invocable: false
source: 10gen/fern
license: Internal
mongodb:
  team: devprod-bv
  owner: srdjan.pajic@mongodb.com
  internal: true
---

# Go Patterns

## Command Structure (Cobra)

```go
// cmd/services/start.go
var startCmd = &cobra.Command{
    Use:   "start [service]",
    Short: "Start a service and its dependencies",
    RunE: func(cmd *cobra.Command, args []string) error {
        ctx := cmd.Context()
        // Command logic here
        return nil
    },
}

func init() {
    // Register flags
    startCmd.Flags().StringVar(&version, "version", "latest", "Version to run")
}
```

## OpenTelemetry Tracing

```go
// Always create spans for significant operations
ctx, span := stats.StartSpan(ctx, "operation_name")
defer span.End()

// Record errors in spans
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return err
}
```

Tracing is only enabled when `workspace.IsStatsEnabled()` returns true (controlled by workspace config).

## Configuration Access

```go
// Access workspace configuration
cfg := workspace.CtxCfg()

// Check if a feature is enabled (pointer fields distinguish "not set" from zero value)
if cfg.Updates != nil && cfg.Updates.SkipCheck != nil && *cfg.Updates.SkipCheck {
    // Feature is enabled
}
```

- Config validation uses JSON Schema (see `pkg/config/config.go`)
- Custom template functions in `pkg/deploy/functions.go`

## Logging

```go
import "log/slog"

// Use structured logging
slog.Debug("operation completed", "service", serviceName, "duration", duration)
slog.Info("starting service", "name", name, "version", version)
slog.Warn("potential issue", "error", err)
slog.Error("operation failed", "error", err)

```

## Template Rendering

```go
// pkg/deploy/templating.go pattern
tmpl, err := template.New("name").Funcs(funcMap).Parse(templateContent)
if err != nil {
    return err
}

var buf bytes.Buffer
if err := tmpl.Execute(&buf, data); err != nil {
    return err
}
```

## Error Handling

```go
// Always wrap errors with context
return fmt.Errorf("failed to start service %s: %w", name, err)

// Record in span
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return err
}

// Never swallow errors silently
```

## Codebase Gotchas

- **PATH Management**: `cmd/root.go:setupPathEnv()` modifies PATH to include `~/fern/bin` and Bazel runfiles
- **Version Resolution** chain: explicit `--version` flag → config file default → service catalog resolver → fallback to "latest"
- **Dependency Graph**: Service dependencies are automatically started; use `--no-deps` to skip (assumes deps running); circular dependencies are detected
- **Module Imports**: Internal packages use `github.com/10gen/fern/...`; always run `make tidy` after adding dependencies

## Adding a New Command

1. Create command file in appropriate `cmd/` subdirectory
2. Define `cobra.Command` with `Use`, `Short`, `Long`, `RunE` (see Cobra pattern above)
3. Register command in parent command's `init()` function
4. Run `make cli-docs` to generate documentation
5. Add tests in `*_test.go` file
6. Run `make verify-commit` before committing

**Architecture rules**: Only `cmd/` files may import `cobra`. Only `cmd/flags/` and `cmd/root.go` may import `pflag`. Business logic belongs in `pkg/`, not in command files. See the `golangci-lint` skill for depguard details.
