package clay

import "cmp"

func CLAY__MAX[T cmp.Ordered](x, y T) T {
	return max(x, y)
}
func CLAY__MIN[T cmp.Ordered](x, y T) T {
	return min(x, y)
}
func Clay__FloatEqual(x, y float32) bool {
	subtracted := x - y
	return subtracted < CLAY__EPSILON && subtracted > -CLAY__EPSILON
}
