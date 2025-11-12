package clay

type Clay_BooleanWarnings struct {
	MaxElementsExceeded           bool
	MaxRenderCommandsExceeded     bool
	MaxTextMeasureCacheExceeded   bool
	TextMeasurementFunctionNotSet bool
}

type Clay__Warning struct {
	BaseMessage    Clay_String
	DynamicMessage Clay_String
}
