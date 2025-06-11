# Schema API Reference

The schema package (`pkg/schema`) provides JSON Schema-compatible validation for structured data. It's the foundation for validating LLM outputs and ensuring data conforms to expected formats.

## Overview

The schema package provides:
- JSON Schema compatible validation
- Type coercion for flexible input handling
- Custom validation functions
- Conditional validation logic
- Schema generation from Go structs
- Integration with structured output processing

## Core Types

### Schema

The main schema definition supporting JSON Schema features.

```go
type Schema struct {
    // Basic schema properties
    Type                 string              `json:"type"`
    Description          string              `json:"description,omitempty"`
    Title                string              `json:"title,omitempty"`
    Default              interface{}         `json:"default,omitempty"`
    Examples             []interface{}       `json:"examples,omitempty"`
    
    // Object properties
    Properties           map[string]*Schema  `json:"properties,omitempty"`
    Required             []string            `json:"required,omitempty"`
    AdditionalProperties interface{}         `json:"additionalProperties,omitempty"`
    
    // Array properties
    Items                *Schema             `json:"items,omitempty"`
    MinItems             *int                `json:"minItems,omitempty"`
    MaxItems             *int                `json:"maxItems,omitempty"`
    UniqueItems          *bool               `json:"uniqueItems,omitempty"`
    
    // String constraints
    MinLength            *int                `json:"minLength,omitempty"`
    MaxLength            *int                `json:"maxLength,omitempty"`
    Pattern              string              `json:"pattern,omitempty"`
    Format               string              `json:"format,omitempty"`
    
    // Numeric constraints
    Minimum              *float64            `json:"minimum,omitempty"`
    Maximum              *float64            `json:"maximum,omitempty"`
    ExclusiveMinimum     *float64            `json:"exclusiveMinimum,omitempty"`
    ExclusiveMaximum     *float64            `json:"exclusiveMaximum,omitempty"`
    MultipleOf           *float64            `json:"multipleOf,omitempty"`
    
    // Enumeration and constants
    Enum                 []interface{}       `json:"enum,omitempty"`
    Const                interface{}         `json:"const,omitempty"`
    
    // Conditional validation
    If                   *Schema             `json:"if,omitempty"`
    Then                 *Schema             `json:"then,omitempty"`
    Else                 *Schema             `json:"else,omitempty"`
    AllOf                []*Schema           `json:"allOf,omitempty"`
    AnyOf                []*Schema           `json:"anyOf,omitempty"`
    OneOf                []*Schema           `json:"oneOf,omitempty"`
    Not                  *Schema             `json:"not,omitempty"`
}
```

### Validator

Interface for validating data against schemas.

```go
type Validator interface {
    // Validate JSON string against schema
    Validate(schema *Schema, data string) (*ValidationResult, error)
    
    // Validate Go struct against schema
    ValidateStruct(schema *Schema, obj interface{}) (*ValidationResult, error)
}

// Validation result
type ValidationResult struct {
    Valid  bool     `json:"valid"`
    Errors []string `json:"errors,omitempty"`
}
```

## Creating Schemas

### Basic Schema Definition

```go
import "github.com/lexlapax/go-llms/pkg/schema/domain"

// Simple object schema
userSchema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "name": {
            Type:        "string",
            MinLength:   intPtr(1),
            MaxLength:   intPtr(100),
            Description: "User's full name",
        },
        "email": {
            Type:   "string",
            Format: "email",
        },
        "age": {
            Type:    "integer",
            Minimum: float64Ptr(0),
            Maximum: float64Ptr(150),
        },
        "active": {
            Type:    "boolean",
            Default: true,
        },
    },
    Required: []string{"name", "email"},
}
```

### Array Schema

```go
// Array of strings with constraints
tagsSchema := &domain.Schema{
    Type:        "array",
    MinItems:    intPtr(1),
    MaxItems:    intPtr(10),
    UniqueItems: boolPtr(true),
    Items: &domain.Schema{
        Type:      "string",
        MinLength: intPtr(2),
        MaxLength: intPtr(20),
        Pattern:   "^[a-z0-9-]+$",
    },
}
```

### Nested Objects

