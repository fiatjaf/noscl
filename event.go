package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
)

func viewEvent(opts docopt.Opts) {
	verbose, _ := opts.Bool("--verbose")
	id := opts["<id>"].(string)
	if id == "" {
		log.Println("provided event ID was empty")
		return
	}
	initNostr()

	_, _, unique := pool.Sub(nostr.Filters{{IDs: []string{id}}})

	for event := range unique {
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

	id := opts["<id>"].(string)
	if id == "" {
		log.Println("Event id is empty! Exiting.")
		return
	}

	event, statuses, err := pool.PublishEvent(&nostr.Event{
		CreatedAt: time.Now(),
		Kind:      nostr.KindDeletion,
		Tags:      nostr.Tags{nostr.Tag{"e", id}},
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}
