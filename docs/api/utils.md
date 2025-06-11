# Utilities API Reference

The util package (`pkg/util`) provides essential utilities that support the core functionality of go-llms, including authentication, JSON handling, LLM-specific helpers, metrics collection, and profiling capabilities.

## Overview

The utilities are organized into focused packages:
- **auth**: HTTP authentication and session management
- **json**: High-performance JSON operations
- **llmutil**: LLM-specific convenience functions
- **metrics**: Performance metrics collection
- **profiling**: CPU and memory profiling

## Authentication (`pkg/util/auth`)

### Core Authentication

The auth package provides a unified interface for various authentication schemes.

```go
import "github.com/lexlapax/go-llms/pkg/util/auth"

// Create auth configuration
config := &auth.AuthConfig{
    Type: "bearer",
    Token: "your-api-token",
}

// Apply to HTTP request
req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
err := auth.ApplyAuth(req, config)
```

### Supported Authentication Types

```go
// API Key authentication
apiKeyAuth := &auth.AuthConfig{
    Type: "apikey",
    APIKey: "your-key",
    APIKeyHeader: "X-API-Key", // or use APIKeyQuery
}

// Basic authentication
basicAuth := &auth.AuthConfig{
    Type: "basic",
    Username: "user",
    Password: "pass",
}

// Bearer token
bearerAuth := &auth.AuthConfig{
    Type: "bearer",
    Token: "jwt-token",
}

// Custom headers
customAuth := &auth.AuthConfig{
    Type: "custom",
    CustomHeaders: map[string]string{
        "X-Custom-Auth": "value",
        "X-Request-ID": "12345",
    },
}
```

### OAuth2 Support

```go
// OAuth2 client credentials flow
oauth2Config := &auth.OAuth2Config{
    ClientID:     "client-id",
    ClientSecret: "client-secret",
    TokenURL:     "https://auth.example.com/token",
    Scopes:       []string{"read", "write"},
}

// Create token manager
tokenManager := auth.NewOAuth2TokenManager(oauth2Config)

// Get token (handles refresh automatically)
token, err := tokenManager.GetToken(context.Background())

// Apply OAuth2 to request
req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("Authorization", "Bearer " + token.AccessToken)
```

### Session Management

```go
// Create session manager
sessionManager := auth.NewSessionManager()

// Store session data
sessionManager.AddCookie("example.com", &http.Cookie{
    Name:  "session",
    Value: "session-id",
})

// Apply cookies to request
sessionManager.ApplyCookies(req)

// Extract cookies from response
sessionManager.ExtractCookies(resp)

// Save/load sessions
data := sessionManager.Export()
sessionManager.Import(data)
```

### Auto-Detection

```go
// Automatically detect auth requirements from URL/state
authConfig := auth.DetectAuthFromState(state, "https://api.github.com")
// Returns appropriate auth config for known services
```

## JSON Operations (`pkg/util/json`)

### High-Performance JSON

The json package uses jsoniter for 2-3x better performance than standard library.

```go
import "github.com/lexlapax/go-llms/pkg/util/json"

// Drop-in replacement for encoding/json
data, err := json.Marshal(obj)
err = json.Unmarshal(data, &result)

// Pretty printing
pretty, err := json.MarshalIndent(obj, "", "  ")

// Streaming
encoder := json.NewEncoder(writer)
err = encoder.Encode(obj)

decoder := json.NewDecoder(reader)
err = decoder.Decode(&result)
```

### Schema-Optimized Operations

```go
// Fast schema marshaling with buffer pooling
jsonStr, err := json.MarshalSchemaFast(schema)

// Pretty-printed schema
jsonStr, err := json.MarshalSchemaIndent(schema, "", "  ")
```

## LLM Utilities (`pkg/util/llmutil`)

### Provider Creation

```go
import "github.com/lexlapax/go-llms/pkg/util/llmutil"

// Create provider from environment
provider, err := llmutil.ProviderFromEnv()

// Create specific provider
provider, err := llmutil.CreateProvider("openai", apiKey, options...)

// With custom configuration
provider, err := llmutil.CreateProvider("anthropic", apiKey,
    llmoptions.WithModel("claude-3-opus"),
    llmoptions.WithTemperature(0.7),
    llmoptions.WithMaxTokens(2000),
)
```

### Batch Operations

