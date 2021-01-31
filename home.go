package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/dustin/go-humanize"
	"github.com/fiatjaf/go-nostr/event"
)

func home(opts docopt.Opts) {
	initNostr()

	for _, follow := range config.Following {
		pool.SubKey(follow.Key)
	}

	pool.ReqFeed(nil)

	seen := make(map[string]bool)
	for em := range pool.Events {
		if em.Event.Kind != event.KindTextNote {
			continue
		}
		if _, ok := seen[em.Event.ID]; ok {
			continue
		}
		seen[em.Event.ID] = true

		fmt.Printf("%s at %s:\n  %s",
			em.Event.PubKey,
			humanize.Time(time.Unix(int64(em.Event.CreatedAt), 0)),
			strings.ReplaceAll(em.Event.Content, "\n", "\n  "),
		)
	}
}
