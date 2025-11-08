# Performance Guide

Go-Clay is designed for high-performance UI layouts with microsecond-level computation times. This guide covers optimization strategies and performance best practices.

## Performance Characteristics

### Layout Algorithm Complexity
- **Time Complexity**: O(n) where n is the number of elements
- **Space Complexity**: O(n) for element storage
- **Memory Usage**: ~3.5MB for 8192 elements (typical)
- **Layout Time**: <1ms for 1000 elements on modern hardware

### Memory Management
Go-Clay uses arena-based allocation to minimize garbage collection pressure:

```go
// Arena allocation example
arena := clay.NewArena(1024 * 1024) // 1MB arena
engine := clay.NewLayoutEngineWithArena(arena)

// All allocations happen within the arena
engine.BeginLayout()
// ... layout code ...
commands := engine.EndLayout()

// Reset arena for next frame
arena.Reset()
```

## Optimization Strategies

### 1. Element Configuration Reuse

**❌ Inefficient:**
```go
// Creates new config each time
for i := 0; i < 1000; i++ {
    clay.Container(clay.IDWithIndex("item", i), clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(100),
                Height: clay.SizingFixed(50),
            },
        },
        BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
    })
}
```

**✅ Efficient:**
```go
// Reuse configuration
itemConfig := clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFixed(100),
            Height: clay.SizingFixed(50),
        },
    },
    BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
}

for i := 0; i < 1000; i++ {
    clay.Container(clay.IDWithIndex("item", i), itemConfig)
}
```

### 2. Minimize Element Hierarchy

**❌ Inefficient:**
```go
// Deep nesting
clay.Container("level1", config).
    Container("level2", config).
        Container("level3", config).
            Container("level4", config).
                Text("Content", textConfig)
```

**✅ Efficient:**
```go
// Flatten hierarchy when possible
clay.Container("content", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingGrow(0),
            Height: clay.SizingGrow(0),
        },
        Padding: clay.PaddingAll(16), // Use padding instead of nested containers
    },
}).
    Text("Content", textConfig)
```

### 3. Use Appropriate Sizing Types

**❌ Inefficient:**
```go
// Fixed sizing for dynamic content
clay.Container("text", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFixed(200), // May cause overflow
            Height: clay.SizingFixed(50),  // May cause overflow
        },
    },
})
```

**✅ Efficient:**
```go
// Use FIT for content-based sizing
clay.Container("text", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFit(),  // Sizes to content
            Height: clay.SizingFit(),  // Sizes to content
        },
    },
})
```

### 4. Batch Similar Operations

**❌ Inefficient:**
```go
// Individual operations
for i := 0; i < 100; i++ {
    clay.Container(clay.IDWithIndex("item", i), config)
}
```

**✅ Efficient:**
```go
// Batch operations
items := make([]clay.ElementConfig, 100)
for i := 0; i < 100; i++ {
    items[i] = config
}

clay.BatchContainers("items", items)
```

### 5. Optimize Text Rendering

**❌ Inefficient:**
```go
// Multiple text elements
clay.Text("Line 1", textConfig).
clay.Text("Line 2", textConfig).
clay.Text("Line 3", textConfig)
```

**✅ Efficient:**
```go
// Single text element with line breaks
clay.Text("Line 1\nLine 2\nLine 3", textConfig)
```

## Memory Optimization

### Arena Allocation
```go
// Create arena with appropriate size
arena := clay.NewArena(2 * 1024 * 1024) // 2MB for large layouts
engine := clay.NewLayoutEngineWithArena(arena)

// Use arena for the entire frame
engine.BeginLayout()
// ... layout code ...
commands := engine.EndLayout()

// Reset for next frame
arena.Reset()
```

### Object Pooling
```go
type LayoutPool struct {
    configs sync.Pool
    texts   sync.Pool
}

func (p *LayoutPool) GetConfig() *clay.ElementConfig {
    if v := p.configs.Get(); v != nil {
        return v.(*clay.ElementConfig)
    }
    return &clay.ElementConfig{}
}

func (p *LayoutPool) PutConfig(config *clay.ElementConfig) {
    // Reset config
    *config = clay.ElementConfig{}
    p.configs.Put(config)
}
```

### String Interning
```go
// Use string interning for repeated text
type StringInterner struct {
    strings map[string]string
    mutex   sync.RWMutex
}

func (si *StringInterner) Intern(s string) string {
    si.mutex.RLock()
    if interned, exists := si.strings[s]; exists {
        si.mutex.RUnlock()
        return interned
    }
    si.mutex.RUnlock()
    
    si.mutex.Lock()
    defer si.mutex.Unlock()
    
    if interned, exists := si.strings[s]; exists {
        return interned
    }
    
    si.strings[s] = s
    return s
}
```

## Profiling and Monitoring

### Performance Metrics
```go
// Get layout statistics
stats := engine.GetStats()
fmt.Printf("Elements: %d\n", stats.ElementCount)
fmt.Printf("Render Commands: %d\n", stats.RenderCommands)
fmt.Printf("Layout Time: %v\n", stats.LayoutTime)
fmt.Printf("Memory Used: %d bytes\n", stats.MemoryUsed)
```

### Benchmarking
```go
func BenchmarkLayout(b *testing.B) {
    engine := clay.NewLayoutEngine()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine.BeginLayout()
        
        // Create test layout
        clay.Container("main", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingGrow(0),
                },
            },
        }).
            Text("Hello World", clay.TextConfig{FontSize: 16})
        
        commands := engine.EndLayout()
        _ = commands
    }
}
```

