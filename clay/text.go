package clay

// Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
type Clay_TextElementConfigWrapMode uint8

const (
	// (default) breaks on whitespace characters.
	CLAY_TEXT_WRAP_WORDS Clay_TextElementConfigWrapMode = iota
	// Don't break on space characters, only on newlines.
	CLAY_TEXT_WRAP_NEWLINES
	// Disable text wrapping entirely.
	CLAY_TEXT_WRAP_NONE
)

// Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
type Clay_TextAlignment uint8

const (
	// (default) Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
	CLAY_TEXT_ALIGN_LEFT Clay_TextAlignment = iota
	// Horizontally aligns wrapped lines of text to the center of their bounding box.
	CLAY_TEXT_ALIGN_CENTER
	// Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
	CLAY_TEXT_ALIGN_RIGHT
)

// Controls various functionality related to text elements.
type Clay_TextElementConfig struct {
	// A pointer that will be transparently passed through to the resulting render command.
	UserData interface{}
	// The RGBA color of the font to render, conventionally specified as 0-255.
	TextColor Clay_Color
	// An integer transparently passed to Clay_MeasureText to identify the font to use.
	// The debug view will pass fontId = 0 for its internal text.
	FontId uint16
	// Controls the size of the font. Handled by the function provided to Clay_MeasureText.
	FontSize uint16
	// Controls extra horizontal spacing between characters. Handled by the function provided to Clay_MeasureText.
	LetterSpacing uint16
	// Controls additional vertical space between wrapped lines of text.
	LineHeight uint16
	// Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
	// CLAY_TEXT_WRAP_WORDS (default) breaks on whitespace characters.
	// CLAY_TEXT_WRAP_NEWLINES doesn't break on space characters, only on newlines.
	// CLAY_TEXT_WRAP_NONE disables wrapping entirely.
	WrapMode Clay_TextElementConfigWrapMode
	// Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
	// CLAY_TEXT_ALIGN_LEFT (default) - Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
	// CLAY_TEXT_ALIGN_CENTER - Horizontally aligns wrapped lines of text to the center of their bounding box.
	// CLAY_TEXT_ALIGN_RIGHT - Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
	TextAlignment Clay_TextAlignment
}

type Clay__MeasureTextCacheItem struct {
	UnwrappedDimensions     Clay_Dimensions
	MeasuredWordsStartIndex int32
	MinWidth                float32
	ContainsNewlines        bool
	// Hash map data
	Id         uint32
	NextIndex  int32
	Generation uint32
}

// Note: Clay_String is not guaranteed to be null terminated. It may be if created from a literal C string,
// but it is also used to represent slices.
type Clay_String struct {
	// Set this boolean to true if the char* data underlying this string will live for the entire lifetime of the program.
	// This will automatically be set for strings created with CLAY_STRING, as the macro requires a string literal.
	IsStaticallyAllocated bool
	Length                int32
	// The underlying character memory. Note: this will not be copied and will not extend the lifetime of the underlying memory.
	Chars []byte
}

func (s Clay_String) String() string {
	return string(s.Chars[:s.Length])
}

type Clay__WrappedTextLine struct {
	Line       Clay_String
	Dimensions Clay_Dimensions
}
type Clay__TextElementData struct {
	Text                Clay_String
	PreferredDimensions Clay_Dimensions
	ElementIndex        int32
	WrappedLines        Clay__Slice[Clay__WrappedTextLine]
}

type Clay__MeasuredWord struct {
	StartOffset int32
	Length      int32
	Width       float32
	Next        int32
}

func DefaultTextElementConfig() *Clay_TextElementConfig {
	return &Clay_TextElementConfig{
		TextColor:     Clay_Color{R: 0, G: 0, B: 0, A: 1},
		FontId:        0,
		FontSize:      16,
		LetterSpacing: 0,
		LineHeight:    16,
	}
}

type TextOption func(*Clay_TextElementConfig)

func TextWithColor(color Clay_Color) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.TextColor = color
	}
}

func TextWithFontId(fontId uint16) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.FontId = fontId
	}
}

func TextWithFontSize(fontSize uint16) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.FontSize = fontSize
	}
}

func TextWithLetterSpacing(letterSpacing uint16) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.LetterSpacing = letterSpacing
	}
}

func TextWithLineHeight(lineHeight uint16) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.LineHeight = lineHeight
	}
}

func TextWithWrapMode(wrapMode Clay_TextElementConfigWrapMode) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.WrapMode = wrapMode
	}
}

func TextWithTextAlignment(textAlignment Clay_TextAlignment) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.TextAlignment = textAlignment
	}
}

func TextWithUserData(userData interface{}) TextOption {
	return func(config *Clay_TextElementConfig) {
		config.UserData = userData
	}
}

func TextWithConfig(configFunc func(config *Clay_TextElementConfig)) TextOption {
	return func(config *Clay_TextElementConfig) {
		configFunc(config)
	}
}
