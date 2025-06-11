# Test Utilities API Reference

The testutils package (`pkg/testutils`) provides mock implementations and helper functions for testing go-llms applications. These utilities enable deterministic testing without external dependencies.

## Overview

Test utilities include:
- Mock LLM providers for testing generation logic
- Mock tools for testing agent-tool interactions
- Helper functions for pointer creation
- Utilities for schema-based data generation

## Mock Providers

### TestMockProvider

A flexible mock provider for testing multi-provider scenarios.

```go
import "github.com/lexlapax/go-llms/pkg/testutils"

// Create mock with custom behavior
mock := &testutils.TestMockProvider{
    Name: "test-provider",
    GenerateFunc: func(ctx context.Context, prompt string, opts ...domain.GenerateOption) (string, error) {
        if strings.Contains(prompt, "error") {
            return "", fmt.Errorf("simulated error")
        }
        return "Mock response to: " + prompt, nil
    },
}

// Use in tests
response, err := mock.Generate(ctx, "test prompt")
```

### CustomMockProvider

Simplified mock provider focused on message-based generation.

```go
mock := &testutils.CustomMockProvider{
    Name: "custom-mock",
    GenerateMessageFunc: func(ctx context.Context, messages []domain.Message, opts ...domain.GenerateOption) (*domain.Message, error) {
        // Custom logic based on messages
        lastMessage := messages[len(messages)-1]
        return &domain.Message{
            Role:    domain.RoleAssistant,
            Content: "Response to: " + lastMessage.Content,
        }, nil
    },
}
```

### MockStructuredProvider

Specialized for testing structured data generation.

```go
// Create with predefined structured response
expectedData := map[string]interface{}{
    "name": "John Doe",
    "age": 30,
    "active": true,
}

mock := &testutils.MockStructuredProvider{
    Name:           "structured-mock",
    StructuredData: expectedData,
}

// Generate with schema
result, err := mock.GenerateWithSchema(ctx, prompt, schema)
// result will contain expectedData
```

### Advanced MockProvider

The provider package includes a sophisticated mock with schema awareness:

```go
import "github.com/lexlapax/go-llms/pkg/llm/provider"

// Create mock with response mappings
mock := provider.NewMockProvider()
mock.AddResponse("What is 2+2?", "4")
mock.AddResponse("Explain AI", "AI is artificial intelligence...")

// Schema-aware generation
schema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "name":  {Type: "string"},
        "email": {Type: "string", Format: "email"},
        "age":   {Type: "integer", Minimum: float64Ptr(0)},
    },
}

// Generates realistic data based on schema
result, err := mock.GenerateWithSchema(ctx, "Generate user", schema)
// Returns: {"name": "John Doe", "email": "john@example.com", "age": 25}
```

## Mock Tools

### MockTool

Complete implementation of the Tool interface for testing.

```go
// Create custom mock tool
mockTool := &testutils.MockTool{
    NameFunc: func() string { return "test-tool" },
    DescriptionFunc: func() string { return "A test tool" },
    ExecuteFunc: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
        // Custom execution logic
        if input, ok := params["input"].(string); ok {
            return "Processed: " + input, nil
        }
        return nil, fmt.Errorf("invalid input")
    },
    SchemaFunc: func() interface{} {
        return map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "input": map[string]interface{}{"type": "string"},
            },
            "required": []string{"input"},
        }
    },
}

// Use with agent
agent.AddTool(mockTool)
```

### Pre-configured Tools

```go
// Calculator tool for testing
calcTool := testutils.CreateCalculatorTool()

// Execute calculator operations
result, err := calcTool.Execute(ctx, map[string]interface{}{
    "operation": "add",
    "a": 10,
    "b": 20,
})
// result: 30

// Generic mock tool helper
customTool := testutils.CreateMockTool(
    "fetcher",
    "Fetches data",
    schema,
    func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
        url := params["url"].(string)
        return fmt.Sprintf("Fetched data from %s", url), nil
    },
)
```

## Helper Functions

### Pointer Helpers

Simplify test data creation with pointer utilities:

```go
// Instead of verbose pointer creation
var age *int
ageVal := 30
age = &ageVal

// Use helper functions
age := testutils.IntPtr(30)
temperature := testutils.Float64Ptr(0.7)
enabled := testutils.BoolPtr(true)
name := testutils.StringPtr("test")

// Useful for schema definitions
schema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "age": {
            Type:    "integer",
            Minimum: testutils.Float64Ptr(0),
            Maximum: testutils.Float64Ptr(150),
        },
        "name": {
            Type:      "string",
            MinLength: testutils.IntPtr(1),
            MaxLength: testutils.IntPtr(100),
        },
    },
}
```

## Testing Patterns

### Testing Multi-Provider Logic

```go
func TestMultiProviderFallback(t *testing.T) {
    // Create mocks with different behaviors
    primary := &testutils.TestMockProvider{
        Name: "primary",
        GenerateFunc: func(ctx context.Context, prompt string, opts ...domain.GenerateOption) (string, error) {
            return "", fmt.Errorf("primary failed")
        },
    }
    
    fallback := &testutils.TestMockProvider{
        Name: "fallback",
        GenerateFunc: func(ctx context.Context, prompt string, opts ...domain.GenerateOption) (string, error) {
            return "fallback response", nil
        },
    }
    
    // Test fallback behavior
    multi := provider.NewMultiProvider(primary, fallback)
    response, err := multi.Generate(ctx, "test")
    
    assert.NoError(t, err)
    assert.Equal(t, "fallback response", response)
}
```

### Testing Structured Generation

