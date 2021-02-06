package main

import (
	"fmt"
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

	tags := make(event.Tags, 0, 1)
	if refid, err := opts.String("--reference"); err == nil {
		tags = append(tags, event.Tag([]interface{}{"e", refid}))
	}

	evt, err := pool.PublishEvent(&event.Event{
		PubKey:    getPubKey(config.PrivateKey),
		CreatedAt: uint32(time.Now().Unix()),
		Kind:      event.KindTextNote,
		Tags:      tags,
		Content:   opts["<content>"].(string),
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
	}

	pool.ReqEvent(evt.ID, nil)
	for em := range pool.Events {
		if em.Event.ID != evt.ID {
			continue
		}
		fmt.Sprint("Seen it on '%s'.", em.Relay)
	}
}
