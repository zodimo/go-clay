package gioui

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"sync"

	"gioui.org/op/paint"
	"gioui.org/font"
)

// ResourceCache manages cached resources for performance optimization
type ResourceCache struct {
	images    map[string]*CachedImage
	fonts     map[string]*font.Font
	colors    map[string]color.NRGBA
	gradients map[string]*CachedGradient
	mutex     sync.RWMutex
	maxSize   int
	currentSize int
}

// CachedImage represents a cached image with metadata
type CachedImage struct {
	Image     image.Image
	ImageOp   paint.ImageOp
	Size      int
	LastUsed  int64
	UseCount  int
}

// CachedGradient represents a cached gradient operation
type CachedGradient struct {
	StartColor color.NRGBA
	EndColor   color.NRGBA
	Vertical   bool
	LastUsed   int64
	UseCount   int
}

// NewResourceCache creates a new resource cache with specified maximum size in bytes
func NewResourceCache(maxSizeBytes int) *ResourceCache {
	return &ResourceCache{
		images:    make(map[string]*CachedImage),
		fonts:     make(map[string]*font.Font),
		colors:    make(map[string]color.NRGBA),
		gradients: make(map[string]*CachedGradient),
		maxSize:   maxSizeBytes,
	}
}

// generateImageKey creates a unique key for image caching
func (rc *ResourceCache) generateImageKey(data interface{}, filter ImageFilter) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v_%d", data, filter)))
	return fmt.Sprintf("img_%x", hash.Sum(nil)[:8])
}

// generateColorKey creates a unique key for color caching
func (rc *ResourceCache) generateColorKey(r, g, b, a float32) string {
	return fmt.Sprintf("color_%.3f_%.3f_%.3f_%.3f", r, g, b, a)
}

// generateGradientKey creates a unique key for gradient caching
func (rc *ResourceCache) generateGradientKey(startColor, endColor color.NRGBA, vertical bool) string {
	return fmt.Sprintf("grad_%d_%d_%d_%d_%d_%d_%d_%d_%t",
		startColor.R, startColor.G, startColor.B, startColor.A,
		endColor.R, endColor.G, endColor.B, endColor.A, vertical)
}

// GetOrCreateImage retrieves or creates a cached image
func (rc *ResourceCache) GetOrCreateImage(data interface{}, filter ImageFilter) (*CachedImage, error) {
	key := rc.generateImageKey(data, filter)
	
	rc.mutex.RLock()
	if cached, exists := rc.images[key]; exists {
		cached.UseCount++
		rc.mutex.RUnlock()
		return cached, nil
	}
	rc.mutex.RUnlock()

	// Convert data to image.Image
	var img image.Image
	switch v := data.(type) {
	case image.Image:
		img = v
	case []byte:
		// Handle byte data - would need image decoding
		return nil, fmt.Errorf("byte slice image data not yet supported")
	case string:
		// Handle file path - would need image loading
		return nil, fmt.Errorf("file path image data not yet supported")
	default:
		return nil, fmt.Errorf("unsupported image data type: %T", data)
	}

	// Calculate image size for cache management
	bounds := img.Bounds()
	imageSize := bounds.Dx() * bounds.Dy() * 4 // Assume RGBA

	// Create image operation
	imageOp := paint.NewImageOp(img)
	switch filter {
	case FilterLinear:
		imageOp.Filter = paint.FilterLinear
	case FilterNearest:
		imageOp.Filter = paint.FilterNearest
	}

	cached := &CachedImage{
		Image:    img,
		ImageOp:  imageOp,
		Size:     imageSize,
		UseCount: 1,
	}

	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Check if we need to evict items to make space
	if rc.currentSize+imageSize > rc.maxSize {
		rc.evictLeastUsed(imageSize)
	}

	rc.images[key] = cached
	rc.currentSize += imageSize

	return cached, nil
}

