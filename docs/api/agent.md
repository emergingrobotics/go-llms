# Agent API Reference

The agent package (`pkg/agent`) provides the core functionality for building LLM-powered agents that can use tools, manage state, and execute complex workflows. This document focuses on core agent concepts, with detailed tool and workflow documentation available separately.

## Overview

Agents in go-llms are autonomous entities that:
- Process inputs using LLM providers
- Execute tools to interact with external systems
- Manage conversation state and context
- Support hooks for monitoring and customization
- Can be composed into complex workflows

## Core Components

### BaseAgent Interface

The foundation of all agents in the system.

```go
type BaseAgent interface {
    // Name returns the agent's identifier
    Name() string
    
    // Description returns the agent's purpose
    Description() string
    
    // Run executes the agent with the given state
    Run(ctx context.Context, state State) (State, error)
    
    // RunAsync executes asynchronously with event streaming
    RunAsync(ctx context.Context, state State) <-chan Event
    
    // Validate checks if the agent is properly configured
    Validate() error
}
```

### LLMAgent

The primary agent implementation that integrates with LLM providers.

```go
import "github.com/lexlapax/go-llms/pkg/agent/core"

// Create an LLM agent
agent := core.NewLLMAgent("assistant", provider)

// Configure the agent
agent.SetSystemPrompt("You are a helpful assistant with access to tools.")
agent.SetMaxIterations(5)  // Limit tool call iterations
agent.AddTool(searchTool)
agent.AddTool(calculatorTool)

// Run the agent
state := domain.NewState().Set("input", "What's the weather in NYC?")
result, err := agent.Run(ctx, state)
```

### State Management

State represents the agent's working memory and context.

```go
// Create state
state := domain.NewState()

// Set values
state.Set("user_input", "Hello")
state.Set("context", contextData)
state.Set("max_tokens", 1000)

// Get values
input, exists := state.Get("user_input")
if !exists {
    // Handle missing value
}

// Clone state for isolation
newState := state.Clone()

// Merge states
state.Merge(otherState)
```

### Event System

Agents emit events during execution for monitoring and debugging.

```go
// Run agent asynchronously to receive events
eventStream := agent.RunAsync(ctx, state)

for event := range eventStream {
    switch e := event.(type) {
    case *domain.ProgressEvent:
        fmt.Printf("Progress: %s - %d/%d\n", e.Message, e.Current, e.Total)
        
    case *domain.ToolCallEvent:
        fmt.Printf("Tool called: %s with args %v\n", e.Tool, e.Arguments)
        
    case *domain.CompletionEvent:
        fmt.Printf("Agent completed: %v\n", e.Result)
        
    case *domain.ErrorEvent:
        fmt.Printf("Error: %v\n", e.Error)
    }
}
```

## Hooks

Hooks provide lifecycle callbacks for agent operations.

### Hook Interface

```go
type Hook interface {
    // Called before LLM generation
    BeforeGenerate(ctx context.Context, state State) error
    
    // Called after LLM generation
    AfterGenerate(ctx context.Context, state State, response *Message, err error) error
    
    // Called before tool execution
    BeforeToolCall(ctx context.Context, tool string, args map[string]interface{}) error
    
    // Called after tool execution
    AfterToolCall(ctx context.Context, tool string, result interface{}, err error) error
}
```

### Built-in Hooks

#### LoggingHook

```go
// Create logging hook
loggingHook := core.NewLoggingHook(slog.Default())

// Configure log level
loggingHook.SetLevel(core.LogLevelDetailed)

// Add to agent
agent.AddHook(loggingHook)
```

Log levels:
- `LogLevelBasic`: Basic operation info
- `LogLevelDetailed`: Include message content
- `LogLevelDebug`: Full details including tool data

#### MetricsHook

```go
// Create metrics hook
metricsHook := core.NewMetricsHook()
agent.AddHook(metricsHook)

// Run agent...

// Get collected metrics
metrics := metricsHook.GetMetrics()
fmt.Printf("Total requests: %d\n", metrics.Requests)
fmt.Printf("Tool calls: %d\n", metrics.ToolCalls)
fmt.Printf("Avg response time: %.2fms\n", metrics.AverageGenTimeMs)
```

### Custom Hooks

```go
type CustomHook struct{}

func (h *CustomHook) BeforeGenerate(ctx context.Context, state domain.State) error {
    // Custom logic before generation
    fmt.Println("Generating response...")
    return nil
}

func (h *CustomHook) AfterGenerate(ctx context.Context, state domain.State, response *domain.Message, err error) error {
    if err != nil {
        // Handle error
        return err
    }
    // Process response
    return nil
}

// Implement other methods...

agent.AddHook(&CustomHook{})
```

## Agent Configuration

### System Prompts

