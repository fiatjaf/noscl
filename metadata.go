package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
)

type Metadata struct {
	Name        string `json:"name,omitempty"`
	About       string `json:"about,omitempty"`
	Picture     string `json:"picture,omitempty"`
    NIP05       string `json:"nip05,omitempty"`
    Banner      string `json:"banner,omitempty"`
    LUD16       string `json:"lud16,omitempty"`
    Website     string `json:"website,omitempty"`

}

func setMetadata(opts docopt.Opts) {
	initNostr()

	name, _ := opts.String("--name")
	about, _ := opts.String("--about")
	picture, _ := opts.String("--picture")
    nip05, _ := opts.String("--nip05")
    banner, _ := opts.String("--banner")
    lud16, _ := opts.String("--lud16")
    website, _ := opts.String("--website")

	jmetadata, _ := json.Marshal(Metadata{
		Name:           name,
		About:          about,
		Picture:        picture,
        NIP05:          nip05,
        Banner:         banner,
        LUD16:          lud16,
        Website:        website,
	})

	event, statuses, err := pool.PublishEvent(&nostr.Event{
		PubKey:    getPubKey(config.PrivateKey),
		CreatedAt: time.Now(),
		Kind:      nostr.KindSetMetadata,
		Tags:      make(nostr.Tags, 0),
		Content:   string(jmetadata),
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}
