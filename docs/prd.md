# Go-Clay Product Requirements Document (PRD)

## Goals and Background Context

### Goals

- Deliver a pure Go implementation of the Clay layout engine with microsecond performance
- Eliminate CGo dependency for Go developers building UIs
- Provide idiomatic Go APIs with full type safety and familiar patterns
- Achieve feature parity with original C Clay library core functionality
- Enable building high-performance UIs in Go without HTML/CSS overhead
- Create foundation for Go-native UI framework ecosystem

### Background Context

The Clay C library revolutionized UI layout by separating computation from rendering, delivering microsecond-level performance for complex layouts using a flexbox-like model. However, it exists only in C, forcing Go developers into CGo with performance penalties, compilation complexity, and loss of Go's type safety.

Go-Clay solves this by implementing the core Clay layout engine in pure Go, maintaining performance while eliminating CGo entirely. This enables Go developers to leverage Clay's algorithmic advantages with idiomatic Go APIs, better tooling, and cross-platform simplicity.

The project targets developers building UI frameworks, desktop applications, terminal tools, and games in Go who need sophisticated, high-performance layouts without external dependencies or HTML/CSS overhead.

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2024-01-XX | 1.0 | Initial PRD creation | System |

## Requirements

### Functional Requirements

- **FR1**: The layout engine MUST support four sizing types: FitToContent (wrap), GrowToFill (flex-grow), PercentOfParent (percentage), FixedPixelSize (absolute).
- **FR2**: The layout engine MUST support two layout directions: LeftToRight (horizontal) and TopToBottom (vertical).
- **FR3**: The layout engine MUST support child alignment options for X (Left, Center, Right) and Y (Top, Center, Bottom) axes.
- **FR4**: The layout engine MUST support per-element padding with independent Left, Right, Top, Bottom values.
- **FR5**: The layout engine MUST support child gap spacing between elements in the layout direction.
- **FR6**: The layout engine MUST provide an ElementID system that can be generated from strings, string+index, or sequential integers.
- **FR7**: The layout engine MUST maintain parent-child relationships in a tree structure for all elements.
- **FR8**: The layout engine MUST compute bounding boxes for all elements based on sizing, padding, and content.
- **FR9**: The layout engine MUST output render commands including: RectangleCommand, TextCommand, ImageCommand, BorderCommand, ClipStartCommand, ClipEndCommand.
- **FR10**: RectangleCommand MUST include bounding box, color, and corner radius (four independent corners).
- **FR11**: TextCommand MUST include bounding box, text content, font ID, font size, color, line height, letter spacing, alignment, and wrap mode.
- **FR12**: ImageCommand MUST include bounding box, image data reference, tint color, and corner radius.
- **FR13**: BorderCommand MUST include bounding box, color, width (four sides), and corner radius.
- **FR14**: The layout engine MUST support arena-based memory allocation with configurable capacity.
- **FR15**: The layout engine MUST provide arena Reset() functionality for reuse between frames.
- **FR16**: The layout engine MUST track and expose memory usage statistics.
- **FR17**: The Gio UI renderer MUST implement full support for all render command types.
- **FR18**: The Gio UI renderer MUST integrate with Gio UI v0.9.0 without version conflicts.
- **FR19**: The Gio UI renderer MUST provide NewRenderer() factory function returning a Renderer interface.
- **FR20**: The layout engine MUST provide fluent API for building layouts (ContainerBuilder pattern).
- **FR21**: The layout engine MUST generate IDs automatically for elements without explicit IDs.
- **FR22**: The layout engine MUST handle nested container hierarchies of arbitrary depth.
- **FR23**: The layout engine MUST support text elements as leaf nodes with measurement.
- **FR24**: The layout engine MUST support image elements with aspect ratio preservation.
- **FR25**: The API MUST provide helper functions for common operations: ColorRGB(), ColorRGBA(), PaddingAll(), SizingFit(), SizingGrow(), CornerRadiusAll().
- **FR26**: The layout engine MUST expose performance statistics: ElementCount, RenderCommands, LayoutTime, MemoryUsed.

### Non-Functional Requirements

- **NFR1**: Layout computation MUST complete in <5μs for 1000 elements on modern CPU (baseline: C implementation).
- **NFR2**: Memory usage MUST stay under 50MB for 1000-element layouts.
- **NFR3**: The core library (excluding renderers) MUST have zero external dependencies.
- **NFR4**: All public APIs MUST be documented with godoc comments.
- **NFR5**: The library MUST pass all tests with coverage >90%.
- **NFR6**: The library MUST be compatible with Go 1.23+ and Go 1.24 toolchain.
- **NFR7**: The API MUST follow Go naming conventions (exported types start with capital letter).
- **NFR8**: The implementation MUST not use unsafe package or undefined behavior.
- **NFR9**: Error conditions MUST be handled gracefully without panics in public APIs (except for programming errors).
- **NFR10**: The layout engine MUST be thread-safe for concurrent BeginLayout()/EndLayout() calls.
- **NFR11**: The fluent API MUST support method chaining for building complex layouts.
- **NFR12**: The renderer interface MUST allow custom renderer implementations.
- **NFR13**: Layout commands MUST be efficiently serializable for potential networking or persistence.
- **NFR14**: The library MUST work on all Go-supported platforms without platform-specific code.
- **NFR15**: API stability MUST be maintained within major versions (1.x.x → 1.y.z).

