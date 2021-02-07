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

var kindNames = map[uint8]string{
	event.KindSetMetadata:            "Profile Metadata",
	event.KindTextNote:               "Text Note",
	event.KindRecommendServer:        "Relay Recommendation",
	event.KindContactList:            "Contact List",
	event.KindEncryptedDirectMessage: "Encrypted Message",
}

func printEvent(evt event.Event) {
	kind, ok := kindNames[evt.Kind]
	if !ok {
		kind = "Unknown Kind"
	}

	fmt.Printf("%s [%s] from %s %s\n  ",
		kind,
		shorten(evt.ID),
		shorten(evt.PubKey),
		humanize.Time(time.Unix(int64(evt.CreatedAt), 0)),
	)

	switch evt.Kind {
	case event.KindSetMetadata:
		var metadata map[string]interface{}
		err := json.Unmarshal([]byte(evt.Content), &metadata)
		if err != nil {
			fmt.Printf("Invalid JSON: '%s',\n  %s",
				err.Error(), evt.Content)
			return
		}
		y, _ := yaml.Marshal(metadata)
		fmt.Print(string(y))
	case event.KindTextNote:
		fmt.Print(strings.ReplaceAll(evt.Content, "\n", "\n  "))
	case event.KindRecommendServer:
	case event.KindContactList:
	case event.KindEncryptedDirectMessage:
	default:
		fmt.Print(evt.Content)
	}

	fmt.Printf("\n")
}

func shorten(id string) string {
	if len(id) < 12 {
		return id
	}
	return id[0:4] + "..." + id[len(id)-4:]
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

		printEvent(em.Event)
	}
}
