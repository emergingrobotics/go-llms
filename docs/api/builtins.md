# Built-in Components API Reference

The builtins package (`pkg/agent/builtins`) provides a comprehensive set of pre-built tools and agent templates ready for use in your applications. All built-in components are designed to be LLM-friendly with rich metadata, examples, and error guidance.

## Overview

Go-LLMs includes 32 built-in tools across 7 categories:
- **Math**: Mathematical operations and calculations
- **Data**: Data processing and transformation
- **DateTime**: Date and time manipulation
- **File**: File system operations
- **System**: System information and command execution
- **Web**: HTTP requests, web scraping, and API clients
- **Feed**: RSS/Atom feed processing

## Registry System

### Tool Registry

The tool registry provides discovery and management of built-in tools.

```go
import (
    "github.com/lexlapax/go-llms/pkg/agent/builtins/tools"
    // Import tool categories to auto-register
    _ "github.com/lexlapax/go-llms/pkg/agent/builtins/tools/math"
    _ "github.com/lexlapax/go-llms/pkg/agent/builtins/tools/data"
    _ "github.com/lexlapax/go-llms/pkg/agent/builtins/tools/web"
)

// Access the global tool registry
registry := tools.Tools
```

### Registry Methods

```go
// List all registered tools
allTools := tools.Tools.List()

// Get a specific tool
calculator, err := tools.GetTool("calculator")

// Search tools by keyword
searchResults := tools.Tools.Search("json")

// Filter by category
fileTools := tools.Tools.ListByCategory("file")

// Filter by tags
apiTools := tools.Tools.ListByTags("api", "rest")

// Get tool with metadata
info := tools.Tools.GetWithMetadata("web_search")
// Access: info.Tool, info.Metadata, info.Category, info.Tags

// Check permissions/resources
resources := tools.Tools.GetResourceUsage("execute")
// Returns: UsesFilesystem, UsesNetwork, UsesConcurrency, etc.
```

## Built-in Tools by Category

### Math Tools

#### calculator
Comprehensive mathematical calculator supporting basic and advanced operations.

```go
calc, _ := tools.GetTool("calculator")

// Example usage
result, _ := calc.Execute(ctx, map[string]interface{}{
    "operation": "add",
    "a": 10,
    "b": 20,
})
// Result: 30

// Advanced operations
result, _ = calc.Execute(ctx, map[string]interface{}{
    "operation": "sin",
    "a": math.Pi / 2,
})
// Result: 1
```

**Supported operations**: add, subtract, multiply, divide, modulo, power, sqrt, abs, sin, cos, tan, asin, acos, atan, log, log10, ln, ceil, floor, round

### Data Tools

#### csv_process
Process CSV data with operations like read, write, filter, and transform.

```go
csvTool, _ := tools.GetTool("csv_process")

// Read CSV
result, _ := csvTool.Execute(ctx, map[string]interface{}{
    "operation": "read",
    "path": "data.csv",
})

// Filter CSV
result, _ = csvTool.Execute(ctx, map[string]interface{}{
    "operation": "filter",
    "data": csvData,
    "column": "age",
    "operator": ">",
    "value": 18,
})
```

#### json_process
Parse, validate, transform, and query JSON data.

```go
jsonTool, _ := tools.GetTool("json_process")

// Parse JSON
result, _ := jsonTool.Execute(ctx, map[string]interface{}{
    "operation": "parse",
    "content": `{"name": "John", "age": 30}`,
})

// Query with JSONPath
result, _ = jsonTool.Execute(ctx, map[string]interface{}{
    "operation": "query",
    "data": jsonData,
    "path": "$.users[?(@.age > 25)]",
})
```

#### xml_process
Process XML data with parsing, transformation, and validation.

```go
xmlTool, _ := tools.GetTool("xml_process")

// Parse XML
result, _ := xmlTool.Execute(ctx, map[string]interface{}{
    "operation": "parse",
    "content": xmlString,
})

// Transform with XSLT
result, _ = xmlTool.Execute(ctx, map[string]interface{}{
    "operation": "transform",
    "data": xmlData,
    "xslt": xsltTemplate,
})
```

#### data_transform
General data transformation operations.

```go
transform, _ := tools.GetTool("data_transform")

// Sort data
result, _ := transform.Execute(ctx, map[string]interface{}{
    "operation": "sort",
    "data": []int{3, 1, 4, 1, 5},
    "order": "asc",
})

// Map transformation
result, _ = transform.Execute(ctx, map[string]interface{}{
    "operation": "map",
    "data": []string{"hello", "world"},
    "function": "uppercase",
})
```

