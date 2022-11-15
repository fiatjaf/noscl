package main

import (
	"errors"
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
)

func publish(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("Can't publish. Private key not set.\n")
		return
	}

	initNostr()

	var tags nostr.Tags

	references, err := optSlice(opts, "--reference")
	if err != nil {
		return
	}

	for _, ref := range references {
		tags = append(tags, nostr.Tag{"e", ref})
	}

	profiles, err := optSlice(opts, "--profile")
	if err != nil {
		return
	}

	for _, profile := range profiles {
		tags = append(tags, nostr.Tag{"p", profile})
	}

	content := opts["<content>"].(string)
	if content == "-" {
		content, err = readContentStdin(4096)
		if err != nil {
			log.Printf("Failed reading content from stdin: %s", err)
		}
	}
	event, statuses, err := pool.PublishEvent(&nostr.Event{
		CreatedAt: time.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      tags,
		Content:   content,
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}

func optSlice(opts docopt.Opts, key string) ([]string, error) {
	if v, ok := opts[key]; ok {
		vals, ok := v.([]string)
		if ok {
			return vals, nil
		}
	}

	return []string{}, errors.New("unable to find opt")
}
