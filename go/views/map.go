package views

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

func MapGet[M ~map[K]V, K comparable, V any](m M, key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

func MapAll[M ~map[K]V, K comparable, V any](m M) iter.Seq2[K, V] {
	return maps.All(m)
}

func MapKeys[M ~map[K]V, K comparable, V any](m M) iter.Seq[K] {
	return maps.Keys(m)
}

func MapValues[M ~map[K]V, K comparable, V any](m M) iter.Seq[V] {
	return maps.Values(m)
}

func MapLen[M ~map[K]V, K comparable, V any](m M) int {
	return len(m)
}

func MapGetFunc[M ~map[K]BaseV, K comparable, BaseV, OutV any](m M, key K, f func(BaseV) OutV) (OutV, bool) {
	if v, ok := m[key]; ok {
		return f(v), true
	}
	var zero OutV
	return zero, false
}

func MapAllFunc[M ~map[K]BaseV, K comparable, BaseV, OutV any](m M, f func(BaseV) OutV) iter.Seq2[K, OutV] {
	return func(yield func(K, OutV) bool) {
		for k, v := range m {
			if !yield(k, f(v)) {
				return
			}
		}
	}
}

func MapValuesFunc[M ~map[K]BaseV, K comparable, BaseV, OutV any](m M, f func(BaseV) OutV) iter.Seq[OutV] {
	return func(yield func(OutV) bool) {
		for _, v := range m {
			if !yield(f(v)) {
				return
			}
		}
	}
}