### Memory Profiling
```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Your layout code here
}
```

## Renderer Optimization

### Command Batching
```go
type OptimizedRenderer struct {
    rectangles []clay.RectangleCommand
    texts      []clay.TextCommand
    images     []clay.ImageCommand
}

func (r *OptimizedRenderer) RenderFrame(commands []clay.RenderCommand) error {
    // Clear previous frame
    r.rectangles = r.rectangles[:0]
    r.texts = r.texts[:0]
    r.images = r.images[:0]
    
    // Batch commands by type
    for _, cmd := range commands {
        switch cmd.CommandType {
        case clay.CommandRectangle:
            r.rectangles = append(r.rectangles, cmd.Data.(clay.RectangleCommand))
        case clay.CommandText:
            r.texts = append(r.texts, cmd.Data.(clay.TextCommand))
        case clay.CommandImage:
            r.images = append(r.images, cmd.Data.(clay.ImageCommand))
        }
    }
    
    // Batch render
    r.renderRectangles()
    r.renderTexts()
    r.renderImages()
    
    return nil
}
```

### GPU Upload Optimization
```go
type GPURenderer struct {
    vertexBuffer []float32
    indexBuffer  []uint32
}

func (r *GPURenderer) UploadGeometry(commands []clay.RenderCommand) {
    // Pre-allocate buffers
    r.vertexBuffer = r.vertexBuffer[:0]
    r.indexBuffer = r.indexBuffer[:0]
    
    // Batch upload to GPU
    for _, cmd := range commands {
        r.addCommandToBuffers(cmd)
    }
    
    // Single GPU upload
    gl.BufferData(gl.ARRAY_BUFFER, len(r.vertexBuffer)*4, 
                  gl.Ptr(r.vertexBuffer), gl.DYNAMIC_DRAW)
}
```

## Common Performance Pitfalls

### 1. Excessive Element Creation
```go
// ❌ Don't create elements in hot paths
func renderFrame() {
    for i := 0; i < 1000; i++ {
        clay.Container(clay.IDWithIndex("item", i), config) // Expensive
    }
}

// ✅ Pre-create elements
var elements []clay.ElementConfig
func init() {
    elements = make([]clay.ElementConfig, 1000)
    for i := range elements {
        elements[i] = config
    }
}

func renderFrame() {
    for i := 0; i < 1000; i++ {
        clay.Container(clay.IDWithIndex("item", i), elements[i])
    }
}
```

### 2. Inefficient Text Measurement
```go
// ❌ Don't measure text repeatedly
func renderText(text string) {
    for i := 0; i < len(text); i++ {
        size := measureText(text[:i]) // Expensive
    }
}

// ✅ Cache text measurements
type TextCache struct {
    cache map[string]clay.Dimensions
}

func (tc *TextCache) MeasureText(text string, config clay.TextConfig) clay.Dimensions {
    key := fmt.Sprintf("%s_%d", text, config.FontSize)
    if size, exists := tc.cache[key]; exists {
        return size
    }
    
    size := measureText(text, config)
    tc.cache[key] = size
    return size
}
```

### 3. Unnecessary Layout Recalculation
```go
// ❌ Don't recalculate layout every frame
func renderFrame() {
    engine.BeginLayout()
    // ... layout code ...
    commands := engine.EndLayout()
    renderer.Render(commands)
}

// ✅ Only recalculate when needed
var lastLayoutHash uint64
func renderFrame() {
    currentHash := calculateLayoutHash()
    if currentHash != lastLayoutHash {
        engine.BeginLayout()
        // ... layout code ...
        commands = engine.EndLayout()
        lastLayoutHash = currentHash
    }
    
    renderer.Render(commands)
}
```

## Performance Testing

### Load Testing
```go
func TestLayoutPerformance(t *testing.T) {
    engine := clay.NewLayoutEngine()
    
    // Test with increasing element counts
    elementCounts := []int{100, 500, 1000, 2000, 5000}
    
    for _, count := range elementCounts {
        t.Run(fmt.Sprintf("Elements_%d", count), func(t *testing.T) {
            start := time.Now()
            
            engine.BeginLayout()
            for i := 0; i < count; i++ {
                clay.Container(clay.IDWithIndex("item", i), config)
            }
            commands := engine.EndLayout()
            
            duration := time.Since(start)
            t.Logf("Layout time for %d elements: %v", count, duration)
            
            // Assert performance requirements
            assert.Less(t, duration, 10*time.Millisecond, 
                       "Layout should complete within 10ms")
        })
    }
}
```

### Memory Testing
```go
func TestMemoryUsage(t *testing.T) {
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    // Create large layout
    for i := 0; i < 10000; i++ {
        clay.Container(clay.IDWithIndex("item", i), config)
    }
    commands := engine.EndLayout()
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    memoryUsed := m2.Alloc - m1.Alloc
    t.Logf("Memory used: %d bytes", memoryUsed)
    
    // Assert memory usage is reasonable
    assert.Less(t, memoryUsed, 10*1024*1024, // 10MB
                "Memory usage should be less than 10MB")
}
```

By following these optimization strategies and best practices, you can achieve excellent performance with Go-Clay even in complex UI scenarios with thousands of elements.
