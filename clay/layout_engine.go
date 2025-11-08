package clay

import (
	"sync"
	"time"
)

// layoutEngine implements the LayoutEngine interface
type layoutEngine struct {
	// Core state
	arena          *Arena
	dimensions     Dimensions
	pointerState   Vector2
	pointerPressed bool
	scrollOffset   Vector2
	debugMode      bool

	// Layout state
	layoutActive   bool
	elementStack   []ElementID
	elementCount   int
	renderCommands []RenderCommand

	// Element storage
	elements      map[ElementID]*Element
	elementBounds map[ElementID]BoundingBox

	// Performance tracking
	startTime time.Time
	stats     LayoutStats

	// Text measurement
	textMeasurer TextMeasurer

	// ID generation
	nextID  uint32
	idMutex sync.Mutex
}

// Element represents a layout element
type Element struct {
	ID           ElementID
	Config       ElementDeclaration
	Parent       ElementID
	Children     []ElementID
	Bounds       BoundingBox
	ComputedSize Dimensions
	ZIndex       int16
	IsFloating   bool
	IsClipping   bool
}

// NewLayoutEngine creates a new layout engine
func NewLayoutEngine() LayoutEngine {
	return &layoutEngine{
		arena:          NewArena(1024 * 1024), // 1MB default
		elements:       make(map[ElementID]*Element),
		elementBounds:  make(map[ElementID]BoundingBox),
		renderCommands: make([]RenderCommand, 0, 1000),
	}
}

// NewLayoutEngineWithArena creates a new layout engine with a custom arena
func NewLayoutEngineWithArena(arena *Arena) LayoutEngine {
	return &layoutEngine{
		arena:          arena,
		elements:       make(map[ElementID]*Element),
		elementBounds:  make(map[ElementID]BoundingBox),
		renderCommands: make([]RenderCommand, 0, 1000),
	}
}

// BeginLayout starts a new layout computation
func (e *layoutEngine) BeginLayout() {
	if e.layoutActive {
		panic("Layout already active. Call EndLayout() first.")
	}

	e.layoutActive = true
	e.startTime = time.Now()
	e.elementStack = e.elementStack[:0]
	e.renderCommands = e.renderCommands[:0]
	e.elementCount = 0

	// Reset arena
	e.arena.Reset()

	// Clear element storage
	for id := range e.elements {
		delete(e.elements, id)
	}
	for id := range e.elementBounds {
		delete(e.elementBounds, id)
	}

	setCurrentEngine(e)
}

// EndLayout completes the layout computation and returns render commands
func (e *layoutEngine) EndLayout() []RenderCommand {
	if !e.layoutActive {
		panic("No active layout. Call BeginLayout() first.")
	}

	e.layoutActive = false

	// Perform two-pass layout algorithm
	e.computeLayout()
	e.generateRenderCommands()

	// Update stats
	e.stats.LayoutTime = time.Since(e.startTime).Nanoseconds()
	e.stats.ElementCount = e.elementCount
	e.stats.RenderCommands = len(e.renderCommands)
	e.stats.MemoryUsed = e.arena.Used()

	setCurrentEngine(nil)

	// Return copy of render commands
	commands := make([]RenderCommand, len(e.renderCommands))
	copy(commands, e.renderCommands)
	return commands
}

// OpenElement opens a new element
func (e *layoutEngine) OpenElement(id ElementID, config ElementDeclaration) {
	if !e.layoutActive {
		panic("No active layout. Call BeginLayout() first.")
	}

	// Create element
	element := &Element{
		ID:     id,
		Config: config,
		ZIndex: 0,
	}

	// Set parent if we have an active element
	if len(e.elementStack) > 0 {
		parentID := e.elementStack[len(e.elementStack)-1]
		element.Parent = parentID
		if parent, exists := e.elements[parentID]; exists {
			parent.Children = append(parent.Children, id)
		}
	}

	// Store element
	e.elements[id] = element
	e.elementStack = append(e.elementStack, id)
	e.elementCount++
}

// CloseElement closes the current element
func (e *layoutEngine) CloseElement() {
	if !e.layoutActive {
		panic("No active layout. Call BeginLayout() first.")
	}

	if len(e.elementStack) == 0 {
		panic("No open element to close.")
	}

	// Remove from stack
	e.elementStack = e.elementStack[:len(e.elementStack)-1]
}

// SetPointerState sets the pointer position and state
func (e *layoutEngine) SetPointerState(pos Vector2, pressed bool) {
	e.pointerState = pos
	e.pointerPressed = pressed
}

