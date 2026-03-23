package views2

type Viewer[T any] interface {
	View() T
}
