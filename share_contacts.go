package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
)

func shareContacts(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("No private key set.\n")
		return
	}

	if len(config.Following) == 0 {
		log.Printf("Contact list empty.\n")
		return
	}

	initNostr()

	var tags nostr.Tags
	for _, follow := range config.Following {
		relay := ""
		if len(follow.Relays) > 0 {
			relay = follow.Relays[0]
		}
		tag := nostr.Tag{"p", follow.Key, relay, follow.Name}
		tags = append(tags, tag)
	}

	event, statuses, err := pool.PublishEvent(&nostr.Event{
		CreatedAt: time.Now(),
		Kind:      nostr.KindContactList,
		Tags:      tags,
		Content:   "",
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}
