package main

import (
	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
)

func home(opts docopt.Opts) {
	initNostr()

	var keys []string
	for _, follow := range config.Following {
		keys = append(keys, follow.Key)
	}

	sub := pool.Sub(nostr.EventFilters{{Authors: keys}})

	for event := range sub.UniqueEvents {
		printEvent(event)
	}
}
