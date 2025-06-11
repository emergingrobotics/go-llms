# Structured Output API Reference

The structured package (`pkg/structured`) enables reliable extraction of structured data from LLM outputs. It provides prompt enhancement for schema-guided generation and robust processing of raw LLM responses.

## Overview

The structured package provides:
- Schema-guided prompt enhancement
- JSON extraction from various output formats
- Validation against schemas
- Direct mapping to Go types
- Caching for performance
- Integration with LLM providers

## Core Interfaces

### Processor

Validates and processes LLM outputs against schemas.

```go
type Processor interface {
    // Process validates output against schema, returns generic result
    Process(schema *Schema, output string) (interface{}, error)

    // ProcessTyped validates and maps output to a specific Go type
    ProcessTyped(schema *Schema, output string, target interface{}) error

    // ToJSON converts an object to JSON string
    ToJSON(obj interface{}) (string, error)
}
```

### PromptEnhancer

Enhances prompts with schema information for better LLM guidance.

```go
type PromptEnhancer interface {
    // Enhance adds schema information to a prompt
    Enhance(prompt string, schema *Schema) (string, error)

    // EnhanceWithOptions adds schema with additional guidance
    EnhanceWithOptions(prompt string, schema *Schema, options map[string]interface{}) (string, error)
}
```

## JSON Processor

The main implementation for processing structured outputs.

```go
import "github.com/lexlapax/go-llms/pkg/structured/processor"

// Create processor
proc := processor.NewJsonProcessor()

// With custom validator
validator := validation.NewValidator(validation.WithCoercion(true))
proc := processor.NewJsonProcessor(processor.WithValidator(validator))
```

### Processing Output

```go
// Define schema
schema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "name":    {Type: "string"},
        "age":     {Type: "integer", Minimum: float64Ptr(0)},
        "email":   {Type: "string", Format: "email"},
        "active":  {Type: "boolean"},
    },
    Required: []string{"name", "email"},
}

// Process raw LLM output
rawOutput := `Here's the user data:
{
    "name": "John Doe",
    "age": 30,
    "email": "john@example.com",
    "active": true
}`

result, err := proc.Process(schema, rawOutput)
if err != nil {
    // Handle validation or extraction error
    log.Fatal(err)
}

// Use generic result
data := result.(map[string]interface{})
fmt.Printf("Name: %s\n", data["name"])
```

### Type-Safe Processing

```go
// Define target struct
type User struct {
    Name   string `json:"name"`
    Age    int    `json:"age"`
    Email  string `json:"email"`
    Active bool   `json:"active"`
}

// Process directly into struct
var user User
err := proc.ProcessTyped(schema, rawOutput, &user)
if err != nil {
    log.Fatal(err)
}

// Use typed result
fmt.Printf("User: %+v\n", user)
```

## Prompt Enhancement

Guide LLMs to generate valid structured output.

### Basic Enhancement

```go
// Create enhancer
enhancer := processor.NewPromptEnhancer()

// Original prompt
prompt := "List the top 3 programming languages with their key features"

// Schema for expected output
schema := &domain.Schema{
    Type: "array",
    Items: &domain.Schema{
        Type: "object",
        Properties: map[string]*domain.Schema{
            "language": {Type: "string"},
            "features": {
                Type:     "array",
                Items:    &domain.Schema{Type: "string"},
                MinItems: intPtr(2),
            },
            "popularity": {
                Type:    "number",
                Minimum: float64Ptr(0),
                Maximum: float64Ptr(100),
            },
        },
        Required: []string{"language", "features"},
    },
    MinItems: intPtr(3),
    MaxItems: intPtr(3),
}

// Enhance prompt
enhanced, err := enhancer.Enhance(prompt, schema)

// The enhanced prompt includes:
// - Original prompt
// - JSON schema definition
// - Clear formatting instructions
// - Field descriptions and requirements
```

### Enhancement with Options

```go
options := map[string]interface{}{
    // Additional instructions
    "instructions": "Focus on modern languages used in 2024",
    
    // Output format hint
    "format": "a JSON array with exactly 3 items",
    
    // Examples to guide the model
    "examples": []interface{}{
        map[string]interface{}{
            "language": "Python",
            "features": []string{"Simple syntax", "Rich ecosystem"},
            "popularity": 85.5,
        },
    },
    
    // Custom guidance
    "guidance": "Popularity should be a percentage from 0-100",
}

enhanced, err := enhancer.EnhanceWithOptions(prompt, schema, options)
```

## JSON Extraction

The processor handles various output formats automatically:

```go
extractor := processor.NewJsonExtractor()

// Extracts JSON from markdown code blocks
output1 := "Here's the data:\n```json\n{\"key\": \"value\"}\n```"
json, _ := extractor.Extract(output1)

// Extracts from mixed text
output2 := "The result is {\"status\": \"success\", \"count\": 42}"
json, _ = extractor.Extract(output2)

// Handles arrays
output3 := "Results: [1, 2, 3, 4, 5]"
json, _ = extractor.Extract(output3)
```

## Complete Workflow

### With LLM Provider

```go
// 1. Define your data structure
type QueryResult struct {
    Query   string   `json:"query"`
    Results []string `json:"results"`
    Count   int      `json:"count"`
    Source  string   `json:"source"`
}

// 2. Create schema
schema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "query":   {Type: "string"},
        "results": {
            Type:  "array",
            Items: &domain.Schema{Type: "string"},
        },
        "count":  {Type: "integer", Minimum: float64Ptr(0)},
        "source": {Type: "string"},
    },
    Required: []string{"query", "results", "count"},
}

