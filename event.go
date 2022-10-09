package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
)

func view(opts docopt.Opts) {
	verbose, _ := opts.Bool("--verbose")
	id := opts["<id>"].(string)
	if id == "" {
		log.Println("provided event ID was empty")
		return
	}
	initNostr()

	sub := pool.Sub(nostr.Filters{{IDs: []string{id}}})

	for event := range sub.UniqueEvents {
		if event.ID != id {
			log.Printf("got unexpected event %s.\n", event.ID)
			continue
		}

		printEvent(event, nil, verbose)
		break
	}
}

func deleteEvent(opts docopt.Opts) {
	initNostr()

	eventid := opts["<eventid>"].(string)
	if eventid == "" {
		log.Println("Event id is empty! Exiting.")
		return
	}

	event, statuses, err := pool.PublishEvent(&nostr.Event{
		CreatedAt: time.Now(),
		Kind:      nostr.KindDeletion,
		Tags:      nostr.Tags{nostr.StringList{"e", eventid}},
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}
