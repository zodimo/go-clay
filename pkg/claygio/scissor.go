package claygio

import (
	"encoding/json"
	"fmt"
)

type ClippingContainer struct {
	X      float32
	Y      float32
	Width  float32
	Height float32

	Horizontal bool
	Vertical   bool
}

func (c *ClippingContainer) String() string {
	jsonString, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling clipping container: %v", err)
	}
	return string(jsonString)
}