// SetLayoutDimensions sets the layout dimensions
func (e *layoutEngine) SetLayoutDimensions(dimensions Dimensions) {
	e.dimensions = dimensions
}

// SetScrollOffset sets the scroll offset
func (e *layoutEngine) SetScrollOffset(offset Vector2) {
	e.scrollOffset = offset
}

// GetElementBounds returns the bounds of an element
func (e *layoutEngine) GetElementBounds(id ElementID) (BoundingBox, bool) {
	bounds, exists := e.elementBounds[id]
	return bounds, exists
}

// IsPointerOver checks if the pointer is over an element
func (e *layoutEngine) IsPointerOver(id ElementID) bool {
	bounds, exists := e.GetElementBounds(id)
	if !exists {
		return false
	}

	return e.pointerState.X >= bounds.X &&
		e.pointerState.X <= bounds.X+bounds.Width &&
		e.pointerState.Y >= bounds.Y &&
		e.pointerState.Y <= bounds.Y+bounds.Height
}

// GetScrollOffset returns the scroll offset for an element
func (e *layoutEngine) GetScrollOffset(id ElementID) Vector2 {
	// For now, return global scroll offset
	// TODO: Implement per-element scroll tracking
	return e.scrollOffset
}

// SetDebugMode enables or disables debug mode
func (e *layoutEngine) SetDebugMode(enabled bool) {
	e.debugMode = enabled
}

// GetStats returns layout statistics
func (e *layoutEngine) GetStats() LayoutStats {
	return e.stats
}

// generateID generates a unique element ID
func (e *layoutEngine) generateID() ElementID {
	e.idMutex.Lock()
	defer e.idMutex.Unlock()
	e.nextID++
	return ElementID(e.nextID)
}

// computeLayout performs the two-pass layout algorithm
func (e *layoutEngine) computeLayout() {
	// Pass 1: Calculate sizes
	e.calculateSizes()

	// Pass 2: Calculate positions
	e.calculatePositions()
}

// calculateSizes performs the first pass of layout computation
func (e *layoutEngine) calculateSizes() {
	// Find root elements (elements with no parent)
	var rootElements []ElementID
	for id, element := range e.elements {
		if element.Parent == 0 {
			rootElements = append(rootElements, id)
		}
	}

	// Calculate sizes for each root element
	for _, rootID := range rootElements {
		e.calculateElementSize(rootID)
	}
}

// calculateElementSize calculates the size of an element and its children
func (e *layoutEngine) calculateElementSize(id ElementID) {
	element, exists := e.elements[id]
	if !exists {
		return
	}

	// Calculate size based on content and constraints
	size := e.determineElementSize(element)
	element.ComputedSize = size

	// Recursively calculate children sizes
	for _, childID := range element.Children {
		e.calculateElementSize(childID)
	}
}

// determineElementSize determines the size of an element
func (e *layoutEngine) determineElementSize(element *Element) Dimensions {
	_ = element.Config.Layout // TODO: Use this

	// Handle text elements
	if element.Config.Text != nil {
		return e.calculateTextSize(element)
	}

	// Handle image elements
	if element.Config.Image != nil {
		return e.calculateImageSize(element)
	}

	// Handle container elements
	return e.calculateContainerSize(element)
}

// calculateTextSize calculates the size of a text element
func (e *layoutEngine) calculateTextSize(element *Element) Dimensions {
	if e.textMeasurer == nil {
		// Fallback to simple calculation
		text := "Sample text" // TODO: Get actual text content
		config := *element.Config.Text
		return Dimensions{
			Width:  float32(len(text)) * config.FontSize * 0.6, // Rough estimate
			Height: config.FontSize * config.LineHeight,
		}
	}

	// Use text measurer
	text := "Sample text" // TODO: Get actual text content
	return e.textMeasurer.MeasureText(text, *element.Config.Text)
}

// calculateImageSize calculates the size of an image element
func (e *layoutEngine) calculateImageSize(element *Element) Dimensions {
	// For now, return a default size
	// TODO: Implement actual image size calculation
	return Dimensions{Width: 100, Height: 100}
}

// calculateContainerSize calculates the size of a container element
func (e *layoutEngine) calculateContainerSize(element *Element) Dimensions {
	config := element.Config.Layout

	// Calculate size based on sizing configuration
	width := e.calculateAxisSize(config.Sizing.Width, e.dimensions.Width)
	height := e.calculateAxisSize(config.Sizing.Height, e.dimensions.Height)

	// If children exist, calculate based on children
	if len(element.Children) > 0 {
		childrenSize := e.calculateChildrenSize(element)

		// Apply sizing constraints
		if config.Sizing.Width.Type == SizingFitToContent {
			width = childrenSize.Width + float32(config.Padding.Left+config.Padding.Right)
		}
		if config.Sizing.Height.Type == SizingFitToContent {
			height = childrenSize.Height + float32(config.Padding.Top+config.Padding.Bottom)
		}
	}

	return Dimensions{Width: width, Height: height}
}