### DateTime Tools

#### datetime_now
Get current date/time in various formats and timezones.

```go
now, _ := tools.GetTool("datetime_now")

// Current time in timezone
result, _ := now.Execute(ctx, map[string]interface{}{
    "timezone": "America/New_York",
    "format": "RFC3339",
})
```

#### datetime_parse
Parse date/time strings with automatic format detection.

```go
parser, _ := tools.GetTool("datetime_parse")

result, _ := parser.Execute(ctx, map[string]interface{}{
    "datetime": "2024-01-15 14:30:00",
    "format": "auto", // Auto-detect format
})
```

#### datetime_calculate
Perform date/time arithmetic.

```go
calc, _ := tools.GetTool("datetime_calculate")

// Add duration
result, _ := calc.Execute(ctx, map[string]interface{}{
    "operation": "add",
    "datetime": "2024-01-15T10:00:00Z",
    "duration": "P1DT2H", // 1 day, 2 hours
})

// Calculate difference
result, _ = calc.Execute(ctx, map[string]interface{}{
    "operation": "diff",
    "from": "2024-01-15T10:00:00Z",
    "to": "2024-01-16T12:00:00Z",
    "unit": "hours",
})
```

### File Tools

#### read_file
Read files with encoding detection and streaming support.

```go
reader, _ := tools.GetTool("read_file")

result, _ := reader.Execute(ctx, map[string]interface{}{
    "path": "/path/to/file.txt",
    "encoding": "auto", // Auto-detect encoding
})
```

#### write_file
Write content to files with various modes.

```go
writer, _ := tools.GetTool("write_file")

result, _ := writer.Execute(ctx, map[string]interface{}{
    "path": "/path/to/file.txt",
    "content": "Hello, World!",
    "mode": "create", // create, append, overwrite
})
```

#### search_files
Search for files by name, content, or attributes.

```go
search, _ := tools.GetTool("search_files")

result, _ := search.Execute(ctx, map[string]interface{}{
    "path": "/project",
    "pattern": "*.go",
    "content": "func main",
    "recursive": true,
})
```

### System Tools

#### system_info
Get comprehensive system information.

```go
sysInfo, _ := tools.GetTool("system_info")

result, _ := sysInfo.Execute(ctx, map[string]interface{}{})
// Returns: OS, architecture, CPU count, memory, Go version, etc.
```

#### execute
Execute system commands with timeout and streaming.

```go
exec, _ := tools.GetTool("execute")

result, _ := exec.Execute(ctx, map[string]interface{}{
    "command": "ls",
    "args": ["-la", "/tmp"],
    "timeout": 30,
    "stream": true,
})
```

#### env_var
Manage environment variables.

```go
envVar, _ := tools.GetTool("env_var")

// Get variable
result, _ := envVar.Execute(ctx, map[string]interface{}{
    "operation": "get",
    "name": "PATH",
})

// List all variables
result, _ = envVar.Execute(ctx, map[string]interface{}{
    "operation": "list",
})
```

### Web Tools

#### web_fetch
Fetch content from URLs with caching support.

```go
fetch, _ := tools.GetTool("web_fetch")

result, _ := fetch.Execute(ctx, map[string]interface{}{
    "url": "https://api.example.com/data",
    "headers": map[string]string{
        "Authorization": "Bearer token",
    },
    "cache": true,
})
```

#### web_search
Search the web using various search engines.

```go
search, _ := tools.GetTool("web_search")

result, _ := search.Execute(ctx, map[string]interface{}{
    "query": "golang best practices",
    "engine": "duckduckgo",
    "limit": 10,
})
```

#### api_client
Advanced API client with OpenAPI and GraphQL support.

```go
apiClient, _ := tools.GetTool("api_client")

// REST API call
result, _ := apiClient.Execute(ctx, map[string]interface{}{
    "method": "POST",
    "url": "https://api.example.com/users",
    "body": map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
    },
    "auth": map[string]interface{}{
        "type": "bearer",
        "token": "your-token",
    },
})

// GraphQL query
result, _ = apiClient.Execute(ctx, map[string]interface{}{
    "url": "https://api.example.com/graphql",
    "graphql": map[string]interface{}{
        "query": `query GetUser($id: ID!) {
            user(id: $id) { name email }
        }`,
        "variables": map[string]interface{}{
            "id": "123",
        },
    },
})
```

