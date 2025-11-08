# Examples

This document provides comprehensive examples of using Go-Clay for various UI scenarios.

## Basic Examples

### Simple Container
```go
package main

import (
    "github.com/zodimo/go-clay"
    "github.com/zodimo/go-clay/renderers/gioui"
)

func main() {
    engine := clay.NewLayoutEngine()
    renderer := gioui.NewRenderer()
    
    engine.BeginLayout()
    
    clay.Container("main", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
            Padding: clay.PaddingAll(16),
        },
        BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
    }).Text("Hello, World!", clay.TextConfig{
        FontSize: 24,
        Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
    })
    
    commands := engine.EndLayout()
    renderer.Render(commands)
}
```

### Sidebar Layout
```go
func createSidebarLayout() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    // Main container
    clay.Container("main", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
            Direction: clay.LeftToRight,
            ChildGap:  16,
        },
    }).
    // Sidebar
    Container("sidebar", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(300),
                Height: clay.SizingGrow(0),
            },
            Direction: clay.TopToBottom,
            Padding:   clay.PaddingAll(16),
        },
        BackgroundColor: clay.Color{R: 0.8, G: 0.8, B: 0.9, A: 1.0},
    }).
        Text("Sidebar", clay.TextConfig{
            FontSize: 20,
            Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
        }).
        Container("nav", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingGrow(0),
                },
                Direction: clay.TopToBottom,
                ChildGap:  8,
            },
        }).
            Text("Home", clay.TextConfig{FontSize: 16}).
            Text("About", clay.TextConfig{FontSize: 16}).
            Text("Contact", clay.TextConfig{FontSize: 16}).
        End().
    // Main content
    Container("content", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
            Padding: clay.PaddingAll(16),
        },
        BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
    }).
        Text("Main Content", clay.TextConfig{
            FontSize: 24,
            Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
        }).
        Text("This is the main content area.", clay.TextConfig{
            FontSize: 16,
            Color:    clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
        })
}
```

## Advanced Examples

### Responsive Grid
```go
func createResponsiveGrid() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    clay.Container("grid", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
            Direction: clay.LeftToRight,
            ChildGap:  8,
        },
    }).
    // Grid items
    Container("item1", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingPercent(0.33),
                Height: clay.SizingFixed(100),
            },
        },
        BackgroundColor: clay.Color{R: 0.9, G: 0.5, B: 0.5, A: 1.0},
    }).
        Text("Item 1", clay.TextConfig{FontSize: 16}).
    Container("item2", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingPercent(0.33),
                Height: clay.SizingFixed(100),
            },
        },
        BackgroundColor: clay.Color{R: 0.5, G: 0.9, B: 0.5, A: 1.0},
    }).
        Text("Item 2", clay.TextConfig{FontSize: 16}).
    Container("item3", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingPercent(0.34),
                Height: clay.SizingFixed(100),
            },
        },
        BackgroundColor: clay.Color{R: 0.5, G: 0.5, B: 0.9, A: 1.0},
    }).
        Text("Item 3", clay.TextConfig{FontSize: 16})
}
```

### Form Layout
```go
func createForm() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    clay.Container("form", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(400),
                Height: clay.SizingGrow(0),
            },
            Direction: clay.TopToBottom,
            ChildGap:  16,
            Padding:  clay.PaddingAll(20),
        },
        BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
    }).
        Text("Contact Form", clay.TextConfig{
            FontSize: 24,
            Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
        }).
        // Name field
        Container("name-field", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingFixed(60),
                },
                Direction: clay.TopToBottom,
                ChildGap:   4,
            },
        }).
            Text("Name", clay.TextConfig{
                FontSize: 14,
                Color:    clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
            }).
            Container("name-input", clay.ElementConfig{
                Layout: clay.LayoutConfig{
                    Sizing: clay.Sizing{
                        Width:  clay.SizingGrow(0),
                        Height: clay.SizingFixed(40),
                    },
                },
                BackgroundColor: clay.Color{R: 0.95, G: 0.95, B: 0.95, A: 1.0},
                Border: &clay.BorderConfig{
                    Width: clay.BorderWidth{All: 1},
                    Color: clay.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
                },
            }).
                Text("Enter your name", clay.TextConfig{
                    FontSize: 16,
                    Color:    clay.Color{R: 0.6, G: 0.6, B: 0.6, A: 1.0},
                }).
            End().
        // Email field
        Container("email-field", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingFixed(60),
                },
                Direction: clay.TopToBottom,
                ChildGap:   4,
            },
        }).
            Text("Email", clay.TextConfig{
                FontSize: 14,
                Color:    clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
            }).
            Container("email-input", clay.ElementConfig{
                Layout: clay.LayoutConfig{
                    Sizing: clay.Sizing{
                        Width:  clay.SizingGrow(0),
                        Height: clay.SizingFixed(40),
                    },
                },
                BackgroundColor: clay.Color{R: 0.95, G: 0.95, B: 0.95, A: 1.0},
                Border: &clay.BorderConfig{
                    Width: clay.BorderWidth{All: 1},
                    Color: clay.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
                },
            }).
                Text("Enter your email", clay.TextConfig{
                    FontSize: 16,
                    Color:    clay.Color{R: 0.6, G: 0.6, B: 0.6, A: 1.0},
                }).
            End().
        // Submit button
        Container("submit-button", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingFixed(120),
                    Height: clay.SizingFixed(40),
                },
                ChildAlignment: clay.ChildAlignment{
                    X: clay.AlignXCenter,
                    Y: clay.AlignYCenter,
                },
            },
            BackgroundColor: clay.Color{R: 0.2, G: 0.6, B: 1.0, A: 1.0},
            CornerRadius: clay.CornerRadius{All: 4},
        }).
            Text("Submit", clay.TextConfig{
                FontSize: 16,
                Color:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
            })
}
```