```go
func TestStructuredOutput(t *testing.T) {
    // Define expected structure
    type User struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }
    
    expectedUser := User{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }
    
    // Create mock that returns structured data
    mock := &testutils.MockStructuredProvider{
        Name:           "test",
        StructuredData: expectedUser,
    }
    
    // Test structured generation
    processor := structured.NewProcessor(mock)
    var result User
    err := processor.Process(ctx, "Generate user", &result)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, result)
}
```

### Testing Agent Tool Interactions

```go
func TestAgentWithTools(t *testing.T) {
    // Track tool calls
    var toolCalls []string
    
    searchTool := &testutils.MockTool{
        NameFunc: func() string { return "search" },
        ExecuteFunc: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
            query := params["query"].(string)
            toolCalls = append(toolCalls, "search:"+query)
            return []string{"result1", "result2"}, nil
        },
    }
    
    analyzeTool := &testutils.MockTool{
        NameFunc: func() string { return "analyze" },
        ExecuteFunc: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
            data := params["data"]
            toolCalls = append(toolCalls, "analyze")
            return map[string]interface{}{"summary": "Analysis complete"}, nil
        },
    }
    
    // Create agent with mock provider and tools
    mockProvider := &testutils.CustomMockProvider{
        GenerateMessageFunc: func(ctx context.Context, messages []domain.Message, opts ...domain.GenerateOption) (*domain.Message, error) {
            // Simulate LLM deciding to use tools
            return &domain.Message{
                Role:    domain.RoleAssistant,
                Content: "I'll search for that information.",
                ToolCalls: []domain.ToolCall{
                    {Tool: "search", Arguments: map[string]interface{}{"query": "test"}},
                    {Tool: "analyze", Arguments: map[string]interface{}{"data": "results"}},
                },
            }, nil
        },
    }
    
    agent := core.NewLLMAgent("test-agent", mockProvider)
    agent.AddTool(searchTool)
    agent.AddTool(analyzeTool)
    
    // Run agent
    result, err := agent.Run(ctx, domain.NewState())
    
    assert.NoError(t, err)
    assert.Equal(t, []string{"search:test", "analyze"}, toolCalls)
}
```

### Testing Error Scenarios

```go
func TestErrorHandling(t *testing.T) {
    // Create mock that fails after N calls
    callCount := 0
    mock := &testutils.TestMockProvider{
        GenerateFunc: func(ctx context.Context, prompt string, opts ...domain.GenerateOption) (string, error) {
            callCount++
            if callCount < 3 {
                return "", fmt.Errorf("temporary error")
            }
            return "success after retries", nil
        },
    }
    
    // Test retry logic
    response, err := llmutil.GenerateWithRetry(ctx, mock, "test",
        llmutil.WithMaxRetries(3),
        llmutil.WithRetryDelay(time.Millisecond),
    )
    
    assert.NoError(t, err)
    assert.Equal(t, "success after retries", response)
    assert.Equal(t, 3, callCount)
}
```

### Testing Streaming

```go
func TestStreaming(t *testing.T) {
    chunks := []string{"Hello", " ", "world", "!"}
    chunkIndex := 0
    
    mock := &testutils.TestMockProvider{
        StreamFunc: func(ctx context.Context, prompt string, opts ...domain.GenerateOption) (<-chan string, error) {
            ch := make(chan string)
            go func() {
                defer close(ch)
                for _, chunk := range chunks {
                    select {
                    case ch <- chunk:
                        time.Sleep(10 * time.Millisecond) // Simulate delay
                    case <-ctx.Done():
                        return
                    }
                }
            }()
            return ch, nil
        },
    }
    
    // Test streaming
    stream, err := mock.Stream(ctx, "test")
    assert.NoError(t, err)
    
    var received []string
    for chunk := range stream {
        received = append(received, chunk)
    }
    
    assert.Equal(t, chunks, received)
}
```

## Best Practices

1. **Use appropriate mock level**: Choose between simple function-based mocks and complex schema-aware mocks based on test needs
2. **Test error paths**: Always test error scenarios and edge cases
3. **Track interactions**: Use closures to track tool/provider calls for verification
4. **Isolate components**: Test each component in isolation before integration testing
5. **Deterministic responses**: Ensure mock responses are predictable and repeatable

## Examples

### Complete Test Setup

```go
package mypackage_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/lexlapax/go-llms/pkg/testutils"
    "github.com/lexlapax/go-llms/pkg/agent/core"
    "github.com/lexlapax/go-llms/pkg/llm/domain"
)

func TestCompleteWorkflow(t *testing.T) {
    ctx := context.Background()
    
    // Setup mock provider
    mockProvider := &testutils.CustomMockProvider{
        Name: "test-llm",
        GenerateMessageFunc: func(ctx context.Context, messages []domain.Message, opts ...domain.GenerateOption) (*domain.Message, error) {
            // Implement test-specific logic
            return &domain.Message{
                Role:    domain.RoleAssistant,
                Content: "Test response",
            }, nil
        },
    }
    
    // Setup mock tools
    tools := []domain.Tool{
        testutils.CreateCalculatorTool(),
        testutils.CreateMockTool("fetcher", "Fetches data", nil, 
            func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
                return "fetched data", nil
            },
        ),
    }
    
    // Create and test agent
    agent := core.NewLLMAgent("test-agent", mockProvider)
    for _, tool := range tools {
        agent.AddTool(tool)
    }
    
    // Run test scenario
    state := domain.NewState().Set("task", "Calculate 2+2 and fetch results")
    result, err := agent.Run(ctx, state)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## See Also

- [Testing Guide](../technical/testing.md) - Comprehensive testing strategies
- [Mock Provider](../technical/testing.md#mock-providers) - Detailed mock provider documentation
- [Agent API](agent.md) - Understanding agents for better testing
- [Tool API](tools.md) - Tool interface for creating mocks