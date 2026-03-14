package main

import (
	"fmt"
	"log"

	thehunted "github.com/krelinga/the-hunted/go"
)

func main() {
	g := thehunted.Game{
		Roller: thehunted.RandomRoller{},
	}
	for !g.IsFinished() {
		f := g.Form()
		var events []thehunted.Event
		var err error
		switch f := f.(type) {
		case thehunted.StartGameForm:
			events, err = handleStartGame(f, &g)
		default:
			log.Fatalf("unexpected form type: %T", f)
		}
		if err != nil {
			log.Fatalf("error advancing game: %v", err)
		}
		for _, event := range events {
			fmt.Printf("event: %s\n", event)
		}
	}
}