// calculateAxisSize calculates the size along one axis
func (e *layoutEngine) calculateAxisSize(axis SizingAxis, parentSize float32) float32 {
	switch axis.Type {
	case SizingFitToContent:
		return 0 // Will be calculated based on content
	case SizingGrowToFillAvailableSpace:
		return parentSize * 0.5 // Default grow behavior
	case SizingPercentOfParent:
		return parentSize * axis.Percent
	case SizingFixedPixelSize:
		return axis.Min
	default:
		return 0
	}
}

// calculateChildrenSize calculates the total size of children
func (e *layoutEngine) calculateChildrenSize(element *Element) Dimensions {
	if len(element.Children) == 0 {
		return Dimensions{}
	}

	config := element.Config.Layout
	var totalWidth, totalHeight float32
	var maxWidth, maxHeight float32

	for _, childID := range element.Children {
		child, exists := e.elements[childID]
		if !exists {
			continue
		}

		childSize := child.ComputedSize

		if config.Direction == LeftToRight {
			totalWidth += childSize.Width
			if childSize.Height > maxHeight {
				maxHeight = childSize.Height
			}
		} else {
			totalHeight += childSize.Height
			if childSize.Width > maxWidth {
				maxWidth = childSize.Width
			}
		}
	}

	// Add gaps between children
	if config.Direction == LeftToRight {
		totalWidth += float32(len(element.Children)-1) * config.ChildGap
	} else {
		totalHeight += float32(len(element.Children)-1) * config.ChildGap
	}

	return Dimensions{
		Width:  totalWidth + maxWidth,
		Height: totalHeight + maxHeight,
	}
}

// calculatePositions performs the second pass of layout computation
func (e *layoutEngine) calculatePositions() {
	// Find root elements
	var rootElements []ElementID
	for id, element := range e.elements {
		if element.Parent == 0 {
			rootElements = append(rootElements, id)
		}
	}

	// Calculate positions for each root element
	for _, rootID := range rootElements {
		e.calculateElementPosition(rootID, Vector2{X: 0, Y: 0})
	}
}

// calculateElementPosition calculates the position of an element and its children
func (e *layoutEngine) calculateElementPosition(id ElementID, position Vector2) {
	element, exists := e.elements[id]
	if !exists {
		return
	}

	// Set element position
	element.Bounds = BoundingBox{
		X:      position.X,
		Y:      position.Y,
		Width:  element.ComputedSize.Width,
		Height: element.ComputedSize.Height,
	}

	// Store bounds for queries
	e.elementBounds[id] = element.Bounds

	// Calculate children positions
	if len(element.Children) > 0 {
		e.calculateChildrenPositions(element)
	}
}

// calculateChildrenPositions calculates positions for all children
func (e *layoutEngine) calculateChildrenPositions(element *Element) {
	config := element.Config.Layout

	// Start position (accounting for padding)
	startX := element.Bounds.X + float32(config.Padding.Left)
	startY := element.Bounds.Y + float32(config.Padding.Top)

	currentX := startX
	currentY := startY

	for i, childID := range element.Children {
		child, exists := e.elements[childID]
		if !exists {
			continue
		}

		// Calculate child position
		childPos := Vector2{X: currentX, Y: currentY}

		// Apply alignment
		childPos = e.applyAlignment(childPos, element, child, i)

		// Set child position
		child.Bounds = BoundingBox{
			X:      childPos.X,
			Y:      childPos.Y,
			Width:  child.ComputedSize.Width,
			Height: child.ComputedSize.Height,
		}

		// Store child bounds
		e.elementBounds[childID] = child.Bounds

		// Recursively calculate child's children
		e.calculateElementPosition(childID, childPos)

		// Advance position for next child
		if config.Direction == LeftToRight {
			currentX += child.ComputedSize.Width + config.ChildGap
		} else {
			currentY += child.ComputedSize.Height + config.ChildGap
		}
	}
}

