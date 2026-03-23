package views2

import (
	"iter"
	"maps"
)

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	Len() int
}

func WrapMap[K comparable, V any](m map[K]V) Map[K, V] {
	return wrappedMap[K, V](m)
}

func WrapViewerMap[K comparable, V1 Viewer[V2], V2 any](m map[K]V1) Map[K, V2] {
	return wrappedViewerMap[K, V1, V2](m)
}

type wrappedMap[K comparable, V any] map[K]V

func (m wrappedMap[K, V]) Get(key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

func (m wrappedMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m)
}

func (m wrappedMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m)
}

func (m wrappedMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(m)
}

func (m wrappedMap[K, V]) Len() int {
	return len(m)
}

type wrappedViewerMap[K comparable, V1 Viewer[V2], V2 any] map[K]V1

func (m wrappedViewerMap[K, V1, V2]) Get(key K) (V2, bool) {
	if v1, ok := m[key]; ok {
		return v1.View(), true
	}
	var zero V2
	return zero, false
}

func (m wrappedViewerMap[K, V1, V2]) All() iter.Seq2[K, V2] {
	return func(yield func(K, V2) bool) {
		for k, v1 := range m {
			if !yield(k, v1.View()) {
				return
			}
		}
	}
}

func (m wrappedViewerMap[K, V1, V2]) Keys() iter.Seq[K] {
	return maps.Keys(m)
}

func (m wrappedViewerMap[K, V1, V2]) Values() iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for _, v1 := range m {
			if !yield(v1.View()) {
				return
			}
		}
	}
}

func (m wrappedViewerMap[K, V1, V2]) Len() int {
	return len(m)
}

func MapClone[T Map[K, V], K comparable, V any](m T) map[K]V {
	result := make(map[K]V, m.Len())
	for k, v := range m.All() {
		result[k] = v
	}
	return result
}

func MapEqual[T Map[K, V], K comparable, V comparable](a, b T) bool {
	if a.Len() != b.Len() {
		return false
	}
	for k, vA := range a.All() {
		vB, ok := b.Get(k)
		if !ok || vA != vB {
			return false
		}
	}
	return true
}