// 3. Enhance prompt
prompt := "Search for 'golang tutorials' and return the top 5 results"
enhanced, _ := processor.EnhancePromptWithSchema(prompt, schema)

// 4. Generate with LLM
provider := llmprovider.NewOpenAIProvider(apiKey, "gpt-4o")
rawOutput, _ := provider.Generate(ctx, enhanced)

// 5. Process output
proc := processor.NewJsonProcessor()
var result QueryResult
err := proc.ProcessTyped(schema, rawOutput, &result)

// 6. Use the structured data
fmt.Printf("Found %d results for '%s'\n", result.Count, result.Query)
for i, r := range result.Results {
    fmt.Printf("%d. %s\n", i+1, r)
}
```

### With Provider's GenerateWithSchema

```go
// Direct structured generation (provider handles enhancement)
provider := llmprovider.NewOpenAIProvider(apiKey, "gpt-4o")

var result QueryResult
err := provider.GenerateWithSchema(ctx, prompt, &result)
if err != nil {
    log.Fatal(err)
}

// Result is already validated and typed
fmt.Printf("Results: %+v\n", result)
```

## Caching

Schema JSON is cached for performance:

```go
// Get the schema cache
cache := processor.GetSchemaCache()

// Check cache statistics
stats := cache.Stats()
fmt.Printf("Cache hits: %d, misses: %d\n", stats.Hits, stats.Misses)

// Clear cache if needed
cache.Clear()

// Cache is used automatically by PromptEnhancer
// No manual management required in normal usage
```

## Error Handling

```go
result, err := proc.Process(schema, output)
if err != nil {
    switch {
    case errors.Is(err, processor.ErrNoJSONFound):
        // No JSON found in output
        fmt.Println("LLM did not return JSON")
        
    case errors.Is(err, processor.ErrInvalidJSON):
        // JSON parsing failed
        fmt.Println("Invalid JSON format")
        
    case errors.Is(err, processor.ErrValidationFailed):
        // Schema validation failed
        var valErr *processor.ValidationError
        if errors.As(err, &valErr) {
            fmt.Printf("Validation errors: %v\n", valErr.Errors)
        }
        
    default:
        // Other error
        fmt.Printf("Processing failed: %v\n", err)
    }
}
```

## Best Practices

1. **Schema Design**: Start with required fields, add optional ones as needed
2. **Prompt Engineering**: Test enhanced prompts with your target model
3. **Error Recovery**: Have fallback strategies for extraction failures
4. **Validation**: Use type coercion for flexibility with user inputs
5. **Performance**: Reuse processors and enhancers across requests

## Advanced Usage

### Custom Validators

```go
// Create processor with custom validator
validator := validation.NewValidator(
    validation.WithCoercion(true),
    validation.WithCustomValidation(true),
)

// Register custom validator
validation.RegisterCustomValidator("sentiment", func(v interface{}) error {
    str, ok := v.(string)
    if !ok {
        return fmt.Errorf("expected string")
    }
    
    validSentiments := []string{"positive", "negative", "neutral"}
    for _, valid := range validSentiments {
        if str == valid {
            return nil
        }
    }
    return fmt.Errorf("invalid sentiment: %s", str)
})

proc := processor.NewJsonProcessor(processor.WithValidator(validator))
```

### Complex Schema Example

```go
// Schema for a code review response
reviewSchema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "summary": {
            Type:        "string",
            MinLength:   intPtr(10),
            MaxLength:   intPtr(200),
            Description: "Brief summary of the code review",
        },
        "issues": {
            Type: "array",
            Items: &domain.Schema{
                Type: "object",
                Properties: map[string]*domain.Schema{
                    "severity": {
                        Type: "string",
                        Enum: []interface{}{"low", "medium", "high", "critical"},
                    },
                    "line":        {Type: "integer", Minimum: float64Ptr(1)},
                    "description": {Type: "string"},
                    "suggestion":  {Type: "string"},
                },
                Required: []string{"severity", "description"},
            },
        },
        "score": {
            Type:    "number",
            Minimum: float64Ptr(0),
            Maximum: float64Ptr(10),
        },
        "approved": {Type: "boolean"},
    },
    Required: []string{"summary", "issues", "score", "approved"},
}

// Use with complex prompt
prompt := "Review this Go code and provide detailed feedback:\n" + codeSnippet
enhanced, _ := enhancer.EnhanceWithOptions(prompt, reviewSchema, map[string]interface{}{
    "instructions": "Be thorough but constructive in your feedback",
})
```

## Integration

The structured package integrates with:

- **Schema Package**: For validation (see [Schema API](schema.md))
- **LLM Providers**: For generation (see [LLM API](llm.md))
- **Agents**: For tool outputs (see [Agent API](agent.md))
- **Utilities**: Helper functions (see [Utils API](utils.md#llm-utilities))

## See Also

- [Schema API Reference](schema.md) - Schema definition and validation
- [LLM API Reference](llm.md) - Provider integration
- [User Guide: Structured Output](../user-guide/structured-output.md) - Practical examples
- [JSON Schema](https://json-schema.org/) - Schema specification