// GetOrCreateColor retrieves or creates a cached color conversion
func (rc *ResourceCache) GetOrCreateColor(r, g, b, a float32) color.NRGBA {
	key := rc.generateColorKey(r, g, b, a)
	
	rc.mutex.RLock()
	if cached, exists := rc.colors[key]; exists {
		rc.mutex.RUnlock()
		return cached
	}
	rc.mutex.RUnlock()

	// Convert Clay color (0-1) to Gio color (0-255)
	gioColor := color.NRGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: uint8(a * 255),
	}

	rc.mutex.Lock()
	rc.colors[key] = gioColor
	rc.mutex.Unlock()

	return gioColor
}

// GetOrCreateGradient retrieves or creates a cached gradient
func (rc *ResourceCache) GetOrCreateGradient(startColor, endColor color.NRGBA, vertical bool) *CachedGradient {
	key := rc.generateGradientKey(startColor, endColor, vertical)
	
	rc.mutex.RLock()
	if cached, exists := rc.gradients[key]; exists {
		cached.UseCount++
		rc.mutex.RUnlock()
		return cached
	}
	rc.mutex.RUnlock()

	cached := &CachedGradient{
		StartColor: startColor,
		EndColor:   endColor,
		Vertical:   vertical,
		UseCount:   1,
	}

	rc.mutex.Lock()
	rc.gradients[key] = cached
	rc.mutex.Unlock()

	return cached
}

// GetOrCreateFont retrieves or creates a cached font
func (rc *ResourceCache) GetOrCreateFont(fontID uint16) (*font.Font, error) {
	key := fmt.Sprintf("font_%d", fontID)
	
	rc.mutex.RLock()
	if cached, exists := rc.fonts[key]; exists {
		rc.mutex.RUnlock()
		return cached, nil
	}
	rc.mutex.RUnlock()

	// For now, return a default font - this would be enhanced to load actual fonts
	fontObj := &font.Font{}

	rc.mutex.Lock()
	rc.fonts[key] = fontObj
	rc.mutex.Unlock()

	return fontObj, nil
}

// evictLeastUsed removes least used items to make space
func (rc *ResourceCache) evictLeastUsed(neededSpace int) {
	// Simple LRU eviction - remove items with lowest use count
	for rc.currentSize+neededSpace > rc.maxSize && len(rc.images) > 0 {
		var lruKey string
		var lruUseCount int = int(^uint(0) >> 1) // Max int

		for key, cached := range rc.images {
			if cached.UseCount < lruUseCount {
				lruUseCount = cached.UseCount
				lruKey = key
			}
		}

		if lruKey != "" {
			cached := rc.images[lruKey]
			rc.currentSize -= cached.Size
			delete(rc.images, lruKey)
		} else {
			break // No more items to evict
		}
	}
}

// Clear removes all cached items
func (rc *ResourceCache) Clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.images = make(map[string]*CachedImage)
	rc.fonts = make(map[string]*font.Font)
	rc.colors = make(map[string]color.NRGBA)
	rc.gradients = make(map[string]*CachedGradient)
	rc.currentSize = 0
}

// GetStats returns cache statistics
func (rc *ResourceCache) GetStats() CacheStats {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	return CacheStats{
		ImageCount:    len(rc.images),
		FontCount:     len(rc.fonts),
		ColorCount:    len(rc.colors),
		GradientCount: len(rc.gradients),
		CurrentSize:   rc.currentSize,
		MaxSize:       rc.maxSize,
	}
}

// CacheStats provides cache usage statistics
type CacheStats struct {
	ImageCount    int
	FontCount     int
	ColorCount    int
	GradientCount int
	CurrentSize   int
	MaxSize       int
}

// String returns a string representation of cache stats
func (cs CacheStats) String() string {
	return fmt.Sprintf("Cache Stats: Images=%d, Fonts=%d, Colors=%d, Gradients=%d, Size=%d/%d bytes",
		cs.ImageCount, cs.FontCount, cs.ColorCount, cs.GradientCount, cs.CurrentSize, cs.MaxSize)
}
