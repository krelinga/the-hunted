package views

type Viewer[T any] interface {
	View() T
}
