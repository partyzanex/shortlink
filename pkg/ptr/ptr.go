package ptr

func Ptr[T any](v T) *T {
	return &v
}

func Value[T any](p *T) (v T) {
	if p != nil {
		v = *p
	}

	return v
}
