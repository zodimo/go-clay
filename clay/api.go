package clay

type CLAY_CONTAINER_FUNC func(elementID Clay_ElementId, elementDeclaration Clay_ElementDeclaration, content ...CLAY_CONTAINER_FUNC)

type ClayContainer interface {
	Run()
}

var _ ClayContainer = (*claContainer)(nil)

type claContainer struct {
	wrapper func()
}

func (c *claContainer) Run() {
	c.wrapper()
}

func CLAY_ROOT(elementID Clay_ElementId, elementDeclaration Clay_ElementDeclaration, content ...ClayContainer) {
	CLAY(elementID, elementDeclaration, content...).Run()
}

func CLAY_ROOT_AUTO_ID(elementDeclaration Clay_ElementDeclaration, content ...ClayContainer) {
	CLAY_AUTO_ID(elementDeclaration, content...).Run()
}

func CLAY(elementID Clay_ElementId, elementDeclaration Clay_ElementDeclaration, content ...ClayContainer) ClayContainer {
	return &claContainer{
		wrapper: func() {
			Clay__OpenElementWithId(elementID)
			Clay__ConfigureOpenElement(elementDeclaration)
			for _, content := range content {
				content.Run()
			}
			Clay__CloseElement()
		},
	}
}

func CLAY_AUTO_ID(elementDeclaration Clay_ElementDeclaration, content ...ClayContainer) ClayContainer {
	return &claContainer{
		wrapper: func() {
			Clay__OpenElement()
			Clay__ConfigureOpenElement(elementDeclaration)
			for _, content := range content {
				content.Run()
			}
			Clay__CloseElement()
		},
	}
}

// #define CLAY_TEXT(text, textConfig) Clay__OpenTextElement(text, textConfig)

func CLAY_TEXT(text string, options ...TextOption) ClayContainer {
	textConfig := &Clay_TextElementConfig{}
	for _, option := range options {
		option(textConfig)
	}
	return &claContainer{
		wrapper: func() {
			Clay__OpenTextElement(CLAY_STRING(text), textConfig)
		},
	}
}

func CLAY_STRING(label string) Clay_String {
	return Clay_String{
		IsStaticallyAllocated: true,
		Length:                int32(len(label)),
		Chars:                 []byte(label),
	}
}

// Note: If a compile error led you here, you might be trying to use CLAY_ID with something other than a string literal.
// To construct an ID with a dynamic string, use CLAY_SID instead.
func CLAY_ID(label string) Clay_ElementId {
	return CLAY_SID(CLAY_STRING(label))
}

func CLAY_SID(label Clay_String) Clay_ElementId {
	return Clay__HashString(label, 0)
}

func CLAY_PADDING_ALL(padding uint16) Clay_Padding {
	return Clay_Padding{
		Left:   padding,
		Right:  padding,
		Top:    padding,
		Bottom: padding,
	}
}

// func CLAY_BORDER_OUTSIDE(widthValue uint16) Clay_BorderWidth {
// 	return Clay_BorderWidth{
// 		Left:            widthValue,
// 		Right:           widthValue,
// 		Top:             widthValue,
// 		Bottom:          widthValue,
// 		BetweenChildren: 0,
// 	}
// }

// #define CLAY_BORDER_ALL(widthValue) {widthValue, widthValue, widthValue, widthValue, widthValue}

func CLAY_CORNER_RADIUS(radius float32) Clay_CornerRadius {
	return Clay_CornerRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomLeft:  radius,
		BottomRight: radius,
	}
}

func CLAY_SIZING_FIT(minMax Clay_SizingMinMax) Clay_SizingAxis {
	return Clay_SizingAxis{
		Type: CLAY__SIZING_TYPE_FIT,
		Size: Clay_SizingAxisSize{MinMax: minMax},
	}
}

func CLAY_SIZING_GROW(minMax Clay_SizingMinMax) Clay_SizingAxis {
	return Clay_SizingAxis{
		Type: CLAY__SIZING_TYPE_GROW,
		Size: Clay_SizingAxisSize{MinMax: minMax},
	}
}
func CLAY_SIZING_FIXED(fixedSize float32) Clay_SizingAxis {
	return Clay_SizingAxis{
		Type: CLAY__SIZING_TYPE_FIXED,
		Size: Clay_SizingAxisSize{MinMax: Clay_SizingMinMax{Min: fixedSize, Max: fixedSize}},
	}
}

func CLAY_SIZING_PERCENT(percentOfParent float32) Clay_SizingAxis {
	if percentOfParent < 0 || percentOfParent > 1 {
		panic("percentOfParent must be between 0 and 1")
	}
	return Clay_SizingAxis{
		Type: CLAY__SIZING_TYPE_PERCENT,
		Size: Clay_SizingAxisSize{Percent: percentOfParent},
	}
}
