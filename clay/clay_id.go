package clay

// Note: If a compile error led you here, you might be trying to use CLAY_ID with something other than a string literal. To construct an ID with a dynamic string, use CLAY_SID instead.
// #define CLAY_ID(label) CLAY_SID(CLAY_STRING(label))

// #define CLAY_SID(label) Clay__HashString(label, 0)

func CLAY_ID(label string) Clay_ElementId {
	return CLAY_SID(CLAY_STRING(label))
}

func CLAY_SID(label Clay_String) Clay_ElementId {
	return Clay__HashString(label)
}
