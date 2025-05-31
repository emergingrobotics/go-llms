# Go-LLMs Project TODOs

## Features
- [ ] Add Model Context Protocol Client support for Agents
- [ ] Add Model Context Protocol Server support for Workflows or Agents

## Testing & Performance
- [ ] Performance profiling and optimization:
  - [ ] Phase 1: Baseline Profiling Infrastructure (Prerequisites)
    - [ ] P1: Create benchmark harness for A/B testing optimizations (REVISIT)
    - [ ] P2: Implement visualization for memory allocation patterns (REVISIT)
    - [ ] P2: Create real-world test scenarios for end-to-end performance (REVISIT)

  - [ ] Phase 2: High-Impact Optimizations (Quick Wins)
    (All P0 and P1 items completed - see TODO-DONE.md)

  - [ ] Phase 3: Advanced Optimizations (After Initial Improvements)
    - [ ] P1: Implement adaptive channel buffer sizing based on usage patterns (REVISIT)
    - [ ] P1: Add pool prewarming for high-throughput scenarios (REVISIT)
    - [ ] P1: Reduce redundant property iterations in schema processing (REVISIT)
    - [ ] P2: Implement more granular locking in cached objects (REVISIT)
    - [ ] P2: Optimize zero-initialization patterns for pooled objects (REVISIT)
    - [ ] P2: Introduce buffer pooling for string builders (REVISIT)

  - [ ] Phase 4: Integration and Validation (Finalization)
    - [ ] P0: Document performance improvements with metrics (REVISIT)
    - [ ] P0: Verify optimizations in high-concurrency scenarios (REVISIT)
    - [ ] P1: Create benchmark comparison charts for before/after (REVISIT)
    - [ ] P1: Implement regression testing to prevent performance degradation (REVISIT)
    - [ ] P2: Add performance acceptance criteria to CI pipeline (REVISIT)

## Architecture & Built-in Components (Immediate - P0)
- [ ] P0: Analyze consistent logging strategy across codebase
  - [ ] Audit current logging approaches (stdlib log, slog, fmt.Printf, etc.)
  - [ ] Define consistent logging strategy (e.g., simple: stdlib log, complex: slog)
  - [ ] Document logging conventions and patterns
  - [ ] Implement consistent logging throughout codebase
  
- [ ] P0: Analyze structure for exposing built-in tools, agents, and workflows
  - [ ] Review current pkg/agent (including workflow subpackage) structure for extensibility
  - [ ] Design pattern for built-in vs user-defined components
  - [ ] Create registry/discovery mechanism for built-in components
  - [ ] Document guidelines for contributing built-in components
  
- [ ] P0: Build useful built-in tools
  - [ ] Research common LLM tool patterns and use cases
  - [ ] Add specific tool tasks to todo.md after research
  - [ ] Review existing tool interface and extend if needed
  - [ ] Implement initial set of built-in tools (list TBD after research)
  - [ ] Add comprehensive examples and documentation
  
- [ ] P0: Build useful built-in agents  
  - [ ] Research common agent patterns (with and without tools)
  - [ ] Add specific agent tasks to todo.md after research
  - [ ] Review existing agent patterns and extend/refactor as needed
  - [ ] Implement initial set of built-in agents (list TBD after research)
  - [ ] Create agent composition patterns
  - [ ] Add comprehensive examples and documentation
  
- [ ] P0: Build useful multi-agent workflows
  - [ ] Research common workflow patterns requiring multiple agents
  - [ ] Add specific workflow tasks to todo.md after research
  - [ ] Review existing workflow patterns in pkg/agent/workflow and extend as needed
  - [ ] Implement workflow coordination mechanisms
  - [ ] Create initial set of built-in workflows (list TBD after research)
  - [ ] Add comprehensive examples and documentation
    
- [ ] Fix identified cross-link issues (path inconsistencies, broken links) (REVISIT)
- [ ] Perform final consistency check across all documentation (REVISIT)
- [ ] API refinement based on usage feedback
- [ ] Final review and preparation for stable release

## Completed Tasks
See TODO-DONE.md for all completed tasks