## User Interface Design Goals

### Overall UX Vision

Go-Clay is a developer-focused library, not an end-user application. The "user interface" is the Go API itself. The design goals prioritize:

**Declarative Layout**: Developers describe layouts declaratively, similar to React or HTML templates, without imperative positioning code. The layout engine handles all computation.

**Composable Abstraction**: Complex layouts built from simple, reusable pieces. Container + child elements compose into trees naturally.

**Performance by Default**: Fast layouts without optimization effort. Core algorithms are optimized for common cases.

**Render Flexibility**: Use any rendering backend (Gio UI, terminal, OpenGL) without changing layout code.

### Key Interaction Paradigms

**Fluent Builder Pattern**: Container().Text().Container().End() chains for readable layout code.

**Element Configuration**: Configure elements via struct composition (ElementConfig, TextConfig, ImageConfig).

**Layout Lifecycle**: Explicit BeginLayout() → declare elements → EndLayout() → render commands pattern.

**Functional Helpers**: Provide convenience functions for common configurations while allowing full struct control.

### Core Screens and Views

Not applicable - this is a library API, not an application UI.

### Accessibility

Not applicable - layout engine responsibility only. Rendering and accessibility handled by renderers.

### Branding

None. Clean, professional API design with Go standard library aesthetic.

### Target Device and Platforms

**Cross-Platform**: Linux, macOS, Windows, WebAssembly, Terminals

## Technical Assumptions

### Repository Structure

**Monorepo**: Single repository containing core library and multiple renderer implementations.

Structure:
```
go-clay/
├── clay/          # Core layout engine
├── renderers/     # Renderer implementations
│   └── gioui/    # Gio UI renderer
├── examples/      # Example applications
├── docs/          # Documentation
└── tests/         # Test suites
```

### Service Architecture

**Library Architecture**: Embedded component, not a service. Applications import and instantiate the layout engine.

### Technology Choices

**Core Library**: Pure Go, standard library only.

**Primary Renderer**: Gio UI v0.9.0.

**Testing**: Standard library testing package, testify for assertions.

**Benchmarking**: Go testing.B, custom benchmarks against C implementation.

### Architecture Style

**Modular Design**: Clear interface boundaries between layout engine and renderers.

**Dependency Injection**: Renderer implementation not hardcoded, configurable per application.

**Arena Allocation**: Predictable memory usage with arena-based allocation.

**Command Pattern**: Layout outputs render commands, renderers consume and execute.

## Out of Scope for MVP

- Terminal renderer (planned for Phase 2)
- SDL2/SDL3 renderers (planned for Phase 2)
- OpenGL renderer (planned for Phase 2)
- WebAssembly renderer (planned for Phase 2)
- Floating elements with z-index management
- Scrolling containers and scroll offset tracking
- Pointer state and hit testing
- Event system and callbacks
- Text measurement caching (basic implementation only)
- Incremental layout updates
- Animation support
- Conditional layout/visibility
- Layout debugging tools or visualization
- Layout inspector UI

## Acceptance Criteria

### Layout Engine

- [ ] All FR1-FR26 requirements implemented
- [ ] NFR1 performance target achieved (benchmark results)
- [ ] NFR2 memory target achieved (memory profiling)
- [ ] 90%+ test coverage with edge case coverage
- [ ] Complete godoc documentation
- [ ] Zero external dependencies verified

### Gio UI Renderer

- [ ] All render command types functional
- [ ] Integration with Gio UI v0.9.0 validated
- [ ] Example application runs successfully
- [ ] Renderer code coverage >80%
- [ ] Performance suitable for 60 FPS applications

### Documentation

- [ ] API reference complete
- [ ] Getting started guide with working example
- [ ] Layout system explanation
- [ ] Renderer guide for custom implementations
- [ ] Performance optimization tips
- [ ] Example gallery with 5+ examples

### Success Metrics

- [ ] Performance within 10% of C implementation
- [ ] Zero known bugs in core layout algorithms
- [ ] API usability validated by 3+ early adopters
- [ ] Documentation enables new users to build layouts in <30 minutes

## Open Questions

1. Should text measurement require a TextMeasurer interface or be built into the engine?
2. How should corner radius affect border rendering in corner cases?
3. What's the strategy for handling images without texture/format abstraction in core?
4. Should the API expose element tree for debugging or keep it internal?
5. How to handle platform-specific behavior (e.g., different text rendering)?

## Next Steps

1. Define detailed architecture document
2. Create initial story/epic breakdown
3. Set up development environment and CI/CD
4. Begin MVP implementation with TDD approach
5. Establish performance benchmarking framework

