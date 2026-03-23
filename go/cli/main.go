package main

import (
	"fmt"
	"log"

	thehunted "github.com/krelinga/the-hunted/go"
)

type EventWriterPrinter struct{}

func (_ EventWriterPrinter) WriteEvent(event thehunted.Event) {
	fmt.Printf("event: %s\n", event)
}

func main() {
	g := thehunted.Game{
		EventWriter: EventWriterPrinter{},
		Roller:      thehunted.RandomRoller{},
		Selector:    selector{},
	}
	for !g.Done() {
		if err := g.Next(); err != nil {
			log.Fatalf("error running game: %v", err)
		}
	}
}

type selector struct{}
