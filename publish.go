package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
)

func publish(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("Can't publish. Private key not set.\n")
		return
	}

	initNostr()

	var event nostr.Event

	if file, _ := opts.String("--file"); file != "" {
		jsonb, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Failed reading content from file: %s", err)
			return
		}
		if err := json.Unmarshal(jsonb, &event); err != nil {
			log.Printf("Failed unmarshaling json from file: %s", err)
			return
		}
	} else {
		references, err := optSlice(opts, "--reference")
		if err != nil {
			return
		}

		var tags nostr.Tags
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

		content, _ := opts.String("<content>")
		if content == "" {
			log.Printf("Content must not be empty")
			return
		}
		if content == "-" {
			content, err = readContentStdin(4096)
			if err != nil {
				log.Printf("Failed reading content from stdin: %s", err)
				return
			}
		}

		event = nostr.Event{
			CreatedAt: time.Now(),
			Kind:      nostr.KindTextNote,
			Tags:      tags,
			Content:   content,
		}
	}

	publishEvent, statuses, err := pool.PublishEvent(&event)
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(publishEvent, statuses)
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
