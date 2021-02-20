package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr/event"
)

func publish(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("No private key set.\n")
		return
	}

	initNostr()

	tags := make(event.Tags, 0, 1)
	if refid, err := opts.String("--reference"); err == nil {
		tags = append(tags, event.Tag([]interface{}{"e", refid}))
	}

	_, statuses, err := pool.PublishEvent(&event.Event{
		PubKey:    getPubKey(config.PrivateKey),
		CreatedAt: uint32(time.Now().Unix()),
		Kind:      event.KindTextNote,
		Tags:      tags,
		Content:   opts["<content>"].(string),
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(statuses)
}
