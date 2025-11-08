# Coding standards and patterns


# Optional and Default Arguments pattern.

Example of an Object called X, with optional options of BufferSize of type int

```golang

//  For Component X prefix with X
type XOptions struct {
	BufferSize int
}
type XOption func(options *XOptions)

func XWithBufferSize(size int) XOption {
	return func(options *XOptions) {
		options.BufferSize = size
	}
}
func defaultXOptions() XOptions{
    return XOptions{
		BufferSize: 0, 
	}
}

type X struct {
    options XOptions
}

func NewX(options ...XOption) (*X, error) {
	opts := defaultXOptions()

	for _, option := range options {
		option(&opts)
	}
    // Validate the Options
	if opts.BufferSize < 0  {
		return nil, fmt.Errorf("bufferSize cannot be a negative number")
	}
	
	return &X{
		options:             opts,
	}, nil
}


```