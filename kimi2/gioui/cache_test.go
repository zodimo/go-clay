package gioui

import (
	"image"
	"image/color"
	"testing"
)

func TestNewResourceCache(t *testing.T) {
	maxSize := 1024 * 1024 // 1MB
	cache := NewResourceCache(maxSize)

	if cache == nil {
		t.Fatal("NewResourceCache returned nil")
	}

	if cache.maxSize != maxSize {
		t.Errorf("Expected maxSize %d, got %d", maxSize, cache.maxSize)
	}

	if cache.currentSize != 0 {
		t.Errorf("Expected currentSize 0, got %d", cache.currentSize)
	}

	if cache.images == nil || cache.fonts == nil || cache.colors == nil || cache.gradients == nil {
		t.Error("Cache maps not initialized properly")
	}
}

func TestGetOrCreateColor(t *testing.T) {
	cache := NewResourceCache(1024 * 1024)

	// Test color conversion
	r, g, b, a := float32(0.5), float32(0.25), float32(0.75), float32(1.0)
	gioColor := cache.GetOrCreateColor(r, g, b, a)

	expectedR := uint8(r * 255)
	expectedG := uint8(g * 255)
	expectedB := uint8(b * 255)
	expectedA := uint8(a * 255)

	if gioColor.R != expectedR {
		t.Errorf("Expected R %d, got %d", expectedR, gioColor.R)
	}
	if gioColor.G != expectedG {
		t.Errorf("Expected G %d, got %d", expectedG, gioColor.G)
	}
	if gioColor.B != expectedB {
		t.Errorf("Expected B %d, got %d", expectedB, gioColor.B)
	}
	if gioColor.A != expectedA {
		t.Errorf("Expected A %d, got %d", expectedA, gioColor.A)
	}

	// Test caching - second call should return same color
	gioColor2 := cache.GetOrCreateColor(r, g, b, a)
	if gioColor != gioColor2 {
		t.Error("Color not properly cached")
	}

	// Test different color
	gioColor3 := cache.GetOrCreateColor(1.0, 0.0, 0.0, 1.0)
	if gioColor3.R != 255 || gioColor3.G != 0 || gioColor3.B != 0 || gioColor3.A != 255 {
		t.Error("Different color not converted correctly")
	}
}

func TestGetOrCreateImage(t *testing.T) {
	cache := NewResourceCache(1024 * 1024)

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	// Test image caching
	cachedImg, err := cache.GetOrCreateImage(img, FilterLinear)
	if err != nil {
		t.Fatalf("GetOrCreateImage failed: %v", err)
	}

	if cachedImg == nil {
		t.Fatal("GetOrCreateImage returned nil cached image")
	}

	if cachedImg.Image != img {
		t.Error("Cached image does not match original")
	}

	if cachedImg.UseCount != 1 {
		t.Errorf("Expected UseCount 1, got %d", cachedImg.UseCount)
	}

	// Test caching - second call should return same cached image
	cachedImg2, err := cache.GetOrCreateImage(img, FilterLinear)
	if err != nil {
		t.Fatalf("Second GetOrCreateImage failed: %v", err)
	}

	if cachedImg2 != cachedImg {
		t.Error("Image not properly cached")
	}

	if cachedImg2.UseCount != 2 {
		t.Errorf("Expected UseCount 2, got %d", cachedImg2.UseCount)
	}

	// Test different filter creates different cache entry
	cachedImg3, err := cache.GetOrCreateImage(img, FilterNearest)
	if err != nil {
		t.Fatalf("GetOrCreateImage with different filter failed: %v", err)
	}

	if cachedImg3 == cachedImg {
		t.Error("Different filter should create different cache entry")
	}
}

func TestGetOrCreateImageUnsupportedTypes(t *testing.T) {
	cache := NewResourceCache(1024 * 1024)

	// Test unsupported data types
	_, err := cache.GetOrCreateImage([]byte{1, 2, 3}, FilterLinear)
	if err == nil {
		t.Error("Expected error for byte slice data")
	}

	_, err = cache.GetOrCreateImage("test.png", FilterLinear)
	if err == nil {
		t.Error("Expected error for string data")
	}

	_, err = cache.GetOrCreateImage(123, FilterLinear)
	if err == nil {
		t.Error("Expected error for int data")
	}

	_, err = cache.GetOrCreateImage(nil, FilterLinear)
	if err == nil {
		t.Error("Expected error for nil data")
	}
}

