package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fiatjaf/go-nostr/event"
	"github.com/fiatjaf/go-nostr/relaypool"
	"gopkg.in/yaml.v2"
)

var kindNames = map[int]string{
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

func printPublishStatus(statuses chan relaypool.PublishStatus) {
	for status := range statuses {
		switch status.Status {
		case relaypool.PublishStatusSent:
			fmt.Printf("Sent to '%s'.\n", status.Relay)
		case relaypool.PublishStatusFailed:
			fmt.Printf("Failed to send to '%s'.\n", status.Relay)
		case relaypool.PublishStatusSucceeded:
			fmt.Printf("Seen it on '%s'.\n", status.Relay)
		}
	}
}
