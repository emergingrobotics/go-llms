# API Reference

This section provides comprehensive API documentation for all go-llms packages, organized by functionality.

## Core Packages

### LLM Integration
- **[LLM API](llm.md)** - Language model provider integration
  - Unified interface for OpenAI, Anthropic, Google Gemini
  - Multi-provider strategies for reliability
  - Streaming and structured generation support

### Data Validation
- **[Schema API](schema.md)** - JSON Schema validation
  - Schema definition and validation
  - Type coercion and custom validators
  - Integration with structured outputs

### Structured Output
- **[Structured API](structured.md)** - Extract structured data from LLMs
  - Prompt enhancement with schemas
  - JSON extraction and validation
  - Type-safe output processing

## Agent Framework

### Core Agent System
- **[Agent API](agent.md)** - Build autonomous agents
  - Agent lifecycle and state management
  - Hook system for monitoring
  - Event-driven architecture

### Tools and Extensions
- **[Tools API](tools.md)** - Create and manage agent tools
  - ToolBuilder pattern for rich metadata
  - Agent-tool bidirectional conversion
  - Performance optimizations

- **[Built-in Tools](builtins.md)** - Pre-built tool library
  - 32 tools across 7 categories
  - MCP compatibility
  - Tool discovery and registry

### Workflows
- **[Workflow API](workflows.md)** - Compose complex agent behaviors
  - Sequential, parallel, conditional, and loop patterns
  - Error handling and recovery
  - State management across steps

## Utilities

### General Utilities
- **[Utils API](utils.md)** - Helper packages
  - Authentication (OAuth2, API keys)
  - High-performance JSON operations
  - LLM convenience functions
  - Metrics and profiling

### Testing
- **[Test Utilities](testutils.md)** - Testing support
  - Mock providers and tools
  - Test data helpers
  - Deterministic testing patterns

## Quick Start

### Basic LLM Usage
```go
import "github.com/lexlapax/go-llms/pkg/llm/provider"

// Create provider
provider := provider.NewOpenAIProvider("api-key", "gpt-4o")

// Generate text
response, err := provider.Generate(ctx, "Hello, world!")
```

### Structured Output
```go
import "github.com/lexlapax/go-llms/pkg/structured/processor"

// Define structure
type Result struct {
    Name  string `json:"name"`
    Score int    `json:"score"`
}

// Generate structured data
var result Result
err := provider.GenerateWithSchema(ctx, prompt, &result)
```

### Agent with Tools
```go
import (
    "github.com/lexlapax/go-llms/pkg/agent/core"
    "github.com/lexlapax/go-llms/pkg/agent/builtins/tools"
)

// Create agent
agent := core.NewLLMAgent("assistant", provider)

// Add built-in tools
calculator, _ := tools.GetTool("calculator")
webSearch, _ := tools.GetTool("web_search")

agent.AddTool(calculator)
agent.AddTool(webSearch)

// Run agent
state := domain.NewState().Set("input", "Search for Python tutorials and calculate 15% of 200")
result, err := agent.Run(ctx, state)
```

## Package Organization

```
pkg/
├── llm/           # LLM providers and interfaces
├── schema/        # JSON Schema validation
├── structured/    # Structured output processing
├── agent/         # Agent framework
│   ├── core/      # Core agent implementation
│   ├── domain/    # Agent interfaces and types
│   ├── tools/     # Tool system
│   ├── builtins/  # Built-in components
│   └── workflow/  # Workflow patterns
├── util/          # Utility packages
│   ├── auth/      # Authentication
│   ├── json/      # JSON operations
│   ├── llmutil/   # LLM helpers
│   ├── metrics/   # Performance metrics
│   └── profiling/ # Profiling tools
└── testutils/     # Testing utilities
```

## Integration Patterns

The packages are designed to work together:

1. **Schema + LLM** → Structured generation
2. **Agent + Tools** → Autonomous task execution
3. **Workflow + Agent** → Complex multi-step processes
4. **Utils + Any** → Enhanced functionality

## Next Steps

- Browse individual API references for detailed documentation
- Check the [User Guide](../user-guide/) for practical examples
- See [Technical Docs](../technical/) for implementation details
- Explore [Examples](../../cmd/examples/) for working code

## Version Information

Current version: v0.3.1

All APIs are subject to change until v1.0. See [CHANGELOG](../../CHANGELOG.md) for version history.