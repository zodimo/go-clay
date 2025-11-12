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