func TestGetOrCreateGradient(t *testing.T) {
	cache := NewResourceCache(1024 * 1024)

	startColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	endColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}

	// Test gradient caching
	gradient := cache.GetOrCreateGradient(startColor, endColor, true)
	if gradient == nil {
		t.Fatal("GetOrCreateGradient returned nil")
	}

	if gradient.StartColor != startColor {
		t.Error("Gradient start color mismatch")
	}

	if gradient.EndColor != endColor {
		t.Error("Gradient end color mismatch")
	}

	if !gradient.Vertical {
		t.Error("Gradient orientation mismatch")
	}

	if gradient.UseCount != 1 {
		t.Errorf("Expected UseCount 1, got %d", gradient.UseCount)
	}

	// Test caching - second call should return same gradient
	gradient2 := cache.GetOrCreateGradient(startColor, endColor, true)
	if gradient2 != gradient {
		t.Error("Gradient not properly cached")
	}

	if gradient2.UseCount != 2 {
		t.Errorf("Expected UseCount 2, got %d", gradient2.UseCount)
	}

	// Test different orientation creates different cache entry
	gradient3 := cache.GetOrCreateGradient(startColor, endColor, false)
	if gradient3 == gradient {
		t.Error("Different orientation should create different cache entry")
	}
}

func TestGetOrCreateFont(t *testing.T) {
	cache := NewResourceCache(1024 * 1024)

	// Test font caching
	font1, err := cache.GetOrCreateFont(1)
	if err != nil {
		t.Fatalf("GetOrCreateFont failed: %v", err)
	}

	if font1 == nil {
		t.Fatal("GetOrCreateFont returned nil font")
	}

	// Test caching - second call should return same font
	font2, err := cache.GetOrCreateFont(1)
	if err != nil {
		t.Fatalf("Second GetOrCreateFont failed: %v", err)
	}

	if font2 != font1 {
		t.Error("Font not properly cached")
	}

	// Test different font ID creates different cache entry
	font3, err := cache.GetOrCreateFont(2)
	if err != nil {
		t.Fatalf("GetOrCreateFont with different ID failed: %v", err)
	}

	if font3 == font1 {
		t.Error("Different font ID should create different cache entry")
	}
}

