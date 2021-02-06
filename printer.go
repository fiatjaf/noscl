package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fiatjaf/go-nostr/event"
	"gopkg.in/yaml.v2"
)

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

		printEvent(em.Event)
	}
}

func printEvent(evt event.Event) {
	fmt.Printf("[%s], from %s at %s\n",
		evt.ID,
		evt.PubKey,
		humanize.Time(time.Unix(int64(evt.CreatedAt), 0)),
	)

	switch evt.Kind {
	case event.KindSetMetadata:
		var metadata map[string]interface{}
		err := json.Unmarshal([]byte(evt.Content), &metadata)
		if err != nil {
			fmt.Printf("Profile Metadata invalid JSON: '%s',\n  %s",
				err.Error(), evt.Content)
			return
		}
		y, _ := yaml.Marshal(metadata)
		fmt.Printf("Profile Metadata:\n%s", string(y))
	case event.KindTextNote:
		fmt.Printf("Text Note:\n  %s", strings.ReplaceAll(evt.Content, "\n", "\n  "))
	case event.KindRecommendServer:
	case event.KindContactList:
	case event.KindEncryptedDirectMessage:
	}

	fmt.Printf("\n")
}
