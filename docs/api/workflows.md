# Workflow API Reference

The workflow package (`pkg/agent/workflow`) provides patterns for composing agents into complex execution flows. It supports sequential, parallel, conditional, and loop-based workflows with comprehensive error handling and state management.

## Overview

Workflows allow you to:
- Chain agents sequentially with state passing
- Execute agents in parallel with result merging
- Create conditional branches based on state
- Loop agent execution with conditions
- Handle errors with retry and recovery strategies
- Monitor execution with events and hooks

## Core Types

### WorkflowStep

Base interface for all workflow steps.

```go
type WorkflowStep interface {
    Name() string
    Execute(ctx context.Context, state domain.State) (domain.State, error)
    Validate() error
}
```

### WorkflowAgent

Interface for workflow agents that extends BaseAgent.

```go
type WorkflowAgent interface {
    domain.BaseAgent
    AddStep(step WorkflowStep) error
    GetSteps() []WorkflowStep
    SetErrorHandler(handler ErrorHandler)
}
```

### WorkflowState

Wraps domain.State with workflow-specific metadata.

```go
type WorkflowState struct {
    State         domain.State
    WorkflowID    string
    CurrentStep   string
    StepResults   map[string]interface{}
    Error         error
    Metadata      map[string]interface{}
}
```

## Workflow Patterns

### Sequential Workflow

Executes agents one after another, passing state between steps.

```go
import "github.com/lexlapax/go-llms/pkg/agent/workflow"

// Create sequential workflow
sequential := workflow.NewSequentialAgent("data-pipeline").
    WithDescription("Processes data through multiple stages").
    WithStopOnError(true).
    WithMaxRetries(3)

// Add steps
sequential.AddAgent(extractAgent)
sequential.AddAgent(transformAgent)
sequential.AddAgent(loadAgent)

// Execute
result, err := sequential.Run(ctx, initialState)
```

#### Sequential Options

```go
sequential := workflow.NewSequentialAgent("name").
    WithStopOnError(true).      // Stop on first error
    WithMaxRetries(3).           // Retry failed steps
    WithRetryDelay(time.Second). // Delay between retries
    WithHooks(beforeHook, afterHook, errorHook)
```

### Parallel Workflow

Executes multiple agents concurrently.

```go
// Create parallel workflow
parallel := workflow.NewParallelAgent("multi-search").
    WithDescription("Search multiple sources simultaneously").
    WithMaxConcurrency(5).
    WithMergeStrategy(workflow.MergeAll).
    WithTimeout(30 * time.Second)

// Add parallel agents
parallel.AddAgent(webSearchAgent)
parallel.AddAgent(databaseAgent)
parallel.AddAgent(cacheAgent)

// Execute - results are merged based on strategy
result, err := parallel.Run(ctx, searchParams)
```

#### Merge Strategies

```go
// Built-in strategies
workflow.MergeAll      // Merge all results into arrays
workflow.MergeFirst    // Use first non-error result
workflow.MergeByKey    // Merge objects by keys

// Custom merge strategy
customMerge := func(results []domain.State) (domain.State, error) {
    merged := domain.NewState()
    // Custom merging logic
    return merged, nil
}

parallel.WithMergeStrategy(workflow.MergeCustom(customMerge))
```

### Conditional Workflow

Executes different branches based on conditions.

```go
// Create conditional workflow
conditional := workflow.NewConditionalAgent("router").
    WithDescription("Routes to different processors based on input type")

// Add branches with conditions
conditional.AddBranch(
    "text-processor",
    func(state domain.State) bool {
        inputType, _ := state.Get("type")
        return inputType == "text"
    },
    textAgent,
)

conditional.AddBranch(
    "image-processor",
    func(state domain.State) bool {
        inputType, _ := state.Get("type")
        return inputType == "image"
    },
    imageAgent,
)

// Add default branch
conditional.SetDefaultBranch(defaultAgent)

// Execute - runs matching branch
result, err := conditional.Run(ctx, inputState)
```

#### Conditional Options

