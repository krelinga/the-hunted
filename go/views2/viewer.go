package views2

type Viewer[T any] interface {
	View() T
}

func Slice[V Viewer[O], O any](s []V) []O {
	return SliceTyped[[]O](s)
}

func SliceTyped[OS ~[]OE, OE any, IS ~[]IE, IE Viewer[OE]](s IS) OS {
	result := make(OS, len(s))
	for i, v := range s {
		result[i] = v.View()
	}
	return result
}

func Map[K comparable, V Viewer[O], O any](m map[K]V) map[K]O {
	return MapTyped[map[K]O](m)
}

func MapTyped[OM ~map[K]OV, OV any, IM ~map[K]IV, IV Viewer[OV], K comparable](m IM) OM {
	result := make(OM, len(m))
	for k, v := range m {
		result[k] = v.View()
	}
	return result
}