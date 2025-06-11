# Go-LLMs Project TODOs - v0.3.x Roadmap

## v0.3.1 Release (Ready)
- [ ] 0.3.1.1: Tag release v0.3.1

## v0.3.2 Documentation simplification and refactoring
### 0.3.2.1 API Documentation (docs/api) (look through code that's been implemented) ✅ COMPLETED
- [x] Create tools.md - Tools API Documentation for pkg/agent/tools
- [x] Create workflows.md - Workflow API Documentation for pkg/agent/workflow
- [x] Create builtins.md - Built-ins API Documentation for pkg/agent/builtins
- [x] Create utils.md - Utilities API Documentation for pkg/util/*
- [x] Create testutils.md - Test Utilities API Documentation for pkg/testutils/*
- [x] Update agent.md - Focus on core agent concepts with cross-references to tools/workflows/builtins
- [x] Update llm.md - Ensure consistency with new structure
- [x] Update schema.md - Ensure consistency with new structure
- [x] Update structured.md - Ensure consistency with new structure
- [x] Update docs/api/README.md - Update index to reflect new modular structure
### 0.3.2.2 Restructure of User Guide Documentation (docs/user-guide) (look through code that's been implemented) ✅ COMPLETED
#### Core Getting Started Flow
- [x] Create/Update getting-started.md - Installation and first program
- [x] Create core-concepts.md - Essential concepts (providers, messages, options)
- [x] Create providers.md - Working with different LLM providers (consolidate from existing)
- [x] Update structured-output.md - Keep focused on extracting structured data
#### Advanced Features
- [x] Create agents.md - Building and using agents (user perspective)
- [x] Create tools.md - Merge builtin-tools.md and built-in-components.md
- [x] Create workflows.md - Composing agent workflows (user perspective)
- [x] Update multimodal-content.md - Keep focused on working with images/content
#### Consolidation and Cleanup
- [x] Merge web-search-tool.md content into tools.md
- [x] Update examples-gallery.md - Make it a quick reference/index
- [x] Delete redundant files after merging
- [x] Move old/redundant files to docs/archives
- [x] Update docs/user-guide/README.md - New structure and learning path
### 0.3.2.3 Restructure of Technical Documentation (docs/technical) (look through code that's been implemented)
#### Core Architecture Documentation
- [ ] Update architecture.md - System design, components, and data flow
- [ ] Create provider-implementation.md - How to add new LLM providers
- [ ] Update tool-development.md - Internal tool architecture and patterns
- [ ] Update performance.md - Performance considerations and optimizations
#### Implementation Details
- [ ] Update concurrency.md - Concurrency patterns used in the library
- [ ] Update caching.md - Caching strategies and implementation
- [ ] Update testing.md - Testing approach and guidelines
- [ ] Update authentication.md - Auth system architecture
#### Cleanup and Organization
- [ ] Remove duplicate content (keep technical details here, user guides in user-guide/)
- [ ] Move multimodal-content.md user aspects to user-guide, keep technical here
- [ ] Move built-in-components.md user aspects to user-guide, keep technical here
- [ ] Delete redundant files after content migration
- [ ] Update docs/technical/README.md - New structure for contributors
### 0.3.2.4 Restructure of archives (docs/archives)
- [ ] move docs/plan documentation to docs/archives, remove docs/plan
- [ ] Review and categorize all files in archives directory
- [ ] Remove outdated design documents that are now fully implemented
- [ ] Keep only historical context valuable for understanding decisions
- [ ] Create archives/README.md explaining what's archived and why
- [ ] Move docs/BETA_DOCUMENTATION_REVIEW.md to archives
- [ ] Move docs/DOCUMENTATION_CONSOLIDATION.md to archives
- [ ] Move docs/MIGRATION_GUIDE_PHASE5.md to archives
- [ ] Ensure consistent file naming (convert underscores to hyphens)
### 0.3.2.5 Root Documentation (README.md and related documentation and root)
#### REFERENCE.md Restructuring
- [ ] perhaps REFERENCE.md should be removed from root and merged with docs/README.md
- [ ] Reorganize by user journey: Getting Started → Core Features → Advanced → Contributing
- [ ] Restructure and Update docs/README.md with relevant docs/ documentation and links and backlinks
- [ ] Update all links based on new documentation structure from 0.3.2.1-0.3.2.4
- [ ] Group documentation by type (API Reference, User Guides, Technical Docs)
- [ ] Add brief descriptions for each linked document
#### README.md Simplification
- [ ] Rewrite opening to clearly state go-llms value proposition
- [ ] Simplify quick start section with minimal example
- [ ] Create clear feature overview with links to detailed docs
- [ ] Reduce code examples to bare essentials
- [ ] Add clear navigation to different documentation sections
#### Other Root Files
- [ ] Create CHANGELOG.md consolidating all release notes
- [ ] Move RELEASE_NOTES_v0.3.1.md content into CHANGELOG.md
- [ ] Delete RELEASE_NOTES_v0.3.1.md after consolidation
- [ ] Ensure only approved markdown files remain in root

## v0.3.4: Built-in Agents Library
### 0.3.4.1: Text Processing Agents 
- [ ] TextSummarize - intelligent summarization using LLM
- [ ] TextExtract - extract structured data from text
- [ ] TextAnalyze - sentiment, entities, keywords
- [ ] TextTranslate - language translation using LLM

### 0.3.4.2: Research Agents 
- [ ] WebResearcher - web research with source tracking
- [ ] DocumentAnalyzer - analyze documents and PDFs
- [ ] FactChecker - verify claims against sources

### 0.3.4.3: Coding Agents 
- [ ] CodeReviewer - review code for issues
- [ ] TestGenerator - generate tests from code
- [ ] DocWriter - generate documentation

### 0.3.4.4: Data Agents 
- [ ] DataAnalyst - analyze datasets and generate insights
- [ ] ReportGenerator - create formatted reports
- [ ] DataCleaner - clean and validate data

### 0.3.4.5: Feed Agents 
- [ ] NewsMonitor - monitor news feeds for keywords and topics using LLM
- [ ] FeedAggregator - aggregate and deduplicate content from multiple feeds
- [ ] FeedSummarizer - summarize feed content using LLM
- [ ] ContentCurator - curate and categorize feed content using LLM

## v0.3.5: Built-in Workflow Patterns
### 0.3.5.1: Core Workflow Patterns 
- [ ] Pipeline - sequential processing workflow
- [ ] MapReduce - parallel processing with aggregation
- [ ] Consensus - multi-agent agreement pattern
- [ ] Retry - with exponential backoff

### 0.3.5.2: Example Workflows 
- [ ] ResearchWorkflow - research → verify → summarize → report
- [ ] CodeReviewWorkflow - analyze → review → suggest → document
- [ ] DataPipeline - ingest → clean → analyze → visualize


## v0.3.6: Enhanced Tool Capabilities
### 0.3.6.1: Web API Client Advanced Features 
- [ ] Auto-Pagination: Automatically follow pagination links
- [ ] Rate Limiting: Respect rate limit headers with intelligent backoff
- [ ] Response Caching: Cache responses with configurable TTL
- [ ] Request Templates: Store and reuse common request patterns
- [ ] Response Transformation: Extract data using JSONPath or JQ-like queries
- [ ] Error Recovery: Smart retries with exponential backoff
- [ ] Mock Mode: Optional mock responses for testing
- [ ] Streaming Responses: Handle large response streaming
- [ ] Request/Response Middleware: Plugin system for custom processing
- [ ] Multi-tenancy: Support multiple API configurations
- [ ] Request Batching: Batch multiple requests for efficiency

### 0.3.6.2: Authentication System Improvements 
- [ ] Create AuthProvider interface with Name(), CanHandle(), Configure() methods
- [ ] Implement AuthRegistry for managing multiple providers
- [ ] Add configuration file support for custom auth mappings (YAML/JSON)
- [ ] Support provider patterns:
  - [ ] URL pattern matching (regex support)
  - [ ] Response-based detection (401 + WWW-Authenticate header)
  - [ ] OpenAPI security scheme integration
  - [ ] OAuth2 discovery via .well-known endpoints
- [ ] Create default providers for common services:
  - [ ] GitHub (including Enterprise)
  - [ ] GitLab (including self-hosted)
  - [ ] Generic Bearer token
  - [ ] Generic API key
  - [ ] Basic auth
- [ ] Ensure backward compatibility with existing detectURLSpecificAuth
- [ ] Add examples and documentation

## v0.3.7: Advanced Agent Features
### 0.3.7.1: State Management 
- [ ] State persistence and serialization design
- [ ] Implement state storage backends (file, database)
- [ ] Add state versioning and migration
- [ ] Create examples with persistent agents

### 0.3.7.2: Agent Discovery
- [ ] Agent discovery and registry design
- [ ] Implement agent metadata and search
- [ ] Add dynamic agent loading
- [ ] Create agent marketplace example

### 0.3.7.3: Advanced Features 
- [ ] Advanced merge strategies for parallel agents
- [ ] Streaming support for long-running agents
- [ ] Agent composition patterns
- [ ] Agent lifecycle management

## v0.3.8: Model Context Protocol (MCP) Support
### 0.3.8.1: MCP Client Support 
- [ ] Research MCP specification and requirements
- [ ] Design MCP client interface for agents
- [ ] Implement MCP client in pkg/agent/mcp/client
- [ ] Add MCP tool discovery and registration
- [ ] Create examples demonstrating MCP client usage
- [ ] Write comprehensive tests
- [ ] Document MCP client usage in user guide

### 0.3.8.2: MCP Server Support 
- [ ] Design MCP server interface for exposing agents/workflows
- [ ] Implement MCP server in pkg/agent/mcp/server
- [ ] Add agent/workflow registration to MCP server
- [ ] Create example MCP server implementations
- [ ] Write comprehensive tests
- [ ] Document MCP server setup and configuration

## v0.3.9: Performance Optimization
### 0.3.9.1: Profiling Infrastructure (REVISIT)
- [ ] Create benchmark harness for A/B testing optimizations
- [ ] Implement visualization for memory allocation patterns
- [ ] Create real-world test scenarios for end-to-end performance

### 0.3.9.2: Advanced Optimizations (REVISIT)
- [ ] Implement adaptive channel buffer sizing based on usage patterns
- [ ] Add pool prewarming for high-throughput scenarios
- [ ] Reduce redundant property iterations in schema processing
- [ ] Implement more granular locking in cached objects
- [ ] Optimize zero-initialization patterns for pooled objects
- [ ] Introduce buffer pooling for string builders

### 0.3.9.3: Performance Validation (REVISIT)
- [ ] Document performance improvements with metrics
- [ ] Verify optimizations in high-concurrency scenarios
- [ ] Create benchmark comparison charts for before/after
- [ ] Implement regression testing to prevent performance degradation
- [ ] Add performance acceptance criteria to CI pipeline

## v0.3.10: Final Polish and Stable Release
### 0.3.10.1: Documentation Polish
- [ ] Fix identified cross-link issues (path inconsistencies, broken links)
- [ ] Perform final consistency check across all documentation
- [ ] Update all examples to showcase v0.3.x features
- [ ] Create migration guide from earlier versions

### 0.3.10.2: API Refinement
- [ ] API refinement based on usage feedback
- [ ] Deprecate old patterns with migration paths
- [ ] Ensure backward compatibility where possible
- [ ] Final review and preparation for v0.4.0 stable release

## Notes
- Tool System Enhancement Phases 1-4: COMPLETED (see TODO-DONE.md)
- Agent Architecture Restructuring Phases 1-7: COMPLETED (see TODO-DONE.md)
- Performance Phase 2 (Quick Wins): COMPLETED (see TODO-DONE.md)

See TODO-DONE.md for all completed tasks