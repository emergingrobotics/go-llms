// Package domain defines the core domain models and interfaces for agents.
package domain

// ABOUTME: Defines the RunContext type for dependency injection in agent workflows
// ABOUTME: Provides type-safe context management for tool execution and state sharing

import "context"

// RunContext carries dependencies for a run
type RunContext[D any] struct {
	ctx  context.Context
	deps D
}

// NewRunContext creates a new run context
func NewRunContext[D any](ctx context.Context, deps D) *RunContext[D] {
	return &RunContext[D]{
		ctx:  ctx,
		deps: deps,
	}
}

// Deps returns the dependencies
func (r *RunContext[D]) Deps() D {
	return r.deps
}

// Context returns the context
func (r *RunContext[D]) Context() context.Context {
	return r.ctx
}
