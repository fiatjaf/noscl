package main

import (
	"log"

	"github.com/docopt/docopt-go"
)

func view(opts docopt.Opts) {
	key := opts["<id>"].(string)
	pool.ReqEvent(key, nil)

	for em := range pool.Events {
		if em.Event.ID != key {
			log.Printf("got unexpected event %s.\n", em.Event.ID)
			continue
		}

		printEvent(em.Event)
		break
	}
}
