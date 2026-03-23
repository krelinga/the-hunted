package views2

import (
	"iter"
	"slices"
)

type Slice[T any] interface {
	Get(index int) T
	All() iter.Seq2[int, T]
	Values() iter.Seq[T]
	Len() int
}

func WrapSlice[T any](s []T) Slice[T] {
	return wrappedSlice[T](s)
}

func WrapViewerSlice[V1 Viewer[V2], V2 any](s []V1) Slice[V2] {
	return wrappedViewerSlice[V1, V2](s)
}

type wrappedSlice[T any] []T

func (s wrappedSlice[T]) Get(index int) T {
	return s[index]
}

func (s wrappedSlice[T]) All() iter.Seq2[int, T] {
	return slices.All(s)
}

func (s wrappedSlice[T]) Values() iter.Seq[T] {
	return slices.Values(s)
}

func (s wrappedSlice[T]) Len() int {
	return len(s)
}

type wrappedViewerSlice[V1 Viewer[V2], V2 any] []V1

func (s wrappedViewerSlice[V1, V2]) Get(index int) V2 {
	return s[index].View()
}

func (s wrappedViewerSlice[V1, V2]) All() iter.Seq2[int, V2] {
	return func(yield func(int, V2) bool) {
		for i, v1 := range s {
			if !yield(i, v1.View()) {
				return
			}
		}
	}
}

func (s wrappedViewerSlice[V1, V2]) Values() iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for _, v1 := range s {
			if !yield(v1.View()) {
				return
			}
		}
	}
}

func (s wrappedViewerSlice[V1, V2]) Len() int {
	return len(s)
}

func SliceClone[T Slice[V], V any](s T) []V {
	result := make([]V, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Get(i)
	}
	return result
}

func SliceEqual[T Slice[V], V comparable](a, b T) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if a.Get(i) != b.Get(i) {
			return false
		}
	}
	return true
}