### Card Layout
```go
func createCard() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    clay.Container("card", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(300),
                Height: clay.SizingGrow(0),
            },
            Direction: clay.TopToBottom,
            Padding:  clay.PaddingAll(16),
        },
        BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
        Border: &clay.BorderConfig{
            Width: clay.BorderWidth{All: 1},
            Color: clay.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
        },
        CornerRadius: clay.CornerRadius{All: 8},
    }).
        // Card header
        Container("header", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingFixed(60),
                },
                Direction: clay.LeftToRight,
                ChildAlignment: clay.ChildAlignment{
                    X: clay.AlignXLeft,
                    Y: clay.AlignYCenter,
                },
            },
        }).
            Container("avatar", clay.ElementConfig{
                Layout: clay.LayoutConfig{
                    Sizing: clay.Sizing{
                        Width:  clay.SizingFixed(40),
                        Height: clay.SizingFixed(40),
                    },
                },
                BackgroundColor: clay.Color{R: 0.7, G: 0.7, B: 0.7, A: 1.0},
                CornerRadius: clay.CornerRadius{All: 20},
            }).
            Container("title-section", clay.ElementConfig{
                Layout: clay.LayoutConfig{
                    Sizing: clay.Sizing{
                        Width:  clay.SizingGrow(0),
                        Height: clay.SizingGrow(0),
                    },
                    Direction: clay.TopToBottom,
                    Padding:  clay.Padding{Left: 12, Right: 0, Top: 0, Bottom: 0},
                },
            }).
                Text("John Doe", clay.TextConfig{
                    FontSize: 18,
                    Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
                }).
                Text("Software Engineer", clay.TextConfig{
                    FontSize: 14,
                    Color:    clay.Color{R: 0.5, G: 0.5, B: 0.5, A: 1.0},
                }).
            End().
        // Card content
        Container("content", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingGrow(0),
                },
                Direction: clay.TopToBottom,
                ChildGap:  12,
            },
        }).
            Text("This is a sample card with some content. It demonstrates how to create a card layout with header and content sections.", clay.TextConfig{
                FontSize: 14,
                Color:    clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
            }).
            Container("tags", clay.ElementConfig{
                Layout: clay.LayoutConfig{
                    Sizing: clay.Sizing{
                        Width:  clay.SizingGrow(0),
                        Height: clay.SizingFixed(30),
                    },
                    Direction: clay.LeftToRight,
                    ChildGap:  8,
                },
            }).
                Container("tag1", clay.ElementConfig{
                    Layout: clay.LayoutConfig{
                        Sizing: clay.Sizing{
                            Width:  clay.SizingFit(),
                            Height: clay.SizingFixed(24),
                        },
                        Padding: clay.Padding{Left: 8, Right: 8, Top: 4, Bottom: 4},
                    },
                    BackgroundColor: clay.Color{R: 0.2, G: 0.6, B: 1.0, A: 1.0},
                    CornerRadius: clay.CornerRadius{All: 12},
                }).
                    Text("Go", clay.TextConfig{
                        FontSize: 12,
                        Color:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
                    }).
                Container("tag2", clay.ElementConfig{
                    Layout: clay.LayoutConfig{
                        Sizing: clay.Sizing{
                            Width:  clay.SizingFit(),
                            Height: clay.SizingFixed(24),
                        },
                        Padding: clay.Padding{Left: 8, Right: 8, Top: 4, Bottom: 4},
                    },
                    BackgroundColor: clay.Color{R: 0.8, G: 0.4, B: 0.8, A: 1.0},
                    CornerRadius: clay.CornerRadius{All: 12},
                }).
                    Text("UI", clay.TextConfig{
                        FontSize: 12,
                        Color:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
                    })
}
```

## Interactive Examples

