package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
)

func viewEvent(opts docopt.Opts) {
	verbose, _ := opts.Bool("--verbose")
    jsonformat, _ := opts.Bool("--json")
	id := opts["<id>"].(string)
	if id == "" {
		log.Println("provided event ID was empty")
		return
	}
	initNostr()

	_, all := pool.Sub(nostr.Filters{{IDs: []string{id}}})
	for event := range nostr.Unique(all) {
		if event.ID != id {
			log.Printf("got unexpected event %s.\n", event.ID)
			continue
		}

		printEvent(event, nil, verbose, jsonformat)
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

// iterEventsWithTimeout returns a channel of events; this channel will be
// closed once events have stopped arriving for timeoutDuration
func iterEventsWithTimeout(events chan nostr.Event, timeoutDuration time.Duration) chan nostr.Event {
	timeout := time.After(timeoutDuration)
	tick := time.Tick(1 * time.Millisecond)

	resulsChan := make(chan nostr.Event)

	go func() {
		for {
			select {
			case <-timeout:
				close(resulsChan)
				return
			case <-tick:
				timeout = time.After(timeoutDuration)
				event := <-events
				resulsChan <- event
			}
		}
	}()

	return resulsChan
}
