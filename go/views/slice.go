package views

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

func SliceGet[S ~[]T, T any](s S, index int) T {
	return s[index]
}

func SliceAll[S ~[]T, T any](s S) iter.Seq2[int, T] {
	return slices.All(s)
}

func SliceValues[S ~[]T, T any](s S) iter.Seq[T] {
	return slices.Values(s)
}

func SliceLen[S ~[]T, T any](s S) int {
	return len(s)
}

func SliceGetFunc[S ~[]BaseT, BaseT, OutT any](s S, index int, f func(BaseT) OutT) OutT {
	return f(s[index])
}

func SliceAllFunc[S ~[]BaseT, BaseT, OutT any](s S, f func(BaseT) OutT) iter.Seq2[int, OutT] {
	return func(yield func(int, OutT) bool) {
		for i, v := range s {
			if !yield(i, f(v)) {
				return
			}
		}
	}
}

func SliceValuesFunc[S ~[]BaseT, BaseT, OutT any](s S, f func(BaseT) OutT) iter.Seq[OutT] {
	return func(yield func(OutT) bool) {
		for _, v := range s {
			if !yield(f(v)) {
				return
			}
		}
	}
}