### Button with Hover State
```go
func createInteractiveButton() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    buttonID := clay.ID("my-button")
    
    // Check if pointer is over button
    isHovered := engine.IsPointerOver(buttonID)
    
    var bgColor clay.Color
    if isHovered {
        bgColor = clay.Color{R: 0.1, G: 0.5, B: 0.9, A: 1.0}
    } else {
        bgColor = clay.Color{R: 0.2, G: 0.6, B: 1.0, A: 1.0}
    }
    
    clay.Container(buttonID, clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(120),
                Height: clay.SizingFixed(40),
            },
            ChildAlignment: clay.ChildAlignment{
                X: clay.AlignXCenter,
                Y: clay.AlignYCenter,
            },
        },
        BackgroundColor: bgColor,
        CornerRadius: clay.CornerRadius{All: 4},
    }).
        Text("Click Me", clay.TextConfig{
            FontSize: 16,
            Color:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
        })
}
```

### Scrollable List
```go
func createScrollableList() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    clay.Container("scroll-container", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(300),
                Height: clay.SizingFixed(200),
            },
        },
        Clip: clay.ClipConfig{
            Horizontal: false,
            Vertical:   true,
        },
    }).
        // List items
        Container("item1", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingFixed(50),
                },
                Padding: clay.PaddingAll(8),
            },
            BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
        }).
            Text("Item 1", clay.TextConfig{FontSize: 16}).
        Container("item2", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingFixed(50),
                },
                Padding: clay.PaddingAll(8),
            },
            BackgroundColor: clay.Color{R: 0.95, G: 0.95, B: 0.95, A: 1.0},
        }).
            Text("Item 2", clay.TextConfig{FontSize: 16}).
        Container("item3", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingGrow(0),
                    Height: clay.SizingFixed(50),
                },
                Padding: clay.PaddingAll(8),
            },
            BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
        }).
            Text("Item 3", clay.TextConfig{FontSize: 16})
}
```

## Component Patterns

### Reusable Button Component
```go
func createButtonComponent(id clay.ElementID, text string, config clay.ElementConfig) *clay.ContainerBuilder {
    return clay.Container(id, clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(120),
                Height: clay.SizingFixed(40),
            },
            ChildAlignment: clay.ChildAlignment{
                X: clay.AlignXCenter,
                Y: clay.AlignYCenter,
            },
        },
        BackgroundColor: clay.Color{R: 0.2, G: 0.6, B: 1.0, A: 1.0},
        CornerRadius: clay.CornerRadius{All: 4},
    }).
        Text(text, clay.TextConfig{
            FontSize: 16,
            Color:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
        })
}

// Usage
func useButtonComponent() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    createButtonComponent(clay.ID("btn1"), "Save", clay.ElementConfig{}).
    createButtonComponent(clay.ID("btn2"), "Cancel", clay.ElementConfig{}).
    createButtonComponent(clay.ID("btn3"), "Delete", clay.ElementConfig{
        BackgroundColor: clay.Color{R: 0.8, G: 0.2, B: 0.2, A: 1.0},
    })
}
```

### Modal Dialog
```go
func createModalDialog() {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    // Backdrop
    clay.Container("backdrop", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
            ChildAlignment: clay.ChildAlignment{
                X: clay.AlignXCenter,
                Y: clay.AlignYCenter,
            },
        },
        BackgroundColor: clay.Color{R: 0, G: 0, B: 0, A: 0.5},
    }).
        // Modal content
        Container("modal", clay.ElementConfig{
            Layout: clay.LayoutConfig{
                Sizing: clay.Sizing{
                    Width:  clay.SizingFixed(400),
                    Height: clay.SizingFixed(300),
                },
            },
            BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
            CornerRadius: clay.CornerRadius{All: 8},
        }).
            Text("Modal Title", clay.TextConfig{
                FontSize: 20,
                Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
            }).
            Text("This is modal content.", clay.TextConfig{
                FontSize: 16,
                Color:    clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
            })
}
```

## Performance Examples

### Efficient List Rendering
```go
func createEfficientList(items []string) {
    engine := clay.NewLayoutEngine()
    engine.BeginLayout()
    
    clay.Container("list", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
            Direction: clay.TopToBottom,
            ChildGap:  1,
        },
    }).
        // Render only visible items
        func() {
            for i, item := range items {
                if i >= 100 { // Limit visible items
                    break
                }
                
                clay.Container(clay.IDWithIndex("item", i), clay.ElementConfig{
                    Layout: clay.LayoutConfig{
                        Sizing: clay.Sizing{
                            Width:  clay.SizingGrow(0),
                            Height: clay.SizingFixed(40),
                        },
                        Padding: clay.PaddingAll(8),
                    },
                    BackgroundColor: clay.Color{R: 0.95, G: 0.95, B: 0.95, A: 1.0},
                }).
                    Text(item, clay.TextConfig{FontSize: 16})
            }
        }()
}
```

These examples demonstrate the flexibility and power of Go-Clay for creating various UI layouts and interactions. The declarative API makes it easy to build complex interfaces while maintaining good performance.
