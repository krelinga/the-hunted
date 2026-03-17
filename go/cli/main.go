package main

import (
	"fmt"
	"log"

	thehunted "github.com/krelinga/the-hunted/go"
)

type EventWriterPrinter struct {}

func (e EventWriterPrinter) WriteEvent(event thehunted.Event) {
	fmt.Printf("event: %s\n", event)
}

func main() {
	g := thehunted.Game{
		Roller:      thehunted.RandomRoller{},
		EventWriter: EventWriterPrinter{},
	}
	for !g.IsFinished() {
		f := g.Form()
		var err error
		switch f := f.(type) {
		case *thehunted.StartGameForm:
			err = handleStartGame(f)
		case *thehunted.SelectLoadoutForm:
			err = handleSelectLoadout(f)
		default:
			log.Fatalf("unexpected form type: %T", f)
		}
		if err != nil {
			log.Fatalf("form population error: %v", err)
		}
		if err := g.Advance(f); err != nil {
			log.Fatalf("error advancing game: %v", err)
		}
	}
}