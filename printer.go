package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/nbd-wtf/go-nostr"
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

func printEvent(evt nostr.Event, nick *string, verbose bool) {
	kind, ok := kindNames[evt.Kind]
	if !ok {
		kind = "Unknown Kind"
	}

	// Don't print encrypted messages that aren't for me
	if evt.Kind == nostr.KindEncryptedDirectMessage {
		if !evt.Tags.ContainsAny("p", nostr.Tag{getPubKey(config.PrivateKey)}) {
			return
		}
	}

	var ID string = shorten(evt.ID)
	var fromField string = shorten(evt.PubKey)

	if nick != nil {
		fromField = fmt.Sprintf("%s (%s)", *nick, shorten(evt.PubKey))
	}

	if verbose {
		ID = evt.ID

		if nick == nil {
			fromField = evt.PubKey
		} else {
			fromField = fmt.Sprintf("%s (%s)", *nick, evt.PubKey)
		}
	}

	fmt.Printf("%s [%s] from %s %s\n  ",
		kind,
		ID,
		fromField,
		humanize.Time(evt.CreatedAt),
	)

	switch evt.Kind {
	case nostr.KindSetMetadata:
		var metadata Metadata
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