```go
// Complex nested structure
orderSchema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "orderId": {
            Type:   "string",
            Format: "uuid",
        },
        "customer": {
            Type: "object",
            Properties: map[string]*domain.Schema{
                "name":  {Type: "string"},
                "email": {Type: "string", Format: "email"},
            },
            Required: []string{"name", "email"},
        },
        "items": {
            Type:     "array",
            MinItems: intPtr(1),
            Items: &domain.Schema{
                Type: "object",
                Properties: map[string]*domain.Schema{
                    "productId": {Type: "string"},
                    "quantity":  {Type: "integer", Minimum: float64Ptr(1)},
                    "price":     {Type: "number", Minimum: float64Ptr(0)},
                },
                Required: []string{"productId", "quantity", "price"},
            },
        },
    },
    Required: []string{"orderId", "customer", "items"},
}
```

## Validation

### Creating a Validator

```go
import "github.com/lexlapax/go-llms/pkg/schema/validation"

// Basic validator
validator := validation.NewValidator()

// With type coercion enabled
validator := validation.NewValidator(
    validation.WithCoercion(true),
)

// With custom validation functions
validator := validation.NewValidator(
    validation.WithCustomValidation(true),
)
```

### Validating Data

```go
// Validate JSON string
jsonData := `{"name": "John Doe", "email": "john@example.com", "age": 30}`
result, err := validator.Validate(userSchema, jsonData)
if err != nil {
    // Handle parsing error
    log.Fatal(err)
}

if !result.Valid {
    fmt.Println("Validation errors:")
    for _, err := range result.Errors {
        fmt.Printf("- %s\n", err)
    }
}

// Validate Go struct
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

user := User{Name: "Jane", Email: "jane@example.com", Age: 25}
result, err = validator.ValidateStruct(userSchema, user)
```

## String Formats

Supported string format validators:

```go
// Standard formats
schema := &domain.Schema{
    Type:   "string",
    Format: "email",     // Email address
}

// Available formats:
// - email: RFC 5322 email address
// - date-time: RFC 3339 date-time
// - date: Full date (YYYY-MM-DD)
// - time: Time with timezone
// - uri: Valid URI
// - uri-reference: URI or relative reference
// - uuid: UUID (any version)
// - hostname: Internet hostname
// - ipv4: IPv4 address
// - ipv6: IPv6 address
// - regex: Valid regular expression
// - json-pointer: JSON Pointer
// - relative-json-pointer: Relative JSON Pointer

// Multiple formats (OR logic)
schema := &domain.Schema{
    Type:   "string",
    Format: "email|phone", // Either email OR phone number
}
```

## Type Coercion

When enabled, the validator attempts to convert values to expected types:

```go
validator := validation.NewValidator(validation.WithCoercion(true))

schema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "count": {Type: "integer"},
        "price": {Type: "number"},
        "active": {Type: "boolean"},
    },
}

// These string values will be coerced
jsonData := `{
    "count": "42",      // Coerced to integer 42
    "price": "19.99",   // Coerced to float 19.99
    "active": "true"    // Coerced to boolean true
}`

result, _ := validator.Validate(schema, jsonData)
// Validation succeeds with coercion
```

## Conditional Validation

### If-Then-Else

```go
// Different validation based on account type
schema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "accountType": {
            Type: "string",
            Enum: []interface{}{"personal", "business"},
        },
    },
    If: &domain.Schema{
        Properties: map[string]*domain.Schema{
            "accountType": {Const: "business"},
        },
    },
    Then: &domain.Schema{
        Properties: map[string]*domain.Schema{
            "companyName": {Type: "string"},
            "taxId":       {Type: "string"},
        },
        Required: []string{"companyName", "taxId"},
    },
    Else: &domain.Schema{
        Properties: map[string]*domain.Schema{
            "firstName": {Type: "string"},
            "lastName":  {Type: "string"},
        },
        Required: []string{"firstName", "lastName"},
    },
}
```

### Logical Operators

```go
// AllOf - must match all schemas
schema := &domain.Schema{
    AllOf: []*domain.Schema{
        {Properties: map[string]*domain.Schema{"a": {Type: "string"}}},
        {Properties: map[string]*domain.Schema{"b": {Type: "number"}}},
    },
}

// AnyOf - must match at least one schema
schema := &domain.Schema{
    AnyOf: []*domain.Schema{
        {Type: "string"},
        {Type: "number"},
    },
}

// OneOf - must match exactly one schema
schema := &domain.Schema{
    OneOf: []*domain.Schema{
        {Type: "string", MinLength: intPtr(5)},
        {Type: "number", Minimum: float64Ptr(0)},
    },
}

// Not - must not match the schema
schema := &domain.Schema{
    Not: &domain.Schema{
        Type: "string",
        Enum: []interface{}{"forbidden", "banned"},
    },
}
```

