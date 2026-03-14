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
		if events, err := g.Advance(f); err != nil {
			log.Fatalf("error advancing game: %v", err)
		} else {
			for _, event := range events {
				fmt.Printf("event: %s\n", event)
			}
		}
	}
}