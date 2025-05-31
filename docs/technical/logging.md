# Logging in Go-LLMs

This document describes the logging approach used throughout the go-llms project. It serves as a reference for developers working with or contributing to the library.

## Core Principles

1. **Library code should not log** - Following Go best practices, library code in `pkg/` returns errors rather than logging them
2. **Logging is the caller's responsibility** - Users of the library decide how to handle errors and what to log
3. **Examples demonstrate logging patterns** - Example programs show recommended logging approaches
4. **Performance over convenience** - Logging is avoided in hot paths; debug builds are used when detailed logging is needed
5. **Thread safety** - All logging approaches used are concurrent-safe

## Logging by Component Type

### Core Library Code (`pkg/`)

The core library follows the Go convention of not imposing logging on users:

- **No direct logging**: Functions return errors with context using `fmt.Errorf()`
- **No stdout/stderr output**: The library never prints directly
- **No logging dependencies**: Users aren't forced to use specific logging libraries

**Exception - Agent Hooks**: The library provides an optional `LoggingHook` in `pkg/agent/workflow/hooks.go` that users can add to agents for structured logging using `slog`.

Example:
```go
// Library code returns errors with context
func (p *Provider) Generate(ctx context.Context, prompt string) (string, error) {
    resp, err := p.client.Complete(ctx, prompt)
    if err != nil {
        return "", fmt.Errorf("provider generate failed: %w", err)
    }
    return resp.Content, nil
}
```

### Command Line Tools (`cmd/`)

CLI tools use direct output for user interaction:

- `fmt.Printf/Println` for normal output
- `fmt.Fprintf(os.Stderr, ...)` for error messages
- Exit codes to indicate success (0) or failure (non-zero)

Example:
```go
func main() {
    result, err := processCommand()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Success: %s\n", result)
}
```

### Example Programs (`cmd/examples/`)

Examples use the standard `log` package for simplicity:

- `log.Printf/Println` for informational output
- `log.Fatalf` for fatal errors
- Agent examples may use `slog` when demonstrating the LoggingHook

Example:
```go
func main() {
    log.Println("Starting example...")
    
    provider := llmutil.NewProvider("openai", apiKey, model)
    result, err := provider.Generate(ctx, "Hello, world!")
    if err != nil {
        log.Fatalf("Generation failed: %v", err)
    }
    
    log.Printf("Result: %s", result)
}
```

### Test Code (`*_test.go`)

Tests use the standard testing package logging:

- `t.Logf()` for debug output (only shown with -v flag)
- `t.Errorf()` for test failures
- `t.Fatalf()` for fatal test errors

### Debug Logging

The library includes a debug logging infrastructure that's only compiled when using the `-tags debug` build flag. This ensures zero overhead in production builds.

#### How to Enable Debug Logging

1. **Build with debug flag**:
   ```bash
   go build -tags debug ./...
   make build-debug
   ```

2. **Test with debug flag**:
   ```bash
   go test -tags debug ./...
   make test-debug
   ```

3. **Control which components log**:
   ```bash
   # Enable all debug logging
   GO_LLMS_DEBUG=all make test-debug
   
   # Enable specific components
   GO_LLMS_DEBUG=param_cache,schema make test-debug
   
   # In your application
   export GO_LLMS_DEBUG=param_cache
   ./bin/go-llms-debug
   ```

#### Using Debug Logging in Code

```go
import "github.com/lexlapax/go-llms/pkg/internal/debug"

// Debug logging is only compiled with -tags debug
debug.Printf("param_cache", "Processing field: %s", fieldName)
debug.Println("schema", "Validation started")
```

When built without the debug tag, these calls compile to no-ops with zero runtime overhead.

## Thread Safety and Concurrency

### Guarantees

- **`slog.Logger`**: Designed to be concurrent-safe
- **`log` package**: Thread-safe by default
- **Immutable loggers**: Logger configurations are not modified after creation
- **No shared state**: Each component manages its own logging

### Best Practices for Concurrent Code

1. **Create loggers once**: Initialize during setup, not during runtime
2. **Use structured logging**: Include context fields for correlation
3. **Avoid logging in hot paths**: Use metrics and periodic summaries instead

Example of structured logging in concurrent operations:
```go
logger.Info("processing item",
    slog.String("worker_id", workerID),
    slog.Int("item_id", itemID),
    slog.String("correlation_id", correlationID))
```

### Anti-Patterns to Avoid

```go
// DON'T: Modify global loggers at runtime
var logger = slog.Default()
func ChangeLogLevel(level slog.Level) {
    logger = slog.New(...) // Race condition!
}

// DON'T: Log while holding locks
mu.Lock()
logger.Info("processing...") // Can cause contention
doWork()
mu.Unlock()

// DON'T: Share logging buffers
var buf bytes.Buffer // Shared across goroutines
fmt.Fprintf(&buf, "log: %v", data) // Race condition!
```

## Using the LoggingHook

The library provides a `LoggingHook` for agent workflows that want structured logging:

```go
import (
    "log/slog"
    "github.com/lexlapax/go-llms/pkg/agent/workflow"
)

// Create a structured logger
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// Create an agent with logging
agent := workflow.NewAgent[MyDeps]("my-agent", provider)
agent.AddHook(workflow.NewLoggingHook(logger, slog.LevelInfo))

// The hook will log:
// - Before/after LLM generation
// - Tool calls and results
// - Errors and retries
```

## Error Handling vs Logging

The library follows these patterns for error handling:

1. **Wrap errors with context**: Use `fmt.Errorf("operation failed: %w", err)`
2. **Return errors up the stack**: Let callers decide how to handle them
3. **Provide detailed error messages**: Include relevant context for debugging
4. **Use error types when appropriate**: For errors that callers might want to handle specially

Example:
```go
// Good: Return error with context
if err := validator.Validate(data); err != nil {
    return fmt.Errorf("validation failed for schema %s: %w", schemaID, err)
}

// Bad: Log and return generic error
if err := validator.Validate(data); err != nil {
    log.Printf("Validation error: %v", err) // Don't log in library
    return errors.New("validation failed")   // Lost context
}
```

## Summary

The go-llms logging approach prioritizes:
- **User control**: Libraries don't impose logging decisions
- **Performance**: No logging overhead in production code
- **Safety**: All logging is thread-safe
- **Clarity**: Clear patterns for each component type
- **Flexibility**: Users can integrate with any logging system

This design allows go-llms to integrate seamlessly into any Go application while maintaining high performance and safety standards.