```go
// Generate responses for multiple prompts in parallel
prompts := []string{
    "Explain quantum computing",
    "What is machine learning?",
    "Define artificial intelligence",
}

responses, err := llmutil.BatchGenerate(ctx, provider, prompts,
    llmoptions.WithMaxConcurrency(3),
)
```

### Retry Logic

```go
// Generate with automatic retry on failure
response, err := llmutil.GenerateWithRetry(ctx, provider, prompt,
    llmutil.WithMaxRetries(3),
    llmutil.WithRetryDelay(time.Second),
    llmutil.WithExponentialBackoff(true),
)
```

### Structured Output

```go
type Analysis struct {
    Sentiment string   `json:"sentiment"`
    Keywords  []string `json:"keywords"`
    Summary   string   `json:"summary"`
}

// Type-safe generation
var result Analysis
err := llmutil.GenerateTyped(ctx, provider, prompt, &result)

// With validation
err = llmutil.ProcessTypedWithProvider(ctx, provider, 
    "Analyze this text: "+text,
    &result,
    llmutil.WithValidation(customValidator),
)
```

### Provider Pools

```go
// Create provider pool for load balancing
pool := llmutil.NewProviderPool(
    []domain.Provider{provider1, provider2, provider3},
    llmutil.StrategyRoundRobin,
)

// Use pool like a regular provider
response, err := pool.Generate(ctx, prompt, options...)

// Strategies available:
// - StrategyRoundRobin: Distribute evenly
// - StrategyFailover: Use backup on failure  
// - StrategyFastest: Route to fastest provider
```

### Environment Configuration

```go
// Get provider-specific options from environment
options := llmutil.GetProviderOptionsFromEnv("openai")

// Standard environment variables:
// OPENAI_API_KEY, ANTHROPIC_API_KEY, GOOGLE_API_KEY
// OPENAI_MODEL, ANTHROPIC_MODEL, GOOGLE_MODEL
// OPENAI_BASE_URL, OPENAI_TIMEOUT, etc.

// Auto-detect API key
apiKey := llmutil.GetAPIKeyFromEnv("anthropic")
```

### Option Factories

```go
// Performance-optimized options
perfOptions := llmutil.WithPerformanceOptions("openai")

// Reliability-optimized options
reliableOptions := llmutil.WithReliabilityOptions("anthropic")

// Streaming-optimized options
streamOptions := llmutil.WithStreamingOptions("gemini")

// Combine with environment
factory := llmutil.CreateOptionFactoryFromEnv("openai")
options := factory()
```

### Model Information

```go
// Create model info service
service := modelinfo.NewService()

// Fetch available models
inventory, err := service.FetchModelInventory(ctx, "openai", apiKey)

// List models
for _, model := range inventory.Models {
    fmt.Printf("Model: %s, Context: %d\n", model.ID, model.ContextLength)
}

// With caching
service := modelinfo.NewService(
    modelinfo.WithCache(cache.NewFileCache("/tmp/models")),
)
```

## Metrics (`pkg/util/metrics`)

### Basic Metrics

```go
import "github.com/lexlapax/go-llms/pkg/util/metrics"

// Counter - monotonically increasing
requestCount := metrics.NewCounter("api.requests")
requestCount.Inc()
requestCount.Add(5)

// Gauge - can go up or down
activeConnections := metrics.NewGauge("connections.active")
activeConnections.Set(10)
activeConnections.Inc()
activeConnections.Dec()

// Timer - track durations
timer := metrics.NewTimer("operation.duration")
stop := timer.Start()
// ... do work ...
stop() // Records duration
```

### Ratio Tracking

```go
// Track ratios like cache hit rates
cacheRatio := metrics.NewRatioCounter("cache.hit_rate")
cacheRatio.IncrementHits()    // Cache hit
cacheRatio.IncrementMisses()   // Cache miss

hitRate := cacheRatio.Ratio()  // Get hit rate (0.0 - 1.0)
```

### Global Registry

```go
// Register metrics globally
metrics.GetRegistry().Register("custom.metric", myCounter)

// Retrieve registered metrics
metric := metrics.GetRegistry().Get("custom.metric")

// List all metrics
snapshot := metrics.GetRegistry().Snapshot()
for name, value := range snapshot {
    fmt.Printf("%s: %v\n", name, value)
}
```

### Specialized Metrics

