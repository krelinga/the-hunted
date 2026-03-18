package thehunted

type Event interface {
	eventIsAClosedType()

	String() string
}

type baseEvent struct {}

func (_ baseEvent) eventIsAClosedType() {}

type EventWriter interface {
	WriteEvent(Event)
}

type NilEventWriter struct {}

func (_ NilEventWriter) WriteEvent(_ Event) {}