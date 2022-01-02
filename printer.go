package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fiatjaf/go-nostr"
	"gopkg.in/yaml.v2"
)

var kindNames = map[int]string{
	nostr.KindSetMetadata:            "Profile Metadata",
	nostr.KindTextNote:               "Text Note",
	nostr.KindRecommendServer:        "Relay Recommendation",
	nostr.KindContactList:            "Contact List",
	nostr.KindEncryptedDirectMessage: "Encrypted Message",
	nostr.KindDeletion:               "Deletion Notice",
}

func printEvent(evt nostr.Event) {
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
	case nostr.KindSetMetadata:
		var metadata map[string]interface{}
		err := json.Unmarshal([]byte(evt.Content), &metadata)
		if err != nil {
			fmt.Printf("Invalid JSON: '%s',\n  %s",
				err.Error(), evt.Content)
			return
		}
		y, _ := yaml.Marshal(metadata)
		fmt.Print(string(y))
	case nostr.KindTextNote:
		fmt.Print(strings.ReplaceAll(evt.Content, "\n", "\n  "))
	case nostr.KindRecommendServer:
	case nostr.KindContactList:
	case nostr.KindEncryptedDirectMessage:
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

func printPublishStatus(event *nostr.Event, statuses chan nostr.PublishStatus) {
	for status := range statuses {
		switch status.Status {
		case nostr.PublishStatusSent:
			fmt.Printf("Sent event %s to '%s'.\n", event.ID, status.Relay)
		case nostr.PublishStatusFailed:
			fmt.Printf("Failed to send event %s to '%s'.\n", event.ID, status.Relay)
		case nostr.PublishStatusSucceeded:
			fmt.Printf("Seen %s on '%s'.\n", event.ID, status.Relay)
		}
	}
}
