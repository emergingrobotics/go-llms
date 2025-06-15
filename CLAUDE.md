# CLAUDE.md

Project guidance for Claude Code when working with go-llms.

## Project: Go-LLMs v0.3.5 ✅ RELEASED

Unified Go interface for LLM providers with agent tooling.

**Status**: v0.3.5 Released (June 15, 2025)
**Providers**: OpenAI, Anthropic, Google (Gemini, Vertex AI), Ollama, OpenRouter  
**Build Status**: All checks passing ✅

## v0.3.5 Features:
Complete scripting engine integration with:
- Schema Package & JSON validation
- Serializable error system
- Dynamic tool discovery & registration  
- Type conversion registry
- Event system with serialization/filtering
- Workflow templates & script steps
- Provider metadata & capabilities
- Structured output parsers
- Comprehensive testing infrastructure (280+ tests)
- Documentation generation (OpenAPI, Markdown, JSON)

## Commands
```bash
make test fmt lint    # Standard workflow
make test-all         # All tests including integration
make generate         # Generate tool metadata
```

## Development Rules
- No backward compatibility until v1.0
- No logging in pkg/ (library code)  
- Run `make fmt lint` before committing
- Put common test fixtures, mocks, helpers and scenarios in pkgt/testutils for reuse


## Architecture
- `pkg/llm/` - Provider implementations
- `pkg/agent/` - Tools, workflows, discovery
- `pkg/schema/` - JSON validation
- `pkg/structured/` - Structured outputs
- `pkg/errors/` - Serializable error system
- `pkg/testutils/` - Testing infrastructure