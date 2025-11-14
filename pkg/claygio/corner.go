package claygio

import "github.com/zodimo/clay-go/clay"

type CornerShapes struct {
	TopStart    CornerShape
	TopEnd      CornerShape
	BottomStart CornerShape
	BottomEnd   CornerShape
}

// Corner shape types for rounded rectangles
type CornerKind int

const (
	CornerKindDefault CornerKind = iota
	CornerKindChamfer
	CornerKindRound
)

type CornerShape struct {
	Kind        CornerKind
	Size        float32
	AdaptToSize bool
}

// isCornerRadiusZero checks if all corner radius values are zero
func IsCornerRadiusZero(radius clay.Clay_CornerRadius) bool {
	return radius.TopLeft == 0 && radius.TopRight == 0 &&
		radius.BottomLeft == 0 && radius.BottomRight == 0
}
