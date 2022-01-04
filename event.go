package main

import (
	"log"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
)

func view(opts docopt.Opts) {
	id := opts["<id>"].(string)
	if id == "" {
		log.Println("provided event ID was empty")
		return
	}
	initNostr()

	sub := pool.Sub(nostr.EventFilters{{IDs: []string{id}}})

	for event := range sub.UniqueEvents {
		if event.ID != id {
			log.Printf("got unexpected event %s.\n", event.ID)
			continue
		}

		printEvent(event)
		break
	}
}
