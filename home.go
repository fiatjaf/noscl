package main

import (
	"encoding/json"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
)

func home(opts docopt.Opts) {
	if len(config.Following) == 0 {
		log.Println("You need to be following someone to run 'home'")
		return
	}

	initNostr()

	verbose, _ := opts.Bool("--verbose")

	var keys []string

	nameMap := map[string]string{}

	for _, follow := range config.Following {
		keys = append(keys, follow.Key)

		if follow.Name != "" {
			nameMap[follow.Key] = follow.Name
		}
	}

	sub := pool.Sub(nostr.Filters{{Authors: keys}})

	for event := range sub.UniqueEvents {
		// Do we have a nick for the author of this message?
		nick, ok := nameMap[event.PubKey]

		if !ok {
			nick = ""
		}

		// If we don't already have a nick for this user, and they are announcing their
		// new name, let's use it.

		if nick == "" {
			if event.Kind == nostr.KindSetMetadata {
				var metadata Metadata
				err := json.Unmarshal([]byte(event.Content), &metadata)
				if err != nil {
					log.Println("Failed to parse metadata.")
					continue
				}

				nick = metadata.Name
				nameMap[nick] = event.PubKey
			}
		}

		printEvent(event, &nick, verbose)
	}
}
