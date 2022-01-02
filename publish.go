package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
)

func publish(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("No private key set.\n")
		return
	}

	initNostr()

	tags := make(nostr.Tags, 0, 1)
	if refid, err := opts.String("--reference"); err == nil {
		tags = append(tags, nostr.Tag([]interface{}{"e", refid}))
	}

	event, statuses, err := pool.PublishEvent(&nostr.Event{
		CreatedAt: uint32(time.Now().Unix()),
		Kind:      nostr.KindTextNote,
		Tags:      tags,
		Content:   opts["<content>"].(string),
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}