```go
agent.SetSystemPrompt(`You are an expert assistant with the following capabilities:
- Access to web search for current information
- Mathematical calculations
- Data analysis

Always cite your sources when using tools.`)
```

### Tool Management

```go
// Add individual tools
agent.AddTool(webSearchTool)
agent.AddTool(calculatorTool)

// Add multiple tools
agent.AddTools(tool1, tool2, tool3)

// Remove a tool
agent.RemoveTool("tool-name")

// Get registered tools
tools := agent.GetTools()
```

For detailed tool creation and management, see [Tools API Reference](tools.md).

### Iteration Control

```go
// Limit the number of LLM/tool iterations
agent.SetMaxIterations(10)

// Set timeout for agent execution
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
result, err := agent.Run(ctx, state)
```

## Advanced Features

### Sub-Agents

Agents can invoke other agents as tools.

```go
// Create a research sub-agent
researchAgent := core.NewLLMAgent("researcher", provider)
researchAgent.AddTool(webSearchTool)
researchAgent.SetSystemPrompt("You are a research specialist.")

// Wrap as tool for main agent
researchTool := tools.NewAgentTool(researchAgent)

// Add to main agent
mainAgent.AddTool(researchTool)
```

### State Validators

Ensure state meets requirements before and after execution.

```go
// Create validator
validator := &domain.StateValidator{
    RequiredFields: []string{"user_input", "context"},
    FieldValidators: map[string]func(interface{}) error{
        "max_tokens": func(v interface{}) error {
            if tokens, ok := v.(int); ok && tokens > 0 && tokens <= 4000 {
                return nil
            }
            return fmt.Errorf("max_tokens must be between 1 and 4000")
        },
    },
}

agent.SetStateValidator(validator)
```

### Error Handling

```go
// Configure retry behavior
agent.SetRetryConfig(domain.RetryConfig{
    MaxRetries: 3,
    RetryDelay: time.Second,
    RetryableErrors: []string{
        "rate_limit",
        "timeout",
    },
})

// Handle errors in execution
result, err := agent.Run(ctx, state)
if err != nil {
    var agentErr *domain.AgentError
    if errors.As(err, &agentErr) {
        fmt.Printf("Agent error: %s (type: %s)\n", 
            agentErr.Message, agentErr.Type)
        
        // Check if retryable
        if agentErr.Retryable {
            // Implement retry logic
        }
    }
}
```

## Example: Complete Agent

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    
    "github.com/lexlapax/go-llms/pkg/agent/core"
    "github.com/lexlapax/go-llms/pkg/agent/domain"
    "github.com/lexlapax/go-llms/pkg/agent/builtins/tools"
    "github.com/lexlapax/go-llms/pkg/llm/provider"
)

func main() {
    // Create provider
    llmProvider := provider.NewOpenAIProvider("api-key", 
        provider.WithModel("gpt-4"),
        provider.WithTemperature(0.7),
    )
    
    // Create agent
    agent := core.NewLLMAgent("assistant", llmProvider)
    
    // Configure agent
    agent.SetSystemPrompt(`You are a helpful assistant that can:
    1. Search the web for information
    2. Perform calculations
    3. Analyze data
    
    Always explain your reasoning and cite sources.`)
    
    // Add built-in tools
    webSearch, _ := tools.GetTool("web_search")
    calculator, _ := tools.GetTool("calculator")
    
    agent.AddTool(webSearch)
    agent.AddTool(calculator)
    
    // Add hooks
    agent.AddHook(core.NewLoggingHook(slog.Default()))
    agent.AddHook(core.NewMetricsHook())
    
    // Set limits
    agent.SetMaxIterations(5)
    
    // Create initial state
    state := domain.NewState()
    state.Set("input", "What's the population of Tokyo and how does it compare to New York?")
    
    // Run agent
    ctx := context.Background()
    result, err := agent.Run(ctx, state)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Get output
    output, _ := result.Get("output")
    fmt.Printf("Response: %s\n", output)
}
```

## Integration with Workflows

Agents can be composed into workflows for complex tasks. See [Workflow API Reference](workflows.md) for details on:
- Sequential workflows
- Parallel execution
- Conditional branching
- Loop patterns

## Best Practices

1. **State Management**: Keep state minimal and well-structured
2. **Tool Selection**: Only add tools the agent actually needs
3. **Error Handling**: Always handle errors appropriately
4. **Monitoring**: Use hooks for observability in production
5. **Testing**: Use mock providers and tools for unit tests (see [Test Utilities](testutils.md))

## See Also

- [Tools API Reference](tools.md) - Creating and managing tools
- [Workflow API Reference](workflows.md) - Composing agents into workflows
- [Built-in Components](builtins.md) - Pre-built tools and agents
- [LLM API Reference](llm.md) - Provider integration
- [User Guide: Custom Agents](../user-guide/custom-agents.md) - Practical agent examples