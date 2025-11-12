package clay

var QueryScrollOffsetFunction Clay__QueryScrollOffsetFunction = nil
var MeasureTextFunction Clay__MeasureTextFunction = nil

func Clay__QueryScrollOffset(elementId uint32, userData interface{}) Clay_Vector2 {
	if QueryScrollOffsetFunction == nil {
		panic("QueryScrollOffsetFunction is not set")
	}
	return QueryScrollOffsetFunction(elementId, userData)
}

func Clay__MeasureText(text Clay_StringSlice, config *Clay_TextElementConfig, userData interface{}) Clay_Dimensions {
	if MeasureTextFunction == nil {
		panic("MeasureTextFunction is not set")
	}
	return MeasureTextFunction(text, config, userData)
}
