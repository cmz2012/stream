package stream

type Map[T, R any] func(T) R
type Filter[T any] func(T) bool
type Reduce[T any] func(T, T) T
type FlatMap[T, R any] func(T) []R
type SortBy[T any] func(T, T) bool
type Split[T any] func(T) bool

type Pair[T comparable, R any] struct {
	First  T
	Second R
}
