package views3

import (
	"iter"
	"maps"
	"slices"
)

type Slice[T any] interface {
	Len() int
	Get(i int) T
	All() iter.Seq2[int, T]
}

type Map[K comparable, V any] interface {
	Len() int
	Get(k K) (V, bool)
	All() iter.Seq2[K, V]
}

type Viewer[T any] interface {
	View() T
}

func NewMap[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	return mapImpl[M, K, V]{m: m}
}

type mapImpl[M ~map[K]V, K comparable, V any] struct {
	m M
}

func (v mapImpl[M, K, V]) Len() int {
	return len(v.m)
}

func (v mapImpl[M, K, V]) Get(k K) (V, bool) {
	val, ok := v.m[k]
	return val, ok
}

func (v mapImpl[M, K, V]) All() iter.Seq2[K, V] {
	return maps.All(v.m)
}

func NewViewerMap[M ~map[K]V, K comparable, V Viewer[V2], V2 any](m M) Map[K, V2] {
	return viewerMapImpl[M, K, V, V2]{m: m}
}

type viewerMapImpl[M ~map[K]V, K comparable, V Viewer[V2], V2 any] struct {
	m M
}

func (v viewerMapImpl[M, K, V, V2]) Len() int {
	return len(v.m)
}

func (v viewerMapImpl[M, K, V, V2]) Get(k K) (V2, bool) {
	val, ok := v.m[k]
	if !ok {
		var zero V2
		return zero, false
	}
	return val.View(), true
}

func (v viewerMapImpl[M, K, V, V2]) All() iter.Seq2[K, V2] {
	return func(yield func(K, V2) bool) {
		for k, v := range v.m {
			if !yield(k, v.View()) {
				return
			}
		}
	}
}

func NewSlice[S ~[]T, T any](s S) Slice[T] {
	return sliceImpl[S, T]{s: s}
}

type sliceImpl[S ~[]T, T any] struct {
	s S
}

func (v sliceImpl[S, T]) Len() int {
	return len(v.s)
}

func (v sliceImpl[S, T]) Get(i int) T {
	return v.s[i]
}

func (v sliceImpl[S, T]) All() iter.Seq2[int, T] {
	return slices.All(v.s)
}

func NewViewerSlice[S ~[]V, V Viewer[V2], V2 any](s S) Slice[V2] {
	return viewerSliceImpl[S, V, V2]{s: s}
}

type viewerSliceImpl[S ~[]V, V Viewer[V2], V2 any] struct {
	s S
}

func (v viewerSliceImpl[S, V, V2]) Len() int {
	return len(v.s)
}

func (v viewerSliceImpl[S, V, V2]) Get(i int) V2 {
	return v.s[i].View()
}

func (v viewerSliceImpl[S, V, V2]) All() iter.Seq2[int, V2] {
	return func(yield func(int, V2) bool) {
		for i, v := range v.s {
			if !yield(i, v.View()) {
				return
			}
		}
	}
}