```go
conditional := workflow.NewConditionalAgent("name").
    WithEvaluateAll(true).        // Evaluate all conditions
    WithAllowMultiple(true).      // Execute multiple matching branches
    WithPriorityEvaluation(true)  // Evaluate in order added
```

### Loop Workflow

Executes agents repeatedly based on conditions.

```go
// Create loop workflow
loop := workflow.NewLoopAgent("data-fetcher").
    WithDescription("Fetches paginated data").
    WithMaxIterations(100).
    WithMaxDuration(5 * time.Minute).
    WithCollectResults(true)

// Set loop body
loop.SetLoopBody(fetchPageAgent)

// Set condition
loop.SetCondition(func(state domain.State, iteration int) bool {
    hasMore, _ := state.Get("hasNextPage")
    return hasMore.(bool)
})

// Execute
result, err := loop.Run(ctx, initialState)
```

#### Loop Helper Functions

```go
// While loop - continues while condition is true
whileLoop := workflow.WhileLoop(
    "while-example",
    func(state domain.State) bool {
        count, _ := state.Get("count")
        return count.(int) < 10
    },
    incrementAgent,
)

// Until loop - continues until condition is true
untilLoop := workflow.UntilLoop(
    "until-example",
    func(state domain.State) bool {
        done, _ := state.Get("completed")
        return done.(bool)
    },
    processAgent,
)

// Count loop - fixed iterations
countLoop := workflow.CountLoop(
    "count-example",
    10, // iterations
    processAgent,
)
```

## Error Handling

### Error Handlers

```go
// Create custom error handler
errorHandler := workflow.NewErrorHandler().
    OnError("network_error", workflow.Retry).
    OnError("validation_error", workflow.Skip).
    OnError("critical_error", workflow.Abort).
    WithDefaultAction(workflow.Continue).
    WithMaxRetries(3).
    WithRetryDelay(time.Second)

// Apply to workflow
workflow.SetErrorHandler(errorHandler)
```

### Error Actions

```go
workflow.Retry    // Retry the failed step
workflow.Skip     // Skip and continue
workflow.Abort    // Stop workflow execution
workflow.Continue // Continue with error in state
```

## State Management

### Working with Workflow State

```go
// Access workflow metadata
state, _ := workflow.Run(ctx, input)
workflowState := state.(*workflow.WorkflowState)

// Check step results
stepResult := workflowState.StepResults["transform"]

// Access workflow metadata
workflowID := workflowState.WorkflowID
currentStep := workflowState.CurrentStep

// Check for errors
if workflowState.Error != nil {
    log.Printf("Workflow failed at step %s: %v", 
        workflowState.CurrentStep, 
        workflowState.Error)
}
```

### State Passing Between Steps

```go
// Each step receives state from previous step
transformStep := func(ctx context.Context, state domain.State) (domain.State, error) {
    // Get data from previous step
    data, _ := state.Get("extracted_data")
    
    // Process data
    transformed := transform(data)
    
    // Update state for next step
    state.Set("transformed_data", transformed)
    return state, nil
}
```

## Hooks and Events

### Workflow Hooks

```go
// Before step execution
beforeHook := func(ctx context.Context, stepName string, state domain.State) error {
    log.Printf("Starting step: %s", stepName)
    return nil
}

// After step execution
afterHook := func(ctx context.Context, stepName string, state domain.State, err error) error {
    if err != nil {
        log.Printf("Step %s failed: %v", stepName, err)
    }
    return nil
}

// Apply hooks
workflow.WithHooks(beforeHook, afterHook, nil)
```

### Event Monitoring

```go
// Async execution with events
eventStream := workflow.RunAsync(ctx, input)

for event := range eventStream {
    switch e := event.(type) {
    case *domain.ProgressEvent:
        log.Printf("Progress: %s - %d/%d", e.Message, e.Current, e.Total)
    case *domain.StepCompletedEvent:
        log.Printf("Completed step: %s", e.StepName)
    case *domain.ErrorEvent:
        log.Printf("Error in step %s: %v", e.StepName, e.Error)
    }
}
```

