package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fiatjaf/go-nostr/event"
	"gopkg.in/yaml.v2"
)

func saveConfig(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("can't open config file " + path + ": " + err.Error())
		return
	}
	yaml.NewEncoder(f).Encode(config)
}

func printIncomingNotes() {
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
