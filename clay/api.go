package clay

// #define CLAY_TEXT(text, textConfig) Clay__OpenTextElement(text, textConfig)
func CLAY_TEXT(text string, options ...TextOption) {
	textConfig := &Clay_TextElementConfig{}
	for _, option := range options {
		option(textConfig)
	}
	Clay__OpenTextElement(CLAY_STRING(text), textConfig)
}

func CLAY(elementID Clay_ElementId, elementDeclaration Clay_ElementDeclaration) {
	Clay__OpenElementWithId(elementID)
	Clay__ConfigureOpenElement(elementDeclaration)
}

func CLAY_STRING(label string) Clay_String {
	return Clay_String{
		IsStaticallyAllocated: true,
		Length:                int32(len(label)),
		Chars:                 []byte(label),
	}
}

// Note: If a compile error led you here, you might be trying to use CLAY_ID with something other than a string literal. To construct an ID with a dynamic string, use CLAY_SID instead.
func CLAY_ID(label string) Clay_ElementId {
	return CLAY_SID(CLAY_STRING(label))
}

func CLAY_SID(label Clay_String) Clay_ElementId {
	return Clay__HashString(label, 0)
}