// applyAlignment applies alignment to a child position
func (e *layoutEngine) applyAlignment(pos Vector2, parent *Element, child *Element, childIndex int) Vector2 {
	config := parent.Config.Layout

	// Calculate available space
	availableWidth := parent.Bounds.Width - float32(config.Padding.Left+config.Padding.Right)
	availableHeight := parent.Bounds.Height - float32(config.Padding.Top+config.Padding.Bottom)

	// Apply horizontal alignment
	switch config.ChildAlignment.X {
	case AlignXCenter:
		pos.X = parent.Bounds.X + (availableWidth-child.ComputedSize.Width)/2 + float32(config.Padding.Left)
	case AlignXRight:
		pos.X = parent.Bounds.X + availableWidth - child.ComputedSize.Width + float32(config.Padding.Left)
	}

	// Apply vertical alignment
	switch config.ChildAlignment.Y {
	case AlignYCenter:
		pos.Y = parent.Bounds.Y + (availableHeight-child.ComputedSize.Height)/2 + float32(config.Padding.Top)
	case AlignYBottom:
		pos.Y = parent.Bounds.Y + availableHeight - child.ComputedSize.Height + float32(config.Padding.Top)
	}

	return pos
}

// generateRenderCommands generates render commands for all elements
func (e *layoutEngine) generateRenderCommands() {
	// Clear previous commands
	e.renderCommands = e.renderCommands[:0]

	// Generate commands for all elements
	for id, element := range e.elements {
		e.generateElementCommands(id, element)
	}

	// Sort commands by Z-index
	e.sortRenderCommands()
}

// generateElementCommands generates render commands for an element
func (e *layoutEngine) generateElementCommands(id ElementID, element *Element) {
	bounds := element.Bounds

	// Generate background rectangle
	if element.Config.BackgroundColor.A > 0 {
		e.renderCommands = append(e.renderCommands, RenderCommand{
			BoundingBox: bounds,
			CommandType: CommandRectangle,
			ZIndex:      element.ZIndex,
			ID:          id,
			Data: RectangleCommand{
				Color:        element.Config.BackgroundColor,
				CornerRadius: element.Config.CornerRadius,
			},
		})
	}

	// Generate border
	if element.Config.Border != nil {
		e.renderCommands = append(e.renderCommands, RenderCommand{
			BoundingBox: bounds,
			CommandType: CommandBorder,
			ZIndex:      element.ZIndex,
			ID:          id,
			Data: BorderCommand{
				Color:        element.Config.Border.Color,
				Width:        element.Config.Border.Width,
				CornerRadius: element.Config.CornerRadius,
			},
		})
	}

	// Generate text
	if element.Config.Text != nil {
		e.renderCommands = append(e.renderCommands, RenderCommand{
			BoundingBox: bounds,
			CommandType: CommandText,
			ZIndex:      element.ZIndex,
			ID:          id,
			Data: TextCommand{
				Text:          "Sample text", // TODO: Get actual text content
				FontID:        element.Config.Text.FontID,
				FontSize:      element.Config.Text.FontSize,
				Color:         element.Config.Text.Color,
				LineHeight:    element.Config.Text.LineHeight,
				LetterSpacing: element.Config.Text.LetterSpacing,
				Alignment:     element.Config.Text.Alignment,
			},
		})
	}

	// Generate image
	if element.Config.Image != nil {
		e.renderCommands = append(e.renderCommands, RenderCommand{
			BoundingBox: bounds,
			CommandType: CommandImage,
			ZIndex:      element.ZIndex,
			ID:          id,
			Data: ImageCommand{
				ImageData:    element.Config.Image.ImageData,
				TintColor:    element.Config.Image.TintColor,
				CornerRadius: element.Config.CornerRadius,
			},
		})
	}

	// Generate clipping commands
	if element.Config.Clip != nil {
		e.renderCommands = append(e.renderCommands, RenderCommand{
			BoundingBox: bounds,
			CommandType: CommandClipStart,
			ZIndex:      element.ZIndex,
			ID:          id,
			Data: ClipStartCommand{
				Horizontal: element.Config.Clip.Horizontal,
				Vertical:   element.Config.Clip.Vertical,
			},
		})
	}
}

// sortRenderCommands sorts render commands by Z-index
func (e *layoutEngine) sortRenderCommands() {
	// Simple bubble sort for now
	// TODO: Implement more efficient sorting
	for i := 0; i < len(e.renderCommands)-1; i++ {
		for j := 0; j < len(e.renderCommands)-i-1; j++ {
			if e.renderCommands[j].ZIndex > e.renderCommands[j+1].ZIndex {
				e.renderCommands[j], e.renderCommands[j+1] = e.renderCommands[j+1], e.renderCommands[j]
			}
		}
	}
}
