package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
)

func home(opts docopt.Opts, inboxMode bool) {
	if len(config.Following) == 0 {
		log.Println("You need to be following someone to run 'home'")
		return
	}

	initNostr()

	verbose, _ := opts.Bool("--verbose")
	jsonformat, _ := opts.Bool("--json")
	noreplies, _ := opts.Bool("--noreplies")
	onlyreplies, _ := opts.Bool("--onlyreplies")
	kinds, kindserr := optSlice(opts, "--kinds")
	if kindserr != nil {
		return
	}
	var intkinds []int
	for _, kind := range kinds {
		if i, e := strconv.Atoi(kind); e == nil {
			intkinds = append(intkinds, i)
		}
	}
	since, _ := opts.Int("--since")
	until, _ := opts.Int("--until")
	limit, _ := opts.Int("--limit")

	var keys []string
	nameMap := map[string]string{}
	for _, follow := range config.Following {
		keys = append(keys, follow.Key)
		if follow.Name != "" {
			nameMap[follow.Key] = follow.Name
		}
	}
	pubkey := getPubKey(config.PrivateKey)
	filters := nostr.Filters{{Limit: limit}}
	if inboxMode {
		// Filter by p tag to me
		filters[0].Tags = nostr.TagMap{"p": {pubkey}}
		// Force kinds to encrypted messages
		intkinds = make([]int, 0)
		intkinds = append(intkinds, nostr.KindEncryptedDirectMessage)
	} else {
		filters[0].Authors = keys
	}
	if since > 0 {
		sinceTime := time.Unix(int64(since), 0)
		filters[0].Since = &sinceTime
	}
	if until > 0 {
		untilTime := time.Unix(int64(until), 0)
		filters[0].Until = &untilTime
	}
	filters[0].Kinds = intkinds
	_, all := pool.Sub(filters)
	for event := range nostr.Unique(all) {
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

		// if only want events referencing another
		if onlyreplies || noreplies {
			var hasReferences bool = false
			for _, tag := range event.Tags {
				if tag[0] == "e" {
					hasReferences = true
					if noreplies {
						continue
					}
				}
			}
			if onlyreplies && !hasReferences {
				continue
			}
		}

		printEvent(event, &nick, verbose, jsonformat)
	}
}