## Composition Patterns

### Nested Workflows

```go
// Workflows can contain other workflows
mainWorkflow := workflow.NewSequentialAgent("main").
    AddAgent(validationWorkflow).  // Another workflow
    AddAgent(processingWorkflow).  // Another workflow
    AddAgent(outputAgent)
```

### Dynamic Workflow Building

```go
// Build workflow based on configuration
func BuildWorkflow(config Config) workflow.WorkflowAgent {
    wf := workflow.NewSequentialAgent("dynamic")
    
    for _, step := range config.Steps {
        switch step.Type {
        case "transform":
            wf.AddAgent(createTransformAgent(step))
        case "validate":
            wf.AddAgent(createValidateAgent(step))
        case "parallel":
            // Add nested parallel workflow
            parallel := workflow.NewParallelAgent(step.Name)
            for _, subStep := range step.SubSteps {
                parallel.AddAgent(createAgent(subStep))
            }
            wf.AddAgent(parallel)
        }
    }
    
    return wf
}
```

## Examples

### Data Processing Pipeline

```go
// Sequential pipeline with error handling
pipeline := workflow.NewSequentialAgent("etl-pipeline").
    WithDescription("Extract, Transform, Load pipeline").
    WithStopOnError(false). // Continue on errors
    WithMaxRetries(3)

// Error handler with specific strategies
errorHandler := workflow.NewErrorHandler().
    OnError("network_error", workflow.Retry).
    OnError("data_error", workflow.Skip).
    WithDefaultAction(workflow.Abort)

pipeline.SetErrorHandler(errorHandler)

// Add ETL steps
pipeline.AddAgent(extractAgent)
pipeline.AddAgent(validateAgent)
pipeline.AddAgent(transformAgent)
pipeline.AddAgent(loadAgent)

// Execute
result, err := pipeline.Run(ctx, sourceConfig)
```

### Multi-Source Aggregation

```go
// Parallel aggregation with custom merge
aggregator := workflow.NewParallelAgent("data-aggregator").
    WithMaxConcurrency(10).
    WithTimeout(30 * time.Second).
    WithMergeStrategy(workflow.MergeCustom(func(results []domain.State) (domain.State, error) {
        merged := domain.NewState()
        allData := []interface{}{}
        
        for _, result := range results {
            if data, exists := result.Get("data"); exists {
                allData = append(allData, data)
            }
        }
        
        merged.Set("aggregated_data", allData)
        merged.Set("source_count", len(results))
        return merged, nil
    }))

// Add data sources
for _, source := range dataSources {
    aggregator.AddAgent(createFetchAgent(source))
}

result, err := aggregator.Run(ctx, queryParams)
```

### Conditional Processing Router

```go
// Route to different processors based on content
router := workflow.NewConditionalAgent("content-router").
    WithEvaluateAll(false). // Stop at first match
    WithAllowMultiple(false)

// Add routing conditions
router.AddBranch("markdown", 
    func(s domain.State) bool {
        content, _ := s.Get("content")
        return strings.HasSuffix(content.(string), ".md")
    },
    markdownProcessor,
)

router.AddBranch("json",
    func(s domain.State) bool {
        content, _ := s.Get("content")
        return strings.HasSuffix(content.(string), ".json")
    },
    jsonProcessor,
)

router.SetDefaultBranch(plainTextProcessor)

result, err := router.Run(ctx, fileState)
```

## Best Practices

1. **State Management**: Keep state minimal and focused on data flow
2. **Error Handling**: Define clear error strategies for each workflow type
3. **Composition**: Build complex workflows from simpler, reusable components
4. **Monitoring**: Use hooks and events for observability
5. **Testing**: Test workflows with mock agents for predictable behavior

## See Also

- [Agent API Reference](agent.md) - Core agent concepts
- [Tools API Reference](tools.md) - Creating tools for workflows
- [Built-in Components](builtins.md) - Pre-built agents and workflows
- [Workflow Examples](../user-guide/workflows.md) - Detailed workflow guides