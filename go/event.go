package thehunted

type Event interface {
	eventIsAClosedType()

	String() string
}

type baseEvent struct {}

func (_ baseEvent) eventIsAClosedType() {}