## Schema Generation

Generate schemas from Go structs:

```go
import "github.com/lexlapax/go-llms/pkg/schema/adapter/reflection"

type Product struct {
    ID          string   `json:"id" schema:"required"`
    Name        string   `json:"name" schema:"required,minLength=1,maxLength=100"`
    Price       float64  `json:"price" schema:"minimum=0"`
    Tags        []string `json:"tags" schema:"minItems=1,uniqueItems=true"`
    InStock     bool     `json:"in_stock"`
    Description *string  `json:"description,omitempty"`
}

// Generate schema from struct
schema := reflection.GenerateSchema(Product{})

// The generated schema will include all constraints from tags
```

## Custom Validators

Register custom validation functions:

```go
// Register a custom validator
validation.RegisterCustomValidator("phone", func(value interface{}) error {
    str, ok := value.(string)
    if !ok {
        return fmt.Errorf("expected string")
    }
    
    // Simple phone validation
    matched, _ := regexp.MatchString(`^\+?[\d\s-()]+$`, str)
    if !matched {
        return fmt.Errorf("invalid phone number")
    }
    return nil
})

// Use in schema
schema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "phone": {
            Type:            "string",
            CustomValidator: "phone",
        },
    },
}
```

## Integration with Other Components

The schema package integrates with:

- **Structured Output**: Validates LLM responses (see [Structured API](structured.md))
- **Agents**: Validates tool parameters and outputs (see [Tools API](tools.md))
- **LLM Providers**: Ensures structured generation (see [LLM API](llm.md))

## Best Practices

1. **Schema Design**: Start simple and add constraints incrementally
2. **Validation Messages**: Use descriptions for better error messages
3. **Type Coercion**: Enable for user inputs, disable for system data
4. **Performance**: Reuse validators and schemas when possible
5. **Testing**: Always test schemas with valid and invalid data

## Examples

### Form Validation

```go
// User registration form schema
registrationSchema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "username": {
            Type:        "string",
            Pattern:     "^[a-zA-Z0-9_]{3,20}$",
            Description: "Alphanumeric username (3-20 chars)",
        },
        "email": {
            Type:   "string",
            Format: "email",
        },
        "password": {
            Type:      "string",
            MinLength: intPtr(8),
            Pattern:   "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).*$",
            Description: "At least 8 chars with uppercase, lowercase, and number",
        },
        "age": {
            Type:    "integer",
            Minimum: float64Ptr(13),
        },
        "terms": {
            Type:  "boolean",
            Const: true,
            Description: "Must accept terms",
        },
    },
    Required: []string{"username", "email", "password", "age", "terms"},
}
```

### API Response Validation

```go
// API response schema with optional fields
apiResponseSchema := &domain.Schema{
    Type: "object",
    Properties: map[string]*domain.Schema{
        "status": {
            Type: "string",
            Enum: []interface{}{"success", "error"},
        },
        "data": {
            Type: "object",
            AdditionalProperties: true, // Allow any properties
        },
        "error": {
            Type: "object",
            Properties: map[string]*domain.Schema{
                "code":    {Type: "string"},
                "message": {Type: "string"},
            },
        },
        "metadata": {
            Type: "object",
            Properties: map[string]*domain.Schema{
                "timestamp": {Type: "string", Format: "date-time"},
                "version":   {Type: "string"},
            },
        },
    },
    Required: []string{"status"},
    If: &domain.Schema{
        Properties: map[string]*domain.Schema{
            "status": {Const: "error"},
        },
    },
    Then: &domain.Schema{
        Required: []string{"error"},
    },
    Else: &domain.Schema{
        Required: []string{"data"},
    },
}
```

## See Also

- [Structured API Reference](structured.md) - Using schemas with LLM outputs
- [Advanced Validation Guide](../user-guide/advanced-validation.md) - Complex validation patterns
- [JSON Schema Specification](https://json-schema.org/) - Official JSON Schema docs
- [Test Utilities](testutils.md) - Testing with schemas