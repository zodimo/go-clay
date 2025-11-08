# Layout System

Go-Clay implements a flexbox-like layout system that provides powerful, responsive UI layouts with minimal configuration.

## Core Concepts

### Sizing System

The layout system uses four sizing types to control how elements take up space:

#### SizingFit
Wraps tightly to content size. Used for text, images, and elements that should size to their content.

```go
clay.Sizing{
    Width:  clay.SizingFit(),
    Height: clay.SizingFit(),
}
```

#### SizingGrow
Expands to fill available space, sharing it with other grow elements. The weight parameter controls how much space this element gets relative to other grow elements.

```go
clay.Sizing{
    Width:  clay.SizingGrow(1.0), // Takes 1/3 of available space
    Height: clay.SizingGrow(2.0), // Takes 2/3 of available space
}
```

#### SizingPercent
Takes a percentage of the parent's size (0.0-1.0).

```go
clay.Sizing{
    Width:  clay.SizingPercent(0.5), // 50% of parent width
    Height: clay.SizingPercent(0.25), // 25% of parent height
}
```

#### SizingFixed
Fixed pixel size, regardless of available space.

```go
clay.Sizing{
    Width:  clay.SizingFixed(200), // Always 200px wide
    Height: clay.SizingFixed(100), // Always 100px tall
}
```

### Layout Directions

#### LeftToRight (Default)
Children are arranged horizontally from left to right.

```go
clay.LayoutConfig{
    Direction: clay.LeftToRight,
    ChildGap:  8, // 8px gap between children
}
```

#### TopToBottom
Children are arranged vertically from top to bottom.

```go
clay.LayoutConfig{
    Direction: clay.TopToBottom,
    ChildGap:  12, // 12px gap between children
}
```

### Alignment

#### Child Alignment
Controls how children are positioned within their parent container.

```go
clay.LayoutConfig{
    ChildAlignment: clay.ChildAlignment{
        X: clay.AlignXCenter, // Center horizontally
        Y: clay.AlignYCenter, // Center vertically
    },
}
```

Available alignments:
- **X**: `AlignXLeft`, `AlignXCenter`, `AlignXRight`
- **Y**: `AlignYTop`, `AlignYCenter`, `AlignYBottom`

### Padding and Gaps

#### Padding
Space between the element's border and its children.

```go
clay.LayoutConfig{
    Padding: clay.PaddingAll(16), // 16px on all sides
}

// Or specify individual sides
clay.LayoutConfig{
    Padding: clay.Padding{
        Left:   8,
        Right:  16,
        Top:    4,
        Bottom: 12,
    },
}
```

#### Child Gap
Space between child elements along the layout direction.

```go
clay.LayoutConfig{
    ChildGap: 8, // 8px between children
}
```

## Layout Algorithm

Go-Clay uses a two-pass layout algorithm for optimal performance:

### Pass 1: Size Calculation
1. Traverse element tree in depth-first order
2. Calculate element sizes based on content and constraints
3. Handle text wrapping and aspect ratios
4. Propagate size changes up the tree

### Pass 2: Position Calculation
1. Traverse element tree again
2. Calculate positions based on layout configuration
3. Handle alignment, padding, and gaps
4. Generate render commands

## Common Layout Patterns

### Sidebar Layout
```go
clay.Container("main", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingGrow(0),
            Height: clay.SizingGrow(0),
        },
        Direction: clay.LeftToRight,
        ChildGap:  16,
    },
}).Container("sidebar", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFixed(300),
            Height: clay.SizingGrow(0),
        },
        Direction: clay.TopToBottom,
        Padding:   clay.PaddingAll(16),
    },
    BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
}).Container("content", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingGrow(0),
            Height: clay.SizingGrow(0),
        },
        Padding: clay.PaddingAll(16),
    },
})
```

### Centered Content
```go
clay.Container("center", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFixed(400),
            Height: clay.SizingFixed(300),
        },
        ChildAlignment: clay.ChildAlignment{
            X: clay.AlignXCenter,
            Y: clay.AlignYCenter,
        },
    },
    BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
})
```

### Responsive Grid
```go
clay.Container("grid", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingGrow(0),
            Height: clay.SizingGrow(0),
        },
        Direction: clay.LeftToRight,
        ChildGap:  8,
    },
}).Container("item1", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingPercent(0.33),
            Height: clay.SizingFixed(100),
        },
    },
}).Container("item2", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingPercent(0.33),
            Height: clay.SizingFixed(100),
        },
    },
}).Container("item3", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingPercent(0.34),
            Height: clay.SizingFixed(100),
        },
    },
})
```

## Advanced Features

### Text Wrapping
```go
clay.Text("Long text that will wrap...", clay.TextConfig{
    FontSize: 16,
    WrapMode: clay.WrapWords,
    Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
})
```

### Aspect Ratios
```go
clay.Container("square", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFixed(200),
            Height: clay.SizingFit(),
        },
    },
    AspectRatio: 1.0, // Square aspect ratio
})
```

### Floating Elements
```go
clay.Container("floating", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFixed(100),
            Height: clay.SizingFixed(100),
        },
    },
    Floating: clay.FloatingConfig{
        AttachTo: clay.AttachToParent,
        ZIndex:   10,
        Offset:   clay.Vector2{X: 10, Y: 10},
    },
})
```

### Scroll Containers
```go
clay.Container("scrollable", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingFixed(300),
            Height: clay.SizingFixed(200),
        },
    },
    Clip: clay.ClipConfig{
        Horizontal: true,
        Vertical:   true,
    },
})
```

## Performance Tips

1. **Use SizingFit for text** - More efficient than fixed sizing
2. **Minimize nested containers** - Flatten hierarchy when possible
3. **Reuse element configurations** - Cache common configs
4. **Use appropriate sizing types** - Fixed when you know the size, Grow for flexible layouts
5. **Limit floating elements** - They require additional computation

## Debugging

Enable debug mode to visualize layout bounds:

```go
engine := clay.NewLayoutEngine()
engine.SetDebugMode(true)
```

This will render colored outlines around all elements, making it easy to see how the layout system is positioning your elements.