```go
// Cache metrics
cacheMetrics := metrics.NewCacheMetrics("schema.cache")
cacheMetrics.RecordHit()
cacheMetrics.RecordMiss()
cacheMetrics.RecordEviction()

// Pool metrics
poolMetrics := metrics.NewPoolMetrics("object.pool")
poolMetrics.RecordAllocation()
poolMetrics.RecordReuse()
poolMetrics.UpdateSize(10)
```

## Profiling (`pkg/util/profiling`)

### Basic Profiling

```go
import "github.com/lexlapax/go-llms/pkg/util/profiling"

// Enable via environment
// export GO_LLMS_ENABLE_PROFILING=true

// Create profiler
profiler := profiling.NewProfiler("/tmp/profiles")

// Profile CPU usage
err := profiler.StartCPUProfile("operation")
// ... do work ...
profiler.StopCPUProfile()

// Profile memory usage
err = profiler.WriteHeapProfile("after_operation")
```

### Operation Profiling

```go
// Profile specific operations
profiling.ProfileOperation("structured_generation", func() error {
    // ... perform operation ...
    return nil
})

// Predefined operation types
profiling.ProfileStructuredOp("extract_data", func() error {
    return processor.Process(data)
})

profiling.ProfileSchemaOp("validate", func() error {
    return validator.Validate(input)
})
```

### Integration Profiling

```go
// Enable component-specific profiling
if profiling.IsEnabled() {
    defer profiling.ProfilePoolOp("get_object", func() error {
        // Pool operation
        return nil
    })()
}
```

## Examples

### Complete Authentication Flow

```go
// OAuth2 with automatic token refresh
oauth2Config := &auth.OAuth2Config{
    ClientID:     os.Getenv("CLIENT_ID"),
    ClientSecret: os.Getenv("CLIENT_SECRET"),
    TokenURL:     "https://oauth.example.com/token",
}

tokenManager := auth.NewOAuth2TokenManager(oauth2Config)
sessionManager := auth.NewSessionManager()

// Make authenticated request
func makeRequest(url string) (*http.Response, error) {
    req, _ := http.NewRequest("GET", url, nil)
    
    // Apply OAuth2 token
    token, _ := tokenManager.GetToken(context.Background())
    req.Header.Set("Authorization", "Bearer " + token.AccessToken)
    
    // Apply session cookies
    sessionManager.ApplyCookies(req)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    
    // Save new cookies
    sessionManager.ExtractCookies(resp)
    
    return resp, nil
}
```

### Performance Monitoring

```go
// Set up metrics for an operation
var (
    requests = metrics.NewCounter("llm.requests")
    errors   = metrics.NewCounter("llm.errors")
    latency  = metrics.NewTimer("llm.latency")
    tokens   = metrics.NewCounter("llm.tokens")
)

// Monitor LLM operation
func generateWithMetrics(provider domain.Provider, prompt string) (string, error) {
    requests.Inc()
    stop := latency.Start()
    defer stop()
    
    response, err := provider.Generate(context.Background(), prompt)
    if err != nil {
        errors.Inc()
        return "", err
    }
    
    tokens.Add(float64(response.TokenCount))
    return response.Content, nil
}

// Report metrics
snapshot := metrics.GetRegistry().Snapshot()
log.Printf("Metrics: %+v", snapshot)
```

### Optimized Batch Processing

```go
// Process multiple items with pooled providers
pool := llmutil.NewProviderPool(providers, llmutil.StrategyFastest)

// Enable profiling
profiler := profiling.NewProfiler("./profiles")
profiler.StartCPUProfile("batch_process")
defer profiler.StopCPUProfile()

// Track metrics
batchTimer := metrics.NewTimer("batch.duration")
stop := batchTimer.Start()
defer stop()

// Process in parallel
results, err := llmutil.BatchGenerate(ctx, pool, prompts,
    llmutil.WithMaxConcurrency(10),
    llmutil.WithTimeout(5*time.Minute),
)

// Log performance
log.Printf("Processed %d items in %v", len(results), batchTimer.Mean())
```

## Best Practices

1. **Authentication**: Store credentials securely, use OAuth2 for long-running processes
2. **JSON**: Use the optimized json package for all JSON operations
3. **Metrics**: Register metrics at startup, use consistent naming
4. **Profiling**: Enable only when debugging performance issues
5. **Provider Pools**: Use for high-throughput applications

## See Also

- [Provider Options](../user-guide/provider-options.md) - Detailed provider configuration
- [Performance Guide](../technical/performance.md) - Optimization techniques
- [Error Handling](../user-guide/error-handling.md) - Handling provider errors