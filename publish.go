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
		log.Printf("No private key set.\n")
		return
	}

	initNostr()

	tags := []nostr.Tag{}

	references, err := optSlice(opts, "--reference")

	if err != nil {
		return
	}

	for _, ref := range references {
		tags = append(tags, nostr.Tag([]interface{}{"e", ref}))
	}

	profiles, err := optSlice(opts, "--profile")

	if err != nil {
		return
	}

	for _, profile := range profiles {
		tags = append(tags, nostr.Tag([]interface{}{"p", profile}))
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

func optSlice(opts docopt.Opts, key string) ([]string, error) {
	if v, ok := opts[key]; ok {
		vals, ok := v.([]string)
		if ok {
			return vals, nil
		}
	}

	return []string{}, errors.New("unable to find opt")
}
