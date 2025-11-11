package clay

// CLAY__DEFAULT_STRUCT = {0}
// typeName typeName##_DEFAULT = CLAY__DEFAULT_STRUCT;                                                             \

var CLAY__STRING_DEFAULT = Clay_String{Length: 0, Chars: ""}
var Clay_LayoutConfig_DEFAULT = Clay_LayoutConfig{}

var Clay_LayoutElementHashMapItem_DEFAULT = Clay_LayoutElementHashMapItem{}

var Clay__ErrorHandlerFunctionDefault = Clay_ErrorHandler{
	ErrorHandlerFunction: func(errorData Clay_ErrorData) {
		// Do nothing
	},
	UserData: nil,
}