### Feed Tools

#### feed_fetch
Fetch and parse RSS/Atom/JSON feeds.

```go
feedFetch, _ := tools.GetTool("feed_fetch")

result, _ := feedFetch.Execute(ctx, map[string]interface{}{
    "url": "https://example.com/feed.xml",
    "format": "auto", // auto, rss, atom, json
})
```

#### feed_aggregate
Combine multiple feeds into one.

```go
aggregate, _ := tools.GetTool("feed_aggregate")

result, _ := aggregate.Execute(ctx, map[string]interface{}{
    "feeds": []string{
        "https://feed1.com/rss",
        "https://feed2.com/atom",
    },
    "sort": "date",
    "limit": 50,
})
```

## Tool Metadata and Discovery

### Accessing Tool Metadata

```go
// Get detailed information about a tool
info := tools.Tools.GetWithMetadata("web_search")

// Access metadata
fmt.Printf("Category: %s\n", info.Category)
fmt.Printf("Tags: %v\n", info.Tags)
fmt.Printf("Version: %s\n", info.Metadata["version"])
fmt.Printf("Permissions: %v\n", info.Metadata["permissions"])

// Check resource usage
resources := tools.Tools.GetResourceUsage("web_search")
if resources.UsesNetwork {
    fmt.Println("This tool requires network access")
}
```

### MCP Export

Export tools for Model Context Protocol compatibility:

```go
// Export single tool
mcpTool := tools.Tools.ExportToMCP("calculator")

// Export all tools in a category
mcpTools := tools.Tools.ExportCategoryToMCP("web")

// Export all tools
allMCPTools := tools.Tools.ExportAllToMCP()
```

## Using Built-in Tools with Agents

```go
// Create an agent with built-in tools
agent := core.NewLLMAgent("research-agent", provider)

// Add tools by name
webSearch, _ := tools.GetTool("web_search")
webFetch, _ := tools.GetTool("web_fetch")
jsonProcess, _ := tools.GetTool("json_process")

agent.AddTool(webSearch)
agent.AddTool(webFetch)
agent.AddTool(jsonProcess)

// Now the agent can use these tools
result, _ := agent.Run(ctx, domain.NewState().
    Set("task", "Search for golang tutorials and extract key points"))
```

## Resource and Permission Management

Built-in tools track their resource usage and permission requirements:

```go
// Check before using a tool
tool, _ := tools.GetTool("execute")
resources := tools.Tools.GetResourceUsage("execute")

if resources.UsesFilesystem {
    // Ensure filesystem access is allowed
}
if resources.RequiresPermissions {
    perms := tools.Tools.GetPermissions("execute")
    // Check permissions: ["execute_commands"]
}
```

## Best Practices

1. **Import only needed categories** to reduce binary size
2. **Check tool availability** before use in production
3. **Review tool permissions** for security-sensitive applications
4. **Use tool metadata** to provide better error messages to users
5. **Leverage MCP export** for integration with other systems

## Examples

### Building a Data Processing Pipeline

```go
// Get required tools
csvTool, _ := tools.GetTool("csv_process")
transform, _ := tools.GetTool("data_transform")
jsonTool, _ := tools.GetTool("json_process")

// Create workflow
pipeline := workflow.NewSequentialAgent("data-pipeline")
pipeline.AddAgent(tools.NewToolAgent(csvTool))
pipeline.AddAgent(tools.NewToolAgent(transform))
pipeline.AddAgent(tools.NewToolAgent(jsonTool))

// Process data
result, _ := pipeline.Run(ctx, domain.NewState().
    Set("csv_path", "input.csv").
    Set("transform_op", "aggregate").
    Set("output_format", "json"))
```

### Web Research Assistant

```go
// Create research agent with web tools
agent := core.NewLLMAgent("researcher", provider)

// Add web tools
agent.AddTool(tools.GetTool("web_search"))
agent.AddTool(tools.GetTool("web_fetch"))
agent.AddTool(tools.GetTool("web_scrape"))
agent.AddTool(tools.GetTool("data_transform"))

// Research a topic
result, _ := agent.Run(ctx, domain.NewState().
    Set("query", "Latest developments in quantum computing"))
```

## See Also

- [Tools API Reference](tools.md) - Creating custom tools
- [Agent API Reference](agent.md) - Using tools with agents
- [Workflow API Reference](workflows.md) - Composing tools in workflows
- [Built-in Tools Guide](../user-guide/builtin-tools.md) - Detailed usage examples