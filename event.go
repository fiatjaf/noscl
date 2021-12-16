package main

import (
	"log"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr/filter"
)

func view(opts docopt.Opts) {
	initNostr()

	id := opts["<id>"].(string)

	sub := pool.Sub(filter.EventFilter{ID: id})

	for event := range sub.UniqueEvents {
		if event.ID != id {
			log.Printf("got unexpected event %s.\n", event.ID)
			continue
		}

		printEvent(event)
		break
	}
}