func TestCacheEviction(t *testing.T) {
	// Create small cache to test eviction
	cache := NewResourceCache(1000) // 1000 bytes cache

	// Create different images that will fit individually but require eviction when both are cached
	// Each 15x15 RGBA image = 15 * 15 * 4 = 900 bytes (fits in 1000 byte cache)
	img1 := image.NewRGBA(image.Rect(0, 0, 15, 15)) // ~900 bytes
	img2 := image.NewRGBA(image.Rect(0, 0, 15, 15)) // ~900 bytes
	
	// Make the images different by setting different pixel values
	// This ensures they get different cache keys
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			img1.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255}) // Red
			img2.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Green
		}
	}

	// Add first image
	_, err := cache.GetOrCreateImage(img1, FilterLinear)
	if err != nil {
		t.Fatalf("Failed to cache first image: %v", err)
	}

	initialSize := cache.currentSize
	if initialSize == 0 {
		t.Error("Cache size should be greater than 0 after adding image")
	}

	// Verify first image is cached
	if cache.currentSize != 900 {
		t.Errorf("Expected cache size to be 900, got %d", cache.currentSize)
	}

	// Add second image (should trigger eviction of first image)
	cachedImg2, err := cache.GetOrCreateImage(img2, FilterLinear)
	if err != nil {
		t.Fatalf("Failed to cache second image: %v", err)
	}

	// Cache should not exceed maximum size after eviction
	if cache.currentSize > cache.maxSize {
		t.Errorf("Cache size (%d) exceeded maximum (%d) after eviction", cache.currentSize, cache.maxSize)
	}

	// Cache should contain only the second image after eviction
	if cache.currentSize != 900 {
		t.Errorf("Expected cache size to be 900 after eviction, got %d", cache.currentSize)
	}

	// Verify we can still access the second image
	if cachedImg2 == nil {
		t.Error("Second cached image should not be nil")
	}

	// Verify first image was evicted by trying to access it again
	// This should create a new cache entry (UseCount = 1)
	cachedImg1Again, err := cache.GetOrCreateImage(img1, FilterLinear)
	if err != nil {
		t.Fatalf("Failed to re-cache first image: %v", err)
	}
	
	if cachedImg1Again.UseCount != 1 {
		t.Errorf("Expected UseCount to be 1 for re-cached image, got %d", cachedImg1Again.UseCount)
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewResourceCache(1024 * 1024)

	// Add some items to cache
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	cache.GetOrCreateImage(img, FilterLinear)
	cache.GetOrCreateColor(1.0, 0.0, 0.0, 1.0)
	cache.GetOrCreateFont(1)

	startColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	endColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	cache.GetOrCreateGradient(startColor, endColor, true)

	// Verify items are in cache
	stats := cache.GetStats()
	if stats.ImageCount == 0 || stats.ColorCount == 0 || stats.FontCount == 0 || stats.GradientCount == 0 {
		t.Error("Cache should contain items before clear")
	}

	// Clear cache
	cache.Clear()

	// Verify cache is empty
	stats = cache.GetStats()
	if stats.ImageCount != 0 || stats.ColorCount != 0 || stats.FontCount != 0 || stats.GradientCount != 0 {
		t.Error("Cache should be empty after clear")
	}

	if cache.currentSize != 0 {
		t.Errorf("Cache size should be 0 after clear, got %d", cache.currentSize)
	}
}

func TestCacheStats(t *testing.T) {
	cache := NewResourceCache(1024 * 1024)

	// Initially empty
	stats := cache.GetStats()
	if stats.ImageCount != 0 || stats.ColorCount != 0 || stats.FontCount != 0 || stats.GradientCount != 0 {
		t.Error("New cache should have zero counts")
	}

	if stats.MaxSize != 1024*1024 {
		t.Errorf("Expected MaxSize %d, got %d", 1024*1024, stats.MaxSize)
	}

	// Add items
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	cache.GetOrCreateImage(img, FilterLinear)
	cache.GetOrCreateColor(1.0, 0.0, 0.0, 1.0)
	cache.GetOrCreateFont(1)

	startColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	endColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	cache.GetOrCreateGradient(startColor, endColor, true)

	// Check updated stats
	stats = cache.GetStats()
	if stats.ImageCount != 1 {
		t.Errorf("Expected ImageCount 1, got %d", stats.ImageCount)
	}
	if stats.ColorCount != 1 {
		t.Errorf("Expected ColorCount 1, got %d", stats.ColorCount)
	}
	if stats.FontCount != 1 {
		t.Errorf("Expected FontCount 1, got %d", stats.FontCount)
	}
	if stats.GradientCount != 1 {
		t.Errorf("Expected GradientCount 1, got %d", stats.GradientCount)
	}

	if stats.CurrentSize == 0 {
		t.Error("CurrentSize should be greater than 0")
	}

	// Test string representation
	statsStr := stats.String()
	if statsStr == "" {
		t.Error("Stats string representation should not be empty")
	}
}

func BenchmarkGetOrCreateColor(b *testing.B) {
	cache := NewResourceCache(1024 * 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetOrCreateColor(0.5, 0.5, 0.5, 1.0)
	}
}

func BenchmarkGetOrCreateImage(b *testing.B) {
	cache := NewResourceCache(1024 * 1024)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetOrCreateImage(img, FilterLinear)
	}
}

func BenchmarkCacheEviction(b *testing.B) {
	cache := NewResourceCache(1000) // Small cache to force eviction

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))
		cache.GetOrCreateImage(img, FilterLinear)
	}
}
